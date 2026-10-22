---
id: studio-std3-runtime-functional-permissions
domain: 09-cmdr-studio
status: draft
owner: CMDR Studio Product Lead
updated: 2026-08-10
source-of-truth: canonical
open_decisions: [OPEN-007, OPEN-013, OPEN-015]
---
# STD-3 Runtime Functional Permission Needs

This is a functional permission-needs inventory, not final RBAC/ABAC. Both historical namespaces `perm.studio.*` and `perm.cmdr-studio.*` remain unresolved and are not bulk-renamed.

| Need | Capability | Risk / class | Masking / step-up | Ownership / dependency |
|---|---|---|---|---|
| Agent read / restricted config read | 034–039 | 0 | mask restricted config | Studio + Security |
| Agent create/update / access configure | 034–037 | 2 | possible step-up; OPEN-013 | Studio; Settings identity refs |
| Agent planning request / proposal review | 038–039 | 1/2 | source masking | Studio |
| Agent Team read/configure | 037 | 0/2 | tenant isolation | Studio |
| Human Gate read/respond/reassign/escalate | 040–041 | 0/2 | sensitive context separate; SoD possible | Studio; Govern boundary |
| Automation Run read/create | 042–043 | 0/2 | restricted context separate | Studio |
| Run schedule / cancel-before-effect | 045–046 | 2 | revalidation | Studio + Shared/Settings |
| Run start/pause/resume/stop/cancel | 046 | 2 or 3 by underlying effect | step-up/SoD/Govern as applicable | Studio + runtime/Govern |
| Retry / manual intervention | 047/049 | 1/2/3 by effect | current authority + idempotency recheck | Studio + Govern |
| Run context / sensitive context read | 048 | 0 | masking; Secret Reference only | Studio + Settings/Security |
| Tool Call / runtime output read | 044/047/050 | 0 | source permission/masking | Studio/source |
| Control Room read / intervention | 049 | 0/2 | step-up for effectful control | Studio |
| provenance export preparation | 051 | 2 | source/export permission | Studio + Shared/Security |
| cross-tenant Run/provenance read | 048/049/051 | 0 | explicit cross-tenant permission required | Security/Settings |

Agent role/objective/team membership never grants any permission. Human Gate reviewer is not automatically a Govern approver. No class-4 destructive authority is created in STD-3.