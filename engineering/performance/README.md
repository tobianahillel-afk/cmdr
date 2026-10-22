# Performance engineering

This directory is the machine-readable performance control plane for CMDR engineering.

## Current runtime status

The architecture registry now contains the bounded E9 pilot `product-runtime` boundary `pilot-context-envelope-runtime`.

Its performance budget is registered **before implementation**, but its benchmark/deep handlers are intentionally non-executable until E9-PILOT-001C supplies the runtime. The target is therefore `on-demand` only during E9-B: registry/deep-policy coverage is real, while no pre-implementation runtime performance pass is fabricated.

`performance-registry-audit` continues to fail for any product-runtime boundary that lacks a registered target.

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

The E9 pilot has one pre-registered deep resource target tied to its performance target. During E9-B it is `on-demand` only and its handler deliberately returns “not executable until E9-C” if manually invoked. PR validation checks policy/coverage without fabricating runtime evidence.

Global caps are machine-readable in `deep-performance-policy.json`: 600 seconds per target, 2048 MiB of Go runtime memory, concurrency 64 and at most two deep targets per run. The separate `performance-deep` workflow is bounded to 25 minutes and runs on nightly, release or explicit on-demand stages.

PR validation remains lightweight: it validates policy and change sensitivity but does not execute deep workloads. Heavy execution happens only in the separate workflow. Profile output is engineering evidence under `engineering/testing/deep-performance-evidence` in the workflow workspace and is never linked into the CMDR product runtime.


## E9 pilot runtime activation

The pilot context-envelope runtime is now implemented. `pilot-context-projection-v1` is executable through a standard-library-only probe compiled inside the product-runtime module. The benchmark engine builds the probe once per target and reuses it for warmup and samples; reported latency is measured inside the probe around `Project(Input)`, excluding build and process-startup time.

The shared GitHub CI target now allows `pr` and `on-demand` stages and evaluates only the absolute `p95 <= 1 ms` budget. Relative regression remains disabled on the shared runner. The deep resource handler uses the same real runtime probe with synthetic fixtures and the registered memory/iteration bounds.
