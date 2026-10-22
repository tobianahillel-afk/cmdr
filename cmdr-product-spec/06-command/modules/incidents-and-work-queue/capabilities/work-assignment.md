---
id: CAP-CMD-102
title: Work Assignment
product: command
module: incidents-and-work-queue
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-003
  - REQ-PROD-013
  - REQ-PROD-021
  - REQ-SEC-001
open_decisions:
  - OPEN-013
source-of-truth: canonical
---

# CAP-CMD-102 — Work Assignment

## 1. Définition
Gère l’affectation explicite d’un Incident ou d’une Task à une équipe et/ou un utilisateur, avec prise en charge, libération, réaffectation, proposition et conflits.

## 2. Problème utilisateur
Un item peut sembler affecté sans acceptation ou être réassigné pendant une action. Sans distinction, le travail reste sans responsable effectif et le SLA dérive.

## 3. Objectifs
Séparer proposition/affectation/acceptation, expliquer les propositions, gérer indisponibilité/conflit et conserver l’historique sans changer l’ownership canonique.

## 4. Non-objectifs
N’administre pas les identités, ne crée pas de policy RH, ne transfère pas l’ownership produit et n’autorise pas un agent à s’accorder des droits.

## 5. Propriétaire
Command / Incidents and Work Queue / Command Product Lead pour l’affectation opérationnelle des objets Command.

## 6. Utilisateurs
Principal : Incident Commander. Secondaires : SOC Analyst, Team Lead, Task owner.

## 7. Conditions d’entrée
Work item accessible, candidat résolu, permission de coordination et version courante.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Work item | Incident ou Task | objet Command | oui | version courante | mutation refusée |
| Candidate | Platform Settings identity projection | principal ou équipe | oui | disponibilité courante ou `unknown` | confirmation explicite requise |
| Rationale | humain, règle, moteur ou agent | raison et provenance | oui pour proposition | run/version visible | proposition non effective |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Incident / Task | Command | assignment, owner et état | consulter et modifier si autorisé |
| Principal / Role | Platform Settings | identité, équipe et disponibilité | consulter et sélectionner en projection |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Incident / Task | affecter, accepter, libérer ou réaffecter | Command | classe 2, version et audit |
| Notification | demander une prise en charge | Shared Capabilities | ne vaut jamais acceptation |
| Identity projection | aucune mutation | Platform Settings | lecture seule |

## 11. Fonctionnalités
Affecter/réaffecter, prendre/libérer, proposer sans appliquer, expliquer les contraintes et détecter/résoudre un conflit de version.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Affecter | coordinateur | Incident/Task | 2 | candidat autorisé | assigned | OPEN-013 |
| Prendre en charge | analyste | Incident/Task | 2 | assignment compatible | accepted | OPEN-013 |
| Libérer ou réaffecter | owner/coordinateur | Incident/Task | 2 | raison et version | nouvel état d’affectation | OPEN-013 |
| Accepter une proposition | coordinateur | proposition | 2 | facteurs visibles | mutation humaine auditée | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Rechercher un candidat | oui | oui | oui | suggestion attribuée | sélection d’équipe/utilisateur et filtres |
| Proposer une affectation | oui | facteurs explicables | oui | oui, non effective | règles de charge/compétence et choix humain |
| Appliquer l’affectation | oui | validation/version | workflow autorisé | jamais silencieusement | action manuelle complète |
| Détecter un conflit | résolution humaine | oui | oui | explication possible | contrôle de version et refresh |

## 14. États fonctionnels
`unassigned`, `assigned`, `accepted`, `declined`, `reassignment-pending`, `unavailable`, `conflict`.

## 15. États d’interface
Empty explique l’absence de candidat ; Partial montre les données manquantes ; Offline bloque mutation ; Permission denied ne divulgue rien ; Stale impose refresh.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Assignment update | événement Incident/Task | Queue, Notifications et Handover | before/after, acteur et raison conservés |
| Proposal disposition | événement d’audit | producteur et analystes | accepted/rejected/expired distinct de l’affectation |
| Assignment conflict | événement | coordinateur | versions concurrentes visibles, aucun écrasement |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Work Assignment | candidat indisponible | Escalation | item, raison, délai et owner courant | même item dans la file |
| Work Assignment | relève acceptée | Handover | items, destinataire et état d’acceptation | acknowledgement séparé |
| Work Assignment | identité sélectionnée | Platform Settings profile | principal, équipe et return origin | item restauré |

## 18. Dépendances
Identity projection, Notification Center, Collaboration Service, CAP-CMD-103 et CAP-CMD-110.

## 19. Source de vérité
Assignment effectif/historique : Command. Identité/disponibilité : Platform Settings. Recommandation : producteur/version/facteurs visibles.

## 20. Provenance et audit
Proposition ou mutation : acteur/producteur, source, version, before/after, raison, résultat, tenant et correlation ID.

## 21. Permissions fonctionnelles
`perm.command.coordinate`, `perm.command.incident.manage`, `perm.command.task.manage`, ABAC team/tenant ; granularité et `OPEN-013` restent ouvertes.

## 22. Limites et erreurs
Candidat inaccessible, disponibilité inconnue, conflit, tenant incompatible, item stale ou permission refusée empêchent une mutation réussie fictive.

## 23. Métriques
Temps item→assignment accepté, propositions acceptées/rejetées, conflits et libérations sans remplacement ; aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native ; preuve documentaire seulement.

## 25. Critères d’acceptation
**Given** un item unassigned, **When** il est affecté, **Then** assignment, acteur, raison et acceptation restent distincts et auditables.

**Given** une proposition IA, **When** elle est consultée, **Then** elle reste non effective et peut être acceptée/rejetée.

**Given** un conflit ou aucun modèle, **When** l’affectation est tentée, **Then** le conflit bloque proprement ou la sélection manuelle reste disponible.

## 26. Questions ouvertes
L’acceptation est-elle obligatoire par tenant et quelle donnée d’indisponibilité est légitime ? — Requirement IDs ci-dessus ; `OPEN-013` reste ouverte.

## 27. Consommateurs documentaires
Unified Work Queue, Incident Detail, Mission Control, Handover, parcours Phase 5, écrans Phase 6, objets Phase 7 et permissions ultérieures.