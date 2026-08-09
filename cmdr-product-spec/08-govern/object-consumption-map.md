---
id: govern-object-consumption-map
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-OBJ-001, REQ-PROD-006, REQ-PROD-015]
open_decisions: [OPEN-007, OPEN-013, OPEN-014, OPEN-015]
---
# Object Consumption Map — Govern GOV-1

This map records functional use only. It creates no complete schema, JSON Schema, final cardinality or final state machine.

| Objet ou concept | Owner actuel | Usage Govern | Opérations locales | Lacune | Phase propriétaire |
|---|---|---|---|---|---|
| Action Recommendation | source product / producer | candidate input explaining proposed action | read/link/compare | no canonical object required by GOV-1 | source product / future Objects review |
| Action Request | Govern | canonical governed request | receive/review/version/transition/withdraw/supersede | final schema/state machine future | Govern / Objects |
| Govern Review Context | Govern local concept | bind request version, reviewer, return origin, review state | create/update/version | local concept, no canonical object yet | Govern / Objects review |
| Govern Queue Item | Govern local projection | Response Inbox projection of an Action Request | assign/reassign/filter/read | projection contract future | Govern / Screens/Objects |
| Target Reference | source owner / Shared linking | identify exact requested target without transferring ownership | read/verify/link | target-type normalization future | source owners / Objects |
| Impact Assessment | Govern local concept | structured impact review | create/update/supersede | schema/scales future | Govern / Objects |
| Risk Assessment | Govern local concept | risk-of-action/inaction and uncertainty | create/update/supersede | final scoring model intentionally absent | Govern / Objects |
| Reversibility Assessment | Govern local concept | rollback/recovery feasibility context | create/update/supersede | final rollback contract belongs GOV-2 | Govern / GOV-2/Objects |
| Policy | Govern | applicable rule/version source | read/link; GOV-1 does not administer policy lifecycle | final policy engine/priority semantics future | Govern / Security/Technique |
| Policy Evaluation | Govern | outcome for request/version/scope | run/read/version/dispute | canonical object decision future | Govern / Objects |
| Policy Conflict | Govern local concept | represent incompatible evaluations/requirements | create/review/resolve-by-disposition | precedence model future | Govern / Security/Objects |
| Policy Exception Candidate | Govern local concept | proposed scoped/time-bound exception | propose/review/withdraw/supersede | active-exception object/state future | Govern / Security/Objects |
| Authority Requirement | Govern local concept + Security rules | required contextual authority | assess/version | final authority model future | Security / Govern / Permissions |
| Authority Context | Govern projection | resolved contextual authority source/scope/expiry | read/assess/link | configured vs contextual contract future | Govern + Settings/Security |
| Approver Candidate | Govern local projection | person/service principal candidate for approval | list/assess eligibility | eligibility algorithm future | Govern / Permissions |
| Approval Request | Govern local concept | request an Approval from eligible authority | create/update/cancel/expire | canonical object decision future | Govern / Objects |
| Approval | Govern | authority expression | create by authorized disposition/read/revoke where allowed/supersede | detailed semantics future | Govern / Objects |
| Delegation | Govern governance concept; identity source Settings | scoped/time-bound transfer of authority usage | assess/create/revoke/expire where permitted | admin source vs contextual use contract future | Govern + Settings/Security |
| Escalation | Govern local concept | route unresolved authority/approval need | create/update/close | destination policy future | Govern / Permissions |
| Emergency Authorization Context | Govern local concept + Security | time-bound exceptional authority context | create/review/expire/revoke | break-glass policy future | Security / Govern |
| Decision Draft | Govern local concept | non-authoritative preparation | create/update/compare/request-changes | canonical relation to Decision future | Govern / Objects |
| Decision | Govern | authoritative disposition | record/read/supersede; no target execution | final state machine future | Govern / Objects |
| Decision Condition | Govern local concept | bounded condition attached to Decision | create/update before finalization/read | condition taxonomy future | Govern / Objects |
| Decision Expiration | Govern local concept | time-bound end of authority effect | calculate/inspect/record expiry | execution-time enforcement future | Govern / GOV-2/Permissions |
| Execution Handoff Package | Govern local concept | exact authorized package for future GOV-2 | prepare/version/withdraw/supersede | not a Response Run; final contract GOV-2 | Govern / GOV-2 |
| Case | Investigate | analytical context projection | read/link | access may be partial | Investigate |
| Finding | Investigate | asserted source context | read/link | Govern never requalifies | Investigate |
| Evidence | Investigate | supporting references/provenance | read restricted projection/link | Govern never requalifies | Investigate |
| Incident | Command | operational context/impact/urgency/return origin | read/link | remains Command-owned | Command |
| Workflow | CMDR Studio | optional routed automation reference | read/invoke only under future permission | Workflow != Playbook | Studio |
| Human Gate | CMDR Studio | automation pause/control projection | read/link | relation to Govern Approval under OPEN-007 | Studio / OPEN-007 |
| Automation Run | CMDR Studio | provenance of automated assistance | read/link/interrupt only by Studio rules | bridge to future Response Run OPEN-015 | Studio / OPEN-015 |
| future Response Run | Govern | future execution record | read placeholder only in GOV-1 | NOT STARTED in GOV-1 | GOV-2 |
| future Result | Govern | future verified execution outcome | read historical/projection only | no Result produced in GOV-1 | GOV-2/GOV-3 |
| Tool | CMDR Studio | optional deterministic/AI helper capability | select/reference by permission | no Govern ownership | Studio |
| Tool Call | CMDR Studio | attributed helper execution | read/link | provenance contract future | Studio |
| Provenance Record | Shared/audit semantics by source | cross-product trace of request/policy/authority/approval/decision | append/reference/export under permission | final physical audit model future | Shared/Security/GOV-3 |

## Invariants

Object references are tenant-scoped and permission-aware. A projection never transfers ownership. Action Request, Approval, Decision and future Response Run/Result remain distinct. Execution Handoff Package is a GOV-1 local functional package, not a Response Run and not proof of execution.
