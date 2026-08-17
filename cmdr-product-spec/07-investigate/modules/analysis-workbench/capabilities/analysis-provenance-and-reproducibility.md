---
id: CAP-INV-312
title: Analysis Provenance and Reproducibility
product: investigate
module: analysis-workbench
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-020
  - REQ-AI-002
  - REQ-OBJ-009
open_decisions:
  - OPEN-005
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-312 — Analysis Provenance and Reproducibility

## 1. Définition
Retracer Case, Analysis Session, Artifact sources et dérivés, Tools, versions, Tool Calls, paramètres, environnement référencé, outputs, erreurs, décisions et dispositions, puis évaluer la reproductibilité.

## 2. Problème utilisateur
Un résultat sans version de Tool, paramètres ou entrée stable ne peut pas être reproduit ni audité. Une reproduction partielle ne doit jamais être présentée comme identique.

## 3. Objectifs
- Construire une chaîne analytique navigable et attribuée.
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
| Contexte Analysis Provenance and Reproducibility | analyste / résultat précédent | paramètres fonctionnels | oui | snapshot de session | demander complétude |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | scope, Hypothesis et return origin | consulter |
| Artifact | Investigate | source, version, type et relations | consulter |
| Analysis Session | Investigate | scope, Artifacts et dispositions | consulter |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Analysis Session | concept Investigate | lier résultat et disposition | aucun schéma final |
| Provenance relation/event | Investigate sémantique / Shared mechanism | émettre/supersede | source et owner visibles |
| Reproducibility assessment | Investigate | créer/versionner/contester | préconditions et écarts obligatoires |

## 11. Fonctionnalités
- Construire une chaîne analytique navigable et attribuée.
- Détecter Tool, version, input ou environnement manquant.
- Reproduire lorsque les préconditions existent et comparer les exécutions.
- Distinguer reproductible, partiel, impossible et disputed.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Ouvrir provenance | Reviewer | Provenance chain | 0 | contexte et permission valides | chaîne affichée | non |
| Demander reproduction | Analyst | Replay request | 1 | contexte et permission valides | préconditions vérifiées | non |
| Comparer exécutions | Reviewer | Analysis Results | 0 | contexte et permission valides | écarts affichés | non |
| Contester résultat | Reviewer | Assessment | 2 | contexte et permission valides | statut disputed | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| Construire une chaîne analytique | oui | règles et relations déterministes | oui | suggestion expliquée | inspection manuelle |
| Résumer Analysis Provenance and Reproducibility | oui | agrégation sourcée | oui | résumé attribué | viewer et filtres |
| Proposer une prochaine étape | oui | checklists/profils | oui | proposition modifiable | catalogue manuel |
| Qualifier Evidence ou Finding | humain | contrôles seulement | workflow de revue | jamais autonome | CAP-INV-107/108/109 |

La provenance automatisée et la disposition humaine restent visibles; aucune fonction essentielle ne dépend de l’IA.

## 14. États fonctionnels
`reproducible`, `partially-reproducible`, `not-reproducible`, `missing-tool`, `missing-version`, `missing-input`, `environment-unavailable`, `disputed`. Machines finales reportées.

## 15. États d’interface
Loading conserve le contexte; Empty explique; Partial nomme les lacunes; Error garde les résultats valides; Offline est stale/read-only; Permission denied ne fuit rien.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Analysis provenance chain | relations/events | Case Replay/Audit | sources, versions et dispositions |
| Reproducibility assessment | Analysis Result | Reviewer | préconditions et écarts explicites |
| Replay context | request context | Tool/Workflow | scope borné et permission-aware |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Analysis Session | ouvrir trace | CAP-INV-312 | session, Artifacts, Tools, runs, outputs | session |
| CAP-INV-312 | reproduire | Tool/Workflow autorisé | inputs, versions, paramètres, environnement, permission | provenance |
| Replay result | comparer | CAP-INV-312 | outputs, erreurs, différences, dispositions | provenance |

Le return origin, les permissions, versions et résultats partiels sont conservés.

## 18. Dépendances
- CAP-INV-105 Artifact Management
- CAP-INV-302 Analysis Session Management
- Shared Trace/Activity/Object Linking/Export
- CAP-INV-008/214
- Studio versioning/Tool Calls
- OPEN-005/015

## 19. Source de vérité
Investigate est source de l’interprétation et des relations analytiques. Les objets sources gardent leur owner ; Studio reste source des Tool/Tool Call et Shared des mécanismes transversaux.

## 20. Provenance et audit
Enregistrer Case, session, Artifact/version, paramètres propres à Analysis Provenance and Reproducibility, Tool/version, Tool Calls, outputs, erreurs, annotations, disposition humaine, timestamp et correlation ID.

## 21. Permissions fonctionnelles
- Artifact read
- Analysis Session read/update
- Tool use/output read
- analysis result annotate
- reproducibility review
- replay request
- Tool/Automation Run projection read
- provenance export

Matrice atomique reportée.

## 22. Limites et erreurs
- Artifact inaccessible, superseded ou restreint.
- Tool indisponible, incompatible ou résultat partiel.
- Contexte tenant/environnement incohérent.
- Tool ou version retiré.
- Input modifié, supprimé ou inaccessible.
- Environnement indisponible ou policy-blocked.

## 23. Métriques
- Utilisations de Analysis Provenance and Reproducibility.
- Résultats partial/failed.
- Temps jusqu’à annotation ou handoff.
- Sorties avec provenance complète.

## 24. Classification de livraison
`defined` / `planned`; aucune preuve d’implémentation, moteur, hyperviseur, API, protocole ou commande.

## 25. Critères d’acceptation
**Given** une Analysis Session ancienne avec ses références historiques
**When** l’analyste ouvre Analysis Provenance and Reproducibility
**Then** la chaîne, les sources, paramètres, versions et limites sont visibles et toute version manquante est signalée

**Given** une entrée manquante ou une permission refusée
**When** l’utilisateur demande Analysis Provenance and Reproducibility
**Then** l’erreur est explicite, les résultats valides sont conservés et aucune reproduction n’est inventée

**Given** aucun modèle IA
**When** l’analyste réalise Analysis Provenance and Reproducibility
**Then** relations, viewers, règles, Tools déterministes, filtres et revue humaine couvrent le workflow essentiel

## 26. Questions ouvertes
- Le schéma final des résultats de Analysis Provenance and Reproducibility est reporté à la phase Objets.
- OPEN-005 reste ouverte.
- OPEN-015 reste ouverte.

## 27. Consommateurs documentaires
- Analysis Workbench et capability map
- Case Workspace, Evidence Board et écrans techniques
- phases Objets, Permissions, Journeys et Technique
