---
id: CAP-GOV-043
title: Verification, Rollback and Recovery Metrics
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
# CAP-GOV-043 — Verification, Rollback and Recovery Metrics

## 1. Définition
Définir les métriques Govern de verification outcome, rollback review/planning/execution et recovery afin d’observer verified-success/partial/failure/inconclusive, rollback/recovery states et residual impact sans assimiler rollback frequency à failure rate ni rollback completion à restauration totale.

## 2. Problème utilisateur
Une hausse de rollback peut indiquer un problème, une bonne capacité de récupération ou un changement de mix d’actions. Un rollback techniquement terminé peut laisser un residual impact. Sans séparation verification/rollback/recovery, les métriques masquent ces nuances.

## 3. Objectifs
- compter verification executed et outcomes ;
- observer insufficient-data/inconclusive/adverse effects ;
- compter rollback-review-required/planned/started/completed/partial/failed ;
- observer recovery-required/completed et residual impact ;
- comparer avec Run reliability sans les fusionner ;
- segmenter par action/Playbook/target/période avec limitations.

## 4. Non-objectifs
Ne pas conclure automatiquement root cause, fixer rollback target rate, modifier Verification/Result, déclencher rollback/recovery, choisir une verification engine ou produire un score universel.

## 5. Propriétaire
Govern / Response Metrics owns these metric semantics. GOV-2 owns verification/rollback/recovery source records; technical owners retain raw primitives; Shared owns generic Metrics/Reporting mechanisms.

## 6. Utilisateurs
Principal : Govern Control Reviewer. Secondaires : Verification Reviewer, Rollback Reviewer, Response Operator, Incident Commander, Security/Compliance Reviewer, Govern Product Lead.

## 7. Conditions d’entrée
Verification/Rollback/Recovery records are available for a scoped snapshot; applicability is known so `not-applicable`, unavailable and true zero are distinguishable.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Verification Plan/Assessment | CAP-GOV-028/029 | verification observations | oui when required | snapshot | verification coverage partial |
| Rollback eligibility/plan | CAP-GOV-030 | rollback observations | when applicable | snapshot | NA/unknown explicit |
| rollback/recovery execution | CAP-GOV-031 | outcome observations | when started | snapshot | state unknown/partial |
| Run/error context | CAP-GOV-022/027 | denominator/context | oui | snapshot | context incomplete |
| Result/residual risk refs | CAP-GOV-032 | outcome context | when available | snapshot | no value fabricated |
| metric definition/snapshot | Shared Metrics | calculation contract | oui | versioned | no metric output |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Verification Plan/Assessment | Govern | criteria/outcome/gaps | metric read |
| Rollback Plan/Response Rollback | Govern | state/outcome | metric read |
| Recovery Context | Govern/technical source | completion/residual state | metric read/link |
| Response Run / Result | Govern | denominator/context | metric read |
| technical rollback refs | Endpoint/Studio/provider | raw outcome refs | restricted aggregate |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Verification Metric Observation | calculate/version/supersede | Govern local concept | verified outcome distinct from runtime success |
| Rollback/Recovery Metric Observation | calculate/version/supersede | Govern local concept | completion ≠ full recovery claim |
| source verification/rollback records | no mutation | GOV-2 owners | aggregate only |

## 11. Fonctionnalités
Aggregate verification coverage/outcomes; count rollback-review-required/planned/started/terminal states; preserve partial rollback/failure; observe recovery-required/completed; track residual impact presence; calculate verification lag conceptually; segment by action/Playbook/target/environment; compare periods; expose applicability/source coverage; correlate with Run reliability and Result outcomes without causal inference.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| read/filter metrics | reviewer | observation | 0 | permission | scoped values | non |
| calculate/reconcile | reviewer/system | observation | 1 | definitions+snapshot | deterministic observation | non |
| annotate interpretation | reviewer | observation | 2 | rationale | attributed note | OPEN-013 |
| prepare comparison | reviewer | observation set | 2 | comparable scope | comparison | OPEN-013 |
| trigger rollback/recovery | none | source action | 3/4 | prohibited from metrics | no action | GOV-2 only |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| aggregate verification outcomes | oui | oui | oui | explain | outcome table |
| aggregate rollback/recovery states | oui | oui | oui | summary | state matrix |
| compare trends | oui | deterministic periods | oui | candidate narrative | comparison table |
| flag residual-impact pattern | oui | grouped observation | oui | hypothesis | grouped table |
| recommend/trigger rollback | human GOV-2 workflow | no automatic from metric | no | proposal only, no action | explicit GOV-2 review |

