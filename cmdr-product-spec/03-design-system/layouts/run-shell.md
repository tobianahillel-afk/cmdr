---
id: layout-run-shell
domain: 03-design-system
status: draft
owner: Design System Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-UX-001
  - REQ-UX-002
  - REQ-AI-007
  - REQ-PROD-015
---

# Run Shell

## 1. Objectif

Fournir la composition canonique pour cette activité sans redéfinir les objets ou permissions.

## 2. Utilisateurs

Response Operator, Incident Commander, approbateurs.

## 3. Anatomy

statut/phase → cible/initiateur/Decision → opérations → progression → résultats/erreurs → vérification/rollback → audit.

## 4. Zones obligatoires

phase, cible, initiateur, Decision source, opérations, statut et trace.

## 5. Zones optionnelles

console, Inspector, comparaison avant/après.

## 6. Navigation

Global Header stable ; navigation locale liée au produit ; aucun changement de produit implicite.

## 7. Context Bar

Sous le Header, affiche tenant, environnement et objets pertinents avec overflow priorisé.

## 8. Canvas

timeline d’exécution et résultat par cible.

## 9. Inspector

étape, cible ou erreur sélectionnée.

## 10. Console

optionnelle, lecture par défaut ; commande soumise à permission.

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

standard ; compact pour logs.

## 17. Responsive

progression et contrôles critiques conservés ; logs plein écran.

## 18. Accessibilité

Landmarks, skip links, focus visible, noms de régions, alternative aux canvases et annonces des mises à jour.

## 19. États

Loading conserve l'anatomie ; Empty explique ; Partial nomme les sources ; Error garde les données valides ; Offline limite les mutations ; Permission denied ne fuit rien.

## 20. IA

run distingue workflow, agent, règle et humain ; pause/interruption visibles.

## 21. Règles de composition

L'action principale est identifiable, les métadonnées d'ownership restent visibles, les overlays ne remplacent pas l'activité.

## 22. Anti-patterns

spinner sans phase, succès global masquant échecs, rollback caché.

## 23. Exemple

Une sélection ouvre l'Inspector sans réinitialiser le canvas, puis `Esc` rend le focus à la source.

## 24. Critère d’acceptation

**Given** un utilisateur autorisé, thème sombre, densité définie et viewport compatible,  
**When** il ouvre, redimensionne et ferme une région au clavier,  
**Then** l'objectif reste visible, le focus est prévisible, les tailles respectent les limites, le contexte et le travail sont conservés et aucune surface concurrente n'apparaît.
