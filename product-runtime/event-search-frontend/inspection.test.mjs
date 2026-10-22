import assert from "node:assert/strict";
import test from "node:test";
import { createEventSearchController } from "./integration.mjs";
import { createEventInspectionJourney, normalizeInspection, renderEventInspection } from "./inspection.mjs";

const context = {
  tenantId: "tenant-a", environmentId: "prod", from: "2026-10-20T10:00:00Z",
  to: "2026-10-20T12:00:00Z", sources: ["edr"], query: "source:edr",
};

function search() {
  return createEventSearchController({
    expectedContext: { tenantId: "tenant-a", environmentId: "prod" },
    transport: {
      async execute(request) {
        return {
          kind: "job", jobId: "job-1", runId: "run-1",
          state: "completed", tenantRef: request.tenantRef,
          environmentRef: request.environmentRef, timeStart: request.timeStart,
          timeEnd: request.timeEnd, sources: [...request.sources],
          correlationId: request.correlationId, queryVersion: request.queryVersion,
          resultRefs: [{ source: "edr", ref: "event-1" }],
          failedSources: [],
          projections: [{
            ref: "event-1", eventTime: "2026-10-20T11:00:00Z",
            source: "edr", summary: { severity: "high" },
          }],
        };
      },
      async retry() { throw new Error("not used"); },
    },
  });
}

function inspection(overrides = {}) {
  return {
    kind: "inspection", tenantRef: "tenant-a", environmentRef: "prod",
    eventRef: "event-1", correlationId: "corr-1", state: "partial",
    source: { sourceRef: "shared-source-1", dataSourceRef: "settings-source-1" },
    parser: { parserRef: "parser-1", version: "v2", status: "degraded" },
    sourceFields: { "source.ip": "192.0.2.10" },
    normalizedFields: { "network.client.ip": "192.0.2.10" },
    missingFields: ["network.server.ip"],
    enrichments: [{ producer: "intel-engine", version: "1", freshness: "stale",
      stale: true, values: { "risk.score": "72" } }],
    raw: { access: "denied" }, rendered: { access: "allowed", text: "safe rendered" },
    ...overrides,
  };
}

async function initializedSearch() {
  const controller = search();
  await controller.execute(context, { correlationId: "corr-1", queryVersion: "query-1" });
  assert.equal(controller.getView().resultWindow.rows[0].ref, "event-1");
  return controller;
}

test("real Event Search controller result opens authorized inspection and restores result window", async () => {
  const controller = await initializedSearch();
  const initial = controller.getView();
  let seen;
  const journey = createEventInspectionJourney({
    searchController: controller,
    transport: { async inspect(request) {
      assert.deepEqual(Object.keys(request).sort(),
        ["tenantRef", "environmentRef", "eventRef", "correlationId", "jobId", "runId", "sourceHint"].sort());
      assert.equal(request.eventRef, "event-1");
      assert.equal(request.runId, "run-1");
      assert.equal(request.sourceHint, "edr");
      return inspection();
    } },
    onView(v) { seen = v; },
  });
  const result = await journey.open(context, "event-1");
  assert.equal(result.kind, "inspected");
  assert.equal(seen.inspection.raw.access, "denied");
  assert.equal(seen.inspection.rendered.text, "safe rendered");
  assert.equal(seen.inspection.enrichments[0].stale, true);
  assert.ok(Object.isFrozen(seen.inspection));
  const back = journey.back();
  assert.equal(back.kind, "search");
  assert.equal(back.resultWindow, initial.resultWindow);
  assert.equal(back.searchState.query, context.query);
  assert.equal(back.searchState.from, new Date(context.from).toISOString());
  assert.deepEqual(controller.getView(), initial);
});

test("missing selected Event and forged scope fail closed before or after transport", async () => {
  const controller = await initializedSearch();
  let calls = 0;
  const journey = createEventInspectionJourney({
    searchController: controller,
    transport: { async inspect() { calls++; return inspection({ tenantRef: "foreign" }); } },
  });
  await assert.rejects(() => journey.open(context, "event-foreign"), /not-authorized-result/);
  assert.equal(calls, 0);
  const invalid = await journey.open(context, "event-1");
  assert.equal(invalid.kind, "unavailable");
  assert.equal(journey.getView().inspection, undefined);
  assert.equal(calls, 1);
});

