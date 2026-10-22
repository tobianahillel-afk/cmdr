---
id: studio-agent-teams
domain: 09-cmdr-studio
status: draft
owner: CMDR Studio Product Lead
updated: 2026-08-03
source-of-truth: canonical
---
# Agent Teams

## Objectif

Composer des équipes d’Automation Agents avec rôles et coordination.

## Périmètre

Module du produit 09-cmdr-studio. Les objets, permissions, composants et transitions partagés sont référencés et non redéfinis.

## Propriétaire fonctionnel

CMDR Studio Product Lead.

## Objets concernés

- agent-team
- automation-agent
- workflow
- human-gate

## Fonctionnalités

- Roles.
- Communication.
- Coordinator.
- Escalation.

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

## STD-3 capability addendum
`CAP-STD-037..039` define Agent Team member/version refs, functional roles, bounded shared context, coordination/conflict behavior, no-effect planning and human oversight. Agent Team ≠ human Team/Workflow and delegation ≠ permission delegation. The historical source above remains preserved; no multi-agent framework is selected.