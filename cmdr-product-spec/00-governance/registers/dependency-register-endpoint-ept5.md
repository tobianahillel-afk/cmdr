---
id: dependency-register-endpoint-ept5
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-11
source-of-truth: registry
---
# Dependency Register — Endpoint EPT-5

| Family | Endpoint capabilities | External dependency | Owner | Boundary |
|---|---|---|---|---|
| EPT5-AUTH | 065..081 | Action Request / Approval / Decision / Response Run | Govern | Endpoint never grants authority |
| EPT5-PRE | 065..066 | Agent/capability/platform/health observations | Endpoint EPT-1/2 | observed support only |
| EPT5-POL | 065..081 | Endpoint Policy / Fleet / tenant-environment | Settings | Policy != authority |
| EPT5-PROC | 067 | process identity/context | EPT-3/4 Endpoint | effect only after Govern handoff |
| EPT5-ISO | 068..069 | network observations and policy | Endpoint/Settings | host isolation != bounded block |
| EPT5-FILE | 070..071 | file/collection context | Endpoint + Investigate refs | quarantine/delete != Evidence verdict |
| EPT5-SVC | 072 | service/system observations | Endpoint | support declared, never assumed |
| EPT5-SES | 081 | local session identity/state | Endpoint + Settings identity refs | directory actions external |
| EPT5-OUT | 073 | EPT-4 technical outputs/errors | Endpoint EPT-4 | Technical Outcome != Result |
| EPT5-VER | 074..076 | Govern Verification Plan/Assessment | Govern | technical verification != Govern verification |
| EPT5-REV | 077..078 | Rollback Plan / Response Rollback / Recovery | Govern | technical reversal != response rollback |
| EPT5-HAND | 079..080 | Result / consumer objects | Govern / source products | Endpoint supplies refs/facts only |
| EPT5-STUDIO | 065/079/080 | Human Gate / Tool Call / Automation Run | Studio | Human Gate != Approval; Run identities distinct |
| EPT5-SHARED | 073/076/079/080 | Jobs / Trace / Activity / Recovery / Reporting | Shared | generic mechanism, not response owner |
| EPT5-SEC | all | least privilege / step-up / SoD / masking | Security | no hidden delegation / raw secret |
| EPT5-EPT6 | all | update/resilience/security depth | future Endpoint EPT-6 | strict stop boundary |

No dependency transfers canonical ownership.