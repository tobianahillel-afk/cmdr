---
id: CAP-EPT-088
title: Update Reversion and Previous-Version Recovery Boundary
product: endpoint-agent
module: updates
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-12
requirement_ids: [REQ-PROD-004, REQ-PROD-005, REQ-PROD-006, REQ-PROD-012, REQ-PROD-016, REQ-PROD-017, REQ-PROD-018, REQ-PROD-019, REQ-PROD-020, REQ-OBJ-007, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004]
open_decisions: [OPEN-008, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-EPT-088 — Update Reversion and Previous-Version Recovery Boundary

## 1. Définition
Définir la reversion technique locale d’un update vers une version précédente disponible et la vérification associée, sans devenir Studio Deployment Reversion ni Govern Response Rollback.

## 2. Problème utilisateur
Restaurer une version précédente peut être confondu avec restauration complète de l’endpoint, rollback de réponse ou réconciliation Fleet.

## 3. Objectifs
Represent previous-version availability, eligibility, request, local reversion, partial/failure, restored version, health verification, state limitations and provenance.

## 4. Non-objectifs
Aucun rollback engine, exact full-state restore, Govern Response Rollback, Studio asset reversion, Fleet-policy reconciliation automatique ou implementation.

## 5. Propriétaire
Endpoint owns local Agent-version reversion and raw recovery facts. Settings owns administrative target/waves; Studio owns Studio deployment reversion; Govern owns response rollback/recovery governance.

## 6. Utilisateurs
Endpoint Operator, Platform Administrator, Release Reviewer, Govern Reviewer when linked, Auditor.

## 7. Conditions d’entrée
Exact update attempt/version state, previous-version availability, current health/target facts, applicable policy/authority and recovery limitations.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| current/previous version refs | CAP-EPT-005/083 | version context | oui | current/pinned | no reversion |
| update attempt + verification | CAP-EPT-085..087 | recovery subject | oui | exact | incomplete context |
| current target/assignment | CAP-EPT-082 | admin projection | oui | current | reconciliation unknown |
| policy/authority context | Settings/Govern/Security | control context | oui if required | current | blocked |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Local Update Operation | Endpoint | failed/current update | read |
| Post-Update Verification | Endpoint | health/mismatch | read |
| Endpoint Agent Fleet/Policy | Settings | target/reversion constraints | read projection |
| Response Rollback | Govern | relation only if governed response context | read ref |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Update Reversion | create/transition | Endpoint | Agent-update scope only |
| Previous-Version Recovery State | derive/update | Endpoint | restored version != full state restored |
| Post-Reversion Verification | create/link | Endpoint | technical facts only |

## 11. Fonctionnalités
Assess previous-version availability/eligibility; bind failed/current update; request/start local reversion; preserve partial/failure/unknown; observe restored version; rerun health/compatibility checks; expose Fleet-target mismatch and manual reconciliation need.

## 12. Actions utilisateur
Inspect Class 0; eligibility Class 1; reversion effect Class 3 where applicable under authority/policy; no autonomous rollback or silent target change.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| assess reversion eligibility | oui | oui if source rules | oui | explain | explicit constraints |
| verify restored version/health | oui | oui | oui | summarize | deterministic checks |
| choose/authorize rollback | accountable owner | policy/Govern | non autonome | interdit | Settings/Govern path |

## 14. États fonctionnels
`not-assessed`, `previous-version-unavailable`, `eligible`, `blocked`, `reversion-requested`, `reverting`, `reversion-partial`, `reversion-failed`, `reversion-unknown`, `previous-version-restored`, `verification-pending`, `verified-technical`, `fleet-target-mismatch`, `manual-recovery-required`.

## 15. États d’interface
No Screen ID. Reversion, restored version, health and Fleet reconciliation remain separate concepts.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Update Reversion | Endpoint technical fact | Settings/Quality/Govern consumer | not Response Rollback |
| Previous-Version Recovery State | Endpoint fact | CAP-EPT-087/093 | state limitations explicit |
| post-reversion verification | technical observation | operator/Settings | no full-state guarantee |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-EPT-087 | degraded/mismatch | CAP-EPT-088 | update/verification refs | recovery context |
| CAP-EPT-088 | previous version restored | CAP-EPT-087 | restored version/health refs | verification |
| CAP-EPT-088 | target mismatch | Settings | local vs desired target | Settings reconciliation only |

## 18. Dépendances
CAP-EPT-005/082..087/093/097; Settings upgrade management; Studio deployment rollback; Govern rollback model; OPEN-008/013/015.

## 19. Source de vérité
Endpoint SOT for local update reversion facts; Settings/Studio/Govern retain their respective administrative/asset/response rollback ownership.

## 20. Provenance et audit
Original update attempt, current/previous/target versions, eligibility inputs, authority/policy refs, reversion attempt/outcome, post-state/health and limitations.

## 21. Permissions fonctionnelles
Reversion support/read/request, verification, sensitive provenance, cross-tenant deny; step-up/SoD according Security; no final RBAC.

## 22. Limites et erreurs
Update reversion != Govern rollback; update reversion != Studio reversion; previous version restored != full endpoint state restored; technical completion != Fleet policy reconciled.

## 23. Métriques
Eligibility/blocked, partial/failure/unknown, restored-version health, target mismatch and manual-recovery counts.

## 24. Classification de livraison
`draft / defined / planned`; no rollback/recovery engine or implementation.

## 25. Critères d’acceptation
**Given** a previous version is restored, **When** Fleet still targets the newer version, **Then** local recovery and Fleet mismatch are both explicit.

**Given** reversion technically completes but health remains degraded, **When** verified, **Then** recovery is not marked healthy/full.

**Given** AI is unavailable, **When** eligibility and verification are assessed, **Then** deterministic/manual paths remain available.

## 26. Questions ouvertes
OPEN-008/013/015 remain open; final runtime bridge and action governance are not selected.

## 27. Consommateurs documentaires
EPT-6 Update/Resilience, Settings, Studio deployment boundary, Govern rollback, Security, Quality, registers and Roadmap.