---
id: 07-investigate-permissions
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements:
  - REQ-PROD-006
  - REQ-SEC-001
  - REQ-SEC-002
  - REQ-SEC-003
  - REQ-INV-006
open_decisions:
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
  - OPEN-017
---
# Functional permission needs — Investigate

This document identifies permission families without final namespaces, RBAC/ABAC, step-up rules or a final separation-of-duties matrix.

| Family | Functional needs | Risk or separation |
|---|---|---|
| Signal / Event Search / Hunt | read, triage, execute/cancel search, raw/field access, create/manage Hunt | runtime priority Command-owned; raw and cross-tenant access distinct |
| Case / Hypothesis / Artifact / Evidence / Finding | lifecycle, author/review, link, qualify, export | author/reviewer and sensitive-content separation |
| Detection Project / Hypothesis | read, create, update, close/reopen, review | contributor and owner separation |
| Data / Schema | readiness read, schema/field read, mapping propose/review | Settings administration remains separate |
| Detection Content | read, create, update, clone, archive, metadata/condition/correlation/sequence/window/enrichment update | reversible Class 2; no runtime mutation |
| Test / Dataset / Expected Outcome | scenario create/update, dataset use, sensitive-data read, oracle review | source-owner permission, masking, OPEN-014 |
| Validation / Historical Replay | run/read/cancel, restricted/cross-tenant preparation | bounded Class 1; cost and future step-up controls |
| Match / Coverage / Gap | candidate TP/FP/FN classify, assessment/gap create/update | uncertain labels; no automatic approval |
| Review Package / Release Candidate | prepare, read, create, assign review, comment, disposition, withdraw | review ≠ Approval; immutable selected version |
| Deployment Readiness / Targets | readiness read/create/update, target projections read/select, Promotion Plan create/update | Settings owns environments/targets; compatibility evaluated per target |
| Govern Handoff | Change Request Draft prepare/submit/cancel; Decision, Approval, Response Run and Result read | requester/reviewer/approver/operator separation; no self-approval |
| Shadow / Canary | plan and request authorized observation, read assessment and per-target results | shadow creates no Signal; canary advancement is governed |
| Runtime / Health | runtime Detection/version/state/health read, reconciliation and assessment | Command/Settings/Endpoint source ownership; health ≠ effectiveness |
| Command Feedback / Production Review | Signal/Alert/Incident feedback read, production match review create/update | operational disposition is evidence, not absolute ground truth |
| Tuning / Suppression / Exception | proposal create/update/withdraw and Govern handoff | proposal ≠ active change; expiry and compensating controls required |
| Drift / Performance | assessment create/update, metrics read, revalidation/capacity request prepare | sensitive infrastructure projections scoped; no automatic active fix |
| Rollback / Recovery | Rollback Plan create/update, request submit, Run/Result read, recovery assessment | execution Class 3 remains Govern/runtime-owned |
| Deactivation / Retirement / Replacement | proposal create/update, consumer/coverage read, Action Request preparation | inactive ≠ retired ≠ deleted; historical lineage preserved |
| Lifecycle provenance | cross-product lineage read/export, Continuous Improvement Package prepare | evidence minimization, owner boundaries and audit integrity |
| Automation | review/tuning/plan proposal request | Tool/Run attribution and human disposition mandatory |

Classes 0–2 are available locally according to policy. Class-3 promotion, activation, deactivation, application of suppression/exception, production rollback and active-version mutation require Govern and the runtime owner. Class 4 is denied by default; provenance and historical Signals/Alerts are never deleted by Investigate.

Atomic namespaces, RBAC/ABAC, final step-up and final separation of duties remain deferred.
