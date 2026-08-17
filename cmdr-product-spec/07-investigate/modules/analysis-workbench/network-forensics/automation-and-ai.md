---
id: investigate-network-forensics-automation-and-ai
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
requirements:
  - REQ-AI-001
  - REQ-AI-002
  - REQ-AI-003
  - REQ-AI-006
open_decisions:
  - OPEN-013
  - OPEN-015
---
# Automation and AI

AI peut proposer protocoles, request/response relations, résumés, groupements, Entity candidates, périodicités, anomalies, corrélations, Hypotheses et packages candidats.

Toutes les fonctions essentielles restent disponibles avec parsers/décodeurs déterministes, viewers, tables, timelines, graphes avec alternative tabulaire, filtres, recherche, comparateurs, règles explicables, extracteurs, checklists, workflows non agentiques et revue humaine.

Interdictions : Tool ou protocole silencieux, fusion Entity, replay, interaction cible, payload révélé sans permission, déchiffrement non autorisé, anomalie/IOC/Evidence/Finding/règle/objet Intelligence automatique, auto-permission, suppression de trace ou chatbot obligatoire.

Toute sortie automatisée expose initiateur, producteur/version, Tool Calls, Automation Run, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et disposition acceptée, modifiée ou rejetée.
