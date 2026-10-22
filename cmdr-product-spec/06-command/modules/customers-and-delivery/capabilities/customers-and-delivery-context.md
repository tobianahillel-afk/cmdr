---
id: CAP-CMD-401
title: Customers and Delivery Context
product: command
module: customers-and-delivery
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-14
requirement_ids: [REQ-PROD-053, REQ-PROD-013, REQ-PROD-019, REQ-PROD-033]
open_decisions: [OPEN-013, OPEN-019]
source-of-truth: canonical
---

# CAP-CMD-401 — Customers and Delivery Context

## 1. Définition
Capability deployment-dependent qui présente depuis Command un contexte de delivery interne, Enterprise ou MSP/MSSP, une projection Customer/engagement externe, un overview multi-tenant read-only borné par Security et des handoffs vers les owners canoniques, sans créer Customer, Portfolio ou autorité cross-tenant.

## 2. Problème utilisateur
Certains déploiements coordonnent plusieurs Tenants, clients ou engagements alors que d'autres sont Internal. L'utilisateur doit comprendre le scope, la source, la fraîcheur et l'audience sans confondre Tenant, Customer, permission ou autorité de réponse.

## 3. Objectifs
- activer le contexte selon le deployment model ;
- préserver un chemin Internal sans Customer ;
- afficher une projection Customer/engagement externe lorsqu'elle existe ;
- fournir un overview MSSP read-only sur un Authorized Tenant Set ;
- exiger un Tenant sélectionné pour Search, Report, Export et toute action tenant-local ;
- réutiliser Tasks Command, Reporting/Export/Search Shared, SLA Command et Govern authority sans duplication.

## 4. Non-objectifs
Ne définit pas Customer lifecycle/admin, Tenant hierarchy, ManagedTenant/TenantGroup/Portfolio/CustomerTenant, cross-tenant mutation/admin/response, delegated administration, multi-tenant Search/Report/Export initial, CRM, billing, customer portal, contract mutation, API, protocole ou implémentation.

## 5. Propriétaire
Command Product Lead est l'unique capability owner. Command possède le contexte local et les Tasks. Security possède autorisation/isolation ; Platform Settings possède Tenant ; Shared possède Search/Reporting/Export ; Govern possède Decision Authority/Response ; la source deployment/customer/contract possède Customer/engagement/obligations contractuelles.

## 6. Utilisateurs
Principal : Service Delivery Manager. Secondaires selon deployment et permissions : Customer Success, Incident Commander, Business Owner, Internal Security Lead et opérateur MSSP autorisé.

## 7. Conditions d’entrée
Deployment model déclaré (`Internal`, `Enterprise multi-tenant` ou `MSP/MSSP`), permissions serveur courantes, Tenant sélectionné pour toute action tenant-local, et Authorized Tenant Set Security lorsque l'overview agrégé est demandé. Les sources Customer/engagement restent optionnelles et deployment-specific.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Deployment model | architecture/deployment owner | Internal / Enterprise multi-tenant / MSP-MSSP | oui | configuration courante | capability `disabled-for-deployment` si contexte non applicable |
| Authorized Tenant Set | Security | projection non canonique de Tenant refs | conditionnelle, MSSP aggregate | courante | overview agrégé indisponible |
| Selected Tenant | Platform Settings / context | Tenant ref | oui pour Search/Report/Export/action | courant | action tenant-local bloquée |
| Command projections | Incident, Task, readiness/SLA | contexte de delivery | non | visible par source | `partial` |
| Result projection | Govern | outcome/residual-risk ref | non | source datée | section absente/partial |
| Customer / engagement context | external deployment/customer/contract source | identité, contrat, obligations, audience | non | version/fraîcheur exposée | aucune obligation client inférée |

## 9. Objets lus
| Objet / projection | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Tenant | Platform Settings | ref/identité de scope | lire le contexte autorisé |
| Authorized Tenant Set | Security | set de Tenant refs | consommer comme résultat d'autorisation, jamais comme objet |
| Incident / Task | Command | état, owner, SLA, outcome de coordination | read et mutation uniquement après Tenant sélectionné selon permission |
| Result | Govern | outcome vérifié et residual risk | read en projection |
| Report | Shared Reporting Engine | draft/snapshot/audience/publication | read/request selon permission et Tenant sélectionné |
| Customer / engagement | source externe | identité, obligations, SLA contractuel, audience, version | read-only projection |

