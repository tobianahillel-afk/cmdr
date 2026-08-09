# Status

## Capability Specification Status

- **Repository architecture:** structural PASS.
- **Historical specification foundations 0, 1, 2, 3:** PASS as previously recorded; this does not assign Delivery Roadmap completion status to Roadmap Phases 1–3.
- **Capability Specification Phase 4A — Command:** PASS; 27 capabilities, 729 sections, 162 mandatory tables. **Current 2026-08-09 revalidation: PASS AFTER POST-PUBLICATION VERIFICATION, 60/60 controls.**
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
- **Govern capability specification:** **PARTIAL**.
  - **GOV-1 — Action Requests, Policy, Authorities and Decisions:** **PENDING POST-PUBLICATION VERIFICATION**; 16 capabilities, 432 sections, 96 mandatory tables; quality report prepared at 169 PASS / 11 PENDING / 0 FAIL before remote publication verification.
  - **GOV-2 — Playbooks, Response Runs, Execution, Verification and Rollback:** **NOT STARTED**.
  - **GOV-3 — Audit Trail, Response Metrics and Govern Closure:** **NOT STARTED**.
- **Global Capability Specification maturity:** PARTIAL because Govern GOV-2/GOV-3, final objects, permissions, screens, technique and implementation remain future.

## Delivery Roadmap Status

The Delivery Roadmap is a separate namespace from Capability Specification. Numeric proximity between the two namespaces has no implicit semantic or dependency meaning. See `18-roadmap-and-releases/phase-numbering-and-namespace-convention.md`.

- **Delivery Roadmap Phase 1 — Foundation:** canonical plan exists; roadmap document status is `draft`; this status file does not invent a delivery-completion verdict.
- **Delivery Roadmap Phase 2 — Command:** canonical plan exists; roadmap document status is `draft`; this status file does not invent a delivery-completion verdict.
- **Delivery Roadmap Phase 3 — Investigate:** canonical plan exists; roadmap document status is `draft`; this status file does not invent a delivery-completion verdict.
- **Delivery Roadmap Phase 4 — Govern:** **PARTIAL**, canonical id `roadmap-phase-4-govern`, canonical title `Phase 4 Govern`; GOV-1 awaits post-publication verification and GOV-2/GOV-3 are NOT STARTED. **No `Phase 4C Govern` exists.**
- **Delivery Roadmap Phase 5 — Studio and Endpoint:** future / not started for capability specification; historical roadmap file retained.
- **Delivery Roadmap Phase 6 — Platform Scale:** future / not started for capability specification; historical roadmap file retained.

`GOV-1` is an execution-lot identifier, not a Roadmap Phase.

## Current totals and invariants

- **Requirements:** 122 total — 99 conform, 20 partial, 3 absent, 0 contradictory.
- **Capabilities:** **286 registered** — 27 Command, 243 Investigate, 16 Govern; **284 defined, 2 proposed; all 286 planned**.
- **Command:** 27 capabilities — 26 defined, 1 proposed, all 27 planned; 729 sections and 162 mandatory tables.
- **Investigate:** 243 capabilities; 6561 sections and 1458 mandatory tables.
- **Govern GOV-1:** 16 capabilities, all defined/planned; 432 sections and 96 mandatory tables.
- **Command + Investigate + Govern GOV-1:** 286 capabilities, 7722 sections and 1716 mandatory tables.
- **CAP-INV-3xx / CAP-INV-4xx / CAP-INV-5xx / CAP-INV-6xx / CAP-INV-7xx:** 97 / 35 / 37 / 18 / 19.
- **Open decisions:** 18; GOV-1 opens 0 and closes 0. OPEN-007, OPEN-013 and OPEN-015 remain directly relevant and open.
- **Screens:** GOV-1 read 9 Govern screen specs, modified 0 detailed screen specs and created 0 Screen IDs.
- **Objects/permissions:** no complete object schema, JSON Schema, final state machine or final RBAC/ABAC matrix is created by GOV-1.
- **Implementation:** no API, protocol, Policy/authority/runtime engine, command, product code, target mutation, Response Run, Result or rollback is created by GOV-1.
- **Command non-regression:** 27 CAP-CMD remain registered; five CAP-CMD Requirements ranges and ten DEP-CMD families remain required; Command capability files modified by GOV-1 = 0.
- **PR #2:** must remain open, Draft, unmerged and not ready for global review through publication verification.
- **Root README:** must remain exactly `# cmdr` on branch and `main`.

## GOV-1 canonical evidence

- `08-govern/README.md`;
- `08-govern/capability-map.md`;
- `08-govern/functional-dependency-map.md`;
- `08-govern/object-consumption-map.md`;
- `08-govern/action-classification.md`;
- `08-govern/automation-and-ai-model.md`;
- `08-govern/cross-product-links.md`;
- `08-govern/permissions.md`;
- `08-govern/screen-capability-map.md`;
- `00-governance/registers/capability-register-govern-gov1.md`;
- `00-governance/registers/capability-register.md`;
- `00-governance/dependency-register.md`;
- `00-governance/source-material/requirements-traceability-matrix.md`;
- `00-governance/source-material/qualitative-baseline.md`;
- `18-roadmap-and-releases/phase-4-govern.md`;
- `16-quality-and-validation/reports/govern-gov1-action-policy-authority-decision-capability-conformance.md`.

## Primary numbering evidence

- `18-roadmap-and-releases/phase-numbering-and-namespace-convention.md`;
- `18-roadmap-and-releases/README.md`;
- `18-roadmap-and-releases/phase-4-govern.md`;
- `16-quality-and-validation/reports/phase-numbering-and-roadmap-namespace-reconciliation.md`.

## Next action

Publish the fifth GOV-1 traceability/quality commit and perform remote verification. Do not mark GOV-1 PASS before that check, and do not start GOV-2 or GOV-3.