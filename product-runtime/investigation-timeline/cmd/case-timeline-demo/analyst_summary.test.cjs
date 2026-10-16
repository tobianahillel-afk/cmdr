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
  "return-to-case", "copy-return-link", "analyze", "analysis", "analysis-actions", "download-report",
  "compare-case", "compare", "comparison"
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
  const downloads = [];
  const LocalURL = class extends URL {};
  LocalURL.createObjectURL = blob => {
    downloads.push({ blob, clicked: false, revoked: false });
    return "blob:http://127.0.0.1:8765/demo-" + downloads.length;
  };
  LocalURL.revokeObjectURL = () => {
    if (downloads.length) downloads[downloads.length - 1].revoked = true;
  };
  const ctx = {
    window: { location: { origin: "http://127.0.0.1:8765", search: "" } },
    document: {
      getElementById: id => fields[id],
      createElement: tag => {
        const element = node();
        if (tag === "a") element.click = () => {
          if (downloads.length) downloads[downloads.length - 1].clicked = true;
        };
        return element;
      },
      createTextNode: text => ({ textContent: text })
    },
    navigator: { clipboard: { writeText: async () => {} } },
    Option: class { constructor(text, value) { this.text = text; this.value = value; } },
    Event: class { constructor(type) { this.type = type; } },
    URL: LocalURL, URLSearchParams, AbortController, Blob,
    setTimeout: callback => callback(),
    fetch: async (address, options) => {
      const url = new URL(String(address));
      requests.push({ url, options });
      return fetchMock(url, options, requests.length);
    }
  };
  vm.runInNewContext(source, ctx);
  fields.token.value = "ephemeral-token";
  fields["case-id"].value = "case-disk";
  return { fields, requests, downloads };
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
            source_ref: options.sourceFor ? options.sourceFor(index) : "source-" + (index % 3),
            source_time: new Date(Date.UTC(2026, 9, 1, 10, 0, index)).toISOString(),
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
  assert.match(env.fields.analysis.textContent, /Plage des horodatages certains \(UTC\)/);
  assert.match(env.fields.analysis.textContent, /Plus grand écart entre événements visibles certains : 2\.000 s/);
  assert.match(env.fields.analysis.textContent, /ne démontre jamais l'absence d'activité/);
  assert.match(env.fields.analysis.textContent, /Sources présentes : 3/);
  assert.match(env.fields.analysis.textContent, /source-0 : 69 événement\(s\)/);
  assert.match(env.fields.analysis.textContent, /source-1 : 68 événement\(s\)/);
  assert.match(env.fields.analysis.textContent, /source-2 : 68 événement\(s\)/);
  assert.match(env.fields.analysis.textContent, /pas une preuve de causalité/);
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

test("source grouping uses only the selected authorized pages, not denied data", async () => {
  const denied = "DO-NOT-SHOW-CROSS-TENANT";
  const env = boot(url => page(Number(url.searchParams.get("offset")), 203, {
    sourceFor: index => index < 201 ? "owner-a" : "owner-b"
  }));
  env.fields.analyze.handlers.click();
  await flush();
  await flush();
  assert.equal(env.requests.length, 2);
  assert.match(env.fields.analysis.textContent, /Sources présentes : 2/);
  assert.match(env.fields.analysis.textContent, /owner-a : 201 événement\(s\)/);
  assert.match(env.fields.analysis.textContent, /owner-b : 2 événement\(s\)/);
  assert.doesNotMatch(env.fields.analysis.textContent, new RegExp(denied));
  assert.ok(env.requests.every(request => request.options.headers.Authorization === "Bearer ephemeral-token"));
});

test("a malformed source reference fails closed without partial correlation output", async () => {
  const env = boot(url => page(Number(url.searchParams.get("offset")), 201, {
    sourceFor: index => index === 200 ? "" : "owner-visible"
  }));
  env.fields.analyze.handlers.click();
  await flush();
  await flush();
  assert.match(env.fields.analysis.textContent, /Aucun total partiel conservé/);
  assert.doesNotMatch(env.fields.analysis.textContent, /owner-visible/);
});

