# Response Metrics

## Objective

Measure response effectiveness, safety and governance quality rather than only speed.

## Scope

This specification owns the page-local behaviour of **Response Metrics** in **Response & Governance**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

Response Governance Lead

## Affected objects

- Metric definition
- Response request
- Decision
- Run
- Playbook
- Tenant

## Features

- Time to decision and execution
- Approval bottlenecks
- Policy failures and exceptions
- Run success, partial success and rollback
- Business-impact outcomes
- Playbook reliability
- Metric drill-down

## UX and interactions

- Every metric links to its definition
- Medians and distributions are preferred to misleading averages
- Tenant comparison is privacy-safe
- Metric changes are versioned

## Permissions

`command.read` or scoped governance analytics permission; exports require `report.manage`.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Metrics derive from canonical object states and immutable timestamps.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Analytics pipeline, decisions, runs, policy, reporting.

## Acceptance criteria

- Metrics reconcile to source objects
- Definition/version is displayed
- Partial success is counted explicitly
- Exceptions are distinguishable from policy failures
- Small cohorts are privacy-protected

## Open questions

- Which metrics are contractual SLAs?
- How are business outcomes attributed to response actions?
