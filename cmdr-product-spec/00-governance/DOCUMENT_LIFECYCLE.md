# Document Lifecycle

## Objective

Define how a CMDR specification moves from draft to an implementation-ready contract.

## Statuses

- **Draft** — useful structure exists; open product decisions remain.
- **Review** — functional owner considers scope complete and requests cross-domain review.
- **Approved** — behaviour and acceptance criteria are authoritative.
- **Implemented** — implementation evidence is linked and criteria pass.
- **Deprecated** — retained only for traceability; not a source of truth.

## Required sections

Every page specification contains: objective, scope, functional owner, affected objects, features, UX and interactions, permissions, states, dependencies, acceptance criteria and open questions. Canonical shared documents may use a structure better suited to their subject but must still state ownership and acceptance conditions.

## Review gates

A document cannot become Approved unless:

- object, state and permission references are canonical;
- dependencies are explicit;
- empty headings and placeholder-only text are absent;
- acceptance criteria are testable;
- open questions are either non-blocking or assigned;
- internal links resolve.

## Change discipline

Use commits grouped by product domain. A broad architectural change should be isolated from page-detail changes. The pull request description states scope, validation performed, unresolved gaps and whether the change is documentation-only.

## Owner

Product Architecture, with domain-owner approval.
