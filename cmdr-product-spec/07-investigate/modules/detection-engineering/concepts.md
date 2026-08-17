---
id: investigate-detection-engineering-concepts
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
---
# Concepts and mandatory distinctions

| Concept | Distinction |
|---|---|
| Detection Content Draft | ≠ runtime Detection, deployed rule, Search Query, Saved Search, Workflow or Tool |
| Detection Engineering Project | ≠ Case and ≠ Automation Run |
| Detection Hypothesis | ≠ Case Hypothesis and ≠ confirmation of attack |
| Data readiness | configured ≠ healthy ≠ complete ≠ sufficient |
| Field mapping | field present ≠ reliable; mapping ≠ certain semantic equivalence |
| Validation | syntax/structure ≠ semantic correctness ≠ useful detection ≠ production performance |
| Test and replay | fixture ≠ ground truth; Expected Outcome ≠ actual; replay ≠ shadow or deployment |
| Review Package | ≠ Release Candidate |
| Release Candidate | ≠ active version and ≠ runtime Detection |
| Detection Review | ≠ Govern Decision or Approval |
| Readiness / Promotion Plan | ≠ deployment or execution |
| Action Request / Decision / Approval / Response Run / Result | distinct Govern objects and stages |
| Deployment | deployed ≠ active; active ≠ healthy; healthy ≠ useful/effective |
| Shadow / Canary | shadow match ≠ Signal; canary success ≠ global success |
| Runtime feedback | Command disposition ≠ absolute ground truth; volume ≠ quality |
| Tuning / Suppression / Exception | proposal ≠ active change; suppression ≠ historical deletion; exception ≠ permanent bypass |
| Drift / Performance | change ≠ failure; latency ≠ accuracy; cost ≠ value |
| Rollback / Retirement | requested ≠ completed; rolled back ≠ recovered; inactive ≠ retired ≠ deleted |
| Continuous Improvement Package | ≠ active modification or canonical Report |

Concepts lacking an Object Register entry remain functional concepts. Schemas, cardinalities, state machines, ASTs, runtime formats and permissions remain deferred.
