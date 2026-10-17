"use strict";

// Local, synthetic, read-only cross-Case text hunt. The browser never decides
// access; discovery and EVERY Timeline page use the existing server ACLs.
(() => {
  const token = document.getElementById("token");
  const caseID = document.getElementById("case-id");
  const search = document.getElementById("search");
  const kind = document.getElementById("kind");
  const from = document.getElementById("from");
  const to = document.getElementById("to");
  const form = document.getElementById("query");
  const query = document.getElementById("hunt-query");
  const button = document.getElementById("hunt-run");
  const status = document.getElementById("hunt-status");
  const results = document.getElementById("hunt-results");
  const maxCases = 8;
  const pageSize = 200;
  const maxEntries = 5000;
  let revision = 0;
  let active = null;

  function clear() {
    ++revision;
    if (active) active.abort();
    active = null;
    button.disabled = false;
    results.replaceChildren();
    status.textContent = "Aucune recherche transversale effectuée.";
  }

  token.addEventListener("input", clear);
  query.addEventListener("input", clear);

  async function authorizedGet(path, bearer, controller) {
    const response = await fetch(new URL(path, window.location.origin), {
      method: "GET", headers: { Authorization: "Bearer " + bearer },
      cache: "no-store", credentials: "omit", redirect: "error",
      mode: "same-origin", signal: controller.signal
    });
    if (!response.ok) throw new Error("case unavailable");
    return response.json();
  }

  // A UI hint only: every pivot executes the existing freshly authorized
  // Case-specific Timeline query. No browser-side source-access decisions.
  function pivotTo(id, text, selectedRevision, bearer, needle) {
    if (selectedRevision !== revision || token.value !== bearer ||
        query.value.trim() !== needle) return;
    // The cross-Case hunt searched every kind and every time. Keeping the
    // existing Case-specific filters would silently hide a proven text hit.
    for (const field of [kind, from, to]) {
      if (field.value) {
        field.value = "";
        field.dispatchEvent(new Event("input", { bubbles: true }));
      }
    }
    caseID.value = id;
    caseID.dispatchEvent(new Event("input", { bubbles: true }));
    search.value = text;
    search.dispatchEvent(new Event("input", { bubbles: true }));
    form.dispatchEvent(new Event("submit", { bubbles: true, cancelable: true }));
  }

  async function runHunt() {
    clear();
    const bearer = token.value;
    const needle = query.value.trim();
    if (!bearer || needle.length < 2 || needle.length > 96) {
      status.textContent = "Saisir un jeton et un indicateur de 2 à 96 caractères.";
      return;
    }
    const selectedRevision = revision;
    const controller = new AbortController();
    active = controller;
    button.disabled = true;
    status.textContent = "Recherche sur les Cases autorisés…";
    try {
      // Never trust a browser-maintained Case list. Discover afresh.
      const discovered = await authorizedGet("/demo-cases", bearer, controller);
      const cases = discovered && discovered.cases;
      if (!Array.isArray(cases) || cases.length > maxCases ||
          cases.some(id => typeof id !== "string" || !id || id.length > 1024) ||
          new Set(cases).size !== cases.length) {
        throw new Error("bounded Case discovery required");
      }
      const matches = [];
      for (const id of cases) {
        let total = null;
        const kinds = { observed: 0, inferred: 0 };
        const seen = new Set();
        for (let offset = 0; offset === 0 || offset < total; offset += pageSize) {
          if (controller.signal.aborted || selectedRevision !== revision)
            throw new Error("cancelled");
          const url = new URL("/case-timeline", window.location.origin);
          url.searchParams.set("case_id", id);
          url.searchParams.set("q", needle);
          url.searchParams.set("limit", String(pageSize));
          url.searchParams.set("offset", String(offset));
          const data = await authorizedGet(url, bearer, controller);
          if (!data || !data.page || !Array.isArray(data.entries) ||
              !["allow", "filtered", "partial", "uncertain", "empty"].includes(data.status) ||
              data.page.limit !== pageSize || data.page.offset !== offset ||
              !Number.isSafeInteger(data.page.total_visible) ||
              data.page.total_visible < 0 || data.page.total_visible > maxEntries ||
              (total !== null && total !== data.page.total_visible)) {
            throw new Error("unstable authorized Case page");
          }
          total = data.page.total_visible;
          const expectedCount = Math.min(pageSize, total - offset);
          if (data.entries.length !== expectedCount ||
              data.page.has_more !== (offset + expectedCount < total))
            throw new Error("incomplete authorized Case page");
          for (const entry of data.entries) {
            if (!entry || typeof entry.id !== "string" || !entry.id ||
                seen.has(entry.id) ||
                !["observed", "inferred", "user-action", "decision", "response"].includes(entry.kind) ||
                typeof entry.source_ref !== "string" || !entry.source_ref ||
                entry.source_ref.length > 1024) {
              throw new Error("invalid authorized entry");
            }
            seen.add(entry.id);
            if (entry.kind === "observed" || entry.kind === "inferred") kinds[entry.kind]++;
          }
        }
        if (total) matches.push({ caseID: id, total, kinds });

      }
      if (controller.signal.aborted || selectedRevision !== revision ||
          token.value !== bearer || query.value.trim() !== needle) return;
      // Commit display only after EVERY Case completed. No partial release.
      results.replaceChildren();
      for (const item of matches) {
        const pivot = document.createElement("button");
        pivot.type = "button";
        pivot.textContent = item.caseID + " — " + item.total +
          " correspondance(s) · observés : " + item.kinds.observed +
          " · inférés : " + item.kinds.inferred;
        pivot.addEventListener("click", () =>
          pivotTo(item.caseID, needle, selectedRevision, bearer, needle));
        results.append(pivot);
      }
      // Source references belong to one canonical Case in the Go owner.
      // A shared source-ID join is unreachable on valid snapshots. Report
      // only text cooccurrence across independently authorized Case results.
      const observed = matches.reduce((sum, item) => sum + item.kinds.observed, 0);
      const inferred = matches.reduce((sum, item) => sum + item.kinds.inferred, 0);
      status.textContent = matches.length
        ? matches.length + " Case(s) avec correspondance parmi " + cases.length +
          " Case(s) accessibles · observés : " + observed +
          " · inférés : " + inferred +
          ". Cooccurrence textuelle partielle, pas identité de source, IOC vérifié ni causalité."
        : "Aucune correspondance dans les Cases accessibles.";
    } catch {
      if (selectedRevision === revision) {
        results.replaceChildren();
        status.textContent = "Recherche incomplète ou accès révoqué : aucun résultat partiel affiché.";
      }
    } finally {
      if (selectedRevision === revision) {
        active = null;
        button.disabled = false;
      }
    }
  }

  button.addEventListener("click", runHunt);
})();
