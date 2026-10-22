---
id: CAP-INV-307
title: Static Binary Analysis
product: investigate
module: analysis-workbench
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-014
  - REQ-INV-002
  - REQ-PROD-020
open_decisions:
  - OPEN-005
source-of-truth: canonical
---
# CAP-INV-307 — Static Binary Analysis

## 1. Définition
Analyser statiquement un binaire en inspectant identité, en-têtes, sections, imports, exports, dépendances, ressources, chaînes, symboles, signatures et caractéristiques structurelles sans désassembler, décompiler ou exécuter.

## 2. Problème utilisateur
Un binaire peut présenter des incohérences structurelles, dépendances inhabituelles ou ressources intégrées importantes. Sans vue consolidée et sourcée, l’analyste ne peut pas préparer correctement les analyses suivantes.

## 3. Objectifs
- Inspecter l’identité et les structures binaires observables.
- Conserver la provenance et le lien au Case.
- Produire des sorties inspectables et réutilisables.
- Préparer un handoff explicite sans qualification automatique.

## 4. Non-objectifs
- Ne pas exécuter l’Artifact.
- Ne pas choisir de moteur, bibliothèque, API, protocole ou format interne.
- Ne pas confirmer automatiquement Evidence ou Finding.

## 5. Propriétaire
Investigate possède le contexte et l’interprétation; Studio, Settings, Govern et Shared conservent leurs objets.

## 6. Utilisateurs
- Malware Analyst
- Case Analyst
- Evidence Reviewer
- Investigation Lead

## 7. Conditions d’entrée
- Case et Artifact accessibles.
- Version et restrictions visibles.
- Permission fonctionnelle réévaluée.
- Session existante ou création autorisée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Artifact source | Investigate / CAP-INV-105 | entrée analysée | oui | version immuable ou référencée | bloquer ou marquer partial |
| Case et objectif | Investigate | contexte analytique | oui | état courant | rester draft |
| Contexte Static Binary Analysis | analyste / résultat précédent | paramètres fonctionnels | oui | snapshot de session | demander complétude |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | scope, Hypothesis et return origin | consulter |
| Artifact | Investigate | source, version, type et relations | consulter |
| Binary structure result | Tool Call | headers, sections, imports/exports | consulter |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Analysis Session | concept Investigate | lier résultat et disposition | aucun schéma final |
| Binary observation | Investigate | créer/annoter | observed ou inferred explicitement |
| Derived Artifact | Investigate concept | créer | parent, Tool et transformation obligatoires |

## 11. Fonctionnalités
- Inspecter l’identité et les structures binaires observables.
- Comparer imports, exports, dépendances, ressources, chaînes et symboles disponibles.
- Conserver signatures et incohérences sans verdict automatique.
- Préparer des transitions futures vers Dynamic Sandbox ou Reverse Engineering.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Inspecter headers/sections | Malware Analyst | Binary result | 0 | contexte et permission valides | résultat visible | non |
| Comparer versions | Malware Analyst | Artifacts | 0 | contexte et permission valides | diff sourcé | non |
| Extraire ressource | Malware Analyst | Derived Artifact | 1 | contexte et permission valides | dérivé créé | non |
| Annoter incohérence | Malware Analyst | Observation | 2 | contexte et permission valides | annotation attribuée | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| Inspecter l’identité et les structures | oui | règles et Tools déterministes | oui | suggestion expliquée | inspection manuelle |
| Résumer Static Binary Analysis | oui | agrégation sourcée | oui | résumé attribué | viewer et filtres |
| Proposer une prochaine étape | oui | checklists/profils | oui | proposition modifiable | catalogue manuel |
| Qualifier Evidence ou Finding | humain | contrôles seulement | workflow de revue | jamais autonome | CAP-INV-107/108/109 |

La provenance automatisée et la disposition humaine restent visibles; aucune fonction essentielle ne dépend de l’IA.

## 14. États fonctionnels
`ready`, `analyzing`, `partial`, `completed`, `unsupported`, `corrupted`, `signature-unknown`, `restricted`. Machines finales reportées.

## 15. États d’interface
Loading conserve le contexte; Empty explique; Partial nomme les lacunes; Error garde les résultats valides; Offline est stale/read-only; Permission denied ne fuit rien.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Static binary result | Analysis Result | session/Case | dimensions et erreurs visibles |
| Binary observations | observations | Hypothesis/Case | aucun score transformé en Finding |
| Extracted resources | Derived Artifacts | CAP-INV-311 | lineage complet |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-305 | format binaire reconnu | CAP-INV-307 | Artifact, structures, restrictions | structure |
| CAP-INV-307 | extraire ressource | CAP-INV-311 | parent, ressource, Tool/version | binary analysis |
| CAP-INV-307 | analyse future | future Dynamic/Reverse | Artifact, dérivés, objectifs, observations, provenance | binary analysis |

Le return origin, les permissions, versions et résultats partiels sont conservés.

## 18. Dépendances
- CAP-INV-105 Artifact Management
- CAP-INV-302 Analysis Session Management
- CAP-INV-312 Analysis Provenance and Reproducibility
- Shared Trace/Activity/Object Linking/Export
- CAP-INV-305/306/311/312
- future Dynamic Sandbox

## 19. Source de vérité
Investigate est source de l’interprétation et des relations analytiques. Les objets sources gardent leur owner ; Studio reste source des Tool/Tool Call et Shared des mécanismes transversaux.

## 20. Provenance et audit
Enregistrer Case, session, Artifact/version, paramètres propres à Static Binary Analysis, Tool/version, Tool Calls, outputs, erreurs, annotations, disposition humaine, timestamp et correlation ID.

## 21. Permissions fonctionnelles
- Artifact read
- Analysis Session read/update
- Tool use/output read
- analysis result annotate
- static binary analysis
- binary raw read
- resource extraction
- future workspace handoff

Matrice atomique reportée.

## 22. Limites et erreurs
- Artifact inaccessible, superseded ou restreint.
- Tool indisponible, incompatible ou résultat partiel.
- Contexte tenant/environnement incohérent.
- Symboles absents ou stripping.
- Signature non vérifiable ou métadonnées conflictuelles.
- Format binaire partiellement supporté.

## 23. Métriques
- Utilisations de Static Binary Analysis.
- Résultats partial/failed.
- Temps jusqu’à annotation ou handoff.
- Sorties avec provenance complète.

## 24. Classification de livraison
`defined` / `planned`; aucune preuve d’implémentation, moteur, hyperviseur, API, protocole ou commande.

## 25. Critères d’acceptation
**Given** un Artifact compatible et une session autorisée
**When** l’analyste utilise Static Binary Analysis
**Then** l’identité, les structures, sources, paramètres, versions et limites sont visibles sans exécuter le contenu

**Given** une entrée manquante ou une permission refusée
**When** l’utilisateur demande Static Binary Analysis
**Then** l’erreur est explicite, les résultats valides sont conservés et aucune conclusion n’est inventée

**Given** aucun modèle IA
**When** l’analyste réalise Static Binary Analysis
**Then** les viewers, règles, Tools déterministes, filtres et revue humaine couvrent le workflow essentiel

## 26. Questions ouvertes
- Le schéma final des résultats de Static Binary Analysis est reporté à la phase Objets.
- OPEN-005 reste ouverte.

## 27. Consommateurs documentaires
- Analysis Workbench et capability map
- Case Workspace, Evidence Board et écrans techniques
- phases Objets, Permissions, Journeys et Technique
