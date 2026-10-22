---
id: CAP-INV-336
title: Data Types, Structures and Memory Layout Analysis
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
  - REQ-OBJ-009
  - REQ-AI-002
open_decisions:
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---

# CAP-INV-336 — Data Types, Structures and Memory Layout Analysis

## 1. Définition
Data Types, Structures and Memory Layout Analysis définit le comportement produit permettant de créer et réviser des interprétations analytiques de types, structures et layouts sans les présenter comme certaines ni altérer l’Artifact.

## 2. Problème utilisateur
de créer et réviser des interprétations analytiques de types, structures et layouts sans les présenter comme certaines ni altérer l’Artifact.

## 3. Objectifs
- afficher types détectés/proposés, tailles, membres, variantes, tableaux et références
- explorer alignements et relations entre structures
- créer, modifier, versionner et comparer une définition analytique
- appliquer localement une interprétation et revenir à une version précédente
- signaler ambiguïtés et conflits

## 4. Non-objectifs
- définir une représentation mémoire interne définitive
- présenter une structure proposée comme correcte automatiquement
- modifier les octets de l’Artifact ou commencer Memory Forensics

## 5. Propriétaire
Investigate possède le contexte analytique, les annotations et les relations au Case. Studio conserve les Tools, Tool Calls et Automation Runs; Platform Settings administre les environnements; Govern conserve l’autorité sur toute cible réelle; Shared conserve les mécanismes transversaux.

## 6. Utilisateurs
Reverse Engineer ou Type Analyst principal; Malware Analyst et Reviewer secondaires.

## 7. Conditions d’entrée
- Reverse Analysis Session active
- usages ou données suffisants ou lacune visible
- permission de connaissance analytique

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Usages code/données | CAP-INV-332/333/334 | accès, variables et références | Oui | session active | définition reste proposed |
| Type/structure proposals | Tool/analyst | noms, tailles, membres, alignements | Non | version sélectionnée | création manuelle possible |
| Runtime snapshots | CAP-INV-341 | valeurs observées pour validation | Non | snapshot attribué | aucune confirmation runtime |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | contexte, Hypothesis et return origin | lecture et liaison uniquement |
| Artifact / Derived Artifact / Runtime Artifact | Investigate | source, versions, restrictions et provenance | lecture; aucun original modifié |
| Tool / Tool Call / Automation Run | CMDR Studio | outil, version, exécution et attribution | lecture et sélection; lifecycle Studio |
| Decompiler/disassembly usages | CAP-INV-332/333 | variables, offsets and access patterns | lecture |
| Type/Structure candidates | Tool results / analyst knowledge | definitions, sizes and relations | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Trace / Activity event | Shared | émission de l’action et de son résultat | append-only; aucune suppression locale |
| Type Definition / Structure Definition | Investigate concepts | create/update/version/supersede | analytical overlay only |
| Local application relation | Investigate | apply/revert | n’altère pas l’Artifact source |

## 11. Fonctionnalités
- voir types, tailles, structures, membres, variantes, tableaux et références
- comparer plusieurs interprétations et sources
- créer/modifier/versionner une définition analytique
- appliquer localement et revert une interprétation
- signaler ambigu/conflicting sans altérer la source

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspecter un type | Type Analyst | Type candidate | 0 | source disponible | usages et confiance visibles | Non |
| Créer/modifier une définition | Type Analyst | Type/Structure Definition | 2 | session modifiable | nouvelle version analytique | OPEN-013 |
| Appliquer localement | Reverse Engineer | Interpretation relation | 2 | compatibilité vérifiée | vue enrichie, source intacte | OPEN-013 |
| Revenir à une version | Reviewer | Definition version | 2 | historique disponible | version précédente active | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Inspecter les types et structures | Oui | Oui | Oui | Oui | viewer déterministe et navigation manuelle de les types et structures |
| Proposer un type, un membre ou un alignement | Oui | Oui | Oui | Oui | règles, heuristiques visibles, comparaison et revue humaine |
| Préparer un handoff | Oui | Oui | Oui | Oui | sélection manuelle, checklist et validation humaine |

