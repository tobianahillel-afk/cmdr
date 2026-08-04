---
id: dependency-register
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-012
  - REQ-PROD-019
  - REQ-PROD-013
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
| DEP-008 | OPEN-005 | workbench forensic | moteur | stratégie initiale | open | limites intégration inconnues | Investigate Product Lead | technique | REQ-PROD-052 | Phase 4B/8 |
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

## Phase 4A — lecture

- les 27 capabilities Command sont fonctionnellement définies mais leur delivery mode reste `planned` ;
- aucune dépendance à une IA n’est bloquante pour un workflow essentiel ;
- `OPEN-006`, `OPEN-010` et `OPEN-013` restent ouvertes ;
- les lacunes Service, Exposure, Report et Audit Record sont routées vers leurs sources propriétaires ou phases ultérieures ;
- la Work Queue n’a aucune dépendance vers un ancien écran déprécié.

## Mise à jour

`resolved` exige que tous les dépendants soient mis à jour. `active` décrit une règle permanente. `partial` décrit un contrat suffisant pour la Phase 4A mais incomplet pour une phase ultérieure.

## Critère d’acceptation

**Given** une capability Command dépendante d’un objet externe ou d’une décision ouverte, **When** le registre est consulté, **Then** owner, impact, caractère bloquant, Requirement IDs et phase cible sont explicites sans choix technique inventé.
