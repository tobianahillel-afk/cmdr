---
id: CAP-GOV-030
title: Rollback Eligibility, Planning and Preconditions
product: govern
module: runs-and-rollback
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-004, REQ-PROD-005, REQ-PROD-008, REQ-PROD-015, REQ-PROD-016, REQ-PROD-020, REQ-OBJ-007, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-GOV-030 — Rollback Eligibility, Planning and Preconditions

## 1. Définition
Évaluer si un Response Run ou une partie de ses effets est éligible à un rollback, puis préparer un **Rollback Plan** explicitant original Run/Decision, rollback action candidate, exact target/scope, trigger, authority, conditions, expected restored state, dependencies, risk/data-loss, executor, verification, recovery alternative et limitations avant toute tentative de rollback.

## 2. Problème utilisateur
Le fait qu’un Playbook ou executor expose une primitive inverse ne signifie ni que le rollback est encore disponible, ni qu’il est sûr, ni qu’il restaure réellement l’état initial. Un rollback mal borné peut amplifier l’impact ou agir après expiration sans autorité explicite.

## 3. Objectifs
- relier rollback à l’original Run, Decision, exact effects and targets ;
- distinguer supported, eligible, safe-enough-for-review and authorized ;
- identifier rollback action candidate et technical executor capability ;
- documenter exact scope, expected restored state, limitations and potential data loss ;
- vérifier authority, trigger, conditions, dependencies and current target state ;
- prévoir post-rollback verification and recovery alternative ;
- produire un plan sans effet pour CAP-GOV-031.

## 4. Non-objectifs
Ne pas garantir la sécurité du rollback, exécuter une primitive, définir une commande/API/protocole, réautoriser automatiquement une Decision, étendre le scope, considérer rollback support comme succès futur, affirmer une restauration exacte ou finaliser la state machine physique.

## 5. Propriétaire
Govern owns rollback eligibility, governance and Rollback Plan semantics. Endpoint/Studio/provider owners retain technical rollback primitives and raw outcomes. Settings owns runtime configuration/secrets. Verification remains Govern-owned through CAP-GOV-028/029/031.

## 6. Utilisateurs
Principal : Response Operator / Rollback Reviewer. Secondaires : Govern Reviewer, Decision Maker, Security/Authority Reviewer, Endpoint/Runtime Operator, Studio Operator, Business/Service Owner as impact contributor, Auditor.

## 7. Conditions d’entrée
Existing Response Run with known/partial effect state, verification/error context indicating rollback review or an explicit authorized rollback requirement, current target state/readiness and visibility into technical rollback support.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| original Run/Decision/Plan | Govern | rollback lineage/authority context | oui | exact pinned versions | no rollback plan |
| affected targets/effects | CAP-GOV-026/027/029 | rollback subject | oui | latest known verified/technical state | scope unknown/block |
| rollback support/action candidate | Playbook/Endpoint/Studio/provider | technical capability metadata | oui for executable rollback | exact capability/version | `not-supported`/unknown |
| rollback authority/conditions | Decision/Policy/Govern authority context | authorization need | oui for effectful rollback | current/effective | authority-required/blocked |
| current target/readiness | CAP-GOV-020/source owner | precondition/current state | oui | fresh enough | re-review/block |
| risk/data-loss/reversibility context | CAP-GOV-005/029 + technical metadata | consequence context | selon rollback | current | unknown explicit |
| verification/recovery alternative | CAP-GOV-028/Playbook/source | post-action safety plan | oui where rollback proceeds | current version | plan incomplete |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Response Run / Step / Result-draft context | Govern | original effects/targets/outcomes | read |
| Decision / conditions / authority context | Govern | rollback authority/bounds | read/reconcile |
| Verification/Residual Risk Assessment | Govern | rollback trigger/need | read |
| Response Rollback | Govern | historical/current rollback relation | read/manage semantics |
| Playbook / Workflow / technical rollback capability | Govern / Studio / Endpoint/provider | candidate reverse action/capability | read/link only |
| target/environment/Secret Reference metadata | source owner/Settings | current readiness/input references | restricted read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Rollback Plan | create/update/version/supersede | Govern local concept | exact original effect/target/scope retained |
| Rollback Eligibility Assessment | create/rerun/version | Govern local concept | supported ≠ eligible ≠ safe ≠ authorized |
| Response Rollback relation | prepare/link | Govern | plan ≠ execution |
| original Run/Decision/target | no rewriting | respective owner | historical effect/authority preserved |

## 11. Fonctionnalités
Determine whether rollback is supported; identify reverse-action candidate; correlate exact affected targets/effects; compare current target state to expected rollback preconditions; evaluate time/authority/Decision conditions; record trigger; set exact rollback scope; document expected restored state and what is not reverted; expose data-loss/business/security risks and limitations; identify executor/dependencies/Secret References; define post-rollback verification and recovery fallback; mark manual/escalation path when unsupported or unsafe.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect rollback support/history | reviewer | eligibility/rollback refs | 0 | read | rollback context | non |
| run eligibility/precondition check | reviewer | Eligibility Assessment | 1 | current sources | eligible/blocked/unknown reasons | non |
| create/update Rollback Plan | Rollback Reviewer | Rollback Plan | 2 | exact original effects/scope | versioned no-effect plan | OPEN-013 |
| request rollback authority/review | reviewer | governance context | 2 | authority missing/changed | explicit review path | OPEN-013 |
| execute rollback | authorized operator/path | Response Rollback | 3/4 per effect | CAP-GOV-031 authority/preconditions | separate execution | governed |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| match effect to rollback capability | oui | catalog/version rules | oui | candidate suggestion | capability matrix |
| check target/precondition/authority | oui | explicit rules | oui | explanation | checklist/diff |
| estimate limitations/data-loss questions | oui | metadata/rules | oui | recommendation only | risk checklist |
| draft Rollback Plan/recovery alternative | oui | templates | oui | draft | structured form |
| authorize/start rollback | accountable governed path | contract-controlled | only when explicitly authorized | prohibited as autonomous AI | explicit controls |

