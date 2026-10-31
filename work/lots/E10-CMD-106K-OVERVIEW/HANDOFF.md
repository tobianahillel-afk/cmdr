# E10-CMD-106K-OVERVIEW — Command Incident Detail canonical business header

## Canonical obligation and operational value

CAP-CMD-106 Incident Coordination §§1–27 and CMD-IWQ-006 Incident Detail §10 require a Command-owned Incident ID/title/priority/severity/business-impact. Prior #209 already delivers a dedicated read-only Class-0 Detail and independently authorized Case navigation, but it visibly **lacked this canonical Command business header**. New behavior provides those fields from a dedicated optional owner-authorized projection in the same existing Incident Detail HTML. It is a useful partial canonical record, **NOT the complete Incident Detail screen**: Alert source, next action, Decision, Run, Result, six UX states, mutation, Govern and OPEN-013 remain outside this task.

## Four identities

- **Product:** Command CAP-CMD-106 Incident Coordination; CAP-CMD-101 Work Queue selects; Investigate CAP-INV-110 Case stays independent. Incident object canonically owned by Command. Linked Case / Timeline remain Investigate/Shared owners, never permission by relation.
- **Engineering:** existing epic E10-CMD-106, task E10-CMD-106K-OVERVIEW, dependency J-DETAILREAD IMPLEMENTED on #209. No new product capability or permission invented.
- **Code:** existing first-party Go `product-runtime/work-queue/incident_detail.go` and new narrow `incident_overview.go` optional `IncidentOverviewOwner`; synthetic local browser owner `product-runtime/command-investigate-host/cmd/local-browser/main.go`. No new service or HTTP route.
- **Git:** child based strictly on historical/frozen PR #209 HEAD `9a81ca6c27fe674d63813f90bba908a911d3b917`, no change to #209 or upstream branches. Git base is NOT parent business capability.

## Trust boundary, expected behavior and exclusions

Existing session, tenant/env, full bounded Command Owner.List and per-Incident CanRead, and fresh selected Incident remain mandatory on EVERY Detail read. Only then a single optional `ReadAuthorizedIncidentOverview(ctx, principal, tenant, environment, incidentID)` is invoked. Returned Incident/tenant/environment must exactly equal authenticated selection; source reference and version and RFC3339Nano observed-at are required for any non-nil value, bounded title, priority, severity, impact must be valid UTF-8 without control bytes. On owner failure/mismatch, generic HTTP 503 occurs **before** HTML buffer emission. Missing optional provider/denied optional projection shows explicitly unavailable; missing individual fields are labelled unknown, never inferred. Source data rendered through Go html/template with original CSP, no-store and no-referrer. No JS, mutation, report, export, case right promotion, new endpoint, third-party component or API.

Synthetic 127.0.0.1 runtime supplies an explicitly invented *demo* title/priority/severity/impact only from its hard-coded **owner adapter** for analyst. Command-only commander can read the original Incident but does not receive the narrower business source projection; Case remains independently denied. Demo owner is not production IAM, tenant resolver, real incident storage, SIEM ingestion, alert provenance or acceptance.

## Evidence and acceptance scope

Go tests validate permitted third-page Incident header, source HTML escaping, field-level unknown, legacy owner without interface, revoke-before-source behavior, wrong tenant/env/Incident, excessive/control source data, malformed provenance/time, source error and no partial information. Existing synthetic local real Go HTTP test checks visible authorized overview and Command-only withheld. Required **exact SHA** checks: Work Queue Linux/macOS/Windows, same-origin Host Linux/macOS/Windows, engine-kernel graphs/manifests/security/coverage/go vet/Go race/SAST/SCA. Historical parent CI does NOT validate child SHA.

Merge prerequisites: native reviews, parent #209 and deeper Git ancestry integration, security controls, and independent L2 if material trust boundary assessment requires it. `E10-CMD-106K-OVERVIEW` is recorded IMPLEMENTED in graph only; no fabricated VERIFIED, ACCEPTED, production or main integration. Any missing or failing check remains a blocker; update this handoff only with confirmed results.
