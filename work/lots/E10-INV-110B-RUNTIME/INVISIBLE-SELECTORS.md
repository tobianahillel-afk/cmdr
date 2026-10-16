# E10-INV-110B — invisible selector hygiene

This security increment closes a trust-boundary ambiguity: an embedded NUL, newline, bidi format mark or zero-width Unicode format mark previously passed the first-party `invalidScope` guard and could reach Case/Timeline/source access providers, their loggers or policy stores. The first-party Go standard library `unicode.IsControl` and `unicode.Cf` is reused instead of adding a sanitizer, external provider, persistence or normalization service. Wildcard and empty rejection remains unchanged.

Tests prove that owner-resolved tenant/environment/Case/Timeline selectors are denied **before any access-provider call**; an invalid owner-provided SourceRef is denied before checking external source permission; legacy caller-snapshot selectors cannot be forwarded to the respective authority; and malformed visible IDs are rejected by pure projection while foreign hidden IDs stay opaque. The legacy behavior for malformed requested IDs remains filtered rather than a full failure.

Not covered: URL-encoded controls, Unicode normalization policy across a real provider, external Case membership correctness, actual production ACLs, navigation, storage, or `OPEN-013`. Product Spec remains immutable. Do not claim CI success or independent L2 approval until exact-SHA GitHub Actions and the distinct review are checked.
