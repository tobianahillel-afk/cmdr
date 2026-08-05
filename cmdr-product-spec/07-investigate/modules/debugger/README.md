---
id: investigate-debugger
domain: 07-investigate
status: deprecated
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: deprecated-pointer
replaced-by: ../analysis-workbench/reverse-engineering-and-debugger/README.md
---
# Debugger — deprecated

This functional source was migrated on 2026-08-05 to the canonical Analysis Workbench submodule:

`../analysis-workbench/reverse-engineering-and-debugger/README.md`

## Migrated needs
- isolated session lifecycle → CAP-INV-339;
- run/pause/step and breakpoints → CAP-INV-340;
- runtime/stack/memory/modules → CAP-INV-341/342;
- trace and snapshots → CAP-INV-343/345.

## Dependents and screens
- `screens/debugger.md` remains active as INV-DBG-001.
- Reverse handoff now targets CAP-INV-339.

## Migration rationale
The previous document was generic and competed with the canonical capability architecture. The active screen remains unchanged and consumes the new capability map.

## Acceptance
- one active functional owner;
- previous needs retained in named capabilities;
- no screen rewrite;
- no engine, command, API, protocol or implementation introduced.
