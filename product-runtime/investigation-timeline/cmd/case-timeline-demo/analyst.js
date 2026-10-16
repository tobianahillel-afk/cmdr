"use strict";

// This page contains no identity provider or persistent credential storage.
// The analyst enters the ephemeral terminal bearer manually; the backend
// authenticates each Case and each source on every request.
(() => {
  const form = document.getElementById("query");
  const token = document.getElementById("token");
  const caseID = document.getElementById("case-id");
  const kind = document.getElementById("kind");
  const search = document.getElementById("search");
  const caseOptions = document.getElementById("case-options");
  const discover = document.getElementById("discover");
  const discoveryStatus = document.getElementById("case-discovery-status");
  const fromInput = document.getElementById("from");
  const toInput = document.getElementById("to");
  const entries = document.getElementById("entries");
  const summary = document.getElementById("summary");
  const message = document.getElementById("message");
  const page = document.getElementById("page");
  const previous = document.getElementById("previous");
  const next = document.getElementById("next");
  const load = document.getElementById("load");
  const sourceDetail = document.getElementById("source-detail");
  const sourceStatus = document.getElementById("source-status");
  const returnToCase = document.getElementById("return-to-case");
  const copyReturnLink = document.getElementById("copy-return-link");
  const limit = 20;
  let offset = 0;
  let pageCount = 0;
  let more = false;
  let busy = false;
  let activeRequest = null;
  let revision = 0;
  let discoveryRevision = 0;
  let activeDiscovery = null;
  let sourceRevision = 0;
  let sourceRequest = null;
  let returnLinkURL = null;
  let restoreOffsetOnSubmit = false;
  let restoredRange = null;

  // A return URL carries only a presentation selection, never credentials or
  // authorization. It is not fetched until the analyst enters a fresh bearer.
  function localTimestamp(iso) {
    if (!/^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}\\.\\d{3}Z$/.test(iso)) return null;
    const date = new Date(iso);
    if (!Number.isFinite(date.getTime()) || date.toISOString() !== iso) return null;
    const pad = n => String(n).padStart(2, "0");
    return date.getFullYear() + "-" + pad(date.getMonth() + 1) + "-" +
      pad(date.getDate()) + "T" + pad(date.getHours()) + ":" +
      pad(date.getMinutes()) + ":" + pad(date.getSeconds());
  }

  function restoreReturnSelection() {
    const values = new URLSearchParams(window.location.search);
    if (!values.size) return;
    const allowed = new Set(["case_id", "kind", "q", "from", "to", "limit", "offset"]);
    if ([...values.keys()].some(key => !allowed.has(key) || values.getAll(key).length !== 1)) {
      message.textContent = "Contexte de retour invalide : aucune donnée n’a été chargée.";
      return;
    }
    const id = values.get("case_id") || "";
    const selectedKind = values.get("kind") || "";
    const query = values.get("q") || "";
    const size = values.get("limit");
    const pageOffset = values.get("offset") || "0";
    const from = values.has("from") ? localTimestamp(values.get("from")) : "";
    const to = values.has("to") ? localTimestamp(values.get("to")) : "";
    if (!id || id.length > 1024 || id.trim() !== id || id === "*" ||
        !["", "observed", "inferred", "user-action", "decision", "response"].includes(selectedKind) ||
        query.length > 96 || query.trim() !== query ||
        size !== String(limit) || !/^(0|[1-9]\\d*)$/.test(pageOffset) ||
        Number(pageOffset) > 5000 ||
        from === null || to === null ||
        (from && to && new Date(values.get("from")) >= new Date(values.get("to")))) {
      message.textContent = "Contexte de retour invalide : aucune donnée n’a été chargée.";
      return;
    }
    caseID.value = id;
    kind.value = selectedKind;
    search.value = query;
    fromInput.value = from;
    toInput.value = to;
    offset = Number(pageOffset);
    restoredRange = { from: values.get("from") || "", to: values.get("to") || "" };
    restoreOffsetOnSubmit = true;
    message.textContent = "Sélection de retour restaurée. Saisissez le jeton temporaire, puis chargez ; tous les accès seront revérifiés.";
  }

  function clearCaseChoices() {
    caseOptions.replaceChildren();
    caseOptions.add(new Option("Aucun Case chargé", ""));
    caseOptions.disabled = true;
    discoveryStatus.textContent = "";
  }

  async function discoverCases() {
    const bearer = token.value;
    const currentRevision = ++discoveryRevision;
    if (activeDiscovery) activeDiscovery.abort();
    clearCaseChoices();
    if (!bearer) {
      discoveryStatus.textContent = "Saisissez d’abord le jeton temporaire.";
      return;
    }
    const controller = new AbortController();
    activeDiscovery = controller;
    discover.disabled = true;
    discoveryStatus.textContent = "Recherche des Cases autorisés…";
    try {
      const response = await fetch(new URL("/demo-cases", window.location.origin), {
        method: "GET",
        headers: { Authorization: "Bearer " + bearer },
        cache: "no-store",
        credentials: "omit",
        redirect: "error",
        mode: "same-origin",
        signal: controller.signal
      });
      if (!response.ok) throw new Error("Case catalogue unavailable");
      const result = await response.json();
      if (currentRevision !== discoveryRevision) return;
      if (!result || !Array.isArray(result.cases) || result.cases.length > 5000 ||
          result.cases.some(id => typeof id !== "string" || !id ||
            id.length > 1024)) {
        throw new Error("Malformed Case catalogue");
      }
      for (const id of result.cases) {
        // Option.textContent is inert even if an owner-supplied ID has markup.
        const option = document.createElement("option");
        option.value = id;
        option.textContent = id;
        caseOptions.add(option);
      }
      caseOptions.disabled = result.cases.length === 0;
      discoveryStatus.textContent = result.cases.length + " Case(s) accessibles.";
    } catch {
      if (currentRevision !== discoveryRevision) return;
      clearCaseChoices();
      discoveryStatus.textContent = "Catalogue indisponible ou accès refusé.";
    } finally {
      if (currentRevision === discoveryRevision) {
        activeDiscovery = null;
        discover.disabled = false;
      }
    }
  }

  function readWindow() {
    // Preserve the exact RFC3339 timestamps on first return: local times
    // during a daylight-saving fold are ambiguous if round-tripped by Date.
    if (restoredRange) return restoredRange;
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

  function clearSource() {
    ++sourceRevision;
    returnLinkURL = null;
    copyReturnLink.disabled = true;
    if (sourceRequest) sourceRequest.abort();
    sourceRequest = null;
    sourceDetail.hidden = true;
    sourceStatus.textContent = "";
    returnToCase.disabled = true;
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
      const open = document.createElement("button");
      open.type = "button";
      open.textContent = "Vérifier la destination de la source";
      open.addEventListener("click", () => openSource(event.id));
      item.append(open);
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


  async function openSource(entryID) {
    const currentRevision = ++sourceRevision;
    if (sourceRequest) sourceRequest.abort();
    sourceDetail.hidden = true;
    returnToCase.disabled = true;
    const selectedCase = caseID.value.trim();
    const bearer = token.value;
    const windowSelection = readWindow();
    if (!selectedCase || !bearer || !windowSelection) {
      sourceStatus.textContent = "Sélection incomplète.";
      return;
    }
    const controller = new AbortController();
    sourceRequest = controller;
    sourceStatus.textContent = "Nouvelle vérification des droits de destination…";
    try {
      const url = new URL("/case-source-navigation", window.location.origin);
      url.searchParams.set("case_id", selectedCase);
      url.searchParams.set("entry_id", entryID);
      url.searchParams.set("limit", String(limit));
      url.searchParams.set("offset", String(offset));
      if (kind.value) url.searchParams.set("kind", kind.value);
      if (search.value.trim()) url.searchParams.set("q", search.value.trim());
      if (windowSelection.from) url.searchParams.set("from", windowSelection.from);
      if (windowSelection.to) url.searchParams.set("to", windowSelection.to);
      const response = await fetch(url, {
        method: "GET",
        headers: { Authorization: "Bearer " + bearer },
        credentials: "omit",
        cache: "no-store",
        redirect: "error",
        mode: "same-origin",
        signal: controller.signal
      });
      if (!response.ok) throw new Error("Destination access denied");
      const data = await response.json();
      if (currentRevision !== sourceRevision) return;
      if (!data || data.status !== "ready" || data.return_is_authority !== false ||
          data.return_case_id !== selectedCase || data.return_entry_id !== entryID ||
          !["evidence", "telemetry-event", "decision", "incident"].includes(data.kind) ||
          typeof data.object_id !== "string" || data.object_id.length === 0 ||
          data.object_id.length > 1024) {
        throw new Error("Invalid destination draft");
      }
      // The server's return_query is untrusted presentation data, not a
      // permission, redirect or arbitrary destination. Bind it to exactly
      // the request selection we just sent, excluding the entry selector.
      const expected = new URLSearchParams(url.searchParams);
      expected.delete("entry_id");
      const received = new URLSearchParams(data.return_query);
      if (received.size !== expected.size ||
          [...received].some(([key, value]) =>
            expected.getAll(key).length !== 1 || expected.get(key) !== value)) {
        throw new Error("Return selection differs from authorized request");
      }
      const localReturn = new URL("/analyst", window.location.origin);
      localReturn.search = expected.toString();
      returnLinkURL = localReturn.href;
      copyReturnLink.disabled = false;
      sourceDetail.replaceChildren();
      addLine(sourceDetail, "Objet source autorisé :", data.kind + " / " + data.object_id);
      addLine(sourceDetail, "Entrée d’origine :", data.return_entry_id);
      addLine(sourceDetail, "Navigation :", "brouillon ; l’application propriétaire doit contrôler l’accès à l’ouverture.");
      sourceDetail.hidden = false;
      returnToCase.disabled = false;
      sourceStatus.textContent = "Destination vérifiée pour cette requête, sans lien externe.";
    } catch {
      if (currentRevision !== sourceRevision) return;
      sourceDetail.hidden = true;
      returnLinkURL = null;
      copyReturnLink.disabled = true;
      sourceStatus.textContent = "Objet source inaccessible ou autorisation révoquée.";
    } finally {
      if (currentRevision === sourceRevision) sourceRequest = null;
    }
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
    clearSource();
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
      if (search.value.trim()) url.searchParams.set("q", search.value.trim());
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

  copyReturnLink.addEventListener("click", async () => {
    if (!returnLinkURL) return;
    try {
      await navigator.clipboard.writeText(returnLinkURL);
      sourceStatus.textContent = "Lien local de retour copié : il contient la sélection du Case, jamais le jeton. Les droits seront revérifiés.";
    } catch {
      sourceStatus.textContent = "Copie indisponible : utilisez le bouton de retour sans quitter la page.";
    }
  });

  returnToCase.addEventListener("click", () => {
    clearSource();
    fetchPage(); // reauthorize the Case and Timeline; keep current filters and offset
  });
  discover.addEventListener("click", discoverCases);
  caseOptions.addEventListener("change", () => {
    if (!caseOptions.value) return;
    caseID.value = caseOptions.value;
    caseID.dispatchEvent(new Event("input", { bubbles: true }));
  });

  form.addEventListener("submit", event => {
    event.preventDefault();
    if (!restoreOffsetOnSubmit) offset = 0;
    restoreOffsetOnSubmit = false;
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
  for (const field of [token, caseID, kind, search, fromInput, toInput]) {
    field.addEventListener("input", () => {
      if (field === token) {
        ++discoveryRevision;
        if (activeDiscovery) activeDiscovery.abort();
        activeDiscovery = null;
        discover.disabled = false;
        clearCaseChoices();
      }
      ++revision;
      clearSource();
      if (activeRequest) activeRequest.abort();
      activeRequest = null;
      busy = false;
      load.disabled = false;
      // Entering a fresh bearer must NOT discard a restored page/time
      // selection; edits to Case and filter controls explicitly do.
      if (field !== token) {
        offset = 0;
        restoreOffsetOnSubmit = false;
        restoredRange = null;
      }
      clearResults("Sélection modifiée. Rechargez la chronologie.");
      message.textContent = "";
    });
  }
  restoreReturnSelection();
})();
