# Status

## Capability Specification Status
- Repository architecture / historical foundations 0–3: PASS as previously recorded.
- **Capability Specification Phase 4A — Command: PASS**, 27 capabilities / 729 sections / 162 tables; 26 defined + 1 proposed.
- **Capability Specification Phase 4B — Investigate: PASS**, 243 / 6561 / 1458.
- **Govern capability specification: PASS**.
  - GOV-1: PASS, historical 180/180; 16 / 432 / 96.
  - GOV-2: PASS, historical 190/190; 17 / 459 / 102.
  - GOV-3: **PASS AFTER POST-PUBLICATION VERIFICATION — 200/200**; 14 / 378 / 84.
- **Studio capability specification: PARTIAL**.
  - STD-1: **PASS AFTER POST-PUBLICATION VERIFICATION — 190/190**, 16 / 432 / 96.
  - STD-2: 17 / 459 / 102, `draft / defined / planned`; **PENDING POST-PUBLICATION VERIFICATION** at build time.
  - STD-3: NOT STARTED.
  - STD-4: NOT STARTED.
- Endpoint capability specification: **NOT STARTED**.
- Global Capability Specification maturity: **PARTIAL**.

## Delivery Roadmap Status
Delivery Roadmap is a separate namespace from Capability Specification.
- Phases 1–3: preserve canonical historical statuses.
- **Delivery Roadmap Phase 4 — Govern: PASS**, id `roadmap-phase-4-govern`.
- **Delivery Roadmap Phase 5 — Studio and Endpoint: PARTIAL**, id `roadmap-phase-5-studio-and-endpoint`; STD-1 PASS, STD-2 functional content prepared/published through its fifth functional commit only after remote verification, STD-3/4 and Endpoint not started.
- Delivery Roadmap Phase 6 — Platform Scale: future.

No `Phase 4C`, `Phase 4D`, `Phase 4E Govern`, `Phase 5A`, `Phase 5B`, `Phase 5C` or `Phase 5D` is introduced. GOV-1/2/3 and STD-1/2 are execution lots, not roadmap phases.

## GOV-3 verified publication — preserved historical evidence
- baseline: `36edacb4eb374e0b56d6c9e9c45931fdb1e0af20`;
- fifth functional SHA: `042f70d3cfd13467acc294bfff726edde9e16cb0`;
- 5 ahead / 0 behind, same merge base;
- GOV-3 gates: **200/200 PASS**.

## STD-1 verified publication — preserved historical evidence
- baseline: `e0c23764df80a3d7109c156d1a2ee0962d19cda6`;
- fifth functional SHA: `d9eb3001482989ab491e1c5446319f410d89d4e8`;
- post-publication correction head before STD-2: `04dcdb43fd7f944a700bf936eebef003546095eb`;
- STD-1 gates: **190/190 PASS**.

## Counts after STD-2 content
- capabilities: **350** — 27 Command / 243 Investigate / 47 Govern / 33 Studio / 0 Endpoint;
- delivery: **348 defined / 2 proposed / 350 planned**;
- Govern: **47 / 1269 / 282**;
- Studio STD-1: **16 / 432 / 96**;
- Studio STD-2: **17 / 459 / 102**;
- Studio cumulative: **33 / 891 / 198**;
- total: **9450 sections / 2100 mandatory tables**;
- Requirements: **122 = 99 conform / 20 partial / 3 absent / 0 contradictory**;
- OPEN decisions: **18**; STD-2 creates/closes 0;
- new Studio Screen IDs / detailed rewrites: **0 / 0**;
- Endpoint capabilities / Screen IDs: **0 / 0**.

## Boundary / maturity
STD-2 is documentary Workflow/Builder/orchestration-definition coverage only. Shared retains generic Trace/Activity/Search/Jobs/Versioning/Recovery; Settings retains providers/integrations/credentials/secrets/tenant/environment administration; Govern retains Playbook/Approval/Decision/Response Run/Result/authority; Endpoint retains technical primitives. No runtime scheduler, Automation Run lifecycle, API/protocol, product code, final language, final JSON Schema, final RBAC/ABAC, detailed screen rewrite or Endpoint capability is introduced.

## Stop line
STD-2 ends at Workflow definition/readiness/pre-publish boundaries. **Do not start STD-3, STD-4 or Endpoint implicitly.**
