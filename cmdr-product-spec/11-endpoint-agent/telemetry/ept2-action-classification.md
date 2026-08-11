# Endpoint EPT-2 — Action Classification

| Class | EPT-2 actions | Boundary |
|---|---|---|
| 0 | inspect source, observation, normalization, health, capability declaration, provenance | read-only technical facts |
| 1 | deterministic normalization/classification, duplicate candidate assessment, freshness/gap/availability calculation, masking | derived state only |
| 2 | bounded telemetry refresh/status request when supported | source/policy administration remains Settings-owned |
| 3 | none | response action excluded |
| 4 | none | high-risk response excluded |

EPT-2 introduces no Detection, acquisition, Live Response, containment, rollback, update or implementation action.
