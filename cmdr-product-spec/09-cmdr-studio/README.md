---
id: 09-cmdr-studio-readme
domain: 09-cmdr-studio
status: draft
owner: CMDR Studio Product Lead
updated: 2026-08-09
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
Trigger produit → Automation Run → Tool Calls → output vers produit source ; action risquée → Action Request/Govern.

## Place de l'IA
Studio possède les capacités agentiques, leur assurance, leurs versions, coûts et traces. Les workflows opérationnels restent utilisables sans IA.

## Delivery classification
Les capabilities Studio sont des cibles produit `planned` tant que leur implémentation et leurs contrats détaillés ne sont pas prouvés. Tool Call et Automation Run restent soumis aux objets/contrats finaux futurs.

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
- no new Screen ID or detailed rewrite;
- permission namespace ambiguity preserved, not normalized;
- no API/protocol/code/runtime implementation.

### Future execution lots
- STD-2 — Workflow Builder & orchestration — NOT STARTED.
- STD-3 — Agents, Human Gates & runtime control — NOT STARTED.
- STD-4 — Assurance & lifecycle — NOT STARTED.
- Endpoint capability specification — NOT STARTED.

## Core non-equivalence
Tool != Tool Call/Skill/Workflow/Agent/Endpoint primitive/Govern Playbook; Tool Call != Automation Run/Response Run/Job/Result; Tool output != Evidence/Finding/Result automatically; Skill composition != Workflow orchestration; Human Gate != Approval; Workflow != Playbook; Automation Run != Response Run.

## Critère d'acceptation
Un module de ce produit ne peut revendiquer un objet ou une capability exclue sans mise à jour des frontières, du registre de propriété et d'une ADR lorsqu'elle est transversale.

---

## STD-2 — Workflow Builder & Orchestration — current canonical addendum

The STD-1 section above is preserved as the exact pre-STD-2 historical snapshot. Its former “STD-2 NOT STARTED” line is historical evidence only; current status is defined below.

- execution lot under the same parent `roadmap-phase-5-studio-and-endpoint`;
- `CAP-STD-017..033`;
- 17 capabilities / 459 sections / 102 mandatory tables;
- all `draft / defined / planned`;
- Workflow definition, Builder Session, I/O/variables, graph, Tool/Skill steps, deterministic conditions/branches, mappings, subworkflows, ordering/parallelism, error paths, retry/idempotency, partial success/compensation, Human Gate/Govern boundary, no-effect readiness, Workflow Version, pre-publish lifecycle and cross-product provenance;
- no Automation Run lifecycle, scheduler, Control Room runtime, publishing/deployment engine, Endpoint capability, API/protocol/code, final language, final JSON Schema or final RBAC/ABAC;
- post-publication verdict: **PASS AFTER POST-PUBLICATION VERIFICATION — 200/200**.

### Current future execution lots
- STD-3 — Agents, Human Gates & Runtime Control — NOT STARTED.
- STD-4 — Assurance & Lifecycle — NOT STARTED.
- Endpoint capability specification — NOT STARTED.

### STD-2 additional non-equivalence
Workflow Definition != Workflow Version; Builder Session != Workflow; draft graph != saved Workflow Version; Tool step != Tool Call; Skill step != Skill; branch != Govern Decision; condition != Govern Policy; Human Gate completion != authorization; variable definition != variable value; mapping != source mutation; retry != authorization renewal; idempotency != exactly-once; compensation != Govern rollback; partial success != success; validation != execution; compatible != authorized; approved-for-publishing-candidate != deployed.

---

## STD-3 — Agents, Human Gates & Runtime Control — current addendum

The STD-1 and STD-2 sections above are preserved as exact historical snapshots. Their earlier STD-3 `NOT STARTED` statements are historical evidence only.

- parent: `roadmap-phase-5-studio-and-endpoint`;
- `CAP-STD-034..051`;
- **18 capabilities / 486 sections / 108 mandatory tables**;
- all `draft / defined / planned`;
- Automation Agent/Agent Team definitions, bounded objectives/access/planning/oversight;
- Human Gate request/runtime lifecycle with `accepted-for-workflow` distinct from Govern Approval/Decision;
- Automation Run creation/lifecycle/steps/attempts, queue/scheduling, controls, failures/retries/partial completion/context;
- Control Room, Studio Runtime Outcome, consumer handoff and runtime provenance;
- no new Screen ID or detailed rewrite;
- no agent framework, model/provider, scheduler/runtime selection, API/protocol/code, final JSON Schema/RBAC, publishing/deployment or Endpoint capability;
- fifth functional/build SHA: `c658168c9de6bd803941116989bc3aaedf154260`;
- post-publication verification is being recorded separately after a documentary history-preservation correction.

### Current future boundary
- STD-4 — Assurance & Lifecycle — NOT STARTED.
- Endpoint capability specification — NOT STARTED.
- Studio capability specification remains PARTIAL.

### STD-3 additional non-equivalence
Automation Agent != Workflow/Skill/Tool/Automation Run/human user; Agent Team != human Team/Workflow; role/objective/access/proposal/plan != authorization/action; Human Gate != Approval/Decision; Automation Run != Workflow/Version/Tool Call/Job/Response Run; created/queued/start-requested != running; paused != stopped; cancellation != rollback; retry attempt != Run; retry != new Run; idempotency != exactly-once; Tool Call/step success != Run success; partial completion != success; Studio outcome != Govern Result/Evidence/Finding; Control Room != Command Work Queue/Govern Runs & Rollback.

## STD-3 final verification addendum
STD-3 final documentary verdict is **PASS AFTER POST-PUBLICATION VERIFICATION — 210/210 gates PASS, 0 PENDING, 0 FAIL**. Canonical evidence: `../16-quality-and-validation/reports/studio-std3-agents-human-gates-runtime-control-post-publication-verification.md`. Studio capability specification remains PARTIAL because STD-4 is NOT STARTED; Endpoint remains NOT STARTED.