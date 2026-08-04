---
id: CAP-INV-325
title: Multi-Run Comparison
product: investigate
module: dynamic-sandbox
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-INV-005
  - REQ-PROD-014
  - REQ-PROD-020
  - REQ-AI-002
open_decisions:
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-325 — Multi-Run Comparison

## 1. Définition
Comparer plusieurs Sandbox Runs selon environnements, profils, versions, timelines, processus, fichiers, changements système, réseau et Runtime Artifacts, en distinguant commun, ajouté, supprimé, modifié et non observé.

## 2. Problème utilisateur
Une différence peut provenir d’un environnement, profil, version, interaction, timeout ou couverture différente. Sans préconditions visibles, elle peut être interprétée à tort comme une évolution de l’Artifact.

## 3. Objectifs
- Sélectionner plusieurs Runs et afficher leurs préconditions.
- Comparer toutes les dimensions comportementales autorisées.
- Distinguer différence et absence d’observation.
- Relier les écarts aux Hypotheses et préparer un Finding Draft sans conclusion automatique.

## 4. Non-objectifs
Ne pas fusionner les Runs ou Artifacts, modifier leurs résultats, produire un score opaque, choisir un moteur, définir forensics réseau ou confirmer automatiquement un Finding.

## 5. Propriétaire
Investigate possède la comparaison et son interprétation; Settings/Studio conservent environnements, profils, Tools et versions.

## 6. Utilisateurs
Malware Analyst; Dynamic Analysis Operator; Case Analyst; Reviewer.

## 7. Conditions d’entrée
Au moins deux Runs accessibles; Artifact et Case cohérents; préconditions et versions disponibles ou lacunes explicites; permission de comparaison.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Sandbox Runs | CAP-INV-317 | exécutions comparées | oui, au moins deux | états terminaux/partiels visibles | comparaison impossible |
| Environnements/profils/versions | Settings/Studio/Session | préconditions | oui | snapshots des Runs | partially-comparable |
| Timelines/process/file/system/network | CAP-INV-318..322 | dimensions comportementales | selon scope | même version de résultat | dimension unavailable |
| Runtime Artifacts | CAP-INV-324 | sorties matérielles | non | versions visibles | dimension vide |
| Objectif de comparaison | analyste | question analytique | oui | version du comparison profile | rester draft |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Dynamic Analysis Session | Investigate concept | Artifact, objectif et Runs | consulter |
| Sandbox Run | Investigate concept | préconditions, résultats et statuts | consulter/comparer |
| Sandbox Environment / Profile / Tool | Settings/Studio | versions et limitations | consulter |
| Runtime Artifact / Observations | Investigate concepts | dimensions et lineage | consulter |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Multi-Run comparison | créer/versionner/supersede | Investigate concept | Runs et dimensions obligatoires |
| Difference observation | annoter/lier | Investigate | différence distincte de conclusion |
| Hypothesis relation | créer/supersede | Investigate | CAP-INV-103 conserve lifecycle |
| Finding Draft context | préparer | Investigate | CAP-INV-109 requis |

## 11. Fonctionnalités
Comparer environnements, profils, versions, timelines, process trees, fichiers, système, réseau et Runtime Artifacts; distinguer commun/ajouté/supprimé/modifié/non observé; conserver filtres et observations.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Sélectionner Runs | Analyst | Run set | 0 | Runs accessibles | sélection conservée | non |
| Configurer dimensions | Analyst | Comparison profile | 2 | préconditions visibles | profil versionné | OPEN-013 |
| Comparer | Analyst | Multi-Run comparison | 0 | dimensions comparables | diff affiché | non |
| Annoter différence | Analyst | Observation | 2 | source visible | annotation attribuée | OPEN-013 |
| Relier Hypothesis | Analyst | Relation | 2 | justification | lien sourcé | OPEN-013 |
| Préparer Finding Draft | Reviewer | Draft context | 2 | résultats/contradictions | CAP-INV-328 | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Aligner dimensions | oui | IDs/versions | oui | explication | matrice de comparaison |
| Détecter différences | oui | comparateurs déterministes | oui | résumé | diff brut |
| Proposer regroupement | oui | règles | oui | suggestion modifiable | filtres |
| Produire conclusion | humain | contrôles seulement | workflow de revue | jamais autonome | Hypothesis/Finding review |

## 14. États fonctionnels
`draft`, `validating`, `comparable`, `partially-comparable`, `completed`, `failed`, `stale`, `superseded`.

## 15. États d’interface
Loading conserve les Runs; Empty explique; Partial nomme dimensions manquantes; Error garde les comparaisons valides; Offline est read-only; Permission denied redacted.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Multi-Run comparison | Analysis Result | analyste/Case | Runs, préconditions et dimensions |
| Difference observations | observations | Hypothesis/CAP-INV-328 | aucune conclusion automatique |
| Comparison context | session relation | Case Replay | filtres, versions et return origin |
| Non-observed explanation | quality result | Reviewer | timeout/coverage distingués d’absence |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-316/317 | sélection multiple | CAP-INV-325 | Runs, environments, profiles, versions, results | session/Run |
| CAP-INV-325 | ouvrir dimension | CAP-INV-318..324 | différence sélectionnée, sources, return origin | comparison |
| CAP-INV-325 | préparer Finding | CAP-INV-328/109 | différences, commun, contradictions, incertitude | comparison |

## 18. Dépendances
CAP-INV-103/316/317/318..324/327/328; Settings versions; Studio Tool versions; Shared Versioning/Trace; OPEN-013/015.

## 19. Source de vérité
Investigate possède la comparaison; chaque Run et résultat source reste inchangé; différence ≠ conclusion; non observé ≠ absent.

## 20. Provenance et audit
Runs, environment/profile/Tool versions, paramètres, dimensions, filtres, résultats partiels, differences, annotations, Hypothesis links et dispositions.

## 21. Permissions fonctionnelles
Multi-Run comparison; Run/results read; comparison save/annotate; Hypothesis link; Finding Draft prepare; sensitive results read; cross-tenant denied.

## 22. Limites et erreurs
Versions incompatibles; dimensions absentes; timeout; couverture différente; résultat partiel; Artifact/version différent; permission retirée; source stale.

## 23. Métriques
Comparaisons complètes/partielles; dimensions indisponibles; différences annotées; non-observed explanations; Findings Draft préparés; provenance complète.

## 24. Classification de livraison
`defined` / `planned`; aucun moteur, score opaque, API, protocole, commande ou implémentation.

## 25. Critères d’acceptation
**Given** deux Runs avec environnements ou profils différents
**When** l’analyste les compare
**Then** préconditions, versions et différences sont visibles avant toute interprétation

**Given** un comportement non observé dans un Run timed-out
**When** la différence est affichée
**Then** elle est marquée non observée avec le timeout, pas absente

**Given** aucun modèle IA
**When** la comparaison est réalisée
**Then** comparateurs, filtres et revue humaine couvrent le workflow

## 26. Questions ouvertes
Multi-Run comparison reste un concept; OPEN-013 et OPEN-015 restent ouvertes.

## 27. Consommateurs documentaires
Dynamic Sandbox, Session Management, Case Replay, Hypothesis, Evidence/Finding handoff, Objets et Permissions.