test("Inspector opening rejects a caller-supplied time/source scope that differs from the actual Search Job", async () => {
  const controller = await initializedSearch();
  let inspectCalls = 0;
  const journey = createEventInspectionJourney({
    searchController: controller,
    transport: { async inspect() { inspectCalls++; return inspection(); } },
  });
  for (const mismatched of [
    { ...context, from: "2026-10-20T09:00:00Z" },
    { ...context, to: "2026-10-20T13:00:00Z" },
    { ...context, sources: ["edr", "firewall"] },
  ]) {
    await assert.rejects(() => journey.open(mismatched, "event-1"), /not-authorized-result/);
  }
  assert.equal(inspectCalls, 0, "no host call should be made for an inconsistent caller scope");
  assert.equal(journey.getView().kind, "search");
});

test("Inspector rejects reused Search Job/run IDs after scope, version or result projection changes", async (t) => {
  const scenarios = [
    { label: "time", change: { ...context, from: "2026-10-20T09:00:00Z" }, version: "query-1" },
    { label: "sources", change: { ...context, sources: ["edr", "firewall"] }, version: "query-1" },
    { label: "version", change: context, version: "query-2" },
  ];
  for (const scenario of scenarios) {
    await t.test(scenario.label, async () => {
      const controller = await initializedSearch();
      const journey = createEventInspectionJourney({
        searchController: controller,
        transport: { async inspect() { return inspection(); } },
      });
      assert.equal((await journey.open(context, "event-1")).kind, "inspected");
      await controller.execute(scenario.change, { correlationId: "corr-1", queryVersion: scenario.version });
      assert.equal(controller.getView().job.jobId, "job-1");
      assert.equal(controller.getView().job.runId, "run-1");
      assert.equal(controller.getView().job.correlationId, "corr-1");
      const invalidated = journey.getView();
      assert.equal(invalidated.kind, "unavailable");
      assert.equal(invalidated.inspection, undefined);
      assert.equal(journey.back().kind, "unavailable");
    });
  }
});

test("raw denial strips raw-derived content, rejects forged raw and keeps rendered accessible", async () => {
  const controller = await initializedSearch();
  const request = {
    tenantRef: "tenant-a", environmentRef: "prod", eventRef: "event-1",
    correlationId: "corr-1",
  };
  assert.throws(() => normalizeInspection(inspection({
    raw: { access: "denied", text: "sensitive" },
  }), request), /denied-payload-present/);
  assert.throws(() => normalizeInspection(inspection({
    rawDerivedSensitive: { "secret.token": "do-not-show" },
  }), request), /raw-derived-without-access/);
  const journey = createEventInspectionJourney({
    searchController: controller,
    transport: { async inspect() { return inspection(); } },
  });
  const result = await journey.open(context, "event-1");
  const mount = new FakeElement("div");
  renderEventInspection(new FakeDocument(), mount, result);
  const visible = allText(mount).join(" ");
  assert.equal(visible.includes("safe rendered"), true);
  assert.equal(visible.includes("Raw unavailable"), true);
  assert.equal(visible.includes("sensitive"), false);
  assert.equal(visible.includes("intel-engine@1"), true);
  assert.equal(visible.includes("192.0.2.10"), true);
  assert.equal(visible.includes("Missing fields"), true);
  assert.equal(visible.includes("network.server.ip"), true);
});

test("raw-only authorized inspection never invents rendered or raw-derived values", async () => {
  const controller = await initializedSearch();
  const journey = createEventInspectionJourney({
    searchController: controller,
    transport: { async inspect() {
      return inspection({ raw: { access: "allowed", text: "raw-safe" },
        rendered: { access: "denied" },
        rawDerivedSensitive: { "raw.trace": "raw-only-value" },
      });
    } },
  });
  const result = await journey.open(context, "event-1");
  assert.equal(result.kind, "inspected");
  assert.equal(result.inspection.raw.text, "raw-safe");
  assert.equal(result.inspection.rendered.access, "denied");
  const mount = new FakeElement("div");
  renderEventInspection(new FakeDocument(), mount, result);
  const visible = allText(mount).join(" ");
  assert.equal(visible.includes("raw-only-value"), true);
  assert.equal(visible.includes("Rendered unavailable"), true);
  assert.equal(visible.includes("safe rendered"), false);
});

test("hostile source fields and stale authorized results do not create HTML or leak stale projection", async () => {
  const controller = await initializedSearch();
  const hostile = "<img src=x onerror=alert(1)>";
  const journey = createEventInspectionJourney({
    searchController: controller,
    transport: { async inspect() { return inspection({ rendered: {
      access: "allowed", text: hostile,
    } }); } },
  });
  const result = await journey.open(context, "event-1");
  const mount = new FakeElement("div");
  renderEventInspection(new FakeDocument(), mount, result);
  assert.equal(allText(mount).some(x => x.includes(hostile)), true);
  assert.equal(allNodes(mount).some(x => x.tagName === "img"), false);
  assert.equal(allNodes(mount).some(x => x.tagName === "script"), false);
  journey.close();
  const reset = journey.getView();
  renderEventInspection(new FakeDocument(), mount, reset);
  assert.equal(allText(mount).some(x => x.includes(hostile)), false);
});

