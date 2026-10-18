# E10-INV-600A — Integrated Wazuh correlation contract stabilization

## Scope and concrete progress

After independently authorizing TWO fingerprint-verified offline Case snapshots, the same first-party collection of rule+agent/time Wazuh observations feeds two pure derived views: bounded deterministic 5-minute windows and source-linked inclusive 60-second pairs.

E10-INV-590A correctly reduced first-party Wazuh cohort collection from four passes to two. However, its two pure formatting functions still had fallible return signatures and forced the authenticated comparator to contain unreachable error branches. This lot **removes those nonexistent error contracts**. It preserves the single existing fallible security boundary in collectAuthorizedWazuhObservations for invalid metadata, uncertain event clocks, malformed source time and first-party typed-source ownership.

## Functional security regression

The added end-to-end Go regression uses two actual Wazuh 4.x JSONL alert files (same 5710 rule/001 agent, 30s offset), independently private-persisted owner snapshots, independently recomputed SHA-256, independent Case authorizations, and the complete compareLoadedAuthorizedSavedCases pipeline. One 5-minute lead and one <=60-second pair must be returned. Unauthorized principal, cross-tenant data, cancelled context and changed Wazuh clock uncertainty or timestamp must be rejected before a complete correlation receipt exists.

Pure derivation views also exercise 25 genuine-sized sorted source cohorts: only 10 deterministic windows and pairs with 30s offset, correct Wazuh rule+agent and scoped event IDs are returned. Same-time records with an unmatching agent never correlate. Security checks remain at the existing owner/time boundary; do not add a second generic validation framework.

## Delivery and limitations

No new parser, search engine, third-party dependency, API or IAM model. No change to existing Case/source authorizations or output schemas. Wazuh source provenance remains synthetic operator input; 5-minute or 60-second co-occurrence does not prove causality, and SHA-256 is not a signature of forensic custody.

This is a child draft of frozen PR #100 SHA f83b2c8fd84dca6d922b8eccfeba1c1d39b25222, and **no ancestor SHA is modified**. Exact-HEAD Linux/Windows/macOS runtime and unchanged mandatory security-critical changed-line >=90% coverage are required. The goal is to stabilize a real integrated Case investigation, not manufacture artificial test pass scores. Only the substantive new security boundary, not each commit, warrants an independent L2; avoid duplicates with the existing PR #100 review family. Prioritize existing green foundation PR #36, product slice #81 and data-safety slice #89 for review and integration after the risk/CI gates are satisfied.
