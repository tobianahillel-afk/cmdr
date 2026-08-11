---
id: endpoint-ept3-cross-product-links
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-11
source-of-truth: canonical
---
# EPT-3 Cross-Product Links

- **Investigate:** owns Detection Engineering authoring/validation/lifecycle plus Case/Evidence/Finding and analyst qualification. Endpoint consumes eligible content and forwards technical context/summary only.
- **Command:** owns canonical Detection/Signal/Alert/Incident. Endpoint Local Detection Match/Signal Candidate are local technical concepts, not competing Command objects.
- **Platform Settings:** owns sources/parsers/configured runtimes/targets/Endpoint Policy/assignments and administrative rollout configuration. Endpoint reads projections only.
- **Govern:** owns Action Request/Decision/Approval/Response Run/Result and response authority. EPT-3 executes no response.
- **Studio:** owns Tool/Tool Call/Skill/Workflow/Automation Run. Endpoint capability/evaluation is not a Tool/Run.
- **Shared:** owns generic Search/Timeline/Linking/Correlation/Trace/Activity/Jobs/Reporting. Endpoint creates only local contextual projections.

All links are permission-aware and tenant-scoped. Navigation/handoff never transfers ownership or permission.