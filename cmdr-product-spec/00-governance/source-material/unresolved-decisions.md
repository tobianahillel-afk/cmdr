---
id: unresolved-decisions
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-06
source-of-truth: canonical
---
# Unresolved decisions

The programme retains **16 open decisions**. OPEN-009 remains the only historically resolved item. Phase 4B.3A.2 closes none.

| Relevant decision | Application in 4B.3A.2 | Disposition |
|---|---|---|
| OPEN-005 | forensic engine selection only | remains open; not reused for Detection runtimes/languages |
| OPEN-008 | source, endpoint, target and runtime availability/support | remains open |
| OPEN-011 / OPEN-012 | Mobile and Cloud scope | remain open and out of scope |
| OPEN-013 | reversible Class-2 changes and future class-3 governance defaults | remains open |
| OPEN-014 | Artifact/Attachment/dataset and lifecycle package relations | remains open |
| OPEN-015 | Automation Run / Response Run and cross-product provenance bridge | remains open |
| **OPEN-017** | Detection runtime, target language and portability strategy | **created open; no option selected** |

## OPEN-017 — Detection runtime, target language and portability strategy
**Status:** open  
**Owner:** Product Architecture  
**Consumers:** Detection Engineering, Command, Platform Settings, Endpoint Agent, Govern, Studio, Shared, Objects, Technique and future implementation.

### Decision question
How will functional Detection Content map to one or more execution runtimes and target representations while preserving versioning, provenance, portability, ownership and governed lifecycle?

### Options retained without selection
1. Portable canonical model with adapters.
2. Native contents per engine.
3. Hybrid model.
4. Future native CMDR runtime.
5. Combination by capability or target.

### Constraints
- no engine, language, vendor syntax, adapter contract, package format or runtime is selected;
- no API, protocol, compiler, AST, storage model or implementation is defined;
- target compatibility remains assessed, not assumed;
- Review, Govern authority and runtime observation remain separate.

All other open decisions retain their prior wording and status.
