---
id: decision-log
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-06
source-of-truth: index
---
# Journal des décisions

Ce fichier indexe les ADR et décisions ouvertes ; il ne remplace ni leur contenu ni les sources canoniques.

| ADR | Sujet | Requirement IDs | Statut | Résumé |
|---|---|---|---|---|
| [ADR-0001](adr/ADR-0001-product-separation.md) | séparation des produits | REQ-PROD-013..018 | draft | six produits et Shared Capabilities avec exclusions |
| [ADR-0002](adr/ADR-0002-single-source-of-truth.md) | source unique | REQ-PROD-006, REQ-PROD-009 | draft | une définition et un propriétaire par concept |
| [ADR-0003](adr/ADR-0003-canonical-object-chain.md) | chaîne nominale | REQ-OBJ-001..010 | draft | progression Telemetry Event vers Result |
| [ADR-0004](adr/ADR-0004-screen-specification-contract.md) | contrat d'écran | REQ-UX-010 | draft | sections substantielles, états adaptés et AC observables |
| [ADR-0005](adr/ADR-0005-page-view-mode-filter-rules.md) | page, vue, mode, filtre | REQ-UX-001, REQ-UX-008, REQ-UX-009 | draft | éviter une page par filtre ou représentation |
| [ADR-0006](adr/ADR-0006-endpoint-agent-ownership.md) | propriété Endpoint Agent | REQ-PROD-018, REQ-OBJ-008 | draft | composant distinct, flotte administrée par Settings |
| [ADR-0007](adr/ADR-0007-agentic-studio-placement.md) | capacités agentiques | REQ-AI-002, REQ-OBJ-009 | draft | Studio possède les objets agentiques, produits opérationnels consommateurs |

## Open decision index update
- `OPEN-017 — Detection runtime, target language and portability strategy` remains open and Detection-only.
- `OPEN-018 — Threat intelligence ontology, interoperability and exchange strategy` remains open; no ontology, standard, protocol, provider, exchange representation or implementation is selected.
- `OPEN-019 — Intelligence dissemination, releasability, sharing and consumer access policy` is created open.
- OPEN-019 selects no final policy, audience model, cross-tenant rule, client-sharing model, external destination or publication authority.
- Options remain: role-based internal dissemination; marking/releasability-based; tenant-isolated controlled sharing; client-specific delivery; external sharing only via Govern; hybrid by classification and audience.
- Open decisions: **18**. No decision is closed; OPEN-009 remains the only historically resolved item.

Aucune ADR ou décision ouverte n'est approuvée ou fermée dans cette phase.
