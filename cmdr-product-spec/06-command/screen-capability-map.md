---
id: command-screen-capability-map
domain: 06-command
status: draft
owner: Command Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-UX-008
  - REQ-UX-010
  - REQ-PROD-013
open_decisions:
  - OPEN-006
  - OPEN-010
---
# Screen readiness and capability map — Command

| Écran actif | Capability principale | Capabilities secondaires | Problèmes connus | Phase de réécriture |
|---|---|---|---|---|
| CMD-MC-001 — Mission Control — Now | CAP-CMD-001 | CAP-CMD-003, CAP-CMD-005, CAP-CMD-006 | screen remains Template-level; content hierarchy later | Phase 6 |
| CMD-MC-002 — Mission Control — Priorities | CAP-CMD-002 | CAP-CMD-104, CAP-CMD-105, CAP-CMD-205 | priority factors and layout later | Phase 6 |
| CMD-MC-003 — Mission Control — Decisions | CAP-CMD-001 | CAP-CMD-006, CAP-CMD-110 | Govern projection only | Phase 6 |
| CMD-MC-004 — Mission Control — Situation | CAP-CMD-003 | CAP-CMD-001, CAP-CMD-006 | timeline grouping later | Phase 6 |
| CMD-MC-005 — Mission Control — Handover | CAP-CMD-004 | CAP-CMD-001, CAP-CMD-005, CAP-CMD-006 | handover record object unresolved | Phase 6 after objects |
| CMD-IWQ-006 — Incident Detail | CAP-CMD-106 | CAP-CMD-102..110 | 27-section screen remains generic | Phase 6 |
| CMD-EXC-001 — Exposure & Coverage | CAP-CMD-202 | CAP-CMD-201, CAP-CMD-203 | legacy module path; canonical Risk and Coverage | Phase 6 migration |
| CMD-RBI-001 — Risk & Business Impact | CAP-CMD-205 | CAP-CMD-201, CAP-CMD-203, CAP-CMD-204 | legacy module path; canonical Risk and Coverage | Phase 6 migration |
| CMD-RDO-001 — Readiness & Operations | CAP-CMD-301 | CAP-CMD-302..305 | screen content later | Phase 6 |
| CMD-CRP-001 — Customer & Reports | CAP-CMD-401 | — | module proposed; OPEN-006 | Phase 6 only if retained |
| Work Queue workspace — No active Screen ID | CAP-CMD-101 | CAP-CMD-102..110 | six views, no page file | Phase 6 screen/workspace contract |

## Phase 4A disposition

- 10 active Command screen files were read and not rewritten;
- five legacy Work Queue files remain `deprecated`;
- no new Screen ID was created;
- Work Queue views remain configurations;
- this matrix does not replace screen specifications.
