---
id: CAP-GOV-007
title: Policy Applicability and Evaluation
product: govern
module: policy-gates
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-004, REQ-PROD-015, REQ-PROD-019, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-013]
source-of-truth: canonical
---
# CAP-GOV-007 — Policy Applicability and Evaluation

## 1. Définition
Identifier les Policies candidates applicables à une Action Request/version, évaluer leurs conditions de manière sourcée et produire des outcomes fonctionnels `pass`, `warn`, `block`, `unknown` ou `not-applicable` sans transformer un outcome en Decision.

## 2. Problème utilisateur
Une Policy peut être versionnée, limitée à un tenant, un environnement, un type d’action, une période ou un scope. Sans applicability explicite, un système peut appliquer la mauvaise version, masquer un `unknown` ou confondre `block` avec une Decision de rejet.

## 3. Objectifs
- identifier Policies candidates et leurs owners/versions/effective periods ;
- déterminer applicability, exclusions et données manquantes ;
- exécuter/rejouer une évaluation déterministe quand possible ;
- exposer outcome, reasons, warnings, blocks, unknowns et provenance ;
- permettre dispute/review sans muter la Policy source.

## 4. Non-objectifs
Ne pas définir le moteur de Policy, le langage de règles, l’authoring/activation administratif, la priorité finale entre Policies, l’exception active, la Decision ou l’exécution.

## 5. Propriétaire
Govern / Policy Gates owns Policy Evaluation and its review semantics. Policy is Govern-owned conceptually; Security/administrative configuration may provide rule/authority sources. GOV-1 does not create a final Policy administration engine.

## 6. Utilisateurs
Principal : Policy Reviewer. Secondaires : Govern Reviewer, Risk Reviewer, Authority Reviewer, Decision Maker, Policy Owner, Auditor.

## 7. Conditions d’entrée
Action Request/version, reviewed context/scope/target, action class, relevant impact/risk facts, tenant/environment and access to candidate Policy metadata/version/scope.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Action Request/version | Govern | evaluation subject | oui | exact review version | no evaluation |
| context/scope/target | CAP-GOV-004 | applicability dimensions | oui | reviewed version | `unknown`/blocked review |
| impact/risk/reversibility | CAP-GOV-005 | conditional factors | according to Policy | current assessment | `unknown` for dependent rule |
| completeness/known unknowns | CAP-GOV-006 | input quality | oui | current | warnings/unknown preserved |
| Policy candidate/version/scope | Govern/Security source | rule metadata | oui for each evaluation | effective version at evaluation time | no silent substitute |
| exclusions/effective period | Policy source | applicability limits | according to Policy | versioned | unknown/not-applicable only if justified |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Action Request | Govern | action/target/scope/version | read |
| Policy | Govern | version/scope/effective period/conditions/exclusions | read/evaluate |
| Context/Risk/Completeness reviews | Govern | factual inputs/unknowns | read |
| Principal/Tenant/Environment | Settings | scope identity | minimal read |
| source objects | Command/Investigate/etc. | facts referenced by Policy | restricted read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Policy Evaluation | create/run/review/version/dispute | Govern | bound to request+Policy versions |
| Policy applicability record | create/update | Govern local concept | explains applicable/not-applicable/unknown |
| Action Request review context | attach evaluation refs | Govern | no Decision inferred |
| Policy | no lifecycle/authoring mutation in GOV-1 | Govern future admin capability | evaluation only |

## 11. Fonctionnalités
Discover candidate Policies; show owner/version/scope/effective period; evaluate applicability; evaluate conditions; distinguish missing data; produce pass/warn/block/unknown/not-applicable; show reasons and input snapshots/refs; rerun after request change; compare versions; dispute/annotate; route conflicts to CAP-GOV-008.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| view Policy/version | reviewer | Policy | 0 | read | source metadata/conditions | non |
| run Policy Evaluation | Policy Reviewer | Policy Evaluation | 1 | request/context + policy version | explainable outcome | non |
| compare/replay evaluation | reviewer | Evaluation versions | 1 | preserved inputs | diff/result | non |
| dispute/annotate outcome | reviewer | Policy Evaluation | 2 | rationale | attributed review state | OPEN-013 |
| edit/activate Policy | none in GOV-1 | Policy lifecycle | — | outside lot | no action | future policy admin/technique |
| reject request automatically | none | Decision | — | forbidden | no Decision | Decision Maker later |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| find candidate Policies | oui | catalog/scope matching | oui | suggestions | Policy catalog/filter |
| evaluate explicit conditions | oui | yes where rules defined | oui | explain outcome | deterministic evaluator/checklist |
| summarize warnings/blocks | oui | structured reasons | oui | sourced summary | outcome table |
| detect missing inputs | oui | rule input checks | oui | suggestion | dependency checklist |
| decide/reject/approve | human Decision capability | validation only | never autonomous | prohibited | CAP-GOV-014/015 |

