const VERSION = 1;
const CONTROL = /[\u0000-\u001f\u007f]/u;
const FILTER_KEY = /^[A-Za-z0-9_.:-]{1,96}$/u;
const DANGEROUS_FILTER_KEYS = new Set(["__proto__", "prototype", "constructor"]);
const DEEP_LINK_KEYS = new Set([
  "version",
  "tenantId",
  "environmentId",
  "from",
  "to",
  "query",
  "sources",
  "filters",
  "activeRef",
  "returnTo",
]);
const URL_SINGLETON_KEYS = new Set([
  "v",
  "tenant",
  "environment",
  "from",
  "to",
  "q",
  "object",
  "return",
]);

export class StateValidationError extends Error {
  constructor(code, message) {
    super(message);
    this.name = "StateValidationError";
    this.code = code;
  }
}

export function normalizeDeepLinkState(input, expected = {}) {
  if (!isPlainObject(input)) {
    throw invalid("invalid-state", "deep-link state must be a plain object");
  }
  for (const key of Object.keys(input)) {
    if (!DEEP_LINK_KEYS.has(key)) {
      throw invalid("unknown-state-field", `unknown deep-link state field: ${key}`);
    }
  }
  if (input.version !== undefined && input.version !== VERSION) {
    throw invalid("unsupported-version", "unsupported deep-link state version");
  }

  const tenantId = requiredScope("tenant", input.tenantId);
  const environmentId = requiredScope("environment", input.environmentId);
  if (expected.tenantId !== undefined && tenantId !== requiredScope("expected tenant", expected.tenantId)) {
    throw invalid("tenant-mismatch", "deep link tenant does not match the selected tenant");
  }
  if (
    expected.environmentId !== undefined &&
    environmentId !== requiredScope("expected environment", expected.environmentId)
  ) {
    throw invalid("environment-mismatch", "deep link environment does not match the selected environment");
  }

  const from = normalizeInstant("from", input.from);
  const to = normalizeInstant("to", input.to);
  if (Date.parse(from) >= Date.parse(to)) {
    throw invalid("invalid-time-range", "deep-link time range must be ordered");
  }

  const query = requiredText("query", input.query);
  const sources = normalizeSources(input.sources);
  const filters = normalizeFilters(input.filters ?? Object.create(null));
  const activeRef = optionalOpaque("active object", input.activeRef);
  const returnTo = normalizeReturnPath(input.returnTo);

  return deepFreeze({
    version: VERSION,
    tenantId,
    environmentId,
    from,
    to,
    query,
    sources,
    filters,
    ...(activeRef === undefined ? {} : { activeRef }),
    ...(returnTo === undefined ? {} : { returnTo }),
  });
}

export function serializeDeepLinkState(input, expected = {}) {
  const state = normalizeDeepLinkState(input, expected);
  const params = new URLSearchParams();
  params.set("v", String(VERSION));
  params.set("tenant", state.tenantId);
  params.set("environment", state.environmentId);
  params.set("from", state.from);
  params.set("to", state.to);
  params.set("q", state.query);
  for (const source of state.sources) {
    params.append("source", source);
  }
  for (const key of Object.keys(state.filters).sort()) {
    for (const value of state.filters[key]) {
      params.append(`filter.${key}`, value);
    }
  }
  if (state.activeRef !== undefined) {
    params.set("object", state.activeRef);
  }
  if (state.returnTo !== undefined) {
    params.set("return", state.returnTo);
  }
  return `?${params.toString()}`;
}

