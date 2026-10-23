# E10-STD-001B-CATALOG — CMDR Studio Library Asset Catalog and Query Registry Contracts

## 1. Executive Summary & Purpose

This document establishes the authoritative architecture, query registry contracts, tenant isolation mechanisms, deterministic pagination controls, permission-aware projection rules, and error envelopes for the **CMDR Studio Library Asset Catalog and Query Registry**. It directly fulfills task **`E10-STD-001B-CATALOG`** under parent epic **`E10-SCALE-007`**, building upon the boundary and metadata contracts established in **`E10-STD-001A-BOUNDARY`**.

This evidence specification operationalizes and governs:
- **`CAP-STD-001`**: "Studio Library and Asset Catalog" — Permission-aware discovery, cataloging, search, filter, grouping, inspection, and owning-surface navigation for reusable Studio assets.
- **`CAP-STD-002`**: "Studio Asset Metadata, Ownership and Classification" — Common metadata contract, ownership integrity, classification taxonomy, and lifecycle-to-functional state projection.
- **`REQ-PROD-006`**: Strict architectural separation of visibility from executability (catalog read access never confers invocation rights).
- **`REQ-PROD-016`**: Unified catalog projection across diverse Studio asset kinds without collapsing or mutating underlying domain lifecycles.
- **`REQ-OBJ-009`**: Complete, deterministic metadata representation, lineage tracking, and provenance preservation.

### Boundary & Governance Invariants
- **Metadata Specification Lot**: This lot establishes authoritative contract definitions, behavioral specifications, and verification proofs.
  - Production runtime files added: `0`
  - Net production LOC: `0`
- **Product Spec Invariance**: Strictly `0 bytes` diff against canonical baseline commit `82e630fce3a7f6236e03709116eca6dd0a1acd40` in `cmdr-product-spec/**`.
- **Integrity Compliance**: Genuine specification and validation; no dummy facades, no hardcoded test shortcuts, and zero credential leakage.

---

## 2. CMDR Studio Library Asset Catalog Registry Core

### 2.1 Architectural Role and Boundary Positioning
The CMDR Studio Library Asset Catalog is a **read-only, tenant-isolated projection engine** and query registry. It sits between canonical Studio asset authoring domains and consuming client surfaces:

```
+-----------------------------------------------------------------------------------+
|                        Consuming Client Surfaces                                  |
|   (Automation Designer, Studio Reviewer, SOC Analyst, Command/Investigate Linking)|
+-----------------------------------------------------------------------------------+
                                         │
                 Class 0 Query Operations (Browse, Search, Filter,
                 Group, Inspect, Open, Compare, Clone Preparation)
                                         ▼
+-----------------------------------------------------------------------------------+
|               CMDR Studio Library Catalog Registry Core (CAP-STD-001)             |
|  - Server-Side Pre-Authorization Engine (perm.studio.library.read / legacy)       |
|  - Tenant Scope Validator (Strict single-tenant & MSSP Authorized Sets)          |
|  - Cursor Pagination & Deterministic Sorting Engine                               |
|  - Functional State Projector (ready, partial, stale, permission-filtered)        |
|  - Owner Boundary Preserver (Explicit "owner unknown" fallback)                   |
|  - Zero-Leakage Error Envelope Normalizer                                         |
+-----------------------------------------------------------------------------------+
                                         │
                        Read-Only Upstream Projections
                                         ▼
+───────────────────────────────────────────────────────────────────────────────────+
|                           Canonical Asset Domains                                 |
|  Skills (CAP-STD-010) │ Workflows (CAP-STD-017) │ Automation Agents (CAP-STD-034) |
|  Agent Teams (CAP-STD-037) │ Tools (CAP-STD-003) │ Secret References (CAP-STD-015) |
+───────────────────────────────────────────────────────────────────────────────────+
```

