# Data Retention and Export

## Objective

Apply retention, legal hold, deletion and export rules by object class and tenant.

## Scope

This specification owns the page-local behaviour of **Data Retention and Export** in **Shared Platform**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

Platform Product Lead

## Affected objects

- Retention policy
- Data object
- Legal hold
- Export job
- Tenant

## Features

- Policy by object/data class
- Legal hold
- Deletion scheduling and evidence
- Portable tenant export
- Evidence/report export controls
- Residency constraints

## UX and interactions

- Users see retention consequences before destructive actions
- Exports show scope and estimated size
- Legal-hold conflicts block deletion with explanation
- Progress and partial failure are visible

## Permissions

`tenant.manage` or designated retention permission; actual exports still require object-specific export permissions.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Export jobs queued/running/completed/failed/expired; retention disposition scheduled/held/disposed.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Storage, audit, tenancy, key management, reporting.

## Acceptance criteria

- Policy is evaluated per object class
- Legal hold always overrides deletion
- Exports cannot broaden access
- Deletion produces auditable evidence
- Expired exports are inaccessible

## Open questions

- What default retention applies per object type?
- Which export formats satisfy portability requirements?
