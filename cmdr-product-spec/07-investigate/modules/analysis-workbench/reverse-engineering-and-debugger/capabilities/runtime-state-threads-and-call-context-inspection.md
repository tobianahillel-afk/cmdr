---
id: CAP-INV-341
title: Runtime State, Threads and Call Context Inspection
product: investigate
module: analysis-workbench
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-05
requirement_ids:
  - REQ-INV-004
  - REQ-PROD-014
  - REQ-PROD-020
  - REQ-OBJ-009
  - REQ-AI-002
open_decisions:
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---

# CAP-INV-341 — Runtime State, Threads and Call Context Inspection

## 1. Définition
Runtime State, Threads and Call Context Inspection définit le comportement produit permettant d’inspecter threads, call stacks, frames, variables, valeurs et état processeur représenté sans compléter silencieusement les données manquantes ni les qualifier comme Evidence.

## 2. Problème utilisateur
d’inspecter threads, call stacks, frames, variables, valeurs et état processeur représenté sans compléter silencieusement les données manquantes ni les qualifier comme Evidence.

## 3. Objectifs
- afficher état runtime courant, thread actif et autres threads
- inspecter call stack, stack frames, paramètres, variables et valeurs disponibles
- naviguer vers disassembly/decompilation
- capturer et comparer des snapshots
- annoter une valeur incertaine et conserver la provenance

## 4. Non-objectifs
- effectuer une Memory Forensics complète
- présenter un frame ou une valeur runtime comme Evidence
- inventer des valeurs ou état processeur absents

## 5. Propriétaire
Investigate possède le contexte analytique, les annotations et les relations au Case. Studio conserve les Tools, Tool Calls et Automation Runs; Platform Settings administre les environnements; Govern conserve l’autorité sur toute cible réelle; Shared conserve les mécanismes transversaux.

## 6. Utilisateurs
Reverse Engineer et Debug Analyst principaux; Reviewer et Evidence Reviewer secondaires.

## 7. Conditions d’entrée
- Debugger Session paused/active with readable state
- permissions runtime read
- Tool result attributed
- sensitive values handled by policy

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Debugger Session/control stop | CAP-INV-339/340 | thread, location and lifecycle | Oui | current stop | stale if resumed |
| Runtime state result | Debugger Tool Call | threads, frames, values and represented processor state | Oui | timestamped | partial state |
| Code context | CAP-INV-332/333/334 | locations, functions and names | Non | session version | raw locations remain |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | contexte, Hypothesis et return origin | lecture et liaison uniquement |
| Artifact / Derived Artifact / Runtime Artifact | Investigate | source, versions, restrictions et provenance | lecture; aucun original modifié |
| Tool / Tool Call / Automation Run | CMDR Studio | outil, version, exécution et attribution | lecture et sélection; lifecycle Studio |
| Debugger Session / control event | CAP-INV-339/340 | current state and stop point | read |
| Runtime state projection | Debugger Tool Call | threads, frames, values and processor state where available | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Trace / Activity event | Shared | émission de l’action et de son résultat | append-only; aucune suppression locale |
| Runtime State Snapshot | Investigate concept | capture/read/compare/supersede | observation, not Evidence |
| Runtime annotation | Investigate | create/update | uncertainty explicit |

## 11. Fonctionnalités
- show current runtime state and active/other threads
- show call stack and frames with availability markers
- show parameters, variables, values and represented processor state when available
- navigate to code views preserving frame selection
- capture/compare snapshots and annotate uncertainty

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspecter runtime/thread/frame | Debug Analyst | Runtime projection | 0 | state readable | values and availability visible | Non |
| Capturer snapshot | Authorized Analyst | Runtime State Snapshot | 1 | policy allows | timestamped snapshot | Non |
| Annoter uncertainty | Analyst | Runtime annotation | 2 | snapshot accessible | attributed note | OPEN-013 |
| Comparer snapshots | Reviewer | Snapshot comparison | 0 | compatible snapshots | differences and staleness visible | Non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Inspecter l’état runtime et les call contexts | Oui | Oui | Oui | Oui | viewer déterministe et navigation manuelle de l’état runtime et les call contexts |
| Résumer une stack ou proposer un frame pertinent | Oui | Oui | Oui | Oui | règles, heuristiques visibles, comparaison et revue humaine |
| Préparer un handoff | Oui | Oui | Oui | Oui | sélection manuelle, checklist et validation humaine |

