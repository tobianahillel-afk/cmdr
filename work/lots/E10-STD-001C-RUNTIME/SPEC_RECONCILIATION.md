# SPEC_RECONCILIATION — CMDR Studio Library & Asset Catalog Runtime

**Document ID**: `work/lots/E10-STD-001C-RUNTIME/SPEC_RECONCILIATION.md`  
**Task ID**: `E10-STD-001C-RUNTIME`  
**Parent Epic**: `E10-SCALE-007`  
**Canonical Capabilities**: `CAP-STD-001`, `CAP-STD-002`  
**Governance Date**: 2026-10-23  
**Status**: APPROVED & NORMATIVE  

---

## 1. Formal Qualification Distinctions

Under the CMDR V2 Quality & Governance Protocol, progress claims must never conflate distinct stages of engineering maturity. The "0 LOC Studio Paradox" in tasks `E10-STD-001A-BOUNDARY` and `E10-STD-001B-CATALOG` arose from treating documentary boundary specification as full implementation.

To maintain absolute technical truth, the following five qualification tiers are formally distinguished:

| Tier | Qualification | Definition & Exit Criteria | Evidence Required | Studio Status |
|---|---|---|---|---|
| **1** | **Documentary Validation (`DOCUMENTARY_PASS`)** | Boundary, contracts, schemas, error codes, and permission rules are fully specified and reconciled with canonical product specs. Zero production files or LOC may be claimed. | Approved `EVIDENCE.md`, valid `manifest.json`, clean architecture audit. | **SATISFIED** (E10-STD-001A, E10-STD-001B) |
| **2** | **Runtime Implementation (`IMPLEMENTED`)** | Functional, executable source code exists under authorized `product-runtime/` paths. All 8 Class 0 operations execute real projection logic with zero external dependencies. | Go source files in `product-runtime/studio-catalog`, green local unit & integration test suites. | **ACTIVE TARGET** (E10-STD-001C) |
| **3** | **Integrated Verification (`VERIFIED`)** | Runtime code is integrated into master integration branch (PR #222) with 100% passing tests, independent L2 reviews, zero regression on existing suites, and green CI. | 18/18 GitHub Actions green checks, L2 gate approvals at exact commit SHA. | **PENDING** |
| **4** | **Product Acceptance (`ACCEPTED`)** | End-to-end user workflows, API contracts, and screen behaviors validated against user-facing criteria. | Forensic audit signoff, operator validation run. | **PENDING** |
| **5** | **Production Release (`RELEASED`)** | Deployed and operational in production infrastructure. | Deployment receipt, telemetry baseline. | **FUTURE** |

### Policy Invariant:
No task may claim `VERIFIED` status in `work/graph.json` or `engineering/state/current.json` unless executable runtime code exists and has passed deterministic automated test verification at the exact commit SHA.

---

## 2. Canonical Permissions Reconciliation

All permissions referenced in Studio Library runtime and documentation must correspond strictly to entries in `cmdr-product-spec/14-security-permissions-and-trust/permission-catalog.md`. Invented or hallucinated permission identifiers are strictly forbidden.

### 2.1 Complete Studio Permission Mapping

| Studio Operation | Canonical Catalog Identifier (`permission-catalog.md`) | Line in Spec | Status / Rule |
|---|---|---|---|
| **Studio Library Browse / Query** | `perm.studio.library.read` | Line 192 | **Canonical**. Sole library-level read permission. |
| **Skill Inspect / Read** | `perm.studio.skill.read` | Line 196 | **Canonical**. Required for detailed skill schema inspection. |
| **Skill Author / Manage** | `perm.studio.skill.manage` | Line 195 | **Canonical**. Governs authoring/mutation. |
| **Workflow Inspect / Read** | `perm.studio.workflow.read` | Line 198 | **Canonical**. Required for workflow graph inspection. |
| **Workflow Author / Manage** | `perm.studio.workflow.manage` | Line 197 | **Canonical**. Governs workflow authoring. |
| **Automation Agent Inspect** | `perm.studio.automation-agent.read` | Line 185 | **Canonical**. Governs agent definition inspection. |
| **Automation Agent Manage** | `perm.studio.automation-agent.manage` | Line 184 | **Canonical**. Governs agent lifecycle management. |
| **Agent Team Inspect** | `perm.studio.agent-team.read` | Line 179 | **Canonical**. Governs team composition inspection. |
| **Agent Team Manage** | `perm.studio.agent-team.manage` | Line 178 | **Canonical**. Governs team configuration. |
| **Version Inspect** | `perm.cmdr-studio.version.read` | Line 28 | **Canonical**. Note: Registered under `cmdr-studio` namespace. |
| **Version Manage** | `perm.cmdr-studio.version.manage` | Line 27 | **Canonical**. Note: Registered under `cmdr-studio` namespace. |

### 2.2 Rectification of Non-Catalog & Hallucinated Permissions

The following identifiers previously claimed in draft evidence are formally disqualified:

1. **`perm.studio.tool.read` & `perm.cmdr-studio.tool.read`**:
   - **Reason**: `tool` is not registered as an independent canonical object in `cmdr-product-spec/05-domain-model/objects/` and has zero entries in `permission-catalog.md`.
   - **Resolution**: Tool assets within Studio Library are governed strictly by the library discovery permission **`perm.studio.library.read`**. No dedicated tool read permission exists or will be enforced.
2. **`perm.cmdr-studio.library.read`**:
   - **Reason**: Absent from `permission-catalog.md`. The catalog strictly registers `perm.studio.library.read` (line 192).
   - **Resolution**: Stripped. The runtime exclusively accepts `perm.studio.library.read`.
3. **`perm.studio.version.read`**:
   - **Reason**: Absent from `permission-catalog.md`. The catalog strictly registers `perm.cmdr-studio.version.read` (line 28).
   - **Resolution**: Mapped to canonical `perm.cmdr-studio.version.read`.

---

## 3. Tenant Isolation & Multi-Tenant Architecture

In strict compliance with `cmdr-product-spec/17-implementation-contracts/tenant-scope-contract.md` and `ADR-0008` (Approved Decisions D4 & D5):

### 3.1 Three Distinct Operational Scopes

1. **Single-Selected Tenant MVP (Primary Operational Mode)**:
   - Every request must include an explicit `tenant_id`.
   - Wildcards (`*`, `%`, `all`, `""`) are strictly rejected with error code `wildcard-tenant`.
   - Cross-tenant requests without explicit multi-tenant authorization return `404 asset-not-found` (see Section 4).
2. **MSSP Read-Only Aggregation (Authorized Tenant Set)**:
   - Authorized Multi-Tenant Service Provider (MSSP) sessions receive an explicit `Authorized Tenant Set` (ATS) from Security authorization projection (`ADR-0008` APPROVED_D4).
   - Requests may specify `tenants: [T1, T2, ...]`.
   - Invariant: `requested_tenants ⊆ authorized_tenant_set`. If any requested tenant falls outside the ATS, the entire request fails closed with `unauthorized-tenant-in-set`.
   - Cross-tenant mutation, administration, and automatic export widening are **strictly forbidden** (`cross-tenant-mutation-forbidden`).
3. **Search Multitenancy (Deferred per ADR-0008 APPROVED_D5)**:
   - `ADR-0008` APPROVED_D5 explicitly mandates: *"Search is single-selected-Tenant initially. Multi-tenant Search is deferred."*
   - Therefore, in ATS aggregation mode, catalog discovery MUST NOT attempt to execute a global multi-tenant search query across tenants. It must execute federated single-tenant queries or aggregate per-tenant catalog projections, preserving the distinct `tenant_id` on every projected asset record.

---

## 4. Existence Oracle Elimination Rules

An existence oracle occurs when an API inadvertently reveals whether a resource exists in a restricted security boundary through differentiating response codes (e.g. 403 Forbidden vs 404 Not Found), response timing, or distinct error payloads.

### 4.1 Cross-Tenant Probing Defense (Uniform 404)
- **Scenario**: Caller authenticated in Tenant A requests asset `skl-corp-secret` belonging to Tenant B.
- **Vulnerability**: Returning `403 Forbidden` confirms that `skl-corp-secret` exists in the system.
- **Rule**: The catalog runtime **MUST NOT** return `403 Forbidden`, `cross-tenant`, or any hint of existence.
- **Mandatory Response**: The runtime MUST return a uniform `404 Not Found` with standard error code `asset-not-found`:
  ```json
  {
    "error": {
      "code": "asset-not-found",
      "message": "Asset skl-corp-secret was not found in tenant tenant-a",
      "correlation_id": "corr-uuid-v4",
      "retryable": false,
      "partial_result": null
    }
  }
  ```
- The response message, status code, headers, and execution latency must be indistinguishable from a request for a completely fictitious asset ID.

### 4.2 Same-Tenant Unpermitted Inspection (Sanitized 200 Stub)
- **Scenario**: Caller holds `perm.studio.library.read` in Tenant A. Asset `skl-admin-audit` exists in Tenant A, but requires `perm.studio.skill.read` for full inspection, which the caller lacks.
- **Listing Behavior (`Browse` / `Search`)**: The asset is silently excluded from catalog listing. Result counts reflect only authorized assets.
- **Direct ID Inspection (`Inspect`)**: To prevent ID confirmation timing side-channels while maintaining auditability, the catalog returns HTTP 200 with a redacted stub:
  - `projection_state`: `"permission-filtered"`
  - `id`, `name`, `display_name`, `type`, `lifecycle_state`: visible.
  - Sensitive details (`description`, `parameters`, `dependencies`, `secret_references`, `owner_id`): set to `null` or empty array `[]`.

---

## 5. Architectural Contracts: Error Envelope & Pagination

### 5.1 Error Envelope (`error-contract.md`)
All API errors return the standard five-field JSON envelope:
```json
{
  "error": {
    "code": "missing-tenant",
    "message": "Tenant identifier is required for all Studio catalog queries",
    "correlation_id": "corr-01j9a8b7c6d5e4f3",
    "retryable": false,
    "partial_result": null
  }
}
```

Standard error codes:
- `missing-tenant`: Tenant parameter was omitted or empty.
- `wildcard-tenant`: Tenant wildcard selector (`*`, `%`, `all`) was rejected.
- `asset-not-found`: Target asset does not exist in the requested tenant (or belongs to a foreign tenant).
- `permission-denied`: Caller lacks `perm.studio.library.read`.
- `unauthorized-tenant-in-set`: Requested tenant not in caller's Authorized Tenant Set.
- `cross-tenant-mutation-forbidden`: Attempted cross-tenant mutation or clone.
- `invalid-cursor`: Cursor token is invalid, expired, or corrupted.
- `invalid-asset-kind`: Asset kind not among supported Studio kinds (`skill`, `workflow`, `automation-agent`, `agent-team`, `tool`).

### 5.2 Deterministic Cursor Pagination (`pagination-filtering-sorting.md`)
- **Ordering Key**: Compound key `(sort_field, id ASC)` to ensure 100% deterministic tie-breaking across repeated queries.
- **Cursor Format**: Opaque, server-generated URL-safe base64 token encoding `(sort_value, id, direction)`.
- **Page Limits**:
  - Minimum: 1 item.
  - Default: 25 items.
  - Maximum: 100 items. Requests exceeding 100 are clamped to 100.
- **Drift Prevention**: Cursors navigate using keyset evaluation (`WHERE (col, id) > (val, last_id)`) rather than volatile numeric offsets.

---

## 6. The 8 Class 0 Operations Blueprint

All operations are read-only Class 0 actions under `REQ-PROD-006` (no implicit execution triggers):

1. **`Browse`**: Paginated listing of assets in tenant scope filtered by server-side authorization.
2. **`Search`**: Keyword and metadata text search over authorized catalog projections (local deterministic matching).
3. **`Filter`**: Multi-dimensional filtering by `asset_kind`, `lifecycle_state`, `risk_classification`, `is_ai_powered`, and `owner_id`.
4. **`Group`**: Aggregation of catalog item counts grouped by `asset_kind`, `lifecycle_state`, or `risk_classification`.
5. **`Open`**: Resolves canonical navigation target URL and return-origin descriptor without executing mutation.
6. **`Inspect`**: Detailed metadata inspection returning full schema or `permission-filtered` sanitized stub.
7. **`Compare`**: Side-by-side read-only diff of two versions or assets within the same tenant.
8. **`ClonePreparation`**: Constructs a pre-populated draft clone descriptor with provenance pointer (`cloned_from`) within the same tenant.
