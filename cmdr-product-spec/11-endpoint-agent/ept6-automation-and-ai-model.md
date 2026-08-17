---
id: endpoint-ept6-automation-and-ai-model
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-12
source-of-truth: canonical
---
# EPT-6 Automation and AI Model

AI is optional for every EPT-6 path. Essential behavior remains manual or deterministic.

| Function | Deterministic/manual baseline | AI may | AI must never |
|---|---|---|---|
| Update failure/health | inspect events, compare expected/observed | summarize/explain | choose arbitrary version, invent success, hide failure |
| Compatibility/readiness | source-backed rules | explain blockers | declare package trusted or platform supported |
| Recovery/replay | state/ref reconciliation | suggest sourced recovery candidate | invent recovery, buffered data, acknowledgements or exactly-once |
| Degraded operation | resource/dependency rules | summarize causes | hide load shedding/loss or declare full health |
| Self-protection/tamper | expected/observed comparison | explain candidate context | declare compromise or authorize response |
| Secret handling | reference validation/masking | summarize non-sensitive metadata | reveal, invent, transform or store raw secrets |
| Local audit/provenance | stable-ref traversal | summarize chain | alter/delete history or invent missing events |
| Closure | deterministic counts/gates | summarize evidence | declare Endpoint/Phase 5 PASS without gates |

No AI path may bypass Settings administration, Security permissions, Govern authority, tenant isolation, masking, OPEN decisions or source-of-truth ownership.