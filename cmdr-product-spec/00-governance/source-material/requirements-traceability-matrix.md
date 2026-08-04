---
id: requirements-traceability-matrix
domain: 00-governance
status: draft
owner: QA and Traceability Lead
updated: 2026-08-04
source-of-truth: source-material
requirements:
  - REQ-PROD-001
  - REQ-PROD-062
  - REQ-UX-001
  - REQ-UX-010
  - REQ-OBJ-001
  - REQ-OBJ-012
---

# Requirements Traceability Matrix

## Interprétation

Les 122 Requirement IDs sources sont conservés. `conform` signifie qu’une phase propriétaire actuelle fournit des règles substantielles et sans contradiction active ; cela ne prouve ni implémentation ni livraison. Phase 4A fournit les preuves Command. Phase 4B.1 ajoute les preuves Investigate pour Signals and Hunt et Cases and Evidence, sans promouvoir artificiellement les exigences dépendantes des Phases 4B.2–4E, Objets, Permissions, Parcours, Écrans ou Technique.

## Couverture

| État | Après Phase 3 | Après Phase 4A | Après Phase 4B.1 |
|---|---:|---:|---:|
| conform | 99 | 99 | 99 |
| partial | 20 | 20 | 20 |
| absent | 3 | 3 | 3 |
| contradictory | 0 | 0 | 0 |
| total | 122 | 122 | 122 |

## Traçabilité Command Phase 4A

Les 27 entrées Command, leurs owners, objets, rôles, dépendances, classes et preuves restent inchangés dans le [Capability Register](../registers/capability-register.md) et les fichiers canoniques sous `06-command/modules/*/capabilities/`. Phase 4B.1 ne modifie aucun document Command.

## Traçabilité Investigate Phase 4B.1

