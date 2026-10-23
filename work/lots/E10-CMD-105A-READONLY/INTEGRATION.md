# E10-CMD-105A-READONLY — CAP-CMD-105 selected operational SLA (Class 0 only)

## Authority and four identities

- **Canonical**: CAP-CMD-105 SLA Tracking (§§1–27), consumer CAP-CMD-101 Unified Work Queue, Command Incidents and Work Queue module. ADR-0008 keeps operational SLA in Command, contractual SLA as optional externally authorized source projection. OPEN-013 blocks Class-2 SLA pause/resume and follow-up writes.
- **Engineering**: new E10-CMD-105 epic, one child E10-CMD-105A-READONLY, depends on existing E10-CMD-101F-SELECTION IMPLEMENTED. Manifest lists exact allowed files and security acceptance.
- **Code**: existing first-party `product-runtime/work-queue` Go server, one `GET /command/work-queue` with the same six system views, native HTML template.
- **Git**: child branch of frozen PR #173 SHA `89a092bade8b2137732ac5da8583a1d321249549`; PR parent is NOT the engineering parent. No mutation of parent SHA. DRAFT only.

## Explicit canonical delta

| Canonical behavior | Before this PR | New demonstrated subset | Still absent |
|---|---|---|---|
| Inspect readable work item's SLA | Boolean `SLAAtRisk` and separated CAP-CMD-104 risk dimension, no state explanation | One optional selected-only host `ReadAuthorizedSLA` claim, serializes independently sourced state, policy version/source, observed/due timestamps and source-owned remaining seconds | Full policy/calendar/pauses computation, policy provider deployment, audit, integrated identity platform |
| Maintain owner permissions | Per object `Owner.CanRead` and selected-page ACL | NEW SLA-owner authorization explicitly independent; called once only after selected object current-page CanRead; denied/off-page no calls | Production IAM and policy-level entitlements implementation |
| Unknown and errors | No detailed SLA | Absent policy = unknown, malformed/unavailable claim = generic 503 with no selected row/metadata | Full partial field source coverage and renewal request |
| Mutations | None | None | Class-2 pause/resume, Task follow-up, SLA event, audited transitions require OPEN-013 and actual owner services |

Host must enforce SLA policy or contract source permissions independently of the selected work-item permission; the new optional adapter contract is NOT an authorization engine. The server validates the bounded source claim at the same selected-item HTTP trust boundary. Do not infer priority, expiration, calendar, pause permission, contractual customer relationship, remaining time or SLA threshold from an untrusted client or from the already owner-provided risk bit. The canonical SLA status remains a host-source projection, not a new calculation engine; without adapter no claim is returned and native UI explicitly says unknown.

## Validation / integration

- Selected authorized `Incident:inc-a` in the existing six-view queue can read one independently source-authorized operational SLA with state `at-risk`, policy version, source, observed/due time and explicit zero remaining seconds.
- Unselected listings remain byte-contract compatible (no selected SLA field/provider request). Foreign tenant, denied, missing or off-current-page selections never invoke SLA policy reader. Invalid, malformed or unavailable provider response returns generic error with no selected object or source metadata; non-existent provider/nil payload remains unknown.
- No new HTTP route, third-party runtime, provider mutation, new ACL model, new Customer object or paid API; native escaped HTML only.
- Compare exact-HEAD Linux/macOS/Windows workflow and mandatory engineering-kernel security changed-line 90% floor, SAST and work complexity. No merge/claim of production readiness from child CI alone, pending reviews/approvals of long parent PR stack.
