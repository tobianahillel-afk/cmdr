---
id: roadmap-phase-6-platform-scale
domain: 18-roadmap-and-releases
status: draft
owner: Product Operations Lead
updated: 2026-08-12
source-of-truth: canonical
---
# Phase 6 Platform Scale

## Objectif

Définir phase 6 platform scale pour CMDR.

## Périmètre

Document canonique du domaine. Il définit uniquement son sujet et renvoie vers les autres sources de vérité pour les concepts partagés.

## Propriétaire fonctionnel

Product Operations Lead.

## Objets concernés

- Concepts du document
- Références canoniques liées

## Fonctionnalités

- MSSP aggregation.
- SLO/resilience.
- Localization.
- Advanced integrations.
- Compliance.

## Capability specification execution

### Tenant, Environment and Administrative Foundations

First functional Phase 6 execution lot. Exactly `CAP-SET-001..004`, owned by Platform Settings Product Lead, cover Tenant lifecycle/isolation, Environment lifecycle/Tenant scope, administrative change validation/provenance, and cross-product Tenant/Environment context semantics.

Structure: **4 capabilities / 108 numbered sections / 24 mandatory tables / at least 12 GWT**. No new Permission ID, Screen ID or canonical object is introduced. `CAP-SET-004` remains Platform Settings-owned; Experience Architecture owns only the propagation mechanism.

Functional build: `90684aaa9badf8ee76e54fdccd11bcd3e7fdde89`. Historical build-time gate state: **154 PASS / 6 PENDING-REMOTE / 0 FAIL**. Post-publication verification: **PASS — 160/160 PASS, 0 PENDING, 0 FAIL**.

### Identity Administration — Principals, Roles and Access Reviews

Second source-audited functional execution lot under the same Delivery Roadmap Phase 6 — Platform Scale. It is **not** Phase 6B and creates no new roadmap phase.

Exactly `CAP-SET-005..007`, owned by Platform Settings Product Lead, cover:

- Principal administrative lifecycle and human/service identity-state projection;
- Role lifecycle, constraints and the bounded typed Principal/Role relation boundary;
- periodic Access Review, evidence/provenance and keep/revoke disposition/handoff.

Structure: **3 capabilities / 81 numbered sections / 18 mandatory tables / at least 9 meaningful GWT**. Settings cumulative build content becomes **7 capabilities / 189 sections / 42 mandatory tables**.

Hard boundaries:

- Principal states remain `pending`, `active`, `suspended`, `revoked`;
- Role states remain `draft`, `active`, `deprecated`; Role expiry is a condition/constraint, never a state;
- Groups, Group Membership, generic Access Assignment and Effective Access are outside this lot;
- CAP-SET-007 uses revocation **disposition/handoff** because no generic canonical assignment-removal mechanic is sourced;
- Security retains Permission Model/RBAC/ABAC/authorization/tenant isolation/SoD/step-up;
- Govern retains Approval/Decision/Decision Authority/Response Run/Result;
- existing `SET-IAM-001` and `SET-AUD-001` are reused;
- new Permission IDs: **0**; new Screen IDs: **0**; new canonical objects: **0**.

Build-time quality state: **168 PASS / 6 PENDING-REMOTE / 0 FAIL**. The functional BUILD SHA is intentionally not predicted inside the fifth functional commit; post-publication evidence remains pending until remote publication and verification.

Global build content is expected and locally indexed as **491 capabilities / 489 defined / 2 proposed / 491 planned / 13257 sections / 2946 mandatory tables**. Requirements remain **122 = 99 conform / 20 partial / 3 absent / 0 contradictory** and OPEN remains **18**.

Settings Capability Specification remains **PARTIAL** and Delivery Roadmap Phase 6 Capability Specification remains **PARTIAL** because later integration/provider/secret, Shared/platform-scale and potentially Customer/MSSP/Delivery scope remain separate source-audited work.

## UX et interactions

- Navigation par liens stables.
- Contenu lisible en thème clair et sombre.
- Aucune duplication des définitions externes.
- Reuse existing Settings modules and screens; no Phase-6 navigation shell is created by Identity Administration.

## Permissions

Les modifications suivent le modèle défini dans `../14-security-permissions-and-trust/permission-model.md` lorsque le document décrit une capacité exécutable. Both completed/current Settings lots create zero new Permission IDs and zero new Screen IDs; if either becomes necessary in a future lot, that prerequisite requires a separate source-owned run.

## États

Le statut documentaire suit `00-governance/document-status-model.md`; les états métier restent dans leurs sources canoniques.

## Dépendances

- 00-governance/source-of-truth-policy.md
- `../10-platform-settings/capabilities/README.md`
- `../16-quality-and-validation/reports/platform-scale-tenant-environment-administrative-foundations-capability-conformance.md`
- `../16-quality-and-validation/reports/platform-scale-tenant-environment-administrative-foundations-post-publication-verification.md`
- `../16-quality-and-validation/reports/platform-scale-identity-administration-principals-roles-access-reviews-capability-conformance.md`

## Critères d’acceptation

- Le document a un propriétaire unique.
- Les liens locaux sont valides.
- Les décisions non tranchées sont attribuées.
- No Phase 6A/6B is created; execution lots remain subordinate to this roadmap phase.
- Documentary PASS does not claim implementation/runtime availability.

## Questions ouvertes

- Quelle date et quel owner doivent être confirmés?
- Quelle dépendance bloque ce jalon?
- Later source-confirmed Platform Scale lots require their own preparation/execution audit and cannot be inferred from this Identity lot.
