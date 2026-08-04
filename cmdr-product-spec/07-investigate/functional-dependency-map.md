---
id: investigate-functional-dependency-map
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-006
  - REQ-PROD-012
  - REQ-PROD-014
---

# Functional dependency map — Investigate Phase 4B.1

| Capability group | Dépendances principales | Type | Comportement d’échec |
|---|---|---|---|
| CAP-INV-001..004 | Command Signal/Incident, Query Engine, data sources, Entity Resolution | projection/shared | Partial explicite ; aucune source ou priorité inventée |
| CAP-INV-005..008 | Query/Search Job, Case, Trace, Saved Views, Export, Studio optional | workspace/shared | Draft et résultats valides conservés ; aucune dépendance IA bloquante |
| CAP-INV-101..104 | Incident, Saved Views, Case, Entity Resolution, Object Linking | ownership/shared | Queue/Case utilisables avec contexte partiel ; aucun merge silencieux |
| CAP-INV-105..108 | Artifact, Evidence, Preview, Trace, Export, collection future | object/trust | Candidate ou qualification incomplete ; aucune Evidence automatique |
| CAP-INV-109..112 | Evidence, Timeline, Trace, Command Readiness, 4B.3 future | reasoning/shared | Finding/replay restent partial ; aucune conclusion inventée |
| CAP-INV-113 | Govern Action Request, Command impact, OPEN-007/013/015 | authority/transition | Draft conservé ; aucune Decision ou exécution locale |
| CAP-INV-114 | Reporting Engine, Export, Evidence/Findings, OPEN-014 | shared/reporting | Content package conservé ; aucune publication locale |

## Règles

- une dépendance ne transfère pas l’ownership ;
- chaque projection porte source, fraîcheur et permission ;
- une défaillance conserve les objets Investigate valides et expose son impact ;
- aucune dépendance IA n’est nécessaire à un workflow essentiel ;
- aucune dépendance 4B.2 ou 4B.3 n’est décrite comme livrée ;
- API, protocole, moteur, stockage et SLO restent hors Phase 4B.1.
