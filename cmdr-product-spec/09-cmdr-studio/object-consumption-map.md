---
id: studio-std1-object-consumption-map
domain: 09-cmdr-studio
status: draft
owner: CMDR Studio Product Lead
updated: 2026-08-09
source-of-truth: canonical
---
# STD-1 Object / Concept Consumption Map

This map is functional only. It creates no final JSON Schema or physical model.

| Object/concept | Canonical owner | Studio usage | Operations | Gap | Owning future lot |
|---|---|---|---|---|---|
| Studio Asset / Library Entry | CMDR Studio | catalog projection | read/link/classify | no final physical schema | STD-1 complete functionally |
| Tool | CMDR Studio | functional definition | read/draft/version/reference | final object file deferred | Objects/Phase 7 |
| Tool Version | CMDR Studio | exact version/compatibility | read/create/deprecate | final physical model deferred | Objects/Phase 7 |
| Tool Input Contract | CMDR Studio | validation contract | read/update draft | no final JSON Schema | Technique |
| Tool Output Contract | CMDR Studio | output/error semantics | read/update draft | no final format | Technique |
| Tool Permission Need | Security + Studio usage | functional need/risk | assess/reference | atomic namespace unresolved | Permissions |
| Tool Call | CMDR Studio | request/lifecycle/provenance | prepare/read/reconcile | final object deferred | Objects/Phase 7 |
| Tool Call Technical Output | runtime owner + Studio projection | attributed output | read/link | raw format/runtime deferred | Technique |
| Skill | CMDR Studio | reusable contract | read/draft/version/reference | canonical object exists | STD-1 |
| Skill Version | CMDR Studio | exact version/lifecycle | read/create/deprecate | deployment detail deferred | STD-4 |
| Skill Dependency | CMDR Studio | typed Tool/Skill refs | read/update/validate | orchestration excluded | STD-2 for orchestration only |
| Skill Input/Output Contract | CMDR Studio | functional mapping | read/update/validate | no final schema | Technique |
| Provider Reference | Platform Settings | dependency projection | read/reference | provider implementation absent | Settings/Technique |
| Runtime Reference | Platform Settings / runtime owner | availability projection | read/reference | final runtime model absent | Technique |
| Integration Reference | Platform Settings | configured ref | read/reference | connection implementation absent | Settings |
| Secret Reference | Platform Settings | opaque ref only | metadata read/reference | raw value prohibited | Settings/Security |
| Workflow Reference | CMDR Studio | future consumer reference | read/link only | functional orchestration not started | STD-2 |
| Automation Agent Reference | CMDR Studio | future consumer reference | read/link only | agent semantics not started | STD-3 |
| Human Gate Reference | CMDR Studio | boundary reference | read/link only | OPEN-007 relation unresolved | STD-3 / Govern |
| Automation Run Reference | CMDR Studio | provenance boundary only | read/link only | final object/OPEN-015 | STD-3 / Phase 7 |
| Decision | Govern | authority projection | read/reference only | no Studio mutation | Govern |
| Response Run | Govern | distinct execution envelope ref | read/reference only | OPEN-015 bridge | Govern/Phase 7 |
| Result | Govern | canonical outcome ref | read only | Tool output never auto-Result | Govern |
| Endpoint Technical Capability | Endpoint Agent | future technical ref | read projection only | Endpoint not started | Endpoint future |
| Trace / Activity / Job | Shared Capabilities | provenance/background refs | consume/link | generic engines remain Shared | Shared/Technique |

## STD-2 workflow concept addendum

The STD-1 table above is preserved verbatim as historical evidence; its Workflow/Human Gate “future” gaps describe the pre-STD-2 state.

