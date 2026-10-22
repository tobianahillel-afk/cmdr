---
id: CAP-EPT-090
title: Reconnect, Replay, Resumption and Duplicate-Control Boundary
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
# CAP-EPT-090 — Reconnect, Replay, Resumption and Duplicate-Control Boundary

## 1. Définition
Définir reconnect, replay, resumption et duplicate-candidate semantics pour les éléments locaux retenus, sans garantir exactly-once, consommation ni ordre causal.

## 2. Problème utilisateur
Après reconnexion, un élément rejoué peut être pris pour consommé, unique ou correctement ordonné alors que l’acknowledgement ou la déduplication reste incertaine.

## 3. Objectifs
Represent disconnected/reconnecting/reconnected, replay pending/replaying/replayed, duplicate candidate, ordering limitations, partial resumption, acknowledgement if sourced and provenance.

## 4. Non-objectifs
Aucun transport, event bus, delivery protocol, exactly-once guarantee, global dedup engine, physical queue or consumer implementation.

## 5. Propriétaire
Endpoint owns local reconnect/replay/resumption facts; consumers and Shared retain their acknowledgement/dedup generic semantics where applicable.

## 6. Utilisateurs
Endpoint Operator, SOC/Investigate consumer, Platform Administrator, Data Quality Reviewer, Auditor.

## 7. Conditions d’entrée
CAP-EPT-089 queue state, connectivity transition, stable item refs where available, ordering/idempotency metadata and consumer acknowledgement only when sourced.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| queued/buffered items | CAP-EPT-089 | replay candidates | oui | retained state | nothing to replay |
| connectivity | CAP-EPT-009/010 | reconnect context | oui | current | reconnect unknown |
| identity/order metadata | CAP-EPT-024 | quality context | non | source time | duplicate/order unknown |
| acknowledgement | consumer/source | receipt context | non | source-defined | consumed unknown |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Buffered Item | Endpoint | item/source identity | read |
| Local Queue State | Endpoint | retained order/limits | read |
| Telemetry quality state | Endpoint | duplicate/order context | read |
| consumer acknowledgement | consumer owner | acknowledgement ref only | read if available |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Replay State | create/transition | Endpoint | replayed != consumed |
| Duplicate Candidate | derive/link | Endpoint | candidate != confirmed duplicate |
| Resumption State | derive/refresh | Endpoint | partial explicit |

## 11. Fonctionnalités
Track connectivity recovery; start bounded replay; preserve item identity/order limitations; expose replay progress/partial state; detect duplicate candidates; attach acknowledgement only if sourced; avoid consumption/exactly-once inference.

## 12. Actions utilisateur
Inspect Class 0; assess replay/order/duplicate candidates Class 1; request bounded resumption/replay Class 2 under policy. No consumer-side mutation by mere acknowledgement read.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| correlate replay items | oui | oui | oui | summarize | stable refs |
| flag duplicate candidates | oui | source rules | oui | explain | deterministic indicators |
| claim exactly-once/consumed | non | interdit | non | interdit | explicit acknowledgement/unknown |

## 14. États fonctionnels
`disconnected`, `reconnecting`, `reconnected`, `replay-pending`, `replaying`, `replayed`, `partially-resumed`, `duplicate-candidate`, `ordering-limited`, `acknowledged`, `ack-unknown`, `replay-failed`, `state-unknown`.

## 15. États d’interface
No Screen ID. Replayed and acknowledged/consumed remain distinct; duplicates remain candidates unless consumer source proves otherwise.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Replay State | Endpoint fact | consumers/Quality | replay != consumption |
| Duplicate Candidate | quality fact | consumer/Shared | candidate only |
| Resumption Projection | Endpoint fact | Settings/operations | partial/limits explicit |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-EPT-089 | reconnect | CAP-EPT-090 | retained items/order/limits | queue refs retained |
| CAP-EPT-090 | replay output | consumer | item/source/provenance | acknowledgement if any |
| consumer ack | receipt | CAP-EPT-090 | ack ref/time | no exactly-once inference |

## 18. Dépendances
CAP-EPT-009/010/024/025/089/091/097; Shared retry/data-quality mechanisms; OPEN-008.

## 19. Source de vérité
Endpoint SOT for local replay/resumption facts; consumer owns its receipt/consumption state; Shared remains generic mechanism owner.

## 20. Provenance et audit
Agent/item refs, disconnect/reconnect times, replay attempt, ordering/duplicate indicators, acknowledgement ref, partial/failure reasons and timestamps.

## 21. Permissions fonctionnelles
Replay-state read, bounded replay request, sensitive item metadata read, provenance read, cross-tenant deny; no final RBAC.

## 22. Limites et erreurs
Replayed != consumed; replay != exactly once; duplicate candidate != confirmed duplicate; reconnect != recovered; acknowledgement absent remains unknown.

## 23. Métriques
Replay backlog/outcomes, partial resumptions, duplicate candidates, acknowledgement unknown, ordering-limited and replay failure counts.

## 24. Classification de livraison
`draft / defined / planned`; no transport/dedup/replay engine implementation.

## 25. Critères d’acceptation
**Given** buffered items replay after reconnect, **When** no consumer acknowledgement exists, **Then** they are replayed but consumption remains unknown.

**Given** two replay items share duplicate indicators, **When** reviewed, **Then** they are duplicate candidates rather than confirmed duplicates.

**Given** AI is unavailable, **When** replay state is reconciled, **Then** deterministic/manual correlation remains possible.

## 26. Questions ouvertes
OPEN-008 remains open; platform/transport-specific ordering and acknowledgement rules are not selected.

## 27. Consommateurs documentaires
EPT-6 Resilience, Shared/Data Quality, Investigate consumers, Settings, Security, Quality, registers and Roadmap.