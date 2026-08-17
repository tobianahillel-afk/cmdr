---
id: investigate-user-goals
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-023
  - REQ-PROD-024
  - REQ-PROD-025
  - REQ-PROD-026
  - REQ-PROD-027
---

# Objectifs utilisateurs Investigate

| Rôle | Objectifs 4B.1 | Informations prioritaires | Actions interdites |
|---|---|---|---|
| SOC Analyst L1/L2 | qualifier Signal, rechercher, lier au Case | source, fraîcheur, confiance, Incident, Case | changer priorité Command silencieusement |
| Threat Hunter | organiser Hunt, queries, résultats et conclusions | question, scope, runs, provenance, limites | transformer résultats en Evidence automatiquement |
| Senior Analyst | coordonner Case, revoir Findings et Action Request | objectif, Evidence pour/contre, reviewer, impact | confirmer un Finding sans preuve ou revue |
| Forensic Analyst | qualifier Artifacts et Evidence | acquisition, provenance, transformations, accès | modifier la source sans version |
| Malware Analyst / Reverse Engineer | fournir Artifacts et résultats futurs | parent, tool/version, derived relations | choisir ou lancer les moteurs 4B.2 dans cette phase |
| Auditor | reconstruire provenance et décisions analytiques | auteur, versions, sources, corrections | accéder à des données hors permission |
