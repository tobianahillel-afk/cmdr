---
id: endpoint-ept3-migration-source-map
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-11
source-of-truth: canonical
---
# EPT-3 Migration and Source Map

## Detection sources
Eight active canonical source documents remain preserved: README, Behavioral Detection, Detection Suppression, Detection Update, Local Correlation, Local Detection Engine, Model Inference and Rule Model. None is deleted or deprecated. `detection-update.md` remains future EPT-6/update evidence and is not promoted into EPT-3 behavior.

## Investigation sources
Eight active canonical source documents remain preserved: README, Host Inspection, Host Timeline, Modules and Drivers, Network Connections, Persistence, Process Tree and User Sessions. None is deleted or deprecated.

## Normative migration
`CAP-EPT-031..046` become the normative capability layer for EPT-3 while the 16 source files remain canonical module/reference sources. No competing physical schema or runtime architecture is introduced.

Potential stale terminology is interpreted through current ownership registers: canonical Command `Detection`/`Signal` remain Command-owned; references to Evidence links do not allow Endpoint to create Evidence; “live host inspection” remains read-only/contextual and escalates to future Collection when new data is required.

Migration strategy is additive only.