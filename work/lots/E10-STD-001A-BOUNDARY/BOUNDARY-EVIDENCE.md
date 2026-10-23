# E10-STD-001A-BOUNDARY — CMDR Studio Library Asset Boundary and Metadata Contracts

## 1. Executive Summary & Purpose

This document establishes the authoritative preimplementation boundary controls, asset metadata contracts, tenant scope invariants, permission enforcement rules, and error handling behaviors for **CMDR Studio Library and Asset Catalog** under canonical capabilities **`CAP-STD-001`** ("Studio Library and Asset Catalog") and **`CAP-STD-002`** ("Studio Asset Metadata, Ownership and Classification").

This boundary satisfies and operationalizes:
- **`REQ-PROD-006`**: Strict separation of visibility from executability (catalog read access never grants execution rights).
- **`REQ-PROD-016`**: Unified catalog projection across all Studio asset families without conflating canonical asset lifecycles.
- **`REQ-OBJ-009`**: Complete, deterministic asset metadata representation with explicit ownership and version lineage traceability.

This task is strictly a **control-plane boundary governance unit**:
- Production files added: `0`
- Net production LOC: `0`
- Product Spec modification: `0 bytes` (strictly immutable against baseline commit `82e630fce3a7f6236e03709116eca6dd0a1acd40`).

---

## 2. Boundary Definition & Architectural Non-Goals

The CMDR Studio Library provides a **permission-aware, tenant-isolated catalog projection** of reusable Studio assets and references.

### 2.1 Architectural Invariants & Non-Goals
1. **NOT a Package Repository**: The Studio Library does not host raw package binaries, tarballs, or package registry endpoints.
2. **NOT a Search Engine**: Search indexing and full-text ranking algorithms are owned by Shared Capabilities. Studio Library consumes Shared Search outputs and maintains safe local browsing/filtering fallbacks.
3. **NOT an Execution Engine**: Studio Library does not execute, invoke, or trigger Tools, Workflows, Skills, or Agents. Execution belongs to specialized execution runtimes and Govern decision authorities.
4. **NOT an Asset Lifecycle Owner**: Studio Library does not own the canonical lifecycle or persistent state of individual assets. Each asset family (`skill`, `workflow`, `automation-agent`, `agent-team`, `tool`) maintains its authoritative domain lifecycle; Studio Library projects authorized views.
5. **Zero Credential Custody**: Studio Library never stores, resolves, or exposes raw credentials or API keys (`CAP-STD-015` §15).

---

## 3. Class 0 Catalog Projection Controls

All operations exposed by the Studio Library boundary are strictly **Class 0 (Read-Only Discovery and Inspection)**:

| Action | Execution Class | Description | Authority / Preconditions |
|---|---|---|---|
| **Browse** | Class 0 | Enumerate authorized asset summaries in tenant scope | `perm.studio.library.read` |
| **Search** | Class 0 | Query catalog projections via Shared Search or local filters | `perm.studio.library.read` |
| **Filter** | Class 0 | Multi-faceted filtering by type, status, owner, risk, AI flag | `perm.studio.library.read` |
| **Inspect** | Class 0 | View detailed asset metadata, schema, dependencies, consumers | `perm.studio.library.read` + asset-specific read |
| **Open** | Class 0 | Navigate to canonical owning surface (Builder, Skill Editor) | Return-origin preservation; no in-situ mutation |
| **Compare** | Class 0 | Side-by-side read-only diff of two versions or assets | Caller must hold read on both assets in same tenant |
| **Clone Preparation** | Class 0 | Generate draft clone descriptor preserving provenance | Single-tenant target only; no cross-tenant copy |
| **Execution / Invoke** | **FORBIDDEN** | Triggering tool calls, workflow runs, or agent tasks | Strictly prohibited within Studio Library |

