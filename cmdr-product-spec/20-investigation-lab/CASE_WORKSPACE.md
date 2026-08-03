# Case Workspace

## Objective

Provide the central investigation workspace linking summary, timeline, evidence, analysis, entities, hypotheses, findings, notes and activity.

## Scope

This specification owns the page-local behaviour of **Case Workspace** in **Investigation Lab**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

DFIR Lead

## Affected objects

- Case
- Incident
- Evidence
- Entity
- Hypothesis
- Finding
- Timeline entry
- Principal

## Features

- Case summary and status
- Structured section navigation
- Evidence and entity overview
- Hypothesis and finding panels
- Case notes and activity log
- Linked incidents and response requests
- Investigator assignments

## UX and interactions

- Sections preserve local state and filters
- Concurrent analytical edits detect conflicts
- Autosave is visible
- Critical conclusions require structured Findings, not free-text notes
- Cross-console links preserve exact case context

## Permissions

`case.read`; editing requires `case.manage`; finding actions use `finding.author` or `finding.approve`.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Uses canonical Case, Evidence, Hypothesis and Finding states.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

All Investigation Lab services, audit, notifications.

## Acceptance criteria

- All case changes appear in activity history
- Case completion checks unresolved evidence and hypotheses
- Notes cannot masquerade as confirmed findings
- Linked objects are bidirectional
- Refresh preserves active section

## Open questions

- Should cases support subcases?
- What conditions are required before archival?
