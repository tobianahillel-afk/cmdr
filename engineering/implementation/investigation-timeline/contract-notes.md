# E10-INV-110A-CONTRACT — first source-backed executable-spec increment

Canonical CAP-INV-110 is *draft defined/planned*. This contract is an engineering-only bounded oracle for later runtime implementation; **it is not the final Shared Timeline Entry schema and not executable product functionality**.

Permission provenance: Case read (`perm.investigate.case.read`) and Shared Timeline Entry read (`perm.shared-capabilities.timeline-entry.read`) are baseline obligations. Evidence, Telemetry Event, Govern Decision and Command Incident reads are CONDITIONAL on an entry referring to such a source; permissions are never unioned into one universal role or inherited across linked objects. Denied or unresolved source access must never expose its identifier or details. Runtime must use authoritative permission evidence; fixture caller-provided `allow` is untrusted and insufficient for production authorization.

The 18 synthetic fixture rows cover same-time observed/inferred, stable ties, missing/wildcard scope, missing environment, denied and unknown access, hidden/cross-tenant sources, missing source, timezone uncertainty, clock skew, immutable supersession, incomplete provenance, duplicate IDs, preserving inference, and empty display. Expected outcomes are *test oracles*, not observed runtime results.

Engineering-only provisional resource limits: up to 5,000 bounded entries and 8 MiB aggregate source input for a future pure local projection; reject oversized/unbounded input prior to allocation and use a stable O(n log n) sort. Limits need stress-test calibration during E10-INV-110B-RUNTIME, not claimed to be canonical contract constants or production SLOs.

The executable schema sanity tests **do not prove the Timeline runtime implements any case**. Full behavior, security, race and resource measurement remain for later lots. No runtime mutations, no model/agent/API/pay service, and no third-party runtime dependencies introduced.