## 14. États fonctionnels
Policy Evaluation: `candidate`, `applicability-review`, `applicable`, `not-applicable`, `evaluating`, `pass`, `warn`, `block`, `unknown`, `disputed`, `superseded`, `stale`.

`block` means the Policy evaluation blocks progression under that Policy context; it is not automatically a Decision disposition of `reject`.

## 15. États d’interface
Loading preserves request/policy versions ; Empty means no candidate Policy found vs no access ; Partial lists unevaluated conditions ; Error preserves valid evaluations ; Offline allows read of cached provenance but no authoritative rerun ; Permission denied masks protected Policy/source ; Stale requires reevaluation after input/version change.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Policy Evaluation | versioned review result | CAP-GOV-008/009/014/015 | Policy/request/input versions and reasons preserved |
| applicability disposition | evaluation context | Policy Gates/Action Center | applicable ≠ satisfied |
| warning/block/unknown set | structured outcomes | reviewer/Decision prep | outcome ≠ Decision or safety claim |
| conflict candidate | handoff context | CAP-GOV-008 | conflicting Policies/versions/reasons visible |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Action Request | policy review required | CAP-GOV-007 | request/context/risk/completeness | Action Center |
| CAP-GOV-007 | conflict/multiple incompatible outcomes | CAP-GOV-008 | Policies/versions/outcomes/reasons | Policy Gates |
| CAP-GOV-007 | evaluation complete | CAP-GOV-009/014 | outcomes/unknowns/warnings/blocks | Policy Gates/Action Center |
| request/context version changes | reevaluation required | CAP-GOV-007 | new version + previous evaluation ref | source review |

## 18. Dépendances
CAP-GOV-003..006/008/009/014/015, Policy object, Security permission/authority sources, Settings tenant/env, future Policy engine/language, Shared Versioning/Trace/Linking, OPEN-013.

## 19. Source de vérité
Policy/version remains canonical Govern/Security-managed source; Govern Policy Evaluation is source of evaluation outcome for the exact request/input versions. A Policy engine output is not a Decision.

## 20. Provenance et audit
Record Policy id/version/owner/scope/effective period, request/version, input refs/versions, evaluation method/tool, outcome/reasons/missing data, reviewer/dispute, timestamps, reruns, supersession, model/tool provenance and correlation ids.

## 21. Permissions fonctionnelles
Policy read; Policy Evaluation run/read; review/dispute; restricted input read; automated recommendation request; provenance export. Policy authoring/activation and final permission namespaces remain future.

## 22. Limites et erreurs
Missing Policy version, unknown condition input, stale request/context, engine unavailable, conflicting scope, permission denial, ambiguous effective period or Policy owner unavailable must yield explicit unknown/partial/error; never silent pass.

## 23. Métriques
Evaluations by outcome; unknown/warn/block rates and causes; reevaluations after request changes; disputes; Policies with missing input dependencies; zero Policy outcomes converted automatically into Decisions.

## 24. Classification de livraison
`defined` / `planned`; target Policy evaluation capability without selecting engine, language, API, protocol or runtime.

## 25. Critères d’acceptation
**Given** an applicable Policy whose explicit condition evaluates to block, **When** evaluation completes, **Then** outcome is `block` with reason/source, progression is controlled by governance, and no automatic rejection Decision is created.

**Given** two candidate Policies with incompatible requirements, **When** both evaluate, **Then** conflict is surfaced and routed to CAP-GOV-008 rather than silently selecting one.

**Given** no AI, **When** Policy review occurs, **Then** catalogs, deterministic conditions, checklists and human review provide the full path.

## 26. Questions ouvertes
OPEN-013 remains open. Final Policy engine/language/precedence/administration may require future technical/product decisions, but GOV-1 does not create a new OPEN where existing roadmap/security dependencies already capture them.

## 27. Consommateurs documentaires
Policy Gates, Action Center, Authority/Approval/Decision capabilities, future Objects/Permissions/Screens/Journeys/Technique and GOV-1 conformance report.
