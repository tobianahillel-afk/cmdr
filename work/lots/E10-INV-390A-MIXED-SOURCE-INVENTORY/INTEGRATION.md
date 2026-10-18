# E10-INV-390A — Source-aware offline Case inventory

## Functional integration

PR #78 already builds and saves ONE synthetic Case containing native Wazuh 4.x alerts and generic normalized JSONL telemetry, then supports authenticated Case/Timeline/source reopening, one first-party Event Search Job, Wazuh-only rule/agent pivots and burst leads. The offline saved-Case verifier `-verify-import-snapshot-case` previously returned ONLY a total and `authorized_events_by_import_slot`. The existing provenance is an operator-selected import namespace, **not an authenticated source-family identifier**. In the current one-step mixed pathway Wazuh events retain `operator-local-jsonl` while generic inputs use `operator-local-jsonl-feed-1`; in other single-source paths generic input can also use `operator-local-jsonl`. The provenance field alone therefore cannot consistently classify evidence type across supported workflows.

This small functional integration improves the EXISTING offline verification receipt, not a new verifier, query engine, parser or API. After exact SHA-256 byte verification and existing Case/principal/tenant/environment/typed-source ACL projection via `loadAuthorizedImportedCase`, count only canonical already-authorized imported events by owner-linked source family:
- `authorized_wazuh_alerts`: typed sources carrying the already-validated synthetic `MemoryTimelineSource.Wazuh` metadata,
- `authorized_generic_events`: remaining authorized imported typed sources.

The two new aggregate numbers must sum to the existing `authorized_imported_events`. The existing `authorized_events_by_import_slot` map stays unchanged for backward-compatible readers, and its key is explicitly a *namespace/ordinal* and NOT proof of a unique evidence tool. Unknown/missing typed source or incomplete traversal refuses the whole verification receipt. No raw Wazuh descriptions, rules, agent IDs, source refs, generic messages, bearer tokens or additional tenants appear in the output. This is an operator-owned synthetic lab owner, NOT an authenticated origin classifier, live SIEM, real tenant IAM or signed evidence custody.

### Integration acceptance

- Real one-step `writePersistedMixedCase` from two Wazuh alerts plus one generic JSONL entry (different legacy provenance namespaces), persisted once, then actual CLI `-verify-import-snapshot-case CASE -expect-snapshot-sha256 SHA` returns 3 total, 2 Wazuh alerts and 1 generic event without a listener, source preview or bearer.
- In a valid previously saved owner graph where one Wazuh typed source has its readers revoked, independently compute a new operator-supplied file digest and verify: the inventory contains 1 Wazuh and 1 generic, with NO metadata about the revoked source.
- All existing Wazuh/generic authorization, Event Search Job, burst and no-overwrite persistence tests still run in Linux, macOS and Windows CI.
- Keep global and changed-security-critical 90% floor unchanged; exact-head green functional tests are necessary but not sufficient for merge.
- Base PR #78 immutable HEAD `b3a83a1b8c74069664778d23341223a78317e717` stays untouched for its already requested L2. Do not open a redundant standalone L2 for a trivial test delta; one targeted independent review if material data classification/ACL risk is raised after CI and the parent L2 findings are accessible. Do not pretend this enables the massive upstream foundation PR #36 to merge.
