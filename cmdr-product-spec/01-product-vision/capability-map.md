---
id: capability-map
domain: 01-product-vision
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-019
  - REQ-PROD-037
  - REQ-PROD-047
  - REQ-PROD-056
  - REQ-INV-001
  - REQ-INV-006
---
# Carte de propriété et de delivery des capabilities

## Taxonomie

- **native** : CMDR possède expérience, contrat fonctionnel et moteur principal.
- **integrated** : moteur externe durable intégré au modèle CMDR.
- **temporary-integration** : moteur externe temporaire avec remplacement ou décision.
- **planned** : capability cible non encore prouvée comme livrée.
- **out-of-scope** : exclusion explicite.

La colonne `Target` décrit la destination. La colonne `Current evidence` décrit ce que la documentation peut affirmer aujourd'hui.

| Capability | Owner | Target | Current evidence | Rationale / limite | Requirement |
|---|---|---|---|---|---|
| Operational coordination | Command | native | planned | expérience Command cible ; implémentation non prouvée | REQ-PROD-013 |
| Investigation lifecycle | Investigate | native | planned | Case et Evidence comme cible | REQ-PROD-014 |
| Governed decision and response | Govern | native | planned | Decision/Run/Result cible | REQ-PROD-015 |
| Agentic automation design | CMDR Studio | native | planned | Studio owner, runtime futur | REQ-PROD-016 |
| Platform administration | Platform Settings | native | planned | workspaces Settings cible | REQ-PROD-017 |
| Endpoint Agent EDR | Endpoint Agent | native | planned | cible complète, pas de preuve de delivery | REQ-PROD-018 |
| Reporting Engine | Shared Capabilities | native | planned | owner décidé, delivery non évaluée | REQ-PROD-045 |
| Generic Saved Views | Shared Capabilities | native | planned | conflits consommateurs à corriger | REQ-OBJ-011 |
| Work Queue Saved Views | Command | native | planned | vues à consolider Phase 6 | REQ-OBJ-012 |
| Search | Shared/Investigate | native | planned | moteur ouvert | REQ-PROD-037 |
| Investigation tools | Investigate | native | planned | expérience à spécifier | REQ-PROD-038 |
| Collection | Investigate + Endpoint | native | planned | contrat futur | REQ-PROD-039 |
| Live Response | Investigate + Endpoint, governed by Govern | native | planned | actions et classes à spécifier | REQ-PROD-040 |
| Endpoint telemetry | Endpoint Agent | native | planned | support ouvert | REQ-PROD-041 |
| Endpoint detection | Endpoint Agent / Detection Engineering | native | planned | moteur et delivery non prouvés | REQ-PROD-042 |
| Orchestration | CMDR Studio | native | planned | workflow runtime futur | REQ-PROD-043 |
| Governance | Govern | native | planned | policies et authority à détailler | REQ-PROD-044 |
| Audit | Shared / Security / owners | native | planned | ledger et immutabilité à détailler | REQ-PROD-046 |
| Automation | CMDR Studio | native | planned | Tool Call et Run à créer | REQ-PROD-047 |
| Forensics | Investigate | native | planned | initial engines OPEN-005 | REQ-INV-001 |
| Static analysis | Investigate | native | planned | moteur ouvert | REQ-INV-002 |
| Reverse engineering | Investigate | native | planned | moteur ouvert | REQ-INV-003 |
| Debugger | Investigate | native | planned | moteur et isolation ouverts | REQ-INV-004 |
| Sandbox | Investigate / Settings admin | native | planned | profiles et infrastructure ouverts | REQ-INV-005 |
| Detection Engineering | Investigate | native | planned | rules remain deterministic | REQ-INV-006 |

## Règles

- Une interface CMDR cohérente ne rend pas automatiquement le moteur natif.
- La classification `planned` reste en place jusqu'à preuve de contrat, moteur et release.
- Un moteur intégré est affiché dans la provenance, pas comme architecture principale de navigation.
- Une temporary integration indique owner, limites, dépendances et replacement plan.
- Les classifications actuelles pourront devenir `integrated` ou `temporary-integration` lorsque les moteurs initiaux seront décidés.

## Autorité de classification

Le capability owner propose la classification. Le Product Lead propriétaire et Product Architecture l'approuvent. Security revoit les capacités de confiance, d'action ou de données sensibles. Engineering valide toute revendication `native` ou `integrated` avant publication. QA and Traceability vérifie les preuves et met à jour les registres. Cette règle résout `OPEN-009`.

## Questions ouvertes

`OPEN-005`, `OPEN-011`, `OPEN-012` et les décisions techniques ultérieures peuvent modifier le moteur ou le statut courant sans modifier la propriété produit.

## Critère d'acceptation

**Given** la capability Reverse Engineering,  
**When** la carte est consultée,  
**Then** Investigate est owner, la cible est native, l'état courant est planned et aucun fournisseur n'est inventé.
