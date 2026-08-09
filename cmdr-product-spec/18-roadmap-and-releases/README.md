---
id: roadmap-readme
domain: 18-roadmap-and-releases
status: draft
owner: Product Operations Lead
updated: 2026-08-09
source-of-truth: canonical
---
# Roadmap and Releases

## Phase-numbering namespaces
CMDR retains two independent namespaces. `Capability Specification Phase 4A — Command`, `Capability Specification Phase 4B — Investigate` and `Delivery Roadmap Phase 4 — Govern` have no numeric parent/child relationship. `Phase 4C/4D/4E Govern` do not exist.

## Capability specification status
- Command: **PASS**, 27 / 729 / 162.
- Investigate: **PASS**, 243 / 6561 / 1458.
- **Govern: PASS**, 47 / 1269 / 282.
  - GOV-1: PASS, historical 180/180, 16 / 432 / 96.
  - GOV-2: PASS, historical 190/190, 17 / 459 / 102.
  - GOV-3: **PASS AFTER POST-PUBLICATION VERIFICATION — 200/200**, 14 / 378 / 84.
- Global Capability Specification maturity: **PARTIAL**.

## Delivery Roadmap phases
1. `phase-1-foundation.md` — Delivery Roadmap Phase 1 — Foundation.
2. `phase-2-command.md` — Delivery Roadmap Phase 2 — Command.
3. `phase-3-investigate.md` — Delivery Roadmap Phase 3 — Investigate.
4. `phase-4-govern.md` — **Delivery Roadmap Phase 4 — Govern: PASS**, id `roadmap-phase-4-govern`.
5. `phase-5-studio-and-endpoint.md` — next verified historical candidate, id `roadmap-phase-5-studio-and-endpoint`; **not started by GOV-3**.
6. `phase-6-platform-scale.md` — future.

## Govern closure
GOV-1 `CAP-GOV-001..016`, GOV-2 `CAP-GOV-017..033` and GOV-3 `CAP-GOV-034..047` together cover all nine Govern modules from Action Request/Decision through governed execution/Result and audit/metrics/continuous-improvement semantics.

Verified GOV-3 publication:
- baseline `36edacb4eb374e0b56d6c9e9c45931fdb1e0af20`;
- fifth functional SHA `042f70d3cfd13467acc294bfff726edde9e16cb0`;
- 5 ahead / 0 behind, same merge base;
- **200/200 gates PASS**.

## Current totals
- global capabilities: **317** — 27 Command / 243 Investigate / 47 Govern;
- defined/proposed/planned: **315 / 2 / 317**;
- Command + Investigate + Govern: **8559 sections / 1902 tables**;
- Requirements: **122 = 99/20/3/0**;
- OPEN: **18**.

## Closure evidence
- `../16-quality-and-validation/reports/govern-gov3-audit-metrics-closure-capability-conformance.md` — 200/200 PASS.
- `../16-quality-and-validation/reports/govern-capability-specification-closure.md` — Govern PASS.
- `../16-quality-and-validation/reports/delivery-roadmap-phase-4-govern-closure.md` — Roadmap Phase 4 Govern PASS.

## Implementation boundary / stop line
Documentary PASS does not prove software delivery. No audit/metrics engine, API/protocol, warehouse/storage schema, final RBAC/retention policy, detailed screen rewrite or product implementation is claimed. **Delivery Roadmap Phase 5 is identified only and has not been started.**