---
id: endpoint-ept3-investigation-context-map
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-11
source-of-truth: canonical
---
# EPT-3 Investigation Context Map

`CAP-EPT-037..045` define local read-only context over already available observations:

Process Tree → File / Network / User-Session / System Context → Endpoint Contextual Timeline → permission-aware Pivots → Detection-to-Investigation Context Expansion → Endpoint Investigation Summary.

Every context preserves snapshot/freshness, source refs, tenant, permission, gaps and provenance. Process ancestry and correlations are relationships, not causal or malicious conclusions. If a pivot requires new file bytes, memory, packet data or other acquisition, EPT-3 returns `collection-required` and stops before EPT-4.

This map creates no Screen ID, wireframe, button, query language, remote shell or Collection workflow.