# Architecture boundaries and path ownership

CMDR uses deny-by-default path ownership for implementation work.

The registry at `engineering/architecture/boundaries.json` is an engineering control-plane artifact. It does not define product behavior.

## Initial boundaries

- `product-spec`: owns `cmdr-product-spec/**` and is never mutable by implementation work;
- `engineering-control-plane`: owns engineering tooling, governance and root autonomous-agent rules;
- `execution-control-plane`: owns `work/**`;
- `ci-control-plane`: owns `.github/**`.

No product-runtime source root is registered yet. This is intentional. A runtime layout must be an explicit architecture decision before any implementation manifest can authorize it.

## Enforcement

`cmdr-dev architecture-audit` validates:

- strict registry decoding;
- unique boundary ownership;
- non-overlapping roots;
- acyclic boundary dependency declarations;
- immutable Product Spec ownership;
- every strict-v2 manifest allowed path belongs to exactly one mutable boundary;
- every strict-v2 manifest explicitly forbids `cmdr-product-spec/**`.

Unknown paths fail closed.
