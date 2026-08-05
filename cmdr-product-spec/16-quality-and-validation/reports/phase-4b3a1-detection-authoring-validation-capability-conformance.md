---
id: phase-4b3a1-detection-authoring-validation-capability-conformance
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
  - REQ-INV-006
  - REQ-AI-002
  - REQ-UX-010
open_decisions:
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
---
# Phase 4B.3A.1 — Detection Engineering Foundations, Authoring and Validation capability conformance

## Verdict
**PASS — 160/160 gates**, subject only to immutable post-publication evidence recorded in PR #2 and the final execution report. Any mismatch in branch SHA, ancestry, PR state or README changes the verdict to PARTIAL.

## Git starting state
- repository: `tobianahillel-afk/cmdr`;
- visibility observed: `public`; this phase did not modify it;
- branch: `docs/cmdr-product-spec-foundation`;
- PR: #2, base `main`, open, Draft and unmerged;
- exact remote start: `56fd1b4decdc323c5d3bfdf5fb3e0bd8fc082f1f`;
- root README on branch and `main`: exactly `# cmdr`;
- published commits after expected SHA: 0;
- temporary Detection or Intelligence branches found: 0;
- previously published CAP-INV-4xx / CAP-INV-5xx: 0 / 0.

## Scope and measures
- canonical capabilities: **17/17**, CAP-INV-401 through CAP-INV-417;
- numbered sections: **459/459**;
- mandatory S8/S9/S10/S13/S16/S17 tables: **102/102**;
- empty, prose-only or generic mandatory tables: **0**;
- duplicate or recycled IDs: **0**;
- missing owner, user, input, output, object, action class, no-AI alternative or GWT: **0**;
- concurrent owners and active contradictions: **0**;
- detailed screen rewrites: **0**; new Screen IDs: **0**;
- canonical object schemas, JSON Schemas, final cardinalities, ASTs and atomic permissions: **0**;
- APIs, protocols, engines, rule languages, vendor syntax, commands and product code: **0**;
- promoted/deployed/activated/deactivated rules, active exceptions and Signal/Alert deletions: **0**;
- CAP-INV-5xx and Threat Intelligence objects: **0**;
- Requirement IDs: **122 before / 122 after**;
- open decisions: **15 before / 15 after**.

## Source audit
The full 1196-path PR manifest was reviewed. At least 60 documentary sources were directly fetched across governance, registers, objects, Command, Signals/Hunt, Cases/Evidence/Findings, Collection and analysis handoffs, Event Search, Endpoint Agent, Platform Settings, Studio, Govern, Shared Capabilities, journeys, Technical Workbench and screens.

Historical Detection-related content was distributed across source owners:
- Command runtime Detection, Signal, Alert and Incident;
- Event Search, Query and Hunt;
- Endpoint Agent local detection, update, suppression and behavioral evaluation;
- Settings sources, parsers, health and retention;
- Shared Query, Correlation, Data Quality and Normalization;
- Studio datasets, evaluations, regression testing and Workflow Builder;
- telemetry-to-detection and detection-to-signal journeys.

No active standalone Detection Engineering authoring module or autonomous Rule Builder screen existed. Competing functional documents deprecated: **0**. Runtime, deployment and source-owner documents remain active for their proper phases.

## Capability-template matrix
| Capability ID | Sections 1–27 | S8 | S9 | S10 | S13 | S16 | S17 | Front matter | Verdict |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| CAP-INV-401 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-402 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-403 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-404 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-405 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-406 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-407 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-408 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-409 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-410 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-411 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-412 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-413 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-414 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-415 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-416 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-417 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |

## Ownership
- Investigate owns Detection Engineering Project, Detection Hypothesis, Detection Content Draft, authoring context, test/review concepts, Coverage/Gap and Review Package as functional concepts.
- Command retains runtime Detection, Signal, Alert, Incident, triage and prioritization.
- Platform Settings retains Data Source, connector, parser, schema, health, retention, environment and secret administration.
- Endpoint Agent retains declared telemetry capabilities, raw local observations and future runtime evaluation.
- Studio retains Tool, Tool Call, Workflow, Automation Agent, Human Gate, Automation Run and generic Dataset/Evaluation mechanisms.
- Govern retains Decision, Approval and future production authority.
- Shared retains Query, Telemetry Event, Jobs, Trace, Activity, Versioning, Linking, Search, Export, Reporting, Collaboration, Comments, Assignments, Comparison, Inspector and Recovery.

## Conceptual invariants
Detection Content ≠ runtime Detection; Draft ≠ deployed rule; Rule ≠ Search Query; Saved Search ≠ Detection Content; Detection Hypothesis ≠ Case Hypothesis; Project ≠ Case/Automation Run; configured source ≠ healthy/sufficient source; field present ≠ reliable; mapping ≠ certain equivalence; schema valid ≠ data correct; syntax valid ≠ logic/usefulness; match/non-match ≠ certain TP/FN; test dataset ≠ certain ground truth; replay ≠ deployment or future performance; coverage mapping ≠ effective coverage; Review Package ≠ Approval/deployment.

