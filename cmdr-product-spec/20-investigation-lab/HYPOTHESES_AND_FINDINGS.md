# Hypotheses and Findings

## Objective

Support explicit analytical reasoning from testable hypotheses to reviewed findings.

## Scope

This specification owns the page-local behaviour of **Hypotheses and Findings** in **Investigation Lab**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

DFIR Lead

## Affected objects

- Hypothesis
- Finding
- Evidence
- Entity
- Case
- Principal

## Features

- Hypothesis statement, confidence and tests
- Evidence for and against
- Alternative hypotheses
- Finding drafting and peer review
- Confidence rationale
- Evidence package and recommendation linkage

## UX and interactions

- Confidence changes require rationale
- Contradictory evidence is first-class
- Findings cannot be confirmed with missing required support
- Review comments are threaded and attributable
- Superseded conclusions remain traceable

## Permissions

`finding.author` to draft; `finding.approve` to confirm; evidence access remains separately enforced.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Uses canonical Hypothesis and Finding states.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Evidence board, case timeline, entity graph, peer review, audit.

## Acceptance criteria

- Confirmed findings link supporting evidence
- Refuted hypotheses remain visible
- Reviewer identity and rationale are recorded
- Recommendations cannot cite rejected findings as established facts
- Concurrent review conflicts are handled

## Open questions

- Is dual approval required for critical findings?
- How is confidence calibrated across teams?
