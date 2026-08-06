---
id: investigate-threat-intelligence-workflows
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
---
# Workflows

## Canonical foundations flow
Case/Finding/Hunt/Incident/technical or Detection handoff → CAP-INV-501 Intake → CAP-INV-502 Requirement → CAP-INV-503 Knowledge Project → CAP-INV-504 Source Catalog → CAP-INV-505 Reliability/Credibility → CAP-INV-506 Material Intake/Normalization → CAP-INV-507..512 candidates and knowledge → CAP-INV-513 Sightings → CAP-INV-514 Relationships → CAP-INV-515 Confidence/Contradictions → CAP-INV-516 Versioning/Deduplication → CAP-INV-517 Lifecycle → CAP-INV-518 Analysis Handoff.

## Required returns
Every flow preserves origin, source owner, tenant, environment, versions, permissions, markings, restrictions, errors, limitations and human disposition. Returned, disputed, rejected and partial states preserve valid work.

## Future-only destinations
CAP-INV-518 can prepare context for future 4B.3B.2, Case, Hunt, Detection Engineering or Analysis Workbench. It does not create a Report, publish externally, activate an Indicator/watchlist or execute a response.
