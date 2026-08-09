# Status

## Capability Specification Status

- **Repository architecture:** structural PASS.
- **Historical specification foundations 0, 1, 2, 3:** PASS as previously recorded; this does not assign Delivery Roadmap completion status to Roadmap Phases 1–3.
- **Capability Specification Phase 4A — Command:** PASS; 27 capabilities, 729 sections, 162 mandatory tables. Current 2026-08-09 revalidation: **PASS AFTER POST-PUBLICATION VERIFICATION, 60/60 controls**.
- **Capability Specification Phase 4B — Investigate:** **PASS**; 243 capabilities, 6561 sections, 1458 mandatory tables. Historical subphase evidence remains unchanged.
- **Govern capability specification:** **PARTIAL**.
  - **GOV-1 — Action Requests, Policy, Authorities and Decisions:** **PASS AFTER POST-PUBLICATION VERIFICATION — 180/180**; 16 capabilities, 432 sections, 96 mandatory tables. Exact verification-record baseline for GOV-2: `b8dd93e03443adb9101c7592094a48e358b460e2`.
  - **GOV-2 — Playbooks, Response Runs, Execution, Verification and Rollback:** **PENDING POST-PUBLICATION VERIFICATION**; 17 capabilities, 459 sections, 102 mandatory tables. Functional contracts are published through the fourth required commit; traceability/quality publication and remote 190-gate verification remain pending.
  - **GOV-3 — Audit Trail, Response Metrics and Govern Closure:** **NOT STARTED**.
- **Global Capability Specification maturity:** **PARTIAL** because GOV-3, final objects, permissions, screens, technique and implementation remain future.

## Delivery Roadmap Status

The Delivery Roadmap is a separate namespace from Capability Specification. Numeric proximity between the namespaces has no implicit semantic or parent-child meaning.

- **Delivery Roadmap Phase 1 — Foundation:** preserve its canonical historical draft/status; no new verdict invented.
- **Delivery Roadmap Phase 2 — Command:** preserve its canonical historical draft/status; no new verdict invented.
- **Delivery Roadmap Phase 3 — Investigate:** preserve its canonical historical draft/status; no new verdict invented.
- **Delivery Roadmap Phase 4 — Govern:** **PARTIAL**, canonical id `roadmap-phase-4-govern`, title `Phase 4 Govern`; GOV-1 PASS, GOV-2 pending remote verification, GOV-3 NOT STARTED.
- **Delivery Roadmap Phase 5 — Studio and Endpoint:** future.
- **Delivery Roadmap Phase 6 — Platform Scale:** future.

`GOV-1`, `GOV-2` and `GOV-3` are execution-lot identifiers, not Roadmap Phases. **No `Phase 4C Govern` or `Phase 4D Govern` exists.**

## Current totals and invariants

- **Requirements:** 122 total — **99 conform, 20 partial, 3 absent, 0 contradictory**.
- **Capabilities:** **303 registered** — 27 Command, 243 Investigate, 33 Govern; **301 defined, 2 proposed; all 303 planned**.
- **Command:** 27 capabilities — 26 defined, 1 proposed; 729 sections / 162 mandatory tables.
- **Investigate:** 243 capabilities; 6561 sections / 1458 mandatory tables.
- **Govern GOV-1:** 16 defined/planned; 432 sections / 96 tables; historical 180/180 PASS preserved.
- **Govern GOV-2:** 17 defined/planned; 459 sections / 102 tables.
- **Govern cumulative:** 33 capabilities; **891 sections / 198 mandatory tables**.
- **Command + Investigate + Govern:** **303 capabilities / 8181 sections / 1818 mandatory tables**.
- **Open decisions:** **18**; GOV-2 opens 0 and closes 0. OPEN-007, OPEN-008, OPEN-013, OPEN-015 and OPEN-019 remain relevant/open.
- **Screens:** nine existing Govern screens retained; GOV-2 primary surfaces are Playbooks and Runs & Rollback; detailed screen rewrites = 0; new Screen IDs = 0.
- **Objects/permissions:** no complete object schema, JSON Schema, final state machine or final RBAC/ABAC matrix is created by GOV-2.
- **Implementation:** no executable command, exploit/bypass, API/protocol, provider/runtime choice, product code or raw secret value is added.
- **GOV-3:** Audit Trail and Response Metrics capabilities created = 0.

## Ownership / execution invariants

- Govern owns Response Playbook semantics, Execution Plan, canonical Response Run/Step governance, verification, rollback/recovery governance and canonical Result.
- Studio retains Workflow, Workflow Version, Tool, Tool Call, Human Gate and Automation Run.
- Endpoint/provider owners retain technical execution primitives, health and raw technical outcomes.
- Platform Settings retains providers, integrations, secrets, credentials and runtime/tenant/environment administration.
- Command retains Incident/general Work Queue; Investigate retains Case/Evidence/Finding/analysis; Shared retains generic Jobs/Trace/Activity/Versioning/Reporting/Recovery.
- Response Run != Automation Run/Job/Tool Call. Technical output != canonical Result. Result never rewrites Decision/Evidence/Finding.

## Non-regression

- GOV-1 `CAP-GOV-001..016`: intact; 16 / 432 / 96 and historical 180/180 PASS preserved.
- Command: 27 CAP-CMD remain registered/unique; 26 defined + 1 proposed; five restored CAP-CMD Requirements ranges and `DEP-CMD-001..010` preserved; Command capability files modified by GOV-2 = 0.
- Investigate: 243 CAP-INV remain; Phase 4B stays PASS; CAP-INV capability files modified by GOV-2 = 0.
- Root README must remain exactly `# cmdr` on canonical branch and `main`; `main` must remain unchanged.
- PR #2 must remain open, Draft, unmerged and not ready for global review.

## GOV-2 evidence

- `08-govern/execution-boundaries.md`;
- `08-govern/capability-map.md`;
- `08-govern/functional-dependency-map.md`;
- `08-govern/object-consumption-map.md`;
- `08-govern/action-classification.md`;
- `08-govern/automation-and-ai-model.md`;
- `08-govern/cross-product-links.md`;
- `08-govern/permissions.md`;
- `08-govern/screen-capability-map.md`;
- `00-governance/registers/capability-register-govern-gov1.md`;
- `00-governance/registers/capability-register-govern-gov2.md`;
- `00-governance/registers/capability-register.md`;
- `00-governance/registers/dependency-register-govern-gov2.md`;
- `00-governance/source-material/requirements-traceability-matrix.md` plus GOV-2 additive evidence;
- `00-governance/source-material/qualitative-baseline.md` plus GOV-2 additive baseline;
- `00-governance/source-material/unresolved-decisions.md` plus GOV-2 audit;
- `18-roadmap-and-releases/phase-4-govern.md`;
- `16-quality-and-validation/reports/govern-gov2-playbooks-response-runs-verification-rollback-capability-conformance.md` once published.

## Next action

Publish the fifth GOV-2 traceability/quality commit, verify the remote state and only then promote GOV-2 to `PASS AFTER POST-PUBLICATION VERIFICATION`. Do **not** begin GOV-3.