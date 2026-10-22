---
id: design-tokens
domain: 03-design-system
status: draft
owner: Design System Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-BRAND-003
  - REQ-BRAND-004
  - REQ-BRAND-007
  - REQ-UX-003
  - REQ-UX-004
  - REQ-UX-005
---

# Architecture des tokens

## Trois niveaux

```text
Brand sources / raw values
        ↓
Primitive tokens
        ↓
Semantic tokens
        ↓
Component tokens
```

## Primitive tokens — inventaire Draft

### Références de marque

| ID | Alias | Statut |
|---|---|---|
| color.brand.cmdr.ink | `../../02-brand/cmdr/palette.md#cmdr-ink` | decided reference |
| color.brand.cmdr.bone | `../../02-brand/cmdr/palette.md#cmdr-bone` | decided reference |
| color.brand.cmdr.moss | `../../02-brand/cmdr/palette.md#cmdr-moss` | decided reference |
| color.brand.cmdr.ember | `../../02-brand/cmdr/palette.md#cmdr-ember` | decided reference |
| color.brand.command.canvas | `../../02-brand/command/palette.md#command-canvas` | decided reference |
| color.brand.command.juniper | `../../02-brand/command/palette.md#command-juniper` | decided reference |

Les autres valeurs CMDR et Command utilisent le même mécanisme d'alias ; aucun hex n'est recopié ici.

### Slots produit non résolus

| ID | Source | Statut |
|---|---|---|
| color.product.investigate.accent | OPEN-001 | unresolved |
| color.product.govern.accent | OPEN-002 | unresolved |
| color.product.studio.accent | OPEN-003 | unresolved |
| font.family.ui | OPEN-004 | unresolved |
| font.family.mono | OPEN-004 | unresolved |

### Spacing

`space.000=0`, `025=2px`, `050=4px`, `100=8px`, `150=12px`, `200=16px`, `300=24px`, `400=32px`, `500=40px`, `600=48px`, `800=64px`, `1000=80px`.

### Dimensions et mouvement

- `size.control.compact=28px`, `standard=36px`, `comfortable=44px`;
- `size.row.compact=32px`, `standard=40px`, `comfortable=48px`;
- `radius.none=0`, `control=4px`, `panel=6px`, `overlay=8px`, `pill=999px`;
- `border.default=1px`, `border.strong=2px`;
- `motion.duration.fast=100ms`, `standard=160ms`, `slow=240ms`;
- `opacity.disabled=0.48`, `scrim=0.56`.

## Semantic tokens — 45 rôles actifs

Surfaces 6, textes 5, bordures/focus/sélection 6, interaction 5, statuts 15, provenance 4, data-viz structure 4. Les valeurs par mode sont définies dans `theme-contract.md`.

Les statuts incluent neutral, information, success, warning, critical, blocked, pending, running, paused, offline, degraded, permission-denied, automated, human-action-required et unknown.

## Component tokens — 24 contrats initiaux

Button 5, table 4, Inspector 3, modal/drawer 4, Context Bar 3, tabs 2, Evidence Card 2, Automation Tray 1. Les composants peuvent étendre cet inventaire par revue, jamais contourner le niveau sémantique.

## Métadonnées exemple

```yaml
id: component.inspector.border.default
category: component
role: inspector separation
value-or-alias: "{color.border.default}"
source: 03-design-system/components/inspector.md
status: draft
modes: [light, dark]
products: [all]
accessibility: non-text boundary; not sole state signal
allowed: [right inspector]
forbidden: [status meaning]
owner: Design System Lead
updated: 2026-08-03
```

## Comptage

- primitives documentées : 41 ;
- semantic roles : 45 ;
- component contracts initiaux : 24 ;
- aliases non résolus : 5 ;
- hex de marque recopiés : 0.
