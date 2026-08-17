---
id: product-pillars
domain: 01-product-vision
status: draft
owner: Head of Product
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-001
  - REQ-PROD-002
  - REQ-PROD-004
  - REQ-PROD-010
  - REQ-AI-007
---
# Piliers produit

| Pilier | Promesse | Capabilities | Propriétaire(s) | Limite | Indicateurs conceptuels | Risque |
|---|---|---|---|---|---|---|
| Operational Command | Maintenir une situation et un ownership actionnables | Mission Control, Work Queue, impact, handover | Command | ne remplace pas l'analyse technique | temps vers owner, handovers complets | surcharge ou faux sentiment de contrôle |
| Evidence and Investigation | Transformer données et artefacts en conclusions vérifiables | Case, collection, Evidence, workbench, Detection Engineering | Investigate | ne décide pas seul la réponse élevée | Case → Finding, provenance | outils fragmentés ou sur-promis |
| Governed Response | Relier preuve, autorité, exécution, vérification et rollback | Action Request, Decision, Approval, Response Run, Result | Govern | ne devient pas file générale | décisions complètes, actions vérifiées | gouvernance trop lente ou contournée |
| Shared Context | Préserver tenant, objets, sélection et historique | Context Bar, linking, Saved Views, Reporting | Shared / Experience | ne transfère pas propriété | transitions sans perte | fuite de contexte ou données |
| Native and Extensible Capabilities | Faire évoluer moteurs natifs et intégrations sans casser l'UX | taxonomy, adapters, replacement plans | Product owners / Shared | wrapper ≠ native | capabilities classifiées | dépendance fournisseur |
| Human and Deterministic Control | Conserver humains, règles et moteurs reproductibles au premier plan | rules, policies, manual paths, APIs | Tous produits | IA non obligatoire | usage sans IA, règles testées | agent-first opaque |
| Accountable Automation | Accélérer avec des runs observables, interrompables et gouvernés | Studio, Tool Calls, Human Gates, Assurance | CMDR Studio | ne possède pas les décisions métier | runs traçables | autonomie cachée |

## Usage

Les piliers servent à arbitrer la roadmap et le positionnement. Ils ne créent pas de nouveau propriétaire. Une feature doit soutenir au moins un pilier, répondre à un problème utilisateur et respecter les frontières.

## Critère d'acceptation

**Given** une feature proposée comme « assistant autonome de réponse »,  
**When** elle est comparée aux piliers,  
**Then** elle doit démontrer Human and Deterministic Control, Accountable Automation et Governed Response ; sinon elle n'entre pas dans le scope.