## 10. Objets créés ou modifiés
| Objet / projection | Opération | Propriétaire | Règle |
|---|---|---|---|
| Report request/draft | demander/contribuer via Reporting Engine | Shared | single-Tenant ; Command ne possède pas Report |
| Delivery follow-up Task | créer/modifier | Command | Class 2 ; Tenant sélectionné ; source/raison liée ; `OPEN-013` préservée |
| Tenant | aucune mutation par cette capability | Platform Settings | context ref seulement |
| Authorized Tenant Set | aucune mutation | Security | projection d'autorisation seulement |
| Customer / contract source | aucune mutation | source externe | projection read-only |

## 11. Fonctionnalités
- reporting/service delivery Internal sans Customer ;
- activation deployment-aware ;
- external Customer/engagement projection ;
- overview MSSP read-only multi-tenant autorisé ;
- Tenant selector et context switching ;
- portfolio-like View sans objet Portfolio ;
- single-Tenant Reporting requests ;
- service-delivery Tasks ;
- contractual SLA projection distincte du SLA opérationnel ;
- source/freshness/audience visibility ;
- handoff tenant-local vers Search, Shared Reporting/Export ou Govern.

## 12. Actions utilisateur
| Action | Rôle | Objet / scope | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter overview agrégé | opérateur autorisé | projections multi-Tenant | 0 | Authorized Tenant Set + object read checks | overview read-only avec Tenant visible | non |
| Sélectionner/changer Tenant | utilisateur autorisé | contexte | 0 | Tenant dans scope autorisé | contexte tenant-local + permission re-evaluation | non |
| Créer report draft | rôle report | Report single-Tenant | 2 | Tenant sélectionné, audience/scope/template | draft Shared | workflow Reporting ; `OPEN-019` si partage externe |
| Créer follow-up | Delivery Manager | Task | 2 | Tenant sélectionné + obligation/gap source | Task Command | `OPEN-013` |
| Publier/exporter | rôle Shared autorisé | Report single-Tenant | 2 | review + permission export/publish + audience | outcome Shared | pas d'autorité Command ; `OPEN-019` si externe |
| Demander une réponse | opérateur autorisé | contexte tenant-local | 0/handoff | Tenant sélectionné + Security re-evaluation | handoff Govern | Govern Decision Authority requise |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Composer overview autorisé | oui | oui | oui | résumé facultatif | Security scope + agrégation déterministe |
| Expliquer source/fraîcheur/audience | oui | oui | oui | oui | métadonnées et règles |
| Préparer un report single-Tenant | oui | templates/snapshots | oui | brouillon attribué | template, citations, édition humaine |
| Détecter une obligation à suivre | oui | règle contractuelle si configurée | oui | suggestion attribuée | revue source + création manuelle |
| Publier/exporter | oui | validation/redaction | workflow possible | jamais automatiquement | revue/action humaine |
| Changer autorité ou scope Tenant | non autonome | Security/Govern seulement | non | interdit | processus propriétaire explicite |

## 14. États fonctionnels
`disabled-for-deployment`, `configured`, `aggregate-read-only`, `partial`, `data-stale`, `audience-restricted`, `suspended`.

`aggregate-read-only` ne donne aucune permission d'action. Les lifecycles Tenant, Report, Task, Customer source et Govern restent chez leurs propriétaires.

## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied et Stale. L'overview agrégé identifie le Tenant de chaque objet. Offline/Permission denied ne réutilisent jamais des données d'un autre Tenant comme fallback. Les mutations sont indisponibles tant qu'un Tenant unique n'est pas sélectionné.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Delivery overview | projection | rôles Command autorisés | read-only multi-Tenant borné, Tenant/source/fraîcheur visibles |
| Selected Tenant context | contexte | Command/Shared/Govern | Tenant unique, permission re-evaluation obligatoire |
| Report request/draft | Report Shared | Reporting Engine/reviewer | single-Tenant, citations/snapshot/audience, ownership Shared |
| Delivery Task | Task Command | Work Queue | tenant-scoped, attribuée et liée à une source d'engagement |
| Response handoff | refs tenant-local | Govern | aucun Response Run ni authority créé par Command |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Customers & Delivery aggregate | select Tenant | contexte Command tenant-local | Tenant ref + return origin | overview agrégé restaurable si toujours autorisé |
| contexte tenant-local | rechercher | Shared Global Search | Tenant sélectionné + query safe | même Tenant restauré |
| contexte tenant-local | créer report | Reporting Engine | Tenant sélectionné, citations, audience, return origin | contexte delivery restauré |
| contexte tenant-local | obligation/gap | Task Coordination | Tenant, engagement ref, due, owner, expected outcome | contexte delivery restauré |
| contexte tenant-local | réponse requise | Govern | Tenant, source refs, request context | retour tenant-local ; Security + Decision Authority réévaluées |
| Customers & Delivery | inspection contractuelle | source customer/contract | référence/version/return origin | même Tenant/scope restauré |

