---
id: CAP-CMD-002
title: Priority Management
product: command
module: mission-control
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-003
  - REQ-PROD-010
  - REQ-PROD-013
  - REQ-PROD-021
open_decisions:
  - OPEN-013
source-of-truth: canonical
---

# CAP-CMD-002 — Priority Management

## 1. Définition
Établit et explique la priorité opérationnelle effective d’un Incident ou d’une Task en séparant impact, urgence, SLA, criticité de service, dépendances et recommandations.

## 2. Problème utilisateur
Les équipes doivent arbitrer sans confondre severity, urgence, impact et priorité. Sans facteurs visibles et historique, la priorité devient arbitraire et vulnérable aux suggestions opaques.

## 3. Objectifs
Montrer la priorité effective et ses facteurs, permettre une modification justifiée, distinguer humain/déterministe/IA et conserver l’historique.

## 4. Non-objectifs
Ne définit pas de score universel, ne modifie pas la severity source, n’applique pas silencieusement une recommandation et ne crée pas une Decision Govern.

## 5. Propriétaire
Command / Mission Control / Command Product Lead pour la priorité opérationnelle des objets Command.

## 6. Utilisateurs
Principal : Incident Commander. Secondaires : SOC Analyst L2, Business Owner, Team Lead.

## 7. Conditions d’entrée
Incident ou Task accessible, facteurs présents ou explicitement incomplets et rôle autorisé à coordonner.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Priority factors | Incident, Task, SLA et Service | facteurs structurés | oui | valeurs courantes | facteurs manquants explicités |
| Current priority | Command | valeur effective | oui | version courante | mutation refusée |
| Recommendation | règle, moteur, workflow ou agent | proposition attribuée | non | run/version visible | aucun effet sur la valeur effective |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Incident / Task | Command | priority, impact, urgency, owner et état | consulter et modifier si autorisé |
| Service | Shared Business Service Catalog | criticité et dépendances | consulter et relier |
| SLA context | source de policy | risque, échéance et version | consulter et comparer |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Incident / Task | modifier priority et justification | Command | classe 2, concurrence optimiste et audit |
| Recommendation disposition | produire acceptation/rejet | Command audit | proposition distincte de la priorité effective |

## 11. Fonctionnalités
Comparer les facteurs, afficher l’origine de la priorité, prévisualiser l’effet, accepter/rejeter une proposition et rétablir une valeur par nouvelle mutation auditée.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspecter les facteurs | lecteur | Incident/Task | 0 | objet accessible | explication visible | non |
| Modifier la priorité | coordinateur | Incident/Task | 2 | justification et version courante | priorité effective | OPEN-013 |
| Accepter/rejeter une recommandation | coordinateur | proposition | 2 | facteurs et producteur visibles | disposition auditée | OPEN-013 pour acceptation |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Comparer les facteurs | oui | oui | oui | résumé attribué | facteurs sources et règles explicables |
| Proposer une priorité | oui | oui | oui | oui, proposition seulement | calcul déterministe ou analyse humaine |
| Modifier la priorité effective | oui | validation/version | workflow autorisé | jamais silencieusement | mutation humaine complète |

## 14. États fonctionnels
`effective`, `proposed`, `disputed`, `awaiting-owner`, `factors-incomplete`, `superseded`.

## 15. États d’interface
Partial nomme les facteurs manquants ; Offline et Stale bloquent les mutations non vérifiables ; Permission denied ne révèle rien ; Error conserve les valeurs valides.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Priority change | événement Incident/Task | Work Queue et Mission Control | acteur, facteurs, avant/après et justification |
| Recommendation disposition | événement d’audit | producteur et analystes | acceptation/rejet distinct de la proposition |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Priority Management | facteurs nécessitant investigation | Investigate | Incident et facteurs incomplets | Incident restauré |
| Priority Management | action nécessitant autorité | Govern | Incident, impact, urgence et action proposée | aucune Decision locale |

## 18. Dépendances
CAP-CMD-104, CAP-CMD-105, CAP-CMD-204, CAP-CMD-205, Metrics Engine et audit hooks ; `OPEN-013` reste ouverte.

## 19. Source de vérité
Priorité effective et justification : Command. Service et SLA : projections sourcées. Recommandations : producteur/version/facteurs visibles.

## 20. Provenance et audit
Acteur/producteur, version/run, facteurs, avant/après, justification, disposition, tenant et correlation ID.

## 21. Permissions fonctionnelles
`perm.command.coordinate`, `perm.command.incident.manage`, `perm.command.task.manage`, ABAC ownership/tenant ; step-up et atomisation reportés.

## 22. Limites et erreurs
Facteurs absents, stale, conflit de version, tenant incompatible ou permission refusée empêchent une priorité présentée comme complète.

## 23. Métriques
Changements justifiés, délai des disputes et dispositions des recommandations ; aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native ; preuve documentaire uniquement.

## 25. Critères d’acceptation
**Given** des facteurs courants, **When** un coordinateur modifie la priorité, **Then** valeur, facteurs, acteur et justification sont auditables.

**Given** une recommandation différente, **When** elle est consultée, **Then** la priorité effective reste inchangée avant acceptation autorisée.

**Given** aucun modèle IA, **When** la capability est utilisée, **Then** règles, moteur déterministe et actions manuelles suffisent.

## 26. Questions ouvertes
Quels facteurs sont obligatoires et quelles actions de classe 2 exigent Govern ? — Requirement IDs ci-dessus ; `OPEN-013` reste ouverte.

## 27. Consommateurs documentaires
Mission Control Priorities, Unified Work Queue, Incident Detail, parcours Phase 5, écrans Phase 6, objets Phase 7 et permissions/contrats ultérieurs.