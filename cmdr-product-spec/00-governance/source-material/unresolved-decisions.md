---
id: source-unresolved-decisions
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-04
source-of-truth: source-material
requirements:
  - REQ-PROD-048
  - REQ-PROD-062
  - REQ-BRAND-008
  - REQ-PROD-057
  - REQ-PROD-053
  - REQ-PROD-060
---

# Unresolved Decisions

Une option, une valeur Draft ou une specification fonctionnelle n’est pas une décision. Une question se ferme uniquement avec preuve approuvée, owner, mise à jour des dépendants et de la matrice.

## Décisions ouvertes — 15

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
required_evidence: three proposed directions; decision, audit, authority and contrast
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
target_phase: Phase 4B.2
blocking: true
requirement_id: REQ-PROD-052
required_evidence: outcomes, licensing, integrity and replacement strategy
affected_files: 07-investigate/modules/collection-and-live-response/*; 07-investigate/modules/analysis-workbench/*
```
### OPEN-006 — Customers and Delivery applicability

```yaml
decision_id: OPEN-006
status: open
owner: Command Product Lead
target_phase: Phase 4A review
blocking: false
requirement_id: REQ-PROD-053
required_evidence: deployment scenarios, business model, generic reusable value, portal/branding/billing exclusions
affected_files: 06-command/modules/customers-and-delivery/*
```
### OPEN-007 — Human Gate and Govern relationship

```yaml
decision_id: OPEN-007
status: open
owner: Security Architecture
target_phase: Phase 4C
blocking: true
requirement_id: REQ-PROD-054
required_evidence: action class, authority, separation of duties and audit
affected_files: 08-govern/*; 09-cmdr-studio/human-gates/*; 07-investigate/modules/cases-and-evidence/capabilities/action-request-preparation.md
```
### OPEN-008 — Endpoint platform support

```yaml
decision_id: OPEN-008
status: open
owner: Endpoint Agent Product Lead
target_phase: Phase 4D
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
target_phase: Phase 4B.2
blocking: false
requirement_id: REQ-PROD-058
required_evidence: demand, legality and engine strategy
affected_files: 07-investigate/modules/analysis-workbench/*
```
### OPEN-012 — Cloud analysis scope

```yaml
decision_id: OPEN-012
status: open
owner: Investigate Product Lead
target_phase: Phase 4B.2
blocking: false
requirement_id: REQ-PROD-059
required_evidence: journeys, sources, APIs and ownership
affected_files: 07-investigate/modules/analysis-workbench/*
```
### OPEN-013 — Default governance for Action Class 2

```yaml
decision_id: OPEN-013
status: open
owner: Security Architecture
target_phase: Phase 4C/7
blocking: true
requirement_id: REQ-PROD-060
required_evidence: blast radius, rollback reliability, tenant policy, step-up and separation of duties
affected_files: 06-command/action-classification.md; 07-investigate/action-classification.md; 08-govern/*; 14-security-permissions-and-trust/*
```
### OPEN-014 — Artifact versus Attachment

```yaml
decision_id: OPEN-014
status: open
owner: Investigate Product Lead
target_phase: Phase 7
blocking: true
requirement_id: REQ-PROD-061
required_evidence: evidence lifecycle, provenance, collaboration, retention, permission, migration and UX tests
affected_files: 05-domain-model/objects/artifact.md; 07-investigate/modules/cases-and-evidence/artifact-versus-attachment.md; 07-investigate/modules/cases-and-evidence/capabilities/attachment-handling.md
```
### OPEN-015 — Automation Run to Response Run bridge

```yaml
decision_id: OPEN-015
status: open
owner: CMDR Studio Product Lead
target_phase: Phase 4C/7
blocking: true
requirement_id: REQ-PROD-062
required_evidence: action classification, Decision, Tool effects and audit
affected_files: 07-investigate/automation-and-ai-model.md; 08-govern/*; 09-cmdr-studio/*
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

## Disposition Phase 4A

- OPEN-006 reste ouverte : CAP-CMD-401 reste `proposed`, `planned` et deployment-dependent.
- OPEN-010 reste ouverte : aucune préférence finale de densité par rôle n’est validée.
- OPEN-013 reste ouverte : aucune politique générale de classe 2 n’est inventée.

## Disposition Phase 4B.1

### OPEN-014

Reste **ouverte et bloquante pour la phase Objets**. `artifact-versus-attachment.md` documente quatre options sans recommandation approuvée : objets distincts ; Attachment comme rôle/relation d’Artifact ; mécanisme Shared avec promotion ; convergence partielle avec lifecycles séparés. CAP-INV-106 reste `proposed` et `planned`.

### OPEN-007

Reste **ouverte**. CAP-INV-113 prépare et soumet une Action Request vers Govern ; Investigate ne crée ni Decision ni Approval et ne décide pas la relation finale Human Gate/Govern.

### OPEN-013

Reste **ouverte**. Les 22 capabilities classent leurs actions ; aucune règle de step-up ou Govern par défaut pour la classe 2 n’est finalisée.

### OPEN-015

Reste **ouverte**. Les sorties automatisées exposent Automation Run, Tool Calls et provenance lorsqu’ils existent ; aucune équivalence avec Response Run n’est décidée.

### OPEN-005, OPEN-011 et OPEN-012

Restent **ouvertes et reportées à Phase 4B.2**. Aucun moteur forensic, périmètre mobile ou cloud, API ou protocole n’est choisi en Phase 4B.1.

Aucun nouvel OPEN n’est créé et aucune décision n’est fermée.

## Historique résolu

`OPEN-009` — capability delivery-classification authority — reste la seule décision historiquement résolue : le capability owner propose, Product Lead et Product Architecture approuvent, Security/Engineering examinent selon l’effet et QA vérifie la preuve.

## Règles de maintenance

- ne jamais réutiliser un identifiant ;
- ne pas fermer une décision parce qu’une phase a commencé ;
- une capability `defined` avec delivery mode `planned` ne prouve aucune livraison ;
- toute nouvelle question nomme owner, target phase, dépendances, preuve et fichiers affectés.