## 18. Dépendances
ADR-0008, Security Tenant Isolation et Permission Model, Platform Settings Tenant/context, Experience Context Preservation, Shared Global Search, Reporting Engine, Export Engine, Metrics Engine, Business Service Catalog, CAP-CMD-105, CAP-CMD-303, `OPEN-013` et `OPEN-019`.

## 19. Source de vérité
Command possède cette capability et les Tasks qu'elle crée. Tenant reste Platform Settings-owned. Authorized Tenant Set et authorization restent Security-owned. Report/Search/Export restent Shared-owned. Result/Decision/Response restent Govern-owned. Customer/engagement/contract reste externe/deployment-specific. Aucune projection ne transfère l'ownership.

## 20. Provenance et audit
Conserver deployment model, acteur, requested/effective scope de manière sûre, Tenant de chaque objet, source Customer/engagement, freshness/version, selected Tenant, audience, citations, snapshot, redaction, reviewer, publication/export outcome et correlation id. Les refus/tentatives cross-scope sont audités sans fuite de donnée.

## 21. Permissions fonctionnelles
Réutilise uniquement les permissions existantes : Command read/task manage selon action, permissions Shared report create/review/publish/export, source read lorsqu'elle existe, et Security RBAC/ABAC/tenant isolation. **0 nouvel ID de permission**. Aucun `cross-tenant.manage`.

## 22. Limites et erreurs
Tenant absent/non autorisé, Customer source ambiguë/stale, audience refusée, permission objet absente, tentative cross-tenant, export widening, ou absence de Tenant sélectionné empêchent l'action concernée. Aucune inférence ne crée Customer, relation de portfolio, obligation contractuelle, permission ou authority.

## 23. Métriques
Nombre de deployments Internal/Enterprise/MSSP utilisant la capability, overview requests autorisées/refusées, context switches, reports single-Tenant avec citations/snapshots, Tasks liées à un engagement, données stale/partial et refus cross-scope. Aucune cible KPI/SLO n'est sélectionnée.

## 24. Classification de livraison
`defined / planned`, deployment-dependent. Cette classification signifie spécification documentaire définie, pas implémentation, déploiement ou disponibilité runtime.

## 25. Critères d’acceptation
**Given** un déploiement Internal sans Customer, **When** un report est demandé, **Then** Reporting Engine fonctionne en single-Tenant sans Customer canonique.

**Given** un opérateur MSSP et un Authorized Tenant Set, **When** l'overview s'ouvre, **Then** seules les projections autorisées sont visibles en read-only et chaque objet conserve son Tenant.

**Given** un contexte agrégé, **When** une mutation ou réponse est demandée, **Then** elle est bloquée jusqu'à sélection d'un Tenant unique et réévaluation Security/Govern applicable.

**Given** un Tenant sélectionné et une permission Report valide, **When** un draft est créé, **Then** il reste single-Tenant et Shared-owned.

**Given** aucun modèle IA, **When** la capability est utilisée, **Then** toutes les lectures, sélections, reports et Tasks autorisés restent possibles par chemins déterministes/manuels.

## 26. Questions ouvertes
`OPEN-006` est résolue par ADR-0008. `OPEN-013` reste ouverte pour le default governance des mutations Class 2. `OPEN-019` reste ouverte pour external/client-facing dissemination, releasability et sharing. Aucun autre choix n'est inféré ici.

## 27. Consommateurs documentaires
`CMD-CRP-001`, Reporting Engine, Global Search, Export Engine, Work Queue, Mission Control, CAP-CMD-105, CAP-CMD-303, Security, Platform Settings, Govern, Experience Architecture, Quality, Roadmap et phases d'implémentation ultérieures.
