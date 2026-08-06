---
id: investigate-threat-intelligence
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements:
  - REQ-PROD-014
  - REQ-PROD-019
  - REQ-PROD-020
  - REQ-INV-006
open_decisions:
  - OPEN-013
  - OPEN-014
  - OPEN-015
  - OPEN-018
---
# Threat Intelligence — Foundations and Knowledge Management

## Mission
Transformer des besoins et matériaux sourcés en connaissances **candidates**, versionnées, contradictoires lorsque nécessaire et entièrement traçables, sans imposer une ontologie externe, confirmer une attribution, déployer un Indicator, créer une watchlist ou publier un Intelligence Report.

## Capability range
`CAP-INV-501..518` — dix-huit capabilities `defined` et `planned`, 486 sections numérotées et 108 tableaux obligatoires.

## Functional chain
Intake → Intelligence Requirement → Knowledge Project → Source Catalog → Reliability/Credibility → Material Intake/Normalization → Observable/Indicator and Threat Entity candidates → Malware/Tool/Infrastructure/Campaign knowledge → TTP/Sightings/Relationships → Confidence/Contradictions → Deduplication/Versioning → Expiration/Revocation → Provenance/Analysis Handoff.

## Ownership
Investigate possède le contexte Threat Intelligence, les Requirements, Projects, candidates, assessments et handoffs fonctionnels. Shared conserve Entity, Graph, Search, Timeline, Linking, Versioning et les mécanismes transverses. Command conserve Detection, Signal, Alert et Incident. Detection Engineering conserve Detection Content. Settings administre les sources, providers, connectors, secrets, accès, rétention et health. Studio conserve Tools et Automation Runs. Govern conserve toute autorité de partage externe et d’action de production.

## Phase boundary
4B.3B.1 couvre fondations et knowledge management. 4B.3B.2 reste non commencée : pas de fusion analytique avancée, attribution, Report, dissemination, watchlist, monitoring continu, external exchange, operationalization ou clôture de Phase 4B.
