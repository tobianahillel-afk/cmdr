# Source of Truth Policy

## Objective

Prevent contradictory product definitions and ensure every CMDR concept has one authoritative location.

## Precedence

When documents conflict, use this order:

1. `00-governance/` for documentation rules and approved decisions.
2. `02-domain-model/` for entities, relationships, states and permissions.
3. `03-experience/` for shared interaction and visual behaviour.
4. `04-cross-console/` for hand-offs between consoles.
5. Console page specifications for page-local behaviour.
6. `40-platform/` for shared technical product capabilities.
7. Examples and reference scenarios, which are illustrative and never normative.

## Ownership rules

- Every concept has one canonical document.
- Other documents link to the canonical definition and describe only their local use.
- Page documents may name fields required for a view, but may not redefine the full object schema.
- State names must come from `STATE_MODELS.md`.
- Permission names and enforcement semantics must come from `PERMISSION_MODEL.md`.
- Shared navigation, search, notifications and accessibility behaviour must not be copied into individual pages.
- Example identifiers such as `INC-2026-0712` or `R-0241` are non-production examples.

## Change protocol

A change that affects several consoles must update the canonical shared document first, then adjust dependent page specifications. The pull request must list impacted documents and acceptance criteria. Unresolved disagreement is recorded in `DECISION_LOG.md` or `50-quality/OPEN_GAPS.md`; it must not be silently resolved by duplication.

## Acceptance criteria

- No two active documents claim canonical ownership of the same concept.
- Every page links to shared definitions it depends on.
- Shared state and permission names are consistent.
- Example data is clearly distinguishable from normative requirements.
- Deleted or superseded concepts remain traceable through Git history rather than copied “legacy” sections.

## Open questions

- Whether formal ADR numbering will be introduced once implementation begins.
- Whether document ownership metadata should later move to machine-readable front matter.
