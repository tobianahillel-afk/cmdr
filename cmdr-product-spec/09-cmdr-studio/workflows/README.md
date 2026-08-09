---
id: studio-workflows
domain: 09-cmdr-studio
status: draft
owner: CMDR Studio Product Lead
updated: 2026-08-10
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

## Fonctionnalités historiques
- Branches.
- Retries.
- Compensation.
- Versioning.

## STD-2 functional specification
`CAP-STD-017..033` détaille Workflow Definition/Version, Builder Session, I/O/variables, nodes/graph, Tool/Skill steps, deterministic conditions/branches, mappings, subworkflows, ordering/parallelism, error paths, retries/idempotency, partial success/compensation, Human Gate boundary, readiness, pre-publish lifecycle and provenance.

STD-2 ne définit pas Automation Agent runtime, Automation Run lifecycle, Control Room runtime, scheduler, publishing/deployment, Endpoint capabilities, API/protocol or orchestration language.

## UX et interactions
- Conserver le contexte de liste, vue et objet.
- Utiliser l’Inspector canonique.
- Afficher les six états obligatoires.
- Préserver navigation clavier et liens profonds.
- STD-2 crée 0 Screen ID et 0 detailed screen rewrite.

## Permissions
Voir `../../14-security-permissions-and-trust/permission-model.md`, le registre et `../workflow-functional-permissions.md`. `perm.studio.*` et `perm.cmdr-studio.*` restent non normalisés.

## États
Les états métier viennent des fichiers d’objets canoniques et des capabilities; la page ajoute Loading, Empty, Partial, Error, Offline, Permission denied/Stale sans changer la sémantique métier.

## Dépendances
- 03-design-system/
- 04-experience-architecture/
- 05-domain-model/
- 17-implementation-contracts/
- STD-1 Tool/Tool Call/Skill contracts
- Govern / Settings / Shared boundaries

## Critères d’acceptation
- aucune duplication des objets/permissions externes;
- 17 STD-2 capabilities × 27 sections et six tables;
- Workflow != Playbook; Human Gate != Approval/Decision; validation != execution; compensation != Govern rollback.

## Questions ouvertes
OPEN-007, OPEN-013 and OPEN-015 remain open. No new OPEN decision is created.