| Capability | Module | Nom | Owner / utilisateurs | Functional status | Delivery mode | Requirement IDs | Objets | Classes | OPEN | Preuve canonique | Écrans/parcours futurs | Dépendances principales |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| CAP-INV-001 | signals-and-hunt | Signal Triage | Investigate Product Lead; SOC Analyst, Threat Hunter | defined | planned | REQ-PROD-002,005,008,014,045 | Signal/Alert/Detection/Event projections, Case link | 0,2 | OPEN-013 | `07-investigate/modules/signals-and-hunt/capabilities/signal-triage.md` | Triage Desk; Signal→Case | CAP-INV-002,004,102; Command; Trace |
| CAP-INV-002 | signals-and-hunt | Event Search | Investigate Product Lead; Analyst, Threat Hunter | defined | planned | REQ-PROD-014,019 | Query, Search Job, Event, Case refs | 0 | — | `07-investigate/modules/signals-and-hunt/capabilities/event-search.md` | Event Search; Search→Case | CAP-INV-003,004,006,007,008; Shared Query/Jobs |
| CAP-INV-003 | signals-and-hunt | Query Authoring and Assistance | Investigate Product Lead; Analyst, Query Author | defined | planned | REQ-PROD-014; REQ-AI-002 | Query draft, schema projection | 0,2 | OPEN-015 | `07-investigate/modules/signals-and-hunt/capabilities/query-authoring-and-assistance.md` | Query editor; IA proposal review | CAP-INV-002,006; Studio optional |
| CAP-INV-004 | signals-and-hunt | Event Inspection and Pivot | Investigate Product Lead; Analyst, Threat Hunter | defined | planned | REQ-PROD-014,019 | Event, Entity, Case, candidate refs | 0,2 | — | `07-investigate/modules/signals-and-hunt/capabilities/event-inspection-and-pivot.md` | Event Inspector/mode; Event→Case | CAP-INV-002,007,102,105,107; Inspector |
| CAP-INV-005 | signals-and-hunt | Hunt Management | Investigate Product Lead; Threat Hunter, Hunt Lead | defined | planned | REQ-PROD-014,020 | Hunt workspace, Query/Search Job, Hypothesis, Case | 0,2 | — | `07-investigate/modules/signals-and-hunt/capabilities/hunt-management.md` | Hunt workspace; Hunt→Case | CAP-INV-002,003,006,008,102,103 |
| CAP-INV-006 | signals-and-hunt | Saved Searches and Query Assets | Investigate Product Lead; Query Author, Threat Hunter | defined | planned | REQ-PROD-014,019 | Query Asset, Query, versions | 0,2 | — | `07-investigate/modules/signals-and-hunt/capabilities/saved-searches-and-query-assets.md` | Saved Search library/mode | CAP-INV-002,003,005; Versioning |
| CAP-INV-007 | signals-and-hunt | Search Result Organization | Investigate Product Lead; Analyst, Threat Hunter | defined | planned | REQ-PROD-014,019 | result selection, annotations, Case links | 0,1,2 | OPEN-013 | `07-investigate/modules/signals-and-hunt/capabilities/search-result-organization.md` | result workspace; Search→Case/export | CAP-INV-002,004,102; Export |
| CAP-INV-008 | signals-and-hunt | Search and Hunt Provenance | Investigate Product Lead; Analyst, Reviewer, Auditor | defined | planned | REQ-PROD-014; REQ-AI-002 | Query/Search Job versions, Trace refs | 0,2 | OPEN-015 | `07-investigate/modules/signals-and-hunt/capabilities/search-and-hunt-provenance.md` | provenance panel; replay | CAP-INV-002,005,006; Trace/Activity |
| CAP-INV-101 | cases-and-evidence | Case Queue | Investigate Product Lead; Case Analyst, Investigation Lead | defined | planned | REQ-PROD-014 | Case, Incident projection, Finding summary | 0 | — | `07-investigate/modules/cases-and-evidence/capabilities/case-queue.md` | Case Queue→Workspace | CAP-INV-102; Search, Saved Views |
| CAP-INV-102 | cases-and-evidence | Case Lifecycle and Coordination | Investigate Product Lead; Investigation Lead, Case Analyst | defined | planned | REQ-PROD-014 | Case, Incident links, contributors | 0,2 | OPEN-013 | `07-investigate/modules/cases-and-evidence/capabilities/case-lifecycle-and-coordination.md` | Case Workspace; Incident→Case | CAP-INV-101,103,104,105,107,109,110,111 |
| CAP-INV-103 | cases-and-evidence | Hypothesis Management | Investigate Product Lead; Analyst, Reviewer | defined | planned | REQ-PROD-014; REQ-AI-002 | Hypothesis, Evidence/results refs | 0,2 | OPEN-013,015 | `07-investigate/modules/cases-and-evidence/capabilities/hypothesis-management.md` | Hypotheses & Findings | CAP-INV-102,107,108,109; Versioning |
| CAP-INV-104 | cases-and-evidence | Entity and Relationship Management | Investigate Product Lead; Analyst, Threat Hunter | defined | planned | REQ-PROD-014 | Shared Entity, relations, Case/Incident refs | 0,2 | OPEN-013 | `07-investigate/modules/cases-and-evidence/capabilities/entity-and-relationship-management.md` | Entity Graph; pivot | CAP-INV-002,004,102; Entity Resolution |
| CAP-INV-105 | cases-and-evidence | Artifact Management | Investigate Product Lead; Analyst, Specialist | defined | planned | REQ-PROD-014,061 | Artifact, versions, derivatives, Case | 0,1,2 | OPEN-014 | `07-investigate/modules/cases-and-evidence/capabilities/artifact-management.md` | Artifact Detail/workbench entry | CAP-INV-102,106,107; 4B.2 boundary |
| CAP-INV-106 | cases-and-evidence | Attachment Handling | Investigate Product Lead; Contributor, Report Author | proposed | planned | REQ-PROD-061 | Attachment, Note/Comment/Report refs | 0,2 | OPEN-014 | `07-investigate/modules/cases-and-evidence/capabilities/attachment-handling.md` | collaboration/reporting | CAP-INV-105,107,111,114; OPEN-014 |
| CAP-INV-107 | cases-and-evidence | Evidence Creation and Management | Investigate Product Lead; Analyst, Evidence Reviewer | defined | planned | REQ-PROD-014,061,062 | Evidence, Artifact/source, Case, Hypothesis/Finding refs | 0,1,2 | OPEN-013 | `07-investigate/modules/cases-and-evidence/capabilities/evidence-creation-and-management.md` | Evidence Board; Artifact→Evidence | CAP-INV-102,103,105,108,109; Provenance |
| CAP-INV-108 | cases-and-evidence | Evidence Review and Qualification | Investigate Product Lead; Reviewer, Analyst | defined | planned | REQ-PROD-014,062 | Evidence dimensions, Hypothesis refs | 0,1,2 | OPEN-013 | `07-investigate/modules/cases-and-evidence/capabilities/evidence-review-and-qualification.md` | Evidence review; collection request future | CAP-INV-103,107,109; 4B.2 boundary |
| CAP-INV-109 | cases-and-evidence | Finding Management | Investigate Product Lead; Analyst, Reviewer | defined | planned | REQ-PROD-014,016 | Finding, Evidence, Incident/Request refs | 0,2 | OPEN-013,015 | `07-investigate/modules/cases-and-evidence/capabilities/finding-management.md` | Findings; Finding→Action Request | CAP-INV-103,107,108,113,114 |
| CAP-INV-110 | cases-and-evidence | Investigation Timeline | Investigate Product Lead; Case Analyst, Reviewer | defined | planned | REQ-PROD-014,018 | Timeline Entry, Case objects, Govern projections | 0,1,2 | — | `07-investigate/modules/cases-and-evidence/capabilities/investigation-timeline.md` | Case Timeline; replay/report | CAP-INV-102,107,109,112,114; Timeline Engine |
| CAP-INV-111 | cases-and-evidence | Case Collaboration and Investigation Notes | Investigate Product Lead; Case Analyst, Contributors | defined | planned | REQ-PROD-014,018 | Note, Comment, Attachment, Case/Task refs | 0,2 | OPEN-013,014 | `07-investigate/modules/cases-and-evidence/capabilities/case-collaboration-and-investigation-notes.md` | Case collaboration; Note→Task | CAP-INV-102,105,106,107; Shared Collaboration |
| CAP-INV-112 | cases-and-evidence | Case Replay and Investigation Review | Investigate Product Lead; Investigation Lead, Reviewer | defined | planned | REQ-PROD-014,020 | Case snapshot, Query/Hypothesis/Evidence/Finding versions | 0,2 | OPEN-013,015 | `07-investigate/modules/cases-and-evidence/capabilities/case-replay-and-investigation-review.md` | replay/review; future Detection/Readiness | CAP-INV-008,102,103,107,109,110 |
| CAP-INV-113 | cases-and-evidence | Action Request Preparation | Investigate Product Lead; Investigation Lead, Senior Analyst | defined | planned | REQ-PROD-014,016 | Govern Action Request draft, Findings, Evidence, Case | 0,2,3 | OPEN-007,013,015 | `07-investigate/modules/cases-and-evidence/capabilities/action-request-preparation.md` | Finding→Govern→Result | CAP-INV-107,108,109; Govern; Command impact |
| CAP-INV-114 | cases-and-evidence | Case Reporting Preparation | Investigate Product Lead; Author, Reviewer | defined | planned | REQ-PROD-014,018 | Report draft, citations, Findings/Evidence/Timeline refs | 0,1,2 | OPEN-014 | `07-investigate/modules/cases-and-evidence/capabilities/case-reporting-preparation.md` | Investigation Report; export future | CAP-INV-109,110,111; Reporting/Export |

