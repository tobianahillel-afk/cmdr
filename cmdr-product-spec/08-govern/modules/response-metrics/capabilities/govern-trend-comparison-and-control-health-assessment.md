---
id: CAP-GOV-046
title: Govern Trend, Comparison and Control Health Assessment
product: govern
module: response-metrics
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-002, REQ-PROD-005, REQ-PROD-006, REQ-PROD-015, REQ-PROD-019, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-013, OPEN-019]
source-of-truth: canonical
---
# CAP-GOV-046 — Govern Trend, Comparison and Control Health Assessment

## 1. Définition
Comparer dans le temps et entre scopes les observations GOV-3 afin d’identifier des trends et produire un **Control Health Assessment** sourcé (`healthy-observed`, `degraded-observed`, `incomplete-data`, `review-required`, `unknown`) sans transformer une anomalie ou corrélation en causalité ou défaillance certaine.

## 2. Problème utilisateur
Une variation peut venir d’un changement de Policy version, d’action mix, de tenant, de source coverage ou d’une vraie dégradation. Sans contexte et comparabilité, un dashboard transforme des différences statistiques en conclusions de contrôle.

## 3. Objectifs
- comparer period-over-period, tenant/environment/action class ;
- comparer Policy/Playbook versions et outcome patterns ;
- intégrer policy/approval/Decision/Run/verification/rollback/outcome/flow observations ;
- exposer definition changes, privacy suppression et source coverage ;
- produire un Control Health Assessment avec rationale/uncertainty ;
- distinguer trend, anomaly candidate, hypothesis et conclusion validée ailleurs.

## 4. Non-objectifs
Ne pas créer un anomaly/ML engine, déclarer automatiquement control failure, fixer benchmark universel, classer des personnes, modifier Policy/Decision/Run, devenir source of truth ou appliquer une amélioration.

## 5. Propriétaire
Govern / Response Metrics owns cross-metric interpretation and Control Health Assessment. Source metric capabilities retain their definitions; Shared owns generic Metrics/Reporting; product owners retain their own KPIs/technical metrics.

## 6. Utilisateurs
Principal : Govern Control Reviewer. Secondaires : Govern Product Lead, Security/Compliance Reviewer, Policy/Authority/Response reviewers, Platform/Studio/Endpoint owners receiving feedback and authorized leadership.

## 7. Conditions d’entrée
At least two comparable snapshots or a defined current-vs-baseline comparison exist; metric definitions/dimensions/source coverage are known enough to establish comparability or explicitly mark it unavailable.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| GOV metric observations | CAP-GOV-039..045 | comparison inputs | oui | pinned snapshots | no assessment |
| definition/version metadata | Shared/GOV metrics | comparability | oui | exact versions | incomparable/unknown |
| tenant/environment/action class | Settings/Govern context | dimensions | according to comparison | snapshot | dimension unavailable |
| Policy/Playbook versions | GOV-1/GOV-2 | change context | when relevant | pinned | explanation limited |
| source coverage/gaps | CAP-GOV-037 | quality context | recommended | same snapshot | uncertainty increased |
| privacy thresholds/masking | Shared/Security | comparison restriction | when sensitive | current policy | suppress/aggregate |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Metric Observations | Govern | values/definitions/dimensions | read/compare |
| Audit Completeness Assessment | Govern | source-quality limitations | read |
| Policy/Playbook versions | Govern | contextual change refs | read |
| technical/product metrics | source product | contextual projection only | restricted read |
| Shared metric/report snapshot | Shared | generic comparison metadata | consume |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Trend Assessment | create/update/version/supersede | Govern local concept | trend ≠ causal explanation |
| Control Health Assessment | create/review/dispute/supersede | Govern local concept | observed state + rationale/limitations |
| source metrics/policies | no mutation | source owner | interpretation only |

