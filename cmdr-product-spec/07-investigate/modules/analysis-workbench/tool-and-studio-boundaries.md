---
id: analysis-workbench-studio-boundaries
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
# Tool and Studio boundaries

Studio possède Tool, Tool Call, Skill, Workflow, Automation Agent, Agent Team, Automation Run, Human Gate, versions et catalogue. Investigate sélectionne un Tool autorisé, transmet son contexte, voit version/statut/permissions/limites, consomme le résultat et conserve sa relation à l’Analysis Session. Investigate ne crée ni n’administre le catalogue et ne transforme pas une intégration en capacité native.
