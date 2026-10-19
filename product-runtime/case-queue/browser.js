// Native Case Queue read-only browser. All data are rendered as inert text;
// the browser never chooses a tenant or grants a Case permission.
(() => {
  "use strict";
  const el = id => document.getElementById(id);
  const tbody = el("cases"), status = el("status"), scope = el("scope");
  const first = el("prev"), next = el("next"), page = el("page");
  const search = el("search"), query = el("query");
  const selected = el("selection"), selectedStatus = el("selection-status");
  const selectedClose = el("selection-close");
  let cursor = "", history = [], generation = 0, controller = null;
  let selectionGeneration = 0, selectionController = null, selectionTrigger = null;
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
  function clearSelection() {
    ++selectionGeneration;
    if (selectionController) selectionController.abort();
    selectionController = null;
    selected.hidden = true;
    selectionTrigger = null;
    for (const id of ["selection-title", "selection-id", "selection-status-value",
      "selection-owner", "selection-updated", "selection-freshness",
      "selection-next-action"]) el(id).textContent = "";
    selectedStatus.textContent = "";
  }
  async function inspectCase(id, trigger) {
    // Selection is non-authorizing; always perform a new independent
    // first-party server Case/Incident/Finding ACL check on THIS page.
    clearSelection();
    const currentSelection = selectionGeneration;
    const currentPage = generation;
    selectionController = new AbortController();
    selectedStatus.textContent = "Rechecking Case access…";
    const args = new URLSearchParams({limit:String(limit),inspect_case:id});
    const term = query.value.trim();
    if (term) args.set("q",term.slice(0,256));
    if (cursor) args.set("cursor",cursor);
    try {
      const response = await fetch("/case-queue?" + args.toString(), {
        credentials:"same-origin",cache:"no-store",redirect:"error",
        referrerPolicy:"no-referrer",signal:selectionController.signal
      });
      if (currentSelection !== selectionGeneration || currentPage !== generation) return;
      if (!response.ok) throw new Error("Case selection no longer accessible");
      const row = await response.json();
      if (currentSelection !== selectionGeneration || currentPage !== generation) return;
      if (!row || !row.case || row.case.id !== id ||
          typeof row.case.title !== "string" || !row.case.id.trim()) {
        throw new Error("invalid authorized Inspector source");
      }
      // Never render raw incident/finding/source object previews. This
      // is a bounded Queue selection summary, NOT a Shared Inspector clone.
      const fields = {
        "selection-title":row.case.title,
        "selection-id":row.case.id,
        "selection-status-value":row.case.status,
        "selection-owner":row.case.owner,
        "selection-updated":row.case.updated_at,
        "selection-freshness":row.stale ? "Stale" : row.partial ? "Partial" : row.case.freshness,
        "selection-next-action":row.case.next_action
      };
      for (const [name, value] of Object.entries(fields)) {
        el(name).textContent = String(value ?? "Unavailable");
      }
      selectedStatus.textContent = "Current Case authorization verified. Read-only selection.";
      selected.hidden = false;
      selectionTrigger = trigger;
      selectedClose.focus();
    } catch (_) {
      if (currentSelection !== selectionGeneration || currentPage !== generation) return;
      clearSelection();
      selectedStatus.textContent = "Case selection unavailable or permission changed.";
    }
  }
  selectedClose.addEventListener("click", () => {
    const trigger = selectionTrigger;
    clearSelection();
    if (trigger && trigger.isConnected !== false) trigger.focus();
  });
  window.addEventListener("keydown", event => {
    if (event.key === "Escape" && !selected.hidden) {
      event.preventDefault();
      const trigger = selectionTrigger;
      clearSelection();
      if (trigger && trigger.isConnected !== false) trigger.focus();
    }
  });
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
      const inspect = document.createElement("button");
      inspect.type = "button";
      inspect.textContent = "Inspect";
      inspect.addEventListener("click", () => inspectCase(row.case.id, inspect));
      tr.children[0].appendChild(inspect);
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
    clearSelection();
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
    clearSelection();
    ++generation;
    if (controller) controller.abort();
    clearRows();
  });
  load();
})();
