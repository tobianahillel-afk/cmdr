---
id: studio-std2-workflow-builder-orchestration-capability-conformance
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-10
source-of-truth: quality-report
---
# Studio STD-2 — Workflow Builder & Orchestration Capability Conformance

Parent: **Delivery Roadmap Phase 5 — Studio and Endpoint**.  
Domain: **Studio**.  
Execution lot: **STD-2 — Workflow Builder & Orchestration**.  
Namespace: **CAP-STD-***.  
STD-2 is not a roadmap phase and creates no Phase 5A/5B.

## Structural result
- capability files: **17/17** — `CAP-STD-017..033`;
- numbered sections: **459/459**;
- mandatory tables: **102/102**;
- empty mandatory tables: **0**;
- generic mandatory tables: **0**;
- duplicate/recycled IDs: **0 / 0**;
- owner conflicts: **0**;
- new Screen IDs / detailed rewrites: **0 / 0**;
- Endpoint capabilities: **0**;
- STD-3/STD-4 capabilities: **0 / 0**;
- APIs / protocols / code / orchestration runtime / final language / final JSON Schema / final RBAC: **0**.

## Metrics
- before STD-2: **333 capabilities / 331 defined / 2 proposed / 333 planned / 8991 sections / 1998 tables**;
- after STD-2 content: **350 capabilities / 348 defined / 2 proposed / 350 planned / 9450 sections / 2100 tables**;
- Studio: **33 capabilities = 16 STD-1 + 17 STD-2**;
- Requirements remain **122 = 99 conform / 20 partial / 3 absent / 0 contradictory**;
- OPEN remains **18**;
- Command 27 PASS / Investigate 243 PASS / Govern 47 PASS preserved.

## Ownership and scope
Studio owns Workflow/Builder/orchestration-definition semantics. Govern retains Playbook, Approval, Decision, Response Run and Result; Shared retains generic Jobs/Trace/Activity/Search/Versioning/Recovery; Settings retains providers/integrations/secrets/environments; Endpoint retains technical primitives. Workflow definitions never grant runtime authority.

## Build-time gate status
Build SHA is not final until the fifth functional commit is published. Publication-dependent checks remain pending by design.

