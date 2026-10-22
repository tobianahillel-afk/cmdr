---
id: endpoint-capability-specification-foundations-preflight
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-11
source-of-truth: quality-report
requirements:
  - REQ-PROD-006
  - REQ-PROD-012
  - REQ-PROD-017
  - REQ-PROD-018
  - REQ-PROD-019
  - REQ-OBJ-008
  - REQ-SEC-001
  - REQ-SEC-002
open_decisions:
  - OPEN-008
  - OPEN-013
  - OPEN-015
---
# Endpoint Capability Specification Foundations — Preflight Rerun / Canonical Closure

## 1. Git baseline
This run is **ENDPOINT PREFLIGHT RERUN / CLOSURE** only. It does not start EPT-1.

Previous attempt:
- verdict: **BLOCKED — 96/100**;
- previous baseline: `c21ea86cde1bea425d7d9233d9973b867f5ef9e8`;
- blocking reason: the exact STD-4 post-publication closure commit did not yet exist at the accepted Endpoint preflight baseline, Studio could not yet be treated as final PASS, and this canonical Endpoint foundations report had not been published.

Rerun baseline:
- repository: `tobianahillel-afk/cmdr`;
- branch: `docs/cmdr-product-spec-foundation`;
- PR: `#2`, open, Draft, unmerged;
- base: `main`;
- exact starting HEAD: `9030186e7aa12990a3d8fb6f30aa107539e2a117`;
- exact starting title: `docs: record Studio STD-4 post-publication verification`;
- parent documentary correction: `c21ea86cde1bea425d7d9233d9973b867f5ef9e8`;
- STD-4 fifth functional/build SHA: `216bff304fa389e4814cb610097571a4a83c1c54`;
- `main`: `bc1ec59e5e79ccbaba291e2984b0eb7d5e54129c`;
- repository visibility: public;
- auto-merge: disabled;
- branch and `main` README: exact `# cmdr`, blob `901c74cda52e28b5ff7fc425ddf28ef89f3ad875`;
- current baseline commit statuses: none;
- current baseline workflow runs: none;
- CI/status disposition: **N/A**.

The starting HEAD is exact and descends linearly from `c21ea86...`; no rebase, reset, force-push or history rewrite is used by this rerun.

## 2. Diff since the blocked Endpoint preflight
Comparison `c21ea86cde1bea425d7d9233d9973b867f5ef9e8` → `9030186e7aa12990a3d8fb6f30aa107539e2a117` is **1 ahead / 0 behind**, with `c21ea86...` as merge base.

Exactly ten documentary files changed:
1. `00-governance/registers/capability-register.md`;
2. `09-cmdr-studio/README.md`;
3. `16-quality-and-validation/README.md`;
4. `16-quality-and-validation/reports/delivery-roadmap-phase-5-studio-domain-closure.md`;
5. `16-quality-and-validation/reports/studio-capability-specification-closure.md`;
6. `16-quality-and-validation/reports/studio-std4-assurance-lifecycle-post-publication-verification.md`;
7. `16-quality-and-validation/validation-status-studio-std4.md`;
8. `18-roadmap-and-releases/phase-5-studio-and-endpoint.md`;
9. `CHANGELOG.md`;
10. `STATUS.md`.

Endpoint-path files changed since the blocked audit: **0**. In particular, no file under `cmdr-product-spec/11-endpoint-agent/` changed. No Settings Endpoint Fleet/Policy source, OPEN-008 source, Object Register, Permission Register, Screen Register or naming convention changed in that diff. The historical read-only Endpoint audit is therefore reusable; only temporal/Git/Studio-closure facts require refresh.

## 3. Studio closure confirmation
STD-4 is **PASS AFTER POST-PUBLICATION VERIFICATION — 220/220 PASS, 0 PENDING, 0 FAIL**. Studio Capability Specification is **PASS** across `CAP-STD-001..068`.

Historical Studio lots remain:
- STD-1: 16 / 432 / 96 — PASS 190/190;
- STD-2: 17 / 459 / 102 — PASS 200/200;
- STD-3: 18 / 486 / 108 — PASS 210/210;
- STD-4: 17 / 459 / 102 — PASS 220/220.

