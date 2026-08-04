---
id: dynamic-sandbox-studio-boundaries
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-AI-002
  - REQ-OBJ-009
open_decisions:
  - OPEN-015
---
# Studio boundaries

Studio possède Tool, Tool Call, Skill, Workflow, Automation Agent, Automation Run, Human Gate, versions et évaluations. Investigate possède Dynamic Analysis Session, le contexte Sandbox Run, les observations et la disposition analytique. Un Tool Call n’est pas une analyse complète ; une Automation Run n’est ni un Sandbox Run ni un Response Run.
