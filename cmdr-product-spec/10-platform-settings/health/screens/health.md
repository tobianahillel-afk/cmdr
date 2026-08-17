---
id: SET-HLT-001
type: screen
product: platform-settings
module: health
workspace: health
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-16
permissions:
  - perm.settings.health.read
source-of-truth: screen
---
# Platform Health

## 1. Objectif

Voir la santé, la fraîcheur, l'impact des dégradations et les projections SLO sourcées sans confondre lecture et exécution.

## 2. Résultats utilisateur

L’utilisateur comprend l’état sourcé, son Tenant, sa fraîcheur et ses limites avant tout handoff.

## 3. Points d’entrée

Navigation produit, lien profond, notification ou transition interproduits autorisée. Une recherche reste dans un Tenant sélectionné.

## 4. Points de sortie

Retour au contexte source, ouverture d’un inspecteur, navigation vers la source canonique ou handoff tenant-local explicitement confirmé.

## 5. Contexte

Tenant, environnement lorsqu'il est sourcé, période/window, typed subject reference, authoritative source, target version/effective period, freshness, objet actif et URL de retour sont visibles et préservés.

## 6. Structure de page

En-tête produit, navigation latérale, barre de contexte, zone principale et inspecteur canonique lorsque nécessaire.

## 7. Hiérarchie de l’information

État, source, freshness et limites précèdent les détails; unknown/partial/stale/conflicting/unsupported restent explicites.

## 8. Actions principales

- Inspecter.
- Ouvrir la source canonique.
- Effectuer un handoff humain tenant-local vers un écran existant lorsque pertinent.

`Acknowledge maintenance` est affiché uniquement comme action non exécutable/désactivée avec raison tant qu'un owner, une action class et une permission explicite ne sont pas sourcés.

## 9. Actions secondaires

Copier un identifiant, partager un lien profond. Report/Export passent par Shared, exigent leurs permissions propres et restent single-Tenant; `perm.settings.health.read` ne les autorise pas.

## 10. Données et objets

- [integration](../../../05-domain-model/objects/integration.md)
- [endpoint-agent-fleet](../../../05-domain-model/objects/endpoint-agent-fleet.md)
- [data-source](../../../05-domain-model/objects/data-source.md)

SLO/Health/Threshold/Availability/Reliability/Resilience/RTO/RPO restent des concepts/projections non canoniques selon ADR-0009.

## 11. Filtres et vues enregistrées

Les filtres sont URL-addressables. Une vue MSSP peut agréger en lecture les Tenants de l'Authorized Tenant Set tout en conservant le Tenant de chaque projection. Search reste single-selected-Tenant.

## 12. Inspector

L’inspecteur suit exclusivement `03-design-system/components/inspector.md`; l’écran ne redéfinit ni sa structure ni ses états.

## 13. UX et interactions

Les panneaux n’effacent pas la position de la liste. Aucun control runtime failover/recovery/monitoring n'est introduit. Toute action tenant-locale exige un Tenant sélectionné.

## 14. Clavier et accessibilité

Parcours clavier complet, focus visible, libellés textuels, alternatives aux graphes et respect de la réduction des animations.

## 15. Permissions

Permission référencée: `perm.settings.health.read`. Elle couvre la lecture seulement et n'implique ni export, ni configuration, ni monitoring runtime, ni failover/recovery. Aucune permission Health/SLO manage n'est créée.

## 16. Audit

Les lectures et handoffs auditables conservent acteur, Tenant, source, action, résultat, timestamp et correlation id selon les mécanismes canoniques. Cette surface ne crée aucune mutation Health.

## 17. État Loading

Afficher le squelette de structure sans inventer de données; annoncer le chargement aux technologies d’assistance.

## 18. État Empty

Expliquer pourquoi aucune donnée n’est visible sans conclure que la plateforme est healthy.

## 19. État Partial

Identifier sources manquantes, target/version manquants, fenêtre/calculation provenance indisponibles, freshness et conséquences. Ne pas synthétiser compliance ou breach.

## 20. État Error

Conserver les données valides, afficher l’erreur, l’identifiant de corrélation et une action de reprise sûre.

## 21. État Offline

Passer en lecture limitée lorsque possible, interdire toute action non garantie et montrer la dernière synchronisation. Offline/Retry n'implique aucun failover automatique.

## 22. État Permission denied

Expliquer la capacité refusée sans révéler de données protégées et fournir le chemin de demande d’accès.

## 23. Comportement responsive

Préserver l’ordre de décision; source, Tenant, freshness et état ne doivent pas être masqués par la réduction de largeur.

## 24. Télémétrie produit

Mesurer ouverture, durée, erreurs et handoffs sans enregistrer de secrets ni inventer des mesures Health/SLO.

## 25. Dépendances

ADR-0009, services d’objet, autorisation Security, Shared Metrics, Shared Notifications, Business Service Catalog et contrats Health/Metrics/Offline-Retry. Les probes et calculs restent source/runtime-owned.

## 26. Critères d’acceptation

Les six états obligatoires sont testés; source/Tenant/freshness sont préservés; aucune valeur compliance/breach n'est inventée; `health.read` ne permet aucune mutation/export/exécution; MSSP reste read-only et cross-tenant effects sont impossibles.

## 27. Questions ouvertes

`OPEN-008`, `OPEN-013`, `OPEN-015` et `OPEN-019` restent ouvertes dans leurs périmètres respectifs. Aucun nouvel ID de permission, objet, écran ou capability n'est introduit.

## Transitions interproduits

- Résumé consommé dans Command.
- SLO breach sourcé peut conduire à un handoff humain tenant-local; il ne crée rien automatiquement.