Studio PASS is documentary only and does not prove implementation or Endpoint delivery.

## 4. Roadmap identity
Canonical parent remains **Delivery Roadmap Phase 5 — Studio and Endpoint**, id `roadmap-phase-5-studio-and-endpoint`.

Studio and Endpoint remain distinct capability domains. STD-1/2/3/4 and EPT-1/2/3/4/5/6 are execution-lot labels only. No Phase 5A/5B/5C/5D/5E1/5E2 or competing Capability Specification Phase 5 is created.

Delivery Roadmap Phase 5 remains **PARTIAL** because Endpoint Capability Specification remains NOT STARTED.

## 5. Global metrics
Recalculated from current canonical registers/evidence:
- Command: **27 — PASS**;
- Investigate: **243 — PASS**;
- Govern: **47 — PASS**;
- Studio: **68 — PASS**;
- Endpoint: **0 — NOT STARTED**;
- global capabilities: **385**;
- defined / proposed / planned: **383 / 2 / 385**;
- sections: **10,395**;
- mandatory tables: **2,310**;
- Requirements: **122 = 99 conform / 20 partial / 3 absent / 0 contradictory**;
- OPEN decisions: **18**;
- Delivery Roadmap Phase 5: **PARTIAL**;
- Global Capability Specification: **PARTIAL**;
- repository maturity: **PARTIAL**.

## 6. Endpoint corpus inventory
The prior read-only audit counted **74 documents** under `11-endpoint-agent/`. Because the rerun diff changes zero Endpoint-path file and the current manifest preserves the same path set, the 74-document inventory remains current.

The corpus covers:
- root/foundations: README, architecture, product definition, platform support, trust boundaries, command contract and resource guardrails;
- Telemetry;
- Detection;
- Investigation;
- Collection;
- Live Response;
- Containment;
- Resilience;
- Security.

Endpoint capability files: **0**. Endpoint capability-register shard: **0**. Endpoint implementation delivered by this run: **0**.

## 7. README information-architecture anomaly
`11-endpoint-agent/README.md` still references `information-architecture.md`.

`11-endpoint-agent/information-architecture.md` is absent on the rerun baseline.

Classification: **D — genuine missing document, non-blocking for capability work**.

This rerun does not redirect the reference to `architecture.md` and does not create the missing Information Architecture. A future EPT-1 may produce it if explicitly justified.

## 8. OPEN-008
`OPEN-008 — platform/source availability and support` remains **open**.

No OPEN decision is closed or created by this rerun. OPEN-008 is non-blocking for provider/platform-neutral capability specification, but may become blocking before actual delivery/support claims.

Unresolved support dimensions include:
- initial Windows support/version;
- initial Linux support/distributions;
- initial macOS support/version;
- workstation versus server scope;
- cloud workload/container scope;
- mobile scope;
- source availability and support evidence.

## 9. Platform support matrix
| Platform/scope | Current documentary state | Delivery claim |
|---|---|---|
| Windows | referenced candidate, to decide by release | **not officially supported/delivered** |
| Linux | referenced candidate, distro/version unresolved | **not officially supported/delivered** |
| macOS | referenced candidate, version unresolved | **not officially supported/delivered** |
| Workstation | scope not finalized | none |
| Server | scope not finalized | none |
| Cloud workload/container | scope not finalized | none |
| Mobile | scope not finalized; Endpoint capability not assumed | none |

`platform-support.md` states that Windows/Linux/macOS are decided by release, capabilities are declared rather than assumed, compatibility is versioned, and unsupported operations must fail explicitly. EPT-1 therefore remains platform-neutral/provider-neutral.

## 10. Capability namespace
Repository-wide namespace evidence is unchanged since the historical audit.

| Family | Current disposition |
|---|---|
| `CAP-EPT-*` | canonical Endpoint capability namespace |
| `CAP-END-*` | non-canonical; no concrete ID |
| `CAP-ENDPOINT-*` | non-canonical; no concrete ID |
| `CAP-AGENT-*` | non-canonical; no concrete ID |
| `CAP-EDR-*` | no concrete match |
| `CAP-SENSOR-*` | no concrete match |

Concrete `CAP-EPT-*` IDs: **0**. Reserved `CAP-EPT-*` IDs/ranges: **0**. Legacy Endpoint capability namespace: **none**.

