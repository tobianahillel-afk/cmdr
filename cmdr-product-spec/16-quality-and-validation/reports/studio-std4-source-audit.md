---
id: studio-std4-source-audit
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-10
source-of-truth: quality-report
---
# Studio STD-4 — Source Audit

## Baseline
- repository: `tobianahillel-afk/cmdr`;
- branch: `docs/cmdr-product-spec-foundation`;
- PR #2: open / Draft / unmerged / base `main`;
- baseline: `9babd679f52f3f28458a5f8f4d9c76698ebf875a`;
- baseline title: `docs: record Studio STD-3 post-publication verification`;
- main: `bc1ec59e5e79ccbaba291e2984b0eb7d5e54129c`;
- branch/main README exact `# cmdr`, same blob `901c74cda52e28b5ff7fc425ddf28ef89f3ad875`;
- auto-merge disabled;
- no `CAP-STD-052+`, no `CAP-EPT-*`, no temporary `std4` branch found before allocation.

## Sources read before ID allocation
Governance: Capability/Object/Permission/Ownership registers, unresolved decisions, Requirements Matrix, qualitative baseline, capability template, source-of-truth/roadmap and Phase-5 preflight.

Studio: Assurance README/evaluations/policy checks/regression/simulations, Evaluations README/model/datasets/scoring, Simulations README/model/scenario library, Versions & Deployment README/versioning/deployment/environments/rollback, canonical Evaluation/Simulation/Version/Deployment objects, Evaluation/Simulation implementation contracts and all STD-1/2/3 capability/evidence boundaries needed for non-regression.

Screens: STD-EVL-001, STD-SIM-001, STD-DEP-001 and the existing Studio screen inventory including STD-ASR-001. Historical verbs such as `Approve result`, `Approve promotion` and `Rollback` are treated as screen/source language, not as authority or final STD-4 semantics.

External boundaries: Govern Approval/Response Run/Result/Rollback; Settings tenants/environments/providers/secrets/health and Endpoint Fleet upgrade administration; Shared Background Jobs/Reporting/Export/Notifications; Endpoint update/rollback. These sources confirm separate ownership.

Quality: STD-1 conformance, STD-2 post-publication verification, STD-3 post-publication verification and the Phase-5 preflight.

## Findings
1. Evaluation, Simulation, Version and Deployment are canonical Studio-owned objects.
2. Existing Assurance/Evaluation/Simulation/Deployment documents are foundation-level and do not provide the 27-section capability contracts required for STD-4.
3. `perm.studio.*` screen references coexist with `perm.cmdr-studio.*` object/catalog permissions; anomaly remains unresolved.
4. Studio historical `rollback` refers to Studio asset/version lifecycle; Govern owns response rollback governance. STD-4 must use explicit **Deployment Reversion** semantics for the Studio lifecycle boundary.
5. Endpoint owns its own technical update/rollback; Settings owns fleet upgrade administration. Studio deployment must not absorb either.
6. Shared Reporting publication is publication of Reports, not Studio asset publication.
7. The Capability Register ends at `CAP-STD-051`; probes for `CAP-STD-052` and `CAP-STD-068` returned no results. No `CAP-EPT-*` capability exists.
8. The 17-capability allocation `CAP-STD-052..068` is therefore available and justified.

## Allocation decision
Allocate exactly `CAP-STD-052..068` as the STD-4 set from the execution specification. No ID is recycled or skipped. Endpoint remains outside this execution.
