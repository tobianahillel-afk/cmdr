import { normalizeDeepLinkState } from "../../../product-runtime/event-search-frontend/state.mjs";
import {
  createResultRenderer,
  createResultWindow,
} from "../../../product-runtime/event-search-frontend/results.mjs";

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
        if (index >= 0) node.parent.children.splice(index, 1);
      }
      node.parent = this;
      this.children.push(node);
    }
  }

  remove() {
    if (this.parent === undefined) return;
    const index = this.parent.children.indexOf(this);
    if (index >= 0) this.parent.children.splice(index, 1);
    this.parent = undefined;
  }
}


const args = parseArgs(process.argv.slice(2));
const iterations = boundedInteger(args.iterations ?? "2000", 1, 1_000_000);
const mode = args.mode ?? "latency";
if (mode !== "latency" && mode !== "resource") {
  throw new Error("mode must be latency or resource");
}

const document = new FakeDocument();
const mount = new FakeElement("div");
const renderer = createResultRenderer(document, mount);
const baseState = Object.freeze({
  tenantId: "tenant-a",
  environmentId: "prod-eu",
  from: "2026-09-24T10:00:00Z",
  to: "2026-09-24T11:00:00Z",
  query: "backend-neutral benchmark query",
  sources: ["edr", "firewall"],
  filters: { severity: ["high"] },
});
const expected = Object.freeze({ tenantId: "tenant-a", environmentId: "prod-eu" });
const rowsA = makeRows("a", 50);
const rowsB = makeRows("b", 50);
const windowA = createResultWindow(rowsA, { size: 50 });
const windowB = createResultWindow(rowsB, { size: 50 });

for (let i = 0; i < 200; i += 1) {
  normalizeDeepLinkState({ ...baseState, query: `warmup-${i & 1}` }, expected);
  renderer.render((i & 1) === 0 ? windowA : windowB);
}

const before = process.memoryUsage();
let peakHeapBytes = before.heapUsed;
let blackhole = 0;
const started = process.hrtime.bigint();
for (let i = 0; i < iterations; i += 1) {
  const state = normalizeDeepLinkState({ ...baseState, query: `q-${i & 1}` }, expected);
  const rendered = renderer.render((i & 1) === 0 ? windowA : windowB);
  blackhole ^= state.query.length + rendered.rendered;
  if ((i & 127) === 0) {
    peakHeapBytes = Math.max(peakHeapBytes, process.memoryUsage().heapUsed);
  }
}
const elapsed = process.hrtime.bigint() - started;
const after = process.memoryUsage();
peakHeapBytes = Math.max(peakHeapBytes, after.heapUsed);
if (blackhole === Number.MIN_SAFE_INTEGER) {
  throw new Error("unreachable benchmark guard");
}

const elapsedNS = Number(elapsed);
process.stdout.write(JSON.stringify({
  operation: "frontend-state-update",
  mode,
  operations: iterations,
  elapsed_ns: elapsedNS,
  ns_per_operation: elapsedNS / iterations,
  peak_heap_bytes: peakHeapBytes,
  total_allocation_bytes: Math.max(0, after.heapUsed - before.heapUsed),
}) + "\n");

function makeRows(prefix, count) {
  return Array.from({ length: count }, (_, index) => ({
    ref: `evt-${prefix}-${String(index).padStart(3, "0")}`,
    eventTime: `2026-09-24T10:${String(index % 60).padStart(2, "0")}:00Z`,
    source: index % 2 === 0 ? "edr" : "firewall",
    summary: {
      severity: ["low", "medium", "high", "critical"][index % 4],
      message: `Synthetic event ${prefix}-${index}`,
    },
  }));
}

function parseArgs(values) {
  const out = Object.create(null);
  for (let i = 0; i < values.length; i += 2) {
    const key = values[i];
    const value = values[i + 1];
    if (!key?.startsWith("--") || value === undefined) {
      throw new Error("arguments must use --key value pairs");
    }
    out[key.slice(2)] = value;
  }
  return out;
}

function boundedInteger(value, min, max) {
  const parsed = Number.parseInt(value, 10);
  if (!Number.isSafeInteger(parsed) || String(parsed) !== String(value) || parsed < min || parsed > max) {
    throw new Error(`integer must be within ${min}..${max}`);
  }
  return parsed;
}