### Invariant: Visibility Does Not Imply Executability (`REQ-PROD-006`)
Under `CAP-STD-001` §2 and §25, holding visibility or read permission on an asset catalog entry **never** conveys execution authority. The catalog projection surfaces execution eligibility metadata (e.g. required privileges, risk level, environment prerequisites) for informational triage, while invocation actions remain permanently blocked or disabled within the catalog boundary.

---

## 4. Strict Tenant Scope Enforcement (`tenant-scope-contract.md`)

In conformance with `cmdr-product-spec/17-implementation-contracts/tenant-scope-contract.md`, `14-security-permissions-and-trust/tenant-isolation.md`, and `ADR-0008`:

### 4.1 Mandatory Explicit Tenant Ref
- Every catalog request MUST supply an explicit, non-empty `tenant_ref`.
- Requests omitting `tenant_ref` MUST be rejected immediately with standard error code `missing-tenant`.
- Default or fallback tenants are strictly prohibited.

### 4.2 Absolute Prohibition of Wildcard Tenant Querying
- Wildcard tenant selectors (`*`, `%`, `all`, `""`) are strictly rejected with error code `wildcard-tenant`.
- Global cross-tenant scans are architecturally blocked.

### 4.3 Fail-Closed Cross-Tenant Isolation
- Assets belonging to Tenant B MUST NEVER be visible or discoverable to a session scoped to Tenant A.
- Any request attempting to access an asset across tenant boundaries MUST fail closed with error code `cross-tenant` (or `asset-not-found`) without disclosing whether the asset exists in any other tenant.

### 4.4 MSSP Read-Only Aggregation on Explicit Authorized Tenant Sets
- In Multi-Tenant Service Provider (MSSP) contexts, an aggregated read-only catalog projection is permitted ONLY when the caller provides an explicit `Authorized Tenant Set`.
- Invariant: `requested_tenants ⊆ authorized_tenant_set`.
- If any requested tenant is not in the caller's authorized set, the entire request is rejected fail-closed with `unauthorized-tenant-in-set`.
- Multi-tenant operations are strictly limited to Class 0 read aggregation. Cross-tenant mutation, asset migration, or cloning across tenants is rejected with `cross-tenant-mutation-forbidden`.

---

## 5. Server-Side Pre-Authorization & Permission Enforcement (`permission-enforcement-contract.md`)

In conformance with `cmdr-product-spec/17-implementation-contracts/permission-enforcement-contract.md` and `14-security-permissions-and-trust/permission-model.md`:

### 5.1 Server-Side Evaluation Only
- Client-supplied authorization flags or bypass tokens are rejected. The server evaluates caller identity against the tenant-scoped authorization policy on every request.

### 5.2 Pre-Authorization Evaluation
- Authorization decisions are evaluated *prior* to reading or projecting asset contents.
- If the caller lacks `perm.studio.library.read`, the entire catalog interface returns `permission-denied`.
- If the caller holds `perm.studio.library.read` but lacks permission for a specific asset type or asset instance:
  - Default: The asset is silently excluded from catalog listing to prevent metadata leakage.
  - Diagnostic / Explicit inspection: Projected as `permission-filtered` with all description, parameters, schema, and dependency details stripped.

### 5.3 Dual-Permission Scheme Compatibility
The boundary handles both canonical `perm.studio.*` permissions and legacy `perm.cmdr-studio.*` permissions:

| Asset Kind / Operation | Canonical Permission | Legacy Compatible Permission |
|---|---|---|
| Library Catalog Browse | `perm.studio.library.read` | `perm.cmdr-studio.library.read` |
| Skill Read / Inspect | `perm.studio.skill.read` | `perm.cmdr-studio.skill.read` |
| Skill Manage / Author | `perm.studio.skill.manage` | `perm.cmdr-studio.skill.manage` |
| Workflow Read / Inspect | `perm.studio.workflow.read` | `perm.cmdr-studio.workflow.read` |
| Workflow Manage / Author | `perm.studio.workflow.manage` | `perm.cmdr-studio.workflow.manage` |
| Automation Agent Read | `perm.studio.automation-agent.read` | `perm.cmdr-studio.automation-agent.read` |
| Automation Agent Manage | `perm.studio.automation-agent.manage` | `perm.cmdr-studio.automation-agent.manage` |
| Version Inspect | `perm.cmdr-studio.version.read` | `perm.studio.version.read` |
| Version Manage | `perm.cmdr-studio.version.manage` | `perm.studio.version.manage` |

