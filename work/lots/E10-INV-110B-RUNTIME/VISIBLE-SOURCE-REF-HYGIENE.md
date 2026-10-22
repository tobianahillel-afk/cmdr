# E10-INV-110B — reject invisible authorized SourceRef values in pure projection

## Verified gap addressed

The server-side `ProjectWithAccessChecker` and `ProjectFromOwnerSnapshots` paths both use the first-party `boundedSourceSelector` check, which invokes `invalidScope` before passing a SourceRef to the external authorization authority. The existing pure `Project` path previously rejected only an empty SourceRef, so a control character, bidi format mark or zero-width format mark could appear in its JSON projection output when callers supplied a visible Entry with `SourcePermission=allow`. The pure function does **not** confer authorization, but this was an inconsistent identifier-validation guarantee.

Minimal fix: use the **same existing** `invalidScope` helper on an authorized visible SourceRef, with no new dependency, Unicode library, encoding, storage or schema change. Do not sanitize or rewrite forensic provenance or timestamps. Unavailable/denied/foreign source fields remain opaque: their malformed references cannot cause a visible projection to deny. The source authorization provider's independent Case/Timeline/source checks stay authoritative. Tests cover missing, wildcard, NUL/newline/zero-width/bidi references and hidden-source noninterference.

Technology: Go standard library `strings` / `unicode` already in the module versus adding a new sanitizer/dependency: reuse is the simpler, safer default. This is a narrow validation correction and does not complete an external identity/permission provider, navigation or `OPEN-013`.

Evidence status: compare exact GitHub Actions run results for this published commit; do not reuse previous source SHA CI or latency evidence as if measured here. L2 results remain unapproved unless an independent report can actually be retrieved and verified.
