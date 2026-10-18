# E10-INV-570A — Single-pass authorized Wazuh cross-Case pair inspection

## Consolidated delivery, not a new engine

The existing operator-local DFIR journey (PR #96) correlates two independently SHA-256-verified Wazuh Cases into deterministic same-rule/same-agent <=60s event pairs and can inspect each original source separately. PR #97 added a complete private **two-source** JSON report for a selected pair, but repeated independent snapshot loading and owner-authorized Timeline traversal for the comparison and again for the actual inspection, plus a JSON marshal/unmarshal roundtrip to recover the correlation result.

This integration **reuses** the existing complete private paired-source capability while simplifying its critical path:

1. Independently load and SHA-256 verify each bounded original snapshot **once** using existing `loadAuthorizedImportedCase`. This returns an already minimized principal/tenant/environment/Case owner and the canonical imported Entry IDs. No owner merge or ACL inheritance.
2. Invoke one shared `compareLoadedAuthorizedSavedCases` in-memory comparison implementation, also called by the ordinary authorized two-Case comparison CLI. It checks matching tenant/environment, independently enumerates complete first-party Case/Timeline/typed-source owner authorization and produces the existing 5-minute cohorts and 60-second Wazuh pairs.
3. Select the requested 1..10 pair, ensure both canonical imported IDs belong to their respective Cases, then invoke the existing first-party `inspectDemoEntryWithData` independently for each original owner to recheck typed-source read access. Both previews are serialized only into the EXISTING private bounded 0600, fsync, atomic no-overwrite publisher after BOTH source opens succeed. No source preview, bearer, partial receipt or listener is emitted on stdout.
4. Preserve original packet and receipt schemas, canonical deterministic tie-breaking, error/denial and hard limits; no new parser/correlation engine, source ACL, storage, IAM or paid API.

This removes **two repeated snapshot reads and two repeated full Timeline triages** per pair bundle and removes the intermediate comparison JSON roundtrip. A single per-source typed-source open remains: it is necessary for final sensitive preview authorization. This also makes the correlation and opened sources use the same in-memory file bytes, rather than a later re-read of an operator-controlled file.

**Regression acceptance:** Existing actual two-Wazuh-snapshot end-to-end tests still exercise valid pair A/B private report, source identity and no overwrite, foreign principal/Case/tenant refusal, bad digests/ordinals and changed source files. A new test keeps the SHA-256 *valid after an owner-file change* but revokes the second correlated event's typed-source reader, requiring a fail-closed private publication (no first-side evidence file created). A matching fingerprint does not grant source ACL.

## PR scope and gates

This branch is built from frozen PR #97 SHA `1b6e7dc1899c1a8c035452ba0b19f815b9eccbcb` and targets frozen, independently CI-green PR #96 SHA `f16a5056875cb19bb9081712beff567675685089`, so the PR only reviews the **five real existing #97 private bundle commits plus this material single-pass refactor and regression**. PR #97 remains untouched; its wider #92 comparison failed the mandatory security gate at 87.25%. Passing a narrower PR **does not** waive the existing red #97/#92 parent dependency or GitHub native reviews. Reconcile upstream approvals, 90% changed-security-critical gate on EVERY required integration diff, L2 risk findings and branch protections before any fusion. The goal is a smaller verifiable functional delivery, not cosmetic PR churn.

**Explicit limitations:** Synthetic operator-controlled local lab, not production IAM, continuous ingestion, durable tenant-separated DB, verified evidence custody, global attribution or causal verdict. SHA-256 is an operator-comparison integrity fingerprint only.
