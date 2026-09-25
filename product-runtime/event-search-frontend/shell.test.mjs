import assert from "node:assert/strict";
import test from "node:test";

import {
  createShellModel,
  normalizeShellState,
  renderEventSearchShell,
  ShellValidationError,
} from "./shell.mjs";

function validDeepLink() {
  return {
    tenantId: "tenant-a",
    environmentId: "prod-eu",
    from: "2026-09-24T10:00:00Z",
    to: "2026-09-24T11:00:00Z",
    query: "sensitive-user-query",
    sources: ["edr", "firewall"],
    filters: { severity: ["high"] },
    returnTo: "/investigate/case/42",
  };
}

function stateFixtures() {
  return [
    { kind: "loading" },
    { kind: "empty", reason: "No matching telemetry", adjustmentHint: "Widen the time range" },
    {
      kind: "partial",
      failedSources: ["firewall"],
      freshness: [{ source: "edr", observedAt: "2026-09-24T10:59:00Z" }],
      consequence: "Decisions may exclude firewall telemetry.",
    },
    { kind: "error", message: "Search failed", correlationId: "corr-42", retryable: true },
    { kind: "offline", lastSync: "2026-09-24T10:58:00Z" },
    {
      kind: "permission-denied",
      capability: "perm.investigate.search.execute",
      reasonCode: "POLICY_DENY",
      requestAccessPath: "/access/request?capability=event-search",
    },
  ];
}

test("six-mandatory-states-normalize-deterministically", () => {
  for (const fixture of stateFixtures()) {
    const first = normalizeShellState(fixture);
    const second = normalizeShellState(structuredClone(fixture));
    assert.deepEqual(second, first);
    assert.ok(Object.isFrozen(first));
  }
});

test("unknown-shell-state-and-fields-fail-closed", () => {
  assert.throws(
    () => normalizeShellState({ kind: "ready" }),
    (error) => error instanceof ShellValidationError && error.code === "unknown-shell-state",
  );
  assert.throws(
    () => normalizeShellState({ kind: "loading", protectedContent: "secret" }),
    (error) => error instanceof ShellValidationError && error.code === "unknown-shell-field",
  );
});

test("partial-state-requires-failed-source-freshness-and-consequence", () => {
  assert.throws(() => normalizeShellState({ kind: "partial", failedSources: [] }));
  assert.throws(() =>
    normalizeShellState({
      kind: "partial",
      failedSources: ["edr"],
      freshness: [],
      consequence: "Incomplete",
    }),
  );
  const state = normalizeShellState({
    kind: "partial",
    failedSources: ["firewall", "firewall"],
    freshness: [
      { source: "firewall", observedAt: "2026-09-24T10:30:00Z" },
      { source: "edr", observedAt: "2026-09-24T10:45:00Z" },
    ],
    consequence: "Incomplete source set",
  });
  assert.deepEqual(state.failedSources, ["firewall"]);
  assert.deepEqual(state.freshness.map((item) => item.source), ["edr", "firewall"]);
});

test("unsafe-access-request-path-is-rejected", () => {
  assert.throws(() =>
    normalizeShellState({
      kind: "permission-denied",
      capability: "perm.investigate.search.execute",
      reasonCode: "DENY",
      requestAccessPath: "https://evil.example/request",
    }),
  );
});

test("shell-model-is-immutable-and-offline-or-denied-cannot-execute", () => {
  const offline = createShellModel(validDeepLink(), {
    kind: "offline",
    lastSync: "2026-09-24T10:00:00Z",
  });
  const denied = createShellModel(validDeepLink(), {
    kind: "permission-denied",
    capability: "perm.investigate.search.execute",
    reasonCode: "DENY",
  });
  assert.equal(offline.canExecute, false);
  assert.equal(denied.canExecute, false);
  assert.ok(Object.isFrozen(offline));
  assert.ok(Object.isFrozen(offline.ui));
});

test("all-six-states-render-with-semantic-status-and-focus-target", () => {
  for (const fixture of stateFixtures()) {
    const document = new FakeDocument();
    const mount = new FakeElement("div");
    const model = createShellModel(validDeepLink(), fixture);
    const rendered = renderEventSearchShell(document, mount, model);
    assert.equal(rendered.focusTarget.attributes.id, "event-search-state-heading");
    assert.equal(rendered.focusTarget.attributes.tabindex, "-1");
    assert.ok(["status", "alert"].includes(findByAttribute(mount, "aria-live")[0].attributes.role));
    assert.equal(findByTag(mount, "main").length, 1);
    assert.equal(findByTag(mount, "form").length, 1);
    assert.equal(findByTag(mount, "label").length, 1);
    assert.equal(findByTag(mount, "textarea").length, 1);
    assert.equal(findByTag(mount, "button").some((node) => node.attributes.type === "submit"), true);
  }
});

