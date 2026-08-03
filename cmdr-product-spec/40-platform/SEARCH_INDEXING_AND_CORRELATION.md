# Search, Indexing and Correlation

## Objective

Provide low-latency, tenant-safe search and relationship discovery across telemetry and product objects.

## Scope

This specification owns the page-local behaviour of **Search, Indexing and Correlation** in **Shared Platform**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

Platform Product Lead

## Affected objects

- Event
- Search document
- Index
- Saved query
- Correlation job
- Entity

## Features

- Normalised schema and raw access
- Field capabilities
- Index lifecycle
- Cross-object search
- Correlation jobs
- Saved queries
- Freshness and completeness indicators

## UX and interactions

- Search always shows time and tenant scope
- Index lag is visible
- Partial results are explicit
- Query history protects sensitive literals
- Pagination is stable

## Permissions

Object read permissions and ABAC are applied before results; export is separately authorised.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Search jobs queued/running/completed/partial/failed/cancelled.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Telemetry ingestion, object services, entity resolution, audit.

## Acceptance criteria

- Tenant filters are enforced at storage/query layer
- Freshness is measurable
- Partial shards cannot appear complete
- Saved queries are versioned
- Result counts reconcile within documented consistency windows

## Open questions

- What search backend and query dialect are selected?
- Which objects support full-text versus structured search?