## AI and safety
AI is optional and may only propose attributed drafts or explanations. Deterministic editors, validators, catalogues, fixtures, replay, comparison, checklists and human review cover all essential workflows. No silent selection, invented truth, automatic certain TP/FN, permission escalation, promotion, deployment, activation, exception or runtime object mutation is allowed.

## Engine and language decision
OPEN-005 remains limited to forensic engines and is not reused. No existing OPEN explicitly covers Detection engine/language selection. The gap is recorded without blocking this functional phase and without creating a new OPEN solely for classification.

## Metrics before and after
| Measure | Before | After |
|---|---:|---:|
| CAP-INV-4xx | 0 | 17 |
| CAP-INV-5xx | 0 | 0 |
| Registered capabilities | 161 | 178 |
| Investigate capabilities | 134 | 151 |
| Defined / proposed / planned | 159 / 2 / 161 | 176 / 2 / 178 |
| Detection Authoring capabilities | 0 | 17 |
| Sections | 0 | 459 / 459 |
| Mandatory tables | 0 | 102 / 102 |
| Empty / prose-only / generic mandatory tables | 0 | 0 |
| Missing owner/user/input/output/object/action/no-AI/GWT | 0 | 0 |
| Active duplicates / concurrent owners | 0 | 0 |
| Module functional files | 0 | 30 |
| Deprecated authoring sources | 0 | 0 |
| Screens read / modified / rewritten / new IDs | 10 / 0 / 0 / 0 | 10 / 0 / 0 / 0 |
| Canonical object files / atomic permissions | 0 / 0 | 0 / 0 |
| APIs / protocols / engines / languages / commands / code | 0 | 0 |
| Deployed rules / active exceptions / Intelligence content | 0 | 0 |
| Requirement IDs / OPEN | 122 / 15 | 122 / 15 |
| Investigate capabilities / sections / tables | 134 / 3618 / 804 | 151 / 4077 / 906 |
| Command + Investigate capabilities / sections / tables | 161 / 4347 / 966 | 178 / 4806 / 1068 |

## Phase status
- Phase 4B.3A.1: PASS after publication verification.
- Phase 4B.3A: PARTIAL.
- Phase 4B.3: PARTIAL.
- Phase 4B: PARTIAL.
- Phase 4 / global maturity: PARTIAL.
- Phase 4B.3A.2: NOT STARTED.
- Phase 4B.3B Threat Intelligence: NOT STARTED.

## 160 gates

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
10. **PASS** — No force-push/history rewrite.
11. **PASS** — Root README unchanged.
12. **PASS** — main unchanged and final remote SHA verified.

### Sources
13. **PASS** — Relevant governance read.
14. **PASS** — Full Capability Register read.
15. **PASS** — Object Register read.
16. **PASS** — Dependency Register read.
17. **PASS** — Command sources read.
18. **PASS** — Signals/Hunt sources read.
19. **PASS** — Cases/Evidence/Findings sources read.
20. **PASS** — Collection sources read.
21. **PASS** — Static/Dynamic sources read.
22. **PASS** — Reverse/Debugger sources read.
23. **PASS** — Memory/Disk/Network sources read.
24. **PASS** — Event Search/SIEM sources read.
25. **PASS** — Historical Detection Engineering sources read.
26. **PASS** — Endpoint Agent/Settings sources read.
27. **PASS** — Studio/Govern sources read.
28. **PASS** — Shared Capabilities and Technical Workbench read.
29. **PASS** — Screens read without rewrite.

### Capabilities
30. **PASS** — CAP-INV-4xx convention respected.
31. **PASS** — No duplicate ID.
32. **PASS** — No recycled ID.
33. **PASS** — Final count justified.
34. **PASS** — One canonical file per capability.
35. **PASS** — One owner per capability.
36. **PASS** — Users defined.
37. **PASS** — User problem defined.
38. **PASS** — Objectives defined.
39. **PASS** — Non-goals defined.
40. **PASS** — Inputs defined.
41. **PASS** — Read objects defined.
42. **PASS** — Created/modified objects defined.
43. **PASS** — Actions classified.
44. **PASS** — S13 matrix present.
45. **PASS** — Specific states.
46. **PASS** — Outputs.
47. **PASS** — Transitions.
48. **PASS** — Source of truth.
49. **PASS** — Provenance.
50. **PASS** — Functional permissions.
51. **PASS** — Limits.
52. **PASS** — Errors.
53. **PASS** — Conceptual metrics.
54. **PASS** — Delivery classification.
55. **PASS** — Given/When/Then.
56. **PASS** — No-AI alternative.
57. **PASS** — Document consumers.

### Template
58. **PASS** — All S8.
59. **PASS** — All S9.
60. **PASS** — All S10.
61. **PASS** — All S13.
62. **PASS** — All S16.
63. **PASS** — All S17.
64. **PASS** — Expected sections reached.
65. **PASS** — Expected tables reached.
66. **PASS** — No empty table.
67. **PASS** — No prose-only replacement or generic table.

