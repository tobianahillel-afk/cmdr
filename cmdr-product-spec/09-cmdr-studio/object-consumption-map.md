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
