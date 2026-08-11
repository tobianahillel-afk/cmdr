---
id: CAP-EPT-080
title: Governed Response Primitive Provenance and Cross-Product Contracts
product: endpoint-agent
module: containment
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-004, REQ-PROD-005, REQ-PROD-006, REQ-PROD-008, REQ-PROD-014, REQ-PROD-015, REQ-PROD-016, REQ-PROD-018, REQ-PROD-019, REQ-PROD-020, REQ-OBJ-007, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004]
open_decisions: [OPEN-007, OPEN-008, OPEN-013, OPEN-014, OPEN-015, OPEN-017]
source-of-truth: canonical
---
# CAP-EPT-080 — Governed Response Primitive Provenance and Cross-Product Contracts

## 1. Définition
Retracer la chaîne EPT-5 complète depuis Action Request/Decision/Response Run jusqu’à primitive request, authority/precheck, target execution, technical outcome, target-state observation, technical verification, optional reversal/release, Govern reconciliation et Result sans fusion d’ownership.

## 2. Problème utilisateur
Une action effectful cross-product peut perdre l’autorité, le scope ou le lien entre ce qui a été demandé, exécuté, observé et finalement jugé par Govern.

## 3. Objectifs
Stable refs/versions/timestamps ; exact requester/origin/target/primitive ; Decision/Run/Step ; Agent/version ; before/requested/observed states ; reversal/release ; verification limits ; owner per hop ; gaps/restrictions ; no identity merger.

## 4. Non-objectifs
Immutable ledger/crypto proof, API/protocol, object identity merger, Result creation, Evidence requalification, final physical schema, or EPT-6 provenance implementation.

## 5. Propriétaire
Endpoint owns provenance of its EPT-5 technical hops. Govern, Studio, Settings, Investigate, Command and Shared retain their canonical objects/mechanisms.

## 6. Utilisateurs
Auditor, Response Operator, Govern Reviewer, Endpoint Operator, Verification Reviewer, Security Reviewer, Investigator/Incident Commander consumers.

## 7. Conditions d’entrée
Stable refs from request/Run/technical stages, tenant/environment, owner/version/time, authority/provenance refs and destination permissions.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Action Request/Decision/Run refs | Govern | authority lineage | oui for governed action | historical/current | chain gap |
| EPT-5 primitive refs | CAP-EPT-065..079/081 | technical chain | oui if action | source time | partial chain |
| Studio refs | Studio | Tool/Human Gate/Automation provenance | when involved | exact version | explicit absent |
| Settings refs | Settings | Policy/Secret/Fleet context | when applicable | versioned | restricted/gap |
| source consumer refs | Investigate/Command | purpose/return context | when involved | stable refs | partial chain |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Action Request/Approval/Decision/Response Run/Result | Govern | provenance refs | no ownership transfer |
| Human Gate/Tool Call/Automation Run | Studio | origin/executor refs | read/link |
| Endpoint Policy/Fleet/Secret Reference | Settings | ref/version only | restricted read |
| Finding/Evidence/Case/Incident | Investigate/Command | source/consumer refs | permission-aware read/link |
| Trace/Activity/Job/Recovery | Shared | generic mechanism refs | read/link only |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| EPT-5 Provenance Chain | create/extend/version | Endpoint | typed owner per hop |
| Cross-Product Handoff Reference | create/link | Endpoint | destination owner retained |
| Provenance Gap/Restriction Marker | derive | Endpoint | missing data never invented |
| external canonical objects | aucune mutation | external owner | references only |

## 11. Fonctionnalités
Reconstruct `Action Request → Decision → Response Run → primitive → authority/precheck → execution → technical outcome → observation/technical verification → optional reversal/release → Govern reconciliation → Result`; preserve exact source owner/status; show missing/restricted hops; distinguish local audit from Shared Trace and audit reconstruction.

## 12. Actions utilisateur
Inspect provenance Class 0; deterministic reconstruction/gap detection Class 1; forward references Class 2 no-effect when permitted. No execution or authority from provenance access.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| correlate stable refs | oui | oui | oui | no need | relation traversal |
| detect gaps/owner mismatch | oui | oui | oui | explain | completeness rules |
| summarize chain | oui | structured | oui | oui, attributed | ordered hops |
| invent authority/state/Result | non | interdit | non | interdit | explicit gap/owner workflow |

## 14. États fonctionnels
`complete`, `partial`, `gap-present`, `restricted`, `stale-reference`, `owner-mismatch-candidate`, `destination-denied`, `superseded`, `unknown`.

## 15. États d’interface
No Screen ID. Facts, technical derivations, Govern assessments and AI summaries remain visually/conceptually distinguishable; navigation grants no permission.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| EPT-5 Provenance Chain | Endpoint audit context | Govern/Audit/Quality | owner/version/time typed |
| handoff refs | references | cross-product consumers | no permission/ownership transfer |
| gap/restriction report | diagnostic | Audit/Quality | no invention |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Investigate/Command context | response request | Govern | source refs/purpose | source owner retained |
| Govern Run/Step | authorized technical handoff | Endpoint EPT-5 | Decision/Run/target/scope | Govern retained |
| Endpoint outcome/verification | reconciliation | Govern | technical refs/limits | Endpoint remains raw fact owner |
| Govern Result | downstream return | consumers | Result ref/classification | Govern owner retained |

## 18. Dépendances
CAP-EPT-001..081, CAP-GOV-003/009/011/015/016/022/025..032/036, Studio, Settings, Shared, Command, Investigate, OPEN-007/008/013/014/015/017.

## 19. Source de vérité
Each product remains SOT of its canonical objects. Endpoint is SOT only for its technical EPT-5 hops and local provenance; Shared Trace is a mechanism, not fact owner.

## 20. Provenance et audit
Requester/origin/target/primitive, Action Request/Approval/Decision/Run/Step, Agent/version, Policy/Secret refs only, precheck, before/requested/observed state, execution/outcome/errors, verification, reversal/release, Result ref, timestamps, masking and correlation.

## 21. Permissions fonctionnelles
Provenance read, restricted response refs, cross-product source refs under source permission, export preparation, cross-tenant deny. No action permission inherited from audit access.

## 22. Limites et erreurs
Audit trail ≠ cryptographic proof automatically; Endpoint action ≠ Tool Call/Automation Run/Response Run; Human Gate ≠ Approval/Decision; local technical reversal ≠ Govern rollback; technical verification ≠ Govern verification.

## 23. Métriques
Chain completeness/gaps, stale/orphaned refs, owner/version coverage, cross-tenant denials, technical objects incorrectly merged/promoted target zero.

## 24. Classification de livraison
`draft / defined / planned`; no ledger/protocol/schema/API or implementation.

## 25. Critères d’acceptation
**Given** Studio Human Gate exists before a response action, **When** chain is reconstructed, **Then** the Gate remains Studio-owned and is not treated as Govern Approval.

**Given** Endpoint technically reverses an action, **When** provenance is reviewed, **Then** Technical Reversal and Govern Response Rollback remain distinct linked hops.

**Given** Govern creates Result after reconciliation, **When** the chain is traversed, **Then** Result links to Endpoint facts without changing their ownership or history.

## 26. Questions ouvertes
OPEN-007/008/013/014/015/017 remain open. No cross-product identity merger or proof mechanism is selected.

## 27. Consommateurs documentaires
Endpoint EPT-5, Govern, Studio, Settings, Shared, Command, Investigate, Security/Trust, Quality, Roadmap and future EPT-6 as boundary only.