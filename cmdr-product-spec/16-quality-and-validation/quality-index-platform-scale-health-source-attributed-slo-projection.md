---
id: quality-index-platform-scale-health-source-attributed-slo-projection
domain: 16-quality-and-validation
status: draft
owner: Product Architecture
updated: 2026-08-16
source-of-truth: canonical
---
# Quality Index — Platform Health and Source-Attributed SLO Projection

## Target
Single functional capability: `CAP-SET-014 — Platform Health and Source-Attributed SLO Projection`.

## Exact gate arithmetic
| Class | Gates | Current result |
|---|---:|---|
| Source/local | 134 | 134 PASS |
| Remote/post-publication | 8 | 8 PENDING-REMOTE |
| **Total** | **142** | **134 PASS / 8 PENDING-REMOTE / 0 FAIL** |

Source/local families are fixed at A8 + B12 + C14 + D16 + E15 + F11 + G12 + H12 + I10 + J12 + K12 = **134**. Section 2B is included inside family A and does not add a gate.

Remote gates R1–R8 require actual publication, remote canonical reread, PR/main/README invariants, CI/status/check/workflow applicability evidence, documentary closure, functional blob invariance and final STOP verification.

## Non-regression invariants
- ADR-0009 unchanged and authoritative;
- `CAP-SET-014` only; `CAP-SET-015+` unallocated/unreserved;
- no new canonical object, Permission ID or Screen ID;
- `perm.settings.health.read` read-only;
- writes none;
- Requirements remain 122 = 99/20/3/0;
- OPEN remains 17;
- Screens remain 56;
- roadmap preservation REMOVED 0 / WEAKENED 0 / UNKNOWN 0;
- Localization, Advanced Integrations and Compliance not started.

Final PASS may be recorded only after all eight remote/post-publication gates are proven. If no CI mechanism applies, classification must be **N/A WITH EVIDENCE**, not “CI PASS”.
