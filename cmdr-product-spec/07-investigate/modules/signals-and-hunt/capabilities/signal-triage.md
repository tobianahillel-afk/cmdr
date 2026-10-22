---
id: CAP-INV-001
title: Signal Triage
product: investigate
module: signals-and-hunt
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-002
  - REQ-PROD-005
  - REQ-PROD-008
  - REQ-PROD-014
  - REQ-PROD-045
open_decisions:
  - OPEN-013
source-of-truth: canonical
---

# CAP-INV-001 — Signal Triage

## 1. Définition

Fournit la qualification analytique d’un Signal dans Investigate en exposant source, Detection et Events liés, confiance, sévérité, fraîcheur, relations, Incident et Case éventuels, sans redéfinir la priorité ou le cycle de vie Command.

## 2. Problème utilisateur

Le SOC Analyst L1/L2 doit décider si un Signal mérite un travail d’investigation, un pivot ou une liaison, mais une vue condensée peut masquer la source, confondre confiance et priorité ou créer implicitement un Case.

## 3. Objectifs

- rendre source, fraîcheur, confiance, sévérité et relations inspectables
- distinguer qualification analytique et priorité opérationnelle Command
- lier ou proposer un Case sans créer automatiquement Evidence ou Finding
- préserver le retour vers le Signal, l’Incident ou le Case source

## 4. Non-objectifs

- ne pas créer ou modifier un Incident
- ne pas transformer une qualification en priorité Command
- ne pas considérer un Signal comme Evidence
- ne pas finaliser la disposition ou machine d’état de Signal

## 5. Propriétaire

Investigate / Signals And Hunt / Investigate Product Lead. Investigate possède uniquement les objets et mutations explicitement listés en section 10 ; les objets consommés restent chez leur propriétaire canonique.

## 6. Utilisateurs

Principal : SOC Analyst L1/L2. Secondaires : Threat Hunter, Senior Analyst et Incident Commander en projection.

## 7. Conditions d’entrée

Tenant et environnement autorisés ; Signal résolvable ; métadonnées de source disponibles ou lacune explicitement déclarée ; permission de lecture et, pour une qualification, permission fonctionnelle de triage.

## 8. Entrées fonctionnelles

| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Signal | Command | objet source de détection | oui | état et version courants | triage impossible ; afficher source indisponible |
| Detection et Events liés | Command / Shared telemetry | projections explicatives | non | fraîcheur par source | Signal consultable mais qualification `insufficient-data` |
| Incident ou Case liés | Command / Investigate | contexte opérationnel et analytique | non | résolution à l’ouverture | proposer liaison ou nouvelle investigation sans l’inventer |

## 9. Objets lus

| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Signal / Alert | Command | source, statut, sévérité, confiance et relations | consulter, filtrer et naviguer |
| Detection | Command | règle ou détection productrice et version | consulter en projection |
| Telemetry Event | Shared Capabilities | événements sources et timestamps | inspecter et pivoter |
| Incident | Command | priorité, owner et contexte | consulter et relier |
| Case | Investigate | statut, owner et objectif | consulter et relier |

## 10. Objets créés ou modifiés

| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Signal qualification | qualifier, rejeter, marquer doublon ou insuffisant selon workflow futur | Command lifecycle avec contribution Investigate | classe 2 ; proposition ou mutation explicitement autorisée, jamais priorité Command |
| Case link | créer ou confirmer une relation | Investigate / Object Linking Service | classe 2 ; aucune Evidence créée |
| Investigation context | préparer recherche, Hunt ou demande de Case | Investigate | conserve source et return origin |

## 11. Fonctionnalités

- inspecter provenance et fraîcheur
- voir Signal, Alert, Detection et Events sans les confondre
- comparer les relations Incident/Case existantes
- préparer un pivot vers Event Search ou Hunt
- qualifier avec raison et historique
- proposer un nouveau Case ou lier un Case existant

## 12. Actions utilisateur

| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter et inspecter | analyste | Signal | 0 | permission et Signal résolu | projection sourcée | non |
| Pivoter vers recherche | analyste | Signal/Event context | 0 | champs autorisés | requête préremplie non exécutée | non |
| Qualifier | analyste autorisé | Signal qualification | 2 | source et raison visibles | qualification auditée | OPEN-013 |
| Lier au Case | analyste autorisé | Case relation | 2 | Case accessible | relation sourcée | OPEN-013 |
| Proposer un Case | analyste autorisé | Case request | 2 | objectif et contexte | brouillon de création Investigate | non |

Une action dont l’effet cible relève des classes 3 ou 4 reste une préparation ou une demande dans Investigate. L’autorité et l’exécution demeurent chez Govern ou le mécanisme propriétaire.

