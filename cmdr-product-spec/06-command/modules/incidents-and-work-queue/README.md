---
id: command-incidents-and-work-queue
domain: 06-command
status: draft
owner: Command Product Lead
updated: 2026-08-03
source-of-truth: canonical
---
# Incidents & Work Queue

## Objectif

Coordonner les incidents, tâches, objets non assignés, risques SLA et charge d’équipe.

## Périmètre

Module du produit 06-command. Les objets, permissions, composants et transitions partagés sont référencés et non redéfinis.

## Propriétaire fonctionnel

Command Product Lead.

## Objets concernés

- incident
- alert
- task
- saved-view
- case

## Fonctionnalités

- Vues Incidents, Tasks, Unassigned, SLA Risk et Team Load.
- Affectation et handover.
- Bulk actions sûres.
- Saved views canoniques.

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
