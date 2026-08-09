---
id: qualitative-baseline-govern-gov2
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-09
source-of-truth: baseline-addendum
---
# Qualitative Baseline — Govern GOV-2

This additive baseline preserves the historical `qualitative-baseline.md` in full and records the current GOV-2 delta without rewriting earlier Command/Investigate/GOV-1 evidence.

## Exact starting point

- repository: `tobianahillel-afk/cmdr`;
- canonical branch: `docs/cmdr-product-spec-foundation`;
- PR #2, base `main`;
- exact GOV-2 baseline: **`b8dd93e03443adb9101c7592094a48e358b460e2`**;
- baseline commit: `docs: record Govern GOV-1 post-publication verification`;
- baseline directly descends from `077e3edb5a6fbfe5513279e061e7b4bbee7c71dd`;
- repository public, PR open/Draft/unmerged, auto-merge disabled;
- branch/main root README: exact `# cmdr`, blob `901c74cda52e28b5ff7fc425ddf28ef89f3ad875`;
- `main`: `bc1ec59e5e79ccbaba291e2984b0eb7d5e54129c`.

## Baseline capability state

| Measure | GOV-2 baseline | After GOV-2 functional specification |
|---|---:|---:|
| Global capabilities | 286 | 303 |
| Command capabilities | 27 | 27 |
| Investigate capabilities | 243 | 243 |
| Govern capabilities | 16 | 33 |
| GOV-1 capabilities | 16 | 16 |
| GOV-2 capabilities | 0 | 17 |
| GOV-3 capabilities | 0 | 0 |
| Defined / proposed / planned | 284 / 2 / 286 | 301 / 2 / 303 |
| Govern sections / mandatory tables | 432 / 96 | 891 / 198 |
| Command + Investigate + Govern sections / tables | 7722 / 1716 | 8181 / 1818 |
| GOV-2 sections / tables | 0 / 0 | 459 / 102 |
| Duplicate / recycled CAP-GOV IDs | 0 / 0 | 0 / 0 |
| Concurrent owner conflicts introduced | 0 | 0 |
| Detailed screen rewrites / new Screen IDs | 0 / 0 | 0 / 0 |
| APIs / protocols / commands / product code | 0 | 0 |
| Requirement IDs | 122 | 122 |
| Requirement states conform / partial / absent / contradictory | 99 / 20 / 3 / 0 | 99 / 20 / 3 / 0 |
| OPEN decisions | 18 | 18 |
| New / closed OPEN by GOV-2 | 0 / 0 | 0 / 0 |

## Source/ownership audit

The GOV-2 run re-read GOV-1 `CAP-GOV-001..016`, the canonical Playbook/Response Run/Response Step/Response Rollback/Result/Secret Reference/Workflow/Human Gate objects, Govern Playbooks and Runs & Rollback modules, the Audit Trail and Response Metrics boundaries, Studio Workflow/Tool/Automation/Human Gate/runtime/error contracts, Endpoint command/retry/verification/rollback sources, Settings secrets/connections/health, Shared Jobs/Trace-related mechanisms, Decision→Run→Result journeys, global registers, Requirements, permissions and all 18 OPEN decisions.

The audit confirms:
- Govern owns Response Playbook semantics, Execution Plan, Response Run/Step governance, verification, rollback/recovery governance and canonical Result;
- Studio owns Workflow/Tool/Tool Call/Human Gate/Automation Run;
- Endpoint/provider owners retain technical execution primitives/raw responses;
- Settings retains providers/integrations/secrets/credentials/runtime/tenant/environment administration;
- Command retains Incident/Work Queue; Investigate retains Case/Evidence/Finding/analysis; Shared retains generic mechanisms.

## Safety baseline

GOV-2 selects no provider/runtime, defines no executable command, exploit, bypass, API/protocol, raw secret value, complete object schema, final state machine or final RBAC/ABAC. AI remains proposal-only and essential functions retain deterministic/manual alternatives. GOV-3 Audit Trail and Response Metrics capabilities remain NOT STARTED.

## Publication status

The first four required functional GOV-2 commits publish the 17 functional capability contracts. The fifth traceability/quality commit records registers and pre-publication gate state. GOV-2 remains **PENDING POST-PUBLICATION VERIFICATION** until the fifth functional commit and remote gates are verified; only a subsequent real verification record may promote GOV-2 to PASS.