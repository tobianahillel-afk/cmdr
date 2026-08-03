# Canonical State Models

## Objective

Provide the only authoritative state names used by CMDR pages and APIs.

## Alert

`new → acknowledged → triaged → linked | closed`

A closed alert requires disposition and rationale. Reopening creates an audit event.

## Incident

`detected → active → contained → eradicated → recovered → closed`

Side states: `monitoring`, `on_hold`. Severity and business impact are independent attributes, not states.

## Case

`intake → triage → investigating → awaiting_input → review → completed → archived`

Archived cases are read-only except for legal-hold and administrative metadata.

## Evidence

`registered → collecting → available → verifying → verified | rejected → superseded`

Integrity failure moves evidence to `rejected`; it is never silently deleted.

## Hypothesis

`proposed → testing → supported | weakened | refuted → closed`

## Finding

`draft → review → confirmed | rejected → superseded`

Only confirmed findings may be treated as established facts in a response request.

## Response request

`draft → submitted → policy_check → awaiting_approval → decided → cancelled`

## Decision

`pending → approved | approved_with_conditions | more_information_required | refused → superseded`

## Run

`queued → preparing → running → succeeded | failed | partially_succeeded | cancelled → rollback_pending | rolling_back → rolled_back | rollback_failed`

## Agent

`unregistered → enrolling → online → degraded → offline → revoked`

## Rules

Transitions require permission, preconditions and audit. UI labels may be localised but stored state keys remain stable. Pages must not invent alternate states.
