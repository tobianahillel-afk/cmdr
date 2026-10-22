---
id: CAP-STD-016
title: Studio Foundations Cross-Product Contracts and Provenance
product: cmdr-studio
module: studio-foundations
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-006, REQ-PROD-009, REQ-PROD-016, REQ-OBJ-009, REQ-AI-002]
open_decisions: [OPEN-007, OPEN-015]
source-of-truth: canonical
---
# CAP-STD-016 — Studio Foundations Cross-Product Contracts and Provenance
## 1. Définition
Contrats de frontière/provenance STD-1. Workflow ≠ Govern Playbook; Human Gate ≠ Approval; Automation Run ≠ Response Run; Tool Call output ≠ Govern Result/Evidence automatiquement.
## 2. Problème utilisateur
Un handoff Studio doit préserver source, owner, version, scope, authority et return-origin sans fusionner les objets des produits.
## 3. Objectifs
Formaliser Studio↔Govern/Investigate/Command/Endpoint/Settings/Shared et provenance bout-en-bout pour les concepts STD-1.
## 4. Non-objectifs
Aucun nouveau object lifecycle externe, response authority, Endpoint primitive, Shared engine ou Workflow/Agent runtime.
## 5. Propriétaire
CMDR Studio possède son handoff context/provenance locale; chaque destination reste owner de ses objets/actions.
## 6. Utilisateurs
Automation Designer, Studio Operator, Investigate Analyst, Incident Commander, Response Operator, Auditor.
## 7. Conditions d’entrée
Exact refs, tenant/env, caller, permission context, return-origin et authority ref si requise.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Studio refs | CMDR Studio | Tool/Call/Skill refs | oui | pinned/current | provenance partial |
| consumer context | product owner | source object/return | conditionnel | current | reject unscoped handoff |
| Govern refs | Govern | Decision/Run/Result refs | conditionnel | exact | keep distinct |
| Shared refs | Shared | Trace/Activity/Job | conditionnel | visible | provenance gap |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Decision/Response Run/Result | Govern | related refs | read |
| Evidence/Finding | Investigate | source/consumer refs | read |
| Trace/Activity/Job | Shared | provenance refs | consume |
| Endpoint capability | Endpoint Agent | future technical ref | read projection |
## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| cross-product handoff context | create reference package | CMDR Studio | ownership unchanged |
| Studio provenance projection | append refs | CMDR Studio | source history immutable |
## 11. Fonctionnalités
Owner validation, reference packaging, return-origin, caller/correlation, exact version, source restrictions and output qualification boundary.
## 12. Actions utilisateur
Class 0 inspect boundary/provenance; Class 2 prepare authorized handoff/reference. No external mutation is performed by this capability.
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| boundary validation | oui | oui | oui | oui | owner/ref checks |
| provenance summary | oui | oui | oui | oui | raw refs/Trace |
## 14. États fonctionnels
complete, partial, blocked, stale; cross-product source states are read-only projections.
## 15. États d’interface
Missing source/permission/authority is explicit; no protected source is leaked in denial.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| handoff package | reference context | destination | owner/exact refs/return preserved |
| STD-1 provenance chain | projection | auditor/consumer | no auto Evidence/Result |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Investigate | authorized Tool use | Studio | Case/Artifact refs + return | technical output to Investigate |
| Govern | execution handoff | Studio/Endpoint | Decision/Response Run refs + scope | technical status to Govern |
| Command | open Studio ref | Library | Incident context + return | Command |
## 18. Dépendances
Govern execution boundaries, Investigate Tool boundary, Settings refs, Shared Trace/Jobs/Linking and Security.
## 19. Source de vérité
Each product retains its canonical objects; this capability owns only Studio-side handoff/provenance semantics.
## 20. Provenance et audit
Actor/system, caller product/object, exact asset/version, tenant/env, authorization/Decision refs, correlation, timestamps, output refs and return-origin.
## 21. Permissions fonctionnelles
Cross-product reference read/prepare, technical output read and provenance export preparation; no authority transfer.
## 22. Limites et erreurs
Unresolvable ref, denied source, stale authority, missing return-origin, provenance gap or owner conflict blocks/marks Partial.
## 23. Métriques
Handoff completeness, owner conflicts, unresolvable refs, provenance gaps, stale authority refs.
## 24. Classification de livraison
`defined / planned`; contracts are provider/runtime-neutral documentation.
## 25. Critères d’acceptation
**Given** Human Gate/Tool Call context, **When** sent to Govern, **Then** no Approval/Decision is implicit.  
**Given** Investigate Tool output, **When** returned, **Then** it is not Evidence/Finding until Investigate qualifies it.  
**Given** IA off, **When** provenance is inspected, **Then** exact refs/correlation/owners remain usable.
## 26. Questions ouvertes
OPEN-007 and OPEN-015 remain open. STD-1 fixes only non-equivalence; it does not select the bridge semantics.
## 27. Consommateurs documentaires
Studio, Govern, Investigate, Command, Endpoint future refs, Settings, Shared, Security, Quality and Roadmap.
