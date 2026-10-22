---
id: CAP-GOV-042
title: Response Run Execution and Reliability Metrics
product: govern
module: response-metrics
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-002, REQ-PROD-005, REQ-PROD-008, REQ-PROD-015, REQ-PROD-016, REQ-PROD-020, REQ-OBJ-007, REQ-SEC-001]
open_decisions: [OPEN-008, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-GOV-042 — Response Run Execution and Reliability Metrics

## 1. Définition
Définir les métriques Govern de création, démarrage, completion, failure, cancellation/stop, partiality, runtime/executor/target availability, retry, step failure, timeout et unknown state des Response Runs, sans confondre fiabilité technique avec efficacité ou succès vérifié.

## 2. Problème utilisateur
Une exécution peut être techniquement « successful » tout en échouant à produire l’effet attendu, et des Runs partiels peuvent être masqués par un taux global. Sans dimensions per-target/per-step et source coverage, reliability metrics deviennent trompeuses.

## 3. Objectifs
- compter Runs created/started/completed/failed/cancelled/stopped/partial ;
- observer runtime/executor/target unavailable ;
- mesurer retries, retry outcomes, timeouts et step failures ;
- préserver per-target partiality et unknown technical state ;
- segmenter par Playbook/action/target type/tenant/environment/période ;
- garder technical reliability distincte de verification/effectiveness.

## 4. Non-objectifs
Ne pas déclarer effectiveness, modifier Run/Result, définir un runtime SLO universel, scorer un provider, choisir retry engine, imposer un executor ou créer un monitoring backend.

## 5. Propriétaire
Govern / Response Metrics owns Response Run governance metric semantics. GOV-2 owns source Run/Step/error records; Studio/Endpoint/provider own raw technical metrics; Shared owns generic Metrics/Reporting infrastructure.

## 6. Utilisateurs
Principal : Govern Control Reviewer. Secondaires : Response Operator, Verification Reviewer, Studio/Endpoint Operator, Platform Operator, Incident Commander and Govern Product Lead.

## 7. Conditions d’entrée
Run/Step/runtime/error records are correlated enough for a scoped snapshot; unknown/unavailable states are represented explicitly; metric definitions distinguish created, requested, confirmed and terminal events.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Response Run lifecycle | CAP-GOV-022/023 | run observations | oui | snapshot | metric partial |
| Step/action outcomes | CAP-GOV-024 | execution observations | when steps exist | snapshot | step coverage unknown |
| runtime reconciliation | CAP-GOV-025/026 | executor status | when execution attempted | correlated source | technical state unknown |
| errors/retries/partiality | CAP-GOV-027 | reliability observations | when present | snapshot | zero vs unavailable distinguished |
| Playbook/target dimensions | CAP-GOV-017..020 | segmentation | when permitted | pinned versions | dimension unknown |
| metric definition/snapshot | Shared Metrics | calculation contract | oui | versioned | no metric publication |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Response Run / Step | Govern | lifecycle/targets/states | metric read |
| Execution Error / Retry | Govern | failure/attempt context | metric read |
| Workflow/Automation Run/Tool Call | Studio | technical correlation only | restricted aggregate/link |
| technical runtime records | Endpoint/provider | availability/status projection | restricted aggregate |
| Metric definition/snapshot | Shared | version/freshness/privacy | consume |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Response Run Metric Observation | calculate/version/supersede | Govern local concept | Run states remain distinct |
| reliability/partiality observation | calculate/version | Govern local concept | technical reliability ≠ effectiveness |
| source Run/technical records | no mutation | source owner | aggregate only |

## 11. Fonctionnalités
Aggregate Run lifecycle counts; calculate start/terminal distributions; segment failure causes; count unavailable runtime/executor/target; aggregate retries and outcomes; count timeouts/step failures/unknown states; preserve per-target partiality; compare Playbook versions/action classes/environments; expose source coverage and missing correlations; feed trend and verification/outcome analysis without collapsing semantics.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect/filter Run metrics | reviewer | observation | 0 | metric read | scoped values | non |
| calculate reliability/partiality | reviewer/system | observation | 1 | snapshot+definition | deterministic observation | non |
| annotate suspected degradation | reviewer | observation | 2 | rationale | attributed note | OPEN-013 |
| prepare comparison | reviewer | observation set | 2 | comparable versions | comparison | OPEN-013 |
| retry/stop/change Run from metric | none | Run | 3/4 | prohibited | no action | GOV-2 workflow only |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| aggregate Run states | oui | oui | oui | explain | pivot/table |
| aggregate failure/retry reasons | oui | mapped classes | oui | summary | reason table |
| detect reliability trend candidate | oui | comparisons | oui | candidate explanation | period comparison |
| summarize partiality | oui | per-target aggregation | oui | sourced summary | target matrix |
| declare effectiveness/root cause | human review | no | no | hypothesis only | cross-metric review |

## 14. États fonctionnels
`not-calculated`, `available`, `partial`, `source-incomplete`, `runtime-unknown`, `privacy-limited`, `definition-changed`, `stale`, `disputed`, `superseded`.

## 15. États d’interface
Loading retains scope ; Empty distinguishes no Runs from missing data ; Partial shows source/correlation coverage ; Error preserves prior snapshot ; Offline read-only ; Permission denied suppresses sensitive target/executor dimensions ; Stale shows latest source time.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Run reliability observation | metric | CAP-GOV-046/Reporting | state definitions/snapshot explicit |
| failure/retry observation | metric | control review/CAP-GOV-047 | reason/source coverage visible |
| partiality observation | metric | CAP-GOV-043/044 | per-target semantics retained |
| reliability trend candidate | derived signal | CAP-GOV-046 | candidate not causal conclusion |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-022..027 | snapshot selected | CAP-GOV-042 | Run/step/runtime/error observations | source Run |
| CAP-GOV-042 | verification comparison | CAP-GOV-043/044 | technical reliability refs | Response Metrics |
| CAP-GOV-042 | trend review | CAP-GOV-046 | metric definitions/snapshots | Response Metrics |
| CAP-GOV-042 | improvement candidate | CAP-GOV-047 | sourced reliability issue | Response Metrics |

## 18. Dépendances
CAP-GOV-017..027/034/043/044/046/047, Studio/Endpoint technical sources, Shared Metrics/Reporting, Settings environment/runtime metadata, OPEN-008/013/015.

## 19. Source de vérité
GOV-2 Run/Step/error records are source for Govern execution governance; technical runtimes remain source for raw status. Govern owns derived Run metric semantics. Reliability metric never becomes canonical Result or verification truth.

## 20. Provenance et audit
Record metric definition/version, Run/step source snapshot, included/excluded states, failure mappings, retry/target denominators, technical source coverage, dimensions, freshness, privacy suppression, annotations and AI summaries.

## 21. Permissions fonctionnelles
`perm.govern.metrics.read`; Run/technical aggregate read; sensitive target/executor dimension read; cross-tenant comparison separately authorized; Reporting/Export preparation. No Run-execution permission implied.

## 22. Limites et erreurs
Unknown executor state, late technical records, partial target visibility, inconsistent states, retention and changed Playbook versions limit comparability. Technical success never automatically means verified outcome success.

## 23. Métriques
Runs created/started/completed/failed/cancelled/stopped/partial; runtime/executor/target unavailable; retry count/outcome; per-target partiality; step failure; timeout; technical unknown state. No universal reliability target.

## 24. Classification de livraison
`defined` / `planned`; no runtime monitoring, SLO, metrics engine or provider implementation selected.

## 25. Critères d’acceptation
**Given** a Run succeeds technically on two targets and fails on one, **When** reliability metrics calculate, **Then** partiality remains explicit and overall success is not fabricated.

**Given** executor status is unknown, **When** metrics calculate, **Then** the Run contributes to unknown/source-incomplete dimensions rather than successful/failed by assumption.

**Given** retry count rises, **When** control review occurs, **Then** the metric does not automatically label effectiveness or root cause.

**Given** no AI, **When** metrics are generated, **Then** deterministic state/reason aggregation and human review are complete.

## 26. Questions ouvertes
OPEN-008/013/015 remain open. Runtime support, retry implementation and SLOs remain future; no new OPEN.

## 27. Consommateurs documentaires
Response Metrics, CAP-GOV-043/044/046/047, Govern closure, Shared Reporting/Metrics and future operational/technical review.