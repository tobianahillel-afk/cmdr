---
id: dependency-roadmap
domain: 18-roadmap-and-releases
status: draft
owner: Product Operations Lead
updated: 2026-08-03
source-of-truth: canonical
---
# Dependency Roadmap

| Séquence | Dépendance | Propriétaire | Débloque |
|---|---|---|---|
| 1 | Approbation des ADR et ownership | Product Architecture | tous domaines |
| 2 | Query language et search backend | Platform Architecture | Event Search, Global Search |
| 3 | Policy engine et authority model | Security Architecture | Govern, Endpoint commands |
| 4 | Workflow/orchestration language | CMDR Studio | Studio et Govern Playbooks |
| 5 | Endpoint PKI, plateformes et update model | Endpoint / Platform Settings | EDR natif |
| 6 | Forensic engines et sandbox profiles | Investigate | chaîne d’analyse |
| 7 | Retention, legal hold, signature et timestamp | Compliance / Security | Evidence, Decisions, Reports, Audit |
| 8 | SLO, RTO, RPO et deployment models | Platform Engineering | release readiness |

Les détails de responsabilité vivent dans `../00-governance/dependency-register.md`. Les questions techniques restent dans les documents propriétaires.
