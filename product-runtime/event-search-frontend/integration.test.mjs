import assert from "node:assert/strict";
import test from "node:test";

import {
  IntegrationValidationError,
  createEventSearchController,
  createExecutionRequest,
  normalizeTransportOutcome,
} from "./integration.mjs";

function deepLink(overrides = {}) {
  return {
    tenantId: "tenant-a",
    environmentId: "prod-eu",
    from: "2026-09-24T10:00:00Z",
    to: "2026-09-24T11:00:00Z",
    query: "backend-neutral opaque query",
    sources: ["firewall", "edr"],
    filters: { severity: ["high"] },
    ...overrides,
  };
}

function execution(overrides = {}) {
  return {
    correlationId: "corr-001",
    queryVersion: "query-v1",
    ...overrides,
  };
}

function request() {
  return createExecutionRequest(deepLink(), execution(), {
    tenantId: "tenant-a",
    environmentId: "prod-eu",
  });
}

function job(overrides = {}) {
  return {
    kind: "job",
    jobId: "job-001",
    runId: "run-001",
    state: "completed",
    tenantRef: "tenant-a",
    environmentRef: "prod-eu",
    timeStart: "2026-09-24T10:00:00.000Z",
    timeEnd: "2026-09-24T11:00:00.000Z",
    sources: ["edr", "firewall"],
    correlationId: "corr-001",
    queryVersion: "query-v1",
    resultRefs: [],
    failedSources: [],
    projections: [],
    ...overrides,
  };
}

function projection(ref, source = "edr") {
  return {
    ref,
    source,
    eventTime: "2026-09-24T10:30:00Z",
    summary: { severity: "high" },
  };
}

test("execution-request-carries-scope-query-and-provenance-but-never-client-authorization", () => {
  const value = createExecutionRequest(
    deepLink(),
    execution(),
    { tenantId: "tenant-a", environmentId: "prod-eu" },
  );
  assert.deepEqual(Object.keys(value), [
    "tenantRef",
    "environmentRef",
    "timeStart",
    "timeEnd",
    "sources",
    "query",
    "filters",
    "correlationId",
    "queryVersion",
  ]);
  assert.equal("permissions" in value, false);
  assert.equal("authorizedTenants" in value, false);
  assert.equal("authorization" in value, false);
  assert.throws(() => createExecutionRequest(
    { ...deepLink(), authorization: ["admin"] },
    execution(),
    { tenantId: "tenant-a", environmentId: "prod-eu" },
  ));
});

test("transport-tenant-mismatch-fails-closed", () => {
  assert.throws(
    () => normalizeTransportOutcome(job({ tenantRef: "tenant-b" }), request()),
    (error) => error instanceof IntegrationValidationError && error.code === "job-scope-mismatch",
  );
});

test("permission-denied-outcome-clears-results-and-exposes-no-protected-projection", async () => {
  const views = [];
  const controller = createEventSearchController({
    transport: {
      async execute(value) {
        return {
          kind: "permission-denied",
          correlationId: value.correlationId,
          capability: "perm.investigate.search.execute",
          reasonCode: "POLICY_DENY",
          requestAccessPath: "/access/request",
        };
      },
      async retry() {
        throw new Error("not used");
      },
    },
    onView(view) {
      views.push(view);
    },
    expectedContext: { tenantId: "tenant-a", environmentId: "prod-eu" },
  });
  const outcome = await controller.execute(deepLink(), execution());
  assert.equal(outcome.kind, "permission-denied");
  const view = controller.getView();
  assert.equal(view.shellState.kind, "permission-denied");
  assert.equal(view.resultWindow, undefined);
  assert.equal(JSON.stringify(view).includes("backend-neutral opaque query"), false);
  assert.equal(JSON.stringify(view).includes("POLICY_DENY"), true);
  assert.equal(views.at(-1), view);
});

test("completed-job-with-authorized-projections-produces-window-without-invented-success-state", async () => {
  const controller = createController(async (value) => job({
    correlationId: value.correlationId,
    queryVersion: value.queryVersion,
    resultRefs: [{ source: "edr", ref: "evt-1" }],
    projections: [projection("evt-1")],
  }));
  const outcome = await controller.execute(deepLink(), execution());
  assert.equal(outcome.job.state, "completed");
  const view = controller.getView();
  assert.equal(view.shellState, null);
  assert.deepEqual(view.resultWindow.rows.map((item) => item.ref), ["evt-1"]);
});

test("completed-job-with-no-results-maps-to-empty-state", async () => {
  const controller = createController(async (value) => job({
    correlationId: value.correlationId,
    queryVersion: value.queryVersion,
  }));
  await controller.execute(deepLink(), execution());
  assert.equal(controller.getView().shellState.kind, "empty");
  assert.equal(controller.getView().resultWindow, undefined);
});