### 2.2 Core Registry Invariants & Non-Goals
1. **Projection-Only Layer**: The Catalog Registry does not persist or own the primary records of assets. Every asset record is projected from its authoritative domain store (e.g., Skill Store, Workflow Definition Store, Agent Registry, Tool Specification Store).
2. **Decoupled Search Backend**: Catalog queries execute primarily against indexed local projections with deterministic fallbacks when Shared Search is unavailable (`CAP-STD-001` §25).
3. **No Execution Authority**: The catalog contains zero execution triggers or runtime runners. Displaying an asset does not make it callable (`REQ-PROD-006`).
4. **No Direct Secret Storage**: External secrets, keys, and tokens are referenced solely by opaque `secret_reference` tokens (`ref_id`, `purpose`, `masked: true`) under `CAP-STD-015`.
5. **Deterministic Immutability**: Versioned catalog records are immutable once published. Updates create new version entries; old versions are marked `superseded` with explicit bi-directional lineage.

---

## 3. Query Interface Operations & Contracts (Class 0)

All operations exposed by the CMDR Studio Library query registry are strictly **Class 0 (Read-Only Discovery and Inspection)**. No state-mutating actions or execution invocations are permitted.

| Operation | Input Signature | Output Envelope | Preconditions & Authority |
|---|---|---|---|
| **Browse** | `tenant_ref`, `kind?`, `limit?`, `cursor?`, `sort?` | `CatalogPageResponse<AssetSummary>` | Caller holds `perm.studio.library.read` |
| **Search** | `tenant_ref`, `query_text`, `kind?`, `limit?`, `cursor?` | `CatalogPageResponse<AssetSummary>` | Caller holds `perm.studio.library.read` |
| **Filter** | `tenant_ref`, `filter_criteria`, `limit?`, `cursor?` | `CatalogPageResponse<AssetSummary>` | Caller holds `perm.studio.library.read` |
| **Group** | `tenant_ref`, `group_by_facet`, `filter_criteria?` | `CatalogGroupResponse<AssetSummary>` | Caller holds `perm.studio.library.read` |
| **Open** | `tenant_ref`, `asset_id`, `version_id?`, `return_origin` | `OpenNavigationDescriptor` | Caller holds read permission on asset |
| **Inspect** | `tenant_ref`, `asset_id`, `version_id?` | `AssetDetailResponse` | Caller holds `perm.studio.library.read` + asset read |
| **Clone-with-provenance** | `tenant_ref`, `source_asset_id`, `version_id?` | `DraftCloneDescriptor` | Target tenant == source tenant; read permission |
| **Compare** | `tenant_ref`, `base_asset_id`, `target_asset_id`, `versions?`| `AssetComparisonResponse` | Single tenant; read permission on both assets |

### 3.1 Detailed Operation Semantics

#### 1. Browse
- Enumerates assets within the caller's tenant scope.
- Supports filtering by asset kind (`skill`, `workflow`, `automation-agent`, `agent-team`, `tool`).
- Returns standardized cursor-paginated summaries with total count and has_more flag.

#### 2. Search
- Accepts free-text queries matching asset name, display name, description, tags, and capability keywords.
- When Shared Search is reachable, ranking and relevance scoring are incorporated.
- When Shared Search is degraded or offline, gracefully degrades to deterministic substring matching on indexed local metadata with functional state reported as `partial`.

#### 3. Filter
- Evaluates multi-faceted structured predicates joined by logical `AND`:
  - `kind`: one of `skill`, `workflow`, `automation-agent`, `agent-team`, `tool`.
  - `lifecycle_state`: one of `draft`, `review`, `approved`, `published`, `deprecated`, `superseded`.
  - `functional_state`: one of `ready`, `partial`, `stale`, `permission-filtered`.
  - `owner_id`: exact match on principal/team identifier, or sentinel `"unknown"`.
  - `risk_level`: one of `low`, `medium`, `high`, `critical`.
  - `is_ai_powered`: boolean flag indicating presence of AI-assisted capabilities.
  - `tags`: subset match against asset classification tag set.

