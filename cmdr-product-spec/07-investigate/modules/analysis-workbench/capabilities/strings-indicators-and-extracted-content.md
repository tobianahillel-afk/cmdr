---
id: CAP-INV-306
title: Strings, Indicators and Extracted Content
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
  - REQ-AI-002
open_decisions:
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-306 — Strings, Indicators and Extracted Content

## 1. Définition
Extraire, rechercher, filtrer, grouper, comparer et annoter des chaînes et contenus candidats tout en conservant encodage, position, origine et contexte.

## 2. Problème utilisateur
Une chaîne isolée peut être bénigne, obfusquée ou hors contexte. La transformer automatiquement en IOC, Evidence ou Finding produirait des conclusions non fiables.

## 3. Objectifs
- Extraire des chaînes selon représentations et encodages déclarés.
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
| Contexte Strings, Indicators and Extracted Content | analyste / résultat précédent | paramètres fonctionnels | oui | snapshot de session | demander complétude |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | scope, Hypothesis et return origin | consulter |
| Artifact | Investigate | source, version, type et relations | consulter |
| Strings result | Tool Call déterministe | valeurs, encodage, position | consulter |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Analysis Session | concept Investigate | lier résultat et disposition | aucun schéma final |
| Extracted content result | Investigate | créer/annoter/exclure | source et position obligatoires |
| Indicator Candidate | concept ouvert | proposer/retirer de l’usage | aucune confirmation automatique |

## 11. Fonctionnalités
- Extraire des chaînes selon représentations et encodages déclarés.
- Conserver position, origine, contexte et transformation.
- Identifier URLs, domaines, IP, chemins et commandes comme candidats seulement.
- Préparer Entity, Indicator Candidate ou Derived Artifact par promotion explicite.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Extraire chaînes | Analyst | Artifact | 1 | contexte et permission valides | résultat sourcé | non |
| Rechercher/filtrer/grouper | Analyst | Strings result | 0 | contexte et permission valides | vue modifiée | non |
| Exclure faux candidat | Analyst | Candidate | 2 | contexte et permission valides | disposition attribuée | OPEN-013 |
| Relier à Entity/Case | Analyst | Relation | 2 | contexte et permission valides | lien sourcé | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| Assister extraire des chaînes selon représentations et encodages déclarés. | oui | règles et Tools déterministes | oui | suggestion expliquée | inspection manuelle |
| Résumer Strings, Indicators and Extracted Content | oui | agrégation sourcée | oui | résumé attribué | viewer et filtres |
| Proposer une prochaine étape | oui | checklists/profils | oui | proposition modifiable | catalogue manuel |
| Qualifier Evidence ou Finding | humain | contrôles seulement | workflow de revue | jamais autonome | CAP-INV-107/108/109 |

La provenance automatisée et la disposition humaine restent visibles; aucune fonction essentielle ne dépend de l’IA.

## 14. États fonctionnels
`extracting`, `available`, `partial`, `encoding-ambiguous`, `candidate`, `excluded`, `restricted`, `superseded`. Machines finales reportées.

## 15. États d’interface
Loading conserve le contexte; Empty explique; Partial nomme les lacunes; Error garde les résultats valides; Offline est stale/read-only; Permission denied ne fuit rien.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Strings collection | Analysis Result | analyste | encodage, position et source |
| Indicator candidates | candidate records | Case/Intelligence future | non confirmés |
| Derived content | Derived Artifact | CAP-INV-311 | transformation et parent |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-305 | extraire contenu | CAP-INV-306 | Artifact, région, encodages, restrictions | structure |
| CAP-INV-306 | promouvoir candidat | Entity/Case or future Intelligence | valeur, contexte, source, incertitude | strings |
| CAP-INV-306 | créer dérivé | CAP-INV-311 | contenu, transformation, parent, Tool | strings |

Le return origin, les permissions, versions et résultats partiels sont conservés.

## 18. Dépendances
- CAP-INV-105 Artifact Management
- CAP-INV-302 Analysis Session Management
- CAP-INV-312 Analysis Provenance and Reproducibility
- Shared Trace/Activity/Object Linking/Export
- Shared Entity/Object Linking
- future Intelligence

## 19. Source de vérité
Investigate est source de l’interprétation et des relations analytiques. Les objets sources gardent leur owner ; Studio reste source des Tool/Tool Call et Shared des mécanismes transversaux.

## 20. Provenance et audit
Enregistrer Case, session, Artifact/version, paramètres propres à Strings, Indicators and Extracted Content, Tool/version, Tool Calls, outputs, erreurs, annotations, disposition humaine, timestamp et correlation ID.

## 21. Permissions fonctionnelles
- Artifact read
- Analysis Session read/update
- Tool use/output read
- analysis result annotate
- strings extraction
- indicator candidate prepare
- Entity link
- Derived Artifact create/export

Matrice atomique reportée.

## 22. Limites et erreurs
- Artifact inaccessible, superseded ou restreint.
- Tool indisponible, incompatible ou résultat partiel.
- Contexte tenant/environnement incohérent.
- Encodage indéterminé ou données binaires.
- Volume excessif ou résultat tronqué.
- Candidat sensible, redacted ou non exportable.

## 23. Métriques
- Utilisations de Strings, Indicators and Extracted Content.
- Résultats partial/failed.
- Temps jusqu’à annotation ou handoff.
- Sorties avec provenance complète.

## 24. Classification de livraison
`defined` / `planned`; aucune preuve d’implémentation, moteur, hyperviseur, API, protocole ou commande.

## 25. Critères d’acceptation
**Given** un Artifact compatible et une session autorisée
**When** l’analyste utilise Strings, Indicators and Extracted Content
**Then** les chaînes, représentations, sources, paramètres, versions et limites sont visibles sans exécuter le contenu

**Given** une entrée manquante ou une permission refusée
**When** l’utilisateur demande Strings, Indicators and Extracted Content
**Then** l’erreur est explicite, les résultats valides sont conservés et aucune conclusion n’est inventée

**Given** aucun modèle IA
**When** l’analyste réalise Strings, Indicators and Extracted Content
**Then** les viewers, règles, Tools déterministes, filtres et revue humaine couvrent le workflow essentiel

## 26. Questions ouvertes
- Le schéma final des résultats de Strings, Indicators and Extracted Content est reporté à la phase Objets.
- OPEN-013 reste ouverte.
- OPEN-015 reste ouverte.

## 27. Consommateurs documentaires
- Analysis Workbench et capability map
- Case Workspace, Evidence Board et écrans techniques
- phases Objets, Permissions, Journeys et Technique
