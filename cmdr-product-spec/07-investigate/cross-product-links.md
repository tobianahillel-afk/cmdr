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
  - OPEN-014
  - OPEN-015
  - OPEN-017
  - OPEN-018
---
# Cross-product links — Detection Engineering and Threat Intelligence

| Transition | Context | Ownership / return |
|---|---|---|
| Finding / Case / Hunt / Incident / technical analysis → Threat Intelligence Intake | need, sources, Evidence, Artifacts, candidates, contradictions, restrictions and return origin | source owner retained; Intake Investigate |
| CAP-INV-435 → CAP-INV-501 | Continuous Improvement candidate context, Indicators/TTP relevance and provenance | Detection Engineering remains owner of its package; Intelligence qualifies a new intake |
| Intake → Requirement → Knowledge Project | question, scope, consumers, priority, criteria and dependencies | Investigate Threat Intelligence |
| Knowledge Project → Source Catalog / Settings request | purpose, tenant, source candidate, access restrictions and gaps | Settings owns providers/sources/secrets/access; request ≠ access |
| Material / technical outputs → candidate knowledge | originals, extraction, ambiguity, source assessment and limitations | Artifact/Tool/source owners retained |
| Threat Entity Candidate → Shared Entity review | candidate, aliases, sources, confidence and contradictions | Shared owns Entity; proposal ≠ canonical Entity |
| Intelligence relationships → Shared Graph | typed sourced candidate edges and table alternative | Shared owns Graph/Linking; edge ≠ causal truth |
| Indicator/TTP candidates → Detection Engineering | sources, Sightings, contradictions, restrictions and gaps | no Detection Content or rule created automatically |
| Sighting / candidate context → Case or Hunt | source event/Artifact, time, environment and uncertainty | destination re-evaluates permission; Sighting ≠ Signal/Incident |
| Analysis Handoff Package → future 4B.3B.2 | questions, candidates, relations, assessments, gaps and provenance | package only; no Report, attribution, dissemination or operationalization |
| Future external sharing/publication | not active | Govern owns release/Approval; phase 4B.3B.2 future |

Tenant, environment, immutable versions, markings, licence, permissions, masking, errors, partiality, timestamps, source owner and return origin are preserved.