### Ownership and concepts
68. **PASS** — Investigate owns Detection Engineering Project.
69. **PASS** — Investigate owns Detection Content.
70. **PASS** — Command retains runtime Detection.
71. **PASS** — Command retains Signal.
72. **PASS** — Command retains Alert.
73. **PASS** — Command retains Incident.
74. **PASS** — Settings owns Data Source administration.
75. **PASS** — Settings owns schema/parser administration.
76. **PASS** — Studio owns Tool.
77. **PASS** — Studio owns Tool Call.
78. **PASS** — Studio owns Automation Run.
79. **PASS** — Govern retains future production authority when required.
80. **PASS** — Detection Content distinct from runtime Detection.
81. **PASS** — Detection Rule distinct from Search Query.
82. **PASS** — Detection Project distinct from Case.
83. **PASS** — Detection Project distinct from Automation Run.
84. **PASS** — Active Data Source distinct from sufficient source.
85. **PASS** — Field present distinct from reliable field.
86. **PASS** — Mapping distinct from certain semantic equivalence.
87. **PASS** — Syntax validation distinct from semantic validation.
88. **PASS** — Semantic validation distinct from production performance.
89. **PASS** — Match distinct from confirmed malicious activity.
90. **PASS** — Non-match distinct from certain false negative.
91. **PASS** — Test Dataset distinct from certain ground truth.
92. **PASS** — Replay Result distinct from runtime Detection.
93. **PASS** — Coverage mapping distinct from effective coverage.
94. **PASS** — Review Package distinct from deployment.

### Functional coverage
95. **PASS** — Intake defined.
96. **PASS** — Project/Workspace defined.
97. **PASS** — Detection Hypothesis defined.
98. **PASS** — Data Readiness defined.
99. **PASS** — Schema/Field Mapping Review defined.
100. **PASS** — Detection Content Authoring defined.
101. **PASS** — Conditions/Correlation defined.
102. **PASS** — Sequence/Threshold/Window defined.
103. **PASS** — Enrichment Requirements defined.
104. **PASS** — Metadata/Ownership/Documentation defined.
105. **PASS** — Structural/Semantic Validation defined.
106. **PASS** — Test Scenario/Dataset Management defined.
107. **PASS** — Expected Outcome/Test Oracle defined.
108. **PASS** — Historical Replay defined.
109. **PASS** — Match and FP/FN Review defined.
110. **PASS** — Coverage and Gap Analysis defined.
111. **PASS** — Provenance and Review Handoff defined.

### AI and product security
112. **PASS** — AI optional.
113. **PASS** — No mandatory chatbot.
114. **PASS** — No silent draft.
115. **PASS** — No invisible field or mapping selection.
116. **PASS** — No automatically accepted test oracle.
117. **PASS** — No invented ground truth.
118. **PASS** — No automatic certain TP.
119. **PASS** — No automatic certain FN.
120. **PASS** — No promoted Detection Content.
121. **PASS** — No deployed rule.
122. **PASS** — No activated rule.
123. **PASS** — No deactivated rule.
124. **PASS** — No active exception.
125. **PASS** — No deleted Signal or Alert.
126. **PASS** — Automated provenance visible.
127. **PASS** — No-AI alternative present.
128. **PASS** — No auto-granted permission.

### Limits
129. **PASS** — No detailed screen rewrite.
130. **PASS** — No new Screen ID.
131. **PASS** — No complete object created.
132. **PASS** — No JSON Schema.
133. **PASS** — No final cardinality.
134. **PASS** — No final object state machine.
135. **PASS** — No atomic permission.
136. **PASS** — No final namespace.
137. **PASS** — No API.
138. **PASS** — No protocol.
139. **PASS** — No engine or language selected.
140. **PASS** — No real command.
141. **PASS** — No product code.
142. **PASS** — No CAP-INV-5xx.
143. **PASS** — 4B.3A.2 and 4B.3B not started.

### Registers, quality and publication
144. **PASS** — Capability Register updated.
145. **PASS** — Dependency Register updated.
146. **PASS** — Object Consumption Map updated.
147. **PASS** — Action Classification updated.
148. **PASS** — Automation and AI Model updated.
149. **PASS** — Cross-product Links updated.
150. **PASS** — Screen Capability Map updated.
151. **PASS** — Requirements Matrix updated.
152. **PASS** — Baseline updated.
153. **PASS** — STATUS updated.
154. **PASS** — CHANGELOG and PR description coherent.
155. **PASS** — No active placeholder.
156. **PASS** — No targeted empty file.
157. **PASS** — No targeted broken link.
158. **PASS** — No concurrent owner or active duplicate.
159. **PASS** — Conformance report published.
160. **PASS** — Commits reachable and construction/remote SHAs identical.

## Publication evidence required
The PR description and final execution report must record the exact five commit SHAs, final construction/remote SHA, merge base, ahead/behind, PR Draft/open/unmerged state, repository visibility, root README equality, unchanged `main`, and CI/workflow result.
