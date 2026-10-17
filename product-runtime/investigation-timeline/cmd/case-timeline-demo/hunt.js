"use strict";

// Local, synthetic, read-only cross-Case text hunt. The browser never decides
// access; discovery and EVERY Timeline page use the existing server ACLs.
(() => {
  const token = document.getElementById("token");
  const caseID = document.getElementById("case-id");
  const search = document.getElementById("search");
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
      // Shared source references are derived ONLY from successful, already
      // authorized Timeline pages, never from host-private catalog metadata.
      const bySource = new Map();
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
            const source = bySource.get(entry.source_ref) ||
              { cases: new Set(), entries: 0 };
            source.cases.add(id);
            source.entries++;
            bySource.set(entry.source_ref, source);
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
      const common = [...bySource].filter(([, source]) => source.cases.size > 1)
        .sort((a, b) => b[1].cases.size - a[1].cases.size ||
          b[1].entries - a[1].entries || a[0].localeCompare(b[0])).slice(0, 5);
      if (common.length) {
        const caption = document.createElement("p");
        caption.textContent = "Sources communes aux Cases autorisés — cooccurrence, pas causalité :";
        results.append(caption);
        for (const [ref, source] of common) {
          const label = ref + " — " + source.cases.size + " Cases, " +
            source.entries + " événements autorisés";
          // The existing single-Case text query has a 96-character bound.
          // Too-long source IDs remain inert text, never truncated selectors.
          const entry = document.createElement(ref.length <= 96 ? "button" : "p");
          entry.textContent = label;
          if (ref.length <= 96) {
            entry.type = "button";
            entry.addEventListener("click", () =>
              pivotTo(source.cases.values().next().value, ref,
                selectedRevision, bearer, needle));
          }
          results.append(entry);
        }
      }
      status.textContent = matches.length
        ? matches.length + " Case(s) avec correspondance parmi " + cases.length +
          " Case(s) accessibles. Recherche textuelle, pas un verdict de sécurité."
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
