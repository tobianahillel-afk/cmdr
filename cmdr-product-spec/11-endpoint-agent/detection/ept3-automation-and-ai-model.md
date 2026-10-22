---
id: endpoint-ept3-automation-ai-model
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-11
source-of-truth: canonical
---
# EPT-3 Automation and AI Model

Essential EPT-3 functions have a deterministic/manual path. AI is optional.

AI may explain a Local Detection Signal Candidate, summarize context/timeline, propose permission-safe pivots, explain coverage gaps, identify related observations as suggestions, and propose an explicitly labelled investigation hypothesis.

AI must not declare maliciousness with authority, create Evidence/Finding/Case/Incident/Decision/Result, execute response, hide/suppress a signal, remove provenance, invent observations/context, claim causality as fact, reveal restricted content or close `OPEN-017`.

| Function | Deterministic/manual path | AI optional use | Forbidden promotion |
|---|---|---|---|
| eligibility/evaluation/match | content conditions + local facts | explanation only | no Finding/Incident |
| context linking/timeline | stable refs/time rules | suggested links/summary | no causal proof |
| investigation summary | structured refs/limits | narrative summary | no Evidence/Finding |
| gap explanation | dependency/reason codes | explanation | no invented visibility |

All AI outputs preserve source refs, uncertainty and attribution.