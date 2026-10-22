---
id: CAP-INV-342
title: Memory, Modules and Loaded Image Inspection
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
  - REQ-PROD-052
  - REQ-OBJ-003
  - REQ-SEC-001
open_decisions:
  - OPEN-005
  - OPEN-013
  - OPEN-014
  - OPEN-015
source-of-truth: canonical
---

# CAP-INV-342 — Memory, Modules and Loaded Image Inspection

## 1. Définition
Memory, Modules and Loaded Image Inspection définit le comportement produit permettant d’inspecter dans une Debugger Session les régions mémoire, modules et images chargées au niveau produit sans devenir une capability de Memory Forensics ni fournir des méthodes d’injection ou d’évasion.

## 2. Problème utilisateur
d’inspecter dans une Debugger Session les régions mémoire, modules et images chargées au niveau produit sans devenir une capability de Memory Forensics ni fournir des méthodes d’injection ou d’évasion.

## 3. Objectifs
- afficher régions, permissions déclarées, relations, modules, images, bases/emplacements, versions et symboles
- rendre zones non mappées et changements observés visibles
- naviguer, annoter et comparer deux snapshots
- extraire un Derived ou Runtime Artifact selon permission et lineage

## 4. Non-objectifs
- analyser une image mémoire complète
- inspecter un Endpoint réel
- définir un moteur, des méthodes d’injection, d’évasion ou de bypass
- présenter un module chargé comme malveillant

## 5. Propriétaire
Investigate possède le contexte analytique, les annotations et les relations au Case. Studio conserve les Tools, Tool Calls et Automation Runs; Platform Settings administre les environnements; Govern conserve l’autorité sur toute cible réelle; Shared conserve les mécanismes transversaux.

## 6. Utilisateurs
Debug Analyst principal; Reverse Engineer et Reviewer secondaires.

## 7. Conditions d’entrée
- Debugger Session active/paused
- snapshot ou state available
- permissions de lecture/extraction
- environnement isolé

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Runtime Snapshot | CAP-INV-341 | state, timestamp and session | Oui | snapshot selected | live projection with stale warning |
| Memory/module result | Debugger Tool Call | regions, permissions, modules, images, symbols | Oui | Tool/version recorded | partial/failed visible |
| Artifact policy | Artifact Management/Settings | extraction and retention restrictions | Oui pour extraction | current | extraction disabled |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | contexte, Hypothesis et return origin | lecture et liaison uniquement |
| Artifact / Derived Artifact / Runtime Artifact | Investigate | source, versions, restrictions et provenance | lecture; aucun original modifié |
| Tool / Tool Call / Automation Run | CMDR Studio | outil, version, exécution et attribution | lecture et sélection; lifecycle Studio |
| Runtime State Snapshot | CAP-INV-341 | session, timestamp, thread and state | read |
| Memory region/module projection | Debugger Tool Call | regions, declared permissions, loaded images and symbols | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Trace / Activity event | Shared | émission de l’action et de son résultat | append-only; aucune suppression locale |
| Memory/Module observation | Investigate concept | annotate/link/compare | debugger-scoped, not Memory Forensics |
| Derived/Runtime Artifact request | Investigate | prepare/create via CAP-INV-311/324 | source and lineage required |

## 11. Fonctionnalités
- show debugger-scoped memory regions and declared permissions
- show loaded modules/images, versions, bases and symbols
- show unmapped regions and observed changes
- navigate and compare snapshots
- prepare Derived/Runtime Artifact extraction with lineage
- keep explicit Memory Forensics boundary

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspecter région/module | Debug Analyst | Memory/Module projection | 0 | state readable | attributed view | Non |
| Comparer snapshots | Reviewer | Memory comparison | 0 | compatible snapshots | changes and gaps visible | Non |
| Annoter observation | Analyst | Memory/Module observation | 2 | session modifiable | versioned note | OPEN-013 |
| Extraire Artifact | Authorized Analyst | Derived/Runtime Artifact | 1 | policy and permission allow | lineage preserved | Non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Inspecter les régions mémoire et modules | Oui | Oui | Oui | Oui | viewer déterministe et navigation manuelle de les régions mémoire et modules |
| Proposer une région ou un module pertinent | Oui | Oui | Oui | Oui | règles, heuristiques visibles, comparaison et revue humaine |
| Préparer un handoff | Oui | Oui | Oui | Oui | sélection manuelle, checklist et validation humaine |

