---
id: command-user-goals
domain: 06-command
status: draft
owner: Command Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-013
  - REQ-PROD-021
  - REQ-PROD-032
  - REQ-PROD-033
---
# User goals — Command

| Rôle | Questions auxquelles Command doit répondre | Capabilities principales |
|---|---|---|
| Incident Commander | que se passe-t-il, pourquoi, qui agit, quelle prochaine action et quelle autorité manque ? | CAP-CMD-001..006, 102..110, 204..205 |
| SOC Analyst L1/L2 | quel travail m’appartient, quel est son contexte, son SLA et son état ? | CAP-CMD-101..109 |
| Team Lead | où sont les items sans owner, bloqués ou à risque, et comment répartir le travail ? | CAP-CMD-001, 002, 005, 102, 103, 105, 108 |
| Business/Service Owner | quel service et quel impact sont concernés, quelle confirmation est attendue ? | CAP-CMD-201, 204, 205 |
| Readiness Coordinator | quelles capacités, plans et exercices sont incomplets ou non testés ? | CAP-CMD-301..305 |
| Service Delivery / Customer Success | quels engagements et reports s’appliquent dans ce déploiement ? | CAP-CMD-401, selon le profil de déploiement et le contexte d’autorisation |

## Résultats communs

Chaque rôle voit owner, source, freshness, permission et prochaine action. Aucun rôle n’obtient implicitement l’autorité Govern ou l’ownership d’un objet externe.
