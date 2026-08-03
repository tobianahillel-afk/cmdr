---
id: layout-global-shell
domain: 03-design-system
status: draft
owner: Design System Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-UX-001
  - REQ-UX-002
  - REQ-UX-006
---

# Global Shell

## 1. Objectif

Fournir la composition canonique pour cette activité sans redéfinir les objets ou permissions.

## 2. Utilisateurs

tous les rôles authentifiés.

## 3. Anatomy

Global Header → Context Bar → navigation locale optionnelle → workspace → Inspector optionnel → overlays globaux.

## 4. Zones obligatoires

Global Header, workspace landmark, skip links, couche Jobs/notifications.

## 5. Zones optionnelles

Context Bar, navigation locale, Inspector, Bottom Console selon shell enfant.

## 6. Navigation

Global Header stable ; navigation locale liée au produit ; aucun changement de produit implicite.

## 7. Context Bar

Sous le Header, affiche tenant, environnement et objets pertinents avec overflow priorisé.

## 8. Canvas

héberge le shell d’activité ; aucun dashboard implicite.

## 9. Inspector

un seul Inspector droit, 400 px par défaut, 320–560 px.

## 10. Console

absente par défaut ; autorisée seulement dans un shell technique ou Run.

## 11. Panneaux

Un seul Inspector ; maximum deux auxiliaires ; séparateurs focusables ; min/max et restauration.

## 12. Ordre de tabulation

Header → Context Bar → navigation locale/explorer → actions du workspace → contenu → Inspector → console → overlays.

## 13. Clavier

`Ctrl/Cmd+K` palette, `[` navigation, `]` Inspector, `Esc` ferme la surface temporaire ; chaque région expose son modèle local.

## 14. Scroll

Le workspace possède le scroll principal ; panneaux techniques peuvent scroller indépendamment sans scroll imbriqué non nommé.

## 15. Redimensionnement

Poignée accessible, flèches ±8 px, Shift ±32 px, Home minimum, End maximum. Les tailles sont mémorisées par shell et viewport.

## 16. Densité

standard par défaut ; suit le shell enfant.

## 17. Responsive

Header condensé et navigation textuelle en menu sous 768 px ; Inspector plein écran.

## 18. Accessibilité

Landmarks, skip links, focus visible, noms de régions, alternative aux canvases et annonces des mises à jour.

## 19. États

Loading conserve l'anatomie ; Empty explique ; Partial nomme les sources ; Error garde les données valides ; Offline limite les mutations ; Permission denied ne fuit rien.

## 20. IA

Command Palette et assistance contextuelle uniquement ; aucun chat global.

## 21. Règles de composition

L'action principale est identifiable, les métadonnées d'ownership restent visibles, les overlays ne remplacent pas l'activité.

## 22. Anti-patterns

actions métier dans le Header, filtres globaux, Endpoint comme onglet principal.

## 23. Exemple

Une sélection ouvre l'Inspector sans réinitialiser le canvas, puis `Esc` rend le focus à la source.

## 24. Critère d’acceptation

**Given** un utilisateur autorisé, thème sombre, densité définie et viewport compatible,  
**When** il ouvre, redimensionne et ferme une région au clavier,  
**Then** l'objectif reste visible, le focus est prévisible, les tailles respectent les limites, le contexte et le travail sont conservés et aucune surface concurrente n'apparaît.