Toute sortie automatisée expose initiateur, producteur et version, Automation Run et Tool Calls lorsqu’ils existent, sources, paramètres, timestamp, statut, incertitude, owner humain et disposition acceptée, modifiée ou rejetée.

## 14. États fonctionnels
- ready
- loading
- partial
- available
- stale
- unmapped
- restricted
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
| Memory/module projection | Debugger result | Debugger Workbench | debugger-scoped and attributed |
| Module/region observation | Observation concept | CAP-INV-343/346 | no malicious qualification |
| Derived/Runtime Artifact | Artifact lineage | CAP-INV-311/324 | source snapshot/session and method preserved |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Runtime Snapshot | open memory/modules | Memory and Modules | session, snapshot, thread, timestamp | return runtime |
| Region/module | navigate code | Disassembly/Decompilation | Artifact/image, location, symbol, selection | return memory view |
| Extraction selection | explicit request | Derived/Runtime Artifact management | source, scope, restrictions, provenance | return debugger |

Chaque transition conserve tenant, environnement, Case, Artifact, sélection, permissions et return origin. Une erreur ne détruit pas la source ni les résultats déjà valides.

## 18. Dépendances
- CAP-INV-311/324
- CAP-INV-341
- CAP-INV-343/346
- Platform Settings retention/storage policies
- Shared File Preview/Export
- Phase 4B.2B.3 boundary

## 19. Source de vérité
Debugger result belongs to Investigate context; extracted Artifact lifecycle follows CAP-INV-311/324; Settings owns retention/storage policies; no forensic object is created.

## 20. Provenance et audit
Trace session, snapshot, Tool/version, region/module identifiers as displayed, declared permissions, selection, comparison, extraction method, restrictions and lineage.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Step-up éventuel | Séparation des tâches | Owner | Phase propriétaire |
|---|---|---:|---|---|---|---|
| memory region/module read | highly sensitive runtime data | 0 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| snapshot comparison | sensitive comparison | 0 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| Derived/Runtime Artifact create/export | capture and diffusion | 1 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| memory annotation | reversible knowledge | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |

La matrice atomique, les namespaces et le modèle RBAC/ABAC restent reportés à la phase Permissions.

## 22. Limites et erreurs
- module loaded not malicious
- memory region view not full forensics
- unmapped/partial regions explicit
- restricted content redacted
- no injection/evasion guidance

## 23. Métriques
- regions/modules viewed
- partial/unmapped rates
- snapshot comparisons
- extractions with complete lineage

## 24. Classification de livraison
Delivery status `defined`; delivery mode `planned`. La promotion exige objets et permissions approuvés, Tools et environnements évalués, contrats techniques, tests de sécurité, écrans et release evidence. Aucun moteur, debugger, produit tiers ou implémentation n’est choisi.

## 25. Critères d’acceptation
### Scénario 1
**Given** un module est chargé dans la session
**When** l’analyste l’inspecte
**Then** version, source and symbols may be shown but no malicious qualification

### Scénario 2
**Given** une zone est non mappée
**When** l’utilisateur navigue
**Then** unmapped is visible and no content is invented

### Scénario 3
**Given** une extraction est autorisée
**When** l’analyste la capture
**Then** Derived/Runtime Artifact has source session/snapshot/method and does not become Evidence

## 26. Questions ouvertes
- OPEN-005 reste ouverte; aucune décision n’est fermée par cette spécification.
- OPEN-013 reste ouverte; aucune décision n’est fermée par cette spécification.
- OPEN-014 reste ouverte; aucune décision n’est fermée par cette spécification.
- OPEN-015 reste ouverte; aucune décision n’est fermée par cette spécification.

Les schémas, cardinalités, machines d’état, formats d’adresse et permissions atomiques sont reportés aux phases propriétaires.

## 27. Consommateurs documentaires
- INV-DBG-001
- CAP-INV-311/324/341/343/346
- future Phase 4B.2B.3 remains separate
- future Memory Region/Module Observation concepts
