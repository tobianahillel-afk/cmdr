---
id: studio-std3-agents-human-gates-runtime-control-capability-conformance
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-10
source-of-truth: quality-report
---
# Studio STD-3 — Agents, Human Gates & Runtime Control Capability Conformance

Parent: **Delivery Roadmap Phase 5 — Studio and Endpoint**.  
Execution lot: **STD-3 — Agents, Human Gates & Runtime Control**.  
Namespace: **CAP-STD-***.  
STD-3 is not a Roadmap Phase or Capability Specification Phase; no Phase 5A/B/C/D is created.

## Structural result
- capability files: **18/18** — `CAP-STD-034..051`;
- numbered sections: **486/486**;
- mandatory tables: **108/108**;
- empty/generic mandatory tables: **0/0**;
- duplicate/recycled IDs: **0/0**;
- owner conflicts: **0**;
- Given/When/Then scenarios across STD-3 contracts: **59**, including all 30 mandatory scenario families;
- new Screen IDs / detailed rewrites: **0/0**;
- Endpoint capabilities: **0**;
- STD-4 capabilities: **0**;
- runtime/scheduler/agent framework/model-provider/API/protocol/product code/final JSON Schema/final RBAC: **0**.

## Metrics
Before STD-3: **350 capabilities / 348 defined / 2 proposed / 350 planned / 9450 sections / 2100 tables**.  
After STD-3 content: **368 capabilities / 366 defined / 2 proposed / 368 planned / 9936 sections / 2208 tables**.  
Studio after STD-3 content: **51 capabilities / 1377 sections / 306 tables**.  
Requirements remain **122 = 99 conform / 20 partial / 3 absent / 0 contradictory**. OPEN remains **18**.

## Ownership and safety
Studio owns Automation Agent/Agent Team/Human Gate/Automation Run/runtime-control/Control Room semantics. Govern retains Approval, Decision, Response Playbook, Response Run, Result and production authority. Settings retains identities/roles, providers, integrations, credentials, raw secrets, tenants/environments and runtime administration. Shared retains generic Jobs/queues/scheduling/Trace/Activity/Notifications/Recovery. Endpoint retains endpoint technical primitives/local runtime.

Key non-equivalence is explicit throughout: Agent objective/role/access/proposal/plan are not authorization; Human Gate is not Approval/Decision; Automation Run is not Workflow/Tool Call/Job/Response Run; Run request states are not confirmations; retry attempt is not a Run; idempotency is not exactly-once; partial completion is not success; Studio outcome is not Govern Result/Evidence/Finding; Control Room is not Command Work Queue or Govern Runs & Rollback.

## Build-time verification
The fifth functional commit is the build publication boundary. At build time, remote-dependent gates and the canonical CHANGELOG/PR metadata remain pending intentionally. No remote fact is pre-declared.

### Git / roadmap — 1–20
1. **PASS** — repository correct
2. **PASS** — visibility
3. **PASS** — branch
4. **PASS** — PR #2
5. **PASS** — base main
6. **PASS** — PR open
7. **PASS** — Draft
8. **PASS** — unmerged
9. **PASS** — auto-merge disabled
10. **PASS** — exact baseline
11. **PASS** — STD-1 history preserved
12. **PASS** — STD-2 history preserved
13. **PASS** — README branch
14. **PASS** — README main
15. **PASS** — main unchanged
16. **PASS** — roadmap Phase 5
17. **PASS** — roadmap ID
18. **PASS** — no Phase 5A/B/C/D
19. **PASS** — STD-3 execution lot only
20. **PASS** — linear publication

