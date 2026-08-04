---
id: 07-investigate-permissions
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-006
  - REQ-SEC-001
  - REQ-SEC-002
  - REQ-SEC-003
---

# Besoins fonctionnels de permissions Investigate

Ce document identifie les familles nécessaires sans créer de namespace final ni remplacer le Permission Model.

| Famille | Besoin fonctionnel | Risque ou séparation |
|---|---|---|
| Signal | read, triage, link | triage ne modifie pas priorité Command |
| Event Search | execute, cancel, raw access, source fields | raw et cross-tenant doivent être distincts |
| Saved Search | create, update, share, deprecate | partage tenant-scoped |
| Hunt | create, manage, contribute, archive | owner et contributeurs distincts |
| Case | read, create, update, close, reopen, assign | création, assignment et clôture séparables |
| Hypothesis | create, update, review | proposal agent jamais autorité |
| Entity relation | create, dispute, merge-propose | merge effectif nécessite contrôle distinct |
| Artifact | import, read, manage, export | contenu sensible et dérivés |
| Evidence | create, read, qualify, dispute, export | intégrité, review et export distincts |
| Finding | author, review, confirm, dispute | author/reviewer separation possible |
| Action Request | prepare, submit, follow | Govern traite ; requester ne s’auto-approuve pas |
| Reporting | prepare, review, publish, export, redact | audience et redaction distinctes |
| Collaboration | note/comment/mention/visibility | Case-scoped et tenant-scoped |

## Lacunes

Raw-event et field-level access, cross-environment/cross-tenant search, bulk linking, Evidence reviewer independence, Finding confirmation, Artifact/Attachment retention, Action Request submit/step-up et sensitive data reveal/copy/export restent à formaliser.

Les permissions de lecture des sources restent requises après liaison. Une projection, un Case ou un Report ne peut pas élargir l’accès. La matrice atomique et les namespaces restent Phase 7.
