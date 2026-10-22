---
id: capability-register-settings-sources-and-parsers
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-14
source-of-truth: registry
---
# Capability Register — Settings Sources & Parsers

Execution lot: **Sources & Parsers — Data Source Administration and Parser Transformation Administration** under Delivery Roadmap Phase 6 — Platform Scale.

| Capability ID | Name | Owner product/module | Status | Delivery status | Delivery mode | Canonical file | Primary roles | Primary objects | Consumers | Requirements | OPEN | Dependencies | Supersedes | Review |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| `CAP-SET-012` | Data Source Administrative Lifecycle, Scope, Freshness and Health Projection | Platform Settings / Sources & Parsers | draft | defined | planned | `../../10-platform-settings/capabilities/cap-set-012-data-source-administrative-lifecycle-scope-freshness-and-health-projection.md` | Platform Settings administrator; authorized operator | `OBJ-DATA_SOURCE`; sourced Integration/Secret Reference refs only | Settings Health/Audit; authorized Command/Investigate/Shared consumers; runtime owner via handoff | REQ-PROD-003,004,005,006,008,009,010,011,012; REQ-AI-001,002,003,004,007,008,009,010,011 | OPEN-008, OPEN-013 | Data Source object; Health/Event contracts; Security; Administrative Audit | none | 2026-08-14 |
| `CAP-SET-013` | Parser Administrative Lifecycle, Versioned Transformation, Fixtures and Validation | Platform Settings / Sources & Parsers | draft | defined | planned | `../../10-platform-settings/capabilities/cap-set-013-parser-administrative-lifecycle-versioned-transformation-fixtures-and-validation.md` | Platform Settings administrator; authorized operator | `OBJ-PARSER` | Settings Audit; authorized downstream consumers; parser runtime owner via handoff | REQ-PROD-003,004,005,006,008,009,010,011,012; REQ-AI-001,002,003,004,007,008,009,010,011 | OPEN-008, OPEN-013 | Parser/Event contracts; Security; Administrative Audit | none | 2026-08-14 |

## Structural evidence

- capabilities: **2**;
- numbered sections: **54**;
- mandatory tables: **12**;
- meaningful GWT: **8**;
- owner conflicts: **0**;
- duplicate/recycled IDs: **0**;
- new Screen IDs: **0**;
- new Permission IDs: **0**;
- new canonical objects: **0**.

## Boundary evidence

`Data Source` and `Parser` are distinct canonical objects and currently have no canonical direct upstream/downstream relation. This shard does not create assignment, compatibility, routing, selection, fallback or precedence. Settings owns administrative lifecycle/configuration and sourced projections only; acquisition, collection, ingestion, connectors, probes, parser execution, normalization, schema-registry implementation and storage remain outside this lot absent a separately sourced owner/implementation.

Requirements remain **122 = 99 conform / 20 partial / 3 absent / 0 contradictory**. OPEN remains **18**. `CAP-SET-014+` remains unallocated and unreserved.
