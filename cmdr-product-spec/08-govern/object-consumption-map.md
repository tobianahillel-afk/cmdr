---
id: govern-object-consumption-map
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-OBJ-001, REQ-OBJ-007, REQ-PROD-006, REQ-PROD-015]
open_decisions: [OPEN-007, OPEN-008, OPEN-013, OPEN-014, OPEN-015, OPEN-019]
---
# Object Consumption Map — Govern GOV-1 + GOV-2

This map records functional use only. It creates no complete schema, JSON Schema, final cardinality or final state machine.

| Objet ou concept | Owner actuel | Usage Govern | Opérations locales | Lacune | Phase propriétaire |
|---|---|---|---|---|---|
| Action Recommendation | source product / producer | candidate input explaining proposed action | read/link/compare | no canonical object required | source product / future Objects review |
| Action Request | Govern | canonical governed request | receive/review/version/transition/withdraw/supersede | final schema/state machine future | Govern / Objects |
| Govern Review Context | Govern local concept | bind request version, reviewer, return origin, review state | create/update/version | local concept, no canonical object yet | Govern / Objects review |
| Govern Queue Item | Govern local projection | Response Inbox projection of an Action Request | assign/reassign/filter/read | projection contract future | Govern / Screens/Objects |
| Target Reference | source owner / Shared linking | identify exact intended target | read/verify/link | target-type normalization future | source owners / Objects |
| Impact Assessment | Govern local concept | structured impact review | create/update/supersede | schema/scales future | Govern / Objects |
| Risk Assessment | Govern local concept | risk-of-action/inaction and uncertainty | create/update/supersede | final scoring model intentionally absent | Govern / Objects |
| Reversibility Assessment | Govern local concept | rollback/recovery feasibility context | create/update/supersede | final schema future | Govern / Objects |
| Policy | Govern | applicable rule/version source | read/link/evaluate | final policy engine/priority semantics future | Govern / Security/Technique |
| Policy Evaluation | Govern | outcome for request/version/scope | run/read/version/dispute | canonical object decision future | Govern / Objects |
| Policy Conflict | Govern local concept | incompatible evaluations/requirements | create/review/disposition | precedence model future | Govern / Security/Objects |
| Policy Exception Candidate | Govern local concept | proposed scoped/time-bound exception | propose/review/withdraw/supersede | active-exception object/state future | Govern / Security/Objects |
| Authority Requirement | Govern local concept + Security rules | required contextual authority | assess/version | final authority model future | Security / Govern / Permissions |
| Authority Context | Govern projection | contextual authority source/scope/expiry | read/assess/link | configured vs contextual contract future | Govern + Settings/Security |
| Approver Candidate | Govern local projection | candidate for Approval | list/assess eligibility | eligibility algorithm future | Govern / Permissions |
| Approval Request | Govern local concept | request Approval from eligible authority | create/update/cancel/expire | canonical object decision future | Govern / Objects |
| Approval | Govern | authority expression | create/read/revoke where allowed/supersede | detailed semantics future | Govern / Objects |
| Delegation | Govern governance concept; identity source Settings | scoped/time-bound delegated authority usage | assess/create/revoke/expire | admin source vs contextual use contract future | Govern + Settings/Security |
| Escalation | Govern local concept | route unresolved authority/approval need | create/update/close | destination policy future | Govern / Permissions |
| Emergency Authorization Context | Govern local concept + Security | time-bound exceptional authority context | create/review/expire/revoke | break-glass policy future | Security / Govern |
| Decision Draft | Govern local concept | non-authoritative preparation | create/update/compare/request-changes | canonical relation to Decision future | Govern / Objects |
| Decision | Govern | authoritative disposition | record/read/supersede; never rewritten by Result | final state machine future | Govern / Objects |
| Decision Condition | Govern local concept | bounded condition attached to Decision | create/update before finalization/read | condition taxonomy future | Govern / Objects |
| Decision Expiration | Govern local concept | time-bound end of authority effect | calculate/inspect/record expiry | execution-time enforcement implementation future | Govern / Permissions |
| Execution Handoff Package | Govern local concept | exact GOV-1 authorized package | prepare/version/withdraw/supersede/read in GOV-2 | not a Response Run | Govern / GOV-1/GOV-2 |
| Response Playbook | Govern | response procedure semantics | search/read/select/version-reference | final schema/editor/runtime future | Govern / GOV-2/Objects |
| Playbook Version | Govern | exact version used for compatibility and execution | compare/pin/review deprecation | physical version model future | Govern / GOV-2/Objects |
| Execution Plan | Govern local concept | no-effect execution intent under exact Decision/Playbook | create/update/version/supersede | final schema/runtime compilation future | Govern / GOV-2/Objects |
| Parameter Binding | Govern local concept | bind definitions to values/references | create/update/remove/version | technical input schemas remain source-owned | Govern / GOV-2/Objects |
| Secret Reference | Platform Settings | sensitive input reference without value | read metadata/bind reference | reveal/use semantics source-owned | Settings / Permissions/Technique |
| Resolved Target | Govern local projection + target source | current resolution of approved Target Reference | resolve/compare/mark drift | canonical target normalization future | source owners / Govern / Objects |
| Readiness Assessment | Govern local concept | current target/executor readiness | run/read/version | implementation/freshness policy future | Govern / GOV-2 |
| Authorization Reconciliation | Govern local concept | execution-time Decision/condition/Approval/Exception recheck | run/read/version/block | enforcement implementation future | Govern / GOV-2/Permissions |
| Response Run | Govern | canonical governed execution envelope | create/manage/transition/link technical runs/close | final persisted state machine future | Govern / Objects |
| Response Step | Govern | functional step/action coordination inside Response Run | instantiate/transition/link executor refs | technical step mapping future | Govern / Objects |
| Execution Action | Govern local concept | exact bounded action intent for a step | prepare/dispatch through owner | no command/protocol chosen | Govern / GOV-2/Technique |
| Executor Handoff | Govern local concept | bounded request to technical owner | prepare/dispatch/correlate/cancel before acceptance where safe | API/protocol future | Govern / GOV-2/Technique |
| Workflow | CMDR Studio | referenced orchestration implementation | read/invoke authorized version/link | Workflow != Playbook | Studio |
| Tool | CMDR Studio | referenced execution/helper capability | read/reference/invoke under source contract | no Govern ownership | Studio |
| Tool Call | CMDR Studio | attributed technical/helper execution | read/link/correlate | Tool Call != Response Run | Studio |
| Automation Run | CMDR Studio | Studio runtime execution/provenance | read/link/interrupt only by Studio rules | bridge to Response Run OPEN-015 | Studio / OPEN-015 |
| Human Gate | CMDR Studio | Studio automation pause/control | read/link | Human Gate != Govern Approval; OPEN-007 | Studio / OPEN-007 |
| technical execution response | Studio/Endpoint/provider owner | raw accepted/status/output/error input | read/link/reconcile | source schemas/provider mapping future | source owner / Technique |
| Runtime Status | Govern local projection | normalize source technical state relative to Run | create/update/version/mark contradiction | final status ontology future | Govern / GOV-2/Objects |
| Retry Record | Govern local concept | bounded retry governance and attempts | propose/check/record attempt/outcome | technical retry implementation source-owned | Govern / GOV-2 |
| Compensation Record | Govern local concept | non-rollback compensating action context | prepare/coordinate/record | compensation != rollback | Govern / GOV-2 |
| Verification Plan | Govern local concept | expected outcome/criteria/source/window | create/update/version | verification engine future | Govern / GOV-2/Objects |
| Verification Assessment | Govern local concept | expected-vs-observed response outcome | create/update/version/dispute | final schema future | Govern / GOV-2/Objects |
| Residual Risk Assessment | Govern local concept | remaining risk/unknowns after response | create/update/version | no universal score | Govern / GOV-2/Objects |
| Rollback Plan | Govern local concept | no-effect rollback intent/preconditions | create/update/version | technical rollback implementation source-owned | Govern / GOV-2/Objects |
| Response Rollback | Govern | governed rollback execution relation/lifecycle | create/manage/link technical refs/verify | final state model future | Govern / Objects |
| Recovery Assessment / Coordination Context | Govern local concept | fallback/recovery status and limitations | create/update/version | recovery primitives source-owned | Govern / GOV-2/Objects |
| Result | Govern | canonical response outcome | create/review/finalize/dispute/supersede | final schema/state permissions future | Govern / Objects |
| Case | Investigate | analysis context/follow-up destination | read/link/handoff prepare | no Govern mutation | Investigate |
| Finding | Investigate | source/follow-up context | read/link/handoff prepare | Govern never requalifies | Investigate |
| Evidence | Investigate | supporting references/provenance | restricted read/link | Govern never requalifies | Investigate |
| Incident | Command | operational source/consumer context | read/link/Result handoff prepare | remains Command-owned | Command |
| Job | Shared | generic background execution/progress mechanism | read/link/correlate | Job != Response Run | Shared |
| GOV-2 Provenance Chain | Govern local composition + Shared Trace | Decision-to-Result cross-record lineage | build/append/reference/export | physical audit store belongs later | Govern / GOV-2/GOV-3 future |
| Cross-Product Handoff Package | Govern local concept | bounded Result/follow-up context to destination owners | prepare/send/withdraw/supersede | destination mutation source-owned; sharing OPEN-019 | Govern / destination owners |
| Provenance Record | Shared/audit semantics by source | cross-product trace primitives | append/reference/export under permission | final physical audit model future | Shared/Security/GOV-3 |

## Invariants

Object references are tenant-scoped and permission-aware. A projection never transfers ownership. Action Request, Approval, Decision, Execution Handoff Package, Execution Plan, Response Run and Result remain distinct. Response Playbook remains distinct from Studio Workflow. Automation Run, Job and Tool Call remain distinct from Response Run. Technical execution output is an input, not canonical Result. Result never rewrites Decision, Evidence or Finding. GOV-2 creates no complete schema, final state machine or GOV-3 capability.