---
id: roadmap-readme
domain: 18-roadmap-and-releases
status: draft
owner: Product Operations Lead
updated: 2026-08-09
source-of-truth: canonical
---
# Roadmap and Releases

## Objectif

Ordonner dépendances, phases, migrations, releases et preuves de readiness sans créer de décisions produit hors de leurs sources.

## Phase numbering namespaces

CMDR maintains two independent phase namespaces. `Capability Specification Phase 4A — Command`, `Capability Specification Phase 4B — Investigate` and `Delivery Roadmap Phase 4 — Govern` are unrelated numerically. `Phase 4C Govern`, `Phase 4D Govern` and `Capability Specification Phase 4C` do not exist.

## Capability Specification execution status

- `Capability Specification Phase 4A — Command`: **PASS AFTER POST-PUBLICATION VERIFICATION**, current revalidation 60/60; 27 capabilities / 729 sections / 162 tables.
- `Capability Specification Phase 4B — Investigate`: **PASS**; 243 capabilities / 6561 sections / 1458 tables.
- `Govern capability specification`: **PARTIAL**.
  - `GOV-1 — Action Requests, Policy, Authorities and Decisions`: **PASS AFTER POST-PUBLICATION VERIFICATION — 180/180**; 16 capabilities / 432 sections / 96 tables.
  - `GOV-2 — Playbooks, Response Runs, Execution, Verification and Rollback`: **PENDING POST-PUBLICATION VERIFICATION**; 17 capabilities / 459 sections / 102 tables.
  - `GOV-3 — Audit Trail, Response Metrics and Govern Closure`: **NOT STARTED**.
- Global Capability Specification maturity: **PARTIAL**.

`GOV-1`, `GOV-2`, `GOV-3` are execution-lot identifiers only.

## Delivery Roadmap phases

1. [`phase-1-foundation.md`](phase-1-foundation.md) — `Delivery Roadmap Phase 1 — Foundation`;
2. [`phase-2-command.md`](phase-2-command.md) — `Delivery Roadmap Phase 2 — Command`;
3. [`phase-3-investigate.md`](phase-3-investigate.md) — `Delivery Roadmap Phase 3 — Investigate`;
4. [`phase-4-govern.md`](phase-4-govern.md) — **Delivery Roadmap Phase 4 — Govern**, canonical id `roadmap-phase-4-govern`, status **PARTIAL**: GOV-1 PASS, GOV-2 pending remote verification, GOV-3 NOT STARTED;
5. [`phase-5-studio-and-endpoint.md`](phase-5-studio-and-endpoint.md) — future;
6. [`phase-6-platform-scale.md`](phase-6-platform-scale.md) — future.

Historical phase files/IDs are not renamed to mimic Capability Specification numbering.

## Govern execution lots

### GOV-1
`CAP-GOV-001..016`; Action Request → Policy/Authority/Approval → Decision → no-effect Execution Handoff Package. Historical remote verification PASS 180/180.

### GOV-2
`CAP-GOV-017..033`; Decision/Handoff → Playbook → Execution Plan → readiness/authority reconciliation → Response Run/Steps → technical executor handoff → runtime/error reconciliation → verification → rollback/recovery if needed → canonical Result → cross-product provenance/handoff.

Key ownership boundaries remain: Workflow/Tool/Automation Run Studio-owned; technical primitives/raw outcomes Endpoint/provider-owned; secrets/providers/integrations Settings-owned; Incident Command-owned; Case/Evidence/Finding Investigate-owned; generic Jobs/Trace/Reporting/Recovery Shared-owned.

### GOV-3
Audit Trail, Response Metrics and Govern closure remain **NOT STARTED**. GOV-2 only emits provenance/conceptual metric inputs.

## Current totals

- capabilities: **303** — 27 Command / 243 Investigate / 33 Govern;
- defined / proposed / planned: **301 / 2 / 303**;
- Govern: **891 sections / 198 mandatory tables**;
- all Command + Investigate + Govern: **8181 sections / 1818 mandatory tables**;
- Requirements: **122 = 99 conform / 20 partial / 3 absent / 0 contradictory**;
- OPEN decisions: **18**.

## Other active roadmap documents

- dependency ordering: [`dependency-roadmap.md`](dependency-roadmap.md);
- decision sequencing: [`decision-sequencing.md`](decision-sequencing.md);
- capability delivery metadata: [`capability-delivery-map.md`](capability-delivery-map.md);
- release/readiness/evidence/migration/deprecation/pilot/launch/versioning guidance.

## Permissions and implementation

The canonical permission source remains `../14-security-permissions-and-trust/permission-model.md`. Capability specification does not prove runtime implementation or grant authority. GOV-2 selects no command, exploit/bypass, API/protocol, provider/runtime, final object schema or final RBAC/ABAC.

## Acceptance

GOV-2 reaches PASS only after five required functional commits are published, 17/459/102 structural conformance is confirmed, all **190 gates** pass remotely, GOV-1/Command/Investigate non-regression is proven and PR/README/main invariants remain intact. Until then it remains PENDING. GOV-3 must not begin as part of GOV-2 closure.