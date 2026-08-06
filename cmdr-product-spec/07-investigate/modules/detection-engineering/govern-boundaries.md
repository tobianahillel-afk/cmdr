---
id: investigate-detection-engineering-govern-boundaries
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements:
  - REQ-INV-006
  - REQ-SEC-001
  - REQ-SEC-002
open_decisions:
  - OPEN-013
  - OPEN-015
  - OPEN-017
---
# Govern boundaries — Detection Engineering lifecycle

Investigate prepares Detection Change Request Drafts, Promotion Plans, Shadow/Canary Plans, Tuning/Suppression/Exception Proposals, Rollback Plans and Retirement Proposals. These are evidence-bearing proposals, not authority.

Govern owns Action Request, Decision, Approval, Response Run and Result. Class-3 promotion, activation, deactivation, active suppression/exception, rollback, retirement and runtime modification require the applicable Govern path. Class 4 is denied by default or strictly governed; provenance is never destroyed.

| Investigate concept | Govern destination | Boundary |
|---|---|---|
| Detection Change Request Draft | Action Request | draft ≠ accepted request |
| Detection Review recommendation | Decision input | review ≠ Decision |
| Approval requirement | Approval | recommendation ≠ Approval |
| Promotion/Canary/Rollback/Retirement plan | Response Run context | plan ≠ execution |
| Runtime observation | Result relation | observation ≠ verified Result |
| Emergency recommendation | emergency policy path | recommendation ≠ emergency authority |

Requester, reviewer, approver and operator separation remains policy-owned. No auto-approval, hidden delegation, implicit emergency authority or Govern bypass.
