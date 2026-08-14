---
id: platform-scale-sources-parsers-data-source-parser-administration-capability-conformance
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-14
source-of-truth: quality-report
---
# Platform Scale — Sources & Parsers — Data Source and Parser Administration — Capability Conformance

## Scope

Execution lot: **Sources & Parsers — Data Source Administration and Parser Transformation Administration** under Delivery Roadmap Phase 6 — Platform Scale. This report validates documentary capability definition only; it does not prove runtime implementation, supported adapters, parser execution, ingestion, connector implementation, schema-standard selection or production deployment.

## Execution baseline

- repository: `tobianahillel-afk/cmdr`;
- branch: `docs/cmdr-product-spec-foundation`;
- baseline: `06a29ddc9ef526ff8a0df3dc2dc52e7918bd02c7` — `docs: restore Phase 6 roadmap historical detail after Models and Providers closure`;
- main: `bc1ec59e5e79ccbaba291e2984b0eb7d5e54129c`;
- PR #2 preflight: open / Draft / unmerged / base `main` / head baseline / `auto_merge=null`;
- branch and main README: exact `# cmdr`, blob `901c74cda52e28b5ff7fc425ddf28ef89f3ad875`;
- roadmap repair baseline: **PASS**.

## Predecessor closures

- Tenant/Environment `CAP-SET-001..004`: **160/160 PASS**;
- Identity `CAP-SET-005..007`: **174/174 PASS**;
- Secrets & Connections `CAP-SET-008..009`: **172/172 PASS**;
- Models & Providers `CAP-SET-010..011`: **204/204 PASS**.

No predecessor capability contract is reopened.

## Capability structure

| Capability | Owner | Status / delivery | Sections | Mandatory tables | Meaningful GWT | Result |
|---|---|---|---:|---:|---:|---|
| `CAP-SET-012` | Platform Settings Product Lead | draft / defined / planned | 27 | 6 | 4 | PASS |
| `CAP-SET-013` | Platform Settings Product Lead | draft / defined / planned | 27 | 6 | 4 | PASS |
| **Total** | one capability owner | — | **54** | **12** | **8** | **PASS** |

The six mandatory non-empty capability tables in each file are S8 Inputs, S9 Objects Read, S10 Objects Created/Modified, S13 Automation/AI, S16 Outputs and S17 Transitions/Handoffs.

## Namespace / object / screen / permission result

- exact allocation: `CAP-SET-012`, `CAP-SET-013`;
- `CAP-SET-014+`: **0 allocated / 0 reserved**;
- duplicate/recycled IDs: **0**;
- canonical objects reused: `OBJ-DATA_SOURCE`, `OBJ-PARSER`; sourced references may include Integration/Secret Reference;
- new canonical objects: **0**;
- active Screen Register remains **56**; reused `SET-SRC-001`, `SET-HLT-001`, `SET-AUD-001`, `SET-SEC-001`;
- new Screen IDs: **0**;
- existing permissions reused: `perm.platform-settings.data-source.read/manage`, `perm.platform-settings.parser.read/manage`, existing `perm.settings.source.read/manage` aliases;
- new Permission IDs: **0**; no bulk normalization.

## Data Source semantics

Canonical Data Source states remain exactly `configured`, `active`, `degraded`, `disabled`; Tenant remains mandatory. CAP-SET-012 owns administrative lifecycle/configuration, sourced scope/freshness, safe administrative disable, strictly local deterministic validation, sourced health projection and provenance.

It does not own acquisition, collection runtime, ingestion engine, connector runtime, external health probe, event-processing engine, parser runtime or storage. Health projection is source-attributed observation, not health-probe execution.

## Parser semantics

Canonical Parser states remain exactly `draft`, `testing`, `active`, `degraded`, `retired`; Tenant remains mandatory. CAP-SET-013 owns administrative lifecycle/versioning, transformation contract metadata, input dialect/output-schema reference, fixtures, error/quality metadata, local configuration validation, source-backed administrative rollback and provenance.

Parser Contract dimensions remain Input dialect / Output schema / Version / Fixtures / Errors / Quality. Schema format, initial schema version and SLO/limits remain unresolved. No ECS, OCSF, CIM, OpenTelemetry or other unsourced standard is selected.

Parser runtime, sandbox/plugin execution, normalization, stream processing, schema-registry implementation, automatic selection, source assignment, precedence and fallback execution are outside the lot.

## Source ↔ Parser relation lock

The canonical Data Source and Parser objects each state no direct upstream/downstream canonical object. This lot creates **no** `ParserAssignment`, `SourceParserAssignment`, `ParserCompatibility`, `ParserRoute`, `ParserSelection`, `ParserFallback`, `ParserPrecedence` or `SourceParserRelation`.

