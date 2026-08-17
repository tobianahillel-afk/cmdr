---
id: CAP-EPT-091
title: Local State Persistence, Restart and Crash Recovery
product: endpoint-agent
module: resilience
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-12
requirement_ids: [REQ-PROD-005, REQ-PROD-006, REQ-PROD-012, REQ-PROD-017, REQ-PROD-018, REQ-PROD-019, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004]
open_decisions: [OPEN-008]
source-of-truth: canonical
---
# CAP-EPT-091 — Local State Persistence, Restart and Crash Recovery

## 1. Définition
Définir les faits fonctionnels de persistance opérationnelle locale, restart et crash recovery de l’Agent sans définir watchdog, journal physique ou service OS.

## 2. Problème utilisateur
Après crash/restart, un Agent redémarré peut avoir perdu un état, conservé un command incomplet ou nécessiter resynchronisation sans être réellement healthy.

## 3. Objectifs
Represent persisted operational state concept, restart/crash, restart recovery, recovered/lost/partial state, dependency restart, resynchronization requirement and provenance.

## 4. Non-objectifs
Aucun watchdog, state journal schema, OS service, recovery engine, automatic replay, full-state guarantee or implementation.

## 5. Propriétaire
Endpoint owns its local restart/crash recovery facts; Shared retains generic Recovery and Settings retains administrative configuration.

## 6. Utilisateurs
Endpoint Operator, Platform Administrator, Support Operator, Security Reviewer, Auditor.

## 7. Conditions d’entrée
Agent identity/version, pre-crash known state if any, restart/crash observation, local persisted-state evidence, incomplete-operation refs and dependency state.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| pre-restart state refs | Endpoint local state | recovery anchor | non | last known | recovery partial/unknown |
| crash/restart event | Endpoint runtime | lifecycle fact | oui | event time | no recovery inference |
| persisted-state evidence | local source | state context | non | recovery time | state unknown |
| incomplete operation refs | CAP-EPT-062/086/090 | continuity context | non | exact attempts | no replay assumption |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Endpoint Agent | Endpoint | version/health/lifecycle | read |
| Local Queue/Replay State | Endpoint | retained work context | read |
| Technical/Update Attempts | Endpoint | incomplete operations | read |
| Endpoint Policy | Settings | effective restrictions | read projection |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Restart/Crash Recovery State | create/transition | Endpoint | source-backed |
| Local Persistent State Projection | recover/mark partial/lost | Endpoint | conceptual only |
| Resynchronization Requirement | derive | Endpoint | restart != synchronized |

## 11. Fonctionnalités
Record restart/crash; inspect retained state; preserve incomplete-operation facts; restore only supported local state; flag lost/partial; assess dependencies; require resynchronization; preserve audit continuity.

## 12. Actions utilisateur
Inspect Class 0; assess recovery/resync Class 1; request bounded restart/recovery Class 2/3 according effect/policy. No silent replay of incomplete effectful operations.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| reconcile pre/post state | oui | oui | oui | explain | state comparison |
| flag incomplete/lost state | oui | oui | oui | summarize | explicit refs/gaps |
| invent recovered state/replay | non | interdit | non | interdit | unknown/manual review |

## 14. États fonctionnels
`running`, `restart-observed`, `crashed`, `recovering`, `recovered-partial`, `recovered-local`, `state-lost`, `state-unknown`, `dependency-restarting`, `resynchronization-required`, `degraded`, `recovery-failed`.

## 15. États d’interface
No Screen ID. Restarted is not healthy; recovered local is not synchronized; unknown state is never coerced to success.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Recovery State | Endpoint fact | CAP-EPT-093/097/Settings | local scope explicit |
| Resynchronization Requirement | derived fact | CAP-EPT-090/Settings | remote state not assumed |
| incomplete-operation context | provenance refs | operators/Govern if relevant | no implicit replay |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| runtime | crash/restart | CAP-EPT-091 | Agent/pre-state/events | Agent context retained |
| CAP-EPT-091 | queue/resync required | CAP-EPT-090 | retained items/state gaps | recovery ref |
| CAP-EPT-091 | dependency issue | CAP-EPT-093 | affected component/capabilities | reassessment ref |

## 18. Dépendances
CAP-EPT-008/010/062/086/089/090/093/097; resilience crash-recovery; Shared Recovery; Settings; OPEN-008.

## 19. Source de vérité
Endpoint SOT for local restart/crash/persistence facts; Shared generic Recovery and Settings admin state remain external.

## 20. Provenance et audit
Agent/version, crash/restart times, prior state refs, recovered/lost state, incomplete attempts, dependency status, resync requirement and actor/requester where applicable.

## 21. Permissions fonctionnelles
Recovery-state read, bounded restart/recovery request, sensitive provenance read, cross-tenant deny; final RBAC not selected.

## 22. Limites et erreurs
Process restarted != service healthy; Agent restarted != Endpoint healthy; local state restored != remote synchronized; replay requires separate eligibility/idempotency.

## 23. Métriques
Crash/restart count, recovery partial/failure, state-loss/unknown, resync-required and incomplete-operation counts.

## 24. Classification de livraison
`draft / defined / planned`; no watchdog/journal/service/recovery implementation.

## 25. Critères d’acceptation
**Given** Agent restarts after crash with incomplete work, **When** recovery is assessed, **Then** incomplete operations remain explicit and are not automatically replayed.

**Given** local state is restored but remote sync is pending, **When** status is shown, **Then** resynchronization-required remains explicit.

**Given** AI is unavailable, **When** pre/post state is reconciled, **Then** deterministic/manual review remains available.

## 26. Questions ouvertes
OPEN-008 remains open; platform-specific persistence/watchdog/service behavior is not selected.

## 27. Consommateurs documentaires
EPT-6 Resilience/Security, Shared Recovery, Settings, Govern consumers where relevant, Quality, registers and Roadmap.