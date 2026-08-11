---
id: endpoint-ept3-action-classification
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-11
source-of-truth: canonical
---
# EPT-3 Action Classification

| Action family | Class | Owner/authority | EPT-3 rule |
|---|---:|---|---|
| inspect signal/context/coverage/timeline/provenance | 0 | source read permission | read-only |
| deterministic eligibility/evaluation/match/group/context/timeline/gap assessment | 1 | Endpoint local technical semantics | no side effect |
| acknowledge local workflow state | 2 | Endpoint where source allows | local state only |
| bounded contextual refresh | 2 | Endpoint where refresh uses already available data | no acquisition |
| suppression-state interaction | 2 request only | canonical suppression/config owner | no implicit policy mutation |
| Collection/Live Response/containment/response | 3/4 future | Investigate/Govern/Endpoint future lots | **not in EPT-3** |

EPT-3 contains no response action. Detection evaluation is never production-response authority.