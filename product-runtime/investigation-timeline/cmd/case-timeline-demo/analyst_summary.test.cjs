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
  "compare-case", "compare", "comparison", "comparison-actions", "download-comparison",
  "download-cross-events", "proximity", "proximity-seconds", "proximity-results",
  "trace-actions", "trace-results", "wazuh-cohorts", "wazuh-status", "wazuh-actions", "wazuh-events",
  "comparison-trace-actions", "comparison-trace-results"
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
            provenance: "synthetic-owner:" + index,
            timezone: "+00:00", version: 1, supersedes: "",
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

test("temporal pivot selects a five-minute UTC intersection and reauthorizes the target Case", async () => {
  const env = boot(url => {
    const id = url.searchParams.get("case_id");
    const at = Number(url.searchParams.get("offset"));
    const response = page(at, id === "case-peer" ? 2 : 3, {
      sourceFor: () => id === "case-peer" ? "peer-only" : "main-only"
    });
    if (url.searchParams.get("limit") === "20") {
      const original = response.json;
      response.json = async () => {
        const data = await original();
        data.page.limit = 20;
        return data;
      };
    }
    return response;
  });
  env.fields["compare-case"].value = "case-peer";
  env.fields.compare.handlers.click();
  for (let i = 0; i < 6; i++) await flush();
  assert.match(env.fields.comparison.textContent, /références communes : 0/);
  assert.match(env.fields.comparison.textContent, /fenêtres UTC de 5 min .* : 1/);
  const actions = env.fields["comparison-actions"].children;
  assert.equal(actions.length, 2);
  assert.match(actions[1].textContent, /Case 2/);
  assert.doesNotMatch(actions[1].textContent, /ephemeral-token|peer-only/);
  actions[1].handlers.click();
  await flush();
  assert.equal(env.fields["comparison-actions"].children.length, 0);
  assert.equal(env.fields["case-id"].value, "case-peer");
  const req = env.requests.at(-1);
  assert.equal(req.url.pathname, "/case-timeline");
  assert.equal(req.url.searchParams.get("case_id"), "case-peer");
  assert.equal(req.url.searchParams.get("limit"), "20");
  assert.equal(req.url.searchParams.get("offset"), "0");
  assert.equal(req.url.searchParams.get("from"), "2026-10-01T10:00:00.000Z");
  assert.equal(req.url.searchParams.get("to"), "2026-10-01T10:05:00.000Z");
  assert.equal(req.options.headers.Authorization, "Bearer ephemeral-token");
  assert.equal(req.options.credentials, "omit");
  assert.equal(req.options.cache, "no-store");
  assert.equal(req.url.href.includes("ephemeral-token"), false);
});

test("uncertain-only second Case has no temporal pivots; revoked token invalidates old buttons", async () => {
  const env = boot(url => {
    const id = url.searchParams.get("case_id");
    const response = page(Number(url.searchParams.get("offset")), 2);
    const original = response.json;
    response.json = async () => {
      const data = await original();
      if (id === "case-peer") for (const entry of data.entries) entry.time_uncertain = true;
      return data;
    };
    return response;
  });
  env.fields["compare-case"].value = "case-peer";
  env.fields.compare.handlers.click();
  for (let i = 0; i < 6; i++) await flush();
  assert.match(env.fields.comparison.textContent, /fenêtres UTC de 5 min .* : 0/);
  // Uncertain clocks prohibit temporal claims, but an already-authorized
  // common source is still actionable through separately checked owner ACLs.
  assert.equal(env.fields["comparison-actions"].children.length, 4);
  assert.equal(env.fields["comparison-actions"].children.filter(button =>
    /^Examiner Case/.test(button.textContent)).length, 0);
  const allowed = boot(url => page(Number(url.searchParams.get("offset")), 2));
  allowed.fields["compare-case"].value = "case-peer";
  allowed.fields.compare.handlers.click();
  for (let i = 0; i < 6; i++) await flush();
  const button = allowed.fields["comparison-actions"].children[0];
  assert.ok(button);
  allowed.fields.token.value = "revoked";
  allowed.fields.token.handlers.input();
  assert.equal(allowed.fields["comparison-actions"].children.length, 0);
  const before = allowed.requests.length;
  button.handlers.click();
  await flush();
  assert.equal(allowed.requests.length, before);
});

