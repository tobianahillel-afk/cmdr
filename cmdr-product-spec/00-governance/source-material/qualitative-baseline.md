---
id: qualitative-baseline
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-07
source-of-truth: canonical
---
# Qualitative Baseline — Investigate programme through Phase 4B.4B pre-publication closure

## Inherited Phase 4B.3B.1 baseline
Remote start for 4B.3B.1 was `c84ea542a0831146b075eb0868dc21326b3b9ad8`: Phase 4B.3A PASS, 196 capabilities, 169 Investigate capabilities, 35 CAP-INV-4xx, no CAP-INV-5xx, 16 open decisions and no canonical Investigate Threat Intelligence module.

### 4B.3B.1 source audit preserved
- Full PR path manifest reviewed: **1259** paths at the starting PR state.
- Direct canonical sources re-read: **at least 55**, spanning governance, registers, Command objects, Signals/Hunt, Cases/Evidence/Entity, Analysis Workbench handoffs, complete Detection Engineering closure and CAP-INV-435, historical Threat Intelligence enrichment, Settings source/secret/health ownership, Studio Tool/Run ownership, Govern authority, Shared Entity/Graph/Linking and required screens/shells.
- Historical Threat Intelligence needs were distributed by owner; no competing Investigate Threat Intelligence functional module or autonomous Intelligence Library screen existed.
- Competing Investigate functional sources deprecated: **0**. Shared enrichment and all owner contracts remained active.
- `OPEN-018` was required because no existing decision governed Intelligence ontology, interoperability and future exchange.

### 4B.3B.1 measures preserved
| Measure | Before 4B.3B.1 | After 4B.3B.1 |
|---|---:|---:|
| Threat Intelligence targeted functional files | 0 | 32 |
| Active Threat Intelligence module documents | 0 | 32 |
| Deprecated functional pointers | 0 | 0 |
| Generic active documents / placeholders | 0 / 0 | 0 / 0 |
| CAP-INV-5xx | 0 | 18 |
| Investigate capabilities | 169 | 187 |
| Registered capabilities | 196 | 214 |
| Defined / proposed / planned | 194 / 2 / 196 | 212 / 2 / 214 |
| Sections expected / present | 0 | 486 / 486 |
| Mandatory tables expected / present | 0 | 108 / 108 |
| Empty / prose-only / generic mandatory tables | 0 | 0 |
| Required surfaces read / specs modified / detailed rewrites / new IDs | 18 / 0 / 0 / 0 | 18 / 0 / 0 / 0 |
| APIs / protocols / standards / providers / commands / code | 0 | 0 |
| Requirement IDs | 122 | 122 |
| Open decisions | 16 | 17 |
| Threat Intelligence capabilities / sections / tables | 0 / 0 / 0 | 18 / 486 / 108 |
| Investigate capabilities / sections / tables | 169 / 4563 / 1014 | 187 / 5049 / 1122 |
| Command + Investigate capabilities / sections / tables | 196 / 5292 / 1176 | 214 / 5778 / 1284 |

## Phase 4B.3B.2 starting state
Remote start: `b24cb7a373562fb96d6087b5c3a9f7fee91deb94`. It contained Phase 4B.3B.1 PASS, 214 capabilities, 187 Investigate capabilities, 18 CAP-INV-5xx, 17 open decisions and no CAP-INV-519..537.

## Phase 4B.3B.2 measures
| Measure | Before 4B.3B.2 | After 4B.3B.2 |
|---|---:|---:|
| CAP-INV-519..537 | 0 | 19 |
| Threat Intelligence capabilities | 18 | 37 |
| Threat Intelligence sections / tables | 486 / 108 | 999 / 222 |
| Investigate capabilities | 187 | 206 |
| Investigate sections / tables | 5049 / 1122 | 5562 / 1236 |
| Global capabilities | 214 | 233 |
| Defined / proposed / planned | 212 / 2 / 214 | 231 / 2 / 233 |
| Command + Investigate sections / tables | 5778 / 1284 | 6291 / 1398 |
| New capability documents | 0 | 19 |
| New analysis/products supporting documents | 0 | 10 |
| Capability sections expected / present | 0 | 513 / 513 |
| Mandatory tables expected / present | 0 | 114 / 114 |
| Empty / prose-only / generic mandatory tables | 0 | 0 |
| Duplicate IDs / concurrent owners | 0 / 0 | 0 / 0 |
| Detailed screen rewrites / new Screen IDs | 0 / 0 | 0 / 0 |
| Canonical object files / atomic permissions | 0 / 0 | 0 / 0 |
| APIs/protocols/standards/providers/commands/code | 0 | 0 |
| Active watchlists/deployed Indicators/rules/blocks/responses | 0 | 0 |
| Actual external sharing/client/public publication | 0 | 0 |
| Cloud/Mobile content | 0 | 0 |
| Requirement IDs | 122 | 122 |
| Open decisions | 17 | 18 |

