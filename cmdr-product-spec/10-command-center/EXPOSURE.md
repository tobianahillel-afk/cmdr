# Exposure

## Objective

Show exploitable attack surface and business-relevant exposure without turning the page into a vulnerability scanner clone.

## Scope

This specification owns the page-local behaviour of **Exposure** in **Command Center**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

SOC Operations Lead

## Affected objects

- Entity
- Asset
- Business service
- Exposure
- Incident
- Finding

## Features

- Exposure overview by service and tenant
- Critical assets and attack paths
- Known vulnerabilities and misconfigurations
- Active exploitation context
- Compensating controls
- Create operational work from confirmed exposure

## UX and interactions

- Users can switch technical and business lenses
- Attack paths offer accessible tabular alternatives
- Risk components are explainable
- Creating work requires review of scope and duplicates

## Permissions

`command.read`; creating incidents/tasks requires `incident.manage`; sensitive asset detail is ABAC-scoped.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Exposure status is separate from Incident state; remediation tracking must not reuse response Run states.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Asset inventory, vulnerability sources, business catalogue, identity graph.

## Acceptance criteria

- Exposure records show source freshness
- Risk scores expose components
- Duplicate work creation is prevented
- Tenant isolation holds in graph queries

## Open questions

- Which exposure scoring model is canonical?
- How are accepted risks represented?
