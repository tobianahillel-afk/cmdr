---
id: CAP-STD-009
title: Tool Call Lifecycle, Technical Outcome and Provenance
product: cmdr-studio
module: studio-foundations
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-005, REQ-PROD-016, REQ-OBJ-009, REQ-SEC-002]
open_decisions: [OPEN-015]
source-of-truth: canonical
---
# CAP-STD-009 — Tool Call Lifecycle, Technical Outcome and Provenance
## 1. Définition
Lifecycle fonctionnel d'un Tool Call et réconciliation de son technical outcome/provenance. Tool Call ≠ Automation Run ≠ Response Run ≠ Job; technical output ≠ canonical Result.
## 2. Problème utilisateur
Le caller doit distinguer prepared/accepted/running/partial/completed/failed et conserver exact Tool version/inputs/runtime sans fabriquer un outcome métier.
## 3. Objectifs
Rendre state, progress, terminal technical outcome, warnings/errors, produced refs et trace attribuables et navigables.
## 4. Non-objectifs
Pas d'Automation Run orchestration, Workflow retries/compensation, Govern Result, business success ou generic Job lifecycle.
## 5. Propriétaire
CMDR Studio possède Tool Call functional lifecycle; runtime owner fournit raw status/output; Shared possède generic Trace/Activity/Job.
## 6. Utilisateurs
Studio Operator, Automation Designer, caller analyst/operator et Auditor.
## 7. Conditions d’entrée
Attributed CAP-STD-008 request, runtime correlation and exact Tool/version.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Tool Call request | CAP-STD-008 | attributed request | oui | exact | reject status |
| runtime status | runtime owner | state/progress | oui | event time | stale/unknown |
| technical output | runtime owner | output/warnings/errors | conditionnel | event time | explicit missing |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Tool | CMDR Studio | exact version/output contract | read |
| Trace | Shared Capabilities | correlation refs | consume/read |
| Response Run | Govern | distinct related ref if any | read only |
## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Tool Call lifecycle projection | advance state | CMDR Studio | not Automation Run orchestration |
| technical outcome/provenance | record refs | CMDR Studio | not Govern Result |
## 11. Fonctionnalités
prepared→validation-pending→ready→queued→accepted→running→partial/completed/failed/timeout/cancelled/rejected/unsupported/permission-blocked/runtime-unavailable with exact provenance.
## 12. Actions utilisateur
Class 0 inspect status/output/provenance; Class 2 request cancel only where runtime contract/authority allows. Retry semantics deferred STD-2/3.
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| state reconciliation | oui | oui | oui | oui | event/status reconciliation |
| error explanation | oui | oui | oui | oui | raw error/status |
## 14. États fonctionnels
prepared, validation-pending, ready, queued, accepted, running, partial, completed, failed, timeout, cancelled, rejected, unsupported, permission-blocked, runtime-unavailable.
## 15. États d’interface
Stale distinguishes last-known runtime state; Error/Partial preserve valid prior data and correlation.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Tool Call state | projection | caller/future Control Room | exact call/version/state |
| technical outcome | technical result | caller/Trace | no automatic Result/Evidence |
| provenance chain | refs | Auditor | caller/inputs/runtime/times retained |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| runtime owner | status event | Studio Tool Call | call/correlation/status/time | return-origin unchanged |
| Studio Tool Call | terminal outcome | caller | technical output/provenance/return | caller decides qualification |
## 18. Dépendances
CAP-STD-005/008/015/016, runtime owner, Shared Trace/Activity and OPEN-015.
## 19. Source de vérité
Studio owns Tool Call lifecycle projection; runtime source owns raw execution state; Govern owns Response Run/Result.
## 20. Provenance et audit
Exact Tool/version, input refs, runtime/provider, caller, timestamps, status history, warnings/errors, produced refs, Trace and correlation are append-preserved.
## 21. Permissions fonctionnelles
Tool Call read, restricted technical output read, cancellation request need and provenance/export preparation; no final RBAC.
## 22. Limites et erreurs
Stale/missing runtime status, duplicate/out-of-order event, partial output, timeout, denied output or unresolvable produced ref remain explicit.
## 23. Métriques
Lifecycle distribution, duration concept, timeout/failure/partial rate, stale status, provenance completeness; no SLO target.
## 24. Classification de livraison
`defined / planned`; no runtime/state engine implemented.
## 25. Critères d’acceptation
**Given** Tool Call technical success, **When** completed, **Then** technical output available and no Govern Result/Evidence auto-created.  
**Given** Tool Call failure, **When** reviewed, **Then** error/provenance visible and no Incident automatic.  
**Given** runtime status stale, **When** inspected, **Then** last-known state/freshness visible without fabricated completion.
## 26. Questions ouvertes
OPEN-015 remains open. Automation Run semantics are not started in STD-1.
## 27. Consommateurs documentaires
Caller products, future STD-2/3 Control Room, Govern handoff, Shared Trace, Security, Quality.