## Résultats Phase 4B.1

- `REQ-PROD-014`, `REQ-PROD-016`, `REQ-PROD-019`, `REQ-PROD-061`, `REQ-PROD-062`, `REQ-AI-002` et plusieurs exigences objets/UX gagnent des preuves fonctionnelles détaillées sans changer leur classement déjà conforme.
- `REQ-PROD-006`, `REQ-PROD-007`, `REQ-PROD-008`, `REQ-PROD-010`, `REQ-PROD-020`, `REQ-AI-007` à `010`, `REQ-OBJ-009`, `REQ-SEC-003` à `005` et `REQ-UX-010` restent partiels car leurs phases propriétaires ne sont pas achevées.
- `REQ-JRN-002`, `REQ-JRN-005` et `REQ-JRN-008` restent absents et relèvent de la phase Parcours.
- aucun Requirement ID nouveau ; aucune contradiction active ; aucune capability `planned` présentée comme livrée.

## Inventaire complet

### conform — 99

`REQ-AI-001`, `REQ-AI-002`, `REQ-AI-003`, `REQ-AI-004`, `REQ-AI-005`, `REQ-AI-006`, `REQ-AI-011`, `REQ-BRAND-001`, `REQ-BRAND-002`, `REQ-BRAND-003`, `REQ-BRAND-004`, `REQ-BRAND-005`, `REQ-BRAND-006`, `REQ-BRAND-008`, `REQ-INV-001`, `REQ-INV-002`, `REQ-INV-003`, `REQ-INV-004`, `REQ-INV-005`, `REQ-INV-006`, `REQ-OBJ-001`, `REQ-OBJ-002`, `REQ-OBJ-003`, `REQ-OBJ-004`, `REQ-OBJ-005`, `REQ-OBJ-006`, `REQ-OBJ-007`, `REQ-OBJ-008`, `REQ-OBJ-010`, `REQ-OBJ-011`, `REQ-OBJ-012`, `REQ-PROD-001`, `REQ-PROD-002`, `REQ-PROD-003`, `REQ-PROD-004`, `REQ-PROD-005`, `REQ-PROD-009`, `REQ-PROD-011`, `REQ-PROD-012`, `REQ-PROD-013`, `REQ-PROD-014`, `REQ-PROD-015`, `REQ-PROD-016`, `REQ-PROD-017`, `REQ-PROD-018`, `REQ-PROD-019`, `REQ-PROD-021`, `REQ-PROD-022`, `REQ-PROD-023`, `REQ-PROD-024`, `REQ-PROD-025`, `REQ-PROD-026`, `REQ-PROD-027`, `REQ-PROD-028`, `REQ-PROD-029`, `REQ-PROD-030`, `REQ-PROD-031`, `REQ-PROD-032`, `REQ-PROD-033`, `REQ-PROD-034`, `REQ-PROD-035`, `REQ-PROD-036`, `REQ-PROD-037`, `REQ-PROD-038`, `REQ-PROD-039`, `REQ-PROD-040`, `REQ-PROD-041`, `REQ-PROD-042`, `REQ-PROD-043`, `REQ-PROD-044`, `REQ-PROD-045`, `REQ-PROD-046`, `REQ-PROD-047`, `REQ-PROD-048`, `REQ-PROD-049`, `REQ-PROD-050`, `REQ-PROD-051`, `REQ-PROD-052`, `REQ-PROD-053`, `REQ-PROD-054`, `REQ-PROD-055`, `REQ-PROD-056`, `REQ-PROD-057`, `REQ-PROD-058`, `REQ-PROD-059`, `REQ-PROD-060`, `REQ-PROD-061`, `REQ-PROD-062`, `REQ-SEC-001`, `REQ-SEC-002`, `REQ-UX-001`, `REQ-UX-002`, `REQ-UX-003`, `REQ-UX-004`, `REQ-UX-005`, `REQ-UX-006`, `REQ-UX-007`, `REQ-UX-008`, `REQ-UX-009`

### partial — 20

`REQ-AI-007`, `REQ-AI-008`, `REQ-AI-009`, `REQ-AI-010`, `REQ-BRAND-007`, `REQ-JRN-001`, `REQ-JRN-003`, `REQ-JRN-004`, `REQ-JRN-006`, `REQ-JRN-007`, `REQ-OBJ-009`, `REQ-PROD-006`, `REQ-PROD-007`, `REQ-PROD-008`, `REQ-PROD-010`, `REQ-PROD-020`, `REQ-SEC-003`, `REQ-SEC-004`, `REQ-SEC-005`, `REQ-UX-010`

### absent — 3

`REQ-JRN-002`, `REQ-JRN-005`, `REQ-JRN-008`

### contradictory — 0

Aucun Requirement ID n’est contradictoire après la consolidation Command et Investigate.

## Règle de maintenance

Une capability peut être fonctionnellement `defined` et techniquement `planned`. Seule une preuve approuvée de moteur, contrat et release permet une promotion de delivery mode.
