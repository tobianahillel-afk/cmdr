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


## Safety-constrained optimization recommendations

E8-C derives advisory recommendations only from a current verified metric snapshot and, when available, an older comparable verified baseline.

There are no absolute performance thresholds in this layer. An opportunity exists only when comparable evidence demonstrates a regression. This prevents the engine from inventing arbitrary targets.

The current recommendation families are:

- validation selectivity: only when validation tier and mandatory-check floor are unchanged while selected work and cost both regress;
- cache reuse: only when a measured cache-reuse rate regresses versus baseline;
- context bounding: only when deterministic context source/dependency counts grow versus baseline.

Before any recommendation is possible, the current snapshot must prove:

- mandatory safety execution is green;
- coordination conflicts are zero;
- recovery does not require revalidation;
- mandatory check floor is known.

Recommendations are advisory artifacts only. Their proof object hard-codes that direct plan mutation, mandatory-check mutation, prerequisite mutation and Product Spec mutation are not allowed. Any suggested change requires a later separately verified lot.

When no baseline exists, when metrics are unavailable, or when snapshots are not comparable, the optimizer returns zero recommendations and explicit withheld reasons. This is a successful safe outcome, not an error.


## Efficiency regression budgets

E8-D separates safety floors from efficiency budgets.

Safety and quality floors are enforced even when no comparable performance baseline exists. The mandatory-check floor is pinned to the verified E8-C floor and may increase but cannot decrease without an explicit reviewed policy change.

Efficiency comparisons use exact profiles composed of validation tier, sorted risk domains and mandatory-check floor. This prevents a new mandatory security gate from being rejected merely because it increases validation cost. For an exact comparable profile, validation selection, execution and relative cost have a zero-regression budget. Cache reuse also has a zero-regression budget when both snapshots contain measured cache evidence.

Context source/dependency growth is reported as advisory rather than blocking because required traceability may legitimately expand the deterministic context graph.

The baseline registry records only verified metadata: source SHA, workflow run, observed time, snapshot digest and metric values. No source code, secrets or product telemetry are stored.
