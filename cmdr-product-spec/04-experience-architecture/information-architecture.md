---
id: information-architecture
domain: 04-experience-architecture
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-013
  - REQ-PROD-014
  - REQ-PROD-015
  - REQ-PROD-016
  - REQ-PROD-017
  - REQ-PROD-018
  - REQ-UX-001
  - REQ-UX-008
---

# Architecture de l’information

## Niveaux

```text
CMDR
├── Command — situation, Incident et Work Queue
├── Investigate — Case, Evidence et Finding
├── Govern — Decision, Response Run et Result
├── CMDR Studio — automatisations et assurance
└── Platform Settings — administration et flotte
```

Endpoint Agent est un composant technique accessible par les workflows propriétaires ; il n'est pas un onglet produit principal. Shared Capabilities fournit recherche, liens, vues, jobs, notifications et collaboration sans devenir un septième produit.

## Contrat structurel

| Niveau | Possède | Routable | Navigation |
|---|---|---:|---|
| Produit | mission et frontières | oui | Global Header |
| Module | activité cohérente | oui si destination stable | navigation locale |
| Page | objectif autonome | oui | navigation ou lien profond |
| Workspace | activité durable et persistante | oui | page ou route dédiée |
| View | sous-ensemble du même travail | paramètre stable | Saved Views |
| Mode | représentation du même objet | paramètre stable | contrôle local |
| Filter | restriction temporaire | sérialisable si sûr | barre de filtres |
| Inspector | détail de sélection | état de route optionnel | panneau droit |
| Drawer | interaction secondaire | état temporaire | déclencheur local |
| Modal | décision courte bloquante | non comme destination | action explicite |

## Routes et deep links

Une route référence le produit, le workspace ou la page, le tenant et les identifiants autorisés. Les secrets, payloads, tokens d'accès et données sensibles ne sont jamais encodés. Une destination réévalue les permissions avant de restaurer le contexte.

## Work Queue

`Incidents & Work Queue` est un workspace Command unique. `All`, `Incidents`, `Tasks`, `Unassigned`, `SLA Risk` et `My Work` sont des vues système, pas six pages.

## Critère d’acceptation

**Given** une proposition de nouvelle destination « Unassigned »,  
**When** elle est classée,  
**Then** elle est rejetée comme page car elle conserve le même objectif, les mêmes objets et le même workspace ; elle devient la vue `Unassigned`.
