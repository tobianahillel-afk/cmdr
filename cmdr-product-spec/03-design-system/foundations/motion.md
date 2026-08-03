---
id: foundation-motion
domain: 03-design-system
status: draft
owner: Design System Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-BRAND-001
  - REQ-BRAND-002
---
# Motion foundation — source boundaries

## Responsibility

The Design System will translate the calm, accountable motion language into exact durations, easing curves and component transitions during Phase 3.

## Canonical source

- [`../../02-brand/motion-and-sound.md`](../../02-brand/motion-and-sound.md)

## Phase 3 responsibility

Phase 3 will specify tokens and behavior for:

- state changes;
- disclosure;
- overlays;
- cross-product continuity;
- progress and live updates;
- interruption and rollback;
- reduced-motion alternatives.

Any duration ranges in the brand document are directional, not implementation tokens.

## Constraints

- Motion explains continuity, hierarchy or state.
- No continuous decorative movement, pulsing urgency or ambient particles.
- Criticality does not rely on blinking.
- `prefers-reduced-motion` preserves information and task completion.
- Sound is opt-in and operationally governed.

## Acceptance criterion

**Given** an Automation Run changes from running to paused,  
**When** the state transition is presented,  
**Then** motion makes the change legible without spectacle, the textual status updates immediately and reduced-motion users receive an equivalent non-animated cue.
