---
id: investigate-memory-forensics-plugin-contract
domain: 07-investigate
status: deprecated
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: deprecated-pointer
replaced_by: ../analysis-workbench/memory-forensics/README.md
---
# Plugin Contract — deprecated functional pointer

Tool and plugin lifecycle belongs to CMDR Studio. Memory Forensics now consumes Tool, Tool Call, version, parameters, partial errors and provenance through the canonical module `../analysis-workbench/memory-forensics/`.

## Migration
- version and parameters → CAP-INV-348/350/361;
- partial failures → all capability error states;
- provenance → CAP-INV-361;
- sensitive output masking → CAP-INV-356.

No plugin framework, engine, offset or command is selected.
