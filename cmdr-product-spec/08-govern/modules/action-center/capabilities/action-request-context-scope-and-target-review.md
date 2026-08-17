---
id: CAP-GOV-004
title: Action Request Context, Scope and Target Review
product: govern
module: action-center
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-004, REQ-PROD-008, REQ-PROD-015, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-013]
source-of-truth: canonical
---
# CAP-GOV-004 — Action Request Context, Scope and Target Review

## 1. Définition
Examiner l’action demandée, son contexte source, ses targets, son scope inclus/exclu, tenant/environnement, time bounds, dépendances et ambiguïtés afin d’établir un contexte de review exploitable sans modifier la cible ni conclure sur l’autorisation.

## 2. Problème utilisateur
Une cible sélectionnée peut être stale, ambiguë, hors tenant, trop large ou différente du contexte du Finding/Incident. Autoriser sans review de scope expose à une action sur le mauvais objet ou un blast radius non compris.

## 3. Objectifs
Résoudre les références de target ; distinguer target selected/verified ; expliciter included/excluded scope ; préserver tenant/env/account/device/resource/endpoint/identity/network context ; exposer freshness/ambiguity/conflicts ; produire un Context Review sourcé.

## 4. Non-objectifs
Ne pas scanner une cible, prouver sa compromission, modifier une permission, sélectionner un runtime, calculer une Decision, exécuter une action ou définir un schéma universel de target.

## 5. Propriétaire
Govern / Action Center possède la review contextuelle de la request. Les targets restent chez leurs owners (Command/Investigate/Settings/Endpoint/cloud/source system) et Govern ne détient qu’une Target Reference permission-aware.

## 6. Utilisateurs
Principal : Govern Reviewer. Secondaires : Risk Reviewer, Authority Reviewer, Decision Maker, source requester, target owner consulté, Auditor.

## 7. Conditions d’entrée
Action Request version courante, source/return origin, target candidate, tenant, action type et permission de lire les projections nécessaires.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| proposed action | Action Request | requested effect | oui | request version | `incomplete` |
| source context | Incident/Case/Finding/Evidence/package | rationale and upstream facts | selon request | source version visible | partial/information-required |
| target references | request + source owners | target identity candidates | oui | resolution timestamp visible | ambiguous/blocked |
| tenant/environment/account/device/resource context | Settings/source owners | scope boundary | tenant oui, others as applicable | current projection | unknown/blocked |
| included/excluded scope and time bounds | requester/policy | bounded authorization candidate | oui for effectful action | request version | information-required |
| dependency/freshness/conflict data | source owners | blast-radius context | selon target | source timestamp | unknown, never inferred |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Action Request | Govern | action/target/scope/version | review/update context |
| Incident | Command | impact/service/urgency | read/link |
| Case/Finding/Evidence | Investigate | source facts/restrictions | read/link; no qualification |
| Tenant/Environment/Principal/Integration | Settings | administrative scope | minimal read |
| Endpoint/Cloud/network/identity target | respective owner | target identity/state/freshness | read-only projection |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Govern Review Context | add context/scope/target review | Govern | versioned against request |
| Target Reference | verify/mark ambiguous/stale | source owner + Govern projection | no target mutation |
| Action Request context | request/update clarification | Govern | class 2, may create new version |
| source targets | no mutation | source owners | read-only |

## 11. Fonctionnalités
Resolve target references; show target type and source; compare request/source target; define included/excluded scope; inspect tenant/env/account/resource/device/identity/network context; expose time bounds, blast-radius candidates, dependencies, freshness and contradictions; mark verified/ambiguous/restricted/stale; request clarification.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect/compare target | reviewer | Target Reference | 0 | source read | sourced comparison | non |
| run bounded scope consistency check | reviewer | Review Context | 1 | request/source refs | discrepancies listed | non |
| update included/excluded scope | authorized reviewer/requester | request context | 2 | version current | new scoped version | OPEN-013 |
| request target clarification | reviewer | Information Request | 2 | ambiguity named | source follow-up | OPEN-013 |
| mutate target | none in GOV-1 | real target | 3/4 | outside GOV-1 | prohibited | future GOV-2/runtime |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| resolve refs and scope keys | oui | oui | oui | no need | deterministic linking |
| detect context mismatch | oui | explicit comparisons | oui | explanation | diff/table |
| summarize blast-radius candidates | oui | graph/relationship aggregation | oui | sourced summary | dependency list/graph |
| suggest missing exclusions/time bounds | oui | checklist | oui | suggestion | form validation |
| authorize target | human later | no | no | prohibited | authority/Decision review |

