---
id: foundation-responsive
domain: 03-design-system
status: draft
owner: Design System Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-UX-004
  - REQ-UX-005
  - REQ-PROD-007
---

# Responsive

## Plages Draft

- compact : `<768 px`;
- standard : `768–1199 px`;
- wide : `1200–1599 px`;
- ultra-wide : `≥1600 px`.

## Comportement

| Région | Compact | Standard | Wide/Ultra |
|---|---|---|---|
| product nav | menu textuel | Header condensé | Header complet |
| local nav | drawer plein écran | collapsed/overlay | rail fixe |
| Context Bar | objet actif + overflow | niveaux prioritaires | chaîne complète utile |
| Inspector | plein écran | drawer 360 px | panneau 320–560 px |
| table | colonnes essentielles + scroll | pinning limité | colonnes configurées |
| console | plein écran/indisponible expliqué | bottom sheet | panneau bas |
| modal | plein écran si nécessaire | max 640 px | max 720 px |
| workbench | stacked/read-only partiel | deux régions | multi-panneaux contractuel |

Un petit écran ne simule pas un workbench complet. Les opérations non sûres sur compact sont en lecture seule ou indisponibles avec explication et lien vers environnement adapté.
