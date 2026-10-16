# E10-INV-110E — bounded, ledger-compatible CAP-INV-110 handoff

The functional deliverable is an executable first-party read-only Case investigation reference journey: discover authorized Case -> bounded Timeline -> navigate to independently reauthorized typed-source preview -> return and recheck Case access, with negative cross-tenant tests. The **D-CLOSURE exact-code** SHA `06d1d47572b9fbea1325ad71309a895d960e1cbd` has four successful push/PR CI workflows (two engineering and two Timeline), Linux race/vet plus macOS/Windows Timeline. B and C were previously VERIFIED; see the lineage in their PROGRESS.json.

Machine-readable bounded handoff: `engineering/implementation/investigation-timeline/bounded-handoff.json`. Its wording mirrors the existing Case Queue handoff and explicitly **does not add an entry to `engineering/implementation/ledger.json`**; a separate scaling or authorization work unit owns that action. The temporary E status is IN_REVIEW until exact-HEAD engineering CI on these handoff files confirms the schema, then the handoff may be finalized.

What this is **NOT**: integrated production SIEM or EDR ingestion, durable tenant repository, production IAM, final browser screen, remote audit pipeline, Timeline Entry mutation, replay/export, full CAP-INV-110 completion, request-level p95 or SLO. `OPEN-013` remains unresolved. The new HTTP integration and navigation ns/op measurements are local synthetic reference evidence only.

Reuse-first technical choice: Go standard-library handlers, injected owner interfaces, existing first-party Case/source permission checks, and existing engineering ledger shape, not a new schema, component framework, authorization package or arbitrary storage vendor.

**Next functional milestone (separately authorize/scope):** replace the copied synthetic Case/Timeline/source owner fixture with a real host-provided, tenant-isolated, server-side authenticated and auditable owner implementation, prove revocation and source access across the complete HTTP journey, then decide persistence/indexing based on representative end-to-end workloads. Neither this document nor E changes an open canonical product decision.
