# E10-INV-120A — export a bounded, operator-selected cross-Case triage summary

## Functional integration

The existing PR #41 analyst viewer already executes a real operator-selected, loopback Case discovery → Timeline → source preview journey and compares two independently authorized Case timelines. The comparison previously existed only as on-screen text. This increment makes the **already completed, fully authorized** cross-Case comparison an explicitly downloadable local JSON summary.

The operator must complete both authorized paginated Case reads and click the new `Enregistrer la comparaison JSON locale` button. Only aggregate counts and at most five common source-reference strings (each <=96 characters, to match the existing search bound), with certain-time envelopes and uncertain counts, enter the file. Long/shared but unexported references are reflected in a count, not silently truncated. No raw event rows, typed-owner preview, bearer, destination URL or source outside authorized Case pages is written.

The export clears on new Case, compare Case, bearer or filter edits; the click guard verifies comparison and selection revisions and the exact in-memory current bearer. Failed authorization, inconsistent pages and cancelled reads never enable a download. This is a client-side local file operation, **not an upload endpoint** or a new authorization decision. The existing backend still rechecks Case, Timeline Entry and typed-source ACLs per read. An earlier successful read cannot prove that a subsequent change of source authority never happened; the document accurately says that the report is a historical nontransactional point-in-time aggregate.

## Reuse-first decision

Reuse `analyst.js`, its existing `/case-timeline` paginated first-party Go authority chain, `Blob` download path already used for one-Case reports and the existing Node VM integration tests. Do **not** add a vendor report engine, indexer, JavaScript library, persistence table, service/API or another redundant authorization layer.

## Tests, risks and proof

The `analyst_summary.test.cjs` behavioral tests verify the successful authenticated two-Case read, exact aggregate JSON, absence of bearer and raw event identifiers, revocation/denial and changed-selection invalidation, and exclusion of overly long references. Existing Go tests separately verify the real operator-selected two-tenant file → loopback bearer → source preview path.

The report is **not** a forensically signed export, live SIEM ingestion, real SSO, validated production ACLs, transactional Case consistency, audit trail or durable evidence store. Treat it as local synthetic triage only. `OPEN-013` and all canonical product decisions remain unchanged. This milestone must be evaluated on its exact commit and independently reviewed before claims of acceptance; PR #41 remains draft stacked on PR #36.
