// A single first-party, same-origin, bounded transport for read inspection and
// separately host-authorized *unexecuted* pivot drafts. No URL from the caller.
const INSPECT_ROUTE = "/investigate/event-inspection";
const PIVOT_ROUTE = "/investigate/event-inspection/pivot";
const MAX_RESPONSE_BYTES = 2 * 1024 * 1024;

async function postAuthorizedJSON(fetchImpl, route, request, signal) {
  const response = await fetchImpl(route, {
    method: "POST",
    headers: { "Content-Type": "application/json", "Accept": "application/json" },
    credentials: "same-origin",
    redirect: "error",
    cache: "no-store",
    body: JSON.stringify(request),
    signal,
  });
  // An HTTP 404 denial is indistinguishable from another unavailable source.
  if (response?.status !== 200 || response.redirected === true ||
      response.type === "opaque" || response.type === "opaqueredirect" ||
      response.headers?.get("content-type")?.split(";", 1)[0].trim().toLowerCase() !== "application/json") {
    throw new Error("inspection-unavailable");
  }
  const declaredLength = response.headers.get("content-length");
  if (declaredLength !== null && (!/^\d+$/.test(declaredLength) ||
      Number(declaredLength) > MAX_RESPONSE_BYTES)) throw new Error("inspection-unavailable");
  if (typeof response.body?.getReader !== "function") throw new Error("inspection-unavailable");
  const reader = response.body.getReader();
  const chunks = [];
  let received = 0;
  try {
    while (true) {
      const { done, value } = await reader.read();
      if (done) break;
      if (!(value instanceof Uint8Array) || received + value.length > MAX_RESPONSE_BYTES) {
        await reader.cancel();
        throw new Error("inspection-unavailable");
      }
      chunks.push(value);
      received += value.length;
    }
  } finally {
    reader.releaseLock();
  }
  const bytes = new Uint8Array(received);
  let offset = 0;
  for (const chunk of chunks) {
    bytes.set(chunk, offset);
    offset += chunk.length;
  }
  try {
    return JSON.parse(new TextDecoder("utf-8", { fatal: true }).decode(bytes));
  } catch {
    throw new Error("inspection-unavailable");
  }
}

export function createInspectionHTTPTransport({ fetchImpl = globalThis.fetch } = {}) {
  if (typeof fetchImpl !== "function") throw new TypeError("fetch unavailable");
  return Object.freeze({
    inspect: (request, { signal } = {}) =>
      postAuthorizedJSON(fetchImpl, INSPECT_ROUTE, request, signal),
    // Callers cannot select an endpoint, permission, return origin or executor.
    draftPivot: (request, { signal } = {}) =>
      postAuthorizedJSON(fetchImpl, PIVOT_ROUTE, request, signal),
  });
}
