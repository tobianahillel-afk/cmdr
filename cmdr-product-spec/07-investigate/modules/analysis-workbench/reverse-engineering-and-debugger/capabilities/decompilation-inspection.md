---
id: CAP-INV-333
title: Decompilation Inspection
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
  - REQ-PROD-052
  - REQ-AI-002
  - REQ-UX-002
open_decisions:
  - OPEN-005
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---

# CAP-INV-333 — Decompilation Inspection

## 1. Définition
Decompilation Inspection définit le comportement produit permettant d’inspecter une représentation décompilée attribuée, partielle et réconciliable avec le désassemblage sans la présenter comme code source original.

## 2. Problème utilisateur
d’inspecter une représentation décompilée attribuée, partielle et réconciliable avec le désassemblage sans la présenter comme code source original.

## 3. Objectifs
- afficher fonctions, variables, types, structures, appels, conditions et boucles proposés
- synchroniser désassemblage et décompilation
- rendre erreurs, ambiguïtés et zones non reconstruites visibles
- comparer plusieurs résultats, annoter et corriger localement une interprétation

## 4. Non-objectifs
- produire ou restaurer le code source original
- masquer le désassemblage sous-jacent
- choisir un décompilateur ou définir son algorithme
- compiler, patcher ou déployer du code

## 5. Propriétaire
Investigate possède le contexte analytique, les annotations et les relations au Case. Studio conserve les Tools, Tool Calls et Automation Runs; Platform Settings administre les environnements; Govern conserve l’autorité sur toute cible réelle; Shared conserve les mécanismes transversaux.

## 6. Utilisateurs
Reverse Engineer principal; Malware Analyst, Reviewer et Type Analyst secondaires.

## 7. Conditions d’entrée
- Reverse Analysis Session active
- Disassembly ou emplacement disponible
- Decompiler Tool/version attribué
- Artifact non modifié

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Disassembly/function context | CAP-INV-332/334 | emplacement, fonction candidate et xrefs | Oui | session active | décompilation non ancrée signalée |
| Decompilation result | Studio Tool Call | représentation, statut et diagnostics | Oui | Tool/version enregistrés | failed/partial/incompatible visible |
| Types/structures | CAP-INV-336 | interprétations analytiques | Non | version sélectionnée | types inconnus restent explicites |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | contexte, Hypothesis et return origin | lecture et liaison uniquement |
| Artifact / Derived Artifact / Runtime Artifact | Investigate | source, versions, restrictions et provenance | lecture; aucun original modifié |
| Tool / Tool Call / Automation Run | CMDR Studio | outil, version, exécution et attribution | lecture et sélection; lifecycle Studio |
| Disassembly result | CAP-INV-332 | ancrage instruction/emplacement | lecture |
| Decompiler result | CMDR Studio Tool Call | pseudo-représentation, erreurs et ambiguïtés | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Trace / Activity event | Shared | émission de l’action et de son résultat | append-only; aucune suppression locale |
| Local interpretation correction | Investigate | create/update/version | corrige l’interprétation, jamais l’Artifact |
| Decompilation annotation | Investigate | create/update | reste distincte du résultat Tool |

## 11. Fonctionnalités
- inspecter la représentation et sa provenance
- synchroniser vue avec désassemblage et historique
- afficher conditions, boucles, appels et zones ambiguës
- comparer plusieurs résultats Tools/versions
- annoter, renommer et corriger une interprétation locale réversible

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspecter la décompilation | Reverse Engineer | Decompiler result | 0 | résultat disponible/partial | représentation attribuée | Non |
| Basculer entre vues | Analyst | Location/function selection | 0 | ancrage disponible | sélection préservée | Non |
| Corriger une interprétation locale | Type Analyst | Interpretation overlay | 2 | session modifiable | overlay versionné, original intact | OPEN-013 |
| Comparer deux résultats | Reviewer | Decompiler comparison | 0 | Tools/versions visibles | différences et limites attribuées | Non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Inspecter la décompilation | Oui | Oui | Oui | Oui | viewer déterministe et navigation manuelle de la décompilation |
| Proposer un nom, un type ou une explication | Oui | Oui | Oui | Oui | règles, heuristiques visibles, comparaison et revue humaine |
| Préparer un handoff | Oui | Oui | Oui | Oui | sélection manuelle, checklist et validation humaine |