test("revocation during slow inspection prevents publication and back cannot restore revoked context", async () => {
  const controller = await initializedSearch();
  let finish;
  const pending = new Promise(resolve => { finish = resolve; });
  const journey = createEventInspectionJourney({
    searchController: controller,
    transport: { inspect() { return pending; } },
  });
  const open = journey.open(context, "event-1");
  await controller.execute(context, { correlationId: "corr-2", queryVersion: "query-1" });
  finish(inspection());
  const result = await open;
  assert.equal(result.kind, "unavailable");
  assert.equal(journey.back().kind, "unavailable");
  assert.equal(journey.getView().inspection, undefined);
});

test("revocation during a pending pivot never republishes an inspection, including on error", async (t) => {
  for (const outcome of ["rejected", "resolved"]) {
    await t.test(outcome, async () => {
      const controller = await initializedSearch();
      let finish, refuse;
      const pending = new Promise((resolve, reject) => { finish = resolve; refuse = reject; });
      const published = [];
      const journey = createEventInspectionJourney({
        searchController: controller,
        transport: {
          async inspect() { return inspection(); },
          draftPivot() { return pending; },
        },
        onView(view) { published.push(view); },
      });
      assert.equal((await journey.open(context, "event-1")).kind, "inspected");
      const inFlight = journey.draftPivot("source.ip", "192.0.2.10");
      assert.equal(journey.getView().pivotStatus, "loading");
      // A new Search run invalidates the original correlation and selected row.
      await controller.execute(context, { correlationId: "corr-2", queryVersion: "query-1" });
      if (outcome === "rejected") refuse(new Error("pivot unavailable"));
      else finish({ kind: "pivot-draft" });
      const outcomeView = await inFlight;
      assert.equal(outcomeView.kind, "unavailable");
      assert.equal(outcomeView.inspection, undefined);
      assert.equal(journey.getView().inspection, undefined);
      assert.equal(published.at(-1).kind, "unavailable");
      assert.equal(journey.back().kind, "unavailable");
      assert.equal(journey.getView().inspection, undefined);
    });
  }
});

test("pivot transport denial retains only a still-authorized inspection", async () => {
  const controller = await initializedSearch();
  const journey = createEventInspectionJourney({
    searchController: controller,
    transport: {
      async inspect() { return inspection(); },
      async draftPivot() { throw new Error("pivot action denied"); },
    },
  });
  assert.equal((await journey.open(context, "event-1")).kind, "inspected");
  const denied = await journey.draftPivot("source.ip", "192.0.2.10");
  assert.equal(denied.kind, "inspected");
  assert.equal(denied.pivotStatus, "unavailable");
  assert.equal(denied.pivotDraft, undefined);
});

test("reading an already-open Inspector after Search run replacement fails closed without a pivot", async () => {
  const controller = await initializedSearch();
  const published = [];
  const journey = createEventInspectionJourney({
    searchController: controller,
    transport: { async inspect() { return inspection(); } },
    onView(next) { published.push(next); },
  });
  assert.equal((await journey.open(context, "event-1")).kind, "inspected");
  assert.equal(journey.getView().kind, "inspected");
  await controller.execute(context, { correlationId: "corr-2", queryVersion: "query-1" });
  const current = journey.getView();
  assert.equal(current.kind, "unavailable");
  assert.equal(current.inspection, undefined);
  assert.equal(journey.getView(), current);
  assert.equal(published.at(-1), current);
  assert.equal(journey.back().kind, "unavailable");
});

test("reading a revoked Inspector during pending pivot aborts its stale response", async () => {
  const controller = await initializedSearch();
  let finish;
  const pending = new Promise(resolve => { finish = resolve; });
  const journey = createEventInspectionJourney({
    searchController: controller,
    transport: {
      async inspect() { return inspection(); },
      draftPivot() { return pending; },
    },
  });
  assert.equal((await journey.open(context, "event-1")).kind, "inspected");
  const pivot = journey.draftPivot("source.ip", "192.0.2.10");
  assert.equal(journey.getView().pivotStatus, "loading");
  await controller.execute(context, { correlationId: "corr-2", queryVersion: "query-1" });
  assert.equal(journey.getView().kind, "unavailable");
  finish({ kind: "pivot-draft" });
  assert.equal((await pivot).kind, "superseded");
  assert.equal(journey.getView().inspection, undefined);
});

