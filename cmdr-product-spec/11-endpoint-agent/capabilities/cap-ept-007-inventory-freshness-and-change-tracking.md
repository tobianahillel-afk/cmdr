---
id: CAP-EPT-007
title: Inventory Freshness and Change Tracking
product: endpoint-agent
module: foundations
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-006, REQ-PROD-012, REQ-PROD-017, REQ-PROD-018, REQ-PROD-019, REQ-OBJ-008, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008]
source-of-truth: canonical
---
# CAP-EPT-007 — Inventory Freshness and Change Tracking

## 1. Définition
Define freshness and source-backed comparison between Endpoint inventory snapshots, distinguishing stale/unknown current state from offline state and ordinary change from malicious activity.

## 2. Problème utilisateur
A snapshot can remain available after it stops representing current state, and a technical change can be overinterpreted as health failure or malicious behavior.

## 3. Objectifs
Track snapshot time/last refresh, freshness, stale state, previous snapshot ref, added/removed/modified facts, missing update and consumer projection.

## 4. Non-objectifs
No threat detection, maliciousness classification, telemetry stream, response, continuous polling implementation, SLA threshold or Fleet lifecycle.

## 5. Propriétaire
Endpoint Agent owns local inventory freshness/change semantics. Source facts retain their canonical owners.

## 6. Utilisateurs
Endpoint Operator, Platform Administrator, Investigate/SOC analyst, Command/Govern consumer and Auditor.

## 7. Conditions d’entrée
At least one source-backed snapshot and its timestamp; comparison needs a prior compatible snapshot and tenant/environment scope.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| current inventory snapshot | CAP-EPT-006 | snapshot ref | oui | explicit timestamp | freshness unknown |
| previous snapshot | CAP-EPT-006 | comparison ref | non | historical | no change diff |
| freshness policy/context | source-owned configuration or rule | assessment input | non | source-defined | stale threshold not inferred |
| connectivity/heartbeat context | CAP-EPT-009 future-in-lot | contextual fact | non | explicit | no offline conclusion |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Inventory Snapshot functional record | Endpoint Agent | current/previous facts | read |
| endpoint-agent | Endpoint Agent | Agent context | read |
| environment | Platform Settings | scope only | read projection |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Inventory Freshness/Change Assessment | derive/supersede | Endpoint Agent | derived facts, no source mutation |
| local-audit-event | record comparison/freshness transition | Endpoint Agent | provenance retained |

## 11. Fonctionnalités
Freshness state; last refresh; missing update; current-state unknown; current/previous comparison; added/removed/modified facts; consumer projection and provenance.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect freshness | authorized user | snapshot | 0 | read | timestamp/state | non |
| compare snapshots | Endpoint Operator | snapshots | 1 | compatible refs | sourced change set | non |
| request bounded refresh | authorized operator | inventory | 2 | source supports it | refresh request only | no response execution |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| calculate freshness | oui | oui when rule supplied | oui | explanation only | timestamps/rules |
| compare snapshots | oui | oui | oui | summary possible | deterministic diff |
| label change malicious | analyst workflow only | not EPT-1 | non | interdit ici | Investigate/Detection future path |

## 14. États fonctionnels
`fresh`, `stale`, `freshness-unknown`, `changed`, `unchanged`, `comparison-unavailable`, `current-state-unknown`.

## 15. États d’interface
Stale names source age; Offline is contextual, not inferred; Partial lists uncomparable facts; Permission denied masks restricted fields; Error keeps prior valid snapshot. No screen design.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Freshness State | derived Endpoint fact | Settings/Investigate/Command/Govern | threshold/source explicit or unknown |
| Inventory Change Set | derived Endpoint fact | Investigate/Settings | change is not maliciousness claim |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| current snapshot | age evaluated | freshness state | timestamp/context | snapshot retained |
| previous + current snapshots | comparison | change set | refs/changed facts | both snapshots retained |
| freshness/change | consumer projection | downstream product | scope/source/limits | no ownership transfer |

## 18. Dépendances
CAP-EPT-006 and contextual CAP-EPT-009/010 semantics, Settings configuration where supplied, Security/privacy, OPEN-008.

## 19. Source de vérité
Snapshots remain the observed records; freshness/change are derived Endpoint facts. Stale never becomes offline automatically and change never becomes Detection/Finding automatically.

## 20. Provenance et audit
Record snapshot refs, timestamps, comparison scope, rules/source, changed facts, unknowns, requester/engine and correlation id.

## 21. Permissions fonctionnelles
Inventory read and sensitive-inventory access apply to comparison inputs; no extra cross-tenant access is granted by a change/freshness projection.

## 22. Limites et erreurs
Stale != offline; freshness != health; change != malicious activity; missing previous snapshot means comparison unavailable, not unchanged.

## 23. Métriques
Fresh/stale/unknown distribution; time since last refresh as descriptive measure; changed snapshot count; comparison failures/restricted facts.

## 24. Classification de livraison
`draft / defined / planned`; no polling engine, telemetry implementation, detection logic or SLA target.

## 25. Critères d’acceptation
**Given** inventory ages beyond an explicit freshness context while heartbeat remains available, **When** evaluated, **Then** inventory can be stale without declaring the Agent offline.

**Given** a hostname/version fact changes, **When** snapshots are compared, **Then** change is recorded without labeling it malicious.

**Given** no prior snapshot exists, **When** comparison is requested, **Then** comparison is unavailable rather than unchanged.

## 26. Questions ouvertes
OPEN-008 remains open; exact source refresh behavior and platform availability are not selected.

## 27. Consommateurs documentaires
Endpoint inventory/health/state, Settings Fleet, Investigate/Command/Govern, registers, Quality, Roadmap and future EPT-2+ sources.