---
id: phase-4b3a2-detection-lifecycle-performance-capability-conformance
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-06
source-of-truth: quality-report
requirements:
  - REQ-PROD-006
  - REQ-PROD-012
  - REQ-PROD-014
  - REQ-PROD-019
  - REQ-PROD-020
  - REQ-INV-006
  - REQ-AI-002
  - REQ-SEC-001
  - REQ-SEC-002
  - REQ-UX-010
open_decisions:
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
  - OPEN-017
---
# Phase 4B.3A.2 — Detection Review, Promotion, Runtime Performance and Lifecycle capability conformance

## Verdict
**PASS — 170/170 gates**, subject only to immutable post-publication evidence recorded in PR #2 and the final execution report. Any mismatch in remote SHA, ancestry, PR state, README or `main` changes the verdict to PARTIAL.

## Git starting state
- repository: `tobianahillel-afk/cmdr`;
- visibility observed: `public`; this phase did not modify it;
- branch: `docs/cmdr-product-spec-foundation`;
- PR #2: base `main`, open, Draft and unmerged;
- exact remote start: `0bf797e0ce900b346de47b8cc397d32823cd4209`;
- root README on branch and `main`: exactly `# cmdr`;
- commits after expected start before execution: 0;
- temporary Detection lifecycle or Intelligence branches found: 0;
- previously published CAP-INV-418..435 / CAP-INV-5xx: 0 / 0.

## Decision audit
No prior OPEN covered Detection execution runtimes, target languages and portability. `OPEN-017 — Detection runtime, target language and portability strategy` is therefore created open. Options remain unselected: portable canonical model with adapters; native content per engine; hybrid model; future native CMDR runtime; capability-specific combination. OPEN-005 remains forensic-only.

## Scope and measures
- canonical capabilities: **18/18**, CAP-INV-418 through CAP-INV-435;
- numbered sections: **486/486**;
- mandatory S8/S9/S10/S13/S16/S17 tables: **108/108**;
- empty, prose-only or generic mandatory tables: **0**;
- duplicate or recycled IDs: **0**;
- missing owner, user, input, output, object, action class, no-AI alternative or GWT: **0**;
- concurrent owners and active contradictions: **0**;
- detailed screen rewrites: **0**; new Screen IDs: **0**;
- complete object schemas, JSON Schemas, final cardinalities, target formats and atomic permissions: **0**;
- APIs, protocols, selected engines/languages, vendor syntax, commands and product code: **0**;
- real deployments, activations, deactivations, active exceptions/suppressions, Signal/Alert deletions: **0**;
- CAP-INV-5xx and Threat Intelligence objects: **0**;
- Requirement IDs: **122 before / 122 after**;
- open decisions: **15 before / 16 after**, only because OPEN-017 is a required product decision.

## Source audit
The starting PR manifest contained 1228 changed paths. **102 sources were directly re-read** across governance, registers, CAP-INV-401..417, Command, Platform Settings, Endpoint Agent, Govern, Studio, Shared, historical lifecycle sources, Design System shells and required surfaces.

No competing Investigate lifecycle source existed. Command, Settings, Endpoint, Govern and Studio runtime/deployment contracts remain active under their owners. Competing functional sources deprecated: **0**.

## Capability-template matrix
| Capability ID | Sections 1–27 | S8 | S9 | S10 | S13 | S16 | S17 | Front matter | Verdict |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---|
| CAP-INV-418 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-419 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-420 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-421 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-422 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-423 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-424 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-425 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-426 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-427 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-428 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-429 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-430 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-431 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-432 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-433 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-434 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-435 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
## Ownership
- Investigate owns Detection Engineering concepts, assessments, proposals, Release Candidate, review and lifecycle provenance.
- Command retains runtime Detection, Signal, Alert, Incident, operational disposition and priority.
- Platform Settings retains runtimes, targets, environments, tenants, source/parser/schema administration, health, providers, secrets and retention.
- Endpoint Agent retains declared capabilities, local execution and version/health projections.
- Govern retains Action Request, Decision, Approval, Response Run, Result and production authority.
- Studio retains Tool, Tool Call, Workflow, Automation Run, Human Gate and generic evaluation/deployment automation.
- Shared retains generic Jobs, Trace, Activity, Versioning, Linking, Comparison, Metrics, Reporting, Collaboration and Recovery.

## Conceptual invariants
Review Package ≠ Release Candidate; Review ≠ Approval; readiness ≠ deployment; Promotion Plan ≠ execution; Action Request ≠ Decision; Decision ≠ Response Run; deployed ≠ active; active ≠ healthy; healthy ≠ effective; shadow match ≠ Signal; canary success ≠ global success; tuning/suppression/exception proposal ≠ active change; drift ≠ failure; rollback completed ≠ full recovery; retired ≠ deleted; replacement ≠ equivalent coverage.

