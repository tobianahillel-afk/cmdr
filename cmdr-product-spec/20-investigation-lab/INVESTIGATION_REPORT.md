# Investigation Report

## Objective

Create a reproducible report separating facts, analysis, uncertainty, findings and recommendations.

## Scope

This specification owns the page-local behaviour of **Investigation Report** in **Investigation Lab**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

DFIR Lead

## Affected objects

- Report
- Case
- Evidence
- Finding
- Hypothesis
- Entity
- Timeline entry

## Features

- Structured report sections
- Evidence and finding citations
- Timeline and entity visualisations
- Confidence and limitations
- Peer review and approval
- Versioned export
- Redaction profiles

## UX and interactions

- Report text never silently changes source objects
- Citations resolve to stable evidence versions
- Preview matches export
- Redactions are visible to the author
- Review comments remain outside final publication unless accepted

## Permissions

`report.manage`; evidence visibility and export permissions are re-evaluated; publication may require `finding.approve` or designated report approval.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Report lifecycle is draft/review/approved/published/archived; linked object states remain independent.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Case, evidence, findings, reporting service, redaction, audit.

## Acceptance criteria

- Every factual claim can cite a source
- Published versions are immutable
- Redaction is tested against export
- Review history is retained
- Report generation cannot broaden data scope

## Open questions

- Which external report standards are required?
- Who approves reports for legal use?