#### 4. Group
- Partitions the authorized tenant catalog into categorical buckets without reading restricted payloads.
- Supported facets: `by_kind`, `by_lifecycle_state`, `by_owner`, `by_risk_classification`.
- Empty buckets are omitted; items in `permission-filtered` state are grouped by kind/owner as redacted counts.

#### 5. Open
- Generates a secure, origin-preserving navigation descriptor directing the client to the canonical authoring surface (e.g., Skill Editor, Workflow Canvas, Agent Configurator).
- Enforces `return_origin` parameter preservation so analysts return to their exact catalog page/filter state upon closing the owning surface.
- Never performs inline mutation inside the catalog viewport.

#### 6. Inspect
- Returns exhaustive asset metadata: parameter schema, prerequisites, secret references (opaque), dependency tree, registered consumers, and supersession links.
- If the caller holds `perm.studio.library.read` but lacks deep inspection rights for a specific asset family (e.g., restricted workflow), sensitive parameters and source code references are stripped, and the response is marked `permission-filtered`.

#### 7. Clone-with-provenance
- Generates a draft clone descriptor intended for initialization in an authoring surface.
- Strictly captures upstream provenance:
  - `cloned_from_id`: immutable ID of original asset.
  - `source_version`: exact semver/version_id of copied asset.
  - `cloned_by`: caller principal ID.
  - `cloned_at`: UTC timestamp.
  - `provenance_chain`: list of preceding clone references.
- Invariant: Clones CANNOT cross tenant boundaries. `target_tenant == source_tenant` is enforced fail-closed.

#### 8. Compare
- Produces a deterministic side-by-side structural diff between two asset records or two versions of the same asset.
- Both assets MUST belong to the caller's active tenant scope. Cross-tenant comparisons are strictly forbidden and return `cross-tenant`.
- Diff components: parameter additions/removals/modifications, dependency version constraint changes, risk classification shifts, and AI capability additions.

---

## 4. Tenant-Scoped Catalog Queries & Isolation Invariants

In strict conformance with `cmdr-product-spec/17-implementation-contracts/tenant-scope-contract.md` and `ADR-0008`:

### 4.1 Strict Single-Tenant Enforcement
1. **Mandatory Explicit Parameter**: Every catalog query MUST include a non-empty `tenant_ref`.
2. **Rejection of Missing Tenant**: Omitting `tenant_ref` immediately halts processing and returns:
   ```json
   {
     "error": {
       "code": "missing-tenant",
       "message": "Explicit tenant reference is mandatory for all catalog queries",
       "correlation_id": "corr-uuid-v4",
       "retryable": false,
       "partial_result": null
     }
   }
   ```
3. **Prohibition of Wildcards**: Wildcard selectors (`*`, `%`, `all`, `""`, `null`) are rejected with `wildcard-tenant`.
4. **No Ambient Fallback**: The catalog NEVER falls back to a default, global, system, or cached tenant.

### 4.2 MSSP Read-Only Aggregation via Authorized Tenant Sets
1. In Multi-Tenant Service Provider (MSSP) environments, catalog browse/search may span multiple tenants ONLY when the caller supplies an explicit `Authorized Tenant Set`.
2. **Subset Inclusion Invariant**: Every queried tenant in `requested_tenants` MUST satisfy:
   $$\text{requested\_tenants} \subseteq \text{authorized\_tenant\_set}$$
3. **Fail-Closed Evaluation**: If even one requested tenant falls outside the caller's authorized set, the query fails closed with `unauthorized-tenant-in-set`. No partial multi-tenant result is leaked.
4. **Partitioned Aggregation**: Aggregate responses retain individual `tenant_id` tags per item. Assets from distinct tenants are never blended or deduplicated across tenant boundaries.

