# E10-INV-008B-COMPARE — bounded independent run-comparison handoff

## Authority and ownership

- **Canonical owner:** `CAP-INV-008 — Search and Hunt Provenance` (Investigate/Signals and Hunt). Canonical section 12 explicitly allows Class-0 **Compare runs**; depends on Shared-owned `CAP-INV-002` Search Job and `CAP-INV-004` Telemetry/Event Inspector. The action compares **two already completed, independently authorized runs**, not traces from an invented store.
- **Engineering:** `E10-INV-008B-COMPARE`, subordinate to existing epic `E10-INV-008`. The code uses implemented but **not yet VERIFIED** `E10-INV-008A-PROVENANCE` on parent PR #131. Graph dependency threshold is `IMPLEMENTED` **for development only**, NEVER a permission to merge with required parent check red.
- **Temporary code owner:** the existing `product-runtime/investigation-timeline/cmd/case-timeline-demo` browser of the local synthetic owner. Production Investigate UI and Shared Trace ownership are not provided.
- **Git parent:** frozen PR #131 `feature/e10-inv-008a-readonly-searchjob-provenance` SHA `19ec74c6d7cf1f4ce4d8087d03b290f39b6f5d08`, not `main` and not the business parent. Preserve all existing L2 SHA references.

## Behavior and exact exclusions

Operator runs one genuine `POST /demo-local-search`, waits for a completed server owner-authorized first-party Search Job (membership and Wazuh/generic lead checks unchanged), clicks **Mémoriser ce run** to pin only bounded metadata, changes search criteria *without changing Case or bearer*, executes a second owner-authorized Search Job and sees a descriptive **difference**: QueryVersion, UTC windows, authorized result counts, declared-source count, failed-source count and audit-event count. The comparison never asks for a new endpoint or retrieves another source. Unknown source/failure metadata is **indisponible**, not zero. Same RunID is explicitly not presented as two separate runs. These are comparison **counts**, not a temporal causal relationship or proof of investigation completeness.

In-memory baseline contains only run identifier, sanitized QueryVersion, UTC bounds, result/audit/source/failure counts, selected Case and numeric scope revision. It does **not** contain the bearer, source labels, event IDs, raw preview, full audit events, workspace selection or serializable Trace. Any Case/bearer or unrelated Timeline control change destroys the baseline and current comparison. An invalid/foreign second Search Job destroys both and displays no partial comparison; all source-clicks still require the existing first-party Case/Timeline/typed-source reauthorization. Browser DOM insertion remains `textContent` only, not HTML injection.

## Proof and gates

Target exact-HEAD GitHub Actions: `investigation-timeline-runtime` with Linux Go/Node browser, macOS and Windows; `engineering-kernel` mandatory changed-security critical >=90% (never waive or adjust). Browser acceptance validates a pair of independently server-authorized completed mocked API responses, per-run Query version and counts, no extra fetch, no labels/refs/bearer, unknown counts explicit, reset on bearer rotation and denial after prior pinned run. Existing source navigation and return-context browser harnesses must pass.

**Expected blocker:** frozen parent PR #131 SHA `19ec74c6d7cf1f4ce4d8087d03b290f39b6f5d08` has exact-head Go security gate **89.30% < mandatory 90%** despite the functional runtime success. This browser-only child does not resolve the earlier Go gate and must not be merged even if its own browser tests pass. #36 foundation and parent chain remain in draft/unmerged.

**Not delivered:** Persistent Shared Trace or Activity Stream, query history across sessions, audit authors, transformations, Studio Tool Calls, redaction/retention, annotation or Case relationship mutations, Export Job, production IAM, forensic custody. Canonical `OPEN-015` and `OPEN-013` remain open; complete `CAP-INV-008` is **not** accepted. The next *integration* priority is to unblock and consolidate parent PRs, not add endless UI features.

## Targeted L2 risk

One independently scheduled, read-only L2 only once this material capability is published to a frozen HEAD and exact CI evidence is known. Focus on same-run reuse, Case/bearer change and abort race, source counts privacy, HTML injection, fake zero for missing metadata and regression of CAP-INV-007 Inspector round-trip. This is not a native GitHub independent account APPROVED review.
