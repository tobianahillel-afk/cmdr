// Native Case Queue read-only browser. All data are rendered as inert text;
// the browser never chooses a tenant or grants a Case permission.
(() => {
  "use strict";
  const el = id => document.getElementById(id);
  const tbody = el("cases"), status = el("status"), scope = el("scope");
  const first = el("prev"), next = el("next"), page = el("page");
  const search = el("search"), query = el("query");
  let cursor = "", history = [], generation = 0, controller = null;
  const limit = 25;

  function clearRows() {
    tbody.replaceChildren();
    scope.textContent = "";
  }
  function addCell(tr, value) {
    const td = document.createElement("td");
    td.textContent = String(value ?? "Unavailable");
    tr.appendChild(td);
  }
  function render(data) {
    if (!data || !Array.isArray(data.rows) || data.rows.length > limit ||
        (data.next_cursor !== undefined &&
         (typeof data.next_cursor !== "string" || data.next_cursor.length > 256))) {
      throw new Error("invalid authorized Queue response");
    }
    clearRows();
    for (const row of data.rows) {
      if (!row || !row.case || typeof row.case.id !== "string" ||
          !row.case.id.trim() || typeof row.case.title !== "string") {
        throw new Error("invalid authorized Case row");
      }
      const tr = document.createElement("tr");
      addCell(tr, row.case.title + " — " + (row.case.human_id || row.case.id));
      addCell(tr, row.case.status);
      addCell(tr, row.case.owner);
      addCell(tr, row.stale ? "Stale" : row.partial ? "Partial" : row.case.freshness);
      addCell(tr, row.case.updated_at);
      tbody.appendChild(tr);
    }
    scope.textContent = [data.tenant_ref, data.environment_ref].filter(Boolean).join(" / ");
    next.disabled = !data.next_cursor;
    next.dataset.cursor = data.next_cursor || "";
    first.disabled = history.length === 0;
    page.textContent = "Page " + (history.length + 1);
    status.textContent = data.rows.length
      ? data.rows.length + " currently accessible Cases on this page."
      : "No accessible Cases on this page.";
  }
  async function load(requestCursor = "") {
    const current = ++generation;
    if (controller) controller.abort();
    controller = new AbortController();
    clearRows(); // never keep old private Case data after a new auth attempt
    status.textContent = "Checking current Case permissions…";
    first.disabled = next.disabled = true;
    const args = new URLSearchParams({ limit: String(limit) });
    const term = query.value.trim();
    if (term) args.set("q", term.slice(0, 256));
    if (requestCursor) args.set("cursor", requestCursor);
    try {
      const reply = await fetch("/case-queue?" + args.toString(), {
        credentials: "same-origin", cache: "no-store", redirect: "error",
        referrerPolicy: "no-referrer", signal: controller.signal
      });
      if (current !== generation) return;
      if (!reply.ok) throw new Error("Case permissions or page have changed");
      const data = await reply.json();
      if (current !== generation) return;
      render(data);
    } catch (err) {
      if (current !== generation) return;
      clearRows();
      history = [];
      cursor = "";
      first.disabled = next.disabled = true;
      page.textContent = "Page 1";
      status.textContent = "Case Queue unavailable. Retry your search; no old Case data is displayed.";
    }
  }
  search.addEventListener("submit", event => {
    event.preventDefault();
    history = [];
    cursor = "";
    load();
  });
  el("clear").addEventListener("click", () => {
    query.value = "";
    history = [];
    cursor = "";
    load();
    query.focus();
  });
  next.addEventListener("click", () => {
    const candidate = next.dataset.cursor;
    if (!candidate || next.disabled) return;
    history.push(cursor);
    cursor = candidate;
    load(cursor);
  });
  first.addEventListener("click", () => {
    history = [];
    cursor = "";
    load();
  });
  window.addEventListener("pagehide", () => {
    ++generation;
    if (controller) controller.abort();
    clearRows();
  });
  load();
})();
