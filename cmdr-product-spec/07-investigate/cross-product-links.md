---
id: investigate-cross-product-links
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements:
  - REQ-PROD-004
  - REQ-PROD-005
  - REQ-PROD-008
  - REQ-INV-006
open_decisions:
  - OPEN-008
  - OPEN-013
  - OPEN-015
---
# Cross-product links — Detection Engineering Authoring

| Transition | Context | Ownership / return |
|---|---|---|
| Finding / Case Hypothesis / Hunt / Incident → Intake | sources, observations, Evidence, contradictions, limits and return origin | source owner retained; Intake Investigate |
| Static / Dynamic / Reverse / Memory / Disk / Network → Intake | Artifact, behavior, conditions, sources, limitations and provenance | technical result owner retained |
| Intake → Project / Detection Hypothesis | objective, scope, consumers, owner and source relations | Investigate |
| Project → Data Readiness → Schema Review | required sources, platforms, environments, freshness, retention, fields and gaps | Settings/Endpoint/Shared projections remain owned |
| Hypothesis / Schema Review → Detection Content Draft | behavior, selected fields, mappings, exclusions and limitations | Investigate; Query/Saved Search remain distinct |
| Draft → Validation / Tests / Replay | selected version, logic, metadata, sources, permissions and parameters | Tool/Run Studio; Query/Events Shared |
| Replay → Match Review → New Draft | matches/non-matches, source events, Expected Outcomes, gaps and dispositions | Investigate; no Signal/Alert |
| Draft / Test / Replay → Coverage / Gap | behavior, platforms, mappings, evidence level and limitations | Investigate assessment; Settings requests prepared only |
| Authoring → Future Review and Promotion | candidate version, validation, tests, replay, FP/FN review, coverage, gaps, risks and provenance | future 4B.3A.2; no Approval or deployment |
| Future runtime Detection → Command | not active in this phase | Command owner; bridge future |
| Future production authority | Review Package to applicable Govern/owner process | no Decision or Approval created here |
| Threat Intelligence | no transition executed | Phase 4B.3B not started |

Tenant, environment, selection, permissions, masking, errors, versions and return origin are preserved.
