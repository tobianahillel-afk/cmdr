---
id: investigate-capability-map
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-012
  - REQ-PROD-014
  - REQ-PROD-019
---

# Capability map — Investigate Phase 4B.1

## Capabilities retenues

| ID | Capability | Module | Owner | Status | Delivery status | Delivery mode | Requirement IDs | OPEN | Dépendances principales | Futurs écrans |
|---|---|---|---|---|---|---|---|---|---|---|
| CAP-INV-001 | Signal Triage | Signals and Hunt | Investigate Product Lead | draft | defined | planned | REQ-PROD-002,005,008,014,045 | OPEN-013 | Command Signal/Alert/Incident, Trace | Triage Desk |
| CAP-INV-002 | Event Search | Signals and Hunt | Investigate Product Lead | draft | defined | planned | REQ-PROD-014,019 | — | Query, Search Job, data sources, Saved Views | Event Search |
| CAP-INV-003 | Query Authoring and Assistance | Signals and Hunt | Investigate Product Lead | draft | defined | planned | REQ-PROD-014; REQ-AI-002 | OPEN-015 | Query schema, autocomplete, Studio optional | Event Search query editor |
| CAP-INV-004 | Event Inspection and Pivot | Signals and Hunt | Investigate Product Lead | draft | defined | planned | REQ-PROD-014,019 | — | Telemetry Event, Entity, Inspector, Trace | Event Inspector / Event Search mode |
| CAP-INV-005 | Hunt Management | Signals and Hunt | Investigate Product Lead | draft | defined | planned | REQ-PROD-014,020 | — | Query/Search Job, Hypothesis, Case, Collaboration | Hunt workspace or mode |
| CAP-INV-006 | Saved Searches and Query Assets | Signals and Hunt | Investigate Product Lead | draft | defined | planned | REQ-PROD-014,019 | — | Query, Versioning, Saved Views distinction | Saved Search library/mode |
| CAP-INV-007 | Search Result Organization | Signals and Hunt | Investigate Product Lead | draft | defined | planned | REQ-PROD-014,019 | OPEN-013 | Search results, Case links, Export | Event Search result workspace |
| CAP-INV-008 | Search and Hunt Provenance | Signals and Hunt | Investigate Product Lead | draft | defined | planned | REQ-PROD-014; REQ-AI-002 | OPEN-015 | Trace, Activity Stream, Query/Search Job versions | Search/Hunt provenance panel |
| CAP-INV-101 | Case Queue | Cases and Evidence | Investigate Product Lead | draft | defined | planned | REQ-PROD-014 | — | Case, Saved Views, Search, Context Bar | Case Queue |
| CAP-INV-102 | Case Lifecycle and Coordination | Cases and Evidence | Investigate Product Lead | draft | defined | planned | REQ-PROD-014 | OPEN-013 | Incident projection, Case objects, Collaboration | Case Workspace |
| CAP-INV-103 | Hypothesis Management | Cases and Evidence | Investigate Product Lead | draft | defined | planned | REQ-PROD-014; REQ-AI-002 | OPEN-013,015 | Evidence, Search results, Versioning | Hypotheses & Findings |
| CAP-INV-104 | Entity and Relationship Management | Cases and Evidence | Investigate Product Lead | draft | defined | planned | REQ-PROD-014 | OPEN-013 | Entity Resolution, Graph, Object Linking | Entity Graph |
| CAP-INV-105 | Artifact Management | Cases and Evidence | Investigate Product Lead | draft | defined | planned | REQ-PROD-014,061 | OPEN-014 | Artifact, Case, Versioning, future workbench | Artifact Detail / workbench entry |
| CAP-INV-106 | Attachment Handling | Cases and Evidence | Investigate Product Lead | draft | proposed | planned | REQ-PROD-061 | OPEN-014 | Attachments, Notes/Comments, Artifact qualification | Case collaboration / reporting |
| CAP-INV-107 | Evidence Creation and Management | Cases and Evidence | Investigate Product Lead | draft | defined | planned | REQ-PROD-014,061,062 | OPEN-013 | Artifact/other source, Case, Provenance, Export | Evidence Board |
| CAP-INV-108 | Evidence Review and Qualification | Cases and Evidence | Investigate Product Lead | draft | defined | planned | REQ-PROD-014,062 | OPEN-013 | Evidence, Hypothesis, collection request future | Evidence Board review mode |
| CAP-INV-109 | Finding Management | Cases and Evidence | Investigate Product Lead | draft | defined | planned | REQ-PROD-014,016 | OPEN-013,015 | Evidence, Incident link, Action Request | Hypotheses & Findings |
| CAP-INV-110 | Investigation Timeline | Cases and Evidence | Investigate Product Lead | draft | defined | planned | REQ-PROD-014,018 | — | Timeline Engine, Case objects, Govern projections | Case Timeline |
| CAP-INV-111 | Case Collaboration and Investigation Notes | Cases and Evidence | Investigate Product Lead | draft | defined | planned | REQ-PROD-014,018 | OPEN-013,014 | Notes, Comments, Presence, Task handoff | Case Workspace collaboration |
| CAP-INV-112 | Case Replay and Investigation Review | Cases and Evidence | Investigate Product Lead | draft | defined | planned | REQ-PROD-014,020 | OPEN-013,015 | Timeline, versions, future Detection/Readiness | Case review/replay mode |
| CAP-INV-113 | Action Request Preparation | Cases and Evidence | Investigate Product Lead | draft | defined | planned | REQ-PROD-014,016 | OPEN-007,013,015 | Govern Action Request, Command impact, Findings/Evidence | Case action-request flow |
| CAP-INV-114 | Case Reporting Preparation | Cases and Evidence | Investigate Product Lead | draft | defined | planned | REQ-PROD-014,018 | OPEN-014 | Reporting Engine, citations, redaction, Export | Investigation Report |

## Comptage

- Signals and Hunt : **8 capabilities**.
- Cases and Evidence : **14 capabilities**.
- Total Phase 4B.1 : **22 capabilities**.
- Delivery status : **20 defined, 1 proposed (CAP-INV-106)** plus les documents de cadre en draft ; aucune capability n’est déclarée livrée.
- Delivery mode : **22 planned**.

## Plages réservées, non commencées

| Plage réservée | Module | Sous-phase | État |
|---|---|---|---|
| CAP-INV-201..299 | Collection and Live Response | 4B.2 | non commencée ; aucun ID attribué |
| CAP-INV-301..399 | Analysis Workbench | 4B.2 | non commencée ; aucun ID attribué |
| CAP-INV-401..499 | Detection Engineering | 4B.3 | non commencée ; aucun ID attribué |
| CAP-INV-501..599 | Intelligence | 4B.3 | non commencée ; aucun ID attribué |

Aucun ID deprecated n’est recyclé et aucune capability future n’est présentée comme définie ou livrée.
