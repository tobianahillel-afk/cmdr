---
id: endpoint-ept1-action-classification
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-11
source-of-truth: canonical
open_decisions: [OPEN-013]
---
# EPT-1 Action Classification

| Class | EPT-1 meaning | Representative actions | Boundary |
|---:|---|---|---|
| 0 | observation | inspect Agent identity, registration, platform, inventory, health, heartbeat, state, capabilities, provenance | read under source permission |
| 1 | no-effect deterministic assessment | identity consistency, inventory comparison, freshness calculation, compatibility assessment, health normalization, heartbeat/state derivation, capability availability normalization | no production effect or authority |
| 2 | bounded reversible technical request/ack | request registration/enrollment handoff, acknowledge/refresh local configuration state, request bounded inventory/health/self-check refresh where source supports it | no Fleet/Policy admin; OPEN-013 remains open |
| 3 | production response effect | none created by EPT-1 | future EPT-4/5 + Govern/technical authority |
| 4 | destructive/irreversible effect | none created by EPT-1 | outside EPT-1 |

EPT-1 creates no Fleet administration, Endpoint Policy assignment or production-response primitive. Class 2 never bypasses Settings ownership or Govern/Security authority.