### Sources — 21–50
21. **PASS** — governance
22. **PASS** — Capability Register
23. **PASS** — Object Register
24. **PASS** — Dependency Register
25. **PASS** — Permission Register
26. **PASS** — OPEN register
27. **PASS** — Requirements
28. **PASS** — baseline
29. **PASS** — Phase 5 preflight
30. **PASS** — STD-1 report
31. **PASS** — STD-2 report
32. **PASS** — CAP-STD-001..016
33. **PASS** — CAP-STD-017..033
34. **PASS** — Automation Agent docs
35. **PASS** — Agent Team docs
36. **PASS** — Workflow docs
37. **PASS** — Tool docs
38. **PASS** — Skill docs
39. **PASS** — Human Gate docs
40. **PASS** — Automation Run docs
41. **PASS** — Control Room docs
42. **PASS** — Runtime docs
43. **PASS** — Evaluations docs boundary
44. **PASS** — Simulations docs boundary
45. **PASS** — Deployment docs boundary
46. **PASS** — Govern docs
47. **PASS** — Settings docs
48. **PASS** — Shared docs
49. **PASS** — Endpoint docs boundary
50. **PASS** — Studio screens/migrations

### Capability/template — 51–80
51. **PASS** — CAP-STD namespace
52. **PASS** — 001..033 preserved
53. **PASS** — 034..051 availability
54. **PASS** — no recycle
55. **PASS** — final count justified
56. **PASS** — one canonical file each
57. **PASS** — owner
58. **PASS** — users
59. **PASS** — problem
60. **PASS** — goals
61. **PASS** — non-goals
62. **PASS** — inputs
63. **PASS** — objects read
64. **PASS** — objects modified
65. **PASS** — actions
66. **PASS** — states
67. **PASS** — outputs
68. **PASS** — transitions
69. **PASS** — source of truth
70. **PASS** — provenance
71. **PASS** — permissions
72. **PASS** — limits
73. **PASS** — errors
74. **PASS** — metrics
75. **PASS** — no-AI
76. **PASS** — ≥3 GWT
77. **PASS** — Requirements
78. **PASS** — OPEN
79. **PASS** — six tables
80. **PASS** — 486 sections / 108 tables if 18

### Ownership/distinctions — 81–135
81. **PASS** — Studio owns Automation Agent semantics
82. **PASS** — Studio owns Human Gate
83. **PASS** — Studio owns Automation Run
84. **PASS** — Studio owns Control Room semantics
85. **PASS** — Govern retains Approval
86. **PASS** — Govern retains Decision
87. **PASS** — Govern retains Response Run
88. **PASS** — Govern retains Result
89. **PASS** — Endpoint retains technical primitives
90. **PASS** — Settings retains secrets/providers
91. **PASS** — Shared retains Jobs/Trace
92. **PASS** — Agent != Workflow
93. **PASS** — Agent != Skill
94. **PASS** — Agent != Tool
95. **PASS** — Agent != Automation Run
96. **PASS** — Agent != human user
97. **PASS** — Agent Team != human Team
98. **PASS** — Agent role != authorization
99. **PASS** — objective != permission
100. **PASS** — Tool access != ownership
101. **PASS** — Tool eligibility != invocation
102. **PASS** — proposal != action
103. **PASS** — plan != Workflow Definition
104. **PASS** — plan != authorization
105. **PASS** — Agent Team coordination != graph
106. **PASS** — Agent delegation != permission delegation
107. **PASS** — Human Gate != Approval
108. **PASS** — Human Gate != Decision
109. **PASS** — Human Gate response != production authority
110. **PASS** — Human Gate reviewer != Govern approver
111. **PASS** — Human Gate complete != authorization
112. **PASS** — Automation Run != Workflow
113. **PASS** — Automation Run != Workflow Version
114. **PASS** — Automation Run != Tool Call
115. **PASS** — Automation Run != Job
116. **PASS** — Automation Run != Response Run
117. **PASS** — Run created != started
118. **PASS** — queued != started
119. **PASS** — start requested != running
120. **PASS** — running != success
121. **PASS** — paused != stopped
122. **PASS** — stop request != stopped
123. **PASS** — cancel != rollback
124. **PASS** — retry != new Run
125. **PASS** — attempt != Run
126. **PASS** — Workflow retry semantics != runtime attempt
127. **PASS** — idempotency != exactly-once
128. **PASS** — Tool Call success != Run success
129. **PASS** — step success != Run success
130. **PASS** — partial completion != success
131. **PASS** — Studio outcome != Govern Result
132. **PASS** — compensation != Govern rollback
133. **PASS** — Control Room != Command Work Queue
134. **PASS** — Control Room != Govern Runs & Rollback
135. **PASS** — AI recommendation != authority

