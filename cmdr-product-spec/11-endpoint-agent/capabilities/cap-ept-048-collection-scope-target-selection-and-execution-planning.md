---
id: CAP-EPT-048
title: Collection Scope, Target Selection and Execution Planning
product: endpoint-agent
module: collection
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-006, REQ-PROD-014, REQ-PROD-018, REQ-PROD-019, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-013, OPEN-014]
source-of-truth: canonical
---
# CAP-EPT-048 — Collection Scope, Target Selection and Execution Planning

## 1. Définition
Transformer un Technical Intake eligible en plan technique borné : target Endpoint/Agent, collection subjects, paths/process/system refs, période éventuelle, limites de volume, prérequis, exclusions, planned items et expected outputs, sans définir de commande finale.

## 2. Problème utilisateur
Une request métier peut rester trop large ou ambiguë pour une exécution locale sûre. Endpoint doit refuser toute extension silencieuse et rendre visible ce qui sera réellement acquis.

## 3. Objectifs
Résoudre target et subjects ; borner scope/time/volume ; vérifier source availability/prerequisites ; lister exclusions et items attendus ; calculer impact conceptuel ; préserver request/version et authority ref.

## 4. Non-objectifs
Aucun wildcard illimité, commande, moteur, chemin système universel, format, transport, Artifact/Evidence qualification ou modification du target.

## 5. Propriétaire
Endpoint possède le Collection Technical Plan local. Investigate reste propriétaire de la request/scope métier ; Settings des policies ; Govern de l’autorité.

## 6. Utilisateurs
DFIR Analyst, Endpoint Operator, Case Analyst, Response Operator, Security Reviewer et Auditor.

## 7. Conditions d’entrée
CAP-EPT-047 accepted-for-technical-planning, target résolu, requested subjects lisibles, capability/policy/authority context suffisamment courant.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Technical Intake/request ref | CAP-EPT-047 | scope source | oui | current assessment | no plan |
| target refs | EPT-1 | Agent/Endpoint | oui | current binding | blocked |
| requested subjects | Investigate | file/process/system/memory/network/triage refs | oui | request version | incomplete |
| policy/limits | Settings | bounds | oui | current | policy-unknown |
| capability/resource estimate | Endpoint | availability/impact | oui si disponible | current | estimate-unknown |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Collection Request | Investigate | scope/purpose/limits | read |
| endpoint-agent | Endpoint | target/resources/capabilities | read |
| Endpoint Policy | Settings | permitted subjects/limits | read |
| Endpoint Investigation Summary | Endpoint | missing-data refs | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Collection Technical Plan | create/version | Endpoint | request-bounded |
| Planned Collection Item | derive | Endpoint | one subject + expected output/limit |
| Plan Exclusion/Prerequisite | derive | Endpoint | explicit, no hidden fallback |

## 11. Fonctionnalités
Resolve exact target, normalize requested subjects into planned items, apply period/depth/count/size/impact bounds where relevant, record prerequisites/exclusions, estimate volume conceptually and mark unsupported items independently.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect plan | Analyst/Operator | Technical Plan | 0 | read | exact scope visible | non |
| validate bounds | Endpoint service/reviewer | plan | 1 | policy/capabilities | pass/block list | non |
| accept bounded plan | Endpoint Operator | plan | 2 | no scope expansion + authority as required | ready-for-operation | selon source class |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| decompose subjects | oui | oui | oui | suggestion | deterministic mapping |
| validate bounds | oui | oui | oui | explain | validators |
| estimate volume | oui | oui si metadata | oui | summarize | raw estimate |
| expand scope silently | non | interdit | non | interdit | explicit request revision |

## 14. États fonctionnels
`draft`, `incomplete`, `bounded`, `partially-supported`, `policy-blocked`, `impact-unknown`, `ready-for-operation`, `superseded`.

## 15. États d’interface
Aucun Screen ID. Partial/Unsupported/Stale doivent montrer items/exclusions/prerequisites séparément ; plan ready ne signifie pas started.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Technical Plan | Endpoint concept | CAP-EPT-049..053 | exact request/target/bounds |
| Planned Items | Endpoint refs | collectors | per-item scope |
| exclusions/prerequisites | diagnostics | Investigate/Govern | no hidden expansion |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-EPT-047 | accepted | CAP-EPT-048 | request/target/authority | intake |
| CAP-EPT-048 | item type selected | CAP-EPT-049..052 | item/bounds/prereqs | same plan |
| CAP-EPT-048 | execution planned | CAP-EPT-053 | items/limits | same request |

## 18. Dépendances
CAP-EPT-047, EPT-1/2/3 context, Settings Policy, Investigate CAP-INV-202/204..208, OPEN-008/013/014.

## 19. Source de vérité
Endpoint SOT du technical plan ; Investigate SOT du requested business scope. A plan cannot broaden its source request.

## 20. Provenance et audit
Request/version, target, subjects, limits, exclusions, prerequisites, capability/policy snapshots, estimates, actor/service, plan versions and correlation IDs.

## 21. Permissions fonctionnelles
Plan read/validate/accept, sensitive target/scope metadata, policy/capability read, cross-tenant deny. No final RBAC.

## 22. Limites et erreurs
Plan ≠ operation ; eligible/authorized ≠ started ; unsupported items remain unsupported ; scope drift requires revalidation/versioning ; no automatic wildcard expansion.

## 23. Métriques
Plans bounded/blocked/partial, excluded items, estimate availability, scope-reduction events, plans changed after authority snapshot.

## 24. Classification de livraison
`draft / defined / planned`; no execution engine or physical plan schema.

## 25. Critères d’acceptation
**Given** une request non bornée, **When** le plan est validé, **Then** il reste incomplete/blocked jusqu’à des limites explicites.

**Given** un item unsupported, **When** le plan contient d’autres items supportés, **Then** l’item unsupported reste explicite et le plan peut être partial sans faux support.

**Given** l’IA est indisponible, **When** le scope est planifié, **Then** mappings/validators déterministes suffisent.

## 26. Questions ouvertes
OPEN-008/013/014 restent ouvertes ; aucun format/collector/tool final n’est choisi.

## 27. Consommateurs documentaires
EPT-4 Collection capabilities, Investigate, Govern, Settings, Security, Quality et future EPT-5 boundary.
