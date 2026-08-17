---
id: CAP-EPT-072
title: Service and System-Component Control Primitive
product: endpoint-agent
module: containment
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-004, REQ-PROD-005, REQ-PROD-006, REQ-PROD-016, REQ-PROD-018, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-EPT-072 — Service and System-Component Control Primitive

## 1. Définition
Définir des primitives de contrôle d’un service ou composant système lorsque le corpus/platform capability les source, avec current/requested state, dependency impact, authority, execution, verification and reversal facts.

## 2. Problème utilisateur
Stop/start/restart/disable-like actions peuvent affecter des dépendances et varier selon plateforme. Un service stopped ne signifie pas remediation complete.

## 3. Objectifs
Pinner service/component identity ; expose current/requested state ; validate dependency/platform support ; support only declared stop/start/restart/disable-like semantics ; track effect/partial/failure ; verify state ; prepare reversal.

## 4. Non-objectifs
Définir noms de services, native commands, universal service manager, platform support, persistent anti-tamper/watchdog, EPT-6 resilience, or remediation verdict.

## 5. Propriétaire
Endpoint owns local service/component technical control facts. Govern owns authority/Run/verification/result; Settings owns Policy/admin config.

## 6. Utilisateurs
Response Operator, Endpoint Operator, System Reviewer, Govern Reviewer, Verification Reviewer, Auditor.

## 7. Conditions d’entrée
CAP-EPT-065/066 ready, stable component reference, declared operation/capability, dependency impact visible enough, authority current.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| service/component ref | EPT observations | target | oui | current | unresolved |
| requested operation/state | Govern | effect intent | oui | pinned | reject |
| current state/dependencies | Endpoint | precheck | oui | timestamped | unknown/block |
| platform capability/policy | Endpoint/Settings | support/restriction | oui | current | unsupported/block |
| authority ref | Govern | authorization | oui | effective | no action |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Response Run/Decision | Govern | operation/scope | read |
| service/system observations | Endpoint | current state/dependencies | read |
| Endpoint Policy | Settings | restrictions | read |
| Technical Capability | Endpoint | supported operation | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Service Control Execution | create/transition | Endpoint | exact component/op |
| Service Control State | observe/update | Endpoint | requested vs observed |
| Technical Outcome | create | Endpoint | no remediation claim |
| service/component | bounded state mutation | target | implementation unspecified |

## 11. Fonctionnalités
Resolve component; assess dependency impact; validate platform capability; request declared control action; track start/status/partial/failure/unknown; observe resulting state; expose reversal capability/limitations.

## 12. Actions utilisateur
Inspect/readiness Class 0/1. Stop/start/restart/disable-like effect Class 3 by default and Govern-dependent.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| resolve current/dependencies | oui | oui | oui | explain | dependency table |
| verify observed state | oui | oui | oui | summary | source observation |
| suggest operation | oui | declared capability | oui | suggestion | operator selection |
| authorize/execute autonomously | non | no | non | interdit | Govern/operator |

## 14. États fonctionnels
`unsupported`, `state-unknown`, `ready`, `control-requested`, `executing`, `stopped-observed`, `running-observed`, `restart-observed`, `disabled-observed`, `partial`, `failed`, `unknown`, `dependency-impact-unknown`, `reversal-eligible`, `drifted`.

## 15. États d’interface
No Screen ID. Unsupported, unknown, stopped, disabled and dependency-impact unknown remain distinct.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Service Control State | technical fact | CAP-EPT-073/075 | component/state/time |
| technical outcome | execution fact | Govern | no remediation verdict |
| reversal/impact context | technical refs | CAP-EPT-077 | limitations explicit |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-EPT-066 | service control ready | CAP-EPT-072 | component/op/dependencies/authority | Run |
| CAP-EPT-072 | outcome | CAP-EPT-073/076 | observed state/limits | same Run |
| CAP-EPT-072 | verify/reverse | CAP-EPT-074/075/077 | expected/current state | Run retained |

## 18. Dépendances
CAP-EPT-021/041/050/058/062/065/066/073..080, Govern, Settings, Security, OPEN-008/013/015.

## 19. Source de vérité
Endpoint SOT of target-side service/component execution/state. Govern retains authority/outcome; Settings retains policy/admin config.

## 20. Provenance et audit
Component identity, operation/current/requested states, dependencies/impact, platform capability, authority/Run/Step, operator, Agent/version, timestamps, outcome/errors/verification/reversal refs.

## 21. Permissions fonctionnelles
Service/system control request/read, verification, sensitive system context, cross-tenant deny, step-up/SoD/Govern dependency; no final RBAC.

## 22. Limites et erreurs
Service stopped ≠ remediation complete; unsupported platform ≠ failed action; restart ≠ rollback; dependency impact unknown blocks universal safety claim.

## 23. Métriques
Unsupported, control outcomes, dependency blocks, partial/failure/unknown, verification mismatch, unauthorized attempts target zero.

## 24. Classification de livraison
`draft / defined / planned`; no native OS command/service manager, API/protocol or EPT-6 implementation.

## 25. Critères d’acceptation
**Given** service mutation is unsupported on the observed platform, **When** eligibility/precheck runs, **Then** no effect occurs and unsupported is not misreported as execution failure.

**Given** service stop is observed, **When** verification completes, **Then** target technical state is recorded but remediation success is not claimed.

**Given** dependencies change before execution, **When** precheck refreshes, **Then** readiness is invalidated until scope/impact are re-reviewed.

## 26. Questions ouvertes
OPEN-008/013/015 remain open; no platform support/runtime semantics are finalized.

## 27. Consommateurs documentaires
EPT-5 verification/reversal/provenance, Govern, Settings, Security, Quality.