test("source grouping renders authorized owner labels as text and caps unusually long labels", async () => {
  const label = "<img onerror=alert(1)>" + "x".repeat(200);
  const env = boot(url => page(Number(url.searchParams.get("offset")), 2, {
    sourceFor: () => label
  }));
  env.fields.analyze.handlers.click();
  await flush();
  assert.match(env.fields.analysis.textContent, /<img onerror=alert\(1\)>/);
  assert.match(env.fields.analysis.textContent, /abrégée/);
  assert.ok(env.fields.analysis.textContent.length < 1100);
});

test("time triage ignores uncertain instants and reports no gaps when only one is certain", async () => {
  const env = boot(url => {
    const response = page(Number(url.searchParams.get("offset")), 3);
    const original = response.json;
    response.json = async () => {
      const data = await original();
      data.entries[1].time_uncertain = true;
      data.entries[2].time_uncertain = true;
      data.entries[2].source_time = "source clock unavailable";
      return data;
    };
    return response;
  });
  env.fields.analyze.handlers.click();
  await flush();
  assert.match(env.fields.analysis.textContent, /horodatage incertain : 2/);
  assert.match(env.fields.analysis.textContent, /Écarts temporels non calculables/);
  assert.doesNotMatch(env.fields.analysis.textContent, /Plus grand écart entre/);
});

test("malformed declared-certain timestamp never publishes a partial forensic analysis", async () => {
  const env = boot(url => {
    const response = page(Number(url.searchParams.get("offset")), 205);
    const original = response.json;
    response.json = async () => {
      const data = await original();
      if (Number(url.searchParams.get("offset")) === 200) {
        data.entries[0].source_time = "not-a-timestamp";
        data.entries[0].time_uncertain = false;
      }
      return data;
    };
    return response;
  });
  env.fields.analyze.handlers.click();
  await flush();
  await flush();
  assert.equal(env.requests.length, 2);
  assert.match(env.fields.analysis.textContent, /Aucun total partiel conservé/);
  assert.doesNotMatch(env.fields.analysis.textContent, /Plage des horodatages/);
});

test("timestamp offsets are normalized for chronological gaps across page boundaries", async () => {
  const env = boot(url => {
    const response = page(Number(url.searchParams.get("offset")), 201);
    const original = response.json;
    response.json = async () => {
      const data = await original();
      for (const item of data.entries) {
        const utc = new Date(item.source_time);
        item.source_time = new Date(utc.getTime() + 3600000).toISOString().replace("Z", "+01:00");
        item.time_uncertain = false;
      }
      return data;
    };
    return response;
  });
  env.fields.analyze.handlers.click();
  await flush();
  await flush();
  assert.match(env.fields.analysis.textContent, /Plage des horodatages certains \(UTC\)/);
  assert.match(env.fields.analysis.textContent, /2026-10-01T10:00:00\.000Z/);
  assert.match(env.fields.analysis.textContent, /2026-10-01T10:03:20\.000Z/);
  assert.match(env.fields.analysis.textContent, /Plus grand écart entre événements visibles certains : 1\.000 s/);
});