This preflight creates and reserves no `CAP-EPT-001` or any other Capability ID.

## 11. Object ownership
### Endpoint Agent owns
- local Endpoint Agent identity and local registration/effective state;
- tenant/environment binding as local technical state;
- platform/OS/architecture observations;
- agent version/build observations;
- local inventory facts and freshness;
- technical health, heartbeat, connectivity, last-seen and degraded state;
- local capability availability/advertisement;
- future technical execution facts through Endpoint-owned primitives such as `agent-command`;
- `local-audit-event`.

### Platform Settings owns
- Endpoint Agent Fleet;
- administrative enrollment;
- Endpoint Policy;
- policy assignment;
- administrative configuration;
- upgrade waves;
- providers/integrations;
- credentials/secrets;
- tenant/environment administration.

### Studio owns
Tool, Tool Call, Skill, Workflow, Automation Agent, Automation Run and Studio Deployment semantics.

### Govern owns
Approval, Decision, Playbook, Response Run, canonical Result, production-response authority, verification and response rollback.

### Shared owns
generic Jobs, Trace, Activity, Search, Reporting, Export, Notifications, Versioning and generic Recovery mechanisms.

No owner conflict is introduced.

## 12. Agent Policy boundary
Canonical `Endpoint Policy` remains **Platform Settings-owned**.

Endpoint may expose only the effective local projection/state needed for execution and diagnostics, for example:
- effective policy/config version;
- applied;
- mismatch;
- stale;
- refused/degraded;
- effective local state.

This does not create a second canonical Endpoint Policy object and does not transfer policy administration to Endpoint.

## 13. Fleet boundary
`Endpoint Agent Fleet` remains **Platform Settings-owned**.

Settings owns enrollment administration, fleet membership/configuration, aggregated health/capability projection and upgrade waves. Endpoint owns the individual agent's technical state.

Invariant: **Fleet administrative configuration ≠ individual Endpoint operational state**.

## 14. Settings boundary
Settings continues to own tenants/environments, principals/roles, providers/integrations, credentials/secrets, sources/parsers, Endpoint Agent Fleet, Endpoint Policies, fleet enrollment administration, administrative version/capability inventory, update-wave administration and administrative health surfaces.

EPT-1 may consume Settings references/projections but cannot redefine these sources.

## 15. Studio boundary
Endpoint technical capability ≠ Studio Tool. Endpoint Agent ≠ Studio Automation Agent. Endpoint execution ≠ Tool Call. Endpoint technical deployment/update primitives ≠ Studio Deployment semantics.

Studio may orchestrate or reference Endpoint primitives only through future authorized contracts; no transfer of Endpoint ownership occurs.

## 16. Govern boundary
Endpoint execution ≠ Govern Response Run. Endpoint technical result/output ≠ Govern Result. Endpoint local rollback ≠ Govern response rollback.

Govern retains production-response authority, Decision/Approval, Response Run, verification, Result and governed rollback/recovery. Endpoint remains a technical executor within its own domain.

## 17. Shared boundary
Endpoint consumes Shared generic mechanisms without creating competing engines. Generic Job ≠ Endpoint execution. Shared Recovery ≠ Endpoint local technical rollback automatically.

## 18. Command / Investigate boundary
Command retains Incident/Alert/Signal/Task coordination. Investigate retains Case/Artifact/Evidence/Finding, endpoint investigation workflow and collection request semantics. Endpoint performs technical observation/execution; it does not acquire Command or Investigate business-object ownership.

## 19. Permission state
Existing Endpoint permission family remains:
- `perm.endpoint-agent.endpoint-agent.read`;
- `perm.endpoint-agent.endpoint-agent.manage`;
- `perm.endpoint-agent.agent-command.read`;
- `perm.endpoint-agent.agent-command.manage`;
- `perm.endpoint-agent.local-audit-event.read`;
- `perm.endpoint-agent.local-audit-event.manage`.

Distinct Settings permissions remain:
- `perm.platform-settings.endpoint-agent-fleet.read/manage`;
- `perm.platform-settings.endpoint-policy.read/manage`.

This rerun performs no permission normalization, creates no new namespace and does not finalize RBAC/ABAC.

