---
id: endpoint-ept5-containment-verification-governed-response-primitives-capability-conformance
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-11
source-of-truth: quality-report
requirements: [REQ-PROD-004, REQ-PROD-005, REQ-PROD-006, REQ-PROD-008, REQ-PROD-014, REQ-PROD-015, REQ-PROD-016, REQ-PROD-018, REQ-PROD-019, REQ-PROD-020, REQ-OBJ-007, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004]
open_decisions: [OPEN-007, OPEN-008, OPEN-013, OPEN-014, OPEN-015, OPEN-017]
---
# Endpoint EPT-5 — Containment, Verification and Governed Response Primitives — Capability Conformance

## Build scope
Exact baseline: `5d576295fa12693ef375a35cfe515d7bdf577f68`. Only EPT-5 is built. EPT-6 is NOT STARTED.

The complete source audit selected **17 independently justified capabilities `CAP-EPT-065..081`**, not a quota: the 16 candidates are supported and `CAP-EPT-081` is additionally required by the distinct local-session lock/termination responsibility in `containment/account-containment.md`; directory identity actions remain external.

## Structural result
- capability files: **17/17**;
- numbered sections: **459/459**;
- mandatory tables: **102/102**;
- GWT scenarios: **at least 51**;
- duplicate IDs: **0**;
- recycled IDs: **0**;
- owner conflicts: **0**;
- empty/generic mandatory tables: **0**;
- Endpoint Screen IDs added: **0**.

Endpoint cumulative: **81 / 2187 / 486**. Global content: **466 capabilities / 464 defined / 2 proposed / 466 planned / 12582 sections / 2796 mandatory tables**. Requirements remain **122 = 99 conform / 20 partial / 3 absent / 0 contradictory**. OPEN remains **18**.

## Build-time gate state
**224 PASS / 6 PENDING-BUILD-OR-REMOTE / 0 FAIL.** Gates **224–229** are publication-dependent and remain PENDING before branch publication. Gate 230 is PASS because EPT-6 is untouched.

