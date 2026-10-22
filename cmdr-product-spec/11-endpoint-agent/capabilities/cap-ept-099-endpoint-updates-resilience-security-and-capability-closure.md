---
id: CAP-EPT-099
title: Endpoint Updates, Resilience, Security and Capability Closure
product: endpoint-agent
module: closure
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-12
requirement_ids: [REQ-PROD-004, REQ-PROD-005, REQ-PROD-006, REQ-PROD-008, REQ-PROD-012, REQ-PROD-014, REQ-PROD-015, REQ-PROD-016, REQ-PROD-017, REQ-PROD-018, REQ-PROD-019, REQ-PROD-020, REQ-OBJ-007, REQ-OBJ-008, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004, REQ-SEC-005]
open_decisions: [OPEN-007, OPEN-008, OPEN-013, OPEN-014, OPEN-015, OPEN-017]
source-of-truth: canonical
---
# CAP-EPT-099 — Endpoint Updates, Resilience, Security and Capability Closure

## 1. Définition
Consolider la provenance et les frontières de capability de l’Endpoint Agent à travers EPT-1..EPT-6, et fournir le contrat de clôture documentaire sans devenir une mega-capability d’implémentation.

## 2. Problème utilisateur
Après six lots, la plateforme doit pouvoir démontrer que toutes les familles Endpoint obligatoires sont couvertes, que les owners restent distincts et que les handoffs sont traçables sans prétendre que l’Agent est implémenté.

## 3. Objectifs
Reconstruct global Endpoint chain, owner per hop, capability-family completeness, cross-product handoffs, unresolved OPEN, delivery status, provenance and closure conditions.

## 4. Non-objectifs
Aucun product code, runtime implementation, API/protocol, physical schema, supported-platform declaration, final RBAC, Screen ID, Phase 6 capability or repository-global completion.

## 5. Propriétaire
Endpoint Agent Product Lead owns Endpoint capability-layer closure. Product Architecture/Quality own global registry/verification; other domains retain their canonical objects/mechanisms.

## 6. Utilisateurs
Product Architecture, Endpoint Product Lead, QA/Traceability Lead, Security Reviewer, Platform Administrator, Command/Investigate/Govern/Studio owners, Auditor.

## 7. Conditions d’entrée
EPT-1..5 PASS evidence, EPT-6 capability set complete, registers/maps/requirements/open decisions available, no blocking owner conflict or capability gap, tenant/product boundaries explicit.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| CAP-EPT-001..098 | Endpoint | capability contracts | oui | current branch | closure blocked |
| EPT-1..5 verification records | Quality | historical PASS evidence | oui | immutable refs | closure blocked |
| EPT-6 conformance/build evidence | Quality | current lot evidence | oui | build-time | final PASS pending remote |
| registers/Requirements/OPEN | Governance | traceability | oui | current | closure blocked/partial |
| cross-product owner sources | Settings/Shared/Security/Studio/Govern/Investigate/Command | ownership boundaries | oui | current | owner gap/conflict |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Endpoint capability contracts | Endpoint | ids/status/dependencies | read |
| Fleet/Policy/Secret References | Settings | admin boundary refs | read projection |
| Tool/Automation/Deployment | Studio | cross-product refs | read/link |
| Decision/Response Run/Result | Govern | authority/result refs | read/link |
| Evidence/Finding/Case/Incident | Investigate/Command | consumer/source refs | read/link |
| Trace/Activity/Jobs/Recovery | Shared | mechanism refs | read/link |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Endpoint Capability Closure State | derive/version | Endpoint/Product Architecture | documentary completeness only |
| Endpoint Provenance Chain | consolidate/link | Endpoint | no foreign object mutation |
| external canonical objects | aucune mutation | external owner | references only |

## 11. Fonctionnalités
Demonstrate chain `administrative configuration → Endpoint identity → platform/version → telemetry → detection → investigation → collection → technical execution → containment → verification → updates/resilience/security → cross-product consumers`; verify namespaces/families/owners; expose gaps; preserve planned delivery; prepare closure evidence.

