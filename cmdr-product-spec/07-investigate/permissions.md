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
  - OPEN-018
---
# Functional permission needs — Investigate

This document identifies permission families without final namespaces, RBAC/ABAC, step-up rules or a final separation-of-duties matrix.

| Family | Functional needs | Risk or separation |
|---|---|---|
| Signal / Event Search / Hunt | read, triage, execute/cancel search, raw/field access, create/manage Hunt | runtime priority Command-owned; raw and cross-tenant access distinct |
| Case / Hypothesis / Artifact / Evidence / Finding | lifecycle, author/review, link, qualify, export | author/reviewer and sensitive-content separation |
| Detection Engineering | Project/Hypothesis/Content/tests/review/readiness/plans/assessments/proposals/handoff | classes 0–2; no runtime production authority |
| Threat Intelligence Intake / Requirement / Project | read, create, update, close/reopen/archive/supersede | origin sensitivity; author/manager/reviewer separation |
| Intelligence Source Catalog / Access | source/access/health/marking read, access context assess, request prepare | Settings administers provider/feed/connector/secret/access |
| Reliability / Credibility / Confidence | assessment read/create/update/review/dispute | methods and limitations visible; no opaque truth score |
| Intelligence Material | existence, metadata, masked preview, read, copy, extraction, relation, export | independent permissions; classification/licence/OPEN-014 |
| Observable / Indicator Candidate | read/create/update/review/withdraw/revoke/supersede | candidate ≠ confirmed/deployed Indicator |
| Threat Entity / Malware / Tool / Infrastructure / Campaign knowledge | read/create/update/review/merge proposal/separate | attribution, reputation and identity risks |
| TTP / Sighting / Relationship | map/create/link/compare/dispute/withdraw | raw events, causality and tenant isolation |
| Dedup / Version / Lifecycle | compare, merge proposal/reject, version, supersede, expire/revoke proposal | no silent merge/deletion; consumer impact |
| Intelligence Handoff / Provenance | prepare Case/Hunt/Detection/Workbench/future-analysis handoff; authorized export | destination permissions re-evaluated; package ≠ Report/publication |
| Automation | extraction, normalization, relation/duplicate/assessment proposal | Tool/Run attribution and human disposition mandatory |

Classes 0–2 are available locally according to policy. External sharing, active Indicator/watchlist, Detection rule creation/deployment, block, response, source administration and irreversible deletion are unavailable and require future owner/Govern decisions. Atomic namespaces, RBAC/ABAC, final step-up and separation of duties remain deferred.
