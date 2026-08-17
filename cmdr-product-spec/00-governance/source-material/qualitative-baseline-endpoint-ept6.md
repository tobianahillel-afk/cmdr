---
id: qualitative-baseline-endpoint-ept6
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-12
source-of-truth: source-material
---
# Qualitative Baseline — Endpoint EPT-6

## Baseline
EPT-6 starts from remote `4d20e8e411aa8919f31e773a9c5ea96efd64791e`, after EPT-5 post-publication verification.

Preserved prior state:
- Command 27 PASS;
- Investigate 243 PASS;
- Govern 47 PASS;
- Studio 68 PASS;
- EPT-1 14 / 378 / 84 — PASS 190/190;
- EPT-2 16 / 432 / 96 — PASS 200/200;
- EPT-3 16 / 432 / 96 — PASS 210/210;
- EPT-4 18 / 486 / 108 — PASS 220/220;
- EPT-5 17 / 459 / 102 — PASS 230/230;
- Endpoint cumulative 81 / 2187 / 486;
- global 466 capabilities / 464 defined / 2 proposed / 466 planned / 12582 sections / 2796 mandatory tables;
- Requirements 122 = 99 conform / 20 partial / 3 absent / 0 contradictory;
- OPEN 18;
- Endpoint Screen IDs 0.

## EPT-6 target
Complete only the remaining Endpoint documentary families for local Agent update lifecycle, resilience, self-protection/security state and Endpoint provenance/closure.

## Ownership invariants
- Settings owns administrative Fleet/update policy and Secret administration.
- Endpoint owns local technical update/resilience/security facts and operations.
- Studio owns Studio asset deployment/reversion.
- Govern owns Response Rollback, response verification and Result.
- Shared owns generic Job/Retry/Recovery/Trace/Activity mechanisms.
- Security owns global permission/privacy/audit-integrity policy.

## Delivery invariant
`defined / planned` is documentary coverage, not implementation evidence. No EPT-6 text may infer platform support, product availability, production readiness or Phase 6 completion.

## Candidate set
The source-driven audit supports `CAP-EPT-082..099` as 18 distinct contracts. Final counts must be recalculated from the actual files and structural checks, not assumed solely from this projection.