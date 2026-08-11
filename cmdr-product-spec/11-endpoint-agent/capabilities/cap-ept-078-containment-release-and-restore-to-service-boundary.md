---
id: CAP-EPT-078
title: Containment Release and Restore-to-Service Boundary
product: endpoint-agent
module: containment
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-004, REQ-PROD-005, REQ-PROD-006, REQ-PROD-016, REQ-PROD-018, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004]
open_decisions: [OPEN-008, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-EPT-078 — Containment Release and Restore-to-Service Boundary

## 1. Définition
Définir la primitive technique de release d’un état de containment et de restauration technique de service/connectivity lorsqu’elle est explicitement autorisée, sans déclarer que la cible est sûre pour retour en production.

## 2. Problème utilisateur
Relâcher isolation/quarantine/network control peut restaurer connectivity ou access sans restaurer l’état pré-incident ni réduire le risque métier.

## 3. Objectifs
Bind release request to existing containment state ; authority ; eligibility/precheck ; execute technical release ; observe connectivity/service/health state ; represent partial/failed/unknown restoration ; verify and return to Govern.

## 4. Non-objectifs
Décider return-to-production, declare secure endpoint, create Govern recovery/Result, re-enable unsupported components, perform EPT-6 update/resilience, or define commands/protocol.

## 5. Propriétaire
Endpoint owns containment-release technical primitive/state. Govern owns release authority, response recovery/verification/Result. Settings owns Policy/admin restrictions.

## 6. Utilisateurs
Response Operator, Endpoint Operator, Verification Reviewer, Govern Reviewer, Service Owner as consumer, Auditor.

## 7. Conditions d’entrée
Existing containment state, explicit Govern release intent/authority, target current enough, release capability/preconditions known, verification path available or limitations explicit.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| current containment state | CAP-EPT-068..072/081 | release subject | oui | current observation | unknown/block |
| release request | Govern | effect intent | oui | pinned | no action |
| authority/conditions | Govern | authorization | oui | current | blocked |
| release capability/readiness | Endpoint | technical | oui | current | unsupported/unknown |
| expected post-release checks | Govern/technical verification | verification | oui when required | pinned | limitation/block |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Isolation/Network/Quarantine/Service/Session state | Endpoint | current containment | read |
| Decision/Response Run/Rollback-Recovery context | Govern | release authority/objective | read/ref |
| Endpoint Policy | Settings | release restrictions | read |
| Agent health/connectivity observations | Endpoint | post-release facts | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Containment Release Execution | create/transition | Endpoint | source containment linked |
| Restore-to-Service Technical State | observe/update | Endpoint | technical only |
| Technical Outcome | create | Endpoint | no return-to-production verdict |
| target containment/control | release effect | target | bounded to authorized scope |

## 11. Fonctionnalités
Revalidate current containment; ensure authority; execute supported release/unblock/unquarantine/re-enable-like inverse only within declared primitive; observe connectivity/service/health; track partial/failure/unknown; verify state; expose residual containment/drift.

## 12. Actions utilisateur
Inspect/precheck Class 0/1. Containment release/restore effect Class 3 by default and governed. No autonomous production-return decision.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| validate containment/release state | oui | oui | oui | explain | state checklist |
| observe post-release health/connectivity | oui | oui | oui | summary | local observations |
| propose release candidate | oui | policy/capability | oui | suggestion | operator selection |
| authorize/declare safe-to-return | non | Govern-owned | non | interdit | Govern/service review |

## 14. États fonctionnels
`not-releasable`, `release-eligible`, `authority-required`, `release-requested`, `releasing`, `released-observed`, `partial-restoration`, `release-failed`, `release-unknown`, `connectivity-restored-observed`, `health-observed`, `residual-containment`, `verification-pending`, `drifted`.

## 15. États d’interface
No Screen ID. Release technical completion and business/security readiness remain separate. Restored connectivity is not displayed as secure by contract.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Containment Release State | technical fact | CAP-EPT-075/079 | source/time/scope |
| restore-to-service outcome | technical outcome | Govern | partial/limits explicit |
| post-release health/connectivity | verification refs | Govern/Settings consumer | no secure verdict |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Govern release/recovery | authorized handoff | CAP-EPT-078 | containment state/authority/desired release | Response Run/Rollback |
| CAP-EPT-078 | technical outcome | CAP-EPT-073/075/076 | release/observed state | same Govern context |
| CAP-EPT-078 | reconciliation | CAP-EPT-079 | outcome/health/limits | Govern |

## 18. Dépendances
CAP-EPT-068..077/079/080/081, CAP-GOV-029..032, Settings Policy/health, Security, OPEN-008/013/015.

## 19. Source de vérité
Endpoint SOT of technical release/restore state. Govern SOT of recovery/verification/Result and authority. Settings SOT of policy/admin config.

## 20. Provenance et audit
Original containment, target, release intent/authority, pre/post states, connectivity/health observations, operator, Agent/version, timestamps, partial/errors, verification/Run/Rollback refs and correlation.

## 21. Permissions fonctionnelles
Containment state read, release request, technical restore, verification read, sensitive context, cross-tenant deny, Govern/step-up/SoD dependency.

## 22. Limites et erreurs
Containment release ≠ restored pre-incident state; connectivity restored ≠ secure; release complete technically ≠ safe-to-return-to-production; release ≠ Govern rollback automatically.

## 23. Métriques
Release attempts/outcomes, partial/unknown, residual containment, health/connectivity mismatch, premature safe-return claims target zero.

## 24. Classification de livraison
`draft / defined / planned`; no release engine, command, API/protocol, final recovery implementation or EPT-6 work.

## 25. Critères d’acceptation
**Given** containment release succeeds technically, **When** connectivity is observed, **Then** restored connectivity is recorded without declaring the endpoint secure.

**Given** only part of the containment can be released, **When** execution completes, **Then** residual-containment and partial-restoration remain explicit.

**Given** release authority is missing or stale, **When** requested, **Then** no effect occurs even if the technical release primitive is available.

## 26. Questions ouvertes
OPEN-008/013/015 remain open; safe-return policy/thresholds remain Govern/operational decisions.

## 27. Consommateurs documentaires
Govern Rollback/Verification/Result, Endpoint health, Settings, Security, Quality, Audit.