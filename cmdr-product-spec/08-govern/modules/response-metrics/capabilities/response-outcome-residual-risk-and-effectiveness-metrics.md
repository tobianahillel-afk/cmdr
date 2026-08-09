---
id: CAP-GOV-044
title: Response Outcome, Residual Risk and Effectiveness Metrics
product: govern
module: response-metrics
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-002, REQ-PROD-005, REQ-PROD-008, REQ-PROD-015, REQ-PROD-020, REQ-OBJ-007, REQ-SEC-001]
open_decisions: [OPEN-008, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-GOV-044 — Response Outcome, Residual Risk and Effectiveness Metrics

## 1. Définition
Définir les métriques Govern sur les outcomes du canonical Result, le residual risk, les unintended effects et les follow-up needs afin d’observer l’efficacité documentée de la réponse sans assimiler `Result success` à business value, `failure` à Decision error ou consumer feedback à ground truth.

## 2. Problème utilisateur
Un Result `success` peut atteindre son objectif technique sans éliminer le risque métier ; un `failed` peut résulter d’un runtime indisponible sans invalider la Decision. Sans distinction, une métrique « success rate » devient un jugement causal erroné.

## 3. Objectifs
- observer Result outcomes `success/partial-success/failed/cancelled/stopped/rolled-back/recovered/verification-failed/inconclusive` ;
- mesurer residual-risk concepts et unintended effects ;
- observer follow-up Investigation/Command actions et repeated-action candidates ;
- relier verification/rollback context sans écraser les sources ;
- comparer outcomes par action class/Playbook/tenant/environment/période ;
- exposer limitations et feedback comme données contextuelles, pas vérité absolue.

## 4. Non-objectifs
Ne pas créer un score universel d’effectiveness, conclure business value, reclasser Result automatiquement, modifier Decision/Finding/Evidence, lancer follow-up, imposer un KPI/SLO ou choisir un analytics engine.

## 5. Propriétaire
Govern / Response Metrics owns response-outcome and residual-risk metric semantics. GOV-2 owns Result/Verification/Residual Risk source records; Command/Investigate own downstream objects; Shared owns Metrics/Reporting mechanisms.

## 6. Utilisateurs
Principal : Govern Control Reviewer. Secondaires : Response/Verification Reviewer, Decision Maker, Incident Commander, Investigator, Business/Service Owner as contextual contributor and Govern Product Lead.

## 7. Conditions d’entrée
Canonical Result records and relevant Verification/Residual Risk sources are available for a scoped snapshot; outcome definitions and follow-up dimensions are versioned; unavailable data is not coerced to zero.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Result/outcome | CAP-GOV-032 | canonical outcome observation | oui | exact Result version | metric unavailable |
| Verification/Residual Risk | CAP-GOV-029 | effectiveness context | when required | Result-related version | partial/inconclusive |
| rollback/recovery | CAP-GOV-030/031 | recovery context | when applicable | snapshot | NA/unknown explicit |
| Run reliability | CAP-GOV-042/043 source observations | technical context | non for Result count | same comparison window | no technical inference |
| downstream follow-up refs | CAP-GOV-033/Command/Investigate | consumer context | when present | correlated snapshot | none/unavailable distinguished |
| metric definition/snapshot | Shared Metrics | calculation contract | oui | versioned | no metric output |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Result | Govern | lifecycle/outcome/follow-up | metric read |
| Verification/Residual Risk Assessment | Govern | verified outcome/remaining risk | metric read |
| Response Run/Rollback | Govern | context only | metric read |
| Incident | Command | downstream follow-up projection | restricted aggregate/link |
| Case/Finding | Investigate | follow-up projection | restricted aggregate/link |
| Metric definition/snapshot | Shared | definition/freshness/privacy | consume |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Result Metric Observation | calculate/version/supersede | Govern local concept | Result classification remains source-owned |
| Residual Risk/Effectiveness Observation | calculate/version | Govern local concept | no universal effectiveness score |
| source Result/Incident/Case/Finding | no mutation | respective owner | metrics do not trigger changes |

## 11. Fonctionnalités
Aggregate Result outcomes; segment by action/Playbook/target/environment; preserve verification-failed/inconclusive; count residual-risk and unintended-effect categories where defined; count follow-up needs/refs; observe repeated-action candidates; compare periods; correlate technical/recovery context without causal assertion; expose missing/late feedback and source coverage.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| read/filter outcome metrics | reviewer | observation | 0 | permission | scoped metrics | non |
| calculate/reconcile outcomes | reviewer/system | observation | 1 | definition+snapshot | deterministic observation | non |
| annotate effectiveness hypothesis | reviewer | observation | 2 | rationale | attributed hypothesis | OPEN-013 |
| prepare comparison/report | reviewer | observation set | 2 | comparable definitions | comparison | OPEN-013 |
| reclassify Result/start follow-up automatically | none | source objects | 3/4 | prohibited | no action | owner workflows only |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| aggregate Result outcomes | oui | oui | oui | explain | outcome table |
| compare verified vs Result classifications | oui | explicit rules | oui | summary | consistency matrix |
| group residual-risk/unintended-effect observations | oui | defined categories | oui | candidate summary | grouped table |
| draft effectiveness hypothesis | oui | source aggregation | oui | hypothesis only | analyst note |
| declare business value/root cause | human product/business review | no | no | prohibited | explicit review |

## 14. États fonctionnels
`not-calculated`, `available`, `partial`, `verification-incomplete`, `feedback-incomplete`, `insufficient-data`, `privacy-limited`, `definition-changed`, `stale`, `disputed`, `superseded`.

## 15. États d’interface
Loading retains dimensions ; Empty distinguishes no Results from missing data ; Partial exposes missing verification/follow-up ; Error preserves last valid snapshot ; Offline read-only ; Permission denied masks target/business dimensions ; Stale exposes Result/feedback lateness.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Result Outcome Observation | metric | CAP-GOV-046/Reporting | canonical Result source/version retained |
| Residual Risk Observation | metric | CAP-GOV-046/047 | no zero-risk implication |
| Effectiveness hypothesis candidate | derived observation | CAP-GOV-046/047 | hypothesis, not causal conclusion |
| follow-up observation | metric/context | Command/Investigate review | no destination mutation |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-029..033 | snapshot selected | CAP-GOV-044 | Result/verification/risk/follow-up refs | source records |
| CAP-GOV-044 | trend/control review | CAP-GOV-046 | observations/definitions/limitations | Response Metrics |
| CAP-GOV-044 | improvement candidate | CAP-GOV-047 | sourced outcome/residual issue | Response Metrics |
| CAP-GOV-044 | report/export | Shared Reporting/Export | metric refs/snapshot | Response Metrics |

## 18. Dépendances
CAP-GOV-029..033/034/042/043/046/047, Command/Investigate follow-up projections, Shared Metrics/Reporting, Settings tenant dimensions, OPEN-008/013/015.

## 19. Source de vérité
Canonical Result and Verification/Residual Risk records remain source. Govern owns derived domain metric observations. Consumer feedback remains attributed context, not ground truth. Metrics never rewrite Result or Decision.

## 20. Provenance et audit
Record definition/version, Result/verification snapshots, outcome categories, residual-risk dimensions, follow-up refs, excluded/unknown items, freshness, privacy suppression, reviewer hypotheses and AI summary provenance.

## 21. Permissions fonctionnelles
`perm.govern.metrics.read`; restricted Result/business dimension read; cross-tenant comparisons separately authorized; Reporting/Export preparation. No source mutation permission is implied.

## 22. Limites et erreurs
Missing verification, late follow-up, inconsistent Result versions, sparse residual-risk categorization, mixed action populations and restricted business context limit effectiveness interpretation. Result success is not automatically business value; failure is not automatically Decision error.

## 23. Métriques
Outcome distribution; residual-risk concepts; unintended effects; verification-failed/inconclusive; follow-up Investigation/Command refs; repeated-action candidates; feedback coverage. No universal effectiveness score or target.

## 24. Classification de livraison
`defined` / `planned`; no effectiveness model, KPI/SLO, causal engine or implementation selected.

## 25. Critères d’acceptation
**Given** Result is `success` but residual risk remains, **When** outcome metrics are reviewed, **Then** success and residual risk are both visible and zero-risk/business-value is not inferred.

**Given** Result is `failed` after executor unavailability, **When** metrics are viewed, **Then** failure is not automatically labelled a Decision error.

**Given** consumer feedback is positive, **When** effectiveness is summarized, **Then** feedback remains an attributed input rather than ground truth.

**Given** no AI, **When** metrics are generated, **Then** deterministic aggregation, source matrices and human review provide full functionality.

## 26. Questions ouvertes
OPEN-008/013/015 remain open. Final effectiveness ontology, business-value model and thresholds remain future; no new OPEN.

## 27. Consommateurs documentaires
Response Metrics, CAP-GOV-046/047, Govern closure, Command/Investigate feedback loops, Shared Reporting/Metrics and future product/technical review.