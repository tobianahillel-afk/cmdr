---
id: investigate-detection-shadow-canary
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
---
# Shadow and canary

CAP-INV-422 defines non-alerting runtime observation. Shadow matches never create operational Signals. CAP-INV-423 defines bounded canary populations, stages, success/stop criteria and rollback.

Shadow ≠ historical replay. Canary success ≠ global success. Stage advancement, stop and rollback remain governed actions, never silent automation.
