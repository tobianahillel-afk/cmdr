---
id: CAP-INV-303
title: Workbench Context and Tool Selection
product: investigate
module: analysis-workbench
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-014
  - REQ-AI-002
  - REQ-OBJ-009
  - REQ-UX-002
open_decisions:
  - OPEN-005
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-303 — Workbench Context and Tool Selection

## 1. Définition
Préserver le contexte du Technical Workbench et permettre la sélection explicite d’un Tool autorisé, de sa version et de ses paramètres fonctionnels sans administrer le catalogue Studio ni les environnements Settings.

## 2. Problème utilisateur
Sans contexte et sélection transparents, un analyste peut perdre le Case actif, lancer un Tool incompatible ou croire qu’une intégration temporaire constitue une capability native.

## 3. Objectifs
- Afficher Case, Artifact actif, sélection, objectif et autres Artifacts liés.
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
| Contexte Workbench Context and Tool Selection | analyste / résultat précédent | paramètres fonctionnels | oui | snapshot de session | demander complétude |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | scope, Hypothesis et return origin | consulter |
| Artifact | Investigate | source, version, type et relations | consulter |
| Tool | CMDR Studio | identité, version, statut, permissions et limites | sélectionner/utiliser |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Analysis Session | concept Investigate | lier résultat et disposition | aucun schéma final |
| Tool selection record | Investigate | créer/supersede | choix, version et justification visibles |
| Tool Call | CMDR Studio | demander/interrompre | ownership Studio conservé |

## 11. Fonctionnalités
- Afficher Case, Artifact actif, sélection, objectif et autres Artifacts liés.
- Présenter Tools compatibles, versions, permissions, limitations et environnements disponibles.
- Lancer explicitement un Tool autorisé et exposer ses Tool Calls.
- Conserver au retour les tabs, filtres, sélection et return origin.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Sélectionner un Tool | Analyst | Tool | 0 | contexte et permission valides | sélection visible | non |
| Modifier les paramètres fonctionnels | Analyst | Tool invocation draft | 2 | contexte et permission valides | nouvelle version du draft | OPEN-013 |
| Lancer l’analyse | Analyst | Tool Call | 1 | contexte et permission valides | Tool Call attribué | non |
| Interrompre une opération | Analyst | Tool Call | 2 | contexte et permission valides | statut interrupt-requested | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| Assister afficher case, artifact actif, sélection, objectif et autres artifacts liés. | oui | règles et Tools déterministes | oui | suggestion expliquée | inspection manuelle |
| Résumer workbench context and tool selection | oui | agrégation sourcée | oui | résumé attribué | viewer et filtres |
| Proposer une prochaine étape | oui | checklists/profils | oui | proposition modifiable | catalogue manuel |
| Qualifier Evidence ou Finding | humain | contrôles seulement | workflow de revue | jamais autonome | CAP-INV-107/108/109 |

La provenance automatisée et la disposition humaine restent visibles; aucune fonction essentielle ne dépend de l’IA.

## 14. États fonctionnels
`context-ready`, `tool-compatible`, `tool-incompatible`, `permission-denied`, `environment-unavailable`, `running`, `partial`, `completed`, `failed`, `cancelled`. Machines finales reportées.

## 15. États d’interface
Loading conserve le contexte; Empty explique; Partial nomme les lacunes; Error garde les résultats valides; Offline est stale/read-only; Permission denied ne fuit rien.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Workbench context | context projection | analyste | Case, Artifact, sélection et objectif conservés |
| Tool invocation | Tool Call | Studio/Trace | version, paramètres et initiateur visibles |
| Tool result | Analysis Result concept | session | output et erreurs distincts |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Analysis Session | ouvrir le Workbench | CAP-INV-303 | session, Case, Artifact, objectif, return origin | session |
| CAP-INV-303 | lancer Tool | Studio Tool Call | Tool/version, Artifact, paramètres, permission | Workbench |
| Tool Call | terminer/échouer | Analysis Session | output, erreurs, statut, provenance | Workbench |

Le return origin, les permissions, versions et résultats partiels sont conservés.

## 18. Dépendances
- CAP-INV-105 Artifact Management
- CAP-INV-302 Analysis Session Management
- CAP-INV-312 Analysis Provenance and Reproducibility
- Shared Trace/Activity/Object Linking/Export
- CMDR Studio Tool catalogue/versioning
- Platform Settings environment health

## 19. Source de vérité
Investigate est source de l’interprétation et des relations analytiques. Les objets sources gardent leur owner ; Studio reste source des Tool/Tool Call et Shared des mécanismes transversaux.

## 20. Provenance et audit
Enregistrer Case, session, Artifact/version, paramètres propres à Workbench Context and Tool Selection, Tool/version, Tool Calls, outputs, erreurs, annotations, disposition humaine, timestamp et correlation ID.

## 21. Permissions fonctionnelles
- Artifact read
- Analysis Session read/update
- Tool use/output read
- analysis result annotate
- Tool use
- Tool version selection
- Tool output read
- Tool interrupt

Matrice atomique reportée.

## 22. Limites et erreurs
- Artifact inaccessible, superseded ou restreint.
- Tool indisponible, incompatible ou résultat partiel.
- Contexte tenant/environnement incohérent.
- Catalogue indisponible ou version retirée.
- Tool temporaire non classifié.
- Automation Tray ou Trace indisponible.

## 23. Métriques
- Utilisations de Workbench Context and Tool Selection.
- Résultats partial/failed.
- Temps jusqu’à annotation ou handoff.
- Sorties avec provenance complète.

## 24. Classification de livraison
`defined` / `planned`; aucune preuve d’implémentation, moteur, hyperviseur, API, protocole ou commande.

## 25. Critères d’acceptation
**Given** un Artifact compatible et une session autorisée
**When** l’analyste utilise Workbench Context and Tool Selection
**Then** afficher Case, Artifact actif, sélection, objectif et autres Artifacts liés; les sources, paramètres, versions et limites sont visibles sans exécuter le contenu

**Given** une entrée manquante ou une permission refusée
**When** l’utilisateur demande Workbench Context and Tool Selection
**Then** l’erreur est explicite, les résultats valides sont conservés et aucune conclusion n’est inventée

**Given** aucun modèle IA
**When** l’analyste réalise Workbench Context and Tool Selection
**Then** les viewers, règles, Tools déterministes, filtres et revue humaine couvrent le workflow essentiel

## 26. Questions ouvertes
- Le schéma final des résultats de Workbench Context and Tool Selection est reporté à la phase Objets.
- OPEN-005 reste ouverte.
- OPEN-015 reste ouverte.

## 27. Consommateurs documentaires
- Analysis Workbench et capability map
- Case Workspace, Evidence Board et écrans techniques
- phases Objets, Permissions, Journeys et Technique
