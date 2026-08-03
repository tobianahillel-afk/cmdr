# Triage Desk

## Objective

Intake new cases, assess urgency, identify missing context and assign the correct investigative path.

## Scope

This specification owns the page-local behaviour of **Triage Desk** in **Investigation Lab**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

DFIR Lead

## Affected objects

- Case
- Incident
- Alert
- Evidence
- Principal
- Triage note

## Features

- New and pending case queues
- Intake completeness check
- Priority and owner assignment
- Evidence preview
- Quick actions for case creation, import and search
- Triage playbook guidance

## UX and interactions

- Notebook-style notes remain structured and searchable
- Previewing evidence never changes integrity state
- Triage decisions show rationale
- Queue position survives panel navigation

## Permissions

`case.read` and `case.create`; assignment requires `case.manage`; evidence import requires `evidence.collect`.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Case moves through intake and triage states only when required metadata and ownership exist.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Incident links, evidence registry, identity, notifications.

## Acceptance criteria

- Duplicate case candidates are surfaced
- Every triage outcome records actor and rationale
- Missing required evidence is explicit
- Tenant context cannot change during intake

## Open questions

- What mandatory fields vary by case type?
- Which triage playbooks ship by default?
