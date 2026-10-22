---
id: workspace-model
domain: 04-experience-architecture
status: draft
owner: UX Architecture Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-UX-001
  - REQ-UX-002
  - REQ-UX-003
  - REQ-UX-004
  - REQ-UX-005
  - REQ-PROD-007
---

# Modèle de workspace

## Définition

Un workspace est une activité longue, multi-étapes ou multi-objet dont l'état survit aux panneaux, aux transitions autorisées et aux interruptions. Il possède une route, un objectif, un propriétaire, une stratégie de persistance et un shell.

## Anatomie commune

```text
Global Header
Context Bar
┌──────────────┬─────────────────────────────┬──────────────┐
│ Local nav /  │ Main Workspace              │ Inspector    │
│ Explorer     │ objectif, canvas, données   │ sélection    │
└──────────────┴─────────────────────────────┴──────────────┘
                         Optional Bottom Console
```

Un workspace peut omettre navigation locale, Inspector ou console, mais ne peut pas recréer leur fonction avec un panneau concurrent.

## État persistant

Route, contexte autorisé, vue, mode, filtres, tri, sélection, onglets, scroll, tailles de panneaux, requête non exécutée non sensible et travail non enregistré récupérable. Les secrets ne sont jamais persistés dans l'URL ou le stockage non protégé.

## Panneaux

- un seul Inspector droit ;
- un seul canvas principal ;
- maximum deux panneaux auxiliaires ouverts ;
- les panneaux ont min/max, raccourci, état réduit et restauration ;
- une modal ne contient jamais un workspace.

## Fenêtres et onglets

Les liens internes s'ouvrent dans le même shell par défaut. Un nouvel onglet explicite reçoit un deep link permission-aware, pas un clone opaque de session. Les onglets de workspace signalent les modifications non enregistrées.

## Critère d’acceptation

**Given** un Technical Workbench avec explorer, canvas, Inspector et console,  
**When** l'utilisateur réduit l'explorer, redimensionne l'Inspector, change de produit puis revient,  
**Then** le contexte autorisé, les tailles, la sélection et les onglets sont restaurés sans empiler de panneau.
