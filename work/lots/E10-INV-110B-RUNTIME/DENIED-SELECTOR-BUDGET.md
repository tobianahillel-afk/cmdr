# E10-INV-110B — denied Timeline IDs cannot consume permitted source budget

## Scope and technical decision

Read-only security correction in both first-party Go authorization adapters. Reuse existing `AccessChecker`, `EntryResolver`, `boundedSourceSelector` and standard-library counters/maps, without a dependency, persistent engine, external provider or new domain owner. The full requested ID set is still validated for duplication and bounded by **8 MiB before ANY remote authority call**.

After all Timeline Entry permission decisions are obtained, the **linked-source authority** budget starts from tenant/environment/Case context and adds *only permitted Timeline Entry IDs*, then every permitted linked `SourceRef`. A denied/unknown/missing entry's requester ID is no longer charged against the source-specific budget. No denied entry is passed to the owner resolver or linked-source checker.

### Reproduction covered by regression tests

One authorized ID `x`, one denied selector with 3 MiB of `h`, and one permitted linked source reference with 6 MiB of `s`. Both identifiers together fit the globally bounded requester envelope. Previously the 3 MiB denied selector was incorrectly included in the source-provider cumulative budget; this changed an otherwise valid projection into a global `deny`. The new tests assert `filtered`, exactly one projected authorized entry, source authority called once for the permitted source, no denied owner resolution. Owner-resolved and legacy caller-snapshot adapter cases are covered independently.

## Proof boundary

No concrete provider authenticates the principal or supplies an atomic Case/Timeline/source ACL snapshot; interfaces and synthetic fixtures do not confer production readiness. Separate source permissions for Evidence, telemetry, Govern or Command remain delegated to owner-side integration, `OPEN-013` remains open, and C-NAV must stay blocked. The prior synthetic p95 is not benchmark evidence for this SHA. CI/L2 results must be checked against the published commit, not inferred from this note.
