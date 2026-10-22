# Bounded implementation wave selection

E10 selects implementation work from the live readiness compiler rather than from a manually maintained backlog snapshot.

## Policy: reuse-first-v1

The first bounded wave contains exactly one capability.

1. Recompute implementation readiness from the canonical Product Spec and verified runtime ledger.
2. Consider only records whose state is exactly `READY`.
3. Derive capability families from verified `IMPLEMENTED` records.
4. Prefer READY capabilities in an already implemented family to reduce unnecessary architectural surface.
5. Within the preferred set, select the lowest stable capability ID.
6. If no READY capability exists in an implemented family, select the lowest stable READY capability ID globally.
7. Require a non-empty canonical capability file and prove it resolves inside the read-only Product Spec before returning PASS.

This policy does not imply that an existing runtime boundary can be reused. The selected canonical capability contract must be reviewed before any architecture boundary or runtime path is created.

Blocked, proposed, or already implemented capabilities can never enter the wave.
