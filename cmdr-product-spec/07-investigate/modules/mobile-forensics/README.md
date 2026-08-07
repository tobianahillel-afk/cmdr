---
id: investigate-mobile-forensics
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-07
source-of-truth: canonical
requirements: [REQ-PROD-012, REQ-PROD-014, REQ-PROD-019, REQ-PROD-020, REQ-PROD-055, REQ-INV-001, REQ-INV-006]
open_decisions: [OPEN-005, OPEN-008, OPEN-011, OPEN-013, OPEN-014, OPEN-015]
---
# Mobile Forensics

## Canonical phase
**Phase 4B.4B — Mobile Forensics Foundations and Mobile Investigation** is the sole Mobile subphase. `CAP-INV-701..719` is the reserved functional range.

## Mission
Provide authorized Mobile Forensics Analysts, DFIR Analysts, Investigation Leads, Evidence Reviewers, SOC Analysts and Sensitive Data Reviewers with a provider-neutral analytical context for mobile evidence representations, while keeping acquisition, device administration, production response and generic mechanisms with their canonical owners.

## Functional chain
Case / Incident / Finding / Hunt / Signal / Collection Job / Artifact / Mobile Evidence Package → Intake → Session → Device/Platform/Scope → Acquisition Context → Integrity/Completeness/Accessibility → Filesystem/Applications/Communications/Media/Location/Sensitive Material/Connectivity/Backups/Deleted Data → Mobile Timeline and Hypotheses → Derived Artifacts and handoffs → Provenance/Reproducibility.

## Ownership
Investigate owns Mobile analytical Sessions, observations, candidates, assessments, Hypotheses, Mobile Timeline semantics, Derived Artifact selection and handoff packages. Collection owns acquisition preparation/execution/results and custody at collection. Platform Settings owns configured mobile sources, MDM/EMM providers, connectors, credentials, secrets, Fleet, storage, retention and policies. Endpoint Agent contributes only declared authorized capabilities when present. Govern owns all real-device actions. Studio owns Tools/Runs. Shared owns Entity, Graph, Timeline engine, Search, Linking, Jobs, Notifications, Trace, Activity, Versioning, Export, Reporting, Collaboration and Recovery.

## Invariants
- acquisition representation ≠ original device;
- Mobile Evidence Package ≠ Mobile Forensics Session ≠ Evidence;
- Derived Artifact ≠ Evidence;
- device/account/contact/number/location records do not establish a person or human presence automatically;
- application/message/call/media records do not establish user intent or authorship automatically;
- partial, encrypted, locked, inaccessible, deleted, recovered and synchronized states remain explicit;
- secret presence never grants reveal, copy, export or use;
- Tool/AI output never becomes an analyst conclusion, Evidence or Finding automatically;
- Mobile network observations do not replace Network Forensics; cloud-backed artifacts do not replace Cloud Analysis.

## Delivery
Documentary and functional only. All capabilities are `draft`, `defined`, `planned`. No platform, tool, engine, API, protocol, proprietary format, command, unlock/bypass/root/jailbreak method, MDM action or product code is delivered.

## Reading order
`scope.md`, `user-questions.md`, `concepts.md`, `workflows.md`, `states.md`, `permissions.md`, `sensitive-data-and-privacy.md`, `shared-capabilities.md`, `automation-and-ai.md`, `collection-and-acquisition-boundaries.md`, `settings-endpoint-studio-govern-boundaries.md`, `capability-map.md`, `source-migration.md`, then `capabilities/`.
