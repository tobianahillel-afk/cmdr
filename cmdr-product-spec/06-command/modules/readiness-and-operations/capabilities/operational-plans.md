---
id: CAP-CMD-304
title: Operational Plans
product: command
module: readiness-and-operations
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids: [REQ-PROD-005, REQ-PROD-013, REQ-PROD-021, REQ-PROD-057]
open_decisions: [OPEN-010, OPEN-013]
source-of-truth: canonical
---

# CAP-CMD-304 — Operational Plans

## 1. Définition
Maintient des plans opérationnels versionnés décrivant scénarios, rôles, responsabilités, conditions d’activation, procédures liées, owner, statut et dernière revue, sans moteur d’exécution.

## 2. Problème utilisateur
Des procédures dispersées n’exposent ni owner, review date ou conditions d’activation. Sans plan versionné, la préparation ne peut être évaluée.

## 3. Objectifs
Rendre plans/owners/conditions visibles, lier procédures/capabilities, versionner revue/supersession et permettre activation conceptuelle sans exécution.

## 4. Non-objectifs
N’exécute pas Playbook/Workflow, ne possède pas les outils, ne remplace pas Studio Workflow/Govern Playbook et ne définit pas tous les champs.

## 5. Propriétaire
Command possède le plan opérationnel de coordination ; Studio/Govern et sources conservent les procédures exécutables.

## 6. Utilisateurs
Principal : Readiness Coordinator. Secondaires : Incident Commander, Business Owner, Response Operator.

## 7. Conditions d’entrée
Scénario, owner, audience, scope, procédures de référence et permission.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Plan definition | owner ou Readiness Coordinator | scénario, rôles et conditions | oui | version courante | plan incomplet, activation bloquée |
| Procedure links | sources documentaires, Govern ou Studio | références de procédures/playbooks/workflows | non | version et dernière revue | `stale-references` ou gap visible |
| Capability requirements | Capability Readiness | capacités et readiness requis | non | assessment datée | activation marquée partielle |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Operational Plan | Command Readiness | scénario, rôles, conditions et version | consulter et modifier |
| Playbook / Workflow / procedure | Govern, Studio ou source | référence, version et statut | consulter, relier et naviguer |
| Capability Readiness | Command | statut, evidence et gaps | consulter et relier |
| Service / Incident | propriétaires sources | contexte d’activation | consulter en projection |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Operational Plan | créer, modifier, revoir, activer ou superséder | Command Readiness | classe 2, version et owner obligatoires |
| Task | créer un suivi de revue ou gap | Command | classe 2, raison et due liés |
| Playbook / Workflow / procedure | aucune mutation | propriétaires sources | références en lecture seule |

## 11. Fonctionnalités
Créer/versionner, définir activation/rôles, lier sans copier, montrer revue/gaps et gérer active/expired/superseded.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Créer ou modifier | coordinateur | Plan | 2 | owner/scénario | draft/version | OPEN-013 |
| Soumettre ou revoir | owner/reviewer | Plan | 2 | complet | review/active | OPEN-013 |
| Superséder | owner | Plan | 2 | replacement lié | superseded | OPEN-013 |
| Créer review Task | coordinateur | Task | 2 | due ou stale | Task | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Vérifier la complétude | oui | oui | oui | explication facultative | checklist versionnée |
| Résoudre les références | oui | oui | oui | non nécessaire | Object Linking/Versioning |
| Préparer un résumé | oui | agrégation possible | oui | brouillon attribué | résumé humain du plan |
| Activer ou superséder | oui | validation/version | workflow possible | jamais automatiquement | décision humaine autorisée |

## 14. États fonctionnels
`draft`, `review`, `active`, `expired`, `superseded`, `withdrawn`, `stale-references`.

## 15. États d’interface
Stale refs et owners manquants sont visibles ; Offline bloque activation/supersession ; Permission denied masque les procédures ; Error conserve la version courante.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Operational Plan | record Command | Incident Commander et Exercise Coordination | versionné, owner et conditions visibles |
| Plan activation event | événement | Mission Control et Audit | acteur, version et conditions conservés |
| Review/gap Task | Task Command | owner et Work Queue | raison, due et plan source liés |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Operational Plan | exécution souhaitée | Studio ou Govern owner | Plan ref, Incident, conditions et return origin | Plan restauré, aucune exécution locale |
| Operational Plan | exercice préparatoire | Exercise Coordination | scénario, rôles et objectifs | Plan restauré |
| Operational Plan | procédure stale | Task Coordination / source owner | référence, version, gap et due | Plan restauré |

## 18. Dépendances
CAP-CMD-301, CAP-CMD-302, CAP-CMD-305, Object Linking, Versioning et Notification Center.

## 19. Source de vérité
Plan record : Command. Procédures/Playbooks/Workflows : propriétaires sources.

## 20. Provenance et audit
Version, author/reviewer, source refs, before/after, activation status, tenant et correlation ID.

## 21. Permissions fonctionnelles
Command read/coordinate et futures permissions review/activate ; `OPEN-010/013` restent ouvertes.

## 22. Limites et erreurs
Procédure stale/inaccessible, capability gap, owner absent, conflit ou refus empêchent un état `active` trompeur.

## 23. Métriques
Plans avec owner/review date, plans expired liés à scénarios actifs et gaps avant activation ; aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native ; aucun moteur d’exécution défini.

## 25. Critères d’acceptation
**Given** un plan complet, **When** un reviewer l’active, **Then** version, owner, conditions et refs sont auditables.

**Given** une procédure stale, **When** le plan est consulté, **Then** `stale-references` est visible et une Task peut être créée.

**Given** aucun modèle, **When** le plan est créé/revu, **Then** édition et validation manuelles suffisent.

## 26. Questions ouvertes
Comment distinguer Operational Plan et Govern Playbook, et quelles transitions exigent un reviewer distinct ? — Requirement IDs ci-dessus ; `OPEN-010/013` restent ouvertes.

## 27. Consommateurs documentaires
Readiness, Exercises, Incident coordination, parcours Phase 5, écrans Phase 6, phase Objets et permissions ultérieures.