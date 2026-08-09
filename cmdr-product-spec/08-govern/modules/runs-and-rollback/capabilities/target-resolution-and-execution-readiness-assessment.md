---
id: CAP-GOV-020
title: Target Resolution and Execution Readiness Assessment
product: govern
module: runs-and-rollback
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-004, REQ-PROD-006, REQ-PROD-008, REQ-PROD-015, REQ-PROD-017, REQ-PROD-019, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-GOV-020 — Target Resolution and Execution Readiness Assessment

## 1. Définition
Résoudre les Target References d’un Execution Plan vers des cibles courantes et évaluer leur **execution readiness** à partir de tenant/environment, identity, state/freshness, reachability projection, executor capability/availability, policy/version/maintenance restrictions et drift depuis la Decision, sans modifier la cible ni étendre le scope.

## 2. Problème utilisateur
Une cible approuvée peut avoir changé, disparu, être devenue ambiguë ou ne plus être accessible par l’executor prévu. Traiter une référence ancienne comme une cible prête peut conduire à agir sur un autre objet ou au-delà du périmètre autorisé.

## 3. Objectifs
- distinguer Target Reference et Resolved Target ;
- vérifier identité, tenant/environment et exact scope ;
- exposer freshness, availability, capability and health ;
- détecter drift, ambiguity, duplicates and missing targets ;
- qualifier readiness sans prétendre à un succès futur ;
- bloquer toute expansion silencieuse du target set.

## 4. Non-objectifs
Ne pas modifier la cible, administrer Fleet/runtime, effectuer une action de probe intrusive, élargir la Decision, conclure qu’une cible résolue est automatiquement approuvée, démarrer un Run ou choisir un executor global.

## 5. Propriétaire
Govern owns request-specific target-resolution/readiness assessment. Source systems/Settings own administrative identity/configuration, Endpoint/runtime owners own technical capability/health, and the target’s source owner remains authoritative for identity/state facts.

## 6. Utilisateurs
Principal : Response Operator / Govern Reviewer. Secondaires : Endpoint/Runtime Operator, Platform Administrator as source owner, Decision Maker, Security Reviewer, Auditor.

## 7. Conditions d’entrée
Execution Plan from CAP-GOV-019, exact Decision/Handoff target bounds, access to authorized target/source metadata and declared executor/runtime capability projections.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| target references + approved bounds | Decision/Handoff/Plan | intended target set | oui | exact plan version | blocked |
| target identity/current state | target/source owner | resolution facts | oui | timestamp visible | ambiguous/unavailable |
| tenant/environment | Settings/source | isolation boundary | oui | current | permission-blocked/unknown |
| executor capability/availability | Endpoint/Studio/provider owner | technical readiness projection | oui for effectful action | latest known | runtime-unavailable/unsupported |
| integration/runtime health | Platform Settings | dependency health | selon executor | latest timestamp | readiness partial/blocked |
| maintenance/policy/version restrictions | Settings/Security/source | execution restrictions | selon target | effective version | unknown/re-review |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Execution Plan | Govern | exact target refs/scope/executor needs | read |
| Decision/Handoff | Govern | approved targets/scope/conditions | read |
| Endpoint/Fleet/Environment/Integration | Settings/Endpoint | identity/capability/health projection | restricted read |
| Workflow/Tool deployment | CMDR Studio | executor availability ref | read/link |
| target object | source owner | stable identity/current state/freshness | read only |
| Policy/Exception context | Govern/Security | current restrictions | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Resolved Target projection | resolve/version/mark ambiguity/drift | Govern local projection | no target mutation |
| Readiness Assessment | create/rerun/version/supersede | Govern local concept | exact plan/target source versions |
| Execution Plan | attach readiness refs only | Govern | no silent target-set rewrite |
| target/Fleet/integration | no mutation | source owners | assessment only |

## 11. Fonctionnalités
Resolve stable identity; compare requested vs resolved target; validate tenant/environment; check target freshness/current state; consume reachability and executor capability projections; check executor/runtime health; check version/policy/maintenance restrictions; detect duplicates/missing/ambiguous targets and drift; classify each target and overall bounded readiness; require re-review/re-decision where drift affects authorization.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect target/current status | operator | target projection | 0 | read | current context | non |
| run deterministic readiness checks | reviewer | Readiness Assessment | 1 | source metadata | explainable readiness | non |
| annotate/accept known limitation for review | reviewer | assessment | 2 | policy permits + rationale | bounded limitation | OPEN-013 |
| request target correction/review | reviewer | Plan/Decision path | 2 | drift/ambiguity | re-review event | OPEN-013 |
| probe/mutate target/start execution | none here | target | — | outside capability | no effect | executor/GOV-2 later |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| resolve stable identifiers | oui | resolver/link rules | oui | no need | deterministic resolution |
| compare approved vs current target | oui | exact diff | oui | explanation | diff table |
| aggregate health/capability facts | oui | structured refs | oui | sourced summary | readiness matrix |
| flag anomalous drift | oui | deterministic changes | oui | suggestion | version/state comparison |
| expand scope/authorize/start | no | no | no | prohibited | CAP-GOV-021..023 |

