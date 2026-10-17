# E10-INV-110B — bind negative-test evidence to declared runtime identity

The existing Go security gate already requires `go test -json` run/pass events from one package. However, the verifier accepted **any** package from `./...` providing matching subtest names, even an unrelated package added under the same module. This permits false attribution of authorization and tenant-isolation evidence.

This increment reuses the **already registered** first-party `runtimeSecurityAdapter.ModuleIdentity` value, and requires each selected negative test's `Package` field to match that exact module identity. No new runner, external API, database, dependency or canonical Product Spec changes. Synthetic tests reject a complete suite from another package, a split across packages, and missing expected-module identity; a real Timeline integration test verifies the correct identity passes while a spoofed identity fails.

The previous test-selection guarantees (run and pass required, parent-only and deeper nested tests rejected, every alternative required) remain. This gate proves execution of named in-module test cases only: **caller-supplied permission flags are not real external authorization-provider enforcement**. `OPEN-013` stays open. B stays READY pending CI, L2 and remaining closure criteria.

CI run references and L2 findings must be appended only after verifying exact published SHA.