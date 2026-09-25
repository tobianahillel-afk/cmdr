# Product-runtime dependency enforcement

CMDR product-runtime dependencies are deny-by-default.

`cmdr-dev dependency-audit` scans only architecture boundaries whose `kind` is `product-runtime`. Engineering/build/test tooling is intentionally outside this runtime gate and is governed separately.

## Supported runtime manifests

The initial native parsers cover:

- Go: `go.mod` `require` entries; modules replaced by local relative paths are treated as repository-local;
- Node: `package.json` `dependencies`, `optionalDependencies` and `peerDependencies`; `devDependencies` are excluded and workspace/file/link dependencies are repository-local;
- C/C++ via vcpkg: `vcpkg.json` dependencies.

The scanner intentionally fails closed on recognized but unsupported package formats inside a product-runtime boundary, including Cargo, Python, JVM, Conan, Composer, Ruby, Swift, Dart and Elixir manifests. A parser must be added before that ecosystem becomes an allowed CMDR product-runtime stack.

## Trusted Base

Every third-party runtime dependency must match an exact `ecosystem + name + declared version` entry in `engineering/security/trusted-base.json`.

The Trusted Base starts empty. An exception requires explicit provenance, license, security review, telemetry behavior and replacement-boundary metadata. Wildcard version approvals are rejected.

This gate is about product runtime. Supply-chain controls for development/build tooling are handled by later security layers.