## 11. Fonctionnalités
Select comparable periods/scopes; normalize only when definitions allow; display definition/source changes; compare policy/exception/approval/Decision/Run/verification/rollback/outcome/flow trends; support privacy-safe tenant comparisons; flag candidate regressions/improvements; attach possible explanations as hypotheses; create Control Health Assessment; compare assessment versions; route review/improvement candidates.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect/compare trends | reviewer | observations | 0 | read permission | comparison | non |
| run deterministic comparison | reviewer/system | Trend Assessment | 1 | comparable snapshots | differences/limitations | non |
| create Control Health Assessment | reviewer | assessment | 2 | sourced comparison + rationale | versioned assessment | OPEN-013 |
| dispute/supersede assessment | reviewer | assessment | 2 | evidence/rationale | preserved history | OPEN-013 |
| change control/policy automatically | none | source | 3/4 | prohibited | no action | owner review only |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| period/dimension comparison | oui | oui | oui | explain | comparison tables |
| detect definition/source changes | oui | oui | oui | summary | metadata diff |
| flag anomaly/trend candidate | oui | explicit statistical/rule method future | possible | candidate only | reviewer-selected comparisons |
| draft Control Health rationale | oui | templates | oui | sourced draft | structured assessment form |
| declare causal control failure | human evidence-based review | no | no | prohibited | explicit investigation/review |

## 14. États fonctionnels
Trend: `not-assessed`, `comparable`, `partially-comparable`, `incomparable`, `candidate-change`, `reviewed`, `superseded`. Control Health: `healthy-observed`, `degraded-observed`, `incomplete-data`, `review-required`, `unknown`, `disputed`, `superseded`.

## 15. États d’interface
Loading retains comparison pair ; Empty means no comparable observations ; Partial identifies missing dimensions/sources ; Error keeps prior assessment ; Offline read-only ; Permission denied suppresses sensitive comparisons ; Stale shows latest snapshot/definition version.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Trend Assessment | comparison object | CAP-GOV-047/Reporting | definitions/scopes/limitations explicit |
| Control Health Assessment | review object | CAP-GOV-047/closure | observed status not automatic truth |
| anomaly/regression candidate | derived observation | review/improvement | candidate only |
| period/scope comparison | comparison artifact | authorized reviewer | privacy/source coverage visible |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-039..045 | comparison requested | CAP-GOV-046 | metric refs/definitions/snapshots | Response Metrics |
| CAP-GOV-037 | source quality issue | CAP-GOV-046 | gap/coverage limitations | Audit Trail |
| CAP-GOV-046 | improvement candidate | CAP-GOV-047 | assessment/evidence/hypotheses | Response Metrics |
| CAP-GOV-046 | report/export | Shared Reporting/Export | assessment refs/classification | Response Metrics |

## 18. Dépendances
CAP-GOV-037/039..045/047, Shared Metrics/Reporting, Settings tenant/environment context, source product metrics, Security privacy, OPEN-008/013/019.

## 19. Source de vérité
Underlying metrics remain source for observations; Govern owns Trend/Control Health interpretations. A dashboard or assessment is not a replacement for canonical source objects and never changes them.

## 20. Provenance et audit
Record compared snapshot IDs, definitions/versions, periods, dimensions, exclusions/suppressions, source coverage, comparison method, detected differences, hypotheses, reviewer rationale, AI assistance and assessment supersession.

## 21. Permissions fonctionnelles
Metric read/comparison; sensitive dimension read; cross-tenant comparison separately authorized; Trend/Control Health create/review; Reporting/Export preparation. No source mutation or external-share right implied.

## 22. Limites et erreurs
Definition drift, policy/playbook changes, source coverage differences, seasonality, sparse data, privacy suppression and mixed populations can make trends incomparable. Anomaly is not control failure; trend is not causal explanation.

## 23. Métriques
Comparison coverage; assessments by health state; incomparable comparisons; candidate regressions/improvements; source-gap-associated changes; disputed/superseded assessments. No universal benchmark.

## 24. Classification de livraison
`defined` / `planned`; no anomaly engine, ML model, benchmark, dashboard or causal model selected.

## 25. Critères d’acceptation
**Given** exception rate rises after a Policy version change, **When** trends are reviewed, **Then** both facts are visible and causation is not asserted automatically.

**Given** two tenants have different privacy/source coverage, **When** compared, **Then** suppression/coverage limitations are explicit and no unfair ranking is produced.

**Given** AI flags an anomaly, **When** Control Health is reviewed, **Then** it remains a candidate with source metrics and can be rejected/disputed by a human.

**Given** no AI, **When** trends/control health are assessed, **Then** deterministic comparisons and structured human review remain complete.

## 26. Questions ouvertes
OPEN-008/013/019 remain open. Final statistical methods, thresholds, benchmarks and cross-tenant sharing rules remain future; no new OPEN.

## 27. Consommateurs documentaires
Response Metrics, CAP-GOV-047, Govern closure, Shared Reporting, product-owner improvement loops and future screen/technical validation.