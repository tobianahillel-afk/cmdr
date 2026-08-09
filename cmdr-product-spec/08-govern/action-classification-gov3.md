---
id: govern-action-classification-gov3
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical-addendum
open_decisions: [OPEN-013, OPEN-019]
---
# Govern Action Classification — GOV-3 Addendum

This file extends the existing Govern action-classification contract without changing GOV-1/GOV-2 classes.

| Class | GOV-3 actions | Boundary |
|---:|---|---|
| 0 | view, search, filter, inspect, compare, trend, read audit/metrics | observation only |
| 1 | deterministic aggregation, reconstruction, completeness assessment, metric calculation, comparison, bounded audit-package validation | no production effect |
| 2 | annotate, assign review, create Audit Evidence Package, create/review Trend or Control Health Assessment, create Continuous Improvement Package, prepare export/handoff, supersede assessments | reversible documentary/governance records; OPEN-013 default remains open |
| 3 | none required as a new GOV-3 production action | any accepted improvement is executed only in its destination-owner lifecycle |
| 4 | no audit/provenance destruction | irreversible trace deletion or uncontrolled external disclosure is outside GOV-3 and prohibited by default |

## Invariants
- Audit review never mutates Decision, Approval, Response Run, Result, Evidence or Finding.
- Audit export preparation does not authorize external sharing.
- Continuous Improvement Package is a proposal, not a change.
- Cross-tenant reads/exports require explicit permission.
- GOV-3 introduces no new target-side production action.