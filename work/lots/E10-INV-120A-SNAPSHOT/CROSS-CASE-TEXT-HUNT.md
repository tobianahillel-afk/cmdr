# E10-INV-120A — bounded cross-Case indicator text hunt

## Functional integration

The offline analyst viewer now offers an explicit **cross-Case text hunt** across up to 8 Cases already discovered by the existing server-authorized `GET /demo-cases` endpoint. For each permitted Case, it uses **the existing** `GET /case-timeline?case_id=...&q=...&limit=200&offset=...` route; every page invokes the canonical Case, Timeline Entry, and typed-source owner authorization chain again. It never requests a private catalog or accepts a client-supplied tenant. The returned Case count and pivot are visible only if **every** selected Case and every page completed successfully.

The investigator enters an IP, domain, hash or other string (2-96 characters), selects **Rechercher dans mes Cases**, sees counts per authorized Case, and can pivot with one click to its already existing Case-specific search. Pivot is a UI hint that forces a fresh server authorization, **not** a source capability or verified indicator match. No additional API, database, detection rule, index, framework, vendor integration or storage is introduced.

## Risk and performance controls

- Hard bound 8 discovered Cases, max 5,000 visible matches per Case, page size 200. Exceeding a bound yields **no partial result**. Case IDs, page totals and duplicate visible IDs are checked.
- Same-origin requests only, authorization bearer in request header only; `credentials: omit`, `cache: no-store`, `redirect: error`.
- Token or hunt query edits abort in-flight work and invalidate stale buttons/results. Never store bearer in localStorage, query string or exported report.
- Unavailable/revoked Case, malformed pages, inconsistent totals, duplicate returned IDs and exceeded limits suppress **all** results; any displayed count includes only server-released authorized entries.
- Browser output uses `textContent`, not HTML injection. The hunt is limited to the selected immutable local synthetic owner snapshot but page-by-page reads are **not transactionally consistent**.
- At most 200 requests at the extreme bound. This is deliberately a local small-dataset analyst tool; an authorized server-side batched search with indexed query, cancellable cursor and audited durable provider would be required for production. No IOC exact matching, malware attribution, causality, real SIEM feed, IAM revocation synchronization, or forensic custody claim.

## Technology reuse

Alternative A (chosen): existing Go Case discovery, Timeline query and typed-source authorization plus a single 140-line first-party browser module; minimal changes, no new trust boundary. Alternative B: add a new multi-Case HTTP endpoint; duplicates trust aggregation and authorization and complicates security review. Alternative C: external search index/DB or vendor SIEM query: inappropriate without source-of-truth and OPEN-013 decisions, and unnecessary for the bounded 8-Case offline proof.

## Verification

Behavioral Node tests cover two authorized Cases and pivot, a second Case being revoked (no partial display), token rotation, over-bound Case discovery, and malformed/duplicate pages. Go integration exercises the real embedded assets and launches the Node harness when available. Verify Go race/vet/gofmt, Node test evidence and GitHub Actions exact SHA. This is a significant functional milestone for risk-based L2 review; earlier L2 analyses of PR41 do not automatically cover this increment. The unit remains READY until a verified decision; PR41 remains draft on PR36 integration base. No changes to canonical Product Spec or OPEN-013.
