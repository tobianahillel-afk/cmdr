# Audit Trail

## Objective

Provide an immutable, searchable record of security-relevant activity across CMDR.

## Scope

This specification owns the page-local behaviour of **Audit Trail** in **Response & Governance**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

Response Governance Lead

## Affected objects

- Audit event
- Principal
- Tenant
- Object reference
- Decision
- Run
- Export

## Features

- Global and object-scoped search
- Actor, action, target and outcome
- Correlation chains
- Policy and permission decisions
- Export and integrity verification
- Retention and legal hold

## UX and interactions

- Audit events are readable without exposing forbidden object content
- Correlation chains navigate across consoles
- Integrity status is visible
- Exports are scoped and watermarked or signed as required

## Permissions

`audit.read`; export requires `audit.export`; legal-hold administration is separately privileged.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Audit events are immutable; retention disposition is metadata, not mutation.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Audit store, identity, all mutating services, retention.

## Acceptance criteria

- Every privileged mutation emits an audit event
- Denied attempts are recorded where policy requires
- Events are time-synchronised
- Integrity verification is reproducible
- Tenant isolation applies to search and export

## Open questions

- Which signing/timestamping standard is required?
- What audit events are customer-visible in MSSP mode?
