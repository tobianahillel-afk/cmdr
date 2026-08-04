---
id: CAP-CMD-103
title: Operational Ownership
product: command
module: incidents-and-work-queue
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-006
  - REQ-PROD-009
  - REQ-PROD-013
  - REQ-OBJ-001
open_decisions:
  - OPEN-013
source-of-truth: canonical
---
# CAP-CMD-103 — Operational Ownership

## 1. Définition
Rend explicites owner principal, équipe, contributeurs, watchers et approbateur éventuel sans modifier le propriétaire canonique du type d’objet.

## 2. Problème utilisateur
Owner, assignee, team, watcher et product owner sont souvent confondus. Rôles : Incident Commander, SOC Analyst, Team Lead, Auditor. Sans séparation, la responsabilité et les droits deviennent ambigus.

## 3. Objectifs
Séparer ownership opérationnel/canonique ; rendre responsabilités et droits visibles ; permettre contributeurs/watchers sans diluer l’owner ; conserver les périodes d’ownership.

## 4. Non-objectifs
Ne pas définir RBAC, faire d’un watcher un owner, transférer Case/Decision à Command ou remplacer Work Assignment.

## 5. Propriétaire
Command / Incidents and Work Queue / Command Product Lead pour les relations opérationnelles sur Incident/Task.

## 6. Utilisateurs
Principal : Incident Commander. Secondaires : SOC Analyst, Team Lead, Auditor. Identités et rôles restent des projections Settings.

## 7. Conditions d’entrée
Incident/Task accessible ; identités résolues ; permission lecture/gestion selon l’action.

## 8. Entrées fonctionnelles
| Entrée | Source | Type | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Canonical ownership | Ownership Register | produit owner | oui | révision courante | unknown et revendication locale bloquée |
| Operational ownership | Incident/Task | owner/team/contributors/watchers | oui | version courante | unassigned |
| Identity context | Platform Settings | principals/teams | oui | courant | ID conservé, unresolved |

## 9. Objets lus
| Objet | Owner | Projection | Droit local |
|---|---|---|---|
| Incident / Task | Command | owner, team, contributors, watchers | lecture/modification |
| Case / Decision / Run | produit propriétaire | owner projeté | lecture |
| Principal / Team | Settings | identité/statut | projection |

## 10. Objets créés ou modifiés
| Objet | Opération | Owner | Règle |
|---|---|---|---|
| Incident / Task | relations owner/team/contributor/watcher | Command | classe 2, permission par relation |
| Audit | période d’ownership | Shared audit | append-only conceptuel |

## 11. Fonctionnalités
Afficher owner produit ; gérer owner principal/équipe ; ajouter/retirer contributeurs et watchers ; montrer l’approbateur Govern distinct ; retracer les changements.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Définir owner principal | coordinateur | Incident/Task | 2 | principal autorisé | owner mis à jour | OPEN-013 |
| Ajouter contributeur | owner/coordinateur | Incident/Task | 2 | principal autorisé | relation ajoutée | OPEN-013 |
| Suivre comme watcher | lecteur | Incident/Task | 2 | lecture autorisée | watcher ajouté | non |
| Consulter owner canonique | lecteur | type objet | 0 | registre disponible | owner produit visible | non |

## 13. Automatisation et IA
| Fonction | Humain | Règle | Moteur | Workflow | Agent | Govern | Sans IA |
|---|---|---|---|---|---|---|---|
| Proposer owner | décision humaine | possible | facteurs charge/compétence | possible | proposition attribuée | OPEN-013 | sélection manuelle |
| Afficher ownership | oui | oui | résolution déterministe | possible | résumé facultatif | non | registres et identités |
Aucun agent ne change un owner silencieusement.

## 14. États fonctionnels
`owner-set`, `team-owned`, `unassigned`, `shared-contribution`, `watcher-only`, `ownership-conflict`, `identity-unresolved`.

## 15. États d’interface
Loading conserve l’objet ; Empty signifie unassigned ; Partial nomme l’identité manquante ; Error garde la version valide ; Offline bloque mutation ; Permission denied ne révèle rien ; Stale force validation. Rendu DS.

## 16. Sorties
| Sortie | Objet/événement | Consommateur | Garantie |
|---|---|---|---|
| Ownership projection | metadata objet | Queue, Inspector, handover | owner opérationnel et owner produit séparés |
| Ownership event | audit | metrics/history | acteur, rôle, période, before/after |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte | Retour |
|---|---|---|---|---|
| Ownership | objet externe | produit owner | référence/owner projeté | aucune mutation locale |
| Conflict | escalade | Command coordination | item, versions, actors | résolution auditée |

## 18. Dépendances
CAP-CMD-102, Ownership Register, Identity projections, Object Linking Service, audit hooks.

## 19. Source de vérité
Relations opérationnelles Incident/Task : Command. Owner canonique : registres/produit source. Identités : Settings.

## 20. Provenance et audit
Acteur/producteur, relation, before/after, période, justification, tenant, environnement, version et correlation ID.

## 21. Permissions fonctionnelles
`perm.command.read`, `perm.command.coordinate`, `perm.command.incident.manage`, `perm.command.task.manage`. Granularité reportée.

## 22. Limites et erreurs
Identité unresolved, conflit de version, owner externe inaccessible, tenant incompatible ou permission refusée ne doivent jamais transférer implicitement l’ownership.

## 23. Métriques
Items avec owner explicite ; durée unassigned ; conflits d’ownership. Aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native ; preuve documentaire seulement.

## 25. Critères d’acceptation
**Given** un Incident sans owner, **When** un coordinateur définit un principal, **Then** owner opérationnel, owner produit et audit restent distincts.

**Given** une proposition IA, **When** elle est affichée, **Then** elle n’est pas effective sans acceptation humaine autorisée.

**Given** une identité inaccessible, **When** l’objet est ouvert, **Then** l’ID reste visible selon permission, l’état est unresolved et aucune donnée n’est inventée.

## 26. Questions ouvertes
Qui peut gérer contributeurs/watchers sans gérer owner ? Comment représenter équipes transverses ? — REQ-PROD-006, REQ-PROD-009, REQ-PROD-013, REQ-OBJ-001. `OPEN-013` reste ouverte.

## 27. Consommateurs documentaires
Work Queue, Inspector, Handover, parcours Phase 5, écrans Phase 6, objets Phase 7, permissions ultérieures.
