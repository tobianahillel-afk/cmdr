# Offline Case/source-owner reauthorization across requests

## Functional increment (PR #40, read-only local evaluation)

With `case-timeline-demo -principal=investigator-42 -dataset=... -reload-dataset`, each HTTP request reopens the **operator-selected trusted local JSON snapshot**, validates its bounded structure and first-party typed owners, and rebuilds the existing Case/Timeline Entry/source authorizers. The existing immutable startup-snapshot mode remains the default; no new dependency, schema, provider protocol, public HTTP upload or persistent write is introduced.

The test changes Case and source-owner reader ACLs **after** the local HTTP server starts: authorized source navigation must become 404 after a source is revoked; Case Timeline and navigation must become 404 after Case access is revoked; malformed/deleted dataset must yield generic 503 without falling back to stale previously authorized data; restoring owner permissions allows access on the following request. Anonymous and foreign tenants remain denied. This is a useful offline forensic rehearsal of changing permissions, **not a production IAM, real-time remote ACL service, atomic snapshot across concurrent requests, durable audit, signed identity or a live ingest pipeline**.

Reuse-first alternatives: (a) rebuild the existing in-memory provider per request — selected for a small 16 MiB local dataset and simple read-only test; (b) add a file watcher/cache — rejected because it creates freshness windows and permission-revocation ambiguity; (c) add database/IdP client/index — deferred pending a separately authorized provider boundary and actual host interface. Request-time costs must be measured before claiming an SLO. The canonical Product Spec and OPEN-013 are unchanged.

PRs #37–#40 remain draft, and this branch targets the engineering integration branch rather than main. Do not merge or imply release readiness from a green unit test. The risk-based L2 should inspect this meaningful change to the source authorization boundary, exact SHA, and the full changed-snapshot HTTP test.
