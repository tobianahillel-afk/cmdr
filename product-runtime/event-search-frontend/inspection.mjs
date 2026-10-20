import { normalizeDeepLinkState } from "./state.mjs";

const REF = /^[A-Za-z0-9_.:/-]{1,256}$/u;
const FIELD = /^[A-Za-z0-9_.:-]{1,128}$/u;
const STATES = new Set(["available", "partial", "restricted", "source-unavailable", "tombstone"]);
const MAX_FIELDS = 256;

export class InspectionJourneyError extends Error {
  constructor(code) {
    super(code);
    this.name = "InspectionJourneyError";
    this.code = code;
  }
}

const fail = (code) => { throw new InspectionJourneyError(code); };
const record = (v) => v !== null && typeof v === "object" &&
  !Array.isArray(v) && (Object.getPrototypeOf(v) === Object.prototype || Object.getPrototypeOf(v) === null);
const ref = (v) => typeof v === "string" && REF.test(v);
const text = (v, max = 4096) => typeof v === "string" && v.length <= max &&
  !v.includes("\u0000");
const immutable = (value) => {
  if (value !== null && typeof value === "object" && !Object.isFrozen(value)) {
    for (const item of Object.values(value)) immutable(item);
    Object.freeze(value);
  }
  return value;
};

function fields(value) {
  if (!record(value) || Object.keys(value).length > MAX_FIELDS) fail("invalid-fields");
  const out = Object.create(null);
  for (const [key, val] of Object.entries(value)) {
    if (!FIELD.test(key) || ["__proto__", "constructor", "prototype"].includes(key) || !text(val)) {
      fail("invalid-field-value");
    }
    out[key] = val;
  }
  return out;
}

function payload(input) {
  if (!record(input) || !["allowed", "denied"].includes(input.access)) fail("invalid-payload");
  if (input.access === "denied") {
    if (Object.hasOwn(input, "text")) fail("denied-payload-present");
    return { access: "denied" };
  }
  if (!text(input.text, 65536)) fail("invalid-payload-text");
  return { access: "allowed", text: input.text };
}

// Only the host's already-authorized projection can be displayed. Nothing in
// the browser response or search URL is a proof of access; reauthorization
// and sensitive-field classification belong to the server adapter.
export function normalizeInspection(input, request) {
  if (!record(input) || input.kind !== "inspection") fail("invalid-inspection");
  for (const key of ["tenantRef", "environmentRef", "eventRef", "correlationId"]) {
    if (!ref(input[key]) || input[key] !== request[key]) fail("inspection-scope-mismatch");
  }
  if (!STATES.has(input.state)) fail("invalid-event-state");
  if (!record(input.source) || !ref(input.source.sourceRef) ||
      !text(input.source.dataSourceRef ?? "", 256) || !record(input.parser)) fail("invalid-provenance");
  for (const key of ["parserRef", "version", "status"]) {
    if (!text(input.parser[key] ?? "", 256)) fail("invalid-parser");
  }
  const raw = payload(input.raw);
  const rendered = payload(input.rendered);
  if (raw.access === "denied" && input.rawDerivedSensitive !== undefined &&
      Object.keys(fields(input.rawDerivedSensitive)).length !== 0) {
    fail("raw-derived-without-access");
  }
  const enrichments = input.enrichments ?? [];
  if (!Array.isArray(enrichments) || enrichments.length > 64) fail("invalid-enrichments");
  const normalized = enrichments.map((enrichment) => {
    if (!record(enrichment) || !ref(enrichment.producer) ||
        !text(enrichment.version, 256) || !text(enrichment.freshness, 256) ||
        typeof enrichment.stale !== "boolean" || enrichment.presentedAsSource === true) {
      fail("invalid-enrichment-provenance");
    }
    return {
      producer: enrichment.producer,
      version: enrichment.version,
      freshness: enrichment.freshness,
      stale: enrichment.stale,
      values: fields(enrichment.values ?? {}),
    };
  });
  if (!Array.isArray(input.missingFields) || input.missingFields.length > MAX_FIELDS ||
      !input.missingFields.every((x) => typeof x === "string" && FIELD.test(x))) {
    fail("invalid-missing-fields");
  }
  return immutable({
    kind: "inspected", tenantRef: request.tenantRef,
    environmentRef: request.environmentRef, eventRef: request.eventRef,
    correlationId: request.correlationId, state: input.state,
    source: { sourceRef: input.source.sourceRef, dataSourceRef: input.source.dataSourceRef ?? "" },
    parser: { parserRef: input.parser.parserRef ?? "", version: input.parser.version ?? "",
      status: input.parser.status ?? "" },
    sourceFields: fields(input.sourceFields ?? {}),
    normalizedFields: fields(input.normalizedFields ?? {}),
    rawDerivedSensitive: raw.access === "allowed" ? fields(input.rawDerivedSensitive ?? {}) : fields({}),
    missingFields: [...input.missingFields], enrichments: normalized, raw, rendered,
  });
}

