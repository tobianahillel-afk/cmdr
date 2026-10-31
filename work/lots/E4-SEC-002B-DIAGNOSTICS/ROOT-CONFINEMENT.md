# E4-SEC-002B-DIAGNOSTICS — repository-confined SAST path inspection

## Evidence and separate ownership

This is engineering-security code and tests, owned by existing `E4-SEC-002B-DIAGNOSTICS` (parent `E4-SEC-002B`), **not** canonical Command CAP-CMD-101. It is grouped in the same cross-workstream **blocking-SAST recovery PR** as a separate, explicitly attributed patch to the Command-owned benchmark cleanup. The Git parent PR #155 SHA `7e199c5a7e5006d16bb230c68fd315cbf05655f2` is frozen.

On first child SHA `0cd99aaa9d9d870fcdd3809eb3526180cf90f0b6`, exact `engineering-kernel` workflow run [39352275146](https://github.com/tobianahillel-afk/cmdr/actions/runs/39352275146) completed **FAILURE** at mandatory `CHK-SAST-GO`: **three high-severity G703 findings**, two from new dynamic-path `os.Remove` calls in the Command temporary-probe cleanup and a third at `tools/cmdr-dev/security_static.go:293` from `os.Lstat(filepath.Join(root, reportedFile))`. This run no longer reports the old `syntax-or-parser` processing error; it is a new exact-SHA observation, **not proof that gosec was already globally green**.

## Design-level remediation

Instead of suppressing or ignoring G703, the code opens `os.OpenRoot(root)` once and obtains a `(*os.Root)` confined to the workspace. The already accepted `normalizeSecurityToolPath` still filters gosec error keys to normalized repository-local Go files; **`workspace.Lstat(filepath.FromSlash(path))`** checks the requested file without allowing intermediate symbolic links to escape the actual open root, and still rejects symlink final targets via `fi.Mode().IsRegular()`. No raw compiler message, attacker-controlled key, external path or source snippet is printed; processing errors remain fatal.

The separate Command owner uses `os.OpenRoot(repoRoot)` then `workspace.OpenRoot(randomTempDirName)`, unlinks ONLY the known compiled benchmark binary with `privateDir.Remove`, and removes the empty directory with `workspace.Remove`; an unexpected file or replacement symlink is never recursively removed and cannot escape the workspace. No new third-party product runtime dependencies, no scanner version change, no modified severity gates, no `#nosec` suppressions, no paid API.

**Acceptance:** the real exact-HEAD gosec scan must reach SUCCESS with no high findings AND no processing errors. If still failing, do not bypass security checks. Existing gosec error-redaction regressions + real Command probe cleanup tests continue. This does NOT prove production IAM or Command queue completion.

## Outstanding

Run exact-HEAD `engineering-kernel` after this code, confirm actual toolchain, tests and SAST status. Until verified, **DRAFT / NO MERGE**. Ancestor PRs #36/#153/#154/#155 remain draft and require independent check/review completion. L2 remains risk-based; ordinary code cleanup does not merit an automated L2 per SHA, but any material change to scanner trust model must receive targeted independent review.
