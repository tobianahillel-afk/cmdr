---
id: CAP-INV-110
title: Investigation Timeline
product: investigate
module: cases-and-evidence
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-014
  - REQ-PROD-018
open_decisions: []
---
# CAP-INV-110 — Investigation Timeline

## 1. Définition

Fournit la chronologie métier d’un Case en distinguant les événements observés, les inférences, les actions utilisateur, les Decisions et les réponses, tout en préservant la source, les corrections et la navigation vers les objets propriétaires.

## 2. Problème utilisateur

Un analyste doit reconstruire ce qui s’est produit, ce qui a été déduit et ce qui a été fait. Sans distinction explicite, une inférence peut être confondue avec un fait et une correction peut effacer l’historique utile.

## 3. Objectifs

- ordonner les événements avec timestamp, timezone, source et provenance ;
- distinguer observed, inferred, user-action, decision et response ;
- afficher les corrections, données manquantes et entrées superseded ;
- permettre filtres, comparaison, navigation, export conceptuel et préparation de replay.

## 4. Non-objectifs

- ne pas posséder le composant Timeline du Design System ;
- ne pas redéfinir Decision, Response Run ou Result ;
- ne pas inventer un moteur technique de replay ;
- ne pas finaliser le schéma de Timeline Entry.

## 5. Propriétaire

Investigate / Cases and Evidence / Investigate Product Lead possède la sélection et la sémantique métier de la timeline du Case. Shared Capabilities reste propriétaire du mécanisme Timeline transversal et de Timeline Entry selon le registre.

## 6. Utilisateurs

Principal : Case Analyst. Secondaires : Investigation Lead, Reviewer, Incident Commander en projection et approbateur Govern consultant le contexte.

## 7. Conditions d’entrée

Case accessible, tenant et environnement conservés, sources temporelles résolues, permissions sur les objets liés et politique de timezone disponible.

## 8. Entrées fonctionnelles

| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Case | Investigate | contexte et périmètre | oui | version courante | timeline indisponible |
| Events et Evidence | Shared / Investigate | faits observés et éléments qualifiés | non | timestamp et source visibles | afficher une timeline partielle |
| Actions, Hypotheses et Findings | Investigate | événements analytiques | non | versions courantes | signaler les périodes sans activité |
| Decisions, Runs et Results | Govern | événements d’autorité et de réponse | non | statut et version courants | conserver le Case sans inventer de réponse |
| Corrections et supersessions | Timeline Engine / auteurs | historique de rectification | non | appliqué à la lecture | afficher l’entrée antérieure comme non courante |

## 9. Objets lus

| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Timeline Entry | Shared Capabilities | type, source, timestamp, relations et supersession | consulter, filtrer et comparer |
| Case / Evidence / Finding / Hypothesis | Investigate | événements et versions | consulter et naviguer |
| Telemetry Event | Shared Capabilities | observation source | inspecter sans qualifier automatiquement |
| Decision / Response Run / Result | Govern | événements d’autorité, exécution et vérification | consulter et relier |
| Incident | Command | contexte opérationnel et impact | consulter et retourner |

## 10. Objets créés ou modifiés

| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Timeline Entry analytique | créer une contribution observée, inférée ou user-action | Shared Capabilities, contenu Investigate | classe 2 ; type, auteur, source et raison obligatoires |
| Timeline correction | supersede une entrée sans l’effacer | Shared Capabilities | l’entrée antérieure reste visible et référencée |
| Case relation | lier une entrée aux objets du Case | Investigate / Object Linking | relation versionnée, sans transfert d’ownership |
| Export context | préparer une sélection temporelle | Investigate / Reporting Engine | non destructif, permission-aware et sourcé |

## 11. Fonctionnalités

- afficher et filtrer la chronologie du Case ;
- distinguer faits observés, inférences, actions, Decisions et réponses ;
- comparer périodes et versions ;
- afficher données manquantes, corrections et supersessions ;
- naviguer vers chaque objet source ;
- ajouter un événement analytique autorisé ;
- fournir une sélection à Case Replay ou Reporting Preparation.

## 12. Actions utilisateur

| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Lire, filtrer ou comparer | utilisateur autorisé | Timeline | 0 | Case accessible | vue sourcée | non |
| Ouvrir l’objet source | analyste | relation | 0 | permission destination | navigation avec retour | non |
| Ajouter une entrée analytique | analyste | Timeline Entry | 2 | type et source explicites | entrée attribuée | OPEN-013 |
| Corriger une entrée | auteur ou reviewer autorisé | supersession | 2 | raison obligatoire | nouvelle version, ancienne conservée | OPEN-013 |
| Préparer un export | analyste | sélection | 1 | permission export | package non destructif | selon données sensibles |

