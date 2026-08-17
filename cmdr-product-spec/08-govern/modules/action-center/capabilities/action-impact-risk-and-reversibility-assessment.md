---
id: CAP-GOV-005
title: Action Impact, Risk and Reversibility Assessment
product: govern
module: action-center
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-003, REQ-PROD-004, REQ-PROD-015, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-013]
source-of-truth: canonical
---
# CAP-GOV-005 — Action Impact, Risk and Reversibility Assessment

## 1. Définition
Évaluer de manière sourcée l’impact, le risque de l’action et de l’inaction, le blast radius, la réversibilité et les exigences conceptuelles de recovery/rollback d’une Action Request, sans produire de score opaque ni prendre la Decision.

## 2. Problème utilisateur
Une action techniquement possible peut avoir un impact métier, sécurité ou disponibilité disproportionné. À l’inverse, l’inaction peut aussi être risquée. Sans assessment explicite, une Decision peut masquer les hypothèses, l’incertitude et l’absence réelle de rollback.

## 3. Objectifs
- distinguer business/security/availability/confidentiality/integrity/user/tenant impact ;
- comparer risk-of-action et risk-of-inaction ;
- expliciter blast radius, duration/persistence et uncertainty ;
- évaluer reversibility, rollback availability connue et recovery expectation ;
- rendre visibles sources, reviewer et changements entre versions.

## 4. Non-objectifs
Ne pas créer un score universel, certifier la sécurité d’une action, exécuter une simulation affectant une cible, inventer un rollback, modifier un Incident/Finding ou décider.

## 5. Propriétaire
Govern / Action Center owns the governance assessment. Source-product risk/impact facts remain source-owned; future GOV-2 owns execution-time rollback/verification behavior.

## 6. Utilisateurs
Principal : Risk Reviewer. Secondaires : Govern Reviewer, Decision Maker, Business Owner/Incident Commander as context contributors, target owner, Auditor.

## 7. Conditions d’entrée
Current Action Request/version, Context/Scope/Target Review, declared effect, source impact/risk context and sufficient permission to inspect relevant target/service projections.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| action class/effect | Action Request + action classification | effect category | oui | request version | blocked/incomplete |
| exact target/scope | CAP-GOV-004 | bounded target context | oui | reviewed version | assessment partial/block |
| business/security impact | Command/source owner/requester | impact context | oui for effectful request | timestamp/source shown | unknown + information request |
| technical/user/tenant impact | target owner/Settings/source | operational context | selon action | current observation | unknown, not zero |
| reversibility/rollback context | requester/procedure/future GOV-2 knowledge | recovery candidate | required where policy says | version visible | `unknown`/not-available |
| legal/policy flags | Policy/Security source | contextual flags | non until evaluated | source version | explicitly unavailable |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Action Request | Govern | effect/action/justification | read/review |
| Context/Target Review | Govern | target/scope/blast-radius candidates | read |
| Incident/Service context | Command/Shared | business impact/criticality | restricted read |
| Finding/Evidence refs | Investigate | support/uncertainty only | read, no requalification |
| future rollback/playbook context | Govern GOV-2 / Studio refs | existence/version candidate | read projection only |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Impact Assessment | create/update/version/supersede | Govern local concept | source dimensions explicit |
| Risk Assessment | create/update/version/supersede | Govern local concept | factors visible; no opaque universal score |
| Reversibility Assessment | create/update/version/supersede | Govern local concept | known/unknown/partial distinguished |
| Action Request | attach review references / request info | Govern | does not silently change requested action |

## 11. Fonctionnalités
Capture impact dimensions; distinguish confirmed/estimated/unknown; compare action/inaction; identify candidate blast radius; record duration/persistence; review rollback/recovery prerequisites; compare assessment versions; expose uncertainty/contradictions; request missing data; prepare decision inputs.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect source impact | reviewer | source projections | 0 | read permitted | sourced context | non |
| run bounded risk/reversibility check | Risk Reviewer | Assessment | 1 | explicit factors | explainable result | non |
| create/update assessment | Risk Reviewer | assessment concepts | 2 | current request/version | versioned assessment | OPEN-013 |
| challenge/annotate uncertainty | reviewer | assessment | 2 | rationale | attributed challenge | OPEN-013 |
| execute rollback/simulation on target | none | target | 3/4 | outside GOV-1 | prohibited | future GOV-2 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| aggregate impact factors | oui | oui | oui | sourced summary | factor table |
| compare action/inaction | oui | explicit matrix | oui | narrative suggestion | comparison matrix |
| evaluate reversibility completeness | oui | checklist | oui | missing-field suggestion | rollback checklist |
| highlight contradictions | oui | deterministic diff/rules | oui | explanation | source comparison |
| choose Decision | human later | no | no | prohibited | CAP-GOV-014/015 |