test("denied Case and malformed certain timestamp cannot create temporal pivot actions", async () => {
  for (const denied of [true, false]) {
    const env = boot(url => {
      const id = url.searchParams.get("case_id");
      if (denied && id === "case-peer") return { ok: false };
      const response = page(Number(url.searchParams.get("offset")), 2);
      const original = response.json;
      response.json = async () => {
        const data = await original();
        if (!denied && id === "case-peer") data.entries[0].source_time = "invalid";
        return data;
      };
      return response;
    });
    env.fields["compare-case"].value = "case-peer";
    env.fields.compare.handlers.click();
    for (let i = 0; i < 6; i++) await flush();
    assert.match(env.fields.comparison.textContent, /Aucun résultat partiel/);
    assert.equal(env.fields["comparison-actions"].children.length, 0);
  }
});

test("source sequence pages all authorized entries and reauthorizes an opened row", async () => {
  const env = boot(url => {
    if (url.pathname === "/case-timeline") return page(Number(url.searchParams.get("offset")), 65, {
      sourceFor: () => "owner-same"
    });
    if (url.pathname === "/case-source-navigation") {
      const query = new URLSearchParams(url.searchParams);
      query.delete("entry_id");
      return { ok: true, async json() { return {
        status: "ready", kind: "telemetry-event", object_id: "local-source",
        return_is_authority: false, return_case_id: "case-disk",
        return_entry_id: url.searchParams.get("entry_id"),
        return_query: query.toString()
      }; } };
    }
    if (url.pathname === "/demo-source-open") return { ok: true, async json() { return {
      status: "ready", synthetic: true, case_id: "case-disk",
      entry_id: url.searchParams.get("entry_id"),
      kind: "telemetry-event", object_id: "local-source", preview: "Authorized local source"
    }; } };
    throw Error("Unexpected source route " + url.pathname);
  });
  env.fields.analyze.handlers.click();
  await flush();
  assert.equal(env.requests.length, 1, "the complete authorized source was read once");
  assert.equal(env.fields["trace-actions"].children.length, 1);
  env.fields["trace-actions"].children[0].handlers.click();
  const result = env.fields["trace-results"];
  assert.match(result.textContent, /1–25 \/ 65/);
  assert.equal(result.children[0].children.length, 25);
  assert.match(result.children[0].children[0].children[0].textContent, /event-0 · observed/);
  result.children[1].handlers.click();
  assert.match(result.textContent, /26–50 \/ 65/);
  result.children[2].handlers.click();
  assert.match(result.textContent, /51–65 \/ 65/);
  assert.equal(result.children[0].children.length, 15);
  assert.match(result.children[0].children[0].children[0].textContent, /event-50/);
  const open = result.children[0].children[0].children[0];
  open.handlers.click();
  await flush();
  await flush();
  assert.deepEqual(env.requests.slice(1).map(r => r.url.pathname),
    ["/case-source-navigation", "/demo-source-open"]);
  for (const req of env.requests.slice(1)) {
    assert.equal(req.url.searchParams.get("entry_id"), "event-50");
    assert.equal(req.options.headers.Authorization, "Bearer ephemeral-token");
    assert.equal(req.options.cache, "no-store");
    assert.equal(req.url.searchParams.has("source_ref"), false);
  }
});

test("sequence actions become inert when bearer or Case changes", async () => {
  const env = boot(url => page(Number(url.searchParams.get("offset")), 28, {sourceFor: () => "owner-a"}));
  env.fields.analyze.handlers.click();
  await flush();
  const trace = env.fields["trace-actions"].children[0];
  trace.handlers.click();
  const oldRow = env.fields["trace-results"].children[0].children[0].children[0];
  assert.equal(env.fields["trace-results"].children[0].children.length, 25);
  env.fields.token.value = "revoked";
  env.fields.token.handlers.input();
  assert.equal(env.fields["trace-actions"].children.length, 0);
  assert.equal(env.fields["trace-results"].children.length, 0);
  assert.match(env.fields["trace-results"].textContent, /Aucune séquence/);
  trace.handlers.click();
  oldRow.handlers.click();
  assert.equal(env.requests.length, 1, "stale source selectors must not make a network request");
});

test("partial analysis never publishes source sequence controls", async () => {
  const env = boot(url => {
    const offset = Number(url.searchParams.get("offset"));
    return page(offset, 201, offset === 200 ? {total: 200} : {});
  });
  env.fields.analyze.handlers.click();
  await flush();
  await flush();
  assert.match(env.fields.analysis.textContent, /Aucun total partiel/);
  assert.equal(env.fields["trace-actions"].children.length, 0);
  assert.equal(env.fields["trace-results"].children.length, 0);
});

