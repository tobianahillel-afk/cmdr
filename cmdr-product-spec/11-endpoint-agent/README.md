---
id: endpoint-agent-readme
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-11
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
- records locaux et résultats d'opérations ;
- EPT-1 local identity/registration, binding, observed platform/version/inventory, health/connectivity/state, capability availability and provenance semantics.

## Consomme
Fleet et policies Settings, workflow Investigate, Decisions/Response Runs Govern, contexte Command.

## Exclusions
- ne possède pas la flotte administrative ;
- ne possède pas Endpoint Policy ou enrollment administration ;
- ne possède pas Incident, Case, Decision ou workflow métier ;
- ne constitue pas à lui seul une preuve de delivery EDR complet.

## Capability Specification
EPT-1 — **Enrollment, Inventory, Health and Platform Foundations** allocates `CAP-EPT-001..014`, all `draft / defined / planned`, with **14 capabilities / 378 sections / 84 mandatory tables**. EPT-1 is an execution lot under Delivery Roadmap Phase 5, not a Roadmap Phase.

Future lots EPT-2 Telemetry, EPT-3 Detection/Investigation, EPT-4 Collection/Live Response, EPT-5 Containment/Verification and EPT-6 Updates/Resilience/Security remain **NOT STARTED**.

See [`capabilities/README.md`](capabilities/README.md), [`capability-map.md`](capability-map.md), [`object-consumption-map.md`](object-consumption-map.md), [`action-classification.md`](action-classification.md), [`automation-and-ai-model.md`](automation-and-ai-model.md), [`cross-product-links.md`](cross-product-links.md) and [`information-architecture.md`](information-architecture.md).

## Transitions principales
Télémétrie → détection ; collecte → Artifact/Evidence ; commande autorisée → résultat vérifiable. These future paths are not started by EPT-1.

## Place de l'IA
Les fonctions essentielles de l'agent ne dépendent pas d'un modèle. Toute inference future reste classifiée, explicable et gouvernée. EPT-1 requires deterministic/manual alternatives.

## Delivery classification
**Target :** EDR natif complet. **Current evidence :** `planned`. EPT-1 defines documentary foundations only; it proves no engine, supported platform or release. `OPEN-008` remains open.

## Sources
- [`../01-product-vision/product-boundaries.md`](../01-product-vision/product-boundaries.md)
- [`../00-governance/ownership-register.md`](../00-governance/ownership-register.md)
- [`information-architecture.md`](information-architecture.md)
- [`product-definition.md`](product-definition.md)

## Critère d'acceptation
Un module de ce produit ne peut revendiquer un objet ou une capability exclue sans mise à jour des frontières, du registre de propriété et d'une ADR lorsqu'elle est transversale.