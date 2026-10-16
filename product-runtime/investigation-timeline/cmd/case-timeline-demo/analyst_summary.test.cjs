"use strict";

const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");
const vm = require("node:vm");
const { test } = require("node:test");

const source = fs.readFileSync(path.join(__dirname, "analyst.js"), "utf8");
const ids = [
  "query", "token", "case-id", "kind", "search", "case-options", "discover",
  "case-discovery-status", "from", "to", "entries", "summary", "message",
  "page", "previous", "next", "load", "source-detail", "source-status",
  "return-to-case", "copy-return-link", "analyze", "analysis"
];
function node() {
  return {
    value: "", textContent: "", disabled: false, hidden: false,
    handlers: {}, children: [],
    addEventListener(name, fn) { this.handlers[name] = fn; },
    add(child) { this.children.push(child); },
    append(...items) { this.children.push(...items); },
    replaceChildren() { this.children = []; },
    dispatchEvent(event) { this.handlers[event.type]?.(); }
  };
}
function boot(fetchMock) {
  const fields = Object.fromEntries(ids.map(id => [id, node()]));
  const requests = [];
  const ctx = {
    window: { location: { origin: "http://127.0.0.1:8765", search: "" } },
    document: {
      getElementById: id => fields[id],
      createElement: () => node(),
      createTextNode: text => ({ textContent: text })
    },
    navigator: { clipboard: { writeText: async () => {} } },
    Option: class { constructor(text, value) { this.text = text; this.value = value; } },
    Event: class { constructor(type) { this.type = type; } },
    URL, URLSearchParams, AbortController,
    fetch: async (address, options) => {
      const url = new URL(String(address));
      requests.push({ url, options });
      return fetchMock(url, options, requests.length);
    }
  };
  vm.runInNewContext(source, ctx);
  fields.token.value = "ephemeral-token";
  fields["case-id"].value = "case-disk";
  return { fields, requests };
}
const flush = () => new Promise(resolve => setImmediate(resolve));
function page(offset, total, options = {}) {
  const count = Math.min(200, Math.max(0, total - offset));
  return {
    ok: true,
    async json() {
      return {
        status: "allow",
        page: {
          limit: 200, offset, total_visible: options.total ?? total,
          has_more: offset + count < total
        },
        entries: Array.from({ length: count }, (_, n) => {
          const index = offset + n;
          return {
            id: "event-" + index,
            kind: index % 2 ? "inferred" : "observed",
            time_uncertain: index === 2
          };
        })
      };
    }
  };
}

test("analyst summarizes authorized Case across two paged HTTP responses", async () => {
  const env = boot(url => page(Number(url.searchParams.get("offset")), 205));
  env.fields.kind.value = "";
  env.fields.search.value = "login";
  env.fields.analyze.handlers.click();
  await flush();
  await flush();
  assert.equal(env.requests.length, 2);
  assert.deepEqual(env.requests.map(x => x.url.searchParams.get("offset")), ["0", "200"]);
  for (const request of env.requests) {
    assert.equal(request.url.pathname, "/case-timeline");
    assert.equal(request.url.searchParams.get("case_id"), "case-disk");
    assert.equal(request.url.searchParams.get("q"), "login");
    assert.equal(request.url.searchParams.get("limit"), "200");
    assert.equal(request.options.headers.Authorization, "Bearer ephemeral-token");
    assert.equal(request.options.credentials, "omit");
    assert.equal(request.options.cache, "no-store");
    assert.ok(!request.url.href.includes("ephemeral-token"));
  }
  assert.match(env.fields.analysis.textContent, /205 événement\(s\)/);
  assert.match(env.fields.analysis.textContent, /observed : 103/);
  assert.match(env.fields.analysis.textContent, /inferred : 102/);
  assert.match(env.fields.analysis.textContent, /horodatage incertain : 1/);
  assert.equal(env.fields.analyze.disabled, false);
});

test("inconsistent page totals fail closed without showing a partial aggregate", async () => {
  const env = boot(url => {
    const offset = Number(url.searchParams.get("offset"));
    return page(offset, 201, offset === 200 ? { total: 200 } : {});
  });
  env.fields.analyze.handlers.click();
  await flush();
  await flush();
  assert.equal(env.requests.length, 2);
  assert.match(env.fields.analysis.textContent, /Aucun total partiel/);
  assert.doesNotMatch(env.fields.analysis.textContent, /observed :/);
});

test("revocation during an in-flight page and token change cancel the summary", async () => {
  let release;
  const second = new Promise(resolve => { release = resolve; });
  const env = boot(url => Number(url.searchParams.get("offset")) === 0
    ? page(0, 201) : second);
  env.fields.analyze.handlers.click();
  await flush();
  assert.equal(env.requests.length, 2);
  env.fields.token.value = "different-token";
  env.fields.token.handlers.input();
  release(page(200, 201));
  await flush();
  assert.match(env.fields.analysis.textContent, /Aucune synthèse calculée/);
  assert.equal(env.fields.analyze.disabled, false);
  assert.equal(env.requests[1].options.signal.aborted, true);
});

test("empty authorized Case yields a zero-only summary with one request", async () => {
  const env = boot(() => page(0, 0));
  env.fields.analyze.handlers.click();
  await flush();
  assert.equal(env.requests.length, 1);
  assert.match(env.fields.analysis.textContent, /0 événement\(s\)/);
});

test("missing bearer never generates an HTTP request", async () => {
  const env = boot(() => { throw Error("request forbidden"); });
  env.fields.token.value = "";
  env.fields.analyze.handlers.click();
  await flush();
  assert.equal(env.requests.length, 0);
  assert.match(env.fields.analysis.textContent, /Sélection incomplète/);
});
