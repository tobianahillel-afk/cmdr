"use strict";

const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");
const vm = require("node:vm");
const { test } = require("node:test");
const script = fs.readFileSync(path.join(__dirname, "hunt.js"), "utf8");

function element() {
  return {
    value: "", textContent: "", disabled: false, children: [], handlers: {},
    addEventListener(event, handler) { this.handlers[event] = handler; },
    dispatchEvent(event) { this.handlers[event.type]?.(event); },
    append(child) { this.children.push(child); },
    replaceChildren() { this.children = []; }
  };
}

function setup(reply) {
  const ids = ["token", "case-id", "search", "query", "hunt-query",
    "hunt-run", "hunt-status", "hunt-results"];
  const fields = Object.fromEntries(ids.map(id => [id, element()]));
  const requests = [];
  let submits = 0;
  fields.query.addEventListener("submit", event => {
    submits++;
    event.preventDefault?.();
  });
  const context = {
    document: {
      getElementById(id) {
        if (!fields[id]) throw new Error("Missing DOM " + id);
        return fields[id];
      },
      createElement() { return element(); }
    },
    window: { location: { origin: "http://127.0.0.1:8765" } },
    URL, AbortController, Set,
    Event: class {
      constructor(type) { this.type = type; }
    },
    fetch: async (address, options) => {
      const url = new URL(String(address));
      requests.push({ url, options });
      return reply(url, options);
    }
  };
  vm.runInNewContext(script, context, { filename: "hunt.js" });
  fields.token.value = "session-bearer";
  fields["hunt-query"].value = "10.0.0.7";
  return { fields, requests, get submits() { return submits; } };
}

const flush = async () => {
  for (let i = 0; i < 7; i++) await new Promise(resolve => setImmediate(resolve));
};
const cases = names => ({ ok: true, async json() { return { cases: names }; } });
const page = (url, ids) => ({
  ok: true,
  async json() {
    return {
      status: "allow",
      page: { offset: Number(url.searchParams.get("offset")),
        limit: 200, total_visible: ids.length, has_more: false },
      entries: ids.map(id => ({ id, kind: "observed", source_ref: "source-" + id }))
    };
  }
});

test("authenticated hunt searches two authorized Cases then reauthorizes Case pivot", async () => {
  const env = setup((url) => {
    if (url.pathname === "/demo-cases") return cases(["case-a", "case-b"]);
    if (url.pathname === "/case-timeline") return page(url, [url.searchParams.get("case_id") + "-hit"]);
    throw new Error("unexpected route");
  });
  env.fields["hunt-run"].dispatchEvent({ type: "click" });
  await flush();
  assert.equal(env.requests.length, 3);
  assert.equal(env.fields["hunt-results"].children.length, 2);
  assert.match(env.fields["hunt-status"].textContent, /2 Case/);
  assert.ok(env.requests.every(x =>
    x.options.headers.Authorization === "Bearer session-bearer" &&
    x.options.credentials === "omit" && x.options.cache === "no-store" &&
    x.options.redirect === "error" && x.options.mode === "same-origin" &&
    !x.url.href.includes("session-bearer")));
  const timelineRequests = env.requests.slice(1);
  assert.deepEqual(timelineRequests.map(x => x.url.searchParams.get("case_id")), ["case-a", "case-b"]);
  assert.ok(timelineRequests.every(x =>
    x.url.searchParams.get("q") === "10.0.0.7" &&
    x.url.searchParams.get("limit") === "200"));
  env.fields["hunt-results"].children[1].dispatchEvent({ type: "click" });
  assert.equal(env.fields["case-id"].value, "case-b");
  assert.equal(env.fields.search.value, "10.0.0.7");
  assert.equal(env.submits, 1);
});

test("revoked second Case refuses all partial results", async () => {
  const env = setup(url => {
    if (url.pathname === "/demo-cases") return cases(["case-a", "case-b"]);
    if (url.searchParams.get("case_id") === "case-a") return page(url, ["visible-hit"]);
    return { ok: false, async json() { return { error: "private-source" }; } };
  });
  env.fields["hunt-run"].dispatchEvent({ type: "click" });
  await flush();
  assert.equal(env.fields["hunt-results"].children.length, 0);
  assert.match(env.fields["hunt-status"].textContent, /aucun résultat partiel/);
  assert.doesNotMatch(env.fields["hunt-status"].textContent, /private-source|visible-hit/);
});

test("bearer rotation cancels deferred discovery and stale results", async () => {
  let resolve;
  const deferred = new Promise(r => { resolve = r; });
  const env = setup(async url => {
    if (url.pathname === "/demo-cases") return deferred;
    throw new Error("stale request escaped discovery");
  });
  env.fields["hunt-run"].dispatchEvent({ type: "click" });
  await flush();
  env.fields.token.value = "new-bearer";
  env.fields.token.dispatchEvent({ type: "input" });
  resolve(cases(["case-a"]));
  await flush();
  assert.equal(env.requests.length, 1);
  assert.equal(env.fields["hunt-results"].children.length, 0);
  assert.equal(env.fields["hunt-run"].disabled, false);
});

