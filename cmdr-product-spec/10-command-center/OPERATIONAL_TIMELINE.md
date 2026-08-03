# Operational Timeline

## Objective

Provide a shared chronological view of incidents, decisions, actions, outages and major security events.

## Scope

This specification owns the page-local behaviour of **Operational Timeline** in **Command Center**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

SOC Operations Lead

## Affected objects

- Timeline entry
- Incident
- Decision
- Run
- Agent
- Integration
- Principal

## Features

- Global and tenant timelines
- Event-type filters
- Observed versus inferred markers
- Decision and run milestones
- Annotations and shift handovers
- Exportable review windows

## UX and interactions

- Zoom from months to seconds where supported
- Dense periods cluster without hiding counts
- Selecting an item opens provenance
- Annotations never alter source events
- Timezone is explicit and switchable

## Permissions

`command.read`; annotation requires `command.coordinate`; export requires `report.manage` or `audit.export` as applicable.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Timeline entries are immutable references; source objects retain their canonical states.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Audit, incidents, runs, platform health, time normalisation.

## Acceptance criteria

- Ordering is deterministic for equal timestamps
- Every item exposes source and ingestion time
- Filters survive sharing
- Export matches the visible scope

## Open questions

- Should executive events have a separate curated layer?
- How long are high-resolution timeline indexes retained?
