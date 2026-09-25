import { createResultWindow, normalizeResultCollection } from "./results.mjs";
import { normalizeShellState } from "./shell.mjs";
import { normalizeDeepLinkState } from "./state.mjs";

const JOB_STATES = new Set(["queued", "running", "completed", "partial", "failed", "cancelled"]);
const TERMINAL_STATES = new Set(["completed", "partial", "failed", "cancelled"]);
const OPAQUE_REF = /^[A-Za-z0-9_.:/-]{1,256}$/u;
const MAX_FAILURES = 256;
const MAX_RESULTS = 10000;

export class IntegrationValidationError extends Error {
  constructor(code, message) {
    super(message);
    this.name = "IntegrationValidationError";
    this.code = code;
  }
}

export function createExecutionRequest(deepLinkInput, executionInput, expected = {}) {
  const deepLink = normalizeDeepLinkState(deepLinkInput, expected);
  if (!isPlainObject(executionInput)) {
    throw invalid("invalid-execution", "execution metadata must be a plain object");
  }
  rejectUnknown(executionInput, new Set(["correlationId", "queryVersion"]));
  return freezeDeep({
    tenantRef: deepLink.tenantId,
    environmentRef: deepLink.environmentId,
    timeStart: deepLink.from,
    timeEnd: deepLink.to,
    sources: Object.freeze([...deepLink.sources]),
    query: deepLink.query,
    filters: deepLink.filters ?? Object.freeze({}),
    correlationId: opaqueRef("correlation id", executionInput.correlationId),
    queryVersion: opaqueRef("query version", executionInput.queryVersion),
  });
}

