---
id: ADR-0008-customers-mssp-delivery-deployment-and-cross-tenant-architecture
domain: 00-governance
status: validated
owner: Product Architecture
updated: 2026-08-14
source-of-truth: canonical
requirements:
  - REQ-PROD-006
  - REQ-PROD-008
  - REQ-PROD-009
  - REQ-PROD-012
  - REQ-PROD-013
  - REQ-PROD-019
  - REQ-PROD-033
  - REQ-PROD-053
  - REQ-SEC-001
---
# ADR-0008 — Customers, MSSP, Delivery Deployment and Cross-Tenant Architecture

## 1. Identifiant

`ADR-0008`

## 2. Statut et approbation

**Validated — 2026-08-14.**

Approval reference: **Hillel Tobiana — explicit project-owner approval in ChatGPT conversation**.

Cette ADR enregistre la décision humaine explicite qui résout `OPEN-006`. Elle approuve une architecture documentaire et fonctionnelle ; elle ne prouve aucune implémentation, disponibilité runtime, intégration, support de plateforme ou déploiement de production.

## 3. Décideur et propriétaires affectés

- décideur projet : Hillel Tobiana ;
- propriétaire de la décision et de l'architecture produit : Product Architecture ;
- Security Architecture conserve l'autorisation, l'isolation Tenant, RBAC/ABAC, SoD et Decision Authority boundary ;
- Command conserve `CAP-CMD-401` et les Tasks de delivery ;
- Platform Settings conserve Tenant, Principal, Role et l'administration tenant-local ;
- Shared conserve Global Search, Reporting Engine et Export Engine ;
- Govern conserve Decision Authority, Action Request, Decision et Response Run.

## 4. Authority input approuvé — D1 à D7, verbatim

### APPROVED_D1_DEPLOYMENT_MODES

CMDR supports Internal, Enterprise multi-tenant, and MSP/MSSP deployment modes.
MSSP capabilities are deployment-dependent.

### APPROVED_D2_CUSTOMER_SEMANTICS

Customer remains an external/deployment/customer/contract projection.
Customer is NOT a canonical CMDR object and is NOT an alias of Tenant.

### APPROVED_D3_MSSP_TENANT_MODEL

MSSP operates over independent Tenants.
No parent/child Tenant hierarchy.
No ManagedTenant, TenantGroup, Portfolio, CustomerTenant, or equivalent canonical object for the MVP.

### APPROVED_D4_CROSS_TENANT_AUTH_MVP

Security resolves an Authorized Tenant Set as a non-canonical authorization projection.
The MVP allows read-only aggregation across that explicitly authorized Tenant set and explicit Tenant context switching.
Cross-tenant mutation, administration, response, delegated administration, and automatic export widening are not allowed.
Existing object read permissions remain subject to server-side RBAC/ABAC and Tenant-set evaluation.
No generic cross-tenant.manage permission is created.

### APPROVED_D5_SEARCH_REPORTING_SCOPE

Search is single-selected-Tenant initially.
Report is single-Tenant initially.
Export is single-Tenant initially and must never widen visibility.
Multi-tenant Search, Reporting, and Export are deferred.
Shared retains Search / Reporting / Export ownership.

### APPROVED_D6_RESPONSE_AUTHORITY

MSSP response requires explicit Tenant selection, Security authorization re-evaluation inside that Tenant, then Govern Decision Authority evaluation inside that Tenant.
No centralized cross-tenant response authority is introduced.

### APPROVED_D7_CAP_CMD_401_DISPOSITION

CAP-CMD-401 keeps the same ID and owner.
After the architecture is canonically recorded, it may move from delivery_status: proposed to delivery_status: defined.
Its scope is narrowed to deployment-aware activation, external Customer/engagement projection, authorized read-only multi-tenant overview, Tenant selector/context switching, portfolio-like UI projection without a Portfolio object, single-Tenant Reporting requests, service-delivery Tasks, contractual SLA projection, and source/freshness/audience visibility.

Explicitly OUT:
canonical Customer lifecycle/admin,
Tenant hierarchy,
cross-tenant mutation/admin/response,
delegated administration,
initial multi-tenant Search/Report/Export,
CRM,
billing,
customer portal,
contract mutation.

OPEN-013 remains OPEN.
OPEN-019 remains OPEN.

