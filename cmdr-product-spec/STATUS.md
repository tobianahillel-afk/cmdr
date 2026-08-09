# Status

## Capability Specification Status
- Historical repository foundations 0–3: PASS as recorded.
- **Command: PASS** — 27 capabilities / 729 sections / 162 tables; 26 defined + 1 proposed.
- **Investigate: PASS** — 243 / 6561 / 1458.
- **Govern: PASS** — 47 / 1269 / 282; GOV-1/2/3 historical gates 180/190/200 preserved.
- **Studio: PARTIAL**.
  - STD-1 `CAP-STD-001..016`: 16 capabilities / 432 sections / 96 mandatory tables; `draft / defined / planned`.
  - STD-2: NOT STARTED.
  - STD-3: NOT STARTED.
  - STD-4: NOT STARTED.
- Endpoint capability specification: **NOT STARTED**.
- Global Capability Specification maturity: **PARTIAL**.

## Delivery Roadmap Status
- Delivery Roadmap Phase 4 — Govern: PASS.
- **Delivery Roadmap Phase 5 — Studio and Endpoint: PARTIAL**; canonical id `roadmap-phase-5-studio-and-endpoint`.
- Delivery Roadmap Phase 6 — Platform Scale: future.

No Phase 5A/5B is created. `STD-1` is an execution lot only.

## Current counts after STD-1 content
- capabilities: **333** — 27 Command / 243 Investigate / 47 Govern / 16 Studio / 0 Endpoint;
- defined / proposed / planned: **331 / 2 / 333**;
- total sections / mandatory tables: **8991 / 1998**;
- Requirements: **122 = 99 conform / 20 partial / 3 absent / 0 contradictory**;
- OPEN decisions: **18**; STD-1 creates/closes 0;
- Studio new Screen IDs / detailed rewrites: **0 / 0**;
- Endpoint capabilities/screens created: **0 / 0**.

## STD-1 boundaries
Library != Shared Search; Tool != Tool Call/Skill/Workflow/Agent/Endpoint primitive/Govern Playbook; Tool Call != Automation Run/Response Run/Job/Result; Skill composition != Workflow orchestration; Secret Reference != secret value; provider/runtime availability != permission; Human Gate != Approval; Workflow != Playbook; Automation Run != Response Run.

Both historical `perm.studio.*` and `perm.cmdr-studio.*` families remain documented; STD-1 makes no bulk rename or final atomic permission decision.

## Implementation boundary
All new Studio capabilities are documentary `defined / planned`. No API, protocol, product code, runtime engine, final Tool/Tool Call schema, final RBAC/ABAC, provider/runtime choice or Endpoint capability is introduced.

## Stop line
STD-1 ends at Studio foundations. **Do not start STD-2, STD-3, STD-4 or Endpoint implicitly.**
