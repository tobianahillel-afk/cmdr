# Event Search runtime boundary

This directory contains the first bounded executable runtime for **CAP-INV-002 Event Search**.

## Implemented in this slice

The runtime validates only the backend-neutral execution envelope required before a later Search Job orchestration layer can run:

- exactly one selected Tenant, with no wildcard or cross-tenant fallback;
- required environment reference;
- ordered RFC3339 time range;
- one or more explicitly authorized sources;
- explicit query presence/validity signal;
- server-evaluated execution permissions;
- opaque correlation and query-version provenance;
- explicit refusal of Case-link mutation while OPEN-013 remains unresolved.

Denied validation returns no populated execution envelope and never widens protected-data visibility.

## Deliberately not implemented

This module does **not** define or select:

- a query language or dialect;
- a search/index engine;
- a storage engine;
- a backend/provider;
- persistence;
- a final production search/index/storage/provider implementation;
- Case-link mutation.

Those concerns remain owned by later bounded tasks and unresolved product decisions.

## Dependency and execution boundary

The module uses the Go standard library only. It performs no network or filesystem I/O during validation. The local `cmd/perf-probe` uses synthetic references only and exists solely to produce deterministic CI performance evidence.


## Search Job orchestration

The bounded runtime now also provides a backend-neutral Search Job orchestrator:

- canonical states: queued, running, completed, partial, failed, cancelled;
- only queued/running cancellation is permitted;
- terminal jobs never restart in place;
- retry creates a new immutable Job ID and Run ID while preserving tenant, environment, time, source, correlation and query-version provenance;
- backend results contain stable references only, never raw event content;
- every requested source must be explicitly accounted for as completed or failed;
- cross-tenant, unaccounted, duplicate or malformed backend results fail closed and do not populate protected result references;
- partial completion preserves valid result references and exposes stable failed-source codes;
- backend error text is never persisted in the Search Job audit trail.

The backend remains a Go interface. No query dialect, search engine, storage engine, index or provider is selected by this implementation.
