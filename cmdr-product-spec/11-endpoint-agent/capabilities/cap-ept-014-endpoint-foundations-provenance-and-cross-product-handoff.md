---
id: CAP-EPT-014
title: Endpoint Foundations Provenance and Cross-Product Handoff
product: endpoint-agent
module: foundations
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-006, REQ-PROD-012, REQ-PROD-017, REQ-PROD-018, REQ-PROD-019, REQ-OBJ-008, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-EPT-014 — Endpoint Foundations Provenance and Cross-Product Handoff

## 1. Définition
Define provenance and cross-product handoff for EPT-1 facts from Settings enrollment/configuration through Endpoint identity/binding/platform/version/inventory/health/connectivity/freshness/capability projections to authorized consumers, preserving every canonical owner.

## 2. Problème utilisateur
Cross-product consumption can erase origin, freshness or ownership and cause Endpoint facts to be mistaken for Evidence/Finding, Govern Result, Response Run, Tool Call, Fleet/Policy ownership or Shared Trace.

## 3. Objectifs
Link source refs, tenant/environment, sequence of foundation facts, consumer handoffs, limitations, return origin, correlation and local audit provenance without semantic promotion.

## 4. Non-objectifs
No implementation protocol, event bus, API, physical schema, response execution, Tool/Automation Run, Evidence qualification, Result creation, Shared Trace engine or final RBAC.

## 5. Propriétaire
Endpoint Agent owns provenance of its local technical foundation facts. Every external object/mechanism retains its canonical owner.

## 6. Utilisateurs
Endpoint Operator, Platform Administrator, Investigate/SOC analyst, Command operator, Govern reviewer, Studio integrator and Auditor.

## 7. Conditions d’entrée
One or more EPT-1 source facts with tenant/environment/source/time plus authorized consumer context. Missing lineage remains explicit.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| EPT-1 fact/projection refs | CAP-EPT-001..013 | provenance refs | oui | per fact | handoff partial |
| Settings refs | Platform Settings | enrollment/Fleet/Policy/scope context | selon handoff | source freshness | source absent/unknown |
| consumer context/return origin | consuming product | navigation/workflow context | oui for handoff | current request | handoff blocked/limited |
| local audit refs | Endpoint Agent | provenance context | non | event history | audit linkage partial |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| endpoint-agent/local-audit-event | Endpoint Agent | identity/fact provenance | Endpoint read |
| endpoint-agent-fleet/endpoint-policy/tenant/environment | Platform Settings | source refs only | Settings read projection |
| Decision/Response Run/Result | Govern | consumer/reference boundaries only | source permission |
| Tool/Automation Agent/Automation Run | Studio | consumer/reference boundaries only | source permission |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Foundation Handoff/Provenance context | derive/link | Endpoint Agent | links facts, does not promote foreign objects |
| local-audit-event | append handoff/source refs | Endpoint Agent | local audit != Shared Trace |

## 11. Fonctionnalités
Source-to-fact lineage; correlation; tenant/environment preservation; freshness/limitations; Settings→Endpoint handoff; Endpoint→consumer projection; return origin; no automatic semantic promotion.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect provenance/handoff | authorized user | provenance context | 0 | source permissions | lineage + limitations | non |
| validate lineage/completeness | Endpoint/Audit reviewer | provenance refs | 1 | refs available | complete/partial/conflict | non |
| prepare bounded consumer handoff | authorized user/system | projection context | 2 | destination permitted | attributed handoff only | no authority transfer |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| inspect/link provenance | oui | oui | oui | sourced summary possible | deterministic refs/trace table |
| explain incomplete/conflicting lineage | oui | oui | oui | oui, attributed | structured limitation review |
| convert technical fact into Evidence/Result/authority | source owner only | validation only | non autonome | interdit | explicit owner workflow |

## 14. États fonctionnels
`provenance-complete`, `provenance-partial`, `provenance-conflict`, `source-restricted`, `handoff-ready`, `handoff-blocked`, `handoff-complete`.

## 15. États d’interface
Loading/Partial/Error/Stale/Permission denied preserve return origin and source refs; restricted provenance is masked, not fabricated. No Screen ID or detailed UX.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Foundation Provenance Context | Endpoint technical record/projection | Audit + all authorized consumers | source/owner/time/limits retained |
| Cross-Product Handoff Package | attributed projection | Settings/Investigate/Command/Govern/Studio | no ownership/authority promotion |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Settings enrollment/config refs | local processing | Endpoint identity/foundation facts | source/tenant/env refs | Settings ownership retained |
| Endpoint foundation facts | authorized consume | Investigate/Command/Govern/Studio/Settings | fact refs/freshness/limits/provenance | Endpoint origin retained |
| consumer workflow | follow-up | source owner | return origin/correlation/source refs | no implicit mutation |

## 18. Dépendances
CAP-EPT-001..013, Settings Fleet/Policy/enrollment, Govern boundaries, Studio boundaries, Investigate/Command consumption, Shared generic mechanisms, Security, OPEN-008/013/015.

## 19. Source de vérité
Each fact/object stays with its canonical owner. Endpoint provenance links them; it does not become a competing source for Fleet, Policy, Evidence, Finding, Decision, Result, Tool, Automation Run, Job or Trace.

## 20. Provenance et audit
Record all source refs, owner, tenant/environment, fact timestamps/versions, actor/engine, consumer, return origin, masking/restrictions, limitations, correlation ids and local-audit refs; history is retained.

## 21. Permissions fonctionnelles
Source permissions govern source visibility; Endpoint read/local-audit read govern local provenance. Handoff never expands permission, and export preparation does not grant export.

## 22. Limites et erreurs
Endpoint facts consumed by Investigate != Evidence/Finding automatically; technical status/output != Govern Result; execution != Response Run/Tool Call; local audit event != Shared Trace; provenance != raw log.

## 23. Métriques
Complete/partial/conflicting provenance; blocked/restricted handoffs; missing source refs; cross-tenant denial; consumer/return-origin coverage.

## 24. Classification de livraison
`draft / defined / planned`; documentary handoff/provenance only, no runtime bus/API/protocol/code, final schema/RBAC or EPT-2+ execution.

## 25. Critères d’acceptation
**Given** Endpoint inventory is consumed by Investigate, **When** handed off, **Then** it remains an attributed technical fact and is not Evidence/Finding automatically.

**Given** Endpoint technical status is visible during a Govern workflow, **When** consumed, **Then** it does not become canonical Result or Response Run state automatically.

**Given** local audit provenance is linked to a Shared Trace, **When** inspected, **Then** the two remain distinct records/mechanisms with their owners preserved.

## 26. Questions ouvertes
OPEN-008/013/015 remain open; EPT-1 selects no cross-product runtime protocol, permission convergence or execution bridge.

## 27. Consommateurs documentaires
Endpoint/Settings/Investigate/Command/Govern/Studio/Shared boundaries, registers, Quality, Roadmap, audit and future implementation contracts.