# E10-INV-110B — bounded native fuzz campaign

## Decision and reuse

The existing first-party Go `FuzzProjectionOpaqueFieldsDoNotInterfere` target contains deterministic seed tests and bounded opaque source mutations (denied, unknown, missing, foreign tenant). Ordinary `go test` runs seeds only. Add an on-demand 25-second native Go fuzz-engine campaign on GitHub Actions, reusing the existing pinned Go setup and `investigation-timeline-evidence.yml`. No third-party fuzz framework, vendor SaaS, paid API, new database, or product runtime dependency. The new job is triggered on an update to the existing evidence workflow or by explicit `workflow_dispatch`; it is **not** scheduled for every PR commit.

## Safety and evidence

The campaign operates on synthetic fields and a 5-minute job timeout; 2 workers, per-process `go test -timeout=3m`. It neither retrieves secret or production records nor uploads discovered fuzz corpus; failures remain observable in GitHub Actions logs. Normal CI still runs seed cases and remains the release gate for code changes. Fuzz PASS means **no issue found during one bounded campaign**, not proof of authorization correctness, complete state coverage or a production SLO. The canonical CAP-INV-110 and `OPEN-013` remain unchanged; B remains READY.

**Publication is not execution evidence:** record the exact GitHub SHA, fuzz job id, `conclusion`, engine execution count (if available), source and fixed limits only after checking the actual workflow log. Cross-platform macOS/Windows smoke jobs already exist in this workflow and their outcomes are independent of fuzz.
