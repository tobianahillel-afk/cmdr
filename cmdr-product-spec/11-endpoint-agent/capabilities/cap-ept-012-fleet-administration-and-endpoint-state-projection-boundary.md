---
id: CAP-EPT-012
title: Fleet Administration and Endpoint State Projection Boundary
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
# CAP-EPT-012 — Fleet Administration and Endpoint State Projection Boundary

## 1. Définition
Define the boundary by which individual Endpoint technical state is projected to the Platform Settings-owned Endpoint Agent Fleet without creating a competing Fleet object or administration capability.

## 2. Problème utilisateur
Fleet views need local Agent state, but aggregation can be mistaken for transfer of ownership or allow Endpoint to administer fleet membership, selectors or upgrade waves.

## 3. Objectifs
Project individual registration, version/build, inventory, health, last-seen, operational state and capability availability with source/freshness while preserving Settings Fleet ownership.

## 4. Non-objectifs
No Fleet creation/admin, selector, bulk action, enrollment administration, upgrade wave, fleet policy, screen redesign or Settings mutation.

## 5. Propriétaire
Endpoint Agent owns only individual technical facts and this projection contract. Platform Settings owns `endpoint-agent-fleet` and all Fleet administration.

## 6. Utilisateurs
Platform Administrator, Endpoint Operator, SOC/Investigate analyst, Command/Govern consumer and Auditor.

## 7. Conditions d’entrée
Individual Agent identity plus source-backed EPT-1 facts, tenant/environment scope and Settings Fleet reference when available.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| registration/enrollment state | CAP-EPT-001/002 | local technical projection | oui | state freshness | Fleet projection partial |
| version/inventory/health/state | CAP-EPT-005..010 | technical projections | selon scope | per fact | explicit unknown/partial |
| capability availability | CAP-EPT-011 | technical projection | non | fact freshness | capability unknown |
| Fleet reference | Platform Settings | administrative ref | non | Settings current | no Fleet association inferred |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| endpoint-agent | Endpoint Agent | individual facts | read |
| endpoint-agent-fleet | Platform Settings | Fleet association/context | Settings read projection |
| endpoint-policy | Platform Settings | policy ref only where projected | Settings read projection |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Endpoint State Projection | derive/publish conceptually | Endpoint Agent | individual facts only |
| endpoint-agent-fleet | none/read-consumer only | Platform Settings | no local Fleet mutation |

## 11. Fonctionnalités
Individual state projection; source/freshness; partial/unknown handling; Fleet association ref; no ownership transfer; aggregation-ready semantics.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect individual projection | authorized user | endpoint-agent | 0 | read | local state set | non |
| validate projection completeness | Endpoint/Settings reviewer | projection | 1 | source facts | complete/partial/unknown | non |
| publish/refresh projection | authorized technical path | projection | 2 | source supports | bounded projection refresh | no Fleet admin |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| inspect/aggregate projection | oui | oui | oui | summary possible | source tables |
| validate freshness/completeness | oui | oui | oui | explanation only | deterministic checks |
| mutate Fleet/selectors/waves | Settings only | external rules | non via Endpoint | interdit ici | Settings-owned administration |

## 14. États fonctionnels
`projection-current`, `projection-partial`, `projection-stale`, `projection-unknown`, `fleet-association-unknown`.

## 15. États d’interface
Settings surfaces own final UI. Endpoint documentation only states projection semantics; Loading/Partial/Stale/Permission denied preserve individual/source state. No new Screen ID.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Individual Endpoint State Projection | projection | Platform Settings Fleet | facts/source/freshness explicit |
| Projection Provenance | local-audit/provenance context | Audit | no Fleet ownership transfer |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| EPT-1 local facts | projection refresh | Settings Fleet consumer | Agent/scope/facts/freshness | Endpoint context retained |
| Settings Fleet | drill/reference | individual Endpoint context | Fleet/Agent refs | Settings origin retained |
| projection | stale/missing fact | partial/unknown projection | affected fact refs | no Fleet mutation |

## 18. Dépendances
CAP-EPT-001..011, Settings Fleet/enrollment/capability inventory/upgrade management, canonical endpoint-agent-fleet, OPEN-008.

## 19. Source de vérité
Endpoint owns individual local facts; Platform Settings owns Fleet membership/configuration/aggregation/admin. Projection never transfers lifecycle ownership.

## 20. Provenance et audit
Record Agent/Fleet refs, tenant/environment, fact refs/timestamps, projection time, missing/restricted facts and correlation id.

## 21. Permissions fonctionnelles
Endpoint read governs local facts; Settings Fleet read/manage remain Settings-owned. Projection does not grant Fleet manage or cross-tenant access.

## 22. Limites et erreurs
Inventory != Fleet; Fleet != local state; missing Fleet ref does not create one; no Fleet selector/upgrade wave/admin action is defined.

## 23. Métriques
Projection current/partial/stale counts; missing fact categories; association-unknown count; denied projection reads.

## 24. Classification de livraison
`draft / defined / planned`; boundary documentation only, no Fleet service/API/screen/implementation.

## 25. Critères d’acceptation
**Given** local health changes, **When** Settings consumes the projection, **Then** Settings may aggregate the new value while Fleet ownership remains unchanged.

**Given** no Fleet association is known, **When** projection is prepared, **Then** association remains unknown and Endpoint creates no Fleet object.

**Given** an upgrade wave exists in Settings, **When** Endpoint state is projected, **Then** EPT-1 does not mutate or execute the wave.

## 26. Questions ouvertes
OPEN-008 remains open; future delivery/aggregation mechanics are not selected.

## 27. Consommateurs documentaires
Settings Fleet/capability inventory, Endpoint capability map, Command/Investigate/Govern consumers, registers, Quality and Roadmap.