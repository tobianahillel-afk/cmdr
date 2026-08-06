---
id: investigate-cloud-analysis-shared-capabilities
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements: [REQ-PROD-019, REQ-PROD-020]
open_decisions: [OPEN-014, OPEN-015]
---
# Shared Capabilities used by Cloud Analysis

Cloud Analysis consumes, but never redefines:
- Entity and Entity Resolution;
- Graph and relationship rendering;
- Timeline and time navigation;
- Search and filtering;
- Object Linking;
- Background Jobs;
- Notifications;
- Trace and Activity;
- Versioning, Comparison and Supersession;
- Export and Reporting;
- Collaboration;
- Recovery.

A Cloud Resource Observation is not a canonical Entity. A Cloud Timeline is a local analytical timeline rendered through Shared, not the Case Timeline. Export does not extend source permissions. Shared Jobs do not grant provider access or execute target mutations.
