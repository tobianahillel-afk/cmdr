---
id: CAP-INV-308
title: Script and Document Static Analysis
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
  - OPEN-014
source-of-truth: canonical
---
# CAP-INV-308 — Script and Document Static Analysis

## 1. Définition
Inspecter scripts et documents sans les exécuter, notamment métadonnées, contenu textuel, scripts intégrés, macros, liens, objets embarqués, références externes et obfuscation candidate.

## 2. Problème utilisateur
Les documents et scripts peuvent embarquer du contenu actif ou externe. Une preview naïve ou un décodage silencieux peut exécuter, altérer ou surinterpréter le contenu.

## 3. Objectifs
- Identifier type, métadonnées et contenu textuel.
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
| Contexte Script and Document Static Analysis | analyste / résultat précédent | paramètres fonctionnels | oui | snapshot de session | demander complétude |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | scope, Hypothesis et return origin | consulter |
| Artifact | Investigate | source, version, type et relations | consulter |
| Document/script structure | Tool Call | contenus, objets, macros et liens | consulter |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Analysis Session | concept Investigate | lier résultat et disposition | aucun schéma final |
| Active content observation | Investigate | créer/annoter | aucune activation |
| Decoded Derived Artifact | Investigate concept | créer | décodage déclaré et traçable |

## 11. Fonctionnalités
- Identifier type, métadonnées et contenu textuel.
- Inspecter scripts, macros, liens, objets et références externes sans activation.
- Afficher les contenus décodés comme Derived Artifacts.
- Comparer original et versions décodées en conservant l’incertitude.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Inspecter contenu | Analyst | Artifact | 0 | contexte et permission valides | contenu rendu sûr | non |
| Décoder une sélection | Analyst | Derived Artifact | 1 | contexte et permission valides | dérivé sourcé | non |
| Extraire objet embarqué | Analyst | Derived Artifact | 1 | contexte et permission valides | objet séparé | non |
| Annoter obfuscation candidate | Analyst | Observation | 2 | contexte et permission valides | incertitude visible | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| Identifier type, métadonnées et contenu textuel | oui | règles et Tools déterministes | oui | suggestion expliquée | inspection manuelle |
| Résumer Script and Document Static Analysis | oui | agrégation sourcée | oui | résumé attribué | viewer et filtres |
| Proposer une prochaine étape | oui | checklists/profils | oui | proposition modifiable | catalogue manuel |
| Qualifier Evidence ou Finding | humain | contrôles seulement | workflow de revue | jamais autonome | CAP-INV-107/108/109 |

La provenance automatisée et la disposition humaine restent visibles; aucune fonction essentielle ne dépend de l’IA.

## 14. États fonctionnels
`recognized`, `active-content-present`, `obfuscation-candidate`, `decoded`, `partial`, `unsupported`, `restricted`, `corrupted`. Machines finales reportées.

## 15. États d’interface
Loading conserve le contexte; Empty explique; Partial nomme les lacunes; Error garde les résultats valides; Offline est stale/read-only; Permission denied ne fuit rien.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Static document result | Analysis Result | session | contenus et limites visibles |
| Embedded objects | Derived Artifacts | CAP-INV-311/309 | parent et position |
| External reference candidates | candidate records | Case/future Intelligence | aucun contact réalisé |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-304/305 | type script/document | CAP-INV-308 | Artifact, preview, structure | preview |
| CAP-INV-308 | décoder/extract | CAP-INV-311 | source, transformation, Tool/version | document analysis |
| CAP-INV-308 | analyse dynamique future | future Dynamic Sandbox | original, dérivés, interactions nécessaires, provenance | document analysis |

Le return origin, les permissions, versions et résultats partiels sont conservés.

## 18. Dépendances
- CAP-INV-105 Artifact Management
- CAP-INV-302 Analysis Session Management
- CAP-INV-312 Analysis Provenance and Reproducibility
- Shared Trace/Activity/Object Linking/Export
- CAP-INV-304/305/306/309/311
- future Dynamic Sandbox

## 19. Source de vérité
Investigate est source de l’interprétation et des relations analytiques. Les objets sources gardent leur owner ; Studio reste source des Tool/Tool Call et Shared des mécanismes transversaux.

## 20. Provenance et audit
Enregistrer Case, session, Artifact/version, paramètres propres à Script and Document Static Analysis, Tool/version, Tool Calls, outputs, erreurs, annotations, disposition humaine, timestamp et correlation ID.

## 21. Permissions fonctionnelles
- Artifact read
- Analysis Session read/update
- Tool use/output read
- analysis result annotate
- script/document raw read
- macro/active content inspection
- embedded content extraction
- Derived Artifact create

Matrice atomique reportée.

## 22. Limites et erreurs
- Artifact inaccessible, superseded ou restreint.
- Tool indisponible, incompatible ou résultat partiel.
- Contexte tenant/environnement incohérent.
- Contenu chiffré ou mot de passe absent.
- Référence externe inaccessible sans téléchargement autorisé.
- Déobfuscation ambiguë ou partielle.

## 23. Métriques
- Utilisations de Script and Document Static Analysis.
- Résultats partial/failed.
- Temps jusqu’à annotation ou handoff.
- Sorties avec provenance complète.

## 24. Classification de livraison
`defined` / `planned`; aucune preuve d’implémentation, moteur, hyperviseur, API, protocole ou commande.

## 25. Critères d’acceptation
**Given** un Artifact compatible et une session autorisée
**When** l’analyste utilise Script and Document Static Analysis
**Then** le type, les métadonnées, le contenu textuel, les sources, paramètres, versions et limites sont visibles sans exécuter le contenu

**Given** une entrée manquante ou une permission refusée
**When** l’utilisateur demande Script and Document Static Analysis
**Then** l’erreur est explicite, les résultats valides sont conservés et aucune conclusion n’est inventée

**Given** aucun modèle IA
**When** l’analyste réalise Script and Document Static Analysis
**Then** les viewers, règles, Tools déterministes, filtres et revue humaine couvrent le workflow essentiel

## 26. Questions ouvertes
- Le schéma final des résultats de Script and Document Static Analysis est reporté à la phase Objets.
- OPEN-005 reste ouverte.
- OPEN-014 reste ouverte.

## 27. Consommateurs documentaires
- Analysis Workbench et capability map
- Case Workspace, Evidence Board et écrans techniques
- phases Objets, Permissions, Journeys et Technique
