---
id: phase-4a-command-current-revalidation
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-09
source-of-truth: quality-report
requirements: [REQ-PROD-003, REQ-PROD-004, REQ-PROD-005, REQ-PROD-006, REQ-PROD-008, REQ-PROD-009, REQ-PROD-010, REQ-PROD-013, REQ-PROD-019, REQ-PROD-021, REQ-OBJ-001, REQ-OBJ-012, REQ-UX-008, REQ-UX-009]
open_decisions: [OPEN-006, OPEN-010, OPEN-013]
---
# Capability Specification Phase 4A — Command Current Revalidation

## Verdict

**PASS AFTER POST-PUBLICATION VERIFICATION — 60/60 mandatory controls PASS, 0 PENDING, 0 FAIL.**

This validates documentary functional conformance only. It validates no software implementation, API, protocol, runtime, final object schema, atomic permission matrix or detailed screen implementation.

## Publication evidence

- repository: `tobianahillel-afk/cmdr`;
- canonical branch: `docs/cmdr-product-spec-foundation`;
- base: `main`;
- PR: #2;
- exact initial SHA: `edbc67be43f971e4c49e91518e3e21517ef69410`;
- initial commit: `docs: reconcile capability and delivery roadmap phase numbering`;
- corrective functional commit: `caa786d8adb8cace93ddfe3f1573a4c79a6dd900` — `docs: revalidate Command Phase 4A traceability and capability register`;
- functional publication ancestry: **1 commit ahead / 0 behind**, same merge base;
- PR #2 after publication: open, Draft, unmerged and mergeable;
- root README on canonical branch and `main`: exact `# cmdr`, same blob `901c74cda52e28b5ff7fc425ddf28ef89f3ad875`;
- `main`: unchanged from `bc1ec59e5e79ccbaba291e2984b0eb7d5e54129c`;
- workflow runs on the functional commit: none;
- commit statuses on the functional commit: none;
- publication used a squash commit onto the canonical documentation branch; no force-push, rebase, reset or history rewrite.

The exact SHA of this post-publication verification-record correction is verified externally against the canonical branch and PR #2 after publication rather than embedded self-referentially in this commit.

## Why this is a corrective revalidation

Capability Specification Phase 4A had already been historically completed and validated before later Investigate work. The current run therefore did not recreate or renumber Command capabilities. It re-read the current Command corpus and corrected only drift in active traceability/register views.

Three material current-state gaps were corrected:

1. the active Requirements Traceability Matrix contained Investigate evidence but no `CAP-CMD-*` evidence ranges;
2. the Command capability-register shard used generic `Command roles` and `canonical specification` dependency values even though every canonical capability file contained specific roles and dependencies;
3. the global Dependency Register had no Command summary families even though the canonical Command dependency map already carried the detailed graph.

## Source audit

### Source material and governance
Directly re-read: all requested source-material files; documentation rules; source-of-truth; ownership; terminology; status lifecycle; dependency register; legacy screen mapping; capability/object/screen registers; and ADR-0001 through ADR-0007 relevant to product separation, source of truth, object chain, screens, page/view/filter rules, Endpoint ownership and Studio agentic placement.

### Product Vision and Phase 3 architecture
Directly re-read: all requested Product Vision files, global/queue/decision shells, table/filter/Saved View/Context Bar/Inspector/Activity/Trace/Timeline/Automation Tray components, cross-product/provenance/dangerous-action patterns and all requested Experience Architecture sources.

### Command
All Command files in the current PR inventory were audited. The repository contains **69 files under `06-command/`**: **61 active documents** and **8 deprecated migration documents/pointers**. The audit directly covered all 27 capability files, five canonical modules, top-level Command maps, ten active Command screen sources and the five deprecated Work Queue aliases.

### Objects and dependent products
Directly re-read canonical Incident, Task, Alert, Signal, Case, Finding, Action Request, Decision, Response Run and Result objects; the shared Business Service Catalog and Reporting Engine; audit semantics; and README/product boundaries for Investigate, Govern, Studio, Platform Settings, Endpoint Agent and Shared Capabilities.

## Canonical capability framework

