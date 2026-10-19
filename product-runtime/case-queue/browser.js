// Native Case Queue read-only browser. All data are rendered as inert text;
// the browser never chooses a tenant or grants a Case permission.
(() => {
  "use strict";
  const el = id => document.getElementById(id);
  const tbody = el("cases"), status = el("status"), scope = el("scope");
  const first = el("prev"), next = el("next"), page = el("page");
  const search = el("search"), query = el("query");
  const caseStatus = el("case-status"), caseOwner = el("case-owner");
  const caseFreshness = el("case-freshness"), caseSort = el("case-sort");
  // Applied server query is committed only by form submission. An unsubmitted
  // edit must NEVER change the ACL/page against which Inspect is reauthorized.
  let applied = {q:"",status:"",owner:"",freshness:"",sort:""};
  function currentControls() {
    return {q:query.value.trim().slice(0,256),
      status:caseStatus.value.trim().slice(0,128),
      owner:caseOwner.value.trim().slice(0,128),
      freshness:caseFreshness.value,sort:caseSort.value};
  }
  function queueParams(extra = {}) {
    const args = new URLSearchParams({limit:String(limit)});
    for (const [name, value] of Object.entries(applied)) {
      if (value) args.set(name,value);
    }
    for (const [name, value] of Object.entries(extra)) {
      if (value) args.set(name,value);
    }
    return args;
  }
  const selected = el("selection"), selectedStatus = el("selection-status");
  const selectedClose = el("selection-close");
  let cursor = "", history = [], generation = 0, controller = null;
  let selectionGeneration = 0, selectionController = null, selectionTrigger = null;
  const limit = 25;
  // A Workspace handoff is transient in-memory navigation state, never a
  // credential, durable bookmark, browser URL or authorization decision.
  const workspace = el("workspace-handoff");
  const workspaceStatus = el("workspace-status"), workspaceReturn = el("workspace-return");
  let workspaceToken = "", workspaceCase = "", workspaceController = null;
  let workspaceGeneration = 0;
  function clearWorkspace() {
    ++workspaceGeneration;
    if (workspaceController) workspaceController.abort();
    workspaceController = null;
    workspaceToken = workspaceCase = "";
    workspace.hidden = true;
    el("workspace-case").textContent = "";
    workspaceStatus.textContent = "";
  }
  async function openWorkspace(id) {
    clearSelection();
    clearWorkspace();
    const currentWorkspace = workspaceGeneration, currentPage = generation;
    workspaceController = new AbortController();
    workspaceStatus.textContent = "Rechecking Case Workspace access…";
    const offset = Number.isSafeInteger(window.scrollY) && window.scrollY >= 0
      ? Math.min(window.scrollY,10000000) : 0;
    const args = queueParams({open_case:id,cursor,scroll:String(offset)});
    try {
      const response = await fetch("/case-queue?" + args.toString(), {
        credentials:"same-origin",cache:"no-store",redirect:"error",
        referrerPolicy:"no-referrer",signal:workspaceController.signal
      });
      if (currentWorkspace !== workspaceGeneration || currentPage !== generation) return;
      if (!response.ok) throw new Error("Case handoff refused");
      const handoff = await response.json();
      if (currentWorkspace !== workspaceGeneration || currentPage !== generation) return;
      if (!handoff || handoff.screen_id !== "INV-CAS-001" ||
          handoff.case_id !== id || typeof handoff.return_token !== "string" ||
          !handoff.return_token || handoff.return_token.length > 16384) {
        throw new Error("invalid current Case Workspace handoff");
      }
      workspaceToken = handoff.return_token;
      workspaceCase = id;
      clearRows();
      first.disabled = next.disabled = true;
      status.textContent = "Transient Case Workspace handoff; use Return to recheck current Queue access.";
      el("workspace-case").textContent = id;
      workspaceStatus.textContent = "Read-only navigation context. Full Case Workspace UI is not hosted here.";
      workspace.hidden = false;
      workspaceReturn.focus();
    } catch (_) {
      if (currentWorkspace !== workspaceGeneration || currentPage !== generation) return;
      clearWorkspace();
      clearRows();
      first.disabled = next.disabled = true;
      status.textContent = "Case Workspace handoff unavailable; refresh the authorized Queue.";
    }
  }
  async function returnFromWorkspace() {
    if (!workspaceToken || workspace.hidden) return;
    const token = workspaceToken, caseID = workspaceCase;
    const oldCursor = cursor, oldApplied = {...applied};
    const current = ++workspaceGeneration;
    if (workspaceController) workspaceController.abort();
    workspaceController = new AbortController();
    workspaceStatus.textContent = "Reauthorizing current Case Queue return…";
    try {
      const response = await fetch("/case-queue/return", {
        method:"POST",headers:{"Content-Type":"application/json"},
        body:JSON.stringify({return_token:token}),
        credentials:"same-origin",cache:"no-store",redirect:"error",
        referrerPolicy:"no-referrer",signal:workspaceController.signal
      });
      if (current !== workspaceGeneration) return;
      if (!response.ok) throw new Error("Case return denied");
      const returned = await response.json();
      if (current !== workspaceGeneration) return;
      if (!returned || returned.origin !== "case-queue" ||
          returned.selected_case_id !== caseID ||
          returned.page_limit !== limit ||
          (returned.page_cursor || "") !== oldCursor ||
          returned.search_query !== oldApplied.q ||
          !Number.isSafeInteger(returned.scroll_offset) ||
          returned.scroll_offset < 0 || returned.scroll_offset > 10000000 ||
          !Array.isArray(returned.filters) || !Array.isArray(returned.sort)) {
        throw new Error("stale Case Queue return");
      }
      const actual = Object.fromEntries(returned.filters.map(f => [f.field,f.values]));
      for (const key of ["status","owner","freshness"]) {
        // The existing Go filter engine trims/lowercases literal values
        // before serializing its return token; compare that same meaning.
        const expected = (oldApplied[key] || "").trim().toLowerCase();
        if (expected !== (actual[key] && actual[key][0] || "")) {
          throw new Error("returned Queue filters differ");
        }
      }
      const sort = returned.sort.map(x => x.field + ":" + x.direction).join(",");
      if (sort !== oldApplied.sort) throw new Error("returned Queue sorting differs");
      clearWorkspace();
      // POST validates the old page, but a NEW GET must still independently
      // reauthorize it before any Case is displayed or the scroll restored.
      await load(oldCursor);
      if (workspaceGeneration !== current + 1 || status.textContent.includes("unavailable")) return;
      if (typeof window.scrollTo === "function") window.scrollTo(0,returned.scroll_offset);
    } catch (_) {
      if (current !== workspaceGeneration) return;
      clearWorkspace();
      clearRows();
      first.disabled = next.disabled = true;
      status.textContent = "Case Queue return denied or stale. Apply filters to reload access.";
    }
  }
  workspaceReturn.addEventListener("click", () => { void returnFromWorkspace(); });

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
    const args = queueParams({inspect_case:id,cursor});
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
      const open = document.createElement("button");
      open.type = "button";
      open.textContent = "Open workspace handoff";
      open.addEventListener("click", () => { void openWorkspace(row.case.id); });
      tr.children[0].appendChild(open);
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
    const args = queueParams({cursor:requestCursor});
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
    clearWorkspace();
    applied = currentControls();
    history = [];
    cursor = "";
    load();
  });
  el("clear").addEventListener("click", () => {
    clearWorkspace();
    query.value = "";
    caseStatus.value = caseOwner.value = "";
    caseFreshness.value = caseSort.value = "";
    applied = currentControls();
    history = [];
    cursor = "";
    load();
    query.focus();
  });
  next.addEventListener("click", () => {
    const candidate = next.dataset.cursor;
    if (!candidate || next.disabled) return;
    clearWorkspace();
    history.push(cursor);
    cursor = candidate;
    load(cursor);
  });
  first.addEventListener("click", () => {
    clearWorkspace();
    history = [];
    cursor = "";
    load();
  });
  window.addEventListener("pagehide", () => {
    clearWorkspace();
    clearSelection();
    ++generation;
    if (controller) controller.abort();
    clearRows();
  });
  load();
})();
