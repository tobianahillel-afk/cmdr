import { normalizeDeepLinkState } from "./state.mjs";

const SHELL_STATES = new Set([
  "loading",
  "empty",
  "partial",
  "error",
  "offline",
  "permission-denied",
]);
const SAFE_PATH_ORIGIN = "https://cmdr.invalid";

export class ShellValidationError extends Error {
  constructor(code, message) {
    super(message);
    this.name = "ShellValidationError";
    this.code = code;
  }
}

export function createShellModel(deepLinkInput, uiInput, expected = {}) {
  const deepLink = normalizeDeepLinkState(deepLinkInput, expected);
  const ui = normalizeShellState(uiInput);
  return freezeDeep({
    deepLink,
    ui,
    focusTargetId: "event-search-state-heading",
    canExecute: ui.kind !== "offline" && ui.kind !== "permission-denied",
  });
}

export function normalizeShellState(input) {
  if (!isPlainObject(input)) {
    throw invalid("invalid-shell-state", "shell state must be a plain object");
  }
  const kind = requiredText("shell state kind", input.kind);
  if (!SHELL_STATES.has(kind)) {
    throw invalid("unknown-shell-state", `unknown shell state: ${kind}`);
  }

  switch (kind) {
    case "loading":
      rejectUnknown(input, new Set(["kind"]));
      return freezeDeep({ kind });
    case "empty":
      rejectUnknown(input, new Set(["kind", "reason", "adjustmentHint"]));
      return freezeDeep({
        kind,
        reason: optionalText(input.reason) ?? "No events match the current search scope.",
        adjustmentHint: optionalText(input.adjustmentHint) ?? "Adjust the time range, sources, or filters.",
      });
    case "partial":
      rejectUnknown(input, new Set(["kind", "failedSources", "freshness", "consequence"]));
      return freezeDeep({
        kind,
        failedSources: normalizeList("failed source", input.failedSources),
        freshness: normalizeFreshness(input.freshness),
        consequence: requiredText("partial-result consequence", input.consequence),
      });
    case "error":
      rejectUnknown(input, new Set(["kind", "message", "correlationId", "retryable"]));
      return freezeDeep({
        kind,
        message: requiredText("error message", input.message),
        correlationId: requiredText("correlation id", input.correlationId),
        retryable: requireBoolean("retryable", input.retryable),
      });
    case "offline":
      rejectUnknown(input, new Set(["kind", "lastSync", "message"]));
      return freezeDeep({
        kind,
        lastSync: normalizeInstant("last synchronization", input.lastSync),
        message: optionalText(input.message) ?? "Offline mode is read-only. New searches and mutations are unavailable.",
      });
    case "permission-denied":
      rejectUnknown(input, new Set(["kind", "capability", "reasonCode", "requestAccessPath"]));
      return freezeDeep({
        kind,
        capability: requiredText("denied capability", input.capability),
        reasonCode: requiredText("deny reason code", input.reasonCode),
        requestAccessPath: normalizeInternalPath(input.requestAccessPath),
      });
  }
  throw invalid("unknown-shell-state", "unsupported shell state");
}