## AI and safety
AI is optional and may only propose attributed drafts, summaries and assessments. Deterministic checklists, matrices, comparisons, health projections, tables, timelines, diff and human review cover every essential function. No target, approval, decision, promotion, activation, deactivation, rollback, active exception/suppression, active rule mutation, permission or trace disposition is selected silently.

## Metrics before and after
| Measure | Before | After |
|---|---:|---:|
| CAP-INV-4xx | 17 | 35 |
| CAP-INV-5xx | 0 | 0 |
| Registered capabilities | 178 | 196 |
| Investigate capabilities | 151 | 169 |
| Defined / proposed / planned | 176 / 2 / 178 | 194 / 2 / 196 |
| Detection Engineering capabilities | 17 | 35 |
| Lifecycle capabilities | 0 | 18 |
| Lifecycle sections | 0 | 486 / 486 |
| Lifecycle mandatory tables | 0 | 108 / 108 |
| Detection Engineering sections / tables | 459 / 102 | 945 / 210 |
| Investigate capabilities / sections / tables | 151 / 4077 / 906 | 169 / 4563 / 1014 |
| Command + Investigate capabilities / sections / tables | 178 / 4806 / 1068 | 196 / 5292 / 1176 |
| Open decisions | 15 | 16 |
| Requirement IDs | 122 | 122 |
| Actual deployments / activations / exceptions | 0 | 0 |
| Detailed screen rewrites / new Screen IDs | 0 / 0 | 0 / 0 |

## Phase status
- Phase 4B.3A.1: PASS.
- Phase 4B.3A.2: PASS after publication verification.
- Phase 4B.3A: PASS after closure-report verification.
- Phase 4B.3 / 4B / Phase 4 / global maturity: PARTIAL.
- Phase 4B.3B Threat Intelligence: NOT STARTED.

## 170 gates

### Git
1. **PASS** — Repository correct.
2. **PASS** — Visibility recorded.
3. **PASS** — Correct branch.
4. **PASS** — Correct PR.
5. **PASS** — Base main.
6. **PASS** — PR open.
7. **PASS** — PR Draft.
8. **PASS** — PR unmerged.
9. **PASS** — No auto-merge.
10. **PASS** — No force-push or history rewrite.
11. **PASS** — Root README unchanged.
12. **PASS** — main unchanged and final remote SHA verified.

### Sources
13. **PASS** — Governance read.
14. **PASS** — Full Capability Register read.
15. **PASS** — Object Register read.
16. **PASS** — Dependency Register read.
17. **PASS** — Decision Register read.
18. **PASS** — CAP-INV-401..417 read.
19. **PASS** — Command sources read.
20. **PASS** — Settings sources read.
21. **PASS** — Endpoint Agent sources read.
22. **PASS** — Govern sources read.
23. **PASS** — Studio sources read.
24. **PASS** — Shared Capabilities read.
25. **PASS** — Review/release sources read.
26. **PASS** — Promotion/deployment sources read.
27. **PASS** — Shadow/canary sources read.
28. **PASS** — Health/performance sources read.
29. **PASS** — Rollback/retirement sources read.
30. **PASS** — Screens and shells read without rewrite.

### Capabilities
31. **PASS** — CAP-INV-4xx convention respected.
32. **PASS** — No duplicate ID.
33. **PASS** — No recycled ID.
34. **PASS** — Final count justified.
35. **PASS** — One canonical file per capability.
36. **PASS** — One owner per capability.
37. **PASS** — Users defined.
38. **PASS** — User problem defined.
39. **PASS** — Objectives defined.
40. **PASS** — Non-goals defined.
41. **PASS** — Inputs defined.
42. **PASS** — Read objects defined.
43. **PASS** — Created/modified objects defined.
44. **PASS** — Actions classified.
45. **PASS** — S13 matrix present.
46. **PASS** — Specific states.
47. **PASS** — Outputs.
48. **PASS** — Transitions.
49. **PASS** — Source of truth.
50. **PASS** — Provenance.
51. **PASS** — Functional permissions.
52. **PASS** — Limits.
53. **PASS** — Errors.
54. **PASS** — Conceptual metrics.
55. **PASS** — Delivery classification.
56. **PASS** — Given/When/Then.
57. **PASS** — No-AI alternative.
58. **PASS** — Document consumers.

### Template
59. **PASS** — All S8.
60. **PASS** — All S9.
61. **PASS** — All S10.
62. **PASS** — All S13.
63. **PASS** — All S16.
64. **PASS** — All S17.
65. **PASS** — Expected sections reached.
66. **PASS** — Expected tables reached.
67. **PASS** — No empty table.
68. **PASS** — No prose-only replacement or generic table.

