---
id: studio-std4-functional-permissions
domain: 09-cmdr-studio
status: draft
owner: CMDR Studio Product Lead
updated: 2026-08-10
source-of-truth: canonical
open_decisions: [OPEN-007, OPEN-013, OPEN-015]
---
# STD-4 Functional Permission Needs

| Need | Risk | Class | Masking | Step-up / SoD | Owner/dependency |
|---|---|---:|---|---|---|
| Evaluation read/create/update/execute | assessment / sensitive inputs | 0–2 | source-sensitive data masked | step-up for sensitive/effectful path | Studio + source permissions |
| Simulation read/create/run | no-effect assurance | 0–2 | synthetic/source restrictions | no production authority | Studio |
| Regression/readiness evaluate | no-effect assessment | 0–2 | source masking | reviewer separation where required | Studio |
| Publishing Candidate create/update/review | lifecycle preparation | 2 | source-sensitive evidence masked | SoD may apply | Studio; Govern dependency where policy requires |
| Publish/release promote | real lifecycle effect | 3 where effectful | no secret exposure | current authority + SoD | Studio + Govern/Settings when applicable |
| Deployment read/create/stage/promote | lifecycle effect | 0–3 | tenant/environment isolation | step-up / Govern dependency if effectful | Studio; Settings config; technical owner |
| Deployment reversion | Studio asset lifecycle effect | 2–3 | source masking | explicit authority | Studio; **not Govern Response Rollback** |
| Deprecate/retire/migration prepare | lifecycle change | 2–3 | consumer-sensitive refs | SoD where required | Studio |
| Cross-tenant evaluation/deployment read | sensitive | 0–1 | strict tenant masking | explicit cross-tenant permission | Security/Settings policy |

`perm.studio.*` and `perm.cmdr-studio.*` remain an unresolved historical namespace anomaly. No bulk rename, atomic catalog redesign or final RBAC/ABAC is selected.