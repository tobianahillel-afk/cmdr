---
id: investigate-detection-engineering-permissions
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements:
  - REQ-SEC-001
  - REQ-SEC-002
open_decisions:
  - OPEN-013
  - OPEN-014
  - OPEN-015
---
# Functional permission needs

| Family | Functional needs | Risk / future control |
|---|---|---|
| Project/Hypothesis | read, create, update, close, reopen, review | ownership, contributor and reviewer separation |
| Data/Schema | readiness read, schema/field read, mapping propose/review | source admin stays Settings; sensitive examples masked |
| Detection Content | read, create, update, clone, archive, metadata/logic update | reversible Class 2; no runtime mutation |
| Tests/Datasets | scenario create/update, dataset read/use, sensitive data read | minimization, source-owner permission, OPEN-014 |
| Validation/Replay | run/read/cancel, cross-tenant preparation | bounded Class 1; step-up and cost controls future |
| Match/Coverage | classify candidates, create/update coverage and gaps | candidate labels only; no guaranteed metrics |
| Review Package | prepare, export, submit to future review | no Approval or deployment; separation of duties future |
| Automation | authoring/test proposal request | Tool/Run attribution and human disposition mandatory |

Atomic namespaces, RBAC/ABAC, final step-up and production permissions are deferred. Deployment, activation, deactivation and rollback permissions are excluded.
