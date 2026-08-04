---
id: CAP-INV-309
title: Archive and Container Inspection
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
  - REQ-SEC-002
open_decisions:
  - OPEN-005
  - OPEN-013
  - OPEN-014
source-of-truth: canonical
---
# CAP-INV-309 — Archive and Container Inspection

## 1. Définition
Inspecter une archive ou un conteneur, sa structure, ses éléments, profondeurs, tailles, types, protections et risques d’expansion, puis extraire de manière bornée des Derived Artifacts sans exécution.

## 2. Problème utilisateur
Une archive imbriquée peut provoquer une expansion excessive, contenir des éléments protégés ou masquer des relations parent/enfant. Une extraction non bornée menace la disponibilité et détruit la provenance.

## 3. Objectifs
- Afficher structure, éléments, tailles, types, profondeur et protections.
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
| Contexte Archive and Container Inspection | analyste / résultat précédent | paramètres fonctionnels | oui | snapshot de session | demander complétude |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | scope, Hypothesis et return origin | consulter |
| Artifact | Investigate | source, version, type et relations | consulter |
| Archive structure | Tool Call | arbre, tailles, types, erreurs | consulter |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Analysis Session | concept Investigate | lier résultat et disposition | aucun schéma final |
| Extraction plan | Investigate | créer/modifier | limites et sélection explicites |
| Derived Artifact | Investigate concept | créer | parent, chemin interne, Tool et paramètres |

## 11. Fonctionnalités
- Afficher structure, éléments, tailles, types, profondeur et protections.
- Évaluer risque de volume ou d’expansion avant extraction.
- Permettre sélection, limites, interruption et extraction bornée.
- Conserver parent/enfant, restrictions et provenance de chaque élément.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Inspecter archive | Analyst | Artifact | 0 | contexte et permission valides | arbre affiché | non |
| Définir limites | Analyst | Extraction plan | 2 | contexte et permission valides | plan versionné | OPEN-013 |
| Extraire sélection | Analyst | Derived Artifacts | 1 | contexte et permission valides | éléments sourcés | non |
| Interrompre extraction | Analyst | Background Job | 2 | contexte et permission valides | partials conservés | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| Afficher structure, éléments, tailles, types et profondeur | oui | règles et Tools déterministes | oui | suggestion expliquée | inspection manuelle |
| Résumer Archive and Container Inspection | oui | agrégation sourcée | oui | résumé attribué | viewer et filtres |
| Proposer une prochaine étape | oui | checklists/profils | oui | proposition modifiable | catalogue manuel |
| Qualifier Evidence ou Finding | humain | contrôles seulement | workflow de revue | jamais autonome | CAP-INV-107/108/109 |

La provenance automatisée et la disposition humaine restent visibles; aucune fonction essentielle ne dépend de l’IA.

## 14. États fonctionnels
`recognized`, `protected`, `expansion-risk`, `ready`, `extracting`, `partial`, `completed`, `failed`, `cancelled`, `depth-limited`. Machines finales reportées.

## 15. États d’interface
Loading conserve le contexte; Empty explique; Partial nomme les lacunes; Error garde les résultats valides; Offline est stale/read-only; Permission denied ne fuit rien.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Archive inventory | Analysis Result | analyste | éléments et erreurs visibles |
| Extraction plan/result | analysis records | session/Trace | limites et disposition |
| Extracted elements | Derived Artifacts | CAP-INV-301/311 | parent/enfant et provenance |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-305 | structure archive reconnue | CAP-INV-309 | Artifact, structure, restrictions | structure |
| CAP-INV-309 | extraire sélection | CAP-INV-311 | parent, chemins, limites, Tool/version | archive |
| Derived Artifact | nouvelle analyse | CAP-INV-301 | parent, transformation, restrictions, Case | archive |

Le return origin, les permissions, versions et résultats partiels sont conservés.

## 18. Dépendances
- CAP-INV-105 Artifact Management
- CAP-INV-302 Analysis Session Management
- CAP-INV-312 Analysis Provenance and Reproducibility
- Shared Trace/Activity/Object Linking/Export
- Shared Background Jobs
- CAP-INV-311

## 19. Source de vérité
Investigate est source de l’interprétation et des relations analytiques. Les objets sources gardent leur owner ; Studio reste source des Tool/Tool Call et Shared des mécanismes transversaux.

## 20. Provenance et audit
Enregistrer Case, session, Artifact/version, paramètres propres à Archive and Container Inspection, Tool/version, Tool Calls, outputs, erreurs, annotations, disposition humaine, timestamp et correlation ID.

## 21. Permissions fonctionnelles
- Artifact read
- Analysis Session read/update
- Tool use/output read
- analysis result annotate
- archive inspect
- archive extraction
- Derived Artifact create
- Background Job cancel

Matrice atomique reportée.

## 22. Limites et erreurs
- Artifact inaccessible, superseded ou restreint.
- Tool indisponible, incompatible ou résultat partiel.
- Contexte tenant/environnement incohérent.
- Mot de passe absent ou protection inconnue.
- Risque d’expansion ou profondeur excessive.
- Élément malformé, duplicate ou extraction partielle.

## 23. Métriques
- Utilisations de Archive and Container Inspection.
- Résultats partial/failed.
- Temps jusqu’à annotation ou handoff.
- Sorties avec provenance complète.

## 24. Classification de livraison
`defined` / `planned`; aucune preuve d’implémentation, moteur, hyperviseur, API, protocole ou commande.

## 25. Critères d’acceptation
**Given** un Artifact compatible et une session autorisée
**When** l’analyste utilise Archive and Container Inspection
**Then** la structure, les éléments, tailles, types, profondeurs, sources, paramètres, versions et limites sont visibles sans exécuter le contenu

**Given** une entrée manquante ou une permission refusée
**When** l’utilisateur demande Archive and Container Inspection
**Then** l’erreur est explicite, les résultats valides sont conservés et aucune conclusion n’est inventée

**Given** aucun modèle IA
**When** l’analyste réalise Archive and Container Inspection
**Then** les viewers, règles, Tools déterministes, filtres et revue humaine couvrent le workflow essentiel

## 26. Questions ouvertes
- Le schéma final des résultats de Archive and Container Inspection est reporté à la phase Objets.
- OPEN-005 reste ouverte.
- OPEN-013 reste ouverte.
- OPEN-014 reste ouverte.

## 27. Consommateurs documentaires
- Analysis Workbench et capability map
- Case Workspace, Evidence Board et écrans techniques
- phases Objets, Permissions, Journeys et Technique