test("source group drilldown reauthorizes a representative Entry with existing owner routes", async () => {
  const env = boot((url) => {
    const offset = Number(url.searchParams.get("offset"));
    if (url.pathname === "/case-timeline") {
      return page(offset, 205, { sourceFor: index => index < 201 ? "owner-a" : "owner-b" });
    }
    if (url.pathname === "/case-source-navigation") {
      const query = new URLSearchParams(url.searchParams);
      query.delete("entry_id");
      return { ok: true, async json() { return {
        status: "ready", kind: "telemetry-event", object_id: "authorized-source-object",
        return_is_authority: false, return_case_id: "case-disk",
        return_entry_id: url.searchParams.get("entry_id"),
        return_query: query.toString()
      }; } };
    }
    if (url.pathname === "/demo-source-open") {
      return { ok: true, async json() { return {
        status: "ready", synthetic: true, kind: "telemetry-event",
        case_id: "case-disk", entry_id: url.searchParams.get("entry_id"),
        object_id: "authorized-source-object", preview: "Synthèse source visible"
      }; } };
    }
    throw Error("Unexpected endpoint " + url.pathname);
  });
  env.fields.analyze.handlers.click();
  await flush();
  await flush();
  assert.equal(env.requests.length, 2);
  const actions = env.fields["analysis-actions"].children;
  assert.equal(actions.length, 2);
  assert.match(actions[0].textContent, /owner-a/);
  assert.match(actions[1].textContent, /owner-b/);
  actions[1].handlers.click();
  await flush();
  await flush();
  assert.equal(env.requests.length, 4);
  assert.deepEqual(env.requests.slice(2).map(request => request.url.pathname),
    ["/case-source-navigation", "/demo-source-open"]);
  assert.equal(env.requests[2].url.searchParams.get("entry_id"), "event-201");
  assert.equal(env.requests[3].url.searchParams.get("entry_id"), "event-201");
  for (const request of env.requests.slice(2)) {
    assert.equal(request.options.headers.Authorization, "Bearer ephemeral-token");
    assert.equal(request.options.credentials, "omit");
    assert.equal(request.options.cache, "no-store");
    assert.equal(request.url.searchParams.has("source_ref"), false);
    assert.equal(request.url.href.includes("ephemeral-token"), false);
  }
  assert.equal(env.fields["source-detail"].hidden, false);
  env.fields.token.value = "revoked";
  env.fields.token.handlers.input();
  assert.equal(env.fields["analysis-actions"].children.length, 0);
  assert.equal(env.fields["source-detail"].hidden, true);
});

test("source group action never exposes a denied destination or partial summary", async () => {
  let deny = false;
  const env = boot(url => {
    if (url.pathname === "/case-timeline") return page(Number(url.searchParams.get("offset")), 205);
    if (url.pathname === "/case-source-navigation") return { ok: !deny, async json() {
      return { status: "ready" };
    } };
    throw Error("Denied source should not reach owner preview");
  });
  env.fields.analyze.handlers.click();
  await flush();
  await flush();
  assert.equal(env.fields["analysis-actions"].children.length, 3);
  deny = true;
  env.fields["analysis-actions"].children[0].handlers.click();
  await flush();
  assert.equal(env.fields["source-detail"].hidden, true);
  assert.equal(env.fields["copy-return-link"].disabled, true);
  assert.match(env.fields["source-status"].textContent, /inaccessible|révoquée/);
  // A failed later analysis MUST NOT leave a source action from a prior result.
  env.fields.analyze.handlers.click();
  await flush();
  await flush();
  assert.equal(env.fields["analysis-actions"].children.length, 3);
  env.fields["case-id"].value = "other-case";
  env.fields["case-id"].handlers.input();
  assert.equal(env.fields["analysis-actions"].children.length, 0);
});


test("explicit authorized JSON report is bounded to summary and contains no bearer or raw rows", async () => {
  const env = boot(url => page(Number(url.searchParams.get("offset")), 205));
  assert.equal(env.fields["download-report"].disabled, true);
  env.fields.analyze.handlers.click();
  await flush();
  await flush();
  assert.equal(env.fields["download-report"].disabled, false);
  assert.equal(env.downloads.length, 0);
  env.fields["download-report"].handlers.click();
  assert.equal(env.downloads.length, 1);
  assert.equal(env.downloads[0].clicked, true);
  assert.equal(env.downloads[0].revoked, true);
  const json = await env.downloads[0].blob.text();
  const report = JSON.parse(json);
  assert.equal(report.schema_version, 1);
  assert.equal(report.type, "cmdr-synthetic-authorized-case-summary");
  assert.equal(report.case_id, "case-disk");
  assert.equal(report.count, 205);
  assert.equal(report.event_kinds.observed, 103);
  assert.equal(report.event_kinds.inferred, 102);
  assert.equal(report.uncertain_timestamps, 1);
  assert.equal(report.unique_sources, 3);
  assert.equal(report.top_sources.length, 3);
  assert.equal(report.largest_visible_certain_gap_ms, 2000);
  assert.equal(report.top_sources[0].count, 69);
  assert.equal(report.filters.q, "");
  assert.ok(report.limitations.some(x => x.includes("not transactionally isolated")));
  assert.equal(json.includes("ephemeral-token"), false);
  assert.equal(json.includes("event-204"), false);
  assert.equal(json.includes("Bearer"), false);
  assert.equal(json.includes("http://127.0.0.1"), false);
});

