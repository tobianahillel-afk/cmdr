---
id: decision-log
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-16
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
| [ADR-0008](adr/ADR-0008-customers-mssp-delivery-deployment-and-cross-tenant-architecture.md) | Customers / MSSP / Delivery deployment et cross-tenant architecture | REQ-PROD-006,008,009,012,013,019,033,053; REQ-SEC-001 | validated | Internal + Enterprise multi-tenant + MSSP sur Tenants indépendants ; Customer externe ; Authorized Tenant Set read-only ; action tenant-local |
| [ADR-0009](adr/ADR-0009-slo-health-resilience-source-ownership-and-runtime-boundary.md) | SLO / Health / Resilience source, ownership et runtime boundary | REQ-PROD-003,005,006,008,009,010,012,013,019,021; REQ-SEC-001,002 | validated | SLO source-attributed non canonique ; Health projection ; aucun calculateur/failover/recovery générique ; visibilité MSSP read-only |

## Phase 6 SLO / Health / Resilience architecture — 2026-08-16

L'architecture SLO / Health / Resilience est approuvée explicitement par le project owner et enregistrée dans `ADR-0009`.

**Approval reference:** Explicit project-owner approval in this conversation.  
**Approved decision-set checksum:** `adb8312c2eb5cb65062177c72ee3b23cbe9f3c165593dab514a5aa53d6ad674a`

Décisions enregistrées :
- SLO utilise un modèle hybride source-attributed sans target store CMDR générique ;
- SLO reste une projection non canonique ; aucun objet SLO/Health Observation/Threshold/Service Level/Availability/Reliability/Resilience/RTO/RPO n'est créé ;
- Platform Architecture possède l'enveloppe contractuelle neutre ; Shared conserve les mécanismes métriques génériques ; les source/runtime owners conservent acquisition et calcul faisant autorité ;
- Platform Health reste une couche Settings de projection/presentation déterministe, pas un runtime de probe/monitoring/calcul ;
- la portée initiale Resilience se limite aux dégradations sourcées, Offline/Retry existant et contraintes de fallback provider comme configuration ;
- aucune configuration SLO générique, Permission ID, exécution failover/recovery générique ou capacité fonctionnelle n'est introduite ;
- visibilité initiale Phase-6 MVP : Tenant-local + MSSP read-only dans l'Authorized Tenant Set ; publication externe/customer reste bloquée par OPEN-019 ;
- Search, Report et Export restent single-selected-Tenant/single-Tenant et Shared-owned.

`OPEN-006` reste resolved. `OPEN-008`, `OPEN-013`, `OPEN-015` et `OPEN-019` restent open. L'enregistrement crée **0 Capability ID, 0 canonical object, 0 Permission ID et 0 Screen ID**. `CAP-SET-014+` reste unallocated/unreserved.

## OPEN-006 resolution — 2026-08-14

`OPEN-006 — Customers and Delivery deployment model` est **resolved** par l'approbation explicite du project owner Hillel Tobiana et appliqué par `ADR-0008`.

Décisions enregistrées :
- Internal, Enterprise multi-tenant et MSP/MSSP sont supportés ; MSSP est deployment-dependent ;
- Customer reste une projection externe et n'est ni objet canonique ni alias Tenant ;
- MSSP opère sur des Tenants indépendants sans hiérarchie ni Portfolio canonique ;
- Security résout un Authorized Tenant Set non canonique ; agrégation read-only et context switching seulement ;
- Search, Report et Export restent single-selected-Tenant initialement ;
- Response exige sélection Tenant, réévaluation Security puis Govern Decision Authority ;
- `CAP-CMD-401` conserve ID/owner et peut devenir `defined` après enregistrement canonique.

Aucun nouveau Capability ID, `CAP-SET-014`, `CAP-CMD-402`, objet canonique, Permission ID ou Screen ID n'est créé par la décision.

## Open decision index update
- `OPEN-008 — platform/source availability and support` remains open.
- `OPEN-013 — default governance/authority for reversible class-2 mutations` remains open.
- `OPEN-015 — Tool Call/Automation Run and Response Run/cross-product provenance bridge` remains open.
- `OPEN-017 — Detection runtime, target language and portability strategy` remains open and Detection-only.
- `OPEN-018 — Threat intelligence ontology, interoperability and exchange strategy` remains open; no ontology, standard, protocol, provider, exchange representation or implementation is selected.
- `OPEN-019 — Intelligence dissemination, releasability, sharing and consumer access policy` remains open.
- OPEN-019 selects no final policy, audience model, client-sharing model, external destination or publication authority.
- Open decisions: **17**.
- Resolved decisions include `OPEN-006` and the historically resolved `OPEN-009`.

Aucune autre OPEN n'est approuvée ou fermée par ADR-0009.