### Functional coverage — 136–153
136. **PASS** — Agent definition
137. **PASS** — objectives/constraints
138. **PASS** — Tool/Skill access
139. **PASS** — Agent Team
140. **PASS** — planning/bounded autonomy
141. **PASS** — oversight/intervention
142. **PASS** — Human Gate request
143. **PASS** — Human Gate lifecycle/boundary
144. **PASS** — Run creation/context
145. **PASS** — Run lifecycle
146. **PASS** — Run steps/attempts
147. **PASS** — queue/scheduling/concurrency
148. **PASS** — start/pause/resume/stop/cancel
149. **PASS** — error/timeout/retry/partial
150. **PASS** — runtime state/context
151. **PASS** — Control Room
152. **PASS** — Run outcome/handoff
153. **PASS** — runtime provenance

### AI / security / limits — 154–180
154. **PASS** — AI optional
155. **PASS** — no mandatory chatbot
156. **PASS** — no auto permission
157. **PASS** — no autonomous scope expansion
158. **PASS** — no automatic Govern Approval
159. **PASS** — no automatic Govern Decision
160. **PASS** — no authority via Human Gate
161. **PASS** — no raw secrets
162. **PASS** — no silent retry
163. **PASS** — no infinite retry
164. **PASS** — no hidden errors
165. **PASS** — no silent lost status
166. **PASS** — no automatic Result
167. **PASS** — no automatic Evidence/Finding
168. **PASS** — no runtime selected
169. **PASS** — no scheduler selected
170. **PASS** — no agent framework selected
171. **PASS** — no model/provider selected
172. **PASS** — no API
173. **PASS** — no protocol
174. **PASS** — no product code
175. **PASS** — no JSON Schema
176. **PASS** — no final RBAC
177. **PASS** — permission namespace anomaly preserved
178. **PASS** — no new Screen ID
179. **PASS** — no detailed screen rewrite
180. **PASS** — no Endpoint capability

### Registers / non-regression — 181–200
181. **PASS** — Capability Register
182. **PASS** — Studio shard
183. **PASS** — Dependency Register
184. **PASS** — Object Map
185. **PASS** — Action Classification
186. **PASS** — Automation/AI model
187. **PASS** — Cross-product links
188. **PASS** — Screen map
189. **PASS** — Requirements
190. **PASS** — baseline
191. **PASS** — OPEN audit
192. **PASS** — STATUS
193. **PENDING-REMOTE** — canonical CHANGELOG
194. **PASS** — roadmap
195. **PENDING-REMOTE** — PR description
196. **PASS** — STD-1 intact
197. **PASS** — STD-2 intact
198. **PASS** — Command intact
199. **PASS** — Investigate intact
200. **PASS** — Govern intact

### Publication — 201–210
201. **PASS** — conformance report
202. **PASS** — build-time status explicit
203. **PASS** — remote gates pending until actual verification
204. **PENDING-REMOTE** — five functional commits reachable
205. **PENDING-REMOTE** — remote verification executed
206. **PENDING-REMOTE** — exact build SHA recorded
207. **PENDING-REMOTE** — exact final SHA recorded
208. **PENDING-REMOTE** — post-publication evidence inside/linked from canonical report
209. **PENDING-REMOTE** — PR remains Draft/main/README intact
210. **PASS** — STD-4 and Endpoint remain NOT STARTED

### Build-time verdict
**202 PASS / 8 PENDING-REMOTE / 0 FAIL.**

Pending gates are exactly **193, 195, 204–209**. Until those are verified against the published fifth functional commit, STD-3 remains **PARTIAL / PENDING POST-PUBLICATION VERIFICATION**.

## Post-publication verification
To be completed after the fifth functional commit is published. This section will record the exact five functional SHAs, build SHA, ancestry, PR/README/main/CI state, remote structural checks and final gate verdict. It must not be inferred from build-time evidence.

## Stop line
STD-4 and Endpoint remain **NOT STARTED**. Documentary STD-3 coverage proves no implemented runtime.