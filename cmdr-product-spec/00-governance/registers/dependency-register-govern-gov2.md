---
id: dependency-register-govern-gov2
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-09
source-of-truth: registry-shard
---
# Dependency Register — Govern GOV-2

This shard is additive to the global Dependency Register and preserves `DEP-CMD-001..010`, all Investigate dependencies and `DEP-GOV-001..010`. It introduces no protocol, provider, runtime or implementation choice.

| ID | Source | Dependent | Type | Reason | Status | Owner | Blocking | Requirement / OPEN | Review |
|---|---|---|---|---|---|---|---|---|---|
| DEP-GOV-011 | GOV-1 Decision + Execution Handoff Package | CAP-GOV-017..021 | authority/planning | GOV-2 begins from exact finalized Decision/handoff; planning never rewrites or silently extends GOV-1 authority | active | Govern | yes | REQ-PROD-004,008,015; OPEN-013 | GOV-2 |
| DEP-GOV-012 | Govern Playbook/version | CAP-GOV-017..024 | response procedure | Playbook selection/version compatibility precedes Plan/Run; Playbook remains distinct from Studio Workflow | active | Govern | yes | REQ-PROD-006,015,016; OPEN-015 | GOV-2 |
| DEP-GOV-013 | Platform Settings Secret Reference/integration/environment/health | CAP-GOV-019..031 | sensitive/runtime projection | Settings retains secrets and runtime configuration; Govern consumes references/metadata only | partial | Platform Settings / Govern | yes where action requires | REQ-PROD-017; REQ-SEC-001,002; OPEN-008 | Settings/Technique |
| DEP-GOV-014 | target source identity/state + Endpoint/runtime capability | CAP-GOV-020/021/023..031 | readiness/execution | exact targets, freshness, drift and executor support must be rechecked before effectful execution/retry/rollback | partial | source owners / Endpoint / Govern | yes | REQ-PROD-004,015,017; OPEN-008/013 | GOV-2/Technique |
| DEP-GOV-015 | Studio Workflow/Tool/Tool Call/Automation Run | CAP-GOV-017..027/031/033 | execution/provenance | Studio executes its own automation; Workflow != Playbook and Automation Run/Tool Call != Response Run | partial | Studio / Govern | yes where selected executor | REQ-AI-002,004; REQ-PROD-016; OPEN-007/015 | Studio/Objects/Technique |
| DEP-GOV-016 | Response Run/Step and Shared Jobs/Trace/Activity | CAP-GOV-022..033 | run/shared mechanisms | Govern owns Run/Step semantics; Shared provides generic background work and trace without becoming Run owner | active/partial | Govern / Shared | yes | REQ-PROD-019,020; REQ-OBJ-007 | GOV-2/Shared |
| DEP-GOV-017 | technical executor statuses/outputs | CAP-GOV-026/027/029/031/032 | runtime reconciliation | raw technical status/output remains source-owned and is normalized/reconciled without becoming canonical Result | active/partial | Studio / Endpoint / provider owners / Govern | yes | REQ-PROD-005,015,020; OPEN-008/015 | GOV-2/Technique |
| DEP-GOV-018 | Verification Plan + authorized observation sources | CAP-GOV-028/029/032 | verification | expected outcome and independent/source-backed verification separate technical completion from verified response outcome | active/partial | Govern / source owners | yes where verification required | REQ-PROD-002,005,020; OPEN-008 | GOV-2 |
| DEP-GOV-019 | rollback/recovery technical capabilities | CAP-GOV-030/031/032 | rollback/recovery | rollback support != eligibility/authority/success; technical reverse primitives remain source-owned | active/partial | Govern / Studio / Endpoint / provider owners | yes when rollback required | REQ-PROD-004,005,015; OPEN-008/013/015 | GOV-2/Technique |
| DEP-GOV-020 | canonical Result + downstream product owners | CAP-GOV-032/033 and Command/Investigate/Settings/Studio | outcome/handoff | Govern owns Result/provenance; downstream handoffs never silently mutate Incident, Case, Evidence, Finding, Workflow or Settings configuration | active | Govern / destination owners | yes | REQ-PROD-005,008,009,015,020; OPEN-019 | GOV-2/GOV-3 future |

## Open-decision boundary

The 18 existing OPEN decisions are preserved. Particularly, OPEN-007 remains for Human Gate/Govern semantics, OPEN-008 for actual platform/runtime support, OPEN-013 for class-2/effect governance, OPEN-015 for Automation Run/Response Run bridging and OPEN-019 for external/cross-tenant sharing. GOV-2 creates and closes no OPEN.