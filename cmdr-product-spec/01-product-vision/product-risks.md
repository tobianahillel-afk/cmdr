---
id: product-risks
domain: 01-product-vision
status: draft
owner: Head of Product
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-006
  - REQ-PROD-012
  - REQ-AI-001
  - REQ-UX-001
---
# Registre des risques produit

| ID | Risque | Probabilité | Impact | Signaux | Mitigation | Propriétaire | Requirement IDs | Statut |
|---|---|---|---|---|---|---|---|---|
| RISK-PROD-001 | Sur-extension du périmètre | high | high | backlog sans owner, products universels | frontières, non-goals, classification | Head of Product | REQ-PROD-013..019 | active |
| RISK-PROD-002 | Promesse native prématurée | high | high | usage de native sans moteur/preuve | capability map, release evidence | Product Architecture | REQ-PROD-012,019 | active |
| RISK-PROD-003 | Duplication entre produits | high | high | objets ou files concurrents | ownership register, ADR | Product Architecture | REQ-PROD-006,009 | active |
| RISK-PROD-004 | Automation Run / Response Run confondus | high | high | audit ou permissions partagés | OPEN-015, modèle Phase 7 | Studio/Govern | REQ-PROD-062 | open |
| RISK-PROD-005 | IA trop centrale | medium | high | workflows bloqués sans provider | non-AI acceptance tests | Head of Product | REQ-PROD-010, REQ-AI-001 | active |
| RISK-PROD-006 | Autonomie opaque | medium | critical | Tool Calls invisibles, auto-approval | Studio assurance et Govern | Security/Studio | REQ-AI-004..010 | active |
| RISK-PROD-007 | Complexité excessive | high | high | navigation par outils, trop de workspaces | UX rules et pilots | UX Architecture | REQ-UX-001..010 | active |
| RISK-PROD-008 | Densité UI incontrôlée | high | medium | panneaux empilés, états illisibles | workbench constraints, role research | Design Lead | REQ-UX-004,005 | active |
| RISK-PROD-009 | Perte de contexte | high | high | ressaisie entre produits | context propagation et analytics | UX Architecture | REQ-PROD-008 | active |
| RISK-PROD-010 | Gouvernance trop lente | medium | high | Action Requests expirées, contournements | policies, classes, metrics | Govern Product Lead | REQ-PROD-004 | active |
| RISK-PROD-011 | Govern contourné | medium | critical | action classe 3/4 lancée localement | permission, Action Request, audit | Security/Govern | REQ-SEC-001..003 | active |
| RISK-PROD-012 | Produits devenant des silos | medium | high | contexte et résultats non partagés | operating model, journeys | Product Architecture | REQ-PROD-008 | active |
| RISK-PROD-013 | Endpoint Agent sur-promis | high | high | README dit complet sans release | planned label et roadmap | Endpoint Product Lead | REQ-PROD-018 | mitigated Phase 1 |
| RISK-PROD-014 | Studio devient produit principal | medium | high | activité opérationnelle déplacée | ADR-0007 et boundaries | Studio Product Lead | REQ-AI-002 | active |
| RISK-PROD-015 | Dette documentaire | high | high | templates génériques, placeholders | lots, quality gates | Documentation Governance | REQ-PROD-006 | active |
| RISK-PROD-016 | Décisions techniques trop tôt | medium | high | UX orientée moteur ou vendor | phase gates | Software Architecture | REQ-PROD-019 | active |
| RISK-PROD-017 | Dépendance excessive aux intégrations | medium | high | pas de replacement plan | temporary-integration metadata | Product Architecture | REQ-PROD-037..047 | active |
| RISK-PROD-018 | Branding générique ou cyberpunk | medium | medium | glow, HUD, shields | Phase 2 and forbidden directions | Brand Lead | REQ-BRAND-001..008 | active |
| RISK-PROD-019 | Mauvais équilibre expert / moins technique | high | medium | trop abstrait ou trop dense | role research, progressive disclosure | UX Research | REQ-PROD-007,057 | active |

## Revue

Les risques `critical` ou `high/high` sont revus avant validation d'une phase. Une mitigation documentaire ne signifie pas que le risque d'implémentation est fermé.

## Critère d'acceptation

**Given** une capability annoncée native,  
**When** aucune preuve de moteur ou release n'existe,  
**Then** `RISK-PROD-002` est déclenché, la capability est reclassée `planned` ou `integrated`, et la formulation publique est corrigée.
