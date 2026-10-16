# E10-INV-120A — authenticated cross-Case event chronology

## Functional increment

PR #41 now offers an **in-page, read-only merged timeline** of entries sharing the most relevant (max 5, max 96 characters) exact source references between two fully authorized Cases. Previously the analyst could see counts and broad time envelopes or switch Case through a pivot but could not inspect event interleaving. The new chronology lists up to 20 individually identified, already authorized events, with Case A/B, observed/inferred kind, shared reference and normalized UTC timestamp. Uncertain-time events appear last and are explicitly labeled, never assigned a fabricated UTC instant. A more precise drill-down reuses the already implemented Case/source pivot, which reauthorizes on every HTTP request.

The existing `sourcesFor` page loop now retains only minimal visible metadata after the existing Case + Timeline Entry + source authorization; comparison requires both complete paginated reads, consistent totals and schema before rendering. A failed or revoked comparison clears the entire list. Input/bearer changes invalidate it immediately. The server continues to own access decisions. No new endpoint, index, database, provider, background job, paid API, framework or state machine. At most 5,000 entries per Case and 20 visible merged records, using existing in-memory maps/arrays and native JavaScript sort.

## Choice and limitations

**Chosen:** reuse `/case-timeline`, the existing owner-backed authorization chain, source grouping and the local analyst DOM. **Rejected:** a new search-index service, cross-tenant query, general correlation engine or persistent evidence trail before production IAM/owner storage exists.

This is strictly a synthetic operator-selected, loopback-local aid. Equal `source_ref` is a correlation hint, not proof of an identical canonical object, causal relation, real ordering among uncertain entries, transactional snapshot, chain of custody, or production SIEM integration. The local JSON export remains aggregate-only; event IDs in the UI are **not** exported. `OPEN-013` is still unresolved.

Acceptance: Node VM integration tests verify merged UTC order, uncertain-last, top shared refs, source visibility, denial and token invalidation. Existing Go integration tests separately exercise actual loopback selected Case → Timeline → typed-source preview, but a browser E2E test and production provider are outstanding. Risk-based independent L2 required on stable exact SHA; no acceptance inferred from green CI alone.
