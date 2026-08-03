# Capabilities

## Objective

Make data coverage, detection, collection, response and analytical capability gaps visible to operators.

## Scope

This specification owns the page-local behaviour of **Capabilities** in **Command Center**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

SOC Operations Lead

## Affected objects

- Integration
- Agent
- Data source
- Detection capability
- Response capability
- Tenant

## Features

- Coverage matrix by tenant and environment
- Data-source health and parsing quality
- Detection and response capability inventory
- Gap and degradation tracking
- Ownership and remediation links
- Readiness for defined incident scenarios

## UX and interactions

- Capabilities are grouped by outcome, not vendor
- Drill-down shows evidence behind coverage
- Stale health is explicit
- Users can compare current and required capability

## Permissions

`command.read`; configuration links require `integration.manage` or `agent.manage`.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Agent state is canonical; integration health uses platform-defined operational status.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Integrations, agent fleet, parser health, detection catalogue, policy catalogue.

## Acceptance criteria

- Coverage claims link to measurable evidence
- Degraded sources update incident confidence where applicable
- No vendor-specific duplicates create false coverage
- Historical changes are available

## Open questions

- How is detection quality separated from simple data presence?
- What constitutes scenario readiness?
