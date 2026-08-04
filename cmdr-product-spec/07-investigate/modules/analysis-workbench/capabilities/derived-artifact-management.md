---
id: CAP-INV-311
title: Derived Artifact Management
product: investigate
module: analysis-workbench
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-014
  - REQ-PROD-020
  - REQ-OBJ-003
open_decisions:
  - OPEN-013
  - OPEN-014
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-311 — Derived Artifact Management

## 1. Définition
Créer, enregistrer, inspecter et relier un Derived Artifact issu d’une extraction, décodage, décompression, conversion, normalisation, découpage, reconstruction ou transformation analytique, sans modifier le source.

## 2. Problème utilisateur
Sans objet dérivé explicite, les transformations deviennent silencieuses, les originaux peuvent être remplacés et la reproductibilité ou la custody deviennent impossibles.

## 3. Objectifs
- Préserver l’Artifact parent et la transformation.
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
| Contexte Derived Artifact Management | analyste / résultat précédent | paramètres fonctionnels | oui | snapshot de session | demander complétude |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | scope, Hypothesis et return origin | consulter |
| Artifact | Investigate | source, version, type et relations | consulter |
| Source Artifact | Investigate | version, restrictions et lineage | consulter |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Analysis Session | concept Investigate | lier résultat et disposition | aucun schéma final |
| Derived Artifact | concept Investigate | créer/versionner/supersede/withdraw | parent et transformation obligatoires |
| Parent-child relation | Investigate | créer | bidirectionnelle et sourcée |

## 11. Fonctionnalités
- Préserver l’Artifact parent et la transformation.
- Exposer auteur ou moteur, Tool/version, paramètres, timestamp et restrictions.
- Permettre descendants, inspection et historique.
- Empêcher toute promotion automatique vers Evidence.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Créer dérivé | Analyst/Tool | Derived Artifact | 1 | contexte et permission valides | dérivé et lineage | non |
| Annoter/reclassifier | Analyst | Derived Artifact | 2 | contexte et permission valides | nouvelle disposition | OPEN-013 |
| Supersede/withdraw | Reviewer | Derived Artifact | 2 | contexte et permission valides | historique conservé | OPEN-013 |
| Inspecter descendant | Analyst | Derived Artifact | 0 | contexte et permission valides | nouvel intake | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| Préserver l’Artifact parent et la transformation | oui | règles et Tools déterministes | oui | suggestion expliquée | inspection manuelle |
| Résumer Derived Artifact Management | oui | agrégation sourcée | oui | résumé attribué | viewer et filtres |
| Proposer une prochaine étape | oui | checklists/profils | oui | proposition modifiable | catalogue manuel |
| Qualifier Evidence ou Finding | humain | contrôles seulement | workflow de revue | jamais autonome | CAP-INV-107/108/109 |

La provenance automatisée et la disposition humaine restent visibles; aucune fonction essentielle ne dépend de l’IA.

## 14. États fonctionnels
`proposed`, `producing`, `available`, `partial`, `invalid`, `restricted`, `superseded`, `withdrawn-from-use`. Machines finales reportées.

## 15. États d’interface
Loading conserve le contexte; Empty explique; Partial nomme les lacunes; Error garde les résultats valides; Offline est stale/read-only; Permission denied ne fuit rien.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Derived Artifact | Artifact concept | CAP-INV-301/105 | parent, transformation, Tool et paramètres |
| Lineage relation | relation | Case/Evidence Review | historique navigable |
| Transformation event | business event | Trace/Activity | auteur, version et statut |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Analysis Result | produire dérivé | CAP-INV-311 | source, transformation, Tool/version, paramètres | source analysis |
| Derived Artifact | inspecter/analyser | CAP-INV-301 | parent, lineage, restrictions, Case | lineage |
| Derived Artifact | préparer Evidence | CAP-INV-313/107 | Artifact, source, transformation, provenance | lineage |

Le return origin, les permissions, versions et résultats partiels sont conservés.

## 18. Dépendances
- CAP-INV-105 Artifact Management
- CAP-INV-302 Analysis Session Management
- CAP-INV-312 Analysis Provenance and Reproducibility
- Shared Trace/Activity/Object Linking/Export
- CAP-INV-105/213/214
- Shared Versioning/Object Linking

## 19. Source de vérité
Investigate est source de l’interprétation et des relations analytiques. Les objets sources gardent leur owner ; Studio reste source des Tool/Tool Call et Shared des mécanismes transversaux.

## 20. Provenance et audit
Enregistrer Case, session, Artifact/version, paramètres propres à Derived Artifact Management, Tool/version, Tool Calls, outputs, erreurs, annotations, disposition humaine, timestamp et correlation ID.

## 21. Permissions fonctionnelles
- Artifact read
- Analysis Session read/update
- Tool use/output read
- analysis result annotate
- Derived Artifact create/read/export
- transformation parameters read
- lineage modify
- withdraw from use

Matrice atomique reportée.

## 22. Limites et erreurs
- Artifact inaccessible, superseded ou restreint.
- Tool indisponible, incompatible ou résultat partiel.
- Contexte tenant/environnement incohérent.
- Transformation échouée ou partielle.
- Parent inaccessible ou superseded.
- Tool/version manquant.

## 23. Métriques
- Utilisations de Derived Artifact Management.
- Résultats partial/failed.
- Temps jusqu’à annotation ou handoff.
- Sorties avec provenance complète.

## 24. Classification de livraison
`defined` / `planned`; aucune preuve d’implémentation, moteur, hyperviseur, API, protocole ou commande.

## 25. Critères d’acceptation
**Given** un Artifact source et une transformation autorisée
**When** l’analyste utilise Derived Artifact Management
**Then** l’Artifact source reste inchangé, la transformation, les sources, paramètres, versions et limites sont visibles, et le dérivé ne devient pas automatiquement Evidence

**Given** une entrée manquante ou une permission refusée
**When** l’utilisateur demande Derived Artifact Management
**Then** l’erreur est explicite, les résultats valides sont conservés et aucune conclusion n’est inventée

**Given** aucun modèle IA
**When** l’analyste réalise Derived Artifact Management
**Then** les viewers, règles, Tools déterministes, filtres et revue humaine couvrent le workflow essentiel

## 26. Questions ouvertes
- Le schéma final des résultats de Derived Artifact Management est reporté à la phase Objets.
- OPEN-013 reste ouverte.
- OPEN-014 reste ouverte.
- OPEN-015 reste ouverte.

## 27. Consommateurs documentaires
- Analysis Workbench et capability map
- Case Workspace, Evidence Board et écrans techniques
- phases Objets, Permissions, Journeys et Technique
