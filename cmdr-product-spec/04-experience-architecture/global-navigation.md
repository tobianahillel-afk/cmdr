---
id: global-navigation
domain: 04-experience-architecture
status: draft
owner: UX Architecture Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-UX-001
  - REQ-UX-006
  - REQ-PROD-008
  - REQ-PROD-013
  - REQ-PROD-018
---

# Navigation globale et locale

## Global Header

Le Global Header est la première zone de l'application et reste visible dans les workspaces principaux. De gauche à droite :

1. wordmark textuel CMDR ;
2. navigation primaire `Command`, `Investigate`, `Govern` ;
3. accès secondaire `Studio` et `Settings` selon permission ;
4. recherche globale ;
5. Command Palette ;
6. Jobs ;
7. notifications ;
8. aide ;
9. profil et préférences.

Il ne contient ni filtres métier, ni actions de l'objet sélectionné, ni chatbot. Endpoint Agent n'y apparaît pas comme produit principal.

## Navigation primaire

- positions stables, intitulés textuels et indicateur actif non dépendant de la couleur ;
- `Tab` atteint chaque destination, `Enter` l'ouvre ;
- le changement de produit transmet le contexte compatible et demande confirmation si un travail non enregistré risque d'être abandonné ;
- une destination interdite est masquée ou annoncée comme indisponible sans révéler de données.

## Navigation locale

Placée à gauche, elle liste modules, workspaces, favoris locaux et vues épinglées. Largeur Draft : `240 px` développée, `56 px` réduite. Elle peut être masquée dans un canvas maximal et restaurée par raccourci. Les actions de l'objet, les réglages globaux et l'Inspector n'y sont pas dupliqués.

## Context Bar

Sous le Global Header, au-dessus du workspace. Il est absent lorsque le contexte est uniquement global. Il montre au plus les niveaux pertinents : tenant, environnement, service, Incident, Case, Endpoint, Artifact/Evidence/Finding, Decision ou Response Run.

Priorité en manque d'espace : tenant → environnement → objet actif → parent direct → autres ancêtres dans overflow. L'identifiant critique reste disponible dans le nom accessible et le menu.

## Raccourcis Draft

- `Ctrl/Cmd+K` : Command Palette ;
- `/` : recherche locale lorsque le focus n'est pas dans un éditeur ;
- `Alt+1..5` : produits autorisés ;
- `[` : navigation locale ;
- `]` : Inspector ;
- `Esc` : fermer la surface temporaire la plus proche.

## Critère d’acceptation

**Given** un viewport de 1024 px, un Incident actif et Studio interdit,  
**When** l'utilisateur parcourt le Header au clavier,  
**Then** Command/Investigate/Govern restent textuels, Studio n'expose aucune donnée, le Context Bar conserve tenant et Incident et le focus suit l'ordre visuel.
