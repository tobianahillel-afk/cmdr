---
id: CAP-INV-112
title: Case Replay and Investigation Review
product: investigate
module: cases-and-evidence
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-014
  - REQ-PROD-020
open_decisions:
  - OPEN-013
  - OPEN-015
---
# CAP-INV-112 — Case Replay and Investigation Review

## 1. Définition

Reconstruit fonctionnellement le déroulement d’une investigation à partir de ses requêtes, Hypotheses, Evidence, Findings, actions et sorties automatisées afin d’identifier les changements, erreurs, éléments manquants et leçons, sans prétendre disposer d’un moteur technique de relecture intégrale.

## 2. Problème utilisateur

Après une investigation, les équipes voient le résultat final mais comprennent difficilement quelles étapes ont conduit aux Findings, quels pivots ont été abandonnés et quelles erreurs ou opportunités ont été manquées.

## 3. Objectifs

- reconstruire une chronologie documentée et versionnée ;
- comparer états, requêtes, Hypotheses, Evidence et Findings ;
- rendre visibles erreurs, lacunes, décisions analytiques et sorties automatisées ;
- préparer un retour d’expérience vers Detection Engineering ou Command Readiness.

## 4. Non-objectifs

- ne pas exécuter à nouveau les outils ou requêtes par défaut ;
- ne pas définir un moteur de replay forensic ;
- ne pas modifier rétroactivement les objets sources ;
- ne pas commencer les capabilities Detection Engineering ou Readiness.

## 5. Propriétaire

Investigate / Cases and Evidence / Investigate Product Lead possède la reconstruction analytique du Case. Les objets sources, Timeline Engine, Automation Run et produits consommateurs conservent leur ownership.

## 6. Utilisateurs

Principal : Investigation Lead ou Reviewer. Secondaires : Case Analyst, Detection Engineer futur, Incident Commander et responsable Readiness.

## 7. Conditions d’entrée

Case lisible, historique et versions disponibles, accès aux Queries/Search Jobs/Automation Runs autorisé, période de revue définie et reviewer identifié.

## 8. Entrées fonctionnelles

| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Case et Timeline | Investigate / Shared | périmètre et séquence | oui | version de clôture ou snapshot choisi | review impossible |
| Queries et Search Jobs | Shared Capabilities | recherches exécutées | non | versions et résultats conservés | marquer la lacune |
| Hypotheses, Evidence et Findings | Investigate | raisonnement et conclusions | oui pour review complète | versions au snapshot | review partielle si incomplet |
| Actions utilisateur | Activity Stream / audit | décisions et modifications | non | horodatées | indiquer les actions inconnues |
| Sorties automatisées | CMDR Studio / Trace | propositions, Tool Calls et runs | non | version et statut visibles | ne pas attribuer à l’automatisation |

## 9. Objets lus

| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case / Hypothesis / Evidence / Finding / Artifact | Investigate | versions, relations et états | consulter et comparer |
| Query / Search Job / Timeline Entry | Shared Capabilities | versions, runs et chronologie | consulter et relier |
| Automation Run / Tool Call | CMDR Studio | initiateur, version, sources et sorties | inspecter en projection |
| Incident / Task | Command | contexte et actions de suivi | consulter et proposer une sortie |
| Result | Govern | résultat vérifié et effet | consulter sans le modifier |

## 10. Objets créés ou modifiés

| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Case review record | créer un snapshot de revue et conclusions | Investigate | classe 2 ; ne modifie pas les objets sources |
| Improvement observation | créer une observation liée | Investigate | doit citer l’étape et la preuve |
| Detection Engineering handoff | préparer un package futur | Investigate → future module | aucune capability 4B.3 créée |
| Readiness handoff | préparer une Task ou recommandation | Command | via transition explicite, sans ownership concurrent |

## 11. Fonctionnalités

- sélectionner un snapshot ou une période ;
- reconstruire la séquence de recherches et décisions ;
- comparer versions de Queries, Hypotheses, Evidence et Findings ;
- identifier erreurs, informations manquantes, pivots abandonnés et actions manquées ;
- documenter leçons et recommandations ;
- préparer des handoffs vers Detection Engineering futur ou Command Readiness.

## 12. Actions utilisateur

| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Lire et comparer | reviewer | snapshots et objets | 0 | permissions sources | reconstruction sourcée | non |
| Créer une review | Investigation Lead | Case review record | 2 | scope et snapshot | review versionnée | OPEN-013 |
| Marquer une erreur ou lacune | reviewer | observation | 2 | source citée | observation attribuée | OPEN-013 |
| Préparer un handoff Detection | reviewer | package futur | 2 | problème et preuves | package non déployé | non |
| Préparer une action Readiness | reviewer | Task/recommandation | 2 | owner cible | handoff Command | OPEN-013 |

