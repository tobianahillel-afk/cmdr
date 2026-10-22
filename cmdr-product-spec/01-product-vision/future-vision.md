---
id: future-vision
domain: 01-product-vision
status: draft
owner: Head of Product
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-012
  - REQ-PROD-018
  - REQ-PROD-019
  - REQ-AI-007
---
# Vision future

## Horizon initial — target

- produits, objets et responsabilités cohérents ;
- Command, Investigate et Govern reliés par le contexte ;
- workspaces pilotes validés ;
- moteurs intégrés ou temporaires classifiés ;
- actions à risque gouvernées ;
- workflows essentiels sans IA.

## Horizon de consolidation — planned

- parcours de bout en bout et écrans majeurs cohérents ;
- Reporting, Saved Views, linking et audit partagés ;
- Detection Engineering relié aux Findings et Results ;
- Endpoint Agent expérience de flotte, inspection, collecte et action clairement spécifiée ;
- Studio Library, Builder, Assurance et Control Room.

## Horizon capacités natives avancées — target

- recherche et investigation natives plus profondes ;
- forensic, static analysis, reverse engineering, debugger et sandbox selon roadmap ;
- Endpoint Agent mature pour télémétrie, détection, Live Response et containment ;
- intégrations temporaires remplacées lorsque la valeur et la sécurité le justifient ;
- contrats et preuves de delivery pour chaque capability.

## Horizon agentique mature — proposed

- Automation Agents évalués, versionnés et supervisés ;
- Tool Calls et Automation Runs pleinement traçables ;
- Human Gates et Govern articulés sans double autorité ;
- coûts, qualité, sécurité et interruption observables ;
- collaboration agents/humains sans dépendance obligatoire à un modèle externe.

## Statuts

- `target` : destination décidée.
- `planned` : capability cible enregistrée mais non livrée.
- `proposed` : direction nécessitant validation ou preuves.
- `open` : décision non prise.

Aucun horizon ne constitue une déclaration de livraison.

## Écosystème d'intégrations

Les intégrations durables ou temporaires doivent adapter leurs données et actions au modèle CMDR, exposer leur provenance et posséder une stratégie de remplacement lorsqu'elles sont temporaires.

## Critère d'acceptation

**Given** une capability forensic listée dans l'horizon avancé,  
**When** un document produit la mentionne aujourd'hui,  
**Then** il indique sa classification courante et ne la présente pas comme native livrée sans preuve.
