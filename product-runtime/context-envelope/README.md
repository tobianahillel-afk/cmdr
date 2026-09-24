# Pilot context-envelope runtime

This directory implements the bounded runtime selected by `PILOT-CONTEXT-001` / `CAP-SET-004`.

The executable contract and canonical pilot fixtures remain under `engineering/pilot/`. This package is a pilot-local implementation, not a claim that its Go API is a canonical public product schema.

Runtime constraints:

- standard-library-only Go module;
- no external product-runtime dependencies;
- no persistence or network I/O;
- no authorization/RBAC/ABAC engine;
- destination authorization is an explicit caller-supplied outcome;
- no raw payload, token, credential, or secret field;
- incompatible tenant/environment context is cleared rather than silently replaced;
- denied/expired/source-missing states mask protected references.

`Project(Input)` is deterministic and side-effect free. Security activation and measured coverage are handled separately by `E9-PILOT-001C-SECURITY`.