## 14. États fonctionnels
`not-calculated`, `available`, `partial`, `applicability-unknown`, `source-incomplete`, `privacy-limited`, `definition-changed`, `stale`, `disputed`, `superseded`.

## 15. États d’interface
Loading retains dimensions ; Empty distinguishes no applicable rollback from missing data ; Partial lists absent verification/recovery sources ; Error preserves prior snapshot ; Offline read-only ; Permission denied suppresses target/runtime detail ; Stale exposes snapshot time.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Verification Metric Observation | metric | CAP-GOV-044/046/Reporting | runtime success remains separate |
| Rollback/Recovery observation | metric | CAP-GOV-044/046 | partial/failure/recovery distinguished |
| residual-impact observation | metric | CAP-GOV-044/047 | no causal conclusion |
| control-review candidate | derived signal | CAP-GOV-046/047 | sourced candidate only |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-028..032 | snapshot selected | CAP-GOV-043 | verification/rollback/recovery/result refs | source records |
| CAP-GOV-043 | effectiveness review | CAP-GOV-044 | observations/residual context | Response Metrics |
| CAP-GOV-043 | trend/control review | CAP-GOV-046 | definitions/snapshots | Response Metrics |
| CAP-GOV-043 | improvement candidate | CAP-GOV-047 | sourced issue | Response Metrics |

## 18. Dépendances
CAP-GOV-022/027..032/034/042/044/046/047, technical rollback sources, Shared Metrics/Reporting, Settings environment metadata, OPEN-008/013/015.

## 19. Source de vérité
GOV-2 verification/rollback/recovery records remain source for Govern governance outcomes; technical owners retain raw primitives. Govern owns derived metric semantics. A metric never changes Result or triggers rollback.

## 20. Provenance et audit
Record metric definition/version, applicability/denominator, source snapshot, verification/rollback/recovery states, residual-impact dimension, excluded/unknown items, technical source coverage, reviewer notes and AI summaries.

## 21. Permissions fonctionnelles
`perm.govern.metrics.read`; restricted verification/target/runtime aggregate read; cross-tenant comparison separately authorized; Reporting/Export preparation. No rollback/run permission implied.

## 22. Limites et erreurs
Missing verification, unknown rollback applicability, partial technical data, late recovery observations, retention and mixed Playbook versions can limit comparisons. Rollback frequency is not automatically failure rate; completion is not proof of exact restoration.

## 23. Métriques
Verification executed; verified-success/partial/failure; inconclusive/insufficient-data; rollback-review-required; rollback planned/started/completed/partial/failed; recovery required/completed; residual impact. No universal targets.

## 24. Classification de livraison
`defined` / `planned`; no verification/rollback metrics engine, target or SLO selected.

## 25. Critères d’acceptation
**Given** runtime reports success but Verification is failed, **When** metrics are calculated, **Then** technical success and verification-failed remain separate observations.

**Given** rollback rate increases, **When** trends are reviewed, **Then** the system does not automatically label overall failure rate as increased.

**Given** rollback completes but residual impact remains, **When** recovery metrics are viewed, **Then** completion is not presented as full restoration.

**Given** no AI, **When** metrics are generated, **Then** deterministic aggregation and human review remain complete.

## 26. Questions ouvertes
OPEN-008/013/015 remain open. Final technical verification/rollback implementations and targets remain future; no new OPEN.

## 27. Consommateurs documentaires
Response Metrics, CAP-GOV-044/046/047, Govern closure, Shared Reporting/Metrics, future screen/technical review and quality reports.