test("over-limit Case discovery refuses to present misleading partial hunt", async () => {
  const env = setup(url => {
    assert.equal(url.pathname, "/demo-cases");
    return cases(Array.from({ length: 9 }, (_, i) => "case-" + i));
  });
  env.fields["hunt-run"].dispatchEvent({ type: "click" });
  await flush();
  assert.equal(env.requests.length, 1);
  assert.equal(env.fields["hunt-results"].children.length, 0);
  assert.match(env.fields["hunt-status"].textContent, /aucun résultat partiel/);
});

test("malformed authorized page or duplicated entries never produces hunt output", async () => {
  for (const mutation of ["wrong-total", "duplicate"]) {
    const env = setup(url => {
      if (url.pathname === "/demo-cases") return cases(["case-a"]);
      const result = page(url, ["one", "two"]);
      return { ok: true, async json() {
        const data = await result.json();
        if (mutation === "wrong-total") data.page.total_visible = 9000;
        else data.entries[1].id = "one";
        return data;
      } };
    });
    env.fields["hunt-run"].dispatchEvent({ type: "click" });
    await flush();
    assert.equal(env.fields["hunt-results"].children.length, 0);
  }
});

test("authorized shared source across two Cases creates reauthorized source pivot", async () => {
  const env = setup(url => {
    if (url.pathname === "/demo-cases") return cases(["case-a", "case-b"]);
    const id = url.searchParams.get("case_id");
    return {
      ok: true,
      async json() {
        return {
          status: "allow", page: { offset: 0, limit: 200, total_visible: 1, has_more: false },
          entries: [{
            id: id + "-hit", kind: id === "case-a" ? "observed" : "inferred",
            source_ref: "ioc-source:10.0.0.7"
          }]
        };
      }
    };
  });
  env.fields["hunt-run"].dispatchEvent({ type: "click" });
  await flush();
  const display = env.fields["hunt-results"].children;
  assert.equal(display.length, 4); // 2 Cases, title, 1 common source
  assert.match(display[0].textContent, /observés : 1.*inférés : 0/);
  assert.match(display[1].textContent, /observés : 0.*inférés : 1/);
  assert.match(display[2].textContent, /cooccurrence, pas causalité/);
  assert.match(display[3].textContent, /ioc-source:10\.0\.0\.7.*2 Cases, 2 événements/);
  display[3].dispatchEvent({ type: "click" });
  assert.equal(env.submits, 1);
  assert.equal(env.fields["case-id"].value, "case-a");
  assert.equal(env.fields.search.value, "ioc-source:10.0.0.7");
  assert.ok(env.requests.every(x => !x.url.href.includes("session-bearer")));
});

test("source correlation does not display a source seen in only one Case", async () => {
  const env = setup(url => {
    if (url.pathname === "/demo-cases") return cases(["case-a", "case-b"]);
    return page(url, [url.searchParams.get("case_id") + "-hit"]);
  });
  env.fields["hunt-run"].dispatchEvent({ type: "click" });
  await flush();
  assert.equal(env.fields["hunt-results"].children.length, 2);
});

test("invalid authorized classification or source reference rejects all Case summaries", async () => {
  for (const mutation of ["unknown-kind", "missing-source", "long-source"]) {
    const env = setup(url => {
      if (url.pathname === "/demo-cases") return cases(["case-a", "case-b"]);
      if (url.searchParams.get("case_id") === "case-a") return page(url, ["authorized"]);
      const response = page(url, ["later"]);
      return {
        ok: true,
        async json() {
          const data = await response.json();
          if (mutation === "unknown-kind") data.entries[0].kind = "private-type";
          if (mutation === "missing-source") data.entries[0].source_ref = "";
          if (mutation === "long-source") data.entries[0].source_ref = "X".repeat(1025);
          return data;
        }
      };
    });
    env.fields["hunt-run"].dispatchEvent({ type: "click" });
    await flush();
    assert.equal(env.fields["hunt-results"].children.length, 0, mutation);
    assert.match(env.fields["hunt-status"].textContent, /aucun résultat partiel/);
  }
});

test("source pivot becomes inert after token rotation", async () => {
  const env = setup(url => {
    if (url.pathname === "/demo-cases") return cases(["case-a", "case-b"]);
    return {
      ok: true, async json() {
        return {
          status: "allow", page: { offset: 0, limit: 200, total_visible: 1, has_more: false },
          entries: [{ id: "hit-" + url.searchParams.get("case_id"),
            kind: "observed", source_ref: "same-source" }]
        };
      }
    };
  });
  env.fields["hunt-run"].dispatchEvent({ type: "click" });
  await flush();
  const pivot = env.fields["hunt-results"].children[3];
  env.fields.token.value = "rotated-bearer";
  env.fields.token.dispatchEvent({ type: "input" });
  pivot.dispatchEvent({ type: "click" });
  assert.equal(env.submits, 0);
  assert.equal(env.fields["hunt-results"].children.length, 0);
});
