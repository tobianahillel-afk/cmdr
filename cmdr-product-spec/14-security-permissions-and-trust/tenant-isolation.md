---
id: security-tenant-isolation
domain: 14-security-permissions-and-trust
status: draft
owner: Security Architecture Lead
updated: 2026-08-14
source-of-truth: canonical
requirements:
  - REQ-PROD-006
  - REQ-PROD-008
  - REQ-PROD-009
  - REQ-SEC-001
---
# Tenant Isolation

## Objectif

Définir l'isolation Tenant et les règles MSSP d'agrégation autorisée sans créer une hiérarchie de Tenants ni transférer permission, ownership ou Decision Authority.

## Périmètre

Ce document est la source Security pour :
- isolation storage/query/cache/export ;
- résolution et consommation d'un `Authorized Tenant Set` ;
- règles d'agrégation MSSP read-only ;
- refus cross-scope ;
- absence de fallback Tenant ;
- conditions de changement explicite de Tenant.

La décision de déploiement Customers/MSSP/Delivery est `../00-governance/adr/ADR-0008-customers-mssp-delivery-deployment-and-cross-tenant-architecture.md`.

## Propriétaire fonctionnel

Security Architecture Lead.

Platform Settings reste owner du Tenant canonique. Security possède l'autorisation et l'isolation. Experience Architecture possède le mécanisme de conservation/changement de contexte. Govern possède Decision Authority et la réponse gouvernée.

## Tenant boundary

Un Tenant reste une frontière d'isolation organisationnelle indépendante.

Invariants :
- chaque objet tenant-scoped conserve son propre Tenant ;
- un Tenant n'est ni Customer ni enfant d'un Customer ;
- aucun parent/child Tenant, `ManagedTenant`, `TenantGroup`, `Portfolio`, `CustomerTenant` ou équivalent canonique n'est introduit pour le MVP ;
- une relation ou projection ne transfère ni propriété ni permission ;
- aucun fallback Tenant silencieux n'est autorisé.

## Authorized Tenant Set

`Authorized Tenant Set` est une **projection d'autorisation non canonique**, résolue par Security pour le Principal, le deployment et le contexte d'autorisation courants.

Il représente un ensemble explicite de références Tenant que le Principal peut tenter de consulter selon ses permissions courantes.

Il n'est pas :
- un objet canonique ;
- un Group ;
- un Role ;
- un AccessAssignment, RoleAssignment ou PermissionAssignment ;
- un Effective Access object ;
- un Portfolio ;
- une Decision Authority ;
- une preuve autonome de droit d'accès à chaque objet.

La règle de set est :

`requested_tenant_set ⊆ authorized_tenant_set`.

Même lorsque cette relation est vraie, chaque objet reste soumis à son permission check serveur RBAC/ABAC.

## MSSP read-only aggregation

Le MVP autorise une agrégation read-only uniquement sur l'Authorized Tenant Set explicitement résolu.

Chaque résultat agrégé conserve au minimum :
- référence Tenant ;
- référence/identité stable de l'objet selon son owner ;
- provenance/source pertinente ;
- état de permission/fraîcheur lorsque nécessaire à une décision sûre.

L'agrégation ne donne aucun droit de mutation.

Sont interdits depuis le contexte agrégé :
- cross-tenant mutation ;
- cross-tenant administration ;
- cross-tenant response ;
- delegated administration ;
- automatic export widening.

## Tenant context switching

Le changement vers un Tenant déjà autorisé est une opération de navigation/contexte, pas un grant.

Après sélection d'un Tenant :
1. tout contexte incompatible est effacé selon Experience Architecture ;
2. Security réévalue les permissions dans le Tenant sélectionné ;
3. les produits consommateurs résolvent ensuite leurs objets dans ce Tenant ;
4. aucune permission issue du contexte agrégé n'est héritée comme droit de mutation.

## Search, Reporting et Export

Pour le MVP :
- Search est single-selected-Tenant ;
- Report est single-Tenant ;
- Export est single-Tenant et ne peut jamais élargir la visibilité.

Multi-tenant Search, Reporting et Export sont différés. Shared conserve l'ownership de ces mécanismes.

Un Authorized Tenant Set ne peut pas être utilisé pour contourner cette limite initiale.

## Administration

L'administration reste tenant-local.

Pour administrer un Tenant, un Principal doit :
- sélectionner ce Tenant explicitement ;
- satisfaire les permissions Settings applicables ;
- satisfaire le contexte RBAC/ABAC courant ;
- respecter les contrôles d'audit et de step-up applicables.

Aucune permission générique `cross-tenant.manage` n'est créée.

## Response Authority

Une réponse MSSP exige :
1. sélection explicite d'un Tenant ;
2. réévaluation Security dans ce Tenant ;
3. évaluation Govern Decision Authority dans ce Tenant ;
4. traitement des Action Request / Decision / Response Run selon leurs propres règles.

Aucune autorité de réponse centralisée cross-tenant n'est créée et l'administration de plateforme n'accorde jamais automatiquement cette autorité.

## Audit

Les évaluations d'accès, refus et tentatives cross-scope sont auditées.

Pour une agrégation multi-tenant, la trace doit permettre d'attribuer sans exposer de données protégées inutiles :
- acteur/Principal ;
- scope Tenant demandé ;
- outcome d'autorisation ;
- action ou type de lecture ;
- timestamp ;
- correlation id ;
- refus et cause sûre lorsque applicable.

## Negative tests obligatoires

Au minimum :
1. un Tenant non présent dans l'Authorized Tenant Set n'apparaît pas dans l'agrégation ;
2. un objet sans permission read reste invisible même si son Tenant est dans le set ;
3. un changement de Tenant ne conserve pas un objet incompatible ;
4. aucune mutation depuis l'overview agrégé n'est acceptée ;
5. aucune administration cross-tenant n'est acceptée ;
6. aucune réponse cross-tenant n'est acceptée ;
7. un export ne peut pas étendre le scope ;
8. aucun fallback Tenant n'est choisi en cas de contexte absent/invalide ;
9. les tentatives cross-scope produisent une trace d'audit sans fuite de donnée.

## Permissions

Les permissions atomiques existantes sont réutilisées avec RBAC/ABAC et scope Tenant. Cette architecture crée **0 nouvel ID de permission**.

La sémantique détaillée reste dans `permission-model.md` et les identifiants restent dans `permission-catalog.md`.

## Dépendances

- `../00-governance/adr/ADR-0008-customers-mssp-delivery-deployment-and-cross-tenant-architecture.md`
- `permission-model.md`
- `decision-authority.md`
- `separation-of-duties.md`
- `../04-experience-architecture/context-preservation.md`
- Tenant canonique Platform Settings
- Shared Global Search / Reporting / Export

## OPEN préservées

- `OPEN-013` reste ouverte pour le default governance des Class-2 mutations ;
- `OPEN-019` reste ouverte pour dissemination/releasability/external or client-facing sharing.

L'agrégation read-only interne ne ferme aucune de ces décisions.

## Critères d'acceptation

**Given** un Principal avec un Authorized Tenant Set explicite, **When** un overview MSSP est demandé, **Then** seuls les objets autorisés des Tenants du set sont retournés et chaque objet conserve son Tenant.

**Given** un contexte agrégé, **When** une mutation, administration ou réponse est demandée, **Then** l'action est bloquée jusqu'à sélection d'un Tenant unique et nouvelle évaluation d'autorisation/authority.

**Given** un Tenant non autorisé, **When** une URL, recherche ou référence tente de le résoudre, **Then** la résolution est refusée sans fallback ni fuite cross-tenant et la tentative est auditable.
