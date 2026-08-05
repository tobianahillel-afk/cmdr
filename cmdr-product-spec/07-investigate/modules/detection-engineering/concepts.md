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
| Match review | match ≠ malicious activity; non-match ≠ certain false negative |
| Test and replay | synthetic fixture ≠ incident; Expected Outcome ≠ actual; replay ≠ deployment |
| Coverage | mapping ≠ effective coverage; test-covered ≠ production-effective; coverage ≠ prevention |
| Review package | ≠ Approval, Decision, deployment or activation |

Detection Engineering concepts lacking an Object Register entry remain functional concepts. Their schemas, cardinalities and state machines are deferred to the Objects phase.
