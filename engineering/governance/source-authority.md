# Engineering source authority

## Purpose

Extend, without replacing, the canonical product source-of-truth policy in `cmdr-product-spec/00-governance/source-of-truth-policy.md`.

## Authority order for implementation work

1. Active canonical product sources under `cmdr-product-spec/**`, using their existing source-of-truth hierarchy.
2. Accepted engineering ADR / Technical Decision Dossier that does not contradict product truth.
3. Engineering architecture and security rules.
4. Canonical work-unit manifest for execution scope.
5. Generated indexes and context bundles.
6. Handoffs and summaries.

Git/GitHub reality determines whether claimed execution actually occurred.

## Conflict rule

- Product versus engineering convenience: product wins.
- Canonical source versus generated index: canonical source wins.
- Accepted decision versus stale summary: accepted decision wins.
- Recorded state versus Git/PR/check evidence: reality wins and state is repaired.
- Equal-authority contradiction: stop only the affected unit, record the contradiction, and select another independent READY unit when possible.

## Non-duplication rule

Engineering documents reference product behavior; they do not re-specify it. A material product change must go through the product governance path rather than being hidden inside an engineering ADR or implementation.