test("common-source pivot opens a representative event in peer Case via both authorization routes", async () => {
  const env = boot(url => {
    const id = url.searchParams.get("case_id");
    if (url.pathname === "/case-timeline") {
      return page(Number(url.searchParams.get("offset")), 1, { sourceFor: () => "shared-ref" });
    }
    if (url.pathname === "/case-source-navigation") {
      const selection = new URLSearchParams(url.searchParams);
      selection.delete("entry_id");
      return { ok: true, async json() { return {
        status: "ready", kind: "telemetry-event", object_id: "owner-" + id,
        return_is_authority: false, return_case_id: id,
        return_entry_id: url.searchParams.get("entry_id"),
        return_query: selection.toString()
      }; } };
    }
    if (url.pathname === "/demo-source-open") {
      return { ok: true, async json() { return {
        status: "ready", synthetic: true, case_id: id,
        entry_id: url.searchParams.get("entry_id"), kind: "telemetry-event",
        object_id: "owner-" + id, preview: "permitted synthetic owner"
      }; } };
    }
    throw Error("unexpected endpoint " + url.pathname);
  });
  env.fields["compare-case"].value = "case-peer";
  env.fields.compare.handlers.click();
  for (let i = 0; i < 6; i++) await flush();
  const actions = env.fields["comparison-actions"].children;
  assert.equal(actions.length, 4); // two common-source and two time-bucket pivots
  assert.match(actions[1].textContent, /source commune.*Case 2/);
  actions[1].handlers.click();
  for (let i = 0; i < 5; i++) await flush();
  assert.equal(env.fields["case-id"].value, "case-peer");
  assert.equal(env.fields["comparison-actions"].children.length, 0);
  assert.deepEqual(env.requests.slice(2).map(r => r.url.pathname),
    ["/case-source-navigation", "/demo-source-open"]);
  for (const req of env.requests.slice(2)) {
    assert.equal(req.url.searchParams.get("case_id"), "case-peer");
    assert.equal(req.url.searchParams.get("entry_id"), "event-0");
    assert.equal(req.options.headers.Authorization, "Bearer ephemeral-token");
    assert.equal(req.options.cache, "no-store");
    assert.equal(req.options.credentials, "omit");
    assert.equal(req.url.searchParams.has("source_ref"), false);
    assert.ok(!req.url.href.includes("ephemeral-token"));
  }
  assert.equal(env.fields["source-detail"].hidden, false);
});

test("shared-source pivot fails closed if owner revokes access", async () => {
  const env = boot(url => {
    if (url.pathname === "/case-timeline") return page(Number(url.searchParams.get("offset")), 1,
      { sourceFor: () => "<img onerror=alert(1)>".repeat(10) });
    if (url.pathname === "/case-source-navigation") return { ok: false };
    throw Error("denied destination owner preview must not be fetched");
  });
  env.fields["compare-case"].value = "case-peer";
  env.fields.compare.handlers.click();
  for (let i = 0; i < 6; i++) await flush();
  const action = env.fields["comparison-actions"].children[1];
  assert.match(action.textContent, /abrégée/);
  action.handlers.click();
  for (let i = 0; i < 5; i++) await flush();
  assert.equal(env.requests.length, 3);
  assert.equal(env.fields["source-detail"].hidden, true);
  assert.match(env.fields["source-status"].textContent, /autorisation révoquée/);
});

test("stale common-source action is inert after revoking bearer", async () => {
  const env = boot(url => page(Number(url.searchParams.get("offset")), 1,
    { sourceFor: () => "shared-source" }));
  env.fields["compare-case"].value = "case-peer";
  env.fields.compare.handlers.click();
  for (let i = 0; i < 6; i++) await flush();
  const action = env.fields["comparison-actions"].children[1];
  const before = env.requests.length;
  env.fields.token.value = "revoked";
  env.fields.token.handlers.input();
  action.handlers.click();
  await flush();
  assert.equal(env.requests.length, before);
  assert.equal(env.fields["comparison-actions"].children.length, 0);
});