export function renderEventSearchShell(document, mount, model, handlers = {}) {
  if (!document || typeof document.createElement !== "function") {
    throw invalid("invalid-document", "a DOM-compatible document is required");
  }
  if (!mount || typeof mount.replaceChildren !== "function") {
    throw invalid("invalid-mount", "a DOM-compatible mount point is required");
  }
  if (!model || !model.deepLink || !model.ui) {
    throw invalid("invalid-model", "a normalized shell model is required");
  }

  const root = element(document, "div", { class: "cmdr-event-search" });
  const skip = element(document, "a", { href: "#event-search-main" }, "Skip to search content");
  root.append(skip);

  const header = element(document, "header");
  header.append(element(document, "p", { "aria-label": "Product" }, "Investigate"));
  header.append(element(document, "h1", {}, "Event Search"));
  root.append(header);

  const nav = element(document, "nav", { "aria-label": "Investigate navigation" });
  nav.append(element(document, "a", { href: "/investigate" }, "Investigate home"));
  root.append(nav);

  const context = element(document, "section", { "aria-labelledby": "event-search-context-heading" });
  context.append(element(document, "h2", { id: "event-search-context-heading" }, "Search context"));
  const contextList = element(document, "dl");
  appendDefinition(document, contextList, "Tenant", model.deepLink.tenantId);
  appendDefinition(document, contextList, "Environment", model.deepLink.environmentId);
  appendDefinition(document, contextList, "From", model.deepLink.from);
  appendDefinition(document, contextList, "To", model.deepLink.to);
  appendDefinition(document, contextList, "Sources", model.deepLink.sources.join(", "));
  context.append(contextList);
  root.append(context);

  const main = element(document, "main", { id: "event-search-main" });
  const form = element(document, "form", { "aria-label": "Event Search query" });
  const label = element(document, "label", { for: "event-search-query" }, "Query");
  const query = element(document, "textarea", {
    id: "event-search-query",
    name: "query",
    rows: "6",
    "aria-describedby": "event-search-query-help",
  });
  query.value = model.ui.kind === "permission-denied" ? "" : model.deepLink.query;
  query.disabled = !model.canExecute;
  form.append(label, query);
  form.append(element(
    document,
    "p",
    { id: "event-search-query-help" },
    "Query syntax remains backend-neutral in this implementation wave.",
  ));

  const execute = element(document, "button", { type: "submit" }, "Run search");
  execute.disabled = !model.canExecute;
  form.append(execute);
  form.addEventListener("submit", (event) => {
    event.preventDefault();
    if (!model.canExecute) {
      return;
    }
    if (typeof handlers.onExecute === "function") {
      handlers.onExecute(model.deepLink);
    }
  });
  main.append(form);

  if (model.deepLink.returnTo !== undefined) {
    main.append(element(document, "a", { href: model.deepLink.returnTo }, "Return to previous context"));
  }

  const status = renderStatus(document, model.ui, handlers);
  main.append(status.region);
  root.append(main);
  mount.replaceChildren(root);
  return { root, focusTarget: status.focusTarget };
}

function renderStatus(document, ui, handlers) {
  const role = ui.kind === "error" || ui.kind === "permission-denied" ? "alert" : "status";
  const live = role === "alert" ? "assertive" : "polite";
  const region = element(document, "section", {
    "aria-labelledby": "event-search-state-heading",
    role,
    "aria-live": live,
  });
  const heading = element(document, "h2", {
    id: "event-search-state-heading",
    tabindex: "-1",
  }, stateTitle(ui.kind));
  region.append(heading);

  switch (ui.kind) {
    case "loading":
      region.setAttribute("aria-busy", "true");
      region.append(element(document, "p", {}, "Loading Event Search. Existing query and scope are preserved."));
      break;
    case "empty":
      region.append(element(document, "p", {}, ui.reason));
      region.append(element(document, "p", {}, ui.adjustmentHint));
      region.append(actionButton(document, "Adjust filters", handlers.onAdjustFilters));
      break;
    case "partial":
      region.append(element(document, "p", {}, ui.consequence));
      region.append(element(document, "p", {}, "Some sources did not return complete data."));
      region.append(textList(document, "Missing sources", ui.failedSources));
      region.append(freshnessList(document, ui.freshness));
      break;
    case "error":
      region.append(element(document, "p", {}, ui.message));
      region.append(element(document, "p", {}, "Correlation ID:"));
      region.append(element(document, "code", {}, ui.correlationId));
      if (ui.retryable) {
        region.append(actionButton(document, "Retry search", handlers.onRetry));
      }
      break;
    case "offline":
      region.append(element(document, "p", {}, ui.message));
      const lastSync = element(document, "time", { datetime: ui.lastSync }, ui.lastSync);
      region.append(element(document, "p", {}, "Last synchronization:"));
      region.append(lastSync);
      break;
    case "permission-denied":
      region.append(element(
        document,
        "p",
        {},
        "You do not have permission to execute Event Search in this context.",
      ));
      appendDefinition(document, region, "Capability", ui.capability);
      appendDefinition(document, region, "Reason code", ui.reasonCode);
      if (ui.requestAccessPath !== undefined) {
        region.append(element(document, "a", { href: ui.requestAccessPath }, "Request access"));
      }
      break;
  }
  return { region, focusTarget: heading };
}

