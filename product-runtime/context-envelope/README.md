# Pilot context-envelope runtime boundary

This directory reserves the uniquely owned runtime boundary for `PILOT-CONTEXT-001` / `CAP-SET-004`.

**E9-PILOT-001B contains no product implementation.** No executable source or dependency manifest belongs here until E9-PILOT-001C is unlocked. The executable pilot contract and fixtures live under `engineering/pilot/`.

Runtime policy for the future implementation:

- standard-library-only unless a Trusted Base exception is separately reviewed;
- no external product-runtime dependency is currently approved;
- no Product Spec file is copied into runtime;
- authorization remains an explicit caller-supplied outcome;
- this pilot-local contract is not a claim of a canonical public API schema.