## 14. États fonctionnels
`not-assessed`, `not-supported`, `supported`, `eligibility-unknown`, `eligible`, `eligible-with-limitations`, `preconditions-failed`, `authority-required`, `authority-invalid`, `scope-mismatch`, `target-drifted`, `plan-draft`, `plan-ready`, `manual-recovery-required`, `re-decision-required`, `superseded`, `cancelled`. Rollback available does not mean safe or successful.

## 15. États d’interface
Loading preserves original Run/effect context ; Empty means no rollback assessment ; Partial names unknown target/capability/risk inputs ; Error preserves prior valid assessment ; Offline prohibits effectful transition and fresh eligibility claims ; Permission denied masks sensitive technical inputs ; Stale requires current target/authority recheck.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Rollback Eligibility Assessment | versioned review | CAP-GOV-031/032 | support/eligibility/authority/preconditions distinguished |
| Rollback Plan | no-effect plan | CAP-GOV-031 | exact original Run/effects/targets/scope/limits/verification |
| manual/recovery requirement | governance condition | CAP-GOV-031/033/operator | unsupported/unsafe rollback not hidden |
| re-decision/authority requirement | governance event | GOV-1/CAP-GOV-021 | no silent reauthorization |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-027/029 | rollback review required | CAP-GOV-030 | original Run/effects/targets/errors/verification/residual risk | Runs & Rollback |
| CAP-GOV-030 | plan ready + authority valid | CAP-GOV-031 | exact Rollback Plan/Eligibility/authority refs | rollback review |
| CAP-GOV-030 | authority/scope changed | CAP-GOV-021/GOV-1 | mismatch and required review | same Run |
| CAP-GOV-030 | unsupported/unsafe | CAP-GOV-031/033 | manual recovery/escalation context | same Run |

## 18. Dépendances
CAP-GOV-005/015/020/021/022/026..029/031..033, Response Rollback object, Studio/Endpoint/provider rollback capability, Settings secrets/health, Shared Trace/Recovery, Security authority, OPEN-008/013/015.

## 19. Source de vérité
Govern owns rollback eligibility and plan semantics. Technical owners remain source of actual rollback capability/current technical constraints. Original Decision/Run remain historical source records and are never rewritten by planning.

## 20. Provenance et audit
Record original Run/Decision/Plan/Playbook versions, affected target/effect refs, verification/error trigger, technical rollback support/version, current target state, authority/condition checks, exact rollback scope, expected restored state, non-reverted elements, risk/data-loss/limitations, executor/Secret Reference ids, verification/recovery plan, reviewer, AI recommendation/disposition and timestamps.

## 21. Permissions fonctionnelles
Rollback plan create/update/read; eligibility run/read; technical rollback capability/restricted metadata read; rollback authority request/reconcile; automated recommendation request; provenance export. Rollback execution is separately controlled by CAP-GOV-031 and source executor permission.

## 22. Limites et erreurs
Unsupported primitive, irreversible effect, current-state drift, unknown effect state, expired authority, unavailable secret/executor, potential destructive data loss, target mismatch or missing verification path must remain explicit and can require manual recovery/re-decision. No unsafe rollback is inferred from capability presence.

## 23. Métriques
Rollback support/eligibility distribution; plans blocked by authority/drift/data-loss; rollback-review frequency; unsupported/manual-recovery cases; plans with undefined verification; automatic rollback from capability presence — target zero.

## 24. Classification de livraison
`defined` / `planned`; no rollback command, executor/API/protocol/provider or state-restoration algorithm selected.

## 25. Critères d’acceptation
**Given** a technical executor advertises rollback support but the current target drifted since the original Run, **When** eligibility is assessed, **Then** rollback is not automatically authorized and the drift/required review is explicit.

**Given** rollback is unavailable for an adverse verified effect, **When** planning occurs, **Then** no false recovery claim is made and manual recovery/escalation options plus residual risk remain visible.

**Given** a rollback may cause data loss, **When** the Plan is prepared, **Then** the risk and expected/not-restored state are explicit before any effectful start.

**Given** no AI, **When** rollback is planned, **Then** capability matrices, precondition/authority checks, risk forms and human review provide full functionality.

## 26. Questions ouvertes
OPEN-008 covers actual rollback support, OPEN-013 class-2 Plan governance and OPEN-015 runtime bridge. Final rollback authority/default trigger and technical implementation remain future; no new OPEN.

## 27. Consommateurs documentaires
Runs & Rollback, CAP-GOV-031..033, GOV-1 authority/Decision flows, Studio/Endpoint/Settings technical owners, future Objects/Permissions/Screens/Technique, GOV-2 report and GOV-3 audit/metrics.
