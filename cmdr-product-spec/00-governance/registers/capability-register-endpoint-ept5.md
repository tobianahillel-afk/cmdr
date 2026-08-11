---
id: capability-register-endpoint-ept5
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-11
source-of-truth: registry
---
# Capability Register — Endpoint EPT-5

Execution lot: **EPT-5 — Containment, Verification and Governed Response Primitives**. All rows are `draft / defined / planned`; documentary definition proves no implementation, supported platform or response authority.

| ID | Capability | Primary role | Primary concepts | Key consumers | OPEN | Delivery |
|---|---|---|---|---|---|---|
| CAP-EPT-065 | Governed Response Primitive Request, Eligibility and Authority Boundary | Response Operator | Response Primitive, Eligibility, Authority Ref | Govern, EPT-5 | 007/008/013/015 | defined/planned |
| CAP-EPT-066 | Response Primitive Preconditions, Readiness and Target-State Planning | Endpoint Operator | Technical Precheck, Readiness, Requested Target State | EPT-5, Govern | 008/013/015 | defined/planned |
| CAP-EPT-067 | Process Control, Suspension, Resume and Termination Primitives | Response Operator | Process Control Execution/State | Govern, Verification | 008/013/015 | defined/planned |
| CAP-EPT-068 | Host and Network Isolation Primitive | Response Operator | Isolation State | Govern, Verification | 008/013/015 | defined/planned |
| CAP-EPT-069 | Network Blocking and Connection Control Primitive | Response Operator | Network Control State | Govern, Verification | 008/013/015 | defined/planned |
| CAP-EPT-070 | File Quarantine and Restricted File-State Primitive | Response Operator | File Quarantine State | Govern, Investigate | 008/013/014/015 | defined/planned |
| CAP-EPT-071 | File Delete, Restore and Recovery Boundary | Response Operator | Delete/Restore State | Govern, Investigate | 008/013/014/015 | defined/planned |
| CAP-EPT-072 | Service and System-Component Control Primitive | Response Operator | Service Control State | Govern, Settings | 008/013/015 | defined/planned |
| CAP-EPT-073 | Response Primitive Execution State and Technical Outcome | Endpoint Operator | Execution State, Technical Outcome | Govern | 008/013/015 | defined/planned |
| CAP-EPT-074 | Technical Verification Request, Scope and Verification Criteria | Verification Reviewer | Technical Verification Request | CAP-EPT-075, Govern | 008/013/015 | defined/planned |
| CAP-EPT-075 | Endpoint Target-State Observation and Technical Effect Verification | Verification Reviewer | Target State Observation, Technical Verification Observation | Govern | 008/013/015 | defined/planned |
| CAP-EPT-076 | Partial Success, Failure, Unknown State and Drift Semantics | Verification Reviewer | Uncertainty State, Response Drift | Govern | 008/013/015 | defined/planned |
| CAP-EPT-077 | Technical Reversal, Compensation and Local Rollback Primitive | Rollback Operator | Technical Reversal | Govern Rollback | 007/008/013/015 | defined/planned |
| CAP-EPT-078 | Containment Release and Restore-to-Service Boundary | Response Operator | Release State, Restore-to-Service State | Govern | 008/013/015 | defined/planned |
| CAP-EPT-079 | Govern Reconciliation, Result Input and Consumer Handoff | Govern Reviewer | Reconciliation Handoff | Govern, Command, Investigate | 007/008/013/014/015/017 | defined/planned |
| CAP-EPT-080 | Governed Response Primitive Provenance and Cross-Product Contracts | Auditor | EPT-5 Provenance Chain | Govern, Studio, Settings, Shared | 007/008/013/014/015/017 | defined/planned |
| CAP-EPT-081 | Local Account Session Lock and Termination Primitive | Response Operator | Local Session Control State | Govern, Settings identity | 008/013/015 | defined/planned |

## Totals
**17 capabilities / 459 numbered sections / 102 mandatory tables / at least 51 GWT.** Duplicate IDs 0; recycled IDs 0; owner conflicts 0; empty/generic mandatory tables 0.

Endpoint cumulative after EPT-5: **81 capabilities / 2187 sections / 486 mandatory tables**. Global content after EPT-5: **466 capabilities / 464 defined / 2 proposed / 466 planned / 12582 sections / 2796 mandatory tables**.

`CAP-EPT-001..064` remain unchanged. `CAP-EPT-082+` is not allocated or reserved by EPT-5. EPT-6 remains **NOT STARTED**. OPEN remains 18.