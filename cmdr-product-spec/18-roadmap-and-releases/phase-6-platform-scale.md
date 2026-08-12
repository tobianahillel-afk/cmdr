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

Settings Capability Specification remains **PARTIAL** and Delivery Roadmap Phase 6 Capability Specification is now **PARTIAL** because later Platform Scale lots remain. Identity Administration, integration/provider/secret, Shared/platform-scale and Customer/MSSP/Delivery scope are **NOT STARTED** by this lot.

## UX et interactions

- Navigation par liens stables.
- Contenu lisible en thème clair et sombre.
- Aucune duplication des définitions externes.

## Permissions

Les modifications suivent le modèle défini dans `../14-security-permissions-and-trust/permission-model.md` lorsque le document décrit une capacité exécutable. This first lot creates zero new Permission IDs and zero new Screen IDs; if either becomes necessary in a future lot, that prerequisite requires a separate run.

## États

Le statut documentaire suit `00-governance/document-status-model.md`; les états métier restent dans leurs sources canoniques.

## Dépendances

- 00-governance/source-of-truth-policy.md
- `../10-platform-settings/capabilities/README.md`
- `../16-quality-and-validation/reports/platform-scale-tenant-environment-administrative-foundations-capability-conformance.md`
- `../16-quality-and-validation/reports/platform-scale-tenant-environment-administrative-foundations-post-publication-verification.md`

## Critères d’acceptation

- Le document a un propriétaire unique.
- Les liens locaux sont valides.
- Les décisions non tranchées sont attribuées.
- No Phase 6A is created; execution lots remain subordinate to this roadmap phase.

## Questions ouvertes

- Quelle date et quel owner doivent être confirmés?
- Quelle dépendance bloque ce jalon?