test("native-form-submit-emits-only-normalized-state", () => {
  const document = new FakeDocument();
  const mount = new FakeElement("div");
  const model = createShellModel(validDeepLink(), { kind: "loading" });
  let received;
  renderEventSearchShell(document, mount, model, {
    onExecute(value) {
      received = value;
    },
  });
  const form = findByTag(mount, "form")[0];
  const event = form.dispatch("submit");
  assert.equal(event.defaultPrevented, true);
  assert.equal(received, model.deepLink);
});

test("offline-state-disables-native-query-and-submit-controls", () => {
  const document = new FakeDocument();
  const mount = new FakeElement("div");
  const model = createShellModel(validDeepLink(), {
    kind: "offline",
    lastSync: "2026-09-24T10:00:00Z",
  });
  renderEventSearchShell(document, mount, model);
  assert.equal(findByTag(mount, "textarea")[0].disabled, true);
  assert.equal(findByTag(mount, "button").find((node) => node.attributes.type === "submit").disabled, true);
});

test("permission-denied-reveals-no-protected-content", () => {
  const document = new FakeDocument();
  const mount = new FakeElement("div");
  const model = createShellModel(
    { ...validDeepLink(), query: "TOP-SECRET-QUERY" },
    {
      kind: "permission-denied",
      capability: "perm.investigate.search.execute",
      reasonCode: "POLICY_DENY",
      requestAccessPath: "/access/request",
    },
  );
  renderEventSearchShell(document, mount, model);
  const textarea = findByTag(mount, "textarea")[0];
  assert.equal(textarea.value, "");
  assert.equal(flattenText(mount).includes("TOP-SECRET-QUERY"), false);
  assert.equal(flattenText(mount).includes("POLICY_DENY"), true);
  assert.equal(findByTag(mount, "a").some((node) => node.attributes.href === "/access/request"), true);
});

test("partial-error-empty-and-offline-details-are-visible-and-safe", () => {
  for (const fixture of stateFixtures().filter((item) =>
    ["partial", "error", "empty", "offline"].includes(item.kind))) {
    const document = new FakeDocument();
    const mount = new FakeElement("div");
    renderEventSearchShell(document, mount, createShellModel(validDeepLink(), fixture));
    const text = flattenText(mount);
    if (fixture.kind === "partial") {
      assert.ok(text.includes("firewall"));
      assert.ok(text.includes("Decisions may exclude firewall telemetry."));
    }
    if (fixture.kind === "error") {
      assert.ok(text.includes("corr-42"));
      assert.ok(text.includes("Retry search"));
    }
    if (fixture.kind === "empty") {
      assert.ok(text.includes("Widen the time range"));
      assert.ok(text.includes("Adjust filters"));
    }
    if (fixture.kind === "offline") {
      assert.ok(text.includes("Last synchronization"));
    }
  }
});

test("focus-target-is-programmatically-focusable-without-forcing-focus", () => {
  const document = new FakeDocument();
  const mount = new FakeElement("div");
  const rendered = renderEventSearchShell(
    document,
    mount,
    createShellModel(validDeepLink(), { kind: "loading" }),
  );
  assert.equal(rendered.focusTarget.focused, false);
  rendered.focusTarget.focus();
  assert.equal(rendered.focusTarget.focused, true);
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
    this.textContent = "";
    this.value = "";
    this.disabled = false;
    this.focused = false;
    this.listeners = new Map();
  }

  set innerHTML(_) {
    throw new Error("unsafe innerHTML sink used");
  }

  setAttribute(name, value) {
    this.attributes[name] = String(value);
  }

  append(...nodes) {
    this.children.push(...nodes);
  }

  replaceChildren(...nodes) {
    this.children = [...nodes];
  }

  addEventListener(type, handler) {
    const handlers = this.listeners.get(type) ?? [];
    handlers.push(handler);
    this.listeners.set(type, handlers);
  }

  dispatch(type) {
    const event = {
      defaultPrevented: false,
      preventDefault() {
        this.defaultPrevented = true;
      },
    };
    for (const handler of this.listeners.get(type) ?? []) {
      handler(event);
    }
    return event;
  }

  focus() {
    this.focused = true;
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

function findByAttribute(root, name) {
  return walk(root, (node) => Object.hasOwn(node.attributes ?? {}, name));
}

function flattenText(root) {
  return [root.textContent, ...(root.children ?? []).map(flattenText)].join(" ");
}
