---
id: quality-readme
domain: 16-quality-and-validation
status: draft
owner: Quality Lead
updated: 2026-08-07
source-of-truth: canonical
---
# Quality and Validation

## Objectif

Définir les contrôles qui comparent l’arborescence réelle au manifeste attendu et valident liens, noms, front matter, sources, permissions, propriétaires, métriques et preuves de publication.

## Périmètre

Document canonique du domaine. Il définit uniquement son sujet et renvoie vers les autres sources de vérité pour les concepts partagés.

## Propriétaire fonctionnel

Quality Lead.

## Objets concernés

- Concepts du document
- Références canoniques liées

## Fonctionnalités

- Expected path manifest.
- Structural audit.
- Screen contract validation.
- Source-of-truth validation.
- Security, accessibility and endpoint validation.
- Release gates.
- Publication-stage verification and evidence preservation.

## Active evidence

- Validation status: `validation-status.md`.
- Structural baseline: `architecture-manifest.md` and `expected-path-manifest.md`.
- Cloud Analysis conformance: `reports/phase-4b4a-cloud-analysis-capability-conformance.md`.
- Cloud Analysis source audit: `reports/phase-4b4a-cloud-analysis-source-audit.md`.
- Prior Threat Intelligence closure reports remain active and preserved.

## Validation rules

- A planned capability is not implementation evidence.
- A report cannot claim post-publication PASS before the remote head, PR, README, `main`, registers and metrics are checked.
- Historical evidence is appended or preserved; closure never condenses earlier traceability.
- Sources are classified as fully read, referenced, identified or absent.
- A failed or pending gate remains visible.

## UX et interactions

- Navigation par liens stables.
- Contenu lisible en thème clair et sombre.
- Aucune duplication des définitions externes.

## Permissions

Les modifications suivent le modèle défini dans `../14-security-permissions-and-trust/permission-model.md` lorsque le document décrit une capacité exécutable. Quality records evidence but does not own product behavior, objects, permissions, roadmap decisions or implementation.

## États

Le statut documentaire suit `00-governance/document-status-model.md`; les états métier restent dans leurs sources canoniques.

## Dépendances

- `../00-governance/source-of-truth-policy.md`
- `../00-governance/source-material/requirements-traceability-matrix.md`
- `../STATUS.md`

## Critères d’acceptation

- Le document a un propriétaire unique.
- Les liens locaux sont valides.
- Les décisions non tranchées sont attribuées.
- Every verdict names the evidence and its verification stage.
- Counts are recalculated from the final registered scope.
- No duplicate ID, owner conflict, active contradiction, broken required link or unsupported implementation claim is hidden.

## Questions ouvertes

- À compléter — décision source non fournie dans le brief canonique.
- Cloud Analysis remains pending until remote verification is complete.
