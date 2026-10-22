---
id: CAP-EPT-082
title: Administrative Update Assignment, Channel and Wave Boundary
product: endpoint-agent
module: updates
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-12
requirement_ids: [REQ-PROD-006, REQ-PROD-012, REQ-PROD-017, REQ-PROD-018, REQ-PROD-019, REQ-OBJ-008, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-013]
source-of-truth: canonical
---
# CAP-EPT-082 — Administrative Update Assignment, Channel and Wave Boundary

## 1. Définition
Projeter sur un Endpoint Agent la cible administrative d’update, son channel/wave et sa source sans transférer à Endpoint l’administration Fleet/Policy ni confondre assignment et exécution locale.

## 2. Problème utilisateur
Une version désirée ou une wave Settings peut être interprétée à tort comme update reçue, téléchargée ou appliquée.

## 3. Objectifs
Exposer target version, channel/wave refs, assignment source, Fleet/Policy refs, received/stale/mismatch state et provenance.

## 4. Non-objectifs
Aucune création de wave/channel, aucun scheduler Fleet, download/install, support claim, API ou implementation.

## 5. Propriétaire
Settings possède l’administration de l’assignment; Endpoint possède seulement sa projection locale/effective et ses faits de réception.

## 6. Utilisateurs
Platform Administrator, Endpoint Operator, Release Reviewer, Auditor, SOC consumer.

## 7. Conditions d’entrée
Agent/tenant/environment connus, assignment Settings référencé, permission de lecture et fraîcheur explicite.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| target version | Settings Fleet | admin target | oui si assigné | versionnée | unassigned/unknown |
| channel/wave ref | Settings | rollout context | non | current admin state | no wave inference |
| Fleet/Policy ref | Settings | ownership context | oui si applicable | current | restricted/unknown |
| observed current version | CAP-EPT-005 | local fact | oui | latest observation | comparison unknown |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Endpoint Agent Fleet | Settings | target/channel/wave | read projection |
| Endpoint Policy | Settings | effective assignment constraints | read projection |
| Endpoint Agent | Endpoint | current version/state | local read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Local Update Assignment Projection | derive/refresh | Endpoint | source Settings preserved |
| Assignment State | derive | Endpoint | assigned/received/stale/mismatch only |
| Settings objects | aucune mutation | Settings | refs only |

## 11. Fonctionnalités
Recevoir et lier assignment/version/channel/wave; comparer current vs target; détecter stale/mismatch; préserver source/version; ne jamais démarrer un update par simple assignment.

## 12. Actions utilisateur
Inspect assignment Class 0; deterministic comparison Class 1; bounded refresh request Class 2 if source supports. Assignment mutation remains Settings-owned.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| inspect/compare assignment | oui | oui | oui | explain only | exact refs/comparison |
| detect stale/mismatch | oui | oui | oui | summarize | deterministic freshness rules |
| choose/change target | Settings only | policy rules | non Endpoint | interdit | Settings workflow |

## 14. États fonctionnels
`unassigned`, `assigned-not-received`, `received`, `stale`, `current-matches-target`, `mismatch`, `superseded`, `unknown`.

## 15. États d’interface
No Endpoint Screen ID. Partial/Stale/Denied states preserve source and do not imply local failure.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| assignment projection | Endpoint fact | EPT-6 update capabilities | admin/local distinction explicit |
| target comparison | derived state | Settings/Endpoint | no execution inference |
| provenance refs | audit context | Quality/Audit | owner/version/time retained |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Settings assignment | publish/change | CAP-EPT-082 | target/channel/wave/Fleet/Policy refs | Settings ownership retained |
| CAP-EPT-082 | local readiness work | CAP-EPT-083/084 | assignment + current version | assignment ref retained |
| local state | mismatch/stale | Settings consumer | observed vs target | no admin mutation |

## 18. Dépendances
CAP-EPT-001/003/005/011/013; Settings Fleet/Policy/upgrade management; Security; OPEN-008/013.

## 19. Source de vérité
Settings is SOT for administrative assignment; Endpoint SOT for locally received/effective projection and observed current state.

## 20. Provenance et audit
Agent, tenant/environment, assignment/Fleet/Policy/channel/wave refs, target/current version, source version, received time, freshness and mismatch reason.

## 21. Permissions fonctionnelles
Update-state read, Settings-assignment projection read, refresh request if supported, cross-tenant deny; no final RBAC.

## 22. Limites et erreurs
Administrative target != current version; assigned != downloaded; Settings wave != local execution; compatible/support not inferred.

## 23. Métriques
Assigned/unassigned, stale, mismatch, projection freshness, cross-tenant denial and orphaned assignment counts.

## 24. Classification de livraison
`draft / defined / planned`; documentary boundary only.

## 25. Critères d’acceptation
**Given** Settings assigns version X, **When** Agent still reports version Y, **Then** mismatch is shown without update execution.

**Given** a wave changes, **When** local projection is stale, **Then** stale is explicit rather than silently refreshed.

**Given** AI is unavailable, **When** target/current are compared, **Then** deterministic comparison still works.

## 26. Questions ouvertes
OPEN-008 and OPEN-013 remain open; no launch platform or class-2 default is selected.

## 27. Consommateurs documentaires
EPT-6 update lifecycle, Settings Fleet/Policy, Endpoint maps/registers, Security, Quality, Roadmap.