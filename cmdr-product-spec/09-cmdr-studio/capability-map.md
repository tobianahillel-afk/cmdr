---
id: studio-capability-map
domain: 09-cmdr-studio
status: draft
owner: CMDR Studio Product Lead
updated: 2026-08-09
source-of-truth: canonical
---
# Studio Capability Map

## STD-1 identity
Parent: **Delivery Roadmap Phase 5 — Studio and Endpoint**. Execution lot: **STD-1 — Studio Foundations — Tools, Skills, Library and Ownership Contracts**. STD-1 is not a roadmap phase and creates no Phase 5A/5B.

## Capabilities
| ID | Capability | Functional family | Delivery |
|---|---|---|---|
| CAP-STD-001 | Studio Library and Asset Catalog | Library | defined / planned |
| CAP-STD-002 | Studio Asset Metadata, Ownership and Classification | Library | defined / planned |
| CAP-STD-003 | Tool Definition and Functional Capability Contract | Tool foundations | defined / planned |
| CAP-STD-004 | Tool Input, Parameter and Validation Contract | Tool foundations | defined / planned |
| CAP-STD-005 | Tool Output, Result and Error Contract | Tool foundations | defined / planned |
| CAP-STD-006 | Tool Versioning, Compatibility and Deprecation | Tool foundations | defined / planned |
| CAP-STD-007 | Tool Permission, Risk and Execution Eligibility | Tool foundations | defined / planned |
| CAP-STD-008 | Tool Call Request and Invocation Context | Tool Call | defined / planned |
| CAP-STD-009 | Tool Call Lifecycle, Technical Outcome and Provenance | Tool Call | defined / planned |
| CAP-STD-010 | Skill Definition and Reusable Capability Contract | Skills | defined / planned |
| CAP-STD-011 | Skill Composition, Dependencies and Preconditions | Skills | defined / planned |
| CAP-STD-012 | Skill Input, Output and Parameter Contract | Skills | defined / planned |
| CAP-STD-013 | Skill Versioning, Lifecycle and Deprecation | Skills | defined / planned |
| CAP-STD-014 | Skill Discovery, Reuse and Consumer Linking | Skills | defined / planned |
| CAP-STD-015 | Provider, Runtime, Integration and Secret Reference Boundaries | Foundation boundary | defined / planned |
| CAP-STD-016 | Studio Foundations Cross-Product Contracts and Provenance | Foundation boundary | defined / planned |

## Structural total
16 capabilities × 27 numbered sections = **432 sections**. Six mandatory tables per capability = **96 mandatory tables**. No STD-2/STD-3/STD-4 or Endpoint capability belongs to this map.

## STD-2 addendum — Workflow Builder & Orchestration

The STD-1 map above is preserved verbatim as the pre-STD-2 snapshot. The historical “No STD-2 ... belongs to this map” sentence applies to the STD-1 snapshot only.

| ID | Capability | Functional family | Delivery |
|---|---|---|---|
| CAP-STD-017 | Workflow Definition and Functional Contract | Workflow foundations | defined / planned |
| CAP-STD-018 | Workflow Builder Session and Editing Context | Builder | defined / planned |
| CAP-STD-019 | Workflow Inputs, Outputs, Variables and Data Context | Workflow data | defined / planned |
| CAP-STD-020 | Workflow Steps, Nodes and Dependency Graph | Graph | defined / planned |
| CAP-STD-021 | Tool and Skill Step Composition | Composition | defined / planned |
| CAP-STD-022 | Conditions, Branches and Deterministic Decision Logic | Control flow | defined / planned |
| CAP-STD-023 | Data Mapping, Transformation and Context Propagation | Mapping | defined / planned |
| CAP-STD-024 | Subworkflow Composition and Reusable Workflow References | Subworkflow | defined / planned |
| CAP-STD-025 | Workflow Ordering, Parallelism and Concurrency Constraints | Ordering/concurrency | defined / planned |
| CAP-STD-026 | Workflow Error Paths and Exception Handling | Errors | defined / planned |
| CAP-STD-027 | Workflow Retry, Idempotency and Duplicate-Execution Protection | Retry/idempotency | defined / planned |
| CAP-STD-028 | Workflow Partial Success and Compensation Semantics | Partial/compensation | defined / planned |
| CAP-STD-029 | Human Gate Step and Govern Boundary | Human Gate boundary | defined / planned |
| CAP-STD-030 | Workflow Validation, Compatibility and Readiness Assessment | Validation/readiness | defined / planned |
| CAP-STD-031 | Workflow Versioning and Compatibility | Versioning | defined / planned |
| CAP-STD-032 | Workflow Draft, Review and Pre-Publishing Lifecycle | Pre-publish lifecycle | defined / planned |
| CAP-STD-033 | Workflow Provenance and Cross-Product Orchestration Contracts | Provenance/boundary | defined / planned |

STD-2 structural total: **17 capabilities / 459 sections / 102 mandatory tables**. Studio cumulative: **33 capabilities / 891 sections / 198 mandatory tables**. STD-3/STD-4 and Endpoint remain NOT STARTED.

