---
id: endpoint-ept1-enrollment-inventory-health-platform-foundations-capability-conformance
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-11
source-of-truth: quality-report
requirements: [REQ-PROD-006, REQ-PROD-012, REQ-PROD-017, REQ-PROD-018, REQ-PROD-019, REQ-OBJ-008, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-013, OPEN-015]
---
# Endpoint EPT-1 — Enrollment, Inventory, Health and Platform Foundations — Capability Conformance

Parent: **Delivery Roadmap Phase 5 — Studio and Endpoint**. Domain: **Endpoint Agent**. Execution lot: **EPT-1**. Capability namespace: **`CAP-EPT-*`**. EPT-1 is not a Roadmap Phase or Capability Specification Phase and creates no Phase 5E1/5E2/etc.

## Structural result
Exactly `CAP-EPT-001..014`: **14/14 files / 378/378 numbered sections / 84/84 mandatory tables / at least 42 GWT scenarios / 0 duplicate ID / 0 recycled ID / 0 empty mandatory table / 0 owner conflict**. Every capability is `draft / defined / planned` and uses the canonical 27-section contract.

## Scope covered
Identity/instance registration; local enrollment handoff/state; tenant/environment binding; observed platform/OS/architecture; Agent version/build/compatibility; inventory snapshot/metadata; inventory freshness/change; health/self-check; heartbeat/connectivity/last-seen; online/degraded/offline/stale/unknown/revoked semantics; technical capability advertisement/availability; Settings Fleet projection boundary; Settings Endpoint Policy effective-local-state boundary; foundations provenance/cross-product handoff.

## Mandatory non-equivalence
Endpoint Agent != Studio Automation Agent; Endpoint != Agent automatically; Device != Agent Instance; Host Identity != Agent Identity; registration != enrollment administration; local enrollment state != Settings administrative record; enrolled != online; registered != healthy; revoked != deleted; platform observed != supported; referenced != delivered; OS identified != compatible; compatible != supported; version observed != approved; inventory != Fleet; snapshot != global source of truth; metadata != Settings configuration; freshness != health; last seen != heartbeat; heartbeat != health; online != fully capable; degraded != offline; offline != revoked; stale != offline; health != security posture; self-check PASS != integrity proof; advertised != authorized; advertised != available; available != global support; technical capability != Studio Tool; Fleet != local state; Endpoint Policy != effective local state; assigned != applied; applied != healthy; local effective config != Settings object; technical result != Govern Result; execution != Response Run; technical execution != Tool Call; local rollback != Govern rollback; OPEN-008 candidate != delivery commitment.

## Technical/safety boundary
No detailed telemetry, detection, endpoint investigation, collection, Live Response, command/script execution, containment, isolation/quarantine, response verification/rollback, update/upgrade, deep resilience/recovery, anti-tamper, crypto/PKI, transport protocol, API, port, certificate/token format, product code, physical schema, final RBAC/ABAC, Endpoint Screen ID or detailed screen design is introduced. EPT-2..EPT-6 remain NOT STARTED.

## AI / no-AI
AI is optional. Essential functions remain manual/deterministic. AI may summarize source-backed health/change/degradation but cannot invent identity/platform/support/heartbeat/capability, enroll autonomously, modify Fleet/Policy, hide stale/offline state, grant permissions or promote technical facts into foreign canonical objects.

## OPEN / Requirements
Requirements remain **122 = 99 conform / 20 partial / 3 absent / 0 contradictory**. OPEN remains **18**; OPEN-008 remains open and no OPEN is closed by EPT-1.

## Build-time gate ledger
Status before fifth-commit publication: **184 PASS / 6 PENDING-BUILD-OR-REMOTE / 0 FAIL**. One FAIL would make EPT-1 PARTIAL.

