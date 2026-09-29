import assert from "node:assert/strict";
import test from "node:test";

import {
  ResultValidationError,
  anchorResultWindow,
  createIncrementalScheduler,
  createResultRenderer,
  createResultWindow,
  normalizeResultCollection,
  normalizeResultProjection,
  reconcileResultWindow,
} from "./results.mjs";

function result(ref, second = 0, summary = {}) {
  return {
    ref,
    eventTime: `2026-09-24T10:00:${String(second).padStart(2, "0")}Z`,
    source: "edr",
    summary,
  };
}

test("result-projection-is-strict-deterministic-and-immutable", () => {
  const normalized = normalizeResultProjection(result("evt-1", 1, {
    user: "alice",
    score: 7,
    active: true,
    optional: null,
  }));
  assert.equal(normalized.eventTime, "2026-09-24T10:00:01.000Z");
  assert.deepEqual(Object.keys(normalized.summary), ["active", "optional", "score", "user"]);
  assert.equal(normalized.summary.score, "7");
  assert.equal(normalized.summary.active, "true");
  assert.equal(normalized.summary.optional, "null");
  assert.ok(Object.isFrozen(normalized));
  assert.ok(Object.isFrozen(normalized.summary));
});

test("result-projection-rejects-unknown-fields-unsafe-keys-and-complex-values", () => {
  assert.throws(
    () => normalizeResultProjection({ ...result("evt-1"), tenantSecret: "x" }),
    (error) => error instanceof ResultValidationError && error.code === "unknown-result-field",
  );
  assert.throws(
    () => normalizeResultProjection(result("evt-1", 0, { constructor: "x" })),
    (error) => error instanceof ResultValidationError && error.code === "unsafe-summary-key",
  );
  assert.throws(
    () => normalizeResultProjection(result("evt-1", 0, { nested: { value: "x" } })),
    (error) => error instanceof ResultValidationError && error.code === "invalid-summary-value",
  );
});

test("collection-rejects-duplicate-stable-references", () => {
  assert.throws(
    () => normalizeResultCollection([result("evt-1"), result("evt-1", 1)]),
    (error) => error instanceof ResultValidationError && error.code === "duplicate-result-ref",
  );
});

test("window-is-bounded-and-reports-before-after-counts", () => {
  const items = Array.from({ length: 8 }, (_, index) => result(`evt-${index}`, index));
  const windowState = createResultWindow(items, { start: 2, size: 3 });
  assert.deepEqual(windowState.rows.map((item) => item.ref), ["evt-2", "evt-3", "evt-4"]);
  assert.equal(windowState.beforeCount, 2);
  assert.equal(windowState.afterCount, 3);
  assert.equal(windowState.anchorRef, "evt-2");
  assert.throws(() => createResultWindow(items, { size: 201 }));
});

test("anchor-preserves-relative-position-when-earlier-results-arrive", () => {
  const first = createResultWindow(
    [result("a"), result("b", 1), result("c", 2), result("d", 3)],
    { start: 1, size: 2 },
  );
  const anchored = anchorResultWindow(first, 1);
  assert.equal(anchored.anchorRef, "c");
  const next = reconcileResultWindow(anchored, [
    result("x", 4),
    result("a"),
    result("b", 1),
    result("c", 2),
    result("d", 3),
  ]);
  assert.deepEqual(next.rows.map((item) => item.ref), ["b", "c"]);
  assert.equal(next.rows[next.anchorOffset].ref, "c");
});

test("missing-anchor-falls-back-to-current-start-without-inventing-position", () => {
  const first = anchorResultWindow(
    createResultWindow([result("a"), result("b", 1), result("c", 2)], { start: 1, size: 2 }),
    0,
  );
  const next = reconcileResultWindow(first, [result("a"), result("c", 2), result("d", 3)]);
  assert.equal(next.start, 1);
  assert.equal(next.anchorRef, "c");
});

