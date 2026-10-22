# Deep security scheduling

CMDR separates **fast PR security gates** from expensive dynamic validation.

The canonical policy is `deep-security-policy.json`. It covers three dynamic gate families:

- `SEC-FUZZ-001`: fuzz targets;
- `SEC-DAST-001`: isolated web/API DAST targets;
- `SEC-RACE-001`: race/concurrency targets and language-specific sanitizer targets.

## Stage policy

The security-gate registry is the source of truth for supported stages.

- **PR**: dynamic gates are considered only when the real change set touches the security/architecture/CI control plane or a registered runtime security-critical path.
- **Nightly**: every active dynamic gate declaring `nightly` is considered.
- **Release**: every active dynamic gate declaring `release` is considered; release evidence must bind an exact artifact SHA-256.
- **On-demand**: active dynamic gates declaring `on-demand` may be selected, optionally by exact gate ID.

No active target means no expensive process is launched. The decision is recorded as deferred, not passed.

## Target registry

A future runtime target is registered directly in `deep-security-policy.json` with:

- stable target ID;
- canonical gate ID;
- target kind (`fuzz`, `dast`, `race`, or `sanitizer`);
- owning `product-runtime` boundary;
- owner;
- concrete target identifier;
- seed/corpus identity;
- configuration identity;
- minimum duration;
- synthetic-data-only requirement.

A registered target whose gate is not active fails closed. An active gate without a registered target also fails closed when the stage schedules that gate.

## Evidence contract

Selected targets require a JSON evidence file at:

`engineering/testing/deep-security-evidence/<target-id>.json`

Evidence must contain the exact gate/target/stage, current Git SHA, target identity, seed, configuration, duration, evidence kind, `synthetic_data=true`, and `result=pass`. Release evidence additionally requires a 64-hex artifact SHA-256.

This lot does not invent fuzz/DAST/race targets because CMDR currently has no product-runtime boundary. The scheduled workflow therefore proves policy/scheduling truthfulness today and becomes fail-closed automatically when targets are later registered.
