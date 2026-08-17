---
id: CAP-INV-103
title: Hypothesis Management
product: investigate
module: cases-and-evidence
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-014
  - REQ-AI-002
open_decisions:
  - OPEN-013
---
# CAP-INV-103 — Hypothesis Management

## 1. Définition
Créer, tester, comparer et supersede des Hypotheses explicites sans les confondre avec des Findings.

## 2. Problème utilisateur
Questions, intuitions et conclusions sont souvent mélangées. Une Hypothesis proposée par IA peut être prise à tort pour une conclusion.

## 3. Objectifs
Définir question, auteur, source, critères de confirmation/infirmation, Evidence pour/contre, statut, comparaison, rejet, supersession et historique.

## 4. Non-objectifs
Ne pas confirmer un Finding, créer une Decision, automatiser la conclusion ou figer la machine d’état.

## 5. Propriétaire
Investigate / Cases and Evidence / Investigate Product Lead. Hypothesis est un objet Investigate.

## 6. Utilisateurs
Principal : Case Analyst. Secondaires : Threat Hunter, Investigation Lead reviewer et automation productrice de propositions.

## 7. Conditions d’entrée
Case accessible, question testable, auteur/producteur identifié et critères suffisamment explicites.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Question/claim | humain ou automation | proposition testable | oui | version courante | rester proposed/draft |
| Source | Signal, Search, Artifact ou analyst | origine | oui | référence résoluble | provenance incomplete |
| Evidence/résultats | Investigate/Shared | supports ou contradictions | non | statut/fraîcheur visibles | testing/inconclusive |
| Critères | analyst/checklist | support/refutation | oui pour revue | versionnés | pas de statut supported |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | scope/objectif | consulter |
| Evidence | Investigate | qualification/source | consulter et relier |
| Event/Search result | Shared Capabilities | faits candidats | consulter sans qualifier |
| Finding | Investigate | conclusions liées | consulter |
| Automation Run | CMDR Studio | provenance proposal | consulter |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Hypothesis | créer, modifier, statut, supersede | Investigate | versions conservées |
| Hypothesis–Evidence relation | créer/qualifier | Investigate | n’altère pas Evidence |
| Case timeline event | émettre | Shared mechanism | sémantique Investigate |
| Finding | aucune confirmation | Investigate | transition séparée CAP-INV-109 |

## 11. Fonctionnalités
Créer Hypothesis, définir critères, lier Evidence pour/contre, comparer, changer états, supersede, conserver historique et recevoir propositions IA attribuées.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Créer/modifier | Analyst | Hypothesis | 2 | Case accessible | version | OPEN-013 |
| Lier Evidence | Analyst | relation | 2 | Evidence accessible | relation qualifiée | OPEN-013 |
| Changer statut | Reviewer | Hypothesis | 2 | raison/éléments | transition | OPEN-013 |
| Supersede | Reviewer | Hypothesis | 2 | remplaçante | lineage | OPEN-013 |
| Comparer | Analyst | Hypotheses | 0 | plusieurs | comparaison | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Proposer Hypothesis | oui | oui | oui | oui | templates et création manuelle |
| Associer résultats | oui | oui | oui | oui | règles explicites |
| Vérifier critères | oui | oui | oui | explication | checklist |
| Changer statut final | reviewer | non autonome | workflow revue | recommandation | décision humaine |

## 14. États fonctionnels
`proposed`, `testing`, `supported`, `contradicted`, `inconclusive`, `rejected`, `superseded`. Machine finale reportée.

## 15. États d’interface
Loading garde la version ; Partial nomme Evidence manquantes ; Error conserve liens ; Offline lecture ; Permission denied sans fuite ; Stale signale les sources dépassées.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Hypothesis version | objet | Case Workspace | attribuée, versionnée, testable |
| Support map | relations | Finding/Replay | pour/contre explicites |
| Status event | Timeline | Case/Review | raison et auteur |
| Finding draft input | contexte | CAP-INV-109 | pas un Finding confirmé |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Search/Hunt | proposition | Hypothesis | question, sources, Query/run, auteur | source |
| Hypothesis | Evidence liée | review | relation, relevance, qualification | Case |
| Hypothesis supported | draft Finding | Finding Management | claim, Evidence pour/contre, uncertainty | Hypothesis |
| Hypothesis superseded | ouvrir remplaçante | Hypothesis | lineage/raisons | originale |

## 18. Dépendances
Case Lifecycle, Evidence, Finding, Search/Hunt provenance, Studio optional et OPEN-013.

## 19. Source de vérité
Hypothesis et statut restent Investigate ; Events/results restent Shared ; Evidence distincte ; proposal agentique reste proposal.

## 20. Provenance et audit
Auteur/producteur, source, critères versions, Evidence/results, changements de statut, agent/version/Tool Calls et disposition humaine.

## 21. Permissions fonctionnelles
Hypothesis create/update, status review, Evidence link, sensitive Case et agent proposal use.

## 22. Limites et erreurs
Source inaccessible, critères manquants, Evidence refusée, conflit de statut, cycle de supersession, Case stale ou provenance agent absente.

## 23. Métriques
Hypotheses avec critères, délai de disposition, Evidence contradictoires visibles, proposals IA acceptées/modifiées/rejetées et Findings avec lineage.

## 24. Classification de livraison
`defined` / `planned`, cible native. Promotion conditionnée par objet versionné, relations et revue humaine.

## 25. Critères d’acceptation
**Given** une proposal agentique **When** elle est ouverte **Then** agent/version/Tool Calls sont visibles, statut `proposed` et aucun Finding confirmé.

**Given** une Evidence contradictoire **When** elle est liée **Then** la contradiction est visible et le statut n’est pas changé automatiquement.

**Given** une remplaçante **When** supersession **Then** lineage et historique restent lisibles.

## 26. Questions ouvertes
Machine finale, confiance et OPEN-013 restent ouverts.

## 27. Consommateurs documentaires
Case Workspace, Evidence Review, Finding, Replay, IA et phases Objets/Permissions.
