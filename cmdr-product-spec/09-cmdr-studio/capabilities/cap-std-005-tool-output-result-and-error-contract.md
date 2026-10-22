---
id: CAP-STD-005
title: Tool Output, Result and Error Contract
product: cmdr-studio
module: studio-foundations
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-002, REQ-PROD-005, REQ-PROD-016, REQ-OBJ-009]
open_decisions: []
source-of-truth: canonical
---
# CAP-STD-005 — Tool Output, Result and Error Contract
## 1. Définition
Sémantique fonctionnelle des technical outputs, structured/human-readable projections, warnings, partial outputs et errors. Tool output ≠ Govern Result/Evidence/Finding automatiquement.
## 2. Problème utilisateur
Un succès technique peut être partiel, stale ou insuffisant et ne doit jamais être promu en outcome métier.
## 3. Objectifs
Rendre outputs/errors attribués, versionnés, consommables et explicitement limités.
## 4. Non-objectifs
Pas de Result Govern, Evidence/Finding creation, Incident auto, output schema final ou moteur runtime.
## 5. Propriétaire
Studio possède le Tool output contract; runtime owner produit le technical output; consuming product qualifie ses objets.
## 6. Utilisateurs
Studio Operator, Automation Designer, Investigate Analyst, consuming product user.
## 7. Conditions d’entrée
Tool/version et Tool Call context connus; runtime output attribuable et correlation présente.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Tool version | Studio | exact ref | oui | pinned | bloque interpretation |
| technical output | runtime owner | raw/structured outcome | oui | execution time | missing/partial explicit |
| call provenance | Studio | caller/correlation | oui | exact call | unattributed error |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Tool | CMDR Studio | output/error contract | read |
| Artifact | Investigate | optional reference | link only |
| Result | Govern | distinct related ref | read only |
## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Tool output contract | define/update | CMDR Studio | technical semantics only |
| technical outcome projection | derive | CMDR Studio | no auto Evidence/Finding/Result |
## 11. Fonctionnalités
Structured/human output, artifact refs, warnings, errors, partial, timeout, unsupported, permission denied, runtime unavailable and provenance.
## 12. Actions utilisateur
Class 0 inspect output/error; Class 1 validate contract; Class 2 link output reference to consumer context without requalification.
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| output validation | oui | oui | oui | oui | contract checks |
| error explanation | oui | oui | oui | oui | raw error/status |
## 14. États fonctionnels
complete, partial, warning, failed, timeout, unsupported, permission-blocked, runtime-unavailable.
## 15. États d’interface
Partial/Error preserve valid output and correlation; no UI state changes semantic qualification.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| technical output | projection | caller | attributed/exact Tool version |
| error/warning classification | technical status | caller/Trace | no business-outcome claim |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| runtime | return outcome | Tool Call | call/version/output refs | caller |
| Tool Call | link output | consumer | technical ref/provenance | consumer qualifies |
## 18. Dépendances
CAP-STD-003/008/009, runtime owner, Shared Trace and consumer owner semantics.
## 19. Source de vérité
Tool contract is Studio; technical raw result is runtime-owned; Evidence/Finding/Result remain external owner objects.
## 20. Provenance et audit
Exact Tool/version, call/caller, input refs, runtime, timestamps, output hashes/refs, warning/error and correlation.
## 21. Permissions fonctionnelles
Output/error read, restricted output masking and provenance read/export preparation.
## 22. Limites et erreurs
Missing output, parse/validation failure, partial, timeout, unsupported, stale or restricted output remain explicit.
## 23. Métriques
Complete/partial/error/timeout rates, output validation failures, provenance completeness; no KPI threshold.
## 24. Classification de livraison
`defined / planned`; no runtime/output engine implemented.
## 25. Critères d’acceptation
**Given** Tool Call succeeds, **When** output consumed, **Then** no Govern Result/Evidence/Finding auto-created.  
**Given** partial output + warning, **When** reviewed, **Then** both remain explicit.  
**Given** IA off, **When** error inspected, **Then** structured error/provenance suffice.
## 26. Questions ouvertes
No new OPEN; Tool Call/Automation Run bridge remains handled by OPEN-015 elsewhere.
## 27. Consommateurs documentaires
Tool Calls, Skills, Investigate/Command consumers, Govern references, Shared Trace, Quality.
