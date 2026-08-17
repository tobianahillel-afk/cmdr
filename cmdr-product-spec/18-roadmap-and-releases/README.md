---
id: roadmap-readme
domain: 18-roadmap-and-releases
status: draft
owner: Product Operations Lead
updated: 2026-08-17
source-of-truth: canonical
---
# Roadmap and Releases

## Phase-numbering namespaces
CMDR retains two independent namespaces. `Capability Specification Phase 4A — Command`, `Capability Specification Phase 4B — Investigate` and `Delivery Roadmap Phase 4 — Govern` have no numeric parent/child relationship. `Phase 4C/4D/4E Govern` do not exist.

## Capability specification status
- Command: **PASS**, 27 capabilities / 729 sections / 162 mandatory tables; **27 defined / 0 proposed / 27 planned**.
- Investigate: **PASS**, 243 / 6561 / 1458.
- Govern: **PASS**, 47 / 1269 / 282.
  - GOV-1: PASS, historical 180/180, 16 / 432 / 96.
  - GOV-2: PASS, historical 190/190, 17 / 459 / 102.
  - GOV-3: **PASS AFTER POST-PUBLICATION VERIFICATION — 200/200**, 14 / 378 / 84.
- Studio: **PASS**, 68 / 1836 / 408.
- Endpoint: **PASS**, 99 / 2673 / 594.
  - EPT-1: PASS 190/190.
  - EPT-2: PASS 200/200.
  - EPT-3: PASS 210/210.
  - EPT-4: PASS 220/220.
  - EPT-5: PASS 230/230.
  - EPT-6: **PASS AFTER POST-PUBLICATION VERIFICATION — 240/240**.
- Global Capability Specification maturity: **PARTIAL**.

## Delivery Roadmap phases
1. `phase-1-foundation.md` — Delivery Roadmap Phase 1 — Foundation; preserve canonical historical status.
2. `phase-2-command.md` — Delivery Roadmap Phase 2 — Command; preserve canonical historical status.
3. `phase-3-investigate.md` — Delivery Roadmap Phase 3 — Investigate; preserve canonical historical status.
4. `phase-4-govern.md` — **Delivery Roadmap Phase 4 — Govern: PASS**, id `roadmap-phase-4-govern`.
5. `phase-5-studio-and-endpoint.md` — **Delivery Roadmap Phase 5 — Studio and Endpoint: PASS — capability specification complete**, id `roadmap-phase-5-studio-and-endpoint`.
6. `phase-6-platform-scale.md` — **Delivery Roadmap Phase 6 — Platform Scale: PARTIAL**, id `roadmap-phase-6-platform-scale`.

STD-1/2/3/4 and EPT-1/2/3/4/5/6 are execution lots only. They are not roadmap subphases and no Phase 5A/5B/5C/5D/E/E6 or Phase 6A/6B is introduced.

## Govern closure — preserved historical evidence
GOV-1 `CAP-GOV-001..016`, GOV-2 `CAP-GOV-017..033` and GOV-3 `CAP-GOV-034..047` together cover all nine Govern modules from Action Request/Decision through governed execution/Result and audit/metrics/continuous-improvement semantics.

Verified GOV-3 publication:
- baseline `36edacb4eb374e0b56d6c9e9c45931fdb1e0af20`;
- fifth functional SHA `042f70d3cfd13467acc294bfff726edde9e16cb0`;
- 5 ahead / 0 behind, same merge base;
- **200/200 gates PASS**.

## Current totals
- global capabilities: **498**;
- defined / proposed / planned: **497 / 1 / 498**;
- only proposed capability: `CAP-INV-106`;
- total structure: **13446 sections / 2988 mandatory tables**;
- Requirements: **122 = 99 conform / 20 partial / 3 absent / 0 contradictory**;
- OPEN: **17**.

## Current Phase 5 / Phase 6 boundary
Delivery Roadmap Phase 5 is **PASS — capability specification complete** after Studio and Endpoint documentary closure. EPT-6's historical BUILD remains **233/240 PASS / 7 PENDING-REMOTE / 0 FAIL**; its current post-publication state is **240/240 PASS**.

Delivery Roadmap Phase 6 remains **PARTIAL**. This synthesis does not close Phase 6 and does not start Phase 7.

## Closure evidence
- `../16-quality-and-validation/reports/govern-gov3-audit-metrics-closure-capability-conformance.md` — historical Govern 200/200 PASS.
- `../16-quality-and-validation/reports/govern-capability-specification-closure.md` — Govern PASS.
- `../16-quality-and-validation/reports/delivery-roadmap-phase-4-govern-closure.md` — Roadmap Phase 4 Govern PASS.
- `phase-5-studio-and-endpoint-ept6-build-addendum.md` — Phase 5 EPT-6 build/closure source.
- `../16-quality-and-validation/reports/endpoint-ept6-updates-resilience-security-provenance-post-publication-verification.md` — EPT-6 240/240 post-publication evidence.
- `../16-quality-and-validation/reports/endpoint-capability-specification-closure.md` — Endpoint capability-specification closure.
- `../16-quality-and-validation/reports/delivery-roadmap-phase-5-studio-and-endpoint-closure.md` — Roadmap Phase 5 capability-specification closure.
- `phase-6-platform-scale.md` — current Phase 6 owner source.

## Implementation boundary / stop line
Documentary PASS proves capability-specification completeness only. It does not prove implementation complete, production readiness, deployment or operational/runtime validation. Repository/global maturity remains **PARTIAL**.
