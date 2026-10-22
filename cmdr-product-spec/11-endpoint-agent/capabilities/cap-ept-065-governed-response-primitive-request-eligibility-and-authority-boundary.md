---
id: CAP-EPT-065
title: Governed Response Primitive Request, Eligibility and Authority Boundary
product: endpoint-agent
module: containment
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-004, REQ-PROD-006, REQ-PROD-014, REQ-PROD-015, REQ-PROD-016, REQ-PROD-018, REQ-PROD-019, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-008, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-EPT-065 — Governed Response Primitive Request, Eligibility and Authority Boundary

## 1. Définition
Définir comment Endpoint reçoit une intention d’exécution issue d’un Govern Response Run/Response Step, la relie à une `Response Primitive` locale, évalue l’éligibilité et la disponibilité techniques et exige une authority/Decision/Run reference appropriée sans créer l’autorité lui-même.

## 2. Problème utilisateur
Une primitive peut exister techniquement tout en étant indisponible, hors scope ou non autorisée. Sans frontière explicite, `eligible`, `available`, `allowed`, `authorized` et `executing` peuvent être confondus.

## 3. Objectifs
Identifier primitive/type/target/origin/requester/purpose/requested state ; calculer eligibility/availability ; vérifier platform/capability/policy projections ; consommer authority/Decision/Response Run refs ; produire des raisons de blocage ; préserver return origin/provenance.

## 4. Non-objectifs
Créer Action Request, Approval, Decision, Response Run ou Result ; décider l’autorité ; exécuter la primitive ; définir API/protocole/commande ; fermer OPEN-013/015 ; commencer EPT-6.

## 5. Propriétaire
Endpoint possède `Response Primitive`, `Primitive Eligibility` et les faits techniques locaux. Govern possède Action Request, Approval, Decision, Response Run, production authority et Result. Settings possède Policy/Fleet/admin configuration.

## 6. Utilisateurs
Response Operator, Endpoint Operator, Govern Reviewer, Verification Reviewer, Security Reviewer, Auditor.

## 7. Conditions d’entrée
Response Run/Step ou autre origine gouvernée résoluble, target tenant-bound, primitive demandée identifiable, capability/platform facts disponibles ou explicitement unknown, authority context fourni selon classe.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Run/Step + Decision refs | Govern | authority lineage | oui pour effet gouverné | current/unexpired | `authority-missing` |
| primitive type + requested state | Govern handoff | technical intent | oui | exact handoff | reject |
| target/Agent identity | EPT-1 | target context | oui | current enough | target-unresolved |
| capability/platform state | EPT-1/EPT-2 | technical eligibility | oui | observed timestamp | eligibility-unknown |
| Endpoint Policy projection | Settings | local restriction | selon primitive | current projection | policy-unknown/block |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Decision / Response Run / Step | Govern | ids, bounds, conditions, expiry | read/reference |
| Endpoint Agent / Technical Capability | Endpoint | target/platform/availability | local read |
| Endpoint Policy | Settings | restriction/config projection | read only |
| Human Gate / Automation Run / Tool Call | Studio | provenance refs only | read/link |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Response Primitive Reference | create/resolve | Endpoint | technical primitive only |
| Primitive Eligibility | derive/version | Endpoint | eligible ≠ authorized |
| Authority Reference Projection | attach/refresh | Endpoint | no authority creation |
| target | aucune mutation | target owner | this capability performs no effect |

## 11. Fonctionnalités
Resolve primitive and target ; classify capability/platform support ; calculate eligibility/availability ; compare requested scope to technical limits ; consume authority/Decision/Run conditions ; expose `eligible`, `ineligible`, `unavailable`, `unknown`, `authority-missing`, `blocked`, `ready-for-precheck`; never infer authorization from capability presence.

