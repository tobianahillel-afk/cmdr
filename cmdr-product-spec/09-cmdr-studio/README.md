---
id: 09-cmdr-studio-readme
domain: 09-cmdr-studio
status: draft
owner: CMDR Studio Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-PROD-016, REQ-AI-001, REQ-AI-002, REQ-OBJ-009]
---
# CMDR Studio

## Mission
Concevoir, versionner, évaluer, déployer et superviser les automatisations déterministes et agentiques.

## Possède
Skill, Tool, Tool Call, Automation Agent, Agent Team, Workflow, Human Gate et Automation Run; Library, Builder, Assurance et Control Room.

## Consomme
Objets/contextes produits sous permission, policies Govern, providers/integrations/secrets administrés par Settings et mécanismes Shared.

## Exclusions
Studio ne possède pas Incident, Case, Finding, Decision, Response Run ou Result, ne remplace pas les produits opérationnels et n'est pas une interface obligatoire.

## Delivery Roadmap Phase 5
Parent: `roadmap-phase-5-studio-and-endpoint`. Studio et Endpoint restent deux capability domains séparés.

### STD-1 — Studio Foundations — Tools, Skills, Library and Ownership Contracts
**Functional content complete; publication verification is recorded in the STD-1 conformance report.**
- `CAP-STD-001..016`;
- 16 capabilities / 432 sections / 96 mandatory tables;
- all `draft / defined / planned`;
- no new Screen ID or detailed screen rewrite;
- permission namespace ambiguity preserved, not normalized;
- OPEN-007/013/015 remain open where consumed;
- no API/protocol/code/runtime implementation.

### Future execution lots
- STD-2 — Workflow Builder & orchestration — NOT STARTED.
- STD-3 — Agents, Human Gates & runtime control — NOT STARTED.
- STD-4 — Assurance & lifecycle — NOT STARTED.
- Endpoint capability specification — NOT STARTED.

## Core boundaries
Tool != Tool Call; Skill != Tool/Workflow/Agent; Workflow != Govern Playbook; Human Gate != Govern Approval; Automation Run != Govern Response Run; Tool output/Tool Call technical outcome != Govern Result or Investigate Evidence/Finding automatically. Settings retains providers/integrations/credentials/secrets/tenant-environment administration; Shared retains generic Search/Trace/Activity/Jobs/Versioning/Reporting/Export/Notifications/Collaboration.

## Sources
- `../01-product-vision/product-boundaries.md`
- `../00-governance/ownership-register.md`
- `information-architecture.md`
- `product-definition.md`
- `capability-map.md`
- `object-consumption-map.md`
- `cross-product-links.md`
- `permissions.md`

## Acceptance
A Studio capability cannot claim excluded object ownership or runtime delivery without updating canonical boundaries/registers and the required ADR/evidence chain.
