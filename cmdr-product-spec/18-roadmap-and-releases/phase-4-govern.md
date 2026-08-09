---
id: roadmap-phase-4-govern
domain: 18-roadmap-and-releases
status: draft
owner: Product Operations Lead
updated: 2026-08-09
source-of-truth: canonical
---
# Phase 4 Govern

## Canonical phase identity

- namespace: **Delivery Roadmap**;
- canonical id: `roadmap-phase-4-govern`;
- canonical title: `Phase 4 Govern`;
- qualified title: **Delivery Roadmap Phase 4 — Govern**;
- previous delivery phase: `Delivery Roadmap Phase 3 — Investigate`;
- next delivery phase: `Delivery Roadmap Phase 5 — Studio and Endpoint`;
- Govern capability specification status: **PARTIAL**.

`Delivery Roadmap Phase 4 — Govern` and `Capability Specification Phase 4` belong to distinct namespaces and have no numeric parent-child relationship. `Phase 4 Govern` MUST NOT be renamed to `Phase 4C Govern`. `Capability Specification Phase 4C` and `Phase 4C Govern` **do not exist**. See [`phase-numbering-and-namespace-convention.md`](phase-numbering-and-namespace-convention.md).

## Capability-specification execution lots

| Execution lot | Scope | Status |
|---|---|---|
| **GOV-1** | Action Requests, Policy, Authorities and Decisions | **PENDING POST-PUBLICATION VERIFICATION** |
| GOV-2 | Playbooks, Response Runs, Execution, Verification and Rollback | **NOT STARTED** |
| GOV-3 | Audit Trail, Response Metrics and Govern Closure | **NOT STARTED** |

`GOV-1`, `GOV-2` and `GOV-3` are execution-lot identifiers. They are not roadmap phases.

## GOV-1 scope

GOV-1 defines the functional chain from governed request intake to a no-effect execution handoff:

`Finding / Incident / Detection Engineering / Threat Intelligence / analysis context → Action Request → Intake → Response Inbox → Action Center → Policy Evaluation → conflict/Exception assessment → authority requirement → approver eligibility/SoD → Approval Request → Approval → Decision Preparation → Decision → Execution Handoff Package`.

GOV-1 covers modules:
1. Response Inbox;
2. Action Center;
3. Decision Register;
4. Policy Gates;
5. Approvals & Authorities.

It reads modules 6–9 only to preserve future boundaries:
6. Playbooks — GOV-2;
7. Runs & Rollback — GOV-2;
8. Audit Trail — GOV-3;
9. Response Metrics — GOV-3.

## GOV-1 capability set

Exactly `CAP-GOV-001..016` are allocated and functionally specified. They are all `draft` / `defined` / `planned`.

Prepared counts before remote closure verification:
- GOV-1 capabilities: **16**;
- numbered capability sections: **432**;
- mandatory S8/S9/S10/S13/S16/S17 tables: **96**;
- empty/generic mandatory tables: **0**;
- GOV-2/GOV-3 capabilities: **0**.

## Ownership and stop line

Govern owns Action Request processing, Policy Evaluation/conflict/exception governance, contextual authority, Approval and Decision. Command keeps Incident/general Work Queue; Investigate keeps Case/Evidence/Finding; Studio keeps Workflow/Human Gate/Automation Run; Settings keeps user/role/group/tenant/secret administration; Endpoint/runtime owners keep technical execution; Shared keeps generic engines.

GOV-1 performs **no target execution**. `Execution Handoff Package` ≠ `Response Run`; `Approval` ≠ `Decision`; `approve` ≠ `execute`. No Result or rollback execution is produced by GOV-1.

## GOV-1 functional commits

Starting baseline:
- `a6adf28aa0fa64b917a0a37be37de2a4cb28b541` — `docs: record Command Phase 4A post-publication revalidation`.

Published GOV-1 functional commits:
1. `587ac4f7f0e437b13c6276e493ad9bf1f3bf6dbd` — `docs: establish Govern request policy and authority boundaries`;
2. `c45f067c6d108f41d29f6140c1298bee4a263712` — `docs: define Govern intake action requests and policy evaluation`;
3. `0dc62e1716909f5702bfbb19488f0f9a0d8530d6` — `docs: specify Govern approvals authorities and separation of duties`;
4. `35f2e5d6f30f4b7cb06c476dc9ea358b4288f349` — `docs: document Govern decisions conditions and execution handoff`;
5. `docs: update Govern foundation traceability and quality gates` — exact squash SHA recorded only after publication.

## Pre-publication status rule

The functional content, counts and local gates can be prepared before the fifth commit, but remote publication-dependent gates remain `PENDING`. GOV-1 MUST NOT be marked `PASS AFTER POST-PUBLICATION VERIFICATION` until the fifth functional commit is on the canonical branch and PR/head/README/main/commit-chain/Command-non-regression checks have been performed remotely.

## Requirements and OPEN decisions

- Requirement IDs remain **122 — 99 conform / 20 partial / 3 absent / 0 contradictory**; GOV-1 adds documentary evidence but no state promotion.
- Open decisions remain **18**; GOV-1 creates/closes **0**.
- `OPEN-007`, `OPEN-013` and `OPEN-015` remain explicitly open and are consumed rather than resolved.

## Permissions / objects / technology

GOV-1 identifies functional permission needs and local object/concept consumption only. It creates no complete object schema, JSON Schema, final state machine, final RBAC/ABAC matrix, API, protocol, Policy engine, authority engine, execution engine or product code.

## Acceptance for GOV-1

GOV-1 can be promoted to PASS only if the published chain contains the five required functional commits, `CAP-GOV-001..016` conform at 432/432 sections and 96/96 mandatory tables, all 180 quality gates pass, Command non-regression remains intact, PR #2 remains open/Draft/unmerged, README/main are unchanged and no GOV-2/GOV-3/implementation scope was started.

## Current next action

Publish the fifth GOV-1 traceability/quality commit, verify remote state, then record post-publication evidence. Do **not** start GOV-2 or GOV-3 as part of GOV-1 closure.