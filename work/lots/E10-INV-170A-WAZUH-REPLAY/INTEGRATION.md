# E10-INV-170A — Local reproducible Wazuh search replay

## Functional result

The authenticated analyst Wazuh Search Job (PR #55) now has an explicit **download and reload query JSON** workflow. After a completed authorized first-party Event Search Job, the operator can save a small JSON file containing:

- format `cmdr-local-wazuh4-query-replay-v1`, synthetic marker;
- Case ID, exact Wazuh selector (rule ID / agent ID), validated numeric value;
- exact UTC lower-inclusive/upper-exclusive search instants;
- original Search Job ID and number of matches.

**No bearer, source payload, arbitrary DSL, raw logs, source previews or hit IDs are stored.** This is an operator-managed local replay receipt, not a persistent CMDR Saved Search, not custody evidence, and not a database. The existing `saved-query-assets` and `case-queue` contracts explicitly leave mutation/persistence unresolved under OPEN-013; those APIs are neither reimplemented nor falsely represented as completed.

## Reopen security

A local browser file input reads exactly one max 8192-byte JSON. A saved Case **must equal the currently selected Case**; the operator must supply the current ephemeral bearer. Import never calls an HTTP endpoint, never restores historical hit identities, never automatically executes, and never treats a stored job ID or match count as an access credential.

Only bounded valid rule/agent selectors, ID values, an original Job identifier, small nonnegative match count, and exact UTC ISO instants in a strictly increasing interval are accepted. Conversion to a native local date/time input must convert **back to the identical UTC instant**; DST gaps/ambiguous local times reject, never silently alter a forensic interval. On successful import only current form **criteria** are filled. The analyst must explicitly click **Execute Search Job** to trigger the existing server-side Case/Timeline/typed-source authorization and current data search. Changing Case/token or criteria invalidates any previous export and search results; async file reads check that the selection remains unchanged.

The downloaded file carries its own limits and warnings: manually saved JSON metadata is sensitive, editable, unsigned, without tamper attestation, server-side storage, retention, recovery or true production provenance. It must not be used for real clients in this synthetic environment.

## Technical scope and decisions

- First-party browser Blob/File/Web APIs only; no new Go backend, persistence engine, external SaaS, third-party package, service worker, browser storage or paid API.
- Reuse already-live `/demo-wazuh-search` and Event Search orchestrator, existing trusted owner for reruns.
- Browser tests verify real successful Search Job -> downloaded JSON absence of secrets -> edited query invalidates earlier results -> re-import does not call network -> user action reruns and reauthenticates. Negative tests cover foreign Case, malformed ID/field/UTC/range, oversize file and untrusted old result claim.
- Keep PR draft stacked on PR #55 exact HEAD `e9b6dd12295d4aedfdbd04acdbee27bdba59b7d6` and keep its independent L2 SHA intact. Upstream engineering-kernel changed security coverage **86.02% < required 90%** is a real merge-blocking failure; do not waive.

No mutation of canonical product specifications or of unresolved OPEN decisions; no merge. Independent L2 only after exact SHA CI is verifiable.
