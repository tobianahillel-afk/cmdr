# Go SAST and SCA

E4-SEC-002B covers the currently implemented Go engineering module only. It does not claim JavaScript, C/C++, container or product-runtime coverage.

## Toolchain

The engineering module uses Go 1.26.8. External security tools are invoked with exact module versions using `go run module@version`; they are deliberately absent from the CMDR module dependency graph.

Reviewed tools:

- gosec v2.28.0, commit `9e75c0576c9878035d4221392108d458abe10fc3`, Apache-2.0;
- govulncheck / `golang.org/x/vuln` v1.8.0, commit `709015412431dd2b5b28a53c06c70bc02d49074c`, BSD-3-Clause.

The machine-readable source of truth is `development-tools.json`.

## SAST policy

gosec runs locally over `tools/cmdr-dev` with:

- medium-or-higher severity;
- medium-or-higher confidence;
- generated files excluded;
- suppression directives required to name rule IDs and carry a justification;
- JSON written to a temporary file;
- `-no-fail` so CMDR, not tool exit-code behavior, decides blocking.

CMDR evidence retains only rule ID, severity, confidence, CWE, repository-relative path and line. gosec's code snippet and descriptive details are intentionally discarded. Any medium/high filtered finding blocks.

No AI/autofix mode is enabled, so source is not sent to AI providers.

## SCA policy

govulncheck v1.8.0 runs in source mode with streaming JSON. CMDR distinguishes:

- module/package-level findings: evidence only;
- symbol-level reachable findings: actionable and blocking.

Actionable evidence must resolve to exact `Go / module / version` identity, plus OSV ID and package/symbol when available. Missing version identity fails closed.

govulncheck does not expose a normalized high/critical severity on its Finding object. The CMDR policy is intentionally stricter than a high/critical floor: **every reachable known vulnerability blocks**. This avoids inventing severity and keeps the evidence source-faithful.

## Network/data boundary

These tools are development-only. They never become CMDR product-runtime dependencies. gosec analyzes source locally. govulncheck analyzes source locally and may query the official Go vulnerability database using module identity/version metadata; it does not upload repository source code.

SEC-SAST-001 and SEC-SCA-001 remain `specified` until this implementation passes both push and pull-request CI.
