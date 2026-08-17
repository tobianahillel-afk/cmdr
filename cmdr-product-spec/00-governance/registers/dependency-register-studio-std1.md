---
id: dependency-register-studio-std1
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-09
source-of-truth: registry
---
# Dependency Register — Studio STD-1

Additive STD-1 dependency evidence; it does not replace the global Dependency Register.

| ID | Source | Dependent | Type | Reason | Status | Owner | Blocking | Requirement / OPEN |
|---|---|---|---|---|---|---|---|---|
| DEP-STD-001 | Studio Library + Shared Search/Linking | CAP-STD-001/002/014 | catalog/shared | discovery without duplicating generic Search | active/partial | Studio / Shared | no silent fallback | REQ-PROD-019 |
| DEP-STD-002 | Ownership/Object registers | CAP-STD-001..016 | ownership | preserve one owner per concept/object | active | Product Architecture | yes | REQ-PROD-006; REQ-OBJ-009 |
| DEP-STD-003 | Tool functional definition | CAP-STD-004..009/011/012 | contract | exact Tool/version/I/O semantics | active | Studio | yes | REQ-PROD-016 |
| DEP-STD-004 | Security permission model | CAP-STD-001..016 | security | visibility/invoke/sensitive/cross-tenant separation | partial | Security | yes | REQ-SEC-001,002 |
| DEP-STD-005 | Platform Settings providers/integrations/secrets/tenant/env/health | CAP-STD-003..009/012/015 | admin projection | Studio consumes refs only | partial | Platform Settings | where dependency required | REQ-PROD-017 |
| DEP-STD-006 | Govern Decision/Approval/Response Run/Result | CAP-STD-007..009/016 | authority/handoff | Studio never owns response authority/outcome | active/partial | Govern | effectful actions | OPEN-007,013,015 |
| DEP-STD-007 | Shared Trace/Activity/Jobs/Versioning | CAP-STD-005/006/008/009/013/016 | shared mechanisms | generic provenance/background/version engines remain Shared | partial | Shared | where source needed | REQ-PROD-019,020 |
| DEP-STD-008 | Investigate Tool consumption | CAP-STD-003..009/014/016 | cross-product | Investigate invokes/consumes but qualifies Evidence/Finding separately | active | Studio / Investigate | yes for handoff | REQ-PROD-014; OPEN-015 |
| DEP-STD-009 | Command Studio consumption | CAP-STD-001/014/016 | cross-product | Command consumes refs/projections; Studio not Work Queue | partial | Studio / Command | no | REQ-PROD-013,016 |
| DEP-STD-010 | Endpoint technical capability | CAP-STD-016 future reference | future boundary | Studio may reference but never own endpoint primitive | planned | Endpoint Agent | before endpoint execution | OPEN-008 |
| DEP-STD-011 | final Tool/Tool Call/Automation Run objects | CAP-STD-003..009/016 | object model | final schemas remain Phase 7 | planned | Product Architecture / Studio | before implementation | OPEN-015 |
| DEP-STD-012 | permission namespace normalization | all STD-1 capabilities/screens | implementation detail | `perm.studio.*` and `perm.cmdr-studio.*` coexist; no STD-1 selection | open-detail | Security / Product Architecture | not blocking functional docs | no new OPEN |
