# E10-INV-530A — One authorized Wazuh cohort extraction for two cross-Case investigation views

## Engineering integration, not another engine

The independent two-Case Wazuh investigation in PR #93 exposes both **fixed five-minute rule+agent cohort counts** and **actual nearest 60-second event-linked navigation leads**. Both were originally implemented with duplicated logic for identifying Case-imported Wazuh entries, restricting to first-party typed `telemetry-event` owners, reading rule/agent identity, rejecting uncertain/non-observed timestamps and parsing UTC times. This multiplied security-critical changed branches and risked divergent interpretations of the same evidence.

One compact first-party helper `collectAuthorizedWazuhObservations` now performs this already-required extraction exactly once PER respective independently authorized Case (the function is called separately for each Case/view). It returns sorted, canonical **per-Case** Wazuh rule+agent observations containing Entry IDs and definite UTC times. Both existing consumers reuse it: five-minute comparison counts now bin these same observations, while 60-second linking performs chronological nearest-event binary search on the same pre-sorted cohorts. No two owner datasets are merged, and no new authorization policy or validator framework is introduced.

### Material functional fix

When the second Case contains multiple authorized Wazuh events with the SAME nearest timestamp before the first Case event, binary search formerly considered only the **last** source ID in the prior equal-time group. It could therefore select a different Entry ID solely because of an implementation detail of the boundary search, contrary to the stable tie-breaking intent. The selected nearest UTC time is now used to find the **first canonical Entry ID at that exact instant**. This gives a reproducible, stable navigation ref even with duplicate timestamps and many events, with no extra owner lookup or access grants. At most 10 pairs; time window remains inclusive ±60s.

### Security/functional acceptance

- Existing PR #93 full CLI test imports two separate real Wazuh 4.x JSONL owner archives, retains two independent SHA-256 bytes fingerprints, independently scopes each authorized Case/typed source and runs the same first-party comparison through its complete owner Timeline. Raw alerts, bearer and source preview are not placed in the comparison result.
- Test adds two SECOND-Case same-rule+agent Wazuh alerts at exactly `2026-10-17T10:00:05Z` and requires the first linked pair to use the **lexicographically lowest authorized canonical Entry ID** among them, independent of raw JSONL input order. Existing chronological first-Case ordering, nearest delta checks, hard 10 limit and stable byte-for-byte report replay remain.
- Both five-minute counts and exact 60-second pairs now share one importer/typed-source/time validation, reducing opportunities for semantic drift. Generic-only Cases still yield no Wazuh cohorts; uncertain timestamps still fail closed. No production IAM, datastore, network Wazuh ingest, signature, attribution or chain-of-custody claim. No new Go modules/third-party packages/paid API.
- Exact-head Go tests on Linux/macOS/Windows, Node browser tests, and **unchanged mandatory 90% changed security-critical coverage** via `engineering-kernel`. No bypass. Parent PR #93 was 84.21% and not mergeable. Any failure retains DRAFT status.

## Integration decision

A short child PR of immutable/frozen PR #93 SHA `e8bab1bfcec485f02e7c63a226ed3e9b363ce0ea` (not `main` or #92). Original PR #93, its ancestors and L2 review SHAs remain untouched. This consolidates shared investigation behavior while correcting a real replay/triage ambiguity, not another governance microtask. The much larger PR #92 candidate against #81 and foundational PR #36 must still pass their individual required checks/reviews before any integration; no premature merge.

One focused independent read-only L2 is warranted because both separately authorized cross-Case correlation outputs now depend on one shared extraction boundary. Review correct owner isolation, source/time parsing, tie selection, sorted invariants, real data paths and regression risks, not cosmetic tests.
