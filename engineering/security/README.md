# CMDR security engine

This directory is the engineering security control plane. It does not replace or reinterpret product security semantics from the canonical Product Spec.

## Security gate registry

`security-gates.json` is the deny-by-default inventory of security control families CMDR engineering must eventually enforce:

- secrets;
- SAST;
- SCA;
- SBOM;
- IaC;
- container security;
- authorization negative tests;
- tenant isolation;
- fuzzing;
- DAST;
- race/concurrency validation;
- supply-chain provenance/integrity.

Every gate declares its intended execution stages, blocking policy, risk domains, path triggers, permitted data handling and required evidence.

Readiness is explicit:

- `specified`: the gate contract exists but the implementation has not yet been verified;
- `active`: a reviewed implementation key exists and the gate is executable;
- `deferred-runtime`: the gate requires runtime/product artifacts that do not exist yet.

A specified or deferred gate must not claim an implementation key. An active gate must use a compiled implementation key known to `cmdr-dev`.

## Data boundary

The registry permits only three data policies:

- `source-local-only`;
- `artifact-local-only`;
- `synthetic-runtime-only`.

This prevents future security-tool integration from silently sending production/customer evidence to external services.

## Current status

E4-SEC-001A defines and validates the gate contracts only. It does **not** claim that secrets scanning, SAST, SCA, SBOM, DAST, fuzzing, tenant isolation or authorization testing are already implemented. Later E4 sublots activate gates only after executable evidence exists.


## Deterministic secret scanner

E4-SEC-002A introduces a local source scanner with redacted SHA-256 finding evidence and a path/fingerprint-scoped allowlist. It runs from the adaptive PR/push validation plan and can also perform a full Git-tracked scan. The implementation commit passed push and pull-request validation with zero findings, so `SEC-SECRETS-001` is now active under compiled implementation key `builtin-secret-scan-v1`.
