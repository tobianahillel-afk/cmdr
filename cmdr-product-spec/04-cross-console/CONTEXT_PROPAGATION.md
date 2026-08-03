# Context Propagation Contract

## Preserved context

Tenant, user scope, source object, selected time range, selected entities, relevant filters, correlation identifier and return URL.

## Explicitly transformed context

- Incident summary becomes Case intake context.
- Evidence-backed Finding becomes Response justification.
- Decision and Run summary become Incident operational updates.

## Excluded context

Secrets, raw malware bytes, unsupported UI-only state, hidden policy internals and data outside the destination user's permissions.

## Conflict behaviour

If target scope differs or access is missing, the transition stops before object creation and explains the blocking condition. If a matching target already exists, the user is offered that object rather than silently duplicating it.

## Acceptance criteria

A transition survives refresh, is traceable in audit history and never broadens tenant or object permissions.
