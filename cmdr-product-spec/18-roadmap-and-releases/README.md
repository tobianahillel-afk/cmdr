---
id: roadmap-readme
domain: 18-roadmap-and-releases
status: draft
owner: Product Operations Lead
updated: 2026-08-09
source-of-truth: canonical
---
# Roadmap and Releases

## Phase-numbering namespaces
CMDR retains two independent namespaces. `Capability Specification Phase 4A — Command`, `Capability Specification Phase 4B — Investigate` and `Delivery Roadmap Phase 4 — Govern` have no numeric parent/child relationship. `Phase 4C/4D/4E Govern` do not exist.

## Capability specification status
- Command: **PASS**, 27 / 729 / 162.
- Investigate: **PASS**, 243 / 6561 / 1458.
- Govern: **PARTIAL pending GOV-3 remote closure**.
  - GOV-1: historical PASS 180/180, 16 / 432 / 96.
  - GOV-2: historical PASS 190/190, 17 / 459 / 102.
  - GOV-3: CAP-GOV-034..047, 14 / 378 / 84; prepublication 192 PASS / 8 PENDING / 0 FAIL.
- Global Capability Specification maturity: **PARTIAL**.

## Delivery Roadmap phases
1. `phase-1-foundation.md` — Delivery Roadmap Phase 1 — Foundation.
2. `phase-2-command.md` — Delivery Roadmap Phase 2 — Command.
3. `phase-3-investigate.md` — Delivery Roadmap Phase 3 — Investigate.
4. `phase-4-govern.md` — **Delivery Roadmap Phase 4 — Govern**, id `roadmap-phase-4-govern`, currently **PARTIAL pending GOV-3 post-publication verification**.
5. `phase-5-studio-and-endpoint.md` — future candidate; not started by GOV-3.
6. `phase-6-platform-scale.md` — future.

## Govern execution lots
- GOV-1 `CAP-GOV-001..016`: Action Request → Policy/Authority/Approval → Decision → Execution Handoff.
- GOV-2 `CAP-GOV-017..033`: Playbook/Plan/Readiness → Response Run → Verification → Rollback/Recovery → Result.
- GOV-3 `CAP-GOV-034..047`: Audit Trail semantics/reconstruction/review → Govern metrics/flow/trends/control health → Continuous Improvement/closure.

The three lots together cover all nine canonical Govern modules. No Phase 5 capability is created.

## Current totals
- global capabilities: **317** — 27 Command / 243 Investigate / 47 Govern;
- defined/proposed/planned: **315 / 2 / 317**;
- Govern: **1269 sections / 282 tables**;
- Command + Investigate + Govern: **8559 sections / 1902 tables**;
- Requirements: **122 = 99/20/3/0**;
- OPEN: **18**.

## Closure evidence
- GOV-3 conformance: `../16-quality-and-validation/reports/govern-gov3-audit-metrics-closure-capability-conformance.md`.
- Govern full closure: `../16-quality-and-validation/reports/govern-capability-specification-closure.md`.
- Delivery Roadmap Phase 4 closure: `../16-quality-and-validation/reports/delivery-roadmap-phase-4-govern-closure.md`.

The closure reports identify Govern as a documentary PASS candidate but do not promote status before fifth-commit publication and remote 200-gate verification.

## Implementation boundary / stop line
Capability documentation does not prove software delivery. GOV-3 introduces no audit/metrics engine, API/protocol, warehouse/storage schema, final RBAC/retention policy, detailed screen rewrite or product code. **Do not start Delivery Roadmap Phase 5 in this execution.**