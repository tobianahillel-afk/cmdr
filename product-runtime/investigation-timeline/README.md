# Investigation Timeline — bounded read-only projection core

This module now contains an **initial first-party Go projection** for CAP-INV-110, not a complete Timeline Engine, source authorization service, backend provider or production-ready UI.

`Project` accepts caller-resolved Case/Shared Timeline Entry snapshots, enforces explicit tenant/environment + Case and Timeline Entry permission outcomes, masks denied/cross-tenant sources, preserves observed/inferred kind, source references, provenance, version and supersession, and sorts known timestamps deterministically. Unknown or clock-uncertain time does not become an observed chronology. Inputs are not mutated. The caller MUST obtain authoritative server-side permission decisions separately; caller-provided permission flags are **not** proof of production authorization.

Executable tests cover the 18 predeclared synthetic contract scenarios and focused negative cases. Limits (5,000 entries and 8 MiB string input budget) are provisional engineering protection, not a canonical limit or production SLO. No persistence, Shared Timeline Entry mutation, replay/export, external runtime dependencies, network or final UI exists here. `OPEN-013` is unresolved.

**Performance is NOT measured yet for Timeline**: the previously registered projection/navigation benchmark and deep-resource handlers remain fail-closed preimplementation handlers. They must not be advertised as passing until actual representative executable probes are introduced and tested. Security coverage and module-specific race/vet are exercised by CI, but no production deployment or full CAP-INV-110 completion is claimed.

Security hardening: `Supersedes` is exposed only when the referenced prior version is also visible and authorized in the bounded snapshot; unavailable/denied predecessor IDs are redacted and status becomes `partial`. A timezone's explicit RFC3339 offset must match the timestamp offset before chronology is considered known. Unknown/named zones conservatively remain uncertain. These are projection policies, not authorization-provider proofs.

Security hardening: trim scope before wildcard checks; uniqueness is enforced among accessible source entries only, preventing denied/cross-tenant ID collisions from denying legitimate projections. Visible supersession links require strictly increasing versions (so cycles cannot be represented); inaccessible predecessor IDs remain redacted. This bounded, caller-supplied projection still is **not** an authoritative authorization provider.