export function parseDeepLink(search, expected = {}) {
  if (typeof search !== "string") {
    throw invalid("invalid-url-state", "deep-link query must be a string");
  }
  const raw = search.startsWith("?") ? search.slice(1) : search;
  const params = new URLSearchParams(raw);
  const seenSingleton = new Set();
  const sources = [];
  const filters = Object.create(null);
  const input = { version: VERSION, sources, filters };

  for (const [key, value] of params) {
    if (URL_SINGLETON_KEYS.has(key)) {
      if (seenSingleton.has(key)) {
        throw invalid("duplicate-url-field", `duplicate deep-link field: ${key}`);
      }
      seenSingleton.add(key);
      switch (key) {
        case "v":
          if (value !== String(VERSION)) {
            throw invalid("unsupported-version", "unsupported deep-link URL version");
          }
          break;
        case "tenant":
          input.tenantId = value;
          break;
        case "environment":
          input.environmentId = value;
          break;
        case "from":
          input.from = value;
          break;
        case "to":
          input.to = value;
          break;
        case "q":
          input.query = value;
          break;
        case "object":
          input.activeRef = value;
          break;
        case "return":
          input.returnTo = value;
          break;
      }
      continue;
    }
    if (key === "source") {
      sources.push(value);
      continue;
    }
    if (key.startsWith("filter.")) {
      const filterKey = key.slice("filter.".length);
      if (!Object.hasOwn(filters, filterKey)) {
        filters[filterKey] = [];
      }
      filters[filterKey].push(value);
      continue;
    }
    throw invalid("unknown-url-field", `unknown deep-link URL field: ${key}`);
  }
  return normalizeDeepLinkState(input, expected);
}

function normalizeSources(value) {
  if (!Array.isArray(value) || value.length === 0) {
    throw invalid("empty-sources", "at least one explicit source is required");
  }
  const sources = new Set();
  for (const item of value) {
    const source = requiredText("source", item);
    if (source === "*") {
      throw invalid("wildcard-source", "wildcard source scope is forbidden");
    }
    sources.add(source);
  }
  return Object.freeze([...sources].sort());
}

function normalizeFilters(value) {
  if (!isPlainObject(value)) {
    throw invalid("invalid-filters", "filters must be a plain object");
  }
  const out = Object.create(null);
  for (const key of Object.keys(value).sort()) {
    if (DANGEROUS_FILTER_KEYS.has(key) || !FILTER_KEY.test(key)) {
      throw invalid("invalid-filter-key", `invalid filter key: ${key}`);
    }
    const rawValues = Array.isArray(value[key]) ? value[key] : [value[key]];
    if (rawValues.length === 0) {
      throw invalid("empty-filter", `filter ${key} has no values`);
    }
    const normalized = new Set();
    for (const raw of rawValues) {
      normalized.add(requiredText(`filter ${key}`, raw));
    }
    out[key] = Object.freeze([...normalized].sort());
  }
  return Object.freeze(out);
}

function requiredScope(label, value) {
  const result = requiredText(label, value);
  if (result === "*") {
    throw invalid("wildcard-scope", `${label} wildcard is forbidden`);
  }
  return result;
}

function requiredText(label, value) {
  if (typeof value !== "string") {
    throw invalid("invalid-text", `${label} must be a string`);
  }
  const normalized = value.trim();
  if (normalized === "" || CONTROL.test(normalized)) {
    throw invalid("invalid-text", `${label} must be non-empty and control-free`);
  }
  return normalized;
}

function optionalOpaque(label, value) {
  if (value === undefined || value === null || value === "") {
    return undefined;
  }
  return requiredText(label, value);
}

function normalizeInstant(label, value) {
  const text = requiredText(label, value);
  const milliseconds = Date.parse(text);
  if (!Number.isFinite(milliseconds)) {
    throw invalid("invalid-time", `${label} must be an RFC3339-compatible instant`);
  }
  return new Date(milliseconds).toISOString();
}

function normalizeReturnPath(value) {
  if (value === undefined || value === null || value === "") {
    return undefined;
  }
  const path = requiredText("return path", value);
  if (!path.startsWith("/") || path.startsWith("//") || path.includes("\\")) {
    throw invalid("unsafe-return-path", "return path must stay inside the CMDR origin");
  }
  const resolved = new URL(path, "https://cmdr.invalid");
  if (resolved.origin !== "https://cmdr.invalid") {
    throw invalid("unsafe-return-path", "return path must stay inside the CMDR origin");
  }
  return `${resolved.pathname}${resolved.search}${resolved.hash}`;
}

function isPlainObject(value) {
  if (value === null || typeof value !== "object" || Array.isArray(value)) {
    return false;
  }
  const prototype = Object.getPrototypeOf(value);
  return prototype === Object.prototype || prototype === null;
}

function deepFreeze(value) {
  if (value && typeof value === "object" && !Object.isFrozen(value)) {
    for (const item of Object.values(value)) {
      deepFreeze(item);
    }
    Object.freeze(value);
  }
  return value;
}

function invalid(code, message) {
  return new StateValidationError(code, message);
}
