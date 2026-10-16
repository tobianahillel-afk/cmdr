# E10-INV-110B-RUNTIME — real on-demand projection probe wiring

## Selection and scope

Reuse the existing first-party Case Queue Go probe pattern and CMDR benchmark/deep-resource control plane. Use Go's standard library and the existing `runPerformanceProcess` subprocess wrapper; no new benchmark framework, database, provider or paid service. A third-party index is unnecessary for this 5,000-entry in-memory bounded projection; it remains a separate future architecture decision.

## Executable source and limitations

- Executable projection: `product-runtime/investigation-timeline/cmd/perf-probe`, 5,000 synthetic already-authorized Timeline Entries.
- On-demand latency target: one preflight, 200 projections per sample, one warmup and five samples. The registry p95 is **across sample means**, not per-request latency and not a production SLO.
- On-demand deep resource target: 1,000 projections and 128 MiB Go heap limit. Reported peak is **periodically sampled Go heap**, not maximum RSS.
- Tests assert the executable identity, reject invalid CLI arguments and disallow source identity in structured measurement output.
- `investigation-timeline-navigation-v1` remains fail-closed **preimplementation**; no navigation runtime exists yet.
- No provider/backend latency, storage, search, network, UI, ACL-provider verification, sustained load or production SLO is represented. Do not confuse synthetic `allow` with authoritative production permission.

## Evidence policy

The registry now makes **on-demand execution possible**, but this document is **not a measured performance result**. Do not declare measured p95 or resource PASS without separately executing on the exact GitHub SHA, checking logs and preserving the evidence. Before the first successful CI, any statements about tests remain unverified. Product Spec and OPEN-013 remain unchanged.
