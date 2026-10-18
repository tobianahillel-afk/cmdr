# E10-INV-730A — Unified authorized Wazuh pair navigation

## Functional integration

The analyst can now use both independently developed **metadata-only** discovery methods in ONE running CMDR CLI: offset-based paging (PR #115) and exact first-Case Entry ID nearest-peer lookup (PR #114). Either result can be passed as the same pair of original Entry IDs to the existing *separately authorized* private paired-source dossier. No raw source preview appears in stdout, no IAM or ingestion provider is added, and no remote API is used.

This integration resolves parallel development on the older shared parent `9f43ceab9f96ac7648429f2361bf8b4324861786`; neither frozen PR #114 nor #115 is changed. Only two independently SHA-256-verified Case archives, owner-authorized typed sources, and matching rule/agent/60-second evidence qualify; neither ID nor a metadata lead is an authorization grant.

## Behavior and trust boundaries

- `-list-saved-wazuh-pair-offset 10` discovers later original Entry IDs without opening private event content.
- `-find-saved-pair-first-entry ENTRY_ID` finds the nearest qualified second-Case event for an exact first-Case ID and returns compatible metadata.
- Combining the two discovery modes, or discovery with a private report/preview, is denied **before evidence emission**. Pair pagination is also incompatible with simultaneous private pair selectors (even when a report path was supplied).
- The existing `-inspect-saved-pair-first-entry` + `-inspect-saved-pair-second-entry` + `-private-report` path re-checks typed-source authorization independently, writes only to a new restrictive private artifact and never prints the sensitive report.
- Evidence remains synthetic/offline and proves correlation by temporal co-occurrence, not causation or production authorization.

## Testing and acceptance

The real subprocess-style Go integration reuses the twelve-event/one-peer two-Case snapshot fixture, proves that offset=10 and exact-ID discovery name the same observed pair and reject unauthorized combinations without a partial stdout or private artifact. Both existing per-feature suites run unchanged, plus repository Go and browser/engineering CI. No native reviewer approval or production release is implied until actual exact-head checks and independent L2.
