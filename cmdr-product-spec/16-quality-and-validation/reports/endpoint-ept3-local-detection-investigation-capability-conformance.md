---
id: endpoint-ept3-local-detection-investigation-capability-conformance
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-11
source-of-truth: quality-report
requirements: [REQ-PROD-006, REQ-PROD-012, REQ-PROD-018, REQ-PROD-019, REQ-INV-001, REQ-INV-006, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004]
open_decisions: [OPEN-008, OPEN-015, OPEN-017]
---
# Endpoint EPT-3 — Local Detection and Endpoint Investigation — Capability Conformance

## Scope
This report validates only **EPT-3 — Local Detection and Endpoint Investigation** under Delivery Roadmap Phase 5 — Studio and Endpoint. EPT-4, EPT-5 and EPT-6 remain NOT STARTED.

## Baseline and source audit
Starting baseline: `5d7c037aff6004984416665e7e188a8700e62b2f` — `docs: record Endpoint EPT-2 post-publication verification`.

All **8/8 Detection** and **8/8 Investigation** Endpoint documents were read before allocation, together with EPT-1/EPT-2 contracts, Investigate Detection Engineering, Case/Evidence/Finding, Command Detection/Signal, Shared Timeline/Linking/Search/Correlation, Settings and Govern boundaries. `OPEN-008` and `OPEN-017` remain open.

## Capability set and structure
`CAP-EPT-031..046`: **16 capability files / 432 numbered sections / 96 mandatory tables / at least 48 Given/When/Then scenarios**, all `draft / defined / planned`.

- duplicate IDs: **0**;
- recycled IDs: **0**;
- owner conflicts: **0**;
- empty/generic mandatory tables: **0**;
- Endpoint Screen IDs added: **0**.

Endpoint cumulative: **46 capabilities / 1242 sections / 276 mandatory tables**. Global content: **431 capabilities / 429 defined / 2 proposed / 431 planned / 11637 sections / 2586 mandatory tables**. Requirements remain **122 = 99 conform / 20 partial / 3 absent / 0 contradictory**. OPEN remains **18**.

## Ownership and semantic invariants
Investigate retains Detection Engineering, Case, Evidence and Finding. Command retains canonical Detection, Signal, Alert and Incident. Endpoint owns only local technical eligibility/evaluation/match/candidate/context/investigation projections. Settings retains administrative sources/runtimes/targets/policies/assignments. Shared retains generic Search/Timeline/Linking/Correlation/Trace/Activity/Jobs/Reporting. Govern retains Decision/Approval/Response Run/Result and response authority. Studio retains Tool/Tool Call/Skill/Automation Run.

EPT-3 explicitly preserves: Detection Content ≠ Match ≠ local candidate; eligible ≠ enabled ≠ applicable ≠ matched; match ≠ malicious verdict/Finding/Evidence/Incident; local candidate ≠ canonical Command Signal/Finding/Evidence/Case/Result/authority; severity ≠ impact; confidence ≠ certainty; rationale ≠ proof; multiple matches ≠ incident; no match ≠ benign; suppressed ≠ deleted; grouped ≠ identical; deduplicated ≠ erased; detection health ≠ endpoint health; coverage ≠ complete visibility; healthy runtime ≠ full coverage; missing signal ≠ benign; Detection ≠ Response; evaluation ≠ side effect; Endpoint Investigation ≠ Case; context ≠ Evidence; summary ≠ Finding; process ancestry/relation/timeline/correlation ≠ causal or malicious proof; file context ≠ acquisition; network/auth/system observations ≠ compromise/persistence conclusions; pivot/context expansion ≠ Collection; investigation query ≠ Live Response; content version ≠ Agent version; Endpoint local detection ≠ SIEM engine.

## Build-time verification — 210 mandatory gates
The authoritative gate ledger is preserved by number and label below. Build-time state is **204 PASS / 6 PENDING-BUILD-OR-REMOTE / 0 FAIL**. Gates **204–209** are explicitly publication-dependent; gate 210 is already PASS.

### Gates 1–20 — Git / roadmap — PASS
1 repo correct; 2 visibility; 3 branch; 4 PR #2; 5 base main; 6 open; 7 Draft; 8 unmerged; 9 auto-merge disabled; 10 exact baseline; 11 EPT-1 PASS; 12 EPT-2 PASS; 13 README branch; 14 README main; 15 main unchanged; 16 Studio PASS; 17 roadmap Phase 5; 18 roadmap ID; 19 EPT-3 only; 20 no Phase 5E3.

### Gates 21–45 — Sources — PASS
21 governance; 22 Capability Register; 23 Object Register; 24 Dependency Register; 25 Permission Register; 26 Screen Register; 27 OPEN; 28 Requirements; 29 EPT-1 report; 30 EPT-2 report; 31 Endpoint README; 32 IA; 33 all Detection docs; 34 all Investigation docs; 35 Telemetry boundaries; 36 Process telemetry; 37 File telemetry; 38 Network telemetry; 39 Auth telemetry; 40 System telemetry; 41 Investigate Evidence/Finding/Case; 42 Detection Engineering; 43 Shared boundaries; 44 Settings/Govern/Studio boundaries; 45 migration sources.

