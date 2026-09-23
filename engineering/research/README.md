# Engineering research evidence

This directory is the machine-readable evidence layer behind Class B/C engineering decisions.

It is deliberately separate from the Product Spec. Research evidence can support an engineering choice, but it cannot silently resolve a product-owned OPEN decision or weaken a product security invariant.

## Packet lifecycle

- `collecting`: evidence is still being gathered; the packet may be incomplete.
- `saturated`: the explicit research stop condition is satisfied.
- `superseded`: historical evidence replaced by another packet for the same decision.

A saturated packet does **not** need an arbitrary minimum number of sources. Instead it must prove that every declared major solution family is covered by sourced claims, critical limitations are known, no critical claim remains unresolved, and additional source search no longer materially changes the candidate set.

## Source provenance

Every source carries a kind, title, publisher, locator, publication/update date and access date. HTTPS, DOI and Git locators are accepted. Critical claims require at least one authoritative or primary source class and at least one explicit limitation.

Community discussions may provide discovery/context, but cannot be the sole authoritative basis for a critical claim.

## Bounded context

The policy caps packet/family/source/claim/limitation counts to keep agent resume context bounded. These are maximums for context safety, not source-count targets.

`cmdr-dev research-context --packet-id RES-PKT-....` emits the validated bounded packet for an agent.

## Decision binding

Once a Class B/C decision is accepted, every `evidence_ref` must resolve to a saturated packet owned by that same decision. Until then a packet may remain collecting or saturated while the decision is still researching.


## Freshness and reuse

`freshness-registry.json` binds accepted B/C decisions to cryptographic snapshots of their decision, research packets, validation evidence, Product Spec baseline and optional local assumption/threat/benchmark inputs. `decision-cache` derives the reusable index dynamically; it never trusts a hand-maintained cache entry.
