---
id: CAP-INV-332
title: Disassembly Inspection
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
  - REQ-SEC-001
  - REQ-UX-002
open_decisions:
  - OPEN-005
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---

# CAP-INV-332 — Disassembly Inspection

## 1. Définition
Disassembly Inspection définit le comportement produit permettant d’inspecter une représentation désassemblée attribuée, navigable et incertaine sans la présenter comme code source ni fournir d’instructions offensives.

## 2. Problème utilisateur
d’inspecter une représentation désassemblée attribuée, navigable et incertaine sans la présenter comme code source ni fournir d’instructions offensives.

## 3. Objectifs
- afficher emplacements, instructions représentées, opérandes, branches, appels, données intégrées et zones non interprétées
- naviguer vers fonctions, références, symboles et données
- annoter ou renommer analytiquement sans modifier la source
- comparer des représentations et préparer une transition contrôlée vers Debugger

## 4. Non-objectifs
- choisir ou implémenter un moteur de désassemblage
- présenter le désassemblage comme code source original
- fournir des commandes, exploits, payloads ou techniques de contournement

## 5. Propriétaire
Investigate possède le contexte analytique, les annotations et les relations au Case. Studio conserve les Tools, Tool Calls et Automation Runs; Platform Settings administre les environnements; Govern conserve l’autorité sur toute cible réelle; Shared conserve les mécanismes transversaux.

## 6. Utilisateurs
Reverse Engineer principal; Malware Analyst et Reviewer secondaires.

## 7. Conditions d’entrée
- Reverse Analysis Session active
- Artifact compatible ou résultat partiel explicitement disponible
- Tool et version attribués
- permission de lecture du contenu

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Artifact et emplacement | CAP-INV-329/331 | source et Code Location | Oui | version active | vue vide ou partielle, aucune instruction inventée |
| Disassembly result | Tool Call Studio | représentation et statut | Oui | Tool/version enregistrés | état failed/incompatible |
| Symbols/xrefs | CAP-INV-334 | repères associés | Non | session courante | navigation par emplacement reste disponible |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | contexte, Hypothesis et return origin | lecture et liaison uniquement |
| Artifact / Derived Artifact / Runtime Artifact | Investigate | source, versions, restrictions et provenance | lecture; aucun original modifié |
| Tool / Tool Call / Automation Run | CMDR Studio | outil, version, exécution et attribution | lecture et sélection; lifecycle Studio |
| Code Location / Function / Symbol candidates | Investigate concepts | emplacements, fonctions, symboles et xrefs | lecture |
| Tool result | CMDR Studio | représentation désassemblée, version et statut | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Trace / Activity event | Shared | émission de l’action et de son résultat | append-only; aucune suppression locale |
| Disassembly annotation / rename relation | Investigate | create/update/version | ne modifie ni instructions ni Artifact |
| View comparison record | Investigate | create/read | conserve Tools et versions comparés |

## 11. Fonctionnalités
- inspecter instructions et opérandes représentées
- afficher branches, appels, données intégrées et zones non interprétées
- synchroniser sélection avec navigation, decompilation et graphes
- annoter, renommer et comparer sans altérer l’Artifact
- préparer contexte de breakpoint sans créer ni exécuter celui-ci

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspecter une représentation | Reverse Engineer | Disassembly result | 0 | résultat disponible ou partiel | vue attribuée et navigable | Non |
| Annoter ou renommer | Analyst | Knowledge record | 2 | session modifiable | interprétation versionnée | OPEN-013 |
| Comparer deux résultats | Reviewer | View comparison | 0 | Tools/versions visibles | différences et limites visibles | Non |
| Préparer le contexte Debugger | Reverse Engineer | Location handoff | 2 | environnement non encore lancé | proposition explicite, aucun Run | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Inspecter le désassemblage | Oui | Oui | Oui | Oui | viewer déterministe et navigation manuelle de le désassemblage |
| Proposer une fonction, un nom ou une explication | Oui | Oui | Oui | Oui | règles, heuristiques visibles, comparaison et revue humaine |
| Préparer un handoff | Oui | Oui | Oui | Oui | sélection manuelle, checklist et validation humaine |

