---
id: capability-register-endpoint-ept6
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-12
source-of-truth: registry
---
# Capability Register — Endpoint EPT-6

Execution lot: **EPT-6 — Updates, Resilience, Security and Endpoint Provenance**. All rows are `draft / defined / planned`; none claims implementation, supported-platform delivery or production readiness.

| ID | Capability | Primary role | Primary concepts | Key consumers | OPEN | Delivery |
|---|---|---|---|---|---|---|
| CAP-EPT-082 | Administrative Update Assignment, Channel and Wave Boundary | Endpoint Operator | Local Update Assignment Projection | Settings, Update | 008/013 | defined/planned |
| CAP-EPT-083 | Update Package, Release Reference and Compatibility Context | Release Reviewer | Package Context, Compatibility | Update, Settings | 008 | defined/planned |
| CAP-EPT-084 | Update Download, Staging and Readiness | Endpoint Operator | Download, Staging, Readiness | Update, Settings | 008/013 | defined/planned |
| CAP-EPT-085 | Update Installation, Activation and Local Lifecycle | Endpoint Operator | Update Operation, Install/Activation | Settings, Audit | 008/013 | defined/planned |
| CAP-EPT-086 | Update Progress, Deferral, Failure and Retry Semantics | Endpoint Operator | Attempt, Progress, Retry Eligibility | Update, Shared | 008/013 | defined/planned |
| CAP-EPT-087 | Post-Update Health, Compatibility and Technical Verification | Release Reviewer | Post-Update Verification | Settings, Quality | 008 | defined/planned |
| CAP-EPT-088 | Update Reversion and Previous-Version Recovery Boundary | Endpoint Operator | Update Reversion, Previous Version | Settings, Govern, Studio | 008/013/015 | defined/planned |
| CAP-EPT-089 | Offline Buffering, Queueing and Local Retention Semantics | Endpoint Operator | Buffered Item, Local Queue | Shared, Investigate | 008 | defined/planned |
| CAP-EPT-090 | Reconnect, Replay, Resumption and Duplicate-Control Boundary | Endpoint Operator | Replay, Resumption, Duplicate Candidate | Shared, consumers | 008 | defined/planned |
| CAP-EPT-091 | Local State Persistence, Restart and Crash Recovery | Support Operator | Restart/Crash Recovery State | Settings, Shared | 008 | defined/planned |
| CAP-EPT-092 | Resource Pressure, Backpressure and Degraded-Operation Semantics | Reliability Reviewer | Resource Pressure, Degraded Operation | Settings, consumers | 008 | defined/planned |
| CAP-EPT-093 | Dependency Failure, Capability Reassessment and Recovery | Endpoint Operator | Dependency State, Capability Reassessment | Settings, Shared | 008 | defined/planned |
| CAP-EPT-094 | Endpoint Self-Protection and Anti-Tamper Boundary | Security Reviewer | Self-Protection, Tamper Candidate | Security, Investigate | 008/013 | defined/planned |
| CAP-EPT-095 | Local Privilege, Service/Process Protection and Security Context | Security Reviewer | Runtime Privilege, Protected Component | Security, Settings | 008/013 | defined/planned |
| CAP-EPT-096 | Sensitive Material, Secret Reference and Credential Handling Boundary | Security Reviewer | Secret Reference, Handling State | Settings, Security | 008 | defined/planned |
| CAP-EPT-097 | Local Audit Event, Security Audit and Provenance Semantics | Auditor | Local Audit Event, Provenance | Security, Shared, Govern | 008 | defined/planned |
| CAP-EPT-098 | Endpoint Security-State, Integrity and Consumer Handoff | Security Reviewer | Security-State Projection, Handoff | Command, Investigate, Settings, Govern | 008/015/017 | defined/planned |
| CAP-EPT-099 | Endpoint Updates, Resilience, Security and Capability Closure | Product Architecture | Endpoint Closure State, Provenance Chain | Roadmap, Quality, all owners | 007/008/013/014/015/017 | defined/planned |

## Totals
**18 capabilities / 486 numbered sections / 108 mandatory tables / at least 54 GWT.**

Endpoint cumulative after EPT-6 content: **99 capabilities / 2673 sections / 594 mandatory tables**.

Global content after EPT-6: **484 capabilities / 482 defined / 2 proposed / 484 planned / 13068 sections / 2904 mandatory tables**.

Requirements baseline remains **122 = 99 conform / 20 partial / 3 absent / 0 contradictory** unless the canonical global matrix records a separately verified change. OPEN remains **18**. Endpoint Screen IDs remain **0**.

`CAP-EPT-001..081` are unchanged. `CAP-EPT-082..099` are the source-driven EPT-6 allocation. No `CAP-EPT-100+` is allocated or reserved by EPT-6.