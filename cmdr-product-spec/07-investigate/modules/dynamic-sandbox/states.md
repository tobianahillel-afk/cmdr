---
id: dynamic-sandbox-states
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-014
---
# Functional states

| Scope | States |
|---|---|
| Dynamic Analysis Intake | draft, incomplete, ready, blocked, environment-unavailable, policy-blocked, unsupported |
| Dynamic Analysis Session | draft, ready, active, paused, blocked, partial, completed, failed, archived, superseded |
| Sandbox Run | validating, queued, preparing, starting, running, idle, stopping, collecting-results, cleaning, resetting, partial, completed, failed, cancelled, timed-out, unsafe-environment |
| Runtime Artifact | proposed, capturing, available, partial, invalid, restricted, superseded, withdrawn-from-use |
| Reproducibility | reproduced, partially-reproduced, not-reproduced, missing-environment, missing-version, missing-input, policy-blocked, disputed |

`queued`, `started`, `Artifact launched`, `completed` et `environment reset` restent des étapes distinctes. Les machines finales sont reportées.
