# CMDR product dependency policy

## Default

Third-party product functionality is **denied by default**.

CMDR should own the implementation of its product behavior, including UI/design-system components, authorization policy logic, workflow semantics, rule processing, orchestration, indexing/search behavior and domain logic.

## Trusted Base exception

A third-party runtime dependency is permitted only when:
- implementing it ourselves would materially reduce security or correctness, or it is an unavoidable platform primitive;
- its exact purpose is narrow;
- no hidden product behavior is delegated to it;
- its version/provenance/license/security posture are reviewed;
- network and telemetry behavior are known;
- transitive dependencies are known;
- replacement boundaries are explicit;
- the dependency is listed in `engineering/security/trusted-base.json`.

Security primitives such as established cryptographic algorithms must not be reimplemented merely to achieve a zero-dependency claim.

## Development tooling

Build, test, SAST, DAST, fuzzing and research tools are outside the product-runtime dependency boundary. They may be external when self-hosted/free use is appropriate, but must not silently become runtime dependencies or receive production/customer data.

## Enforcement target

Future dependency gates must fail a product-runtime dependency change that is not present in the Trusted Base allowlist.
