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
Owner, assignee, team, watcher et product owner sont souvent confondus. Sans séparation, responsabilités et droits deviennent ambigus.

## 3. Objectifs
Séparer ownership opérationnel/canonique, rendre responsabilités visibles, permettre contributeurs/watchers et conserver les périodes d’ownership.

## 4. Non-objectifs
Ne définit pas RBAC, ne fait pas d’un watcher un owner, ne transfère pas Case/Decision à Command et ne remplace pas Work Assignment.

## 5. Propriétaire
Command / Incidents and Work Queue / Command Product Lead pour les relations opérationnelles sur Incident/Task.

## 6. Utilisateurs
Principal : Incident Commander. Secondaires : SOC Analyst, Team Lead, Auditor.

## 7. Conditions d’entrée
Incident/Task accessible, identités résolues et permission lecture/gestion selon l’action.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Canonical ownership | Ownership Register | owner produit | oui | dernière révision | `unknown` et revendication locale bloquée |
| Operational ownership | Incident/Task | owner, équipe, contributeurs et watchers | oui | version courante | `unassigned` |
| Identity context | Platform Settings | principals et teams | oui | résolution courante | ID conservé, état `identity-unresolved` |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Incident / Task | Command | owner, équipe, contributeurs et watchers | consulter et modifier selon relation |
| Case / Decision / Response Run | produit propriétaire | owner projeté | consulter en lecture seule |
| Principal / Team | Platform Settings | identité et statut | consulter et sélectionner |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Incident / Task | définir owner/équipe/contributeur/watcher | Command | classe 2, permission par relation |
| Ownership event | ajouter période et avant/après | Shared Audit | append-only conceptuel |
| Objet externe | aucune mutation | produit propriétaire | ownership projeté seulement |

## 11. Fonctionnalités
Afficher owner produit, gérer owner principal/équipe, ajouter/retirer contributeurs/watchers, montrer l’approbateur Govern distinct et retracer les changements.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Définir owner principal | coordinateur | Incident/Task | 2 | principal autorisé | owner mis à jour | OPEN-013 |
| Ajouter contributeur | owner/coordinateur | Incident/Task | 2 | principal autorisé | relation ajoutée | OPEN-013 |
| Suivre comme watcher | lecteur | Incident/Task | 2 | lecture autorisée | watcher ajouté | non |
| Consulter owner canonique | lecteur | type objet | 0 | registre disponible | owner produit visible | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Résoudre l’owner canonique | oui | oui | oui | non nécessaire | Ownership Register |
| Proposer un owner opérationnel | oui | facteurs charge/compétence | oui | proposition attribuée | recherche et sélection manuelles |
| Modifier une relation d’ownership | oui | validation/version | workflow possible | jamais silencieusement | action humaine complète |
| Afficher l’historique | oui | oui | oui | résumé facultatif | événements d’audit |

## 14. États fonctionnels
`owner-set`, `team-owned`, `unassigned`, `shared-contribution`, `watcher-only`, `ownership-conflict`, `identity-unresolved`.

## 15. États d’interface
Empty signifie unassigned ; Partial nomme l’identité manquante ; Offline bloque mutation ; Permission denied ne révèle rien ; Stale impose validation.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Ownership projection | métadonnée objet | Queue, Inspector et Handover | owner opérationnel et owner produit séparés |
| Ownership event | événement d’audit | historique et métriques | acteur, relation, période et before/after |
| Unresolved identity state | état fonctionnel | coordinateur | ID stable conservé sans donnée inventée |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Operational Ownership | ouverture owner externe | produit propriétaire | référence objet et owner projeté | objet Command restauré |
| Ownership conflict | escalade | Command coordination | item, versions et acteurs | résolution auditée |
| Identity unresolved | inspection | Platform Settings | principal ID et return origin | même relation restaurée |

## 18. Dépendances
CAP-CMD-102, Ownership Register, Identity projections, Object Linking Service et audit hooks.

## 19. Source de vérité
Relations Incident/Task : Command. Owner canonique : registre/produit source. Identités : Platform Settings.

## 20. Provenance et audit
Acteur/producteur, relation, before/after, période, justification, tenant, version et correlation ID.

## 21. Permissions fonctionnelles
`perm.command.read`, `perm.command.coordinate`, `perm.command.incident.manage`, `perm.command.task.manage` ; granularité reportée.

## 22. Limites et erreurs
Identité unresolved, conflit, owner externe inaccessible, tenant incompatible ou refus ne transfèrent jamais implicitement l’ownership.

## 23. Métriques
Items avec owner explicite, durée unassigned et conflits d’ownership ; aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native ; preuve documentaire seulement.

## 25. Critères d’acceptation
**Given** un Incident sans owner, **When** un principal est défini, **Then** owner opérationnel, owner produit et audit restent distincts.

**Given** une proposition IA, **When** elle est affichée, **Then** elle n’est pas effective sans acceptation autorisée.

**Given** une identité inaccessible, **When** l’objet est ouvert, **Then** l’état unresolved est visible sans donnée inventée.

## 26. Questions ouvertes
Qui peut gérer contributeurs/watchers et comment représenter les équipes transverses ? — Requirement IDs ci-dessus ; `OPEN-013` reste ouverte.

## 27. Consommateurs documentaires
Work Queue, Inspector, Handover, parcours Phase 5, écrans Phase 6, objets Phase 7 et permissions ultérieures.