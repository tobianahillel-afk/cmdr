---
id: analysis-workbench-capability-map
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
requirements:
  - REQ-PROD-012
  - REQ-PROD-014
---
# Capability map — Analysis Workbench

## Static Analysis foundation
CAP-INV-301..313 remain canonical and unchanged for intake, Analysis Session, Tool selection, preview, format/structure, strings/content, static binary/script/archive analysis, comparison, Derived Artifacts, provenance and Evidence/Finding handoff.

## Reverse Engineering and Debugger
| ID | Capability | Family | Primary local concepts | Classes |
|---|---|---|---|---|
| CAP-INV-329 | Reverse Engineering Intake and Preconditions | reverse-intake | Trace / Activity event; Reverse Intake | 0,2 |
| CAP-INV-330 | Reverse Analysis Session Management | reverse-session | Trace / Activity event; Reverse Analysis Session | 0,2 |
| CAP-INV-331 | Binary Navigation and Address Space Orientation | navigation | Trace / Activity event; Bookmark / location annotation | 0,2 |
| CAP-INV-332 | Disassembly Inspection | code-analysis | Trace / Activity event; Disassembly annotation / rename relation | 0,2 |
| CAP-INV-333 | Decompilation Inspection | code-analysis | Trace / Activity event; Local interpretation correction | 0,2 |
| CAP-INV-334 | Functions, Symbols and Cross-References | code-analysis | Trace / Activity event; Function interpretation | 0,2 |
| CAP-INV-335 | Control Flow, Call Graph and Data Flow Analysis | code-analysis | Trace / Activity event; Graph annotation / path selection | 0,2 |
| CAP-INV-336 | Data Types, Structures and Memory Layout Analysis | code-analysis | Trace / Activity event; Type Definition / Structure Definition | 0,2 |
| CAP-INV-337 | Analyst Annotation, Renaming and Knowledge Capture | knowledge | Trace / Activity event; Annotation / Comment / Bookmark / Rename | 0,2 |
| CAP-INV-338 | Binary Diffing and Version Comparison | comparison | Trace / Activity event; Binary comparison result | 0,1,2 |
| CAP-INV-339 | Debugger Session Management | debugger | Trace / Activity event; Debugger Session | 0,1,2 |
| CAP-INV-340 | Breakpoint and Execution Control | debugger | Trace / Activity event; Breakpoint | 1,2 |
| CAP-INV-341 | Runtime State, Threads and Call Context Inspection | debugger | Trace / Activity event; Runtime State Snapshot | 0,1,2 |
| CAP-INV-342 | Memory, Modules and Loaded Image Inspection | debugger | Trace / Activity event; Memory/Module observation | 0,1,2 |
| CAP-INV-343 | Debug Events, Exceptions and Trace Analysis | debugger | Trace / Activity event; Debug Event / Exception Event / Debug Trace | 0,1,2 |
| CAP-INV-344 | Reversible Analytical Experiments and Patch Hypotheses | experiments | Trace / Activity event; Patch Hypothesis | 0,1,2 |
| CAP-INV-345 | Reverse and Debugger Provenance and Reproducibility | provenance | Trace / Activity event; Provenance relation / Reproducibility Assessment | 0,1,2 |
| CAP-INV-346 | Reverse and Debugger Handoff to Evidence, Findings and Detection Engineering | handoff | Trace / Activity event; Evidence candidate package | 0,2 |

Dynamic Sandbox remains CAP-INV-314..328 in its own module. No forensic capability is included.
