---
id: dependency-register-studio-std4
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-10
source-of-truth: registry
---
# Dependency Register — Studio STD-4

| Family | Source | Dependency | Ownership/boundary |
|---|---|---|---|
| DEP-STD4-001 | CAP-STD-052..060 | Evaluation/Simulation/Regression evidence | Studio assurance semantics only |
| DEP-STD4-002 | CAP-STD-052..068 | Tool/Skill/Workflow/Agent/Run exact versions | consumes STD-1/2/3; no redefinition |
| DEP-STD4-003 | CAP-STD-058/062/065/066/068 | Approval/Decision/Response Run/Result | Govern-owned; no implicit authority |
| DEP-STD4-004 | CAP-STD-058/064..068 | tenant/environment/provider/integration/Secret Reference/health | Settings-owned administration |
| DEP-STD4-005 | CAP-STD-052..068 | Jobs/Trace/Activity/Reporting/Export/Notifications/Versioning/Recovery | Shared generic mechanisms |
| DEP-STD4-006 | CAP-STD-064..068 | Endpoint technical deployment/update/runtime | Endpoint-owned; no CAP-EPT created |
| DEP-STD4-007 | CAP-STD-052..067 | Security permission model | functional needs only; no final RBAC |
| DEP-STD4-008 | CAP-STD-061..067 | exact Version lineage | Studio Version object; no package registry |
| DEP-STD4-009 | CAP-STD-066 | prior Studio asset version | reversion != Govern Response Rollback |
| DEP-STD4-010 | CAP-STD-068 | STD-1/2/3/4 evidence + registers | documentary closure only |

OPEN-003/007/008/013/015 remain open where consumed. No dependency is asserted implemented merely because it is referenced.