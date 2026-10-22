---
id: CAP-EPT-066
title: Response Primitive Preconditions, Readiness and Target-State Planning
product: endpoint-agent
module: containment
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-004, REQ-PROD-005, REQ-PROD-006, REQ-PROD-015, REQ-PROD-016, REQ-PROD-018, REQ-PROD-019, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-EPT-066 — Response Primitive Preconditions, Readiness and Target-State Planning

## 1. Définition
Évaluer les préconditions techniques d’une primitive EPT-5, observer l’état cible courant, expliciter l’état cible demandé, les dépendances, impact/reversibility information et la stratégie de vérification technique avant toute exécution.

## 2. Problème utilisateur
Une action autorisée peut encore être dangereuse ou impossible si la cible a changé, est stale/offline, si une dépendance manque ou si l’état demandé n’est plus pertinent.

## 3. Objectifs
Capturer current/requested target state ; vérifier freshness/availability/capability/policy ; identifier prerequisites/dependencies ; exposer expected technical impact ; documenter reversibility ; référencer verification strategy ; produire ready/blocked/unknown sans prétendre garantir la sécurité universelle.

## 4. Non-objectifs
Autoriser l’action, modifier la cible, décider le rollback Govern, créer Verification Plan Govern, définir commandes/engine, garantir que precheck PASS signifie safe dans toutes circonstances.

## 5. Propriétaire
Endpoint possède `Technical Precheck`, readiness et target-state planning technique. Govern possède Execution Plan, Decision conditions, response verification/rollback. Settings possède Policy/admin config.

## 6. Utilisateurs
Response Operator, Endpoint Operator, Verification Reviewer, Govern Reviewer, Security Reviewer, Auditor.

## 7. Conditions d’entrée
Primitive/eligibility CAP-EPT-065, exact target, requested state, authority context, latest technical observations and declared dependencies.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| primitive/eligibility | CAP-EPT-065 | subject | oui | current version | no precheck |
| current target facts | EPT-1..4 | observed state | oui | freshness visible | state-unknown |
| requested target state | Govern handoff | intended technical state | oui | exact Run/Step | blocked |
| dependencies/capabilities | Endpoint/Settings | readiness inputs | selon primitive | current projection | dependency-missing |
| verification requirement ref | Govern Decision/Plan | downstream constraint | selon action | pinned | explicit unknown/block |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Response Primitive/Eligibility | Endpoint | type/target/limits | read |
| Endpoint Agent/observations | Endpoint | current health/state | read |
| Decision/Execution Plan/Run | Govern | requested state/conditions | read/reference |
| Endpoint Policy/Fleet | Settings | configuration/restrictions | read only |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Technical Precheck | create/rerun/version | Endpoint | no target effect |
| Requested Target State projection | attach/version | Endpoint | mirrors approved intent, no authority |
| Readiness Assessment | derive | Endpoint | PASS ≠ universal safety guarantee |
| target | aucune mutation | target owner | precheck only |

## 11. Fonctionnalités
Observe target state/freshness ; compare current vs requested ; verify platform/capability/dependency/policy/authority freshness ; identify expected impact and reversibility metadata ; determine verification-observation path ; produce `ready`, `blocked`, `stale`, `unknown`, `already-in-desired-state`, `dependency-missing`, `scope-mismatch`.

## 12. Actions utilisateur
Class 0 inspect precheck/state. Class 1 run/re-run deterministic precheck/readiness. Updating requested effect or authority is outside Endpoint and returns to Govern.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| compare current/requested state | oui | oui | oui | explain | state diff |
| evaluate dependencies/freshness | oui | oui | oui | summary | readiness checklist |
| summarize impact/limits | oui | sourced metadata | oui | oui | source matrix |
| declare safe/authorize/start | non | interdit | non | interdit | Govern + explicit execution path |

## 14. États fonctionnels
`not-checked`, `checking`, `ready`, `blocked`, `target-stale`, `target-unknown`, `target-unavailable`, `scope-mismatch`, `dependency-missing`, `policy-blocked`, `authority-stale`, `already-in-desired-state`, `verification-path-missing`, `superseded`.

## 15. États d’interface
Aucun Screen ID. Stale/Offline/Partial/Permission denied distinguent absence de données, absence de cible et absence d’autorité ; le dernier état observé reste timestampé.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Technical Precheck | Endpoint assessment | EPT-5 primitive execution | exact sources/reasons/time |
| Readiness Assessment | local readiness | Govern/Response Operator | ready ≠ authorized/safe universally |
| state/impact/limit summary | technical context | verification/reversal | observed vs requested explicit |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-EPT-065 | eligible handoff | CAP-EPT-066 | primitive/target/requested state/authority refs | origin Run |
| CAP-EPT-066 | ready | CAP-EPT-067..072/081 | exact precheck + state + limits | Run retained |
| stale/blocked | precheck outcome | Govern | reason/current observation/required re-review | no effect |

## 18. Dépendances
CAP-EPT-001..065, CAP-GOV-019/020/021/022/025/028, Settings Policy/Fleet, Security least privilege, OPEN-008/013/015.

## 19. Source de vérité
Endpoint est SOT des observations/prechecks locaux. Govern reste SOT de requested response objective/Execution Plan/authority. Settings reste SOT des policies/admin settings.

## 20. Provenance et audit
Primitive/request refs, target, current/requested states, source timestamps, capability/platform, dependencies, policy/authority refs, impact/reversibility metadata, verification-strategy ref, actor/service, reasons, correlation id.

## 21. Permissions fonctionnelles
Precheck execute/read, target-state read, capability/policy projection read, sensitive response-context read, cross-tenant deny. Effet et authority nécessitent permissions distinctes.

## 22. Limites et erreurs
requested state ≠ achieved state ; precheck PASS ≠ action safe in all circumstances ; stale target blocks fresh readiness claim ; current state unknown ≠ failed ; Policy ≠ authority ; readiness ≠ permission.

## 23. Métriques
Ready/blocked/unknown, stale-target rate, dependency missing, already-desired, authority/policy blocks, execution attempted after failed/stale precheck target zero.

## 24. Classification de livraison
`draft / defined / planned`; aucun precheck engine, action runtime, API, protocol ou implementation n’est livré.

## 25. Critères d’acceptation
**Given** la target state est stale avant exécution, **When** precheck s’exécute, **Then** readiness n’est pas déclarée fraîche et la primitive ne démarre pas.

**Given** le current state correspond déjà au requested state, **When** readiness est évaluée, **Then** `already-in-desired-state` est exposé sans prétendre que l’objectif de réponse métier est satisfait.

**Given** l’IA est absente, **When** readiness est évaluée, **Then** state diff, freshness et dependency checks déterministes fournissent le résultat.

## 26. Questions ouvertes
OPEN-008/013/015 restent ouvertes. Aucun seuil universel de freshness/safety, class-2 default ou run bridge n’est choisi.

## 27. Consommateurs documentaires
EPT-5 primitive families, Govern Runs/Verification/Rollback, Settings, Security, Quality et future EPT-6 uniquement comme frontière.