---

## 6. Standard Error Handling & Zero Credential Leakage (`error-contract.md`)

In conformance with `cmdr-product-spec/17-implementation-contracts/error-contract.md` and `CAP-STD-015`:

### 6.1 Standard Error Envelope
All boundary errors return a uniform, deterministic envelope:
```json
{
  "error": {
    "code": "missing-tenant",
    "message": "Tenant reference is required for all Studio Library queries",
    "correlation_id": "corr-uuid-v4",
    "retryable": false,
    "partial_result": null
  }
}
```

### 6.2 Deterministic Error Code Registry
- `missing-tenant`: Tenant parameter was omitted or empty.
- `wildcard-tenant`: Tenant wildcard selector (`*`) was rejected.
- `cross-tenant`: Attempted access to an asset outside caller's tenant boundary.
- `unauthorized-tenant-in-set`: Requested tenant not present in authorized set.
- `cross-tenant-mutation-forbidden`: Attempted mutation across tenant boundaries.
- `missing-asset-id`: Target asset identifier not provided.
- `invalid-asset-kind`: Asset kind not among supported Studio asset kinds.
- `asset-access-denied`: Caller not authorized to access requested asset.
- `permission-denied`: Caller lacks general library read permissions.
- `invalid-cursor`: Cursor token is invalid, expired, or malformed.
- `stale-dependency`: Asset depends on a deprecated or invalid dependency.

### 6.3 Zero Credential Leakage Invariant
1. Assets requiring external API keys, service accounts, or webhook tokens reference them via opaque `secret_reference` descriptors (`ref_id`, `purpose`, `scope`, `rotation_metadata`).
2. Plaintext secrets are NEVER fetched, exposed, logged, or returned in catalog projections.
3. System stack traces, internal IP addresses, database schemas, and provider authentication tokens are strictly stripped from all error responses.

---

## 7. Asset Classification, Metadata Contract, and States

### 7.1 Supported Asset Kinds
The Studio Library strictly categorizes assets into five supported kinds:
1. **`skill`**: Reusable capability knowledge contract (`OBJ-SKILL`, `CAP-STD-010`).
2. **`workflow`**: Step orchestration and decision graph (`OBJ-WORKFLOW`, `CAP-STD-017`).
3. **`automation-agent`**: Autonomous agent definition (`OBJ-AUTOMATION_AGENT`, `CAP-STD-034`).
4. **`agent-team`**: Multi-agent collaborative team definition (`OBJ-AGENT_TEAM`, `CAP-STD-037`).
5. **`tool`**: Reusable executable capability contract (`CAP-STD-003`).

### 7.2 Canonical Lifecycle States
Source assets progress through six authoritative lifecycle states:
- `draft`: Work-in-progress, unvalidated, private to author/contributors.
- `review`: Submitted for validation, linting, and policy review.
- `approved`: Validated by testing and policy checks, ready for publication.
- `published`: Authoritative, immutable version available for tenant consumption.
- `deprecated`: Active but discouraged; successor version recommended.
- `superseded`: Fully replaced by successor; read-only historical reference.

### 7.3 Projected Functional States
The catalog projects assets into four deterministic functional states:
- **`ready`**: Validated, all required dependencies available and healthy, caller fully authorized.
- **`partial`**: Asset metadata visible, but one or more optional dependencies degraded, or Shared Search backend unavailable.
- **`stale`**: Asset version is superseded, or dependencies are deprecated/unvalidated.
- **`permission-filtered`**: Asset exists in tenant catalog, but caller lacks detailed read rights; rendered as redacted stub.

