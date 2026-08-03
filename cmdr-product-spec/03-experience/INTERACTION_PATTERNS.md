# Shared Interaction Patterns

## Lists and queues

Support saved views, filters, sort, column selection, bulk selection where safe, pagination or virtualisation, ownership and empty/error/loading states. Bulk destructive actions require explicit review.

## Detail panels

Use side panels for quick inspection and full pages for sustained work. Opening a panel must not reset the underlying queue.

## Timelines

Timelines distinguish observed events, inferred assertions, user notes, decisions and executed actions. Every item exposes provenance and timezone.

## Evidence interactions

Preview never alters original evidence. Download/export requires explicit permission and creates an audit event. Derived artefacts retain parent provenance.

## Confirmation

Irreversible or high-impact operations show target, scope, impact, prerequisites and rollback availability. Typed confirmation is reserved for exceptional risk, not routine friction.

## Autosave and conflicts

Notes, hypotheses and draft findings autosave with visible status. Concurrent edits show version conflict and provide merge/reload options; last-write-wins is not acceptable for analytical conclusions.

## Errors

Errors explain what failed, what remains unchanged, correlation identifier and safe next action. Partial success is represented explicitly.
