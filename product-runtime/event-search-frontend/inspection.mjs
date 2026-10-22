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

// A chosen field must exist in the authorized visible Inspector projection.
function pivotOrigins(inspected, field, value) {
  const result = [];
  for (const [name, values] of [["source", inspected.sourceFields],
    ["normalized", inspected.normalizedFields],
    ["raw-derived-sensitive", inspected.rawDerivedSensitive]]) {
    if (Object.hasOwn(values, field) && values[field] === value) result.push(name);
  }
  for (const enrichment of inspected.enrichments) {
    if (Object.hasOwn(enrichment.values, field) && enrichment.values[field] === value) {
      result.push(`enrichment:${enrichment.producer}@${enrichment.version}`);
    }
  }
  return result;
}

// Go normalizes RFC3339 UTC dates ("Z"), while JS commonly serializes
// the same instant with ".000Z". Compare exact nanoseconds, not text or
// lossy millisecond-only Date objects; never widen the chosen Search range.
const RFC3339_TIME = /^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2})(?:\.(\d{1,9}))?(Z|[+-]\d{2}:\d{2})$/u;
function exactInstant(value) {
  if (!text(value, 64)) return null;
  const parts = RFC3339_TIME.exec(value);
  if (!parts) return null;
  const seconds = Date.parse(parts[1] + parts[3]);
  if (!Number.isFinite(seconds)) return null;
  const fraction = BigInt((parts[2] ?? "").padEnd(9, "0") || "0");
  return BigInt(seconds) * 1_000_000n + fraction;
}
function sameInstant(first, second) {
  const left = exactInstant(first);
  return left !== null && left === exactInstant(second);
}

