# Engineering decision registry

This directory is the engineering source of truth for **recorded technical decisions**. It does not replace product decisions.

## Classes

The machine-readable class policy is `decision-policy.json`.

- **A — local engineering**: equivalent local naming/refactor/internal representation. Recording is optional; external research is not mandatory.
- **B — significant technical**: data structures, cache/serialization strategy, internal concurrency, meaningful dependency/build choices. Recording and alternatives are mandatory. Representative benchmarks are mandatory when performance-sensitive.
- **C — critical**: security/architecture/algorithm/durability/high-impact performance choices. Current research, alternatives and adversarial falsification are mandatory. Performance-sensitive Class C choices also require representative CMDR benchmarks.

## Product boundary

A recorded decision declares exactly one boundary:

- `engineering-only`
- `product-open-decision`
- `product-behavior`
- `product-security-invariant`

Only `engineering-only` decisions can be accepted by the engineering engine. A choice that crosses the product boundary must become `blocked-product`, `rejected`, or historical `superseded`; the agent must not silently invent a product default.

## Acceptance

The registry deliberately starts empty. Earlier engineering work is not retroactively labelled as an accepted research decision without explicit evidence.

For a decision to be `accepted`:

- owner, decision key, question, rationale, constraints and affected work units are valid;
- Class B/C entries contain research evidence references;
- performance-sensitive Class B/C entries contain benchmark evidence;
- Class C entries contain adversarial-review evidence;
- the product boundary is `engineering-only`;
- no other live decision exists for the same decision key.

Evidence references are identifiers in this lot. E5-RES-001B/C make their backing artifacts and stronger semantics mandatory.

## Supersession

Historical entries may be `superseded`. They must point to another registered decision with the same `decision_key`. At most one live entry may exist for a decision key.

This makes a fresh agent able to distinguish active, blocked, rejected and historical choices without relying on chat history.


## Validation evidence

`decision-validation.json` holds the engineering proof required before Class B/C acceptance: candidate trade-offs, representative CMDR benchmarks, prototypes when external evidence is insufficient, and adversarial falsification for Class C. These objects are validated by `CHK-DECISION-GATES`; opaque references alone are no longer sufficient once a decision is accepted.
