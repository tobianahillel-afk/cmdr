# Policy Gates

## Objective

Define and evaluate conditions that must hold before response actions proceed.

## Scope

This specification owns the page-local behaviour of **Policy Gates** in **Response & Governance**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

Response Governance Lead

## Affected objects

- Policy
- Gate
- Response request
- Decision
- Condition
- Business service

## Features

- Gate catalogue and versions
- Rule authoring and simulation
- Inputs and evaluation explanation
- Required approvals by action/impact
- Exceptions and emergency paths
- Policy test cases

## UX and interactions

- Evaluation explains pass, fail and unknown
- Policy changes show affected actions
- Simulation cannot be mistaken for enforcement
- Emergency override requires explicit rationale and authority

## Permissions

`policy.manage`; gate evaluation is service-controlled; exception approval uses `response.approve` with designated authority.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Policy states: draft, review, active, retired. Gate results are pass/fail/unknown/not_applicable.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Authority model, business catalogue, response actions, audit.

## Acceptance criteria

- Active policy versions are immutable
- Every decision stores evaluated policy snapshot
- Unknown does not silently pass
- Exceptions are time-bound and audited
- Tests cover critical action classes

## Open questions

- Which policy engine or expression language is used?
- How are jurisdiction-specific policies layered?
