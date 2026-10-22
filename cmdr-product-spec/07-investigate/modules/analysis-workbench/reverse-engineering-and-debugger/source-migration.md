---
id: reverse-debugger-source-migration
domain: 07-investigate
status: draft
owner: Product Architecture
updated: 2026-08-05
source-of-truth: canonical
---
# Source migration

| Previous source | Status after migration | Canonical replacement | Needs retained | Dependents |
|---|---|---|---|---|
| `modules/reverse-engineering/README.md` | deprecated pointer | this module README and CAP-INV-329..338 | disassembly, decompilation, xrefs, CFG, symbols | INV-REV-001 retained |
| `modules/reverse-engineering/reverse-workspace.md` | deprecated pointer | CAP-INV-330/331/337 plus Technical Workbench | synchronized views, history, versioned annotations, original immutability | INV-REV-001 retained |
| `modules/debugger/README.md` | deprecated pointer | this module README and CAP-INV-339..344 | run/pause/step, breakpoints, runtime views, trace/snapshots | INV-DBG-001 retained |
| `modules/debugger/debug-session-model.md` | deprecated pointer | CAP-INV-339/340/343/345 | lifecycle, audit, failure preserves prior results | INV-DBG-001 retained |

The two screens remain active sources of needs. Migration date: 2026-08-05. Acceptance: no competing active functional owner, links point to the canonical module, and no screen specification is rewritten.
