## Follow-up simplification (shared authorization path)

The original dossier duplicated the `/demo-source-open` owner response parsing from `inspect.go`. A focused extract now exposes `inspectAuthorizedDemoEntry` to both workflows. The HTTP handler, bearer and typed source identity checks are unchanged. There is one authoritative parser, rather than two copies of the trust-boundary code. No additional provider, database, framework or network access.
