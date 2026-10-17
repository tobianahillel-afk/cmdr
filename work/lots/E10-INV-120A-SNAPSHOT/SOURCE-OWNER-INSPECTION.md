# E10-INV-120A — owner-authorized local source inspection from a Timeline Entry

## Functional capability

Following an offline Case/Timeline or two-Case five-minute chronological drilldown, the analyst can explicitly inspect the **owner-authorized preview** of one selected event with no new API or provider:

```sh
go run ./cmd/case-timeline-demo -principal=analyst-a -snapshot=/absolute/path/owner.json -open-source-case=case-c -open-source-entry=event-c
```

The resulting local JSON receipt contains Case ID, Timeline Entry ID, typed source object ID, preview and a conspicuous synthetic marker. It is emitted only after the existing trusted in-memory owner has performed fresh Case → Timeline Entry → typed-source permission checks via `GET /demo-source-open`. The operator cannot provide or spoof an object ID or permission flag; the CLI never trusts an earlier comparison receipt as authority. No network listener starts in receipt mode. Denied/foreign/revoked sources return a generic error with no receipt.

## Technology choice

Reuse existing `loadDemoSnapshot`, `newDemoHandlerWithData`, `queryDemoOwner` and the owner `sourcePreviewHandler`; Go stdlib `url.Values` and `encoding/json`. A new source index, second permission policy, HTTP endpoint, database or third-party framework would duplicate logic without improving the bounded offline investigation.

## Security and production exclusions

The output may contain sensitive **authorized local source preview content**; use only with explicitly supplied operator-owned synthetic datasets. It is not a signed evidence object, full source attachment, chain-of-custody certificate, live SIEM connector, external IAM integration, durable Case store or production-ready remote source API. No preview is returned on partial failure. `OPEN-013` and canonical product model are unchanged. Tests must verify actual disk-only preview, denied tenant and source, malformed input, cancellation and writer error. An independent risk-based L2 is appropriate for this new source data exposure boundary.
