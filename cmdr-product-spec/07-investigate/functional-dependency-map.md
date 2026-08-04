---
id: investigate-functional-dependency-map
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-006
  - REQ-PROD-012
  - REQ-PROD-014
open_decisions:
  - OPEN-007
  - OPEN-013
  - OPEN-014
  - OPEN-015
---

# Functional dependency map — Investigate Phase 4B.1

| Capability | Capability dependencies | Product / Shared dependencies | Object dependencies | OPEN / future dependency | Failure behavior |
|---|---|---|---|---|---|
| CAP-INV-001 | 002,004,102 | Command, Trace, Inspector | Signal, Alert, Detection, Event, Incident, Case | OPEN-013 | triage partial ; aucune priorité ou Evidence inventée |
| CAP-INV-002 | 003,004,006,007,008 | Query/Search Job, data sources, Jobs, Saved Views | Query, Search Job, Event, Case | moteur/langage futurs | résultats partiels et erreurs par source visibles |
| CAP-INV-003 | 002,006 | schema discovery, Studio optional, Documentation | Query, Query Asset | OPEN-015 | saisie manuelle et validation déterministe restent disponibles |
| CAP-INV-004 | 002,007,102,105,107 | Inspector, Entity Resolution, Trace | Event, Entity, Case, Artifact/Evidence candidates | 4B.2 pour analyse technique | Event reste inspectable ; aucune qualification automatique |
| CAP-INV-005 | 002,003,006,008,102,103 | Collaboration, Presence, Trace | Query, Search Job, Hypothesis, Case, Incident | décision future objet Hunt | workspace reste draft/partial sans objet Hunt canonique |
| CAP-INV-006 | 002,003,005 | Versioning, sharing, Saved Views distinction | Query, Query Asset | future Detection Rule 4B.3 | asset non exécutable si prérequis manquant |
| CAP-INV-007 | 002,004,102 | Export, Object Linking, Inspector | Event results, Case | OPEN-013 | sélection locale conservée ; export bloqué si permission absente |
| CAP-INV-008 | 002,005,006 | Trace, Activity Stream, Audit Hooks | Query, Search Job, Automation Run refs | OPEN-015 | provenance partielle explicitement signalée |
| CAP-INV-101 | 102 | Search, Saved Views, Queue Shell, Context Bar | Case, Incident projection, Finding summary | — | file spécialisée reste utilisable avec projections partielles |
| CAP-INV-102 | 101,103,104,105,107,109,110,111 | Command Incident, Collaboration, Object Linking | Case, Incident, Task refs | OPEN-013, phase Objets | état draft conceptuel ; aucun cycle final imposé |
| CAP-INV-103 | 102,107,108,109 | Versioning, Trace, Studio optional | Hypothesis, Evidence, results | OPEN-013,015 | Hypothesis reste proposed/inconclusive ; aucun Finding automatique |
| CAP-INV-104 | 002,004,102 | Entity Resolution, Graph, Object Linking | Entity, Event, Case, Incident | OPEN-013, identity future | conflit visible ; aucune fusion silencieuse |
| CAP-INV-105 | 102,106,107 | Versioning, preview, future workbench | Artifact, Attachment relation, Case | OPEN-014, 4B.2 | Artifact enregistré sans analyse ; dérivé indisponible signalé |
| CAP-INV-106 | 105,107,111,114 | Notes/Comments, Reporting, Attachments Shared | Attachment, Artifact, Evidence | OPEN-014 | Attachment reste documentaire ; promotion bloquée si contrat absent |
| CAP-INV-107 | 102,103,105,108,109 | Provenance, Trace, Export, Audit Hooks | Evidence, Artifact/source, Case | OPEN-013, phase Trust | candidate ou Evidence partial ; aucune mutation silencieuse |
| CAP-INV-108 | 103,107,109 | review, Versioning, future collection request | Evidence, Hypothesis | OPEN-013, 4B.2 demande collecte | qualification distincte par dimension ; collecte non exécutée |
| CAP-INV-109 | 103,107,108,113,114 | review, Trace, Object Linking | Finding, Evidence, Incident refs | OPEN-013,015 | Finding reste draft/disputed ; aucune Decision créée |
| CAP-INV-110 | 102,107,109,112,114 | Timeline Engine, Trace, Context preservation | Timeline Entry, Case objects, Govern projections | phase Objets | observed/inferred séparés ; source inaccessible masquée |
| CAP-INV-111 | 102,105,106,107 | Notes, Comments, Presence, Notifications, Task handoff | Note, Comment, Attachment, Task ref | OPEN-013,014 | brouillon conservé ; aucune Task locale concurrente |
| CAP-INV-112 | 008,102,103,107,109,110 | Versioning, Activity Stream, Studio projections | Case review inputs, Automation Run refs | OPEN-013,015 ; 4B.3/Readiness futurs | review partial et handoff conceptuel seulement |
| CAP-INV-113 | 107,108,109 | Govern, Command impact, Studio workflow optional | Action Request, Finding, Evidence, Case | OPEN-007,013,015 | draft conservé ; aucune Decision, Approval ou exécution locale |
| CAP-INV-114 | 109,110,111 | Reporting Engine, Versioning, Export, Redaction | Report, Finding, Evidence, Timeline refs | OPEN-014, export futur | draft conservé ; aucune publication ou format final simulé |

## Règles de dépendance

- une dépendance ou projection ne transfère jamais l’ownership ;
- source, fraîcheur, permission et version accompagnent chaque projection ;
- l’indisponibilité d’une dépendance conserve les objets Investigate valides et produit un état Partial, Stale ou Error explicite ;
- aucune IA ni Automation Agent n’est nécessaire à un workflow essentiel ;
- les dépendances vers 4B.2 et 4B.3 sont des frontières et handoffs, pas des capabilities définies ou livrées ;
- aucun protocole, API, moteur, stockage, SLO ou cardinalité n’est défini dans cette carte.
