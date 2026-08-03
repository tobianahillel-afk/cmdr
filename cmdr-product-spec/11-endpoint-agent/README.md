---
id: endpoint-agent-readme
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-018
  - REQ-PROD-012
  - REQ-SEC-004
  - REQ-SEC-005
---
# Endpoint Agent

## Mission

Servir de composant local sécurisé pour les capacités endpoint visées : télémétrie, détection, inspection, collecte, Live Response, containment, résilience et audit.

## Possède

- exécution locale selon contrats, policy, permission et autorité ;
- records locaux et résultats d'opérations.

## Consomme

Fleet et policies Settings, workflow Investigate, Decisions/Response Runs Govern, contexte Command.

## Exclusions

- ne possède pas la flotte administrative ;
- ne possède pas Incident, Case, Decision ou workflow métier ;
- ne constitue pas à lui seul une preuve de delivery EDR complet.

## Transitions principales

Télémétrie → détection ; collecte → Artifact/Evidence ; commande autorisée → résultat vérifiable.

## Place de l'IA

Les fonctions essentielles de l'agent ne dépendent pas d'un modèle. Toute inference future reste classifiée, explicable et gouvernée.

## Delivery classification

**Target :** EDR natif complet. **Current evidence :** `planned`. Les documents existants décrivent un périmètre cible ; ils ne prouvent ni moteur, ni plateforme supportée, ni release. `OPEN-008` conserve la question du support initial.

## Sources

- [`../01-product-vision/product-boundaries.md`](../01-product-vision/product-boundaries.md)
- [`../00-governance/ownership-register.md`](../00-governance/ownership-register.md)
- `information-architecture.md`
- `product-definition.md`

## Critère d'acceptation

Un module de ce produit ne peut revendiquer un objet ou une capability exclue sans mise à jour des frontières, du registre de propriété et d'une ADR lorsqu'elle est transversale.
