---
id: CAP-INV-003
title: Query Authoring and Assistance
product: investigate
module: signals-and-hunt
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-AI-001
  - REQ-AI-010
  - REQ-AI-011
  - REQ-PROD-010
  - REQ-PROD-045
open_decisions:
  - OPEN-015
source-of-truth: canonical
---

# CAP-INV-003 — Query Authoring and Assistance

## 1. Définition

Fournit l’écriture manuelle, le builder, l’autocomplete, la validation, l’explication, l’historique, les paramètres et l’assistance IA facultative pour une Query lisible et modifiable.

## 2. Problème utilisateur

Un analyste doit écrire des recherches reproductibles sans mémoriser tout le schéma. Une aide qui masque la requête finale ou exécute automatiquement une proposition coûteuse retire le contrôle et la traçabilité.

## 3. Objectifs

- rendre toute requête finale visible et éditable
- fournir autocomplete, validation, documentation et snippets sans IA
- attribuer chaque proposition IA avec sources et run
- séparer authoring, validation et exécution

## 4. Non-objectifs

- ne pas définir le dialecte final
- ne pas transformer le builder en langage concurrent
- ne pas exécuter automatiquement une proposition
- ne pas faire d’un Automation Agent le propriétaire de Query

## 5. Propriétaire

Investigate / Signals And Hunt / Investigate Product Lead. Investigate possède uniquement les objets et mutations explicitement listés en section 10 ; les objets consommés restent chez leur propriétaire canonique.

## 6. Utilisateurs

Principal : SOC Analyst L2 et Threat Hunter. Secondaires : Senior Analyst, Detection Engineer et utilisateurs guidés par templates.

## 7. Conditions d’entrée

Scope de données connu ; métadonnées de champs accessibles ; éventuel modèle autorisé mais jamais requis ; permission d’utiliser les sources ou exemples correspondants.

## 8. Entrées fonctionnelles

| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Authoring context | utilisateur / Case / Hunt | objectif, champs, période et paramètres | oui | contexte courant | éditeur manuel disponible |
| Field catalog | Query Engine / data sources | champs, types et exemples autorisés | oui | version du schéma | documentation statique et validation partielle |
| Assistance request | utilisateur | question et contraintes | non | run courant | aucun impact sur l’authoring manuel |

## 9. Objets lus

| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Query | Shared Capabilities | texte, paramètres et versions | consulter, modifier un brouillon et comparer |
| Field metadata | sources / Query Engine | schéma et disponibilité | consulter |
| Case / Hypothesis / Hunt context | Investigate | objectif et objets liés | utiliser pour préparer, sans mutation |
| Automation Run | CMDR Studio | run, sources et Tool Calls | inspecter si assistance utilisée |

## 10. Objets créés ou modifiés

| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Query draft | créer ou modifier un brouillon local | Shared Query object / Investigate usage | classe 2 ; versionné avant partage |
| Query proposal | produire un brouillon attribué | CMDR Studio / Investigate proposal | jamais effective ou exécutée automatiquement |
| Authoring history | enregistrer versions et dispositions | Versioning / Trace | append-only conceptuel |

## 11. Fonctionnalités

- éditeur texte et builder déterministe
- autocomplete et validation de schéma
- paramètres et variables
- historique et diff
- snippets et exemples
- explication déterministe
- proposition IA attribuée
- acceptation, modification ou rejet

## 12. Actions utilisateur

| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Écrire ou construire | analyste | Query draft | 2 | scope connu | brouillon versionné | OPEN-013 |
| Valider | analyste | Query draft | 0 | catalogue disponible | diagnostics | non |
| Demander une proposition | analyste | Query proposal | 0 | assistance autorisée | brouillon attribué | non |
| Accepter comme brouillon | analyste | Query draft | 2 | proposition inspectée | nouvelle version humaine | OPEN-013 |
| Rejeter | analyste | proposal | 2 | proposition active | disposition auditée | non |

Une action dont l’effet cible relève des classes 3 ou 4 reste une préparation ou une demande dans Investigate. L’autorité et l’exécution demeurent chez Govern ou le mécanisme propriétaire.

## 13. Automatisation et IA

| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Autocomplete | oui | catalogue et grammaire | oui | non nécessaire | catalogue déterministe |
| Valider | oui | parse/type checks | oui | explication facultative | diagnostics déterministes |
| Expliquer | oui | plan/AST conceptuel | oui | résumé attribué | documentation et structure de la Query |
| Proposer une Query | oui | templates/règles | oui | brouillon seulement | écriture manuelle, snippets et builder |

Toute proposition automatisée expose initiateur, producteur, version ou run, sources, facteurs, éventuels Tool Calls, incertitude et disposition. Elle ne devient pas silencieusement un état effectif.

## 14. États fonctionnels

`empty`, `draft`, `valid`, `invalid`, `proposed`, `modified`, `accepted-as-draft`, `rejected`, `superseded`.

Ces états décrivent le travail de la capability ; ils ne finalisent pas la machine d’état canonique des objets, reportée à la phase Objets.

## 15. États d’interface

Loading ne remplace pas le texte ; diagnostics restent associés à une version ; Permission denied retire les exemples sensibles ; Offline conserve le draft local récupérable ; conflit propose diff et nouvelle version. Le Design System conserve la propriété du rendu et des interactions génériques.

## 16. Sorties

| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Query draft | objet ou brouillon versionné | Event Search et Saved Searches | lisible, modifiable et non exécuté |
| Diagnostics | résultat déterministe | auteur | champs, positions et limites visibles |
| AI proposal disposition | événement | Trace / Studio | agent, run, sources et accept/reject conservés |

## 17. Transitions

| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Case ou Hunt | demande d’aide | Query Authoring | objectif, Hypothesis, champs et contraintes | retour au contexte source |
| Query Authoring | exécuter explicitement | Event Search | version, paramètres, scope et return origin | éditeur restauré |
| Query Authoring | sauvegarder | Saved Searches | Query draft, metadata et dépendances de champs | retour à l’éditeur |

Chaque transition réévalue les permissions, conserve le tenant et l’environnement et ne transfère jamais l’ownership du produit destination.

## 18. Dépendances

- CAP-INV-002
- CAP-INV-006
- CAP-INV-008
- Shared Query Engine
- Schema catalog
- CMDR Studio optional assistance
- OPEN-015

## 19. Source de vérité

Query et Search Job restent Shared ; Investigate possède le contexte d’auteur, les assets de recherche et la disposition de propositions ; Automation Run reste Studio.

## 20. Provenance et audit

Auteur humain, source du draft, versions, diff, field catalog/version, Case/Hunt/Hypothesis context, agent/run/Tool Calls, acceptation ou rejet.

## 21. Permissions fonctionnelles

Besoins : query create/edit, field metadata read, model assistance opt-in, sensitive example masking et saved-search share. Les namespaces atomiques, le step-up et la séparation des tâches définitive restent à la phase Permissions.

## 22. Limites et erreurs

Schéma absent, champ interdit, proposition non fondée, contexte trop large, modèle indisponible ou conflit de version ne bloquent pas l’écriture manuelle et sont explicités.

## 23. Métriques

- part des queries validées avant exécution
- propositions acceptées après modification versus telles quelles
- erreurs de champ détectées avant run
- usage des voies sans IA

Aucune cible chiffrée définitive n’est fixée en Phase 4B.1.

## 24. Classification de livraison

`defined` / `planned`; assistance IA et moteurs restent dépendants de futures implémentations, sans fournisseur obligatoire.

## 25. Critères d’acceptation

**Given** un Case actif, une Hypothesis et une assistance autorisée  
**When** l’utilisateur demande une Query  
**Then** le brouillon est attribué, lisible, modifiable, non exécuté, rejetable et les sources sont visibles

**Given** aucun modèle n’est disponible  
**When** l’utilisateur écrit une Query  
**Then** autocomplete, documentation, builder, validation et historique restent disponibles

**Given** une proposition référence un champ non autorisé  
**When** elle est validée  
**Then** le diagnostic est visible et aucune exécution ni divulgation n’a lieu

## 26. Questions ouvertes

Quel dialecte et quel catalogue de schéma seront canoniques ? Comment relier un Automation Run d’assistance sans préjuger `OPEN-015` ?

## 27. Consommateurs documentaires

Event Search, Hunt, Saved Searches, Case Workspace, futurs écrans d’authoring et phases Technique/Permissions.
