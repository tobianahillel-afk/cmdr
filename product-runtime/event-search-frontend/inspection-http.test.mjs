import assert from "node:assert/strict";
import test from "node:test";
import { spawn, execFile } from "node:child_process";
import { randomBytes } from "node:crypto";
import { mkdtemp, rm } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { promisify } from "node:util";
import { createEventSearchController } from "./integration.mjs";
import { createEventInspectionJourney, normalizePivotDraft } from "./inspection.mjs";
import { createInspectionHTTPTransport } from "./inspection-http.mjs";

const exec = promisify(execFile);
const fixtureDir = resolve(fileURLToPath(new URL("../event-inspection-host/", import.meta.url)));
const context = {
  tenantId: "tenant-a", environmentId: "prod",
  from: "2026-10-20T10:00:00Z", to: "2026-10-20T12:00:00Z",
  sources: ["edr"], query: "source:edr",
};

function searchController() {
  return createEventSearchController({
    expectedContext: { tenantId: "tenant-a", environmentId: "prod" },
    transport: {
      async execute(request) {
        return {
          kind: "job", jobId: "job-1", runId: "run-1", state: "completed",
          tenantRef: request.tenantRef, environmentRef: request.environmentRef,
          timeStart: request.timeStart, timeEnd: request.timeEnd,
          sources: [...request.sources], correlationId: request.correlationId,
          queryVersion: request.queryVersion,
          resultRefs: [{ source: "edr", ref: "event-1" }],
          failedSources: [],
          projections: [{ ref: "event-1", eventTime: "2026-10-20T11:00:00Z",
            source: "edr", summary: { severity: "high" } }],
        };
      },
      async retry() { throw new Error("unexpected retry"); },
    },
  });
}

async function loopbackHost(t) {
  const directory = await mkdtemp(join(tmpdir(), "cmdr-inspection-"));
  const binary = join(directory, process.platform === "win32" ? "fixture.exe" : "fixture");
  t.after(async () => { await rm(directory, { recursive: true, force: true }); });
  await exec("go", ["build", "-o", binary, "./testdata/fixture"], { cwd: fixtureDir, timeout: 120000 });
  const token = randomBytes(32).toString("hex");
  const child = spawn(binary, [], {
    cwd: fixtureDir, env: { ...process.env, CMDR_TEST_BEARER: token },
    stdio: ["ignore", "pipe", "pipe"],
  });
  t.after(() => { child.kill(); });
  const endpoint = await new Promise((accept, reject) => {
    let stdout = "";
    let stderr = "";
    const timeout = setTimeout(() => {
      child.kill();
      reject(new Error("loopback fixture timeout"));
    }, 10000);
    child.stderr.on("data", (chunk) => { stderr += String(chunk).slice(0, 4096); });
    child.once("error", (error) => { clearTimeout(timeout); reject(error); });
    child.once("exit", (code) => {
      clearTimeout(timeout);
      reject(new Error(`loopback fixture exited ${code}: ${stderr}`));
    });
    child.stdout.on("data", (chunk) => {
      stdout += String(chunk);
      const match = stdout.match(/http:\/\/127\.0\.0\.1:\d+/u);
      if (match) {
        clearTimeout(timeout);
        accept(match[0]);
      }
    });
  });
  return { endpoint, token };
}

function realLocalFetch(endpoint, token) {
  return (route, options) => {
    assert.ok(["/investigate/event-inspection", "/investigate/event-inspection/pivot"].includes(route));
    assert.equal(options.credentials, "same-origin");
    assert.equal(options.redirect, "error");
    return fetch(new URL(route, endpoint), {
      ...options,
      headers: { ...options.headers, Authorization: `Bearer ${token}` },
    });
  };
}