## 13. Automatisation et IA

| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Reconstituer la séquence | oui | oui, relations et timestamps | oui | non nécessaire | timeline et audit déterministes |
| Comparer des versions | oui | diff déterministe | oui | explication facultative | diff et inspection manuels |
| Détecter une étape manquante | oui | règles de complétude | oui | suggestion | checklist de revue |
| Préparer des leçons | oui | agrégation de faits | workflow de revue | brouillon attribué | synthèse humaine structurée |

Toute sortie automatisée expose initiateur, producteur, version, Automation Run ou Tool Calls, sources, statut, incertitude, owner humain, possibilité de rejet et trace. Une suggestion ne devient pas une conclusion de review.

## 14. États fonctionnels

`preparing`, `ready-for-review`, `in-review`, `partial`, `completed`, `disputed`, `superseded`. Les objets sources conservent leurs propres états et ne sont jamais réécrits par le replay.

## 15. États d’interface

Loading conserve snapshot et comparaison ; Empty distingue Case sans historique et permission ; Partial liste les sources manquantes ; Error préserve les observations ; Offline reste en lecture ; Permission denied masque les détails protégés ; Stale signale un Case modifié après le snapshot.

## 16. Sorties

| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Review de Case | record versionné | Investigation Lead et audit | snapshot, auteur et sources explicites |
| Liste d’erreurs/lacunes | observations | équipe Case | attribuée et reliée aux étapes |
| Handoff Detection | package conceptuel | future Detection Engineering | non déployé, exigences et preuves visibles |
| Handoff Readiness | Task ou recommandation | Command Readiness | owner, priorité contextuelle et return origin |

## 17. Transitions

| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Case Workspace | ouvrir review | Case Replay | Case, snapshot, timeline et objets liés | Case restauré |
| Review | ouvrir étape | Query, Timeline, Hypothesis, Evidence ou Finding | version, timestamp, raison et return origin | même point de review |
| Review | préparer détection | future Detection Engineering | problème, données, Queries, lacune et preuves | review inchangée |
| Review | créer action Readiness | Command | leçon, impact, owner proposé et Case source | retour à la review |

## 18. Dépendances

CAP-INV-008, CAP-INV-102, CAP-INV-103, CAP-INV-107, CAP-INV-109, CAP-INV-110, Trace, Activity Stream, Versioning, Automation Run projection, Command Task/Readiness et futures phases 4B.3.

## 19. Source de vérité

Les objets et événements sources restent canoniques chez leurs propriétaires. La review est une reconstruction versionnée Investigate ; elle ne remplace ni Timeline, ni audit, ni Result.

## 20. Provenance et audit

Case ID, snapshot, période, reviewer, objets et versions lus, comparaisons, erreurs, exclusions, assistance automatisée, handoffs et correlation IDs.

## 21. Permissions fonctionnelles

Case review read/create, audit access, Query/Search Job history, Automation Run inspect, sensitive Evidence read et handoff create. Step-up et séparation des tâches restent à finaliser.

## 22. Limites et erreurs

Historique incomplet, objet supprimé ou inaccessible, version expirée, horloge divergente, Automation Run non disponible, Case modifié pendant la review ou handoff refusé produisent un état partiel explicite.

## 23. Métriques

- reviews avec snapshot complet ;
- lacunes reliées à une étape précise ;
- actions d’amélioration acceptées ;
- différences entre Findings initiaux et finaux ;
- sorties automatisées correctement attribuées.

## 24. Classification de livraison

`defined` / `planned`, cible native pour la reconstruction fonctionnelle. Aucun moteur de replay technique, Detection Engine ou simulateur n’est déclaré livré.

## 25. Critères d’acceptation

**Given** un Case avec plusieurs versions de Query et Findings  
**When** le reviewer compare deux snapshots  
**Then** les différences, auteurs, dates et objets sources sont visibles sans modifier les versions d’origine.

**Given** un historique incomplet  
**When** la review est clôturée  
**Then** elle reste `partial`, les lacunes sont nommées et aucune conclusion complète n’est prétendue.

**Given** une leçon pour la détection  
**When** un handoff est préparé  
**Then** un package conceptuel est créé sans démarrer ni spécifier une capability 4B.3.

## 26. Questions ouvertes

La forme canonique du review record, sa rétention, la gouvernance classe 2 et le bridge Automation Run/Response Run restent ouverts via OPEN-013 et OPEN-015.

## 27. Consommateurs documentaires

Case Workspace, futures revues post-investigation, Detection Engineering 4B.3, Command Readiness, phase Objets, Permissions et Reporting.
