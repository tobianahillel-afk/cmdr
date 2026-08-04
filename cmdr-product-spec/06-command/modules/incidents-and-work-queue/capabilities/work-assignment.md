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
Un item peut sembler affecté sans acceptation ou être réassigné pendant une action. Rôles : Incident Commander, SOC Analyst, Team Lead, Task owner. Sans la capacité, le travail reste sans responsable effectif et le SLA dérive.

## 3. Objectifs
Séparer proposition, affectation et acceptation ; expliquer toute proposition ; gérer indisponibilité et conflit ; conserver l’historique sans modifier l’ownership canonique du type d’objet.

## 4. Non-objectifs
Ne pas administrer les identités, créer une politique RH, transférer l’ownership produit ou permettre à un agent de s’accorder des droits.

## 5. Propriétaire
Command / Incidents and Work Queue / Command Product Lead pour l’affectation opérationnelle des objets Command.

## 6. Utilisateurs
Principal : Incident Commander. Secondaires : SOC Analyst, Team Lead, Task owner. Les projections d’identité restent Platform Settings.

## 7. Conditions d’entrée
Work item accessible, candidat résolu, permission de coordination et version courante.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Work item | Incident/Task | objet | oui | version courante | refuser mutation |
| Candidate | identity projection | principal/équipe | oui | disponibilité connue ou unknown | confirmation explicite |
| Rationale | humain/règle/moteur/agent | raison/provenance | oui si proposition | run/version | proposition non effective |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Incident / Task | Command | assignment, owner, state | lecture/modification |
| Principal / Role | Platform Settings | identité, équipe, disponibilité | projection |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Incident / Task | assignment et acceptance state | Command | classe 2, version et audit |
| Notification | demande de prise en charge | Shared | ne vaut pas acceptation |

## 11. Fonctionnalités
Affecter/réaffecter ; prendre/libérer ; proposer sans appliquer ; expliquer contraintes ; détecter et résoudre un conflit de version.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Affecter | coordinateur | Incident/Task | 2 | candidat autorisé | assigned | OPEN-013 |
| Prendre en charge | analyste | Incident/Task | 2 | assignment compatible | accepted | OPEN-013 |
| Libérer | owner actuel | Incident/Task | 2 | raison | unassigned/team-owned | OPEN-013 |
| Réaffecter | coordinateur | Incident/Task | 2 | version courante | pending/assigned | OPEN-013 |
| Accepter proposition | coordinateur | proposal | 2 | provenance visible | mutation humaine auditée | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Humain | Règle | Moteur déterministe | Workflow | Agent | Govern | Alternative sans IA |
|---|---|---|---|---|---|---|---|
| Affecter | oui | oui si policy | validation/conflit | possible | proposition seulement | OPEN-013 | sélection manuelle |
| Expliquer candidat | correction | critères explicables | score/facteurs visibles | possible | résumé attribué | non | données équipe/charge |
Aucune proposition ne devient effective sans contrat autorisé et trace.

## 14. États fonctionnels
`unassigned`, `assigned`, `accepted`, `declined`, `reassignment-pending`, `unavailable`, `conflict`.

## 15. États d’interface
Loading conserve l’item ; Empty explique l’absence de candidat ; Partial montre données manquantes ; Error garde version valide ; Offline bloque mutation ; Permission denied ne divulgue rien ; Stale force refresh. Rendu : Design System.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Assignment update | Incident/Task event | Queue, notifications, handover | before/after, actor, reason |
| Proposal disposition | audit event | producer/analystes | accepted/rejected/expired |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Assignment | candidat indisponible | Escalation | item, raison, délai | file |
| Accepted | transfert de handover | Handover | items/destinataire | acknowledgement séparé |

## 18. Dépendances
Identity projection, Notification Center, Collaboration Service, CAP-CMD-103 et CAP-CMD-110. Aucun ownership externe n’est transféré.

## 19. Source de vérité
Assignment effectif et historique sont Command ; identité/disponibilité sont des projections ; recommendation conserve producteur, version/run et facteurs.

## 20. Provenance et audit
Chaque proposition/mutation enregistre acteur/producteur, source, version, before/after, raison, résultat, tenant, environnement et correlation ID.

## 21. Permissions fonctionnelles
`perm.command.coordinate`, `perm.command.incident.manage`, `perm.command.task.manage`, ABAC team/tenant. Granularité et `OPEN-013` restent ouvertes.

## 22. Limites et erreurs
Candidat inaccessible, disponibilité inconnue, conflit de version, changement de tenant/environnement, item stale ou permission refusée empêchent une mutation présentée comme réussie. Le draft est récupérable.

## 23. Métriques
Temps item→assignment accepté ; proposals acceptées/rejetées ; conflits ; libérations sans remplacement. Aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native. Preuve documentaire seulement ; promotion après objets, permissions, parcours, écrans, contrats, implémentation et validation.

## 25. Critères d’acceptation
**Given** un item unassigned et un candidat autorisé, **When** le coordinateur l’affecte, **Then** assignment, acteur, raison et version sont auditables et l’acceptation reste distincte.

**Given** une proposition IA, **When** elle est consultée, **Then** elle reste non effective, facteurs et producer sont visibles et l’utilisateur peut accepter/rejeter.

**Given** aucun modèle IA ou un conflit de version, **When** l’affectation est tentée, **Then** la voie manuelle fonctionne ou le conflit bloque proprement sans écraser l’autre mutation.

## 26. Questions ouvertes
L’acceptation est-elle obligatoire par tenant ? Quelle donnée d’indisponibilité est légitime sans devenir RH ? — REQ-PROD-003, REQ-PROD-013, REQ-PROD-021, REQ-SEC-001. `OPEN-013` reste ouverte.

## 27. Consommateurs documentaires
Unified Work Queue, Incident Detail, Mission Control, Handover, parcours Phase 5, écrans Phase 6, objets Phase 7 et permissions/contrats ultérieurs.
