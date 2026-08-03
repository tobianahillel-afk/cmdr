---
id: dependency-register
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-012
  - REQ-PROD-019
---
# Registre des dépendances produit et documentaires

## Règles

Ce registre capture les dépendances nécessaires à une décision produit ou documentaire. Les choix de protocole, base de données, file de messages ou infrastructure finale restent hors Phase 1.

| ID | Source | Dépendant | Type | Raison | Statut | Impact | Propriétaire | Bloquant | Requirement IDs | Revue |
|---|---|---|---|---|---|---|---|---:|---|---|
| DEP-001 | `source-material/` | `01-product-vision/` | décision | appliquer mission, principes et frontières | active | vision incohérente si non propagée | Head of Product | oui | REQ-PROD-001..020 | Phase 1 |
| DEP-002 | `product-boundaries.md` | READMEs produits | frontière | empêcher la revendication concurrente d'une capability | active | doublons de responsabilité | Product Architecture | oui | REQ-PROD-013..018 | Phase 1 |
| DEP-003 | `ownership-register.md` | modèle d'objets | propriété | guider les futurs schémas et machines d'état | active | objets concurrents | Product Architecture | oui | REQ-OBJ-001..012 | Phase 7 |
| DEP-004 | `terminology-rules.md` | UX et contenu | vocabulaire | distinguer page, vue, mode, filtre, agent et run | active | navigation et objets ambigus | Content Design | oui | REQ-UX-001, REQ-OBJ-009 | Phase 3 |
| DEP-005 | ADR-0005 | Work Queue et Mission Control | UX | consolider les variantes en vues ou modes | open | écrans clones | UX Architecture | oui | REQ-UX-008, REQ-UX-009 | Phase 3/6 |
| DEP-006 | Phase 2 marque | identités Investigate/Govern/Studio | décision visuelle | palettes non décidées | open | design non validable | Brand Design Lead | non | REQ-PROD-048..051 | Phase 2 |
| DEP-007 | Capability map | modules fonctionnels | classification | ne pas annoncer une capacité planned comme native livrée | active | promesse trompeuse | Product Architecture | oui | REQ-PROD-012, REQ-PROD-019 | Phase 4 |
| DEP-008 | `OPEN-005` | workbench forensic | moteur | choisir la stratégie initiale sans bloquer l'UX activité | open | limites d'intégration inconnues | Investigate Product Lead | non pour UX, oui pour technique | REQ-PROD-052 | Phase 4/8 |
| DEP-009 | `OPEN-007` | Studio et Govern | autorité | préciser Human Gate versus Decision | open | contournement ou double approbation | Govern + Studio | oui | REQ-PROD-054 | Phase 4/7 |
| DEP-010 | `OPEN-008` | Endpoint Agent | plateforme | définir le support initial | open | expérience flotte et actions incomplètes | Endpoint Agent Product Lead | oui avant implémentation | REQ-PROD-055 | Phase 4/8 |
| DEP-011 | `OPEN-014` | Investigate | modèle | distinguer Artifact et Attachment | open | preuve et collaboration ambiguës | Investigate Product Lead | oui pour modèle détaillé | REQ-PROD-061 | Phase 7 |
| DEP-012 | `OPEN-015` | Studio / Govern | modèle | distinguer Automation Run et Response Run | open | audit et autorité confondus | Studio + Govern | oui | REQ-PROD-062, REQ-OBJ-009 | Phase 7 |
| DEP-013 | Phase 3 | six écrans pilotes | UX | établir le niveau de qualité avant les 61 écrans | planned | propagation de templates insuffisants | UX Architecture | oui | REQ-UX-010 | Phase 6 |
| DEP-014 | Phase 5 | parcours | écrans et transitions | valider le contexte de bout en bout | planned | navigation locale non cohérente | Product Architecture | oui | REQ-JRN-001..008 | Phase 5 |
| DEP-015 | Phase 7 | permissions détaillées | Security | unifier namespaces et risque après frontières fonctionnelles | planned | contrôle d'accès incohérent | Security Architecture | oui | REQ-SEC-001..005 | Phase 7 |

## Mise à jour

Une dépendance change de statut avec preuve. `resolved` signifie que les documents dépendants ont été mis à jour, pas seulement que la décision a été prise.

## Critère d'acceptation

**Given** une exigence bloquée par une décision ultérieure,  
**When** le registre est consulté,  
**Then** la source, le dépendant, le propriétaire, l'impact, la phase cible et le caractère bloquant sont explicites sans inventer un choix technique.
