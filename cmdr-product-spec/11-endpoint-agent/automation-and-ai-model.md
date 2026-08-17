---
id: endpoint-ept1-automation-and-ai-model
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-11
source-of-truth: canonical
---
# Endpoint EPT-1 Automation and AI Model

EPT-1 is deterministic/manual first. Every essential function has a non-AI path using source facts, deterministic comparisons or human review.

## Optional AI
AI may summarize sourced health, explain stale/degraded states, summarize inventory changes or suggest an anomaly candidate for later analyst review. AI output is attributed and carries source/provenance/uncertainty.

## Prohibited AI behavior
AI cannot invent Agent identity, platform/OS facts, support status, heartbeat, capability availability or provenance; cannot enroll autonomously, grant permission, alter Fleet/Policy, hide offline/stale state, promote technical facts into Evidence/Finding/Decision/Result, or create cross-tenant access.

## Deterministic guarantees
Unknown stays unknown; missing source never becomes a positive assertion; platform observed != supported; advertised != authorized; assigned Policy != applied; state derivation preserves source and freshness.

CMDR Studio owns agentic capabilities. Endpoint Agent is not a Studio Automation Agent.