## 14. États fonctionnels
`context-unreviewed`, `context-partial`, `target-unresolved`, `target-ambiguous`, `target-verified`, `scope-incomplete`, `scope-bounded`, `scope-conflict`, `stale`, `restricted`, `reviewed`, `blocked`.

## 15. États d’interface
Loading keeps target selection ; Empty means no target provided ; Partial lists missing source dimensions ; Error retains resolved refs ; Offline prohibits scope mutation ; Permission denied masks sensitive refs ; Stale displays last-resolved timestamp and prevents silent verification.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Context/Scope Review | Review Context result | CAP-GOV-005/006/007/009/014 | sources, gaps, version and uncertainty preserved |
| verified Target Reference | projection state | risk/authority/Decision review | verified identity ≠ safe/authorized target |
| clarification request | Information Request | source product | exact ambiguity and return origin |
| scope change diff | request version diff | reviewers/Approvers | previous scope retained |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Action Request | context review opened | CAP-GOV-004 | request/version/source refs/targets | Action Center |
| CAP-GOV-004 | scope established | CAP-GOV-005 | exact targets/scope/time/dependencies | context review |
| CAP-GOV-004 | clarification needed | source product | unresolved refs/questions | same request/version lineage |
| CAP-GOV-004 | review complete | CAP-GOV-006/007/009/014 | context result + gaps | Action Center |

## 18. Dépendances
CAP-GOV-001/003/005/006/007/009/014, Shared Linking/Search/Trace, Settings tenant/env/identity projections, source target owners, OPEN-013.

## 19. Source de vérité
Real target state remains source-owned. Govern is source of the review disposition and exact target/scope references used for the request version. Target selected does not become target verified automatically.

## 20. Provenance et audit
Record source object/version, target refs/resolution sources/timestamps, included/excluded scope, tenant/env, time bounds, dependencies, ambiguities, restrictions, reviewer, automation checks/suggestions, before/after scope diffs and return origin.

## 21. Permissions fonctionnelles
Action Request read/review/context update, restricted context read, target projection read, cross-tenant review if explicit, automated recommendation request. Target mutation permissions are absent from GOV-1.

## 22. Limites et erreurs
Unresolvable/multiple targets, stale state, permission denial, cross-tenant mismatch, changed environment, conflicting source IDs, target deleted/unavailable, partial inventory or hidden dependency produce explicit unknown/blocked states, never a fabricated verified target.

## 23. Métriques
Requests with unresolved/ambiguous/stale targets; scope changes after submission; cross-scope blocks; time to resolve target clarification; zero target mutations performed by GOV-1.

## 24. Classification de livraison
`defined` / `planned`; provider/runtime-neutral functional review only.

## 25. Critères d’acceptation
**Given** a request whose selected endpoint resolves to a different environment than the Incident, **When** context review runs, **Then** the conflict is visible, target is not marked verified and Decision readiness is blocked or requires clarification.

**Given** a verified target, **When** review completes, **Then** verification proves only identity/scope consistency, not safety, Approval or authorization.

**Given** no AI, **When** scope is reviewed, **Then** links, deterministic comparisons, forms and human inspection provide the full function.

## 26. Questions ouvertes
OPEN-013 remains open for class-2 context mutations. Final target-reference schema and cross-tenant authorization matrix are future Objects/Permissions work.

## 27. Consommateurs documentaires
Action Center, Risk/Policy/Authority/Decision capabilities, source-product handoffs, future Objects/Permissions/Screens/Journeys/GOV-2 and conformance gates.
