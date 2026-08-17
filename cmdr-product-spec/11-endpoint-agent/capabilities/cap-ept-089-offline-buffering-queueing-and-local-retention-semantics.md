---
id: CAP-EPT-089
title: Offline Buffering, Queueing and Local Retention Semantics
product: endpoint-agent
module: resilience
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-12
requirement_ids: [REQ-PROD-005, REQ-PROD-006, REQ-PROD-012, REQ-PROD-017, REQ-PROD-018, REQ-PROD-019, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004]
open_decisions: [OPEN-008]
source-of-truth: canonical
---
# CAP-EPT-089 — Offline Buffering, Queueing and Local Retention Semantics

## 1. Définition
Définir les faits fonctionnels de buffering, queueing et rétention locale en mode connecté ou offline sans définir moteur de stockage ni garantie de durabilité infinie.

## 2. Problème utilisateur
Des données locales peuvent être retenues, sous pression ou perdues sans que consommateurs puissent distinguer buffered, queued, retained, dropped et delivered.

## 3. Objectifs
Represent buffered item, queue state, queued/retained/drop indication, capacity pressure, conceptual oldest/newest timestamps, retention limitation, degraded state and provenance.

## 4. Non-objectifs
Aucun moteur de stockage, schéma physique de queue, chiffrement concret, capacité/SLO final, delivery guarantee, replay ou exactly-once.

## 5. Propriétaire
Endpoint owns local buffer/queue facts; Shared retains generic Jobs/Retry mechanisms and Security owns global storage/privacy requirements.

## 6. Utilisateurs
Endpoint Operator, SOC/Investigate consumer, Platform Administrator, Security Reviewer, Auditor.

## 7. Conditions d’entrée
Agent/source identity, connectivity state, local item/source refs, queue/retention policy projection where applicable and capacity observations.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| local item/source | Endpoint telemetry/output | retained candidate | oui | event time | no queue item |
| connectivity state | CAP-EPT-009/010 | online/offline context | oui | fresh | connectivity unknown |
| resource/capacity facts | Endpoint local state | pressure context | non | current | pressure unknown |
| policy/retention projection | Settings/Security | restriction context | non | current | local limits explicit/unknown |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Telemetry/Technical Output refs | Endpoint | payload reference/metadata | read within permission |
| Endpoint Policy | Settings | retention/queue projection | read projection |
| Agent connectivity | Endpoint | offline/degraded context | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Buffered Item | create/state | Endpoint | reference + provenance, no physical schema |
| Local Queue State | derive/refresh | Endpoint | bounded conceptual state |
| Drop/Retention Limitation | record | Endpoint | source-backed, never hidden |

## 11. Fonctionnalités
Queue local items; preserve source/time; expose retained/queued/dropped indication; track capacity pressure and conceptual age bounds; mark degraded/unknown; never claim delivered until handoff confirms it.

## 12. Actions utilisateur
Inspect Class 0; calculate pressure/freshness Class 1; bounded queue-status refresh Class 2. No arbitrary deletion or retention policy mutation.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| calculate queue state | oui | oui | oui | summarize | counters/timestamps |
| classify pressure/limits | oui | oui | oui | explain | threshold/rule source |
| invent retained/dropped data | non | interdit | non | interdit | explicit unknown |

## 14. États fonctionnels
`empty`, `buffered`, `queued`, `retained`, `capacity-pressure`, `degraded`, `drop-indicated`, `retention-expiring`, `offline`, `state-unknown`, `stale`.

## 15. États d’interface
No Screen ID. Buffered/retained never renders as delivered; unknown storage facts remain unknown.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Local Queue Projection | Endpoint fact | CAP-EPT-090/092 | conceptual state only |
| retention/drop limitation | technical quality fact | Investigate/Settings/Quality | no hidden loss |
| queue provenance | local audit context | CAP-EPT-097 | source/time retained |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| endpoint source | delivery unavailable | CAP-EPT-089 | item/source/time | source retained |
| CAP-EPT-089 | reconnect | CAP-EPT-090 | queued items/order/limits | queue state retained |
| CAP-EPT-089 | pressure | CAP-EPT-092 | capacity/drop/degraded facts | queue context retained |

## 18. Dépendances
CAP-EPT-009/010/024/025/090/092/097; Endpoint resilience sources; Shared Jobs/Retry; Settings/Security; OPEN-008.

## 19. Source de vérité
Endpoint SOT for local buffer/queue facts; Shared and Settings/Security retain generic/admin policy ownership.

## 20. Provenance et audit
Agent/source/item refs, timestamps, queue state, connectivity, pressure, retention/drop indication, policy ref and consumer handoff refs.

## 21. Permissions fonctionnelles
Buffer/queue state read, sensitive queued-metadata read, provenance read, cross-tenant deny; no final RBAC.

## 22. Limites et erreurs
Offline != data lost automatically; buffered != delivered; retained != durable forever; backpressure != storage guarantee; no physical queue defined.

## 23. Métriques
Queued/retained/drop-indicated items, age/pressure distributions, degraded duration and unknown/stale queue state.

## 24. Classification de livraison
`draft / defined / planned`; no storage/queue implementation.

## 25. Critères d’acceptation
**Given** connectivity is lost, **When** an item is retained locally, **Then** it is buffered/queued and not marked delivered.

**Given** capacity pressure causes a drop indication, **When** quality is reviewed, **Then** the indication is explicit rather than silently hidden.

**Given** AI is unavailable, **When** queue state is assessed, **Then** deterministic/manual inspection remains available.

## 26. Questions ouvertes
OPEN-008 remains open; platform storage mechanism and retention values are not selected.

## 27. Consommateurs documentaires
EPT-6 Resilience, Endpoint telemetry/live response, Shared, Settings, Security, Quality, registers and Roadmap.