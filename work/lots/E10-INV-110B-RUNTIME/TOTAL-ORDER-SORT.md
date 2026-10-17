# E10-INV-110B — strictly ordered Timeline sort

## Bounded reuse choice

The existing Go `sort.SliceStable` implementation was correctly deterministic but used a stable sort even though the comparator already provides a **strict total ordering** for all permitted visible rows: (1) known timestamps precede unknown/uncertain timestamps, (2) known instants ascend, and (3) the unique Timeline Entry ID is the final tie breaker. A visible duplicate ID is rejected *before* sorting, making comparison ties between distinct rows impossible.

Use Go standard library `sort.Slice` with the **unchanged comparator**. It avoids introducing a new index, database, timeline engine, sorting dependency or custom algorithm and should reduce redundant work in this bounded 5,000-entry in-memory projection. This is a hypothesis to be measured, **not** a claimed speedup, and does not change the canonical Product Spec or ownership.

## Non-regression proof and known limits

`TestProjectionSortTotalOrderAcrossClockAndTimestampTies` exercises chronologically unequal entries, equal instants written with different numeric timezone offsets, ID ties, two forms of unknown time, sixty seeded input permutations, caller immutability and duplicate visible ID rejection. Existing 30-permutation mixed observed/inferred property checks, contract fixtures and benchmark full-row checks remain in force.

Follow the two GitHub CI workflows on the exact code SHA, then read their Go race/vet/coverage and local synthetic benchmark results before claiming any improvement. The historical 4.700ms batch-mean p95 / 5.15MB Go heap result belongs to an **older commit and different workload**; it must NOT be attributed to this sort change. Tests alone do not establish provider latency, production request-level p95, RSS, ACL-provider authenticity or full CAP-INV-110 readiness. OPEN-013 remains unresolved; B remains READY until the complete review gates are met.

No additional L2 task should be created for this same immutable SHA if one has already been scheduled.