| Object/concept | Canonical owner | Studio usage | Operations | Gap / future boundary |
|---|---|---|---|---|
| Workflow Definition | CMDR Studio | canonical functional contract | read/create/edit | final physical graph schema absent |
| Workflow Version | CMDR Studio | exact version/dependency snapshot | read/create/compare/deprecate | promotion/deployment deferred STD-4 |
| Workflow Draft | CMDR Studio | mutable pre-version context | read/edit | physical persistence unspecified |
| Builder Session | CMDR Studio | editing/base-version/conflict context | open/edit/save/discard | UI/protocol unspecified |
| Workflow Step / Node | CMDR Studio | typed graph element | read/edit/validate | runtime step object absent |
| Tool Step | CMDR Studio | Tool reference/config | read/edit/validate | not Tool Call |
| Skill Step | CMDR Studio | Skill reference/config | read/edit/validate | not Skill mutation |
| Human Gate Step | CMDR Studio | Human Gate reference/config | read/edit/validate | runtime pause/OPEN-007 deferred STD-3 |
| Subworkflow Reference | CMDR Studio | exact Workflow/version ref | read/edit/validate | not copied Workflow |
| Workflow Input / Output / Variable | CMDR Studio | data contract/context | read/edit/validate | no final schema |
| Condition / Branch | CMDR Studio | deterministic control definition | read/edit/preview | not Govern Policy/Decision |
| Data Mapping | CMDR Studio | source-target mapping | read/edit/validate | no transform language |
| Dependency Edge / Parallel Group | CMDR Studio | ordering/concurrency constraint | read/edit/validate | scheduler absent |
| Retry Policy / Idempotency Requirement | CMDR Studio | bounded retry expectation | read/edit/validate | exactly-once not claimed |
| Compensation Step | CMDR Studio | candidate compensation definition | read/edit/validate | not Govern rollback |
| Validation / Compatibility Assessment | CMDR Studio | no-effect readiness evidence | create/read | not execution/deployment |

## STD-3 runtime concept addendum

| Object/concept | Canonical owner | Studio usage | Local operations | Gap | Owning lot/phase |
|---|---|---|---|---|---|
| Automation Agent | CMDR Studio | executable-agent definition ref | read/create/update/reference | canonical object exists; final technical runtime absent | STD-3 functional / Technique |
| Agent Version/reference | CMDR Studio | exact reproducible Agent ref | read/reference/compare | dedicated physical version representation unresolved | Objects/Technique |
| Agent Objective / Constraint | CMDR Studio | bounded task/scope context | define/bind/validate | not permission/authority | STD-3 |
| Agent Team | CMDR Studio | member/role/coordination definition | read/create/update/reference | canonical object exists; framework absent | STD-3 |
| Agent Role | CMDR Studio | functional coordination role | configure/read | not authorization/Settings Role | STD-3 |
| Agent Plan / Step Proposal | CMDR Studio | no-effect candidate plan/step | create/review/reject | not Workflow Definition/action | STD-3 |
| Human Oversight Context | CMDR Studio | supervision/intervention record | inspect/intervene/escalate | no canonical object required | STD-3 |
| Human Gate | CMDR Studio | runtime review checkpoint | request/read/respond/expire/escalate | physical historical states remain; OPEN-007 | STD-3 / Govern |
| Human Gate Request / Response | CMDR Studio | runtime review context/outcome | create/respond/link | not Govern Approval Request/Approval | STD-3 |
| Automation Run | CMDR Studio | business runtime identity/lifecycle | create/read/control/reconcile | physical object file deferred Phase 7 | STD-3 / Objects Phase 7 |
| Automation Run State | CMDR Studio | source-backed lifecycle projection | read/transition on evidence | no final technical state machine | STD-3 / Technique |
| Run Context | CMDR Studio | temporary execution context | bind/read/mask/expire | not canonical source object/memory | STD-3 |
| Run Step / Run Attempt | CMDR Studio | runtime projection/attempt | instantiate/read/retry-link | no final physical schema | STD-3 / Objects/Technique |
| Runtime Queue Item / Schedule | CMDR Studio semantics + Shared mechanism | Run queue/schedule projection | request/read/cancel-before-start | generic engine not Studio-owned | Shared/Technique |
| Runtime Intervention | CMDR Studio | attributed operator control context | prepare/request/record | no Decision mutation | STD-3 |
| Studio Runtime Outcome | CMDR Studio | technical Run outcome | create/read/handoff | not Govern Result/Evidence/Finding | STD-3 |
| Workflow / Workflow Version | CMDR Studio | source definition for Run | read/reference only at runtime | definition remains STD-2-owned | STD-2 |
| Tool / Tool Call / Skill | CMDR Studio | dependencies and runtime calls | reference/call via existing contracts | Tool Call remains distinct from Run | STD-1 |
| Secret / Provider / Runtime Reference | Platform Settings | opaque runtime dependency projection | read/reference | raw secret/admin config prohibited | Settings |
| Approval / Decision / Response Run / Result | Govern | authority/correlation/outcome refs | read/link/handoff only | OPEN-007/015 bridges remain | Govern |
| Endpoint Capability | Endpoint Agent | future executor reference | read/reference only | no Endpoint capability created | Endpoint future |
| Job / Trace / Activity | Shared Capabilities | generic execution/provenance mechanisms | consume/link | not Automation Run business object | Shared |

No physical schema or JSON Schema is introduced by this map.