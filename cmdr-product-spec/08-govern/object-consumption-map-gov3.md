---
id: govern-object-consumption-map-gov3
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical-addendum
---
# Object Consumption Map — Govern GOV-3 Addendum

This additive map extends `object-consumption-map.md` without rewriting GOV-1/GOV-2 evidence. It defines concepts only; no complete schema, JSON Schema, final cardinality or state machine.

| Objet ou concept | Owner actuel | Usage Govern | Opérations locales | Lacune | Phase propriétaire |
|---|---|---|---|---|---|
| Govern Audit Event | Govern local concept | interpreted lifecycle event over source records | create/interpret/version/supersede | technical event/storage model future | Govern / future Objects/Technique |
| Audit Reconstruction | Govern local concept | historical Decision/Run chain review | create/update/version | reconstruction engine future | Govern |
| Audit Gap | Govern local concept | expected/observed provenance gap candidate | create/review/resolve/supersede | gap taxonomy/storage future | Govern |
| Audit Contradiction | Govern local concept | conflicting source/state candidate | create/review/dispute/supersede | contradiction model future | Govern |
| Audit Completeness Assessment | Govern local concept | scoped completeness/source-quality review | create/update/version | no truth/cryptographic claim | Govern |
| Audit Evidence Package | Govern local concept | sourced package for review/export/handoff | create/version/supersede | package ≠ Investigate Evidence | Govern / Investigate handoff |
| Metric Definition concept | Govern semantics + Shared mechanism | define Govern metric meaning/dimensions | reference/version | generic engine Shared-owned | Shared + Govern semantics |
| Metric Observation | Govern local concept | sourced metric snapshot | calculate/version/supersede | physical metric store future | Govern / Shared Metrics |
| Policy Metric | Govern | Policy/Exception/Emergency observations | calculate/compare | no universal threshold | Govern |
| Approval Metric | Govern | Approval lifecycle observations | calculate/compare | privacy/final dimensions future | Govern |
| Authority Metric | Govern | authority/SoD observations | calculate/compare | no role-to-authority inference | Govern/Security |
| Decision Metric | Govern | disposition/timeliness observations | calculate/compare | no universal SLO | Govern |
| Response Run Metric | Govern | Run reliability/partiality observations | calculate/compare | raw technical metrics source-owned | Govern |
| Verification Metric | Govern | verification outcome observations | calculate/compare | verification engine future | Govern |
| Rollback Metric | Govern | rollback state observations | calculate/compare | rollback primitive source-owned | Govern |
| Recovery Metric | Govern | recovery outcome observations | calculate/compare | no full-restoration inference | Govern |
| Result Metric | Govern | canonical Result outcome observations | calculate/compare | effectiveness ontology future | Govern |
| Flow Metric | Govern | Inbox/lifecycle ageing/stage durations | calculate/compare | final SLO/state model future | Govern |
| Trend Assessment | Govern local concept | cross-period/scope comparison | create/review/supersede | statistical method future | Govern |
| Control Health Assessment | Govern local concept | sourced control-health interpretation | create/review/dispute/supersede | not automatic truth | Govern |
| Continuous Improvement Package | Govern local concept | no-effect improvement proposal | create/review/route/supersede | destination change owner-specific | Govern / destination owners |
| Govern Closure Assessment | QA/Product Architecture + Govern context | documentary closure evidence | create/update/finalize | implementation validation future | Quality/Roadmap |
| Trace | Shared | generic provenance infrastructure | consume/link | no Govern engine | Shared |
| Activity | Shared | generic activity infrastructure | consume/link | no Govern clone | Shared |
| Report | Shared Reporting | rendering/publication of authorized snapshots | request/link | Report remains Shared-owned | Shared |
| Export Job | Shared | authorized export execution | request/observe | external sharing separate | Shared / OPEN-019 |
| Decision | Govern | source for reconstruction/metrics | read/link | GOV-1 immutable by default | Govern GOV-1 |
| Approval | Govern | source for reconstruction/metrics | read/link | GOV-1 immutable by default | Govern GOV-1 |
| Response Run | Govern | source for audit/metrics | read/link | GOV-2 immutable by default | Govern GOV-2 |
| Result | Govern | source for outcome metrics | read/link | GOV-2 immutable by default | Govern GOV-2 |
| Incident | Command | downstream/context projection | read/link/aggregate under permission | no owner transfer | Command |
| Case / Finding | Investigate | audit/improvement handoff context | read/link only | package ≠ Evidence/Finding | Investigate |
| Workflow / Automation Run | Studio | source technical provenance/metrics projection | read/link only | Automation Run ≠ Response Run | Studio |
| Provenance | source owners + Shared mechanisms | chain GOV-1→GOV-3 | read/correlate/export-prep | physical model future | Shared/Security/Govern |

## Invariants
A Govern audit or metric concept never mutates its source object. Audit Trail ≠ Trace/Activity. Audit Evidence Package ≠ canonical Evidence. Metric ≠ objective/Policy/SLO/KPI automatically. Trend ≠ causal explanation. Control Health Assessment ≠ automatic truth. Continuous Improvement Package applies no change. Raw secrets are excluded.