---
id: quality-index-platform-scale-secrets-and-connections
domain: 16-quality-and-validation
status: draft
owner: Quality Engineering Lead
updated: 2026-08-13
source-of-truth: canonical
---
# Quality Index — Platform Scale Secrets & Connections

## Current lot
**Secrets & Connections — Integration and Secret Reference Administration**

- capability contracts: `CAP-SET-008`, `CAP-SET-009`;
- structure: **2 capabilities / 54 numbered sections / 12 mandatory tables / 12 meaningful GWT**;
- build quality: **166 PASS / 6 PENDING-REMOTE / 0 FAIL**;
- final target: **172/172** only after actual remote post-publication verification;
- conformance report: `reports/platform-scale-secrets-connections-integration-secret-reference-capability-conformance.md`;
- validation status: `validation-status-platform-scale-secrets-and-connections.md`;
- roadmap: `../18-roadmap-and-releases/phase-6-platform-scale.md`.

## Hard invariants
- Platform Settings Product Lead is the unique capability owner;
- Tenant required; Environment source-dependent only;
- 0 new canonical objects;
- 0 new Screen IDs;
- 0 new Permission IDs;
- `Integration.capabilities` metadata != CMDR Capability / CAP-*;
- Connection/Connector/Credential/raw Secret objects not introduced;
- technical `Test Connection` executor not invented;
- Secret Reference rotation/revocation != underlying external secret/credential operation unless a canonical owner/outcome exists;
- Requirements remain 122 = 99/20/3/0;
- OPEN remains 18;
- documentary PASS != implementation.

## Maturity
At functional BUILD, Platform Settings Capability Specification, Delivery Roadmap Phase 6 Capability Specification, Global Capability Specification and repository maturity remain **PARTIAL**. Remote-dependent gates remain pending until publication and verification.
