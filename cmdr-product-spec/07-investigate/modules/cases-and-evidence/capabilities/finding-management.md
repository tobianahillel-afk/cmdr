---
id: CAP-INV-109
title: Finding Management
product: investigate
module: cases-and-evidence
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
---
# CAP-INV-109 — Finding Management

## 1. Définition
Créer, revoir, confirmer, contester, supersede ou retirer un Finding soutenu par des Evidence, sans le confondre avec Hypothesis, Decision ou Result.

## 2. Problème utilisateur
Les conclusions peuvent être publiées sans Evidence contradictoire, auteur ou revue. Un résumé IA peut être pris pour un Finding confirmé.

## 3. Objectifs
Documenter assertion, Evidence pour/contre, auteur, revue et statut ; gérer draft/proposed/under-review/confirmed/disputed/superseded/withdrawn ; relier Incident, Action Request et Report.

## 4. Non-objectifs
Ne pas créer Decision, représenter Result, confirmer automatiquement ou définir un score universel.

## 5. Propriétaire
Investigate / Cases and Evidence / Investigate Product Lead. Finding est une conclusion canonique Investigate.

## 6. Utilisateurs
Principal : Case Analyst. Reviewer : Finding Reviewer/Lead. Consommateurs : Command, Govern, Reporting. Automation peut seulement proposer.

## 7. Conditions d’entrée
Case et assertion, au moins une Evidence ou raison explicite pour draft, auteur/producteur identifié et permission de review pour confirmation.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Assertion | analyst/automation | claim | oui | versionnée | rester draft |
| Supporting Evidence | Investigate | support | oui pour confirm | qualification courante | confirmation interdite |
| Contradictory Evidence | Investigate | counterevidence | non mais recherchée | statut courant | none/unknown explicite |
| Case/Incident context | Investigate/Command | scope/impact | oui | fraîcheur visible | Case-only draft |
| Review method | policy/workflow | validation | oui pour confirm | version courante | under-review |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Evidence | Investigate | qualification/source | lire/lier |
| Hypothesis | Investigate | lineage/statut | lire/lier |
| Case | Investigate | scope | lire |
| Incident | Command | impact/priority context | consulter |
| Decision/Result | Govern | downstream/outcome | consulter |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Finding | créer, versionner, statut, supersede, withdraw | Investigate | confirmation exige reviewer autorisé |
| Finding–Evidence relation | créer/qualifier | Investigate | pour/contre explicite |
| Incident relation | créer | cross-product relation | n’altère pas Incident |
| Action Request input | préparer | Govern lifecycle | Finding reste Investigate |

## 11. Fonctionnalités
Créer draft/proposal, définir assertion/scope, lier Evidence favorable/contradictoire, afficher auteur/méthode, revoir/confirmer/contester, supersede/withdraw et alimenter Action Request/Report.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Créer/éditer draft | Case Analyst | Finding | 2 | Case access | version draft | OPEN-013 |
| Soumettre review | Analyst | Finding | 2 | Evidence links | under-review | OPEN-013 |
| Confirmer/contester | Reviewer | Finding | 2 | permission review | status/reason | OPEN-013 |
| Supersede/withdraw | Lead | Finding | 2 | reason/replacement | lineage | OPEN-013 |
| Préparer request | Analyst | request draft | 3/4 context | Finding allowed | Govern draft | Govern |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Draft assertion | oui | templates/rules | workflow | oui | auteur humain |
| Vérifier Evidence refs | oui | oui | oui | explication | validation déterministe |
| Suggest confidence | oui | oui | oui | oui | reviewer voit facteurs |
| Confirm Finding | reviewer | non autonome | workflow revue | non autonome | reviewer humain |

## 14. États fonctionnels
`draft`, `proposed`, `under-review`, `confirmed`, `disputed`, `superseded`, `withdrawn`. États Draft.

## 15. États d’interface
Loading garde draft ; Partial nomme Evidence absentes ; Error conserve version ; Offline lecture ; Permission denied sans fuite ; Stale signale Evidence superseded.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Finding | objet | Case/Command/Govern/Reporting | assertion, auteur, Evidence, statut |
| Review event | Timeline/Trace | Case/Audit | reviewer, méthode, raison |
| Action Request input | contexte structuré | CAP-INV-113 | ownership sources conservé |
| Report content | bloc cité | CAP-INV-114 | citations/redactions |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Hypothesis | draft conclusion | Finding Management | claim, lineage, Evidence pour/contre | Hypothesis |
| Finding | submit review | Review | Evidence, method, uncertainty | Finding |
| Finding confirmed | prepare request | Action Request Preparation | Finding/Evidence, target, impact | Case |
| Govern Result | link outcome | Finding/Case | Request, Decision, Run, Result | Finding |

## 18. Dépendances
Evidence Creation/Review, Hypothesis, Case, Govern Action Request, Reporting Engine, review permissions et OPEN-013.

## 19. Source de vérité
Finding/review restent Investigate ; Evidence Investigate ; Incident Command ; Action Request/Decision/Result Govern ; Reporting Engine Shared.

## 20. Provenance et audit
Auteur/producteur, assertion versions, Evidence pour/contre, review method/reviewer, reasons, request/result links et agent proposal metadata.

## 21. Permissions fonctionnelles
Finding create/edit, review/confirm/dispute, supersede/withdraw, sensitive Evidence, submit request et séparation des tâches future.

## 22. Limites et erreurs
Missing Evidence, contradiction cachée, reviewer conflict, concurrent version, Evidence superseded, Case closed, permission revoked ou agent provenance absente.

## 23. Métriques
Confirmed Findings avec Evidence pour/contre, review time, dispute/supersession, agent proposal disposition et requests returned.

## 24. Classification de livraison
`defined` / `planned`, cible native. Promotion conditionnée par rôles review, dependency checks, versioning et transition Govern.

## 25. Critères d’acceptation
**Given** une proposal agentique **When** ouverte **Then** statut `proposed`, agent/run visible et review obligatoire.

**Given** une Evidence contradictoire **When** review **Then** elle reste visible et le statut est raisonné.

**Given** un Finding confirmé et action classe 3 **When** request soumise **Then** Finding/Evidence restent Investigate et aucune Decision locale.

## 26. Questions ouvertes
Modèle de confiance/validation, séparation auteur/reviewer et OPEN-013 restent ouverts.

## 27. Consommateurs documentaires
Case, Command projections, Govern, Reporting, Replay/Readiness et phases Objets/Permissions.
