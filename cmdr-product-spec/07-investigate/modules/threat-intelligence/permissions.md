---
id: investigate-threat-intelligence-permissions
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
  - OPEN-018
---
# Functional permission needs

| Family | Functional needs | Risk / future control |
|---|---|---|
| Intake / Requirement / Project | read, create, update, close, reopen, archive, supersede | origin sensitivity; author/manager/reviewer separation |
| Source Catalog / Access | catalog read, access-context assess, access request prepare | Settings administration and secret access remain separate |
| Reliability / Credibility | assessment create/update/review/dispute | reputation, source identity and method visibility |
| Material | reference/import, restricted read, masked preview, extraction run/read | classification, licence, victim/customer data and OPEN-014 |
| Candidates / Knowledge | create/update/review/withdraw/revoke | false qualification, attribution and consumer impact |
| Sightings / Relations | create, link, compare, dispute, withdraw | tenant isolation, raw event access and causal overclaim |
| Confidence / Contradiction | create/update/review/supersede | opaque score and source-independence errors |
| Dedup / Version / Lifecycle | compare, merge proposal, version, supersede, expire/revoke proposal | no silent merge/deletion; consumer review |
| Handoff / Export | prepare, cross-product handoff, provenance export | markings, minimization, destination permission recheck |
| Automation | extraction/enrichment/dedup suggestion request | Tool/Run attribution and human disposition mandatory |

Classes 0–2 are local according to policy. External sharing, active Indicator/watchlist, deployment, block, source administration or destructive deletion are classes 3/4 and unavailable in this phase. Atomic namespaces, RBAC/ABAC, final step-up and final separation of duties are deferred.