### Git / roadmap — 1–20
1 PASS repo correct; 2 PASS visibility; 3 PASS branch; 4 PASS PR #2; 5 PASS base main; 6 PASS PR open; 7 PASS Draft; 8 PASS unmerged; 9 PASS auto-merge disabled; 10 PASS exact baseline; 11 PASS preflight 100/100; 12 PASS README branch; 13 PASS README main; 14 PASS main unchanged at build baseline; 15 PASS Studio PASS preserved; 16 PASS Phase 5 roadmap preserved; 17 PASS roadmap ID preserved; 18 PASS EPT-1 execution-lot only; 19 PASS no Phase 5E1; 20 PASS linear five-commit construction.

### Sources — 21–45
21 PASS governance; 22 PASS Capability Register; 23 PASS Object Register; 24 PASS Dependency Register; 25 PASS Permission Register; 26 PASS Screen Register; 27 PASS OPEN; 28 PASS Requirements; 29 PASS Endpoint preflight; 30 PASS Phase 5 roadmap; 31 PASS Endpoint README; 32 PASS architecture; 33 PASS platform support; 34 PASS enrollment/registration; 35 PASS identity semantics; 36 PASS inventory sources/gap; 37 PASS health/self-check; 38 PASS heartbeat/connectivity; 39 PASS capability inventory/availability sources; 40 PASS Settings Fleet; 41 PASS Settings Policy; 42 PASS Studio boundary; 43 PASS Govern boundary; 44 PASS Shared boundary; 45 PASS migration sources.

### Capability/template — 46–75
46 PASS CAP-EPT namespace; 47 PASS no preexisting concrete IDs; 48 PASS no reserved IDs; 49 PASS 001..014 available; 50 PASS no recycle; 51 PASS final count justified; 52 PASS one file/capability; 53 PASS owner; 54 PASS named users; 55 PASS problem; 56 PASS goals; 57 PASS non-goals; 58 PASS inputs; 59 PASS objects read; 60 PASS objects modified; 61 PASS actions; 62 PASS states; 63 PASS outputs; 64 PASS transitions; 65 PASS source-of-truth; 66 PASS provenance; 67 PASS functional permissions; 68 PASS limits; 69 PASS errors; 70 PASS metrics; 71 PASS no-AI path; 72 PASS >=3 GWT each; 73 PASS Requirements; 74 PASS OPEN; 75 PASS 378 sections / 84 tables.

### Ownership / distinctions — 76–125
76 PASS Endpoint owns local Agent state; 77 PASS Settings owns Fleet; 78 PASS Settings owns Endpoint Policy; 79 PASS Settings owns enrollment admin; 80 PASS Settings owns tenant/env admin; 81 PASS Studio owns Automation Agent; 82 PASS Studio owns Tool; 83 PASS Govern owns Response Run; 84 PASS Govern owns Result; 85 PASS Shared owns generic Jobs/Trace; 86 PASS Endpoint Agent != Automation Agent; 87 PASS Endpoint != Agent automatically; 88 PASS Device != Agent Instance; 89 PASS Host Identity != Agent Identity; 90 PASS registration != enrollment admin; 91 PASS local enrollment != admin record; 92 PASS enrolled != online; 93 PASS registered != healthy; 94 PASS revoked != deleted; 95 PASS platform observed != supported; 96 PASS referenced != delivered; 97 PASS OS identified != compatible; 98 PASS compatible != supported; 99 PASS version observed != approved; 100 PASS inventory != Fleet; 101 PASS snapshot != global source of truth; 102 PASS metadata != Settings config; 103 PASS inventory freshness != health; 104 PASS last seen != heartbeat; 105 PASS heartbeat != health; 106 PASS online != fully capable; 107 PASS degraded != offline; 108 PASS offline != revoked; 109 PASS stale != offline; 110 PASS health != security posture; 111 PASS self-check PASS != integrity proof; 112 PASS advertised != authorized; 113 PASS advertised != available; 114 PASS available != global support; 115 PASS technical capability != Tool; 116 PASS Fleet != local state; 117 PASS Policy != effective local state; 118 PASS assigned != applied; 119 PASS applied != healthy; 120 PASS local effective config != Settings object; 121 PASS technical result != Govern Result; 122 PASS execution != Response Run; 123 PASS technical execution != Tool Call; 124 PASS local rollback != Govern rollback; 125 PASS OPEN-008 != delivery commitment.

