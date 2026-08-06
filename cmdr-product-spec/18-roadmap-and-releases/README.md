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

## Périmètre

Document canonique du domaine. Il définit uniquement son sujet et renvoie vers les autres sources de vérité pour les concepts partagés.

## Propriétaire fonctionnel

Product Operations Lead.

## Objets concernés

- Concepts du document
- Références canoniques liées

## Fonctionnalités

- Dependency roadmap.
- Release strategy.
- Phase plans.
- Migration and deprecation.
- Readiness and release evidence.

## Phase plans actifs

- Command: `phase-2-command.md` and related release evidence.
- Investigate foundation: `phase-3-investigate.md`.
- Cloud Analysis: `phase-4b4a-cloud-analysis.md` — `CAP-INV-601..618`, provider-neutral functional scope, post-publication verification pending in the closure commit.
- Mobile Forensics: not started; no capability or phase implementation is created by Cloud Analysis.

## UX et interactions

- Navigation par liens stables.
- Contenu lisible en thème clair et sombre.
- Aucune duplication des définitions externes.

## Permissions

Les modifications suivent le modèle défini dans `../14-security-permissions-and-trust/permission-model.md` lorsque le document décrit une capacité exécutable.

## États

Le statut documentaire suit `00-governance/document-status-model.md`; les états métier restent dans leurs sources canoniques. A roadmap `PASS` is documentary evidence only and never proves implementation.

## Dépendances

- `../00-governance/source-of-truth-policy.md`
- `../00-governance/dependency-register.md`
- `../STATUS.md`
- `../16-quality-and-validation/validation-status.md`

## Critères d’acceptation

- Le document a un propriétaire unique.
- Les liens locaux sont valides.
- Les décisions non tranchées sont attribuées.
- Chaque phase cite son scope, ses dépendances, ses preuves et ses décisions ouvertes.
- Une phase n’est marquée PASS qu’après les vérifications exigées.
- Mobile Forensics n’est pas démarrée implicitement par la clôture Cloud.

## Questions ouvertes

- À compléter — décision source non fournie dans le brief canonique.
