---
id: studio-workflows
domain: 09-cmdr-studio
status: draft
owner: CMDR Studio Product Lead
updated: 2026-08-03
source-of-truth: canonical
---
# Workflows

## Objectif

Orchestrer étapes déterministes, Skills, Agents et Human Gates.

## Périmètre

Module du produit 09-cmdr-studio. Les objets, permissions, composants et transitions partagés sont référencés et non redéfinis.

## Propriétaire fonctionnel

CMDR Studio Product Lead.

## Objets concernés

- workflow
- skill
- automation-agent
- human-gate

## Fonctionnalités

- Branches.
- Retries.
- Compensation.
- Versioning.

## UX et interactions

- Conserver le contexte de liste, vue et objet.
- Utiliser l’Inspector canonique.
- Afficher les six états obligatoires.
- Préserver navigation clavier et liens profonds.

## Permissions

Voir `../../14-security-permissions-and-trust/permission-model.md` et le registre des permissions.

## États

Les états métier viennent des fichiers d’objets canoniques; la page ajoute uniquement Loading, Empty, Partial, Error, Offline et Permission denied.

## Dépendances

- 03-design-system/
- 04-experience-architecture/
- 05-domain-model/
- 17-implementation-contracts/

## Critères d’acceptation

- Aucune définition d’objet ou de permission locale.
- Tous les écrans du module ont un front matter et 27 sections.
- Les transitions sont auditées et idempotentes.

## Questions ouvertes

- À compléter — contenu source non fourni dans le brief canonique.

## STD-2 canonical addendum

`CAP-STD-017..033` now details Workflow Definition/Version, Builder Session, I/O/variables, nodes/graph, Tool/Skill steps, deterministic conditions/branches, mappings, subworkflows, ordering/parallelism, error paths, retries/idempotency, partial success/compensation, Human Gate/Govern boundary, readiness, pre-publish lifecycle and provenance.

The historical generic “Questions ouvertes” statement above is preserved as source history; current STD-2 explicitly keeps OPEN-007, OPEN-013 and OPEN-015 open and creates no new OPEN.

STD-2 creates 0 Screen IDs and 0 detailed screen rewrites. It defines no Automation Agent runtime, Automation Run lifecycle, Control Room runtime, scheduler, publishing/deployment engine, Endpoint capability, API/protocol or orchestration language. `perm.studio.*` and `perm.cmdr-studio.*` remain non-normalized.
