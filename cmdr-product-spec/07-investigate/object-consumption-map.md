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
---
# Object consumption map — Investigate through Phase 4B.3A

| Objet ou concept | Owner actuel | Usage local | Opérations locales | Lacune | Phase propriétaire |
|---|---|---|---|---|---|
| Case / Finding / Evidence / Hypothesis / Hunt / Artifact | Investigate | need, reasoning, sources and return origin | read/link | final relation models | Objects |
| Incident / runtime Detection / Signal / Alert | Command | source, runtime and feedback projections | read/link only | authoring/runtime bridge | Command/Objects |
| Telemetry Event / Query / Search Job | Shared | source events, historical replay and search context | read/run by permission | technical contracts | Shared/Technique |
| Detection Engineering Project / Hypothesis | Investigate concepts | durable context and objective | create/update/review/supersede | canonical objects absent | Objects |
| Detection Content / Draft / Version | Investigate concepts | functional content and lineage | create/update/clone/withdraw/supersede | schemas/format absent | Objects/Technique |
| Readiness / Schema / Field / Mapping / Logic concepts | Investigate using Settings/Shared projections | authoring dependencies | read/propose/review | final models absent | Objects/Settings |
| Test Scenario / Dataset relation / Expected Outcome | Investigate with Studio/Artifact/Shared projections | controlled evidence | create/use/review | ownership/cardinality; OPEN-014 | Objects/Studio |
| Validation / Replay / Match Review / Coverage / Gap | Investigate concepts | preproduction evidence | run/read/review/supersede | canonical objects absent | Objects/Trust |
| Review Package | Investigate concept | authoring handoff | prepare/version/withdraw | ≠ Release Candidate | Objects |
| Detection Release Candidate | Investigate concept | immutable candidate for review | create/review/withdraw/supersede | canonical object absent | Objects |
| Detection Review | Investigate concept | human recommendation and comments | assign/comment/dispose | ≠ Govern Decision | Objects/Trust |
| Deployment Readiness Assessment | Investigate concept | target-specific preconditions | create/update/dispute | target contract unresolved | Objects/Settings |
| Promotion Plan | Investigate concept | ordered targets and conditions | create/update/compare | execution format absent | Objects/Govern |
| Detection Change Request Draft | Investigate concept | Govern handoff preparation | prepare/submit/cancel | ≠ Action Request accepted | Objects/Govern |
| Action Request / Decision / Approval / Response Run / Result | Govern | authority and execution projections | prepare/read/link only | final cross-run bridge; OPEN-015 | Govern/Objects |
| Environment / Deployment Target / configured runtime | Platform Settings | target selection and administrative state | read/select proposal/request | runtime strategy; OPEN-017 | Settings/Technique |
| Runtime Version Observation | Investigate concept using Command/Settings/Endpoint | expected/observed reconciliation | create/compare/dispute | observation schema absent | Objects |
| Detection Health Assessment | Investigate concept | health interpretation | create/compare/dispute | health contract future | Objects/Settings |
| Shadow Evaluation Plan / Shadow Assessment | Investigate concepts | non-alerting runtime observation | prepare/request/read/assess | runtime execution contract future | Objects/Govern |
| Canary Plan / Canary Assessment | Investigate concepts | phased target evidence | prepare/request/read/assess | rollout model future | Objects/Govern |
| Runtime Quality Assessment / Production Match Review | Investigate concepts consuming Command | quality and outcome review | create/update/reconcile | ground-truth limits | Objects/Command |
| Tuning Proposal | Investigate concept | proposed new Draft changes | create/update/withdraw | ≠ active tuning | Objects |
| Suppression Proposal / Exception Proposal | Investigate concepts | bounded Govern requests | create/update/renew/revoke request | active mechanism owner future | Objects/Govern |
| Drift Assessment / Performance Assessment | Investigate concepts | compatibility and resource review | create/update/compare | metric/runtime contracts future | Objects/Settings |
| Rollback Plan / Recovery Assessment | Investigate concepts | Govern rollback preparation and verification | create/update/read Result | ≠ rollback execution | Objects/Govern |
| Retirement Proposal / Replacement Relation | Investigate concepts | deactivation/retirement transition | create/update/compare | replacement equivalence unresolved | Objects/Govern |
| Continuous Improvement Package / Lifecycle Provenance Assessment | Investigate concepts | return to new cycle and lineage | create/update/export/handoff | no canonical schema | Objects/Shared |
| Tool / Tool Call / Workflow / Automation Run | Studio | execution and automation provenance | select/invoke/read/link | final contracts; OPEN-015 | Studio/Objects |
| Jobs / Notifications / Trace / Activity / Versioning / Linking / Comparison / Search / Export / Reporting / Collaboration / Recovery | Shared | generic mechanisms | consume/emit semantics | contracts future | Shared/Trust |
| Threat Intelligence concepts | future 4B.3B | future candidate handoff only | none canonical | phase not started | 4B.3B |

No complete schema, JSON Schema, final cardinality, object state machine, rule AST, target format, deployment package or atomic permission is defined. All new lifecycle entities remain functional concepts until the Objects phase.
