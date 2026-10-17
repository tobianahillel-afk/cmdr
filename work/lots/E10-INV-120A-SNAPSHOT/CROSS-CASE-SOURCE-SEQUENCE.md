# E10-INV-120A-SNAPSHOT — Cross-Case common-source event sequence

## Functional integration

The existing local, authenticated analyst comparison previously stopped at
counts of shared source references and a representative event button. This
increment reuses the **already authorized Case/Timeline HTTP page results**
from each of two independent server-side Case/Timeline/source access checks
and exposes a paginated chronological union for a selected common source.

- Each Case is fetched independently through existing `/case-timeline`, no new
  API, storage engine, source metadata authority, or backend search index.
- Grouped entries retain their owning Case, classification
  (observed/inferred/other), original timestamp and uncertainty flag.
- Cross-Case sort uses certain parseable timestamps first; uncertain
  timestamps sort after certain ones. Case/entry IDs give a deterministic
  secondary key and are never treated as proof of causality.
- Only complete, consistent, source-authorized page traversals create
  controls; 25 event selectors per page, bounded at 5000 entries per Case.
- Clicking an event invokes the existing two-step source navigation and
  owner preview HTTP chain, independently rechecking Case, Entry, source
  rights. No original bearer or source ref is passed in a URL.
- Changing token, Case, comparison Case, search, kind or time invalidates
  the controls; label and record text use inert `textContent`.

## Alternatives and limits

**Selected:** existing browser in-memory grouped `Map`, stable native sorting,
first-party DOM, existing source-open owner checks. **Deferred/rejected:**
new source-centric API, separate correlation index, third-party search engine
and persistent write layer. No new dependency or product contract.

This is **operator-local synthetic** analyst triage, NOT production IAM,
real source ingestion, transactionally consistent Case history, an evidentiary
chain-of-custody export, shared-source identity assurance or causal proof.
OPEN-013 remains unresolved. Tests/CI and risk-based independent L2 must
be checked on the exact published SHA.