## 14. États fonctionnels
Per target / aggregate outcomes include `ready`, `ready-with-limitations`, `stale`, `drifted`, `ambiguous`, `unavailable`, `unsupported`, `permission-blocked`, `runtime-unavailable`, `re-review-required`, `re-decision-required`, `superseded`. Readiness is not success.

## 15. États d’interface
Loading preserves target list ; Empty means no resolvable target vs no access ; Partial lists unresolved/limited targets ; Error keeps last valid resolution but marks freshness ; Offline never declares fresh readiness ; Permission denied avoids leaking target existence where required ; Stale requires rerun before start.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Resolved Target set | bounded projections | CAP-GOV-021/022/023 | source id + current identity + approved relation visible |
| Readiness Assessment | versioned assessment | CAP-GOV-021/022 | per-target outcomes, sources, timestamps |
| drift/mismatch blocker | governance event | GOV-1/CAP-GOV-021 | no scope expansion |
| dependency/runtime blocker | readiness event | operator/Settings/Studio/Endpoint | exact unavailable dependency named |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-019 | plan ready for resolution | CAP-GOV-020 | targets/scope/executor/dependency refs | Execution Plan |
| CAP-GOV-020 | ready/limited | CAP-GOV-021 | resolved set + readiness + timestamps | readiness review |
| CAP-GOV-020 | target drift impacts authority | CAP-GOV-021/GOV-1 | approved vs current diff | same plan/handoff |
| dependency owner | state changes | CAP-GOV-020 | new health/capability/version | rerun assessment |

## 18. Dépendances
CAP-GOV-015..019/021/022/023, target source owners, Settings Fleet/environments/integrations/health, Endpoint technical capabilities, Studio deployments, Security/Policy context, Shared Linking/Versioning/Trace, OPEN-008/013/015.

## 19. Source de vérité
Source owners remain authoritative for target identity/state and executor capability. Govern owns the resolved-target relation to the Decision and the readiness assessment for the current plan. Readiness never becomes source of target truth or execution success.

## 20. Provenance et audit
Record plan/Decision versions, original refs, resolved ids, resolution source/time, tenant/environment, state/freshness, capability/health sources, restrictions, approved-vs-current diff, readiness outcome/reason, reviewer, reruns, automation/AI provenance and timestamps.

## 21. Permissions fonctionnelles
Target resolve/read, readiness run/read, restricted target/runtime metadata, cross-tenant review only if explicitly authorized, automated recommendation request and provenance export. No target/Fleet/runtime administration is granted.

## 22. Limites et erreurs
Missing target, duplicate/ambiguous identity, stale source, cross-tenant mismatch, executor unavailable, unsupported action, maintenance restriction, permission denial or conflicting current state yields explicit blocked/limited/re-review status; no guessed target or silent scope change.

## 23. Métriques
Targets by readiness state; drift before start; unresolved/ambiguous rate; runtime-unavailable blocks; stale reassessments; scope expansions prevented; readiness-passed executions later failing verification tracked by downstream metrics, not reclassified here.

## 24. Classification de livraison
`defined` / `planned`; no target resolver implementation, network probe, provider, runtime or command selected.

## 25. Critères d’acceptation
**Given** the approved Target Reference now resolves to a different current resource identity, **When** readiness runs, **Then** drift is visible, scope is not expanded and review/re-decision is required where authority changed.

**Given** an executor is unavailable, **When** readiness is assessed, **Then** the target/Run preparation cannot claim started or ready without limitation and the exact dependency is recorded.

**Given** no AI, **When** readiness is assessed, **Then** deterministic identity/version/health comparisons and human review provide full functionality.

## 26. Questions ouvertes
OPEN-008 remains open for actual platform/executor support, OPEN-013 for class-2 assessment dispositions, OPEN-015 for run bridge. No new OPEN is required.

## 27. Consommateurs documentaires
Runs & Rollback, CAP-GOV-021/022/023/025, Decision/Playbook reviewers, Settings/Endpoint/Studio source owners, future Objects/Permissions/Screens/Technique and GOV-2 report.
