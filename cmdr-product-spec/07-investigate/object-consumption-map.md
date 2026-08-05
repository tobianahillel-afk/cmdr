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
---
# Object consumption map — Investigate through Phase 4B.3A.1

| Objet ou concept | Owner actuel | Usage local | Opérations locales | Lacune | Phase propriétaire |
| --- | --- | --- | --- | --- | --- |
| Case | Investigate | need source, context and return origin | read/link | final cardinality | Objects |
| Incident | Command | operational source and consumer context | read/link only | cross-product contract | Command/Objects |
| Finding / Evidence / Hypothesis / Hunt / Artifact | Investigate | reasoning, sources and technical handoffs | read/link | final relation models | Objects |
| Telemetry Event / Query / Search Job | Shared | source events, historical replay and search context | read/run according to permission | technical contracts | Shared/Technique |
| runtime Detection / Signal / Alert | Command | future runtime projection only | no local create/update/delete | authoring-to-runtime bridge | 4B.3A.2/Command |
| Detection Engineering Project | Investigate concept | durable authoring workspace | create/update/close/reopen/archive/supersede | canonical object absent | Objects |
| Detection Hypothesis | Investigate concept | behavior and authoring objective | create/update/review/dispute | canonical object absent | Objects |
| Detection Content / Draft / Version | Investigate concepts | functional content and version lineage | create/update/clone/withdraw/supersede | schemas and format absent | Objects/Technique |
| Data Source Requirement / Telemetry Readiness Assessment | Investigate concepts using Settings projections | required sources and sufficiency review | create/update/compare | objects absent | Objects/Settings |
| Schema Reference / Field Reference / Field Mapping | Settings/Shared projections plus Investigate review | semantic dependency selection | read/propose/review | final model absent | Objects/Settings |
| Detection Condition / Correlation / Sequence / Threshold | Investigate concepts | functional logic | create/update/version | no AST or language | Objects/Technique |
| Enrichment Requirement | Investigate concept | required/optional context and fallback | create/update/request owner | object absent | Objects |
| Detection Test Scenario / Expected Outcome | Investigate concepts | controlled expectations | create/update/review/supersede | objects absent | Objects/Trust |
| Detection Test Dataset | Studio/Artifact/Shared projection plus Investigate usage | versioned test data relation | read/use/link | ownership/cardinality unresolved; OPEN-014 | Studio/Objects |
| Validation Result / Replay Result / Match Review | Investigate concepts | authoring quality results | run/read/review/supersede | objects absent | Objects/Trust |
| Detection Coverage Assessment / Detection Gap | Investigate concepts | coverage evidence and missing capability | create/update/compare | objects absent | Objects |
| Review Package | Investigate concept | future 4B.3A.2 handoff | prepare/version/submit/withdraw | not Approval/Decision | Objects/4B.3A.2 |
| Tool / Tool Call / Automation Run / Workflow | CMDR Studio | execution and automation provenance | select/invoke/read/link | final contracts; OPEN-015 | Studio/Objects |
| Data Source / Parser / Schema / Health / Retention / Environment | Platform Settings | administrative projections | read/request only | support and SLOs; OPEN-008 | Settings/Objects |
| Endpoint telemetry capability | Endpoint Agent / Settings | declared platform and field support | read projection | support list; OPEN-008 | Endpoint/Settings |
| Decision / Approval / Response Run / Result | Govern | future production authority | read/link/prepare only | future review bridge | Govern |
| Jobs / Notifications / Trace / Activity / Versioning / Linking / Search / Export / Reporting / Collaboration / Recovery | Shared | generic mechanisms | consume/emit semantics | contracts future | Shared/Trust |
| Threat Intelligence concepts | future 4B.3B | not consumed or created in 4B.3A.1 | none | phase not started | 4B.3B |

No complete schema, JSON Schema, final cardinality, object state machine, rule AST, field model, package format or atomic permission is defined. Detection Engineering concepts remain functional until the Objects phase decides canonical status.
