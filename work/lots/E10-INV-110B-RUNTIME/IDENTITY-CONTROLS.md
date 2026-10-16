# E10-INV-110B-RUNTIME — first security prerequisite: measurement identity binding

Independent source audit identified that the generic benchmark catalog previously allowed one known handler to masquerade as another workload by editing JSON metadata, despite the allowlist. In particular a Case Queue probe could be relabeled as an Investigation Timeline probe. That is unacceptable evidence substitution even for preimplementation targets.

This first B increment binds the two CAP-INV-110 workload identities to their exact target/boundary/owner and compiled benchmark/deep handler identities, with three adversarial Go tests covering mismatches. Existing performance adapters and Go standard library reused; no implementation of final Timeline runtime, storage index or backend, and no performance/security results invented. The new `go test` test file is control-plane code, not product runtime.

Next task within B remains to implement and test the actual bounded read-only projection using authentic fixtures and measured coverage before activating either registered probe. This note is not a VERIFIED handoff; CI and L2 review must examine the exact published SHA.
