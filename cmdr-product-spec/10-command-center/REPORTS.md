# Command Reports

## Objective

Produce operational, leadership and customer-facing reports from governed data without creating parallel facts.

## Scope

This specification owns the page-local behaviour of **Command Reports** in **Command Center**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

SOC Operations Lead

## Affected objects

- Report
- Incident
- Metric definition
- Business service
- Tenant

## Features

- Scheduled and on-demand reports
- Operational and executive templates
- Metric definition links
- Tenant branding and scope
- Versioning and approval workflow
- Export history

## UX and interactions

- Preview matches final export
- Users can trace every figure to filtered data
- Reports freeze a data snapshot
- Sensitive sections can be permission-gated
- Scheduled delivery shows recipients and next run

## Permissions

`report.manage` to create or schedule; `command.read` to consume; exports additionally enforce tenant and data-classification policy.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Report draft/review/published/archived lifecycle is platform-defined and must not replace Incident states.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Reporting service, metrics catalogue, notifications, audit.

## Acceptance criteria

- Generated report records source snapshot and template version
- Recipient scope is validated at send time
- Failed delivery is visible and retryable
- Numbers reconcile with dashboard definitions

## Open questions

- Which report types require formal approval?
- What formats are required at launch?