Toute sortie automatisée expose initiateur, producteur et version, Automation Run et Tool Calls lorsqu’ils existent, sources, paramètres, timestamp, statut, incertitude, owner humain et disposition acceptée, modifiée ou rejetée.

## 14. États fonctionnels
- detected
- proposed
- applied-locally
- ambiguous
- conflicting
- superseded
- withdrawn

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
| Versioned type/structure definition | Analytical knowledge | Disassembly/Decompilation/Debugger | source, confiance et version visibles |
| Local interpretation mapping | Relation record | Reverse Session | réversible et limitée à la session |
| Conflict/ambiguity notice | Review event | Reviewer | aucune résolution silencieuse |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Decompilation/Disassembly | ouvrir type | Types/Structures | variables, usages, offsets, sources | retour code |
| Type definition | appliquer localement | Code views | version, scope, locations | retour definition |
| Runtime snapshot | comparer interprétation | Types/Structures | valeurs observées, timestamp, incertitude | retour debugger |

Chaque transition conserve tenant, environnement, Case, Artifact, sélection, permissions et return origin. Une erreur ne détruit pas la source ni les résultats déjà valides.

## 18. Dépendances
- CAP-INV-332/333/334
- CAP-INV-337 versioned knowledge
- CAP-INV-341 runtime context
- Shared Versioning
- OPEN-013/015

## 19. Source de vérité
Les propositions Tool sont attribuées; les définitions et mappings analytiques appartiennent à Investigate; aucun layout final ou objet canonique n’est créé.

## 20. Provenance et audit
Tracer auteur/agent, sources, Tool/version, membres et tailles proposés, versions, application locale, revert, conflits, confiance et disposition.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Step-up éventuel | Séparation des tâches | Owner | Phase propriétaire |
|---|---|---:|---|---|---|---|
| type definition read | contenu analytique | 0 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| type definition create/update | mutation réversible | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| local interpretation apply/revert | impact sur les vues | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| structure export | diffusion | 1 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |

La matrice atomique, les namespaces et le modèle RBAC/ABAC restent reportés à la phase Permissions.

## 22. Limites et erreurs
- taille inconnue visible
- alignement proposé non certain
- conflit de définition préservé
- runtime value n’est pas preuve
- aucune écriture dans l’Artifact

## 23. Métriques
- définitions proposed/applied/superseded
- reverts
- conflits et ambiguïtés
- suggestions acceptées/modifiées/rejetées

## 24. Classification de livraison
Delivery status `defined`; delivery mode `planned`. La promotion exige objets et permissions approuvés, Tools et environnements évalués, contrats techniques, tests de sécurité, écrans et release evidence. Aucun moteur, debugger, produit tiers ou implémentation n’est choisi.

## 25. Critères d’acceptation
### Scénario 1
**Given** une structure proposée par un Tool
**When** l’analyste l’examine
**Then** source, version et incertitude sont visibles; aucune application automatique

### Scénario 2
**Given** une interprétation appliquée localement
**When** l’analyste revert
**Then** les vues reviennent à la version précédente et l’Artifact reste inchangé

### Scénario 3
**Given** aucun modèle IA
**When** une définition est créée
**Then** éditeur analytique, usages, comparaison et versioning restent disponibles

## 26. Questions ouvertes
- OPEN-013 reste ouverte; aucune décision n’est fermée par cette spécification.
- OPEN-015 reste ouverte; aucune décision n’est fermée par cette spécification.

Les schémas, cardinalités, machines d’état, formats d’adresse et permissions atomiques sont reportés aux phases propriétaires.

## 27. Consommateurs documentaires
- INV-REV-001
- Types view
- Disassembly/Decompilation/Debugger
- future Type/Structure object specifications
