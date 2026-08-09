---
id: roadmap-phase-5-studio-and-endpoint
domain: 18-roadmap-and-releases
status: draft
owner: Product Operations Lead
updated: 2026-08-09
source-of-truth: canonical
---
# Phase 5 Studio And Endpoint

## Objectif
Définir phase 5 Studio and Endpoint pour CMDR tout en conservant Studio et Endpoint comme deux capability domains séparés.

## Périmètre
Document canonique du domaine Roadmap. Il ne remplace pas les sources de vérité propriétaires et ne transforme aucun execution lot en sous-phase numérotée.

## Propriétaire fonctionnel
Product Operations Lead.

## Objets concernés
- Concepts du document
- Références canoniques liées

## Fonctionnalités historiques prévues
- Skills/Agents/workflows.
- Assurance/control room.
- EDR telemetry/detection/response.
- Fleet administration comme dépendance Settings, sans transfert d'ownership vers Endpoint.

## UX et interactions
- Navigation par liens stables.
- Contenu lisible en thème clair et sombre.
- Aucune duplication des définitions externes.

## Permissions
Les modifications suivent `../14-security-permissions-and-trust/permission-model.md` lorsque le document décrit une capacité exécutable.

## États
Le statut documentaire suit la gouvernance canonique; les états métier restent dans leurs sources propriétaires.

## Dépendances
- `../00-governance/source-of-truth-policy.md`
- `dependency-roadmap.md`
- Studio/Settings/Govern/Shared/Endpoint ownership sources.

## STD-1 current execution state
- Preflight: COMPLETE — 80/80.
- STD-1 — Studio Foundations — Tools, Skills, Library and Ownership Contracts: **functional content published; post-publication evidence recorded separately**.
- `CAP-STD-001..016`: 16 capabilities / 432 sections / 96 mandatory tables.
- Studio capability specification: PARTIAL.
- STD-2 Workflow Builder & Orchestration: NOT STARTED.
- STD-3 Agents, Human Gates & Runtime Control: NOT STARTED.
- STD-4 Assurance & Lifecycle: NOT STARTED.
- Endpoint capability specification: NOT STARTED.
- Delivery Roadmap Phase 5 overall: PARTIAL.

STD-1/2/3/4 are execution-lot planning labels only. No Phase 5A/5B/5C/5D exists.

## Implementation boundary
No API, protocol, product code, final Tool/Tool Call physical schema, final RBAC/ABAC, provider/runtime selection, detailed screen rewrite or Endpoint implementation is introduced by STD-1.

## Questions ouvertes
OPEN-007/013/015 remain open where consumed. OPEN-008 remains Endpoint/platform support and is not resolved by STD-1. Permission namespace normalization remains an implementation/detail issue, not a new blocking OPEN.

## Next candidate
STD-2 — Workflow Builder & Orchestration. **Do not begin STD-2 implicitly.**
