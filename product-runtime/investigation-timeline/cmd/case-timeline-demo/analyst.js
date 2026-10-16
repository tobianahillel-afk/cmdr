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

  const analyze = document.getElementById("analyze");
  const analysis = document.getElementById("analysis");
  const analysisActions = document.getElementById("analysis-actions");
  const compareCase = document.getElementById("compare-case");
  const compareButton = document.getElementById("compare");
  const comparison = document.getElementById("comparison");
  const comparisonActions = document.getElementById("comparison-actions");
  const comparisonChronology = document.getElementById("comparison-chronology");
  const downloadComparison = document.getElementById("download-comparison");
  let preparedComparison = null;
  downloadComparison.disabled = true;
  let compareRevision = 0;
  let compareRequest = null;
  function clearComparison() {
    ++compareRevision;
    if (compareRequest) compareRequest.abort();
    compareRequest = null;
    compareButton.disabled = false;
    preparedComparison = null;
    downloadComparison.disabled = true;
    comparison.textContent = "Aucune comparaison calculée.";
    comparisonActions.replaceChildren();
    comparisonChronology.replaceChildren();
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
            { ref: entry.source_ref, entryID: entry.id, count: 0, observed: 0, inferred: 0 };
          source.count++;
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
      const timeRanges = new Map();
      const visibleRows = [];
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
          groups.set(entry.source_ref, (groups.get(entry.source_ref) || 0) + 1);
          if (entry.kind === "observed" || entry.kind === "inferred") kinds[entry.kind]++;
          // Derive a temporal envelope only from timestamps the authorized
          // Timeline projection marks certain. Unknown time is never guessed.
          const range = timeRanges.get(entry.source_ref) ||
            { first: null, last: null, certain: 0, uncertain: 0 };
          let instant = null;
          if (entry.time_uncertain) {
            range.uncertain++;
            uncertain++;
          } else {
            if (typeof entry.source_time !== "string" ||
                !/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d{1,9})?(?:Z|[+-]\d{2}:\d{2})$/.test(entry.source_time)) {
              throw Error("invalid certain comparison timestamp");
            }
            instant = Date.parse(entry.source_time);
            if (!Number.isFinite(instant)) throw Error("invalid certain comparison instant");
            range.certain++;
            certainMin = certainMin === null ? instant : Math.min(certainMin, instant);
            certainMax = certainMax === null ? instant : Math.max(certainMax, instant);
            range.first = range.first === null ? instant : Math.min(range.first, instant);
            range.last = range.last === null ? instant : Math.max(range.last, instant);
          }
          timeRanges.set(entry.source_ref, range);
          // Retain only the already authorized response fields needed for
          // a read-only, bounded chronology; never cache owner previews.
          visibleRows.push({
            id: entry.id, kind: entry.kind, ref: entry.source_ref,
            instant, uncertain: entry.time_uncertain
          });
        }
      }
      return { total, groups, timeRanges, kinds, certainMin, certainMax, uncertain, visibleRows };
    }
    try {
      const left = await sourcesFor(selectedCase);
      const right = await sourcesFor(otherCase);
      if (controller.signal.aborted || activeComparison !== compareRevision ||
          activeSelection !== revision) return;
      const common = [...left.groups].filter(([ref]) => right.groups.has(ref))
        .map(([ref, count]) => ({
          ref, left: count, right: right.groups.get(ref),
          leftTime: left.timeRanges.get(ref), rightTime: right.timeRanges.get(ref)
        }))
        .sort((a, b) => Math.min(b.left, b.right) - Math.min(a.left, a.right) ||
          (a.ref < b.ref ? -1 : a.ref > b.ref ? 1 : 0));
      // Interpret only fully authorized, certain timestamp envelopes.
      // The relation describes *visible ranges*, never causality or a
      // shared underlying object; uncertain events cannot be ordered.
      const temporalRelation = group => {
        const a = group.leftTime;
        const b = group.rightTime;
        if (!a.certain || !b.certain) {
          return { relation: "indeterminate", gapMS: null,
            label: "ordre temporel indéterminé" };
        }
        if (a.last < b.first) {
          const gapMS = b.first - a.last;
          return { relation: "case_a_before_case_b", gapMS,
            label: "plage certaine A avant B · écart entre plages : " +
              (gapMS / 1000).toFixed(3) + " s" };
        }
        if (b.last < a.first) {
          const gapMS = a.first - b.last;
          return { relation: "case_b_before_case_a", gapMS,
            label: "plage certaine B avant A · écart entre plages : " +
              (gapMS / 1000).toFixed(3) + " s" };
        }
        return { relation: "certain_envelopes_overlap", gapMS: 0,
          label: "plages certaines recoupées (pas nécessairement d'événements simultanés)" };
      };
      const rangeLabel = range => range.certain
        ? new Date(range.first).toISOString() + " → " +
          new Date(range.last).toISOString() + " (" + range.certain + " certain(s), " +
          range.uncertain + " incertain(s))"
        : "aucun horaire certain (" + range.uncertain + " incertain(s))";
      const top = common.slice(0, 5).map(group => {
        const temporalGap = group.leftTime.certain && group.rightTime.certain
          ? " · écart des premiers instants certains : " +
            (Math.abs(group.leftTime.first - group.rightTime.first) / 1000).toFixed(3) + " s"
          : " · écart temporel indéterminé";
        return (group.ref.length > 120 ? group.ref.slice(0, 120) + "… (abrégée)" : group.ref) +
          " : " + group.left + " / " + group.right +
          " · A : " + rangeLabel(group.leftTime) +
          " · B : " + rangeLabel(group.rightTime) + temporalGap +
          " · " + temporalRelation(group).label;
      }).join(" · ");
      const caseRange = value => value.certainMin === null ?
        "aucune plage certaine" : new Date(value.certainMin).toISOString() +
        " → " + new Date(value.certainMax).toISOString();
      const envelopeOverlap = left.certainMin === null || right.certainMin === null
        ? null : left.certainMin <= right.certainMax && right.certainMin <= left.certainMax;
      comparison.textContent = "Comparaison des événements autorisés : " +
        left.total + " / " + right.total +
        " · observés : " + left.kinds.observed + " / " + right.kinds.observed +
        " · inférés : " + left.kinds.inferred + " / " + right.kinds.inferred +
        " · horodatages incertains : " + left.uncertain + " / " + right.uncertain +
        " · plage certaine Case 1 (UTC) : " + caseRange(left) +
        " · plage certaine Case 2 (UTC) : " + caseRange(right) +
        " · chevauchement des plages certaines : " +
        (envelopeOverlap === null ? "indéterminable" : envelopeOverlap ? "oui" : "non") +
        " · références communes : " + common.length + (top ? " · " + top : "") +
        ". Les écarts reposent sur les seuls horodatages certains et visibles. " +
        "Présence commune d'une référence exacte, pas une preuve de causalité, " +
        "d'identité d'objet ou d'exhaustivité. Lectures non transactionnelles.";
      // An explicit local report contains only the already authorized
      // aggregate: no raw event rows, token, owner preview or URL. Exclude
      // references beyond the search bound instead of silently truncating.
      const exportable = common.filter(group => group.ref.length <= 96).slice(0, 5);
      const safeEnvelope = range => ({
        certain_count: range.certain,
        uncertain_count: range.uncertain,
        first_certain_utc: range.certain ? new Date(range.first).toISOString() : null,
        last_certain_utc: range.certain ? new Date(range.last).toISOString() : null
      });
      const receipt = {
        schema_version: 1,
        type: "cmdr-synthetic-authorized-case-comparison",
        generated_at_utc: new Date().toISOString(),
        case_ids: [selectedCase, otherCase],
        filters: {
          kind: kind.value, q: search.value.trim(),
          from: windowSelection.from, to: windowSelection.to
        },
        visible_event_counts: [left.total, right.total],
        visible_observed_counts: [left.kinds.observed, right.kinds.observed],
        visible_inferred_counts: [left.kinds.inferred, right.kinds.inferred],
        uncertain_timestamp_counts: [left.uncertain, right.uncertain],
        certain_case_envelopes_utc: [left, right].map(value => ({
          first: value.certainMin === null ? null : new Date(value.certainMin).toISOString(),
          last: value.certainMax === null ? null : new Date(value.certainMax).toISOString()
        })),
        certain_case_envelopes_overlap: envelopeOverlap,
        shared_source_references: common.length,
        exported_source_references: exportable.map(group => ({
          source_ref: group.ref,
          visible_counts: [group.left, group.right],
          case_a_time: safeEnvelope(group.leftTime),
          case_b_time: safeEnvelope(group.rightTime),
          temporal_relation: temporalRelation(group).relation,
          certain_envelope_gap_ms: temporalRelation(group).gapMS
        })),
        excluded_or_unlisted_references: common.length - exportable.length,
        limitations: [
          "Local synthetic offline owner snapshot; no production identity or real SIEM source",
          "Both Cases were authorized independently when pages were read",
          "Comparisons are nontransactional and may become stale after rights or data change",
          "Shared source reference is a triage indicator, not proof of causation or common object",
          "Temporal ordering compares certain visible envelopes only; uncertain events may change real ordering",
          "Case-wide observed/inferred counts and envelopes include authorized events regardless of shared sources; ranges do not prove continuous activity",
          "Unsigned local aggregate; no raw events, forensic custody proof or complete export"
        ]
      };
      preparedComparison = {
        serialized: JSON.stringify(receipt, null, 2) + "\n",
        revision: activeComparison, selection: activeSelection,
        caseA: selectedCase, caseB: otherCase, bearer
      };
      downloadComparison.disabled = false;
      // Correlate only common, bounded source references across BOTH fully
      // authorized Case reads. No new query/ACL, source preview, or storage.
      const shared = new Set(exportable.map(group => group.ref));
      const merged = [
        ...left.visibleRows.filter(row => shared.has(row.ref)).map(row => ({ ...row, side: "A" })),
        ...right.visibleRows.filter(row => shared.has(row.ref)).map(row => ({ ...row, side: "B" }))
      ].sort((a, b) => {
        if (a.uncertain !== b.uncertain) return a.uncertain ? 1 : -1;
        if (!a.uncertain && a.instant !== b.instant) return a.instant - b.instant;
        return a.side < b.side ? -1 : a.side > b.side ? 1 :
          a.id < b.id ? -1 : a.id > b.id ? 1 : 0;
      });
      const heading = document.createElement("p");
      heading.textContent = "Chronologie corrélée des références communes autorisées : " +
        Math.min(merged.length, 20) + " / " + merged.length +
        " événement(s), horodatages incertains à la fin. " +
        "Ordre de triage uniquement, ni causalité ni preuve de conservation.";
      comparisonChronology.append(heading);
      const list = document.createElement("ol");
      for (const row of merged.slice(0, 20)) {
        const item = document.createElement("li");
        item.textContent = "Case " + row.side + " · " + row.id + " · " +
          row.kind + " · " + row.ref + " · " +
          (row.uncertain ? "heure incertaine" : new Date(row.instant).toISOString());
        list.append(item); // textContent prevents operator-controlled HTML injection
      }
      comparisonChronology.append(list);
      // A correlation is a triage hint, NOT an authority. Pivot into the
      // existing Case search; fetchPage performs a fresh server-side Case,
      // Timeline Entry and source authorization before showing any result.
      // Limit source references to the existing 96-character search field.
      for (const group of common.slice(0, 5)) {
        if (group.ref.length > 96) continue;
        for (const targetCase of [selectedCase, otherCase]) {
          const button = document.createElement("button");
          button.type = "button";
          button.textContent = "Examiner dans " +
            (targetCase === selectedCase ? "Case A" : "Case B") +
            " : " + group.ref;
          button.addEventListener("click", () => {
            if (activeComparison !== compareRevision ||
                activeSelection !== revision || controller.signal.aborted ||
                token.value !== bearer || caseID.value.trim() !== selectedCase ||
                compareCase.value.trim() !== otherCase) return;
            caseID.value = targetCase;
            search.value = group.ref;
            // The input listener invalidates the aggregate, aborts old
            // requests, clears owner previews and resets pagination.
            caseID.dispatchEvent(new Event("input", { bubbles: true }));
            fetchPage();
          });
          comparisonActions.append(button);
        }
      }
    } catch {
      if (activeComparison === compareRevision) {
        comparison.textContent = "Comparaison indisponible ou droits modifiés. Aucun résultat partiel.";
        comparisonActions.replaceChildren();
        comparisonChronology.replaceChildren();
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
  compareCase.addEventListener("input", clearComparison);
  downloadComparison.addEventListener("click", () => {
    const receipt = preparedComparison;
    if (!receipt || receipt.revision !== compareRevision ||
        receipt.selection !== revision || token.value !== receipt.bearer ||
        caseID.value.trim() !== receipt.caseA ||
        compareCase.value.trim() !== receipt.caseB) return;
    let objectURL = null;
    try {
      objectURL = URL.createObjectURL(new Blob([receipt.serialized], { type: "application/json" }));
      const link = document.createElement("a");
      link.href = objectURL;
      link.download = "cmdr-local-case-comparison.json";
      link.click();
      setTimeout(() => URL.revokeObjectURL(objectURL), 1000);
    } catch {
      if (objectURL) URL.revokeObjectURL(objectURL);
      comparison.textContent = "Export local indisponible. Relancez la comparaison si nécessaire.";
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
