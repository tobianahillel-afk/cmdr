// First-party same-origin transport. A search row and a JSON response are never
// authority; the HTTP host independently authorizes Search Job and Event reads.
const INSPECT_ROUTE = "/investigate/event-inspection";
const MAX_RESPONSE_BYTES = 2 * 1024 * 1024;

export function createInspectionHTTPTransport({ fetchImpl = globalThis.fetch } = {}) {
  if (typeof fetchImpl !== "function") throw new TypeError("fetch unavailable");
  return Object.freeze({
    async inspect(request, { signal } = {}) {
      const response = await fetchImpl(INSPECT_ROUTE, {
        method: "POST",
        headers: { "Content-Type": "application/json", "Accept": "application/json" },
        credentials: "same-origin",
        redirect: "error",
        cache: "no-store",
        body: JSON.stringify(request),
        signal,
      });
      // Never propagate a backend denial's details or a cross-origin redirect.
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
        // Reject malformed UTF-8 before handing data to the existing strict
        // frontend projection validator; no DOM interpolation occurs here.
        return JSON.parse(new TextDecoder("utf-8", { fatal: true }).decode(bytes));
      } catch {
        throw new Error("inspection-unavailable");
      }
    },
  });
}