test("partial-job-preserves-valid-results-and-names-failed-sources", async () => {
  const controller = createController(async (value) => job({
    state: "partial",
    correlationId: value.correlationId,
    queryVersion: value.queryVersion,
    resultRefs: [{ source: "edr", ref: "evt-1" }],
    failedSources: [{ source: "firewall", code: "SOURCE_TIMEOUT" }],
    projections: [projection("evt-1")],
    freshness: [{ source: "edr", observedAt: "2026-09-24T10:59:00Z" }],
    consequence: "Firewall telemetry is incomplete.",
  }));
  const outcome = await controller.execute(deepLink(), execution());
  assert.equal(outcome.job.state, "partial");
  const view = controller.getView();
  assert.equal(view.shellState.kind, "partial");
  assert.deepEqual(view.shellState.failedSources, ["firewall"]);
  assert.deepEqual(view.resultWindow.rows.map((item) => item.ref), ["evt-1"]);
});

test("projection-must-match-server-result-ref-and-source", () => {
  assert.throws(
    () => normalizeTransportOutcome(job({
      resultRefs: [{ source: "edr", ref: "evt-1" }],
      projections: [projection("evt-2")],
    }), request()),
    (error) => error instanceof IntegrationValidationError && error.code === "projection-ref-mismatch",
  );
  assert.throws(
    () => normalizeTransportOutcome(job({
      resultRefs: [{ source: "firewall", ref: "evt-1" }],
      projections: [projection("evt-1", "edr")],
    }), request()),
    (error) => error instanceof IntegrationValidationError && error.code === "projection-ref-mismatch",
  );
});

test("failed-job-never-carries-results-and-maps-to-safe-error", () => {
  const outcome = normalizeTransportOutcome(job({
    state: "failed",
    failedSources: [
      { source: "edr", code: "BACKEND_FAILURE" },
      { source: "firewall", code: "BACKEND_FAILURE" },
    ],
  }), request());
  assert.equal(outcome.shellState.kind, "error");
  assert.equal(outcome.shellState.correlationId, "corr-001");
  assert.equal(JSON.stringify(outcome).includes("raw backend stack"), false);
  assert.throws(() => normalizeTransportOutcome(job({
    state: "failed",
    resultRefs: [{ source: "edr", ref: "evt-1" }],
    projections: [projection("evt-1")],
  }), request()));
});

test("queued-and-running-jobs-cannot-prematurely-expose-results", () => {
  for (const state of ["queued", "running"]) {
    const outcome = normalizeTransportOutcome(job({ state }), request());
    assert.equal(outcome.terminal, false);
    assert.throws(() => normalizeTransportOutcome(job({
      state,
      resultRefs: [{ source: "edr", ref: "evt-1" }],
      projections: [projection("evt-1")],
    }), request()));
  }
});

test("AbortController-cancellation-aborts-injected-transport-without-selecting-a-provider", async () => {
  let signal;
  let release;
  const controller = createEventSearchController({
    transport: {
      execute(_request, options) {
        signal = options.signal;
        return new Promise((_resolve, reject) => {
          release = () => {
            const error = new Error("aborted transport");
            error.name = "AbortError";
            reject(error);
          };
          signal.addEventListener("abort", release, { once: true });
        });
      },
      async retry() {
        throw new Error("not used");
      },
    },
    expectedContext: { tenantId: "tenant-a", environmentId: "prod-eu" },
  });
  const pending = controller.execute(deepLink(), execution());
  assert.equal(signal.aborted, false);
  assert.equal(controller.cancel(), true);
  assert.equal(signal.aborted, true);
  assert.equal(controller.cancel(), false);
  const outcome = await pending;
  assert.equal(outcome.kind, "cancelled");
});

test("new-execution-supersedes-stale-response", async () => {
  const pending = [];
  const controller = createEventSearchController({
    transport: {
      execute(value) {
        return new Promise((resolve) => pending.push({ value, resolve }));
      },
      async retry() {
        throw new Error("not used");
      },
    },
    expectedContext: { tenantId: "tenant-a", environmentId: "prod-eu" },
  });
  const first = controller.execute(deepLink({ query: "first" }), execution({ correlationId: "corr-first" }));
  const second = controller.execute(deepLink({ query: "second" }), execution({ correlationId: "corr-second" }));
  pending[1].resolve(job({
    correlationId: "corr-second",
    resultRefs: [{ source: "edr", ref: "evt-second" }],
    projections: [projection("evt-second")],
  }));
  const secondOutcome = await second;
  pending[0].resolve(job({ correlationId: "corr-first" }));
  const firstOutcome = await first;
  assert.equal(secondOutcome.job.jobId, "job-001");
  assert.equal(firstOutcome.kind, "superseded");
  assert.deepEqual(controller.getView().resultWindow.rows.map((item) => item.ref), ["evt-second"]);
});

