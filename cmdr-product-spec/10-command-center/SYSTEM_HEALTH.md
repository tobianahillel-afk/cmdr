# System Health

## Objective

Show whether CMDR can be trusted operationally: ingestion, search, agents, integrations, queues and critical dependencies.

## Scope

This specification owns the page-local behaviour of **System Health** in **Command Center**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

SOC Operations Lead

## Affected objects

- Integration
- Agent
- Data pipeline
- Service
- Health signal
- Incident

## Features

- Current health and SLO view
- Ingestion lag and parsing failures
- Search/index freshness
- Agent connectivity
- Dependency incidents
- Maintenance and degradation banners

## UX and interactions

- Health is separated from security severity
- Affected product functions are explained
- Users can subscribe to material degradation
- Historical charts expose gaps without visual noise

## Permissions

`command.read`; detailed administrative diagnostics require `platform.admin`, `agent.manage` or `integration.manage`.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Uses canonical Agent states; service health uses shared platform operational statuses.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Observability, integrations, agent fleet, notification service.

## Acceptance criteria

- Freshness calculations are documented
- Health banners appear on affected pages
- Partial outages identify lost capabilities
- No sensitive infrastructure detail leaks to unauthorised users

## Open questions

- Which SLOs are customer-visible?
- How are third-party outages represented?
