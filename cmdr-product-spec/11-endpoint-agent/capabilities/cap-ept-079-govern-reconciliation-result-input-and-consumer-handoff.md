---
id: CAP-EPT-079
title: Govern Reconciliation, Result Input and Consumer Handoff
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
# CAP-EPT-079 — Govern Reconciliation, Result Input and Consumer Handoff

## 1. Définition
Handoff Endpoint vers Govern des faits d’exécution EPT-5, Technical Outcomes, Target-State/Verification Observations, partial/failure/drift/reversal/release state et limitations, afin que Govern réconcilie le Response Run et crée éventuellement le canonical Result.

## 2. Problème utilisateur
Si Endpoint retourne seulement `success/fail`, Govern perd les faits de cible/partialité/verification. Si Endpoint crée directement Result, la frontière d’autorité est violée.

## 3. Objectifs
Return exact action/target/requested/observed state ; technical outcome ; uncertainty ; reversal/release ; verification limitations ; source timestamps ; Run/Step correlation ; consumer permissions ; preserve Result boundary.

## 4. Non-objectifs
Créer/finaliser Result, mark Response Run success, modify Decision/Incident/Case/Finding/Evidence, define cross-product API/protocol, or resolve OPEN-015.

## 5. Propriétaire
Endpoint owns the technical handoff payload-by-reference and its source facts. Govern owns reconciliation, verification interpretation, Response Run, rollback and Result.

## 6. Utilisateurs
Govern Response Operator/Reviewer, Verification Reviewer, Endpoint Operator, Incident Commander/Investigator as downstream consumers, Auditor.

## 7. Conditions d’entrée
A correlated Govern Run/Step, one or more Endpoint EPT-5 facts/outcomes, stable provenance refs and destination permission context.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Run/Step/Decision correlation | Govern | return origin | oui | stable/pinned | no canonical handoff |
| Technical Outcome | CAP-EPT-073 | execution fact | oui if action occurred | exact attempt | incomplete |
| Target/Verification Observation | CAP-EPT-075 | observed effect | when available/required | timestamped | limitation explicit |
| uncertainty/drift | CAP-EPT-076 | technical limitation | when present | current | none invented |
| reversal/release | CAP-EPT-077/078 | recovery facts | when applicable | exact attempt | explicit NA/unknown |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| EPT-5 technical records | Endpoint | source facts | read |
| Response Run/Step/Decision | Govern | destination/correlation | read/ref |
| Result | Govern | destination type/status only | no create/mutate |
| Finding/Evidence/Incident/Case | source products | return-context refs only | permission-aware link |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Govern Reconciliation Handoff | create/version | Endpoint | references facts, no Result |
| Consumer Handoff Ref | create | Endpoint | destination owner preserved |
| handoff restriction/gap | derive | Endpoint | missing/restricted fact explicit |
| Response Run/Result/source objects | aucune mutation directe | respective owner | consumer decides |

## 11. Fonctionnalités
Assemble exact technical refs; preserve per-target action/status; include requested vs observed state; include verification freshness/limits; include reversal/release facts; include partial/unknown/drift; mask restricted content; send to Govern; receive acknowledgement/reference without outcome merger; support consumer return origin.

## 12. Actions utilisateur
Inspect handoff Class 0; validate completeness/correlation Class 1; transmit refs Class 2 no target effect. Result creation/finalization stays Govern and may carry its own class/authority.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| assemble/correlate refs | oui | oui | oui | no need | structured handoff |
| validate completeness | oui | oui | oui | explain | checklist |
| summarize technical outcome | oui | source-backed | oui | oui, attributed | outcome matrix |
| create Govern Result/declare success | non | interdit | non | interdit | CAP-GOV-032 |

## 14. États fonctionnels
`draft`, `incomplete`, `ready`, `restricted`, `gap-present`, `handoff-requested`, `acknowledged`, `destination-denied`, `stale`, `superseded`.

## 15. États d’interface
No Screen ID. Technical facts, Govern normalized status, verification interpretation and Result remain separately attributed; inaccessible details are not leaked.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Reconciliation Handoff | Endpoint ref set | Govern CAP-GOV-026/029/031/032 | source facts/limits preserved |
| consumer handoff | references | Command/Investigate as allowed | no ownership/permission transfer |
| gap/restriction marker | limitation | Govern/Audit | no invented facts |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| EPT-5 technical execution | outcome available | CAP-EPT-079 | action/target/status/observations | source primitive |
| CAP-EPT-079 | handoff | Govern Response Run | technical refs/limits/provenance | same Run |
| Govern reconciliation | verification/Result | consumers | Govern-owned objects/refs | Endpoint remains source of raw facts |

## 18. Dépendances
CAP-EPT-065..078/080/081, CAP-GOV-022/025/026/029/031/032/033, Command/Investigate consumers, Shared Trace, OPEN-007/008/013/014/015/017.

## 19. Source de vérité
Endpoint remains SOT of raw technical EPT-5 facts. Govern is SOT of normalized response semantics, verification assessment, rollback and Result. Destination consumers retain their objects.

## 20. Provenance et audit
Run/Step/Decision, primitive/execution/target, requested/observed states, outputs/errors, verification source/time, uncertainty/drift, reversal/release, destination, masking/restriction, sender/ack timestamps and correlation.

## 21. Permissions fonctionnelles
Technical outcome/verification/provenance read, reconciliation handoff prepare/transmit, restricted context, export preparation, cross-tenant deny. No Result-create permission is granted to Endpoint.

## 22. Limites et erreurs
Technical output ≠ Result; technical verification ≠ Govern verification; technical success ≠ Response Run success; rollback technical success ≠ Govern reconciliation; destination acknowledgement ≠ outcome acceptance.

## 23. Métriques
Handoff completeness/gaps, destination denied, stale facts, orphaned executions, technical facts promoted directly to Result target zero.

## 24. Classification de livraison
`draft / defined / planned`; no cross-product API/protocol/schema or Result implementation.

## 25. Critères d’acceptation
**Given** Endpoint technical outcome is success but verification mismatches, **When** handoff reaches Govern, **Then** both facts are transmitted and Endpoint does not create Result.

**Given** Govern later creates Result after reconciliation, **When** provenance is inspected, **Then** Result remains Govern-owned and links back to raw Endpoint facts.

**Given** consumer lacks permission for sensitive technical output, **When** handoff is rendered, **Then** restricted refs remain masked without changing the underlying outcome.

## 26. Questions ouvertes
OPEN-007/008/013/014/015/017 remain open. No final Run bridge, result schema or Artifact relation is selected.

## 27. Consommateurs documentaires
Govern Runs/Verification/Rollback/Result, Command, Investigate, Shared Reporting/Trace, Quality, Audit.