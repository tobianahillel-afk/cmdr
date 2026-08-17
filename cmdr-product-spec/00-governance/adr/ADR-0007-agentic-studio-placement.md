---
id: ADR-0007-agentic-studio-placement
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-AI-001
  - REQ-AI-002
  - REQ-AI-003
  - REQ-AI-004
  - REQ-OBJ-009
---
# ADR-0007 — Placement des capacités agentiques dans CMDR Studio

## 1. Identifiant

`ADR-0007`

## 2. Titre

Placement des capacités agentiques dans CMDR Studio

## 3. Statut

Draft. Cette ADR n'est pas approuvée et ne remplace aucune décision source au-delà de ce qu'elle applique explicitement.

## 4. Date

2026-08-03

## 5. Propriétaire

Product Architecture.

## 6. Décideurs attendus

Head of Product et propriétaires des domaines affectés ; Security, UX ou Engineering selon les effets décrits.

## 7. Requirement IDs

`REQ-AI-001`, `REQ-AI-002`, `REQ-AI-003`, `REQ-AI-004`, `REQ-OBJ-009`

## 8. Contexte

Command, Investigate et Govern peuvent bénéficier d'agents sans créer trois moteurs, modèles ou gouvernances concurrents.

## 9. Problème

Si chaque produit possède ses agents, les permissions, outils, évaluations, coûts et traces divergent et l'IA peut devenir propriétaire implicite du workflow.

## 10. Forces en présence

Optionnalité de l'IA, responsabilité humaine, réutilisation, isolation tenant, assurance et expérience opérationnelle.

## 11. Options étudiées

1. Agents propres à chaque produit.
2. Studio comme produit opérationnel central.
3. Studio propriétaire des capacités agentiques, produits opérationnels consommateurs.

## 12. Décision

CMDR Studio possède Skill, Tool, Tool Call, Automation Agent, Agent Team, Workflow, Human Gate et Automation Run. Les autres produits déclenchent ou consomment ces capacités sans céder leurs objets ou décisions.

## 13. Justification

Un propriétaire unique permet assurance et observabilité communes tout en maintenant les workflows essentiels sans IA.

## 14. Conséquences positives

- Évaluations centralisées.
- Traçabilité des Tool Calls.
- Coûts et versions visibles.
- Pas de moteur concurrent.

## 15. Conséquences négatives

- Dépendance à Studio pour les capacités agentiques optionnelles.
- Contrats de contexte interproduit nécessaires.

## 16. Risques

- Studio devient l'interface principale.
- Human Gate confondue avec Decision.
- Agent s'auto-approuve.

## 17. Effets sur la navigation

Studio fournit Library, Builder, Assurance et Control Room ; les utilisateurs restent dans le produit opérationnel pour leur travail principal.

## 18. Effets sur les objets

Les objets métier restent chez Command, Investigate et Govern.

## 19. Effets sur les permissions

Une Automation Agent reçoit uniquement des permissions explicites ; elle ne s'accorde rien et ne contourne pas Govern.

## 20. Effets sur les parcours

Agentic Investigation crée une Automation Run et des Tool Calls, puis revient à l'analyste ; une réponse risquée crée une Action Request.

## 21. Effets sur les autres documents

README Studio, AI constraints, ownership, product boundaries et futurs objets Studio.

## 22. Migration

Ajouter Tool Call et Automation Run en Phase 7 ; supprimer les formulations laissant croire qu'un produit opérationnel possède son moteur d'agents.

## 23. Critères de réévaluation

Réévaluer si une capacité agentique ne peut être utilisée sans déplacer l'activité principale dans Studio.

## 24. Questions encore ouvertes

- OPEN-007 — relation Human Gate / Govern.
- OPEN-015 — transition Automation Run / Response Run.
