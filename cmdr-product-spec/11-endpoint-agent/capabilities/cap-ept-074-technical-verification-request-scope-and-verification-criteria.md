---
id: CAP-EPT-074
title: Technical Verification Request, Scope and Verification Criteria
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
# CAP-EPT-074 — Technical Verification Request, Scope and Verification Criteria

## 1. Définition
Définir une demande de **Endpoint Technical Verification** portant sur l’état cible technique attendu après une primitive, avec scope, observations requises, freshness/window, critères techniques et comportement lorsque la vérification est indisponible.

## 2. Problème utilisateur
Technical completion ne prouve pas l’effet. À l’inverse, Endpoint ne doit pas reproduire le Govern Verification Plan qui décide si l’objectif de réponse est satisfait.

## 3. Objectifs
Relier primitive execution et expected technical target state ; définir scope/checks conceptuels ; observations/freshness/time window ; source limitations ; partial/unavailable verification ; consumer/Run refs ; provenance.

## 4. Non-objectifs
Créer Govern Verification Plan/Assessment/Result, inventer un score, modifier Finding/Incident, define verification engine/API/query, or infer business success.

## 5. Propriétaire
Endpoint owns Technical Verification Request and target-side criteria limited to technical state. Govern owns response-level Verification Plan/Assessment and residual-risk interpretation.

## 6. Utilisateurs
Verification Reviewer, Response Operator, Endpoint Operator, Govern Reviewer, Auditor.

## 7. Conditions d’entrée
Primitive execution/outcome and expected technical state known; verification requirement from Decision/Run available or explicit; authorized local observations identified.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| primitive execution/outcome | CAP-EPT-073 | subject | oui | exact attempt | no verification |
| expected technical state | CAP-EPT-066/Govern | criterion anchor | oui | pinned | incomplete |
| required local observations | primitive/capability | technical checks | oui if verifiable | declared source | unverifiable |
| freshness/window | Govern requirement/technical rule | temporal bound | selon check | current | explicit unknown |
| Run/consumer ref | Govern | return origin | oui | stable | orphaned request |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Technical Outcome | Endpoint | action/target/status | read |
| Requested Target State | Endpoint projection | expected technical state | read |
| Govern Verification Plan/Run | Govern | requirement/return context only | read/ref |
| local observations | Endpoint | verification source facts | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Technical Verification Request | create/version/cancel | Endpoint | scope bounded to target-side facts |
| Technical Criterion Set | create/version | Endpoint | no business criterion |
| observation request refs | create/link | Endpoint | no target mutation unless separate authorized check |
| Govern Verification | aucune mutation | Govern | consumes later |

## 11. Fonctionnalités
Pin execution/expected state; define technical checks and sources; apply freshness/window; distinguish independent observation from executor status where available; mark unverifiable/partial/source-unavailable; dispatch no-effect local checks; hand off observed facts without verdict inflation.

## 12. Actions utilisateur
Inspect Class 0; deterministic technical verification request/check Class 1 when no effect; any check with effect must use separate governed primitive and cannot be smuggled into verification.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| derive checks from primitive | oui | catalog/rules | oui | suggestions | check catalog |
| validate source/freshness | oui | oui | oui | explain | source matrix |
| summarize criteria | oui | structured | oui | oui | criterion table |
| decide Govern success/rollback | non | interdit | non | interdit | Govern verification |

## 14. États fonctionnels
`draft`, `criteria-incomplete`, `ready`, `verification-requested`, `observations-pending`, `partial`, `source-unavailable`, `unverifiable`, `stale`, `completed-for-technical-review`, `cancelled`, `superseded`.

## 15. États d’interface
No Screen ID. Verification unavailable, stale, partial and source permission denied remain explicit; no green-success UI contract is defined.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Technical Verification Request | Endpoint contract | CAP-EPT-075 | exact execution/expected state/checks |
| criterion/source matrix | technical context | Govern | no business verdict |
| unavailable/partial marker | limitation | CAP-EPT-076/079 | missing data preserved |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-EPT-073 | terminal/partial state | CAP-EPT-074 | execution/expected state | Run |
| CAP-EPT-074 | observations due | CAP-EPT-075 | criteria/sources/window | verification request |
| source unavailable | no observation | CAP-EPT-076/079 | limitation/affected criteria | Govern |

## 18. Dépendances
CAP-EPT-066/073/075/076/079/080, CAP-GOV-028/029, EPT-2 observations, Shared Trace, OPEN-008/013/015.

## 19. Source de vérité
Endpoint SOT of technical verification-request semantics. Govern SOT of response-level verification meaning. Source observations retain their owners.

## 20. Provenance et audit
Execution/primitive/target, expected state, criteria/version, source refs, freshness/window, requester, Run/Decision refs, unavailable reasons, timestamps, AI proposal/disposition, correlation.

## 21. Permissions fonctionnelles
Technical verification request/read, local observation read, restricted state, provenance export, cross-tenant deny; source permissions independently enforced.

## 22. Limites et erreurs
Technical verification ≠ Govern verification; verification unavailable ≠ action failed automatically; executor success is not sufficient criterion unless explicitly sourced; stale observations cannot prove current state.

## 23. Métriques
Requests ready/partial/unverifiable, source unavailable/stale, criterion coverage, verifications incorrectly producing Result target zero.

## 24. Classification de livraison
`draft / defined / planned`; no verification engine, query/API/protocol or final criterion schema.

## 25. Critères d’acceptation
**Given** verification source is unavailable, **When** verification is due, **Then** source-unavailable is recorded and no action-failed or success verdict is fabricated.

**Given** verification state is stale, **When** criteria are evaluated, **Then** freshness mismatch is explicit and desired-state match is not asserted as current.

**Given** AI is absent, **When** verification is requested, **Then** deterministic criterion/source rules provide the complete technical path.

## 26. Questions ouvertes
OPEN-008/013/015 remain open; final verification engine/thresholds remain unselected.

## 27. Consommateurs documentaires
CAP-EPT-075/076/079/080, Govern Verification, Security, Quality, Audit.