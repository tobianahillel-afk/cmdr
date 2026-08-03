# Audit and Immutability

## Objective

Capture tamper-evident, queryable history for security-relevant activity.

## Scope

This specification owns the page-local behaviour of **Audit and Immutability** in **Shared Platform**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

Platform Product Lead

## Affected objects

- Audit event
- Correlation identifier
- Principal
- Object reference
- Integrity proof

## Features

- Central event contract
- Append-only storage
- Correlation across services
- Integrity verification
- Retention and legal hold
- Scoped export
- Clock-health monitoring

## UX and interactions

- Product pages link to relevant audit history
- Integrity failures are prominent
- Audit search remains usable at scale
- Redaction never destroys the original protected record

## Permissions

`audit.read` and `audit.export`; producers cannot suppress mandatory events.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Audit events are immutable; storage partitions may be active/sealed/retained/disposed under policy.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

All services, time sync, storage, key management.

## Acceptance criteria

- Mandatory mutation events are complete
- Integrity chain verification is repeatable
- Actor and tenant are always present
- Failed writes trigger operational alerts
- Retention actions are audited

## Open questions

- What tamper-evidence mechanism is required?
- Which events require external timestamping?
