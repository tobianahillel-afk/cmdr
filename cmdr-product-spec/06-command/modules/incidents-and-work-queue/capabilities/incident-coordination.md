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
Un Signal ou Alert peut devenir Incident sans déduplication ni contexte, puis investigation et réponse se fragmentent. Sans coordination, ownership et retours interproduits se perdent.

## 3. Objectifs
Préserver la source, appliquer une déduplication conceptuelle, coordonner owner/priority/impact/next action, lier plusieurs Cases, initier une Action Request et consommer Result.

## 4. Non-objectifs
Ne conduit pas le forensic, ne confirme pas Finding, ne crée pas Decision, n’exécute pas containment/Response Run et ne fixe pas la machine d’état finale.

## 5. Propriétaire
Command / Incidents and Work Queue / Command Product Lead. Incident reste Command ; Case/Finding restent Investigate ; objets de réponse restent Govern.

## 6. Utilisateurs
Principal : Incident Commander. Secondaires : SOC Analyst L1/L2, Business Owner, Response Operator.

## 7. Conditions d’entrée
Signal/Alert ou Incident existant, tenant/environnement, permission et déduplication conceptuelle vérifiée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Triage source | Signal, Alert ou saisie manuelle | provenance et contexte initial | non si création manuelle | état source courant | justification et origine manuelle obligatoires |
| Coordination context | utilisateur, Service et Task | impact, owner et prochaine action | oui | version courante | Incident `partial` ou en triage |
| Cross-product links | Object Linking Service | Case, Decision, Run et Result refs | non | résolution à l’ouverture | relation inaccessible sans fuite |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Signal / Alert | Command ou source | provenance, severity et disposition | consulter et relier |
| Case / Finding | Investigate | statut et résumé autorisé | consulter, relier et naviguer |
| Decision / Response Run / Result | Govern | statut, conditions et résultat | consulter et relier en lecture seule |
| Incident / Task | Command | coordination opérationnelle | consulter et modifier si autorisé |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Incident | créer, mettre à jour, relier ou transitionner | Command | classe 2, machine d’état finale reportée |
| Task | créer et relier | Command | classe 2, résultat attendu explicite |
| Action Request | demander création ou liaison | Govern | Command transmet le contexte ; Govern possède l’objet |
| Case / Finding / Result | aucune mutation | Investigate ou Govern | projections en lecture seule |

## 11. Fonctionnalités
Créer/recevoir Incident, préserver source, dédupliquer/lier, gérer coordination, ouvrir/lier plusieurs Cases, préparer Action Request et consommer Run/Result.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Créer Incident | triage autorisé | Incident | 2 | source/justification et dedup | Incident créé | OPEN-013 |
| Modifier coordination | owner/coordinateur | Incident | 2 | version courante | mutation auditée | OPEN-013 |
| Ouvrir ou lier Case | analyste autorisé | Case ref | 2 | permission Investigate | Case créé/lié par Investigate | non |
| Préparer Action Request | requester | Action Request | 2 | contexte, impact et action | draft Govern | oui |
| Fermer ou rouvrir | Incident Commander | Incident | 2 | conditions futures | transition auditée | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Dédupliquer Signal/Alert | oui | oui | oui | suggestion attribuée | recherche et comparaison manuelles |
| Créer ou promouvoir Incident | oui | validation possible | workflow possible | proposition seulement | triage humain complet |
| Préparer le contexte Case | oui | agrégation sourcée | oui | brouillon attribué | formulaire et références manuelles |
| Préparer Action Request | oui | règles de complétude | oui | brouillon attribué | package manuel vers Govern |
| Consommer Result | oui | liaison/validation | oui | résumé facultatif | lecture Govern et mise à jour humaine |

## 14. États fonctionnels
`new`, `triaged`, `assigned`, `in-progress`, `blocked`, `pending-external`, `pending-decision`, `monitoring`, `resolved`, `closed`, `reopened`. Ils restent Draft.

## 15. États d’interface
Partial nomme les projections absentes ; Error conserve les données Command ; Offline bloque mutation ; Permission denied ne divulgue rien ; Stale expose source/date.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Incident | objet Command | Queue, Mission Control, Investigate et Govern | provenance, owner et contexte conservés |
| Case link | relation typée | Investigate et Command | owner Investigate et permissions réévaluées |
| Action Request draft | objet Govern | Govern | origine Command, action class et impact visibles |
| Incident update from Result | événement Command | Mission Control et Work Queue | Result source inchangé et lié |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Signal ou Alert | promotion ou liaison | Incident | source IDs, entities, severity, confidence et disposition | triage ou Incident Detail |
| Incident | Open in Investigate | Case workspace | tenant, environnement, Incident, période, entities et Alerts | Incident Detail restauré |
| Incident ou Finding context | demande d’action | Govern | refs, impact, urgence, targets, action et alternatives | Incident Detail restauré |
| Result | association autorisée | Incident | Decision, Run, Result et residual risk | prochaine action Command |

## 18. Dépendances
CAP-CMD-102, CAP-CMD-104, CAP-CMD-107, CAP-CMD-110, Object Linking Service, Timeline Engine, Investigate Case lifecycle et Govern Action Request lifecycle.

## 19. Source de vérité
Incident/Task : Command. Case/Finding : Investigate. Action Request/Decision/Run/Result : Govern.

## 20. Provenance et audit
Création, liaison, transition, package et proposition enregistrent acteur/producteur, source, version/run, before/after, justification, tenant et correlation ID.

## 21. Permissions fonctionnelles
`perm.command.incident.read/manage`, `perm.investigate.case.create/read`, `perm.govern.action-request.create`, ABAC tenant/env ; atomisation reportée.

## 22. Limites et erreurs
Source absente, duplicate possible, Case inaccessible, plusieurs Cases, destination indisponible, conflit ou refus conservent Incident et contexte sans transfert d’ownership.

## 23. Métriques
Disposition Signal/Alert→Incident, Incident→owner, Incident→Case si nécessaire et clôtures justifiées ; aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native ; preuve documentaire uniquement.

## 25. Critères d’acceptation
**Given** un Alert sans Incident équivalent, **When** il est promu, **Then** provenance, état initial et acteur sont conservés.

**Given** un Incident sans Case, **When** Open in Investigate est choisi, **Then** Investigate possède le Case et le retour restaure Incident Detail.

**Given** une action classe 3, **When** elle est demandée, **Then** Command prépare Action Request, n’exécute rien et Govern reçoit le contexte.

## 26. Questions ouvertes
Quels états appartiennent à Incident et quelles conditions de fermeture/réouverture s’appliquent ? — Requirement IDs ci-dessus ; `OPEN-013` reste ouverte.

## 27. Consommateurs documentaires
Incident Detail, Queue, Mission Control, parcours Alert→Incident/Incident→Case/Action Request, écrans Phase 6, objets Phase 7 et permissions ultérieures.