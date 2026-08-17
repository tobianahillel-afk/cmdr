---
id: investigate-debugger-debug-session-model
domain: 07-investigate
status: deprecated
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: deprecated-pointer
replaced-by: ../analysis-workbench/reverse-engineering-and-debugger/README.md
---
# Debug Session Model — deprecated

This functional source was migrated on 2026-08-05 to the canonical Analysis Workbench submodule:

`../analysis-workbench/reverse-engineering-and-debugger/README.md`

## Migrated needs
- session lifecycle and recovery → CAP-INV-339;
- audited controls → CAP-INV-340;
- failure preserves prior results → CAP-INV-343/345;
- no real Endpoint target → CAP-INV-339.

## Dependents and screens
- `screens/debugger.md` remains active.
- Platform Settings and Studio retain environment/Tool ownership.

## Migration rationale
The previous document was generic and competed with the canonical capability architecture. The active screen remains unchanged and consumes the new capability map.

## Acceptance
- one active functional owner;
- previous needs retained in named capabilities;
- no screen rewrite;
- no engine, command, API, protocol or implementation introduced.
