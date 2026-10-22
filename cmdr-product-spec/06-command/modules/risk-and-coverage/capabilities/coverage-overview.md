---
id: CAP-CMD-203
title: Coverage Overview
product: command
module: risk-and-coverage
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids: [REQ-PROD-005, REQ-PROD-013, REQ-PROD-032, REQ-PROD-037]
open_decisions: [OPEN-013]
source-of-truth: canonical
---

# CAP-CMD-203 — Coverage Overview

## 1. Définition
Fournit une lecture sourcée des couvertures détection, endpoint, réponse, données et procédures sans posséder les moteurs, règles ou flotte.

## 2. Problème utilisateur
Un Service peut sembler protégé alors qu’une source, un agent, une procédure ou capacité de réponse est absente. Sans séparation par famille, un total masque les lacunes.

## 3. Objectifs
Séparer les familles, afficher définition/population/source/fraîcheur, relier les gaps aux Services/work items et créer une Task d’amélioration sans modifier les moteurs.

## 4. Non-objectifs
Ne gère pas les règles de détection, la flotte Endpoint, les agents Studio et ne produit pas un pourcentage global opaque.

## 5. Propriétaire
Command possède la lecture opérationnelle ; Investigate, Settings, Endpoint Agent, Govern, Studio et Readiness possèdent leurs sources respectives.

## 6. Utilisateurs
Principal : Incident Commander. Secondaires : Detection Engineer, Platform Administrator et Readiness Coordinator en consultation.

## 7. Conditions d’entrée
Service ou scope, projections autorisées et définition de métrique ou état `unknown`.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Detection coverage | Investigate | projection de couverture | non | version et date | famille `unavailable` |
| Endpoint/data/response coverage | Settings, Endpoint Agent et Govern | projections par famille | non | timestamp de chaque source | gap propre à la famille |
| Procedural coverage | Readiness and Operations | plans et tests | non | date de revue/test | `not-tested` ou `unknown` |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Coverage projection | produit source | définition, population, état et fraîcheur | consulter, filtrer et naviguer |
| Service conceptuel | Business Service Catalog | criticité et relations | consulter et relier |
| Task | Command | action d’amélioration | consulter et modifier si autorisé |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Task | créer une action d’amélioration | Command | classe 2, gap et validation attendue liés |
| Coverage source | aucune mutation | produit propriétaire | projection en lecture seule |
| Coverage object | aucun objet concurrent créé | phase Objets future | définitions et cardinalités restent sources-owned |

## 11. Fonctionnalités
Afficher cinq familles, montrer définition/population/données manquantes, lier gaps, ouvrir le propriétaire source et créer une action de suivi.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter | lecteur | projection | 0 | source autorisée | limites visibles | non |
| Ouvrir la source | lecteur | produit source | 0 | permission | transition | non |
| Signaler un gap | coordinateur | Task | 2 | gap sourcé | Task | OPEN-013 |
| Suivre la validation | Readiness Coordinator | Task | 2 | owner/due | follow-up | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Résoudre les familles | oui | oui | oui | non nécessaire | projections sources |
| Calculer les indicateurs | revue humaine | oui, définitions versionnées | oui | résumé facultatif | métriques déterministes et tableau source |
| Détecter un gap | oui | règles possibles | oui | suggestion attribuée | comparaison manuelle et règles |
| Créer une amélioration | oui | validation/déduplication | oui | brouillon Task | création manuelle de Task |

## 14. États fonctionnels
`covered`, `partial`, `gap`, `degraded`, `not-tested`, `unavailable`, `unknown`.

## 15. États d’interface
Chaque famille peut être Partial indépendamment ; source/date/définition restent visibles ; Offline bloque mutation mais conserve la dernière projection.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Coverage context | projection multi-source | Risk, Mission Control et Readiness | familles séparées, sources et fraîcheur visibles |
| Coverage gap | événement fonctionnel | Readiness et owner source | attribué et non confondu avec un score global |
| Improvement Task | Task Command | Work Queue | gap, Service, owner et validation attendue liés |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Coverage Overview | ouverture d’une famille | Investigate, Settings, Govern ou Studio | Service, famille, source et return origin | même vue Coverage restaurée |
| Coverage Overview | gap opérationnel | Improvement Actions | source, gap, Service et résultat attendu | Coverage Overview restauré |
| Coverage Overview | comparaison Readiness | Readiness Overview | Service, familles et dates | Risk and Coverage restauré |

## 18. Dépendances
CAP-CMD-201, CAP-CMD-301, CAP-CMD-303, Metrics Engine, Business Service Catalog et Data Quality Service.

## 19. Source de vérité
Chaque produit source possède sa couverture ; Command agrège uniquement les projections.

## 20. Provenance et audit
Définition/version/source/population/fraîcheur, Task et acteur sont conservés.

## 21. Permissions fonctionnelles
Lecture Command/source et `perm.command.task.manage` ; atomisation reportée.

## 22. Limites et erreurs
Définitions incompatibles, données manquantes, double comptage, source stale ou refus empêchent toute comparaison opaque.

## 23. Métriques
Familles avec définition/fraîcheur, gaps sans owner/action et délai gap→validation ; aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native ; sources techniques et objets détaillés non définis ici.

## 25. Critères d’acceptation
**Given** détection couverte et endpoint indisponible, **When** Coverage est consultée, **Then** les familles restent séparées et aucun total complet n’est affiché.

**Given** un gap sourcé, **When** une Task est créée, **Then** source, Service, owner et validation attendue sont liés.

**Given** aucun modèle, **When** la capability est utilisée, **Then** projections et métriques déterministes suffisent.

## 26. Questions ouvertes
Quelles définitions sont comparables et comment éviter le double comptage ? — Requirement IDs ci-dessus ; `OPEN-013` reste ouverte.

## 27. Consommateurs documentaires
Risk, Mission Control, Readiness, parcours Phase 5, écrans Phase 6, phase Objets et permissions ultérieures.