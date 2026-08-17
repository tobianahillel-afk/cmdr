---
id: CMD-CRP-001
type: screen
product: command
module: customers-and-delivery
workspace: overview
status: draft
owner: Command Product Lead
updated: 2026-08-14
permissions:
  - perm.command.read
  - perm.shared.report.create
source-of-truth: screen
---
# Customers & Delivery

## 1. Objectif

Présenter le contexte de delivery Internal/Enterprise/MSSP, permettre un overview MSSP read-only des Tenants explicitement autorisés, sélectionner un Tenant et lancer des workflows tenant-local sans créer Customer ou Portfolio canonique.

## 2. Résultats utilisateur

L’utilisateur distingue deployment, Tenant, Customer/engagement externe, source, fraîcheur et audience ; il peut changer explicitement de Tenant et comprend pourquoi toute mutation ou réponse exige un Tenant unique.

## 3. Points d’entrée

Navigation Command, lien profond tenant-safe, retour depuis Reporting/Work Queue/Govern, ou contexte MSSP autorisé. Une entrée agrégée ne transmet jamais de droit de mutation.

## 4. Points de sortie

Retour au contexte source, sélection d'un Tenant, ouverture d'un inspecteur, Reporting Engine, Work Queue ou Govern après réévaluation des permissions dans le Tenant sélectionné.

## 5. Contexte

Deployment model, Authorized Tenant Set comme projection Security lorsque pertinent, Tenant sélectionné, environnement compatible, période, source/fraîcheur et URL de retour sont visibles. Chaque objet agrégé expose son Tenant. Le set n'est jamais un objet ni un grant autonome.

## 6. Structure de page

En-tête Command, navigation latérale, Context Bar, zone principale et Inspector canonique. Le portfolio-like overview est une View/projection de ce workspace, pas une Page ni un nouvel objet.

## 7. Hiérarchie de l’information

1. deployment et mode Internal/Enterprise/MSSP ;
2. scope Tenant autorisé et Tenant sélectionné ;
3. état delivery/readiness/SLA ;
4. Customer/engagement externe si présent ;
5. actions tenant-local disponibles seulement après sélection.

## 8. Actions principales

- Consulter l'overview autorisé en lecture seule ;
- Sélectionner/changer de Tenant ;
- Créer un report single-Tenant via Shared ;
- Créer un follow-up Task tenant-local lorsque permis.

## 9. Actions secondaires

Copier un identifiant non sensible, ouvrir la source canonique, partager un lien profond tenant-safe. Export est single-Tenant et exige la permission Shared correspondante ; `perm.command.read` ne l'accorde pas.

## 10. Données et objets

- Tenant — Platform Settings projection ;
- Incident / Task — Command ;
- Result — Govern projection ;
- Report — Shared projection ;
- Authorized Tenant Set — projection Security non canonique ;
- Customer / engagement / contractual SLA — projection d'une source deployment/customer/contract externe.

Aucun Customer, Portfolio, ManagedTenant, TenantGroup ou CustomerTenant canonique n'est créé.

## 11. Filtres et vues enregistrées

Le portfolio-like overview est une View du workspace. Les filtres sont URL-addressables sans secret. Les Saved Views utilisent la capability canonique. Search reste single-selected-Tenant initialement.

## 12. Inspector

L’inspecteur suit exclusivement `03-design-system/components/inspector.md`; il expose le Tenant et la source sans redéfinir Customer, Report ou objets propriétaires.

## 13. UX et interactions

Dans le contexte agrégé, les actions de mutation/admin/response sont indisponibles. Une action tenant-local déclenche d'abord une sélection explicite du Tenant, un nettoyage du contexte incompatible puis une réévaluation serveur. Aucun fallback Tenant silencieux.

## 14. Clavier et accessibilité

Parcours clavier complet, focus visible, libellés textuels, Tenant identifié autrement que par la couleur, alternatives aux graphes et réduction des animations respectée.

## 15. Permissions

Permissions référencées : `perm.command.read`, `perm.shared.report.create`. Les permissions complémentaires de Task/report publish/export restent celles de leurs owners et ne sont pas créées par l'écran. Security évalue RBAC/ABAC/Authorized Tenant Set côté serveur. Aucun `cross-tenant.manage`.

## 16. Audit

Les context switches, actions tenant-local, refus cross-scope, report requests et mutations Task conservent acteur, Tenant, action, résultat, justification/source et correlation id selon leur owner.

## 17. État Loading

Afficher la structure et le scope attendu sans inventer Customer, Tenant ou données d'un autre contexte.

## 18. État Empty

Distinguer : aucun Tenant autorisé, aucun objet dans le Tenant, contexte Customer absent ou module non applicable au deployment. Ne pas masquer ces cas derrière un Empty générique.

## 19. État Partial

Identifier les Tenants/sources manquants ou stale, la projection Customer/engagement absente, les impacts et le Tenant sélectionné. Ne jamais combler une lacune par des données d'un autre Tenant.

## 20. État Error

Conserver uniquement les données valides/encore autorisées, afficher correlation id et action de reprise tenant-safe. Aucune tentative de fallback Tenant.

## 21. État Offline

Passer en lecture limitée uniquement si le cache respecte encore l'isolation et la fraîcheur affichée ; bloquer mutations, publication/export et réponse non garanties.

## 22. État Permission denied

Expliquer que le Tenant, l'objet ou l'action n'est pas autorisé sans révéler Customer, engagement ou données d'un autre Tenant. Fournir un chemin de demande d'accès sûr lorsque disponible.

## 23. Comportement responsive

Préserver en priorité deployment, Tenant, source/fraîcheur et état read-only. Les colonnes secondaires deviennent onglets/panneaux sans masquer le Tenant actif.

## 24. Télémétrie produit

Mesurer ouverture, mode deployment, context switches, erreurs/refus, usage read-only, report requests et actions critiques sans enregistrer de secrets, contrat brut ou contenu d'Evidence.

## 25. Dépendances

ADR-0008, `CAP-CMD-401`, Tenant Isolation/Permission Model, Context Preservation, Shared Reporting/Global Search/Export, CAP-CMD-105, Work Queue et Govern Decision Authority.

## 26. Critères d’acceptation

**Given** un MSSP autorisé sur plusieurs Tenants, **When** l'overview est affiché, **Then** il reste read-only et chaque objet expose son Tenant.

**Given** une mutation demandée depuis l'overview, **When** aucun Tenant unique n'est sélectionné, **Then** l'action est bloquée.

**Given** un Tenant sélectionné, **When** un report est créé, **Then** il est single-Tenant et reste Shared-owned.

**Given** un deployment Internal sans Customer, **When** l'écran est utilisé pour Reporting, **Then** aucun Customer canonique n'est requis.

## 27. Questions ouvertes

`OPEN-013` reste ouverte pour les mutations Class 2. `OPEN-019` reste ouverte pour client/external sharing. `OPEN-006` est résolue par ADR-0008.

## Transitions interproduits

- Vers Shared Reporting Engine : Tenant sélectionné uniquement.
- Vers Shared Global Search : Tenant sélectionné uniquement.
- Vers Govern : Tenant sélectionné, Security re-evaluation puis Decision Authority.
