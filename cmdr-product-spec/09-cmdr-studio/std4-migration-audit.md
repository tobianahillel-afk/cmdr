---
id: studio-std4-migration-audit
domain: 09-cmdr-studio
status: draft
owner: CMDR Studio Product Lead
updated: 2026-08-10
source-of-truth: migration-audit
---
# STD-4 Migration Audit

Sources for testing, evaluation, simulation, assurance, regression, policy checks, versioning, promotion, deployment, rollback, deprecation and lifecycle were read before canonical capability allocation.

Findings:
- existing Assurance/Evaluations/Simulations/Versions & Deployment modules remain active and are **not deprecated**;
- their requirements are consumed by `CAP-STD-052..068`;
- historical UI verbs such as “Approve result”, “Approve promotion” and “Rollback” remain source evidence but do not establish authority;
- normative capability semantics use Publishing Review / Quality Gate and Studio Deployment Reversion, preserving Govern Approval/Decision/Response Rollback ownership;
- no competing functional source requires deprecation;
- no STD-1/2/3 module is deprecated or rewritten;
- no physical test format, sandbox engine, deployment engine, package registry or Endpoint implementation is selected.