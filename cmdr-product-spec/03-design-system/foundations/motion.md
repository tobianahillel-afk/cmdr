---
id: foundation-motion
domain: 03-design-system
status: draft
owner: Design System Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-BRAND-002
  - REQ-UX-004
  - REQ-AI-001
---

# Mouvement

Durées Draft : fast 100 ms, standard 160 ms, slow 240 ms, maximum 320 ms. Easing : enter `cubic-bezier(.2,.8,.2,1)`, exit `cubic-bezier(.4,0,1,1)`.

- panneau : translation courte + opacity ;
- sélection : transition de surface sans déplacement ;
- progress : déterminé ou indéterminé accessible, jamais décoratif ;
- live update : highlight ≤800 ms, sans déplacer le focus ;
- error/rollback : état textuel persistant, pas de secousse ;
- pause/reprise : changement d'icône, label et annonce.

Reduced motion supprime translation/zoom et ramène les durées à 1 ms sauf progression nécessaire. Glow, pulsation permanente, rebond, ticker, flux lumineux et animation IA magique sont interdits.
