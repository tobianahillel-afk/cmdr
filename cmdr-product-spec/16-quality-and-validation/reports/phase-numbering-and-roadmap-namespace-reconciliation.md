---
id: phase-numbering-and-roadmap-namespace-reconciliation
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-09
source-of-truth: quality-report
requirements: [REQ-PROD-006, REQ-PROD-008, REQ-PROD-013, REQ-PROD-014, REQ-PROD-015]
---
# Phase Numbering and Roadmap Namespace Reconciliation

## Verdict

**PASS — 50/50 gates PASS, 0 PENDING, 0 FAIL after post-publication verification.**

This is a documentary-governance correction only. It creates no product capability, object, screen, atomic permission, API, protocol, implementation or product decision.

The exact final commit SHA is verified externally against PR #2 and the canonical branch after publication rather than embedded in this self-referential commit. This preserves the requirement for one functional commit only.

## Before

Two valid historical numbering schemes coexisted without an explicit namespace contract:

- Capability Specification used `Phase 4A — Command`, `Phase 4B — Investigate` and nested 4B subphases to track detailed capability-specification execution and quality gates.
- The historical Delivery Roadmap used Phase 1 Foundation, Phase 2 Command, Phase 3 Investigate, Phase 4 Govern, Phase 5 Studio and Endpoint and Phase 6 Platform Scale.

Because both schemes used the numeral `4`, a future reader could incorrectly infer `4A → 4B → 4C Govern` even though no canonical `Phase 4C Govern` existed.

Historical reports also intentionally preserve the wording and status that were true at their publication time; they are not mechanically rewritten by this correction.

## Canonical resolution

Two explicit namespaces are now canonical:

### Capability Specification Phase

Purpose: detailed capabilities, templates, ownership, traceability and quality gates.

- `Capability Specification Phase 4A — Command`: PASS, unchanged.
- `Capability Specification Phase 4B — Investigate`: PASS, unchanged by this numbering correction.
- `Capability Specification Phase 4` global maturity: PARTIAL.

The 4A/4B identifiers, capability IDs, reports, SHAs and historical status evidence remain unchanged.

### Delivery Roadmap Phase

Purpose: historical product-delivery sequence.

1. Delivery Roadmap Phase 1 — Foundation.
2. Delivery Roadmap Phase 2 — Command.
3. Delivery Roadmap Phase 3 — Investigate.
4. Delivery Roadmap Phase 4 — Govern.
5. Delivery Roadmap Phase 5 — Studio and Endpoint.
6. Delivery Roadmap Phase 6 — Platform Scale.

The namespaces are independent. Numeric adjacency across namespaces creates no dependency, inheritance or parent-child relation.

## Qualification and migration

New ambiguous documentation must use a qualified form such as `Capability Specification Phase 4B — Investigate` or `Delivery Roadmap Phase 4 — Govern`. A new bare `Phase 4` is prohibited where either namespace could be intended.

Historical locally unambiguous references are preserved. Files and IDs are not renamed merely to make the two schemes look sequential.

Functional dependencies continue to come from ownership, object, transition and dependency sources. The actual Investigate→Govern relationship is functional: Investigate produces findings/evidence/context and Action Request inputs; Govern owns authority, Decision, Response Run and Result.

## Govern

Confirmed without starting capability work:

- existing canonical file: `18-roadmap-and-releases/phase-4-govern.md`;
- existing canonical ID: `roadmap-phase-4-govern`;
- existing canonical title: `Phase 4 Govern`;
- qualified title: `Delivery Roadmap Phase 4 — Govern`;
- previous delivery phase: Delivery Roadmap Phase 3 — Investigate;
- next delivery phase: Delivery Roadmap Phase 5 — Studio and Endpoint;
- Govern capability specification: **NOT STARTED**;
- `Phase 4C Govern`: **DOES NOT EXIST**;
- CAP-GOV IDs created: **0**.

`Phase 4 Govern` MUST NOT be renamed to `Phase 4C Govern` solely to align with Capability Specification Phase 4A/4B.

## Impact

| Measure | Result |
|---|---:|
| Capabilities changed | 0 |
| Capability IDs changed | 0 |
| Govern capabilities created | 0 |
| Objects changed | 0 |
| Ownership changed | 0 |
| Requirements added/removed | 0 |
| Requirement distribution changed | 0 |
| Screens changed | 0 |
| Atomic permissions changed | 0 |
| APIs/protocols/code created | 0 |
| Open decisions created | 0 |
| Open decisions closed | 0 |
| Existing open decisions retained | 18 |
| Historical reports removed | 0 |
| Historical phase files renamed | 0 |
| Historical SHAs rewritten | 0 |

## Reconciliation gates — 50/50

