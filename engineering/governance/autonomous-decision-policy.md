# Autonomous decision policy

The goal is high autonomy without silent product mutation.

## Class A — local engineering decision

Examples: local naming, small refactor, equivalent internal representation.

Action: decide, implement, test, record only if useful. No user interruption.

## Class B — significant technical decision

Examples: data structure, cache strategy, serialization strategy, internal concurrency design.

Action:
1. search existing accepted decisions;
2. define explicit CMDR constraints;
3. compare credible alternatives;
4. benchmark when performance-sensitive;
5. record an ADR/TDD;
6. continue autonomously.

## Class C — critical architecture/security/algorithm decision

Examples: storage/index architecture, cryptographic integration, authorization architecture, hostile parser design, durability semantics, high-scale concurrency.

Action:
1. current web/scientific research is mandatory;
2. identify solution families before products/implementations;
3. evaluate security, performance, failure modes, maturity and replaceability;
4. prototype finalists when evidence is insufficient;
5. benchmark representative CMDR workloads;
6. perform adversarial review intended to falsify the preferred choice;
7. record decision, evidence date and revisit triggers;
8. continue only when the decision is sufficiently supported.

## Product-decision boundary

If a choice would alter canonical product behavior, weaken a product security invariant, or resolve a product-owned OPEN decision without authority, mark only the affected work unit `BLOCKED_DECISION`. Do not fabricate a default. Continue with another independent READY unit when possible.

## Research stop condition

Research stops at evidence saturation: major credible solution families are covered, new sources no longer materially change the candidate set, critical limitations are known, and the selected candidate survives required CMDR-specific validation.