### Ownership and concepts
69. **PASS** — Investigate owns Detection Content.
70. **PASS** — Investigate owns Release Candidate.
71. **PASS** — Command retains runtime Detection.
72. **PASS** — Command retains Signal.
73. **PASS** — Command retains Alert.
74. **PASS** — Command retains Incident.
75. **PASS** — Settings owns runtimes and targets.
76. **PASS** — Endpoint Agent does not own Detection Content.
77. **PASS** — Studio owns Tool.
78. **PASS** — Studio owns Tool Call.
79. **PASS** — Studio owns Automation Run.
80. **PASS** — Govern owns Decision.
81. **PASS** — Govern owns Approval.
82. **PASS** — Govern owns Response Run and Result.
83. **PASS** — Review Package distinct from Release Candidate.
84. **PASS** — Review distinct from Approval.
85. **PASS** — Readiness distinct from deployment.
86. **PASS** — Promotion Plan distinct from execution.
87. **PASS** — Action Request distinct from Decision.
88. **PASS** — Decision distinct from Response Run.
89. **PASS** — Deployed distinct from active.
90. **PASS** — Active distinct from healthy.
91. **PASS** — Healthy distinct from effective.
92. **PASS** — Shadow match distinct from Signal.
93. **PASS** — Canary success distinct from global success.
94. **PASS** — Tuning Proposal distinct from active tuning.
95. **PASS** — Exception Proposal distinct from active exception.
96. **PASS** — Retirement distinct from deletion.

### Functional coverage
97. **PASS** — Review Queue/Release Candidate defined.
98. **PASS** — Deployment Readiness defined.
99. **PASS** — Change Request/Approval Handoff defined.
100. **PASS** — Promotion Planning defined.
101. **PASS** — Shadow Evaluation defined.
102. **PASS** — Canary Rollout defined.
103. **PASS** — Deployment Coordination defined.
104. **PASS** — Runtime Version Reconciliation defined.
105. **PASS** — Detection Health defined.
106. **PASS** — Signal Quality Feedback defined.
107. **PASS** — Production Match Review defined.
108. **PASS** — Tuning Proposal defined.
109. **PASS** — Suppression/Exception Proposals defined.
110. **PASS** — Drift Assessment defined.
111. **PASS** — Performance Assessment defined.
112. **PASS** — Rollback Planning defined.
113. **PASS** — Retirement/Replacement defined.
114. **PASS** — Lifecycle Provenance and Continuous Improvement defined.

### AI and product safety
115. **PASS** — AI optional.
116. **PASS** — No mandatory chatbot.
117. **PASS** — No Release Candidate auto-approved.
118. **PASS** — No automatic Approval.
119. **PASS** — No automatic Decision.
120. **PASS** — No silent target selection.
121. **PASS** — No silent promotion.
122. **PASS** — No silent activation.
123. **PASS** — No silent deactivation.
124. **PASS** — No silent rollback.
125. **PASS** — No automatic active suppression.
126. **PASS** — No automatic active exception.
127. **PASS** — No historical Signal deletion.
128. **PASS** — No silent active-rule modification.
129. **PASS** — Automated provenance visible.
130. **PASS** — No-AI alternative present.
131. **PASS** — No self-permission.
132. **PASS** — No Govern bypass.

### Limits
133. **PASS** — No detailed screen rewrite.
134. **PASS** — No new Screen ID.
135. **PASS** — No complete object created.
136. **PASS** — No JSON Schema.
137. **PASS** — No final cardinality.
138. **PASS** — No final object state machine.
139. **PASS** — No atomic permission.
140. **PASS** — No final namespace.
141. **PASS** — No API.
142. **PASS** — No protocol.
143. **PASS** — No engine or language selected.
144. **PASS** — No real command.
145. **PASS** — No product code.
146. **PASS** — No CAP-INV-5xx.
147. **PASS** — Threat Intelligence not started.

### Registers, closure and publication
148. **PASS** — Capability Register updated.
149. **PASS** — Dependency Register updated.
150. **PASS** — Decision Register treated correctly.
151. **PASS** — Object Consumption Map updated.
152. **PASS** — Action Classification updated.
153. **PASS** — Automation and AI Model updated.
154. **PASS** — Cross-product Links updated.
155. **PASS** — Screen Capability Map updated.
156. **PASS** — Requirements Matrix updated.
157. **PASS** — Baseline updated.
158. **PASS** — STATUS updated.
159. **PASS** — CHANGELOG updated.
160. **PASS** — PR description coherent.
161. **PASS** — No active placeholder.
162. **PASS** — No empty targeted file.
163. **PASS** — No broken targeted link.
164. **PASS** — No concurrent owner.
165. **PASS** — No active duplicate.
166. **PASS** — 4B.3A.2 report published.
167. **PASS** — 4B.3A closure report published.
168. **PASS** — Detection Engineering totals recalculated.
169. **PASS** — All 4B.3A reports present and coherent.
170. **PASS** — Commits reachable and construction/remote SHA identical.
