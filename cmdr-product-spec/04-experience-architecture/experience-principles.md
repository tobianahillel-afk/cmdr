---
id: experience-principles
domain: 04-experience-architecture
status: draft
owner: UX Architecture Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-003
  - REQ-PROD-007
  - REQ-PROD-008
  - REQ-UX-001
  - REQ-UX-006
  - REQ-AI-001
---

# Principes d’expérience

| Principe | Règle observable | Anti-pattern |
|---|---|---|
| Objectif avant structure | chaque workspace annonce l'activité et l'action principale | page créée pour un filtre |
| Contexte explicite | tenant, environnement et objet actif sont visibles | changement silencieux de tenant |
| Retour réel | Back restaure l'état de travail précédent | retour vers l'accueil produit |
| Une sélection, un Inspector | le panneau droit suit l'objet sélectionné | panneau droit spécifique par module |
| Densité contrôlée | les informations sont hiérarchisées en trois niveaux | dashboard de cartes surdimensionnées |
| Action attribuable | owner, source, acteur/run et prochaine étape restent visibles | production automatisée anonyme |
| IA optionnelle | chemins manuels et déterministes restent disponibles | chat central bloquant |
| Permission à destination | le contexte ne contourne jamais l'autorisation | préchargement de données interdites |
| Stable sous interruption | filtres, sélection et travail non enregistré sont protégés | refresh complet lors d'une mise à jour |
| Accessibilité structurelle | clavier, focus, alternatives et annonces sont inclus dans le contrat | correctif accessibilité après conception |

## Arbitrage

Intégrité et responsabilité précèdent la vitesse ; préservation du travail précède la simplification de route ; une décision ouverte reste visible plutôt qu'incorporée dans un défaut silencieux.

## Critère d’acceptation

**Given** un utilisateur sans fournisseur de modèle, au clavier, dans un Case dense,  
**When** il sélectionne une Evidence, ouvre l’Inspector puis revient à l’Incident source,  
**Then** le Case reste utilisable, le focus est restauré, le contexte est conservé, aucune donnée non autorisée n’apparaît et aucun chat n’est nécessaire.
