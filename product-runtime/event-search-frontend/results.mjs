const MAX_RESULT_WINDOW = 200;
const MAX_SUMMARY_FIELDS = 24;
const MAX_TEXT_LENGTH = 4096;
const SAFE_FIELD = /^[A-Za-z0-9_.:-]{1,128}$/u;
const BLOCKED_FIELD_KEYS = new Set(["__proto__", "prototype", "constructor"]);

export class ResultValidationError extends Error {
  constructor(code, message) {
    super(message);
    this.name = "ResultValidationError";
    this.code = code;
  }
}

export function normalizeResultProjection(input) {
  if (!isPlainObject(input)) {
    throw invalid("invalid-result", "result projection must be a plain object");
  }
  rejectUnknown(input, new Set(["ref", "eventTime", "source", "summary"]));
  const summary = normalizeSummary(input.summary ?? {});
  return freezeDeep({
    ref: requiredText("result ref", input.ref, 512),
    eventTime: normalizeInstant(input.eventTime),
    source: requiredText("result source", input.source, 256),
    summary,
  });
}

export function normalizeResultCollection(items) {
  if (!Array.isArray(items)) {
    throw invalid("invalid-results", "results must be an array");
  }
  const refs = new Set();
  const normalized = [];
  for (const item of items) {
    const result = normalizeResultProjection(item);
    if (refs.has(result.ref)) {
      throw invalid("duplicate-result-ref", `duplicate result ref: ${result.ref}`);
    }
    refs.add(result.ref);
    normalized.push(result);
  }
  return Object.freeze(normalized);
}

export function createResultWindow(items, options = {}) {
  const normalized = normalizeResultCollection(items);
  if (!isPlainObject(options)) {
    throw invalid("invalid-window", "window options must be a plain object");
  }
  rejectUnknown(options, new Set(["start", "size", "anchorRef", "anchorOffset"]));
  const size = boundedInteger("window size", options.size ?? 50, 1, MAX_RESULT_WINDOW);
  let start = boundedInteger("window start", options.start ?? 0, 0, Number.MAX_SAFE_INTEGER);
  let anchorRef = optionalText(options.anchorRef, 512);
  let anchorOffset = boundedInteger("anchor offset", options.anchorOffset ?? 0, 0, size - 1);

  if (anchorRef !== undefined) {
    const anchorIndex = normalized.findIndex((item) => item.ref === anchorRef);
    if (anchorIndex >= 0) {
      start = Math.max(0, anchorIndex - anchorOffset);
    } else {
      anchorRef = undefined;
      anchorOffset = 0;
    }
  }
  start = Math.min(start, normalized.length);
  const end = Math.min(normalized.length, start + size);
  const rows = Object.freeze(normalized.slice(start, end));
  if (anchorRef === undefined && rows.length > 0) {
    anchorRef = rows[0].ref;
    anchorOffset = 0;
  } else if (anchorRef !== undefined) {
    const visibleIndex = rows.findIndex((row) => row.ref === anchorRef);
    if (visibleIndex < 0) {
      anchorRef = rows[0]?.ref;
      anchorOffset = 0;
    } else {
      anchorOffset = visibleIndex;
    }
  }

  return freezeDeep({
    total: normalized.length,
    start,
    end,
    size,
    beforeCount: start,
    afterCount: normalized.length - end,
    anchorRef,
    anchorOffset,
    rows,
  });
}

export function anchorResultWindow(windowState, visibleOffset) {
  validateWindowState(windowState);
  if (windowState.rows.length === 0) {
    return windowState;
  }
  const offset = boundedInteger("visible anchor offset", visibleOffset, 0, windowState.rows.length - 1);
  return freezeDeep({
    ...windowState,
    anchorRef: windowState.rows[offset].ref,
    anchorOffset: offset,
  });
}

export function reconcileResultWindow(previousWindow, nextItems) {
  validateWindowState(previousWindow);
  return createResultWindow(nextItems, {
    start: previousWindow.start,
    size: previousWindow.size,
    anchorRef: previousWindow.anchorRef,
    anchorOffset: previousWindow.anchorOffset,
  });
}

export function createIncrementalScheduler(schedule = queueMicrotask) {
  if (typeof schedule !== "function") {
    throw invalid("invalid-scheduler", "schedule must be a function");
  }
  let pending = false;
  let latest;

  return Object.freeze({
    submit(task) {
      if (typeof task !== "function") {
        throw invalid("invalid-task", "scheduled update must be a function");
      }
      latest = task;
      if (pending) {
        return;
      }
      pending = true;
      schedule(() => {
        pending = false;
        const run = latest;
        latest = undefined;
        if (run !== undefined) {
          run();
        }
      });
    },
    cancel() {
      latest = undefined;
    },
    isPending() {
      return pending;
    },
  });
}