## STD-3 current addendum — Agents, Human Gates & Runtime Control

The STD-1 and STD-2 sections above remain exact historical snapshots. Their pre-STD-3 `NOT STARTED` statements are historical evidence only.

| ID | Capability | Functional family | Delivery |
|---|---|---|---|
| CAP-STD-034 | Automation Agent Definition and Functional Contract | Automation Agent | defined / planned |
| CAP-STD-035 | Agent Objectives, Constraints and Execution Context | Agent context | defined / planned |
| CAP-STD-036 | Agent Tool, Skill and Resource Access Governance | Agent access | defined / planned |
| CAP-STD-037 | Agent Team Composition, Roles and Coordination | Agent Team | defined / planned |
| CAP-STD-038 | Agent Planning, Step Proposal and Bounded Autonomy | Agent planning | defined / planned |
| CAP-STD-039 | Agent Oversight, Intervention and Escalation | Oversight | defined / planned |
| CAP-STD-040 | Human Gate Request and Review Context | Human Gate | defined / planned |
| CAP-STD-041 | Human Gate Lifecycle, Response, Expiration and Govern Boundary | Human Gate | defined / planned |
| CAP-STD-042 | Automation Run Creation and Execution Context | Automation Run | defined / planned |
| CAP-STD-043 | Automation Run Lifecycle and State Management | Automation Run | defined / planned |
| CAP-STD-044 | Automation Run Steps, Attempts and Tool Call Coordination | Run execution | defined / planned |
| CAP-STD-045 | Runtime Queueing, Scheduling and Concurrency Control | Runtime control | defined / planned |
| CAP-STD-046 | Runtime Start, Pause, Resume, Stop and Cancellation Control | Runtime control | defined / planned |
| CAP-STD-047 | Runtime Error, Timeout, Retry and Partial Completion Handling | Runtime errors | defined / planned |
| CAP-STD-048 | Runtime Context, State and Data Propagation | Runtime context | defined / planned |
| CAP-STD-049 | Control Room Monitoring and Runtime Intervention | Control Room | defined / planned |
| CAP-STD-050 | Automation Run Outcome, Consumer Handoff and Follow-up | Runtime outcome | defined / planned |
| CAP-STD-051 | Studio Runtime Provenance and Cross-Product Contracts | Provenance/boundary | defined / planned |

STD-3 content total: **18 capabilities / 486 sections / 108 mandatory tables**. Studio cumulative: **51 capabilities / 1377 sections / 306 mandatory tables**. Post-publication documentary verification is in progress; STD-4 and Endpoint remain NOT STARTED.

## STD-4 addendum — Assurance & Lifecycle

The STD-1/2/3 text above is preserved as historical evidence; earlier `STD-4 NOT STARTED` statements are historical snapshots only.

| ID | Capability | Functional family | Delivery |
|---|---|---|---|
| CAP-STD-052 | Studio Evaluation Definition and Success Criteria | Evaluation | defined / planned |
| CAP-STD-053 | Evaluation Suite, Cases, Inputs and Expected Outcomes | Evaluation suites | defined / planned |
| CAP-STD-054 | Deterministic and Assisted Evaluation Execution | Evaluation execution | defined / planned |
| CAP-STD-055 | Studio Simulation and No-Effect Scenario Execution | Simulation | defined / planned |
| CAP-STD-056 | Workflow, Agent and Tool Regression Testing | Regression | defined / planned |
| CAP-STD-057 | Reliability, Failure-Mode and Safety Assessment | Reliability/safety | defined / planned |
| CAP-STD-058 | Permission, Governance and Boundary Assurance | Boundary assurance | defined / planned |
| CAP-STD-059 | Evaluation Result, Comparison and Regression Analysis | Evaluation result | defined / planned |
| CAP-STD-060 | Studio Quality Gates and Readiness Assessment | Readiness | defined / planned |
| CAP-STD-061 | Publishing Candidate and Release Package Preparation | Publishing candidate | defined / planned |
| CAP-STD-062 | Publishing Review, Approval Boundary and Publication Lifecycle | Publication | defined / planned |
| CAP-STD-063 | Studio Asset Release, Version Promotion and Channel Management | Release/promotion | defined / planned |
| CAP-STD-064 | Deployment Target and Compatibility Assessment | Deployment compatibility | defined / planned |
| CAP-STD-065 | Studio Deployment, Promotion and Rollout Lifecycle | Deployment lifecycle | defined / planned |
| CAP-STD-066 | Deployment Health, Reversion and Previous-Version Recovery | Health/reversion | defined / planned |
| CAP-STD-067 | Studio Asset Deprecation, Retirement and Migration | Deprecation/migration | defined / planned |
| CAP-STD-068 | Studio Assurance, Lifecycle Provenance and Capability Closure | Provenance/closure | defined / planned |

STD-4 content total: **17 capabilities / 459 sections / 102 mandatory tables**. Studio cumulative after STD-4 content: **68 capabilities / 1836 sections / 408 mandatory tables**. Final Studio PASS is pending post-publication 220-gate verification. Endpoint remains NOT STARTED.