# E10-INV-110B — strictly total sort and authorization adapters

This supplement verifies that both read-only, server-authority adapters maintain the same deterministic chronology under different selected-ID orders. The new test checks two visible entries representing the same instant using distinct UTC offsets, mixed observed/inferred kinds, Case/Timeline/source checks for both entries, caller-supplied permission flags ignored, and equal results after swapping requested IDs.

The projection implementation was separately changed in commit `b9d419954931644bebd535e143165f3e8516b788` to use Go standard-library `sort.Slice`, preserving the strict comparison keys and rejection of duplicate visible IDs. No new third-party code or owner mutation is involved.

The new test is independent of the exact snapshot under one-time L2 review; this later commit does **not** redefine that review's SHA. CI, race, format and benchmark conclusions must be verified against the **actual** resulting commit. No local timing improvement, production SLO, identity-provider integration, Case membership guarantee by a real provider, or completion of `OPEN-013` is asserted. The B lot remains READY pending release gates.
