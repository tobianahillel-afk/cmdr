---
id: cognitive-load
domain: 04-experience-architecture
status: draft
owner: UX Research Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-007
  - REQ-UX-004
  - REQ-UX-005
  - REQ-PROD-057
---

# Charge cognitive

## Budget d'attention Draft

- une action principale par région ; une seconde seulement si les résultats sont mutuellement exclusifs ;
- maximum deux panneaux auxiliaires ouverts ;
- maximum six onglets techniques visibles ;
- maximum trois niveaux de navigation affichés simultanément ;
- statuts critiques regroupés, jamais répétés dans chaque carte ;
- la page n'utilise pas de cartes lorsque table, section ou règle suffit.

## Densité et rôle

Les valeurs par défaut suivent l'activité : Queue/Workbench `compact`, Case/Run `standard`, Decision/Settings `standard`, lecture longue `comfortable`. L'utilisateur peut changer la densité dans les limites d'accessibilité. `OPEN-010` reste ouverte pour les préférences finales par rôle.

## Réduction du bruit

Actions secondaires dans menus nommés ; métadonnées regroupées ; mises à jour live annoncées sans déplacer la sélection ; aucune animation ambiante ; avertissements proportionnés à la conséquence.

## Critère d’acceptation

**Given** une Queue en densité compacte avec 50 lignes visibles,  
**When** trois mises à jour arrivent,  
**Then** les lignes sont patchées sans refresh, la sélection ne saute pas, la fraîcheur est annoncée et aucune carte ou toast répétitif n'est ajouté.
