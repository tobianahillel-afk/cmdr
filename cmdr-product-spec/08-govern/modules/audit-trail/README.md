---
id: govern-audit-trail
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-PROD-002, REQ-PROD-004, REQ-PROD-005, REQ-PROD-008, REQ-PROD-015, REQ-PROD-019, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-008, OPEN-013, OPEN-015, OPEN-019]
---
# Audit Trail — GOV-3

## Mission

Interpret and reconstruct the Govern lifecycle from Action Request through Result for review, completeness assessment and evidence-package preparation while preserving source ownership and relying on Shared Trace/Activity/Search/Reporting/Export mechanisms rather than creating duplicate infrastructure.

## Owned GOV-3 capabilities

- `CAP-GOV-034` — Govern Audit Trail Intake and Event Semantics.
- `CAP-GOV-035` — Decision, Approval and Authority Audit Reconstruction.
- `CAP-GOV-036` — Response Run, Verification, Rollback and Result Audit Reconstruction.
- `CAP-GOV-037` — Audit Completeness, Integrity, Gap and Contradiction Assessment.
- `CAP-GOV-038` — Govern Audit Review, Search and Evidence Package Preparation.

## Ownership boundary

Govern owns **Govern audit semantics and review outputs**. Shared retains Trace, Activity, Search, Reporting, Export and generic audit plumbing. Platform Settings retains retention/storage/export-destination configuration. Security retains permission, privacy, legal-hold and integrity policy. Endpoint/Studio/other runtimes retain raw technical records.

## Mandatory distinctions

- Audit Trail ≠ Trace infrastructure.
- Audit Trail ≠ Activity feed.
- Govern Audit Event ≠ raw log/SIEM event.
- audit reconstruction ≠ original execution.
- audit record present ≠ action correct.
- audit record absent ≠ action did not happen.
- audit completeness ≠ truth completeness.
- gap ≠ malicious behavior.
- contradiction ≠ automatic falsification.
- timestamp order ≠ causality.
- correlation ≠ causation.
- immutable requirement ≠ implemented immutable storage.
- integrity ≠ cryptographic proof automatically.
- Audit Evidence Package ≠ canonical Investigate Evidence.

## Source lifecycle

Audit Trail consumes GOV-1 `CAP-GOV-001..016` and GOV-2 `CAP-GOV-017..033` as source semantics. It never mutates or redefines Action Request, Approval, Decision, Playbook, Response Run, Verification Assessment, Rollback/Recovery or Result.

## Sensitive data

Actor identities, target context, exceptions, emergency paths and technical references are permission-aware. Existence, metadata, masked preview, read, export and external sharing are separate operations. Raw secret material is prohibited.

## AI/no-AI

AI may summarize sourced sequences or suggest missing-link/contradiction candidates. It cannot invent an event, declare fraud/violation, alter history or hide conflict. Deterministic correlation, source tables, filters and human review provide the full non-AI path.

## Screen

`GOV-AUD-001` remains the existing Audit Trail surface. GOV-3 creates no Screen ID and does not define final columns, filters, buttons, animations, shortcuts or wireframes. Shared Inspector/Trace/Reporting mechanisms remain source-owned.

## GOV-3 outputs

Audit reconstruction, completeness/gap/contradiction assessments, Review Summary and permission-aware Audit Evidence Packages. These outputs can feed Investigate or Reporting by explicit handoff, but do not become canonical Evidence automatically.

## Dependencies

GOV-1/GOV-2 provenance, Shared Trace/Activity/Search/Reporting/Export/Versioning, Settings retention/storage/tenant context, Security permission/privacy/integrity/legal-hold rules and OPEN-007/008/013/015/019.