### 4.3 Elimination of Existence Oracles (Anti-Leakage)
- When a caller requests an asset ID that belongs to Tenant B while scoped to Tenant A:
  - The registry MUST NOT disclose that the asset exists in Tenant B.
  - The registry MUST return a uniform error response: `asset-not-found` (HTTP 404 equivalent).
  - Error messages, timing channels, and correlation responses MUST be indistinguishable from a non-existent asset ID.

---

## 5. Projection and Pagination Contracts

In strict conformance with `cmdr-product-spec/17-implementation-contracts/pagination-filtering-sorting.md`:

### 5.1 Opaque Cursor Specification
Pagination uses deterministic, server-generated, opaque cursor tokens. The decoded cursor payload consists of:
```json
{
  "version": 1,
  "tenant_id": "tenant-corp-prod",
  "sort_field": "updated_at",
  "sort_direction": "desc",
  "last_value": "2026-10-15T12:00:00Z",
  "tiebreaker_id": "skl-01j9a8b7c6d5e4f3",
  "page_size": 25,
  "query_fingerprint": "a3f5b78c..."
}
```

### 5.2 Deterministic Ordering Invariant
- To prevent duplicate entries, skipped rows, or pagination drift across updates:
  - **Primary Sort**: User-selected field (`name`, `updated_at`, `risk_classification`).
  - **Mandatory Secondary Sort (Tiebreaker)**: Canonical asset identifier `id ASC`.
- The compound sort key `(sort_field, id)` guarantees strict total ordering across all pages.

### 5.3 Bounded Page Limits
- Default page size: `limit = 25`.
- Minimum allowed: `limit = 1`.
- Maximum allowed: `limit = 100`.
- Limits exceeding `100` are clamped to `100` or rejected with `page-limit-exceeded`.

### 5.4 Standard Catalog Page Response Envelope
```json
{
  "items": [
    {
      "id": "skl-01j9a8b7c6d5e4f3",
      "tenant_id": "tenant-corp-prod",
      "type": "skill",
      "name": "triage-suspicious-login",
      "display_name": "Triage Suspicious Login",
      "description": "Deterministic and AI-assisted triage of suspicious authentication events.",
      "version": "1.2.0",
      "lifecycle_state": "published",
      "functional_state": "ready",
      "risk_classification": "low",
      "is_ai_powered": true,
      "has_non_ai_fallback": true,
      "owner": {
        "owner_id": "usr-soc-lead",
        "owner_name": "SOC Lead Analyst",
        "owner_type": "user",
        "is_unknown": false
      },
      "updated_at": "2026-10-15T12:00:00Z"
    }
  ],
  "pagination": {
    "limit": 25,
    "has_more": true,
    "cursor_next": "eyJ2ZXJzaW9uIjoxLCJzb3J0X2ZpZWxkIjoidXBkYXRlZF9hdCIs...",
    "cursor_prev": null,
    "total_count": 142,
    "is_total_exact": true
  }
}
```

---

## 6. Permission-Aware Filtering & Projection Enforcement

In strict conformance with `cmdr-product-spec/17-implementation-contracts/permission-enforcement-contract.md` and `REQ-PROD-006`:

### 6.1 Server-Side Pre-Authorization Evaluation
- Authorization decisions are evaluated on the server **prior** to projecting asset data.
- Client headers, query parameters, or token assertions cannot override server-evaluated policy.

### 6.2 Dual-Permission Namespace Resolution
The registry resolves both canonical and legacy permission tokens identically:

| Operation | Canonical Permission | Legacy Permission | Fail-Closed Error |
|---|---|---|---|
| Catalog Query / Browse | `perm.studio.library.read` | `perm.cmdr-studio.library.read` | `permission-denied` |
| Skill Metadata Inspect | `perm.studio.skill.read` | `perm.cmdr-studio.skill.read` | Redacted or `permission-filtered` |
| Workflow Metadata Inspect | `perm.studio.workflow.read` | `perm.cmdr-studio.workflow.read` | Redacted or `permission-filtered` |
| Automation Agent Inspect | `perm.studio.automation-agent.read` | `perm.cmdr-studio.automation-agent.read` | Redacted or `permission-filtered` |
| Agent Team Inspect | `perm.studio.agent-team.read` | `perm.cmdr-studio.agent-team.read` | Redacted or `permission-filtered` |
| Tool Contract Inspect | `perm.studio.tool.read` | `perm.cmdr-studio.tool.read` | Redacted or `permission-filtered` |

