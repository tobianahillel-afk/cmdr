# Status

## Capability Specification Status

- **Repository architecture:** structural PASS.
- **Historical specification foundations 0, 1, 2, 3:** PASS as previously recorded; this does not assign Delivery Roadmap completion status to Roadmap Phases 1–3.
- **Capability Specification Phase 4A — Command:** PASS; 27 capabilities, 729 sections, 162 mandatory tables.
- **Capability Specification Phase 4B.1:** PASS; 22 capabilities, 594 sections, 132 tables.
- **Capability Specification Phase 4B.2:** PASS; 112 capabilities, 3024 sections, 672 tables.
- **Capability Specification Phase 4B.3A — Detection Engineering:** PASS; 35 capabilities, 945 sections, 210 tables.
- **Capability Specification Phase 4B.3B.1 — Threat Intelligence Foundations:** PASS; 18 capabilities, 486 sections, 108 tables.
- **Capability Specification Phase 4B.3B.2 — Intelligence Analysis, Dissemination and Operationalization:** PASS after remote publication verification; 19 capabilities, 513 sections, 114 tables and 180 gates.
- **Capability Specification Phase 4B.3B — Threat Intelligence:** PASS; 37 capabilities, 999 sections and 222 tables.
- **Capability Specification Phase 4B.3:** PASS; Detection Engineering and Threat Intelligence both PASS.
- **Capability Specification Phase 4B.4A — Cloud Analysis:** PASS AFTER POST-PUBLICATION VERIFICATION; 18 capabilities, 486 sections, 108 mandatory tables and 243/243 gates.
- **Capability Specification Phase 4B.4B — Mobile Forensics:** PASS AFTER POST-PUBLICATION VERIFICATION; 19 capabilities, 513 sections, 114 mandatory tables and 250/250 gates.
- **Capability Specification Phase 4B.4:** PASS; Cloud and Mobile both verified; combined 37 capabilities, 999 sections and 222 tables.
- **Capability Specification Phase 4B — Investigate:** PASS.
- **Capability Specification Phase 4 global maturity:** PARTIAL because final objects, permissions, screens, technique and implementation remain future.

## Delivery Roadmap Status

The Delivery Roadmap is a separate namespace from Capability Specification. Numeric proximity between the two namespaces has no implicit semantic or dependency meaning. See `18-roadmap-and-releases/phase-numbering-and-namespace-convention.md`.

- **Delivery Roadmap Phase 1 — Foundation:** canonical plan exists; roadmap document status is `draft`; this status file does not invent a delivery-completion verdict.
- **Delivery Roadmap Phase 2 — Command:** canonical plan exists; roadmap document status is `draft`; this status file does not invent a delivery-completion verdict.
- **Delivery Roadmap Phase 3 — Investigate:** canonical plan exists; roadmap document status is `draft`; this status file does not invent a delivery-completion verdict.
- **Delivery Roadmap Phase 4 — Govern:** **NOT STARTED** for capability specification; canonical id `roadmap-phase-4-govern`; canonical title `Phase 4 Govern`; no `Phase 4C Govern` exists.
- **Delivery Roadmap Phase 5 — Studio and Endpoint:** **NOT STARTED** for capability specification; historical roadmap file retained.
- **Delivery Roadmap Phase 6 — Platform Scale:** **NOT STARTED** for capability specification; historical roadmap file retained.

The next product candidate for a separate capability-specification execution is Govern. No Govern capability is created by the numbering reconciliation.

## Current totals and invariants

- **Requirements:** 122 total — 99 conform, 20 partial, 3 absent, 0 contradictory.
- **Capabilities:** 270 registered — 27 Command and 243 Investigate; 268 defined, 2 proposed; all 270 planned.
- **CAP-INV-3xx / CAP-INV-4xx / CAP-INV-5xx / CAP-INV-6xx / CAP-INV-7xx:** 97 / 35 / 37 / 18 / 19.
- **Investigate:** 6561 sections and 1458 mandatory tables.
- **Command + Investigate:** 7290 sections and 1620 mandatory tables.
- **Open decisions:** 18; none opened or closed by the numbering reconciliation. OPEN-011 remains open for Mobile delivery strategy; OPEN-012 remains open for Cloud delivery strategy; OPEN-018 and OPEN-019 remain open.
- **Screens:** existing surfaces mapped conceptually; 0 screen spec modified, 0 detailed rewrite, 0 new Screen ID by this reconciliation.
- **Implementation:** no code, API, protocol, provider/platform selection, connector, physical schema, scanner, device/target mutation, secret use, deployed rule or response is created by this reconciliation.
- **Govern capabilities:** 0; Govern capability specification remains NOT STARTED.
- **PR #2:** must remain open, Draft, unmerged and not ready for global review.
- **Root README:** must remain exactly `# cmdr` on branch and `main`.

## Primary numbering evidence

- `18-roadmap-and-releases/phase-numbering-and-namespace-convention.md`;
- `18-roadmap-and-releases/README.md`;
- `18-roadmap-and-releases/phase-4-govern.md`;
- `06-command/capability-map.md`;
- `07-investigate/capability-map.md`;
- `16-quality-and-validation/reports/phase-numbering-and-roadmap-namespace-reconciliation.md`.
