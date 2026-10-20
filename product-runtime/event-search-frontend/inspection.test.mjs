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
  controller.getView().resultWindow.rows[0].ref === "event-1";
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

class FakeDocument {
  createElement(tag) { return new FakeElement(tag); }
}
class FakeElement {
  constructor(tag) {
    this.tagName = tag; this.children = []; this.textContent = "";
  }
  append(...children) { this.children.push(...children); }
  replaceChildren(...children) { this.children = [...children]; this.textContent = ""; }
}
function allNodes(node) {
  return [node, ...node.children.flatMap(allNodes)];
}
function allText(node) {
  return allNodes(node).map(x => x.textContent);
}