1. **PASS** — repo correct
2. **PASS** — visibility
3. **PASS** — branch
4. **PASS** — PR #2
5. **PASS** — base
6. **PASS** — open
7. **PASS** — Draft
8. **PASS** — unmerged
9. **PASS** — auto-merge
10. **PASS** — exact baseline
11. **PASS** — STD-1 history preserved
12. **PASS** — README branch
13. **PASS** — README main
14. **PASS** — main unchanged
15. **PASS** — roadmap Phase 5
16. **PASS** — roadmap ID
17. **PASS** — no 5A
18. **PASS** — no 5B
19. **PASS** — STD-2 execution lot only
20. **PASS** — linear publication
21. **PASS** — governance
22. **PASS** — Capability Register
23. **PASS** — Object Register
24. **PASS** — Dependency Register
25. **PASS** — Permission Register
26. **PASS** — OPEN
27. **PASS** — Requirements
28. **PASS** — baseline
29. **PASS** — Phase 5 preflight
30. **PASS** — STD-1 report
31. **PASS** — CAP-STD-001..016
32. **PASS** — Workflow docs
33. **PASS** — Builder docs
34. **PASS** — Tool docs
35. **PASS** — Tool Call docs
36. **PASS** — Skill docs
37. **PASS** — Human Gate docs
38. **PASS** — Automation Run docs boundary
39. **PASS** — Agent docs boundary
40. **PASS** — Evaluations/Simulation boundary
41. **PASS** — Govern
42. **PASS** — Settings
43. **PASS** — Shared
44. **PASS** — Studio screens
45. **PASS** — migration sources
46. **PASS** — CAP-STD namespace
47. **PASS** — 001..016 preserved
48. **PASS** — 017..033 availability
49. **PASS** — no recycle
50. **PASS** — count justified
51. **PASS** — one file/cap
52. **PASS** — owner
53. **PASS** — users
54. **PASS** — problem
55. **PASS** — goals
56. **PASS** — non-goals
57. **PASS** — inputs
58. **PASS** — objects read
59. **PASS** — objects modified
60. **PASS** — actions
61. **PASS** — states
62. **PASS** — outputs
63. **PASS** — transitions
64. **PASS** — source of truth
65. **PASS** — provenance
66. **PASS** — permissions
67. **PASS** — limits
68. **PASS** — errors
69. **PASS** — metrics
70. **PASS** — no-AI
71. **PASS** — GWT
72. **PASS** — Requirements
73. **PASS** — OPEN
74. **PASS** — six tables
75. **PASS** — 459 sections / 102 tables if 17
76. **PASS** — Studio owns Workflow
77. **PASS** — Studio owns Builder semantics
78. **PASS** — Studio owns Workflow Version
79. **PASS** — Studio owns orchestration graph
80. **PASS** — Studio owns conditions/branches
81. **PASS** — Studio owns mappings
82. **PASS** — Studio owns retry semantics
83. **PASS** — Studio owns compensation semantics
84. **PASS** — Govern retains Playbook
85. **PASS** — Govern retains Approval
86. **PASS** — Govern retains Decision
87. **PASS** — Govern retains Response Run
88. **PASS** — Shared retains Jobs
89. **PASS** — Shared retains Trace
90. **PASS** — Settings retains secrets
91. **PASS** — Workflow != Tool
92. **PASS** — Workflow != Skill
93. **PASS** — Workflow != Agent
94. **PASS** — Workflow != Automation Run
95. **PASS** — Workflow != Playbook
96. **PASS** — Definition != Version
97. **PASS** — Builder Session != Workflow
98. **PASS** — draft != saved version
99. **PASS** — saved != published
100. **PASS** — published != deployed
101. **PASS** — node != Tool
102. **PASS** — Tool step != Tool Call
103. **PASS** — Skill step != Skill
104. **PASS** — branch != Decision
105. **PASS** — condition != Policy
106. **PASS** — Human Gate != Approval
107. **PASS** — Human Gate != Decision
108. **PASS** — variable definition != value
109. **PASS** — Secret Reference != secret
110. **PASS** — mapping != source mutation
111. **PASS** — retry != authorization renewal
112. **PASS** — retry != infinite retry
113. **PASS** — idempotency != exactly once
114. **PASS** — compensation != Govern rollback
115. **PASS** — partial success != success
116. **PASS** — validation != execution
117. **PASS** — static validation != runtime success
118. **PASS** — compatible != authorized
119. **PASS** — graph valid != deployable
120. **PASS** — Tool Call != Automation Run
121. **PASS** — Automation Run != Response Run
122. **PASS** — Workflow != Endpoint primitive
123. **PASS** — Subworkflow ref != copy
124. **PASS** — version change visible
125. **PASS** — provenance preserved
126. **PASS** — Workflow definition
127. **PASS** — Builder session
128. **PASS** — inputs/outputs/variables
129. **PASS** — graph
130. **PASS** — Tool/Skill steps
131. **PASS** — conditions/branches
132. **PASS** — mapping
133. **PASS** — subworkflows
134. **PASS** — parallelism
135. **PASS** — errors
136. **PASS** — retry/idempotency
137. **PASS** — partial success/compensation
138. **PASS** — Human Gate
139. **PASS** — validation/readiness
140. **PASS** — versioning
141. **PASS** — pre-publish lifecycle
142. **PASS** — provenance/contracts
143. **PASS** — AI optional
144. **PASS** — no mandatory chatbot
145. **PASS** — no auto publish
146. **PASS** — no auto deploy
147. **PASS** — no auto execute
148. **PASS** — no auto Approval
149. **PASS** — no auto Decision
150. **PASS** — no permission invention
151. **PASS** — no raw secret
152. **PASS** — no hidden errors
153. **PASS** — no silent scope expansion
154. **PASS** — no API
155. **PASS** — no protocol
156. **PASS** — no code
157. **PASS** — no orchestration runtime selected
158. **PASS** — no final language
159. **PASS** — no JSON Schema
160. **PASS** — no final RBAC
161. **PASS** — permission anomaly preserved
162. **PASS** — no Screen ID
163. **PASS** — no detailed screen rewrite
164. **PASS** — no STD-3/4 capability
165. **PASS** — no Endpoint capability
166. **PASS** — Capability Register
167. **PASS** — Studio shard
168. **PASS** — Dependency Register
169. **PASS** — Object Map
170. **PASS** — Action Classification
171. **PASS** — Automation/AI
172. **PASS** — cross-links
173. **PASS** — screen map
174. **PASS** — Requirements
175. **PASS** — baseline
176. **PASS** — OPEN
177. **PASS** — STATUS
178. **PENDING-DOCUMENTARY** — CHANGELOG — canonical CHANGELOG entry will be recorded with post-publication evidence
179. **PASS** — roadmap
180. **PASS** — PR
181. **PASS** — STD-1 16 capabilities intact
182. **PASS** — STD-1 432 sections intact
183. **PASS** — STD-1 96 tables intact
184. **PASS** — Command PASS intact
185. **PASS** — Investigate/Govern PASS intact
186. **PASS** — conformance report
187. **PASS** — build-time status recorded
188. **PASS** — remote-dependent gates remain pending before publication
189. **PENDING-REMOTE** — five functional commits reachable
190. **PENDING-REMOTE** — remote verification actually executed
191. **PENDING-REMOTE** — final remote SHA recorded
192. **PENDING-REMOTE** — post-publication evidence linked from report
193. **PENDING-REMOTE** — build SHA recorded
194. **PENDING-REMOTE** — PR remains Draft
195. **PENDING-REMOTE** — README unchanged
196. **PENDING-REMOTE** — main unchanged
197. **PASS** — metrics recalculated
198. **PASS** — 17 new caps if final list
199. **PASS** — Endpoint remains 0
200. **PASS** — final parent statuses justified

Build-time verdict: **191 PASS / 9 PENDING / 0 FAIL**. STD-2 remains **PARTIAL / PENDING POST-PUBLICATION VERIFICATION** until the fifth functional commit and canonical post-publication evidence are verified.

## Post-publication verification
Pending. After publication, this same report must record the exact fifth functional SHA, baseline ancestry, PR/README/main/CI checks, canonical CHANGELOG evidence and final gate result. No final 200/200 verdict is valid until this section is updated.