Toute sortie automatisée expose initiateur, producteur et version, Automation Run et Tool Calls lorsqu’ils existent, sources, paramètres, timestamp, statut, incertitude, owner humain et disposition acceptée, modifiée ou rejetée.

## 14. États fonctionnels
- queued
- generating
- partial
- available
- failed
- ambiguous
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
| Decompilation representation | Tool result projection | Reverse Workbench | jamais étiquetée code source original |
| Interpretation overlay | Analyst knowledge record | CAP-INV-336/337 | réversible, versionnée et attribuée |
| Comparison result | Analysis Result concept | Reviewer/CAP-INV-338 | Tools, versions et préconditions visibles |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Disassembly | ouvrir vue liée | Decompilation | Artifact, emplacement, fonction, sélection, annotations | retour même emplacement |
| Decompilation | naviguer vers instruction | Disassembly | emplacement source et correspondance incertaine | retour vue décompilée |
| Decompilation | préparer type/structure | CAP-INV-336 | variables, usages, hypothèse et provenance | retour fonction |

Chaque transition conserve tenant, environnement, Case, Artifact, sélection, permissions et return origin. Une erreur ne détruit pas la source ni les résultats déjà valides.

## 18. Dépendances
- CAP-INV-332
- CAP-INV-334
- CAP-INV-336
- CAP-INV-337
- CMDR Studio Tool/Tool Call
- OPEN-005/013/015

## 19. Source de vérité
Le résultat décompilé est une projection Tool attribuée; le désassemblage reste accessible; les overlays analytiques appartiennent à Investigate.

## 20. Provenance et audit
Tracer Tool/version, paramètres, Artifact/version, fonction/emplacement, diagnostics, ambiguïtés, corrections locales, annotations et disposition humaine.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Step-up éventuel | Séparation des tâches | Owner | Phase propriétaire |
|---|---|---:|---|---|---|---|
| decompilation view | contenu sensible | 0 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| interpretation correction | mutation réversible | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| variable/type rename | connaissance partagée | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| result export | risque de diffusion | 1 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |

La matrice atomique, les namespaces et le modèle RBAC/ABAC restent reportés à la phase Permissions.

## 22. Limites et erreurs
- ambiguïté toujours visible
- échec n’efface pas le désassemblage
- correspondance partielle non complétée silencieusement
- aucun code source original revendiqué
- aucun patch ou compilation

## 23. Métriques
- résultats available/partial/failed
- temps de bascule synchronisée
- corrections locales versionnées
- suggestions acceptées/modifiées/rejetées

## 24. Classification de livraison
Delivery status `defined`; delivery mode `planned`. La promotion exige objets et permissions approuvés, Tools et environnements évalués, contrats techniques, tests de sécurité, écrans et release evidence. Aucun moteur, debugger, produit tiers ou implémentation n’est choisi.

## 25. Critères d’acceptation
### Scénario 1
**Given** une décompilation partielle ou ambiguë
**When** l’analyste ouvre la vue
**Then** ambiguïté visible, désassemblage accessible, Tool/version affichés et annotation possible

### Scénario 2
**Given** une correspondance instruction-ligne absente
**When** l’utilisateur navigue
**Then** aucun emplacement n’est inventé et l’absence est signalée

### Scénario 3
**Given** aucun fournisseur IA
**When** la vue est utilisée
**Then** décompilation déterministe disponible selon Tool et toutes les opérations manuelles restent accessibles

## 26. Questions ouvertes
- OPEN-005 reste ouverte; aucune décision n’est fermée par cette spécification.
- OPEN-013 reste ouverte; aucune décision n’est fermée par cette spécification.
- OPEN-015 reste ouverte; aucune décision n’est fermée par cette spécification.

Les schémas, cardinalités, machines d’état, formats d’adresse et permissions atomiques sont reportés aux phases propriétaires.

## 27. Consommateurs documentaires
- INV-REV-001
- Code Viewer / Diff Viewer
- CAP-INV-332/334/336/337/338
- future decompiler Tool contracts
