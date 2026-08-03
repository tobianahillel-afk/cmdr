# Decision Register

## Objective

Maintain the authoritative record of response decisions and their rationale.

## Scope

This specification owns the page-local behaviour of **Decision Register** in **Response & Governance**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

Response Governance Lead

## Affected objects

- Decision
- Response request
- Principal
- Condition
- Policy
- Run

## Features

- Searchable decision ledger
- Decision document with numbered sections
- Rationale and authority
- Conditions and expiry
- Linked run and outcome
- Supersession history
- Export for audit

## UX and interactions

- Published decisions are immutable
- Corrections create superseding records
- Users can compare versions
- Sensitive rationale is access-controlled
- The register is usable as a formal record, not merely a table

## Permissions

`audit.read` or `response.approve` to view according to scope; export requires `audit.export`.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Uses canonical Decision states; published records may be superseded but not overwritten.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Audit store, policy, identity, reporting.

## Acceptance criteria

- Every decision has actor, authority and policy snapshot
- Supersession chains are complete
- Exports are tamper-evident
- Search cannot reveal inaccessible decision metadata
- Linked runs are resolvable

## Open questions

- What electronic-signature standard is required?
- How long must decisions be retained?
