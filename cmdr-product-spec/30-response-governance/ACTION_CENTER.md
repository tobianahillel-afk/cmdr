# Action Center

## Objective

Review action details, impact, conditions, rollback and evidence before a decision.

## Scope

This specification owns the page-local behaviour of **Action Center** in **Response & Governance**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

Response Governance Lead

## Affected objects

- Response request
- Finding
- Evidence package
- Decision
- Business service
- Playbook

## Features

- Situation and established finding
- Proposed action and targets
- Risk of action and inaction
- Business impact preview
- Conditions and rollback plan
- Policy and authority summary
- Decision controls

## UX and interactions

- Decision controls remain fixed and visible
- Justification is mandatory for all outcomes
- Evidence opens read-only in Investigation context
- Impact assumptions are editable only by authorised owners
- A dry-run preview is clearly labelled

## Permissions

`response.approve`; editing the request requires `response.request`; executing after approval requires `response.execute`.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Decision outcomes use canonical states and may create a Run only after conditions are met.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Policy engine, business catalogue, evidence package, playbooks, audit.

## Acceptance criteria

- Decision cannot be submitted without rationale
- Displayed targets exactly match execution scope
- Conditions are machine-checkable where possible
- Rollback availability is verified
- Approval does not imply successful execution

## Open questions

- Which actions require IT-owner approval?
- Can decision authority be delegated temporarily?
