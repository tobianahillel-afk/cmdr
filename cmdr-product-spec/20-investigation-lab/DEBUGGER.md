# Debugger

## Objective

Support controlled step/run debugging of suspicious code inside isolated analysis environments.

## Scope

This specification owns the page-local behaviour of **Debugger** in **Investigation Lab**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

DFIR Lead

## Affected objects

- Evidence
- Debug session
- Breakpoint
- Memory snapshot
- Process/thread entity
- Finding

## Features

- Launch and attach in isolated targets
- Run, pause, step into/over/out
- Breakpoint and watch management
- Registers, stack, memory and modules
- Thread and exception views
- Snapshot and trace capture
- Promote observations to evidence

## UX and interactions

- Execution state is always prominent
- Dangerous resume/attach actions show target profile
- Keyboard controls are discoverable but protected against accidental activation
- Session loss and worker failure preserve collected evidence
- Debugger is visually distinct from response execution

## Permissions

`debugger.execute`; profile restrictions use ABAC; evidence export remains separate.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Debug session states: provisioning, paused, running, terminated, failed. These are local tool states, not Run states for response orchestration.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Isolated debugger workers, reverse engineering, evidence storage.

## Acceptance criteria

- Debugger cannot attach to production endpoints
- Every command is timestamped and attributable
- Snapshots retain provenance
- Session termination is explicit
- Collected evidence remains available after failure

## Open questions

- Which debugger backends are supported?
- How are anti-debug countermeasures surfaced?
