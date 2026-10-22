---
id: govern-permissions-gov3
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical-addendum
---
# Govern Functional Permission Needs — GOV-3 Addendum

This document describes functional permission needs only. Security retains the final permission catalog/RBAC/ABAC/step-up model. Existing `perm.govern.metrics.read` may be referenced; no new atomic identifier here is promoted to canonical implementation.

| Besoin fonctionnel | Capability | Risk | Class | Masking / scope | Step-up potential | SoD | Owner / technical owner |
|---|---|---|---:|---|---|---|---|
| Audit Trail read | 034..038 | sensitive history | 0 | tenant/environment/object | possible for restricted context | source-dependent | Govern / Security |
| restricted audit read | 034..038 | identities/actions | 0 | field masking/classification | yes | source-dependent | Security |
| actor identity read | 034..040 | personal/sensitive | 0 | identity minimization | possible | yes for reviewer conflict | Settings/Security |
| audit search | 038 | query over sensitive corpus | 0 | scope filters cannot widen access | possible | no | Shared Search + Govern semantics |
| reconstruction create/read | 035/036 | derived historical interpretation | 1/2 | source permissions inherited, not expanded | possible | reviewer separation when required | Govern |
| gap assessment create/review | 037 | audit conclusion candidate | 1/2 | source refs/masking | possible | review separation possible | Govern/Security |
| contradiction review | 037 | sensitive dispute | 2 | source visibility | possible | reviewer independence possible | Govern |
| Audit Evidence Package create | 038 | data aggregation | 2 | classification/masking required | yes | package reviewer possible | Govern |
| Audit Evidence Package export preparation | 038 | disclosure | 2 | Shared Export rechecks visibility | yes | external release separate | Govern + Shared/Security |
| metric read | 039..046 | aggregate sensitive data | 0 | `perm.govern.metrics.read` family | possible | no | Govern/Security |
| sensitive metric dimension read | 039..046 | identity/tenant/target | 0 | suppression/aggregation | yes | possible | Security/Settings |
| cross-tenant metric read | 039..046 | isolation | 0 | explicit tenant scope | yes | possible | Security/Settings |
| metric comparison | 039..046 | derived observation | 1 | same visibility constraints | no/possible | no | Govern + Shared Metrics |
| Trend Assessment create | 046 | derived interpretation | 2 | source refs preserved | possible | reviewer separation possible | Govern |
| Control Health Assessment create/review | 046 | governance conclusion | 2 | limitations visible | possible | reviewer separation possible | Govern |
| Continuous Improvement Package create | 047 | cross-product proposal | 2 | destination-safe context | possible | owner review distinct | Govern |
| provenance export preparation | 034..047 | sensitive lineage | 2 | Shared Export/classification | yes | external sharing separate | Govern + Shared/Security |
| external-sharing preparation | 038/047 | high disclosure risk | 2 | no publication without authority | yes | OPEN-019 | Govern/Security/Shared |

## Invariants
- A dashboard, metric or Audit Review never grants source-data permission.
- Cross-tenant access is never inferred from aggregate access.
- Audit Evidence Package export ≠ external-sharing authorization.
- No audit read implies permission to mutate Decision, Run, Result or source records.
- Final atomic permissions remain future Security work.