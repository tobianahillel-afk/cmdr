---
id: dependency-register
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-012
  - REQ-PROD-013
  - REQ-PROD-014
  - REQ-PROD-019
---

# Registre des dépendances produit et documentaires

## Règles

Le registre capture les dépendances nécessaires à une décision produit ou documentaire. Il ne choisit aucun protocole, API, stockage, moteur, fournisseur ou infrastructure. Une dépendance `partial` peut permettre une specification fonctionnelle tout en bloquant la phase Objets, Permissions, Technique ou Implémentation.

| ID | Source | Dépendant | Type | Raison | Statut | Impact | Propriétaire | Bloquant | Requirement IDs | Revue |
|---|---|---|---|---|---|---|---|---|---|---|
| DEP-001 | source-material/ | 01-product-vision/ | décision | appliquer mission, principes et frontières | active | vision incohérente si non propagée | Head of Product | oui | REQ-PROD-001..020 | Phase 1 |
| DEP-002 | product-boundaries.md | READMEs produits | frontière | empêcher la revendication concurrente | active | doublons de responsabilité | Product Architecture | oui | REQ-PROD-013..018 | Phase 1 |
| DEP-003 | ownership-register.md | modèle d’objets | propriété | guider schémas et machines d’état | active | objets concurrents | Product Architecture | oui | REQ-OBJ-001..012 | Phase 7 |
| DEP-004 | terminology-rules.md | UX et contenu | vocabulaire | distinguer page, vue, mode, filtre, agent et run | active | navigation ambiguë | Content Design | oui | REQ-UX-001, REQ-OBJ-009 | Phase 3 |
| DEP-005 | ADR-0005 | Work Queue | UX | consolider variantes en vues | resolved | écrans clones | UX Architecture | oui | REQ-UX-008, REQ-UX-009 | Phase 3 |
| DEP-006 | Phase 2 marque | identités produit | visuel | palettes/typographie ouvertes | open | design final non validable | Brand Design Lead | non | REQ-PROD-048..051 | Phase 2 review |
| DEP-007 | Capability map | modules fonctionnels | classification | ne pas annoncer planned comme livré | active | promesse trompeuse | Product Architecture | oui | REQ-PROD-012, REQ-PROD-019 | Phase 4 |
| DEP-008 | OPEN-005 | workbench forensic | moteur | stratégie initiale | open | limites intégration inconnues | Investigate Product Lead | technique | REQ-PROD-052 | Phase 4B.2/8 |
| DEP-009 | OPEN-007 | Studio et Govern | autorité | Human Gate versus Decision | open | double approbation | Govern + Studio | oui | REQ-PROD-054 | Phase 4C/7 |
| DEP-010 | OPEN-008 | Endpoint Agent | plateforme | support initial | open | flotte/actions incomplètes | Endpoint Product Lead | avant implémentation | REQ-PROD-055 | Phase 4D/8 |
| DEP-011 | OPEN-014 | Investigate | modèle | Artifact versus Attachment | open | preuve ambiguë | Investigate Product Lead | oui modèle | REQ-PROD-061 | Phase 7 |
| DEP-012 | OPEN-015 | Studio / Govern | modèle | Automation Run versus Response Run | open | audit confondu | Studio + Govern | oui | REQ-PROD-062, REQ-OBJ-009 | Phase 7 |
| DEP-013 | Phase 3 | écrans pilotes | UX | niveau de qualité avant écrans | planned | templates insuffisants | UX Architecture | oui | REQ-UX-010 | Phase 6 |
| DEP-014 | Phase 5 | parcours | transitions | contexte bout en bout | planned | navigation incohérente | Product Architecture | oui | REQ-JRN-001..008 | Phase 5 |
| DEP-015 | Phase 7 | permissions détaillées | security | unifier namespaces après frontières | planned | contrôle d’accès incohérent | Security Architecture | oui | REQ-SEC-001..005 | Phase 7 |
| DEP-CMD-001 | Capability template | 27 capabilities Command | documentaire | contrat à 27 sections et ID immuable | active | capabilities incomplètes | Product Architecture | oui | REQ-PROD-006, REQ-PROD-013 | Phase 4A |
| DEP-CMD-002 | Incident/Task | Work Queue | objet | coordination, assignment, priority, SLA et relations | partial | schémas et états Draft | Command Product Lead | non Phase4A; oui Phase7 | REQ-OBJ-001, REQ-OBJ-012 | Phase 7 |
| DEP-CMD-003 | Business Service Catalog | Risk and Coverage | shared capability | Service n’existe pas comme objet canonique | partial | criticité/dépendances comme projection | Shared Capabilities Lead | non | REQ-PROD-021 | Phase 4D/7 |
| DEP-CMD-004 | Sources d’exposition | Exposure Overview | intégration | Exposure n’existe pas comme objet canonique | open | source et freshness variables | Command + source owner | non Phase4A | REQ-PROD-037..042 | Phase 4D/7 |
| DEP-CMD-005 | Saved Views | Unified Work Queue | shared capability | stockage et permission-aware application | active | aucune seconde source | Shared Capabilities Lead | oui | REQ-OBJ-011, REQ-OBJ-012 | Phase 4A/4D |
| DEP-CMD-006 | OPEN-006 | Customers and Delivery | scope | applicabilité par déploiement | open | module non universel | Command Product Lead | oui avant activation | REQ-PROD-053 | Phase 4A review |
| DEP-CMD-007 | OPEN-010 | densité des activités | Command defaults | UX research | défauts Draft, validation future | open | préférences de rôle non finales | UX Architecture | non fonctionnel | REQ-PROD-057 | Phase 3 review |
| DEP-CMD-008 | OPEN-013 | actions classe 2 | Command mutations | governance | valeur par défaut non décidée | open | step-up/Govern selon contexte | Security Architecture | oui avant politique finale | REQ-PROD-060, REQ-SEC-002 | Phase 4C/7 |
| DEP-CMD-009 | Investigate Case | Incident Coordination | transition | création/liaison idempotente | partial | contrat détaillé futur | Investigate Product Lead | non Phase4A | REQ-PROD-008, REQ-OBJ-002 | Phase 4B/5 |
| DEP-CMD-010 | Govern Action Request | Escalation | transition | package contexte/impact/action/alternatives | partial | autorité reste Govern | Govern Product Lead | oui classe 3/4 | REQ-PROD-015, REQ-SEC-002 | Phase 4C/5 |
| DEP-CMD-011 | Result | Situation/Incident | projection | réinjecter vérification et risque résiduel | partial | Result reste Govern | Govern Product Lead | non | REQ-PROD-008, REQ-OBJ-007 | Phase 4C/5 |
| DEP-CMD-012 | Permission catalog | actions Command | permission | familles assignment/SLA/bulk/handover à atomiser | partial | aucune permission finale inventée | Security Architecture | oui avant implémentation | REQ-SEC-001..005 | Phase 7 |
| DEP-INV-001 | Capability template | 22 capabilities Investigate | documentaire | 27 sections, six tableaux et ID immuable | active | capability incomplète si contrat absent | Product Architecture | oui | REQ-PROD-006,014 | Phase 4B.1 |
| DEP-INV-002 | Query et Search Job Shared | Event Search / provenance | objet partagé | authoring, exécution, annulation, résultats et versions | partial | moteur/langage/index restent ouverts | Shared Capabilities Lead | non 4B.1; oui technique | REQ-PROD-014,019 | 4B.2/8 |
| DEP-INV-003 | Signal/Alert/Incident Command | Signal Triage et Case intake | projection | préserver source, priorité et coordination Command | active | aucune ownership concurrente | Command + Investigate | oui frontière | REQ-PROD-008,013,014 | 4B.1/5 |
| DEP-INV-004 | Incident Command | Case Lifecycle | transition | lier un ou plusieurs Incidents sans transfert | partial | cardinalités et workflow final ouverts | Investigate Product Lead | non 4B.1; oui Objets | REQ-OBJ-002, REQ-PROD-014 | Phase 7 |
| DEP-INV-005 | Artifact, Case, Provenance | Evidence | objet/trust | qualification explicite, versions, transformations et accès | partial | stockage/intégrité/custody détaillés ouverts | Investigate + Trust | oui avant implémentation | REQ-PROD-061,062 | Phase 7/8 |
| DEP-INV-006 | Evidence | Finding | raisonnement | support favorable/contradictoire et revue | active | aucun Finding sans références | Investigate Product Lead | oui | REQ-PROD-014,016 | 4B.1/7 |
| DEP-INV-007 | Govern Action Request | Action Request Preparation | autorité/transition | Investigate produit, Govern possède lifecycle | partial | permissions et workflow Govern futurs | Govern Product Lead | oui pour soumission | REQ-PROD-015,016; REQ-SEC-002 | 4C/5/7 |
| DEP-INV-008 | Reporting Engine Shared | Case Reporting Preparation | shared capability | citations, versions, redactions et export request | partial | objet Report et formats futurs | Shared Capabilities Lead | non 4B.1 | REQ-PROD-018 | 4D/8 |
| DEP-INV-009 | CMDR Studio | assistance automatisée Investigate | automation | provenance, Workflow, Agent et Tool Calls | partial | aucun agent requis ; bridge run ouvert | Studio Product Lead | non essentiel | REQ-AI-002,006,010 | 4C/7 |
| DEP-INV-010 | OPEN-014 | Attachment Handling | décision objet | relation Attachment/Artifact/Evidence | open | permissions, rétention et migration non finales | Investigate + Product Architecture | oui modèle | REQ-PROD-061 | Phase 7 |
| DEP-INV-011 | OPEN-013 | actions classe 2 Investigate | gouvernance | step-up et Govern par défaut non décidés | open | mutations réversibles non finalisées | Security Architecture | oui avant politique finale | REQ-PROD-060, REQ-SEC-002 | 4C/7 |
| DEP-INV-012 | OPEN-015 | sorties automatisées / Result bridge | provenance | distinguer Automation Run et Response Run | open | audit et relations inter-run incomplets | Studio + Govern | oui modèle | REQ-PROD-062, REQ-OBJ-009 | Phase 7 |
| DEP-INV-013 | OPEN-007 | Action Request / Human Gate | autorité | éviter double approbation ou contournement | open | parcours Govern final incomplet | Govern + Studio | oui classes 3/4 | REQ-PROD-054 | 4C/7 |
| DEP-INV-014 | phases 4B.2 et 4B.3 | handoffs futurs | frontière | collecte, workbench, détection et intelligence | planned | fonctions futures non spécifiées comme livrées | Investigate Product Lead | non 4B.1 | REQ-PROD-014,052 | 4B.2/4B.3 |

## Lecture de phase

- Phase 4A : 27 capabilities Command fonctionnellement définies, delivery mode `planned`.
- Phase 4B.1 : 22 capabilities Investigate fonctionnellement définies, 21 `defined`, CAP-INV-106 `proposed`, toutes `planned`.
- aucune dépendance IA n’est bloquante pour un workflow essentiel ;
- OPEN-005,007,011,012,013,014 et 015 restent ouvertes ;
- aucune dépendance 4B.2 ou 4B.3 n’est décrite comme capability livrée.

## Mise à jour et acceptation

`resolved` exige la mise à jour de tous les dépendants. `active` décrit une règle permanente. `partial` décrit un contrat suffisant pour la phase courante mais incomplet pour une phase ultérieure.

**Given** une capability dépendante d’un objet externe ou d’une décision ouverte, **When** le registre est consulté, **Then** owner, impact, caractère bloquant, Requirement IDs, comportement de phase et revue cible sont explicites sans choix technique inventé.
