---
id: command-module-risk-and-coverage
domain: 06-command
status: draft
owner: Command Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-005
  - REQ-PROD-013
  - REQ-PROD-021
  - REQ-UX-005
---
# Risk and Coverage

## Mission
Relier Service, Exposure, couverture, impact et contexte de risque pour la coordination, sans scanner, Detection Engineering ni score opaque.

## Utilisateurs
Incident Commander, Business/Service Owner, SOC Analyst L2, Risk/Readiness Coordinator.

## Ownership
Command possède la lecture opérationnelle et les mutations d’impact/Task autorisées. Business Service Catalog et produits sources possèdent Service, Exposure et couvertures. Investigate possède la logique détaillée de Detection Engineering.

## Capabilities
| ID | Capability | Status | Mode |
|---|---|---|---|
| CAP-CMD-201 | Service Context | defined | planned |
| CAP-CMD-202 | Exposure Overview | defined | planned |
| CAP-CMD-203 | Coverage Overview | defined | planned |
| CAP-CMD-204 | Business Impact Context | defined | planned |
| CAP-CMD-205 | Risk Prioritization Context | defined | planned |

## Shared Capabilities consommées
| Shared capability | Usage local | Source canonique |
|---|---|---|
| Business Service Catalog | Service, owner, criticité, dépendances | `../../../12-shared-capabilities/business-service-catalog.md` |
| Object Linking | relations Service/Exposure/Incident/Coverage | `../../../12-shared-capabilities/object-linking-service.md` |
| Metrics / Data Quality | définitions, fraîcheur, lacunes | Shared sources canoniques |
| Inspector / Context Bar / Trace | inspection, contexte et provenance | Design System sources |
| Reporting / Export | projections sourcées et redaction | Shared engines |

## Frontière
Les anciens modules `exposure-and-coverage` et `risk-and-business-impact` sont migrés vers cette responsabilité. Leurs écrans restent inchangés jusqu’à la phase écrans.

## Critère
**Given** une exposition affectant un Service, **When** Risk and Coverage est consulté, **Then** source, fraîcheur, couverture, impact et owner sont distincts et Command n’exécute aucun scan ni modification de règle.
