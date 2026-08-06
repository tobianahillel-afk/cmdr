---
id: roadmap-readme
domain: 18-roadmap-and-releases
status: draft
owner: Product Operations Lead
updated: 2026-08-07
source-of-truth: canonical
---
# Roadmap and Releases

## Objectif
Ordonner les dépendances, phases, migrations, releases et preuves de readiness sans créer de décisions produit hors de leurs sources.

## Phase plans actifs
- Command: `phase-2-command.md` and related release evidence.
- Investigate foundation: `phase-3-investigate.md`.
- Cloud Analysis: `phase-4b4a-cloud-analysis.md` — `CAP-INV-601..618`, provider-neutral functional scope, post-publication verification pending in the closure commit.
- Mobile Forensics: not started; no capability or phase implementation is created by Cloud Analysis.

## Functionalités
- Dependency roadmap.
- Release strategy.
- Phase plans.
- Migration and deprecation.
- Readiness and release evidence.

## Ownership and boundaries
Product Operations owns roadmap sequencing and evidence. Product, object, permission, screen and technical sources remain with their canonical owners. A roadmap `PASS` is documentary evidence only and never proves implementation.

## Permissions and states
Executable permissions remain in `../14-security-permissions-and-trust/permission-model.md`. Documentary status follows governance lifecycle; business states remain in owner sources.

## Dépendances
- `../00-governance/source-of-truth-policy.md`
- `../00-governance/dependency-register.md`
- `../STATUS.md`
- `../16-quality-and-validation/validation-status.md`

## Critères d’acceptation
- Le document a un propriétaire unique.
- Chaque phase cite son scope, ses dépendances, ses preuves et ses décisions ouvertes.
- Une phase n’est marquée PASS qu’après les vérifications exigées.
- Mobile Forensics n’est pas démarrée implicitement par la clôture Cloud.
