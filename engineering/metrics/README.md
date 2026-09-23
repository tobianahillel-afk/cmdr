# Engineering engine metrics and optimization

This directory measures the engineering control plane, not the CMDR product runtime.

Metrics are local repository evidence. No external telemetry service, analytics SDK or product-runtime instrumentation is introduced.

## Safety model

Every metric has a closed category, unit, aggregation, source key, safety class and direction.

Only metrics with `safety_class = efficiency` may set `optimization_eligible = true`.

`safety-floor` and `quality-floor` metrics are deliberately non-optimizable. An optimization that reduces runtime by removing mandatory checks, suppressing revalidation, allowing ownership conflict or weakening evidence is invalid regardless of its measured speed benefit.

Metric observations bind an exact registered metric and source, a finite bounded value, canonical UTC timestamp, full lowercase Git SHA and a short evidence reference. Raw source code, secrets, prompts and external telemetry destinations are not metric evidence.

E8-A validates the model only. E8-B compiles snapshots from already-produced CMDR engineering evidence; it does not scrape conversations or create a second source of truth.


## Runtime snapshot compilation

E8-B compiles a metric snapshot once, at the end of a successful adaptive validation execution. It reuses the already computed validation plan and execution accounting, then reads only metadata-scale context, cache, recovery and coordination state.

Every registered metric must appear exactly once as either:

- an observation bound to the current full Git SHA and one canonical source key; or
- an explicit unavailable item with a reason.

For example, an empty performance cache is recorded as `MET-CACHE-REUSE-RATE = unavailable / performance-cache:not-applicable`; it is never guessed as 0% or 100%.

`MET-SAFETY-GREEN` is derived from the selected mandatory checks. Missing or failed mandatory execution produces value 0. A snapshot compilation failure fails the validation run; metrics cannot silently disappear to make the engine look faster.
