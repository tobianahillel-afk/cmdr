# E10-INV-520A — Actionable 60-second event links across independent Wazuh Cases

## User-visible capability

The existing offline comparison of two independently saved and SHA-256-verified operator-owned Cases already reports aggregate Timeline overlap and shared Wazuh **rule + agent + five-minute UTC window** counts. An analyst cannot inspect the precise alerts underlying those aggregate leads without manually hunting both archives.

This increment reuses the **same first-party read and authorization flow**, without a second search engine or source merger, to add up to **10 navigation leads** with actual Wazuh Entry IDs and UTC times from BOTH previously authorized Cases when matching **Wazuh rule ID + agent ID** events are no more than **60 seconds** apart (inclusive). The standard `-compare-saved-case` / `-compare-saved-with-case` command and the existing optional private report export format remain otherwise unchanged.

Example new receipt property `shared_authorized_wazuh_60s_event_pairs`:

```json
[{
  "rule_id": "5710",
  "agent_id": "001",
  "first_entry_id": "local-jsonl-entry-<first-case-id-digest>",
  "second_entry_id": "local-jsonl-entry-<second-case-id-digest>",
  "first_time_utc": "2026-10-17T10:02:00Z",
  "second_time_utc": "2026-10-17T10:03:00Z",
  "offset_seconds": 60
}]
```

Each Entry ID is a **Case-scoped navigation reference**. Reopening the event requires the original selected Case and a fresh first-party typed-source owner authorization. The pair DOES NOT prove a global event identity, causality, a common physical host, malicious activity or forensic custody. Raw message, alert description, hostname, IP, file path and source preview are NOT copied into pairs. Rule/agent IDs are authorized metadata and are visible only to an operator already entitled to both archives.

## Existing boundaries reused, not weakened

- The two **separately retained operator SHA-256** digests are checked against each file's exact bytes before parsing.
- Each private Case is scoped independently to principal, tenant, environment, Case, Timeline entries and authorized typed sources. Cross-tenant or cross-environment comparison remains forbidden. There is no common owner graph.
- The existing complete first-party Case triage traversal and `buildCaseComparison` must succeed for BOTH files before any output is written. Wazuh correlation uses only imported events tied to first-party Wazuh typed-source metadata in these narrowly scoped owners.
- Only definite observed Wazuh event times participate; missing rule or agent metadata cannot be synthesized from generic JSONL.
- Group by rule+agent, sort by UTC time, use binary search to find the nearest event from second Case per event in first Case. Bounded input is inherited (<=5000 Case entries per file); runtime roughly O(n log n), deterministic ordering, maximum 10 pairs. Standard five-minute aggregates remain available.
- Generic-only Cases produce no Wazuh pairs. An alert 60s apart is included; >60s is excluded by design. Incomplete/revoked owner or modified second archive fails before any receipt.
- This is **synthetic offline forensic investigation**, not continuous Wazuh ingestion, production IAM, signed evidence or a persisted global correlation index. No paid APIs, new Go modules, new governance machinery or CI-policy changes.

## Tests, CI, and integration constraints

Existing on-disk integration test actually imports real Wazuh 4.x JSONL into TWO distinct authorized private owner snapshots, removes original feeds, verifies both digests and compares Cases via the real CLI. The test now sets exact +60s rule+agent overlap, requires precisely one pair linked to valid Entry IDs in their respective saved Case catalogs, verifies original five-minute aggregates, and checks no raw Wazuh descriptions/credential/source preview appear in pair metadata. Existing generic-only tests now reject fabricated pairs; unauthorized or corrupt second archive tests still fail closed.

**GitHub parent is PR #92 frozen HEAD `fced2c5c89f0458247be9cc0644e516118af2ae6`**, whose functional CI is green but consolidated changed security-critical coverage was **88.72% < mandatory 90%**. The new branch/PR must stay DRAFT and unmerged until exact-head Linux/macOS/Windows and engineering-kernel mandatory 90% checks pass, as must all upstream merge requirements. Do not lower or waive this 90% floor, and do not mutate any L2-frozen ancestor SHA. One independent READ-ONLY L2 review of this material source-linking increment, never a second review of parent #92.
