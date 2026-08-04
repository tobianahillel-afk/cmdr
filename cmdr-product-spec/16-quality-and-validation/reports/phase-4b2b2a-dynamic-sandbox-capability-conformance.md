---
id: report-phase-4b2b2a-dynamic-sandbox
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-006
  - REQ-PROD-012
  - REQ-PROD-014
  - REQ-INV-005
  - REQ-UX-010
open_decisions:
  - OPEN-005
  - OPEN-013
  - OPEN-014
  - OPEN-015
---
# Phase 4B.2B.2A — Dynamic Sandbox and Behavioral Analysis conformance

## Verdict
**PASS — 123/123 gates.**

## Starting Git state
The exact remote head read before any modification was `2bc2d37e2d726d15e1ee6af1ce3ae633703b4cc3`. It did not contain the previously reported Phase 4B.2B.1. The missing sub-phase was restored first through an ordinary fast-forward chain, then this sub-phase was built on its final restoration commit `df06eeaa94a7a3bc982efc9f5e5c6b4571b279c0`.

## Scope and measures
- Dynamic capabilities: **15/15**, CAP-INV-314 through CAP-INV-328.
- Numbered sections: **405/405**.
- Mandatory tables: **90/90**.
- Empty, prose-only or generic mandatory tables: **0**.
- Duplicate IDs, concurrent owners and active contradictions: **0**.
- Reverse/Debugger/Forensics capabilities: **0**.
- Detailed screens rewritten: **0**.
- APIs, protocols, engines, hypervisors, commands and product code: **0**.
- Requirement IDs: **122 before / 122 after**.
- Open decisions: **15 before / 15 after**.

## Capability conformance
| Capability ID | Sections 1–27 | S8 | S9 | S10 | S13 | S16 | S17 | Front matter | Verdict |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---|
| CAP-INV-314 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-315 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-316 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-317 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-318 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-319 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-320 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-321 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-322 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-323 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-324 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-325 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-326 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-327 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-328 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |

## Required distinctions
Artifact source, Derived Artifact, Runtime Artifact, Evidence and Finding remain distinct. Behavioral Observation remains distinct from Finding. Sandbox Run remains distinct from Analysis Session, Automation Run and Response Run. Simulated interaction/network remain distinct from real user/production. Timeout or absence of observation never proves absence of threat.

## Ownership
Investigate owns Dynamic Analysis context and interpretation. Platform Settings owns Sandbox Environment administration. Studio owns Tool, Tool Call, Workflow and Automation Run. Govern owns authority over real targets. Shared owns Background Jobs, Notifications, Trace, Activity, Export, Object Linking, Versioning and Audit Hooks.

