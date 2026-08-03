# Mission Control

## Objective

Give SOC leaders and analysts a truthful global overview of active security work, priority, business impact, team load and pending governed decisions.

## Scope

This specification owns the page-local behaviour of **Mission Control** in **Command Center**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

SOC Operations Lead

## Affected objects

- Incident
- Alert
- Business service
- Decision
- Run
- Agent health summary

## Features

- Operational KPI cards with trend and freshness
- Priority queue with ownership and SLA pressure
- Incident heatmap and business-impact summary
- Decisions awaiting attention
- Operational timeline and system-health strip
- Tenant and time-range filtering

## UX and interactions

- Opening an incident preserves queue position
- Map markers use accessible shape-aware focus
- Live updates patch cards without full-page refresh
- Users can freeze the time window for review
- Metrics expose definitions and drill-downs

## Permissions

`command.read`; coordination actions additionally require `command.coordinate` or `incident.manage`.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Incident and Run states are shown independently; freshness is an attribute, not a state.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Normalised data, incident service, business-service catalogue, governance summaries, platform health.

## Acceptance criteria

- Counts reconcile with filtered queues
- Every metric exposes source and last update
- No device or incident without usable location pollutes map metrics
- Keyboard users can inspect every map item through a list alternative
- Cross-console links preserve tenant and incident context

## Open questions

- Which executive metrics are tenant-configurable?
- What latency qualifies a card as stale?
