---
id: endpoint-ept4-permissions
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-11
source-of-truth: canonical
---
# EPT-4 Functional Permission Needs

No final RBAC/ABAC namespace is selected.

| Need | Owner/source | Risk | Class | Step-up / SoD | Masking | Govern dependency |
|---|---|---:|---:|---|---|---|
| Collection request read | Investigate | low | 0 | source permission | sensitive scope | no |
| collection technical intake/plan | Endpoint | medium | 1/2 | tenant/target checks | paths/subjects | according action |
| file/system collection execute | Endpoint | medium/high | 2/3 | impact dependent | content/paths | impact dependent |
| memory collection execute | Endpoint | high | 2/3 | step-up likely | memory highly sensitive | impact dependent |
| network capture start/stop | Endpoint | medium/high | 2/3 | bounded scope | network content | impact dependent |
| collection cancel/retry/resume | Endpoint | medium | 2 | OPEN-013 | status/output | according source action |
| collection output read | Endpoint/Investigate destination | high | 0 | source permission | content | no automatic qualification |
| Live Response session read/request/open/close | Endpoint + Investigate context | medium/high | 0/2 | session/target/authority | transcript | policy/Govern as required |
| command request/non-mutating execute | Endpoint | high | 2 | target/session/authority | params/output | policy dependent |
| effectful command/script | Endpoint technical + Govern authority | very high | 3 | mandatory authority/SoD | params/output/secrets | mandatory |
| script execute | Endpoint | high/very high | 2/3 | runtime/source/authority | content/params/output | effect dependent |
| file read/download | Endpoint | high | 2 | exact path/source | content | policy dependent |
| temporary upload/collision override | Endpoint technical + Govern authority | very high | 3 | mandatory | file/path | mandatory |
| execution output read | Endpoint | high | 0 | source permission | sensitive output | no automatic Result |
| operator takeover/control | Endpoint/Investigate | high | 2 | participant/control policy | transcript | policy dependent |
| provenance read | source owners | variable | 0 | source permissions | restricted refs | no |
| cross-tenant read/execute | none by default | critical | — | deny | no leakage | explicit authority cannot bypass tenant isolation |

Raw secrets are never arbitrary display fields; use Secret References when applicable.