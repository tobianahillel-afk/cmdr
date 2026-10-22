---
id: capability-register-settings-health-and-slo
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-16
source-of-truth: registry
---
# Capability Register — Settings Platform Health and Source-Attributed SLO Projection

This shard registers the single source-audited Phase-6 Platform Settings capability defined after validated ADR-0009. It records documentary capability definition only and does not claim runtime implementation, monitoring, measurement acquisition, SLO calculation, failover, recovery or production availability.

| Capability ID | Name | Owner product/module | Status | Delivery status | Delivery mode | Canonical file | Primary roles | Primary objects | Consumers | Requirements | OPEN decisions | Dependencies | Supersedes | Review date |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| `CAP-SET-014` | Platform Health and Source-Attributed SLO Projection | Platform Settings / Health | draft | defined | planned | `../../10-platform-settings/capabilities/cap-set-014-platform-health-and-source-attributed-slo-projection.md` | Platform Administrator | Tenant; Environment (conditional); Integration; Data Source; Model Provider; Endpoint Agent Fleet | platform-settings; command; shared; security | REQ-PROD-003,005,006,008,009,010,012,013,019,021; REQ-SEC-001,002 | OPEN-008; OPEN-013; OPEN-015; OPEN-019 | ADR-0009; ADR-0008; CAP-SET-001; CAP-SET-002; CAP-SET-010; CAP-SET-012; CAP-CMD-401 | none | 2026-11-16 |

## Structural evidence
- allocation: exactly `CAP-SET-014`;
- owner: Platform Settings Product Lead;
- documentary state: `draft / defined / planned`;
- structure: **1 capability / 27 numbered sections / 6 mandatory substantive tables / 5 meaningful GWT**;
- primary existing screen: `SET-HLT-001`;
- existing permission: `perm.settings.health.read`;
- new canonical objects / Permission IDs / Screen IDs: **0 / 0 / 0**;
- writes: **none**;
- `CAP-SET-015+` remains unallocated and unreserved.

## Boundary
SLO remains source-attributed and non-canonical. Settings owns bounded deterministic Health/SLO projection and presentation only. Platform Architecture retains neutral Health/Metrics contract semantics; Shared retains generic metric/Search/Reporting/Export/Notification mechanisms; source/runtime owners retain measurement acquisition and authoritative source-specific calculation; Security retains Tenant isolation/Authorized Tenant Set; Command and Govern retain their existing operational/response ownership.

No row in this shard may be interpreted as implementation, deployment, source support or runtime availability.
