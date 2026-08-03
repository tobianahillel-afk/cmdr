# Playbooks

## Objective

Create and govern versioned response procedures and atomic actions.

## Scope

This specification owns the page-local behaviour of **Playbooks** in **Response & Governance**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

Response Governance Lead

## Affected objects

- Playbook
- Step
- Action connector
- Policy
- Rollback definition
- Version

## Features

- Visual and code/config views
- Versioning and review
- Input and output contracts
- Preconditions and policy requirements
- Rollback steps
- Dry run and test evidence
- Deprecation and compatibility

## UX and interactions

- Editor validates before publication
- Production and draft versions are unmistakable
- Step failures and branching are explicit
- Secrets are referenced, never embedded
- Changes show diff and affected approvals

## Permissions

`playbook.manage`; testing may require `response.execute` in an isolated scope; viewing is separately scoped.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Playbook states: draft, review, approved, published, deprecated. Runs use canonical Run states.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Connector catalogue, secrets manager, policy engine, agents, audit.

## Acceptance criteria

- Published versions are immutable
- Every action has timeout and failure policy
- Rollback is defined or explicitly unavailable
- Test evidence is attached
- Existing runs retain original version

## Open questions

- Which playbook language is canonical?
- What review is required for high-impact connectors?
