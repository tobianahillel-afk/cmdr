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
  - OPEN-017
---
# Functional permission needs

| Family | Functional needs | Risk / future control |
|---|---|---|
| Authoring | Project/Hypothesis/Draft/validation/test/replay/coverage permissions from CAP-INV-401..417 | reversible classes 0–2; no runtime mutation |
| Review Candidate | read/create/update/withdraw, assign reviewer, comment and disposition | immutable version, reviewer separation |
| Readiness/Targets | assessment read/create/update, environment/target read/select | Settings remains administrator; cross-tenant and compatibility step-up future |
| Govern handoff | Change Request prepare/submit/cancel; Decision/Approval/Run/Result read | no auto-approval; requester/approver/operator separation |
| Shadow/Canary | plan, request and assessment read | bounded class 1/2 locally; execution requires authority/runtime owner |
| Runtime projections | Detection/version/health/Signal/Alert/Incident feedback read | source ownership, minimization and tenant scope |
| Quality/Tuning | production review and Tuning Proposal create/update | proposal only; new Draft and tests required |
| Suppression/Exception | proposal create/update/renew/revoke request | expiry, scope, alternatives and compensating monitoring; Govern required |
| Drift/Performance | assessment create/update and source/metric projections read | sensitive infrastructure data masking |
| Rollback/Retirement | Plan/Proposal create/update; Govern execution projections read | class 3 execution unavailable locally |
| Provenance | lifecycle read/export and Continuous Improvement Package prepare | cross-product masking, audit and trace immutability |
| Automation | review/tuning/plan proposal requests | Tool/Run attribution and human disposition mandatory |

Atomic namespaces, RBAC/ABAC, final step-up and separation-of-duties rules are deferred. Investigate has no direct permission to deploy, activate, deactivate, apply an exception/suppression, roll back, retire, delete runtime content or delete provenance.