## Gates
| # | Group | Control | Verdict |
|---:|---|---|---|
| 1 | Git | Repository correct | PASS |
| 2 | Git | Branch target correct | PASS |
| 3 | Git | PR #2 correct | PASS |
| 4 | Git | Base main | PASS |
| 5 | Git | PR open | PASS |
| 6 | Git | PR Draft | PASS |
| 7 | Git | PR unmerged | PASS |
| 8 | Git | No force-push | PASS |
| 9 | Git | History preserved | PASS |
| 10 | Git | Root README unchanged | PASS |
| 11 | Git | main unchanged | PASS |
| 12 | Git | Final remote SHA verified | PASS |
| 13 | Sources | Canonical template read | PASS |
| 14 | Sources | Governance/ownership/registers read | PASS |
| 15 | Sources | 4B.1 Artifact/Evidence/Finding read | PASS |
| 16 | Sources | 4B.2A custody/provenance read | PASS |
| 17 | Sources | 4B.2B.1 static capabilities read | PASS |
| 18 | Sources | Dynamic Sandbox active/deprecated sources read | PASS |
| 19 | Sources | Studio Tool/Run/Human Gate sources read | PASS |
| 20 | Sources | Settings Sandbox Environment/admin sources read | PASS |
| 21 | Sources | Govern authority/Run/Result sources read | PASS |
| 22 | Sources | Technical Workbench and Run Shell read | PASS |
| 23 | Sources | Dynamic/Static/Reverse/Debugger/Case/Evidence/Settings/Studio screens read | PASS |
| 24 | Sources | Future phase boundaries verified | PASS |
| 25 | Capabilities | CAP-INV-314..328 convention respected | PASS |
| 26 | Capabilities | 15 unique IDs | PASS |
| 27 | Capabilities | Final count justified | PASS |
| 28 | Capabilities | 15 canonical files | PASS |
| 29 | Capabilities | One owner each | PASS |
| 30 | Capabilities | Users specified | PASS |
| 31 | Capabilities | User problem specific | PASS |
| 32 | Capabilities | Goals specific | PASS |
| 33 | Capabilities | Non-goals specific | PASS |
| 34 | Capabilities | Inputs specific | PASS |
| 35 | Capabilities | Objects read specific | PASS |
| 36 | Capabilities | Objects written specific | PASS |
| 37 | Capabilities | Actions classified | PASS |
| 38 | Capabilities | S13 present | PASS |
| 39 | Capabilities | Functional states specific | PASS |
| 40 | Capabilities | Outputs specific | PASS |
| 41 | Capabilities | Transitions specific | PASS |
| 42 | Capabilities | Source of truth specific | PASS |
| 43 | Capabilities | Provenance specific | PASS |
| 44 | Capabilities | Functional permissions identified | PASS |
| 45 | Capabilities | Limits/errors specific | PASS |
| 46 | Capabilities | Metrics specific | PASS |
| 47 | Capabilities | Delivery classification defined/planned | PASS |
| 48 | Capabilities | At least three GWT criteria | PASS |
| 49 | Capabilities | No-AI alternative present | PASS |
| 50 | Capabilities | Document consumers identified | PASS |
| 51 | Template | All S8 present | PASS |
| 52 | Template | All S9 present | PASS |
| 53 | Template | All S10 present | PASS |
| 54 | Template | All S13 present | PASS |
| 55 | Template | All S16 present | PASS |
| 56 | Template | All S17 present | PASS |
| 57 | Template | 405 sections present | PASS |
| 58 | Template | 90 mandatory tables present | PASS |
| 59 | Template | No empty mandatory table | PASS |
| 60 | Template | No generic or prose-only mandatory table | PASS |
| 61 | Ownership | Investigate owns analytical context | PASS |
| 62 | Ownership | Investigate owns observation interpretation | PASS |
| 63 | Ownership | Studio owns Tool/Tool Call | PASS |
| 64 | Ownership | Studio owns Automation Run | PASS |
| 65 | Ownership | Settings owns Sandbox Environment administration | PASS |
| 66 | Ownership | Govern owns real-target authority | PASS |
| 67 | Concepts | Sandbox Run distinct from Analysis Session | PASS |
| 68 | Concepts | Sandbox Run distinct from Automation Run | PASS |
| 69 | Concepts | Sandbox Run distinct from Response Run | PASS |
| 70 | Concepts | Runtime Artifact distinct from source Artifact | PASS |
| 71 | Concepts | Runtime Artifact distinct from Evidence | PASS |
| 72 | Concepts | Observation distinct from Finding | PASS |
| 73 | Concepts | Simulated interaction/network distinct from real user/production | PASS |
| 74 | Concepts | Dynamic analysis distinct from Reverse/Debugger | PASS |
| 75 | Coverage | Dynamic intake/preconditions defined | PASS |
| 76 | Coverage | Environment selection defined | PASS |
| 77 | Coverage | Dynamic Analysis Session defined | PASS |
| 78 | Coverage | Sandbox Run lifecycle defined | PASS |
| 79 | Coverage | Behavioral timeline defined | PASS |
| 80 | Coverage | Process tree analysis defined | PASS |
| 81 | Coverage | File/system changes defined | PASS |
| 82 | Coverage | Network behavior defined | PASS |
| 83 | Coverage | Persistence candidates defined | PASS |
| 84 | Coverage | Interaction profiles defined | PASS |
| 85 | Coverage | Runtime Artifact management defined | PASS |
| 86 | Coverage | Multi-Run comparison defined | PASS |
| 87 | Coverage | Safety/containment of analysis defined | PASS |
| 88 | Coverage | Provenance/reproducibility defined | PASS |
| 89 | Coverage | Evidence/Finding handoff defined | PASS |
| 90 | Coverage | Technical Workbench/Run Shell constraints respected | PASS |
| 91 | AI/Safety | AI optional | PASS |
| 92 | AI/Safety | No mandatory chatbot | PASS |
| 93 | AI/Safety | No silent Run | PASS |
| 94 | AI/Safety | No invisible environment selection | PASS |
| 95 | AI/Safety | No silent real-network access | PASS |
| 96 | AI/Safety | No silent duration/scope expansion | PASS |
| 97 | AI/Safety | No execution outside authorized environment | PASS |
| 98 | AI/Safety | No opaque score as truth | PASS |
| 99 | AI/Safety | No automatic Evidence qualification | PASS |
| 100 | AI/Safety | No automatic Finding confirmation | PASS |
| 101 | AI/Safety | Automation provenance visible | PASS |
| 102 | AI/Safety | No-AI path complete | PASS |
| 103 | AI/Safety | No permission bypass | PASS |
| 104 | Limits | No detailed screen rewrite | PASS |
| 105 | Limits | No canonical object schema | PASS |
| 106 | Limits | No final cardinality | PASS |
| 107 | Limits | No atomic permission matrix | PASS |
| 108 | Limits | No API | PASS |
| 109 | Limits | No protocol | PASS |
| 110 | Limits | No engine/hypervisor/product chosen | PASS |
| 111 | Limits | No command or product code | PASS |
| 112 | Limits | No Reverse Engineering or Debugger capability | PASS |
| 113 | Limits | No Memory/Disk/advanced Network Forensics | PASS |
| 114 | Publication | Capability/Dependency registers updated | PASS |
| 115 | Publication | Object/Action/AI/Cross-product maps updated | PASS |
| 116 | Publication | Screen map updated without rewrite | PASS |
| 117 | Publication | Requirements Matrix/baseline/OPEN updated | PASS |
| 118 | Publication | STATUS/CHANGELOG/PR coherent | PASS |
| 119 | Publication | No placeholder or empty file | PASS |
| 120 | Publication | No broken modified link | PASS |
| 121 | Publication | No concurrent owner or duplicate | PASS |
| 122 | Publication | At least 120 gates recorded | PASS |
| 123 | Publication | Commits reachable from remote branch | PASS |

## Publication rule
The final commit is verified after fast-forward. A mismatch in remote SHA, ancestry, PR state, README or any gate changes Phase 4B.2B.2A to PARTIAL.
