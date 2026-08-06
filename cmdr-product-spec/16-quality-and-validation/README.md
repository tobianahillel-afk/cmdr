---
id: quality-readme
domain: 16-quality-and-validation
status: draft
owner: Quality Lead
updated: 2026-08-07
source-of-truth: canonical
---
# Quality and Validation

## Objectif
Définir les contrôles qui comparent l’arborescence réelle aux contrats attendus et valident liens, noms, front matter, sources, permissions, propriétaires, métriques et preuves de publication.

## Active evidence
- Validation status: `validation-status.md`.
- Structural baseline: `architecture-manifest.md` and `expected-path-manifest.md`.
- Cloud Analysis conformance: `reports/phase-4b4a-cloud-analysis-capability-conformance.md`.
- Cloud Analysis source audit: `reports/phase-4b4a-cloud-analysis-source-audit.md`.
- Prior Threat Intelligence closure reports remain active and preserved.

## Validation rules
- A planned capability is not implementation evidence.
- A report cannot claim post-publication PASS before the remote head, PR, README, `main`, registers and metrics are checked.
- Historical evidence is appended or preserved; closure never condenses earlier traceability.
- Sources are classified as fully read, referenced, identified or absent.
- A failed or pending gate remains visible.

## Permissions and ownership
Quality reads canonical sources and records evidence; it does not own product behavior, objects, permissions, roadmap decisions or implementation.

## Critères d’acceptation
- Every verdict names the evidence and its verification stage.
- Counts are recalculated from the final registered scope.
- No duplicate ID, owner conflict, active contradiction, broken required link or unsupported implementation claim is hidden.
- Cloud Analysis remains pending until remote verification is complete.
