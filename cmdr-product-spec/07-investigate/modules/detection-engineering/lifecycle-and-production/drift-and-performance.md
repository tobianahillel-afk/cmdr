---
id: investigate-detection-drift-performance
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
---
# Drift and performance

CAP-INV-431 assesses source, schema, field, mapping, normalization, enrichment, runtime, policy, retention, quality, volume and dependency changes without assuming failure. CAP-INV-432 assesses execution cost and resource impact without equating latency, cost, accuracy or business value.

Drift Assessment and Performance Assessment consume Settings, Endpoint, Command and Shared projections. They create no parser, runtime configuration, threshold, capacity change or active tuning. Revalidation, a new Draft, a Settings request, a capacity request, a Tuning Proposal or a Rollback Plan may be prepared with full provenance.

`drift-candidate` is not incompatibility; `schema change` is not certain rule failure; high performance is not high quality; high cost is not automatically low value. Definitions, freshness, missing periods and per-target limitations remain visible.
