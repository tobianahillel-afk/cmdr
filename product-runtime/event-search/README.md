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
- Search Job lifecycle transitions, cancellation or partial-result orchestration;
- Case-link mutation.

Those concerns remain owned by later bounded tasks and unresolved product decisions.

## Dependency and execution boundary

The module uses the Go standard library only. It performs no network or filesystem I/O during validation. The local `cmd/perf-probe` uses synthetic references only and exists solely to produce deterministic CI performance evidence.
