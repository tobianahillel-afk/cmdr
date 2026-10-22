---
id: capability-register-endpoint-ept4
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-11
source-of-truth: registry
---
# Capability Register — Endpoint EPT-4

Execution lot: **EPT-4 — Collection and Live Response Technical Execution** under Delivery Roadmap Phase 5 — Studio and Endpoint. All rows are `draft / defined / planned`. Documentary definition does not prove implementation.

| ID | Capability | Owner | Primary roles | Primary concepts | Key consumers | Requirements | OPEN | Delivery |
|---|---|---|---|---|---|---|---|---|
| CAP-EPT-047 | Collection Request, Eligibility and Authority Boundary | Endpoint Agent | DFIR Analyst, Endpoint Operator | technical intake, eligibility | Investigate, Govern | REQ-PROD-006/014/018/019/020, REQ-SEC-001/002 | OPEN-008/013/014/015 | defined / planned |
| CAP-EPT-048 | Collection Scope, Target Selection and Execution Planning | Endpoint Agent | DFIR Analyst, Endpoint Operator | technical plan, planned item | EPT-4 collectors | REQ-PROD-006/014/018/019, REQ-SEC-001/002 | OPEN-008/013/014 | defined / planned |
| CAP-EPT-049 | File and Filesystem Collection | Endpoint Agent | DFIR Analyst, Evidence Reviewer | file attempt, Collection Item | Investigate | REQ-PROD-014/018/020, REQ-SEC-001/004 | OPEN-008/013/014 | defined / planned |
| CAP-EPT-050 | Process and System-State Collection | Endpoint Agent | Case Analyst, Response Operator | fresh snapshot | Investigate | REQ-PROD-014/018/019, REQ-SEC-001 | OPEN-008/013/014 | defined / planned |
| CAP-EPT-051 | Memory Acquisition and Collection Boundary | Endpoint Agent | DFIR Analyst, Govern Reviewer | memory attempt/item | Investigate Memory | REQ-PROD-014/018/020/052/055, REQ-SEC-001/004 | OPEN-005/008/013/014 | defined / planned |
| CAP-EPT-052 | Network and Packet Capture Collection Boundary | Endpoint Agent | DFIR Analyst, Response Operator | capture attempt/item/loss | Investigate Network | REQ-PROD-014/018/020/055, REQ-SEC-001/004 | OPEN-008/013/014 | defined / planned |
| CAP-EPT-053 | Collection Progress, Partial Completion, Failure and Cancellation | Endpoint Agent | Analyst, Endpoint Operator | operation state/progress | Investigate, Shared Jobs | REQ-PROD-014/018/019, REQ-SEC-002 | OPEN-008/013/015 | defined / planned |
| CAP-EPT-054 | Collection Packaging, Completeness, Integrity and Metadata | Endpoint Agent | Evidence Reviewer, Auditor | package/completeness/integrity metadata | Investigate, Trust | REQ-PROD-014/020, REQ-OBJ-004, REQ-SEC-004 | OPEN-014 | defined / planned |
| CAP-EPT-055 | Collection Transfer, Delivery and Consumer Handoff | Endpoint Agent | DFIR Analyst, Endpoint Operator | transfer/ack/handoff | Investigate, Shared | REQ-PROD-014/018/020, REQ-SEC-001/004 | OPEN-008/014/015 | defined / planned |
| CAP-EPT-056 | Live Response Session Definition, Eligibility and Target Binding | Endpoint Agent | Response Operator, Endpoint Operator | technical session/binding | Investigate, Govern | REQ-PROD-006/014/018/019/020, REQ-SEC-001/002 | OPEN-008/013/015/017 | defined / planned |
| CAP-EPT-057 | Live Response Session Lifecycle and Runtime State | Endpoint Agent | Response Operator, Session Supervisor | technical session state | Investigate | REQ-PROD-014/018/019, REQ-SEC-002 | OPEN-008/013/015 | defined / planned |
| CAP-EPT-058 | Technical Command Request and Invocation Boundary | Endpoint Agent | Response Operator, Govern Reviewer | command request/eligibility | Govern, Studio | REQ-PROD-006/014/018/019/020, REQ-SEC-001/002 | OPEN-008/013/015/017 | defined / planned |
| CAP-EPT-059 | Interactive Command and Shell Execution | Endpoint Agent | Response Operator, Endpoint Operator | execution attempt/runtime context | Investigate, Govern | REQ-PROD-014/018/019, REQ-SEC-001/002 | OPEN-008/013/015/017 | defined / planned |
| CAP-EPT-060 | Script Execution and Runtime Context Boundary | Endpoint Agent | Response Operator, Govern Reviewer | script ref/runtime/attempt | Investigate, Studio, Govern | REQ-PROD-014/018/019, REQ-SEC-001/002/004 | OPEN-008/013/014/015/017 | defined / planned |
| CAP-EPT-061 | Live Response File Read, Download and Transfer Operations | Endpoint Agent | Response Operator, DFIR Analyst | file transfer/read/download/temp upload | Investigate, Govern | REQ-PROD-014/018/020, REQ-SEC-001/004 | OPEN-008/013/014/015 | defined / planned |
| CAP-EPT-062 | Technical Execution Output, Error, Timeout and Cancellation Semantics | Endpoint Agent | Response Operator, Case Analyst | technical output/error/termination knowledge | Investigate, Govern, Studio | REQ-PROD-014/018/019/020, REQ-SEC-001/002/004 | OPEN-013/015 | defined / planned |
| CAP-EPT-063 | Live Response Operator Control, Session Closure and Audit Context | Endpoint Agent | Session Supervisor, Auditor | operator control/closure | Investigate, Govern, Shared | REQ-PROD-014/018/019, REQ-SEC-001/002/004 | OPEN-008/013/015 | defined / planned |
| CAP-EPT-064 | Collection and Live Response Provenance and Cross-Product Execution Contracts | Endpoint Agent | Auditor, DFIR Analyst, Govern Reviewer | provenance chain/handoff refs | Investigate, Govern, Studio, Settings, Shared | REQ-PROD-006/014/018/019/020, REQ-OBJ-004/007, REQ-SEC-001/002/004 | OPEN-008/013/014/015/017 | defined / planned |

## Totals
**18 capabilities / 486 numbered sections / 108 mandatory tables / at least 54 GWT.** Duplicate IDs 0; recycled IDs 0; owner conflicts 0; empty mandatory tables 0. Endpoint cumulative after EPT-4: **64 / 1728 / 384**.

`CAP-EPT-001..046` remain unchanged. `CAP-EPT-065+` is not allocated or reserved by EPT-4. EPT-5/EPT-6 remain NOT STARTED.

## Final verification addendum
Functional/build SHA: `32082487b48779434c9ac730b43efd498009825d`. Remote verification confirms PR #2 open/Draft/unmerged, main/README unchanged, CI/status N/A, exact structure and namespace, OPEN-008/014/015/017 still open, Endpoint Screen IDs 0 and EPT-5/EPT-6 NOT STARTED. After publication and final recheck of the documentation-only verification record, **EPT-4 = PASS AFTER POST-PUBLICATION VERIFICATION — 220/220 PASS, 0 PENDING, 0 FAIL**. Exact final verification-record SHA is recorded in PR #2 to avoid self-reference.