### 6.3 Projection Modes: Invisible vs Permission-Filtered
1. **Catalog Browse/Search Listing**:
   - Assets for which the caller lacks type-level read permission are **silently excluded** from browse results.
   - Total counts reflect only assets the caller is authorized to view.
2. **Direct Identifier Inspection (`Inspect`)**:
   - If caller requests inspection of an asset where they hold library read but lack asset-specific read:
     - The registry returns a `permission-filtered` stub.
     - Payload displays only public identifiers: `id`, `type`, `display_name`, `"functional_state": "permission-filtered"`.
     - Sensitive fields are redacted: `description: null`, `parameters: []`, `dependencies: []`, `secret_references: []`.

### 6.4 Invariant: Visibility Does Not Imply Executability (`REQ-PROD-006`)
- Every catalog projection contains execution metadata purely for informational classification (`execution_class`, `required_privileges`, `risk_classification`).
- The catalog projection interface permanently lacks invoke endpoints, trigger handlers, or execution dispatchers.
- Any attempt to send execution verbs to the catalog registry endpoint is rejected with HTTP 405 / `execution-not-supported`.

---

## 7. Metadata, Functional States & Owner Boundary Preservation

### 7.1 Asset Kinds
The registry supports five canonical Studio asset kinds:
1. `skill`: Reusable capability contract (`OBJ-SKILL`, `CAP-STD-010`).
2. `workflow`: Step orchestration and decision graph (`OBJ-WORKFLOW`, `CAP-STD-017`).
3. `automation-agent`: Autonomous agent definition (`OBJ-AUTOMATION_AGENT`, `CAP-STD-034`).
4. `agent-team`: Multi-agent collaborative team definition (`OBJ-AGENT_TEAM`, `CAP-STD-037`).
5. `tool`: Executable integration and capability contract (`CAP-STD-003`).

### 7.2 Lifecycle States to Functional States Projection Matrix

Source assets maintain their domain lifecycle states independently. The Catalog Registry projects them into four deterministic functional states:

| Source Lifecycle State | Dependency Health | Caller Authorization | Projected Functional State | Catalog Action Available |
|---|---|---|---|---|
| `published` | All healthy | Fully authorized | **`ready`** | Browse, Search, Inspect, Open, Clone, Compare |
| `published` | Degraded / missing optional | Fully authorized | **`partial`** | Inspect (warning displayed), Open |
| `approved` | Unvalidated / staging | Fully authorized | **`partial`** | Inspect, Open (preview only) |
| `draft` / `review` | Author private | Fully authorized | **`ready`** (author only) | Inspect, Open author surface |
| `deprecated` | Any | Fully authorized | **`stale`** | Inspect, Compare (warning: deprecation notice) |
| `superseded` | Any | Fully authorized | **`stale`** | Inspect, View superseding link |
| Any | Any | Library read only (no asset read) | **`permission-filtered`** | Stub inspect only; redacted fields |

### 7.3 Explicit `owner unknown` Enforcement (`CAP-STD-002` §8, §25)
- Owner information is structured as:
  ```json
  {
    "owner_id": null,
    "owner_name": null,
    "owner_type": "unknown",
    "owner_display": "owner unknown",
    "is_unknown": true
  }
  ```
- **Strict Invariant**: When an asset record lacks an owner attribute, or the designated owner is unregistered/deleted, the catalog MUST project `owner unknown`.
- The system MUST NEVER fabricate, guess, assign a default owner, or attribute ownership to the requesting user.

