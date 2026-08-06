---
id: qualitative-baseline
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-06
source-of-truth: canonical
---
# Qualitative Baseline — Threat Intelligence programme through Phase 4B.3B.2

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

## Maturity
Phase 4B.3B.1, 4B.3B.2, 4B.3B and 4B.3 pass their documentary functional gates. Phase 4B remains PARTIAL because OPEN-011 and OPEN-012 remain open and are not explicitly deferred outside its approved scope. Phase 4 and global maturity remain PARTIAL.