The existing canonical template remains `templates/capability-specification-template.md`; no competing `capability-template.md` was created. It requires stable front matter, 27 substantive sections, functional inputs/outputs, object reads/mutations, classes 0–4, automation/AI alternatives, states, provenance, permissions, limits, delivery classification and at least three Given/When/Then criteria.

Command ID ranges remain:
- `CAP-CMD-001..099` — Mission Control;
- `CAP-CMD-100..199` — Incidents and Work Queue;
- `CAP-CMD-200..299` — Risk and Coverage;
- `CAP-CMD-300..399` — Readiness and Operations;
- `CAP-CMD-400..499` — Customers and Delivery.

IDs are immutable and never recycled.

## Command capability conformance

| Measure | Result |
|---|---:|
| Command capability files | 27 / 27 |
| Numbered sections | 729 / 729 |
| Mandatory S8/S9/S10/S13/S16/S17 tables | 162 / 162 |
| Unique IDs | 27 / 27 |
| Defined | 26 |
| Proposed | 1 |
| Planned delivery mode | 27 |
| Native current claims | 0 |
| Integrated current claims | 0 |
| Temporary-integration current claims | 0 |
| Capabilities without owner | 0 |
| Capabilities without named users | 0 |
| Capabilities without inputs | 0 |
| Capabilities without outputs | 0 |
| Capabilities without object semantics | 0 |
| Capabilities without classified actions | 0 |
| Capabilities without alternative without AI | 0 |
| Capabilities without Given/When/Then criteria | 0 |
| Duplicate Capability IDs | 0 |
| Concurrent active owners | 0 |

## Command modules

### Mission Control — 6 capabilities
`CAP-CMD-001..006`: Situation Overview, Priority Management, Situation Timeline, Handover, Operational Blockers, Recent Results and Outcomes. Command owns operational situation/coordination; Govern projections remain read-only; blockers reuse Incident/Task/relations rather than a new canonical object.

### Incidents and Work Queue — 10 capabilities
`CAP-CMD-101..110`: one Work Queue workspace, assignment, operational ownership, priority/severity, SLA, Incident coordination, Task coordination, bulk operations, freshness/staleness and escalation.

The exact six system views remain `All`, `Incidents`, `Tasks`, `Unassigned`, `SLA Risk`, `My Work`. `Team Load` is not a system view. Deprecated aliases CMD-IWQ-001..005 remain migration pointers; CMD-IWQ-006 Incident Detail remains active.

### Risk and Coverage — 5 capabilities
`CAP-CMD-201..205`: Service Context, Exposure Overview, Coverage Overview, Business Impact Context and Risk Prioritization Context. Command consumes external/shared Service, Exposure and Coverage context without becoming a CMDB, scanner or Detection Engineering owner; no universal opaque risk score is introduced.

### Readiness and Operations — 5 capabilities
`CAP-CMD-301..305`: Readiness Overview, Exercise Coordination, Improvement Actions, Operational Plans and Capability Readiness. Operational exercises are distinct from Studio/Sandbox simulation; readiness is distinct from document status and delivery mode; improvement actions reuse Task.

### Customers and Delivery — 1 proposed capability
`CAP-CMD-401` remains `proposed`, `planned`, deployment-dependent under OPEN-006. Reporting Engine remains Shared; internal deployments do not require a Customer concept.

## Objects and ownership

Command owns Incident and operational Task coordination. It consumes but does not own Case, Evidence, Finding, Action Request, Decision, Response Run, Result, Workflow, Automation Run or Endpoint Agent Fleet. `service.md`, `exposure.md`, `report.md` and `audit-record.md` are not created by this phase; their current shared/source semantics and future object work remain explicit.

## Actions and AI

All Command capability actions are classified C0–C4. C0 observation is local. C1 collection generally hands off to Investigate. C2 reversible changes remain subject to OPEN-013. C3/C4 real response actions are request/context only and require Govern authority/execution.

Every essential capability has a manual and/or deterministic path without a model provider. AI may summarize or propose but cannot silently modify priority/assignment, create a Decision, self-approve, own the Work Queue, hide source data, grant itself permission or bypass Govern.

## Dependencies

