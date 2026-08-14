---
id: permission-model
domain: 14-security-permissions-and-trust
status: draft
owner: Security Architecture Lead
updated: 2026-08-14
source-of-truth: canonical
requirements:
  - REQ-PROD-004
  - REQ-PROD-006
  - REQ-PROD-008
  - REQ-SEC-001
  - REQ-SEC-002
---
# Permission Model

## Modèle

CMDR combine :

- **RBAC** pour les capacités de base ;
- **ABAC** pour tenant, environnement, ownership, classification, criticité, legal hold et contexte ;
- **Decision Authority** pour l’approbation d’actions, indépendante du CRUD ;
- **Separation of Duties** pour empêcher l’auto-approbation ou les combinaisons interdites ;
- **Step-up Authentication** pour les actions à risque.

## Format

`perm.<product-or-domain>.<resource>.<action>`

Une permission absente du catalogue ne doit pas être inventée par un consommateur.

## Règles générales

- Le Tenant est appliqué avant toute résolution d’objet.
- Lire n’implique ni exporter, ni administrer, ni exécuter, ni approuver, ni répondre.
- L’administration de plateforme n’accorde pas automatiquement autorité de réponse.
- Les permissions UI ne remplacent jamais l’enforcement serveur.
- Les décisions d’accès, refus et tentatives cross-scope sont auditées.
- L’Action Request, la Decision et le Response Run réévaluent permissions et conditions à chaque étape.
- Les commandes Endpoint Agent exigent permission, politique, signature, cible et expiry.
- Une projection, une relation, une View ou un contexte de navigation ne crée aucune permission.

## Authorized Tenant Set

Conformément à `tenant-isolation.md` et ADR-0008, Security peut résoudre un `Authorized Tenant Set` comme résultat d'autorisation non canonique pour le Principal et le contexte courants.

Ce résultat :
- contient uniquement des références Tenant ;
- n'est ni Role, ni Group, ni Assignment, ni Permission, ni objet canonique ;
- ne remplace jamais le contrôle de permission objet ;
- ne crée pas d'autorité Govern ;
- ne crée pas de droit d'administration.

Condition nécessaire pour une demande multi-tenant read-only :

`requested_tenant_set ⊆ authorized_tenant_set`.

Condition supplémentaire pour chaque objet : permission `read` existante + RBAC/ABAC + scope Tenant côté serveur.

## Actions cross-tenant — MVP

| Primitive | Scope MVP | Règle |
|---|---|---|
| READ | Authorized Tenant Set explicite | agrégation read-only autorisée si chaque objet passe son permission check |
| SEARCH | single-selected-Tenant | multi-tenant Search différé |
| EXPORT | single-selected-Tenant | export explicite et jamais d'élargissement de visibilité |
| ADMIN | single-selected-Tenant | permission administrative tenant-local + RBAC/ABAC |
| RESPONSE | single-selected-Tenant | Security re-evaluation puis Govern Decision Authority |
| DELEGATED ADMINISTRATION | non inclus | aucune délégation cross-tenant implicite |

Aucune permission générique `cross-tenant.manage` n'est créée.

## Tenant context switching

Passer d'un Tenant autorisé à un autre est Class 0 de navigation/contexte lorsqu'aucun effet métier n'est produit.

Le changement :
- n'accorde aucun droit ;
- efface le contexte incompatible ;
- réévalue les permissions dans le Tenant destination ;
- ne transporte jamais une permission d'administration, export, execution, approval ou response depuis le contexte précédent.

## Administration

L'administration est tenant-local. Un Principal qui peut lire plusieurs Tenants n'est pas automatiquement administrateur de ces Tenants.

Pour toute mutation Settings, les permissions administratives existantes et les contraintes RBAC/ABAC/SoD/step-up continuent de s'appliquer après sélection du Tenant.

## Response Authority

La réponse MSSP suit obligatoirement :

`select Tenant → Security authorization re-evaluation → Govern Decision Authority evaluation → governed response workflow`.

Un rôle MSSP, un Authorized Tenant Set, un rôle Platform Administrator ou un droit de lecture multi-tenant ne constituent jamais une Decision Authority.

## Search / Reporting / Export boundary

Shared reste owner de Search, Reporting et Export.

Pour le MVP :
- Search est single-selected-Tenant ;
- Report est single-Tenant ;
- Export est single-Tenant et ne peut jamais élargir la visibilité ;
- multi-tenant Search/Report/Export sont différés.

Le droit de lecture d'un objet ou Tenant ne confère jamais automatiquement le droit `export` ou `publish`.

## UX

Une action non pertinente peut être cachée. Une action pertinente mais interdite est désactivée avec raison, permission requise et chemin de demande d’accès sans révéler de donnée protégée.

Dans une vue MSSP agrégée :
- le Tenant de chaque objet reste visible ;
- les mutations et réponses cross-tenant sont indisponibles ;
- une action tenant-local exige une sélection explicite de Tenant avant affichage de l'action applicable.

## Audit

Les décisions allow/deny et tentatives cross-scope conservent acteur, contexte Tenant, action, outcome, justification ou policy reference lorsque disponible, timestamp et correlation id, sans enregistrer inutilement de donnée protégée.

## Identity non-regression

Cette architecture ne crée ni Group, ni AccessAssignment, ni RoleAssignment, ni PermissionAssignment, ni Effective Access. Une relation Principal/Role ne vaut pas permission effective et Role reste distinct de Decision Authority.

## OPEN préservées

- `OPEN-013` reste ouverte pour la politique par défaut des mutations réversibles Class 2 ;
- `OPEN-019` reste ouverte pour dissemination/releasability/external sharing.

Aucune de ces OPEN n'est résolue par le simple droit de lecture multi-tenant.

## Sources liées

- Catalogue : `permission-catalog.md`
- Registre : `../00-governance/registers/permission-register.md`
- Isolation : `tenant-isolation.md`
- Autorité : `decision-authority.md`
- Séparation des tâches : `separation-of-duties.md`
- Décision de déploiement : `../00-governance/adr/ADR-0008-customers-mssp-delivery-deployment-and-cross-tenant-architecture.md`

## Critères d'acceptation

**Given** un Principal autorisé sur plusieurs Tenants, **When** une lecture agrégée est demandée, **Then** l'ensemble demandé est borné par l'Authorized Tenant Set et chaque objet est encore contrôlé côté serveur.

**Given** un contexte MSSP agrégé, **When** une action d'administration ou de réponse est demandée, **Then** elle n'est pas autorisée par le scope agrégé et exige un Tenant unique avec nouvelle évaluation.

**Given** une permission `read`, **When** un export, une exécution ou une approbation est demandée, **Then** le droit de lecture seul est insuffisant.
