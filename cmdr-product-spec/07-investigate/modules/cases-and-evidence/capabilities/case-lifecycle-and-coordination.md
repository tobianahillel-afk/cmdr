---
id: CAP-INV-102
title: Case Lifecycle and Coordination
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
open_decisions:
  - OPEN-013
---
# CAP-INV-102 — Case Lifecycle and Coordination

## 1. Définition
Créer, recevoir, organiser, affecter et faire évoluer un Case Investigate tout en conservant les liens Incident et la provenance.

## 2. Problème utilisateur
Une investigation durable exige objectif, scope, owner, contributeurs et prochaine action. Sans lifecycle explicite, le Case devient un dossier statique ou une seconde coordination Incident.

## 3. Objectifs
Créer ou recevoir un Case depuis Signal, Incident, Hunt ou action humaine ; gérer owner, contributeurs, objectif, scope, prochaine action et états Draft ; lier plusieurs Incidents lorsque permis.

## 4. Non-objectifs
Ne pas posséder Incident ; ne pas figer la machine d’état finale ; ne pas créer une Task générale ; ne pas approuver une réponse.

## 5. Propriétaire
Investigate / Cases and Evidence / Investigate Product Lead. Case est un objet canonique Investigate.

## 6. Utilisateurs
Principaux : Investigation Lead et Case Analyst. Secondaires : SOC Analyst et Incident Commander consommateur.

## 7. Conditions d’entrée
Tenant/environnement autorisés, objectif ou demande, owner ou stratégie d’affectation, source accessible et permissions create/update/assign.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Demande d’ouverture | Command, Signal, Hunt ou humain | contexte initial | oui pour création | état au moment de la demande | draft incomplet, aucune création silencieuse |
| Incident(s) parent(s) | Command | contexte opérationnel | non | version/fraîcheur visibles | Case standalone avec raison |
| Objectif et scope | Investigation Lead | contrat d’investigation | oui pour active | versionnés | rester requested/open |
| Owner/contributeurs | Identity/Presence | responsabilité | oui pour active | disponibilité visible | rester unassigned |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Incident | Command | priorité, impact, services, owner | consulter et lier |
| Signal/Alert | Command | source et triage | consulter et lier |
| Case | Investigate | version, statut, scope, relations | coordonner |
| Task | Command | action opérationnelle liée | consulter et relier |
| Decision/Result | Govern | statut/outcome | consulter |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Case | créer, modifier, suspendre, reprendre, clôturer, rouvrir | Investigate | transitions versionnées et permission-aware |
| Case–Incident relation | créer/supersede | Investigate relation | n’altère pas Incident |
| Case assignment | modifier | Investigate | owner et contributeurs distincts |
| Task | demander ou lier | Command | aucun lifecycle concurrent |

## 11. Fonctionnalités
Créer/recevoir Case, lier plusieurs Incidents, définir objectif/scope/owner/contributeurs/next action, gérer états, versionner et exposer une projection à Command.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Créer Case | Lead | Case | 2 | contexte/permission | requested/open | OPEN-013 |
| Affecter | Lead | owner | 2 | Case accessible | assignment auditée | OPEN-013 |
| Modifier scope | Analyst | Case | 2 | état modifiable | nouvelle version | OPEN-013 |
| Clôturer/réouvrir | Lead | état | 2 | critères/raison | transition | OPEN-013 |
| Lier Incident | Analyst | relation | 2 | Incident accessible | lien sourcé | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Préremplir depuis Incident | oui | oui | oui | oui | mapping déterministe |
| Proposer owner | oui | oui | oui | oui | availability/règles |
| Vérifier complétude | oui | oui | oui | explication | checklist |
| Transition d’état | oui | contrat explicite | workflow | proposition | action humaine/règle approuvée |

## 14. États fonctionnels
`requested`, `open`, `active`, `blocked`, `pending-data`, `pending-review`, `ready-for-action`, `closed`, `reopened`. Machine finale reportée.

## 15. États d’interface
Loading conserve le draft ; Partial nomme les sources ; Error garde la version ; Offline lecture seule ; Permission denied sans fuite ; Stale bloque ou confirme les mutations.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Case | objet | Investigate/Command projection | versionné, tenant-scoped, audité |
| Case status projection | projection | Command | owner, next action, fraîcheur |
| Case relation | relation | produits | sourcée sans transfert |
| Lifecycle event | Timeline/Activity | Replay/Audit | acteur, raison, version |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Command Incident | Open in Investigate | Case Lifecycle | tenant, env, Incident, service, Signals, objectif | Incident Detail |
| Signal/Hunt | create/link | Case Lifecycle | source, scope, période, provenance | source |
| Case | ready for action | Action Request Preparation | Findings, Evidence, impact | Case |
| Closed Case | reopen | Case Workspace | raison, prior state, unresolved items | Queue/Workspace |

## 18. Dépendances
Command projections, CAP-INV-101/103/107/109, Shared Collaboration/Timeline/Object Linking/Versioning, future Case state model et OPEN-013.

## 19. Source de vérité
Case, objectif, scope, owner et statut d’investigation sont Investigate ; Incident/Task restent Command ; Decision/Result restent Govern.

## 20. Provenance et audit
Source/requester, versions objectif/scope, assignments, transitions/raisons, liens Incident/Signal/Hunt et suggestions automatisées.

## 21. Permissions fonctionnelles
Case read/create/update, assign, close/reopen, link Incident/Signal, sensitive Case et cross-tenant interdit.

## 22. Limites et erreurs
Duplicate Case, source inaccessible, conflit de version, owner indisponible, transition invalide, Case fermé, tenant mismatch ou permission révoquée.

## 23. Métriques
Temps request→active, complétude objectif/scope/owner, reopen rate, durée blocked et Cases avec contexte parent.

## 24. Classification de livraison
`defined` / `planned`, cible native. Promotion conditionnée par modèle objet, versioning, duplicate handling, permissions et projection Command.

## 25. Critères d’acceptation
**Given** un Incident sans Case **When** Open in Investigate **Then** contexte conservé, Case Investigate créé/proposé et retour Incident disponible.

**Given** un second Incident lié **When** l’analyste l’ajoute **Then** deux Incidents restent Command et le scope Case est versionné.

**Given** un Case prêt **When** le lead clôture **Then** raison/version sont enregistrées et réouverture possible.

## 26. Questions ouvertes
Machine finale, cardinalités Incident–Case, duplicate policy et OPEN-013 restent ouvertes.

## 27. Consommateurs documentaires
Case Queue/Workspace, Command projection, Evidence/Finding, Action Request, phases Objets/Permissions.
