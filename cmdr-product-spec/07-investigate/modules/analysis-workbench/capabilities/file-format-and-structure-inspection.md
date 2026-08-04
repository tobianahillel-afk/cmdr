---
id: CAP-INV-305
title: File Format and Structure Inspection
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
# CAP-INV-305 — File Format and Structure Inspection

## 1. Définition
Identifier et inspecter fonctionnellement le format, les en-têtes, sections, segments, ressources, tables, dépendances et structures imbriquées d’un Artifact sans imposer de parseur.

## 2. Problème utilisateur
Un format mal déclaré, une structure incohérente ou imbriquée peut masquer du contenu important. Sans inspection structurée, l’analyste ne sait pas quelles parties sont fiables ou extractibles.

## 3. Objectifs
- Identifier format et structures principales.
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
| Contexte File Format and Structure Inspection | analyste / résultat précédent | paramètres fonctionnels | oui | snapshot de session | demander complétude |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | scope, Hypothesis et return origin | consulter |
| Artifact | Investigate | source, version, type et relations | consulter |
| Structure result | Tool Call / parser déterministe | arbre, offsets et erreurs | consulter |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Analysis Session | concept Investigate | lier résultat et disposition | aucun schéma final |
| Structure observation | Investigate | créer/annoter | observed distinct de conclusion |
| Derived Artifact | Investigate concept | proposer/créer | parent et transformation obligatoires |

## 11. Fonctionnalités
- Identifier format et structures principales.
- Comparer structure déclarée et observée.
- Inspecter headers, sections, segments, ressources, tables et dépendances.
- Extraire des éléments comme Derived Artifacts avec provenance.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Inspecter structure | Analyst | Artifact | 0 | contexte et permission valides | structure affichée | non |
| Filtrer sections/ressources | Analyst | Structure result | 0 | contexte et permission valides | sélection conservée | non |
| Extraire un élément | Analyst | Derived Artifact | 1 | contexte et permission valides | dérivé sourcé | non |
| Annoter une incohérence | Analyst | Observation | 2 | contexte et permission valides | annotation attribuée | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| Assister identifier format et structures principales. | oui | règles et Tools déterministes | oui | suggestion expliquée | inspection manuelle |
| Résumer File Format and Structure Inspection | oui | agrégation sourcée | oui | résumé attribué | viewer et filtres |
| Proposer une prochaine étape | oui | checklists/profils | oui | proposition modifiable | catalogue manuel |
| Qualifier Evidence ou Finding | humain | contrôles seulement | workflow de revue | jamais autonome | CAP-INV-107/108/109 |

La provenance automatisée et la disposition humaine restent visibles; aucune fonction essentielle ne dépend de l’IA.

## 14. États fonctionnels
`format-recognized`, `format-ambiguous`, `structure-partial`, `nested`, `corrupted`, `unsupported`, `restricted`. Machines finales reportées.

## 15. États d’interface
Loading conserve le contexte; Empty explique; Partial nomme les lacunes; Error garde les résultats valides; Offline est stale/read-only; Permission denied ne fuit rien.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Structure map | Analysis Result | Workbench/Inspector | source, offsets et erreurs visibles |
| Extracted element | Derived Artifact | CAP-INV-311 | parent/transformation/version |
| Structure inconsistency | Observation | Case/Hypothesis | aucun verdict automatique |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-304 | inspection approfondie | CAP-INV-305 | Artifact, types, métadonnées | preview |
| CAP-INV-305 | extraire élément | CAP-INV-311 | parent, structure, transformation, Tool | structure |
| CAP-INV-305 | analyser contenu | CAP-INV-306/307/308/309 | élément, type et provenance | structure |

Le return origin, les permissions, versions et résultats partiels sont conservés.

## 18. Dépendances
- CAP-INV-105 Artifact Management
- CAP-INV-302 Analysis Session Management
- CAP-INV-312 Analysis Provenance and Reproducibility
- Shared Trace/Activity/Object Linking/Export
- CAP-INV-311
- Studio Tool versioning

## 19. Source de vérité
Investigate est source de l’interprétation et des relations analytiques. Les objets sources gardent leur owner ; Studio reste source des Tool/Tool Call et Shared des mécanismes transversaux.

## 20. Provenance et audit
Enregistrer Case, session, Artifact/version, paramètres propres à File Format and Structure Inspection, Tool/version, Tool Calls, outputs, erreurs, annotations, disposition humaine, timestamp et correlation ID.

## 21. Permissions fonctionnelles
- Artifact read
- Analysis Session read/update
- Tool use/output read
- analysis result annotate
- raw content read
- embedded content extraction
- Derived Artifact create
- structure result export

Matrice atomique reportée.

## 22. Limites et erreurs
- Artifact inaccessible, superseded ou restreint.
- Tool indisponible, incompatible ou résultat partiel.
- Contexte tenant/environnement incohérent.
- Format polyglotte ou imbriqué.
- Structure tronquée, malformée ou chiffrée.
- Résultats divergents entre Tools.

## 23. Métriques
- Utilisations de File Format and Structure Inspection.
- Résultats partial/failed.
- Temps jusqu’à annotation ou handoff.
- Sorties avec provenance complète.

## 24. Classification de livraison
`defined` / `planned`; aucune preuve d’implémentation, moteur, hyperviseur, API, protocole ou commande.

## 25. Critères d’acceptation
**Given** un Artifact compatible et une session autorisée
**When** l’analyste utilise File Format and Structure Inspection
**Then** le format, les structures, les sources, paramètres, versions et limites sont visibles sans exécuter le contenu

**Given** une entrée manquante ou une permission refusée
**When** l’utilisateur demande File Format and Structure Inspection
**Then** l’erreur est explicite, les résultats valides sont conservés et aucune conclusion n’est inventée

**Given** aucun modèle IA
**When** l’analyste réalise File Format and Structure Inspection
**Then** les viewers, règles, Tools déterministes, filtres et revue humaine couvrent le workflow essentiel

## 26. Questions ouvertes
- Le schéma final des résultats de File Format and Structure Inspection est reporté à la phase Objets.
- OPEN-005 reste ouverte.

## 27. Consommateurs documentaires
- Analysis Workbench et capability map
- Case Workspace, Evidence Board et écrans techniques
- phases Objets, Permissions, Journeys et Technique
