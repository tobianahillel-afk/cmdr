# E10-INV-101O-INCIDENT — Command Incident → authorized Case Queue

## Four ownership layers

- **Canonical behavior:** `CAP-INV-101` Case Queue §17, §21, §25; REQ-PROD-014/020; INV-CAS-001. Consumes the Command-owned `Incident` source projection and Investigate-owned `Case` in the Case Queue. Saved Views stay Shared-owned.
- **Engineering unit:** `E10-INV-101O-INCIDENT` under already existing `E10-INV-101` epic, dependent on `E10-INV-101N-RETURN`, which remains IMPLEMENTED pending native approvals.
- **Code owner:** `product-runtime/case-queue` Go. NO new separate Command relation service, projection runtime, IAM provider, Saved View provider or storage.
- **Git:** new child branch based on frozen PR #145 exact SHA `d31c326a785134e06b86fa79c9d8701d0a287a27`, NOT main; parent SHA unchanged.

## Requirement → existing → integrated → remaining

| Requirement | Existing runtime | This delivered behavior | Remaining |
|---|---|---|---|
| CAP-INV-101 §17 Command Incident → linked Cases in Queue/Workspace | First-party Queue handles one optional independently ACL-authorized Incident projection, but no `incident_id` selector | Optional `incident_id=EXACT_OPAQUE_ID` selects Case rows **only** when the projected current Command Incident was authorized; keeps old filters, Saved View and Case/Incident owner checks | Canonical many-to-many typed `linked-incidents` graph/provider, real Command Incident Detail integration and host IAM |
| CAP-INV-101 §§17/25 filtered view, selection, scroll preserved | #145 added current fully authorized paged Queue → Workspace → POST return | Existing filter validator and projector preserve opaque Incident ID across current pagination cursor, native handoff and reauthorized return | Deployed Case UI/Command link, production persistence, audits and approvals |
| REQ-PROD-020 no cross-tenant/unauthorized leaks | Existing per-request Case/Incident/Finding source ACL and Shared Saved View version policy | Filter runs **after** source owner checks and only over visible `row.Case.Incident.ID`; denied/unavailable source cannot match. Revocation after handoff makes return fail closed. | Independent live provider/IAM, policy/CSRF/TLS implementation |

## Acceptance

Real Go `httptest` (source owners injected by host, no network):
- Visible Case matches exactly the command incident ID. A hidden Case, one with denied Incident, unrelated Incident and a Case with none do not match or disclose source text. Incident ID comparison preserves case and whitespace and does not normalize identity.
- Link-scoped page 1/page 2 → open selected Case C on its present authorized second page → `POST /case-queue/return` restores exact `incident_id`, sort, page cursor, scroll and selection. Revoking the Command Incident read right before returning refuses access.
- Empty, duplicated, tenant/principal injection, control characters and too-long selectors are rejected before the session or provider is invoked.
- Existing Queue, Shared Saved View, paging, case-owner and return suites must remain green on exact HEAD.

## Security and limitations

A client can supply an Incident ID filter but cannot supply an Incident projection or change the host's authenticated scope. This does **not** imply that an Incident belongs to a Case if the host's candidate object graph is untrusted; authoritative object relations, provenance and tenant enforcement remain the responsibility of the host-owned Command and Investigate providers. No new permissions transferred by relationships, no mutation of `OPEN-013`, no work queue aggregation, no signed link token, no production deployment or product acceptance. Security-critical changed-code 90% floor and all GitHub review requirements remain mandatory.

## Handoff / recovery

After GitHub Actions, record the exact `HEAD`, `BASE`, commits/diff, all required CI runs including failure/pass, and one risk-based read-only L2. Until independent L2 and native reviews/merges, work unit status is at most IMPLEMENTED, never VERIFIED/ACCEPTED. Do not modify previous L2-frozen SHAs or claim parent/main merged.