export function normalizeTransportOutcome(input, request, options = {}) {
  validateRequest(request);
  if (!isPlainObject(options)) {
    throw invalid("invalid-options", "transport normalization options must be a plain object");
  }
  rejectUnknown(options, new Set(["retryOfRunId"]));
  const retryOfRunId = optionalOpaqueRef("retry parent run id", options.retryOfRunId);

  if (!isPlainObject(input)) {
    throw invalid("invalid-transport-outcome", "transport outcome must be a plain object");
  }
  const kind = requiredText("transport outcome kind", input.kind, 64);
  if (kind === "permission-denied") {
    rejectUnknown(input, new Set([
      "kind",
      "correlationId",
      "capability",
      "reasonCode",
      "requestAccessPath",
    ]));
    const correlationId = opaqueRef("permission correlation id", input.correlationId);
    if (correlationId !== request.correlationId) {
      throw invalid("correlation-mismatch", "permission denial correlation does not match request");
    }
    const shellState = normalizeShellState({
      kind: "permission-denied",
      capability: requiredText("denied capability", input.capability, 256),
      reasonCode: opaqueRef("deny reason code", input.reasonCode),
      requestAccessPath: input.requestAccessPath,
    });
    return freezeDeep({ kind, correlationId, shellState });
  }
  if (kind !== "job") {
    throw invalid("unknown-transport-outcome", `unknown transport outcome: ${kind}`);
  }

  const allowed = new Set([
    "kind",
    "jobId",
    "runId",
    "state",
    "tenantRef",
    "environmentRef",
    "timeStart",
    "timeEnd",
    "sources",
    "correlationId",
    "queryVersion",
    "parentRunId",
    "resultRefs",
    "failedSources",
    "projections",
    "freshness",
    "consequence",
  ]);
  rejectUnknown(input, allowed);

  const job = {
    jobId: opaqueRef("job id", input.jobId),
    runId: opaqueRef("run id", input.runId),
    state: requiredText("job state", input.state, 64),
    tenantRef: requiredText("job tenant", input.tenantRef, 256),
    environmentRef: requiredText("job environment", input.environmentRef, 256),
    timeStart: normalizeInstant("job time start", input.timeStart),
    timeEnd: normalizeInstant("job time end", input.timeEnd),
    sources: normalizeOpaqueList("job source", input.sources, 512),
    correlationId: opaqueRef("job correlation id", input.correlationId),
    queryVersion: opaqueRef("job query version", input.queryVersion),
    parentRunId: optionalOpaqueRef("job parent run id", input.parentRunId),
    resultRefs: normalizeResultRefs(input.resultRefs ?? []),
    failedSources: normalizeFailedSources(input.failedSources ?? []),
  };
  if (!JOB_STATES.has(job.state)) {
    throw invalid("unknown-job-state", `unknown Search Job state: ${job.state}`);
  }
  assertJobScopeMatchesRequest(job, request);
  if (retryOfRunId === undefined) {
    if (job.parentRunId !== undefined) {
      throw invalid("unexpected-parent-run", "initial execution must not claim a parent run");
    }
  } else {
    if (job.parentRunId !== retryOfRunId) {
      throw invalid("retry-parent-mismatch", "retry must identify the previous run as parent");
    }
    if (job.runId === retryOfRunId) {
      throw invalid("retry-run-reuse", "retry must create a new run id");
    }
  }

  const projections = normalizeResultCollection(input.projections ?? []);
  validateOutcomeCardinality(job, projections);
  validateProjectionsMatchRefs(job.resultRefs, projections);

  let shellState = null;
  if (job.state === "partial") {
    const freshness = normalizeFreshness(input.freshness);
    shellState = normalizeShellState({
      kind: "partial",
      failedSources: job.failedSources.map((failure) => failure.source),
      freshness,
      consequence: requiredText("partial consequence", input.consequence, 1024),
    });
  } else {
    if (input.freshness !== undefined || input.consequence !== undefined) {
      throw invalid("unexpected-partial-metadata", "freshness/consequence are only valid for partial jobs");
    }
    if (job.state === "failed") {
      shellState = normalizeShellState({
        kind: "error",
        message: "Search failed. Retry creates a new Search Job.",
        correlationId: job.correlationId,
        retryable: true,
      });
    } else if (job.state === "completed" && projections.length === 0) {
      shellState = normalizeShellState({
        kind: "empty",
        reason: "No events match the current search scope.",
        adjustmentHint: "Adjust the time range, sources, or filters.",
      });
    }
  }

  return freezeDeep({
    kind,
    job: freezeDeep(job),
    projections,
    shellState,
    terminal: TERMINAL_STATES.has(job.state),
  });
}

