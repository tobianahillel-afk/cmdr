# Entity Graph

## Objective

Explore relationships between hosts, accounts, processes, files, networks, cloud resources and services.

## Scope

This specification owns the page-local behaviour of **Entity Graph** in **Investigation Lab**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

DFIR Lead

## Affected objects

- Entity
- Relationship
- Evidence
- Event
- Finding
- Case

## Features

- Graph and table views
- Relationship provenance
- Temporal filtering
- Path search
- Entity merge/split review
- Pivot to evidence and event search
- Business-service overlay

## UX and interactions

- Graph never implies causality from proximity alone
- Every edge exposes source and confidence
- Keyboard/list alternative is complete
- High-density graphs support progressive expansion
- Merge operations require review

## Permissions

`case.read`; merge/split requires `case.manage`; sensitive entity types are ABAC-scoped.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Entities may be active/inactive/merged; analytical confidence belongs to relationships or Findings.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Entity resolution, event search, evidence, asset inventory.

## Acceptance criteria

- Every edge is source-backed
- Graph and table counts reconcile
- Merges are reversible and audited
- Cross-tenant edges are prohibited
- Temporal filters affect nodes and edges consistently

## Open questions

- What entity-resolution rules are automatic?
- How are privacy-sensitive identities masked?