OPEN-006:
APPROVED AS RESOLVED BY D1–D7 ABOVE.

No new capability ID.
No CAP-SET-014.
No CAP-CMD-402.
No new canonical object.
No new Permission ID.
No new Screen ID.

The approved decision set above must be used verbatim as the authority input for the architecture-recording run.

## 5. Décision de déploiement

CMDR supporte trois modes : Internal, Enterprise multi-tenant et MSP/MSSP. Les capacités MSSP sont deployment-dependent. Aucun mode ne transforme Customer en prérequis universel de CMDR.

## 6. Customer semantics

`Customer` reste une projection externe issue d'une source deployment/customer/contract. CMDR peut afficher ses références, son engagement, ses obligations, son audience, sa version et sa fraîcheur lorsque la source existe, mais ne possède pas son lifecycle.

`Tenant ≠ Customer`. Aucun Customer canonique ni alias Tenant n'est créé.

## 7. Modèle Tenant / MSSP

Les Tenants restent des frontières d'isolation indépendantes. Un MSSP ne crée ni parent Tenant, ni child Tenant, ni `ManagedTenant`, `TenantGroup`, `Portfolio`, `CustomerTenant` ou relation équivalente pour le MVP.

Une vue de portfolio peut exister comme projection UI d'un ensemble de Tenants autorisés ; elle n'est pas un objet canonique et ne transfère aucune propriété ni permission.

## 8. Authorized Tenant Set

Security résout un `Authorized Tenant Set` comme projection d'autorisation non canonique pour le Principal, le deployment et le contexte d'autorisation courants.

Invariants :
- l'ensemble demandé doit être un sous-ensemble de l'ensemble autorisé ;
- chaque objet conserve son Tenant ;
- chaque lecture reste soumise à ses permissions objet et au RBAC/ABAC côté serveur ;
- aucun fallback Tenant silencieux ;
- la projection n'est ni Group, ni Role, ni Assignment, ni Effective Access, ni Portfolio ;
- aucune permission `cross-tenant.manage` n'est créée.

## 9. Cross-tenant MVP

Le MVP autorise seulement :
- agrégation read-only sur l'Authorized Tenant Set explicite ;
- affichage de l'identité Tenant de chaque objet ;
- sélection et changement explicites de Tenant.

Le MVP interdit en contexte agrégé :
- mutation ;
- administration ;
- réponse ;
- delegated administration ;
- élargissement automatique d'export.

Une mutation nécessite de revenir à un Tenant unique et de réévaluer l'autorisation dans ce Tenant.

## 10. Search, Reporting et Export

Initialement :
- Search = un Tenant sélectionné ;
- Report = un Tenant sélectionné ;
- Export = un Tenant sélectionné et ne peut jamais élargir la visibilité.

Multi-tenant Search, Reporting et Export sont différés. Shared conserve l'ownership de Search, Reporting et Export.

## 11. Response Authority

Pour une réponse MSSP :
1. sélectionner explicitement un Tenant ;
2. réévaluer l'autorisation Security dans ce Tenant ;
3. évaluer Govern Decision Authority dans ce Tenant ;
4. seulement ensuite créer/traiter les objets de réponse applicables.

Aucune autorité de réponse cross-tenant centralisée n'est introduite.

## 12. CAP-CMD-401

`CAP-CMD-401` conserve son ID et son owner Command Product Lead. Après enregistrement canonique de cette architecture, son `delivery_status` peut devenir `defined` tout en restant `draft / planned` et deployment-dependent.

Périmètre autorisé :
- deployment-aware activation ;
- projection Customer/engagement externe ;
- overview multi-tenant read-only autorisé ;
- Tenant selector/context switching ;
- projection UI type portfolio sans objet Portfolio ;
- single-Tenant Reporting requests ;
- service-delivery Tasks ;
- contractual SLA projection ;
- source/freshness/audience visibility.

Périmètre exclu : canonical Customer lifecycle/admin, Tenant hierarchy, cross-tenant mutation/admin/response, delegated administration, initial multi-tenant Search/Report/Export, CRM, billing, customer portal et contract mutation.

## 13. Permissions et Identity

Aucun nouvel ID de permission n'est créé. Les permissions de lecture existantes restent évaluées avec le scope Tenant via RBAC/ABAC côté serveur.

