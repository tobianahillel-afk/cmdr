---
id: roadmap-phase-6-platform-scale
domain: 18-roadmap-and-releases
status: draft
owner: Product Operations Lead
updated: 2026-08-13
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

Structure: **3 capabilities / 81 numbered sections / 18 mandatory tables / 12 meaningful GWT**. Settings cumulative state at Identity closure is **7 capabilities / 189 sections / 42 mandatory tables**.

Hard boundaries:
- Principal states remain `pending`, `active`, `suspended`, `revoked`;
- Role states remain `draft`, `active`, `deprecated`; Role expiry is a condition/constraint, never a state;
- Groups, Group Membership, generic Access Assignment and Effective Access are outside this lot;
- CAP-SET-007 uses revocation **disposition/handoff** because no generic canonical assignment-removal mechanic is sourced;
- Security retains Permission Model/RBAC/ABAC/authorization/tenant isolation/SoD/step-up;
- Govern retains Approval/Decision/Decision Authority/Response Run/Result;
- existing `SET-IAM-001` and `SET-AUD-001` are reused;
- new Permission IDs: **0**; new Screen IDs: **0**; new canonical objects: **0**.

Functional BUILD: `75a1fdeff9acc589c13773e95f8953ceeb29edd3`.

Historical build-time quality state is **168 PASS / 6 PENDING-REMOTE / 0 FAIL**. Remote publication verification completed all six remote gates: **PASS AFTER POST-PUBLICATION VERIFICATION — 174/174 PASS, 0 PENDING, 0 FAIL**.

Remote closure evidence confirms baseline `d605265f5b8a2e4350388b4ec9cfe51920a4aa50` → BUILD = **5 ahead / 0 behind**, same merge-base; Requirements **122 = 99/20/3/0**; OPEN **18**; PR #2 open/Draft/unmerged on `main`; `main` and branch/main README unchanged; CI N/A with its recorded evidence.

### Secrets & Connections — Integration and Secret Reference Administration

Third source-audited functional execution lot under the same Delivery Roadmap Phase 6 — Platform Scale. It is **not** Phase 6B and creates no new roadmap phase.

Exactly `CAP-SET-008..009`, owned by Platform Settings Product Lead, cover:
- `CAP-SET-008` — Integration administrative lifecycle, deterministic local validation, connection-test administrative request/preconditions/status/result projection, safe disable and provenance;
- `CAP-SET-009` — Secret Reference reference-only administration, lifecycle, rotation/expiry/revocation of the reference, Security handoff and provenance.

Structure: **2 capabilities / 54 numbered sections / 12 mandatory tables / 12 meaningful GWT**. Settings cumulative state becomes **9 capabilities / 243 sections / 54 mandatory tables**.

Hard boundaries:
- canonical Integration states remain `draft`, `validating`, `active`, `degraded`, `disabled`, `error`;
- canonical Secret Reference states remain `pending`, `active`, `rotating`, `expired`, `revoked`;
- Integration `capabilities` metadata is not the CMDR Capability object and cannot allocate `CAP-*` IDs;
- `Connection` remains module/functional terminology; no Connection, Connector or Credential canonical object is created;
- `Secret Reference` remains reference-only and no raw secret value is read, logged or exposed by the contracts;
- current canonical sources do **not** assign Platform Settings the technical external connection-probe executor; CAP-SET-008 is limited to administrative request/preconditions/status/result projection/provenance and handoff for external testing;
- current canonical sources do **not** assign Platform Settings underlying-secret generation/write/rotation or external credential revocation; CAP-SET-009 defines Secret Reference lifecycle/reference mutation and administrative handoff/result projection only;
- Tenant is mandatory; Environment is source-dependent only;
- existing `SET-SEC-001` and `SET-AUD-001` are reused;
- new Permission IDs: **0**; new Screen IDs: **0**; new canonical objects: **0**;
- Models & Providers and Sources & Parsers remain separate later source-audited lots;
- `CAP-SET-010+` is neither allocated nor reserved.

Exact functional chain:
1. `9b49b5822a21ecee4dec1a4a8c0614024123d129` — `docs: establish Settings Secrets and Connections capability ownership and boundaries`;
2. `641389f1a3432e4162fc7cc623c5982b2252b584` — `docs: define Settings integration lifecycle validation and connection state`;
3. `f3b75ddd4c8efc939db701a980703389d917ab17` — `docs: specify Settings secret reference lifecycle rotation revocation and security boundaries`;
4. `7d3e49b54e42709abe62f41ecb31b85c4b824d87` — `docs: update Settings Secrets and Connections traceability and quality gates`.

