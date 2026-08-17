---
id: investigate-cloud-analysis-permissions
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements: [REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-012, OPEN-013]
---
# Cloud Analysis permissions

## Functional permission needs
- Cloud scope read/select;
- source and restricted-source read;
- Session create/update/close/reopen;
- inventory, identity, IAM policy, permission path and audit-log read;
- configuration, compute/workload, container/orchestration, serverless, network and storage metadata read;
- sensitive-material presence, metadata, masked preview, reveal, copy and export as distinct rights;
- Timeline read and correlation create;
- anomaly/Hypothesis create, review, dispute, withdraw and supersede;
- Derived Artifact create/export;
- Evidence Candidate, Finding Draft and Detection/TI/Collection handoff prepare;
- provenance read/export;
- cross-tenant analysis.

## Rules
A Session permission does not grant source permission. Metadata permission does not grant content permission. Reveal does not grant copy, export or use. Cross-tenant correlation requires explicit scope and authority. Classes 3 and 4 are not executable by Investigate.

## Deferred
Final namespaces, atomic permissions, RBAC/ABAC, step-up thresholds and separation-of-duties matrix remain future Security/Permissions work.
