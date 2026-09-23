# Performance engineering

This directory is the machine-readable performance control plane for CMDR engineering.

## Current runtime status

The architecture registry currently contains **no** `product-runtime` boundary. The performance registry therefore reports runtime coverage as `not-applicable`; it does not claim that product performance has passed.

As soon as a `product-runtime` boundary exists, `performance-registry-audit` fails unless that boundary has at least one registered performance target.

## Measurement contract

Every future target binds:

- an architecture boundary and owner path;
- a closed workload key;
- a declared execution environment;
- one or more PR/nightly/release/on-demand stages;
- change-trigger paths;
- metric kind, unit, aggregation, comparator and absolute budget;
- an optional relative-regression limit only on environments explicitly approved for relative comparisons.

Throughput metrics use lower bounds. Latency, CPU, memory, allocation and utilization metrics use upper bounds.

## Noise policy

Shared CI environments are useful for coarse bounded smoke budgets, but they are not authoritative for tight relative regression claims. Relative baselines require a dedicated or explicitly calibrated environment.

The Product Spec remains read-only. Performance evidence is an engineering implementation concern unless a product-owned requirement explicitly states otherwise.


## Deep performance validation

`deep-performance-audit` adds bounded profiling, load, soak and resource validation without making performance tooling a product-runtime dependency. Every deep target must reference an already registered performance target and a compiled handler key.

The current repository still has no product-runtime boundary and therefore no deep target is fabricated: the real repository result is `not-applicable`.

Global caps are machine-readable in `deep-performance-policy.json`: 600 seconds per target, 2048 MiB of Go runtime memory, concurrency 64 and at most two deep targets per run. The separate `performance-deep` workflow is bounded to 25 minutes and runs on nightly, release or explicit on-demand stages.

PR validation remains lightweight: it validates policy and change sensitivity but does not execute deep workloads. Heavy execution happens only in the separate workflow. Profile output is engineering evidence under `engineering/testing/deep-performance-evidence` in the workflow workspace and is never linked into the CMDR product runtime.
