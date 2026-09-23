# Decision freshness and reusable evidence cache

Accepted Class B/C decisions are reusable only while their evidence basis remains identical and their explicit freshness window and revisit signals remain clear.

## Snapshot

A freshness record binds the decision to:

- the current Product Spec baseline commit;
- the SHA-256 digest of the decision record;
- each referenced saturated research packet and its SHA-256 digest;
- the validated benchmark/prototype/adversarial record and its SHA-256 digest;
- optional repository-local tracked inputs such as assumptions, threat models, benchmark fixtures, dependency locks, architecture contracts or security policy.

`cmdr-dev decision-freshness-snapshot --decision-id ENG-DEC-...` generates the canonical automatic portion of this snapshot. Optional tracked inputs may then be added and reviewed before acceptance.

## Revisit triggers

Every record declares automatic triggers for expiry and evidence drift. Additional external triggers are required when material:

- Class C: material new research;
- performance-sensitive: benchmark regression;
- security-sensitive: security advisory and threat-model change;
- dependency decisions: upstream version change;
- tracked assumption/threat-model/benchmark inputs: the corresponding trigger.

External triggers require a dated signal observation with evidence. A fired signal makes the decision stale.

## Fail-closed behavior

A stale decision cannot remain silently accepted. `CHK-DECISION-FRESHNESS` fails and requires either:

1. move the decision to `revisit-required`; or
2. repeat the necessary research/benchmark/adversarial work and write a new reviewed freshness snapshot.

Validation does not mutate the decision automatically.

## Reusable cache

`cmdr-dev decision-cache --as-of-date YYYY-MM-DD` returns only accepted, fresh decisions. Each entry includes its decision key, class, freshness horizon, aggregate basis digest, research packet IDs and validation ID. A fresh agent can therefore reuse previously proven work without replaying research until a revisit trigger fires.

The cache is generated from canonical registries rather than hand-maintained, so stale evidence cannot remain silently cached.
