# Event Search

## Objective

Search and correlate SIEM/telemetry data with precise field, time and tenant control.

## Scope

This specification owns the page-local behaviour of **Event Search** in **Investigation Lab**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

DFIR Lead

## Affected objects

- Event
- Alert
- Entity
- Case
- Saved query
- Search job

## Features

- Structured and free-text query modes
- Field browser and schema hints
- Time-range and tenant filters
- Aggregation and correlation
- Saved queries and history
- Add selected results to case as evidence references
- Raw and normalised event views

## UX and interactions

- Queries show execution state and cost guardrails
- Large results stream or paginate safely
- Adding to case stores immutable references and query context
- Raw data is escaped and copyable
- Search errors identify invalid fields and correlation IDs

## Permissions

`case.read` for case context; telemetry access is ABAC-scoped; adding evidence requires `evidence.collect`; exports require `evidence.export`.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Search job states are queued/running/completed/failed/cancelled; evidence created from results follows canonical Evidence states.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Search/index platform, parsers, schema registry, evidence service.

## Acceptance criteria

- Tenant filters are enforced server-side
- Saved queries are versioned
- Results expose ingestion and event timestamps
- Case references remain resolvable after index rollover
- Export matches query and time scope

## Open questions

- Which query language is canonical?
- How are expensive queries limited without blocking DFIR?