### 7.4 Explicit Ownership & `owner unknown` Enforcement (`CAP-STD-002` §8, §25)
- Owner identity is captured via structured object: `owner_id`, `owner_name`, `owner_type` (`user`, `team`, `system`), `is_unknown`.
- Invariant: When owner metadata is absent from the underlying asset record, the catalog MUST explicitly project `owner unknown` (`"owner_display": "owner unknown"`, `"is_unknown": true`).
- The system MUST NEVER guess, default, synthesize, or fabricate an owner identity.

### 7.5 Canonical Metadata Structure (`AssetMetadata`)
```json
{
  "id": "skl-01j9a8b7c6d5e4f3",
  "tenant_id": "tenant-corp-prod",
  "environment_id": "env-default",
  "type": "skill",
  "name": "triage-suspicious-login",
  "display_name": "Triage Suspicious Login",
  "description": "Deterministic and AI-assisted triage of suspicious authentication events.",
  "owner": {
    "owner_id": "usr-sec-ops-lead",
    "owner_name": "SOC Lead Analyst",
    "owner_type": "user",
    "is_unknown": false
  },
  "contributors": ["usr-analyst-1", "usr-analyst-2"],
  "version": {
    "version_id": "ver-1.2.0",
    "semver": "1.2.0",
    "created_at": "2026-10-15T12:00:00Z",
    "is_immutable": true
  },
  "lifecycle_state": "published",
  "projection_state": "ready",
  "execution_class": "class_0_inspect",
  "risk_classification": "low",
  "is_ai_powered": true,
  "has_non_ai_fallback": true,
  "dependencies": [
    {
      "kind": "tool",
      "ref_id": "tool-ip-reputation",
      "version_constraint": "^2.0.0",
      "is_optional": false,
      "health_status": "available"
    },
    {
      "kind": "secret_reference",
      "ref_id": "sec-ref-threat-intel-api",
      "purpose": "Threat intelligence lookup",
      "is_optional": true,
      "is_masked": true
    }
  ],
  "consumers": [
    {
      "consumer_type": "workflow",
      "consumer_ref": "wfl-incident-triage-v1"
    }
  ],
  "supersession": {
    "is_superseded": false,
    "superseded_by": null,
    "supersedes": "skl-01j8oldversion"
  },
  "provenance": {
    "created_by": "usr-author",
    "created_at": "2026-09-01T10:00:00Z",
    "updated_by": "usr-editor",
    "updated_at": "2026-10-15T12:00:00Z",
    "cloned_from": null
  }
}
```

---

## 8. Verification Status & Governance Rules

1. **Preimplementation Status**: This document registers the normative boundary contracts for CMDR Studio Library. No production runtime executable code is claimed or created in this task (`production_files: 0`, `net_production_loc: 0`).
2. **Product Spec Invariance**: Verified 0 bytes modified in `cmdr-product-spec/**` against baseline commit `82e630fce3a7f6236e03709116eca6dd0a1acd40`.
3. **Toolchain Health**:
   - `cmdr-dev validate-manifests`: PASS (201/201 manifests).
   - `cmdr-dev architecture-audit`: PASS (15/15 boundaries, 1466 path claims).
   - `cmdr-dev doctor`: PASS (`healthy=true`).
4. **Test Suite Non-Regression**:
   - Node.js test suites: 143/143 tests pass (61 in timeline demo, 82 in event-search frontend).
   - Go modules: 11/11 modules vet clean.
5. **Next Unlocked Work Unit**: Promotes `E10-STD-001A-BOUNDARY` to `IMPLEMENTED` and unlocks `E10-STD-001B-CATALOG` for catalog registry implementation.
