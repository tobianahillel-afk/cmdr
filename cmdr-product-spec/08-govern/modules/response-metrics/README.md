---
id: govern-response-metrics
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-PROD-002, REQ-PROD-004, REQ-PROD-005, REQ-PROD-008, REQ-PROD-015, REQ-PROD-019, REQ-PROD-020, REQ-OBJ-007, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-008, OPEN-010, OPEN-013, OPEN-015, OPEN-019]
---
# Response Metrics — GOV-3

## Mission

Define Govern-specific metric meanings over the GOV-1/GOV-2 lifecycle, preserve definition/snapshot/freshness/privacy and produce sourced comparisons/control-health/improvement candidates without duplicating the Shared Metrics or Reporting engines and without turning metrics into policy, SLO or automatic decisions.

## Owned GOV-3 capabilities

- `CAP-GOV-039` — Policy, Exception and Emergency Governance Metrics.
- `CAP-GOV-040` — Approval, Authority and Separation-of-Duties Metrics.
- `CAP-GOV-041` — Decision Flow, Disposition and Timeliness Metrics.
- `CAP-GOV-042` — Response Run Execution and Reliability Metrics.
- `CAP-GOV-043` — Verification, Rollback and Recovery Metrics.
- `CAP-GOV-044` — Response Outcome, Residual Risk and Effectiveness Metrics.
- `CAP-GOV-045` — Govern Queue, Ageing and Lifecycle Flow Metrics.
- `CAP-GOV-046` — Govern Trend, Comparison and Control Health Assessment.
- `CAP-GOV-047` — Govern Continuous Improvement, Closure and Provenance.

## Ownership boundary

Govern owns metric semantics, dimensions and interpretation for Govern. Shared Capabilities retains the generic Metrics Engine, Reporting Engine, Export, Jobs and dashboard primitives. Command, Investigate, Studio and Endpoint retain their own product/technical metrics. Platform Settings retains tenant/environment/retention/storage configuration.

## Mandatory distinctions

- metric ≠ objective, Policy, SLO or KPI automatically;
- metric target ≠ universal truth;
- count ≠ quality;
- throughput ≠ effectiveness;
- faster Decision ≠ better Decision;
- Approval latency ≠ Approval quality;
- low/high exception rate ≠ governance health automatically;
- Policy block count ≠ prevented incident count;
- Run technical success ≠ verified outcome success;
- verified success ≠ zero residual risk;
- rollback rate ≠ failure rate automatically;
- Result success ≠ business value automatically;
- Result failure ≠ Decision error automatically;
- feedback ≠ ground truth;
- trend ≠ causal explanation;
- anomaly ≠ control failure;
- dashboard ≠ source of truth.

## Definitions and snapshots

Every observation identifies its metric definition/version, source snapshot/time window, dimensions, denominators/exclusions, source coverage, freshness and privacy suppression where relevant. Zero, not-applicable, unavailable and insufficient-data remain distinct.

## Privacy and cross-tenant use

Requester/approver identity, tenant, target, exceptions, emergency paths and sensitive operational dimensions are permission-aware. Cross-tenant comparisons require explicit authorization and may require suppression/aggregation. No dashboard or metric grants access to source data that the user cannot otherwise read.

## AI/no-AI

AI may explain sourced observations, summarize trends or draft hypotheses/improvement packages. It cannot mutate a Decision/Result, change thresholds, declare causality/control failure as fact, publish externally or apply an improvement. Deterministic aggregation, tables, comparisons and human review form the complete non-AI path.

## Screen

`GOV-MET-001` remains the existing Response Metrics surface. GOV-3 creates no new Screen ID and defines no final columns, filters, buttons, shortcuts, animations or wireframes. Shared Reporting/Inspector primitives remain source-owned.

## Closure boundary

`CAP-GOV-047` may assemble evidence for documentary Govern closure after all mandatory gates and post-publication checks complete. Govern PASS never means product implementation complete and does not start Delivery Roadmap Phase 5.