"use strict";

const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");
const vm = require("node:vm");
const { test } = require("node:test");

const source = fs.readFileSync(path.join(__dirname, "analyst.js"), "utf8");

function node(tag = "div") {
  return {
    tag, value: "", disabled: false, hidden: false, textContent: "",
    children: [], handlers: {},
    addEventListener(name, listener) { this.handlers[name] = listener; },
    append(...children) { this.children.push(...children); },
    add(child) { this.children.push(child); },
    replaceChildren() { this.children = []; }
  };
}

function boot(search = "", navigationQuery, denySource = false) {
  const ids = [
    "query", "token", "case-id", "kind", "search", "case-options",
    "discover", "case-discovery-status", "from", "to", "entries",
    "summary", "message", "page", "previous", "next", "load",
    "source-detail", "source-status", "return-to-case", "copy-return-link",
    "analyze", "analysis", "analysis-actions", "download-report",
    "wazuh-cohorts", "wazuh-status", "wazuh-actions", "wazuh-events",
    "wazuh-previous", "wazuh-next",
    "wazuh-search-field", "wazuh-search-value", "wazuh-search-from",
    "wazuh-search-to", "wazuh-search-run", "wazuh-search-status", "wazuh-search-results",
  "wazuh-search-bursts",
  "wazuh-replay-save", "wazuh-replay-file",
  "local-search-field", "local-search-family", "local-search-value", "local-search-from", "local-search-to", "local-search-run", "local-search-status", "local-search-results",
    "compare-case", "compare", "comparison", "comparison-actions", "download-comparison",
    "trace-actions", "trace-results",
    "comparison-trace-actions", "comparison-trace-results", "download-cross-events",
    "proximity", "proximity-seconds", "proximity-results"
  ];
  const fields = Object.fromEntries(ids.map(id => [id, node(id)]));
  const requests = [];
  const copied = [];
  const scope = {
    window: { location: { origin: "http://127.0.0.1:8765", search } },
    document: {
      getElementById(id) {
        if (!fields[id]) throw new Error("Unknown DOM id " + id);
        return fields[id];
      },
      createElement: tag => node(tag),
      createTextNode: text => ({ tag: "text", textContent: text })
    },
    navigator: { clipboard: { async writeText(value) { copied.push(value); } } },
    URL, URLSearchParams, AbortController,
    Option: class { constructor(text, value) { this.text = text; this.value = value; } },
    Event: class { constructor(type) { this.type = type; } },
    fetch: async (address, options) => {
      const url = new URL(String(address));
      requests.push({ url, options });
      if (url.pathname === "/case-timeline") {
        const offset = Number(url.searchParams.get("offset"));
        return {
          ok: true,
          async json() {
            return {
              status: "allow",
              page: { total_visible: 1, has_more: false, offset },
              entries: offset > 0 ? [] : [{
                id: "event-new", kind: "observed",
                source_time: "2026-10-01T11:00:00Z",
                source_ref: "demo-telemetry-a", provenance: "synthetic-telemetry",
                version: 2
              }]
            };
          }
        };
      }
      if (url.pathname === "/case-source-navigation") {
        const expected = new URLSearchParams(url.searchParams);
        expected.delete("entry_id");
        expected.sort();
        return {
          ok: true,
          async json() {
            return {
              status: "ready", kind: "telemetry-event", object_id: "demo-event-a",
              return_is_authority: false,
              return_case_id: url.searchParams.get("case_id"),
              return_entry_id: url.searchParams.get("entry_id"),
              return_query: navigationQuery || expected.toString()
            };
          }
        };
      }
      if (url.pathname === "/demo-source-open") {
        return {
          ok: !denySource,
          async json() {
            return {
              status: "ready", synthetic: true,
              case_id: url.searchParams.get("case_id"),
              entry_id: url.searchParams.get("entry_id"),
              kind: "telemetry-event", object_id: "demo-event-a",
              preview: "Événement réseau synthétique autorisé."
            };
          }
        };
      }
      throw new Error("Unexpected outbound route: " + url.pathname);
    }
  };
  vm.runInNewContext(source, scope, { filename: "analyst.js" });
  return { fields, requests, copied };
}

function flush() {
  return new Promise(resolve => setImmediate(resolve));
}