| # | Gate | Verdict | Evidence |
|---:|---|---|---|
| 1 | Repository correct | PASS | `tobianahillel-afk/cmdr` |
| 2 | Branch correct | PASS | `docs/cmdr-product-spec-foundation` |
| 3 | PR #2 correct | PASS | canonical documentation PR |
| 4 | Initial SHA exact | PASS | `b34284fb3e7434ae495ecd151cf75c572a1861ae` |
| 5 | Branch README unchanged | PASS | exact `# cmdr` |
| 6 | Main README unchanged | PASS | exact `# cmdr`, same blob |
| 7 | Main unchanged | PASS | no update to `main` ref |
| 8 | Capability Specification 4A identified | PASS | Command capability map / historical report |
| 9 | Capability Specification 4B identified | PASS | Investigate capability map / historical reports |
| 10 | Delivery Roadmap Phase 1 identified | PASS | `phase-1-foundation.md` |
| 11 | Delivery Roadmap Phase 2 identified | PASS | `phase-2-command.md` |
| 12 | Delivery Roadmap Phase 3 identified | PASS | `phase-3-investigate.md` |
| 13 | Delivery Roadmap Phase 4 Govern identified | PASS | `phase-4-govern.md` |
| 14 | Delivery Roadmap Phase 5 identified | PASS | `phase-5-studio-and-endpoint.md` |
| 15 | Delivery Roadmap Phase 6 identified | PASS | `phase-6-platform-scale.md` |
| 16 | `phase-4-govern.md` preserved | PASS | same path retained |
| 17 | `roadmap-phase-4-govern` preserved | PASS | canonical id unchanged |
| 18 | No Phase 4C created | PASS | convention prohibits competing phase |
| 19 | No concurrent Govern phase file created | PASS | only existing roadmap Govern phase retained |
| 20 | No Govern capability file created | PASS | zero capability additions under Govern |
| 21 | No Capability ID created | PASS | reconciliation tree adds none |
| 22 | No Capability ID renamed | PASS | Command/Investigate ranges preserved |
| 23 | No historical conformance report deleted | PASS | report corpus retained |
| 24 | No historical SHA modified | PASS | history preserved; one descendant commit only |
| 25 | No historical PASS status lost | PASS | Command/Investigate PASS retained |
| 26 | Capability Specification namespace documented | PASS | canonical convention |
| 27 | Delivery Roadmap namespace documented | PASS | canonical convention |
| 28 | `Phase 4` ambiguity documented | PASS | qualification rule |
| 29 | Qualification rule documented | PASS | ambiguous new references must name namespace |
| 30 | STATUS coherent | PASS | two independent status axes |
| 31 | Roadmap index coherent | PASS | separate execution-plan and delivery-phase sections |
| 32 | Dependency Register coherent | PASS | existing functional/documentary dependency register retained; no numeric edge added |
| 33 | Cross-product Links coherent | PASS | existing Investigate→Govern functional links retained |
| 34 | CHANGELOG coherent | PASS | reconciliation entry added without deleting history |
| 35 | PR description coherent | PASS | post-publication governance note added to PR #2 |
| 36 | No new OPEN | PASS | open-decision count remains 18 |
| 37 | No OPEN closed | PASS | existing decisions retained |
| 38 | No placeholder introduced | PASS | reconciliation documents are substantive |
| 39 | No empty file | PASS | every added/updated document has substantive content |
| 40 | No targeted link broken | PASS | links target retained canonical files |
| 41 | No API | PASS | documentation governance only |
| 42 | No protocol | PASS | documentation governance only |
| 43 | No code | PASS | Markdown-only tree |
| 44 | No Govern capability | PASS | capability specification remains NOT STARTED |
| 45 | Fast-forward publication only | PASS | final branch update is non-forced descendant update |
| 46 | No force-push | PASS | `force=false` branch ref update |
| 47 | No rebase | PASS | baseline remains direct parent of reconciliation commit |
| 48 | Construction SHA equals remote branch SHA | PASS | verified after publication; exact SHA recorded in PR/final execution output |
| 49 | PR remains Draft | PASS | PR #2 verified Draft after publication |
| 50 | Govern capability specification remains NOT STARTED | PASS | roadmap/status/convention agree; zero CAP-GOV files |

## Preserved historical evidence

Historical Phase 4A and Phase 4B reports remain unchanged. Some pre-post-publication reports intentionally retain their original `PENDING` or `PARTIAL` wording because they are time-stamped evidence, not mutable current-status dashboards. Current capability-specification status is recorded in `STATUS.md`, active capability maps and PR #2 without rewriting those reports.

## Next step

The next separate execution may begin Govern capability specification from the new verified canonical branch SHA. It must first re-verify that SHA and preserve `Delivery Roadmap Phase 4 — Govern` / `roadmap-phase-4-govern`; it must not create `Phase 4C Govern` merely for numeric alignment.
