---
id: requirements-traceability-matrix
domain: 00-governance
status: draft
owner: QA and Traceability Lead
updated: 2026-08-03
source-of-truth: source-material
requirements:
  - REQ-PROD-001
  - REQ-PROD-062
  - REQ-UX-001
  - REQ-UX-010
  - REQ-OBJ-011
  - REQ-OBJ-012
---

# Requirements Traceability Matrix

## Interprétation

Les 122 Requirement IDs sources sont conservés. `conform` signifie que la phase propriétaire actuelle fournit des règles substantielles, testables et sans contradiction active ; cela ne prouve ni écrans détaillés ni implémentation.

## Couverture

| État | Après Phase 1 | Après Phase 2 | Après Phase 3 |
|---|---:|---:|---:|
| conform | 81 | 88 | 99 |
| partial | 32 | 25 | 20 |
| absent | 6 | 6 | 3 |
| contradictory | 3 | 3 | 0 |
| total | 122 | 122 | 122 |

## Traçabilité détaillée Phase 3

| Requirement | Source | Canonical evidence | Avant | Après | Couverture | Dépendants | Preuve | Question |
|---|---|---|---|---|---|---|---|---|
| REQ-UX-001 | ux-and-navigation-decisions.md | 04-experience-architecture/information-architecture.md; global-navigation.md; page-view-mode-filter-rules.md | partial | conform | normative classification and routing rules | products, screens, screen register | page/workspace/view/mode/filter decision tree | — |
| REQ-UX-002 | ux-and-navigation-decisions.md | 03-design-system/components/inspector.md; layouts/* | partial | conform | one Inspector and eight shell contracts | all products and pilot screens | single right Inspector, panel limits, keyboard | — |
| REQ-UX-003 | ux-and-navigation-decisions.md | 03-design-system/foundations/tokens.md; components/* | partial | conform | three-level token architecture and shared components | Design System and products | 41 primitives, 45 semantics, 24 component contracts | OPEN-001..004 |
| REQ-UX-004 | ux-and-navigation-decisions.md | 03-design-system/foundations/accessibility.md; responsive.md; density.md | absent | conform | exact Draft density/responsive/accessibility rules | all shells/components | viewport ranges, sizes and keyboard contracts | OPEN-010 |
| REQ-UX-005 | ux-and-navigation-decisions.md | 03-design-system/foundations/theme-contract.md; data-visualization.md | absent | conform | light/dark, contrast and visualization contracts | all UI consumers | semantic themes and non-color alternatives | — |
| REQ-UX-006 | ux-and-navigation-decisions.md | 04-experience-architecture/context-preservation.md; 03-design-system/components/context-bar.md | partial | conform | context model separated from visual component | all products and transitions | transmission table, overflow and permission rules | — |
| REQ-UX-007 | ux-and-navigation-decisions.md | 04-experience-architecture/history-and-back.md; cross-product-transitions.md | partial | conform | real Back and return-origin model | products, deep links, journeys | focus/filters/scroll restored; errors covered | — |
| REQ-UX-008 | ux-and-navigation-decisions.md | 06-command/modules/incidents-and-work-queue/README.md; saved-views.md; deprecated screen aliases | contradictory | conform | one Work Queue workspace and six views | Command navigation/screens/register | five legacy screen specs deprecated | — |
| REQ-UX-009 | ux-and-navigation-decisions.md | 12-shared-capabilities/saved-views.md; 03-design-system/components/saved-views.md | partial | conform | generic persistence separated from product catalog | all products and Command | permission re-evaluation and URL-safe state | — |
| REQ-UX-010 | ux-and-navigation-decisions.md | 04-experience-architecture/screen-section-contract.md; screen-state-requirements.md | absent | partial | general contract substantive; 61 screens still later-phase | Phase 6 screens and validation | GWT contract and six states | OPEN-010 |
| REQ-OBJ-011 | canonical ownership decisions | 12-shared-capabilities/saved-views.md; ownership-register.md | contradictory | conform | Shared owns generic view mechanism | all products | single generic source; engine path deprecated | — |
| REQ-OBJ-012 | canonical ownership decisions | 06-command/modules/incidents-and-work-queue/saved-views.md; ownership-register.md | contradictory | conform | Command owns Work Queue catalog | Command | six system views; none are pages | — |
| REQ-AI-007 | ai-and-automation-constraints.md | 03-design-system/components/automation-tray.md; layouts/run-shell.md; builder-shell.md | partial | partial | UI interruption and trace defined; runtime later | Studio, Govern, implementation | pause/stop/resume and Tool Calls visible | OPEN-007, OPEN-015 |
| REQ-AI-010 | ai-and-automation-constraints.md | 03-design-system/patterns/attribution-and-provenance.md; components/trace.md | partial | partial | visual provenance substantive; functional model later | all products | human/rule/engine/workflow/agent/external distinction | — |
| REQ-BRAND-007 | brand-and-visual-decisions.md | 03-design-system/foundations/typography.md | partial | partial | metrics defined; families unresolved | all components | 13 roles and aliases OPEN-004 | OPEN-004 |
| REQ-PROD-057 | unresolved-decisions.md | 04-experience-architecture/role-based-defaults.md; 03-design-system/foundations/density.md | conform | conform | open decision preserved with Draft activity defaults | all products | compact/standard/comfortable and overrides | OPEN-010 |

## Inventaire complet — état Phase 3

### conform — 99

`REQ-AI-001`, `REQ-AI-002`, `REQ-AI-003`, `REQ-AI-004`, `REQ-AI-005`, `REQ-AI-006`, `REQ-AI-011`, `REQ-BRAND-001`, `REQ-BRAND-002`, `REQ-BRAND-003`, `REQ-BRAND-004`, `REQ-BRAND-005`, `REQ-BRAND-006`, `REQ-BRAND-008`, `REQ-INV-001`, `REQ-INV-002`, `REQ-INV-003`, `REQ-INV-004`, `REQ-INV-005`, `REQ-INV-006`, `REQ-OBJ-001`, `REQ-OBJ-002`, `REQ-OBJ-003`, `REQ-OBJ-004`, `REQ-OBJ-005`, `REQ-OBJ-006`, `REQ-OBJ-007`, `REQ-OBJ-008`, `REQ-OBJ-010`, `REQ-OBJ-011`, `REQ-OBJ-012`, `REQ-PROD-001`, `REQ-PROD-002`, `REQ-PROD-003`, `REQ-PROD-004`, `REQ-PROD-005`, `REQ-PROD-009`, `REQ-PROD-011`, `REQ-PROD-012`, `REQ-PROD-013`, `REQ-PROD-014`, `REQ-PROD-015`, `REQ-PROD-016`, `REQ-PROD-017`, `REQ-PROD-018`, `REQ-PROD-019`, `REQ-PROD-021`, `REQ-PROD-022`, `REQ-PROD-023`, `REQ-PROD-024`, `REQ-PROD-025`, `REQ-PROD-026`, `REQ-PROD-027`, `REQ-PROD-028`, `REQ-PROD-029`, `REQ-PROD-030`, `REQ-PROD-031`, `REQ-PROD-032`, `REQ-PROD-033`, `REQ-PROD-034`, `REQ-PROD-035`, `REQ-PROD-036`, `REQ-PROD-037`, `REQ-PROD-038`, `REQ-PROD-039`, `REQ-PROD-040`, `REQ-PROD-041`, `REQ-PROD-042`, `REQ-PROD-043`, `REQ-PROD-044`, `REQ-PROD-045`, `REQ-PROD-046`, `REQ-PROD-047`, `REQ-PROD-048`, `REQ-PROD-049`, `REQ-PROD-050`, `REQ-PROD-051`, `REQ-PROD-052`, `REQ-PROD-053`, `REQ-PROD-054`, `REQ-PROD-055`, `REQ-PROD-056`, `REQ-PROD-057`, `REQ-PROD-058`, `REQ-PROD-059`, `REQ-PROD-060`, `REQ-PROD-061`, `REQ-PROD-062`, `REQ-SEC-001`, `REQ-SEC-002`, `REQ-UX-001`, `REQ-UX-002`, `REQ-UX-003`, `REQ-UX-004`, `REQ-UX-005`, `REQ-UX-006`, `REQ-UX-007`, `REQ-UX-008`, `REQ-UX-009`

### partial — 20

`REQ-AI-007`, `REQ-AI-008`, `REQ-AI-009`, `REQ-AI-010`, `REQ-BRAND-007`, `REQ-JRN-001`, `REQ-JRN-003`, `REQ-JRN-004`, `REQ-JRN-006`, `REQ-JRN-007`, `REQ-OBJ-009`, `REQ-PROD-006`, `REQ-PROD-007`, `REQ-PROD-008`, `REQ-PROD-010`, `REQ-PROD-020`, `REQ-SEC-003`, `REQ-SEC-004`, `REQ-SEC-005`, `REQ-UX-010`

### absent — 3

`REQ-JRN-002`, `REQ-JRN-005`, `REQ-JRN-008`

### contradictory — 0

Aucun Requirement ID n'est classé contradictoire après la migration Saved Views / Work Queue.

## Mouvement Phase 3

- partial → conform : REQ-UX-001, 002, 003, 006, 007, 009 ;
- absent → conform : REQ-UX-004, 005 ;
- contradictory → conform : REQ-UX-008, REQ-OBJ-011, REQ-OBJ-012 ;
- absent → partial : REQ-UX-010.

Les trois exigences absentes restantes sont `REQ-JRN-002`, `REQ-JRN-005`, `REQ-JRN-008`, possédées par les phases parcours. Les décisions de palettes, typographie, densité finale et logo restent ouvertes.

## Règle de mise à jour

Une valeur Draft réversible peut rendre la fondation testable ; elle ne ferme pas une décision humaine. Une exigence écran reste partielle tant que les écrans propriétaires ne sont pas substantiels.
