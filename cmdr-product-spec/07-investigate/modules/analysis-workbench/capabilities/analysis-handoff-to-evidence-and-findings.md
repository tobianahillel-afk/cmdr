---
id: CAP-INV-313
title: Analysis Handoff to Evidence and Findings
product: investigate
module: analysis-workbench
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-014
  - REQ-OBJ-004
  - REQ-AI-002
open_decisions:
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-313 — Analysis Handoff to Evidence and Findings

## 1. Définition
Sélectionner observations, résultats, Derived Artifacts, contradictions et sources d’une analyse pour préparer une Evidence candidate ou un Finding Draft, sans qualification ou confirmation automatique.

## 2. Problème utilisateur
Les outputs de Tools et observations ne sont pas des Evidence ou Findings. Sans handoff explicite, les résultats automatisés peuvent acquérir une autorité qu’ils n’ont pas.

## 3. Objectifs
- Conserver la relation avec Case et Hypothesis.
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
| Contexte Analysis Handoff to Evidence and Findings | analyste / résultat précédent | paramètres fonctionnels | oui | snapshot de session | demander complétude |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | scope, Hypothesis et return origin | consulter |
| Artifact | Investigate | source, version, type et relations | consulter |
| Analysis Session/Result | Investigate | observations, outputs, erreurs et dispositions | consulter |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Analysis Session | concept Investigate | lier résultat et disposition | aucun schéma final |
| Evidence candidate | Investigate | préparer/modifier/retirer | qualification CAP-INV-107/108 obligatoire |
| Finding Draft | Investigate | préparer/modifier | confirmation CAP-INV-109 obligatoire |

## 11. Fonctionnalités
- Conserver la relation avec Case et Hypothesis.
- Sélectionner sources, observations, dérivés et contradictions.
- Préparer Evidence candidate et Finding Draft avec provenance.
- Revenir au Workbench après revue ou rejet.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Sélectionner résultats | Analyst | Analysis Results | 0 | contexte et permission valides | sélection conservée | non |
| Préparer Evidence candidate | Analyst | Evidence candidate | 2 | contexte et permission valides | draft sourcé | OPEN-013 |
| Préparer Finding Draft | Analyst | Finding Draft | 2 | contexte et permission valides | draft avec contradictions | OPEN-013 |
| Transmettre à revue | Analyst/Lead | Handoff | 2 | contexte et permission valides | contexte remis | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| Conserver la relation avec Case et Hypothesis | oui | règles et relations déterministes | oui | suggestion expliquée | inspection manuelle |
| Résumer Analysis Handoff to Evidence and Findings | oui | agrégation sourcée | oui | résumé attribué | viewer et filtres |
| Proposer une prochaine étape | oui | checklists/profils | oui | proposition modifiable | catalogue manuel |
| Qualifier Evidence ou Finding | humain | contrôles seulement | workflow de revue | jamais autonome | CAP-INV-107/108/109 |

La provenance automatisée et la disposition humaine restent visibles; aucune fonction essentielle ne dépend de l’IA.

## 14. États fonctionnels
`draft`, `incomplete`, `ready-for-review`, `submitted`, `returned`, `accepted-as-candidate`, `rejected`, `superseded`. Machines finales reportées.

## 15. États d’interface
Loading conserve le contexte; Empty explique; Partial nomme les lacunes; Error garde les résultats valides; Offline est stale/read-only; Permission denied ne fuit rien.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Evidence candidate package | draft context | CAP-INV-107/108 | source, provenance et qualification nécessaire |
| Finding Draft package | draft context | CAP-INV-109 | observations, Evidence et contradictions |
| Handoff disposition | business event | session/Trace | accepté, modifié ou rejeté |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Analysis Session | préparer Evidence | CAP-INV-313 | résultats, sources, dérivés, Hypothesis, provenance | session |
| CAP-INV-313 | soumettre candidate | CAP-INV-107/108 | candidate, Artifact, Case, qualification requise | Workbench |
| CAP-INV-313 | soumettre Finding Draft | CAP-INV-109 | draft, Evidence existantes, contradictions, incertitude | Workbench |

Le return origin, les permissions, versions et résultats partiels sont conservés.

## 18. Dépendances
- CAP-INV-105 Artifact Management
- CAP-INV-302 Analysis Session Management
- CAP-INV-312 Analysis Provenance and Reproducibility
- Shared Trace/Activity/Object Linking/Export
- CAP-INV-103/107/108/109
- Case Replay/Timeline

## 19. Source de vérité
Investigate est source de l’interprétation et des relations analytiques. Les objets sources gardent leur owner ; Studio reste source des Tool/Tool Call et Shared des mécanismes transversaux.

## 20. Provenance et audit
Enregistrer Case, session, Artifact/version, paramètres propres à Analysis Handoff to Evidence and Findings, Tool/version, Tool Calls, outputs, erreurs, annotations, disposition humaine, timestamp et correlation ID.

## 21. Permissions fonctionnelles
- Artifact read
- Analysis Session read/update
- Tool use/output read
- analysis result annotate
- Evidence candidate prepare
- Finding Draft prepare
- handoff submit

Matrice atomique reportée.

## 22. Limites et erreurs
- Artifact inaccessible, superseded ou restreint.
- Tool indisponible, incompatible ou résultat partiel.
- Contexte tenant/environnement incohérent.
- Source inaccessible ou provenance incomplète.
- Evidence candidate rejetée ou renvoyée.
- Finding Draft sans Evidence suffisante.

## 23. Métriques
- Utilisations de Analysis Handoff to Evidence and Findings.
- Résultats partial/failed.
- Temps jusqu’à annotation ou handoff.
- Sorties avec provenance complète.

## 24. Classification de livraison
`defined` / `planned`; aucune preuve d’implémentation, moteur, hyperviseur, API, protocole ou commande.

## 25. Critères d’acceptation
**Given** des résultats statiques reliés à un Case
**When** l’analyste utilise Analysis Handoff to Evidence and Findings
**Then** les observations, sources, contradictions, paramètres, versions et limites sont visibles et la qualification reste explicitement nécessaire

**Given** une entrée manquante ou une permission refusée
**When** l’utilisateur demande Analysis Handoff to Evidence and Findings
**Then** l’erreur est explicite, les résultats valides sont conservés et aucune Evidence ou Finding n’est confirmée

**Given** aucun modèle IA
**When** l’analyste réalise Analysis Handoff to Evidence and Findings
**Then** les viewers, règles, formulaires et revue humaine couvrent le workflow essentiel

## 26. Questions ouvertes
- Le schéma final des résultats de Analysis Handoff to Evidence and Findings est reporté à la phase Objets.
- OPEN-013 reste ouverte.
- OPEN-015 reste ouverte.

## 27. Consommateurs documentaires
- Analysis Workbench et capability map
- Case Workspace, Evidence Board et écrans techniques
- phases Objets, Permissions, Journeys et Technique