## 20. Screen state
Endpoint Agent product-specific Screen IDs: **0**.

Adjacent Settings-owned screens remain:
- `SET-EAF-001` — Endpoint Agent Fleet;
- `SET-EPL-001` — Endpoint Policies.

No Endpoint Screen ID is created or reserved by this preflight.

## 21. Migration state
The read-only Endpoint corpus remains active foundation documentation. No destructive migration, renaming, source deletion or ownership transfer is performed.

Known stale/missing link: Endpoint README → missing `information-architecture.md`; classified as non-blocking and preserved for future explicit correction.

No Endpoint functional source is modified by this preflight.

## 22. Endpoint execution-lot decomposition
The six-lot decomposition remains justified and unchanged:
1. **EPT-1 — Enrollment, Inventory, Health and Platform Foundations**;
2. **EPT-2 — Telemetry, Observation and Technical Capability Declaration**;
3. **EPT-3 — Local Detection and Endpoint Investigation**;
4. **EPT-4 — Collection and Live Response Technical Execution**;
5. **EPT-5 — Containment, Verification and Governed Response Primitives**;
6. **EPT-6 — Updates, Resilience, Security and Endpoint Provenance**.

These labels are execution lots, not Roadmap Phases.

The six-lot split prevents the first lot from mixing foundational identity/health with high-risk live response, containment, security or update semantics.

## 23. Future EPT-1 scope
If separately authorized after this preflight, EPT-1 covers only foundations:
- Endpoint Agent identity and instance registration;
- enrollment handoff and local enrollment state;
- unenrollment/revocation local effect;
- tenant/environment binding;
- platform, OS and architecture identification;
- support-status representation without declaring support;
- agent version/build and compatibility context;
- inventory snapshot and endpoint metadata;
- inventory freshness and change tracking;
- technical health and self-checks;
- heartbeat, connectivity and last-seen;
- online/offline/degraded/stale semantics;
- technical capability advertisement summary;
- Fleet projection boundary;
- Endpoint Policy effective-state boundary;
- provenance and cross-product handoff.

## 24. EPT-1 out of scope
EPT-1 does not include detailed telemetry, Detection, detailed Endpoint Investigation, Collection/acquisition, Live Response, command/script execution, containment/isolation/quarantine, local rollback, detailed Result reconciliation, update/upgrade, deep resilience/buffering/recovery, anti-tamper, crypto, PKI, transport protocol, API, Fleet administration, Policy administration or Screen design.

## 25. Fourteen future capability titles — no IDs allocated
Candidate titles remain:
1. Endpoint Agent Identity and Instance Registration
2. Enrollment Handoff and Endpoint Enrollment State
3. Tenant and Environment Binding
4. Platform, Operating System and Architecture Identification
5. Agent Version, Build and Compatibility Context
6. Endpoint Inventory Snapshot and Metadata
7. Inventory Freshness and Change Tracking
8. Agent Health and Self-Check Assessment
9. Heartbeat, Connectivity and Last-Seen State
10. Degraded, Offline and Stale State Semantics
11. Technical Capability Advertisement and Availability
12. Fleet Administration and Endpoint State Projection Boundary
13. Endpoint Policy Assignment and Effective Local State Boundary
14. Endpoint Foundations Provenance and Cross-Product Handoff

**Suggested future allocation if the namespace remains free immediately before EPT-1:**
- 001 → title 1;
- 002 → title 2;
- 003 → title 3;
- 004 → title 4;
- 005 → title 5;
- 006 → title 6;
- 007 → title 7;
- 008 → title 8;
- 009 → title 9;
- 010 → title 10;
- 011 → title 11;
- 012 → title 12;
- 013 → title 13;
- 014 → title 14.

This is **not a reservation**. The Capability Register contains no `CAP-EPT-001..014` rows after this preflight.

## 26. Blockers and non-blockers
### Historical blockers now resolved
- Gate 10 — exact baseline: resolved by real starting HEAD `9030186e...`.
- Gate 11 — baseline title: resolved by the actual remote commit `docs: record Studio STD-4 post-publication verification`.
- Gate 15 — Studio PASS preserved: resolved by STD-4 220/220 and Studio PASS.
- Gate 99 — report publication: resolved only when the commit containing this canonical report is fast-forward published and the report is fetched back remotely.