export function createEventSearchController(options) {
  if (!isPlainObject(options)) {
    throw invalid("invalid-controller-options", "controller options must be a plain object");
  }
  rejectUnknown(options, new Set(["transport", "onView", "windowSize", "expectedContext"]));
  const transport = validateTransport(options.transport);
  const onView = options.onView ?? (() => {});
  if (typeof onView !== "function") {
    throw invalid("invalid-view-handler", "onView must be a function");
  }
  const windowSize = boundedInteger("result window size", options.windowSize ?? 50, 1, 200);
  const expectedContext = isPlainObject(options.expectedContext) ? options.expectedContext : {};
  let generation = 0;
  let activeAbort;
  let lastJob;
  let lastRequest;
  let resultWindow;
  let currentView = freezeDeep({ shellState: null, resultWindow: undefined, job: undefined });

  function publish(shellState, windowValue, job) {
    currentView = freezeDeep({
      shellState,
      resultWindow: windowValue,
      job,
    });
    onView(currentView);
    return currentView;
  }

  async function execute(deepLinkInput, executionInput) {
    const request = createExecutionRequest(deepLinkInput, executionInput, expectedContext);
    const preserve = sameExecutionIntent(lastRequest, request) ? resultWindow : undefined;
    return runOperation("execute", request, preserve, undefined);
  }

  async function retry() {
    if (lastJob === undefined || lastRequest === undefined || !lastJob.terminal) {
      throw invalid("retry-unavailable", "retry requires a terminal Search Job");
    }
    return runOperation("retry", lastRequest, resultWindow, lastJob.job.runId);
  }

  async function runOperation(mode, request, preservedWindow, retryOfRunId) {
    if (activeAbort !== undefined) {
      activeAbort.abort();
    }
    const token = ++generation;
    const controller = new AbortController();
    activeAbort = controller;
    publish(normalizeShellState({ kind: "loading" }), preservedWindow, lastJob?.job);

    try {
      const raw = mode === "execute"
        ? await transport.execute(request, { signal: controller.signal })
        : await transport.retry(
          freezeDeep({
            jobId: lastJob.job.jobId,
            runId: lastJob.job.runId,
            correlationId: lastJob.job.correlationId,
            queryVersion: lastJob.job.queryVersion,
          }),
          { signal: controller.signal },
        );

      if (token !== generation || controller.signal.aborted) {
        return freezeDeep({ kind: "superseded" });
      }
      const outcome = normalizeTransportOutcome(raw, request, { retryOfRunId });
      if (outcome.kind === "permission-denied") {
        lastJob = undefined;
        lastRequest = request;
        resultWindow = undefined;
        publish(outcome.shellState, undefined, undefined);
        return outcome;
      }

      if (retryOfRunId !== undefined && outcome.job.jobId === lastJob.job.jobId) {
        throw invalid("retry-job-reuse", "retry must create a new job id");
      }
      lastRequest = request;
      if (outcome.terminal) {
        lastJob = outcome;
      }

      if (outcome.job.state === "cancelled") {
        publish(null, preservedWindow, outcome.job);
        return outcome;
      }
      if (outcome.job.state === "queued" || outcome.job.state === "running") {
        publish(normalizeShellState({ kind: "loading" }), preservedWindow, outcome.job);
        return outcome;
      }

      if (outcome.projections.length > 0) {
        resultWindow = createResultWindow(outcome.projections, { size: windowSize });
      } else {
        resultWindow = undefined;
      }
      publish(outcome.shellState, resultWindow, outcome.job);
      return outcome;
    } catch (error) {
      if (token !== generation || controller.signal.aborted || isAbortError(error)) {
        return freezeDeep({ kind: "cancelled" });
      }
      const safe = normalizeShellState({
        kind: "error",
        message: "Search failed. Retry when the service is available.",
        correlationId: request.correlationId,
        retryable: true,
      });
      publish(safe, preservedWindow, lastJob?.job);
      return freezeDeep({ kind: "error", code: "transport-or-response-rejected" });
    } finally {
      if (token === generation) {
        activeAbort = undefined;
      }
    }
  }

  function cancel() {
    if (activeAbort === undefined) {
      return false;
    }
    generation += 1;
    activeAbort.abort();
    activeAbort = undefined;
    return true;
  }

  return Object.freeze({
    execute,
    retry,
    cancel,
    getView() {
      return currentView;
    },
  });
}

function validateTransport(value) {
  if (!isPlainObject(value)) {
    throw invalid("invalid-transport", "transport must be a plain object");
  }
  rejectUnknown(value, new Set(["execute", "retry"]));
  if (typeof value.execute !== "function" || typeof value.retry !== "function") {
    throw invalid("invalid-transport", "transport requires execute and retry functions");
  }
  return value;
}

function validateRequest(request) {
  if (!isPlainObject(request)) {
    throw invalid("invalid-request", "normalized execution request is required");
  }
  for (const key of [
    "tenantRef",
    "environmentRef",
    "timeStart",
    "timeEnd",
    "sources",
    "query",
    "filters",
    "correlationId",
    "queryVersion",
  ]) {
    if (!(key in request)) {
      throw invalid("invalid-request", `normalized request is missing ${key}`);
    }
  }
}

function assertJobScopeMatchesRequest(job, request) {
  if (
    job.tenantRef !== request.tenantRef ||
    job.environmentRef !== request.environmentRef ||
    job.timeStart !== request.timeStart ||
    job.timeEnd !== request.timeEnd ||
    job.correlationId !== request.correlationId ||
    job.queryVersion !== request.queryVersion ||
    !sameStringArray(job.sources, request.sources)
  ) {
    throw invalid("job-scope-mismatch", "server Search Job scope/provenance does not match request");
  }
}

