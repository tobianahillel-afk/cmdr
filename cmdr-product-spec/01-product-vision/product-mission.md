---
id: product-mission
domain: 01-product-vision
status: draft
owner: Head of Product
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-001
  - REQ-PROD-008
  - REQ-PROD-013
  - REQ-PROD-018
---
# Mission produit

## Mission

CMDR aide les équipes de cybersécurité et les responsables métier à transformer des signaux fragmentés en travail coordonné, preuves utilisables, décisions responsables et résultats vérifiables, sans perdre le contexte entre outils, rôles et étapes.

## Pour qui

La mission sert principalement :

- les équipes SOC qui priorisent et coordonnent ;
- les analystes et experts qui recherchent et établissent les faits ;
- les opérateurs et approbateurs qui évaluent et exécutent les réponses ;
- les administrateurs qui maintiennent la plateforme et l'Endpoint Agent ;
- les concepteurs d'automatisation qui construisent des capacités gouvernées ;
- les responsables métier et auditeurs qui exigent impact, autorité et preuve.

## Problèmes traités aujourd'hui au niveau produit

CMDR cherche à réduire :

- la fragmentation des consoles ;
- la ressaisie et les handovers incomplets ;
- la séparation entre preuve technique et impact métier ;
- les actions risquées déclenchées sans contexte suffisant ;
- l'opacité des automatisations et productions IA ;
- la difficulté à relier résultats, détections et amélioration.

## Manière d'accomplir la mission

- Command maintient la situation et l'ownership opérationnel.
- Investigate conduit la recherche, la collecte, l'analyse et la preuve.
- Govern encadre décision, autorité, exécution, vérification et rollback.
- Studio conçoit et supervise les automatisations optionnelles.
- Platform Settings administre le contexte de plateforme et la flotte.
- Endpoint Agent fournit progressivement les capacités locales visées.
- Shared Capabilities assure les services communs sans devenir un produit parallèle.

## Limites

La mission actuelle est produit et documentaire. Elle ne prouve pas que chaque capability est implémentée, ni que les moteurs, protocoles et infrastructures finaux sont choisis.

## Critère d'acceptation

**Given** un Incident actif, plusieurs outils sources et une Action Request,  
**When** les rôles passent de Command à Investigate puis Govern,  
**Then** chaque rôle retrouve le contexte nécessaire, comprend sa responsabilité, sait quelle action est attendue et peut retracer l'origine du résultat.