## 12. Actions utilisateur
Class 0 inspect primitive/eligibility/provenance. Class 1 recompute eligibility/availability. No class 2/3 effect is performed here. Bypass of missing authority is prohibited.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| resolve primitive/capability | oui | oui | oui | explain only | catalog + declared capabilities |
| compare target/scope | oui | oui | oui | summary | deterministic diff |
| explain blockers | oui | reason codes | oui | oui, attributed | structured reasons |
| grant authority/start effect | non | interdit | non | interdit | Govern + explicit operator path |

## 14. États fonctionnels
`received`, `target-unresolved`, `primitive-unresolved`, `eligible`, `ineligible`, `available`, `unavailable`, `eligibility-unknown`, `authority-missing`, `authority-stale`, `policy-blocked`, `ready-for-precheck`, `rejected`, `superseded`.

## 15. États d’interface
Aucun Screen ID. Loading/Partial/Error/Offline/Permission denied/Stale préservent la distinction capability, permission, authority et dernière observation valide.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Primitive Eligibility | Endpoint assessment | CAP-EPT-066/Govern | reasons/source/version |
| Response Primitive Ref | technical reference | EPT-5 primitive families | exact type/target |
| Authority blocker/ref | boundary projection | Govern | no local authority inference |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Govern Response Run/Step | effect handoff | CAP-EPT-065 | Run/Decision/target/action/conditions | same Run |
| CAP-EPT-065 | eligible + authority context | CAP-EPT-066 | primitive/target/requested state/limits | Run retained |
| blocked/ineligible | assessment | Govern | reasons/limitations | no effect |

## 18. Dépendances
CAP-EPT-001..064, CAP-GOV-015/016/022/025, Settings Endpoint Policy/Fleet, Security authority/step-up, OPEN-007/008/013/015.

## 19. Source de vérité
Endpoint est SOT de l’éligibilité/availability technique locale. Govern est SOT de l’autorité/Decision/Response Run. Settings est SOT de Policy/admin configuration.

## 20. Provenance et audit
Primitive/type/version, target, requester/origin, purpose, requested state, Decision/Run/Step refs, capability/platform observations, Policy ref/version, authority ref, eligibility reasons, actor/service, timestamps, correlation id.

## 21. Permissions fonctionnelles
Primitive read, eligibility read/recompute, capability/policy read, authority-reference read, sensitive response context read, cross-tenant deny. Aucun final RBAC/ABAC.

## 22. Limites et erreurs
`Response Primitive ≠ Action Request`; eligible ≠ authorized ; available ≠ allowed ; authority ref present ≠ authority valid automatically ; Endpoint Policy ≠ authority ; capability availability ≠ permission ; operator ≠ approver.

## 23. Métriques
Eligibility outcomes, unsupported/unavailable, authority-missing/stale, policy blocks, cross-tenant denials, accepted-to-precheck without authority target zero. Aucun SLO numérique.

## 24. Classification de livraison
`draft / defined / planned`. Contrat documentaire ; aucune primitive, commande, API, protocole ou plateforme livrée n’est prouvée.

## 25. Critères d’acceptation
**Given** une primitive est techniquement disponible sans Decision/authority requise, **When** eligibility est calculée, **Then** elle reste `authority-missing` et aucun effet ne démarre.

**Given** une primitive n’est pas déclarée sur la plateforme observée, **When** la request arrive, **Then** elle est ineligible/unavailable avec raison et aucun fallback effectful n’est tenté.

**Given** l’IA est indisponible, **When** l’éligibilité est évaluée, **Then** les règles déterministes et projections canoniques fournissent le même blocage.

## 26. Questions ouvertes
OPEN-007/008/013/015 restent ouvertes. Ce contrat ne décide ni support release, class-2 default, Human Gate bridge ni run identity merger.

## 27. Consommateurs documentaires
Endpoint EPT-5, Govern Runs & Rollback, Investigate containment handoff, Settings, Studio, Security, Shared, Quality et Roadmap Phase 5.