## 13. Automatisation et IA

| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Résumer le Signal | oui | agrégation sourcée | oui | résumé attribué | lecture des champs, Events et relations |
| Proposer une qualification | oui | règles explicables | oui | proposition seulement | checklist et décision humaine |
| Détecter un doublon potentiel | oui | matching déterministe | oui | explication facultative | comparaison des identifiants et relations |
| Construire un pivot | oui | mapping de champs | oui | brouillon de requête | sélection manuelle de champs |

Toute proposition automatisée expose initiateur, producteur, version ou run, sources, facteurs, éventuels Tool Calls, incertitude et disposition. Elle ne devient pas silencieusement un état effectif.

## 14. États fonctionnels

`new`, `under-review`, `qualified`, `linked-to-incident`, `linked-to-case`, `dismissed`, `duplicate`, `stale`, `insufficient-data`.

Ces états décrivent le travail de la capability ; ils ne finalisent pas la machine d’état canonique des objets, reportée à la phase Objets.

## 15. États d’interface

Loading conserve le Signal et le return origin ; Partial nomme les Events ou relations manquants ; Stale affiche source/date ; Permission denied ne révèle aucun champ protégé ; Offline reste en lecture. Le Design System conserve la propriété du rendu et des interactions génériques.

## 16. Sorties

| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Triage disposition | événement de qualification | Command et audit | attribué, justifié et distinct de la priorité |
| Search pivot | contexte de requête | Event Search | non exécuté, permission-aware et sourcé |
| Case relation ou proposal | relation ou brouillon | Cases and Evidence | tenant, Signal, Incident et objectif conservés |

## 17. Transitions

| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Signal Triage | pivot utilisateur | Event Search | tenant, environnement, Signal, champs, période et return origin | Signal restauré |
| Signal Triage | commencer une investigation | Case Lifecycle | Signal, Incident éventuel, objectif, service et auteur | retour au Signal ou Incident |
| Signal Triage | Case existant choisi | Case Workspace | Signal, relations et annotation | retour exact au triage |

Chaque transition réévalue les permissions, conserve le tenant et l’environnement et ne transfère jamais l’ownership du produit destination.

## 18. Dépendances

- CAP-INV-002
- CAP-INV-004
- CAP-INV-005
- CAP-INV-102
- Object Linking Service
- Trace
- Command Signal/Incident ownership

## 19. Source de vérité

Signal, Alert, Detection et Incident restent Command ; Telemetry Event reste Shared ; Case et liens d’investigation sont Investigate. La qualification ne transfère aucun ownership.

## 20. Provenance et audit

Signal ID, Detection/version, Event refs, source, fraîcheur, acteur ou producteur, qualification avant/après, raison, Case/Incident refs et correlation ID.

## 21. Permissions fonctionnelles

Besoins : Signal read/triage, raw-event access, Case read/create/link et cross-environment navigation. Les namespaces atomiques, le step-up et la séparation des tâches définitive restent à la phase Permissions.

## 22. Limites et erreurs

Signal inaccessible, source stale, données partielles, doublon incertain, conflit de version ou permission refusée produisent un état explicite ; aucune priorité, Evidence ou Finding n’est inventé.

## 23. Métriques

- temps Signal → qualification
- part des qualifications avec sources inspectées
- liaisons Case/Incident restaurables
- propositions automatisées acceptées, modifiées ou rejetées

Aucune cible chiffrée définitive n’est fixée en Phase 4B.1.

## 24. Classification de livraison

`defined` / `planned`, cible fonctionnelle native ; aucune implémentation, moteur ou workflow livré n’est prouvé.

## 25. Critères d’acceptation

**Given** un Signal visible, un Incident lié, aucun Case et un analyste autorisé  
**When** l’analyste commence une investigation  
**Then** le contexte est conservé, un Case est créé ou proposé par Investigate, le Signal reste sourcé et aucune Evidence n’est créée automatiquement

**Given** aucun modèle IA n’est disponible  
**When** le Signal est trié  
**Then** les champs, règles, pivots et actions manuelles permettent le résultat essentiel

**Given** une source Event est interdite  
**When** le Signal est consulté  
**Then** la lacune est visible sans divulgation et aucune qualification complète n’est présentée

## 26. Questions ouvertes

Quelles dispositions de Signal sont réellement modifiables par Investigate et lesquelles restent Command ? La machine d’état finale et `OPEN-013` restent ouvertes.

## 27. Consommateurs documentaires

Signals and Hunt, Case intake, futurs parcours Signal→Case, écrans pilote Triage/Event Inspector, phase Objets et Permissions.
