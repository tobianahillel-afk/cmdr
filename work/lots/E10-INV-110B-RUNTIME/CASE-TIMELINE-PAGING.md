# E10-INV-110B — bounded authorized Case Timeline HTTP pages

## Functional increment

The already implemented authenticated read-only Case Timeline HTTP endpoint now accepts optional `limit` (1..200) and `offset` (0..5000), with `offset` valid only if a `limit` is present. It keeps the existing `case_id`, `kind`, `from` and `to` selectors. The default response `{"status":...,"entries":...}` is unchanged. When pagination is explicitly requested, the response adds a `page` object with `offset`, `limit`, `total_visible` and `has_more`. Query selectors are parsed before any provider read.

A request always performs **complete** Case/Timeline Entry/canonical linked-source authorization and deterministic projection before filtering and pagination. Neither page indices nor query headers are granted authority; `total_visible` counts only authorized visible entries *after optional kind/time filtering*. Inaccessible source IDs, source metadata or their counts are not included. The status remains that of the full filtered authorized projection, so an empty page does not hide previous incompleteness.

## Reuse-first technical decision

Selected: Go standard-library `strconv`, simple slice ranges, and existing `CaseTimelineHTTPHandler` and `selectAuthorizedTimeline`. Rejected for this bounded in-memory and provider-injected scope: a new indexing service, database, GraphQL framework, signed cursor protocol, and UI-specific ACL engine. The entire Case is at most 5000 entries and the existing owner authorization pipeline already resolves them. Server-checked simple offset pages avoid new state and dependencies.

## Important delivery limits

This is a **live offset**, not a cryptographically bound or transactionally stable cursor: changes to real owner repositories between requests can shift items. The underlying Case and source access authorities are checked afresh for each request. The memory-backed demo is executable but not a production identity or persistent ACL service. No provider read batching, durable Case store, final navigation UI, Shared Timeline Entry mutation, canonical schema changes or production SLO is implied. `OPEN-013` remains open.

Acceptance: existing unpaged responses unchanged; 2+ consecutive pages provide all authorized chronologically ordered entries without duplicates within an unchanged fixture; combined kind/time/paging filters work; missing/invalid selectors are rejected before owner I/O; cross-tenant and denied source contents remain concealed. Tests and CI must pass at the exact published SHA, with independent L2 assigned only at the material functional milestone.
