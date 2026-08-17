---
id: layout-decision-shell
domain: 03-design-system
status: draft
owner: Design System Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-UX-001
  - REQ-UX-002
  - REQ-PROD-015
  - REQ-SEC-002
---

# Decision Shell

## 1. Objectif

Fournir la composition canonique pour cette activité sans redéfinir les objets ou permissions.

## 2. Utilisateurs

approbateurs techniques, Business Owner, auditeurs.

## 3. Anatomy

demande → Evidence → impact/risque → alternatives → conditions → autorité → Decision → trace.

## 4. Zones obligatoires

source, portée, options, conditions, séparation des rôles, rollback.

## 5. Zones optionnelles

Inspector, comparaison, policy explanation.

## 6. Navigation

Global Header stable ; navigation locale liée au produit ; aucun changement de produit implicite.

## 7. Context Bar

Sous le Header, affiche tenant, environnement et objets pertinents avec overflow priorisé.

## 8. Canvas

document décisionnel structuré, pas mosaïque de cartes.

## 9. Inspector

détail source sans déplacer la décision.

## 10. Console

interdite.

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

standard, comfortable pour justification longue.

## 17. Responsive

sections empilées ; actions sticky non masquantes ; Evidence en drawer.

## 18. Accessibilité

Landmarks, skip links, focus visible, noms de régions, alternative aux canvases et annonces des mises à jour.

## 19. États

Loading conserve l'anatomie ; Empty explique ; Partial nomme les sources ; Error garde les données valides ; Offline limite les mutations ; Permission denied ne fuit rien.

## 20. IA

résumé attribué seulement ; aucune apparence d’autorité.

## 21. Règles de composition

L'action principale est identifiable, les métadonnées d'ownership restent visibles, les overlays ne remplacent pas l'activité.

## 22. Anti-patterns

rouge permanent, bouton Approve isolé, théâtre judiciaire, recommandation=Decision.

## 23. Exemple

Une sélection ouvre l'Inspector sans réinitialiser le canvas, puis `Esc` rend le focus à la source.

## 24. Critère d’acceptation

**Given** un utilisateur autorisé, thème sombre, densité définie et viewport compatible,  
**When** il ouvre, redimensionne et ferme une région au clavier,  
**Then** l'objectif reste visible, le focus est prévisible, les tailles respectent les limites, le contexte et le travail sont conservés et aucune surface concurrente n'apparaît.