function validateOutcomeCardinality(job, projections) {
  switch (job.state) {
    case "queued":
    case "running":
      if (job.resultRefs.length !== 0 || job.failedSources.length !== 0 || projections.length !== 0) {
        throw invalid("premature-results", "non-terminal jobs cannot expose terminal results");
      }
      break;
    case "completed":
      if (job.failedSources.length !== 0) {
        throw invalid("completed-with-failures", "completed job cannot name failed sources");
      }
      break;
    case "partial":
      if (job.failedSources.length === 0 || job.resultRefs.length === 0) {
        throw invalid("invalid-partial", "partial job requires valid results and failed sources");
      }
      break;
    case "failed":
      if (job.resultRefs.length !== 0 || projections.length !== 0) {
        throw invalid("failed-with-results", "failed job cannot expose result projections");
      }
      break;
    case "cancelled":
      if (job.resultRefs.length !== 0 || job.failedSources.length !== 0 || projections.length !== 0) {
        throw invalid("cancelled-with-results", "cancelled job cannot expose result projections");
      }
      break;
  }
}

function validateProjectionsMatchRefs(refs, projections) {
  if (refs.length !== projections.length) {
    throw invalid("projection-ref-mismatch", "every server result ref requires exactly one authorized projection");
  }
  const expected = new Set(refs.map((item) => `${item.source}\u0000${item.ref}`));
  for (const projection of projections) {
    const key = `${projection.source}\u0000${projection.ref}`;
    if (!expected.delete(key)) {
      throw invalid("projection-ref-mismatch", "projection is not authorized by server result refs");
    }
  }
  if (expected.size !== 0) {
    throw invalid("projection-ref-mismatch", "server result refs are missing presentation projections");
  }
}

function normalizeResultRefs(value) {
  if (!Array.isArray(value) || value.length > MAX_RESULTS) {
    throw invalid("invalid-result-refs", "result refs must be a bounded array");
  }
  const seen = new Set();
  const out = value.map((item) => {
    if (!isPlainObject(item)) {
      throw invalid("invalid-result-ref", "result ref must be a plain object");
    }
    rejectUnknown(item, new Set(["source", "ref"]));
    const normalized = {
      source: opaqueRef("result source", item.source),
      ref: opaqueRef("result ref", item.ref),
    };
    const key = `${normalized.source}\u0000${normalized.ref}`;
    if (seen.has(key)) {
      throw invalid("duplicate-result-ref", "duplicate server result ref");
    }
    seen.add(key);
    return freezeDeep(normalized);
  });
  return Object.freeze(out);
}

function normalizeFailedSources(value) {
  if (!Array.isArray(value) || value.length > MAX_FAILURES) {
    throw invalid("invalid-failed-sources", "failed sources must be a bounded array");
  }
  const seen = new Set();
  const out = value.map((item) => {
    if (!isPlainObject(item)) {
      throw invalid("invalid-source-failure", "source failure must be a plain object");
    }
    rejectUnknown(item, new Set(["source", "code"]));
    const source = opaqueRef("failed source", item.source);
    if (seen.has(source)) {
      throw invalid("duplicate-failed-source", `duplicate failed source: ${source}`);
    }
    seen.add(source);
    return freezeDeep({ source, code: opaqueRef("failed source code", item.code) });
  });
  return Object.freeze(out);
}