### Non-blockers retained
- `OPEN-008` remains open for support/delivery choices;
- missing Endpoint `information-architecture.md` remains a genuine non-blocking documentation gap;
- platform-specific implementation choices remain unresolved by design;
- no final RBAC/ABAC, API, protocol, physical schema or implementation exists.

## 27. Prior 96/100 history
The previous Endpoint preflight result **BLOCKED — 96/100** is intentionally retained. It is not rewritten as if it had passed historically.

Previous failing gates were exactly **10, 11, 15 and 99**. The rerun revalidates all 100 controls rather than changing the historical record.

## 28. Rerun quality gates — 1–100
### Git / roadmap — 1–20
1. **PASS** — repository correct
2. **PASS** — visibility recorded as public
3. **PASS** — branch correct
4. **PASS** — PR #2 correct
5. **PASS** — base `main`
6. **PASS** — PR open
7. **PASS** — PR Draft
8. **PASS** — PR unmerged
9. **PASS** — auto-merge disabled
10. **PASS** — exact rerun baseline `9030186e7aa12990a3d8fb6f30aa107539e2a117`
11. **PASS** — baseline title `docs: record Studio STD-4 post-publication verification`
12. **PASS** — README branch unchanged, exact `# cmdr`
13. **PASS** — README main unchanged, exact `# cmdr`
14. **PASS** — `main` unchanged at `bc1ec59e5e79ccbaba291e2984b0eb7d5e54129c`
15. **PASS** — Studio PASS preserved; STD-4 220/220
16. **PASS** — Roadmap Phase 5 preserved
17. **PASS** — roadmap id `roadmap-phase-5-studio-and-endpoint` preserved
18. **PASS** — Phase 5 remains PARTIAL
19. **PASS** — Endpoint remains NOT STARTED
20. **PASS** — no history rewrite; blocked-baseline → rerun-baseline is fast-forward 1 ahead / 0 behind

### Source audit — 21–40
21. **PASS** — governance sources valid
22. **PASS** — Capability Register valid
23. **PASS** — Object Register valid
24. **PASS** — Dependency Register valid; STD-4 additions transfer no Endpoint ownership
25. **PASS** — Permission Register valid
26. **PASS** — Screen Register valid
27. **PASS** — OPEN Register/source valid; 18 open
28. **PASS** — Requirements valid; 122 = 99/20/3/0
29. **PASS** — Phase 5 roadmap valid
30. **PASS** — previous Phase 5 preflight evidence valid as historical source
31. **PASS** — Studio closure valid
32. **PASS** — Endpoint corpus inventory confirmed at 74 documents
33. **PASS** — active Endpoint docs covered by prior read-only audit and unchanged path set
34. **PASS** — deprecated/historical Endpoint sources classified; no competing active source introduced
35. **PASS** — Settings boundary valid
36. **PASS** — Studio boundary valid
37. **PASS** — Govern boundary valid
38. **PASS** — Shared boundary valid
39. **PASS** — Command/Investigate boundary valid
40. **PASS** — historical Endpoint sources covered; rerun required no full 74-doc re-audit because relevant diff is zero

### Endpoint foundations — 41–60
41. **PASS** — Endpoint Agent concept resolved
42. **PASS** — Endpoint/Device/Host terminology treated as bounded context without inventing a Device owner
43. **PASS** — Agent identity ownership resolved to Endpoint Agent
44. **PASS** — shared Endpoint administration boundary resolved to Settings where applicable
45. **PASS** — enrollment semantics inventoried
46. **PASS** — installation/registration state inventoried without implementation claim
47. **PASS** — tenant/environment association inventoried
48. **PASS** — platform sources inventoried
49. **PASS** — OPEN-008 fully analyzed and remains open
50. **PASS** — platform support matrix produced
51. **PASS** — agent version/build semantics inventoried
52. **PASS** — inventory semantics inventoried
53. **PASS** — metadata semantics inventoried
54. **PASS** — health semantics inventoried
55. **PASS** — heartbeat semantics inventoried
56. **PASS** — connectivity semantics inventoried
57. **PASS** — freshness/staleness semantics inventoried
58. **PASS** — offline/degraded states inventoried
59. **PASS** — provenance requirements inventoried
60. **PASS** — first-lot boundaries established

