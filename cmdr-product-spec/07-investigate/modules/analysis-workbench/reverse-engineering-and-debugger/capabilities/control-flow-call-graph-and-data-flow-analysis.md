---
id: CAP-INV-335
title: Control Flow, Call Graph and Data Flow Analysis
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
  - REQ-UX-002
  - REQ-UX-006
open_decisions:
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---

# CAP-INV-335 — Control Flow, Call Graph and Data Flow Analysis

## 1. Définition
Control Flow, Call Graph and Data Flow Analysis définit le comportement produit permettant d’explorer control flow, relations d’appel et flux fonctionnels de données avec alternatives structurées, chemins non résolus et incertitude explicite.

## 2. Problème utilisateur
d’explorer control flow, relations d’appel et flux fonctionnels de données avec alternatives structurées, chemins non résolus et incertitude explicite.

## 3. Objectifs
- visualiser blocs, branches, chemins, boucles et appels
- explorer relations entre fonctions et comparer des graphes
- suivre sources et utilisations d’une information au niveau fonctionnel
- annoter, relier à une Hypothesis et préparer une analyse runtime
- fournir une alternative tabulaire aux graphes

## 4. Non-objectifs
- présenter un graphe statique comme ordre réel d’exécution
- masquer les données sous-jacentes
- définir les algorithmes d’analyse ou commencer Network/Memory Forensics

## 5. Propriétaire
Investigate possède le contexte analytique, les annotations et les relations au Case. Studio conserve les Tools, Tool Calls et Automation Runs; Platform Settings administre les environnements; Govern conserve l’autorité sur toute cible réelle; Shared conserve les mécanismes transversaux.

## 6. Utilisateurs
Reverse Engineer principal; Malware Analyst, Reviewer et Accessibility user secondaires.

## 7. Conditions d’entrée
- fonctions ou emplacements disponibles
- données de graphe attribuées ou absence visible
- alternative structurée disponible

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Functions and references | CAP-INV-334 | nodes, calls and data refs | Oui | session active | graphe partiel |
| Graph result | Tool/Shared Graph Engine | CFG/call/data-flow projection | Oui | Tool/version active | fallback tabulaire |
| Hypothesis | Investigate | question analytique | Non | Case courant | annotation sans relation |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | contexte, Hypothesis et return origin | lecture et liaison uniquement |
| Artifact / Derived Artifact / Runtime Artifact | Investigate | source, versions, restrictions et provenance | lecture; aucun original modifié |
| Tool / Tool Call / Automation Run | CMDR Studio | outil, version, exécution et attribution | lecture et sélection; lifecycle Studio |
| Function/Basic Block/Reference candidates | Investigate concepts | nœuds et relations statiques | lecture |
| Graph projection | Shared Graph Engine / Tool results | CFG, call graph et data-flow représentation | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Trace / Activity event | Shared | émission de l’action et de son résultat | append-only; aucune suppression locale |
| Graph annotation / path selection | Investigate | create/update/version | n’altère pas le graphe source |
| Hypothesis relation | Investigate | link/unlink | relation analytique, pas conclusion |

## 11. Fonctionnalités
- afficher CFG avec blocs, branches, boucles et chemins conditionnels
- afficher call graph avec relations non résolues
- suivre sources/utilisations d’information sans prétendre à une sémantique certaine
- comparer des graphes et annoter une sélection
- offrir une représentation tabulaire et préserver sélection/focus

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspecter un graphe | Reverse Engineer | Graph projection | 0 | données disponibles | vue et alternative structurée | Non |
| Filtrer ou comparer | Analyst | Graph selection | 0 | deux sources compatibles | différences et limites visibles | Non |
| Annoter un chemin | Analyst | Graph annotation | 2 | session modifiable | annotation versionnée | OPEN-013 |
| Relier à une Hypothesis | Reviewer | Hypothesis relation | 2 | Case accessible | relation explicite, aucune conclusion | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Inspecter les graphes et flux fonctionnels | Oui | Oui | Oui | Oui | viewer déterministe et navigation manuelle de les graphes et flux fonctionnels |
| Proposer un chemin, une relation ou une explication | Oui | Oui | Oui | Oui | règles, heuristiques visibles, comparaison et revue humaine |
| Préparer un handoff | Oui | Oui | Oui | Oui | sélection manuelle, checklist et validation humaine |

