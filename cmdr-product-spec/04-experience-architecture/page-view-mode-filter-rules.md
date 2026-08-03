---
id: page-view-mode-filter-rules
domain: 04-experience-architecture
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: canonical
---
# Règles Page, Vue, Mode et Filtre

## Page

Une **Page** est une destination routable avec identifiant d’écran, objectif utilisateur et contrat de permissions. Elle peut contenir plusieurs vues.

## Vue

Une **Vue** change l’ensemble ou l’organisation des objets sans changer le but principal de la page. Exemple: `Incidents`, `Tasks`, `Unassigned`, `SLA Risk`, `Team Load` dans la Work Queue.

## Mode

Un **Mode** change la manière d’interagir ou de représenter le même ensemble: table, graph, canvas, diff, compact, lecture. Un mode ne crée pas une nouvelle source de vérité.

## Filtre

Un **Filtre** réduit le jeu d’objets. Il doit être visible, supprimable, sérialisable de manière sûre et ne jamais élargir les permissions.

## Règles

- Une nouvelle Page exige un identifiant d’écran et une entrée au screen register.
- Une Vue n’est pas un écran distinct sauf si elle possède un but, une URL et des critères autonomes.
- Un Mode conserve la sélection et le contexte.
- Un Filtre ne modifie pas les objets.
- Les vues enregistrées de Work Queue sont définies uniquement dans `../06-command/modules/incidents-and-work-queue/saved-views.md`.
