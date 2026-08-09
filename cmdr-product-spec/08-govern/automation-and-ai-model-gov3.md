---
id: govern-automation-and-ai-model-gov3
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical-addendum
---
# Automation and AI Model — GOV-3 Addendum

GOV-3 essential behavior remains available without AI.

| Fonction | Humain | Déterministe | Automatisable | IA | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| event classification/correlation | oui | oui | oui | candidate explanation | source/event tables and IDs |
| audit reconstruction | oui | relation traversal/diffs | oui | sourced summary | reconstruction tables |
| completeness/gap detection | oui | expected-vs-observed rules | oui | candidate explanation | completeness matrix |
| metric aggregation | oui | oui | oui | explanation | pivots/tables |
| trend comparison | oui | oui | oui | summary/hypothesis | comparison tables |
| Control Health Assessment | accountable reviewer | comparisons support | workflow possible | draft only | structured review form |
| improvement package | accountable reviewer | owner mapping/checks | workflow possible | sourced draft | template/checklist |
| closure checks | QA/reviewer | registry/gate checks | oui | summary only | validation matrices |

## Prohibited AI behavior
AI cannot invent missing events, delete/alter history, declare fraud or policy violation as established fact, mutate Decision/Result, hide contradictions, change thresholds without human governance, publish audit content externally, expand tenant access or apply a Continuous Improvement Package.

Every AI output is attributed, source-linked and distinguishable from audit facts, metric observations and human assessments.