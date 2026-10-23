# Case Search Workbench — read-only offline integration

This small Go application **composes two existing first-party runtimes** rather than creating a new Search Engine:

1. `investigation-timeline`: synthetic principal -> Case -> authorized Timeline entries -> typed source ACL; complete paginated traversal.
2. `event-search`: validated single-tenant Search Job envelope -> lifecycle -> backend completion and audit -> stable event references.

The operator supplies a **trusted synthetic** Case owner snapshot JSON; no browser, API, external data source or database is involved. The application searches already released **telemetry-event** Timeline references and provenance only, using a case-insensitive literal substring query and a precise half-open UTC interval. A Search Job includes its state, provenance, validated completed sources, sorted result references and operation audit trail.

From the repository root:

```sh
cd product-runtime/case-search-workbench
go run . -snapshot=/absolute/path/to/owner.json -principal=analyst-a -case=case-a -query=telemetry -from=2026-10-01T00:00:00Z -to=2026-10-02T00:00:00Z
```

No source preview, raw log text, foreign Case or denied source is searched or included in the output. At most 5,000 authorized Case rows are traversed; read-only searches do not mutate Cases, Events or relations. Source read permissions are checked again when executing the Search Job, and missing/revoked sources cause the entire job to fail closed. Real disk-to-Job, cross-tenant, source ACL and cancellation tests run on Linux/macOS/Windows.

**No production claims:** principal selection and four search execution permissions are synthetic, local and operator-controlled. This is *not* enterprise IAM, persistent Search Jobs, a live ingestion/index provider, a general query language, full-text payload search, an Event Inspector, transactional ACL snapshot, or production SLO evidence. The original CAP-INV-002 and CAP-INV-110 domain ownership stays unchanged; OPEN-013 still blocks governed Case mutation. No paid API or third-party library is used.

The isolated composition boundary is explicitly allowed to depend only on the two existing first-party modules, and both Go modules are resolved from local checkout paths.
