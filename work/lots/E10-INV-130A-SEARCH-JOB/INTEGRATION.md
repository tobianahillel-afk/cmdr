# E10-INV-130A — integrate Event Search Job with authorized Case Timeline

## Functional product slice
A working **offline read-only** search job path now composes the existing Event Search orchestrator (CAP-INV-002) and Investigation Timeline owner-authorized memory provider (CAP-INV-110). It is not another isolated validator: the same process reads a trusted bounded owner snapshot, reuses the existing Case→Timeline→typed source authorization flow, selects telemetry-event references for a user-entered literal query and half-open RFC3339 interval, executes the real Search Job lifecycle and returns sorted event references plus generated job/audit identities.

## Dependencies / architecture
- Upstream verified work: E10-INV-002G and E10-INV-110E on engineering foundation SHA ea1770b48a4bc1b2e850a7ffab15d5839f9e9e49.
- New separately registered composition boundary `product-runtime/case-search-workbench/**`; first-party imports only from `event-search` and `investigation-timeline` via local Go replacements.
- Neither source runtime's contracts or canonical Product Spec changed. No new index/store, orchestration engine, protocol server, third-party dependency or payable API. OPEN-013 remains open for governed mutations.
- A Case ID selects authorized context, **not** an Event Search data source. Only already authorized telemetry-event source references become the Search Job's source set.

## Behavioral acceptance
1. Actual disk snapshot → per-principal Case/Timeline/source owner → source discovery → created/running/completed first-party Search Job → audit and stable ref report. Source previews/payloads never searched or emitted.
2. Query is case-insensitive literal text over **authorized** Timeline ID, source reference and provenance only; exact timestamps use [from,to), uncertain timestamps excluded, no query language or raw payload search invented.
3. Cases across tenants fail closed. Denied typed source refs and non-telemetry evidence are never indexed, searched, or counted. The backend re-reads the complete authorized owner chain on job execution, and a revoked/missing previously eligible source blocks the entire job.
4. A complete no-match search is completed with zero result refs. Invalid query/scope, cancelled context, invalid source and inconsistent pagination emit no partial JSON report.
5. Linux Go tests, race, vet and gofmt; macOS/Windows Go test and vet, all at exact PR HEAD.

## Risks / qualification
Only synthetic operator-selected identity and demonstration Search Job execution permissions exist in this prototype. No production IAM, actual central Search permission provider, durable Search Job repository, live telemetry ingestion, search engine indexing, customer data, persistent audit trail, transactional Case snapshot, end-to-end UI navigation, SLO, or mutation is implied. The Case owner demo re-authorizes full Case snapshots for each page, so worst-case search performance remains unqualified beyond 5,000 records; no production-scale claim.

Status: IMPLEMENTED candidate; CI and independent risk-based L2 must be evaluated at the exact published HEAD. Keep PR draft pending review and parent engineering foundation integration; no merge to main.
