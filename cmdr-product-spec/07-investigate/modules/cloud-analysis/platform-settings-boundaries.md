---
id: investigate-cloud-analysis-settings-boundaries
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements: [REQ-PROD-014, REQ-PROD-017, REQ-PROD-055]
open_decisions: [OPEN-008, OPEN-012]
---
# Platform Settings boundaries

Platform Settings owns providers, connectors, credentials, secrets, configured organizations, tenants, accounts, subscriptions, projects, ingestion, schemas, parsers, health, retention, storage, access policies and administrative configuration.

Cloud Analysis consumes authorized projections only:
- configured versus accessible status;
- source and connector identity;
- scope units and provider terminology;
- coverage, health, schema/parser and retention context;
- policy and restriction projections;
- secret references, never values by default.

Cloud Analysis cannot create, update, test, repair or delete a connector; grant access; rotate credentials; change retention; add a provider; or mark coverage complete. Gaps become proposals or handoffs to Settings.