test("explicit Inspector-to-Search handoff reauthorizes the unexecuted pivot and restores the original window", async () => {
  const controller = await initializedSearch();
  const originalWindow = controller.getView().resultWindow;
  const transportCalls = [];
  let inspectionCalls = 0;
  const journey = createEventInspectionJourney({
    searchController: controller,
    transport: {
      async inspect() { inspectionCalls++; return inspection(); },
      async draftPivot(request) {
        transportCalls.push(request);
        return {
          kind: "pivot-draft", targetCapability: "CAP-INV-002",
          auditRequired: true, tenantRef: request.tenantRef,
          environmentRef: request.environmentRef, eventRef: request.eventRef,
          correlationId: request.correlationId, jobId: request.jobId,
          runId: request.runId, timeRange: { start: request.timeStart, end: request.timeEnd },
          field: request.field, value: request.value, fieldOrigin: "source",
          sourceRef: "shared-source-1", dataSourceRef: "settings-source-1",
          returnContext: { origin: "event-search", tenantRef: request.tenantRef,
            environmentRef: request.environmentRef, eventRef: request.eventRef,
            correlationId: request.correlationId },
        };
      },
    },
  });
  await journey.open(context, "event-1");
  assert.equal((await journey.draftPivot("source.ip", "192.0.2.10")).kind, "inspected");
  const surface = new FakeElement("main");
  let continuation;
  renderEventInspection(new FakeDocument(), surface, journey.getView(), {
    onContinuePivotInSearch() { continuation = journey.continuePivotInSearch(); },
  });
  const next = allNodes(surface).find(node => node.tagName === "button" &&
    node.textContent.includes("Continue with draft in Event Search"));
  assert.ok(next, "visible explicit handoff action");
  next.click();
  const result = await continuation;
  assert.equal(result.kind, "search");
  assert.equal(result.queryDraft.kind, "pivot-draft");
  assert.equal(result.queryDraft.fieldOrigin, "source");
  assert.equal(result.queryDraft.returnContext.eventRef, "event-1");
  assert.equal(result.resultWindow, originalWindow);
  assert.equal(result.searchState.query, context.query);
  assert.equal(transportCalls.length, 2, "a fresh host authorization is required at handoff");
  assert.deepEqual(Object.keys(result).sort(),
    ["kind", "searchState", "resultWindow", "queryDraft"].sort());
  const searchPanel = new FakeElement("main");
  let returned;
  renderEventInspection(new FakeDocument(), searchPanel, result, {
    onReturnToInspectedEvent() { returned = journey.returnToInspectedEvent(); },
  });
  assert.ok(allText(searchPanel).some(x => x.includes("unexecuted Query draft")));
  const returnAction = allNodes(searchPanel).find(node => node.tagName === "button" &&
    node.textContent === "Return to inspected Event");
  assert.ok(returnAction, "explicit return to the same Event");
  returnAction.click();
  const reinspected = await returned;
  assert.equal(reinspected.kind, "inspected");
  assert.equal(reinspected.inspection.eventRef, "event-1");
  assert.equal(reinspected.pivotDraft, undefined, "old draft is not a grant for this inspection");
  assert.equal(inspectionCalls, 2, "return MUST reauthorize against the Telemetry owner");
  const restoredSearch = journey.back();
  assert.equal(restoredSearch.kind, "search");
  assert.equal(restoredSearch.resultWindow, originalWindow);
  assert.equal(restoredSearch.searchState.query, context.query);
  assert.equal(restoredSearch.queryDraft, undefined);
});