export function createEventInspectionJourney({ searchController, transport, onView = () => {} }) {
  if (typeof searchController?.getView !== "function" || typeof transport?.inspect !== "function" ||
      typeof onView !== "function") fail("invalid-journey");
  let sequence = 0;
  let abort;
  let view = immutable({ kind: "search" });
  let openedFrom;

  function publish(next) {
    view = immutable(next);
    onView(view);
    return view;
  }

  function validSearch(start) {
    const now = searchController.getView();
    return now.job?.jobId === start.job.jobId && now.job?.runId === start.job.runId &&
      now.job?.correlationId === start.job.correlationId &&
      now.job?.tenantRef === start.job.tenantRef &&
      now.job?.environmentRef === start.job.environmentRef &&
      !["permission-denied", "loading", "error", "offline"].includes(now.shellState?.kind) &&
      now.resultWindow?.rows.some((row) => row.ref === start.row.ref && row.source === start.row.source);
  }

  async function open(searchState, eventRef) {
    const scope = normalizeDeepLinkState(searchState);
    const snapshot = searchController.getView();
    const row = snapshot.resultWindow?.rows.find((item) => item.ref === eventRef);
    if (!row || !snapshot.job || !ref(eventRef) || !ref(snapshot.job.jobId) ||
        !ref(snapshot.job.runId) || !ref(snapshot.job.correlationId) ||
        snapshot.job.tenantRef !== scope.tenantId ||
        snapshot.job.environmentRef !== scope.environmentId ||
        ["permission-denied", "loading", "error", "offline"].includes(snapshot.shellState?.kind)) fail("not-authorized-result");
    abort?.abort();
    const generation = ++sequence;
    const control = new AbortController();
    abort = control;
    const request = immutable({
      tenantRef: scope.tenantId, environmentRef: scope.environmentId,
      eventRef, correlationId: snapshot.job.correlationId,
      jobId: snapshot.job.jobId, runId: snapshot.job.runId,
      sourceHint: row.source,
    });
    const start = { job: snapshot.job, row };
    openedFrom = undefined;
    publish({ kind: "loading" });
    try {
      const response = await transport.inspect(request, { signal: control.signal });
      if (generation !== sequence || control.signal.aborted) return { kind: "superseded" };
      if (!validSearch(start)) return publish({ kind: "unavailable" });
      if (record(response) && response.kind === "permission-denied") return publish({ kind: "unavailable" });
      const inspected = normalizeInspection(response, request);
      openedFrom = { start, searchState: scope, resultWindow: snapshot.resultWindow };
      return publish({ kind: "inspected", inspection: inspected });
    } catch {
      if (generation !== sequence || control.signal.aborted) return { kind: "superseded" };
      openedFrom = undefined;
      return publish({ kind: "unavailable" });
    } finally {
      if (generation === sequence) abort = undefined;
    }
  }

  function back() {
    sequence += 1;
    abort?.abort();
    abort = undefined;
    const previous = openedFrom;
    openedFrom = undefined;
    if (!previous || !validSearch(previous.start)) return publish({ kind: "unavailable" });
    return publish({ kind: "search", searchState: previous.searchState,
      resultWindow: previous.resultWindow });
  }

  function close() {
    sequence += 1;
    abort?.abort();
    abort = undefined;
    openedFrom = undefined;
    return publish({ kind: "search" });
  }

  return Object.freeze({ open, back, close, getView: () => view });
}

// Render with native DOM text-only sinks. This is a view of server-authorized
// content, never an authorization mechanism or a parser/decoding engine.
export function renderEventInspection(document, mount, current) {
  if (typeof document?.createElement !== "function" ||
      typeof mount?.replaceChildren !== "function") fail("invalid-dom");
  const make = (tag, value = "") => {
    const el = document.createElement(tag);
    el.textContent = value;
    return el;
  };
  const section = make("section");
  const heading = make("h2", current.kind === "inspected" ? "Event inspection" :
    current.kind === "loading" ? "Loading event inspection" : "Event inspection unavailable");
  section.append(heading);
  if (current.kind === "inspected") {
    const p = current.inspection;
    section.append(make("p", `Event: ${p.eventRef} | State: ${p.state}`));
    section.append(make("p", `Source: ${p.source.sourceRef} | Parser: ${p.parser.parserRef} ${p.parser.version}`));
    for (const [label, source] of [["Source fields", p.sourceFields],
      ["Normalized fields", p.normalizedFields],
      ["Raw-derived fields", p.rawDerivedSensitive]]) {
      const group = make("section");
      group.append(make("h3", label));
      for (const [field, value] of Object.entries(source)) {
        group.append(make("p", `${field}: ${value}`));
      }
      section.append(group);
    }
    for (const enrichment of p.enrichments) {
      const group = make("section");
      group.append(make("h3", `Enrichment: ${enrichment.producer}@${enrichment.version} (${enrichment.freshness})`));
      for (const [field, value] of Object.entries(enrichment.values)) {
        group.append(make("p", `${field}: ${value}`));
      }
      section.append(group);
    }
    for (const [label, item] of [["Raw", p.raw], ["Rendered", p.rendered]]) {
      section.append(make("h3", label));
      section.append(make("pre", item.access === "allowed" ? item.text : `${label} unavailable`));
    }
  }
  mount.replaceChildren(section);
}
