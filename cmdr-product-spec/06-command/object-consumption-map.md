---
id: command-object-consumption-map
domain: 06-command
status: draft
owner: Command Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-OBJ-001
  - REQ-OBJ-002
  - REQ-OBJ-005
  - REQ-PROD-006
  - REQ-PROD-013
---
# Object consumption and dependency map — Command

| Capabilities | Objet/projection | Owner | Besoin fonctionnel | Lacune actuelle | Phase propriétaire |
|---|---|---|---|---|---|
| CAP-CMD-109 | Acknowledgement | local/session ou audit selon effet | write: marquer staleness vue | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-106 | Action Request | Govern | write: initier le draft via Govern | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-204 | Action Request context | Govern | write: transmettre impact | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-103 | Audit | Shared audit hook | write: période d’ownership | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-002 | Audit trail | source audit partagée | write: ajouter historique de changement | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-305 | Capability | Capability Register | read: ID, owner, status, delivery | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-302, CAP-CMD-304 | Capability Readiness | Command | read: baseline | record fonctionnel non encore formalisé comme objet | Phase Objets / produit owner |
| CAP-CMD-301 | Capability projection | product owners | read: availability/config/test state | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-101 | Case / Decision / Response Run | produits propriétaires | read: statut et lien | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-103 | Case / Decision / Run | produit propriétaire | read: owner projeté | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-109 | Case / Decision / Run projections | produits propriétaires | read: freshness metadata | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-004, CAP-CMD-106 | Case / Finding | Investigate | read: références et résumé autorisé | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-110 | Case / Finding / Decision | produit propriétaire | read: liens et statut | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-110 | Case or Action Request | Investigate/Govern | write: demander création/liaison | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-205 | Confidence / severity | Signal/Alert/source | read: context | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-203 | Coverage projections | source products | read: scope, status, definition, freshness | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-203 | Coverage source | source product | write: aucune mutation | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-401 | Customer/contract source | deployment owner | write: none | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-401 | Customer/engagement context | deployment-specific source | read: portfolio/obligations | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-201 | Data quality follow-up | Command | write: Task | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-005 | Decision / Case / external dependency | produit propriétaire | read: statut pertinent | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-001, CAP-CMD-004, CAP-CMD-106 | Decision / Response Run / Result | Govern | read: statut, portée, résultat résumé | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-003 | Decision / Run / Result | Govern | read: événements gouvernés | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-110 | Escalation record/state | Command coordination | write: créer, envoyer, accepter, décliner, retourner | record fonctionnel non encore formalisé comme objet | Phase Objets / produit owner |
| CAP-CMD-302 | Exercise record | Command readiness | write: create/update/run/complete | record fonctionnel non encore formalisé comme objet | Phase Objets / produit owner |
| CAP-CMD-303 | Exercise/readiness/coverage | Command/source | read: gap and context | record fonctionnel non encore formalisé comme objet | Phase Objets / produit owner |
| CAP-CMD-108 | Export job | Shared Export Engine | write: créer un job | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-201 | Exposure / coverage projections | sources propriétaires | read: relation au service | objet canonique absent; source fonctionnelle/projection seulement | Phase Objets / produit owner |
| CAP-CMD-202 | Exposure projection | external/shared source | read: type, status, source, freshness | objet canonique absent; source fonctionnelle/projection seulement | Phase Objets / produit owner |
| CAP-CMD-202 | Exposure source | external owner | write: aucune opération | objet canonique absent; source fonctionnelle/projection seulement | Phase Objets / produit owner |
| CAP-CMD-204 | Finding/Result | Investigate/Govern | read: technical evidence/outcome | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-004 | Handover record | Command coordination | write: créer, mettre à jour, envoyer, reconnaître, superséder | record fonctionnel non encore formalisé comme objet | Phase Objets / produit owner |
| CAP-CMD-109 | Health | Settings/Shared | read: dependency status | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-305 | Health/Assurance/Exercise | source products | read: evidence | record fonctionnel non encore formalisé comme objet | Phase Objets / produit owner |
| CAP-CMD-001, CAP-CMD-003, CAP-CMD-006, CAP-CMD-101, CAP-CMD-202, CAP-CMD-204 | Incident | Command | read: priorité, état, owner, impact, prochaine action | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-001, CAP-CMD-006, CAP-CMD-106, CAP-CMD-204 | Incident | Command | write: mise à jour de prochaine action ou accusé de situation | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-002, CAP-CMD-004, CAP-CMD-005, CAP-CMD-102, CAP-CMD-103, CAP-CMD-104, CAP-CMD-105, CAP-CMD-106, CAP-CMD-108, CAP-CMD-109, CAP-CMD-110, CAP-CMD-201, CAP-CMD-205, CAP-CMD-401 | Incident / Task | Command | read: priority, impact, urgency, owner, state | owner Task à aligner sur Command; schéma détaillé reporté | Phase Objets / produit owner |
| CAP-CMD-002, CAP-CMD-005, CAP-CMD-101, CAP-CMD-102, CAP-CMD-103, CAP-CMD-104, CAP-CMD-105, CAP-CMD-108, CAP-CMD-201, CAP-CMD-202 | Incident / Task | Command | write: modifier priority et justification | owner Task à aligner sur Command; schéma détaillé reporté | Phase Objets / produit owner |
| CAP-CMD-004 | Incident / Task ownership | Command | write: transfert séparé après acceptation | owner Task à aligner sur Command; schéma détaillé reporté | Phase Objets / produit owner |
| CAP-CMD-205 | Incident/Task | Command | write: acceptation via CAP-CMD-002 | owner Task à aligner sur Command; schéma détaillé reporté | Phase Objets / produit owner |
| CAP-CMD-102 | Notification | Shared | write: demande de prise en charge | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-302 | Operational Plan | Command readiness | read: roles, procedures, activation | record fonctionnel non encore formalisé comme objet | Phase Objets / produit owner |
| CAP-CMD-304 | Operational Plan record | Command readiness | write: create/update/review/supersede | record fonctionnel non encore formalisé comme objet | Phase Objets / produit owner |
| CAP-CMD-301 | Plan/Exercise records | Command readiness | read: status and results | record fonctionnel non encore formalisé comme objet | Phase Objets / produit owner |
| CAP-CMD-304 | Plan/Playbook references | Command/Govern/Studio/source | read: metadata only | record fonctionnel non encore formalisé comme objet | Phase Objets / produit owner |
| CAP-CMD-102 | Principal / Role | Platform Settings | read: identité, équipe, statut disponible | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-103 | Principal / Team | Settings | read: identité et statut | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-205 | Priority recommendation | Command coordination | write: proposal record | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-301, CAP-CMD-305 | Readiness assessment | Command readiness | write: update status/rationale | record fonctionnel non encore formalisé comme objet | Phase Objets / produit owner |
| CAP-CMD-104 | Recommendation disposition | Shared audit | write: audit event | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-401 | Report | Reporting Engine | read: draft/review/published | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-401 | Report request/draft context | Shared | write: create via Reporting Engine | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-303, CAP-CMD-401 | Result | Govern | read: outcome/residual risk | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-006 | Result / Response Run / Decision | Govern | read: résultat, vérification, portée, conditions | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-002 | SLA context | Command/policy source | read: risque et échéance | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-105 | SLA policy / engagement | Settings ou deployment-specific source | read: durée, calendrier, pause | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-108 | Saved View / filters | Shared/Command | read: scope de sélection | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-101 | Saved View application | Shared/Command catalog | write: appliquer configuration | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-002, CAP-CMD-202, CAP-CMD-203, CAP-CMD-204 | Service | Shared | read: criticité et dépendances | objet canonique absent; source fonctionnelle/projection seulement | Phase Objets / produit owner |
| CAP-CMD-001, CAP-CMD-201 | Service | Shared Business Service Catalog | read: owner, criticité, dépendances | objet canonique absent; source fonctionnelle/projection seulement | Phase Objets / produit owner |
| CAP-CMD-205 | Service / Exposure / Coverage | Shared/source products | read: context | objet canonique absent; source fonctionnelle/projection seulement | Phase Objets / produit owner |
| CAP-CMD-104 | Service / SLA | Shared/policy | read: criticité, risk | objet canonique absent; source fonctionnelle/projection seulement | Phase Objets / produit owner |
| CAP-CMD-304 | Service/Incident | Shared/Command | read: scope/context | objet canonique absent; source fonctionnelle/projection seulement | Phase Objets / produit owner |
| CAP-CMD-106 | Signal / Alert | Command | read: source, severity, disposition | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-104 | Signal / Alert | Command/source | read: severity, confidence | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-107 | Source object | produit propriétaire | read: résumé et relation | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-107, CAP-CMD-303 | Source relation | Object Linking Service | write: ajouter lien | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-001, CAP-CMD-101, CAP-CMD-107, CAP-CMD-203, CAP-CMD-301, CAP-CMD-303, CAP-CMD-305 | Task | Command | read: urgence, owner, échéance, blocage | owner Task à aligner sur Command; schéma détaillé reporté | Phase Objets / produit owner |
| CAP-CMD-005, CAP-CMD-105, CAP-CMD-106, CAP-CMD-107, CAP-CMD-109, CAP-CMD-203, CAP-CMD-301, CAP-CMD-302, CAP-CMD-303, CAP-CMD-304, CAP-CMD-305, CAP-CMD-401 | Task | Command | write: créer une action de déblocage si nécessaire | owner Task à aligner sur Command; schéma détaillé reporté | Phase Objets / produit owner |
| CAP-CMD-006 | Timeline | Shared/Command content | write: ajouter événement de consommation | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-003 | Timeline Entry | Shared | read: type, timestamps, source, objet | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |
| CAP-CMD-302 | Workflow/Simulation | Studio | read: optional references only | schéma/cardinalités reportés à Phase Objets | Phase Objets / produit owner |

## Confirmed missing canonical object paths

`service.md`, `exposure.md`, `report.md` and `audit-record.md` do not exist. Service uses Business Service Catalog; Report uses Reporting Engine; audit uses event-and-audit semantics. No object is created in Phase 4A.

## Ownership correction

Operational Task is owned by Command according to source decisions and the ownership register. Shared Task Inbox remains a consumer/projection service.
