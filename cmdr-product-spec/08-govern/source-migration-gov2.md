---
id: govern-source-migration-gov2
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: migration-audit
requirements: [REQ-PROD-006, REQ-PROD-015, REQ-PROD-016, REQ-PROD-020]
open_decisions: [OPEN-007, OPEN-008, OPEN-013, OPEN-015, OPEN-019]
---
# Source Migration Audit — Govern GOV-2

## Purpose

Record the audit of historical Govern execution material used by GOV-2, identify the new canonical capability sources, and prevent duplicate active ownership. This is a documentary migration map, not a runtime migration or deprecation programme.

## Baseline note

At the exact GOV-2 baseline `b8dd93e03443adb9101c7592094a48e358b460e2`, there were no standalone canonical files named `08-govern/concepts.md`, `08-govern/workflows.md`, `08-govern/states.md` or `08-govern/source-migration.md`. GOV-2 therefore audited equivalent semantics in the existing product/module/object/journey/security sources rather than claiming those absent files were read.

## Historical execution sources audited

| Historical/source document | Need retained | GOV-2 canonical destination | Migration treatment |
|---|---|---|---|
| `08-govern/modules/playbooks/README.md` | Playbook purpose, versioning, preconditions, rollback/testing evidence | CAP-GOV-017..019 + updated module README | active supporting source aligned; not deprecated |
| `08-govern/modules/playbooks/playbook-versioning.md` | immutable versions, active version, rollback-to-version, diffs | CAP-GOV-017/018 | retained as supporting historical contract; no duplicate owner |
| `05-domain-model/objects/playbook.md` | Govern owner, lifecycle/version relation | CAP-GOV-017/018 and future Objects phase | canonical object remains active; no schema rewrite |
| `08-govern/modules/runs-and-rollback/README.md` | live progress, partial success, manual intervention, verification | CAP-GOV-020..033 + updated module README | active supporting source aligned; not deprecated |
| `08-govern/modules/runs-and-rollback/rollback-model.md` | rollback trigger/state/scope/authority | CAP-GOV-030/031 + updated support contract | active supporting source aligned; not deprecated |
| `05-domain-model/objects/response-run.md` | Govern-owned canonical Run | CAP-GOV-022..033 | object remains active; lifecycle in capability spec is functional, not final state machine |
| `05-domain-model/objects/response-step.md` | Govern-owned response step | CAP-GOV-024..027 | object remains active; technical step mapping future |
| `05-domain-model/objects/response-rollback.md` | Govern-owned rollback relation | CAP-GOV-030/031 | object remains active; no command/implementation added |
| `05-domain-model/objects/result.md` | Govern-owned canonical Result | CAP-GOV-032/033 | object remains active; raw technical outputs remain source-owned |
| `13-user-journeys/decision-to-response-run.md` | Decision → conditions → Playbook/version → Run | CAP-GOV-017..023 | journey retained; capability files become detailed functional specification |
| `13-user-journeys/response-run-to-result.md` | steps → partial success → rollback → verification → Result | CAP-GOV-024..032 | journey retained; no duplicate execution owner |
| `13-user-journeys/result-to-command.md` | Result handoff to Command | CAP-GOV-033 | journey retained; Command owns Incident transition |
| Studio Workflow/orchestration/error/runtime sources | technical workflow, compensation, retry, Human Gate, Automation Run | CAP-GOV-017/018/024..027/033 boundaries | **not deprecated**; Studio retains ownership |
| Endpoint command/retry/verification/rollback sources | technical execution and reverse primitives | CAP-GOV-020/025..031 boundaries | **not deprecated**; Endpoint retains ownership |
| Settings secrets/connections/health | secret/config/runtime source | CAP-GOV-019/020/025/030/031 boundaries | **not deprecated**; Settings retains ownership |
| Shared Jobs/Trace/Recovery mechanisms | generic async/provenance/recovery mechanisms | CAP-GOV-022..033 boundaries | **not deprecated**; Shared retains ownership |
| GOV-3 Audit Trail / Response Metrics modules | future audit/metrics needs | future GOV-3 | explicitly retained and NOT STARTED |

## Duplicate-resolution rules

No historical source is deprecated merely because GOV-2 adds detailed capability contracts. The new `CAP-GOV-017..033` files are canonical for GOV-2 functional capability semantics, while existing domain objects, journeys, Studio/Endpoint/Settings/Shared sources remain canonical for their own ownership domains.

The following must not be merged or replaced:
- Response Playbook vs Studio Workflow;
- Response Run vs Automation Run / Tool Call / Shared Job;
- Govern Result vs technical executor output;
- compensation vs rollback;
- Govern verification vs source-owned observations;
- Govern rollback governance vs technical rollback primitive.

## Active-source outcome

- competing active GOV-2 capability owner: **0**;
- deprecated Studio Workflow/Automation Run/Tool source: **0**;
- deprecated Endpoint execution source: **0**;
- deprecated Shared Job/Trace/Recovery source: **0**;
- deprecated GOV-1 Decision/Handoff source: **0**;
- deprecated GOV-3 Audit Trail/Response Metrics source: **0**;
- GOV-3 capability created: **0**.

## Implementation boundary

This audit migrates documentary responsibility only. It selects no runtime, provider, connector, command, API/protocol, schema or data migration procedure.