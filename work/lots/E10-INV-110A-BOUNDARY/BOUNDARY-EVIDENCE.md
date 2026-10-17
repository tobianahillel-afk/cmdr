# E10-INV-110A-BOUNDARY — preimplementation controls

This lot registers exactly one `investigation-timeline-runtime` first-party boundary and atomic path ownership in `engineering/architecture/boundaries.json`. `product-runtime/investigation-timeline/README.md` is the sole Timeline product-runtime content; subsequent child manifests are allowed to use the owned path only within their own scope. Shared owns Timeline Entry and Timeline Engine, not this new Investigate projection.

Security obligations in `runtime-security-policy.json` require tenant isolation and source-permission negative tests. Current status is **preimplementation**: no test coverage or authorization/tenant runtime PASS is asserted.

Two performance targets and deep-resource targets are registered for future immutable projection and navigation. Both are `on-demand`; their compiled handlers intentionally return a `preimplementation` error, checked by `TestInvestigationTimelinePreimplementationMustNotEmitMeasurements`. No measured speed, memory improvement, or production SLO is asserted; later executable work must install real probes before changing this gate.

Technology choice: reuse existing Go standard library and CMDR authorization/snapshot patterns, with stable ordering and explicit provenance. No vendor index/event-sourcing engine or additional runtime dependency is needed for **bounded local projection**; storage and scale decisions require separately benchmarked evidence.

`OPEN-013` remains unresolved. Case mutation, analytic correction, source writes, export and final UI remain outside this lot. `cmdr-product-spec/**` and existing verified product runtimes remain unchanged.

**Verification status:** never infer VERIFIED from these notes. Promote only after exact-HEAD push and PR CI succeed and their IDs are recorded in `HANDOFF.json`/`PROGRESS.json`.