Functional BUILD tree: `dd40d275d72c943486edc72a8155001a060cc7cc`. Baseline `a5c450528ad25f090832adcda4ae37c79bea072b` → BUILD is **4 ahead / 0 behind**, same merge-base. Publication used non-forced fast-forward.

Global state after the functional lot is **493 capabilities / 491 defined / 2 proposed / 493 planned / 13,311 numbered sections / 2,958 mandatory tables**. Requirements remain **122 = 99 conform / 20 partial / 3 absent / 0 contradictory**. OPEN remains **18**; `OPEN-008`, `OPEN-012` and `OPEN-013` remain unresolved where relevant.

Historical build-time conformance is preserved as **166 PASS / 6 PENDING-REMOTE / 0 FAIL**.

Post-publication verification then executed all six remote-dependent gates:
- remote branch exactly BUILD;
- ancestry/diff = 4 ahead / 0 behind, same merge-base and exact expected 10 functional surfaces;
- remote CAP-SET contracts, shard/global register, Requirements, OPEN, Screen Register, Permission Register, canonical Integration/Secret Reference objects and quality/roadmap surfaces re-read;
- PR #2 remained open / Draft / unmerged, base `main`, head BUILD, auto-merge disabled;
- `main` remained `bc1ec59e5e79ccbaba291e2984b0eb7d5e54129c` and both README remained exact `# cmdr` with blob `901c74cda52e28b5ff7fc425ddf28ef89f3ad875`;
- CI/check/workflow was **N/A with evidence**: 0 statuses, 0 check runs, 0 check suites, 0 workflow runs and no `.github/workflows` directory.

Final lot verdict: **PASS AFTER POST-PUBLICATION VERIFICATION — 172/172 PASS, 0 PENDING, 0 FAIL**.

Settings Capability Specification remains **PARTIAL** and Delivery Roadmap Phase 6 Capability Specification remains **PARTIAL** because later Models & Providers, Sources & Parsers and other source-audited Platform Scale work remain separate. Global Capability Specification and repository maturity remain **PARTIAL**. Documentary PASS does not claim implementation/runtime availability.

## UX et interactions

- Navigation par liens stables.
- Contenu lisible en thème clair et sombre.
- Aucune duplication des définitions externes.
- Reuse existing Settings modules and screens; no Phase-6 navigation shell is created by these Settings lots.

## Permissions

Les modifications suivent le modèle défini dans `../14-security-permissions-and-trust/permission-model.md` lorsque le document décrit une capacité exécutable. The three Settings execution lots through Secrets & Connections create zero new Permission IDs and zero new Screen IDs; if either becomes necessary in a future lot, that prerequisite requires a separate source-owned run.

## États

Le statut documentaire suit `00-governance/document-status-model.md`; les états métier restent dans leurs sources canoniques.

## Dépendances

- 00-governance/source-of-truth-policy.md
- `../10-platform-settings/capabilities/README.md`
- `../16-quality-and-validation/reports/platform-scale-tenant-environment-administrative-foundations-capability-conformance.md`
- `../16-quality-and-validation/reports/platform-scale-tenant-environment-administrative-foundations-post-publication-verification.md`
- `../16-quality-and-validation/reports/platform-scale-identity-administration-principals-roles-access-reviews-capability-conformance.md`
- `../16-quality-and-validation/validation-status-platform-scale-identity-administration.md`
- `../16-quality-and-validation/quality-index-platform-scale-identity-administration.md`
- `../16-quality-and-validation/reports/platform-scale-secrets-connections-integration-secret-reference-capability-conformance.md`
- `../16-quality-and-validation/validation-status-platform-scale-secrets-and-connections.md`
- `../16-quality-and-validation/quality-index-platform-scale-secrets-and-connections.md`

## Critères d’acceptation

- Le document a un propriétaire unique.
- Les liens locaux sont valides.
- Les décisions non tranchées sont attribuées.
- No Phase 6A/6B is created; execution lots remain subordinate to this roadmap phase.
- Documentary PASS does not claim implementation/runtime availability.

## Questions ouvertes

- Quelle date et quel owner doivent être confirmés?
- Quelle dépendance bloque ce jalon?
- Later source-confirmed Platform Scale lots require their own preparation/execution audit and cannot be inferred from Secrets & Connections.
