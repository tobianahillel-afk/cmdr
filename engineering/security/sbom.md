# Deterministic CycloneDX SBOM

CMDR generates SBOM evidence internally with `cmdr-dev sbom`. The generator targets CycloneDX JSON 1.7 and introduces no product runtime dependency.

## Determinism

The document intentionally omits a timestamp and random serial number. Its root component is versioned by the exact Git commit SHA and the complete JSON is SHA-256 digested. Component ordering, properties and BOM references are deterministic.

The generator writes only inside the repository boundary. A direct invocation without `--output` writes `engineering/security/evidence/cmdr.cdx.json`; adaptive CI writes to its repository-local temporary validation directory.

## Inventory scope

The generator inventories:

- direct dependencies found under registered `product-runtime` architecture boundaries as runtime scope;
- direct supported dependencies under `tools/cmdr-dev` as development scope;
- reviewed pinned development tools from `development-tools.json` as development scope.

Runtime components use CycloneDX `scope=required`. Development-only components use `scope=excluded`. Every component also carries explicit CMDR properties for ecosystem, manifest, directness and runtime/development scope.

Because the current source manifests do not resolve transitive closure, the BOM composition is explicitly marked `incomplete`. This prevents a direct-dependency inventory from being presented as a complete resolved graph.

## Fail-closed identity rules

The generator refuses:

- dependency manifest ecosystems that are recognized but unsupported;
- missing, wildcard, `latest` or otherwise unresolved Go versions;
- npm ranges such as `^`, `~`, wildcard or comparator ranges;
- vcpkg constraints when no resolved lock/baseline supplies an exact package identity.

Unknown ecosystems therefore cannot silently disappear.

## Provenance

The SBOM records:

- CycloneDX schema/spec version;
- generator name/version;
- exact source Git SHA;
- Go engineering toolchain version;
- reviewed development-tool module/version/commit/license metadata;
- SHA-256 of the serialized SBOM.

`SEC-SBOM-001` is active for the current source/development inventory. Evidence: pull-request run `35850552860` generated CycloneDX 1.7 with 2 development components, 0 runtime components, 0 unsupported manifests, exact source-SHA binding and deterministic SHA-256 digest. `SEC-SUPPLY-001` remains specified because no releasable CMDR product artifact/signature pipeline exists yet; SBOM generation alone is not artifact-signing/provenance enforcement.
