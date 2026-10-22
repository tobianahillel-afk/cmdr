---
id: endpoint-security-ept6-capability-links
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-12
source-of-truth: canonical
---
# Security — EPT-6 Capability Links

The historical Endpoint Security documents remain source/background material. Normative EPT-6 security contracts are:
- `CAP-EPT-094` — Endpoint Self-Protection and Anti-Tamper Boundary;
- `CAP-EPT-095` — Local Privilege, Service/Process Protection and Security Context;
- `CAP-EPT-096` — Sensitive Material, Secret Reference and Credential Handling Boundary;
- `CAP-EPT-097` — Local Audit Event, Security Audit and Provenance Semantics;
- `CAP-EPT-098` — Endpoint Security-State, Integrity and Consumer Handoff.

Global Security retains Permission Model, privacy/minimization, audit/integrity, tenant isolation and secrets/key policy. Platform Settings retains Secret/credential administration. Endpoint owns only local technical observations/handling facts.

Tamper candidate != compromise verdict; privilege != permission; Local Audit Event != Shared Trace; append-only != cryptographic proof; Secret Reference != raw secret; Endpoint security state != Finding/Incident/Result.