# Cross-boundary dependency edges

CMDR architecture boundaries are not only path owners; their `may_depend_on` lists are enforced against observed local module dependencies.

`cmdr-dev boundary-edge-audit` currently extracts repository-local module edges from the product-runtime ecosystems already supported by the runtime dependency gate:

- Node workspace dependencies resolve by local package name;
- Node `file:`, `link:` and relative dependencies resolve by repository path;
- Go local `replace ... => ./path` dependencies resolve by repository path.

Self-boundary edges are allowed. A cross-boundary edge is allowed only when the source boundary explicitly lists the target in `may_depend_on`.

A local target that cannot be mapped to a registered `product-runtime` boundary fails closed. Duplicate package/module identities also fail.

File-level import analyzers are added when a concrete runtime language is selected; the package/module boundary gate exists before product code starts so the architecture cannot silently drift.
