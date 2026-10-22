---
id: endpoint-ept6-functional-permissions
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-12
source-of-truth: canonical
open_decisions: [OPEN-008, OPEN-013, OPEN-015]
---
# EPT-6 Functional Permissions

This file identifies functional permission needs; Security retains the canonical Permission Model and no final RBAC/ABAC matrix is selected.

| Need | Owner | Risk | Action class | Step-up / SoD | Masking | External dependency |
|---|---|---|---:|---|---|---|
| update state/readiness read | Endpoint/Security policy | low | 0 | normal policy | package metadata as needed | Settings target refs |
| update execution request | Endpoint effect under policy | high | 3 candidate | step-up/SoD where policy requires | sensitive refs masked | Settings/Security/Govern context |
| update activate | Endpoint effect | high | 3 candidate | step-up/SoD | masked sensitive context | Settings target/policy |
| update defer/retry | Endpoint local request | medium | 2/3 source-dependent | OPEN-013 | error details masked | Settings/Shared |
| update reversion | Endpoint local effect | high | 3 candidate | step-up/SoD | provenance protected | Settings/Govern boundary |
| buffer/recovery state read | Endpoint | low/medium | 0 | normal policy | payload metadata masked | Shared mechanisms |
| local recovery request | Endpoint effect | medium/high | 2/3 source-dependent | current authority/policy | sensitive state masked | Shared/Settings |
| self-protection state read | Endpoint | sensitive | 0 | security role where needed | component details masked | Security/Settings |
| self-protection mutation | Endpoint if separately sourced | high | 3 candidate | step-up/SoD | sensitive controls hidden | Security/Govern/Settings |
| sensitive security-state read | Endpoint/Security policy | high confidentiality | 0 | restricted role | mandatory masking | Security |
| local audit read | Endpoint | medium | 0 | role/tenant enforcement | secret/sensitive fields masked | Security/Shared |
| sensitive audit read | Endpoint/Security | high confidentiality | 0 | restricted + step-up where policy | mandatory masking | Security |
| Secret Reference use | Settings/Security owner, consumed by Endpoint | high | follows parent action | least privilege/SoD | raw secret never exposed | Settings/Security |
| provenance read | Endpoint + source owners | medium | 0 | source permission | protected fields masked | all source owners |
| cross-tenant read/mutate | source owner | prohibited unless explicit | — | deny by default | no leakage | Security tenant isolation |

Permission possession never creates Decision Authority; runtime privilege never creates product permission.