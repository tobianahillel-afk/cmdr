---
id: govern-cross-product-links-gov3
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical-addendum
---
# Govern Cross-Product Links — GOV-3 Addendum

GOV-3 consumes and returns permission-aware projections; it never transfers object ownership.

| Source | GOV-3 use | Destination / feedback | Boundary |
|---|---|---|---|
| Govern GOV-1 Action Request/Policy/Approval/Decision | audit reconstruction + metrics | GOV-1 owner/reviewer | audit/metric does not mutate source |
| Govern GOV-2 Run/Verification/Rollback/Result | reconstruction + reliability/outcome metrics | GOV-2 owner/reviewer | raw technical output remains distinct from Result |
| Command Incident/operational context | outcome/follow-up projection | Command | Command owns Incident/KPIs; GOV-3 only links/feeds improvement context |
| Investigate Case/Evidence/Finding | source/follow-up context | Investigate | Audit Evidence Package ≠ Evidence; Investigate decides qualification |
| Studio Workflow/Automation Run/Tool Call | provenance and contextual metric projection | Studio | Workflow/Automation metrics remain Studio-owned |
| Endpoint technical runtime | raw technical refs/technical metrics projection | Endpoint | technical metrics remain Endpoint-owned |
| Platform Settings | tenant/environment/retention/storage/access context | Settings | GOV-3 cannot administer retention/storage/access |
| Shared Trace/Activity | correlation/presentation | Shared | Audit Trail ≠ Shared Trace/Activity |
| Shared Metrics/Reporting/Export | generic calculation/render/export | Shared | Govern owns meaning only; export cannot widen visibility |
| Security | permission/privacy/integrity/legal-hold policy | Security | GOV-3 cannot claim cryptographic proof or final RBAC |
| CAP-GOV-047 Improvement Package | proposal | appropriate owner | package applies no change; destination owner decides |

## Handoff rules
- Audit Review → Investigate: package/provenance only; no automatic Evidence.
- Audit/Metric → Reporting/Export: exact snapshot/classification/masking; external release remains separately authorized.
- Govern Result/metric → Command: operational context only; no silent Incident mutation.
- Govern finding candidate → Studio/Settings/Endpoint: improvement proposal only.
- GOV-3 closure → roadmap: status evidence only; no Phase 5 capability creation.