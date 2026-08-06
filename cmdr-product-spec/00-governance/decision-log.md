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

| ADR | Sujet | Statut | Résumé |
|---|---|---|---|
| ADR-0001 | séparation des produits | draft | six produits et Shared Capabilities avec exclusions |
| ADR-0002 | source unique | draft | une définition et un propriétaire par concept |
| ADR-0003 | chaîne nominale | draft | Telemetry Event vers Result |
| ADR-0004 | contrat d’écran | draft | sections substantielles et états testables |
| ADR-0005 | page, vue, mode, filtre | draft | éviter une page par filtre/représentation |
| ADR-0006 | propriété Endpoint Agent | draft | composant distinct, flotte administrée par Settings |
| ADR-0007 | capacités agentiques | draft | Studio possède les objets agentiques |

## Open decision index update
- `OPEN-017 — Detection runtime, target language and portability strategy` remains open and Detection-only.
- `OPEN-018 — Threat intelligence ontology, interoperability and exchange strategy` remains open; no ontology, standard, protocol, provider or representation is selected.
- `OPEN-019 — Intelligence dissemination, releasability, sharing and consumer access policy` is created open.
- OPEN-019 selects no final policy, audience model, cross-tenant rule, client-sharing model, external destination or publication authority.
- Options remain: role-based internal dissemination; marking/releasability-based; tenant-isolated controlled sharing; client-specific delivery; external sharing only via Govern; hybrid by classification and audience.
- Open decisions: **18**. No decision is closed; OPEN-009 remains the only historically resolved item.