`06-command/functional-dependency-map.md` contains **161 detailed dependency edges**. The global Dependency Register contains **10 Command dependency families** for registry-level traceability without duplicating the detailed graph. Dependencies do not transfer ownership.

## Screen status

- active Command screen specs read: **10**;
- detailed screen rewrites in this corrective run: **0**;
- screen files modified: **0**;
- new Screen IDs: **0**;
- deprecated Work Queue aliases preserved: **5**;
- existing generic/deferred Phase-3 screen specs with Q27 placeholders: **10**.

Those ten screen placeholders are explicitly deferred to the later screen phase and are outside the Phase 4A capability/module rewrite scope. Active capability/module specifications have **0 generic placeholders**.

## Before / after corrective metrics

| Measure | Before | After |
|---|---:|---:|
| Command files | 69 | 69 |
| Active Command docs | 61 | 61 |
| Deprecated migration docs | 8 | 8 |
| Command Capability IDs | 27 | 27 |
| Registered Command capabilities | 27 | 27 |
| Missing owner/user/input/output/object/action/no-AI/AC | 0 | 0 |
| Defined / proposed | 26 / 1 | 26 / 1 |
| Native / integrated / temporary / planned / out-of-scope | 0 / 0 / 0 / 27 / 0 | 0 / 0 / 0 / 27 / 0 |
| Detailed functional dependency edges | 161 | 161 |
| Global Command dependency families | 0 | 10 |
| Command evidence ranges in active Requirements Matrix | 0 | 5 |
| Generic role summaries in Command registry | 27 | 0 |
| Generic dependency summaries in Command registry | 27 | 0 |
| Duplicate IDs / concurrent owners | 0 / 0 | 0 / 0 |
| Screens modified | 0 | 0 |
| Objects modified | 0 | 0 |
| Permission files modified | 0 | 0 |
| Product code / APIs / protocols / fonts | 0 / 0 / 0 / 0 | 0 / 0 / 0 / 0 |
| Requirements conform / partial / absent / contradictory | 99 / 20 / 3 / 0 | 99 / 20 / 3 / 0 |
| Requirement IDs | 122 | 122 |
| Open decisions | 18 | 18 |
| New OPEN | 0 | 0 |
| Targeted broken links introduced | 0 | 0 |
| Targeted empty files | 0 | 0 |

## Mandatory controls — 60/60 PASS

