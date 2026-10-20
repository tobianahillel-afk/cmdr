# E10-INV-004I-PIVOTHTTP — Class-0 pivot HTTP handoff

Product owner Investigate CAP-INV-004; only output is an unexecuted CAP-INV-002 Query draft. Shared owns Telemetry Event, Search Job, Query; Settings owns parser/source. Engineering parent E10-INV-004, child I depends on G-HOST IMPLEMENTED and C-PIVOT VERIFIED; code owner existing event-inspection core plus existing Go inspection host; Git base PR #184 is not an engineering parent.

Source-identical blobs pivot.go 938b1d4ae090e2a4080604d8b7d155f24bec69de and pivot_presence_test.go 736a2b0a85dd3219e77cdd9aae6b8320f8a0bd53 copied unchanged from PR #182 SHA 82eab4cade8e20763f80509c560f61971001b7ec. No mutation of audited #181/#182/#184/#185 SHA or canonical product contract.

New existing-handler route: POST /investigate/event-inspection/pivot. It shares existing SessionResolver, SearchOwner, TelemetryOwner and immutable Project authorization, then calls PreparePivot. No browser-provided valueAuthorized, arbitrary returnOrigin, IAM decision or unclassified raw field is accepted. Source/projector server has to classify fields. Return origin is fixed server-side to event-search. Query draft preserves authorized field/provenance, Event, tenant/env, Search Job/run, correlation, time range; no execution, persistence or mutation.

Go httptest positives: authorized selected source field; enriched field producer/version; true present empty string; raw-derived value explicitly permitted by owner. Denials: missing empty field, raw-restricted field, wrong tenant/environment/Job, missing presence, mismatch, forged returnOrigin, invalid interval and class-2 mutation property. Existing Inspector POST unchanged. Workflow now tests core and host modules across three OS and race on Linux; existing security gates and coverage limits not relaxed.

BOUNDARIES: real IAM, native provider, Search Job storage, execution engine, complete browser pivot UI, TLS/CSRF host binding and persisted audit are NOT delivered. AuditRequired indicates obligation, not event saved. OPEN-013/014 remain blocked as scoped. User-provided time window is syntactically validated only; executable search must independently check Query scope and authorization. No native GitHub review/merge/acceptance and no production claim.

Merge dependency: consolidate source-identical #182 fix into this sibling candidate and review delta; #185 current L2 and SHA remain frozen and cannot be claimed to cover this new branch. This is one material Class-0 HTTP behavior, not a new pivot engine.
