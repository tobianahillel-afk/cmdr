# E9 bounded pilot

The selected pilot slice is `PILOT-CONTEXT-001`: the deterministic semantics of cross-product Tenant/Environment context preservation defined by `CAP-SET-004`.

This is intentionally smaller than Case/Evidence lifecycle work. The selected capability has no open decision, performs Class 0/1 behavior only, and can exercise the complete engineering path without silently resolving `OPEN-013`, `OPEN-014`, the global query/search backend, or forensic-engine choices.

E9-A contains **scope evidence only**. It does not establish a runtime directory and contains no product implementation.

The machine-readable source is `scope.json`. `cmdr-dev pilot-scope-audit` verifies it against the canonical Product Spec inventory and product graph, checks the permission catalog, and requires the E9-A work manifest to carry the same Product Spec references.

The planned behavior is deliberately narrow:

- preserve authorized Tenant context;
- preserve Environment only when compatible;
- preserve a non-sensitive return-origin;
- consume a fresh destination-authorization outcome;
- clear incompatible context on Tenant change;
- never transport raw payloads or secrets.

Tenant/Environment administration, the authorization engine itself, UI rendering, persistence, cross-tenant mutation, new permissions and new screens are excluded.