Toute sortie automatisée expose initiateur, producteur et version, Automation Run et Tool Calls lorsqu’ils existent, sources, paramètres, timestamp, statut, incertitude, owner humain et disposition acceptée, modifiée ou rejetée.

## 14. États fonctionnels
- capturing
- available
- partial
- stale
- invalid
- restricted
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
| Runtime State Snapshot | Snapshot concept | CAP-INV-342/343/345/346 | timestamp, session, thread and Tool attribution |
| Frame/thread selection | Context reference | Disassembly/Decompilation | location and selection preserved |
| Snapshot comparison | Analysis Result concept | Reviewer | missing/stale values explicit |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Execution Control | pause/reach event | Runtime State | session, thread, location, stack, timestamp | return control |
| Runtime frame | navigate code | Disassembly/Decompilation | Artifact/copy, location, function, frame selection | return snapshot |
| Runtime Snapshot | prepare handoff | CAP-INV-346 | selected values, frames, sources, uncertainty | return debugger |

Chaque transition conserve tenant, environnement, Case, Artifact, sélection, permissions et return origin. Une erreur ne détruit pas la source ni les résultats déjà valides.

## 18. Dépendances
- CAP-INV-339/340
- CAP-INV-332/333/334
- CAP-INV-342/343
- Shared Comparison/Trace
- OPEN-013/015

## 19. Source de vérité
Runtime values and snapshots are attributed Tool results in Investigate context; code knowledge remains Investigate; Tool lifecycle remains Studio.

## 20. Provenance et audit
Trace session, Tool/version, timestamp, stop reason, thread/frame IDs as displayed, value availability, snapshot creator, restrictions, comparisons and annotations.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Step-up éventuel | Séparation des tâches | Owner | Phase propriétaire |
|---|---|---:|---|---|---|---|
| runtime state/thread/call stack read | sensitive runtime data | 0 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| runtime snapshot create/export | capture and diffusion | 1 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| runtime annotation | reversible knowledge | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| processor state view | sensitive low-level state | 0 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |

La matrice atomique, les namespaces et le modèle RBAC/ABAC restent reportés à la phase Permissions.

## 22. Limites et erreurs
- state stale after resume
- missing values remain missing
- restricted values redacted
- frame-location mismatch visible
- snapshot not Evidence automatically

## 23. Métriques
- snapshots available/partial/stale
- thread/frame navigation success
- missing value rates
- snapshot comparison use

## 24. Classification de livraison
Delivery status `defined`; delivery mode `planned`. La promotion exige objets et permissions approuvés, Tools et environnements évalués, contrats techniques, tests de sécurité, écrans et release evidence. Aucun moteur, debugger, produit tiers ou implémentation n’est choisi.

## 25. Critères d’acceptation
### Scénario 1
**Given** un runtime state incomplet
**When** l’analyste inspecte
**Then** missing fields are visible and no silent completion occurs

### Scénario 2
**Given** un snapshot is selected for Evidence candidate
**When** handoff is prepared
**Then** snapshot remains distinct from Evidence and qualification is required

### Scénario 3
**Given** aucune IA
**When** threads/stacks are inspected
**Then** manual viewers, navigation, capture and comparison remain available

## 26. Questions ouvertes
- OPEN-013 reste ouverte; aucune décision n’est fermée par cette spécification.
- OPEN-015 reste ouverte; aucune décision n’est fermée par cette spécification.

Les schémas, cardinalités, machines d’état, formats d’adresse et permissions atomiques sont reportés aux phases propriétaires.

## 27. Consommateurs documentaires
- INV-DBG-001
- Disassembly/Decompilation views
- CAP-INV-342/343/345/346
- future Runtime Snapshot/Thread/Stack Frame concepts
