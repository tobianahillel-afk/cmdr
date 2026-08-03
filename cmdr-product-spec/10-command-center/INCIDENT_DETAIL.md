# Incident Detail

## Objective

Provide one operational record for a security situation, its impact, timeline, affected assets, ownership and next decisions.

## Scope

This specification owns the page-local behaviour of **Incident Detail** in **Command Center**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

SOC Operations Lead

## Affected objects

- Incident
- Alert
- Entity
- Timeline entry
- Case
- Finding
- Response request
- Decision
- Run

## Features

- Situation summary and business impact
- Attack/operational timeline
- Affected entities and services
- Linked alerts and detections
- Team activity and ownership
- Recommended next actions
- Open Investigation and Response links

## UX and interactions

- Timeline zoom and provenance inspection
- Inline assignment and status changes
- Related objects open in panels before full navigation
- Unsaved operational notes are protected
- Critical actions route to governance rather than executing locally

## Permissions

`command.read` to view; `incident.manage` to edit; `case.create` to open a case; `response.request` to submit a governed action.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Uses canonical Incident state; linked Case, Decision and Run states display separately.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Event correlation, case service, business impact, response summaries, audit.

## Acceptance criteria

- Incident history is complete and ordered
- Status cannot become contained solely from an approval; execution evidence is required
- Linked case and response objects are bidirectional
- Every mutation records actor and rationale
- Refresh preserves active tab and timeline position

## Open questions

- Can incidents aggregate across tenants in MSSP views?
- Which fields are mandatory before closure?
