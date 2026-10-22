---
id: investigate-threat-intelligence-permissions
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements: [REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-013, OPEN-014, OPEN-015, OPEN-018, OPEN-019]
---
# Functional permission needs

| Family | Functional needs | Risk / future control |
|---|---|---|
| Session/questions/hypotheses | read, create, update, close/reopen, review/dispute | source exposure and author/reviewer separation |
| Fusion/assessments | run/read/create/update/review | dependent sources, opaque confidence, false attribution |
| Product planning/authoring | create/update/version/clone/archive/comment | sensitive-source inclusion and version confusion |
| Review/release | assign/comment/disposition/recommend | quality passed ≠ Approval; separation of duties |
| Markings/releasability | read/update/assess | product permission ≠ source permission |
| Dissemination/internal publication | prepare/publish/suspend/withdraw by policy | tenant/audience leakage; OPEN-019 |
| Consumer access | read, grant-request, revoke-request | no implicit source/raw access |
| Watchlist/operationalization | create/update/review/prepare handoff | activation/runtime remain external owner |
| Monitoring/change | create/update/run-request/read | monitoring ≠ active target surveillance |
| Feedback/effectiveness | create/read/compare/assess | personal/customer data and bias |
| Satisfaction/collection feedback | assess/reopen/prepare package | package does not execute collection |
| Correction/retraction/supersession | create/update/propose/review | no history deletion; external recall governed |
| External sharing | prepare package/Action Request only | actual sharing class 3 under Govern/destination owner |
| Automation/provenance | request proposals/read/export trace | Tool/Run attribution and human disposition |

Classes 0–2 are local subject to policy. Class 3/4 actions are unavailable locally. Atomic namespaces, RBAC/ABAC, final step-up, cross-tenant policy and final separation of duties remain future. Existence, metadata, masked preview, read, copy, extraction, product inclusion, internal publication, export, external preparation, future sharing and future recall are distinct rights.
