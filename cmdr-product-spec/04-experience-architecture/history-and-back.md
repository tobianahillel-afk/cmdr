---
id: history-and-back
domain: 04-experience-architecture
status: draft
owner: UX Architecture Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-008
  - REQ-UX-007
  - REQ-UX-009
---

# Historique et retour

## Ordre de fermeture

`Esc` ferme d'abord le menu, puis la modal, le drawer, l'Inspector non épinglé ou la console maximisée. Il ne quitte jamais silencieusement le workspace.

## Back

Le bouton navigateur revient au dernier état adressable réel : route, vue, filtres, tri, sélection et scroll. Une transition interproduit ajoute une entrée d'historique et un `Return to source` explicite. Le retour ne redirige pas vers l'accueil produit sauf si l'origine n'existe plus.

## Deep link

Un deep link sans origine retourne à la destination précédente du navigateur. Après connexion ou expiration de session, la route est restaurée seulement après tenant et permissions. Après changement de tenant, l'historique incompatible est invalidé avec explication.

## Travail non enregistré

Navigation produit, Back, fermeture d'onglet ou changement de tenant déclenchent autosave lorsque contractuel ; sinon confirmation avec `Rester`, `Enregistrer`, `Abandonner`. Les modales imbriquées et drawers empilés sont interdits.

## Critère d’acceptation

**Given** une Work Queue filtrée, une ligne sélectionnée et l'Inspector ouvert,  
**When** l'utilisateur ferme l'Inspector puis utilise Back après avoir ouvert Incident Detail,  
**Then** le focus revient à la ligne, les filtres et le scroll sont restaurés et aucun retour vers l'accueil n'a lieu.
