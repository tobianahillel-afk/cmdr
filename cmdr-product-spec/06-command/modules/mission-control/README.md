---
id: command-module-mission-control
domain: 06-command
status: draft
owner: Command Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-003
  - REQ-PROD-005
  - REQ-PROD-006
  - REQ-PROD-008
  - REQ-PROD-009
  - REQ-PROD-010
  - REQ-PROD-013
  - REQ-PROD-021
  - REQ-UX-007
---

# Mission Control

## Mission

Maintenir une conscience partagée de la situation, des priorités, des blocages, des handovers et des résultats sans devenir une seconde Work Queue.

## Utilisateurs

- Incident Commander
- SOC Analyst L2
- Business Owner
- Response Operator

## Ownership

**Le module possède :**

- composition de la situation Command
- priorité opérationnelle effective
- handover de coordination
- blocages portés par Incident/Task ou Task de suivi

**Le module ne possède pas :**

- files de travail dupliquées
- investigation technique
- Decision/Response Run/Result ownership
- chatbot global

## Capabilities

| ID | Capability | Delivery status | Delivery mode | Source |
|---|---|---|---|---|
| CAP-CMD-001 | Situation Overview | defined | planned | capabilities/situation-overview.md |
| CAP-CMD-002 | Priority Management | defined | planned | capabilities/priority-management.md |
| CAP-CMD-003 | Situation Timeline | defined | planned | capabilities/situation-timeline.md |
| CAP-CMD-004 | Handover | defined | planned | capabilities/handover.md |
| CAP-CMD-005 | Operational Blockers | defined | planned | capabilities/operational-blockers.md |
| CAP-CMD-006 | Recent Results and Outcomes | defined | planned | capabilities/recent-results-and-outcomes.md |

## Entrées et sorties du module

Les entrées détaillées sont possédées par chaque capability. Au niveau module, les entrées sont des objets Command, des références permission-aware et des projections sourcées ; les sorties sont des mutations Command auditées, des Tasks, des relations ou des context packages vers le produit propriétaire.

## Shared Capabilities consommées

| Shared Capability | Usage local | Données locales | Source canonique |
|---|---|---|---|
| Global Search | retrouver les objets autorisés et ouvrir des liens profonds | objets/références Command du module | ../../../12-shared-capabilities/global-search.md |
| Command Palette | naviguer et lancer des actions non dangereuses autorisées | objets/références Command du module | ../../../04-experience-architecture/command-palette.md |
| Notifications | signaler assignment, SLA, handover, escalation et résultats | objets/références Command du module | ../../../12-shared-capabilities/notification-center.md |
| Activity Stream | afficher les mutations et collaborations du contexte | objets/références Command du module | ../../../03-design-system/components/activity-stream.md |
| Trace | attribuer humain, règle, moteur, workflow, agent ou source externe | objets/références Command du module | ../../../03-design-system/components/trace.md |
| Timeline | ordonner les événements métier fournis par Command | objets/références Command du module | ../../../12-shared-capabilities/timeline-engine.md |
| Reporting | composer des reports à partir de citations Command | objets/références Command du module | ../../../12-shared-capabilities/reporting-engine.md |
| Export | export permission-aware, redaction et audit | objets/références Command du module | ../../../12-shared-capabilities/export-engine.md |
| Inspector | inspecter la sélection sans composant local concurrent | objets/références Command du module | ../../../03-design-system/components/inspector.md |
| Context Bar | conserver tenant, environnement et objets liés | objets/références Command du module | ../../../03-design-system/components/context-bar.md |
| Cross-product Links | références stables sans transfert d’ownership | objets/références Command du module | ../../../12-shared-capabilities/object-linking-service.md |
| Collaboration | comments, mentions et conflits | objets/références Command du module | ../../../12-shared-capabilities/collaboration-service.md |
| Presence | présence informative, jamais permission | objets/références Command du module | ../../../12-shared-capabilities/collaboration-service.md |
| Localization | formats et libellés localisés | objets/références Command du module | ../../../12-shared-capabilities/localization-engine.md |
| Audit Hooks | enregistrer acteur, changement, résultat et correlation | objets/références Command du module | ../../../05-domain-model/event-and-audit-semantics.md |
| Metrics | définitions, fraîcheur et réconciliation des métriques | objets/références Command du module | ../../../12-shared-capabilities/metrics-engine.md |

## États et erreurs

Les états fonctionnels propres à chaque capability sont définis dans les fichiers individuels. Loading, Empty, Partial, Error, Offline, Permission denied et Stale suivent le Design System ; aucune source manquante n’est masquée.

## IA et automatisation

Toutes les fonctions essentielles disposent d’une voie manuelle ou déterministe. Une règle, un workflow ou un agent conserve provenance, version/run, facteurs et alternative non-IA. Une suggestion n’est jamais une priorité, un Finding ou une Decision effective.

## Documents du module

- [`capabilities.md`](capabilities.md) — index du module ;
- `capabilities/` — specifications canoniques individuelles ;
- [`../../capability-map.md`](../../capability-map.md) — registre Command ;
- [`../../functional-dependency-map.md`](../../functional-dependency-map.md) — dépendances transversales ;
- [`../../screen-capability-map.md`](../../screen-capability-map.md) — préparation des écrans.

## Critère d’acceptation

**Given** le module Mission Control et un utilisateur autorisé, **When** une capability est utilisée, **Then** l’owner, les objets, entrées, sorties, classes d’action, états, dépendances et alternative sans IA sont résolus depuis sa source canonique; aucun objet d’un autre produit n’est redéfini.
