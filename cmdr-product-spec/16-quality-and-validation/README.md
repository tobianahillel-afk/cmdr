---
id: quality-readme
domain: 16-quality-and-validation
status: draft
owner: Quality Lead
updated: 2026-08-09
source-of-truth: canonical
---
# Quality and Validation

## Objectif

Définir les contrôles qui comparent l’arborescence réelle au manifeste attendu et valident liens, noms, front matter, sources, permissions, propriétaires, métriques et preuves de publication.

## Périmètre

Quality records evidence and verification stages. It does not own product behavior, objects, permissions, roadmap decisions, technology choices or implementation.

## Propriétaire fonctionnel

Quality Lead.

## Active evidence

- Validation status: `validation-status.md`.
- Structural baseline: `architecture-manifest.md` and `expected-path-manifest.md`.
- Current Command revalidation: `reports/phase-4a-command-current-revalidation.md` — **PASS AFTER POST-PUBLICATION VERIFICATION, 60/60**.
- Phase 4B Investigate closure: `reports/phase-4b-investigate-capability-closure.md` — PASS.
- Cloud Analysis: `reports/phase-4b4a-cloud-analysis-capability-conformance.md` — **243/243 PASS**.
- Cloud source audit: `reports/phase-4b4a-cloud-analysis-source-audit.md`.
- Mobile Forensics: `reports/phase-4b4b-mobile-forensics-capability-conformance.md` — **250/250 PASS after publication verification**.
- Phase 4B.4 closure: `reports/phase-4b4-cloud-and-mobile-analysis-closure.md`.
- **Govern GOV-1:** `reports/govern-gov1-action-policy-authority-decision-capability-conformance.md` — **PENDING POST-PUBLICATION VERIFICATION, 169 PASS / 11 PENDING / 0 FAIL before the fifth functional commit is remotely verified**.
- Prior Detection Engineering and Threat Intelligence closure reports remain active and preserved.

## GOV-1 quality contract

Parent roadmap: **Delivery Roadmap Phase 4 — Govern**. Canonical roadmap id: `roadmap-phase-4-govern`. Execution lot: **GOV-1**. GOV-1 is not a roadmap phase. `Phase 4C Govern` does not exist.

GOV-1 requires:
- 16 unique `CAP-GOV-001..016` capability files;
- 432 numbered capability sections;
- 96 mandatory S8/S9/S10/S13/S16/S17 tables;
- zero empty/generic mandatory tables;
- zero duplicate/recycled IDs or owner conflicts;
- 180/180 gates after remote verification;
- no GOV-2/GOV-3 capability or execution implementation;
- non-regression of the five CAP-CMD Requirements ranges, ten DEP-CMD families and detailed Command registry rows.

## Validation rules

- A planned capability is not implementation evidence.
- A report cannot claim post-publication PASS before the remote head, PR, README, `main`, commit chain, registers, metrics and Command non-regression are checked.
- Historical evidence is appended or preserved; closure never condenses earlier traceability.
- Sources are classified by what was actually read/used; unverified sources are not silently claimed.
- A failed or pending gate remains visible.
- Any single GOV-1 FAIL makes the execution lot PARTIAL.
- Remote-dependent GOV-1 gates stay PENDING until the fifth functional commit is published.

## UX et interactions

- Navigation par liens stables.
- Contenu lisible en thème clair et sombre.
- No detailed Govern screen rewrite is performed by GOV-1.

## Permissions

The canonical permission source remains `../14-security-permissions-and-trust/permission-model.md`. GOV-1 validates functional permission needs and authority boundaries only; it does not finalize RBAC/ABAC.

## États

Documentary status follows governance lifecycle. GOV-1 `PENDING POST-PUBLICATION VERIFICATION` is a quality stage, not a product runtime status.

## Dépendances

- `../00-governance/source-of-truth-policy.md`
- `../00-governance/source-material/requirements-traceability-matrix.md`
- `../00-governance/dependency-register.md`
- `../18-roadmap-and-releases/phase-4-govern.md`
- `../STATUS.md`

## Critères d’acceptation

- Every verdict names evidence and verification stage.
- Counts are recalculated from the final registered scope.
- No duplicate ID, owner conflict, active contradiction, broken required link or unsupported implementation claim is hidden.
- Command evidence restored immediately before GOV-1 remains intact.

## Questions ouvertes

OPEN-007, OPEN-013 and OPEN-015 remain open and must not be silently resolved by a quality verdict.