## Gate ledger
| # | Gate | Build status |
|---:|---|---|
| 1 | repo correct | PASS |
| 2 | visibility | PASS |
| 3 | branch | PASS |
| 4 | PR #2 | PASS |
| 5 | base main | PASS |
| 6 | open | PASS |
| 7 | Draft | PASS |
| 8 | unmerged | PASS |
| 9 | auto-merge disabled | PASS |
| 10 | exact baseline | PASS |
| 11 | EPT-1 PASS | PASS |
| 12 | EPT-2 PASS | PASS |
| 13 | EPT-3 PASS | PASS |
| 14 | EPT-4 PASS | PASS |
| 15 | README branch | PASS |
| 16 | README main | PASS |
| 17 | main unchanged | PASS |
| 18 | Studio PASS | PASS |
| 19 | EPT-5 only | PASS |
| 20 | no Phase 5E5 | PASS |
| 21 | governance | PASS |
| 22 | Capability Register | PASS |
| 23 | Object Register | PASS |
| 24 | Dependency Register | PASS |
| 25 | Permission Register | PASS |
| 26 | OPEN | PASS |
| 27 | Requirements | PASS |
| 28 | EPT-1 report | PASS |
| 29 | EPT-2 report | PASS |
| 30 | EPT-3 report | PASS |
| 31 | EPT-4 report | PASS |
| 32 | Endpoint README | PASS |
| 33 | IA | PASS |
| 34 | all Containment docs | PASS |
| 35 | effectful historical Live Response sources | PASS |
| 36 | Govern Action Request | PASS |
| 37 | Govern Approval | PASS |
| 38 | Govern Decision | PASS |
| 39 | Govern Response Run | PASS |
| 40 | Govern Result | PASS |
| 41 | Govern verification | PASS |
| 42 | Govern rollback | PASS |
| 43 | Studio Human Gate/Open-015 | PASS |
| 44 | Settings Policy/admin boundaries | PASS |
| 45 | EPT-6 Resilience/Security boundary | PASS |
| 46 | CAP-EPT namespace | PASS |
| 47 | 001..064 preserved | PASS |
| 48 | 065+ availability | PASS |
| 49 | no recycle | PASS |
| 50 | final count justified | PASS |
| 51 | one file/capability | PASS |
| 52 | owner | PASS |
| 53 | users | PASS |
| 54 | problem | PASS |
| 55 | goals | PASS |
| 56 | non-goals | PASS |
| 57 | inputs | PASS |
| 58 | objects read | PASS |
| 59 | objects modified | PASS |
| 60 | actions | PASS |
| 61 | states | PASS |
| 62 | outputs | PASS |
| 63 | transitions | PASS |
| 64 | source-of-truth | PASS |
| 65 | provenance | PASS |
| 66 | permissions | PASS |
| 67 | errors | PASS |
| 68 | limits | PASS |
| 69 | metrics | PASS |
| 70 | no-AI | PASS |
| 71 | GWT | PASS |
| 72 | Requirements | PASS |
| 73 | OPEN | PASS |
| 74 | six mandatory tables | PASS |
| 75 | recalculated structure for justified final count | PASS — 17 → 459/102 |
| 76 | Govern owns authority | PASS |
| 77 | Govern owns Approval | PASS |
| 78 | Govern owns Decision | PASS |
| 79 | Govern owns Response Run | PASS |
| 80 | Govern owns Result | PASS |
| 81 | Govern owns response verification | PASS |
| 82 | Govern owns response rollback | PASS |
| 83 | Endpoint owns technical primitives | PASS |
| 84 | Endpoint owns target-side facts | PASS |
| 85 | Settings owns policy/admin config | PASS |
| 86 | Studio Human Gate preserved | PASS |
| 87 | Shared generic mechanisms preserved | PASS |
| 88 | primitive != Action Request | PASS |
| 89 | eligible != authorized | PASS |
| 90 | available != allowed | PASS |
| 91 | Decision != technical command | PASS |
| 92 | Response Run != Endpoint execution | PASS |
| 93 | requested state != achieved state | PASS |
| 94 | precheck PASS != universally safe | PASS |
| 95 | started != effective | PASS |
| 96 | technically effective != incident contained | PASS |
| 97 | terminate request != terminated | PASS |
| 98 | process absent != threat resolved | PASS |
| 99 | suspended != terminated | PASS |
| 100 | resume != rollback | PASS |
| 101 | isolation != offline | PASS |
| 102 | isolation != network failure | PASS |
| 103 | network block != isolation | PASS |
| 104 | block != threat removed | PASS |
| 105 | quarantine requested != applied | PASS |
| 106 | quarantine != delete | PASS |
| 107 | delete != guaranteed recoverability | PASS |
| 108 | restored != safe | PASS |
| 109 | service stopped != remediation complete | PASS |
| 110 | technical success != Response Run success | PASS |
| 111 | technical failure != response failure automatically | PASS |
| 112 | technical output != Result | PASS |
| 113 | technical verification != Govern verification | PASS |
| 114 | desired state observed != business objective achieved | PASS |
| 115 | verification unavailable != technical failure automatically | PASS |
| 116 | unknown != fail automatically | PASS |
| 117 | partial != complete | PASS |
| 118 | drift != execution failure automatically | PASS |
| 119 | retry != rollback | PASS |
| 120 | reversal != Govern rollback | PASS |
| 121 | compensation != Response Run rollback | PASS |
| 122 | release != restored pre-incident state | PASS |
| 123 | connectivity restored != secure | PASS |
| 124 | cancel != rollback | PASS |
| 125 | stop request != stop confirmed | PASS |
| 126 | rollback requested != rollback complete | PASS |
| 127 | rollback technical success != Govern reconciliation | PASS |
| 128 | operator != approver | PASS |
| 129 | Endpoint Policy != authority | PASS |
| 130 | availability != permission | PASS |
| 131 | Human Gate != Approval | PASS |
| 132 | Human Gate != Decision | PASS |
| 133 | Endpoint action != Tool Call | PASS |
| 134 | Endpoint action != Automation Run | PASS |
| 135 | audit trail != crypto proof | PASS |
| 136 | OPEN-007 preserved | PASS |
| 137 | OPEN-008 preserved | PASS |
| 138 | OPEN-013 preserved | PASS |
| 139 | OPEN-014 preserved | PASS |
| 140 | OPEN-015 preserved | PASS |
| 141 | OPEN-017 preserved | PASS |
| 142 | no owner conflict | PASS |
| 143 | EPT-4 technical execution reused | PASS |
| 144 | Govern authority preserved | PASS |
| 145 | Govern verification preserved | PASS |
| 146 | Govern Result preserved | PASS |
| 147 | technical reversal boundary preserved | PASS |
| 148 | EPT-6 boundary preserved | PASS |
| 149 | provenance preserved | PASS |
| 150 | no false implementation | PASS |
| 151 | primitive request/eligibility | PASS |
| 152 | authority boundary | PASS |
| 153 | precheck/readiness | PASS |
| 154 | process control | PASS |
| 155 | host isolation | PASS |
| 156 | network control justified | PASS |
| 157 | quarantine | PASS |
| 158 | delete/restore justified | PASS |
| 159 | service/system control justified | PASS |
| 160 | execution outcome | PASS |
| 161 | verification request | PASS |
| 162 | target-state verification | PASS |
| 163 | partial/failure/drift | PASS |
| 164 | technical reversal | PASS |
| 165 | containment release | PASS |
| 166 | Govern reconciliation/provenance | PASS |
| 167 | AI optional | PASS |
| 168 | no AI authority | PASS |
| 169 | no AI Decision | PASS |
| 170 | no AI Approval | PASS |
| 171 | no auto effectful primitive | PASS |
| 172 | no fake target state | PASS |
| 173 | no fake verification | PASS |
| 174 | no fake Result | PASS |
| 175 | no provenance destruction | PASS |
| 176 | no raw secret | PASS |
| 177 | tenant isolation | PASS |
| 178 | sensitive target data controls | PASS |
| 179 | least privilege | PASS |
| 180 | step-up documented | PASS |
| 181 | SoD documented | PASS |
| 182 | no commands | PASS |
| 183 | no PowerShell | PASS |
| 184 | no shell | PASS |
| 185 | no API | PASS |
| 186 | no protocol | PASS |
| 187 | no physical schema | PASS |
| 188 | no final policy engine | PASS |
| 189 | no final verification engine | PASS |
| 190 | no final rollback engine | PASS |
| 191 | no final RBAC | PASS |
| 192 | no Screen ID | PASS |
| 193 | no detailed response UX | PASS |
| 194 | no EPT-6 | PASS |
| 195 | no implementation | PASS |
| 196 | Capability Register | PASS |
| 197 | Endpoint shard | PASS |
| 198 | Dependency Register | PASS |
| 199 | Object Map | PASS |
| 200 | Action Classification | PASS |
| 201 | Automation/AI | PASS |
| 202 | cross-links | PASS |
| 203 | Requirements | PASS |
| 204 | OPEN | PASS |
| 205 | STATUS build addendum | PASS |
| 206 | CHANGELOG EPT-5 | PASS |
| 207 | roadmap build addendum | PASS |
| 208 | quality/validation | PASS |
| 209 | IA | PASS |
| 210 | Containment docs/maps | PASS |
| 211 | Govern links | PASS |
| 212 | CAP-EPT-001..064 intact in build diff | PASS |
| 213 | EPT-1 preserved | PASS |
| 214 | EPT-2 preserved | PASS |
| 215 | EPT-3 preserved | PASS |
| 216 | EPT-4 preserved | PASS |
| 217 | Command preserved | PASS |
| 218 | Investigate/Govern preserved | PASS |
| 219 | Studio preserved | PASS |
| 220 | Endpoint preflight preserved | PASS |
| 221 | conformance report | PASS |
| 222 | build-time status explicit | PASS |
| 223 | remote-dependent gates pending before publication | PASS |
| 224 | five functional commits reachable from remote branch | PENDING-BUILD-OR-REMOTE |
| 225 | remote verification actually executed | PENDING-BUILD-OR-REMOTE |
| 226 | exact build SHA | PENDING-BUILD-OR-REMOTE |
| 227 | exact final SHA | PENDING-BUILD-OR-REMOTE |
| 228 | post-publication evidence | PENDING-BUILD-OR-REMOTE |
| 229 | PR/main/README intact after publication | PENDING-BUILD-OR-REMOTE |
| 230 | EPT-6 NOT STARTED | PASS |

## Mandatory GWT coverage
The 17 capability files each contain at least three Given/When/Then cases, covering the supplied minimums across authority absence, availability/staleness, process identity changes, suspend/terminate, isolation partiality, network partial enforcement, missing/changed quarantine target, delete authority, restore safety uncertainty, unsupported service mutation, verification mismatch/unavailability/staleness, drift, partial/unknown execution, Govern-requested reversal, partial reversal failure, containment release, Govern handoff/Result boundary, Human Gate non-equivalence, AI non-authority, cross-tenant denial, OPEN-013 and EPT-6 stop.

## Final build verdict
**EPT-5 remains PENDING POST-PUBLICATION VERIFICATION.** Documentary build is structurally conformant, but final PASS is forbidden until gates 224–229 close on the remotely published five-commit chain.