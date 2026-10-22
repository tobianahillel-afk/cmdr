---
id: investigate-detection-promotion-deployment
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
---
# Promotion and deployment coordination

CAP-INV-421 plans target order and conditions. CAP-INV-424 observes Govern-owned Response Runs and Results. Platform Settings owns environments, targets, configured runtimes and administrative health.

Promotion Plan ≠ execution; deployed ≠ active; active ≠ healthy. Every target retains its own state, error and version. Investigate never directly administers the runtime.