test("real Node fetch to Go authorized HTTP host opens and returns to identical Search window", async (t) => {
  const { endpoint, token } = await loopbackHost(t);
  const controller = searchController();
  await controller.execute(context, { correlationId: "corr-1", queryVersion: "query-1" });
  const windowBefore = controller.getView().resultWindow;
  const transport = createInspectionHTTPTransport({ fetchImpl: realLocalFetch(endpoint, token) });
  const journey = createEventInspectionJourney({ searchController: controller, transport });
  const inspected = await journey.open(context, "event-1");
  assert.equal(inspected.kind, "inspected");
  assert.equal(inspected.inspection.eventRef, "event-1");
  assert.equal(inspected.inspection.source.sourceRef, "edr");
  assert.equal(inspected.inspection.parser.parserRef, "parser-a");
  assert.equal(inspected.inspection.sourceFields["source.ip"], "192.0.2.4");
  assert.equal(inspected.inspection.enrichments[0].freshness, "stale");
  assert.equal(inspected.inspection.raw.access, "denied");
  assert.equal(inspected.inspection.rendered.text, "rendered authorized");
  assert.deepEqual(Object.keys(inspected.inspection.rawDerivedSensitive), []);
  assert.ok(!JSON.stringify(inspected).includes("RAW-PRIVATE-STRING"));
  const probeRequest = {
    tenantRef: "tenant-a", environmentRef: "prod", eventRef: "event-1",
    correlationId: "corr-1", jobId: "job-1", runId: "run-1", sourceHint: "edr",
    field: "source.ip", value: "192.0.2.4", valuePresent: true,
    timeStart: controller.getView().job.timeStart, timeEnd: controller.getView().job.timeEnd,
  };
  const serverDraft = await transport.draftPivot(probeRequest);
  assert.equal(serverDraft.kind, "pivot-draft", JSON.stringify(serverDraft));
  assert.equal(normalizePivotDraft(serverDraft, probeRequest, inspected.inspection).fieldOrigin, "source");
  const pivot = await journey.draftPivot("source.ip", "192.0.2.4");
  assert.equal(pivot.kind, "inspected");
  assert.equal(pivot.pivotDraft.kind, "pivot-draft");
  assert.equal(pivot.pivotDraft.targetCapability, "CAP-INV-002");
  assert.equal(pivot.pivotDraft.fieldOrigin, "source");
  assert.equal(pivot.pivotDraft.jobId, "job-1");
  assert.equal(pivot.pivotDraft.runId, "run-1");
  assert.equal(pivot.pivotDraft.returnContext.origin, "event-search");
  assert.equal(pivot.pivotDraft.auditRequired, true);
  assert.equal(pivot.pivotDraft.timeRange.start, controller.getView().job.timeStart);
  assert.equal(pivot.pivotDraft.timeRange.end, controller.getView().job.timeEnd);
  assert.ok(!JSON.stringify(pivot.pivotDraft).includes("RAW-PRIVATE-STRING"));
  // A visible enrichment is not automatically allowed to become a pivot.
  const deniedPivot = await journey.draftPivot("risk.score", "72");
  assert.equal(deniedPivot.kind, "inspected");
  assert.equal(deniedPivot.pivotStatus, "unavailable");
  assert.equal(deniedPivot.pivotDraft, undefined);
  // Raw-derived fields are absent when independent raw permission is denied.
  await assert.rejects(() => journey.draftPivot("secret.token", "RAW-PRIVATE-STRING"), /pivot-not-available/);
  const prior = journey.back();
  assert.equal(prior.kind, "search");
  assert.equal(prior.resultWindow, windowBefore);

  const wrongBearer = createInspectionHTTPTransport({ fetchImpl: realLocalFetch(endpoint, "wrong-bearer") });
  await assert.rejects(() => wrongBearer.inspect({
    tenantRef: "tenant-a", environmentRef: "prod", eventRef: "event-1",
    jobId: "job-1", runId: "run-1", correlationId: "corr-1", sourceHint: "edr",
  }), /inspection-unavailable/);

  const deniedJourney = createEventInspectionJourney({
    searchController: controller,
    transport: createInspectionHTTPTransport({ fetchImpl: realLocalFetch(endpoint, token) }),
  });
  await assert.rejects(() => deniedJourney.open(context, "foreign-event"), /not-authorized-result/);
  assert.equal(deniedJourney.getView().kind, "search");

  const forgedScope = await fetch(new URL("/investigate/event-inspection", endpoint), {
    method: "POST", headers: { "Content-Type": "application/json", Authorization: `Bearer ${token}` },
    body: JSON.stringify({ tenantRef: "other-tenant", environmentRef: "prod",
      eventRef: "event-1", jobId: "job-1", runId: "run-1", correlationId: "corr-1", sourceHint: "edr" }),
  });
  assert.equal(forgedScope.status, 404);
  assert.ok(!(await forgedScope.text()).includes("RAW-PRIVATE-STRING"));
});

