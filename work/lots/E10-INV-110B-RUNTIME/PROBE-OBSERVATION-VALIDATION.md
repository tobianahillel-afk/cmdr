# E10-INV-110B — reject forged / malformed Timeline probe observations

## Selection and trust boundary

The existing Go standard-library child-probe protocol is reused: `encoding/json`, `math`, `io.EOF` and the first-party `pilotProjectionProbeObservation`. No extra third-party runner or benchmark framework. A full source index, database or provider implementation would not address evidence spoofing at this boundary.

The parent now requires one JSON object (no trailing JSON or unknown fields), exact requested operation/mode/iteration count, positive elapsed, sample mean, Go heap and allocations, and mathematical consistency between `elapsed_ns / operations` and `ns_per_operation` within numeric round-off tolerance. Malformed child data causes a fixed error that never interpolates its contents. New unit tests deliberately falsify these conditions, including fabricated mean and an embedded sensitive marker, while accepting correctly measured latency and resource samples.

This hardens **synthetic first-party performance evidence**, not the true p95 of user requests. Existing on-demand measurements from older GitHub SHAs cannot be attributed to this code. It does not provide a production authentication/permission owner adapter, atomic ACL snapshot, navigation, production SLO or resolution of OPEN-013. Lot B stays READY until its outstanding integration, security and independent-review gates close.

The exact CI and independent L2 verdict must be appended only after GitHub verification; this note itself is not evidence of passed tests.
