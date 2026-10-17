# E10-INV-110B — subprocess diagnostic confidentiality

## Change and reason

The first-party Timeline benchmark runner formerly embedded stderr (or stdout on an empty stderr) verbatim into Go errors on non-zero child exits. Go build, test helper and runtime probe output can contain machine-local paths or unexpected source metadata. The CI kernel reports tool errors; forwarding raw child content risks recording those values in GitHub Actions logs.

The bounded existing `runTimelinePerformanceProcess` now returns a wrapped `exec.ExitError` with a fixed generic description and **without child stdout/stderr** on failure. Error unwrapping still reports the process exit status, while context cancellation/deadline and output-size failure retain their distinct error branches. No extra library, external provider or new data model is introduced. This is simpler than inventing a redaction regex, since arbitrary diagnostics cannot be reliably sanitized.

A subprocess regression test emits explicit distinctive markers to **both** streams, fails, and asserts neither marker appears in the returned error, and no result is returned. The prior 64 KiB per-stream bounded buffers remain. The success path still returns bounded structured JSON and existing identity validators reject invalid observations.

## Limits and following actions

This patch only covers the Timeline helper and its Go error return. It cannot prevent a child from writing directly to an external system, nor guarantee confidentiality of GitHub Actions command tracing, other runners or the separately uploaded synthetic JSON measurements. It is not new evidence for a deployed ACL provider or production SLO. Re-check exact-SHA kernel/Timeline CI and obtain independent L2 conclusions before promoting B to VERIFIED. PR #36 remains the development vehicle; OPEN-013 remains unresolved.
