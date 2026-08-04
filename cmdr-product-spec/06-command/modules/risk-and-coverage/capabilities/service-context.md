---
id: CAP-CMD-201
title: Service Context
product: command
module: risk-and-coverage
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids: [REQ-PROD-005, REQ-PROD-013, REQ-PROD-021, REQ-PROD-032]
open_decisions: [OPEN-013]
source-of-truth: canonical
---

# CAP-CMD-201 — Service Context

## 1. Définition
Consomme le Business Service Catalog pour relier un Incident ou une Task au service, owner, criticité, environnement et dépendances, sans créer un objet Service local.

## 2. Problème utilisateur
La coordination technique est incomplète sans service métier ni owner. Sans cette capacité, impact, priorité et escalade reposent sur des hypothèses non sourcées.

## 3. Objectifs
Présenter identité/owner/criticité/dépendances, lier les work items, montrer source/fraîcheur et signaler les lacunes au propriétaire du catalogue.

## 4. Non-objectifs
Ne devient pas une CMDB, ne crée pas Service, ne modifie pas les dépendances du catalogue et n’infère pas silencieusement un owner.

## 5. Propriétaire
Command possède la relation opérationnelle ; Shared Business Service Catalog possède le concept Service. Le fichier objet `service.md` reste absent et relève d’une future phase Objets.

## 6. Utilisateurs
Principal : Incident Commander. Secondaires : Business Owner, SOC Analyst L2, Service Delivery Manager.

## 7. Conditions d’entrée
Incident/Task ou contexte risque, référence Service explicite ou candidats à confirmer et permission catalogue.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Service reference | Incident, Task ou utilisateur | ID ou candidat | oui ou état `unknown` | version courante du work item | impact conservé sans service confirmé |
| Service metadata | Business Service Catalog | owner, criticité, dépendances et environnement | non | date de révision du catalogue | projection `partial`, aucun owner inventé |
| Operational links | Object Linking Service | Incidents, expositions et couvertures liées | non | résolution à l’ouverture | relations indisponibles signalées |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Service conceptuel | Shared Business Service Catalog | identité, owner, criticité et dépendances | consulter, filtrer, relier et naviguer |
| Incident / Task | Command | service links et impact | consulter et modifier la relation |
| Exposure / Coverage | propriétaires sources à formaliser | relation au Service | consulter en projection |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Incident / Task | ajouter, confirmer ou corriger le lien Service | Command | classe 2 ; Service reste inchangé |
| Task | créer un suivi de qualité de catalogue | Command | classe 2, lacune et résultat attendu liés |
| Service | aucune mutation | Shared Business Service Catalog | projection conceptuelle ; aucun schéma créé en Phase 4A |

## 11. Fonctionnalités
Rechercher/confirmer un Service, afficher owner/criticité/environnement/dépendances, montrer incidents/expositions, distinguer assumed/confirmed et créer un suivi de correction.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter | lecteur | Service projection | 0 | permission | contexte sourcé | non |
| Lier | coordinateur | Incident/Task | 2 | candidat choisi | relation sourcée | OPEN-013 |
| Confirmer ou contester | Business Owner/coordinateur | lien | 2 | source visible | statut du lien | OPEN-013 |
| Créer correction | coordinateur | Task | 2 | lacune identifiée | Task | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Rechercher des candidats | oui | matching explicable | oui | suggestion attribuée | recherche catalogue et sélection manuelle |
| Confirmer le lien | oui | validation de compatibilité | workflow possible | jamais silencieusement | confirmation humaine |
| Détecter une lacune | oui | règle de complétude | oui | résumé possible | contrôles de champs et Task manuelle |
| Présenter les dépendances | oui | résolution déterministe | oui | non nécessaire | projection catalogue |

## 14. États fonctionnels
`confirmed`, `assumed`, `partial`, `stale`, `owner-unknown`, `dependency-unknown`, `not-linked`.

## 15. États d’interface
Partial/Stale nomment les champs et sources ; Offline bloque la mutation ; Permission denied ne révèle pas Service ; Error conserve le lien Command valide.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Service context | projection | Mission Control, Incident Detail et Risk | source, fraîcheur et statut du lien visibles |
| Service link | relation typée | Incident/Task et Object Linking Service | auditée, permission-aware et sans transfert d’ownership |
| Catalog correction Task | Task Command | data owner | attribuée, liée à la lacune et non destructive |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Service Context | correction nécessaire | propriétaire du catalogue | Service ref, champ, source et justification | work item restauré |
| Service Context | inspection du Service | Business Service Catalog | Service ID, tenant et return origin | même Incident/Task restauré |
| Service Context | impact nécessitant autorité | Govern | Incident, Service, impact et owner | Command restauré |

## 18. Dépendances
Business Service Catalog, Object Linking Service, CAP-CMD-204, CAP-CMD-106 et CAP-CMD-107.

## 19. Source de vérité
Service et métadonnées : catalogue partagé. Relation/Task : Command. Projections d’Exposure/Coverage : sources propriétaires futures.

## 20. Provenance et audit
Acteur, candidat, source, statut assumed/confirmed, before/after, tenant et correlation ID.

## 21. Permissions fonctionnelles
Lecture Command/catalogue et gestion Incident/Task pour la relation ; atomisation reportée.

## 22. Limites et erreurs
Service absent, multiple, stale, inaccessible ou conflictuel ne produit jamais owner/impact inventé ; aucune cardinalité ou machine d’état Service n’est définie.

## 23. Métriques
Incidents avec Service confirmed/assumed, Services sans owner/criticité et délai de correction ; aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native ; preuve documentaire uniquement et dépendance à la future phase Objets.

## 25. Critères d’acceptation
**Given** deux Services candidats, **When** l’un est confirmé, **Then** source et statut sont audités sans modifier le catalogue.

**Given** aucun modèle, **When** la recherche est utilisée, **Then** catalogue, règles et sélection manuelle suffisent.

**Given** le catalogue indisponible, **When** l’Incident s’ouvre, **Then** le lien reste partial/unknown et les données Command restent utilisables.

## 26. Questions ouvertes
Objet Service canonique ou catalogue fonctionnel, et qui confirme service-impact ? — Requirement IDs ci-dessus ; `OPEN-013` reste ouverte.

## 27. Consommateurs documentaires
Mission Control, Incident Detail, Risk, parcours Phase 5, écrans Phase 6, phase Objets et permissions ultérieures.