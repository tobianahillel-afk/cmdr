# Command Center Settings

## Objective

Configure user-level Command Center preferences without duplicating tenant administration.

## Scope

This specification owns the page-local behaviour of **Command Center Settings** in **Command Center**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

SOC Operations Lead

## Affected objects

- Principal
- Saved view
- Dashboard layout
- Notification preference

## Features

- Default tenant and time range
- Saved views and dashboard layout
- Display density and theme shortcut
- Operational notification preferences
- Keyboard shortcut reference

## UX and interactions

- Changes preview immediately and can reset to default
- Settings distinguish personal from tenant-wide
- Accessibility preferences are easy to find
- Destructive reset is confirmable

## Permissions

`command.read`; personal settings are self-service; tenant-wide changes route to platform administration.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Settings have saved/saving/error states only; business object states are not redefined.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Identity, preferences, notifications, design system.

## Acceptance criteria

- Personal settings follow the user across sessions
- Tenant-wide controls are absent from this page
- Reset restores documented defaults
- Preference failures do not block core operations

## Open questions

- Which layouts are allowed to be user-customised?
- Should teams share saved views?
