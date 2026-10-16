# E10-INV-110B — negative-test evidence must be from one package and actual execution

This incremental correction strengthens the existing Go standard-library security test gate. Previously each regex alternative merely needed a `pass` event somewhere in `go test -json`, so distinct Go packages could collectively satisfy the list without any single package having executed all the expected denial/tenant tests. A `pass` event also could count without a corresponding `run`.

The verifier now groups events by `Package`, pairs ordered `run` + `pass` for exact named test cases, and requires **one package** to have all alternatives. Package-level pass, parent-only pass, omitted Package, skipped cases, missing alternatives, deeper nested cases and incomplete cross-package unions cannot satisfy the gate. Test fixtures attempt each bypass.

**Reuse-first comparison:** Go `test -json` and `encoding/json` + `regexp` are already in CMDR; a third-party runner adds unnecessary dependencies and more trust surface. This is a bounded CI evidence-hardening change, **not** a real Case/Timeline/external-source authorization-provider integration. Product Spec, runtime behavior and `OPEN-013` remain unchanged.

**Status:** actual CI verdict and independent L2 must be checked against the published SHA before asserting success. Existing passed security cases are not retroactively proof of end-to-end provider behavior.