## Phase 4B.4A starting state
Remote start: `6cc7bf2426c12a6f2f983c62581f9d44e9829990`. The four published functional Cloud commits advanced the canonical branch to `59f69203574a6b1e3dc8c631b3a8a2e4bf5afe42`, adding `CAP-INV-601..618`, the Cloud module documents and the canonical `4B.4A` roadmap while leaving the closing registers and quality evidence incomplete.

### 4B.4A corrective source audit
- The canonical branch, PR #2, repository visibility, README files, five functional commits, 18 capability files and temporary closure branch were verified directly.
- Governance and source-material documents were re-read: master brief, product boundaries, canonical object/ownership decisions, AI constraints, native strategy, non-goals, source-of-truth policy, documentation rules, ownership register, object register, screen register, Requirements Matrix, qualitative baseline and OPEN decisions.
- Investigate family sources were re-read across Signals/Hunt, Cases/Evidence, Collection/Live Response, Analysis Workbench, Dynamic Sandbox, Reverse/Debugger, Memory, Disk, Network, Detection Engineering and Threat Intelligence.
- All 18 `CAP-INV-601..618` contracts were read directly and checked for front matter, sections 1–27, six mandatory tables, permissions, non-AI paths and non-implementation boundaries.
- Platform Settings sources were re-read for sources/parsers, secrets/connections, tenants/environments, health and users/roles.
- Command, Studio, Govern, Shared, permission, tenant-isolation, secret-management and Evidence-trust sources were re-read.
- Canonical object sources were re-read for Case, Artifact, Evidence, Finding, Hypothesis, Principal, Role, Secret Reference, Data Source, Integration, Environment, Tenant and Timeline Entry.
- Existing surfaces were re-read for Event Search, Case Workspace, Evidence Board, Hypotheses/Findings and Settings source, secret and tenant administration.
- No contradiction requiring a change to `CAP-INV-601..618` was found. No historical source was deprecated or condensed.

## Phase 4B.4A measures
| Measure | Before 4B.4A | After verified 4B.4A |
|---|---:|---:|
| CAP-INV-601..618 | 0 | 18 |
| Cloud Analysis capabilities | 0 | 18 |
| Cloud Analysis sections / tables | 0 / 0 | 486 / 108 |
| Investigate capabilities | 206 | 224 |
| Investigate sections / tables | 5562 / 1236 | 6048 / 1344 |
| Global capabilities | 233 | 251 |
| Defined / proposed / planned | 231 / 2 / 233 | 249 / 2 / 251 |
| Command + Investigate sections / tables | 6291 / 1398 | 6777 / 1506 |
| New capability documents | 0 | 18 |
| New Cloud supporting and roadmap documents | 0 | 15 |
| Capability sections expected / present | 0 | 486 / 486 |
| Mandatory tables expected / present | 0 | 108 / 108 |
| Empty / prose-only / generic mandatory tables | 0 | 0 |
| Duplicate/recycled IDs / concurrent owners | 0 / 0 / 0 | 0 / 0 / 0 |
| Detailed screen rewrites / new Screen IDs | 0 / 0 | 0 / 0 |
| Canonical Cloud object files / atomic permissions | 0 / 0 | 0 / 0 |
| APIs/protocols/providers/connectors/commands/code | 0 | 0 |
| Cloud target mutations / secret use / deployed rules / responses | 0 | 0 |
| Mobile capabilities | 0 | 0 |
| Requirement IDs | 122 | 122 |
| Requirement states conform / partial / absent / contradictory | 99 / 20 / 3 / 0 | 99 / 20 / 3 / 0 |
| Open decisions | 18 | 18 |
| Quality gates | not applicable | 243 / 243 PASS |