Cette décision ne crée ni Group, ni AccessAssignment, ni RoleAssignment, ni PermissionAssignment, ni Effective Access. Role reste distinct de Permission et de Decision Authority.

## 14. Screens et expérience

Aucun nouveau Screen ID n'est créé. `CMD-CRP-001` peut être aligné sur Customers & Delivery en conservant son ID. Le portfolio-like overview est une View/projection, pas une nouvelle Page par simple convenance.

Le changement de Tenant est une navigation autorisée, jamais un grant. Les permissions sont réévaluées après changement de contexte.

## 15. Delivery et SLA

- service-delivery follow-up = Task Command ;
- operational SLA tracking = Command ;
- contractual SLA = projection d'une source deployment/customer/contract ;
- contract mutation = hors scope ;
- SLO/resilience = autre travail Phase 6, non résolu ici.

## 16. OPEN préservées

`OPEN-013` reste ouverte pour la politique par défaut des mutations Class 2.

`OPEN-019` reste ouverte pour dissemination/releasability/sharing/external or client-facing delivery. Cette ADR n'autorise aucune publication externe implicite.

## 17. Alternatives rejetées ou différées

Rejetées pour le MVP :
- Internal-only ;
- Customer canonique ;
- Customer = Tenant ;
- hiérarchie parent/child de Tenants ;
- Managed Tenant Portfolio canonique ;
- permission générique `cross-tenant.manage` ;
- réponse centralisée cross-tenant ;
- cross-tenant mutation/admin/delegation.

Différées :
- multi-tenant Search ;
- multi-tenant Reporting ;
- multi-tenant Export ;
- delegated administration ;
- customer portal ;
- CRM/billing.

## 18. Conséquences positives

- Tenant reste une frontière d'isolation stable ;
- aucun objet client ou portfolio n'est inventé ;
- Internal reste pleinement viable sans Customer ;
- MSSP peut disposer d'un overview read-only sans créer d'autorité globale ;
- Shared et Govern conservent leurs responsabilités ;
- `CAP-CMD-401` peut être défini sans nouveau Capability ID.

## 19. Conséquences et coûts

- les opérations sur plusieurs Tenants exigent des changements de contexte explicites ;
- les reports consolidés multi-tenant sont différés ;
- les sources customer/contract restent deployment-specific ;
- le mécanisme runtime qui calcule l'Authorized Tenant Set reste une dépendance d'implémentation Security et n'est pas défini physiquement par cette ADR.

## 20. Non-objectifs

Cette ADR ne définit aucune API, protocole, schéma physique, index de recherche, moteur d'autorisation, stockage, connecteur, runtime de reporting, mécanisme de billing, CRM, portail client, SLO, failover, DR, Localization, Advanced Integrations ou Compliance engine.

## 21. Migration documentaire

- `OPEN-006` passe de open à resolved ;
- le module Customers & Delivery devient la source Command active pour ce contexte ;
- l'ancien module Customer & Reports reste un chemin deprecated/migration ;
- `CAP-CMD-401` peut passer de proposed à defined sans changement d'ID ;
- les registres, Requirements et roadmap sont mis à jour sans réécrire les snapshots historiques.

## 22. Critères d'acceptation

**Given** un Principal autorisé sur plusieurs Tenants, **When** il ouvre l'overview MSSP, **Then** seules des projections read-only des Tenants explicitement autorisés apparaissent et chaque objet expose son Tenant.

**Given** une action de mutation ou de réponse depuis un contexte agrégé, **When** l'utilisateur tente de l'exécuter, **Then** elle est refusée jusqu'à sélection d'un Tenant unique et réévaluation des permissions/authority.

**Given** un déploiement Internal sans Customer, **When** Reporting est utilisé, **Then** Reporting Engine reste disponible sans Customer canonique ni module client obligatoire.

**Given** un report ou export, **When** il est produit dans le MVP, **Then** il est single-Tenant et n'élargit jamais la visibilité.

## 23. Réévaluation

Réévaluer uniquement si une exigence source approuvée nécessite un objet Customer/Portfolio, une délégation d'administration, une autorité de réponse cross-tenant ou un Report/Search/Export réellement multi-tenant. Ces extensions exigent un nouveau source audit et ne sont pas déduites de cette ADR.
