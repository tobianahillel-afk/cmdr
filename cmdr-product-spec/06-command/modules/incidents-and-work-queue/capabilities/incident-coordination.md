---
id: CAP-CMD-106
title: Incident Coordination
product: command
module: incidents-and-work-queue
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-OBJ-001
  - REQ-PROD-003
  - REQ-PROD-008
  - REQ-PROD-013
  - REQ-SEC-001
open_decisions:
  - OPEN-013
source-of-truth: canonical
---
# CAP-CMD-106 — Incident Coordination

## 1. Définition
Crée, reçoit et coordonne un Incident en préservant provenance, owner, priorité, impact, Tasks, prochaine action et liens vers Case, Decision, Response Run et Result.

## 2. Problème utilisateur
Un signal peut devenir Incident sans déduplication ni contexte, puis investigation et réponse se fragmentent. Principal : Incident Commander ; secondaires : SOC Analyst L1/L2, Business Owner, Response Operator. Sans capacité, ownership et retours interproduits se perdent.

## 3. Objectifs
Préserver la source ; appliquer une déduplication conceptuelle ; coordonner owner/priority/impact/next action ; lier plusieurs Cases ; initier Action Request sans Decision locale ; consommer Result.

## 4. Non-objectifs
Ne pas conduire forensic, confirmer Finding, créer Decision, exécuter containment/Response Run ou définir la machine d’état finale.

## 5. Propriétaire
Command / Incidents and Work Queue / Command Product Lead. Incident reste Command ; Case/Finding restent Investigate ; objets de réponse restent Govern.

## 6. Utilisateurs
Principal : Incident Commander. Secondaires : SOC Analyst L1/L2, Business Owner, Response Operator. Chaque transition réévalue permissions destination.

## 7. Conditions d’entrée
Signal/Alert ou Incident existant ; tenant/environnement ; permission ; déduplication conceptuelle vérifiée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Triage source | Signal/Alert ou manuel | source/provenance | non si manuel | source courante | justification/origine manuelle |
| Coordination context | utilisateur/Service/Task | impact, owner, next action | oui | version courante | Incident partial/triage |
| Cross-product links | Object Linking Service | Case/Decision/Run/Result refs | non | résolution courante | inaccessible sans fuite |

## 9. Objets lus
| Objet | Owner | Projection | Droit local |
|---|---|---|---|
| Signal / Alert | Command/source | source, severity, disposition | lecture |
| Case / Finding | Investigate | statut/résumé | projection |
| Decision / Response Run / Result | Govern | statut, conditions, résultat | projection |
| Incident / Task | Command | coordination | lecture/modification |

## 10. Objets créés ou modifiés
| Objet | Opération | Owner | Règle |
|---|---|---|---|
| Incident | créer, mettre à jour, relier, transitionner | Command | classe 2, future state machine |
| Task | créer/lier | Command | classe 2 |
| Action Request | initier via Govern | Govern | Command fournit contexte seulement |

## 11. Fonctionnalités
Créer/recevoir Incident ; préserver source ; dédupliquer/lier ; gérer coordination ; ouvrir/lier plusieurs Cases ; préparer Action Request ; consommer Run/Result et prochaine action.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Créer Incident | triage autorisé | Incident | 2 | source/justification + dedup | new/triaged | OPEN-013 |
| Modifier coordination | owner/coordinateur | Incident | 2 | version courante | audité | OPEN-013 |
| Lier/ouvrir Case | analyste autorisé | Case ref | 2 | permission Investigate | Case créé/lié par Investigate | non |
| Préparer Action Request | requester | Action Request | 2 | contexte/impact/action | draft Govern | oui |
| Fermer/réouvrir | Incident Commander | Incident | 2 | conditions futures | transition auditée | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Humain | Règle | Moteur | Workflow | Agent | Govern | Sans IA |
|---|---|---|---|---|---|---|---|
| Promotion Alert→Incident | décision/validation | possible | dedup/validation | possible | proposition | OPEN-013 | triage manuel |
| Package Case/Action Request | correction | règles de complétude | agrégation sourcée | possible | brouillon attribué | oui pour action | formulaire manuel |
Aucun agent ne confirme un Finding ou crée une Decision.

## 14. États fonctionnels
`new`, `triaged`, `assigned`, `in-progress`, `blocked`, `pending-external`, `pending-decision`, `monitoring`, `resolved`, `closed`, `reopened`. Ce sont des états de travail Draft, pas la machine finale.

## 15. États d’interface
Loading conserve l’Incident ; Empty concerne l’absence de source/relations ; Partial nomme les projections absentes ; Error garde les données Command ; Offline bloque mutation ; Permission denied ne divulgue rien ; Stale montre source/date. Rendu DS.

## 16. Sorties
| Sortie | Objet/événement | Consommateur | Garantie |
|---|---|---|---|
| Incident | objet Command | Queue, Mission Control, Investigate/Govern | provenance, owner, contexte |
| Case link | relation | Investigate/Command | owner Investigate |
| Action Request draft | objet Govern | Govern | origine Command et classe visibles |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte | Retour |
|---|---|---|---|---|
| Signal/Alert | promote/link | Incident | source IDs, entities, severity, confidence, disposition | triage |
| Incident | Open in Investigate | Case | tenant, env, Incident, time range, entities, alerts | Incident Detail |
| Incident/Finding context | request action | Govern | refs, impact, urgence, targets, alternatives | Incident Detail |
| Result | associated | Incident | Decision, Run, Result, residual risk | prochaine action |

## 18. Dépendances
CAP-CMD-102, 104, 107, 110, Object Linking Service, Timeline Engine, Investigate Case lifecycle, Govern Action Request lifecycle.

## 19. Source de vérité
Incident/Task : Command. Case/Finding : Investigate. Action Request/Decision/Run/Result : Govern. Toutes les projections conservent source et permissions.

## 20. Provenance et audit
Création, liaison, transition, package et proposition enregistrent acteur/producteur, source, version/run, before/after, justification, tenant et correlation ID.

## 21. Permissions fonctionnelles
`perm.command.incident.read/manage`, `perm.investigate.case.create/read`, `perm.govern.action-request.create`, ABAC tenant/env. Atomisation reportée.

## 22. Limites et erreurs
Source absente, duplicate possible, Case inaccessible, plusieurs Cases, destination indisponible, conflit de version ou permission refusée conservent Incident et contexte sans ownership implicite.

## 23. Métriques
Disposition Signal/Alert→Incident ; Incident→owner ; Incident→Case si nécessaire ; fermetures avec Result ou justification. Aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native ; preuve documentaire uniquement.

## 25. Critères d’acceptation
**Given** un Alert autorisé et aucun Incident équivalent, **When** il est promu, **Then** source/provenance, état initial et acteur sont conservés.

**Given** un Incident sans Case, **When** Open in Investigate est choisi, **Then** tenant/env/Incident sont transmis, Investigate possède le Case et le retour restaure Incident Detail.

**Given** une action classe 3, **When** elle est demandée, **Then** Command prépare Action Request, n’exécute rien et Govern reçoit le contexte.

## 26. Questions ouvertes
Quels états appartiennent à Incident versus coordination ? Quelles conditions de fermeture/réouverture ? — REQ-OBJ-001, REQ-PROD-003, REQ-PROD-008, REQ-PROD-013, REQ-SEC-001. `OPEN-013` reste ouverte.

## 27. Consommateurs documentaires
Incident Detail, Queue, Mission Control, parcours Alert→Incident/Incident→Case/Action Request, écrans Phase 6, objets Phase 7 et permissions ultérieures.
