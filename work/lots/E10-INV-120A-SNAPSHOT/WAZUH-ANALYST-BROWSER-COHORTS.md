# E10-INV-120A — analyst Wazuh agent/rule pivot, end to end

## Functional result

When starting `case-timeline-demo -snapshot=<absolute owner JSON> -wazuh-alerts=<absolute Wazuh 4.x JSONL> -into-case=<authorized Case> -serve-import`, the imported events are already attached to the existing in-memory first-party owner and shown in the browser. A new analyst button invokes `GET /demo-wazuh-cohorts?case_id=<Case>` under the existing ephemeral Authorization bearer. It aggregates the real imported source's `agent.id`, `rule.id`, counts, distinct rules per agent and highest known source-reported severity, then shows top-ten cohorts with browser `textContent`.

This is not a second Wazuh parser, search engine or tenant authority. `writeDemoWazuh4RuleTriageWithData` reuses the same offline CLI cohort calculator and the existing authenticated Case/Timeline/typed-source HTTP traversal before releasing data. Wazuh input bytes come from the SAME bounded local file read already used for import; no reread of a mutable input path, no owner snapshot export, no external backend.

The optional route is installed only for Wazuh viewer mode, never for ordinary Case views. A wrong bearer returns 401, wrong/multiple Case IDs return generic 404, POST is refused, and cancelled/denied owner reads produce no cohort JSON. The owner data is static, copied at construction; each request rechecks the existing in-memory Case, Timeline and source permissions, but no remote IAM or live revocation provider exists. The browser discards stale responses on bearer or Case/filter changes, and never uses HTML interpretation for Wazuh IDs.

## Simplicity and scope

Reuse Go standard-library `net/http`, `encoding/json`, existing demo owner/ACL chain and first-party Wazuh aggregation; reuse the embedded analyst UI rather than introducing React, a new SIEM store, index, parser SDK or persistence. All are unnecessary for a bounded local 8 MiB/5000-entry offline laboratory journey. Source `agent.id` is not a verified host identity, rule levels are descriptive not a causal or incident score, and `OPEN-013` remains unresolved.

## Verification and risk

Functional tests use a local operator owner snapshot, actual Wazuh JSONL, bearer-protected `httptest` calls, multiple Case rejection, no route for ordinary sessions, cancellation, confidential source text exclusion, and JavaScript browser VM tests for success, denial, malformed data and changed bearer. CI exact-SHA and independent L2 are required; no release-SLO, production IAM, SIEM integration, forensic custody or main-branch delivery is claimed.