Toute sortie automatisée expose initiateur, producteur et version, Automation Run et Tool Calls lorsqu’ils existent, sources, paramètres, timestamp, statut, incertitude, owner humain et disposition acceptée, modifiée ou rejetée.

## 14. États fonctionnels
- ready
- computing
- partial
- available
- unresolved-paths
- ambiguous
- failed
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
| CFG/call/data-flow projection | Graph result | Reverse Workbench | statique, attribuée et non assimilée au runtime |
| Structured graph table | Accessible representation | all users | mêmes nœuds, relations et incertitudes |
| Selected path/context | Analysis selection | CAP-INV-339/346 | sources et limites conservées |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Function detail | ouvrir graphe | Control/Call/Data Flow | fonction, blocs, refs, sélection | retour fonction |
| Graph | sélectionner chemin | Debugger preparation | Artifact/copie, locations, conditions candidates | retour graph |
| Graph | préparer handoff | CAP-INV-346 | nœuds, relations, sources, contradictions | retour graph |

Chaque transition conserve tenant, environnement, Case, Artifact, sélection, permissions et return origin. Une erreur ne détruit pas la source ni les résultats déjà valides.

## 18. Dépendances
- CAP-INV-334
- Shared Graph Engine
- Technical Workbench Graph component
- CAP-INV-339/346
- OPEN-013/015

## 19. Source de vérité
Les nœuds et relations viennent des résultats Tool; Shared fournit le mécanisme de graphe; Investigate conserve sélection, annotations et liens Hypothesis.

## 20. Provenance et audit
Tracer Tool/version, graphe et version, filtres, nœuds/relations sélectionnés, chemins unresolved, annotations, comparaisons et handoffs.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Step-up éventuel | Séparation des tâches | Owner | Phase propriétaire |
|---|---|---:|---|---|---|---|
| control-flow/call/data-flow view | contenu technique sensible | 0 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| graph annotation | mutation réversible | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| Hypothesis link | relation analytique | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| graph export | diffusion | 1 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |

La matrice atomique, les namespaces et le modèle RBAC/ABAC restent reportés à la phase Permissions.

## 22. Limites et erreurs
- graphe statique jamais présenté comme runtime
- chemin non résolu visible
- graphe trop grand conserve progressive disclosure
- alternative tabulaire obligatoire
- données manquantes non inventées

## 23. Métriques
- graphes available/partial/failed
- chemins unresolved
- utilisation de l’alternative structurée
- comparaisons et annotations

## 24. Classification de livraison
Delivery status `defined`; delivery mode `planned`. La promotion exige objets et permissions approuvés, Tools et environnements évalués, contrats techniques, tests de sécurité, écrans et release evidence. Aucun moteur, debugger, produit tiers ou implémentation n’est choisi.

## 25. Critères d’acceptation
### Scénario 1
**Given** un call graph avec appels indirects non résolus
**When** l’analyste inspecte
**Then** les chemins non résolus sont visibles et aucune cible n’est inventée

### Scénario 2
**Given** un utilisateur active l’alternative structurée
**When** le graphe est consulté
**Then** nœuds, relations, sélection et incertitude restent accessibles

### Scénario 3
**Given** aucune IA
**When** l’analyse est effectuée
**Then** graphes déterministes, tables, filtres et annotations manuelles restent disponibles

## 26. Questions ouvertes
- OPEN-013 reste ouverte; aucune décision n’est fermée par cette spécification.
- OPEN-015 reste ouverte; aucune décision n’est fermée par cette spécification.

Les schémas, cardinalités, machines d’état, formats d’adresse et permissions atomiques sont reportés aux phases propriétaires.

## 27. Consommateurs documentaires
- INV-REV-001
- Graph and Tree components
- CAP-INV-334/339/346
- future graph concept specifications
