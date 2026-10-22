---
id: CAP-INV-338
title: Binary Diffing and Version Comparison
product: investigate
module: analysis-workbench
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-05
requirement_ids:
  - REQ-INV-003
  - REQ-PROD-014
  - REQ-PROD-020
  - REQ-AI-002
  - REQ-UX-006
open_decisions:
  - OPEN-005
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---

# CAP-INV-338 — Binary Diffing and Version Comparison

## 1. Définition
Binary Diffing and Version Comparison définit le comportement produit permettant de comparer plusieurs Artifacts ou versions avec préconditions et résultats partiels visibles sans fusionner les originaux ni qualifier automatiquement une différence.

## 2. Problème utilisateur
de comparer plusieurs Artifacts ou versions avec préconditions et résultats partiels visibles sans fusionner les originaux ni qualifier automatiquement une différence.

## 3. Objectifs
- vérifier la relation et les préconditions des Artifacts
- comparer métadonnées, sections, fonctions, symboles, chaînes, graphes, structures, ressources et annotations
- classer ajouté, supprimé, modifié, déplacé ou non résolu
- enregistrer le résultat, relier des différences aux Hypotheses et préparer un Finding Draft

## 4. Non-objectifs
- fusionner ou modifier les Artifacts
- présenter une différence comme vulnérabilité ou Finding confirmé
- choisir un moteur de diff ou définir son algorithme

## 5. Propriétaire
Investigate possède le contexte analytique, les annotations et les relations au Case. Studio conserve les Tools, Tool Calls et Automation Runs; Platform Settings administre les environnements; Govern conserve l’autorité sur toute cible réelle; Shared conserve les mécanismes transversaux.

## 6. Utilisateurs
Reverse Engineer principal; Malware Analyst, Reviewer et Vulnerability Researcher viewer secondaires.

## 7. Conditions d’entrée
- au moins deux Artifacts accessibles
- relation ou absence de relation explicitement évaluée
- Tools/versions et résultats disponibles ou partiels

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Artifact set | Artifact Management | versions et lineage | Oui | versions sélectionnées | comparaison bloquée si moins de deux |
| Analysis results | Static/Reverse capabilities | sections, functions, graphs and knowledge | Oui | pour chaque version | résultat partiel avec préconditions |
| Comparison Tool | CMDR Studio | capabilities et version | Oui | courante | manual/other deterministic comparison only |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | contexte, Hypothesis et return origin | lecture et liaison uniquement |
| Artifact / Derived Artifact / Runtime Artifact | Investigate | source, versions, restrictions et provenance | lecture; aucun original modifié |
| Tool / Tool Call / Automation Run | CMDR Studio | outil, version, exécution et attribution | lecture et sélection; lifecycle Studio |
| Artifacts and versions | Artifact Management | sources, lineage and restrictions | lecture |
| Static/reverse results | CAP-INV-305..337 | metadata, sections, functions, symbols, graphs, structures, annotations | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Trace / Activity event | Shared | émission de l’action et de son résultat | append-only; aucune suppression locale |
| Binary comparison result | Investigate concept | create/review/supersede | n’altère ni ne fusionne les originaux |
| Hypothesis/Finding Draft relation | Investigate | link/prepare | différence n’est pas conclusion |

## 11. Fonctionnalités
- sélectionner et vérifier plusieurs Artifacts
- comparer chaque catégorie avec sources et Tools visibles
- distinguer added/removed/modified/moved/unresolved
- afficher limites et résultats partiels par version
- enregistrer observations et préparer Hypothesis/Finding Draft sans confirmation

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Configurer la comparaison | Reverse Engineer | Comparison request | 2 | Artifacts accessibles | scope versionné | OPEN-013 |
| Lancer le traitement isolé | Analyst | Comparison job | 1 | Tool autorisé | résultat progressif | Non |
| Inspecter/filtrer différences | Reviewer | Comparison result | 0 | résultat disponible | différences attribuées | Non |
| Relier à Hypothesis/Finding Draft | Investigation Lead | Analytical relation | 2 | justification fournie | draft uniquement | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Inspecter les différences binaires | Oui | Oui | Oui | Oui | viewer déterministe et navigation manuelle de les différences binaires |
| Proposer un regroupement ou une explication des différences | Oui | Oui | Oui | Oui | règles, heuristiques visibles, comparaison et revue humaine |
| Préparer un handoff | Oui | Oui | Oui | Oui | sélection manuelle, checklist et validation humaine |

