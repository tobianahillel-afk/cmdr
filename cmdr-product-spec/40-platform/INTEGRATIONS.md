# Integrations

## Objective

Connect CMDR to telemetry, identity, ticketing, threat intel, EDR, cloud and response systems through explicit capability contracts.

## Scope

This specification owns the page-local behaviour of **Integrations** in **Shared Platform**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

Platform Product Lead

## Affected objects

- Integration
- Connector
- Credential reference
- Capability
- Health signal
- Tenant

## Features

- Connector catalogue
- Configuration and validation
- Capability declaration
- Credential rotation
- Health and rate-limit visibility
- Schema/version compatibility
- Test and disable workflows

## UX and interactions

- Secrets are never displayed after entry
- Test results distinguish auth, connectivity and capability failures
- Disabling explains affected workflows
- Vendor branding does not replace capability language

## Permissions

`integration.manage`; usage additionally requires the domain permission for the requested operation.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Integrations are draft/validating/active/degraded/disabled/error.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Secrets manager, agents, networking, audit, capability catalogue.

## Acceptance criteria

- Credentials are referenced securely
- Health reflects real checks
- Capabilities are versioned
- Disabling stops new work safely
- Tenant assignment is enforced

## Open questions

- Which connectors are launch scope?
- How are connector SDKs versioned?
