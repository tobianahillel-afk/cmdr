---
id: source-ux-and-navigation-decisions
domain: 00-governance
status: draft
owner: UX Architecture
updated: 2026-08-03
source-of-truth: source-material
---
# UX and Navigation Decisions

| Requirement | Decision |
|---|---|
| REQ-UX-001 | A page or workspace represents a genuinely different user objective. |
| REQ-UX-002 | A view is a saved subset or organization of the same objects. |
| REQ-UX-003 | A mode is another representation of the same space. |
| REQ-UX-004 | A filter is a temporary restriction of data. |
| REQ-UX-005 | There is one canonical Inspector. |
| REQ-UX-006 | Cross-product context is preserved. |
| REQ-UX-007 | Back returns to the real previous state. |
| REQ-UX-008 | Work Queue variants are saved views, not six pages. |
| REQ-UX-009 | Mission Control Now, Priorities, Situation and Timeline are views or modes. |
| REQ-UX-010 | Screen sections must be substantive and screen-specific. |

## Canonical shells

Global Shell, Queue Shell, Case Shell, Technical Workbench Shell, Decision Shell, Run Shell, Builder Shell and Settings Shell.

## Technical Workbench

Global header, Context Bar, left explorer, central canvas, one canonical right Inspector, optional bottom console, artifact tabs, Tool Dock and Automation Tray.

Constraints: one main canvas, one right Inspector, optional bottom console, maximum two auxiliary panels and maximum six visible technical tabs.

## Context preservation

Tenant, environment, active Incident, Case, Endpoint and Artifact, filters, selected view, tabs, scroll, display mode, selection and navigation history are preserved when relevant.