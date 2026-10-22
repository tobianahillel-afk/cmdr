---
id: govern-functional-dependency-map
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-PROD-006, REQ-PROD-015, REQ-PROD-019, REQ-PROD-020]
open_decisions: [OPEN-007, OPEN-008, OPEN-013, OPEN-015, OPEN-019]
---
# Functional Dependency Map — Govern GOV-1 + GOV-2

## GOV-1 dependencies preserved

| Capability | Dependency | Type | Failure behavior |
|---|---|---|---|
| CAP-GOV-001 | Action Request source/version/return origin | source | incomplete/blocked; never invent context |
| CAP-GOV-001..006 | Command/Investigate/target + Settings/Security projections | cross-product | restricted/stale/missing remains explicit |
| CAP-GOV-002/003 | Shared Search/Notifications/Assignments/Versioning/Linking/Trace | shared | degraded state exposed; no silent overwrite |
| CAP-GOV-007/008 | Policy/version/scope + future evaluator | policy | unknown/conflict explicit; no automatic rejection/bypass |
| CAP-GOV-009..013 | Security authority/SoD/step-up/emergency + Settings identities | authority | missing/expired authority blocks; Role never proves authority |
| CAP-GOV-014/015 | CAP-GOV-004..013 + Action Request/Approval/Policy/authority | Decision | unresolved/stale input prevents silent finalization |
| CAP-GOV-016 | CAP-GOV-015 + Studio provenance | handoff | rejected/expired/scope mismatch produces no execution handoff |
| CAP-GOV-010..016 | OPEN-007/013/015 | open governance | no silent Human Gate equivalence, class-2 default or run bridge |

## GOV-2 dependencies

| Capability | Dependency | Type | Failure behavior |
|---|---|---|---|
| CAP-GOV-017 | GOV-1 Execution Handoff + Playbook catalog | planning | no candidate invented; selection remains no-effect |
| CAP-GOV-017/018 | Studio Workflow/Tool + Settings/Endpoint dependency metadata | procedure/runtime projection | missing/deprecated/incompatible dependency is explicit |
| CAP-GOV-018 | exact Decision/Playbook versions | compatibility | version change forces recheck; no silent substitution |
| CAP-GOV-019 | Playbook/Decision/Handoff + parameter definitions | plan | incomplete/mismatched binding blocks plan readiness |
| CAP-GOV-019 | Settings Secret Reference | sensitive input | missing reference blocks; raw secret never copied into Govern |
| CAP-GOV-020 | target source + Settings/Endpoint/Studio health/capability | readiness | stale/drifted/ambiguous/unavailable target never expands scope |
| CAP-GOV-021 | GOV-1 Decision/Approval/Exception + current Plan/readiness | authority reconciliation | expired/mismatched authority blocks Run preparation |
| CAP-GOV-022 | CAP-GOV-018..021 | canonical Run creation | Run not created without pinned lineage; creation != start |
| CAP-GOV-023 | CAP-GOV-020..022 + executor controls/time | Run control | request != confirmation; expired/stale state blocks effectful start/resume |
| CAP-GOV-024 | Plan/Run/step dependencies | execution coordination | failed dependency prevents silent downstream effect |
| CAP-GOV-025 | Studio Workflow/Tool/Automation Run or Endpoint/provider primitive | execution handoff | source rejection/timeout remains explicit; owner not transferred |
| CAP-GOV-026 | technical executor raw status/output | runtime reconciliation | unknown/stale/contradictory state never maps silently to success |
| CAP-GOV-027 | CAP-GOV-020/021/026 + retry/idempotency constraints | error/retry | expiry/drift/limit/non-idempotent uncertainty blocks automatic retry |
| CAP-GOV-028 | Decision/Run expected outcome + authorized observation sources | verification planning | missing criterion/source yields incomplete plan, not success |
| CAP-GOV-029 | Verification Plan + runtime/target/source observations | verification | insufficient/conflicting evidence yields inconclusive/failed/partial state |
| CAP-GOV-030 | original Run/Decision + verification/error + rollback capability | rollback planning | unsupported/unsafe/expired/drifted state blocks effectful rollback |
| CAP-GOV-031 | Rollback Plan + technical rollback/recovery owner + verification | rollback/recovery | partial/failure remains explicit; no exact-restoration claim |
| CAP-GOV-032 | Run + runtime/error + verification + rollback/recovery | Result | missing/contradictory sources prevent fabricated `success` |
| CAP-GOV-033 | CAP-GOV-016..032 + destination owners + Shared Trace/Linking | provenance/handoff | missing/restricted links remain explicit; no destination mutation |
| CAP-GOV-017..033 | Shared Jobs/Notifications/Trace/Activity/Versioning/Reporting/Linking/Recovery | shared mechanisms | shared degradation exposed; generic mechanisms never become Run owner |
| CAP-GOV-017..033 | OPEN-008/013/015 | implementation/authority bridge | no assumed runtime, class-2 default or Automation Run/Response Run equivalence |
| CAP-GOV-033 | OPEN-019 | dissemination | cross-tenant/external handoff blocked unless explicitly authorized |

## Dependency rules

Dependencies never transfer ownership. Missing data, stale sources, permission denial or unavailable Shared/Studio/Settings/Endpoint/provider services produce explicit partial/blocked/unknown behavior. GOV-2 remains provider/runtime neutral, defines no command/API/protocol and leaves GOV-3 Audit Trail/Response Metrics capabilities untouched.