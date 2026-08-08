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

## Capability Specification execution plans

This namespace tracks detailed capability specification, ownership, traceability, templates and quality gates.

- `Capability Specification Phase 4A — Command`: **PASS**; historical IDs and reports preserved.
- `Capability Specification Phase 4B — Investigate`: **PASS** after Mobile post-publication verification; historical IDs and reports preserved.
- `Capability Specification Phase 4` global maturity: **PARTIAL** because later object, permission, screen, technique and implementation work remains future.
- Cloud Analysis execution evidence: `phase-4b4a-cloud-analysis.md` — `CAP-INV-601..618`, **PASS AFTER POST-PUBLICATION VERIFICATION**.
- Mobile Forensics execution evidence: `phase-4b4b-mobile-forensics.md` — `CAP-INV-701..719`, **PASS AFTER POST-PUBLICATION VERIFICATION**.

The next product candidate for a separate capability-specification execution is Govern. Govern capability specification is **NOT STARTED** by this reconciliation.

## Delivery Roadmap phases

This namespace preserves the historical product-delivery sequence:

1. [`phase-1-foundation.md`](phase-1-foundation.md) — `Delivery Roadmap Phase 1 — Foundation`;
2. [`phase-2-command.md`](phase-2-command.md) — `Delivery Roadmap Phase 2 — Command`;
3. [`phase-3-investigate.md`](phase-3-investigate.md) — `Delivery Roadmap Phase 3 — Investigate`;
4. [`phase-4-govern.md`](phase-4-govern.md) — `Delivery Roadmap Phase 4 — Govern`, canonical id `roadmap-phase-4-govern`, Govern capability specification **NOT STARTED**;
5. [`phase-5-studio-and-endpoint.md`](phase-5-studio-and-endpoint.md) — `Delivery Roadmap Phase 5 — Studio and Endpoint`;
6. [`phase-6-platform-scale.md`](phase-6-platform-scale.md) — `Delivery Roadmap Phase 6 — Platform Scale`.

The Delivery Roadmap plan files remain historical canonical files and are not renamed to mimic Capability Specification numbering.

## Other active roadmap documents

- dependency ordering: [`dependency-roadmap.md`](dependency-roadmap.md);
- decision sequencing: [`decision-sequencing.md`](decision-sequencing.md);
- capability delivery metadata: [`capability-delivery-map.md`](capability-delivery-map.md);
- release strategy, readiness, evidence, notes and post-release validation;
- backlog, migration, deprecation, pilot, launch and versioning guidance.

The existing active roadmap corpus is retained. This reconciliation adds one convention document and changes no historical phase filename or ID.

## UX et interactions

- Navigation par liens stables.
- Contenu lisible en thème clair et sombre.
- Aucune duplication des définitions externes.

## Permissions

Les modifications suivent le modèle défini dans `../14-security-permissions-and-trust/permission-model.md` lorsque le document décrit une capacité exécutable.

## États

Le statut documentaire suit `00-governance/document-status-model.md`; les états métier restent dans leurs sources canoniques. A roadmap `PASS` or capability-specification `PASS` is documentary evidence only and never proves implementation.

Delivery Roadmap Phase 1–3 plan files currently carry document status `draft`; this index does not invent a delivery-completion verdict for them. Delivery Roadmap Phase 4–6 capability-specification work is not started unless a separate execution explicitly begins it.

## Dépendances

- `../00-governance/source-of-truth-policy.md`
- `../00-governance/dependency-register.md`
- [`dependency-roadmap.md`](dependency-roadmap.md)
- [`phase-numbering-and-namespace-convention.md`](phase-numbering-and-namespace-convention.md)
- `../STATUS.md`
- `../16-quality-and-validation/validation-status.md`

## Critères d’acceptation

- Le document a un propriétaire unique.
- Les liens locaux sont valides.
- Les décisions non tranchées sont attribuées.
- Chaque nouvelle référence de phase ambiguë qualifie son namespace.
- Aucun identifiant historique n’est renommé pour harmoniser les nombres.
- `roadmap-phase-4-govern` reste canonique et aucune `Phase 4C Govern` n’est créée.
- Les dépendances fonctionnelles viennent des sources d’ownership/dependency/transition, jamais de la proximité numérique entre namespaces.

## Questions ouvertes

La réconciliation de namespace ne crée ni ne ferme aucune décision produit. Les décisions existantes restent dans `../00-governance/source-material/unresolved-decisions.md`.