test("complete authenticated two-Case comparison exports a bounded local JSON investigation lead", async () => {
  const env = boot(url => {
    const id = url.searchParams.get("case_id");
    const total = id === "case-disk" ? 205 : 3;
    return page(Number(url.searchParams.get("offset")), total,
      { sourceFor: index => index % 2 ? "source-common" : "source-local" });
  });
  env.fields["compare-case"].value = "case-peer";
  assert.equal(env.fields["download-comparison"].disabled, true);
  env.fields.compare.handlers.click();
  for (let i = 0; i < 7; i++) await flush();
  assert.equal(env.requests.length, 3); // 2 pages in first Case, 1 in second
  assert.equal(env.fields["download-comparison"].disabled, false);
  assert.equal(env.downloads.length, 0, "only explicit user click may save");
  env.fields["download-comparison"].handlers.click();
  assert.equal(env.downloads.length, 1);
  assert.equal(env.downloads[0].clicked, true);
  assert.equal(env.downloads[0].revoked, true);
  const json = await env.downloads[0].blob.text();
  const receipt = JSON.parse(json);
  assert.equal(receipt.schema_version, 1);
  assert.equal(receipt.type, "cmdr-synthetic-authorized-cross-case-comparison");
  assert.equal(receipt.case_ids[0], "case-disk");
  assert.equal(receipt.case_ids[1], "case-peer");
  assert.equal(receipt.cases[0].visible_events, 205);
  assert.equal(receipt.cases[1].visible_events, 3);
  assert.equal(receipt.cases[0].observed, 103);
  assert.equal(receipt.cases[0].inferred, 102);
  assert.equal(receipt.common_reference_count, 2);
  assert.equal(receipt.top_common_sources.length, 2);
  assert.ok(receipt.common_five_minute_windows >= 1);
  assert.ok(receipt.top_common_windows_utc.length >= 1);
  assert.equal(receipt.certain_range_overlap, true);
  assert.ok(receipt.limitations.some(value => value.includes("not causality")));
  for (const secret of ["ephemeral-token", "Bearer", "event-204", "http://127.0.0.1"]) {
    assert.equal(json.includes(secret), false, "export must not contain " + secret);
  }
  for (const req of env.requests) {
    assert.equal(req.options.headers.Authorization, "Bearer ephemeral-token");
    assert.equal(req.options.credentials, "omit");
    assert.equal(req.options.cache, "no-store");
    assert.ok(!req.url.href.includes("ephemeral-token"));
  }
});

test("cross-Case report never becomes available after denied second Case or stale selection", async () => {
  const denied = boot(url => url.searchParams.get("case_id") === "case-peer"
    ? { ok: false } : page(Number(url.searchParams.get("offset")), 1));
  denied.fields["compare-case"].value = "case-peer";
  denied.fields.compare.handlers.click();
  for (let i = 0; i < 6; i++) await flush();
  assert.equal(denied.fields["download-comparison"].disabled, true);
  denied.fields["download-comparison"].handlers.click();
  assert.equal(denied.downloads.length, 0);
  assert.match(denied.fields.comparison.textContent, /Aucun résultat partiel/);

  const env = boot(url => page(Number(url.searchParams.get("offset")), 1));
  env.fields["compare-case"].value = "case-peer";
  env.fields.compare.handlers.click();
  for (let i = 0; i < 6; i++) await flush();
  assert.equal(env.fields["download-comparison"].disabled, false);
  env.fields["compare-case"].value = "changed-peer";
  env.fields["compare-case"].handlers.input();
  assert.equal(env.fields["download-comparison"].disabled, true);
  env.fields["download-comparison"].handlers.click();
  assert.equal(env.downloads.length, 0);

  env.fields["compare-case"].value = "case-peer";
  env.fields.compare.handlers.click();
  for (let i = 0; i < 6; i++) await flush();
  assert.equal(env.fields["download-comparison"].disabled, false);
  env.fields.token.value = "revoked";
  env.fields.token.handlers.input();
  assert.equal(env.fields["download-comparison"].disabled, true);
  env.fields["download-comparison"].handlers.click();
  assert.equal(env.downloads.length, 0);
});

test("inconsistent authorized cross-Case page never leaves a downloadable receipt", async () => {
  const env = boot(url => {
    const id = url.searchParams.get("case_id");
    if (id === "case-peer") return { ok: false };
    const offset = Number(url.searchParams.get("offset"));
    return page(offset, 201, offset === 200 ? { total: 200 } : {});
  });
  env.fields["compare-case"].value = "case-peer";
  env.fields.compare.handlers.click();
  for (let i = 0; i < 6; i++) await flush();
  assert.equal(env.fields["download-comparison"].disabled, true);
  assert.equal(env.downloads.length, 0);
});

