---
id: ownership-register
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-OBJ-001
  - REQ-OBJ-002
  - REQ-OBJ-003
  - REQ-OBJ-004
  - REQ-OBJ-005
  - REQ-OBJ-006
  - REQ-OBJ-007
  - REQ-OBJ-008
  - REQ-OBJ-009
  - REQ-OBJ-010
  - REQ-OBJ-011
  - REQ-OBJ-012
---

# Registre de propriété canonique

Le propriétaire définit sémantique, états et changements. Un consommateur affiche, filtre, lie ou demande une transition autorisée ; il ne crée aucune seconde définition.

| Concept | Propriétaire | Fichier canonique | Consommateurs | Actions consommateurs | Restrictions | Requirement IDs |
|---|---|---|---|---|---|---|
| Incident | Command | `../05-domain-model/objects/incident.md` | Investigate, Govern, Reporting | read/link/project | no competing lifecycle | REQ-OBJ-001 |
| Task | Command unless declared specialized | `../05-domain-model/objects/task.md` | all products | create/assign/link under owner workflow | specialized tasks name owner | REQ-PROD-009 |
| Case | Investigate | `../05-domain-model/objects/case.md` | Command, Govern, Reporting | open/read/link | Command does not own Case | REQ-OBJ-002 |
| Hypothesis | Investigate | `../05-domain-model/objects/hypothesis.md` | Studio support, Reporting | propose/read/link | AI proposal attributed | REQ-OBJ-002 |
| Artifact | Investigate | `../05-domain-model/objects/artifact.md` | Command, Govern, Studio | preview/link/process | Artifact ≠ Evidence | REQ-PROD-061 |
| Evidence | Investigate | `../05-domain-model/objects/evidence.md` | Command, Govern, Reporting | read/cite/verify | no retroactive mutation | REQ-OBJ-003 |
| Finding | Investigate | `../05-domain-model/objects/finding.md` | Command, Govern, Studio | read/request action | Finding ≠ Decision/Result | REQ-OBJ-004 |
| Action Request | Govern lifecycle; Command/Investigate producers | `../05-domain-model/objects/action-request.md` | Command, Investigate, Govern | create/complete/follow | producer cannot self-approve | REQ-PROD-004 |
| Decision | Govern | `../05-domain-model/objects/decision.md` | Command, Investigate, Studio, Endpoint | provide context/read | no implicit agent decision | REQ-OBJ-005 |
| Approval | Govern | `../05-domain-model/objects/approval.md` | Studio Human Gate, Audit | submit/read | Approval ≠ Decision | REQ-OBJ-006 |
| Response Run | Govern | `../05-domain-model/objects/response-run.md` | Command, Investigate, Endpoint, Studio | observe/execute authorized/verify | Response Run ≠ Automation Run | REQ-OBJ-007 |
| Result | Govern | `../05-domain-model/objects/result.md` | Command, Investigate, Reporting | consume/link/measure | Result ≠ Finding | REQ-OBJ-007 |
| Endpoint | shared model; administered by Settings | Phase 7 canonical file | Command, Investigate, Govern, Endpoint | display/target under permission | no local redefinition | REQ-PROD-017 |
| Endpoint Agent Fleet | Platform Settings | `../05-domain-model/objects/endpoint-agent-fleet.md` | Command, Investigate, Govern | read/select projection | Investigate does not administer | REQ-OBJ-008 |
| Skill | CMDR Studio | `../05-domain-model/objects/skill.md` | operational products | invoke under permission | Skill ≠ Tool | REQ-OBJ-009 |
| Tool / Tool Call | CMDR Studio | Phase 7 canonical files | operational products, Govern, Audit | call/observe/interrupt under contract | no hidden call | REQ-OBJ-009 |
| Automation Agent / Team | CMDR Studio | `../05-domain-model/objects/automation-agent.md` | operational products | trigger/use | no business ownership | REQ-AI-002, REQ-OBJ-009 |
| Workflow | CMDR Studio | `../05-domain-model/objects/workflow.md` | all products | execute authorized version | does not replace owner workflow | REQ-OBJ-009 |
| Human Gate | CMDR Studio | `../05-domain-model/objects/human-gate.md` | Govern, operational products | provide validation | relation to Decision OPEN-007 | REQ-AI-004, REQ-OBJ-009 |
| Automation Run | CMDR Studio | Phase 7 canonical file | operational products, Audit | observe/interrupt/resume | bridge to Response Run OPEN-015 | REQ-PROD-062, REQ-OBJ-009 |
| Reporting Engine | Shared Capabilities | `../12-shared-capabilities/reporting-engine.md` | Command, Investigate, Govern | compose/review/publish | no local engine | REQ-OBJ-010 |
| Saved Views génériques | Shared Capabilities | `../12-shared-capabilities/saved-views.md` | all products | store/apply/version/share configuration | product owns system views | REQ-OBJ-011 |
| Work Queue Saved Views | Command | `../06-command/modules/incidents-and-work-queue/saved-views.md` | Command | define/apply six system views | not six pages; no generic engine duplication | REQ-OBJ-012, REQ-UX-008 |
| Permission Model | Security, Permissions and Trust | `../14-security-permissions-and-trust/permission-model.md` | all products | reference/enforce | no local namespace | REQ-PROD-006 |
| Inspector | Design System | `../03-design-system/components/inspector.md` | all products | configure sections/actions | one right Inspector | REQ-UX-002 |
| Context Bar component | Design System | `../03-design-system/components/context-bar.md` | all products | render authorized context | does not own propagation | REQ-UX-006 |
| Context preservation | Experience Architecture | `../04-experience-architecture/context-preservation.md` | all products | transmit/restore references | never expands permission | REQ-UX-006, REQ-UX-007 |

## Résolution Phase 3 — Saved Views

La contradiction est résolue : Shared Capabilities possède le mécanisme générique dans `../12-shared-capabilities/saved-views.md`; Command possède le catalogue Work Queue dans `../06-command/modules/incidents-and-work-queue/saved-views.md`; le Design System possède l'interaction dans `../03-design-system/components/saved-views.md`. Les six vues ne sont pas des écrans.

## Règle de changement

Tout changement de propriétaire exige ADR, registre, sources, permissions, consommateurs et matrice. Une projection dans l'Inspector ou le Context Bar ne transfère pas la propriété.

## Critère d’acceptation

**Given** `My Work` dans la Work Queue,  
**When** la chaîne de propriété est inspectée,  
**Then** Command possède le contenu de la vue, Shared Capabilities sa persistance, le Design System son composant, et aucun fichier d'écran ne revendique la vue.
