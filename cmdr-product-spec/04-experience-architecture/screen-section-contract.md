---
id: screen-section-contract
domain: 04-experience-architecture
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: canonical
---
# Contrat de spécification d’écran

Chaque fichier sous un dossier `screens/` possède un front matter valide et exactement les sections obligatoires suivantes:

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

Une section `Transitions interproduits` complète le contrat lorsqu’une transition existe. L’écran référence les objets, composants et permissions canoniques au lieu de les redéfinir.
