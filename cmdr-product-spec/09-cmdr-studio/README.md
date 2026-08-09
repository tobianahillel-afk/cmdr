---
id: 09-cmdr-studio-readme
domain: 09-cmdr-studio
status: draft
owner: CMDR Studio Product Lead
updated: 2026-08-10
source-of-truth: canonical
requirements:
  - REQ-PROD-016
  - REQ-AI-001
  - REQ-AI-002
  - REQ-OBJ-009
---
# CMDR Studio

## Mission
Concevoir, versionner, évaluer, déployer et superviser les automatisations déterministes et agentiques.

## Possède
- Skill, Tool, Tool Call, Automation Agent, Agent Team, Workflow, Human Gate et Automation Run ;
- Library, Builder, Assurance et Control Room.

## Consomme
Objets et contextes des produits sous permission, policies Govern, providers et secrets administrés par Settings.

## Exclusions
- ne possède pas Incident, Case, Finding, Decision ou Response Run ;
- ne remplace pas les produits opérationnels ;
- n'est pas l'interface obligatoire.

## Transitions principales
Trigger produit → Automation Run → Tool Calls → output vers produit source ; action risquée → Action Request/Govern. STD-2 documente uniquement la définition du Workflow avant ce runtime futur.

## Place de l'IA
Studio possède les capacités agentiques, leur assurance, leurs versions, coûts et traces. Les workflows opérationnels restent utilisables sans IA.

## Delivery classification
Les capabilities Studio sont `planned` tant que leur implémentation et leurs contrats techniques ne sont pas prouvés. Tool Call et Automation Run restent soumis aux objets/contrats finaux futurs.

## Sources
- `../01-product-vision/product-boundaries.md`
- `../00-governance/ownership-register.md`
- `information-architecture.md`
- `product-definition.md`

## STD-1 — Studio Foundations — Tools, Skills, Library and Ownership Contracts
Parent: **Delivery Roadmap Phase 5 — Studio and Endpoint** (`roadmap-phase-5-studio-and-endpoint`). STD-1 est un execution lot, pas une Roadmap Phase et ne crée pas Phase 5A/5B.

- `CAP-STD-001..016`;
- 16 capabilities / 432 sections / 96 mandatory tables;
- all `draft / defined / planned`;
- Library, asset metadata/ownership, Tool/Tool Call, Skills, Settings dependency-reference boundaries and cross-product provenance;
- post-publication verification: PASS 190/190;
- no API/protocol/code/runtime implementation.

## STD-2 — Workflow Builder & Orchestration
STD-2 est un execution lot du même parent roadmap et ne crée aucune Phase 5A/5B.

- `CAP-STD-017..033`;
- 17 capabilities / 459 sections / 102 mandatory tables;
- all `draft / defined / planned`;
- Workflow definition, Builder Session, I/O/variables, graph, Tool/Skill steps, deterministic branches, mappings, subworkflows, parallelism, errors, retry/idempotency, partial/compensation, Human Gate boundary, readiness, versioning, pre-publish lifecycle and provenance;
- no Automation Run lifecycle, runtime scheduler, Control Room runtime, detailed publishing/deployment, Endpoint capability, API/protocol/code or final RBAC/ABAC;
- build-time status: PENDING POST-PUBLICATION VERIFICATION.

### Future execution lots
- STD-3 — Agents, Human Gates & runtime control — NOT STARTED.
- STD-4 — Assurance & lifecycle — NOT STARTED.
- Endpoint capability specification — NOT STARTED.

## Core non-equivalence
Tool != Tool Call/Skill/Workflow/Agent/Endpoint primitive/Govern Playbook; Tool Call != Automation Run/Response Run/Job/Result; Tool output != Evidence/Finding/Result automatically; Skill composition != Workflow orchestration; Workflow Definition != Workflow Version; Builder Session != Workflow; Tool step != Tool Call; branch != Govern Decision; condition != Govern Policy; Human Gate != Approval/Decision; retry != authorization renewal; compensation != Govern rollback; validation != execution; approved-for-publishing-candidate != deployed; Automation Run != Response Run.

## Critère d'acceptation
Un module de ce produit ne peut revendiquer un objet, une authority ou une capability exclue sans mise à jour des frontières, du registre de propriété et d'une ADR lorsqu'elle est transversale.