function normalizeFreshness(value) {
  if (!Array.isArray(value) || value.length === 0) {
    throw invalid("invalid-freshness", "partial outcome requires freshness evidence");
  }
  const seen = new Set();
  const out = value.map((item) => {
    if (!isPlainObject(item)) {
      throw invalid("invalid-freshness", "freshness item must be a plain object");
    }
    rejectUnknown(item, new Set(["source", "observedAt"]));
    const source = opaqueRef("freshness source", item.source);
    if (seen.has(source)) {
      throw invalid("duplicate-freshness-source", `duplicate freshness source: ${source}`);
    }
    seen.add(source);
    return freezeDeep({ source, observedAt: normalizeInstant("freshness observedAt", item.observedAt) });
  });
  return Object.freeze(out.sort((a, b) => a.source.localeCompare(b.source)));
}

function normalizeOpaqueList(label, value, maxItems) {
  if (!Array.isArray(value) || value.length === 0 || value.length > maxItems) {
    throw invalid("invalid-list", `${label} list must be non-empty and bounded`);
  }
  const seen = new Set();
  const out = [];
  for (const item of value) {
    const normalized = opaqueRef(label, item);
    if (seen.has(normalized)) {
      throw invalid("duplicate-list-value", `duplicate ${label}: ${normalized}`);
    }
    seen.add(normalized);
    out.push(normalized);
  }
  return Object.freeze(out.sort());
}

function sameExecutionIntent(left, right) {
  if (left === undefined || right === undefined) {
    return false;
  }
  return (
    left.tenantRef === right.tenantRef &&
    left.environmentRef === right.environmentRef &&
    left.timeStart === right.timeStart &&
    left.timeEnd === right.timeEnd &&
    left.query === right.query &&
    JSON.stringify(left.filters) === JSON.stringify(right.filters) &&
    sameStringArray(left.sources, right.sources)
  );
}

function sameStringArray(left, right) {
  return (
    Array.isArray(left) &&
    Array.isArray(right) &&
    left.length === right.length &&
    left.every((value, index) => value === right[index])
  );
}

function normalizeInstant(label, value) {
  const text = requiredText(label, value, 128);
  const millis = Date.parse(text);
  if (!Number.isFinite(millis)) {
    throw invalid("invalid-time", `${label} must be an RFC3339-compatible instant`);
  }
  return new Date(millis).toISOString();
}

function opaqueRef(label, value) {
  const text = requiredText(label, value, 256);
  if (!OPAQUE_REF.test(text)) {
    throw invalid("invalid-reference", `${label} is not a safe opaque reference`);
  }
  return text;
}

function optionalOpaqueRef(label, value) {
  if (value === undefined || value === null || value === "") {
    return undefined;
  }
  return opaqueRef(label, value);
}

function requiredText(label, value, maxLength) {
  if (typeof value !== "string") {
    throw invalid("invalid-text", `${label} must be a string`);
  }
  const normalized = value.trim();
  if (
    normalized === "" ||
    normalized.length > maxLength ||
    /[\u0000-\u001f\u007f]/u.test(normalized)
  ) {
    throw invalid("invalid-text", `${label} must be non-empty, bounded and control-free`);
  }
  return normalized;
}

function boundedInteger(label, value, min, max) {
  if (!Number.isSafeInteger(value) || value < min || value > max) {
    throw invalid("invalid-integer", `${label} must be an integer from ${min} to ${max}`);
  }
  return value;
}

function rejectUnknown(input, allowed) {
  for (const key of Object.keys(input)) {
    if (!allowed.has(key)) {
      throw invalid("unknown-field", `unknown field: ${key}`);
    }
  }
}

function isPlainObject(value) {
  if (value === null || typeof value !== "object" || Array.isArray(value)) {
    return false;
  }
  const prototype = Object.getPrototypeOf(value);
  return prototype === Object.prototype || prototype === null;
}

function isAbortError(error) {
  return error instanceof Error && error.name === "AbortError";
}

function freezeDeep(value) {
  if (value && typeof value === "object" && !Object.isFrozen(value)) {
    for (const child of Object.values(value)) {
      freezeDeep(child);
    }
    Object.freeze(value);
  }
  return value;
}

function invalid(code, message) {
  return new IntegrationValidationError(code, message);
}
