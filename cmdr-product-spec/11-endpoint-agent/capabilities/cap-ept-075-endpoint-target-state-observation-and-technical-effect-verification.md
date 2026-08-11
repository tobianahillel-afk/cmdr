---
id: CAP-EPT-075
title: Endpoint Target-State Observation and Technical Effect Verification
product: endpoint-agent
module: containment
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-004, REQ-PROD-005, REQ-PROD-006, REQ-PROD-016, REQ-PROD-018, REQ-PROD-019, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004]
open_decisions: [OPEN-008, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-EPT-075 — Endpoint Target-State Observation and Technical Effect Verification

## 1. Définition
Observer l’état technique cible après une primitive et comparer `expected technical state` à `observed target state` comme fait Endpoint, sans déclarer l’objectif de réponse Govern satisfait.

## 2. Problème utilisateur
Un executor peut répondre success alors que l’état réel mismatch, est partiel, stale ou unverifiable. Inversement desired state observed ne signifie pas incident contained.

## 3. Objectifs
Collect current target observation ; timestamp/freshness ; compare matched/mismatched/partial/unknown/stale/unverifiable ; attach supporting observations ; expose drift seed ; hand off to Govern with limits.

## 4. Non-objectifs
Créer Govern Verification Assessment/Result, attribuer causalité, qualifier Incident/Finding, choose risk score, or define observation protocol/engine.

## 5. Propriétaire
Endpoint owns target-side observations and technical match semantics. Govern owns expected-response interpretation, residual risk and canonical verification/Result.

## 6. Utilisateurs
Verification Reviewer, Endpoint Operator, Response Operator, Govern Reviewer, Auditor, Investigate/Command consumers.

## 7. Conditions d’entrée
CAP-EPT-074 request/criteria, target accessible or explicit unavailable, permitted observations and exact expected technical state.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| verification request/criteria | CAP-EPT-074 | verification contract | oui | exact version | no assessment |
| expected technical state | CAP-EPT-066/074 | comparator | oui | pinned | incomplete |
| local target observations | Endpoint/EPT-2 | observed facts | oui per criterion | timestamped | unknown/unverifiable |
| primitive outcome | CAP-EPT-073 | execution context | oui | exact attempt | context gap |
| source limits/permissions | source/Security | interpretation constraint | selon source | current | restricted/partial |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Technical Verification Request | Endpoint | criteria/expected state | read |
| Technical Outcome | Endpoint | execution context | read |
| observations/capability state | Endpoint | current target facts | read |
| Govern Run/Verification Plan | Govern | correlation only | read/ref |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Target State Observation | create/version | Endpoint | source/time/limits mandatory |
| Technical Verification Observation | derive/version | Endpoint | technical match only |
| supporting observation links | add | source owners | no ownership transfer |
| Govern Verification/Result | aucune mutation | Govern | downstream only |

## 11. Fonctionnalités
Gather declared observations; preserve source and freshness; compare criteria; produce matched/mismatched/partially-matched/unknown/stale/unverifiable; expose supporting facts and contradictions; never overwrite executor status; make verification limits explicit.

## 12. Actions utilisateur
Inspect Class 0; deterministic observation/comparison Class 1; no target mutation or Govern verdict in this capability.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| collect allowed observations | oui | oui | oui | no need | source checks |
| compare expected/observed | oui | oui | oui | explain | criterion table |
| summarize mismatch/limits | oui | structured | oui | oui | raw comparison |
| declare incident contained/Result | non | interdit | non | interdit | Govern |

## 14. États fonctionnels
`observation-pending`, `matched`, `mismatched`, `partially-matched`, `unknown`, `stale`, `unverifiable`, `source-unavailable`, `contradictory`, `drift-detected`, `superseded`.

## 15. États d’interface
No Screen ID. Observed fact, derived technical match and any AI explanation remain visually/conceptually distinct in future presentation.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Target State Observation | technical fact | CAP-EPT-076/079 | source/time/limits |
| Technical Verification Observation | Endpoint assessment | Govern CAP-GOV-029 | match ≠ response success |
| drift/mismatch marker | technical condition | CAP-EPT-076 | no causal claim |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-EPT-074 | checks ready | CAP-EPT-075 | criteria/expected/sources | verification request |
| CAP-EPT-075 | mismatch/unknown/drift | CAP-EPT-076 | observed state/limits | Run |
| CAP-EPT-075 | observation ready | CAP-EPT-079 | technical verification refs | Govern |

## 18. Dépendances
CAP-EPT-015..030/066/073/074/076/079/080, CAP-GOV-028/029/032, Shared Trace, OPEN-008/013/015.

## 19. Source de vérité
Endpoint is SOT for local observed target state and technical match. Govern is SOT for response-level verification interpretation.

## 20. Provenance et audit
Target/primitive/execution, expected state, observations/source/version/time, comparison rule/version, match status, limitations, Run/Decision refs, reviewer/automation/AI and correlation.

## 21. Permissions fonctionnelles
Target-state/verification read, no-effect verification execute, restricted local observations, provenance read/export, cross-tenant deny.

## 22. Limites et erreurs
Observed desired state ≠ business objective achieved; process absent ≠ malicious activity resolved; isolation observed ≠ incident contained; restored file ≠ safe; unknown ≠ failed automatically.

## 23. Métriques
Matched/mismatch/partial/unknown/stale/unverifiable, runtime-success→verification-mismatch, drift after match, false business-verdict promotions target zero.

## 24. Classification de livraison
`draft / defined / planned`; no telemetry/verification engine, API/protocol or implementation.

## 25. Critères d’acceptation
**Given** primitive reports success but observed target state mismatches, **When** comparison runs, **Then** technical verification is mismatched and the raw success remains separately preserved.

**Given** desired state is observed, **When** handoff reaches Govern, **Then** Endpoint reports the technical match but does not claim response objective satisfied.

**Given** state cannot be observed, **When** verification runs, **Then** result is unknown/unverifiable rather than invented failure or success.

## 26. Questions ouvertes
OPEN-008/013/015 remain open; no final response-verification thresholds or runtime bridge is selected.

## 27. Consommateurs documentaires
CAP-EPT-076/079/080, Govern Verification/Result, Command/Investigate consumers, Quality/Audit.