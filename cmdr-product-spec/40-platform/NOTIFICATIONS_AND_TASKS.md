# Notifications and Tasks

## Objective

Deliver actionable, deduplicated attention signals without becoming a second unstructured queue.

## Scope

This specification owns the page-local behaviour of **Notifications and Tasks** in **Shared Platform**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

Platform Product Lead

## Affected objects

- Notification
- Task
- Principal
- Subscription
- Incident
- Decision
- Run

## Features

- In-product inbox
- Email/chat integrations
- Subscriptions and routing
- Escalation and reminders
- Task ownership and due dates
- Digest and quiet-hours policy
- Deduplication

## UX and interactions

- Notifications deep-link to exact context
- Acknowledging a notification does not change source object state
- Users control non-mandatory channels
- Critical governance reminders cannot be silently muted

## Permissions

Viewing is self-scoped; routing administration requires tenant or platform permissions; source actions retain their own permissions.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Notification states unread/read/dismissed; Task states open/in_progress/blocked/done/cancelled.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Identity, all domain events, integration connectors, scheduler.

## Acceptance criteria

- Duplicate events collapse predictably
- Delivery failures are visible
- Deep links preserve tenant context
- Mandatory notifications are policy-controlled
- Task completion does not silently close source objects

## Open questions

- Which channels are required at launch?
- How are on-call schedules integrated?
