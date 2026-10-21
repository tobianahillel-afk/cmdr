# E10-CMD-106E-HOST — Same-origin Command → Investigate read-only host

## Product gap and real output

The prior Command CAP-CMD-106D PR #195 produced a source-authorized HTML link from a selected Incident to `/case-timeline?case_id=...`. Its handoff explicitly states a standalone Work Queue host returned 404 because the independent Investigate route was not mounted. This new bounded integration closes **that specific, demonstrable navigation gap**: two **existing** runtime `http.Handler` implementations are mounted on exactly `/command/work-queue` and `/case-timeline` under one HTTP origin. The Case selector conveyed by Command is still only a read pointer; no source authority/permission transfer occurs.

## Four owners

- Canonical pilot: Command CAP-CMD-106 (Incident → authorized linked Investigate Case navigation) and CAP-CMD-101 Work Queue, consuming Investigate CAP-INV-110 (read-only Timeline). Incident remains Command; Case remains Investigate; Timeline Entry remains Shared; relation remains Shared Object Linking.
- Engineering parent: EXISTING `E10-CMD-106` epic, new child `E10-CMD-106E-HOST`, depends E10-CMD-106D-NAVIGATION IMPLEMENTED and E10-INV-110E VERIFIED. No fictional new product capacity or transferred owner. Existing graph has zero READY tasks; this task was selected from the specifically documented #195 host gap rather than blindly creating a new epic.
- Code owner: NEW, explicitly registered **cross-product routing composition** `command-investigate-host-runtime` at `product-runtime/command-investigate-host/**`, allowed to depend on existing Command Work Queue and Investigate Timeline runtime boundaries. Neither source implementation nor canonical contract is edited; Go modules both pinned to 1.26.10.
- Git parent: #197 `87bcb604fd388d26f4919d5c6bb6c3c489528bda` (security-updated Go 1.26.10 stack), DRAFT; this PR is a child branch and does NOT modify the L2-frozen #197 SHA. This Git parent differs from the product/engineering parent.

## Concrete verification and trust boundary

Go integration tests construct the REAL first-party `workqueue.NewHTTPHandler` and `investigationtimeline.NewMemoryCaseTimelineHTTPHandler` against independently constructed synthetic owner services, mount them through `commandinvestigatehost.NewHandler`, request a selected Command Incident HTML page, extract/verify its preexisting Case href, and issue the **subsequent GET through the same host**. An Investigate-authorized analyst sees the Case Timeline entry and source-only metadata; a Command Incident-only reader can see the already independently source-authorized Case link but receives 404 from Investigate owner validation on navigation. Invalid token, unsupported route, forged tenant query, mutation and missing handler fail. GitHub executes Go tests, vet, Linux race and >=90% host statement coverage across Linux/macOS/Windows; engineering kernel gates remain mandatory.

The router deliberately owns NO IAM, tenant selector, source access data, filesystem or database. It does NOT read or forward principal/tenant/environment from user query fields and does not consume the link as authorization. The two handlers remain separate and re-evaluate owners per request. No login, TLS termination, identity provider, production data source, deployed server binary, native browser renderer, return-to-Command context, durable audit or Class-2 Case/Incident mutation is delivered. Synthetic Go owner adapters **are not** production IAM or actual Wazuh/Defender source connectors. This milestone demonstrates combined **in-process HTTP functionality**, not production acceptance.

## Handoff

Check exact PR HEAD and base, diff, Go 1.26.10 module replacement and architecture/security registrations, host test and workflow, and global engineering checks. Treat native GitHub review approval, prior L2 results and merged ancestors separately. If CI reveals extra required registry changes, correct the new branch only and keep old SHA frozen. Next product work should prioritize actual trusted IAM/Shared source adapters and deployed application host when authority and implementation architecture are accepted; do not fabricate OPEN-013 approval.