test("cross-Case source sequence paginates real authorized rows and rechecks peer source owner", async () => {
  const env = boot(url => {
    const caseName = url.searchParams.get("case_id");
    if (url.pathname === "/case-timeline") {
      return page(Number(url.searchParams.get("offset")), 27, { sourceFor: () => "shared-correlation" });
    }
    if (url.pathname === "/case-source-navigation") {
      const selection = new URLSearchParams(url.searchParams);
      selection.delete("entry_id");
      return { ok: true, async json() { return {
        status: "ready", kind: "telemetry-event", object_id: "owner-" + caseName,
        return_is_authority: false, return_case_id: caseName,
        return_entry_id: url.searchParams.get("entry_id"),
        return_query: selection.toString()
      }; } };
    }
    if (url.pathname === "/demo-source-open") return { ok: true, async json() { return {
      status: "ready", synthetic: true, case_id: caseName,
      entry_id: url.searchParams.get("entry_id"), kind: "telemetry-event",
      object_id: "owner-" + caseName, preview: "authorized synthetic owner"
    }; } };
    throw Error("unexpected source route " + url.pathname);
  });
  env.fields["compare-case"].value = "case-peer";
  env.fields.compare.handlers.click();
  for (let i = 0; i < 6; i++) await flush();
  assert.equal(env.requests.length, 2, "both complete Case page traversals are required");
  assert.equal(env.fields["comparison-trace-actions"].children.length, 1);
  env.fields["comparison-trace-actions"].children[0].handlers.click();
  const result = env.fields["comparison-trace-results"];
  assert.match(result.textContent, /1–25 \/ 54/);
  assert.equal(result.children[0].children.length, 25);
  assert.match(result.children[0].children[0].children[0].textContent, /Case 1 · observed/);
  assert.match(result.children[0].children[1].children[0].textContent, /Case 2 · observed/);
  assert.equal(env.requests.length, 2, "local trace must not invent an extra source/index endpoint");
  result.children.find(el => /suivants/.test(el.textContent)).handlers.click();
  assert.match(result.textContent, /26–50 \/ 54/);
  result.children.find(el => /suivants/.test(el.textContent)).handlers.click();
  assert.match(result.textContent, /51–54 \/ 54/);
  assert.equal(result.children[0].children.length, 4);
  assert.match(result.children[0].children.at(-1).children[0].textContent, /incertain/);
  // Return to first page and open the peer Case's authorized event.
  result.children.find(el => /précédents/.test(el.textContent)).handlers.click();
  result.children.find(el => /précédents/.test(el.textContent)).handlers.click();
  result.children[0].children[1].children[0].handlers.click();
  for (let i = 0; i < 6; i++) await flush();
  assert.equal(env.fields["case-id"].value, "case-peer");
  assert.equal(env.requests.length, 4, "the two first-party owner routes must be rechecked");
  assert.deepEqual(env.requests.slice(2).map(r => r.url.pathname),
    ["/case-source-navigation", "/demo-source-open"]);
  for (const req of env.requests.slice(2)) {
    assert.equal(req.url.searchParams.get("case_id"), "case-peer");
    assert.equal(req.url.searchParams.get("entry_id"), "event-0");
    assert.equal(req.url.searchParams.has("source_ref"), false);
    assert.equal(req.options.headers.Authorization, "Bearer ephemeral-token");
  }
  assert.equal(env.fields["comparison-trace-results"].children.length, 0);
  assert.equal(env.fields["source-detail"].hidden, false);
});

test("cross-Case source sequence becomes inert on bearer invalidation", async () => {
  const env = boot(url => page(Number(url.searchParams.get("offset")), 2,
    { sourceFor: () => "common-owner" }));
  env.fields["compare-case"].value = "case-peer";
  env.fields.compare.handlers.click();
  for (let i = 0; i < 6; i++) await flush();
  const oldAction = env.fields["comparison-trace-actions"].children[0];
  oldAction.handlers.click();
  const oldRow = env.fields["comparison-trace-results"].children[0].children[0].children[0];
  const requestsBefore = env.requests.length;
  env.fields.token.value = "revoked-token";
  env.fields.token.handlers.input();
  oldRow.handlers.click();
  oldAction.handlers.click();
  await flush();
  assert.equal(env.requests.length, requestsBefore, "stale local correlation cannot invoke owners");
  assert.equal(env.fields["comparison-trace-actions"].children.length, 0);
  assert.equal(env.fields["comparison-trace-results"].children.length, 0);
});

