---
id: investigate-object-consumption-map
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements:
  - REQ-PROD-006
  - REQ-PROD-014
  - REQ-PROD-061
  - REQ-PROD-062
  - REQ-INV-006
open_decisions:
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
  - OPEN-017
  - OPEN-018
---
# Object consumption map — Investigate through Phase 4B.3B.1

| Objet ou concept | Owner actuel | Usage local | Opérations locales | Lacune | Phase propriétaire |
|---|---|---|---|---|---|
| Case / Finding / Evidence / Hypothesis / Hunt / Artifact | Investigate | need, reasoning, sources and return origin | read/link | final relation models | Objects |
| Incident / runtime Detection / Signal / Alert | Command | source, runtime and feedback projections | read/link only | cross-product bridge | Command/Objects |
| Telemetry Event / Query / Search Job | Shared | source events, replay/search and Sighting context | read/run by permission | technical contracts | Shared/Technique |
| Detection Engineering Project / Hypothesis / Content / Versions | Investigate concepts | detection lifecycle and Intelligence handoff | read/link; local Detection operations remain in its module | canonical schemas/formats absent | Objects/Technique |
| Review/Readiness/Promotion/Runtime assessments and proposals | Investigate concepts | Detection lifecycle | read/link/prepare in Detection module | canonical objects absent | Objects/Govern/Settings |
| Action Request / Decision / Approval / Response Run / Result | Govern | authority and execution projections | prepare/read/link only | cross-run bridge; OPEN-015 | Govern/Objects |
| Environment / Deployment Target / configured runtime / Data Source / Parser / Provider | Platform Settings | source, target and administrative projections | read/select/request only | support and strategy; OPEN-008/017/018 | Settings/Technique |
| Tool / Tool Call / Workflow / Automation Run | Studio | execution and automation provenance | select/invoke/read/link | final contracts; OPEN-015 | Studio/Objects |
| Jobs / Notifications / Trace / Activity / Versioning / Linking / Comparison / Search / Export / Reporting / Collaboration / Recovery | Shared | generic mechanisms | consume/emit domain semantics | contracts future | Shared/Trust |
| Intelligence Requirement | Investigate concept | question, objective, priority, scope and satisfaction | create/update/pause/close/reopen/supersede | canonical object absent | Objects |
| Collection Priority | Investigate concept | priority for expected knowledge | create/update/version | not Command priority; schema absent | Objects |
| Knowledge Project | Investigate concept | durable Intelligence workspace | create/update/pause/archive/reopen/supersede | canonical object absent | Objects |
| Intelligence Source / Source Access Context | Settings projection + Investigate assessment | catalog, restrictions, health and permitted use | read/assess/request access | provider/access contracts; OPEN-008/018 | Settings/Objects |
| Source Reliability Assessment | Investigate concept | historical source-quality assessment | create/update/review/dispute/supersede | method/scale not final | Objects/Trust |
| Information Credibility Assessment | Investigate concept | claim-specific credibility | create/update/review/dispute/supersede | method/scale not final | Objects/Trust |
| Intelligence Material / Material Reference | source owner / Investigate concept | authorized material and original relation | reference/import/read/withdraw/supersede | Artifact/Material relation; OPEN-014 | Objects |
| Normalized Knowledge Record | Investigate concept | functional normalized representation | create/version/withdraw | no physical schema/exchange format | Objects/Technique |
| Observable Candidate | Investigate concept | observed/extracted value candidate | create/update/review/withdraw/supersede | canonical object/status absent | Objects |
| Indicator Candidate | Investigate concept | indicator hypothesis with sources and lifecycle | create/update/review/revoke/supersede | confirmed Indicator model future; OPEN-018 | Objects/4B.3B.2 |
| Threat Entity Candidate | Investigate concept | candidate identity/actor/persona/group | create/update/merge proposal/separate/supersede | Shared Entity relation and ontology future | Objects/Shared |
| Entity | Shared | canonical identity projection | read/propose relation only | candidate-to-Entity contract; OPEN-018 | Shared/Objects |
| Malware Family / Variant Candidate | Investigate concept | family/variant knowledge distinct from sample | create/update/review/supersede | canonical object absent | Objects |
| Malware Sample / Artifact / Derived Artifact | Investigate | technical sources | read/link | relation/cardinality; OPEN-014 | Objects |
| Tool Knowledge / Capability Knowledge | Investigate concepts | observed software/capability knowledge | create/update/version/supersede | canonical objects absent | Objects |
| Infrastructure Knowledge | Investigate concept | addresses/domains/URLs/certs/services/clusters | create/update/review/expire/revoke | physical model/ontology absent | Objects |
| Activity Cluster / Campaign / Intrusion Set Candidate | Investigate concepts | candidate grouping and context | create/update/compare/merge proposal/separate | advanced analysis/attribution future | Objects/4B.3B.2 |
| TTP Mapping / Behavior Candidate | Investigate concept | versioned sourced mapping | create/update/review/supersede | taxonomy/ontology choice; OPEN-018 | Objects |
| Sighting | Investigate concept | occurrence linked to candidate and source | create/update/dispute/withdraw/supersede | canonical object absent | Objects |
| Intelligence Relationship | Investigate concept using Shared Graph | typed sourced candidate relation | propose/review/dispute/withdraw/supersede | physical graph/cardinality absent | Objects/Shared |
| Confidence Assessment | Investigate concept | justified analytic confidence | create/update/review/supersede | no calibrated universal score | Objects/Trust |
| Contradiction Record | Investigate concept | opposing sources/interpretations | create/update/review/supersede | canonical object absent | Objects |
| Knowledge Version / Supersession relation | Investigate concept using Shared Versioning | candidate and assessment history | create/diff/supersede/restore context | final version model absent | Objects/Shared |
| Expiration Assessment / Revocation Assessment | Investigate concepts | stale/expiry/revocation lifecycle | create/review/dispute/supersede | consumer lifecycle future | Objects/4B.3B.2 |
| Analysis Handoff Package | Investigate concept | future-analysis or cross-product handoff | prepare/version/withdraw/supersede | not Report/publication; schema absent | Objects/4B.3B.2 |
| Provenance Record | Shared mechanisms + product semantics | complete source, Tool/Run and human lineage | read/emit/link/export by permission | final audit contracts | Shared/Trust |

No complete schema, JSON Schema, final cardinality, object state machine, physical graph, matching algorithm, exchange format, technical identifier, rule AST, target format, deployment package or atomic permission is defined. New Threat Intelligence entities remain functional concepts until the Objects and Technique phases.