// A draft is not a Query execution or permission token. Trust remains on Go.
export function normalizePivotDraft(input, request, inspected) {
  if (!record(input) || input.kind !== "pivot-draft" ||
      input.targetCapability !== "CAP-INV-002" || input.auditRequired !== true ||
      !FIELD.test(input.field ?? "") || input.field !== request.field ||
      !text(input.value) || input.value !== request.value ||
      input.jobId !== request.jobId || input.runId !== request.runId ||
      !sameInstant(input.timeRange?.start, request.timeStart) ||
      !sameInstant(input.timeRange?.end, request.timeEnd)) fail("invalid-pivot-draft");
  for (const key of ["tenantRef", "environmentRef", "eventRef", "correlationId"]) {
    if (input[key] !== request[key]) fail("pivot-scope-mismatch");
  }
  if (!record(input.returnContext) || input.returnContext.origin !== "event-search" ||
      ["tenantRef", "environmentRef", "eventRef", "correlationId"].some(
        key => input.returnContext[key] !== request[key])) fail("invalid-pivot-return");
  if (input.sourceRef !== inspected.source.sourceRef ||
      input.dataSourceRef !== inspected.source.dataSourceRef) fail("invalid-pivot-source");
  const origins = pivotOrigins(inspected, request.field, request.value);
  if (origins.length !== 1 || origins[0] !== input.fieldOrigin) fail("ambiguous-pivot-provenance");
  return immutable({
    kind: "pivot-draft", targetCapability: input.targetCapability,
    field: input.field, value: input.value, fieldOrigin: input.fieldOrigin,
    tenantRef: input.tenantRef, environmentRef: input.environmentRef,
    eventRef: input.eventRef, jobId: input.jobId, runId: input.runId,
    correlationId: input.correlationId,
    sourceRef: input.sourceRef, dataSourceRef: input.dataSourceRef,
    timeRange: { start: input.timeRange.start, end: input.timeRange.end },
    returnContext: {
      tenantRef: input.returnContext.tenantRef,
      environmentRef: input.returnContext.environmentRef,
      eventRef: input.returnContext.eventRef,
      correlationId: input.returnContext.correlationId,
      origin: "event-search",
    }, auditRequired: true,
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
      openedFrom = { start, request, searchState: scope, resultWindow: snapshot.resultWindow };
      return publish({ kind: "inspected", inspection: inspected });
    } catch {
      if (generation !== sequence || control.signal.aborted) return { kind: "superseded" };
      openedFrom = undefined;
      return publish({ kind: "unavailable" });
    } finally {
      if (generation === sequence) abort = undefined;
    }
  }

  async function draftPivot(field, value) {
    const prior = view;
    const origin = openedFrom;
    if (prior.kind !== "inspected" || !origin ||
        typeof transport.draftPivot !== "function" ||
        !FIELD.test(field ?? "") || !text(value) ||
        pivotOrigins(prior.inspection, field, value).length !== 1 ||
        !validSearch(origin.start)) fail("pivot-not-available");
    const job = origin.start.job;
    if (!text(job.timeStart, 64) || !text(job.timeEnd, 64)) fail("pivot-invalid-time");
    abort?.abort();
    const generation = ++sequence;
    const control = new AbortController();
    abort = control;
    const request = immutable({
      ...origin.request, field, value, valuePresent: true,
      timeStart: job.timeStart, timeEnd: job.timeEnd,
    });
    publish({ ...prior, pivotStatus: "loading", pivotDraft: undefined });
    try {
      const response = await transport.draftPivot(request, { signal: control.signal });
      if (generation !== sequence || control.signal.aborted) return { kind: "superseded" };
      if (!validSearch(origin.start)) {
        openedFrom = undefined;
        return publish({ kind: "unavailable" });
      }
      const draft = normalizePivotDraft(response, request, prior.inspection);
      return publish({ kind: "inspected", inspection: prior.inspection, pivotDraft: draft });
    } catch {
      if (generation !== sequence || control.signal.aborted) return { kind: "superseded" };
      if (!validSearch(origin.start)) {
        openedFrom = undefined;
        return publish({ kind: "unavailable" });
      }
      return publish({ kind: "inspected", inspection: prior.inspection, pivotStatus: "unavailable" });
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

  return Object.freeze({ open, draftPivot, back, close, getView: () => view });
}

// Render with native DOM text-only sinks. This is a view of server-authorized
// content, never an authorization mechanism or a parser/decoding engine.
export function renderEventInspection(document, mount, current, { onDraftPivot } = {}) {
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
    const missing = make("section");
    missing.append(make("h3", "Missing fields"));
    for (const name of p.missingFields) missing.append(make("p", name));
    section.append(missing);
    for (const enrichment of p.enrichments) {
      const group = make("section");
      group.append(make("h3", `Enrichment: ${enrichment.producer}@${enrichment.version} (${enrichment.freshness})`));
      for (const [field, value] of Object.entries(enrichment.values)) {
        group.append(make("p", `${field}: ${value}`));
      }
      section.append(group);
    }
    if (typeof onDraftPivot === "function") {
      const pivot = make("section");
      pivot.append(make("h3", "Prepare an authorized, unexecuted Search pivot draft"));
      const choices = [
        ...Object.entries(p.sourceFields),
        ...Object.entries(p.normalizedFields),
        ...Object.entries(p.rawDerivedSensitive),
        ...p.enrichments.flatMap(e => Object.entries(e.values)),
      ];
      for (const [field, value] of choices) {
        if (pivotOrigins(p, field, value).length !== 1) continue;
        const button = make("button", `Prepare pivot: ${field} = ${value}`);
        button.type = "button";
        button.addEventListener("click", () => onDraftPivot(field, value));
        pivot.append(button);
      }
      section.append(pivot);
    }
    if (current.pivotDraft) {
      section.append(make("p", `Read-only draft, not executed: ${current.pivotDraft.field} = ${current.pivotDraft.value} | origin ${current.pivotDraft.fieldOrigin}`));
    } else if (current.pivotStatus === "unavailable") {
      section.append(make("p", "Pivot draft unavailable; no search was executed"));
    }
    for (const [label, item] of [["Raw", p.raw], ["Rendered", p.rendered]]) {
      section.append(make("h3", label));
      section.append(make("pre", item.access === "allowed" ? item.text : `${label} unavailable`));
    }
  }
  mount.replaceChildren(section);
}