test("JSON export is disabled for incomplete, revoked or changed analysis", async () => {
  const env = boot(url => page(Number(url.searchParams.get("offset")), 201));
  env.fields.analyze.handlers.click();
  await flush();
  await flush();
  assert.equal(env.fields["download-report"].disabled, false);
  env.fields.token.value = "revoked";
  env.fields.token.handlers.input();
  assert.equal(env.fields["download-report"].disabled, true);
  env.fields["download-report"].handlers.click();
  assert.equal(env.downloads.length, 0);
  env.fields.analyze.handlers.click();
  await flush();
  await flush();
  assert.equal(env.fields["download-report"].disabled, false);
  env.fields["case-id"].value = "different-case";
  env.fields["case-id"].handlers.input();
  assert.equal(env.fields["download-report"].disabled, true);
  env.fields["download-report"].handlers.click();
  assert.equal(env.downloads.length, 0);

  const bad = boot(url => {
    const offset = Number(url.searchParams.get("offset"));
    return page(offset, 201, offset === 200 ? {total: 200} : {});
  });
  bad.fields.analyze.handlers.click();
  await flush();
  await flush();
  assert.equal(bad.fields["download-report"].disabled, true);
  bad.fields["download-report"].handlers.click();
  assert.equal(bad.downloads.length, 0);
});

test("two authorized Cases yield exact shared-source counts across HTTP page boundaries", async () => {
  const env = boot(url => {
    const at = Number(url.searchParams.get("offset"));
    const id = url.searchParams.get("case_id");
    assert.equal(url.pathname, "/case-timeline");
    if (id === "case-disk") return page(at, 205, {
      sourceFor: n => n < 201 ? "shared-ref" : "only-left"
    });
    if (id === "case-peer") return page(at, 3, {
      sourceFor: n => n < 2 ? "shared-ref" : "only-right"
    });
    throw Error("unexpected Case");
  });
  env.fields["compare-case"].value = "case-peer";
  env.fields.compare.handlers.click();
  for (let i = 0; i < 6; i++) await flush();
  assert.deepEqual(env.requests.map(r => [r.url.searchParams.get("case_id"),r.url.searchParams.get("offset")]),
    [["case-disk","0"],["case-disk","200"],["case-peer","0"]]);
  assert.match(env.fields.comparison.textContent, /205 \/ 3/);
  assert.match(env.fields.comparison.textContent, /observés : 103 \/ 2/);
  assert.match(env.fields.comparison.textContent, /inférés : 102 \/ 1/);
  assert.match(env.fields.comparison.textContent, /horodatages incertains : 1 \/ 1/);
  assert.match(env.fields.comparison.textContent, /plage certaine Case 1 \(UTC\)/);
  assert.match(env.fields.comparison.textContent, /plage certaine Case 2 \(UTC\)/);
  assert.match(env.fields.comparison.textContent, /chevauchement des plages certaines : oui/);
  assert.match(env.fields.comparison.textContent, /références communes : 1/);
  assert.match(env.fields.comparison.textContent, /shared-ref : 201 \/ 2/);
  assert.doesNotMatch(env.fields.comparison.textContent, /only-left|only-right/);
  assert.ok(env.requests.every(r => r.options.headers.Authorization === "Bearer ephemeral-token" &&
    r.options.credentials === "omit" && r.options.cache === "no-store" &&
    !r.url.href.includes("ephemeral-token")));
});