test("explicit two-Case sensitive chronology exports all authorized rows in known-first order", async () => {
  const env = boot(url => {
    const id = url.searchParams.get("case_id");
    const total = id === "case-disk" ? 205 : 3;
    return page(Number(url.searchParams.get("offset")), total);
  });
  env.fields["compare-case"].value = "case-peer";
  assert.equal(env.fields["download-cross-events"].disabled, true);
  env.fields.compare.handlers.click();
  for (let i = 0; i < 8; i++) await flush();
  assert.equal(env.requests.length, 3, "two authorized Case traversals required");
  assert.equal(env.fields["download-cross-events"].disabled, false);
  assert.equal(env.downloads.length, 0, "sensitive export must require an explicit click");
  env.fields["download-cross-events"].handlers.click();
  assert.equal(env.downloads.length, 1);
  assert.equal(env.downloads[0].clicked, true);
  assert.equal(env.downloads[0].revoked, true);
  const raw = await env.downloads[0].blob.text();
  const receipt = JSON.parse(raw);
  assert.equal(receipt.type, "cmdr-synthetic-authorized-cross-case-chronology");
  assert.deepEqual(receipt.case_ids, ["case-disk", "case-peer"]);
  assert.equal(receipt.events.length, 208);
  assert.equal(receipt.uncertain_entries, 2);
  assert.equal(receipt.certain_entries, 206);
  assert.equal(receipt.events.at(-1).time_uncertain, true);
  assert.equal(receipt.events[0].case_id, "case-disk");
  assert.equal(receipt.events[0].id, "event-0");
  assert.equal(receipt.events[0].provenance, "synthetic-owner:0");
  assert.equal(receipt.events[0].version, 1);
  assert.equal(receipt.events[0].source_ref, "source-0");
  assert.ok(receipt.limitations.some(x => x.includes("Sensitive")));
  for (const forbidden of ["ephemeral-token", "Bearer ephemeral-token", "http://127.0.0.1", "owner preview"]) {
    assert.equal(raw.includes(forbidden), false, "private credential/preview in export: " + forbidden);
  }
  for (const request of env.requests) {
    assert.equal(request.options.headers.Authorization, "Bearer ephemeral-token");
    assert.equal(request.options.credentials, "omit");
    assert.equal(request.options.cache, "no-store");
  }
});

test("sensitive chronology fails closed on denied Case, invalid metadata and changed credential", async () => {
  const denied = boot(url => url.searchParams.get("case_id") === "case-peer"
    ? { ok: false } : page(Number(url.searchParams.get("offset")), 2));
  denied.fields["compare-case"].value = "case-peer";
  denied.fields.compare.handlers.click();
  for (let i = 0; i < 6; i++) await flush();
  assert.equal(denied.fields["download-cross-events"].disabled, true);
  denied.fields["download-cross-events"].handlers.click();
  assert.equal(denied.downloads.length, 0);

  const invalid = boot(url => {
    const result = page(Number(url.searchParams.get("offset")), 2);
    const original = result.json;
    result.json = async () => {
      const data = await original();
      if (url.searchParams.get("case_id") === "case-peer") delete data.entries[0].provenance;
      return data;
    };
    return result;
  });
  invalid.fields["compare-case"].value = "case-peer";
  invalid.fields.compare.handlers.click();
  for (let i = 0; i < 6; i++) await flush();
  assert.equal(invalid.fields["download-comparison"].disabled, false, "aggregate remains available");
  assert.equal(invalid.fields["download-cross-events"].disabled, true, "incomplete event metadata cannot be exported");
  invalid.fields["download-cross-events"].handlers.click();
  assert.equal(invalid.downloads.length, 0);

  const authorized = boot(url => page(Number(url.searchParams.get("offset")), 2));
  authorized.fields["compare-case"].value = "case-peer";
  authorized.fields.compare.handlers.click();
  for (let i = 0; i < 6; i++) await flush();
  assert.equal(authorized.fields["download-cross-events"].disabled, false);
  authorized.fields.token.value = "revoked";
  authorized.fields.token.handlers.input();
  assert.equal(authorized.fields["download-cross-events"].disabled, true);
  authorized.fields["download-cross-events"].handlers.click();
  assert.equal(authorized.downloads.length, 0);
});

