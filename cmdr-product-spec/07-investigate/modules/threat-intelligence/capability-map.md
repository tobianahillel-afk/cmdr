---
id: investigate-threat-intelligence-capability-map
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements:
  - REQ-PROD-012
  - REQ-PROD-014
  - REQ-INV-006
open_decisions:
  - OPEN-018
  - OPEN-019
---
# Capability map — Threat Intelligence

## Foundations and knowledge management
CAP-INV-501..518 cover Intake, Requirements, Knowledge Projects, sources/materials, candidates, Sightings, relationships, confidence, versioning/lifecycle and Analysis Handoff Packages.

## Analysis, products, dissemination and operationalization
| ID range | Functional family | Count |
|---|---|---:|
| CAP-INV-519..521 | Analysis Sessions, questions, competing hypotheses, fusion and structured analysis | 3 |
| CAP-INV-522..524 | Actor/attribution, Campaign/activity and malware/tool/infrastructure assessments | 3 |
| CAP-INV-525..527 | Product planning, authoring, quality review and release recommendation | 3 |
| CAP-INV-528..533 | Releasability, internal publication, watchlists, operationalization, monitoring and external-sharing preparation | 6 |
| CAP-INV-534..537 | Feedback, effectiveness, Requirement satisfaction, correction and lifecycle provenance | 4 |
| **Phase 4B.3B.2** | **CAP-INV-519..537** | **19** |

## Totals and boundaries
- Threat Intelligence: **37 capabilities, 999 sections, 222 mandatory tables**.
- All capabilities are `defined` / `planned`; implementation remains absent.
- Internal publication is reversible and policy-bound; it grants no source access.
- Watchlist Definition is not an active watchlist.
- Operationalization Package is not a deployed Indicator, Detection Content, Signal, block or runtime mutation.
- External sharing is prepared only; Govern and destination owners retain execution authority.
- OPEN-018 retains ontology/interoperability; OPEN-019 retains dissemination/releasability/sharing policy.
