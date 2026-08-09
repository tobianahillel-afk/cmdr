---
id: CAP-STD-023
title: Data Mapping, Transformation and Context Propagation
product: cmdr-studio
module: workflows
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-10
requirement_ids: [REQ-PROD-006, REQ-PROD-016, REQ-PROD-019, REQ-OBJ-009, REQ-SEC-001, REQ-AI-002]
open_decisions: [OPEN-013]
source-of-truth: canonical
---
# CAP-STD-023 — Data Mapping, Transformation and Context Propagation
## 1. Définition
Contrat fonctionnel de mapping entre source output et target input : field selection/rename/projection, deterministic transformation concept, validation, missing/incompatible handling, sensitive propagation and provenance. Mapping ≠ source-data mutation.
## 2. Problème utilisateur
Sans contrat, un Workflow peut altérer implicitement une source, perdre provenance ou utiliser un mapping non compatible/sensible.
## 3. Objectifs
Définir mappings et transformations déterministes composables, validables et traçables sans choisir de langage.
## 4. Non-objectifs
Aucun transformation language, ETL engine, source mutation, secret reveal, JSON schema, Evidence qualification ou runtime engine.
## 5. Propriétaire
CMDR Studio Product Lead. Studio owns only the Workflow/Builder orchestration semantics described here; referenced products retain their canonical objects and authority.
## 6. Utilisateurs
Automation Designer principal; Studio Reviewer et, selon le handoff, Studio Operator, Command/Investigate analyst, Response Operator ou Auditor autorisés.
## 7. Conditions d’entrée
Tenant et environnement résolus, Workflow/version ou draft identifiable, actor authentifié, permission context disponible et références requises explicitement résolues.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| source output contract | Tool/Skill/step/Workflow | typed functional output | oui | exact version | mapping blocked |
| target input contract | Tool/Skill/step/Workflow | typed functional input | oui | exact version | mapping blocked |
| mapping definition | Workflow draft | selection/rename/transform | oui | draft | no implicit transform |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Workflow | CMDR Studio | mapping context | read/write |
| Tool/Skill I-O | CMDR Studio | source/target contracts | read |
| Secret Reference | Platform Settings | opaque reference metadata | restricted read |
## 10. Objets créés ou modifiés
| Objet/concept | Opération | Propriétaire | Règle |
|---|---|---|---|
| Data Mapping definition | create/update conceptual | CMDR Studio | part of Workflow, no physical schema |
| mapping validation assessment | derive | CMDR Studio | no source mutation |
## 11. Fonctionnalités
Select/project/rename fields; deterministic transform concept; target binding; missing/incompatible handling; sensitive metadata propagation; lineage.
## 12. Actions utilisateur
| Action | Rôle | Objet/concept | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Inspect mapping | Automation Designer | Data Mapping | 0 | read | mapping visible | non |
| Validate mapping | Studio Reviewer | Data Mapping | 1 | source/target contracts | assessment | non |
| Edit mapping | Automation Designer | Workflow draft | 2 | manage | draft mapping changed | OPEN-013 |
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| suggest mapping | oui | oui | oui | oui | manual mapping |
| validate types | oui | oui | oui | oui | deterministic contract checks |

AI is optional. No essential capability in this contract requires a chatbot or model provider; AI suggestions remain reviewable, attributable and non-authorizing.
## 14. États fonctionnels
unmapped, partial, valid, missing-source, missing-target, incompatible-type, restricted, deprecated-dependency.
## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied and Stale expose missing source/freshness/permission explicitly. UI state never changes functional ownership or authorization.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| validated mapping | Workflow mapping | next step/CAP-STD-030 | source-target lineage preserved |
| mapping error | assessment | Builder | exact field/ref/reason |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| step output | map | target step input | source ref/transform/target ref | Workflow |
| mapping validation | fail | Builder error context | source/target/reason | Workflow |
| Workflow output mapping | return | consumer | output ref + lineage | consumer |
## 18. Dépendances
CAP-STD-004/005/012/019/021/030; Settings sensitive refs; Security source/target access.
## 19. Source de vérité
Source objects remain source-owned; Studio owns only mapping definition and derived Workflow context. Consumer owns later qualification such as Evidence.
## 20. Provenance et audit
Record source step/output/version, target step/input/version, transform descriptor concept, sensitivity, actor, validation result and correlation.
## 21. Permissions fonctionnelles
Workflow mapping read/edit/validate; source/target metadata visibility; sensitive-binding metadata read. No permission inheritance through mapping. `perm.cmdr-studio.*` and `perm.studio.*` remain coexisting historical namespaces; STD-2 performs no bulk rename and final RBAC/ABAC remains future.
## 22. Limites et erreurs
Missing field, incompatible type, denied metadata, sensitive leak risk, stale version or unsupported transform concept blocks/marks Partial.
## 23. Métriques
Mapping validation failures, incompatible types, missing fields, sensitive-flow blocks, provenance completeness.
## 24. Classification de livraison
`defined / planned`. Documentary definition does not prove implementation, publishing, deployment, runtime availability or production execution.
## 25. Critères d’acceptation
**Given** a source output type is incompatible with the target input, **When** mapping validation runs, **Then** the mapping is blocked with explicit source/target types.

**Given** mapping handles sensitive data, **When** it is inspected, **Then** only authorized metadata/references are shown and raw secrets remain hidden.

**Given** AI is disabled, **When** a mapping is authored, **Then** manual field mapping and deterministic validation remain available.
## 26. Questions ouvertes
OPEN-013 remain open and are not resolved by this capability.
## 27. Consommateurs documentaires
Workflow/Builder, Tool/Skill I-O contracts, CAP-STD-019/021/030/033, Security/Settings, Quality.