test("browser compares two authorized Cases and renders nearest certain-time leads", async () => {
  const env = boot(async url => {
    const result = await page(Number(url.searchParams.get("offset")), 4).json();
    if (url.searchParams.get("case_id") === "case-peer") {
      result.entries[0].source_time = "2026-10-01T10:00:02.250Z";
      result.entries[1].source_time = "2026-10-01T10:00:05.000Z";
      result.entries[3].source_time = "2026-10-01T10:00:50.000Z";
    }
    return { ok: true, async json() { return result; } };
  });
  env.fields["compare-case"].value = "case-peer";
  env.fields["proximity-seconds"].value = "3";
  assert.equal(env.fields.proximity.disabled, true);
  env.fields.compare.handlers.click();
  for (let i = 0; i < 8; i++) await flush();
  assert.equal(env.requests.length, 2);
  assert.equal(env.fields.proximity.disabled, false);
  env.fields.proximity.handlers.click();
  assert.match(env.fields["proximity-results"].textContent, /3 piste\(s\)/);
  assert.match(env.fields["proximity-results"].textContent, /sans causalité/);
  const list = env.fields["proximity-results"].children[0];
  assert.equal(list.children.length, 3);
  assert.match(list.children[0].textContent, /750 ms/);
  assert.match(list.children[1].textContent, /1250 ms/);
  assert.match(list.children[2].textContent, /2250 ms/);
  assert.equal(list.children[0].children.length, 2);
  assert.equal(env.requests.every(r => r.options.headers.Authorization === "Bearer ephemeral-token"), true);
  assert.equal(env.requests.some(r => r.url.pathname !== "/case-timeline"), false);
  assert.equal(JSON.stringify(env.fields["proximity-results"]).includes("ephemeral-token"), false);
  env.fields["proximity-seconds"].value = "0";
  env.fields["proximity-seconds"].handlers.input();
  env.fields.proximity.handlers.click();
  assert.match(env.fields["proximity-results"].textContent, /Fenêtre invalide/);
});

test("temporal leads cannot survive bearer changes or a denied second Case", async () => {
  const env = boot(url => page(Number(url.searchParams.get("offset")), 2));
  env.fields["compare-case"].value = "case-peer";
  env.fields["proximity-seconds"].value = "30";
  env.fields.compare.handlers.click();
  for (let i = 0; i < 6; i++) await flush();
  env.fields.proximity.handlers.click();
  assert.match(env.fields["proximity-results"].textContent, /piste\(s\)/);
  env.fields.token.value = "revoked";
  env.fields.token.handlers.input();
  assert.equal(env.fields.proximity.disabled, true);
  assert.equal(env.fields["proximity-results"].children.length, 0);
  assert.match(env.fields["proximity-results"].textContent, /Aucune piste/);

  const denied = boot(url => url.searchParams.get("case_id") === "case-peer" ?
    { ok: false } : page(Number(url.searchParams.get("offset")), 2));
  denied.fields["compare-case"].value = "case-peer";
  denied.fields.compare.handlers.click();
  for (let i = 0; i < 6; i++) await flush();
  assert.equal(denied.fields.proximity.disabled, true);
  assert.equal(denied.fields["proximity-results"].children.length, 0);
});

test("Wazuh authorized agent/rule cohorts appear only after bearer HTTP response", async () => {
  const env = boot(url => {
    assert.equal(url.pathname, "/demo-wazuh-cohorts");
    assert.equal(url.searchParams.get("case_id"), "case-disk");
    return { ok: true, async json() { return {
      format: "cmdr-local-wazuh4-case-triage-v1", synthetic: true,
      imported_entries: 3,
      top_agents: [{ agent_id: "001", count: 2, distinct_rule_ids: 1, highest_known_rule_level: 9 }],
      top_rules: [{ rule_id: "5710", count: 3, highest_known_rule_level: 9 }],
      multi_feed_signals: [{ agent_id: "001", rule_id: "5710", alert_count: 3,
        distinct_feeds: 2, highest_known_rule_level: 9 }]
    }; } };
  });
  env.fields["wazuh-cohorts"].handlers.click();
  await flush();
  await flush();
  assert.equal(env.requests.length, 1);
  assert.equal(env.requests[0].options.headers.Authorization, "Bearer ephemeral-token");
  assert.equal(env.requests[0].options.credentials, "omit");
  assert.equal(env.requests[0].options.cache, "no-store");
  assert.ok(!env.requests[0].url.href.includes("ephemeral-token"));
  assert.match(env.fields["wazuh-status"].textContent, /agent 001 : 2 alerte/);
  assert.match(env.fields["wazuh-status"].textContent, /règle 5710 : 3 alerte/);
  assert.match(env.fields["wazuh-status"].textContent, /agent 001 \/ règle 5710 : 3 alerte/);
  assert.match(env.fields["wazuh-status"].textContent, /2 flux distinct/);
});

