---
id: CAP-INV-322
title: Persistence and Execution Mechanism Analysis
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
  - REQ-AI-002
open_decisions:
  - OPEN-013
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-322 — Persistence and Execution Mechanism Analysis

## 1. Définition
Regrouper les comportements candidats liés au démarrage, à la relance ou à des mécanismes d’exécution observés, en reliant processus, fichiers, changements, temporalité et contradictions sans confirmer automatiquement une persistance.

## 2. Problème utilisateur
Des changements similaires à des mécanismes de persistance peuvent être bénins, incomplets ou non activés. Les qualifier automatiquement produit un Finding trompeur et peut divulguer des méthodes offensives exécutables.

## 3. Objectifs
- Regrouper des candidats sans confirmer intention ni persistance.
- Relier processus, fichiers, changements, temporalité et contradictions.
- Comparer plusieurs Runs et documenter les écarts.
- Préparer Hypothesis ou Evidence candidate par revue humaine.

## 4. Non-objectifs
Ne pas fournir d’instructions de création de persistance, commandes, méthodes offensives exécutables, reverse/debugger, moteur, instrumentation ou confirmation automatique.

## 5. Propriétaire
Investigate possède l’interprétation candidate; Studio, Settings, Govern et Shared conservent leurs objets et mécanismes.

## 6. Utilisateurs
Malware Analyst; Case Analyst; Dynamic Analysis Operator; Evidence Reviewer.

## 7. Conditions d’entrée
Process et file/system observations disponibles; même Run identifiable; sources, timestamps et contradictions accessibles ou lacunes explicites.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Process observations | CAP-INV-319 | contexte d’exécution | oui | même Run | candidate incomplete |
| File/system changes | CAP-INV-320 | changements associés | oui | même Run | candidate incomplete |
| Behavioral Timeline | CAP-INV-318 | temporalité et source | oui | ordonnée | partial |
| Runs de comparaison | CAP-INV-325 | reproduction/contradiction | non | préconditions visibles | single-run only |
| Hypothesis existante | CAP-INV-103 | raisonnement | non | version courante | préparer nouvelle Hypothesis |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Sandbox Run | Investigate concept | environnement, profil et résultat | consulter |
| Process Observation | Investigate concept | processus et relations | consulter/lier |
| System Change | Investigate concept | fichier/configuration/service | consulter/lier |
| Hypothesis / Evidence | Investigate | raisonnement et qualification | consulter |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Persistence candidate | créer/annoter/supersede | Investigate concept | candidat, jamais confirmation automatique |
| Hypothesis context | préparer/lier | Investigate | CAP-INV-103 conserve le lifecycle |
| Evidence candidate context | préparer | Investigate | CAP-INV-107/108 requis |

## 11. Fonctionnalités
Regrouper comportements candidats de démarrage/relance/exécution; relier processus, fichiers et changements; afficher sources, temporalité et contradictions; comparer Runs; annoter et préparer Hypothesis/Evidence candidate.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter candidats | Analyst | Candidate set | 0 | sources accessibles | vue sourcée | non |
| Relier observations | Analyst | Relation | 2 | mêmes Run/context | relation attribuée | OPEN-013 |
| Annoter contradiction | Reviewer | Candidate | 2 | source visible | contradiction conservée | OPEN-013 |
| Préparer Hypothesis | Analyst | Hypothesis draft | 2 | contexte suffisant | CAP-INV-103 | OPEN-013 |
| Préparer Evidence candidate | Reviewer | Candidate package | 2 | provenance complète | CAP-INV-328 | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Regrouper candidats | oui | règles explicables | oui | regroupement suggéré | filtres/checklists |
| Rechercher contradictions | oui | comparaison multi-Run | oui | explication | diff déterministe |
| Proposer Hypothesis | oui | templates | oui | suggestion modifiable | formulaire humain |
| Confirmer persistance/Finding | humain | contrôles seulement | workflow de revue | jamais autonome | Evidence/Finding review |

## 14. États fonctionnels
`collecting`, `candidate`, `weakly-supported`, `contradicted`, `not-reproduced`, `partial`, `disputed`, `superseded`.

## 15. États d’interface
Loading conserve les sources; Empty distingue aucun candidat et données absentes; Partial nomme les lacunes; Error garde les candidats valides; Offline est read-only; Permission denied redacted.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Persistence candidate set | Analysis Result | analyste/Case | sources et contradictions visibles |
| Hypothesis context | draft package | CAP-INV-103 | aucune confirmation implicite |
| Evidence candidate context | draft package | CAP-INV-328 | qualification nécessaire |
| Comparison input | candidate snapshot | CAP-INV-325 | Run et préconditions |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-319/320 | regrouper comportements | CAP-INV-322 | processes, changes, timestamps, sources | source views |
| CAP-INV-322 | préparer raisonnement | CAP-INV-103 | candidat, supporting/contradicting observations | candidates |
| CAP-INV-322 | préparer qualification | CAP-INV-328 | candidats, Runtime Artifacts, contradictions, provenance | candidates |

## 18. Dépendances
CAP-INV-103/317/318/319/320/325/328; Evidence/Finding lifecycles; Shared Trace/Linking; OPEN-013/015.

## 19. Source de vérité
Investigate possède le candidat et son interprétation; mécanisme observé ≠ intention; candidat ≠ persistance confirmée; Finding reste CAP-INV-109.

## 20. Provenance et audit
Run, environment/profile, processes, changes, timestamps, sources, rules/Tools, contradictions, annotations, comparisons, dispositions et handoffs.

## 21. Permissions fonctionnelles
Behavioral results read; candidate annotate/link; Hypothesis prepare; Evidence candidate prepare; multi-Run comparison; cross-tenant denied.

## 22. Limites et erreurs
Source process/change manquante; temporalité incohérente; comportement non reproduit; Run partiel; classification ambiguë; permission révoquée.

## 23. Métriques
Candidats créés/contradictoires/non reproduits; Hypotheses préparées; Evidence candidates; dispositions humaines; provenance complète.

## 24. Classification de livraison
`defined` / `planned`; aucune méthode offensive, commande, moteur, instrumentation ou implémentation.

## 25. Critères d’acceptation
**Given** des changements associés à un mécanisme de démarrage
**When** l’analyste ouvre Persistence and Execution Mechanism Analysis
**Then** ils sont présentés comme candidats avec processus, fichiers, temporalité et contradictions

**Given** un second Run qui ne reproduit pas le comportement
**When** les Runs sont comparés
**Then** le candidat devient not-reproduced ou contradicted, pas confirmé

**Given** aucun modèle IA
**When** le workflow est réalisé
**Then** règles, comparateurs, templates et revue humaine couvrent les fonctions essentielles

## 26. Questions ouvertes
Persistence candidate reste un concept; OPEN-013 et OPEN-015 restent ouvertes.

## 27. Consommateurs documentaires
Dynamic Sandbox, Hypothesis Management, Multi-Run Comparison, Evidence/Finding handoff, Objets, Permissions et Journeys.
