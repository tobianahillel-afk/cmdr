---
id: investigate-reverse-engineering
domain: 07-investigate
status: deprecated
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: deprecated-pointer
replaced-by: ../analysis-workbench/reverse-engineering-and-debugger/README.md
---
# Reverse Engineering — deprecated

This functional source was migrated on 2026-08-05 to the canonical Analysis Workbench submodule:

`../analysis-workbench/reverse-engineering-and-debugger/README.md`

## Migrated needs
- disassembly and decompilation → CAP-INV-332/333;
- cross-references and symbols → CAP-INV-334;
- CFG → CAP-INV-335;
- annotations and derived output → CAP-INV-337/344.

## Dependents and screens
- `screens/reverse-engineering.md` remains active as INV-REV-001.
- Static and Dynamic handoffs now target CAP-INV-329.

## Migration rationale
The previous document was generic and competed with the canonical capability architecture. The active screen remains unchanged and consumes the new capability map.

## Acceptance
- one active functional owner;
- previous needs retained in named capabilities;
- no screen rewrite;
- no engine, command, API, protocol or implementation introduced.