test("cross-tenant Case denial fails closed with no partial cross-Case correlation", async () => {
  const env = boot(url => url.searchParams.get("case_id") === "case-peer"
    ? { ok: false, async json() { throw Error("must not decode denied Case"); } }
    : page(Number(url.searchParams.get("offset")), 2));
  env.fields["compare-case"].value = "case-peer";
  env.fields.compare.handlers.click();
  for (let i = 0; i < 5; i++) await flush();
  assert.equal(env.requests.length, 2);
  assert.match(env.fields.comparison.textContent, /Aucun résultat partiel/);
  assert.doesNotMatch(env.fields.comparison.textContent, /source-0/);
});

test("Case or token change invalidates an in-flight cross-Case comparison", async () => {
  let release;
  const pending = new Promise(resolve => { release = resolve; });
  const env = boot(url => url.searchParams.get("case_id") === "case-peer"
    ? pending : page(Number(url.searchParams.get("offset")), 1));
  env.fields["compare-case"].value = "case-peer";
  env.fields.compare.handlers.click();
  await flush();
  await flush();
  assert.equal(env.requests.length, 2);
  env.fields.token.value = "revoked-bearer";
  env.fields.token.handlers.input();
  release(page(0, 1));
  await flush();
  await flush();
  assert.equal(env.requests[1].options.signal.aborted, true);
  assert.match(env.fields.comparison.textContent, /Aucune comparaison calculée/);
  assert.equal(env.fields.compare.disabled, false);
});

test("second Case cannot reuse same Case or absent bearer as a comparison", async () => {
  const env = boot(() => { throw Error("HTTP must not be called"); });
  env.fields["compare-case"].value = "case-disk";
  env.fields.compare.handlers.click();
  await flush();
  assert.equal(env.requests.length, 0);
  env.fields["compare-case"].value = "case-peer";
  env.fields.token.value = "";
  env.fields.compare.handlers.click();
  await flush();
  assert.equal(env.requests.length, 0);
  assert.match(env.fields.comparison.textContent, /Choisissez deux Cases/);
});

test("two independently authorized Cases with disjoint sources still produce useful temporal triage", async () => {
  const env = boot(url => {
    const at = Number(url.searchParams.get("offset"));
    const id = url.searchParams.get("case_id");
    const response = page(at, id === "case-peer" ? 2 : 3, {
      sourceFor: () => id === "case-peer" ? "peer-only" : "main-only"
    });
    const original = response.json;
    response.json = async () => {
      const data = await original();
      if (id === "case-peer") {
        // No certain timestamps for this Case: the overlap is unknowable.
        for (const e of data.entries) e.time_uncertain = true;
      }
      return data;
    };
    return response;
  });
  env.fields["compare-case"].value = "case-peer";
  env.fields.compare.handlers.click();
  for (let i = 0; i < 5; i++) await flush();
  assert.match(env.fields.comparison.textContent, /observés : 2 \/ 1/);
  assert.match(env.fields.comparison.textContent, /inférés : 1 \/ 1/);
  assert.match(env.fields.comparison.textContent, /références communes : 0/);
  assert.match(env.fields.comparison.textContent, /plage certaine Case 2 \(UTC\) : aucune plage certaine/);
  assert.match(env.fields.comparison.textContent, /chevauchement des plages certaines : indéterminable/);
  assert.doesNotMatch(env.fields.comparison.textContent, /peer-only|main-only/);
});

test("malformed late certain timestamp discards entire cross-Case triage", async () => {
  const env = boot(url => {
    const id = url.searchParams.get("case_id");
    const result = page(Number(url.searchParams.get("offset")), id === "case-peer" ? 1 : 2);
    const original = result.json;
    result.json = async () => {
      const data = await original();
      if (id === "case-peer") data.entries[0].source_time = "not-a-timestamp";
      return data;
    };
    return result;
  });
  env.fields["compare-case"].value = "case-peer";
  env.fields.compare.handlers.click();
  for (let i = 0; i < 5; i++) await flush();
  assert.match(env.fields.comparison.textContent, /Aucun résultat partiel/);
  assert.doesNotMatch(env.fields.comparison.textContent, /observés :|plage certaine Case 1/);
});
