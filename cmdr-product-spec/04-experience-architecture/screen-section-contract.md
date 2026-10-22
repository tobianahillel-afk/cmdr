---
id: screen-section-contract
domain: 04-experience-architecture
status: draft
owner: UX Architecture Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-UX-010
---

# Contrat de spécification d’écran

Chaque écran actif conserve exactement les sections suivantes :

1. Objectif
2. Résultats utilisateur
3. Points d’entrée
4. Points de sortie
5. Contexte
6. Structure de page
7. Hiérarchie de l’information
8. Actions principales
9. Actions secondaires
10. Données et objets
11. Filtres et vues enregistrées
12. Inspector
13. UX et interactions
14. Clavier et accessibilité
15. Permissions
16. Audit
17. État Loading
18. État Empty
19. État Partial
20. État Error
21. État Offline
22. État Permission denied
23. Comportement responsive
24. Télémétrie produit
25. Dépendances
26. Critères d’acceptation
27. Questions ouvertes

`Transitions interproduits` est ajoutée lorsqu'applicable. Chaque section contient une application locale substantive ou `Non applicable` avec justification ; répéter le template n'est pas une preuve.

Les critères utilisent Given/When/Then et nomment rôle, tenant, environnement, viewport, thème, densité, permissions, données, état initial, action, entrée clavier/pointeur, focus, résultat, conservation, accessibilité, trace et erreur.

Un écran ne redéfinit ni objet, ni permission, ni composant, ni palette. Une View, un Mode ou un Filter ne reçoit pas d'identifiant d'écran actif.

**Given** une spécification pilote, **When** elle est évaluée, **Then** chaque section est spécifique ou justifiée, les six états sont testables et les transitions conservent contexte et retour.