## 12. Actions utilisateur
Inspect closure/provenance Class 0; deterministic completeness/owner/traceability checks Class 1; prepare review package Class 2 no-effect. No implementation or response execution.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| count/validate capability structure | oui | oui | oui | no need | registry/structural checks |
| detect owner/gap conflicts | oui | oui | oui | explain | dependency/ownership maps |
| summarize closure | oui | structured | oui | oui, attributed | conformance report |
| declare PASS without gates | non | interdit | non | interdit | QA gate process |

## 14. États fonctionnels
`incomplete`, `build-complete`, `pending-remote-verification`, `blocking-gap`, `owner-conflict`, `traceability-gap`, `closure-pass`, `superseded-evidence`.

## 15. États d’interface
No Endpoint Screen ID. Any future status surface must distinguish documentary PASS from implementation/release readiness.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Endpoint Closure State | documentary result | Roadmap/Quality/Product Architecture | capability completeness only |
| Endpoint Provenance Chain | typed refs | Audit/cross-product owners | owner per hop retained |
| gap/conflict report | validation evidence | reviewers | no hidden missing family |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| EPT-1..6 capability sets | closure audit | Endpoint closure report | ids/counts/owners/dependencies | gap/pass state |
| Endpoint closure PASS | Phase 5 evaluation | Roadmap Phase 5 | Endpoint PASS + Studio PASS | Phase 5 capability-spec status |
| Phase 5 closure | next roadmap planning | Phase 6 preflight | historical closure only | no Phase 6 capability started |

## 18. Dépendances
CAP-EPT-001..098; EPT-1..5 post-publication reports; EPT-6 source/conformance reports; Governance registers; Settings/Studio/Govern/Shared/Security/Command/Investigate; OPEN-007/008/013/014/015/017.

## 19. Source de vérité
Each domain remains SOT for its objects/mechanisms. Endpoint closure is authoritative only for Endpoint capability coverage and local technical provenance; global repository maturity remains governed by roadmap/quality.

## 20. Provenance et audit
Capability IDs/files/statuses, source audits, owner/dependency/Requirement/OPEN refs, validation gates, commit/baseline/build/final refs, counts and historical PASS evidence are preserved.

## 21. Permissions fonctionnelles
Capability/provenance read, sensitive source refs according owner permissions, quality review, cross-tenant deny. Closure does not grant action authority or final RBAC.

## 22. Limites et erreurs
Endpoint Capability Specification PASS != Endpoint implementation complete; Phase 5 PASS != repository complete; local provenance != immutable ledger; no OPEN is closed to manufacture PASS.

## 23. Métriques
Capability/file/section/table/GWT counts, missing/duplicate/recycled IDs, owner conflicts, traceability gaps, placeholders/generic tables, OPEN/Requirement coverage and gate status.

## 24. Classification de livraison
`draft / defined / planned`; closure is documentary capability specification only, not implementation, deployment or production readiness.

## 25. Critères d’acceptation
**Given** all 99 Endpoint capabilities are structurally valid and owner-consistent, **When** closure is audited, **Then** Endpoint may become documentary PASS only after EPT-6 post-publication gates also pass.

**Given** OPEN-008 remains unresolved, **When** closure is evaluated, **Then** Endpoint may still pass if contracts remain platform-neutral and no support claim is made.

**Given** Endpoint and Studio both pass, **When** Phase 5 is evaluated, **Then** Phase 5 may pass as capability-specification complete while global repository maturity remains PARTIAL and Phase 6 stays unstarted.

## 26. Questions ouvertes
OPEN-007/008/013/014/015/017 remain open; EPT-6 closes none. Phase 6 questions are deferred to a separate preflight.

## 27. Consommateurs documentaires
Endpoint README/maps/registers, Governance, Quality, Delivery Roadmap Phase 5, Command/Investigate/Govern/Studio/Settings/Shared/Security and future Phase 6 preflight.