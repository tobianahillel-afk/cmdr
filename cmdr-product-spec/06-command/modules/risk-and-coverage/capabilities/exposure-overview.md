---
id: CAP-CMD-202
title: Exposure Overview
product: command
module: risk-and-coverage
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids: [REQ-PROD-006, REQ-PROD-013, REQ-PROD-032, REQ-PROD-037]
open_decisions: [OPEN-013]
source-of-truth: canonical
---

# CAP-CMD-202 — Exposure Overview

## 1. Définition
Présente les expositions utiles à la coordination avec source, état, fraîcheur et relations Service/Incident, sans scanner ni créer une source d’exposition.

## 2. Problème utilisateur
Les coordinateurs doivent comprendre le risque sans se substituer aux outils de vulnérabilité ou posture. Sans capacité, exposition et contexte opérationnel restent séparés.

## 3. Objectifs
Agréger les projections autorisées, relier Exposure/Service/Incident/Coverage, rendre source/fraîcheur visibles et créer Task/lien sans modifier la source.

## 4. Non-objectifs
Ne scanne pas, ne définit pas un objet Exposure complet, ne confirme pas une vulnérabilité et ne calcule pas de score universel.

## 5. Propriétaire
Command possède la projection opérationnelle ; la source externe ou Shared Capabilities possède le concept Exposure. Le fichier `exposure.md` reste absent et relève de la phase Objets.

## 6. Utilisateurs
Principal : Incident Commander. Secondaires : SOC Analyst L2, Business Owner et owner de la source d’exposition.

## 7. Conditions d’entrée
Source configurée ou état partial, tenant/Service et permission de la source.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Exposure projection | intégration ou source owner | résumé d’exposition | non | timestamp de la source | `source-unavailable` ou `unknown` |
| Service links | Business Service Catalog | relation et criticité | non | date du catalogue | exposition sans Service confirmé |
| Incident links | Command | relations opérationnelles | non | version courante | aucune relation Incident affichée |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Exposure conceptuelle | source externe ou Shared à confirmer | type, état, source et fraîcheur | consulter, filtrer et ouvrir la source |
| Service conceptuel | Shared Business Service Catalog | relation et criticité | consulter et relier |
| Incident / Task | Command | relation et état | consulter et modifier le lien |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Incident / Task | ajouter un lien Exposure ou un suivi | Command | classe 2, référence stable et audit |
| Exposure source | aucune mutation | propriétaire externe | projection en lecture seule |
| Exposure object | aucun objet créé | phase Objets future | aucun schéma, état ou cardinalité inventé |

## 11. Fonctionnalités
Filtrer par Service/source/état/fraîcheur, inspecter la définition, lier Incident, créer Task et comparer Exposure/Coverage sans fusion.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter | lecteur | Exposure projection | 0 | permission | détail sourcé | non |
| Lier Incident | coordinateur | relation | 2 | objets compatibles | lien audité | OPEN-013 |
| Créer Task | coordinateur | Task | 2 | action attendue | Task liée | non |
| Ouvrir la source | lecteur | source externe | 0 | deep link autorisé | transition avec retour | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Résoudre les projections | oui | oui | oui | non nécessaire | intégration et lecture de source |
| Matcher Exposure et Service | confirmation humaine | règle explicable | oui | suggestion attribuée | recherche et liaison manuelles |
| Proposer un Incident lié | oui | matching sourcé | oui | proposition seulement | filtres et sélection manuelle |
| Créer un suivi | oui | validation/déduplication | oui | brouillon Task | création manuelle de Task |

## 14. États fonctionnels
`active`, `mitigated-source`, `accepted-source`, `unknown`, `stale`, `source-unavailable`, `unlinked`.

## 15. États d’interface
Partial/Stale/Unavailable sont explicites ; Offline bloque le lien ; Permission denied masque la projection ; Error conserve les liens Command valides.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Exposure context | projection sourcée | Risk, Mission Control et Incident Detail | source, fraîcheur et limites visibles |
| Exposure link | relation typée | Incident/Task | auditée, permission-aware et sans mutation source |
| Follow-up Task | Task Command | Work Queue | attribuée et liée à l’Exposure source |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Exposure Overview | ouverture source | outil ou produit propriétaire | source ID, Service, tenant et return origin | même vue Risk restaurée |
| Exposure Overview | investigation requise | Investigate | Incident, Exposure ref, période et question | Risk/Incident restauré |
| Exposure Overview | suivi opérationnel | Task Coordination | source, Service et résultat attendu | Exposure Overview restauré |

## 18. Dépendances
Business Service Catalog, Object Linking Service, Data Quality Service, CAP-CMD-201 et CAP-CMD-203.

## 19. Source de vérité
La source externe possède Exposure ; Command possède uniquement liens Incident/Task et projection opérationnelle.

## 20. Provenance et audit
Source/version/time, relation, acteur, justification et correlation ID.

## 21. Permissions fonctionnelles
Lecture Command et source-specific ; gestion Incident/Task pour le lien ; atomisation reportée.

## 22. Limites et erreurs
Source down, projection stale, lien ambigu, Service inconnu ou permission refusée ne deviennent jamais une Exposure confirmée ; aucun objet complet n’est créé.

## 23. Métriques
Projections avec source/fraîcheur, liens Service/Incident et sources indisponibles ; aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native ; aucun scanner ou objet Exposure livré n’est prouvé.

## 25. Critères d’acceptation
**Given** une Exposure externe autorisée, **When** elle est liée à un Incident, **Then** la relation est auditée et la source reste owner.

**Given** la source indisponible, **When** la vue s’ouvre, **Then** état/fraîcheur restent visibles et aucune donnée n’est inventée.

**Given** aucun modèle, **When** la capability est utilisée, **Then** filtres, sources et actions manuelles suffisent.

## 26. Questions ouvertes
Quel objet/source canonique et quelles expositions sont pertinentes sans devenir une liste de vulnérabilités ? — Requirement IDs ci-dessus ; `OPEN-013` reste ouverte.

## 27. Consommateurs documentaires
Risk screens, Incident Detail, transition Investigate, parcours Phase 5, écrans Phase 6, phase Objets et permissions ultérieures.