## 14. États fonctionnels
`not-assessed`, `partial`, `impact-reviewed`, `risk-reviewed`, `reversibility-unknown`, `reversible-candidate`, `partially-reversible`, `irreversible-candidate`, `conflicted`, `information-required`, `reviewed`, `superseded`.

## 15. États d’interface
Loading preserves factors ; Empty means assessment not started ; Partial names missing dimensions ; Error keeps previous valid assessment ; Offline read-only ; Permission denied masks protected source ; Stale exposes source timestamps and invalidates silent “current” claim.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Impact Assessment | review result | CAP-GOV-007/009/014 | dimensions/source/uncertainty visible |
| Risk Assessment | review result | Policy/Authority/Decision review | factors visible, no universal score claim |
| Reversibility Assessment | review result | CAP-GOV-014/015/016 | rollback/recovery requirements explicit or unknown |
| information need | review event | CAP-GOV-003/source product | missing dimension named |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-004 | target/scope reviewed | CAP-GOV-005 | exact scope/targets/dependencies | context review |
| CAP-GOV-005 | policy factors ready | CAP-GOV-007 | action class/impact/risk/reversibility | assessment |
| CAP-GOV-005 | authority factors ready | CAP-GOV-009 | impact/blast radius/reversibility | assessment |
| CAP-GOV-005 | decision prep | CAP-GOV-014 | assessments/uncertainty/sources | Action Center |
| CAP-GOV-005 | missing information | CAP-GOV-003/source | question and affected dimension | assessment |

## 18. Dépendances
CAP-GOV-003/004/007/009/014/015/016, Command impact/service projections, target owners, Security policy, future GOV-2 rollback semantics, Shared Linking/Trace, OPEN-013.

## 19. Source de vérité
Govern owns the GOV-1 assessment records; source facts remain at Command/Investigate/Settings/target owners. Future rollback execution/result remains GOV-2-owned. An assessment is evidence for a Decision, not authority.

## 20. Provenance et audit
Store request/version, target/scope, factor/source/version, confirmed/estimated/unknown status, reviewer, methodology label, uncertainty, contradictions, change diffs, automation/tool provenance and human disposition.

## 21. Permissions fonctionnelles
Risk Assessment create/update/read, restricted source read, automated recommendation request, cross-tenant review if authorized, provenance export. Step-up/SoD may apply to material assessment changes but final rules are future.

## 22. Limites et erreurs
Missing telemetry, stale service context, unsupported target, absent rollback plan, contradictory business/technical impact, no recovery estimate, inaccessible Evidence or model failure keep uncertainty explicit and cannot be converted to zero risk.

## 23. Métriques
Assessments with unknown dimensions; requests missing rollback/recovery context; material changes between versions; action-vs-inaction comparisons completed; zero target-effect simulations/executions in GOV-1.

## 24. Classification de livraison
`defined` / `planned`; no definitive risk engine, threshold, score formula, simulation engine or rollback implementation selected.

## 25. Critères d’acceptation
**Given** a high-impact action with no rollback information, **When** assessment is performed, **Then** reversibility remains `unknown`/incomplete and Decision preparation cannot silently claim reversibility.

**Given** conflicting impact estimates from Command and a target owner, **When** review compares them, **Then** both sources and uncertainty remain visible; no opaque score hides the conflict.

**Given** no AI, **When** risk is assessed, **Then** explicit factor matrices, checklists and human review provide full functionality.

## 26. Questions ouvertes
OPEN-013 remains open. Final risk model, thresholds, rollback contract and object schemas remain future Security/Objects/GOV-2 decisions; no new OPEN is required by GOV-1.

## 27. Consommateurs documentaires
Action Center, Policy/Authority/Decision capabilities, Execution Handoff, future Objects/Permissions/GOV-2/Screens/Journeys/Technique, quality gates and reports.