## 13. Automatisation et IA

| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Ordonner les événements | oui | oui, timestamp et relations | oui | non nécessaire | tri et filtres déterministes |
| Classer observed/inferred | oui | règles de type | oui | proposition seulement | sélection explicite du type |
| Résumer une période | oui | agrégation factuelle | oui | brouillon attribué | lecture filtrée et synthèse humaine |
| Détecter une lacune temporelle | oui | comparaison déterministe | oui | explication possible | indicateurs de couverture temporelle |

Toute sortie automatisée expose initiateur, moteur ou agent, version, sources, timestamp, statut, incertitude, Automation Run ou Tool Calls lorsqu’ils existent, owner humain, rejet possible et trace.

## 14. États fonctionnels

Dimensions fonctionnelles : `observed`, `inferred`, `user-action`, `decision`, `response`, `missing-data`, `corrected`, `superseded`, `disputed`. Elles ne constituent pas une machine d’état unique de Timeline Entry.

## 15. États d’interface

Loading conserve période et filtres ; Empty distingue absence réelle et permission ; Partial nomme les sources manquantes ; Error conserve correlation ID ; Offline reste en lecture ; Permission denied masque les objets protégés ; Stale expose la dernière synchronisation.

## 16. Sorties

| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Timeline de Case | projection ordonnée | Case Workspace | timezone, source et type visibles |
| Entrée analytique | Timeline Entry | Case, replay et audit | attribuée, versionnée et inspectable |
| Correction | supersession event | reviewers et audit | historique antérieur conservé |
| Sélection temporelle | context package | CAP-INV-112 / CAP-INV-114 | permission-aware, sourcée et non destructive |

## 17. Transitions

| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Case Workspace | ouvrir Timeline | Investigation Timeline | Case, période, tenant, environnement et sélection | Case restauré |
| Timeline Entry | ouvrir source | Event, Evidence, Finding, Decision ou Result owner | objet, timestamp, relation et return origin | même position temporelle |
| Timeline | lancer review | CAP-INV-112 | période, entrées, corrections, lacunes et versions | filtres restaurés |
| Timeline | préparer rapport | CAP-INV-114 | sélection, citations et objets sources | timeline inchangée |

## 18. Dépendances

Timeline Engine, Trace, Object Linking, Context Bar, Inspector, CAP-INV-102, CAP-INV-107, CAP-INV-109, CAP-INV-112, CAP-INV-114 et projections Command/Govern.

## 19. Source de vérité

Les objets propriétaires restent leurs sources canoniques. Timeline Entry et mécanisme temporel restent Shared ; Investigate est source de vérité pour la sélection métier et les contributions analytiques liées au Case.

## 20. Provenance et audit

Case ID, Timeline Entry ID, type, auteur ou producteur, timestamp source et ingestion, timezone, objet source, version, correction, raison, supersession, export et correlation ID.

## 21. Permissions fonctionnelles

Case read, Timeline read/filter, analytic-entry create, correction/supersede, sensitive-object read et export. Les permissions atomiques, step-up et séparation des tâches restent à la phase Permissions.

## 22. Limites et erreurs

Horloges divergentes, timezone absente, objet inaccessible, source stale, entrée dupliquée, correction concurrente, donnée manquante ou changement d’environnement produisent un état explicite ; aucune inférence n’est présentée comme observation.

## 23. Métriques

- proportion d’entrées avec source et timezone ;
- corrections conservant l’entrée antérieure ;
- lacunes temporelles détectées ;
- retours vers objets sources réussis ;
- sélections réutilisées par replay ou reporting.

## 24. Classification de livraison

`defined` / `planned`, cible native pour la sémantique Investigate et intégration au Timeline Engine Shared. Aucun mécanisme de stockage ou replay technique n’est déclaré livré.

## 25. Critères d’acceptation

**Given** une observation et une inférence au même instant  
**When** l’utilisateur ouvre la timeline  
**Then** les deux types sont distincts, leurs sources sont visibles et l’inférence n’apparaît pas comme fait observé.

**Given** une entrée incorrecte déjà référencée  
**When** un reviewer la corrige  
**Then** une nouvelle version est créée, l’ancienne reste visible comme superseded et la raison est auditée.

**Given** un objet Govern inaccessible  
**When** la timeline est affichée  
**Then** l’existence autorisée de l’événement est indiquée sans divulgation et le Case reste utilisable.

## 26. Questions ouvertes

Le schéma final de Timeline Entry, la politique de correction, la précision temporelle et les permissions de contribution restent à la phase Objets et Permissions.

## 27. Consommateurs documentaires

Case Workspace, Evidence Board, Case Replay, Reporting Preparation, futurs écrans Timeline, phase Objets, Permissions et parcours interproduits.