test("pivot handoff denied by the owner on reauthorization must not enter Event Search", async () => {
  const controller = await initializedSearch();
  let calls = 0;
  const journey = createEventInspectionJourney({
    searchController: controller,
    transport: {
      async inspect() { return inspection(); },
      async draftPivot(request) {
        calls++;
        if (calls === 2) throw new Error("pivot authorization revoked");
        return {
          kind: "pivot-draft", targetCapability: "CAP-INV-002", auditRequired: true,
          tenantRef: request.tenantRef, environmentRef: request.environmentRef,
          eventRef: request.eventRef, correlationId: request.correlationId,
          jobId: request.jobId, runId: request.runId,
          timeRange: { start: request.timeStart, end: request.timeEnd },
          field: request.field, value: request.value, fieldOrigin: "source",
          sourceRef: "shared-source-1", dataSourceRef: "settings-source-1",
          returnContext: { origin: "event-search", tenantRef: request.tenantRef,
            environmentRef: request.environmentRef, eventRef: request.eventRef,
            correlationId: request.correlationId },
        };
      },
    },
  });
  await journey.open(context, "event-1");
  await assert.rejects(() => journey.continuePivotInSearch(), /pivot-not-available/);
  assert.equal((await journey.draftPivot("source.ip", "192.0.2.10")).kind, "inspected");
  const rejected = await journey.continuePivotInSearch();
  assert.equal(rejected.kind, "inspected");
  assert.equal(rejected.pivotStatus, "unavailable");
  assert.equal(rejected.pivotDraft, undefined);
  assert.equal(journey.getView().queryDraft, undefined);
  assert.equal(calls, 2);
});

test("returning from a pivot draft cannot reuse an Inspector when host revokes Event.read", async () => {
  const controller = await initializedSearch();
  let inspected = 0;
  const journey = createEventInspectionJourney({
    searchController: controller,
    transport: {
      async inspect() {
        inspected++;
        return inspected === 1 ? inspection() : { kind: "permission-denied" };
      },
      async draftPivot(request) {
        return {
          kind: "pivot-draft", targetCapability: "CAP-INV-002", auditRequired: true,
          tenantRef: request.tenantRef, environmentRef: request.environmentRef,
          eventRef: request.eventRef, correlationId: request.correlationId,
          jobId: request.jobId, runId: request.runId,
          timeRange: { start: request.timeStart, end: request.timeEnd },
          field: request.field, value: request.value, fieldOrigin: "source",
          sourceRef: "shared-source-1", dataSourceRef: "settings-source-1",
          returnContext: { origin: "event-search", tenantRef: request.tenantRef,
            environmentRef: request.environmentRef, eventRef: request.eventRef,
            correlationId: request.correlationId },
        };
      },
    },
  });
  await assert.rejects(() => journey.returnToInspectedEvent(), /event-return-not-available/);
  await journey.open(context, "event-1");
  await journey.draftPivot("source.ip", "192.0.2.10");
  assert.equal((await journey.continuePivotInSearch()).kind, "search");
  const denied = await journey.returnToInspectedEvent();
  assert.equal(denied.kind, "unavailable");
  assert.equal(journey.getView().inspection, undefined);
  assert.equal(journey.getView().queryDraft, undefined);
  assert.equal(inspected, 2);
});

test("revocation after entering Search with a draft removes the draft before next read", async () => {
  const controller = await initializedSearch();
  const journey = createEventInspectionJourney({
    searchController: controller,
    transport: {
      async inspect() { return inspection(); },
      async draftPivot(request) {
        return {
          kind: "pivot-draft", targetCapability: "CAP-INV-002", auditRequired: true,
          tenantRef: request.tenantRef, environmentRef: request.environmentRef,
          eventRef: request.eventRef, correlationId: request.correlationId,
          jobId: request.jobId, runId: request.runId,
          timeRange: { start: request.timeStart, end: request.timeEnd },
          field: request.field, value: request.value, fieldOrigin: "source",
          sourceRef: "shared-source-1", dataSourceRef: "settings-source-1",
          returnContext: { origin: "event-search", tenantRef: request.tenantRef,
            environmentRef: request.environmentRef, eventRef: request.eventRef,
            correlationId: request.correlationId },
        };
      },
    },
  });
  await journey.open(context, "event-1");
  await journey.draftPivot("source.ip", "192.0.2.10");
  assert.equal((await journey.continuePivotInSearch()).kind, "search");
  await controller.execute(context, { correlationId: "corr-2", queryVersion: "query-1" });
  assert.equal(journey.getView().kind, "unavailable");
  assert.equal(journey.getView().queryDraft, undefined);
  assert.equal(journey.back().kind, "unavailable");
});

class FakeDocument {
  createElement(tag) { return new FakeElement(tag); }
}
class FakeElement {
  constructor(tag) {
    this.tagName = tag; this.children = []; this.textContent = "";
    this.handlers = Object.create(null);
  }
  append(...children) { this.children.push(...children); }
  addEventListener(event, handler) { this.handlers[event] = handler; }
  click() { this.handlers.click?.(); }
  replaceChildren(...children) { this.children = [...children]; this.textContent = ""; }
}
function allNodes(node) {
  return [node, ...node.children.flatMap(allNodes)];
}
function allText(node) {
  return allNodes(node).map(x => x.textContent);
}