test("transport rejects non-JSON responses and oversized streams without rendering anything", async () => {
  const request = { tenantRef: "tenant-a" };
  const invalid = createInspectionHTTPTransport({
    fetchImpl: async () => new Response("not-json", { headers: { "Content-Type": "text/plain" } }),
  });
  await assert.rejects(() => invalid.inspect(request), /inspection-unavailable/);
  const large = createInspectionHTTPTransport({
    fetchImpl: async () => new Response("x".repeat(2 * 1024 * 1024 + 1),
      { headers: { "Content-Type": "application/json" } }),
  });
  await assert.rejects(() => large.inspect(request), /inspection-unavailable/);
});


test("browser refuses forged pivot draft provenance and keeps Inspector available", async (t) => {
  const { endpoint, token } = await loopbackHost(t);
  const controller = searchController();
  await controller.execute(context, { correlationId: "corr-1", queryVersion: "query-1" });
  const sourceFetch = realLocalFetch(endpoint, token);
  const transport = createInspectionHTTPTransport({
    fetchImpl: async (route, options) => {
      const response = await sourceFetch(route, options);
      if (route !== "/investigate/event-inspection/pivot" || response.status !== 200) return response;
      const claimed = await response.json();
      claimed.fieldOrigin = "normalized"; // cannot override server-visible source origin
      return new Response(JSON.stringify(claimed), {
        status: 200, headers: { "Content-Type": "application/json" },
      });
    },
  });
  const journey = createEventInspectionJourney({ searchController: controller, transport });
  assert.equal((await journey.open(context, "event-1")).kind, "inspected");
  const denied = await journey.draftPivot("source.ip", "192.0.2.4");
  assert.equal(denied.kind, "inspected");
  assert.equal(denied.pivotStatus, "unavailable");
  assert.equal(denied.pivotDraft, undefined);
  assert.equal(journey.back().kind, "search");
});

test("returning to Search revokes a pending pivot draft without resurrecting it", async (t) => {
  const { endpoint, token } = await loopbackHost(t);
  const controller = searchController();
  await controller.execute(context, { correlationId: "corr-1", queryVersion: "query-1" });
  const sourceFetch = realLocalFetch(endpoint, token);
  let release;
  const transport = createInspectionHTTPTransport({
    fetchImpl: (route, options) =>
      route.endsWith("/pivot")
        ? new Promise(resolve => { release = () => resolve(sourceFetch(route, options)); })
        : sourceFetch(route, options),
  });
  const journey = createEventInspectionJourney({ searchController: controller, transport });
  assert.equal((await journey.open(context, "event-1")).kind, "inspected");
  const pending = journey.draftPivot("source.ip", "192.0.2.4");
  assert.equal(journey.back().kind, "search");
  release();
  assert.equal((await pending).kind, "superseded");
  assert.equal(journey.getView().kind, "search");
});

test("same-origin transport never sends draft on a client-chosen URL", async () => {
  const calls = [];
  const transport = createInspectionHTTPTransport({
    fetchImpl: async (route, options) => {
      calls.push({ route, options });
      return new Response(JSON.stringify({ kind: "denied" }), {
        status: 404, headers: { "Content-Type": "application/json" },
      });
    },
  });
  await assert.rejects(() => transport.draftPivot({
    tenantRef: "tenant-a", environmentRef: "prod", field: "source.ip",
    value: "192.0.2.4", url: "https://evil.invalid/",
  }), /inspection-unavailable/);
  assert.equal(calls.length, 1);
  assert.equal(calls[0].route, "/investigate/event-inspection/pivot");
  assert.equal(calls[0].options.credentials, "same-origin");
  assert.equal(calls[0].options.redirect, "error");
  assert.equal(calls[0].options.cache, "no-store");
});