Toute sortie automatisée expose initiateur, producteur et version, Automation Run et Tool Calls lorsqu’ils existent, sources, paramètres, timestamp, statut, incertitude, owner humain et disposition acceptée, modifiée ou rejetée.

## 14. États fonctionnels
- draft
- validating
- queued
- running
- partial
- available
- failed
- incompatible
- superseded

Ces états sont fonctionnels et ne constituent pas une machine d’état objet définitive.

## 15. États d’interface
- **Loading** conserve le Workbench, l’Artifact, la sélection et le return origin.
- **Empty** explique l’absence de résultat sans simuler une analyse.
- **Partial** identifie les sources, vues ou événements manquants et leurs conséquences.
- **Error** conserve les résultats valides, l’erreur et une reprise sûre.
- **Offline** limite les mutations et affiche la dernière synchronisation.
- **Permission denied** ne révèle aucune donnée protégée.
- **Stale** distingue la dernière observation connue de l’état courant.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Binary comparison result | Analysis Result concept | Reverse Workbench/Case | préconditions, Tools, versions et partials visibles |
| Difference selections | Comparison observations | CAP-INV-337/346 | catégorie, sources et incertitude |
| Finding Draft input | Draft relation | CAP-INV-109 | aucune vulnérabilité confirmée automatiquement |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Reverse Session | sélectionner versions | Binary Diff | Artifacts, lineage, objectifs, restrictions | retour session |
| Difference | ouvrir source | Disassembly/Decompilation/Functions | Artifact/version, location/function, selection | retour diff |
| Comparison | préparer handoff | CAP-INV-346 | differences, partials, contradictions, provenance | retour diff |

Chaque transition conserve tenant, environnement, Case, Artifact, sélection, permissions et return origin. Une erreur ne détruit pas la source ni les résultats déjà valides.

## 18. Dépendances
- CAP-INV-310 Multi-Artifact Comparison
- CAP-INV-305..337
- Shared Background Jobs/Diff Viewer
- CAP-INV-103/109
- OPEN-005/013/015

## 19. Source de vérité
Chaque Artifact et résultat conserve son owner et sa version; le comparison result appartient au contexte Investigate; l’exécution Tool reste Studio.

## 20. Provenance et audit
Tracer Artifacts/versions, relation supposée, préconditions, Tool/version, catégories, limites, résultats partiels, sélections et handoffs.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Step-up éventuel | Séparation des tâches | Owner | Phase propriétaire |
|---|---|---:|---|---|---|---|
| binary comparison | exposition et traitement de plusieurs Artifacts | 0 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| comparison request/update | scope réversible | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| comparison processing | traitement isolé | 1 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| Finding Draft prepare | risque de conclusion prématurée | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |

La matrice atomique, les namespaces et le modèle RBAC/ABAC restent reportés à la phase Permissions.

## 22. Limites et erreurs
- relation entre versions inconnue visible
- résultat partiel pour une version signalé
- incompatibilité Tool explicite
- originaux immuables
- différence non vulnérabilité

## 23. Métriques
- comparaisons available/partial/failed
- catégories unresolved
- préconditions différentes
- handoffs préparés et rejetés

## 24. Classification de livraison
Delivery status `defined`; delivery mode `planned`. La promotion exige objets et permissions approuvés, Tools et environnements évalués, contrats techniques, tests de sécurité, écrans et release evidence. Aucun moteur, debugger, produit tiers ou implémentation n’est choisi.

## 25. Critères d’acceptation
### Scénario 1
**Given** deux versions avec résultat partiel pour l’une
**When** l’analyste compare
**Then** préconditions différentes et partial visible; originaux inchangés; aucune vulnérabilité confirmée

### Scénario 2
**Given** deux Artifacts non liés
**When** la comparaison est préparée
**Then** la relation inconnue est visible et exige confirmation du scope

### Scénario 3
**Given** aucune IA
**When** la comparaison est lancée
**Then** comparateurs déterministes, filtres et revue humaine restent disponibles

## 26. Questions ouvertes
- OPEN-005 reste ouverte; aucune décision n’est fermée par cette spécification.
- OPEN-013 reste ouverte; aucune décision n’est fermée par cette spécification.
- OPEN-015 reste ouverte; aucune décision n’est fermée par cette spécification.

Les schémas, cardinalités, machines d’état, formats d’adresse et permissions atomiques sont reportés aux phases propriétaires.

## 27. Consommateurs documentaires
- INV-REV-001
- Diff Viewer
- CAP-INV-310/337/346
- Case Workspace and future detection handoff