### Functional coverage — 126–139
126 PASS Agent identity; 127 PASS enrollment/local registration; 128 PASS tenant/environment; 129 PASS platform/OS/architecture; 130 PASS version/build/compatibility; 131 PASS inventory; 132 PASS freshness/change; 133 PASS health/self-check; 134 PASS heartbeat/connectivity; 135 PASS degraded/offline/stale; 136 PASS capability advertisement; 137 PASS Fleet boundary; 138 PASS Policy effective-state boundary; 139 PASS provenance/handoff.

### AI / security / technical limits — 140–165
140 PASS AI optional; 141 PASS no mandatory chatbot; 142 PASS no identity invention; 143 PASS no platform invention; 144 PASS no support invention; 145 PASS no autonomous enrollment; 146 PASS no Fleet mutation; 147 PASS no Policy mutation; 148 PASS no hidden offline/stale; 149 PASS no fake heartbeat; 150 PASS no fake capability; 151 PASS no raw secret; 152 PASS tenant isolation; 153 PASS restricted metadata; 154 PASS no protocol; 155 PASS no PKI design; 156 PASS no certificate format; 157 PASS no token format; 158 PASS no API; 159 PASS no port; 160 PASS no code; 161 PASS no physical schema; 162 PASS no final RBAC; 163 PASS no Endpoint Screen ID; 164 PASS no detailed screen design; 165 PASS no EPT-2+ capability.

### Registers / non-regression — 166–180
166 PASS Capability Register prepared; 167 PASS Endpoint shard; 168 PASS Dependency Register addendum; 169 PASS Object Map; 170 PASS Action Classification; 171 PASS Automation/AI; 172 PASS cross-links; 173 PASS Requirements addendum; 174 PASS qualitative baseline; 175 PASS OPEN preserved; 176 PASS STATUS build addendum prepared; 177 PASS EPT-1 changelog record; 178 PASS roadmap addendum prepared; 179 PASS Command/Investigate/Govern/Studio content unchanged; 180 PASS Endpoint preflight history preserved.

### Publication — 181–190
181 PASS conformance report exists; 182 PASS build-time state explicit; 183 PASS remote-dependent gates explicitly pending before publication; 184 PENDING-COMMIT-OBJECT five functional commits reachable; 185 PENDING-REMOTE remote verification actually executed; 186 PENDING-REMOTE exact build SHA recorded; 187 PENDING-REMOTE exact final SHA recorded; 188 PENDING-REMOTE post-publication evidence linked; 189 PENDING-REMOTE PR/main/README intact after publication; 190 PASS EPT-2..EPT-6 remain NOT STARTED.

## Post-publication verification contract
After the fifth functional commit is created, gate 184 must be verified by exact ancestry from baseline. After fast-forward publication, gates 185–189 must be verified against the remote branch/PR/main/README and the exact fifth SHA must be recorded in PR #2 post-publication evidence. If all six pending gates close and no divergence is found, final verdict is **PASS AFTER POST-PUBLICATION VERIFICATION — 190/190 PASS, 0 PENDING, 0 FAIL**. If a real divergence exists, EPT-1 remains PARTIAL/PENDING and only a documentary fast-forward correction may address that divergence.

## Final-state consequences on successful verification
Endpoint Capability Specification becomes **PARTIAL** with EPT-1 PASS and EPT-2..EPT-6 NOT STARTED. Endpoint has 14 capabilities; global totals become **399 capabilities / 397 defined / 2 proposed / 399 planned / 10773 sections / 2394 mandatory tables**. Command, Investigate, Govern and Studio remain PASS; Delivery Roadmap Phase 5, Global Capability Specification and repository maturity remain PARTIAL.