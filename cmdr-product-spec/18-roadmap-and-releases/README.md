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

Ordonner les dépendances, phases, migrations, releases et preuves de readiness sans créer de décisions produit hors de leurs sources.

## Phase numbering namespaces

CMDR maintains two independent phase namespaces. The canonical convention is [`phase-numbering-and-namespace-convention.md`](phase-numbering-and-namespace-convention.md).

A bare `Phase 4` must not be used in a new ambiguous context. `Capability Specification Phase 4B — Investigate` and `Delivery Roadmap Phase 4 — Govern` are different phases in different namespaces; the shared number does not create a parent-child relationship and does not imply a `Phase 4C Govern`.

## Capability Specification execution status

This namespace tracks detailed capability specification, ownership, traceability, templates and quality gates.

- `Capability Specification Phase 4A — Command`: **PASS AFTER POST-PUBLICATION VERIFICATION**; 27 capabilities, 729 sections, 162 tables; current revalidation 60/60.
- `Capability Specification Phase 4B — Investigate`: **PASS**; historical IDs and reports preserved.
- `Govern capability specification`: **PARTIAL**.
  - `GOV-1 — Action Requests, Policy, Authorities and Decisions`: **PENDING POST-PUBLICATION VERIFICATION**; 16 capabilities, 432 sections, 96 mandatory tables.
  - `GOV-2 — Playbooks, Response Runs, Execution, Verification and Rollback`: **NOT STARTED**.
  - `GOV-3 — Audit Trail, Response Metrics and Govern Closure`: **NOT STARTED**.
- Global Capability Specification maturity: **PARTIAL** because later Govern lots, objects, permissions, screens, technique and implementation remain future.
- Cloud Analysis execution evidence: `phase-4b4a-cloud-analysis.md` — `CAP-INV-601..618`, **PASS AFTER POST-PUBLICATION VERIFICATION**.
- Mobile Forensics execution evidence: `phase-4b4b-mobile-forensics.md` — `CAP-INV-701..719`, **PASS AFTER POST-PUBLICATION VERIFICATION**.

`GOV-1`, `GOV-2` and `GOV-3` are execution-lot identifiers, not roadmap phases.

## Delivery Roadmap phases

This namespace preserves the historical product-delivery sequence:

1. [`phase-1-foundation.md`](phase-1-foundation.md) — `Delivery Roadmap Phase 1 — Foundation`;
2. [`phase-2-command.md`](phase-2-command.md) — `Delivery Roadmap Phase 2 — Command`;
3. [`phase-3-investigate.md`](phase-3-investigate.md) — `Delivery Roadmap Phase 3 — Investigate`;
4. [`phase-4-govern.md`](phase-4-govern.md) — **`Delivery Roadmap Phase 4 — Govern`**, canonical id `roadmap-phase-4-govern`, current status **PARTIAL** because GOV-1 awaits remote verification and GOV-2/GOV-3 are NOT STARTED;
5. [`phase-5-studio-and-endpoint.md`](phase-5-studio-and-endpoint.md) — `Delivery Roadmap Phase 5 — Studio and Endpoint`, future;
6. [`phase-6-platform-scale.md`](phase-6-platform-scale.md) — `Delivery Roadmap Phase 6 — Platform Scale`, future.

The Delivery Roadmap plan files remain historical canonical files and are not renamed to mimic Capability Specification numbering. `Phase 4C Govern` and `Capability Specification Phase 4C` **do not exist**.

## Govern GOV-1

Parent roadmap: **Delivery Roadmap Phase 4 — Govern**. Canonical roadmap id: `roadmap-phase-4-govern`.

GOV-1 defines `CAP-GOV-001..016` for Action Request intake/lifecycle, scope/risk/completeness, Policy evaluation/conflicts/exceptions, authority/eligibility/SoD, Approval/delegation/emergency, Decision preparation/recording/conditions/expiration and a no-effect Execution Handoff Package. It creates no Response Run, Result, target execution or rollback.

The GOV-1 fifth functional commit and remote 180-gate verification are required before GOV-1 can be promoted to PASS. GOV-2 must not start as part of that closure.

## Other active roadmap documents

- dependency ordering: [`dependency-roadmap.md`](dependency-roadmap.md);
- decision sequencing: [`decision-sequencing.md`](decision-sequencing.md);
- capability delivery metadata: [`capability-delivery-map.md`](capability-delivery-map.md);
- release strategy, readiness, evidence, notes and post-release validation;
- backlog, migration, deprecation, pilot, launch and versioning guidance.

The existing active roadmap corpus and all historical phase filenames/IDs remain retained.

## UX et interactions

- Navigation par liens stables.
- Contenu lisible en thème clair et sombre.
- Aucune duplication des définitions externes.

## Permissions

The canonical permission source remains `../14-security-permissions-and-trust/permission-model.md`. A capability-specification PASS never grants runtime authority or proves implementation.

## États

Roadmap, capability-specification, document and runtime states remain distinct. Delivery Roadmap Phase 1–3 plan files currently carry document status `draft`; this index does not invent delivery-completion verdicts for them.

## Dépendances

- `../00-governance/source-of-truth-policy.md`
- `../00-governance/dependency-register.md`
- [`dependency-roadmap.md`](dependency-roadmap.md)
- [`phase-numbering-and-namespace-convention.md`](phase-numbering-and-namespace-convention.md)
- `../STATUS.md`
- `../16-quality-and-validation/validation-status.md`

## Critères d’acceptation

- Chaque nouvelle référence de phase ambiguë qualifie son namespace.
- Aucun identifiant historique n’est renommé pour harmoniser les nombres.
- `roadmap-phase-4-govern` reste canonique et aucune `Phase 4C Govern` n’est créée.
- GOV-1 remains an execution lot, not a roadmap phase.
- Les dépendances fonctionnelles viennent des sources d’ownership/dependency/transition, jamais de la proximité numérique entre namespaces.
- GOV-2/GOV-3 remain NOT STARTED until separate executions.

## Questions ouvertes

GOV-1 closes no product decision. `OPEN-007`, `OPEN-013` and `OPEN-015` remain explicitly open; all other current OPEN decisions retain their dispositions.