| # | Control | Verdict | Evidence |
|---:|---|---|---|
| 1 | Exact initial head SHA recovered | PASS | `edbc67be43f971e4c49e91518e3e21517ef69410` |
| 2 | Correct PR | PASS | PR #2 |
| 3 | Correct branch | PASS | `docs/cmdr-product-spec-foundation` |
| 4 | Correct base | PASS | `main` |
| 5 | PR open after canonical publication | PASS | PR #2 state open |
| 6 | PR Draft after canonical publication | PASS | PR #2 draft true |
| 7 | PR unmerged after canonical publication | PASS | PR #2 merged false |
| 8 | Root README unchanged after canonical publication | PASS | exact `# cmdr`, unchanged blob |
| 9 | `main` unchanged after canonical publication | PASS | `main` still at `bc1ec59e5e79ccbaba291e2984b0eb7d5e54129c` |
| 10 | All active Command files read | PASS | 61 active documents; deprecated pointers also audited |
| 11 | Required governance files read | PASS | source material, governance, registers and ADRs audited |
| 12 | Required Command objects read | PASS | Incident/Task/Alert/Signal and cross-product owners audited |
| 13 | All active Command screens read without detailed rewrite | PASS | 10 read; 0 rewritten |
| 14 | Five Work Queue screens remain deprecated | PASS | CMD-IWQ-001..005 preserved |
| 15 | Canonical capability model substantive | PASS | existing canonical template retained |
| 16 | Capability ID convention defined | PASS | stable CAP-CMD ranges |
| 17 | Capability Register substantive | PASS | global contract + detailed Command shard |
| 18 | All Command capabilities registered | PASS | 27/27 |
| 19 | Every capability has owner | PASS | 27/27 |
| 20 | Every capability has users | PASS | 27/27 |
| 21 | Every capability has inputs | PASS | 27/27 S8 |
| 22 | Every capability has outputs | PASS | 27/27 S16 |
| 23 | Read and modified objects distinguished | PASS | 27/27 S9/S10 |
| 24 | Every capability has classified actions | PASS | 27/27 S12; classes 0–4 |
| 25 | Every capability has capability-specific functional states | PASS | 27/27 S14 |
| 26 | Human/deterministic/workflow/agent behavior distinguished | PASS | canonical S13 + Command AI model |
| 27 | Every capability has no-AI alternative | PASS | 27/27 |
| 28 | Every capability has Given/When/Then criteria | PASS | at least three per file |
| 29 | Command retains Incident/coordination ownership | PASS | ownership register + capability boundaries |
| 30 | Command does not own Case | PASS | Investigate projection only |
| 31 | Command does not own Evidence | PASS | Investigate-owned |
| 32 | Command does not own Finding | PASS | Investigate-owned |
| 33 | Command does not own Decision | PASS | Govern-owned |
| 34 | Command does not own Response Run | PASS | Govern-owned |
| 35 | Work Queue remains one workspace | PASS | CAP-CMD-101 + Saved Views |
| 36 | Exactly six system views | PASS | All, Incidents, Tasks, Unassigned, SLA Risk, My Work |
| 37 | Team Load not reintroduced as system view | PASS | deprecated alias only |
| 38 | Customers and Delivery remains proposed | PASS | CAP-CMD-401 proposed/planned |
| 39 | OPEN-006 remains open | PASS | no decision change |
| 40 | OPEN-010 remains open | PASS | no decision change |
| 41 | OPEN-013 remains open | PASS | no decision change |
| 42 | Shared Capabilities not redefined | PASS | consumed by reference/source ownership |
| 43 | Screens not rewritten in detail | PASS | 0 screen modifications |
| 44 | Objects not massively rewritten | PASS | 0 object modifications |
| 45 | Detailed permissions not finalized | PASS | functional needs only; 0 permission-file change |
| 46 | No technical protocol created | PASS | documentation-only corrective diff |
| 47 | No API created | PASS | documentation-only corrective diff |
| 48 | No product code added | PASS | Markdown-only |
| 49 | No font file added | PASS | none in corrective diff |
| 50 | No new generic placeholder | PASS | new/updated capability-governance docs substantive |
| 51 | No targeted empty file | PASS | all corrective docs substantive |
| 52 | No targeted local link introduced broken | PASS | referenced canonical paths audited |
| 53 | Referenced Requirement IDs exist | PASS | source inventory/matrix retained at 122 IDs |
| 54 | No duplicate Capability ID | PASS | 27 unique CAP-CMD IDs |
| 55 | No concurrent active owner | PASS | ownership register and 27 contracts agree |
| 56 | No document artificially Validated | PASS | documentary statuses remain draft; phase verdict separate |
| 57 | Traceability Matrix updated | PASS | five Command evidence ranges restored |
| 58 | Baseline updated | PASS | current Phase 4A metrics appended with history preserved |
| 59 | Capability Register updated | PASS | global contract + specific Command shard evidence |
| 60 | Phase 4B not started or functionally modified by this run | PASS | historical Phase 4B content preserved; corrective diff contains no CAP-INV file |

## Open decisions

The programme retains **18 open decisions**; this Phase 4A corrective run creates and closes **0**. Command-specific emphasis:

- **OPEN-006 — Customers and Delivery**: remains open; CAP-CMD-401 stays proposed/deployment-dependent.
- **OPEN-010 — final density by role and activity**: remains open; detailed screen density belongs to later UX/screen work.
- **OPEN-013 — default governance for class-2 actions**: remains open; Phase 4A identifies C2 actions but does not choose the final step-up/approval policy.

Other current OPEN decisions retain their canonical dispositions; no identifier is reused.

## Final status

- Capability Specification Phase 4A — Command: **PASS**.
- Capability Specification Phase 4B — Investigate: **PASS and unchanged by this run**.
- Capability Specification Phase 4 global maturity: **PARTIAL**.
- Global repository maturity: **PARTIAL**.
- Delivery Roadmap Phase 4 — Govern capability specification: **NOT STARTED**.

No later capability-specification phase is started by this verification record.