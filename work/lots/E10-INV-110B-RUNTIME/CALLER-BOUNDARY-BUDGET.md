# E10-INV-110B — second independent boundary: caller-snapshot access adapter

The owner-derived snapshot path is protected by the separate immutable commit `0f4cd39e0e20fae81273c515d4e437b8468c0f9e` (which is being independently reviewed). This commit completes the same **input-envelope limit** for the older `ProjectWithAccessChecker` entrypoint, which otherwise still forwarded a caller-provided over-8-MiB SourceRef into an authoritative source access request before the final `Project` budget denied it.

Reuse first: the existing Go `MaxInputBytes` engineering threshold, string `len`, existing Case/Timeline/source access obligations, and standard tests; no added dependency or reinvented policy machinery. It refuses scopes/CaseID exceeding the engineering budget even for an empty entry list and refuses overlong source references before the provider sees them. It does not claim the legacy caller-snapshot entrypoint is a production provider. Its production replacement remains the independently authorized owner-resolver adapter with real server integration still outstanding.

Scope: no changes to Product Spec, owner boundaries, navigation, release status, performance benchmarks, or `OPEN-013`. CI and L2 approval must be separately verified per exact SHA.