### Gates 46–75 — Capability/template — PASS
46 namespace; 47 `CAP-EPT-001..030` preserved; 48 `031+` availability; 49 no recycle; 50 count justified; 51 one file/cap; 52 owner; 53 users; 54 problem; 55 goals; 56 non-goals; 57 inputs; 58 objects read; 59 objects modified; 60 actions; 61 states; 62 outputs; 63 transitions; 64 SOT; 65 provenance; 66 permissions; 67 errors; 68 limits; 69 metrics; 70 no-AI; 71 GWT; 72 Requirements; 73 OPEN; 74 six tables; 75 432 sections / 96 tables.

### Gates 76–140 — Ownership/distinctions — PASS
76 Endpoint local detection semantics; 77 Investigate owns Evidence; 78 Investigate owns Finding; 79 Investigate owns Case; 80 Govern owns Result; 81 Govern owns response authority; 82 Shared generic Search/Timeline preserved; 83 Settings admin boundary preserved; 84 Detection Content != Match; 85 Content != Signal; 86 eligible != enabled; 87 enabled != applicable; 88 applicable != matched; 89 match != malicious verdict; 90 match != Finding; 91 match != Evidence; 92 match != Incident; 93 signal != Finding; 94 signal != Evidence; 95 signal != Case; 96 signal != Result; 97 signal != authority; 98 severity != impact; 99 confidence != certainty; 100 rationale != proof; 101 multiple matches != incident; 102 no match != benign; 103 suppressed != deleted; 104 grouped != identical; 105 dedup != erase; 106 detection health != endpoint health; 107 coverage != complete visibility; 108 runtime healthy != full coverage; 109 missing signal != benign; 110 Detection != Response; 111 evaluation != side effect; 112 Endpoint Investigation != Case; 113 context != Evidence; 114 summary != Finding; 115 ancestry != causal proof; 116 relation != malicious; 117 hash hit != malware proof; 118 file context != acquisition; 119 network context != compromise; 120 reputation hit != compromise; 121 auth anomaly != account compromise; 122 system observation != persistence proof; 123 timeline != causal proof; 124 correlation != causation; 125 pivot != Collection; 126 context expansion != acquisition; 127 investigation query != Live Response; 128 content version != Agent version; 129 local detection != SIEM engine; 130 OPEN-017 preserved; 131 OPEN-008 preserved; 132 no final detection language; 133 no final runtime; 134 no owner conflict; 135 EPT-2 observations reused; 136 CAP-EPT-011/027/028 not duplicated; 137 local summary handoff correct; 138 Investigate qualification retained; 139 Govern boundary retained; 140 provenance retained.

### Gates 141–156 — Functional coverage — PASS
141 content consumption; 142 eligibility; 143 evaluation/match; 144 signal lifecycle; 145 context/severity/confidence; 146 grouping/suppression boundary; 147 coverage/health/gaps; 148 process investigation; 149 file investigation; 150 network investigation; 151 user/session investigation justified; 152 system investigation justified; 153 timeline/correlation; 154 pivots/context retrieval; 155 detection-investigation handoff; 156 provenance.

### Gates 157–180 — AI/security/limits — PASS
157 AI optional; 158 no malicious authoritative verdict; 159 no Evidence creation; 160 no Finding creation; 161 no Case creation; 162 no Incident creation; 163 no response execution; 164 no raw secret; 165 sensitive context controls; 166 tenant isolation; 167 no API; 168 no protocol; 169 no final rule language; 170 no final runtime; 171 no final ML model/runtime; 172 no physical schema; 173 no collection; 174 no Live Response; 175 no containment; 176 no final RBAC; 177 no Screen ID; 178 no detailed UX; 179 no EPT-4/5/6; 180 no implementation.

### Gates 181–200 — Registers / non-regression — PASS
181 Capability Register; 182 Endpoint shard; 183 Dependency Register; 184 Object Map; 185 Action Classification; 186 Automation/AI; 187 cross-links; 188 Requirements; 189 OPEN; 190 STATUS; 191 CHANGELOG; 192 roadmap; 193 quality/validation; 194 IA; 195 CAP-EPT-001..030 intact; 196 EPT-1 190 preserved; 197 EPT-2 200 preserved; 198 Command/Investigate/Govern intact; 199 Studio intact; 200 preflight preserved.

### Gates 201–210 — Publication
- **201 PASS** — conformance report exists.
- **202 PASS** — build-time state is explicit.
- **203 PASS** — remote gates are explicitly pending until publication.
- **204 PENDING-BUILD-OR-REMOTE** — five functional commits reachable.
- **205 PENDING-BUILD-OR-REMOTE** — remote verification executed.
- **206 PENDING-BUILD-OR-REMOTE** — exact build SHA.
- **207 PENDING-BUILD-OR-REMOTE** — exact final SHA.
- **208 PENDING-BUILD-OR-REMOTE** — post-publication evidence.
- **209 PENDING-BUILD-OR-REMOTE** — PR/main/README intact after publication.
- **210 PASS** — EPT-4/5/6 NOT STARTED.

## Build-time verdict
**204 PASS / 6 PENDING-BUILD-OR-REMOTE / 0 FAIL.** A single FAIL keeps EPT-3 PARTIAL. Final **210/210 PASS** can only be declared after exact five-commit reachability, publication, real remote verification, exact build/final SHA recording, canonical post-publication evidence and final PR/main/README recheck.

Documentary PASS never proves implementation or supported-platform delivery.