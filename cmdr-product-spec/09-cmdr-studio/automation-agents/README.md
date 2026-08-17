---
id: studio-automation-agents
domain: 09-cmdr-studio
status: draft
owner: CMDR Studio Product Lead
updated: 2026-08-03
source-of-truth: canonical
---
# Automation Agents

## Objectif

Définir agents, objectifs, outils, limites, mémoire et politiques.

## Périmètre

Module du produit 09-cmdr-studio. Les objets, permissions, composants et transitions partagés sont référencés et non redéfinis.

## Propriétaire fonctionnel

CMDR Studio Product Lead.

## Objets concernés

- automation-agent
- skill
- workflow
- human-gate

## Fonctionnalités

- Role and goal.
- Tool access.
- Memory/context.
- Guardrails.

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
`CAP-STD-034..036`, `CAP-STD-038` and `CAP-STD-039` now define Automation Agent identity, bounded objectives/context, Tool/Skill/resource access, bounded plan proposals and human oversight. Agent objective/role/access/proposal never grants permission or authority. The historical source above remains preserved; no framework, model or provider is selected.