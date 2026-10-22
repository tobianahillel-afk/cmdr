---
id: investigate-cloud-analysis-source-migration
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements: [REQ-PROD-006, REQ-PROD-009, REQ-PROD-012]
open_decisions: [OPEN-012]
---
# Cloud Analysis source migration

## Audit result
No prior canonical Investigate Cloud Analysis module, Cloud capability family or Cloud subphase identifier was found. Historical Cloud needs are distributed across Platform Settings, Endpoint telemetry, Event Search, Network Forensics, Security, Shared capabilities and roadmap placeholders.

## Migration rule
- keep owner documents active;
- create no competing provider, connector, secret, Entity, Graph, Timeline, Detection, Decision or Tool source;
- reference historical generic Cloud material through its owner;
- do not deprecate Settings, Endpoint, Shared, Command, Govern, Studio or forensic sources;
- create no Mobile Forensics source.

## Canonical adaptation
`Phase 4B.4A` becomes the Cloud Analysis identifier because no existing identifier conflicts. `CAP-INV-601..618` is the first free family after `CAP-INV-5xx`.
