import assert from "node:assert/strict";
import test from "node:test";

import { createEventSearchController } from "./integration.mjs";
import { createResultRenderer, createResultWindow } from "./results.mjs";
import { createShellModel, renderEventSearchShell } from "./shell.mjs";
import { parseDeepLink, serializeDeepLinkState } from "./state.mjs";

const EXPECTED = Object.freeze({ tenantId: "tenant-a", environmentId: "prod-eu" });

function deepLink(overrides = {}) {
  return {
    tenantId: "tenant-a",
    environmentId: "prod-eu",
    from: "2026-09-24T10:00:00Z",
    to: "2026-09-24T11:00:00Z",
    query: 'message:"<script>not executable</script>"',
    sources: ["edr", "firewall"],
    filters: { severity: ["high"] },
    ...overrides,
  };
}

function execution(correlationId = "corr-e2e-001") {
  return { correlationId, queryVersion: "query-v1" };
}

function partialJob(request, payload = '<img src=x onerror="boom">') {
  return {
    kind: "job",
    jobId: "job-e2e-001",
    runId: "run-e2e-001",
    state: "partial",
    tenantRef: request.tenantRef,
    environmentRef: request.environmentRef,
    timeStart: request.timeStart,
    timeEnd: request.timeEnd,
    sources: [...request.sources],
    correlationId: request.correlationId,
    queryVersion: request.queryVersion,
    resultRefs: [{ source: "edr", ref: "evt-e2e-001" }],
    failedSources: [{ source: "firewall", code: "TIMEOUT" }],
    projections: [{
      ref: "evt-e2e-001",
      source: "edr",
      eventTime: "2026-09-24T10:30:00Z",
      summary: { message: payload, severity: "high" },
    }],
    freshness: [{ source: "edr", observedAt: "2026-09-24T10:59:00Z" }],
    consequence: "Firewall data is incomplete; valid EDR results are preserved.",
  };
}

function permissionDenied(request) {
  return {
    kind: "permission-denied",
    correlationId: request.correlationId,
    capability: "perm.investigate.search.execute",
    reasonCode: "POLICY_DENY",
    requestAccessPath: "/access/request",
  };
}

function createHarness(state) {
  const document = new FakeDocument();
  const shellMount = new FakeElement("div");
  const resultMount = new FakeElement("div");
  const renderer = createResultRenderer(document, resultMount);
  const views = [];

  function render(view) {
    views.push(view);
    if (view.shellState !== null) {
      renderEventSearchShell(
        document,
        shellMount,
        createShellModel(state, view.shellState, EXPECTED),
      );
    }
    renderer.render(view.resultWindow ?? createResultWindow([], { size: 50 }));
  }

  return { document, shellMount, resultMount, renderer, views, render };
}

test("e2e-partial-refresh-accessibility-and-hostile-content-remain-safe", async () => {
  const original = deepLink();
  const url = serializeDeepLinkState(original, EXPECTED);
  const refreshed = parseDeepLink(url, EXPECTED);
  assert.equal(serializeDeepLinkState(refreshed, EXPECTED), url);

  const harness = createHarness(refreshed);
  const controller = createEventSearchController({
    transport: {
      async execute(request) {
        return partialJob(request);
      },
      async retry() {
        throw new Error("not used");
      },
    },
    onView: harness.render,
    expectedContext: EXPECTED,
    windowSize: 50,
  });

  const outcome = await controller.execute(refreshed, execution());
  assert.equal(outcome.job.state, "partial");
  assert.deepEqual(outcome.job.failedSources.map((item) => item.source), ["firewall"]);
  assert.deepEqual(harness.renderer.visibleRefs(), ["evt-e2e-001"]);

  const hostile = '<img src=x onerror="boom">';
  assert.equal(flattenText(harness.resultMount).includes(hostile), true);
  assert.equal(findByTag(harness.resultMount, "img").length, 0);
  assert.equal(findByTag(harness.shellMount, "main").length, 1);
  assert.equal(findByTag(harness.shellMount, "form").length, 1);
  assert.equal(findByTag(harness.shellMount, "label").length, 1);
  assert.equal(findByTag(harness.shellMount, "textarea").length, 1);
  assert.equal(findByAttribute(harness.shellMount, "aria-live").length > 0, true);
  assert.equal(flattenText(harness.shellMount).includes("firewall"), true);
  assert.equal(flattenText(harness.shellMount).includes("valid EDR results are preserved"), true);
});

