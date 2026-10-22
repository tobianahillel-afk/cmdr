---
id: phase-numbering-and-namespace-convention
domain: 18-roadmap-and-releases
status: draft
owner: Product Operations Lead
updated: 2026-08-09
source-of-truth: canonical
---
# Phase Numbering and Namespace Convention

## 1. Purpose

CMDR uses two historical phase-numbering schemes that serve different purposes. This convention makes both schemes explicit without renaming, recycling or re-parenting any existing phase identifier.

The shared numeral `4` is valid because the schemes live in separate namespaces. Numeric proximity across namespaces has no semantic, dependency or parent-child meaning.

## 2. Namespace A — Capability Specification Phase

**Name:** Capability Specification Phase.

**Purpose:** track the transformation of the product reference set into detailed capabilities, templates, ownership, traceability and quality gates.

Existing canonical execution labels include:

- `Capability Specification Phase 4A — Command`;
- `Capability Specification Phase 4B — Investigate`;
- `Capability Specification Phase 4` for the global capability-specification maturity axis.

Historical 4A/4B identifiers, reports, SHAs, Capability IDs and PASS verdicts are preserved. They are not renumbered to match the Delivery Roadmap.

## 3. Namespace B — Delivery Roadmap Phase

**Name:** Delivery Roadmap Phase.

**Purpose:** represent the historical product-delivery sequence.

Canonical sequence:

1. `Delivery Roadmap Phase 1 — Foundation` — [`phase-1-foundation.md`](phase-1-foundation.md);
2. `Delivery Roadmap Phase 2 — Command` — [`phase-2-command.md`](phase-2-command.md);
3. `Delivery Roadmap Phase 3 — Investigate` — [`phase-3-investigate.md`](phase-3-investigate.md);
4. `Delivery Roadmap Phase 4 — Govern` — [`phase-4-govern.md`](phase-4-govern.md);
5. `Delivery Roadmap Phase 5 — Studio and Endpoint` — [`phase-5-studio-and-endpoint.md`](phase-5-studio-and-endpoint.md);
6. `Delivery Roadmap Phase 6 — Platform Scale` — [`phase-6-platform-scale.md`](phase-6-platform-scale.md).

The existing file `phase-4-govern.md`, identifier `roadmap-phase-4-govern` and title `Phase 4 Govern` remain canonical.

## 4. Authority

The Capability Specification namespace is authoritative only for capability-specification execution status, capability-family sequencing, template conformance and related quality gates.

The Delivery Roadmap namespace is authoritative only for the historical product-delivery sequence and its phase plans.

Neither namespace silently renumbers or overrides the other.

## 5. Explicit non-equivalence

`Capability Specification Phase 4B — Investigate` is **not** the same phase as `Delivery Roadmap Phase 4 — Govern`.

The shared number `4` is a namespace overlap, not a parent-child relation. In particular:

- `Capability Specification Phase 4A — Command` does not imply `Delivery Roadmap Phase 4A`;
- `Capability Specification Phase 4B — Investigate` does not imply `Delivery Roadmap Phase 4B`;
- `Delivery Roadmap Phase 4 — Govern` is not a child of `Capability Specification Phase 4`;
- there is no inferred `Capability Specification Phase 4C — Govern`.

## 6. Qualification rule

New documentation MUST use a qualified form whenever the word `Phase` plus a number could be interpreted in more than one namespace.

Allowed forms:

- `Capability Specification Phase 4`;
- `Capability Specification Phase 4A — Command`;
- `Capability Specification Phase 4B — Investigate`;
- `Delivery Roadmap Phase 4 — Govern`;
- `Roadmap Phase 4 — Govern` when the Delivery Roadmap context is already explicit.

A new ambiguous use of bare `Phase 4` is prohibited.

Historical, locally unambiguous occurrences do not require mechanical rewriting.

## 7. Non-competition rule

No phase may be created or renamed solely to make numbering look sequential across namespaces.

Specifically, `Phase 4 Govern` MUST NOT be renamed to `Phase 4C Govern` solely to align it with Capability Specification Phase 4A/4B.

No concurrent `Phase 4C Govern`, competing Govern roadmap identifier or duplicate Govern phase file may be introduced while `roadmap-phase-4-govern` remains canonical.

## 8. Migration rule

Migration is documentary and incremental:

1. preserve historical filenames, IDs, SHAs, reports and verdicts;
2. qualify future ambiguous references;
3. update active indexes, status documents and sequencing guidance;
4. add links to this convention where future phase interpretation could otherwise diverge;
5. never rewrite a historical report merely to harmonize numbering.

## 9. Dependency rule

Functional or documentary dependencies MUST NOT be inferred from numeric adjacency between different phase namespaces.

For example, `Capability Specification Phase 4B — Investigate` and `Delivery Roadmap Phase 4 — Govern` have no implicit numeric dependency. Their real functional relationship is defined by product ownership and transitions:

`Investigate outputs → Action Request → Govern decision/authority → Response Run → Result`.

Dependency evidence comes from canonical ownership, object, transition and dependency sources, not from phase-number proximity.

## 10. Govern relation

Govern remains `Delivery Roadmap Phase 4 — Govern` with canonical identifier `roadmap-phase-4-govern`.

Its capability specification is **NOT STARTED** at the time of this reconciliation. No Govern capability, CAP-GOV identifier, new object, atomic permission, screen, API, protocol or implementation is created by this convention.

The next separate capability-specification execution may target Govern, but it must start from the then-current canonical branch SHA and must not create `Phase 4C Govern` merely for numbering symmetry.

## 11. Studio and Endpoint relation

`Delivery Roadmap Phase 5 — Studio and Endpoint` remains the historical delivery phase after Govern. Its number does not derive from Capability Specification Phase 4A/4B and does not prescribe a future capability-specification identifier.

No Studio or Endpoint capability-specification phase is started by this reconciliation.

## 12. Platform Scale relation

`Delivery Roadmap Phase 6 — Platform Scale` remains the historical delivery phase after Studio and Endpoint. Its number does not prescribe a future capability-specification identifier.

No Platform Scale capability-specification phase is started by this reconciliation.

## Acceptance

A reader must be able to determine from any new ambiguous phase reference which namespace is authoritative. Existing identifiers remain stable, `roadmap-phase-4-govern` remains unique, and no numeric relation is inferred across namespaces without an explicit functional dependency.
