# E10-INV-120A — select an authorized Case at offline startup

A host operating the authenticated loopback analyst demo can supply an explicit Case from the operator-selected read-only JSON snapshot:

```sh
go run ./cmd/case-timeline-demo -snapshot="$(pwd)/cmd/case-timeline-demo/testdata/example-owners.json" -principal=investigator-42 -case=case-incident-42
```

Without `-case`, the existing first readable Case selection is unchanged. No Case ID may be provided as an authorization decision through HTTP: the CLI resolves the Case **before** the local listener binds, matching its tenant/environment and reader list to the operator-selected synthetic principal. A requested foreign or missing Case fails before an ephemeral bearer is printed. After startup, existing per-request bearer authentication and Case/Timeline/source owner ACLs are rechecked on every route.

This extends the real operator path from startup into the Case selector, Timeline and source preview; it does not add a data store, service, dependency, identity provider or new Case entity. All imported files and permissions are operator controlled, the source preview remains synthetic, and production IAM, live ingestion and durable evidence persistence are still absent. `OPEN-013` is not closed. The approved base version of main.go retains quoted terminal identifier output.

Verification is through the actual CLI loopback integration test, positive and negative Case selection, zero credential leakage on startup rejection, and required GitHub Actions at the exact snapshot SHA. No CI or L2 acceptance may be claimed until independently checked.
