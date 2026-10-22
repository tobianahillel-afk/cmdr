# CMDR work-unit complexity policy

This policy sizes work for autonomous and human implementation. It applies to all executable non-test code, including CMDR engineering tooling; it is not limited to product runtime code.

## Soft target envelope

A lot, sublot or task should normally remain at or below:

- 2 primary components;
- 8 materially changed executable/non-test files;
- 400 net non-test LOC.

Exceeding a soft target produces a warning. It does not automatically block READY if no hard split trigger applies.

Planning uses the larger of the estimated and refined production values so an optimistic `production_files` or `net_production_loc` value cannot hide an oversized estimate.

## Mandatory split triggers

A non-epic work unit must remain below or equal to:

- 12 materially changed executable/non-test files;
- 800 net non-test LOC;
- 1 bounded context;
- 1 independent migration;
- 1 distinct security model;
- 1 separately testable behavior.

Any exceeded hard limit yields `split-required`.

A `split-required` unit may remain DRAFT, DECOMPOSED or BLOCKED while it is being redesigned. It must not enter READY or any later lifecycle state.

## Counting

Count executable/non-test implementation code and engineering-control-plane code.

Do not count:

- tests;
- fixtures;
- generated code;
- vendored code;
- purely generated evidence artifacts.

Configuration or schema files that materially implement runtime or enforcement behavior may be counted in the planning estimate even when they are not source-code files. The goal is bounded cognitive and review complexity, not gaming a LOC metric.

Epics are decomposition containers and are not directly evaluated against lot budgets.
