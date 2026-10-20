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
  const wazuhButton = document.getElementById("wazuh-cohorts");
  const wazuhStatus = document.getElementById("wazuh-status");
  const wazuhActions = document.getElementById("wazuh-actions");
  const wazuhEvents = document.getElementById("wazuh-events");
  const wazuhPrevious = document.getElementById("wazuh-previous");
  const wazuhNext = document.getElementById("wazuh-next");
  let wazuhSelection = null;
  let wazuhOffset = 0;
  let wazuhRevision = 0;
  let wazuhRequest = null;
  function clearWazuh() {
    ++wazuhRevision;
    if (wazuhRequest) wazuhRequest.abort();
    wazuhRequest = null;
    wazuhButton.disabled = false;
    wazuhSelection = null;
    wazuhOffset = 0;
    wazuhPrevious.disabled = true;
    wazuhNext.disabled = true;
    wazuhActions.replaceChildren();
    wazuhEvents.replaceChildren();
    wazuhEvents.textContent = "Aucun événement Wazuh sélectionné.";
    wazuhStatus.textContent = "Aucune corrélation Wazuh demandée.";
  }
  async function loadWazuhCohorts() {
    const bearer = token.value;
    const selectedCase = caseID.value.trim();
    clearWazuh();
    if (!bearer || !selectedCase) {
      wazuhStatus.textContent = "Sélectionnez un Case et son jeton temporaire.";
      return;
    }
    const currentRevision = wazuhRevision;
    const currentSelection = revision;
    const controller = new AbortController();
    wazuhRequest = controller;
    wazuhButton.disabled = true;
    wazuhStatus.textContent = "Vérification des Case, Timeline et sources Wazuh…";
    try {
      const url = new URL("/demo-wazuh-cohorts", window.location.origin);
      url.searchParams.set("case_id", selectedCase);
      const response = await fetch(url, {
        method: "GET", headers: { Authorization: "Bearer " + bearer },
        credentials: "omit", cache: "no-store", signal: controller.signal
      });
      if (!response.ok) throw Error("Wazuh unavailable");
      const data = await response.json();
      if (controller.signal.aborted || currentRevision !== wazuhRevision ||
          currentSelection !== revision || token.value !== bearer ||
          caseID.value.trim() !== selectedCase) return;
      if (!data || data.format !== "cmdr-local-wazuh4-case-triage-v1" ||
          data.synthetic !== true || !Number.isSafeInteger(data.imported_entries) ||
          data.imported_entries < 1 || data.imported_entries > 5000 ||
          !Array.isArray(data.top_agents) || !Array.isArray(data.top_rules) ||
          data.top_agents.length > 10 || data.top_rules.length > 10) {
        throw Error("Invalid cohort schema");
      }
      const validID = value => typeof value === "string" && /^[0-9]{1,16}$/.test(value);
      const validCount = value => Number.isSafeInteger(value) && value >= 0 && value <= 5000;
      const validLevel = value => value === undefined || (Number.isInteger(value) && value >= 0 && value <= 16);
      if (data.top_agents.some(a => !a || !validID(a.agent_id) ||
          !validCount(a.count) || !validCount(a.distinct_rule_ids) ||
          !validLevel(a.highest_known_rule_level)) ||
          data.top_rules.some(a => !a || !validID(a.rule_id) ||
          !validCount(a.count) || !validLevel(a.highest_known_rule_level))) {
        throw Error("Invalid cohort values");
      }
      // The optional additive source-cohort schema is absent in older snapshots.
      // When present, validate the COMPLETE group before rendering anything.
      const hasTechniques = data.top_techniques !== undefined;
      if (hasTechniques) {
        const validTechnique = value =>
          typeof value === "string" && /^T[0-9]{4}(?:\.[0-9]{3})?$/.test(value);
        const mapped = data.imported_entries - data.unmapped_technique_alerts;
        const top = data.top_techniques;
        if (!Array.isArray(top) ||
            !validCount(data.unmapped_technique_alerts) ||
            data.unmapped_technique_alerts > data.imported_entries ||
            !Number.isSafeInteger(data.distinct_technique_ids) ||
            data.distinct_technique_ids < 0 ||
            data.distinct_technique_ids > mapped * 16 ||
            top.length !== Math.min(10, data.distinct_technique_ids) ||
            top.some(a => !a || !validTechnique(a.technique_id) ||
              !Number.isSafeInteger(a.count) || a.count < 1 || a.count > mapped ||
              !validLevel(a.highest_known_rule_level)) ||
            new Set(top.map(a => a.technique_id)).size !== top.length ||
            top.some((a, i) => i > 0 &&
              (a.count > top[i - 1].count ||
               (a.count === top[i - 1].count && a.technique_id <= top[i - 1].technique_id))) ||
            (data.distinct_technique_ids <= 10 &&
             top.reduce((sum, a) => sum + a.count, 0) < mapped)) {
          throw Error("Invalid Wazuh ATT&CK source cohorts");
        }
      } else if (data.distinct_technique_ids !== undefined ||
                 data.unmapped_technique_alerts !== undefined) {
        throw Error("Partial Wazuh ATT&CK source cohorts");
      }
      // Always use textContent: source labels are never interpreted as HTML.
      const agentLines = data.top_agents.map(a => "agent " + a.agent_id + " : " + a.count +
        " alerte(s), " + a.distinct_rule_ids + " règle(s), niveau maximum connu " +
        (a.highest_known_rule_level ?? "inconnu"));
      const ruleLines = data.top_rules.map(a => "règle " + a.rule_id + " : " + a.count +
        " alerte(s), niveau maximum connu " + (a.highest_known_rule_level ?? "inconnu"));
      const techniqueLines = hasTechniques ? data.top_techniques.map(a =>
        "étiquette " + a.technique_id + " : " + a.count +
        " alerte(s), niveau maximum connu " + (a.highest_known_rule_level ?? "inconnu")) : [];
      wazuhStatus.textContent = data.imported_entries + " alerte(s) Wazuh importée(s) et autorisée(s)." +
        "\nAgents les plus touchés :\n" + (agentLines.join("\n") || "aucun agent connu") +
        "\nRègles les plus fréquentes :\n" + (ruleLines.join("\n") || "aucune règle connue") +
        (hasTechniques ?
          "\nÉtiquettes MITRE ATT&CK déclarées par Wazuh, non vérifiées :\n" +
            (techniqueLines.join("\n") || "aucune étiquette connue") +
            "\nAlertes sans étiquette MITRE : " + data.unmapped_technique_alerts : "") +
        "\nDonnées locales synthétiques, sans preuve de causalité.";
      for (const [selector, cohorts] of [
        ["agent_id", data.top_agents.map(a => a.agent_id)],
        ["rule_id", data.top_rules.map(a => a.rule_id)]
      ]) {
        for (const value of cohorts) {
          const button = document.createElement("button");
          button.type = "button";
          button.textContent = "Voir les événements Wazuh — " +
            (selector === "agent_id" ? "agent " : "règle ") + value;
          button.addEventListener("click", () => loadWazuhDrilldown(selector, value));
          wazuhActions.append(button);
        }
      }
      // Intersection is a server-owned authorized filter, not a permission.
      // Reuse the cohort response; do not perform client-side union or join.
      if (data.top_agents.length && data.top_rules.length) {
        const agentPicker = document.createElement("select");
        agentPicker.ariaLabel = "Agent Wazuh à croiser";
        const rulePicker = document.createElement("select");
        rulePicker.ariaLabel = "Règle Wazuh à croiser";
        for (const agent of data.top_agents) {
          agentPicker.add(new Option("Agent " + agent.agent_id, agent.agent_id));
        }
        for (const rule of data.top_rules) {
          rulePicker.add(new Option("Règle " + rule.rule_id, rule.rule_id));
        }
        const intersect = document.createElement("button");
        intersect.type = "button";
        intersect.textContent = "Croiser l'agent et la règle";
        intersect.addEventListener("click", () =>
          loadWazuhDrilldown("agent_rule", agentPicker.value, 0, rulePicker.value));
        wazuhActions.append(agentPicker, rulePicker, intersect);
      }
    } catch {
      if (currentRevision === wazuhRevision) {
        wazuhStatus.textContent = "Corrélation indisponible ou accès refusé. Aucune donnée partielle affichée.";
      }
    } finally {
      if (currentRevision === wazuhRevision) {
        wazuhRequest = null;
        wazuhButton.disabled = false;
      }
    }
  }
  // A cohort is a triage selector, never a permission. This route reruns the
  // server-owned Case/Timeline/typed-source traversal before returning IDs.
  async function loadWazuhDrilldown(selector, value, requestedOffset = 0, secondValue = "") {
    const bearer = token.value;
    const selectedCase = caseID.value.trim();
    if (!bearer || !selectedCase || !/^[0-9]{1,16}$/.test(value) ||
        !["agent_id", "rule_id", "agent_rule"].includes(selector) ||
        (selector === "agent_rule" ? !/^[0-9]{1,16}$/.test(secondValue) : secondValue !== "") ||
        !Number.isInteger(requestedOffset) || requestedOffset < 0 ||
        requestedOffset > 5000 || requestedOffset % 25 !== 0) return;
    if (wazuhRequest) wazuhRequest.abort();
    const controller = new AbortController();
    wazuhRequest = controller;
    wazuhPrevious.disabled = true;
    wazuhNext.disabled = true;
    const selectedRevision = wazuhRevision;
    const caseRevision = revision;
    wazuhEvents.replaceChildren();
    wazuhEvents.textContent = "Nouvelle vérification des événements Wazuh…";
    try {
      const url = new URL("/demo-wazuh-drilldown", window.location.origin);
      url.searchParams.set("case_id", selectedCase);
      url.searchParams.set(selector === "rule_id" ? "rule_id" : "agent_id", value);
      if (selector === "agent_rule") url.searchParams.set("rule_id", secondValue);
      url.searchParams.set("offset", String(requestedOffset));
      const response = await fetch(url, {
        method: "GET", headers: { Authorization: "Bearer " + bearer },
        credentials: "omit", cache: "no-store", redirect: "error",
        mode: "same-origin", signal: controller.signal
      });
      if (!response.ok) throw Error("Wazuh drilldown denied");
      const data = await response.json();
      if (controller.signal.aborted || wazuhRequest !== controller ||
          selectedRevision !== wazuhRevision ||
          caseRevision !== revision || token.value !== bearer ||
          caseID.value.trim() !== selectedCase) return;
      if (!data || data.format !== "cmdr-local-wazuh4-drilldown-v1" ||
          data.synthetic !== true || data.selector !== selector ||
          data.value !== value ||
          data.secondary_value !== (selector === "agent_rule" ? secondValue : undefined) ||
          !Number.isSafeInteger(data.total_matches) ||
          data.total_matches < 0 || data.total_matches > 5000 ||
          !Array.isArray(data.entries) || data.entries.length > 25 ||
          data.entries.length !== Math.min(25, Math.max(0, data.total_matches - requestedOffset)) ||
          (data.offset !== undefined && data.offset !== requestedOffset) ||
          (data.has_more !== undefined &&
            (typeof data.has_more !== "boolean" ||
             data.has_more !== (requestedOffset + data.entries.length < data.total_matches))) ||
          data.entries.some(e => !e || typeof e.id !== "string" ||
            e.id.length < 1 || e.id.length > 1024 ||
            !["observed", "inferred", "user-action", "decision", "response"].includes(e.kind) ||
            typeof e.source_time !== "string" || e.source_time.length > 64 ||
            typeof e.time_uncertain !== "boolean")) {
        throw Error("Invalid Wazuh drilldown");
      }
      wazuhSelection = { selector, value, secondValue };
      wazuhOffset = requestedOffset;
      wazuhPrevious.disabled = requestedOffset === 0;
      wazuhNext.disabled = requestedOffset + data.entries.length >= data.total_matches;
      wazuhEvents.replaceChildren();
      wazuhEvents.textContent = data.total_matches +
        " événement(s) Wazuh autorisé(s), lignes " +
        (data.entries.length ? requestedOffset + 1 : 0) + "–" +
        (requestedOffset + data.entries.length) +
        ". Chaque page et ouverture revérifie les droits.";
      const list = document.createElement("ol");
      for (const entry of data.entries) {
        const item = document.createElement("li");
        const button = document.createElement("button");
        button.type = "button";
        button.textContent = entry.kind + " · " + entry.source_time +
          (entry.time_uncertain ? " (incertain)" : "") + " · " + entry.id.slice(0, 120);
        button.addEventListener("click", () => {
          if (controller.signal.aborted || selectedRevision !== wazuhRevision ||
              caseRevision !== revision || token.value !== bearer ||
              caseID.value.trim() !== selectedCase) return;
          openSource(entry.id); // existing authoritative navigation ACL chain
        });
        item.append(button);
        list.append(item);
      }
      wazuhEvents.append(list);
    } catch {
      if (selectedRevision === wazuhRevision) {
        wazuhEvents.replaceChildren();
        wazuhEvents.textContent = "Événements indisponibles ou droits modifiés. Aucun résultat partiel.";
        wazuhPrevious.disabled = true;
        wazuhNext.disabled = true;
      }
    } finally {
      if (wazuhRequest === controller) wazuhRequest = null;
    }
  }
  wazuhButton.addEventListener("click", loadWazuhCohorts);
  wazuhPrevious.addEventListener("click", () => {
    if (wazuhSelection && wazuhOffset > 0 && !wazuhRequest)
      loadWazuhDrilldown(wazuhSelection.selector, wazuhSelection.value, wazuhOffset - 25, wazuhSelection.secondValue);
  });
  wazuhNext.addEventListener("click", () => {
    if (wazuhSelection && !wazuhNext.disabled && !wazuhRequest)
      loadWazuhDrilldown(wazuhSelection.selector, wazuhSelection.value, wazuhOffset + 25, wazuhSelection.secondValue);
  });
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

  const analyze = document.getElementById("analyze");
  const analysis = document.getElementById("analysis");
  const analysisActions = document.getElementById("analysis-actions");
  const traceActions = document.getElementById("trace-actions");
  const traceResults = document.getElementById("trace-results");
  const compareCase = document.getElementById("compare-case");
  const compareButton = document.getElementById("compare");
  const comparison = document.getElementById("comparison");
  const comparisonActions = document.getElementById("comparison-actions");
  const comparisonTraceActions = document.getElementById("comparison-trace-actions");
  const comparisonTraceResults = document.getElementById("comparison-trace-results");
  const downloadComparison = document.getElementById("download-comparison");
  const downloadCrossEvents = document.getElementById("download-cross-events");
  const proximityButton = document.getElementById("proximity");
  const proximitySeconds = document.getElementById("proximity-seconds");
  const proximityResults = document.getElementById("proximity-results");
  let preparedProximity = null;
  proximityButton.disabled = true;
  let preparedComparison = null;
  let preparedCrossEvents = null;
  downloadCrossEvents.disabled = true;
  downloadComparison.disabled = true;
  let compareRevision = 0;
  let compareRequest = null;
  function clearComparison() {
    ++compareRevision;
    if (compareRequest) compareRequest.abort();
    compareRequest = null;
    compareButton.disabled = false;
    comparison.textContent = "Aucune comparaison calculée.";
    comparisonActions.replaceChildren();
    comparisonTraceActions.replaceChildren();
    comparisonTraceResults.replaceChildren();
    comparisonTraceResults.textContent = "Aucune séquence croisée sélectionnée.";
    preparedComparison = null;
    preparedCrossEvents = null;
    preparedProximity = null;
    proximityButton.disabled = true;
    proximityResults.replaceChildren();
    proximityResults.textContent = "Aucune piste temporelle calculée.";
    downloadComparison.disabled = true;
    downloadCrossEvents.disabled = true;
  }
  const downloadReport = document.getElementById("download-report");
  let analysisRequest = null;
  let analysisRevision = 0;
  let preparedReport = null;
  // Match the disabled HTML default even in minimal embedded DOM hosts.
  downloadReport.disabled = true;

  function clearAnalysis() {
    ++analysisRevision;
    if (analysisRequest) analysisRequest.abort();
    analysisRequest = null;
    analyze.disabled = false;
    preparedReport = null;
    downloadReport.disabled = true;
    analysis.textContent = "Aucune synthèse calculée.";
    analysisActions.replaceChildren();
    traceActions.replaceChildren();
    traceResults.replaceChildren();
    traceResults.textContent = "Aucune séquence de source sélectionnée.";
  }

  // Aggregate ONLY entries returned by the existing server-side Case/Entry/
  // source-authorized endpoint. Each page is freshly authorized, so reject
  // changing totals and never display a partial summary as definitive.
  async function analyzeAuthorizedCase() {
    const bearer = token.value;
    const selectedCase = caseID.value.trim();
    const selectedWindow = readWindow();
    const selectedKind = kind.value;
    const selectedSearch = search.value.trim();
    clearAnalysis();
    if (!bearer || !selectedCase || !selectedWindow) {
      analysis.textContent = "Sélection incomplète pour la synthèse.";
      return;
    }
    const currentRevision = analysisRevision;
    const currentSelection = revision;
    const controller = new AbortController();
    analysisRequest = controller;
    analyze.disabled = true;
    analysis.textContent = "Analyse des pages autorisées en cours…";
    const pageSize = 200;
    const kinds = ["observed", "inferred", "user-action", "decision", "response"];
    const counts = Object.fromEntries(kinds.map(value => [value, 0]));
    const seen = new Set();
    const bySource = new Map();
    const certainTimes = [];
    let expectedTotal = null;
    let uncertain = 0;
    try {
      for (let offset = 0; offset === 0 || offset < expectedTotal; offset += pageSize) {
        if (controller.signal.aborted) throw new Error("cancelled");
        const url = new URL("/case-timeline", window.location.origin);
        url.searchParams.set("case_id", selectedCase);
        url.searchParams.set("limit", String(pageSize));
        url.searchParams.set("offset", String(offset));
        if (selectedKind) url.searchParams.set("kind", selectedKind);
        if (selectedSearch) url.searchParams.set("q", selectedSearch);
        if (selectedWindow.from) url.searchParams.set("from", selectedWindow.from);
        if (selectedWindow.to) url.searchParams.set("to", selectedWindow.to);
        const response = await fetch(url, {
          method: "GET",
          headers: { Authorization: "Bearer " + bearer },
          cache: "no-store", credentials: "omit", redirect: "error",
          mode: "same-origin", signal: controller.signal
        });
        if (!response.ok) throw new Error("Case unavailable");
        const data = await response.json();
        if (controller.signal.aborted || currentRevision !== analysisRevision ||
            currentSelection !== revision) return;
        if (!data || !data.page || !Array.isArray(data.entries) ||
            !Number.isSafeInteger(data.page.total_visible) ||
            data.page.total_visible < 0 || data.page.total_visible > 5000 ||
            data.page.limit !== pageSize || data.page.offset !== offset ||
            !["allow", "partial", "filtered", "uncertain", "empty"].includes(data.status) ||
            (expectedTotal !== null && expectedTotal !== data.page.total_visible)) {
          throw new Error("Inconsistent authorized pagination");
        }
        expectedTotal = data.page.total_visible;
        const count = Math.min(pageSize, Math.max(0, expectedTotal - offset));
        if (data.entries.length !== count ||
            data.page.has_more !== (offset + count < expectedTotal)) {
          throw new Error("Incomplete authorized page");
        }
        for (const entry of data.entries) {
          if (!entry || typeof entry.id !== "string" || !entry.id ||
              seen.has(entry.id) || !kinds.includes(entry.kind) ||
              typeof entry.time_uncertain !== "boolean" ||
              typeof entry.source_ref !== "string" || !entry.source_ref) {
            throw new Error("Invalid authorized timeline entry");
          }
          seen.add(entry.id);
          counts[entry.kind]++;
          if (entry.time_uncertain) {
            uncertain++;
          } else {
            // Trust the already authorized Entry, not the browser clock.
            // Only unambiguous RFC3339 timestamps contribute to elapsed time.
            if (typeof entry.source_time !== "string" ||
                !/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d{1,9})?(?:Z|[+-]\d{2}:\d{2})$/.test(entry.source_time)) {
              throw new Error("Invalid certain timestamp");
            }
            const millis = Date.parse(entry.source_time);
            if (!Number.isFinite(millis)) throw new Error("Unparseable certain timestamp");
            certainTimes.push(millis);
          }
          // Correlate only source references already released after the
          // existing server-side Case, Timeline Entry and source ACL chain.
          const source = bySource.get(entry.source_ref) ||
            { ref: entry.source_ref, entryID: entry.id, count: 0, observed: 0, inferred: 0, rows: [] };
          source.count++;
          source.rows.push({
            id: entry.id, kind: entry.kind, time: entry.source_time,
            uncertain: entry.time_uncertain
          });
          if (entry.kind === "observed") source.observed++;
          if (entry.kind === "inferred") source.inferred++;
          bySource.set(entry.source_ref, source);
        }
      }
      if (currentRevision !== analysisRevision || currentSelection !== revision) return;
      const ranked = [...bySource.values()].sort((a, b) =>
        b.count - a.count || (a.ref < b.ref ? -1 : a.ref > b.ref ? 1 : 0));
      // The HTTP pages may contain same-time observations and inferred
      // entries. Rank only certain instants, never infer missing activity.
      certainTimes.sort((a, b) => a - b);
      let largestGapMS = 0;
      for (let i = 1; i < certainTimes.length; i++) {
        largestGapMS = Math.max(largestGapMS, certainTimes[i] - certainTimes[i - 1]);
      }
      const chronology = certainTimes.length >= 2
        ? " · Plage des horodatages certains (UTC) : " +
          new Date(certainTimes[0]).toISOString() + " → " +
          new Date(certainTimes[certainTimes.length - 1]).toISOString() +
          " · Plus grand écart entre événements visibles certains : " +
          (largestGapMS / 1000).toFixed(3) + " s"
        : " · Écarts temporels non calculables (moins de deux horodatages certains)";
      const topSources = ranked.slice(0, 3).map(source => {
        // Cap visible labels; do not place owner identifiers in HTML or URLs.
        const label = source.ref.length > 120 ?
          source.ref.slice(0, 120) + "… (abrégée)" : source.ref;
        return label + " : " + source.count + " événement(s), " +
          source.observed + " observé(s), " + source.inferred + " inféré(s)";
      });
      analysis.textContent = expectedTotal + " événement(s) autorisé(s) · " +
        kinds.map(value => value + " : " + counts[value]).join(" · ") +
        " · horodatage incertain : " + uncertain + chronology +
        ". Sources présentes : " + ranked.length +
        (topSources.length ? " · Principales références : " +
          topSources.join(" · ") : "") +
        ". Regroupement par référence, pas une preuve de causalité. " +
        "Un écart entre événements visibles ne démontre jamais l'absence d'activité " +
        "(données masquées, non collectées ou horodatages incertains). " +
        "Synthèse locale à l'instant de lecture ; les pages ne sont pas transactionnelles.";
      // Only a bounded aggregate is prepared for explicit local download:
      // no bearer, no hidden entries, no raw event rows, no URLs or owner preview.
      // This is an unsigned, nontransactional demo summary, not a legal
      // chain-of-custody attestation or production forensic export.
      const report = {
        schema_version: 1,
        type: "cmdr-synthetic-authorized-case-summary",
        case_id: selectedCase,
        filters: {
          kind: selectedKind, q: selectedSearch,
          from: selectedWindow.from, to: selectedWindow.to
        },
        generated_at_utc: new Date().toISOString(),
        count: expectedTotal,
        event_kinds: counts,
        uncertain_timestamps: uncertain,
        certain_time_range_utc: certainTimes.length ? {
          from: new Date(certainTimes[0]).toISOString(),
          to: new Date(certainTimes[certainTimes.length - 1]).toISOString()
        } : null,
        largest_visible_certain_gap_ms: certainTimes.length >= 2 ? largestGapMS : null,
        unique_sources: ranked.length,
        top_sources: ranked.slice(0, 3).map(source => ({
          source_ref: source.ref,
          count: source.count,
          observed: source.observed,
          inferred: source.inferred
        })),
        limitations: [
          "Local synthetic Case snapshot; not a production IAM or verified source feed",
          "Only server-authorized entries visible at paginated read time",
          "Pages are not transactionally isolated; hidden or uncollected events may exist",
          "Unsigned aggregate without raw events, cryptographic custody proof or audit backend"
        ]
      };
      preparedReport = { serialized: JSON.stringify(report, null, 2) + "\n",
        revision: currentRevision, selection: currentSelection, caseID: selectedCase };
      downloadReport.disabled = false;
      // A representative authorized Entry is only a UI selector, never a
      // source capability. Both existing routes independently reauthorize
      // Case, Entry and typed source when the analyst clicks a source.
      for (const source of ranked.slice(0, 3)) {
        const label = source.ref.length > 120 ?
          source.ref.slice(0, 120) + "… (abrégée)" : source.ref;
        const button = document.createElement("button");
        button.type = "button";
        button.textContent = "Ouvrir une entrée de cette source : " + label;
        button.addEventListener("click", () => {
          if (currentRevision !== analysisRevision ||
              currentSelection !== revision || controller.signal.aborted) return;
          openSource(source.entryID);
        });
        analysisActions.append(button);
        const traceButton = document.createElement("button");
        traceButton.type = "button";
        traceButton.textContent = "Parcourir la séquence de cette source : " + label;
        traceButton.addEventListener("click", () => {
          if (currentRevision !== analysisRevision ||
              currentSelection !== revision || controller.signal.aborted ||
              token.value !== bearer || caseID.value.trim() !== selectedCase) return;
          // The sequence is a local view of already returned authorized rows.
          // Any opened source is reauthorized via the existing owner routes.
          const tracePageSize = 25;
          function renderTrace(start) {
            if (currentRevision !== analysisRevision ||
                currentSelection !== revision || controller.signal.aborted ||
                token.value !== bearer || caseID.value.trim() !== selectedCase) return;
            traceResults.replaceChildren();
            const end = Math.min(start + tracePageSize, source.rows.length);
            traceResults.textContent = "Séquence de source (ordre des pages autorisées) : " +
              label + " · " + (start + 1) + "–" + end + " / " +
              source.rows.length + ". Lecture non transactionnelle ; ouvrir une entrée revérifie ses droits.";
            const list = document.createElement("ol");
            for (const row of source.rows.slice(start, end)) {
              const item = document.createElement("li");
              const open = document.createElement("button");
              open.type = "button";
              const idLabel = row.id.length > 120 ? row.id.slice(0, 120) + "… (abrégé)" : row.id;
              open.textContent = idLabel + " · " + row.kind + " · " +
                row.time + (row.uncertain ? " (horodatage incertain)" : "");
              open.addEventListener("click", () => {
                if (currentRevision !== analysisRevision ||
                    currentSelection !== revision || controller.signal.aborted ||
                    token.value !== bearer || caseID.value.trim() !== selectedCase) return;
                openSource(row.id);
              });
              item.append(open);
              list.append(item);
            }
            traceResults.append(list);
            if (start > 0) {
              const prevTrace = document.createElement("button");
              prevTrace.type = "button";
              prevTrace.textContent = "Entrées précédentes de cette source";
              prevTrace.addEventListener("click", () => renderTrace(Math.max(0, start - tracePageSize)));
              traceResults.append(prevTrace);
            }
            if (end < source.rows.length) {
              const nextTrace = document.createElement("button");
              nextTrace.type = "button";
              nextTrace.textContent = "Entrées suivantes de cette source";
              nextTrace.addEventListener("click", () => renderTrace(end));
              traceResults.append(nextTrace);
            }
          }
          renderTrace(0);
        });
        traceActions.append(traceButton);
      }
    } catch {
      if (currentRevision === analysisRevision) {
        analysis.textContent = "Synthèse indisponible ou droits modifiés. Aucun total partiel conservé.";
      }
    } finally {
      if (currentRevision === analysisRevision) {
        analysisRequest = null;
        analyze.disabled = false;
      }
    }
  }

  // Compare two fully reauthorized Case timelines, without an extra API,
  // frontend trust decision, owner-source read or unscoped search index.
  async function compareAuthorizedCases() {
    const selectedCase = caseID.value.trim();
    const otherCase = compareCase.value.trim();
    const bearer = token.value;
    const windowSelection = readWindow();
    clearComparison();
    if (!bearer || !selectedCase || !otherCase || selectedCase === otherCase ||
        !windowSelection || otherCase.length > 1024) {
      comparison.textContent = "Choisissez deux Cases différents et un jeton local.";
      return;
    }
    const activeComparison = compareRevision;
    const activeSelection = revision;
    const controller = new AbortController();
    compareRequest = controller;
    compareButton.disabled = true;
    comparison.textContent = "Comparaison des pages autorisées en cours…";
    const pageSize = 200;
    const allowedKinds = new Set(["observed", "inferred", "user-action", "decision", "response"]);
    // Each fetch goes through the existing server-side Case/Timeline/source
    // ACL chain. One inaccessible Case invalidates the ENTIRE comparison.
    async function sourcesFor(id) {
      const groups = new Map();
      const buckets = new Map();
      const events = [];
      const seen = new Set();
      const kinds = { observed: 0, inferred: 0 };
      let certainMin = null;
      let certainMax = null;
      let uncertain = 0;
      let total = null;
      for (let at = 0; at === 0 || at < total; at += pageSize) {
        if (controller.signal.aborted) throw Error("comparison cancelled");
        const url = new URL("/case-timeline", window.location.origin);
        url.searchParams.set("case_id", id);
        url.searchParams.set("limit", String(pageSize));
        url.searchParams.set("offset", String(at));
        if (kind.value) url.searchParams.set("kind", kind.value);
        if (search.value.trim()) url.searchParams.set("q", search.value.trim());
        if (windowSelection.from) url.searchParams.set("from", windowSelection.from);
        if (windowSelection.to) url.searchParams.set("to", windowSelection.to);
        const response = await fetch(url, {
          method: "GET", headers: { Authorization: "Bearer " + bearer },
          cache: "no-store", credentials: "omit", redirect: "error",
          mode: "same-origin", signal: controller.signal
        });
        if (!response.ok) throw Error("Case unavailable");
        const data = await response.json();
        if (controller.signal.aborted || activeComparison !== compareRevision ||
            activeSelection !== revision) throw Error("comparison invalidated");
        if (!data || !data.page || !Array.isArray(data.entries) ||
            !Number.isSafeInteger(data.page.total_visible) ||
            data.page.total_visible < 0 || data.page.total_visible > 5000 ||
            data.page.limit !== pageSize || data.page.offset !== at ||
            !["allow", "filtered", "partial", "uncertain", "empty"].includes(data.status) ||
            (total !== null && total !== data.page.total_visible)) {
          throw Error("inconsistent authorized comparison page");
        }
        total = data.page.total_visible;
        const expected = Math.min(pageSize, Math.max(0, total - at));
        if (data.entries.length !== expected ||
            data.page.has_more !== (at + expected < total)) {
          throw Error("incomplete authorized comparison");
        }
        for (const entry of data.entries) {
          if (!entry || typeof entry.id !== "string" || !entry.id ||
              seen.has(entry.id) || !allowedKinds.has(entry.kind) ||
              typeof entry.source_ref !== "string" || !entry.source_ref ||
              typeof entry.time_uncertain !== "boolean") {
            throw Error("invalid authorized comparison entry");
          }
          seen.add(entry.id);
          // Retain ONLY the projected, source-authorized fields. The optional
          // detailed download is never assembled from raw owner payloads.
          events.push({
            case_id: id, id: entry.id, kind: entry.kind,
            source_time: entry.source_time, timezone: entry.timezone,
            time_uncertain: entry.time_uncertain, source_ref: entry.source_ref,
            provenance: entry.provenance, version: entry.version,
            supersedes: entry.supersedes
          });
          const group = groups.get(entry.source_ref) || { count: 0, firstEntryID: entry.id, rows: [] };
          group.count++;
          group.rows.push({
            caseID: id, id: entry.id, kind: entry.kind,
            time: entry.source_time, uncertain: entry.time_uncertain
          });
          groups.set(entry.source_ref, group);
          if (entry.kind === "observed" || entry.kind === "inferred") kinds[entry.kind]++;
          if (entry.time_uncertain) {
            uncertain++;
          } else {
            // Only declared-certain, valid RFC3339 instants may define a
            // forensic comparison range; uncertain entries remain counted.
            if (typeof entry.source_time !== "string" ||
                !/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d{1,9})?(?:Z|[+-]\d{2}:\d{2})$/.test(entry.source_time)) {
              throw Error("invalid certain comparison timestamp");
            }
            const instant = Date.parse(entry.source_time);
            if (!Number.isFinite(instant)) throw Error("unparseable certain comparison timestamp");
            certainMin = certainMin === null ? instant : Math.min(certainMin, instant);
            certainMax = certainMax === null ? instant : Math.max(certainMax, instant);
            // A shared five-minute interval is a triage clue, not causal proof.
            const bucket = Math.floor(instant / 300000);
            const counts = buckets.get(bucket) || { total: 0, observed: 0, inferred: 0 };
            counts.total++;
            if (entry.kind === "observed" || entry.kind === "inferred") counts[entry.kind]++;
            buckets.set(bucket, counts);
          }
        }
      }
      return { total, groups, buckets, kinds, uncertain, certainMin, certainMax, events };
    }
    try {
      const left = await sourcesFor(selectedCase);
      const right = await sourcesFor(otherCase);
      if (controller.signal.aborted || activeComparison !== compareRevision ||
          activeSelection !== revision) return;
      const common = [...left.groups].filter(([ref]) => right.groups.has(ref))
        .map(([ref, group]) => ({
          ref, left: group.count, right: right.groups.get(ref).count,
          leftEntryID: group.firstEntryID, rightEntryID: right.groups.get(ref).firstEntryID,
          rows: [...group.rows, ...right.groups.get(ref).rows]
        }))
        .sort((a, b) => Math.min(b.left, b.right) - Math.min(a.left, a.right) ||
          (a.ref < b.ref ? -1 : a.ref > b.ref ? 1 : 0));
      const top = common.slice(0, 5).map(group =>
        (group.ref.length > 120 ? group.ref.slice(0, 120) + "… (abrégée)" : group.ref) +
        " : " + group.left + " / " + group.right).join(" · ");
      const range = value => value.certainMin === null ? "aucune plage certaine" :
        new Date(value.certainMin).toISOString() + " → " +
        new Date(value.certainMax).toISOString();
      const overlap = left.certainMin === null || right.certainMin === null
        ? "indéterminable" :
        (left.certainMin <= right.certainMax && right.certainMin <= left.certainMax
          ? "oui" : "non");
      const simultaneous = [...left.buckets].filter(([bucket]) => right.buckets.has(bucket))
        .map(([bucket, leftCount]) => ({ bucket, left: leftCount, right: right.buckets.get(bucket) }))
        .sort((a, b) => Math.min(b.left.total, b.right.total) -
          Math.min(a.left.total, a.right.total) || a.bucket - b.bucket);
      // Complete both Case page traversals before exposing any local receipt.
      // Case/source owner checks are on the existing server endpoints; this
      // report is a manually saved, unsigned summary, not evidence custody.
      const comparisonReport = {
        schema_version: 1,
        type: "cmdr-synthetic-authorized-cross-case-comparison",
        generated_at_utc: new Date().toISOString(),
        case_ids: [selectedCase, otherCase],
        filters: {
          kind: kind.value, q: search.value.trim(),
          from: windowSelection.from, to: windowSelection.to
        },
        cases: [left, right].map(item => ({
          visible_events: item.total,
          observed: item.kinds.observed,
          inferred: item.kinds.inferred,
          uncertain_timestamps: item.uncertain,
          certain_time_range_utc: item.certainMin === null ? null : {
            from: new Date(item.certainMin).toISOString(),
            to: new Date(item.certainMax).toISOString()
          }
        })),
        certain_range_overlap: overlap === "indéterminable" ? null : overlap === "oui",
        common_reference_count: common.length,
        top_common_sources: common.slice(0, 3).map(group => ({
          source_ref: group.ref.slice(0, 120),
          case_1_count: group.left,
          case_2_count: group.right
        })),
        common_five_minute_windows: simultaneous.length,
        top_common_windows_utc: simultaneous.slice(0, 3).map(item => ({
          from: new Date(item.bucket * 300000).toISOString(),
          to: new Date((item.bucket + 1) * 300000).toISOString(),
          case_1_events: item.left.total,
          case_2_events: item.right.total
        })),
        limitations: [
          "Synthetic operator-selected local Cases, not production IAM or live SIEM",
          "Each Case separately authorized by server-side Case/Entry/source checks",
          "Paginated reads are not transactional; permissions and data can change",
          "Overlapping times and source labels are triage hints, not causality or identity proof",
          "Unsigned aggregates; not chain-of-custody evidence, no raw events"
        ]
      };
      preparedComparison = {
        serialized: JSON.stringify(comparisonReport, null, 2) + "\n",
        comparison: activeComparison, selection: activeSelection,
        left: selectedCase, right: otherCase
      };
      downloadComparison.disabled = false;
      // Explicit, SENSITIVE operator-local chronology export. Both Case
      // traversals have already completed and passed pagination/ACL checks.
      // Uncertain instants follow certain instants; no causality is inferred.
      const fullRows = [...left.events, ...right.events];
      const safeRows = fullRows.length === left.total + right.total &&
        fullRows.length <= 10000 && fullRows.every(row =>
          typeof row.id === "string" && row.id.length > 0 &&
          typeof row.source_ref === "string" && row.source_ref.length > 0 &&
          typeof row.provenance === "string" && row.provenance.length > 0 &&
          typeof row.timezone === "string" && row.timezone.length > 0 &&
          typeof row.source_time === "string" &&
          Number.isSafeInteger(row.version) && row.version >= 1 &&
          (row.supersedes === undefined || typeof row.supersedes === "string"));
      if (safeRows) {
        preparedProximity = {
          rows: fullRows, left: selectedCase, right: otherCase,
          comparison: activeComparison, selection: activeSelection, bearer
        };
        proximityButton.disabled = false;
        fullRows.sort((a, b) => {
          if (a.time_uncertain !== b.time_uncertain) return a.time_uncertain ? 1 : -1;
          if (!a.time_uncertain) {
            const leftTime = Date.parse(a.source_time);
            const rightTime = Date.parse(b.source_time);
            if (leftTime !== rightTime) return leftTime - rightTime;
          }
          return a.case_id < b.case_id ? -1 : a.case_id > b.case_id ? 1 :
            a.id < b.id ? -1 : a.id > b.id ? 1 : 0;
        });
        const detailed = JSON.stringify({
          schema_version: 1,
          type: "cmdr-synthetic-authorized-cross-case-chronology",
          generated_at_utc: new Date().toISOString(),
          case_ids: [selectedCase, otherCase],
          filters: comparisonReport.filters,
          certain_entries: fullRows.filter(row => !row.time_uncertain).length,
          uncertain_entries: fullRows.filter(row => row.time_uncertain).length,
          events: fullRows,
          limitations: [
            "Sensitive operator-local export of authorized entry IDs, provenance and source references",
            "Synthetic snapshot and bearer; not a real IAM, live SIEM or custody record",
            "Each page authorized independently; not a transactionally consistent snapshot",
            "Known times precede uncertain times; coincidence does not prove event identity or causation",
            "No raw source-owner payload, remote destination, token or cryptographic signature"
          ]
        }, null, 2) + "\n";
        if (detailed.length <= 16 * 1024 * 1024) {
          preparedCrossEvents = {
            serialized: detailed, comparison: activeComparison,
            selection: activeSelection, left: selectedCase,
            right: otherCase, bearer
          };
          downloadCrossEvents.disabled = false;
        }
      }
      comparison.textContent = "Comparaison des événements autorisés : " +
        left.total + " / " + right.total +
        " · observés : " + left.kinds.observed + " / " + right.kinds.observed +
        " · inférés : " + left.kinds.inferred + " / " + right.kinds.inferred +
        " · horodatages incertains : " + left.uncertain + " / " + right.uncertain +
        " · plage certaine Case 1 (UTC) : " + range(left) +
        " · plage certaine Case 2 (UTC) : " + range(right) +
        " · chevauchement des plages certaines : " + overlap +
        " · fenêtres UTC de 5 min avec événements certains dans les deux Cases : " + simultaneous.length +
        " · références communes : " + common.length + (top ? " · " + top : "") +
        ". Ces plages ne sont pas des preuves d'activité continue, et une référence " +
        "commune n'atteste ni causalité ni identité d'objet. Lectures non transactionnelles.";
      // A shared source is a triage lead, never a client-side permission.
      // The existing owner routes reauthorize Case, Entry and typed source.
      for (const group of common.slice(0, 5)) {
        const label = group.ref.length > 120 ?
          group.ref.slice(0, 120) + "… (abrégée)" : group.ref;
        for (const target of [
          { id: selectedCase, label: "Case 1", entryID: group.leftEntryID },
          { id: otherCase, label: "Case 2", entryID: group.rightEntryID }
        ]) {
          const button = document.createElement("button");
          button.type = "button";
          button.textContent = "Ouvrir un événement de source commune — " +
            target.label + " · " + label;
          button.addEventListener("click", () => {
            if (controller.signal.aborted || compareRevision !== activeComparison ||
                revision !== activeSelection || token.value !== bearer ||
                caseID.value.trim() !== selectedCase ||
                compareCase.value.trim() !== otherCase) return;
            clearComparison();
            if (target.id !== selectedCase) {
              caseID.value = target.id;
              caseID.dispatchEvent(new Event("input", { bubbles: true }));
            }
            openSource(target.entryID);
          });
          comparisonActions.append(button);
        }
        // Provide an actual cross-Case investigative view, not just one
        // representative event. Rows were supplied by the two separately
        // authenticated Timeline page traversals above; source identity is
        // NEVER used as a new server-side authorization capability.
        const trace = document.createElement("button");
        trace.type = "button";
        trace.textContent = "Examiner la séquence croisée de source : " + label;
        trace.addEventListener("click", () => {
          if (controller.signal.aborted || compareRevision !== activeComparison ||
              revision !== activeSelection || token.value !== bearer ||
              caseID.value.trim() !== selectedCase ||
              compareCase.value.trim() !== otherCase) return;
          const chronological = [...group.rows].sort((a, b) => {
            const aTime = a.uncertain ? null : Date.parse(a.time);
            const bTime = b.uncertain ? null : Date.parse(b.time);
            if (aTime === null || bTime === null) {
              if (aTime === null && bTime !== null) return 1;
              if (bTime === null && aTime !== null) return -1;
            } else if (aTime !== bTime) return aTime - bTime;
            return a.caseID < b.caseID ? -1 : a.caseID > b.caseID ? 1 :
              a.id < b.id ? -1 : a.id > b.id ? 1 : 0;
          });
          const pageSize = 25;
          function renderCrossCase(start) {
            if (controller.signal.aborted || compareRevision !== activeComparison ||
                revision !== activeSelection || token.value !== bearer ||
                caseID.value.trim() !== selectedCase ||
                compareCase.value.trim() !== otherCase) return;
            comparisonTraceResults.replaceChildren();
            const end = Math.min(start + pageSize, chronological.length);
            comparisonTraceResults.textContent = "Source commune : " + label +
              " · " + (start + 1) + "–" + end + " / " + chronological.length +
              " événements visibles (deux Cases). Horodatages incertains en fin de liste. " +
              "Lecture non transactionnelle ; toute ouverture revérifie les droits.";
            const list = document.createElement("ol");
            for (const row of chronological.slice(start, end)) {
              const item = document.createElement("li");
              const open = document.createElement("button");
              open.type = "button";
              open.textContent = (row.caseID === selectedCase ? "Case 1" : "Case 2") +
                " · " + row.kind + " · " + String(row.time).slice(0, 64) +
                (row.uncertain ? " (incertain)" : "") +
                " · " + row.id.slice(0, 120);
              open.addEventListener("click", () => {
                if (controller.signal.aborted || compareRevision !== activeComparison ||
                    revision !== activeSelection || token.value !== bearer ||
                    caseID.value.trim() !== selectedCase ||
                    compareCase.value.trim() !== otherCase) return;
                clearComparison();
                if (row.caseID !== selectedCase) {
                  caseID.value = row.caseID;
                  caseID.dispatchEvent(new Event("input", { bubbles: true }));
                }
                openSource(row.id); // Existing server-side Case/Entry/source ACLs.
              });
              item.append(open);
              list.append(item);
            }
            comparisonTraceResults.append(list);
            if (start > 0) {
              const previousPage = document.createElement("button");
              previousPage.type = "button";
              previousPage.textContent = "Événements croisés précédents";
              previousPage.addEventListener("click", () => renderCrossCase(Math.max(0, start - pageSize)));
              comparisonTraceResults.append(previousPage);
            }
            if (end < chronological.length) {
              const nextPage = document.createElement("button");
              nextPage.type = "button";
              nextPage.textContent = "Événements croisés suivants";
              nextPage.addEventListener("click", () => renderCrossCase(end));
              comparisonTraceResults.append(nextPage);
            }
          }
          renderCrossCase(0);
        });
        comparisonTraceActions.append(trace);
      }
      for (const item of simultaneous.slice(0, 5)) {
        const from = new Date(item.bucket * 300000).toISOString();
        const to = new Date((item.bucket + 1) * 300000).toISOString();
        for (const target of [
          { id: selectedCase, label: "Case 1", count: item.left },
          { id: otherCase, label: "Case 2", count: item.right }
        ]) {
          const button = document.createElement("button");
          button.type = "button";
          button.textContent = "Examiner " + target.label + " · " + from +
            " → " + to + " · " + target.count.total + " événement(s) certain(s), " +
            target.count.observed + " observé(s), " + target.count.inferred + " inféré(s)";
          button.addEventListener("click", () => {
            if (controller.signal.aborted || compareRevision !== activeComparison ||
                revision !== activeSelection || token.value !== bearer ||
                caseID.value.trim() !== selectedCase || compareCase.value.trim() !== otherCase) return;
            // A pivot is a UI selector only; fetchPage reauthorizes all records.
            clearComparison();
            caseID.value = target.id;
            offset = 0;
            restoreOffsetOnSubmit = false;
            fromInput.value = localTimestamp(from);
            toInput.value = localTimestamp(to);
            restoredRange = { from, to }; // preserve exact UTC through DST folds
            fetchPage();
          });
          comparisonActions.append(button);
        }
      }
    } catch {
      if (activeComparison === compareRevision) {
        preparedComparison = null;
        preparedCrossEvents = null;
        downloadComparison.disabled = true;
        downloadCrossEvents.disabled = true;
        comparisonActions.replaceChildren();
        comparisonTraceActions.replaceChildren();
        comparisonTraceResults.replaceChildren();
        comparisonTraceResults.textContent = "Aucune séquence croisée sélectionnée.";
        comparison.textContent = "Comparaison indisponible ou droits modifiés. Aucun résultat partiel.";
      }
    } finally {
      if (activeComparison === compareRevision) {
        compareRequest = null;
        compareButton.disabled = false;
      }
    }
  }

  // A return URL carries only a presentation selection, never credentials or
  // authorization. It is not fetched until the analyst enters a fresh bearer.
  function localTimestamp(iso) {
    if (!/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$/.test(iso)) return null;
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
        size !== String(limit) || !/^(0|[1-9]\d*)$/.test(pageOffset) ||
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

  // Reuse BOTH owner-authorized page traversals; do not infer permissions
  // from the browser or create a second source-of-truth endpoint.
  function nearestCaseLeads(rows, firstCase, secondCase, seconds) {
    if (!Number.isSafeInteger(seconds) || seconds < 1 || seconds > 300 ||
        rows.length > 10000) return null;
    const certain = id => rows.filter(e => e.case_id === id && !e.time_uncertain)
      .map(row => ({ row, ms: Date.parse(row.source_time) }))
      .sort((a, b) => a.ms - b.ms || (a.row.id < b.row.id ? -1 : a.row.id > b.row.id ? 1 : 0));
    const left = certain(firstCase), right = certain(secondCase);
    if ([...left, ...right].some(e => !Number.isFinite(e.ms))) return null;
    const leads = [], limitMS = seconds * 1000;
    for (const event of left) {
      let lo = 0, hi = right.length;
      while (lo < hi) {
        const mid = (lo + hi) >>> 1;
        if (right[mid].ms < event.ms) lo = mid + 1;
        else hi = mid;
      }
      let prior = lo - 1;
      if (prior >= 0) {
        const priorTime = right[prior].ms;
        while (prior > 0 && right[prior - 1].ms === priorTime) prior--;
      }
      let best = null;
      for (const at of [prior, lo]) {
        if (at < 0 || at >= right.length) continue;
        const peer = right[at], gap = Math.abs(peer.ms - event.ms);
        if (gap > limitMS) continue;
        if (!best || gap < best.gap || (gap === best.gap && peer.row.id < best.right.id)) {
          best = { left: event.row, right: peer.row, gap };
        }
      }
      if (best) leads.push(best);
    }
    leads.sort((a, b) => a.gap - b.gap ||
      (a.left.id < b.left.id ? -1 : a.left.id > b.left.id ? 1 : 0) ||
      (a.right.id < b.right.id ? -1 : a.right.id > b.right.id ? 1 : 0));
    return { leads: leads.slice(0, 25), leftCount: left.length, rightCount: right.length };
  }

  function showCaseProximity() {
    const snap = preparedProximity;
    proximityResults.replaceChildren();
    if (!snap || snap.comparison !== compareRevision || snap.selection !== revision ||
        snap.bearer !== token.value || caseID.value.trim() !== snap.left ||
        compareCase.value.trim() !== snap.right) {
      proximityResults.textContent = "Comparaison expirée ; rechargez les deux Cases.";
      return;
    }
    const seconds = Number(proximitySeconds.value);
    const result = nearestCaseLeads(snap.rows, snap.left, snap.right, seconds);
    if (!result) {
      proximityResults.textContent = "Fenêtre invalide ; entrez un nombre entier de 1 à 300 secondes.";
      return;
    }
    proximityResults.textContent = result.leads.length + " piste(s) dans " + seconds +
      " s (maximum 25) · événements certains : " + result.leftCount +
      " / " + result.rightCount +
      ". Proximité temporelle non exhaustive, sans causalité démontrée. Toute ouverture revérifie les droits.";
    const list = document.createElement("ol");
    for (const lead of result.leads) {
      const item = document.createElement("li");
      item.textContent = "Écart " + lead.gap + " ms · " +
        lead.left.kind + " / " + lead.right.kind + " · ";
      for (const target of [
        { label: "Case 1", case: snap.left, entry: lead.left },
        { label: "Case 2", case: snap.right, entry: lead.right }
      ]) {
        const button = document.createElement("button");
        button.type = "button";
        button.textContent = target.label + " : " + target.entry.id.slice(0, 120);
        button.addEventListener("click", () => {
          if (snap.comparison !== compareRevision || snap.selection !== revision ||
              snap.bearer !== token.value || caseID.value.trim() !== snap.left ||
              compareCase.value.trim() !== snap.right) return;
          clearComparison();
          if (caseID.value.trim() !== target.case) {
            caseID.value = target.case;
            caseID.dispatchEvent(new Event("input", { bubbles: true }));
          }
          openSource(target.entry.id);
        });
        item.append(button);
      }
      list.append(item);
    }
    proximityResults.append(list);
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
      // Opening the synthetic source is a SEPARATE HTTP operation:
      // the owner reauthenticates the Case/Entry/typed source instead of
      // treating the navigation draft as a bearer capability.
      const sourceURL = new URL("/demo-source-open", window.location.origin);
      sourceURL.searchParams.set("case_id", selectedCase);
      sourceURL.searchParams.set("entry_id", entryID);
      const openedResponse = await fetch(sourceURL, {
        method: "GET",
        headers: { Authorization: "Bearer " + bearer },
        credentials: "omit",
        cache: "no-store",
        redirect: "error",
        mode: "same-origin",
        signal: controller.signal
      });
      if (!openedResponse.ok) throw new Error("Source owner access denied");
      const owner = await openedResponse.json();
      if (currentRevision !== sourceRevision) return;
      if (!owner || owner.status !== "ready" || owner.synthetic !== true ||
          owner.case_id !== selectedCase || owner.entry_id !== entryID ||
          owner.kind !== data.kind || owner.object_id !== data.object_id ||
          typeof owner.preview !== "string" || owner.preview.length > 2048) {
        throw new Error("Source owner response differs from navigation");
      }
      const localReturn = new URL("/analyst", window.location.origin);
      localReturn.search = expected.toString();
      returnLinkURL = localReturn.href;
      copyReturnLink.disabled = false;
      sourceDetail.replaceChildren();
      addLine(sourceDetail, "Objet source autorisé :", data.kind + " / " + data.object_id);
      addLine(sourceDetail, "Entrée d’origine :", data.return_entry_id);
      addLine(sourceDetail, "Aperçu propriétaire synthétique :", owner.preview || "Aucun aperçu disponible");
      addLine(sourceDetail, "Ouverture :", "source vérifiée par le propriétaire à l’ouverture ; démonstration locale, aucune donnée réelle.");
      sourceDetail.hidden = false;
      returnToCase.disabled = false;
      sourceStatus.textContent = "Objet synthétique ouvert après une nouvelle vérification côté serveur.";
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
    clearAnalysis();
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
  analyze.addEventListener("click", analyzeAuthorizedCase);
  compareButton.addEventListener("click", compareAuthorizedCases);
  proximityButton.addEventListener("click", showCaseProximity);
  proximitySeconds.addEventListener("input", () => {
    proximityResults.replaceChildren();
    proximityResults.textContent = "Fenêtre modifiée ; relancez les pistes.";
  });
  compareCase.addEventListener("input", clearComparison);
  downloadComparison.addEventListener("click", () => {
    const receipt = preparedComparison;
    if (!receipt || receipt.comparison !== compareRevision ||
        receipt.selection !== revision || !token.value ||
        caseID.value.trim() !== receipt.left ||
        compareCase.value.trim() !== receipt.right) return;
    let objectURL = null;
    try {
      const blob = new Blob([receipt.serialized], { type: "application/json" });
      objectURL = URL.createObjectURL(blob);
      const link = document.createElement("a");
      link.href = objectURL;
      link.download = "cmdr-local-cross-case-comparison.json";
      link.click();
      setTimeout(() => URL.revokeObjectURL(objectURL), 1000);
    } catch {
      if (objectURL) URL.revokeObjectURL(objectURL);
      comparison.textContent = "Export local indisponible. Recomparez les Cases si nécessaire.";
    }
  });
  downloadCrossEvents.addEventListener("click", () => {
    const receipt = preparedCrossEvents;
    if (!receipt || receipt.comparison !== compareRevision ||
        receipt.selection !== revision || token.value !== receipt.bearer ||
        caseID.value.trim() !== receipt.left ||
        compareCase.value.trim() !== receipt.right) return;
    let objectURL = null;
    try {
      const blob = new Blob([receipt.serialized], { type: "application/json" });
      if (blob.size > 16 * 1024 * 1024) throw Error("export exceeds size bound");
      objectURL = URL.createObjectURL(blob);
      const link = document.createElement("a");
      link.href = objectURL;
      link.download = "cmdr-local-cross-case-chronology-sensitive.json";
      link.click();
      setTimeout(() => URL.revokeObjectURL(objectURL), 1000);
    } catch {
      if (objectURL) URL.revokeObjectURL(objectURL);
      comparison.textContent = "Export chronologique indisponible. Recomparez les Cases.";
    }
  });
  downloadReport.addEventListener("click", () => {
    const receipt = preparedReport;
    if (!receipt || receipt.revision !== analysisRevision ||
        receipt.selection !== revision || !token.value ||
        caseID.value.trim() !== receipt.caseID) return;
    let objectURL = null;
    try {
      const blob = new Blob([receipt.serialized], { type: "application/json" });
      objectURL = URL.createObjectURL(blob);
      const link = document.createElement("a");
      link.href = objectURL;
      link.download = "cmdr-local-case-summary.json";
      link.click();
      // Download starts asynchronously in browsers: do not revoke before it
      // consumes the Blob URL. No remote endpoint or credential is involved.
      setTimeout(() => URL.revokeObjectURL(objectURL), 1000);
    } catch {
      if (objectURL) URL.revokeObjectURL(objectURL);
      analysis.textContent = "Export local indisponible. Relancez l’analyse si nécessaire.";
    }
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
      clearAnalysis();
      clearComparison();
      clearWazuh();
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