## Security / action / AI boundaries

Security retains RBAC, ABAC, tenant-first resolution, server enforcement, Decision Authority distinct from CRUD, SoD, step-up and audited deny/cross-scope. Read does not imply export/execute/approve; Settings `.manage` does not create runtime execution authority.

Class 0 = observation; Class 1 = strictly local deterministic no-effect validation; Class 2 = reversible/versioned Settings administration with `OPEN-013` unresolved. `Test Source` / `Test Parser` are not automatically Class 1 when they invoke external/runtime effects; Settings owns request/preconditions/handoff/observed-result projection only.

AI remains optional. It may explain sourced configuration, metadata, fixtures and validation errors, but may not enable ingestion, execute parser/runtime, fabricate health/validation, invent schema fields/relations, expose secrets, bypass authorization, silently mutate production configuration or resolve OPEN decisions.

## Requirements / OPEN / counters

Requirements remain **122 = 99 conform / 20 partial / 3 absent / 0 contradictory**. No Requirement ID is created/deleted and no state is changed merely because the capability definition exists. Where individual status is not exposed by the active matrix, no status is fabricated.

OPEN remains **18**; zero created, zero closed. `OPEN-008` and `OPEN-013` are directly relevant and remain open. `OPEN-011/012/014/015/017/018/019` remain implementation dependencies where applicable. Parser/Event/Health contract questions for schema format/version and SLO/limits remain unresolved.

Functional BUILD content:
- global: **497 capabilities / 495 defined / 2 proposed / 497 planned / 13,419 numbered sections / 2,982 mandatory tables**;
- Settings: **13 capabilities / 351 sections / 78 mandatory tables**;
- Sources & Parsers: **2 / 54 / 12 / 8 GWT**.

## Mandatory quality-gate matrix

| Family | Count | Build result |
|---|---:|---|
| A Git/baseline/PR/main/README | 13 | PASS |
| B repaired roadmap baseline | 16 | PASS |
| C predecessor closures | 8 | PASS |
| D source corpus | 10 | PASS |
| E namespace | 6 | PASS |
| F two structural contracts | 18 | PASS |
| G ownership | 7 | PASS |
| H Data Source semantics | 10 | PASS |
| I Parser semantics | 12 | PASS |
| J no Source→Parser relation | 6 | PASS |
| K object boundaries | 8 | PASS |
| L Secrets & Connections | 5 | PASS |
| M Models & Providers | 4 | PASS |
| N Investigate/Collection | 5 | PASS |
| O Command | 3 | PASS |
| P Endpoint | 4 | PASS |
| Q Studio | 4 | PASS |
| R Detection Engineering | 3 | PASS |
| S Threat Intelligence | 3 | PASS |
| T Security | 8 | PASS |
| U tenant/environment | 5 | PASS |
| V runtime separation | 8 | PASS |
| W schema/Event/Parser constraints | 9 | PASS |
| X health projection | 6 | PASS |
| Y screens | 6 | PASS |
| Z permissions | 8 | PASS |
| AA action classification | 7 | PASS |
| AB AI | 8 | PASS |
| AC audit/provenance | 8 | PASS |
| AD Requirements | 20 | PASS |
| AE all OPEN | 18 | PASS |
| AF terminology/migration | 8 | PASS |
| AG registers | 6 | PASS |
| AH counts | 7 | PASS |
| AI roadmap preservation | 7 | PASS at candidate BUILD after semantic audit |
| AJ local diff/non-regression | 8 | PASS at candidate BUILD after diff audit |
| AK publication | 4 | PENDING-REMOTE |
| AL remote verification | 6 | PENDING-REMOTE |
| AM CI/status/check/workflow | 5 | PENDING-REMOTE |
| AN documentary closure | 9 | PENDING-REMOTE |
| **Total** | **316** | **292 PASS / 24 PENDING-REMOTE / 0 FAIL** |

## BUILD verdict

**292 PASS / 24 PENDING-REMOTE / 0 FAIL**.

Final `316/316` and `PASS AFTER POST-PUBLICATION VERIFICATION` are forbidden until the candidate BUILD is published by non-forced fast-forward, remote ancestry/content/PR/main/README/namespace are re-read, CI/status/check/workflow applicability is actually inspected, and the one allowed documentary closure commit is published and reverified.

Platform Settings Capability Specification, Delivery Roadmap Phase 6 Capability Specification, Global Capability Specification and repository maturity remain **PARTIAL**.