export function createResultRenderer(document, mount) {
  if (!document || typeof document.createElement !== "function") {
    throw invalid("invalid-document", "a DOM-compatible document is required");
  }
  if (!mount || typeof mount.append !== "function") {
    throw invalid("invalid-mount", "a DOM-compatible mount point is required");
  }

  const section = element(document, "section", { "aria-labelledby": "event-search-results-heading" });
  const heading = element(document, "h2", { id: "event-search-results-heading" }, "Search results");
  const status = element(document, "p", { "aria-live": "polite" });
  const table = element(document, "table", { "aria-label": "Event Search results" });
  const caption = element(document, "caption", {}, "Windowed telemetry result projections");
  const head = element(document, "thead");
  const headerRow = element(document, "tr");
  for (const label of ["Event time", "Source", "Event reference", "Summary"]) {
    headerRow.append(element(document, "th", { scope: "col" }, label));
  }
  head.append(headerRow);
  const body = element(document, "tbody");
  table.append(caption, head, body);
  section.append(heading, status, table);
  mount.append(section);

  const nodes = new Map();

  return Object.freeze({
    render(windowState) {
      validateWindowState(windowState);
      const scrollTop = Number.isFinite(mount.scrollTop) ? mount.scrollTop : undefined;
      table.setAttribute("aria-rowcount", String(windowState.total));
      status.textContent = windowState.total === 0
        ? "No result projections are available."
        : `Showing results ${windowState.start + 1}–${windowState.end} of ${windowState.total}.`;

      const desired = new Set();
      for (let index = 0; index < windowState.rows.length; index += 1) {
        const result = windowState.rows[index];
        desired.add(result.ref);
        let row = nodes.get(result.ref);
        if (row === undefined) {
          row = createResultRow(document, result);
          nodes.set(result.ref, row);
        } else {
          updateResultRow(row, result);
        }
        row.setAttribute("aria-rowindex", String(windowState.start + index + 1));
        body.append(row);
      }

      for (const [ref, row] of [...nodes.entries()]) {
        if (!desired.has(ref)) {
          row.remove();
          nodes.delete(ref);
        }
      }
      if (scrollTop !== undefined) {
        mount.scrollTop = scrollTop;
      }
      return freezeDeep({
        rendered: windowState.rows.length,
        total: windowState.total,
        beforeCount: windowState.beforeCount,
        afterCount: windowState.afterCount,
      });
    },
    getRowNode(ref) {
      return nodes.get(ref);
    },
    visibleRefs() {
      return Object.freeze([...nodes.keys()]);
    },
  });
}

function createResultRow(document, result) {
  const row = element(document, "tr", { "data-result-ref": result.ref });
  row.append(
    element(document, "td"),
    element(document, "td"),
    element(document, "td"),
    element(document, "td"),
  );
  updateResultRow(row, result);
  return row;
}

function updateResultRow(row, result) {
  if (!Array.isArray(row.children) && typeof row.children?.length !== "number") {
    throw invalid("invalid-row", "result row must expose child cells");
  }
  const cells = [...row.children];
  if (cells.length !== 4) {
    throw invalid("invalid-row", "result row must contain four cells");
  }
  cells[0].textContent = result.eventTime;
  cells[1].textContent = result.source;
  cells[2].textContent = result.ref;
  cells[3].textContent = formatSummary(result.summary);
}

function formatSummary(summary) {
  return Object.entries(summary)
    .map(([key, value]) => `${key}=${value}`)
    .join(" · ");
}

function normalizeSummary(value) {
  if (!isPlainObject(value)) {
    throw invalid("invalid-summary", "result summary must be a plain object");
  }
  const entries = Object.entries(value);
  if (entries.length > MAX_SUMMARY_FIELDS) {
    throw invalid("summary-too-wide", `result summary exceeds ${MAX_SUMMARY_FIELDS} fields`);
  }
  const normalized = {};
  entries.sort(([left], [right]) => left.localeCompare(right));
  for (const [key, raw] of entries) {
    if (BLOCKED_FIELD_KEYS.has(key) || !SAFE_FIELD.test(key)) {
      throw invalid("unsafe-summary-key", `unsafe summary field: ${key}`);
    }
    normalized[key] = normalizeSummaryValue(raw);
  }
  return freezeDeep(normalized);
}

function normalizeSummaryValue(value) {
  if (value === null) {
    return "null";
  }
  if (typeof value === "string") {
    return requiredText("summary value", value, MAX_TEXT_LENGTH);
  }
  if (typeof value === "number" && Number.isFinite(value)) {
    return String(value);
  }
  if (typeof value === "boolean") {
    return String(value);
  }
  throw invalid("invalid-summary-value", "summary values must be scalar strings, finite numbers, booleans, or null");
}

function normalizeInstant(value) {
  const text = requiredText("event time", value, 128);
  const millis = Date.parse(text);
  if (!Number.isFinite(millis)) {
    throw invalid("invalid-event-time", "event time must be an RFC3339-compatible instant");
  }
  return new Date(millis).toISOString();
}

function optionalText(value, maxLength) {
  if (value === undefined || value === null || value === "") {
    return undefined;
  }
  return requiredText("text", value, maxLength);
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

function validateWindowState(value) {
  if (
    !value ||
    typeof value !== "object" ||
    !Array.isArray(value.rows) ||
    !Number.isSafeInteger(value.start) ||
    !Number.isSafeInteger(value.end) ||
    !Number.isSafeInteger(value.total) ||
    !Number.isSafeInteger(value.size)
  ) {
    throw invalid("invalid-window", "normalized result window is required");
  }
}

function rejectUnknown(input, allowed) {
  for (const key of Object.keys(input)) {
    if (!allowed.has(key)) {
      throw invalid("unknown-result-field", `unknown result field: ${key}`);
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

function invalid(code, message) {
  return new ResultValidationError(code, message);
}
