"use strict";

// This page contains no identity provider or persistent credential storage.
// The analyst enters the ephemeral terminal bearer manually; the backend
// authenticates each Case and each source on every request.
(() => {
  const form = document.getElementById("query");
  const token = document.getElementById("token");
  const caseID = document.getElementById("case-id");
  const kind = document.getElementById("kind");
  const fromInput = document.getElementById("from");
  const toInput = document.getElementById("to");
  const entries = document.getElementById("entries");
  const summary = document.getElementById("summary");
  const message = document.getElementById("message");
  const page = document.getElementById("page");
  const previous = document.getElementById("previous");
  const next = document.getElementById("next");
  const load = document.getElementById("load");
  const limit = 20;
  let offset = 0;
  let pageCount = 0;
  let more = false;
  let busy = false;
  let activeRequest = null;
  let revision = 0;

  function readWindow() {
    const from = fromInput.value ? new Date(fromInput.value) : null;
    const to = toInput.value ? new Date(toInput.value) : null;
    if ((from && !Number.isFinite(from.getTime())) ||
        (to && !Number.isFinite(to.getTime())) ||
        (from && to && from.getTime() >= to.getTime())) {
      return null;
    }
    return { from: from ? from.toISOString() : "",
             to: to ? to.toISOString() : "" };
  }

  function clearResults(description) {
    entries.replaceChildren();
    summary.textContent = description;
    page.textContent = "—";
    pageCount = 0;
    more = false;
    previous.disabled = true;
    next.disabled = true;
  }

  function addLine(parent, label, value) {
    const line = document.createElement("p");
    const name = document.createElement("strong");
    name.textContent = label + " ";
    line.append(name, document.createTextNode(String(value)));
    parent.append(line);
  }

  function render(data) {
    if (!data || !Array.isArray(data.entries) || !data.page ||
        !Number.isSafeInteger(data.page.total_visible) ||
        data.page.total_visible < 0 || data.page.total_visible > 5000 ||
        typeof data.page.has_more !== "boolean" ||
        !Number.isSafeInteger(data.page.offset) ||
        data.page.offset !== offset || data.entries.length > limit) {
      throw new Error("Malformed authorized timeline response");
    }
    entries.replaceChildren();
    for (const event of data.entries) {
      if (!event || typeof event.id !== "string" ||
          typeof event.kind !== "string" ||
          typeof event.source_time !== "string") {
        throw new Error("Malformed timeline entry");
      }
      const item = document.createElement("li");
      if (event.kind === "inferred") item.className = "inferred";
      const name = document.createElement("h3");
      name.textContent = event.id;
      const badge = document.createElement("p");
      badge.className = "type";
      badge.textContent = event.kind;
      item.append(name, badge);
      addLine(item, "Horodatage :", event.source_time +
        (event.time_uncertain ? " (incertain)" : ""));
      addLine(item, "Source :", event.source_ref ?? "—");
      addLine(item, "Provenance :", event.provenance ?? "—");
      addLine(item, "Version :", event.version ?? "—");
      if (event.supersedes) addLine(item, "Remplace :", event.supersedes);
      entries.append(item);
    }
    pageCount = data.entries.length;
    more = data.page.has_more;
    const total = data.page.total_visible;
    summary.textContent = total + " entrée(s) visibles après autorisation";
    page.textContent = total === 0 ? "0 / 0" :
      (offset + 1) + "–" + (offset + pageCount) + " / " + total;
    message.textContent = pageCount === 0 ? "Aucun événement visible pour cette sélection." : "";
  }

  async function fetchPage() {
    const windowSelection = readWindow();
    if (!windowSelection) {
      clearResults("Intervalle non valide");
      message.textContent = "La fin doit être strictement postérieure au début.";
      return;
    }
    const secret = token.value;
    const selectedCase = caseID.value.trim();
    if (!secret || !selectedCase) {
      clearResults("Aucune donnée affichée");
      message.textContent = "Saisir le jeton temporaire et l’identifiant du Case.";
      return;
    }
    const currentRevision = ++revision;
    if (activeRequest) activeRequest.abort();
    const controller = new AbortController();
    activeRequest = controller;
    busy = true;
    load.disabled = true;
    next.disabled = true;
    previous.disabled = true;
    message.textContent = "Chargement des événements autorisés…";
    try {
      const url = new URL("/case-timeline", window.location.origin);
      url.searchParams.set("case_id", selectedCase);
      url.searchParams.set("limit", String(limit));
      url.searchParams.set("offset", String(offset));
      if (kind.value) url.searchParams.set("kind", kind.value);
      if (windowSelection.from) url.searchParams.set("from", windowSelection.from);
      if (windowSelection.to) url.searchParams.set("to", windowSelection.to);
      const response = await fetch(url, {
        method: "GET",
        headers: { Authorization: "Bearer " + secret },
        cache: "no-store",
        credentials: "omit",
        redirect: "error",
        mode: "same-origin",
        signal: controller.signal
      });
      if (!response.ok) throw new Error("Case timeline unavailable");
      const data = await response.json();
      if (currentRevision !== revision) return;
      render(data);
    } catch {
      if (currentRevision !== revision) return;
      clearResults("Aucune donnée affichée");
      message.textContent = "Accès refusé ou requête indisponible. Vérifiez la sélection et le jeton.";
    } finally {
      if (currentRevision === revision) {
        busy = false;
        activeRequest = null;
        load.disabled = false;
        previous.disabled = offset === 0 || pageCount === 0;
        next.disabled = !more;
      }
    }
  }

  form.addEventListener("submit", event => {
    event.preventDefault();
    offset = 0;
    fetchPage();
  });
  previous.addEventListener("click", () => {
    if (busy || offset === 0) return;
    offset = Math.max(0, offset - limit);
    fetchPage();
  });
  next.addEventListener("click", () => {
    if (busy || !more) return;
    offset += limit;
    fetchPage();
  });
  // An old authorized reply must never overwrite the result of a changed
  // Case, token, event kind or time window, even if abort races completion.
  for (const field of [token, caseID, kind, fromInput, toInput]) {
    field.addEventListener("input", () => {
      ++revision;
      if (activeRequest) activeRequest.abort();
      activeRequest = null;
      busy = false;
      load.disabled = false;
      offset = 0;
      clearResults("Sélection modifiée. Rechargez la chronologie.");
      message.textContent = "";
    });
  }
})();