### Ownership / namespace — 61–80
61. **PASS** — Settings ownership resolved
62. **PASS** — Agent/Endpoint Policy ownership assessed as Platform Settings
63. **PASS** — Fleet ownership assessed as Platform Settings
64. **PASS** — Studio boundary explicit
65. **PASS** — Govern boundary explicit
66. **PASS** — Shared boundary explicit
67. **PASS** — technical result != Govern Result
68. **PASS** — Endpoint Agent != Studio Automation Agent
69. **PASS** — Endpoint technical capability != Studio Tool
70. **PASS** — Endpoint execution != Tool Call
71. **PASS** — Endpoint execution != Response Run
72. **PASS** — namespace families globally revalidated against current PR evidence
73. **PASS** — canonical namespace determined as `CAP-EPT-*`
74. **PASS** — concrete Endpoint Capability IDs inventoried: 0
75. **PASS** — reserved Endpoint Capability IDs/ranges inventoried: 0
76. **PASS** — no ID allocated
77. **PASS** — no ID reserved
78. **PASS** — Endpoint and Settings permission families inventoried
79. **PASS** — no RBAC normalization/finalization
80. **PASS** — no owner conflict introduced

### Documentation / screens / migration — 81–90
81. **PASS** — README anomaly reproduced
82. **PASS** — `information-architecture.md` classified genuine missing/non-blocking
83. **PASS** — screen inventory completed
84. **PASS** — no new Endpoint Screen ID
85. **PASS** — migration sources classified
86. **PASS** — stale link identified: Endpoint README → missing information architecture
87. **PASS** — no destructive migration
88. **PASS** — no Endpoint capability file created
89. **PASS** — no object schema created
90. **PASS** — no implementation created

### Conclusion — 91–100
91. **PASS** — global baseline metrics recalculated
92. **PASS** — six-lot Endpoint decomposition justified
93. **PASS** — EPT-1 scope justified
94. **PASS** — EPT-1 out-of-scope justified
95. **PASS** — fourteen future capability titles proposed without IDs
96. **PASS** — blockers classified
97. **PASS** — non-blocking OPENs/gaps identified
98. **PASS** — exact next functional execution identified as EPT-1, but not started
99. **PASS ON PUBLICATION** — this canonical report is included in the single preflight-closure commit; final PASS is valid only after the report is fetched from the remote branch
100. **PASS** — Endpoint Capability Specification remains NOT STARTED with 0 capabilities

Pre-publication logical state: **99 controls immediately verifiable + gate 99 bound to publication of this same canonical record**. The post-publication verification section below is mandatory before the final external verdict is stated.

## 29. Exact next functional run
Only after successful post-publication verification of this 100/100 preflight, the next functional run is:

**ENDPOINT LOT EPT-1 — Enrollment, Inventory, Health and Platform Foundations**.

Before allocating any ID, EPT-1 must immediately revalidate the then-current remote HEAD, `CAP-EPT-*` namespace, zero concrete IDs, zero reserved IDs, OPEN-008 and Endpoint/Settings ownership. Only then may `CAP-EPT-001..014` be considered for real allocation.

## 30. Final verdict and stop line
Rerun target verdict: **PASS — 100/100** **only after** the single report commit is fast-forward published and the following remote checks succeed:
- exact new remote HEAD is the report commit;
- commit reachable from `9030186e...` by one fast-forward commit;
- PR #2 remains open/Draft/unmerged on `main`;
- branch/main README remain unchanged;
- `main` remains unchanged;
- Studio remains PASS;
- Endpoint capabilities remain 0;
- concrete `CAP-EPT-*` IDs remain 0;
- reserved `CAP-EPT-*` IDs remain 0;
- this report exists remotely and records the 100-gate closure;
- Delivery Roadmap Phase 5 remains PARTIAL.

If any post-publication invariant diverges, the rerun is **PARTIAL/BLOCKED**, not PASS.

This report creates no capability, no reserved ID, no Screen ID, no object schema, no permission, no API/protocol, no implementation and no support-platform delivery claim. **Do not start EPT-1 in this commit.**