---
id: govern-screen-capability-map-gov3
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical-addendum
---
# Screen Capability Map — Govern GOV-3 Addendum

No new Screen ID and no detailed screen rewrite is introduced.

| Existing surface | Primary GOV-3 capabilities | Secondary/linked capabilities | Ownership / boundary |
|---|---|---|---|
| `GOV-AUD-001` Audit Trail | CAP-GOV-034..038 | CAP-GOV-047 | Govern business semantics; Shared Trace/Activity/Search/Reporting/Export remain source-owned |
| `GOV-MET-001` Response Metrics | CAP-GOV-039..047 | CAP-GOV-037/038 | Govern metric meanings; Shared Metrics/Reporting remain source-owned |
| `GOV-DEC-001` Decision Register | CAP-GOV-035/041 | CAP-GOV-034/037/046 | existing GOV-1 screen; audit/metrics projection only |
| `GOV-RUN-001` Runs & Rollback | CAP-GOV-036/042..044 | CAP-GOV-034/037/046 | existing GOV-2 screen; no detailed rewrite |
| `GOV-INB-001` Response Inbox | CAP-GOV-045 | CAP-GOV-039..041 | existing GOV-1 screen; flow projection only |
| `GOV-ACT-001` Action Center | CAP-GOV-035/039..041 | CAP-GOV-047 | source context only |
| Shared Reporting surfaces | CAP-GOV-038/039..047 | — | Shared owns rendering/publication/export |
| Shared Inspector/Trace | CAP-GOV-034..038 | — | Design System/Shared owns component mechanics |

## GOV-3 screen constraints
- new Screen IDs: **0**;
- detailed rewrites: **0**;
- wireframes: **0**;
- final buttons/columns/filters: **0**;
- animations/shortcuts: **0**.

Only capability references, ownership and cross-product links are specified.