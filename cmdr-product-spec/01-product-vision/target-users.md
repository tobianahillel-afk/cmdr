---
id: target-users
domain: 01-product-vision
status: draft
owner: Head of Product
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-021
  - REQ-PROD-022
  - REQ-PROD-023
  - REQ-PROD-024
  - REQ-PROD-025
  - REQ-PROD-026
  - REQ-PROD-027
  - REQ-PROD-028
  - REQ-PROD-029
  - REQ-PROD-030
  - REQ-PROD-031
  - REQ-PROD-032
  - REQ-PROD-033
  - REQ-PROD-034
  - REQ-PROD-035
  - REQ-PROD-036
---
# Utilisateurs cibles

| Requirement | Rôle | Objectif | Responsabilité | Niveau | Fréquence | Décisions | Objets | Workspaces | Information prioritaire | Risque d'erreur | Collaboration |
|---|---|---|---|---|---|---|---|---|---|---|---|
| REQ-PROD-021 | Incident Commander | Stabiliser la situation et coordonner | Owner opérationnel | avancé | continu en crise | priorité, ownership, escalade | Incident, Task, Decision, Result | Mission Control, Incident Detail | impact, blocages, prochaine action | priorisation erronée | SOC, métier, approbateurs |
| REQ-PROD-022 | SOC Analyst L1 | Trier et qualifier rapidement | Triage initial | intermédiaire | quotidienne | disposition, escalade | Signal, Alert, Incident | Signals, Work Queue | contexte minimal fiable | fermer trop tôt | L2, Incident Commander |
| REQ-PROD-023 | SOC Analyst L2 | Approfondir et coordonner l'investigation | Triage avancé | avancé | quotidienne | ouvrir Case, demander collecte | Incident, Case, Evidence | Incident Detail, Case Workspace | timeline, entités, couverture | sur-investigation | L1, Hunter, DFIR |
| REQ-PROD-024 | Senior Analyst | Résoudre les cas complexes et guider | Expertise transverse | expert | fréquente | hypothèses, findings, escalade | Case, Hypothesis, Finding | Case Workspace, Workbench | contradictions, provenance | biais d'expertise | analystes, approbateurs |
| REQ-PROD-025 | Threat Hunter | Chercher des comportements non détectés | Hunts et hypothèses | avancé | récurrente | pivots, promotion en Case | Query, Signal, Hypothesis | Signals and Hunt | qualité des données, coverage | confirmation biaisée | Detection Engineer |
| REQ-PROD-026 | Forensic Analyst | Acquérir et analyser des preuves | Intégrité et reproductibilité | expert | par incident | acquisition, interprétation | Artifact, Evidence, Finding | Collection, Forensic Workbench | hashes, timeline, chain of custody | altération ou attribution abusive | Case owner, Auditor |
| REQ-PROD-027 | Malware Analyst | Caractériser un artefact malveillant | Analyse statique/dynamique | expert | par besoin | classification, IOC, behavior | Artifact, Evidence, Finding | Static, Sandbox | provenance, environnement | sur-confiance sandbox | Reverse Engineer, Detection |
| REQ-PROD-028 | Reverse Engineer | Comprendre le comportement binaire | Analyse approfondie | expert | spécialisée | annotation, fonction, comportement | Artifact, Finding | Reverse Workbench | symbols, xrefs, versions | conclusion non reproductible | Malware Analyst |
| REQ-PROD-029 | Detection Engineer | Créer et améliorer les détections | Rules, tests, deployment | expert | récurrente | test, approval, rollback | Detection, Signal, Rule assets | Detection Engineering | FP, coverage, data quality | remplacer règle par score opaque | Hunter, SOC |
| REQ-PROD-030 | Response Operator | Exécuter et vérifier la réponse | Run execution | avancé | par incident | suspendre, reprendre, rollback | Response Run, Result, Endpoint | Run Workspace | cible, classe, rollback | agir hors scope | Approver, Incident Commander |
| REQ-PROD-031 | Technical Approver | Évaluer faisabilité et risque technique | Autorité technique | expert | par demande | approve, condition, refuse | Action Request, Decision, Approval | Decision Workspace | blast radius, evidence, policy | approbation sans preuve | Business Owner, Operator |
| REQ-PROD-032 | Business Owner | Assumer impact et arbitrage métier | Autorité métier | métier | par crise | accepter impact ou interruption | Incident, Decision, Result | Decision Workspace, Mission Control | services, clients, alternatives | décision sur jargon technique | Incident Commander |
| REQ-PROD-033 | Customer Success Manager | Relier opérations et engagement client | Contexte et livraison | intermédiaire | récurrente selon déploiement | reporting, communication | Incident, Result, Report | Customers and Delivery | engagements, statut, audience | exposer des données non approuvées | Command, Reporting |
| REQ-PROD-034 | Auditor | Vérifier conformité et attribution | Contrôle indépendant | avancé | périodique | constater, demander preuve | Evidence, Decision, Approval, Audit | Audit Ledger | immutabilité, acteur, policy | audit dépendant d'un résumé | Govern, Security |
| REQ-PROD-035 | Platform Administrator | Maintenir plateforme et flotte | Administration | avancé | quotidienne | configurer, enrôler, corriger | Tenant, Role, Fleet, Policy | Settings Shell | health, version, scope | utiliser Settings pour enquêter | Security, Endpoint |
| REQ-PROD-036 | Automation Designer | Concevoir des automatisations sûres | Build, tests, assurance | avancé | récurrente | version, simulation, deploy | Skill, Tool, Workflow, Agent | Builder, Assurance, Control Room | permissions, coût, trace | automatiser sans owner métier | Product owners, Approvers |

## Principes transverses

- Les rôles décrivent des responsabilités, pas nécessairement des intitulés d'entreprise.
- Un utilisateur peut cumuler des rôles, sans cumuler automatiquement leurs permissions.
- La densité et les informations par défaut peuvent varier ; les sources canoniques restent identiques.
- Un rôle moins technique reçoit des explications et impacts, pas une preuve masquée.
- Un expert obtient les détails nécessaires sans transformer tous les écrans en workbench.

## Question ouverte

`OPEN-010` doit déterminer les préférences de densité par rôle après recherche utilisateur. La Phase 1 ne choisit pas de profils visuels définitifs.

## Critère d'acceptation

**Given** un Business Owner et un Reverse Engineer consultant la même Decision,  
**When** chacun ouvre son contexte,  
**Then** ils voient la même source canonique et le même statut, mais une hiérarchie adaptée : impact, alternatives et responsabilité pour le premier ; Evidence, méthode et détails techniques pour le second.
