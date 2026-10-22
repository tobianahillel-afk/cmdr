# Functional permissions — Collection and Live Response

This document identifies permission needs without defining the final atomic matrix or global namespaces.

| Need | Primary capabilities | Risk | Class | Step-up | Govern | Separation of duties | Owner phase |
|---|---|---|---:|---|---|---|---|
| Endpoint read | 201..215 | wrong tenant or stale target | 0 | no normally | no | no | Permissions |
| Agent capability read | 201,202,207..210 | unsupported operation presented | 0 | no | no | no | Permissions/Endpoint |
| Collection prepare | 202,204..208 | overbroad scope | 2 | possible | by impact | author/reviewer optional | Permissions |
| Collection submit | 202 | unauthorized acquisition | 1/2 | by policy | by class | possible | Permissions/Govern |
| Collection cancel/retry | 203 | duplicate or incomplete effects | 1/2 | possible | by original class | possible | Permissions |
| Raw result read | 203,206,212 | sensitive data exposure | 0 | possible | no | reviewer may differ | Permissions |
| Artifact receive/export | 203..208,211..213 | custody or disclosure | 1/2 | by sensitivity | by policy | possible | Permissions/Trust |
| Live Session request/open/join/extend/close | 209 | persistent endpoint access | 2 | likely | OPEN-007/013 | requester/operator separation possible | Permissions/Govern |
| Endpoint operation execute/interrupt | 210 | mutation, containment or irreversible effect | 0..4 | by class | mandatory 3/4 | strong for 3/4 | Permissions/Govern |
| File download/upload | 205,211 | sensitive data, collision, deployment misuse | 1/2 | possible | by policy | possible | Permissions |
| Process/system inspection | 206 | sensitive snapshot or active effect | 0/1/2 | by effect | by class | possible | Permissions |
| Memory acquisition request | 207 | high impact and sensitive data | 1/2 | likely | by policy | possible | Permissions/Govern |
| Network capture request | 208 | privacy, resource and loss | 1/2 | likely | by policy | possible | Permissions/Govern |
| Containment request | 215/113 | major operational effect | 3/4 | mandatory | mandatory | requester ≠ approver | Govern/Permissions |
| Transcript read | 209,210,214 | secrets/sensitive output | 0 | likely | no | audit viewer may differ | Permissions |
| Operation result verify | 212 | false success or conflation with Govern Result | 2 | possible | no normally | operator/reviewer separation possible | Permissions |
| Custody review | 213 | false integrity claim or history rewrite | 0/2 | possible | no normally | collector/reviewer separation | Trust/Permissions |
| Cross-tenant/environment operation | all | data leakage or wrong target | any | mandatory | by class | strong | Permissions/Govern |

No permission automatically creates a Decision, Response Run, Result, policy assignment, fleet administration right or self-approval.
