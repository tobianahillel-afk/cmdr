---
id: page-view-mode-filter-rules
domain: 04-experience-architecture
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-UX-001
  - REQ-UX-008
  - REQ-UX-009
---

# Règles Page, Workspace, Vue, Mode, Filtre et surfaces

## Définitions normatives

- **Page** : objectif autonome, route stable et place justifiée dans la navigation.
- **Workspace** : activité durable, complexe ou multi-étapes avec état de travail.
- **View** : sélection ou configuration enregistrée du même objectif.
- **Mode** : représentation ou interaction différente du même ensemble.
- **Filter** : restriction temporaire et visible.
- **Inspector** : explication et actions de l'objet sélectionné.
- **Drawer** : interaction secondaire courte sans perdre la page.
- **Modal** : décision courte ou confirmation bloquante.

## Arbre de décision

1. L'objectif utilisateur change-t-il réellement ? Sinon, pas de Page.
2. Le sous-ensemble change-t-il seulement ? Utiliser une View.
3. La représentation change-t-elle seulement ? Utiliser un Mode.
4. Le contenu est-il seulement réduit ? Utiliser un Filter.
5. S'agit-il du détail de la sélection ? Utiliser l'Inspector.
6. S'agit-il d'une action courte ? Utiliser Drawer ou Modal.
7. Une activité longue ou multi-objet exige un Workspace.

## Sérialisation

Vue, mode, filtres sûrs, tri et sélection peuvent être adressables. Une vue enregistre une configuration, jamais des données d'objet ou des permissions. L'ouverture réévalue l'autorisation.

## Décision Work Queue

Une route `Incidents & Work Queue`, six vues système : `All`, `Incidents`, `Tasks`, `Unassigned`, `SLA Risk`, `My Work`. Les anciens fichiers d'écran sont des alias de migration dépréciés.

## Critère d’acceptation

**Given** une demande de « page Active Runs »,  
**When** l'objectif reste la supervision des runs,  
**Then** `Active Runs` est une View du Run workspace ; aucun nouvel écran propriétaire n'est créé.