test("e2e-permission-denial-clears-previous-results-and-masks-query", async () => {
  const state = parseDeepLink(serializeDeepLinkState(deepLink()), EXPECTED);
  const harness = createHarness(state);
  let calls = 0;
  const controller = createEventSearchController({
    transport: {
      async execute(request) {
        calls += 1;
        return calls === 1 ? partialJob(request, "AUTHORIZED-RESULT") : permissionDenied(request);
      },
      async retry() {
        throw new Error("not used");
      },
    },
    onView: harness.render,
    expectedContext: EXPECTED,
  });

  await controller.execute(state, execution("corr-e2e-allow"));
  assert.equal(flattenText(harness.resultMount).includes("AUTHORIZED-RESULT"), true);

  const denied = await controller.execute(state, execution("corr-e2e-deny"));
  assert.equal(denied.kind, "permission-denied");
  assert.equal(controller.getView().resultWindow, undefined);
  assert.deepEqual(harness.renderer.visibleRefs(), []);
  assert.equal(flattenText(harness.resultMount).includes("AUTHORIZED-RESULT"), false);
  assert.equal(flattenText(harness.shellMount).includes(state.query), false);
  assert.equal(findByTag(harness.shellMount, "textarea")[0].value, "");
});

test("e2e-cancellation-prevents-stale-result-publication", async () => {
  const state = parseDeepLink(serializeDeepLinkState(deepLink()), EXPECTED);
  const harness = createHarness(state);
  let aborted = false;
  const controller = createEventSearchController({
    transport: {
      execute(request, { signal }) {
        return new Promise((resolve, reject) => {
          signal.addEventListener("abort", () => {
            aborted = true;
            reject(new DOMException("cancelled", "AbortError"));
          }, { once: true });
          setTimeout(() => resolve(partialJob(request, "STALE-RESULT")), 50);
        });
      },
      async retry() {
        throw new Error("not used");
      },
    },
    onView: harness.render,
    expectedContext: EXPECTED,
  });

  const pending = controller.execute(state, execution("corr-e2e-cancel"));
  assert.equal(controller.cancel(), true);
  const result = await pending;
  assert.equal(aborted, true);
  assert.equal(result.kind, "cancelled");
  assert.equal(flattenText(harness.resultMount).includes("STALE-RESULT"), false);
});

test("e2e-offline-state-remains-semantic-and-non-executable", () => {
  const state = parseDeepLink(serializeDeepLinkState(deepLink()), EXPECTED);
  const document = new FakeDocument();
  const mount = new FakeElement("div");
  const model = createShellModel(state, {
    kind: "offline",
    lastSync: "2026-09-24T10:58:00Z",
  }, EXPECTED);
  renderEventSearchShell(document, mount, model);
  assert.equal(model.canExecute, false);
  assert.equal(findByTag(mount, "textarea")[0].disabled, true);
  assert.equal(findByTag(mount, "button").find((node) => node.attributes.type === "submit").disabled, true);
  assert.equal(findByTag(mount, "main").length, 1);
  assert.equal(findByAttribute(mount, "aria-live").length > 0, true);
});

class FakeDocument {
  createElement(tag) {
    return new FakeElement(tag);
  }
}

class FakeElement {
  constructor(tag) {
    this.tagName = String(tag).toLowerCase();
    this.attributes = Object.create(null);
    this.children = [];
    this.parent = undefined;
    this.textContent = "";
    this.value = "";
    this.disabled = false;
    this.focused = false;
    this.scrollTop = 0;
    this.listeners = new Map();
  }

  set innerHTML(_) {
    throw new Error("unsafe innerHTML sink used");
  }

  setAttribute(name, value) {
    this.attributes[name] = String(value);
  }

  append(...nodes) {
    for (const node of nodes) {
      if (node.parent !== undefined) {
        const index = node.parent.children.indexOf(node);
        if (index >= 0) node.parent.children.splice(index, 1);
      }
      node.parent = this;
      this.children.push(node);
    }
  }

  replaceChildren(...nodes) {
    for (const child of this.children) child.parent = undefined;
    this.children = [];
    this.append(...nodes);
  }

  remove() {
    if (this.parent === undefined) return;
    const index = this.parent.children.indexOf(this);
    if (index >= 0) this.parent.children.splice(index, 1);
    this.parent = undefined;
  }

  addEventListener(type, handler) {
    const handlers = this.listeners.get(type) ?? [];
    handlers.push(handler);
    this.listeners.set(type, handlers);
  }

  focus() {
    this.focused = true;
  }
}

function walk(root, predicate, out = []) {
  if (predicate(root)) out.push(root);
  for (const child of root.children ?? []) walk(child, predicate, out);
  return out;
}

function findByTag(root, tag) {
  return walk(root, (node) => node.tagName === tag);
}

function findByAttribute(root, name) {
  return walk(root, (node) => Object.hasOwn(node.attributes ?? {}, name));
}

function flattenText(root) {
  return [root.textContent, ...(root.children ?? []).map(flattenText)].join(" ");
}
