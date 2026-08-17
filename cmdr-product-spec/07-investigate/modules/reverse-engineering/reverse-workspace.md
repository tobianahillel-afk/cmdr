---
id: investigate-reverse-engineering-reverse-workspace
domain: 07-investigate
status: deprecated
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: deprecated-pointer
replaced-by: ../analysis-workbench/reverse-engineering-and-debugger/README.md
---
# Reverse Workspace — deprecated

This functional source was migrated on 2026-08-05 to the canonical Analysis Workbench submodule:

`../analysis-workbench/reverse-engineering-and-debugger/README.md`

## Migrated needs
- synchronized views → CAP-INV-331/332/333;
- navigation history → CAP-INV-331;
- versioned annotations → CAP-INV-337;
- original immutability → CAP-INV-344.

## Dependents and screens
- `screens/reverse-engineering.md` remains active.
- Technical Workbench Shell remains canonical.

## Migration rationale
The previous document was generic and competed with the canonical capability architecture. The active screen remains unchanged and consumes the new capability map.

## Acceptance
- one active functional owner;
- previous needs retained in named capabilities;
- no screen rewrite;
- no engine, command, API, protocol or implementation introduced.
