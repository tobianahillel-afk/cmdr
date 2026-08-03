---
id: command-risk-and-business-impact
domain: 06-command
status: draft
owner: Command Product Lead
updated: 2026-08-03
source-of-truth: canonical
---
# Risk & Business Impact

## Objectif

Relier incidents, services métier, risque et impact sans score opaque.

## Périmètre

Module du produit 06-command. Les objets, permissions, composants et transitions partagés sont référencés et non redéfinis.

## Propriétaire fonctionnel

Command Product Lead.

## Objets concernés

- incident
- entity
- result

## Fonctionnalités

- Impact métier.
- Risque d’action et d’inaction.
- Services critiques.
- Explication des facteurs.

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
