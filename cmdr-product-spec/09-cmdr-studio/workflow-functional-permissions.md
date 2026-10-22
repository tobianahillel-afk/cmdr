---
id: studio-std2-workflow-functional-permissions
domain: 09-cmdr-studio
status: draft
owner: CMDR Studio Product Lead
updated: 2026-08-10
source-of-truth: canonical
open_decisions: [OPEN-013]
---
# STD-2 Workflow Functional Permission Needs

No final RBAC/ABAC or atomic namespace is selected. Both historical families `perm.studio.*` and `perm.cmdr-studio.*` remain documented.

| Need | Class/risk | Masking | Step-up / SoD consideration | Functional owner | Admin/runtime owner |
|---|---|---|---|---|---|
| Workflow read / restricted read | 0 | restricted fields masked | scope/tenant | Studio | Security/Settings context |
| Workflow create/edit | 2 | no secret values | OPEN-013 may require step-up | Studio | Studio |
| Builder collaborate | 2 | collaborator scope | conflict/SoD where configured | Studio | Settings identities |
| Workflow Version create/compare | 0/2 | restricted dependencies masked | version owner/reviewer separation possible | Studio | Studio |
| validation run | 1 | denied source stays denied | no authority grant | Studio | source owners |
| Tool/Skill/Human Gate step configure | 2 | dependency metadata only | underlying invoke/authority remains separate | Studio | Tool/Skill/Govern owners |
| condition/mapping/retry/compensation configure | 2 | sensitive metadata masked | OPEN-013; effectful future path separately governed | Studio | future runtime/Govern |
| Secret Reference bind | 2 | opaque reference only | secret-use permission independent | Studio binding | Settings/Security |
| Workflow review / pre-publish candidate | 0/2 | review scope | author/reviewer SoD can apply | Studio | future STD-4 |
| provenance export preparation | 0 | source restrictions preserved | export authorization separate | Studio | Shared/Security |

A Workflow reference never grants Tool invoke, Secret use, Govern Approval/Decision or Endpoint execution permission.
