---
id: theme-contract
domain: 03-design-system
status: draft
owner: Design System Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-UX-003
  - REQ-UX-004
  - REQ-UX-005
  - REQ-BRAND-003
---

# Contrat de thème

## Light — valeurs sémantiques Draft

| Token | Valeur |
|---|---|
| color.surface.canvas | `#F5F6F3` |
| color.surface.primary | `#FFFFFF` |
| color.surface.secondary | `#ECEFEB` |
| color.surface.elevated | `#FAFBF8` |
| color.text.primary | `#1D2523` |
| color.text.secondary | `#4F5B57` |
| color.text.muted | `#68736F` |
| color.border.default | `#C5CBC7` |
| color.border.strong | `#909B96` |
| color.focus.ring | `#176B5D` |
| color.selection.background | `#DDEBE6` |

## Dark — valeurs sémantiques Draft

| Token | Valeur |
|---|---|
| color.surface.canvas | `#171C1B` |
| color.surface.primary | `#202725` |
| color.surface.secondary | `#28302E` |
| color.surface.elevated | `#303936` |
| color.text.primary | `#F0F1EC` |
| color.text.secondary | `#C3CAC5` |
| color.text.muted | `#9AA49F` |
| color.border.default | `#46514D` |
| color.border.strong | `#68746F` |
| color.focus.ring | `#78C5B4` |
| color.selection.background | `#29433C` |

Aucun noir absolu. Le thème sombre distingue canvas, code, panneau, Inspector et overlay ; il n'ajoute ni glow ni saturation cyberpunk.

## Statuts — couples Draft

| Statut | Light foreground/background | Dark foreground/background |
|---|---|---|
| critical | `#9F2D2D` / `#F7E8E6` | `#FFB4AC` / `#4A2624` |
| warning | `#7A5000` / `#FFF3D9` | `#FFD27A` / `#493814` |
| success | `#245F45` / `#E4F2EA` | `#8FD3AE` / `#1F3C30` |
| information | `#245A73` / `#E3EFF5` | `#91CCE5` / `#203946` |
| pending | `#67582A` / `#F3EEDB` | `#DED095` / `#403B24` |
| running | `#1C6571` / `#E2F1F2` | `#8ED6DE` / `#1E3C40` |
| paused | `#65546F` / `#EFE8F1` | `#CFB9D9` / `#392F40` |
| offline | `#555B59` / `#ECEEED` | `#BEC5C1` / `#323836` |
| degraded | `#714A15` / `#F8EAD6` | `#E9BD7D` / `#44321D` |
| blocked | `#7B3832` / `#F4E4E2` | `#E8AAA4` / `#452925` |
| permission-denied | `#6C3F51` / `#F3E7EC` | `#DFB0C2` / `#402933` |

`automated` est une provenance, pas un statut coloré : label + type + trace. `human-action-required` utilise warning avec icône/personne et texte.

## Modes et préférences

`light`, `dark`, `system`; high contrast est une adaptation ultérieure. Le changement de thème ne modifie ni statut, ni ordre, ni données. Les réglages suivent l'utilisateur et respectent le système au premier lancement.
