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
---
# Functional permission needs — Investigate

This document identifies permission families without final namespaces or RBAC/ABAC.

| Family | Functional needs | Risk or separation |
|---|---|---|
| Signal / Event Search / Hunt | read, triage, execute/cancel search, raw/field access, create/manage Hunt | runtime priority Command-owned; raw and cross-tenant access distinct |
| Case / Hypothesis / Artifact / Evidence / Finding | lifecycle, author/review, link, qualify, export | author/reviewer and sensitive content separation |
| Detection Project / Hypothesis | read, create, update, close/reopen, review | contributor and owner separation |
| Data / Schema | readiness read, schema/field read, mapping propose/review | Settings administration remains separate |
| Detection Content | read, create, update, clone, archive, metadata/condition/correlation/sequence/window/enrichment update | reversible Class 2; no runtime mutation |
| Test / Dataset / Expected Outcome | scenario create/update, dataset use, sensitive data read, oracle review | source-owner permission, masking, OPEN-014 |
| Validation / Replay | run/read/cancel, restricted/cross-tenant preparation | Class 1 guardrails, cost and step-up future |
| Match / Coverage / Gap | candidate TP/FP/FN classify, assessment/gap create/update | labels remain uncertain; no automatic approval |
| Review Package | prepare, export and submit to future owner | no Approval, promotion or deployment |
| Automation | authoring/test proposal request | Tool/Run attribution and human disposition mandatory |

Deployment, activation, deactivation, rollback, active exception and production permissions are excluded. Atomic namespaces, final step-up and separation of duties remain future.