test("restored Case selection requires new bearer and preserves UTC time and page", async () => {
  const url = new URL("http://127.0.0.1:8765/analyst");
  url.searchParams.set("case_id", "case-a");
  url.searchParams.set("kind", "observed");
  url.searchParams.set("limit", "20");
  url.searchParams.set("offset", "20");
  url.searchParams.set("from", "2026-10-25T01:30:00.000Z");
  url.searchParams.set("to", "2026-10-25T03:30:00.000Z");
  const env = boot(url.search);
  assert.equal(env.requests.length, 0, "return link may not fetch automatically");
  assert.equal(env.fields["case-id"].value, "case-a");
  assert.equal(env.fields["token"].value, "", "bearer must not be restored");
  env.fields["token"].value = "ephemeral-value";
  env.fields["token"].handlers.input();
  env.fields.query.handlers.submit({ preventDefault() {} });
  await flush();
  assert.equal(env.requests.length, 1);
  const request = env.requests[0];
  assert.equal(request.url.pathname, "/case-timeline");
  assert.equal(request.url.searchParams.get("offset"), "20");
  assert.equal(request.url.searchParams.get("from"), "2026-10-25T01:30:00.000Z");
  assert.equal(request.url.searchParams.get("to"), "2026-10-25T03:30:00.000Z");
  assert.equal(request.options.headers.Authorization, "Bearer ephemeral-value");
  assert.equal(request.url.href.includes("ephemeral-value"), false);
});

test("source draft produces explicit same-origin return without authority", async () => {
  const env = boot();
  env.fields["token"].value = "ephemeral-value";
  env.fields["case-id"].value = "case-a";
  env.fields.kind.value = "observed";
  env.fields.query.handlers.submit({ preventDefault() {} });
  await flush();
  assert.equal(env.fields.entries.children.length, 1);
  const item = env.fields.entries.children[0];
  const open = item.children.find(child => child.tag === "button");
  assert.ok(open, "each visible event needs a source action");
  open.handlers.click();
  await flush();
  assert.equal(env.requests.length, 3);
  assert.equal(env.requests[1].url.pathname, "/case-source-navigation");
  assert.equal(env.requests[2].url.pathname, "/demo-source-open");
  assert.equal(env.requests[2].options.headers.Authorization, "Bearer ephemeral-value");
  assert.equal(env.requests[2].url.searchParams.has("object_id"), false);
  assert.equal(env.fields["source-detail"].hidden, false);
  assert.equal(env.fields["copy-return-link"].disabled, false);
  await env.fields["copy-return-link"].handlers.click();
  assert.equal(env.copied.length, 1);
  const url = new URL(env.copied[0]);
  assert.equal(url.origin, "http://127.0.0.1:8765");
  assert.equal(url.pathname, "/analyst");
  assert.equal(url.searchParams.get("case_id"), "case-a");
  assert.equal(url.searchParams.has("entry_id"), false);
  assert.equal(url.href.includes("ephemeral-value"), false);
  assert.equal(url.href.includes("demo-event-a"), false);
  assert.equal(boot(url.search).requests.length, 0);
});

test("tampered source return selection never enables copy", async () => {
  const env = boot("", "case_id=case-b&kind=observed&limit=20&offset=0");
  env.fields["token"].value = "ephemeral-value";
  env.fields["case-id"].value = "case-a";
  env.fields.kind.value = "observed";
  env.fields.query.handlers.submit({ preventDefault() {} });
  await flush();
  env.fields.entries.children[0].children.find(n => n.tag === "button").handlers.click();
  await flush();
  assert.equal(env.fields["copy-return-link"].disabled, true);
  await env.fields["copy-return-link"].handlers.click();
  assert.deepEqual(env.copied, []);
});

test("unknown query selector is inert and cannot auto-authorize", () => {
  const env = boot("?case_id=case-a&tenant=other&limit=20");
  assert.equal(env.fields["case-id"].value, "");
  assert.equal(env.requests.length, 0);
  assert.match(env.fields.message.textContent, /invalide/);
});

test("denied owner opening cannot be treated as a ready source draft", async () => {
  const env = boot("", undefined, true);
  env.fields.token.value = "ephemeral-value";
  env.fields["case-id"].value = "case-a";
  env.fields.query.handlers.submit({ preventDefault() {} });
  await flush();
  env.fields.entries.children[0].children.find(n => n.tag === "button").handlers.click();
  await flush();
  assert.equal(env.requests.length, 3);
  assert.equal(env.fields["source-detail"].hidden, true);
  assert.equal(env.fields["copy-return-link"].disabled, true);
  assert.equal(env.fields["return-to-case"].disabled, true);
  assert.deepEqual(env.copied, []);
});
