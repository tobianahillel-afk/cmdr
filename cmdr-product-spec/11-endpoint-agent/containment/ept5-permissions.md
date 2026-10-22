---
id: endpoint-ept5-permissions
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-11
source-of-truth: canonical
---
# Endpoint EPT-5 — Functional Permissions

No final RBAC/ABAC namespace is selected. Every effect permission is distinct from CRUD/read access and from Govern authority.

| Need | Owner | Risk | Action class | Step-up / SoD | Govern dependency | Settings dependency |
|---|---|---|---:|---|---|---|
| response primitive read | Endpoint | low | 0 | source access | no | tenant/env |
| primitive eligibility/precheck | Endpoint | low | 1 | least privilege | no effect | Policy/capability projection |
| process suspend/resume/terminate | Endpoint technical | high | 3 | expected | required | Policy |
| host isolation/release | Endpoint technical | high | 3 | expected | required | Policy/Fleet |
| network block/unblock | Endpoint technical | high | 3 | expected | required | Policy |
| quarantine/release | Endpoint technical | high | 3 | expected | required | Policy |
| file delete | Endpoint technical | destructive | 3/4 | strong step-up/SoD | required | Policy |
| file restore/recovery | Endpoint technical | high | 3 | expected | required | Policy |
| service/system control | Endpoint technical | high | 3 | expected | required | Policy |
| local session lock/terminate | Endpoint technical | high | 3 | expected | required | identity/Policy refs |
| technical verification | Endpoint | low/no-effect | 1 | source permission | no authority | source availability |
| technical reversal | Endpoint technical | high | 3/4 | expected | Govern Rollback required | Policy |
| technical outcome/provenance read | Endpoint | sensitive | 0 | masking/source permission | no | tenant/env |
| cross-tenant response execute | none by default | critical | 4/deny | explicit authority + isolation | mandatory | tenant/env boundary |

Raw secrets are never ordinary parameters; only authorized Secret References are consumed. Sensitive process, network, file, memory/session and response output metadata require masking and least privilege.