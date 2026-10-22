---
id: shared-capabilities-readme
domain: 12-shared-capabilities
status: draft
owner: Shared Capabilities Product Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-OBJ-010
  - REQ-OBJ-011
  - REQ-PROD-019
---
# Shared Capabilities

## Mission

Fournir des services et moteurs communs à plusieurs produits sans dupliquer leur définition ni prendre possession de leurs workflows métier.

## Possède notamment

- Reporting Engine ;
- Saved Views génériques ;
- linking, timeline, graph, search, notifications, export et autres services communs selon leur classification.

## Ne possède pas

- Incident, Case, Decision ou Response Run ;
- Work Queue Saved Views propres à Command ;
- l'autorité métier des produits consommateurs ;
- une navigation de produit parallèle.

## Règle de delivery

Chaque capability doit indiquer `native`, `integrated`, `temporary-integration`, `planned` ou `out-of-scope`. La présence d'un fichier ne prouve pas un moteur livré.

## Consommation

Les produits configurent l'usage local et conservent leurs décisions. Ils ne recréent pas un reporting engine, saved-view engine ou Inspector concurrent.

## Critère d'acceptation

**Given** une Saved View hors Work Queue,  
**When** son propriétaire est recherché,  
**Then** Shared Capabilities est la source générique ; les règles Work Queue restent exclusivement Command.
