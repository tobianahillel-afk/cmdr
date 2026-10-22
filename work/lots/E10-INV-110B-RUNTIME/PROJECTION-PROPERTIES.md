# E10-INV-110B — independent deterministic runtime properties

Additional regression tests in `product-runtime/investigation-timeline/projection_properties_test.go` exercise invariants **without relying on the 18 predeclared contract fixtures**.

- Thirty fixed-seed permutations of forty same-instant mixed observed/inferred entries must preserve exact stable timeline order, classification, provenance and input immutability.
- For four masked-source classes (denied, unknown, missing, foreign), seven variations of confidential entry IDs, references, time, source type, version and provenance must not influence authorized result bytes or status when hidden-source presence is held fixed.
- These tests deliberately do **not** claim that exposing the `filtered` or `partial` status is acceptable for every production threat model, nor that caller-provided access strings are decisions from actual authoritative providers.

This work is independent from the same lot's `go test -json` negative-gate package-binding hardening, whose L2 is pinned to SHA `ac9c86bf68bfd0f3586c044308435b6617f26c61`. The tests do not modify that audited code or the canonical Product Spec. A green CI only attests these executable properties, **not** production Case/Timeline/source authorization. Performance targets remain synthetic, `OPEN-013` remains unresolved.
