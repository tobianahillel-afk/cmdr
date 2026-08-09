---
id: govern-source-migration-gov3
domain: 08-govern
status: draft
owner: Product Architecture
updated: 2026-08-09
source-of-truth: migration-audit
---
# Govern GOV-3 Source Migration Audit

## Audited historical responsibilities
Audit/audit ledger/governance history, compliance reporting, approval/Decision/Run/response/rollback/effectiveness metrics, KPI/SLA/dashboard and continuous-improvement references were reviewed against the current product boundaries.

## Canonical disposition
- Govern Audit Trail owns Govern-domain audit interpretation/reconstruction/review semantics.
- Shared Trace/Activity retain generic provenance/activity infrastructure.
- Shared Reporting/Export retain Report and export mechanics.
- Shared Metrics retains generic metric definition/version/reconciliation/privacy/freshness mechanics.
- Govern Response Metrics owns Govern-specific metric definitions/interpretation only.
- Platform Settings retains retention/storage/export-destination administration.
- Security retains permission/privacy/integrity/legal-hold policy.
- Command, Investigate, Studio and Endpoint retain their own operational/investigation/automation/technical metrics.

## Existing Govern surfaces
`GOV-AUD-001` and `GOV-MET-001` remain active. No legacy screen is replaced or deprecated by GOV-3. Their generic module descriptions are aligned to CAP-GOV-034..047; detailed screen specifications remain future work.

## No competing architecture
No second Trace, Activity, Reporting, Export, Metrics, SIEM, audit store, warehouse or retention architecture is created. Existing generic historical `À compléter` content in Audit Trail/Response Metrics is superseded by the GOV-3 module/capability semantics, not by a competing module.

## Explicit non-deprecations
GOV-3 does not deprecate Shared Trace, Shared Activity, Shared Reporting, Shared Export, Command metrics, Investigate metrics, Studio metrics or Endpoint technical metrics.

## Implementation boundary
No audit engine, metrics engine, storage schema, event format, API, protocol, query language, dashboard implementation, KPI/SLO threshold or external compliance claim is selected.