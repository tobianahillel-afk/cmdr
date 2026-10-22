---
id: CAP-INV-304
title: Artifact Preview and Metadata Inspection
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
  - REQ-OBJ-003
  - REQ-UX-002
open_decisions:
  - OPEN-014
source-of-truth: canonical
---
# CAP-INV-304 — Artifact Preview and Metadata Inspection

## 1. Définition
Prévisualiser de manière sûre un Artifact et inspecter ses métadonnées, signatures disponibles, versions, relations, restrictions et incohérences sans activer de contenu.

## 2. Problème utilisateur
Une preview active ou une métadonnée présentée comme conclusion peut exécuter un contenu dangereux ou induire l’analyste en erreur.

## 3. Objectifs
- Fournir une preview sûre lorsque le format le permet.
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
| Contexte Artifact Preview and Metadata Inspection | analyste / résultat précédent | paramètres fonctionnels | oui | snapshot de session | demander complétude |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | scope, Hypothesis et return origin | consulter |
| Artifact | Investigate | source, version, type et relations | consulter |
| File and Evidence Preview | Shared Capabilities | rendu sûr et statut | consulter |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Analysis Session | concept Investigate | lier résultat et disposition | aucun schéma final |
| Annotation | Investigate | créer/supersede | auteur et source obligatoires |
| Extraction request draft | Investigate | préparer | aucun contenu actif exécuté |

## 11. Fonctionnalités
- Fournir une preview sûre lorsque le format le permet.
- Afficher nom, taille, types déclaré/détecté, timestamps, source et acquisition.
- Rendre visibles signatures, incohérences, versions et restrictions.
- Permettre annotations, copie autorisée et préparation d’extraction.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Ouvrir preview | Analyst | Artifact | 0 | contexte et permission valides | rendu sûr | non |
| Copier une valeur | Analyst | Metadata | 0 | contexte et permission valides | valeur copiée | non |
| Annoter une incohérence | Analyst | Annotation | 2 | contexte et permission valides | annotation attribuée | OPEN-013 |
| Préparer extraction | Analyst | Extraction draft | 1 | contexte et permission valides | scope explicite | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| Assister fournir une preview sûre lorsque le format le permet. | oui | règles et Tools déterministes | oui | suggestion expliquée | inspection manuelle |
| Résumer Artifact Preview and Metadata Inspection | oui | agrégation sourcée | oui | résumé attribué | viewer et filtres |
| Proposer une prochaine étape | oui | checklists/profils | oui | proposition modifiable | catalogue manuel |
| Qualifier Evidence ou Finding | humain | contrôles seulement | workflow de revue | jamais autonome | CAP-INV-107/108/109 |

La provenance automatisée et la disposition humaine restent visibles; aucune fonction essentielle ne dépend de l’IA.

## 14. États fonctionnels
`preview-available`, `preview-unavailable`, `metadata-partial`, `type-mismatch`, `restricted`, `corrupted`, `stale`. Machines finales reportées.

## 15. États d’interface
Loading conserve le contexte; Empty explique; Partial nomme les lacunes; Error garde les résultats valides; Offline est stale/read-only; Permission denied ne fuit rien.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Safe preview | preview projection | analyste | aucun contenu actif |
| Metadata set | Analysis Result concept | CAP-INV-305/306 | source et disponibilité visibles |
| Annotation | investigation annotation | Case | aucun verdict implicite |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-301 | ouvrir aperçu | CAP-INV-304 | Artifact/version, restrictions, type | intake |
| CAP-INV-304 | inspection approfondie | CAP-INV-305 | Artifact, métadonnées, divergence | preview |
| CAP-INV-304 | préparer extraction | CAP-INV-306/309 | scope, source et restrictions | preview |

Le return origin, les permissions, versions et résultats partiels sont conservés.

## 18. Dépendances
- CAP-INV-105 Artifact Management
- CAP-INV-302 Analysis Session Management
- CAP-INV-312 Analysis Provenance and Reproducibility
- Shared Trace/Activity/Object Linking/Export
- Shared File and Evidence Preview
- Inspector

## 19. Source de vérité
Investigate est source de l’interprétation et des relations analytiques. Les objets sources gardent leur owner ; Studio reste source des Tool/Tool Call et Shared des mécanismes transversaux.

## 20. Provenance et audit
Enregistrer Case, session, Artifact/version, paramètres propres à Artifact Preview and Metadata Inspection, Tool/version, Tool Calls, outputs, erreurs, annotations, disposition humaine, timestamp et correlation ID.

## 21. Permissions fonctionnelles
- Artifact read
- Analysis Session read/update
- Tool use/output read
- analysis result annotate
- Artifact preview
- raw content read
- metadata copy
- analysis annotation

Matrice atomique reportée.

## 22. Limites et erreurs
- Artifact inaccessible, superseded ou restreint.
- Tool indisponible, incompatible ou résultat partiel.
- Contexte tenant/environnement incohérent.
- Preview non supportée ou renderer indisponible.
- Métadonnées conflictuelles ou timezone inconnue.
- Signature absente, invalide ou non vérifiable.

## 23. Métriques
- Utilisations de Artifact Preview and Metadata Inspection.
- Résultats partial/failed.
- Temps jusqu’à annotation ou handoff.
- Sorties avec provenance complète.

## 24. Classification de livraison
`defined` / `planned`; aucune preuve d’implémentation, moteur, hyperviseur, API, protocole ou commande.

## 25. Critères d’acceptation
**Given** un Artifact compatible et une session autorisée
**When** l’analyste utilise Artifact Preview and Metadata Inspection
**Then** la preview sûre, les sources, paramètres, versions et limites sont visibles sans exécuter le contenu

**Given** une entrée manquante ou une permission refusée
**When** l’utilisateur demande Artifact Preview and Metadata Inspection
**Then** l’erreur est explicite, les résultats valides sont conservés et aucune conclusion n’est inventée

**Given** aucun modèle IA
**When** l’analyste réalise Artifact Preview and Metadata Inspection
**Then** les viewers, règles, Tools déterministes, filtres et revue humaine couvrent le workflow essentiel

## 26. Questions ouvertes
- Le schéma final des résultats de Artifact Preview and Metadata Inspection est reporté à la phase Objets.
- OPEN-014 reste ouverte.

## 27. Consommateurs documentaires
- Analysis Workbench et capability map
- Case Workspace, Evidence Board et écrans techniques
- phases Objets, Permissions, Journeys et Technique