test("retry-as-new-run-requires-new-job-run-and-parent-link", async () => {
  let retryInput;
  const controller = createEventSearchController({
    transport: {
      async execute(value) {
        return job({
          state: "partial",
          correlationId: value.correlationId,
          queryVersion: value.queryVersion,
          resultRefs: [{ source: "edr", ref: "evt-1" }],
          failedSources: [{ source: "firewall", code: "TIMEOUT" }],
          projections: [projection("evt-1")],
          freshness: [{ source: "edr", observedAt: "2026-09-24T10:59:00Z" }],
          consequence: "Firewall data is incomplete.",
        });
      },
      async retry(value) {
        retryInput = value;
        return job({
          jobId: "job-002",
          runId: "run-002",
          parentRunId: "run-001",
          resultRefs: [{ source: "edr", ref: "evt-2" }],
          projections: [projection("evt-2")],
        });
      },
    },
    expectedContext: { tenantId: "tenant-a", environmentId: "prod-eu" },
  });
  await controller.execute(deepLink(), execution());
  const retried = await controller.retry();
  assert.equal(retryInput.jobId, "job-001");
  assert.equal(retryInput.runId, "run-001");
  assert.equal(retried.job.jobId, "job-002");
  assert.equal(retried.job.runId, "run-002");
  assert.equal(retried.job.parentRunId, "run-001");
});

test("retry-rejects-reused-job-or-run-identities", async () => {
  const requestValue = request();
  assert.throws(() => normalizeTransportOutcome(job({
    parentRunId: "run-001",
    runId: "run-001",
  }), requestValue, { retryOfRunId: "run-001" }));

  const controller = createEventSearchController({
    transport: {
      async execute(value) {
        return job({ correlationId: value.correlationId, queryVersion: value.queryVersion });
      },
      async retry() {
        return job({ parentRunId: "run-001", runId: "run-002" });
      },
    },
    expectedContext: { tenantId: "tenant-a", environmentId: "prod-eu" },
  });
  await controller.execute(deepLink(), execution());
  const outcome = await controller.retry();
  assert.equal(outcome.kind, "error");
});

test("raw-transport-errors-are-not-reflected-to-view", async () => {
  const controller = createController(async () => {
    throw new Error("postgres password=do-not-leak");
  });
  const outcome = await controller.execute(deepLink(), execution());
  assert.equal(outcome.kind, "error");
  const serialized = JSON.stringify(controller.getView());
  assert.equal(serialized.includes("do-not-leak"), false);
  assert.equal(serialized.includes("corr-001"), true);
});

test("permission-denial-cannot-carry-result-data-as-unknown-fields", () => {
  assert.throws(() => normalizeTransportOutcome({
    kind: "permission-denied",
    correlationId: "corr-001",
    capability: "perm.investigate.search.execute",
    reasonCode: "DENY",
    projections: [projection("secret")],
  }, request()));
});

test("partial-outcome-requires-freshness-consequence-and-valid-failure-scope", () => {
  assert.throws(() => normalizeTransportOutcome(job({
    state: "partial",
    resultRefs: [{ source: "edr", ref: "evt-1" }],
    projections: [projection("evt-1")],
    failedSources: [{ source: "firewall", code: "TIMEOUT" }],
  }), request()));
  assert.throws(() => normalizeTransportOutcome(job({
    state: "partial",
    resultRefs: [{ source: "edr", ref: "evt-1" }],
    projections: [projection("evt-1")],
    failedSources: [{ source: "outside", code: "TIMEOUT" }],
    freshness: [{ source: "edr", observedAt: "2026-09-24T10:59:00Z" }],
    consequence: "Incomplete",
  }), request()));
});

test("invalid-controller-and-retry-before-terminal-fail-closed", async () => {
  assert.throws(() => createEventSearchController({ transport: {} }));
  const controller = createController(async (value) => job({
    state: "running",
    correlationId: value.correlationId,
    queryVersion: value.queryVersion,
  }));
  await controller.execute(deepLink(), execution());
  await assert.rejects(() => controller.retry(), IntegrationValidationError);
});

function createController(execute) {
  return createEventSearchController({
    transport: {
      execute,
      async retry() {
        throw new Error("retry not configured");
      },
    },
    expectedContext: { tenantId: "tenant-a", environmentId: "prod-eu" },
  });
}
