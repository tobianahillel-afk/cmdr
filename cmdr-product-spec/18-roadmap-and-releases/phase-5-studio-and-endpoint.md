---
id: roadmap-phase-5-studio-and-endpoint
domain: 18-roadmap-and-releases
status: draft
owner: Product Operations Lead
updated: 2026-08-10
source-of-truth: canonical
---
# Phase 5 Studio And Endpoint

## Objectif
Définir Phase 5 Studio and Endpoint pour CMDR tout en conservant Studio et Endpoint comme deux capability domains séparés.

## Périmètre
Document canonique Roadmap. Il ne remplace pas les sources propriétaires et ne transforme aucun execution lot en sous-phase numérotée.

## Propriétaire fonctionnel
Product Operations Lead.

## Objets concernés
- Concepts du document
- Références canoniques liées

## Fonctionnalités historiques prévues
- Skills/Agents/workflows.
- Assurance/control room.
- EDR telemetry/detection/response.
- Fleet administration comme dépendance Settings, sans transfert d'ownership vers Endpoint.

## UX et interactions
- Navigation par liens stables.
- Contenu lisible en thème clair et sombre.
- Aucune duplication des définitions externes.

## Permissions
Les modifications suivent `../14-security-permissions-and-trust/permission-model.md`; aucun execution lot ne finalise implicitement RBAC/ABAC.

## Dépendances
- `../00-governance/source-of-truth-policy.md`
- `dependency-roadmap.md`
- Studio/Settings/Govern/Shared/Endpoint ownership sources.

## STD-1 — verified
- **PASS AFTER POST-PUBLICATION VERIFICATION — 190/190**.
- `CAP-STD-001..016`: 16 / 432 / 96.

## STD-2 — verified
- **PASS AFTER POST-PUBLICATION VERIFICATION — 200/200**.
- `CAP-STD-017..033`: 17 / 459 / 102.
- fifth functional/build SHA `655e9ce0ade2d64a7738a6a479572fef9b6f0e2f`.

## STD-3 — current execution
- `CAP-STD-034..051`: **18 capabilities / 486 sections / 108 mandatory tables**.
- scope: Automation Agents/Teams, bounded autonomy/oversight, Human Gate runtime, Automation Run lifecycle/control, queue/scheduling semantics, failures/retries/partial, runtime context, Control Room, outcome/handoff and provenance.
- build-time status: **PENDING POST-PUBLICATION VERIFICATION — 202 PASS / 8 PENDING-REMOTE / 0 FAIL**.
- no new Screen IDs/detailed rewrites.
- no agent framework/model/provider/runtime/scheduler/API/protocol/code/final JSON Schema/RBAC.
- STD-3 remains documentary and PARTIAL until actual remote verification.

## Future lots
- STD-4 — Assurance & Lifecycle: NOT STARTED.
- Endpoint capability specification: NOT STARTED.
- Studio capability specification: PARTIAL.
- Delivery Roadmap Phase 5 overall: PARTIAL.

STD-1/2/3/4 are execution-lot labels only. No Phase 5A/5B/5C/5D exists.

## Ownership boundary
Automation Run ≠ Govern Response Run; Human Gate ≠ Govern Approval/Decision; Studio Runtime Outcome ≠ Govern Result. Shared owns generic Jobs/queue/scheduling mechanisms; Settings owns provider/secret/runtime administration; Endpoint owns technical endpoint primitives.

## Questions ouvertes
OPEN-007/013/015 remain open where consumed. OPEN-003 is palette-only and OPEN-008 remains Endpoint/platform support. Permission namespace ambiguity remains unresolved.

## Next candidate
STD-4 — Assurance & Lifecycle, **but it must not start until STD-3 post-publication verification is complete and explicitly authorized**.