Toute sortie automatisée expose initiateur, producteur et version, Automation Run et Tool Calls lorsqu’ils existent, sources, paramètres, timestamp, statut, incertitude, owner humain et disposition acceptée, modifiée ou rejetée.

## 14. États fonctionnels
- queued
- generating
- partial
- available
- failed
- incompatible
- ambiguous
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
| Disassembly representation | Tool result projection | Reverse Workbench | Tool/version, statut et incertitude visibles |
| Selected code locations | Location selection | CAP-INV-334/335/339 | source et système de repérage conservés |
| Analyst interpretation | Annotation/rename record | CAP-INV-337 | auteur, confiance et historique conservés |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Binary Navigation | ouvrir représentation | Disassembly | Artifact, emplacement, sélection, provenance | retour navigation |
| Disassembly | basculer vers décompilation | CAP-INV-333 | fonction/emplacement, xrefs, annotations, historique | retour même emplacement |
| Disassembly | préparer Debugger | CAP-INV-339 | Artifact/copie isolée, emplacements, objectif, restrictions | retour Reverse Session |

Chaque transition conserve tenant, environnement, Case, Artifact, sélection, permissions et return origin. Une erreur ne détruit pas la source ni les résultats déjà valides.

## 18. Dépendances
- CAP-INV-307 Static Binary Analysis
- CAP-INV-331 navigation
- CAP-INV-333 decompilation
- CAP-INV-334 functions/symbols/xrefs
- CMDR Studio Tool/Tool Call
- OPEN-005/013/015

## 19. Source de vérité
L’Artifact reste source; la représentation est un résultat Tool attribué; les annotations et renommages appartiennent à Investigate; Tool lifecycle reste Studio.

## 20. Provenance et audit
Tracer Artifact/version, Tool/version, paramètres, emplacements, statut, erreurs, sélection, annotations, renommages, comparaisons et acceptation/rejet de suggestions.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Step-up éventuel | Séparation des tâches | Owner | Phase propriétaire |
|---|---|---:|---|---|---|---|
| disassembly view | contenu sensible et propriété intellectuelle | 0 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| function annotation | mutation réversible | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| symbol rename | interprétation partagée | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| automated reverse request | exécution de Tool | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |

La matrice atomique, les namespaces et le modèle RBAC/ABAC restent reportés à la phase Permissions.

## 22. Limites et erreurs
- résultat partiel ne masque pas les zones absentes
- instruction non interprétée reste visible
- aucune commande réelle générée
- Tool indisponible n’efface pas les résultats antérieurs
- aucune transition vers cible réelle

## 23. Métriques
- couverture disponible/partielle
- erreurs et incompatibilités
- navigations vers xrefs/fonctions
- annotations acceptées/rejetées

## 24. Classification de livraison
Delivery status `defined`; delivery mode `planned`. La promotion exige objets et permissions approuvés, Tools et environnements évalués, contrats techniques, tests de sécurité, écrans et release evidence. Aucun moteur, debugger, produit tiers ou implémentation n’est choisi.

## 25. Critères d’acceptation
### Scénario 1
**Given** un résultat partiel
**When** l’analyste inspecte la zone absente
**Then** la lacune et le Tool sont visibles; aucune instruction n’est inventée

### Scénario 2
**Given** une suggestion de fonction IA
**When** l’analyste l’examine
**Then** elle reste candidate, attribuée et accept/reject possible

### Scénario 3
**Given** aucun modèle IA
**When** l’analyste utilise le désassemblage
**Then** navigation, xrefs, annotations et comparaison déterministes restent disponibles

## 26. Questions ouvertes
- OPEN-005 reste ouverte; aucune décision n’est fermée par cette spécification.
- OPEN-013 reste ouverte; aucune décision n’est fermée par cette spécification.
- OPEN-015 reste ouverte; aucune décision n’est fermée par cette spécification.

Les schémas, cardinalités, machines d’état, formats d’adresse et permissions atomiques sont reportés aux phases propriétaires.

## 27. Consommateurs documentaires
- INV-REV-001
- Technical Workbench Code Viewer
- CAP-INV-331/333/334/335/339
- future Tool and object contracts
