---
id: CAP-EPT-006
title: Endpoint Inventory Snapshot and Metadata
product: endpoint-agent
module: foundations
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-006, REQ-PROD-012, REQ-PROD-017, REQ-PROD-018, REQ-PROD-019, REQ-OBJ-008, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008]
source-of-truth: canonical
---
# CAP-EPT-006 — Endpoint Inventory Snapshot and Metadata

## 1. Définition
Define a point-in-time Endpoint-local inventory snapshot of source-backed technical facts and metadata, without turning it into a global CMDB, Fleet object or administrative Settings configuration.

## 2. Problème utilisateur
Consumers need coherent endpoint context, but individual facts collected at different times or from different owners can be mistaken for a complete current global source of truth.

## 3. Objectifs
Bind Agent/Endpoint refs, hostname when available, platform/OS/architecture, Agent version/build, sourced network/technical metadata, local capability summary, timestamp, completeness, limitations and provenance.

## 4. Non-objectifs
No global CMDB, Fleet administration, software/hardware exhaustive discovery, detailed telemetry, network scan, physical schema, response action or screen design.

## 5. Propriétaire
Endpoint Agent owns the local observed snapshot semantics. Settings retains Fleet and administrative configuration; Shared/other products retain their canonical objects.

## 6. Utilisateurs
Endpoint Operator, Platform Administrator as consumer, Investigate/SOC analyst, Command/Govern consumer and Auditor.

## 7. Conditions d’entrée
Agent identity, authorized observable sources, tenant/environment binding and per-fact source/freshness information.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Agent/Endpoint refs | Endpoint + shared/Settings projection | identity context | oui | current | snapshot incomplete |
| platform/OS/architecture | CAP-EPT-004 | technical facts | selon source | fact freshness | explicit unknown |
| Agent version/build | CAP-EPT-005 | technical facts | selon source | fact freshness | explicit unknown |
| local technical/network metadata | authorized local sources | metadata facts | non | per source | omitted + limitation |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| endpoint-agent | Endpoint Agent | identity/version/local state | read |
| Endpoint shared model | Settings-administered shared model | association/context | read if available |
| tenant/environment | Platform Settings | scope refs | read projection |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Inventory Snapshot functional record | create/refresh/supersede | Endpoint Agent | local observed record, not global CMDB/Fleet |
| local-audit-event | append snapshot provenance/change ref | Endpoint Agent | no secret value |

## 11. Fonctionnalités
Point-in-time snapshot; per-fact source/freshness; completeness/limitations; technical metadata; local capability summary; tenant-scoped provenance; explicit unknown/restricted facts.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect inventory | authorized user | Inventory Snapshot | 0 | read | snapshot + limits | non |
| validate completeness | Endpoint Operator | Inventory Snapshot | 1 | source map | complete/partial/unknown facts | non |
| request bounded refresh | authorized operator | Inventory Snapshot | 2 | source supports refresh | refresh request/derived snapshot only | no response authority |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| inspect/compare inventory | oui | oui | oui | sourced summary possible | tables/diff |
| assess completeness | oui | oui | oui | explanation only | deterministic source checks |
| invent missing metadata | non | non | non | interdit | unknown/omitted state |

## 14. États fonctionnels
`complete-for-declared-scope`, `partial`, `stale`, `restricted`, `unknown`, `superseded`.

## 15. États d’interface
Empty means no usable snapshot, not no endpoint; Partial lists missing facts; Offline displays last snapshot with freshness; Permission denied masks restricted metadata; Stale exposes timestamp. No final UI.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Inventory Snapshot | Endpoint technical record | Settings Fleet, Investigate, Command, Govern | point-in-time/source/limits explicit |
| Inventory Metadata Projection | projection | authorized consumers | no global system-of-record claim |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Agent/platform/version facts | snapshot requested/derived | Inventory Snapshot | refs/facts/source/time | source context retained |
| snapshot | new facts observed | newer snapshot | changed facts + previous ref | prior snapshot retained |
| snapshot | consumer projection | Settings/Investigate/etc. | scope/freshness/limitations | return origin retained |

## 18. Dépendances
CAP-EPT-001/003/004/005, Settings Fleet projections, Security/privacy, local-audit-event, Shared links and OPEN-008.

## 19. Source de vérité
Each underlying fact remains source-owned; the Endpoint Inventory Snapshot is authoritative only for its point-in-time local observation and declared scope, not global enterprise truth.

## 20. Provenance et audit
Record snapshot id/time, Agent/Endpoint refs, each fact source/freshness, completeness, restrictions, prior snapshot ref and correlation id.

## 21. Permissions fonctionnelles
Endpoint read for ordinary inventory; sensitive/restricted inventory access is a documented need subject to Security, not a new arbitrary atomic permission.

## 22. Limites et erreurs
Inventory != Fleet; snapshot != global source of truth; metadata != Settings configuration; missing facts remain unknown; no active discovery is implied.

## 23. Métriques
Snapshot completeness distribution; restricted/unknown facts; freshness; supersession/change counts; access denials.

## 24. Classification de livraison
`draft / defined / planned`; provider/platform-neutral documentary capability, no collector/CMDB implementation or detailed telemetry.

## 25. Critères d’acceptation
**Given** an inventory source is unavailable, **When** a snapshot is built, **Then** it is partial with the missing source named.

**Given** Settings has administrative metadata, **When** Endpoint inventory is viewed, **Then** that metadata is referenced rather than re-owned.

**Given** AI is unavailable, **When** completeness is assessed, **Then** deterministic source checks still produce the result.

## 26. Questions ouvertes
OPEN-008 remains open. Exact future inventory depth/source support belongs to later delivery decisions and EPT-2+ where relevant.

## 27. Consommateurs documentaires
Endpoint freshness/health/capabilities, Settings Fleet, Investigate/Command/Govern, registers, Quality, Roadmap and future contracts.