test("renderer-uses-text-safe-native-table-primitives", () => {
  const document = new FakeDocument();
  const mount = new FakeElement("div");
  const renderer = createResultRenderer(document, mount);
  const payload = '<img src=x onerror="boom">';
  const windowState = createResultWindow([
    result("evt-1", 1, { message: payload, severity: "high" }),
  ], { size: 10 });
  const summary = renderer.render(windowState);
  assert.equal(summary.rendered, 1);
  assert.equal(findByTag(mount, "table").length, 1);
  assert.equal(findByTag(mount, "th").length, 4);
  assert.equal(flattenText(mount).includes(payload), true);
  assert.equal(findByTag(mount, "img").length, 0);
  assert.equal(renderer.getRowNode("evt-1").attributes["data-result-ref"], "evt-1");
});

test("incremental-render-reuses-stable-row-and-preserves-scroll-position", () => {
  const document = new FakeDocument();
  const mount = new FakeElement("div");
  mount.scrollTop = 420;
  const renderer = createResultRenderer(document, mount);
  renderer.render(createResultWindow([result("a"), result("b", 1)], { size: 2 }));
  const originalB = renderer.getRowNode("b");
  renderer.render(createResultWindow([
    result("b", 1, { updated: true }),
    result("c", 2),
  ], { size: 2 }));
  assert.equal(renderer.getRowNode("b"), originalB);
  assert.equal(renderer.getRowNode("a"), undefined);
  assert.deepEqual(renderer.visibleRefs(), ["b", "c"]);
  assert.equal(mount.scrollTop, 420);
  assert.ok(flattenText(originalB).includes("updated=true"));
});

test("renderer-handles-empty-window-without-full-fragment-replacement", () => {
  const document = new FakeDocument();
  const mount = new FakeElement("div");
  const renderer = createResultRenderer(document, mount);
  renderer.render(createResultWindow([result("a")], { size: 1 }));
  const table = findByTag(mount, "table")[0];
  const section = findByTag(mount, "section")[0];
  const summary = renderer.render(createResultWindow([], { size: 1 }));
  assert.equal(summary.rendered, 0);
  assert.equal(findByTag(mount, "table")[0], table);
  assert.equal(findByTag(mount, "section")[0], section);
  assert.ok(flattenText(mount).includes("No result projections are available."));
});

test("scheduler-coalesces-bursts-to-latest-update", () => {
  const queued = [];
  const scheduler = createIncrementalScheduler((task) => queued.push(task));
  const calls = [];
  scheduler.submit(() => calls.push("first"));
  scheduler.submit(() => calls.push("latest"));
  assert.equal(scheduler.isPending(), true);
  assert.equal(queued.length, 1);
  queued.shift()();
  assert.deepEqual(calls, ["latest"]);
  assert.equal(scheduler.isPending(), false);
});

test("scheduler-cancel-drops-pending-work", () => {
  const queued = [];
  const scheduler = createIncrementalScheduler((task) => queued.push(task));
  let called = false;
  scheduler.submit(() => { called = true; });
  scheduler.cancel();
  queued.shift()();
  assert.equal(called, false);
  assert.throws(() => scheduler.submit("not-a-function"));
});

test("renderer-and-window-fail-closed-on-invalid-input", () => {
  assert.throws(() => createResultRenderer({}, new FakeElement("div")));
  assert.throws(() => createResultRenderer(new FakeDocument(), {}));
  assert.throws(() => createResultWindow("not-an-array"));
  assert.throws(() => anchorResultWindow({}, 0));
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
    this.scrollTop = 0;
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
        if (index >= 0) {
          node.parent.children.splice(index, 1);
        }
      }
      node.parent = this;
      this.children.push(node);
    }
  }

  remove() {
    if (this.parent === undefined) {
      return;
    }
    const index = this.parent.children.indexOf(this);
    if (index >= 0) {
      this.parent.children.splice(index, 1);
    }
    this.parent = undefined;
  }
}

function walk(root, predicate, out = []) {
  if (predicate(root)) {
    out.push(root);
  }
  for (const child of root.children ?? []) {
    walk(child, predicate, out);
  }
  return out;
}

function findByTag(root, tag) {
  return walk(root, (node) => node.tagName === tag);
}

function flattenText(root) {
  return [root.textContent, ...(root.children ?? []).map(flattenText)].join(" ");
}
