---
id: endpoint-ept5-action-classification
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-11
source-of-truth: canonical
---
# Endpoint EPT-5 — Action Classification

| Action family | Default class | Endpoint authority | Govern dependency | Notes |
|---|---:|---|---|---|
| inspect primitive/state/verification/provenance | 0 | read only | none beyond source access | no effect |
| eligibility/precheck/readiness/technical verification | 1 | deterministic local assessment | no effect authority | observed facts only |
| handoff/reference preparation | 2 | no target effect | OPEN-013 remains open | no authority creation |
| process suspend/resume/terminate | 3 | none autonomously | required | exact target identity |
| host isolation/release | 3 | none autonomously | required | significant effect |
| bounded network block/unblock | 3 by default | none autonomously | required | OPEN-013 prevents arbitrary downgrade |
| quarantine/release | 3 | none autonomously | required | quarantined != malicious |
| file delete | 3/4 | none autonomously | required + irreversibility review | provenance retained |
| file restore/recovery | 3 by default | none autonomously | required | restored != safe |
| service/system control | 3 | none autonomously | required | platform-dependent |
| local session lock/termination | 3 | none autonomously | required | directory actions external |
| technical reversal/release | 3/4 according to effect | none autonomously | Govern Rollback/recovery authority | technical reversal != Govern rollback |

Class 2 is used only for genuinely no-effect or bounded governance preparation where existing rules permit it. OPEN-013 remains open. Class 4 grants no Endpoint autonomous authority.