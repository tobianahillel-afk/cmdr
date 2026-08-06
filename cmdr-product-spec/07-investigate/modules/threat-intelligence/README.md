---
id: investigate-threat-intelligence
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements: [REQ-PROD-014, REQ-PROD-019, REQ-PROD-020, REQ-INV-006]
open_decisions: [OPEN-013, OPEN-014, OPEN-015, OPEN-018, OPEN-019]
---
# Threat Intelligence

## Mission
Transformer des besoins et matériaux sourcés en connaissances candidates, analyses structurées et produits Intelligence versionnés, puis préparer une diffusion interne contrôlée, des handoffs opérationnels et un cycle d’amélioration, sans perdre provenance, restrictions, alternatives ou responsabilité humaine.

## Capability ranges
- `CAP-INV-501..518` — Foundations and Knowledge Management: **18 capabilities, 486 sections, 108 tables**.
- `CAP-INV-519..537` — Analysis, Products, Dissemination and Operationalization: **19 capabilities, 513 sections, 114 tables**.
- Total Threat Intelligence: **37 capabilities, 999 sections, 222 tables**; all `defined` / `planned`.

## Functional chain
Intake → Requirement → Knowledge Project → Sources/Materials → Candidates/Sightings/Relationships → Confidence/Lifecycle → Analysis Handoff → Analysis Session → Questions/Competing Hypotheses → Fusion/Structured Assessment → Actor/Campaign/Malware/Infrastructure Assessments → Product Plan/Draft → Quality Review/Release Recommendation → Releasability/Dissemination Plan → Internal Publication or External-Sharing Preparation → Watchlist/Operationalization Handoffs → Monitoring/Feedback → Requirement Satisfaction/Correction → Lifecycle Provenance/Continuous Improvement.

## Ownership
Investigate owns Intelligence analytical and product concepts. Shared retains Entity, Graph, Timeline, Search, Linking, Versioning, Reporting, Export, Notifications, Collaboration, Trace, Activity and Recovery. Command retains runtime Detection, Signal, Alert, Incident and operational dispositions. Detection Engineering retains Detection Content/Hypothesis/Coverage/Gap/lifecycle. Settings retains providers, sources, connectors, secrets, access, destinations, storage, retention and health. Studio retains Tool/Tool Call/Workflow/Automation Run/Agent/Human Gate. Govern retains Decision, Approval, external release/sharing and production authority.

## Boundaries
Internal publication is controlled and reversible; it grants no raw-source permission. Watchlist Definition is not activation. Operationalization Package is not a rule, deployed Indicator, Signal, block or runtime change. External sharing is prepared only. No API, protocol, STIX/TAXII-like contract, provider, physical schema, runtime, code, Cloud Analysis, Mobile Forensics or detailed screen rewrite is added.
