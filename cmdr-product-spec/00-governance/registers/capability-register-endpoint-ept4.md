---
id: capability-register-endpoint-ept4
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-11
source-of-truth: registry
---
# Capability Register — Endpoint EPT-4

Execution lot: **EPT-4 — Collection and Live Response Technical Execution** under Delivery Roadmap Phase 5 — Studio and Endpoint. All rows are `draft / defined / planned`. Documentary definition does not prove implementation, platform support or response authority.

| ID | Capability | Owner | Primary roles | Primary concepts | Key consumers | Requirements | OPEN | Delivery |
|---|---|---|---|---|---|---|---|---|
| CAP-EPT-047 | Collection Request, Eligibility and Authority Boundary | Endpoint Agent | DFIR Analyst, Endpoint Operator | technical intake, eligibility, authority ref | Investigate, Govern | REQ-PROD-006/014/018/019/020, REQ-SEC-001/002 | OPEN-008/013/014/015 | defined / planned |
| CAP-EPT-048 | Collection Scope, Target Selection and Execution Planning | Endpoint Agent | DFIR Analyst, Endpoint Operator | technical plan, planned item | Collection executors | REQ-PROD-006/014/018/019, REQ-SEC-001/002 | OPEN-008/013/014 | defined / planned |
| CAP-EPT-049 | File and Filesystem Collection | Endpoint Agent | DFIR Analyst, Evidence Reviewer | file attempt, Collection Item | Investigate | REQ-PROD-014/018/020, REQ-SEC-001/004 | OPEN-008/013/014 | defined / planned |
| CAP-EPT-050 | Process and System-State Collection | Endpoint Agent | Case Analyst, Response Operator | fresh process/system snapshot | Investigate | REQ-PROD-014/018/019, REQ-SEC-001 | OPEN-008/013/014 | defined / planned |
| CAP-EPT-051 | Memory Acquisition and Collection Boundary | Endpoint Agent | DFIR Analyst, Govern Reviewer | memory attempt/item | Investigate Memory | REQ-PROD-014/018/020/052/055, REQ-SEC-001/004 | OPEN-005/008/013/014 | defined / planned |
| CAP-EPT-052 | Network and Packet Capture Collection Boundary | Endpoint Agent | DFIR Analyst, Response Operator | capture attempt/item/loss | Investigate Network | REQ-PROD-014/018/020/055, REQ-SEC-001/004 | OPEN-008/013/014 | defined / planned |
| CAP-EPT-053 | Collection Progress, Partial Completion, Failure and Cancellation | Endpoint Agent | Analyst, Endpoint Operator | operation state, progress, retry/resume eligibility | Investigate, Shared Jobs | REQ-PROD-014/018/019, REQ-SEC-002 | OPEN-008/013/015 | defined / planned |
| CAP-EPT-054 | Collection Packaging, Completeness, Integrity and Metadata | Endpoint Agent | Evidence Reviewer, Auditor | Collection Package, completeness, integrity metadata | Investigate, Trust | REQ-PROD-014/020, REQ-OBJ-004, REQ-SEC-004 | OPEN-014 | defined / planned |
| CAP-EPT-055 | Collection Transfer, Delivery and Consumer Handoff | Endpoint Agent | DFIR Analyst, Endpoint Operator | transfer, acknowledgement, handoff | Investigate, Shared | REQ-PROD-014/018/020, REQ-SEC-001/004 | OPEN-008/014/015 | defined / planned |
| CAP-EPT-056 | Live Response Session Definition, Eligibility and Target Binding | Endpoint Agent | Response Operator, Endpoint Operator | technical session, binding, eligibility | Investigate, Govern | REQ-PROD-006/014/018/019/020, REQ-SEC-001/002 | OPEN-008/013/015/017 | defined / planned |
| CAP-EPT-057 | Live Response Session Lifecycle and Runtime State | Endpoint Agent | Response Operator, Session Supervisor | technical session state | Investigate | REQ-PROD-014/018/019, REQ-SEC-002 | OPEN-008/013/015 | defined / planned |
| CAP-EPT-058 | Technical Command Request and Invocation Boundary | Endpoint Agent | Response Operator, Govern Reviewer | command request, invocation eligibility | Govern, Studio | REQ-PROD-006/014/018/019/020, REQ-SEC-001/002 | OPEN-008/013/015/017 | defined / planned |
| CAP-EPT-059 | Interactive Command and Shell Execution | Endpoint Agent | Response Operator, Endpoint Operator | execution attempt, transcript/runtime context | Investigate, Govern | REQ-PROD-014/018/019, REQ-SEC-001/002 | OPEN-008/013/015/017 | defined / planned |
| CAP-EPT-060 | Script Execution and Runtime Context Boundary | Endpoint Agent | Response Operator, Govern Reviewer | script ref, runtime, execution attempt | Investigate, Studio, Govern | REQ-PROD-014/018/019, REQ-SEC-001/002/004 | OPEN-008/013/014/015/017 | defined / planned |
| CAP-EPT-061 | Live Response File Read, Download and Transfer Operations | Endpoint Agent | Response Operator, DFIR Analyst | file transfer/read/download/temp upload | Investigate, Govern | REQ-PROD-014/018/020, REQ-SEC-001/004 | OPEN-008/013/014/015 | defined / planned |
| CAP-EPT-062 | Technical Execution Output, Error, Timeout and Cancellation Semantics | Endpoint Agent | Response Operator, Case Analyst | technical output/error/termination knowledge | Investigate, Govern, Studio | REQ-PROD-014/018/019/020, REQ-SEC-001/002/004 | OPEN-013/015 | defined / planned |
| CAP-EPT-063 | Live Response Operator Control, Session Closure and Audit Context | Endpoint Agent | Session Supervisor, Auditor | operator control, closure context | Investigate, Govern, Shared | REQ-PROD-014/018/019, REQ-SEC-001/002/004 | OPEN-008/013/015 | defined / planned |
| CAP-EPT-064 | Collection and Live Response Provenance and Cross-Product Execution Contracts | Endpoint Agent | Auditor, DFIR Analyst, Govern Reviewer | EPT-4 provenance chain, handoff refs | Investigate, Govern, Studio, Settings, Shared | REQ-PROD-006/014/018/019/020, REQ-OBJ-004/007, REQ-SEC-001/002/004 | OPEN-008/013/014/015/017 | defined / planned |

## Totals
**18 capabilities / 486 numbered sections / 108 mandatory tables / at least 54 GWT.** Duplicate IDs 0; recycled IDs 0; owner conflicts 0; empty/generic mandatory tables 0. Endpoint cumulative after EPT-4: **64 capabilities / 1728 sections / 384 mandatory tables**.

`CAP-EPT-001..046` remain unchanged. `CAP-EPT-065+` is not allocated or reserved by EPT-4. EPT-5 and EPT-6 remain **NOT STARTED**. OPEN-008, OPEN-014, OPEN-015 and OPEN-017 remain open.