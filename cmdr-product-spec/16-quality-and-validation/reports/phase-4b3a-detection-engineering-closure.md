---
id: phase-4b3a-detection-engineering-closure
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-06
source-of-truth: quality-report
requirements:
  - REQ-PROD-006
  - REQ-PROD-012
  - REQ-PROD-014
  - REQ-PROD-019
  - REQ-PROD-020
  - REQ-INV-006
open_decisions:
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
  - OPEN-017
---
# Phase 4B.3A — Detection Engineering closure

## Verdict
**PASS**, subject to final remote publication verification. CAP-INV-401..417 and CAP-INV-418..435 form one coherent Detection Engineering lifecycle. Threat Intelligence is not included.

## Sub-phases
| Sub-phase | Capability range | Capabilities | Sections | Tables | Gates | Status |
|---|---|---:|---:|---:|---:|---|
| 4B.3A.1 Foundations, Authoring and Validation | CAP-INV-401..417 | 17 | 459 | 102 | 160 | PASS |
| 4B.3A.2 Review, Promotion, Runtime Performance and Lifecycle | CAP-INV-418..435 | 18 | 486 | 108 | 170 | PASS after publication verification |
| **Detection Engineering total** | **CAP-INV-401..435** | **35** | **945** | **210** | — | **PASS** |

## Reports
- `phase-4b3a1-detection-authoring-validation-capability-conformance.md`;
- `phase-4b3a2-detection-lifecycle-performance-capability-conformance.md`;
- this closure report.

## End-to-end functional coverage
1. Need intake, Project and Detection Hypothesis.
2. Data-source readiness, schema/field review and functional authoring.
3. Conditions, correlation, sequences, thresholds, windows, enrichment and metadata.
4. Structural/semantic validation, scenarios, datasets, Expected Outcomes and historical replay.
5. Candidate TP/FP/FN review, coverage, gaps and Review Package.
6. Immutable Release Candidate, formal human review and readiness.
7. Target/promotion planning, Govern handoff, shadow and canary.
8. Govern-owned execution observation and runtime version reconciliation.
9. Health, Command feedback, runtime quality and production match review.
10. Tuning, bounded suppression/exception proposals, drift and performance.
11. Rollback, recovery, deactivation, retirement, replacement and lifecycle provenance.
12. Continuous Improvement Package returning to a new authoring cycle.

## Ownership closure
Investigate owns business context, content, assessments and proposals. Command owns runtime Detection, Signal, Alert and Incident. Settings owns environments, targets, configured runtimes and administration. Endpoint owns local capability/execution projections. Govern owns authority and production execution objects. Studio owns Tools and automation. Shared owns generic mechanisms. Concurrent owners: **0**.

## Transition closure
- CAP-INV-417 → CAP-INV-418 preserves immutable version, evidence, gaps and return origin.
- CAP-INV-420 → Govern preserves action class, risks, target scope, rollout and rollback.
- Govern Result → CAP-INV-424/425 preserves per-target partiality and observed state.
- Command feedback → CAP-INV-427/428 preserves source ownership and uncertainty.
- CAP-INV-429/431/432/433/434 → new Draft or request preserves evidence and never mutates active content silently.
- CAP-INV-435 returns to CAP-INV-401/402/403/406 and preserves the complete lineage.

## Decision closure
OPEN-017 is created open because runtime, target language and portability are required future product decisions. No option is selected. OPEN-005 remains forensic-only. Open decisions total: **16**. No decision is closed; OPEN-009 remains the only historically resolved item.

## Requirements and maturity
Requirement IDs remain 122: 99 conform, 20 partial, 3 absent, 0 contradictory. REQ-INV-006 gains complete functional Detection Engineering evidence but remains globally partial until Threat Intelligence and implementation are complete.

## Screens and migration
Required surfaces read: 16. Screen specs modified: 0. Detailed rewrites: 0. New Screen IDs: 0. Competing Investigate lifecycle sources deprecated: 0. Runtime and administrative contracts remain active under their owners.

## Quality
- 35/35 Detection Engineering capability files;
- 945/945 sections;
- 210/210 mandatory tables;
- empty/prose-only/generic mandatory tables: 0;
- duplicate/recycled IDs: 0;
- active duplicates/concurrent owners: 0;
- placeholders/empty targeted files/broken targeted links: 0;
- complete schemas/APIs/protocols/engines/languages/commands/code: 0;
- actual deployments/activations/active exceptions/Signal deletions: 0;
- CAP-INV-5xx/Threat Intelligence content: 0.

## Final status
- Phase 4B.3A.1: PASS.
- Phase 4B.3A.2: PASS after publication verification.
- Phase 4B.3A: PASS.
- Phase 4B.3: PARTIAL.
- Phase 4B: PARTIAL.
- Phase 4: PARTIAL.
- Global maturity: PARTIAL.
- Next phase: 4B.3B.1 Threat Intelligence Foundations and Knowledge Management — NOT STARTED.