## Verified 4B.4A maturity
Phase 4B.3 remains PASS. Phase 4B.4A is **PASS AFTER POST-PUBLICATION VERIFICATION**. Before Mobile work, Phase 4B.4, Phase 4B, Phase 4 and global maturity remained PARTIAL and Mobile Forensics was NOT STARTED.

## Phase 4B.4B starting state
Remote start: `ed874ea414fc57f24fa61f410f91b7345f4a868a`. At that exact SHA, Phase 4B.4A was verified PASS; there were 251 global capabilities, 224 Investigate capabilities, no `CAP-INV-7xx`, no canonical Mobile module, no Mobile roadmap file and no Mobile Screen ID. The four published Mobile functional commits advance the canonical branch to `78d49e3fb8895a4fb9f46bd1f4f7a28fcb4d8a52` before the closure commit.

### 4B.4B source and ownership audit
- Repository, PR #2, public visibility, canonical branch, `main`, root READMEs and start SHA were verified directly before modification.
- All active capability-register shards through `CAP-INV-618` were read; `CAP-INV-701..719` was confirmed free and unrecycled.
- Governance/source material read included master brief, capability inventory, product boundaries, native strategy, non-goals, AI constraints, source-of-truth/documentation/terminology rules, ownership/status, open decisions, registers, requirements and baseline.
- Investigate sources read included Collection/Live Response, Static/Dynamic/Reverse, Memory/Disk/Network, Detection Engineering, Threat Intelligence and Cloud Analysis.
- Collection acquisition/request/job/custody boundaries were read and retained; Mobile consumes results and never owns acquisition execution.
- Platform Settings, Endpoint Agent, Studio, Govern, Shared, permissions, privacy/minimization, secrets/key management and Evidence trust sources were read and retained.
- Existing mapped screens were read for Disk/Artifact, Case, Evidence, Event Search, Static, Reverse, Endpoint Fleet, Settings Sources, Studio Control Room and Govern Action Center; no active Mobile screen existed.
- No competing canonical Mobile functional source was found; deprecated Mobile sources: 0; active owner documents remain active.
- No Cloud capability, owner or functional scope was changed.

## Phase 4B.4B measures before final remote verification
| Measure | Before 4B.4B | After closure content, before remote verification |
|---|---:|---:|
| CAP-INV-701..719 | 0 | 19 |
| Mobile Forensics capabilities | 0 | 19 |
| Mobile sections / tables | 0 / 0 | 513 / 114 |
| Phase 4B.4 Cloud + Mobile capabilities | 18 | 37 |
| Phase 4B.4 sections / tables | 486 / 108 | 999 / 222 |
| Investigate capabilities | 224 | 243 |
| Investigate sections / tables | 6048 / 1344 | 6561 / 1458 |
| Global capabilities | 251 | 270 |
| Defined / proposed / planned | 249 / 2 / 251 | 268 / 2 / 270 |
| Command + Investigate sections / tables | 6777 / 1506 | 7290 / 1620 |
| New capability documents | 0 | 19 |
| New Mobile supporting + roadmap documents | 0 | 15 |
| Capability sections expected / present | 0 | 513 / 513 |
| Mandatory tables expected / present | 0 | 114 / 114 |
| Empty / prose-only / generic mandatory tables | 0 | 0 |
| Duplicate/recycled IDs / concurrent owners / active contradictions | 0 / 0 / 0 / 0 | 0 / 0 / 0 / 0 |
| Detailed screen rewrites / new Screen IDs | 0 / 0 | 0 / 0 |
| Canonical Mobile object files / atomic permissions | 0 / 0 | 0 / 0 |
| APIs/protocols/platforms/tools/connectors/commands/code | 0 | 0 |
| Unlock/bypass/root/jailbreak/device mutation/secret use | 0 | 0 |
| Requirement IDs | 122 | 122 |
| Requirement states conform / partial / absent / contradictory | 99 / 20 / 3 / 0 | 99 / 20 / 3 / 0 |
| Open decisions | 18 | 18 |

## Maturity before final Mobile publication verification
Phase 4B.3 and verified Cloud Phase 4B.4A remain PASS. Mobile functional and closure content is complete, but Phase 4B.4B remains **PENDING POST-PUBLICATION VERIFICATION** until the fifth functional commit is published and remote branch/PR/README/register/report checks pass. Therefore Phase 4B.4 and Phase 4B remain PARTIAL at this publication stage. Phase 4 and global maturity remain PARTIAL regardless because later objects, permissions, screens, technique and implementation remain future.
