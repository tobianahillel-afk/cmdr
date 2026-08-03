---
id: source-unresolved-decisions
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: source-material
requirements:
  - REQ-PROD-048
  - REQ-PROD-062
  - REQ-BRAND-008
  - REQ-PROD-057
---

# Unresolved Decisions

Une option, une valeur Draft ou un prototype n'est pas une décision. Une question se ferme uniquement avec preuve approuvée, owner, mise à jour des fichiers affectés et de la matrice.

## Décisions ouvertes

### OPEN-001 — Investigate palette direction

```yaml
decision_id: OPEN-001
status: open
owner: Brand Design Lead
target_phase: Phase 2 review
blocking: false
requirement_id: REQ-PROD-048
required_evidence: three proposed directions; light/dark, contrast, dense workbench and family review
affected_files: 02-brand/investigate/palette-proposals.md; 03-design-system/foundations/tokens.md
```


### OPEN-002 — Govern palette direction

```yaml
decision_id: OPEN-002
status: open
owner: Brand Design Lead
target_phase: Phase 2 review
blocking: false
requirement_id: REQ-PROD-049
required_evidence: three proposed directions; decision, audit, non-punitive authority and contrast
affected_files: 02-brand/govern/palette-proposals.md; 03-design-system/foundations/tokens.md
```


### OPEN-003 — CMDR Studio palette direction

```yaml
decision_id: OPEN-003
status: open
owner: Brand Design Lead
target_phase: Phase 2 review
blocking: false
requirement_id: REQ-PROD-050
required_evidence: three proposed directions; Builder, Control Room, Assurance and anti-purple-AI review
affected_files: 02-brand/studio/palette-proposals.md; 03-design-system/foundations/tokens.md
```


### OPEN-004 — Final typography stack and licensing

```yaml
decision_id: OPEN-004
status: open
owner: Design System Lead
target_phase: Phase 2 review
blocking: false
requirement_id: REQ-PROD-051
required_evidence: license, languages, dense tables, code, multi-OS and performance
affected_files: 02-brand/cmdr/typography.md; 03-design-system/foundations/typography.md
```


### OPEN-005 — Initial forensic engines

```yaml
decision_id: OPEN-005
status: open
owner: Investigate Product Lead
target_phase: Phase 4
blocking: true
requirement_id: REQ-PROD-052
required_evidence: outcomes, licensing, integrity and replacement strategy
affected_files: 07-investigate/modules/*
```


### OPEN-006 — Customers and Delivery applicability

```yaml
decision_id: OPEN-006
status: open
owner: Command Product Lead
target_phase: Phase 4
blocking: false
requirement_id: REQ-PROD-053
required_evidence: managed-service and internal-deployment evidence
affected_files: 06-command/modules/customer-and-reports/*
```


### OPEN-007 — Human Gate and Govern relationship

```yaml
decision_id: OPEN-007
status: open
owner: Security Architecture
target_phase: Phase 4
blocking: true
requirement_id: REQ-PROD-054
required_evidence: action class, authority, SoD and audit
affected_files: 08-govern/*; 09-cmdr-studio/human-gates/*
```


### OPEN-008 — Endpoint platform support

```yaml
decision_id: OPEN-008
status: open
owner: Endpoint Agent Product Lead
target_phase: Phase 4
blocking: true
requirement_id: REQ-PROD-055
required_evidence: demand, sensor feasibility, support and compatibility
affected_files: 11-endpoint-agent/platform-support.md
```


### OPEN-010 — Density by role and activity

```yaml
decision_id: OPEN-010
status: open
owner: UX Architecture
target_phase: Phase 3 review
blocking: false
requirement_id: REQ-PROD-057
required_evidence: usability studies across queue, case, workbench, decision and settings
affected_files: 04-experience-architecture/role-based-defaults.md; 03-design-system/foundations/density.md
```


### OPEN-011 — Mobile forensic scope

```yaml
decision_id: OPEN-011
status: open
owner: Investigate Product Lead
target_phase: Phase 4
blocking: false
requirement_id: REQ-PROD-058
required_evidence: demand, legality and engine strategy
affected_files: 07-investigate/modules/*
```


### OPEN-012 — Cloud analysis scope

```yaml
decision_id: OPEN-012
status: open
owner: Investigate Product Lead
target_phase: Phase 4
blocking: false
requirement_id: REQ-PROD-059
required_evidence: journeys, sources, APIs and ownership
affected_files: 07-investigate/modules/*
```


### OPEN-013 — Default governance for Action Class 2

```yaml
decision_id: OPEN-013
status: open
owner: Security Architecture
target_phase: Phase 4
blocking: true
requirement_id: REQ-PROD-060
required_evidence: blast radius, rollback reliability and tenant policy
affected_files: 08-govern/*; 14-security-permissions-and-trust/*
```


### OPEN-014 — Artifact versus Attachment

```yaml
decision_id: OPEN-014
status: open
owner: Investigate Product Lead
target_phase: Phase 7
blocking: true
requirement_id: REQ-PROD-061
required_evidence: evidence lifecycle, provenance, collaboration, retention
affected_files: 05-domain-model/objects/artifact.md
```


### OPEN-015 — Automation Run to Response Run bridge

```yaml
decision_id: OPEN-015
status: open
owner: CMDR Studio Product Lead
target_phase: Phase 4
blocking: true
requirement_id: REQ-PROD-062
required_evidence: action classification, Decision, Tool effects and audit
affected_files: 08-govern/*; 09-cmdr-studio/*
```


### OPEN-016 — Final wordmark construction and optional symbol

```yaml
decision_id: OPEN-016
status: open
owner: Brand Design Lead
target_phase: Phase 2 review
blocking: false
requirement_id: REQ-BRAND-008
required_evidence: small-size, monochrome, legal, spacing and application tests
affected_files: 02-brand/logo-and-wordmark.md; assets/brand/*
```

## Disposition Phase 3

- `OPEN-001`, `OPEN-002`, `OPEN-003` : les slots produit restent `unresolved`; aucune valeur n'alimente les tokens actifs.
- `OPEN-004` : les métriques typographiques sont Draft, les aliases de familles restent non résolus et aucun fichier de police n'est ajouté.
- `OPEN-010` : la Phase 3 adopte des défauts par activité avec override utilisateur, mais la préférence finale par rôle attend les études ; la question reste ouverte.
- `OPEN-016` : aucun asset de logo n'est produit.
- `OPEN-005` à `OPEN-008` et `OPEN-011` à `OPEN-015` restent hors du pouvoir de décision de Phase 3.

## Historique résolu

`OPEN-009` — capability delivery-classification authority — résolue en Phase 1 : capability owner propose, Product Lead et Product Architecture approuvent, Security/Engineering examinent selon l'effet, QA vérifie la preuve.

## Règles de maintenance

Aucun placeholder générique. Commencer la phase cible ne ferme pas la question. Un alias non résolu doit échouer explicitement plutôt que prendre une valeur arbitraire.
