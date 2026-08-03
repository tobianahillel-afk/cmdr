# Runs and Rollback

## Objective

Supervise execution, step outcomes, partial failures and rollback.

## Scope

This specification owns the page-local behaviour of **Runs and Rollback** in **Response & Governance**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

Response Governance Lead

## Affected objects

- Run
- Step
- Decision
- Playbook
- Agent
- Rollback
- Entity

## Features

- Run queue and live status
- Step timeline and logs
- Target and scope verification
- Pause/cancel controls where safe
- Rollback eligibility and execution
- Partial-success resolution
- Return operational outcome

## UX and interactions

- Live updates do not reset investigation context
- Partial success is never shown as success
- Rollback explains what will and will not revert
- Logs redact secrets
- Manual intervention is structured and audited

## Permissions

`response.execute`; rollback requires `response.rollback`; viewing may require `response.approve` or scoped read permission.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Uses canonical Run states exactly.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Execution engine, agents, connectors, audit, incident updates.

## Acceptance criteria

- Every step is attributable and ordered
- Target scope matches approved decision
- Timeouts produce explicit results
- Rollback status is independent and visible
- Final outcome returns to Command Center

## Open questions

- Which actions support pause/resume?
- How are orphaned runs recovered?