test("Wazuh denied, malformed or stale replies never render a partial cohort", async () => {
  for (const answer of [
    { ok: false },
    { ok: true, async json() { return { format: "cmdr-local-wazuh4-case-triage-v1",
      synthetic: true, imported_entries: 1,
      top_agents: [{agent_id: "<img>", count: 1, distinct_rule_ids: 1}],
      top_rules: [] }; } }
  ]) {
    const env = boot(() => answer);
    env.fields["wazuh-cohorts"].handlers.click();
    await flush();
    await flush();
    assert.match(env.fields["wazuh-status"].textContent, /indisponible|refusé/);
    assert.doesNotMatch(env.fields["wazuh-status"].textContent, /<img>|agent 001/);
  }
  let resolve;
  const env = boot(() => new Promise(done => { resolve = done; }));
  env.fields["wazuh-cohorts"].handlers.click();
  env.fields.token.value = "changed";
  env.fields.token.handlers.input();
  resolve({ ok: true, async json() { return {
    format: "cmdr-local-wazuh4-case-triage-v1", synthetic: true,
    imported_entries: 1, top_agents: [{agent_id:"001", count:1, distinct_rule_ids:1}],
    top_rules: []
  }; } });
  await flush();
  await flush();
  assert.doesNotMatch(env.fields["wazuh-status"].textContent, /agent 001/);
});

test("Wazuh cohort drilldown opens a real authorized event through the existing source route", async () => {
  const requestsSeen = [];
  const env = boot((url, options) => {
    requestsSeen.push(url.pathname);
    if (url.pathname === "/demo-wazuh-cohorts") return { ok: true, async json() {
      return { format: "cmdr-local-wazuh4-case-triage-v1", synthetic: true,
        imported_entries: 2,
        top_agents: [{ agent_id: "001", count: 1, distinct_rule_ids: 1, highest_known_rule_level: 9 }],
        top_rules: [{ rule_id: "5710", count: 2, highest_known_rule_level: 9 }] };
    } };
    if (url.pathname === "/demo-wazuh-drilldown") return { ok: true, async json() {
      return { format: "cmdr-local-wazuh4-drilldown-v1", synthetic: true,
        selector: url.searchParams.has("agent_id") ? "agent_id" : "rule_id",
        value: url.searchParams.get("agent_id") || url.searchParams.get("rule_id"),
        total_matches: 1, entries: [{
          id: "local-jsonl-entry-abcdef", kind: "observed",
          source_time: "2026-10-01T10:00:00Z", time_uncertain: false
        }] };
    } };
    if (url.pathname === "/case-source-navigation") return { ok: true, async json() {
      return { status: "ready", kind: "telemetry-event", object_id: "authorized-event",
        return_is_authority: false, return_case_id: "case-disk",
        return_entry_id: "local-jsonl-entry-abcdef",
        return_query: "case_id=case-disk&limit=20&offset=0" };
    } };
    throw Error("Unexpected endpoint " + url.pathname);
  });
  env.fields["wazuh-cohorts"].handlers.click();
  await flush();
  await flush();
  assert.equal(env.fields["wazuh-actions"].children.length, 2);
  env.fields["wazuh-actions"].children[0].handlers.click();
  await flush();
  await flush();
  assert.deepEqual(requestsSeen.slice(0, 2), ["/demo-wazuh-cohorts", "/demo-wazuh-drilldown"]);
  assert.equal(env.fields["wazuh-events"].children.length, 1);
  const list = env.fields["wazuh-events"].children[0];
  assert.equal(list.children.length, 1);
  list.children[0].children[0].handlers.click();
  await flush();
  await flush();
  assert.ok(requestsSeen.includes("/case-source-navigation"), "drilldown did not reauthorize source");
  const drift = env.requests.find(r => r.url.pathname === "/demo-wazuh-drilldown");
  assert.equal(drift.url.searchParams.get("case_id"), "case-disk");
  assert.equal(drift.url.searchParams.get("agent_id"), "001");
  assert.equal(drift.options.headers.Authorization, "Bearer ephemeral-token");
  assert.ok(!drift.url.href.includes("ephemeral-token"));
});

test("Wazuh drilldown denial reveals no partial entry", async () => {
  const env = boot(url => {
    if (url.pathname === "/demo-wazuh-cohorts") return { ok: true, async json() {
      return { format: "cmdr-local-wazuh4-case-triage-v1", synthetic: true,
        imported_entries: 1,
        top_agents: [{ agent_id: "001", count: 1, distinct_rule_ids: 1 }],
        top_rules: [] };
    } };
    return { ok: false, async json() { throw Error("denied"); } };
  });
  env.fields["wazuh-cohorts"].handlers.click();
  await flush();
  await flush();
  env.fields["wazuh-actions"].children[0].handlers.click();
  await flush();
  await flush();
  assert.equal(env.fields["wazuh-events"].children.length, 0);
  assert.match(env.fields["wazuh-events"].textContent, /droits modifiés/);
});