function stateTitle(kind) {
  return {
    loading: "Loading",
    empty: "No results",
    partial: "Partial results",
    error: "Search error",
    offline: "Offline",
    "permission-denied": "Permission denied",
  }[kind];
}

function actionButton(document, label, handler) {
  const button = element(document, "button", { type: "button" }, label);
  if (typeof handler === "function") {
    button.addEventListener("click", handler);
  }
  return button;
}

function textList(document, label, values) {
  const section = element(document, "section", { "aria-label": label });
  const list = element(document, "ul");
  for (const value of values) {
    list.append(element(document, "li", {}, value));
  }
  section.append(list);
  return section;
}

function freshnessList(document, freshness) {
  const section = element(document, "section", { "aria-label": "Source freshness" });
  const list = element(document, "dl");
  for (const item of freshness) {
    list.append(element(document, "dt", {}, item.source));
    list.append(element(document, "dd", {}, item.observedAt));
  }
  section.append(list);
  return section;
}

function appendDefinition(document, parent, label, value) {
  parent.append(element(document, "dt", {}, label));
  parent.append(element(document, "dd", {}, value));
}

function element(document, tag, attributes = {}, text) {
  const node = document.createElement(tag);
  for (const [name, value] of Object.entries(attributes)) {
    node.setAttribute(name, String(value));
  }
  if (text !== undefined) {
    node.textContent = String(text);
  }
  return node;
}

function normalizeFreshness(value) {
  if (!Array.isArray(value) || value.length === 0) {
    throw invalid("invalid-freshness", "partial results require source freshness");
  }
  const seen = new Set();
  const result = [];
  for (const item of value) {
    if (!isPlainObject(item)) {
      throw invalid("invalid-freshness", "freshness entries must be plain objects");
    }
    rejectUnknown(item, new Set(["source", "observedAt"]));
    const source = requiredText("freshness source", item.source);
    if (seen.has(source)) {
      throw invalid("duplicate-freshness-source", `duplicate freshness source: ${source}`);
    }
    seen.add(source);
    result.push({ source, observedAt: normalizeInstant("freshness observedAt", item.observedAt) });
  }
  result.sort((a, b) => a.source.localeCompare(b.source));
  return freezeDeep(result);
}

function normalizeList(label, value) {
  if (!Array.isArray(value) || value.length === 0) {
    throw invalid("invalid-list", `${label} list must be non-empty`);
  }
  const set = new Set(value.map((item) => requiredText(label, item)));
  return Object.freeze([...set].sort());
}

function normalizeInternalPath(value) {
  if (value === undefined || value === null || value === "") {
    return undefined;
  }
  const path = requiredText("internal path", value);
  if (!path.startsWith("/") || path.startsWith("//") || path.includes("\\")) {
    throw invalid("unsafe-internal-path", "internal path must stay inside the CMDR origin");
  }
  const resolved = new URL(path, SAFE_PATH_ORIGIN);
  if (resolved.origin !== SAFE_PATH_ORIGIN) {
    throw invalid("unsafe-internal-path", "internal path must stay inside the CMDR origin");
  }
  return `${resolved.pathname}${resolved.search}${resolved.hash}`;
}

function normalizeInstant(label, value) {
  const text = requiredText(label, value);
  const millis = Date.parse(text);
  if (!Number.isFinite(millis)) {
    throw invalid("invalid-time", `${label} must be an RFC3339-compatible instant`);
  }
  return new Date(millis).toISOString();
}

function optionalText(value) {
  if (value === undefined || value === null || value === "") {
    return undefined;
  }
  return requiredText("text", value);
}

function requiredText(label, value) {
  if (typeof value !== "string") {
    throw invalid("invalid-text", `${label} must be a string`);
  }
  const normalized = value.trim();
  if (normalized === "" || /[\u0000-\u001f\u007f]/u.test(normalized)) {
    throw invalid("invalid-text", `${label} must be non-empty and control-free`);
  }
  return normalized;
}

function requireBoolean(label, value) {
  if (typeof value !== "boolean") {
    throw invalid("invalid-boolean", `${label} must be boolean`);
  }
  return value;
}

function rejectUnknown(input, allowed) {
  for (const key of Object.keys(input)) {
    if (!allowed.has(key)) {
      throw invalid("unknown-shell-field", `unknown shell-state field: ${key}`);
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
  return new ShellValidationError(code, message);
}
