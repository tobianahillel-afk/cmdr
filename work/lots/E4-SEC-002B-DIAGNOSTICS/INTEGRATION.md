# E4-SEC-002B-DIAGNOSTICS — SAST gosec processing errors remain blocking and become locatable

## Trigger and evidence

Engineering-kernel at exact PR #154 SHA `fb2988df301d3d6e531d3c85f2f9dc6bbe22557f`, run [39346089085](https://github.com/tobianahillel-afk/cmdr/actions/runs/39346089085), FAILED at `CHK-SAST-GO`: `gosec reported 1 processing error(s); details withheld from CI evidence`. The concurrent [command-work-queue-runtime 39346089080](https://github.com/tobianahillel-afk/cmdr/actions/runs/39346089080) passed Linux/macOS/Windows. PR #153 had a distinct `CHK-DEEP-PERFORMANCE` failure; the #154 probe addressed that historical deficiency. Do not assume the current SAST failure comes from newly added Command code: `runGoSAST` scans all registered Go module roots in a deterministic sequence and failed without identifying the offending module.

## Owner and targeted change

Canonical business capability: **none**, this is an existing engineering security control. Engineering owner: existing `E4-SEC-002B` (Go SAST/SCA), child task `E4-SEC-002B-DIAGNOSTICS`. Code owner: `tools/cmdr-dev/security_static.go` and tests, no changes to `product-runtime/` or to canonical Command CAP-CMD-101 definitions. Git parent is frozen #154 HEAD, but this engineering ownership **does not inherit** Command's business parent.

Changes:
- When a gosec scan fails, append the reviewed **module root as repository-relative metadata**.
- For gosec `Golang errors`, continue to fail immediately; derive a deterministic list of **up to four existing, normalized, in-repository .go file paths** if supplied by gosec. Suppress compiler output, analyzer messages, source snippets, external paths, package names, and untrusted arbitrary values. No safe path -> report processing-error count and `source paths unavailable` while still returning FAILURE.
- Keep the exact pinned gosec module/version, scanners, scan roots, severity/confidence filters, processing-error failure, vulnerability gate, and all coverage/performance/security thresholds untouched.

## Verification and next gate

Behavioral regression test creates a real local repo-relative Go path and a synthetic gosec JSON processing error containing sensitive compiler content plus off-repository and nonexistent paths. It must fail and reveal only the known safe path and count. A pathless external error must still fail with no attacker-controlled text. The exact branch/PR SHA's engineering-kernel can now reveal which registered module (and safe Go file, if supplied) is preventing SAST; **this is a diagnostic unblocking step, not a waiver or proof that SAST passed**.

Material next action is to inspect the exposed module and fix its actual parser/package-loading cause, not fabricate green checks. Do not merge #154/this child or claim product acceptance until CI and required upstream approvals pass. No paid API, new dependency or product capability.
