# Decision validation evidence

This registry turns opaque benchmark/adversarial references into validated engineering evidence.

It is intentionally empty until a real Class B/C decision reaches acceptance.

## Candidate trade-offs

An accepted significant/critical decision must compare at least two candidates. Every candidate is tied to a solution family from the decision's research packets and carries strengths, weaknesses and research packet evidence.

## Benchmarks

Performance-sensitive accepted decisions must reference representative CMDR benchmarks by exact `BENCH-*` IDs. A referenced benchmark must:

- record date, workload, environment, dataset scale and artifact identity;
- compare the preferred candidate with at least one alternative;
- contain metric results for every compared candidate.

This is CMDR measurement, distinct from external research claims.

## Prototype trigger

`external_evidence_sufficient` is explicit. When false, the accepted decision requires a conclusive prototype of the preferred candidate with no unresolved uncertainty.

## Adversarial falsification

Class C decisions must reference `ADV-*` reviews that actively attempt to falsify the preferred candidate. An accepted decision is blocked if the review is inconclusive/rejected, has unresolved blocking findings, targets a different candidate, or lacks research evidence.

## Counter-evidence

An accepted decision cannot retain `blocking_counter_evidence` or `unresolved_material_uncertainty`. This prevents a preferred candidate from being accepted while material contradictory evidence is still open.
