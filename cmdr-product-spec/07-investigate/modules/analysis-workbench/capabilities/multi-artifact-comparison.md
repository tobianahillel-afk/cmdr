---
id: CAP-INV-310
title: Multi-Artifact Comparison
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
  - REQ-AI-002
open_decisions:
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-310 — Multi-Artifact Comparison

## 1. Définition
Comparer plusieurs Artifacts ou résultats selon métadonnées, structures, chaînes, dépendances, ressources et contenus extraits, en distinguant identique, ajouté, supprimé, modifié et non comparable.

## 2. Problème utilisateur
Des différences non contextualisées peuvent être prises pour des conclusions. L’analyste doit connaître les versions, préconditions et dimensions réellement comparées.

## 3. Objectifs
- Sélectionner plusieurs Artifacts et dimensions comparables.
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
| Contexte Multi-Artifact Comparison | analyste / résultat précédent | paramètres fonctionnels | oui | snapshot de session | demander complétude |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | scope, Hypothesis et return origin | consulter |
| Artifact | Investigate | source, version, type et relations | consulter |
| Artifact set | Investigate | versions et relations | consulter |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Analysis Session | concept Investigate | lier résultat et disposition | aucun schéma final |
| Comparison result | Investigate concept | créer/versionner | sources et profil obligatoires |
| Comparison observation | Investigate | annoter/supersede | différence distincte de conclusion |

## 11. Fonctionnalités
- Sélectionner plusieurs Artifacts et dimensions comparables.
- Conserver filtres, ordre, versions et préconditions.
- Enregistrer observations et relations aux Hypotheses.
- Préparer un Finding Draft sans fusionner ni modifier les sources.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Sélectionner Artifacts | Analyst | Artifact set | 0 | contexte et permission valides | sélection conservée | non |
| Configurer comparaison | Analyst | Comparison profile | 2 | contexte et permission valides | profil versionné | OPEN-013 |
| Comparer | Analyst | Comparison result | 0 | contexte et permission valides | diff affiché | non |
| Annoter différence | Analyst | Observation | 2 | contexte et permission valides | annotation attribuée | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| Sélectionner les Artifacts et dimensions | oui | règles et Tools déterministes | oui | suggestion expliquée | inspection manuelle |
| Résumer Multi-Artifact Comparison | oui | agrégation sourcée | oui | résumé attribué | viewer et filtres |
| Proposer une prochaine étape | oui | checklists/profils | oui | proposition modifiable | catalogue manuel |
| Qualifier Evidence ou Finding | humain | contrôles seulement | workflow de revue | jamais autonome | CAP-INV-107/108/109 |

La provenance automatisée et la disposition humaine restent visibles; aucune fonction essentielle ne dépend de l’IA.

## 14. États fonctionnels
`draft`, `validating`, `comparable`, `partially-comparable`, `completed`, `failed`, `stale`, `superseded`. Machines finales reportées.

## 15. États d’interface
Loading conserve le contexte; Empty explique; Partial nomme les lacunes; Error garde les résultats valides; Offline est stale/read-only; Permission denied ne fuit rien.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Comparison result | Analysis Result | analyste/Case | sources, dimensions et préconditions |
| Difference observations | observations | Hypothesis/Finding draft | aucune conclusion automatique |
| Saved comparison context | session relation | Case Replay | filtres et versions conservés |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Analysis Session | sélection multiple | CAP-INV-310 | Artifacts, résultats, profils, return origin | session |
| CAP-INV-310 | lier différence | Hypothesis | observation, sources, incertitude | comparison |
| CAP-INV-310 | préparer Finding | CAP-INV-313/109 | différences, Evidence existantes, contradictions | comparison |

Le return origin, les permissions, versions et résultats partiels sont conservés.

## 18. Dépendances
- CAP-INV-105 Artifact Management
- CAP-INV-302 Analysis Session Management
- CAP-INV-312 Analysis Provenance and Reproducibility
- Shared Trace/Activity/Object Linking/Export
- CAP-INV-302/312/313
- Shared Versioning/Saved Views

## 19. Source de vérité
Investigate est source de l’interprétation et des relations analytiques. Les objets sources gardent leur owner ; Studio reste source des Tool/Tool Call et Shared des mécanismes transversaux.

## 20. Provenance et audit
Enregistrer Case, session, Artifact/version, paramètres propres à Multi-Artifact Comparison, Tool/version, Tool Calls, outputs, erreurs, annotations, disposition humaine, timestamp et correlation ID.

## 21. Permissions fonctionnelles
- Artifact read
- Analysis Session read/update
- Tool use/output read
- analysis result annotate
- multi-Artifact comparison
- comparison save
- comparison annotate
- Finding Draft prepare

Matrice atomique reportée.

## 22. Limites et erreurs
- Artifact inaccessible, superseded ou restreint.
- Tool indisponible, incompatible ou résultat partiel.
- Contexte tenant/environnement incohérent.
- Versions ou dimensions incompatibles.
- Résultat Tool manquant ou non reproductible.
- Artifact inaccessible après création de la comparaison.

## 23. Métriques
- Utilisations de Multi-Artifact Comparison.
- Résultats partial/failed.
- Temps jusqu’à annotation ou handoff.
- Sorties avec provenance complète.

## 24. Classification de livraison
`defined` / `planned`; aucune preuve d’implémentation, moteur, hyperviseur, API, protocole ou commande.

## 25. Critères d’acceptation
**Given** plusieurs Artifacts compatibles et une session autorisée
**When** l’analyste utilise Multi-Artifact Comparison
**Then** les sources, dimensions, préconditions, paramètres, versions et limites sont visibles et aucune différence n’est présentée comme conclusion automatique

**Given** une entrée manquante ou une permission refusée
**When** l’utilisateur demande Multi-Artifact Comparison
**Then** l’erreur est explicite, les résultats valides sont conservés et aucune conclusion n’est inventée

**Given** aucun modèle IA
**When** l’analyste réalise Multi-Artifact Comparison
**Then** les viewers, règles, Tools déterministes, filtres et revue humaine couvrent le workflow essentiel

## 26. Questions ouvertes
- Le schéma final des résultats de Multi-Artifact Comparison est reporté à la phase Objets.
- OPEN-013 reste ouverte.
- OPEN-015 reste ouverte.

## 27. Consommateurs documentaires
- Analysis Workbench et capability map
- Case Workspace, Evidence Board et écrans techniques
- phases Objets, Permissions, Journeys et Technique