---

## 8. Standard Error Envelopes & Zero Credential Leakage

In strict conformance with `cmdr-product-spec/17-implementation-contracts/error-contract.md`:

### 8.1 Uniform Error Envelope Schema
Every error generated by the catalog registry conforms to the canonical envelope:
```json
{
  "error": {
    "code": "cross-tenant",
    "message": "Access to requested resource across tenant boundaries is forbidden",
    "correlation_id": "corr-8f3a9e2d-4b1c-4f7e-a9d2-3c1e5b7a9f0e",
    "retryable": false,
    "partial_result": null
  }
}
```

### 8.2 Catalog Error Code Registry

| Error Code | HTTP Status | Retryable | Description |
|---|---|---|---|
| `missing-tenant` | 400 Bad Request | No | Request omitted mandatory `tenant_ref`. |
| `wildcard-tenant` | 400 Bad Request | No | Wildcard tenant selector was rejected. |
| `cross-tenant` | 404 Not Found / 403 Forbidden | No | Resource belongs to a different tenant. |
| `unauthorized-tenant-in-set` | 403 Forbidden | No | Requested tenant not in caller's Authorized Tenant Set. |
| `cross-tenant-mutation-forbidden` | 403 Forbidden | No | Attempted cross-tenant cloning or copying. |
| `missing-asset-id` | 400 Bad Request | No | Operation requires `asset_id` parameter. |
| `asset-not-found` | 404 Not Found | No | Asset does not exist in the requested tenant scope. |
| `permission-denied` | 403 Forbidden | No | Caller lacks required catalog or asset read permissions. |
| `invalid-cursor` | 400 Bad Request | No | Pagination cursor is expired, malformed, or tampered. |
| `page-limit-exceeded` | 400 Bad Request | No | Requested page limit exceeds maximum (100). |
| `unsupported-grouping-facet` | 400 Bad Request | No | Group facet is not supported by the catalog registry. |
| `stale-dependency` | 409 Conflict | No | Referenced dependency or version is invalid/stale. |
| `service-degraded` | 503 Service Unavailable | Yes | Upstream Shared Search degraded; fallback available. |

### 8.3 Zero Credential Leakage Guarantee
1. **Opaque Secret References**: Secret assets reference credentials via opaque keys (`sec-ref-ip-reputation-token`). The catalog registry never contacts secret vaults or resolves raw API keys.
2. **Error Scrubbing**: Error messages and stack traces are sanitized before transmission. No database connection strings, server IPs, filesystem paths, or internal tokens are ever reflected.
3. **Correlation Tracking**: Every error includes a unique `correlation_id` (UUID v4) for forensic tracking without leaking internal system state.

---

## 9. Verification Status & Governance Rules

1. **Preimplementation Contract Status**: This evidence document establishes the canonical CMDR Studio Library Asset Catalog and Query Registry contract under `E10-STD-001B-CATALOG`.
2. **Production Code Neutrality**:
   - Production runtime files added: `0`
   - Net production LOC: `0`
3. **Product Spec Immutability**:
   - `git diff 82e630fce3a7f6236e03709116eca6dd0a1acd40 HEAD -- cmdr-product-spec` asserts strictly **`0 bytes`** modification.
4. **Toolchain Verification**:
   - `cmdr-dev validate-manifests`: **PASS** (201/201 manifests, strict v2 valid).
   - `cmdr-dev architecture-audit`: **PASS** (15/15 boundaries valid, 1466 path claims).
   - `cmdr-dev doctor`: **PASS** (`healthy=true`).
5. **Test Suite Non-Regression**:
   - Node.js test suites: **143/143 tests pass** (61 in timeline demo + 82 in event-search frontend).
   - Go modules: **11/11 modules vet clean**.
6. **Graph Promotion**: Promotes work unit `E10-STD-001B-CATALOG` from `READY` to `IMPLEMENTED`.
