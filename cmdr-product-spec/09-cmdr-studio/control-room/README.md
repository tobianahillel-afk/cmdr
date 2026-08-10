---
id: studio-control-room
domain: 09-cmdr-studio
status: draft
owner: CMDR Studio Product Lead
updated: 2026-08-03
source-of-truth: canonical
---
# Control Room

## Objectif

Superviser les déploiements et exécutions d’automatisation.

## Périmètre

Module du produit 09-cmdr-studio. Les objets, permissions, composants et transitions partagés sont référencés et non redéfinis.

## Propriétaire fonctionnel

CMDR Studio Product Lead.

## Objets concernés

- deployment
- automation-agent
- agent-team
- workflow
- result

## Fonctionnalités

- Runtime status.
- Queues.
- Pause/stop.
- Human escalation.

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

## STD-3 runtime capability addendum
`CAP-STD-042..051` define Automation Run creation/lifecycle/steps/attempts, queue/scheduling semantics, control requests, runtime error/retry realization, transient context, Control Room monitoring/intervention, Studio Runtime Outcome and provenance. Control Room ≠ Command Work Queue/Govern Runs & Rollback/Endpoint console. Deployment promotion remains STD-4; generic Job/queue/scheduling engines remain Shared.