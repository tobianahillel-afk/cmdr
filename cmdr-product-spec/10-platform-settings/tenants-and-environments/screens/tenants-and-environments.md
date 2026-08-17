---
id: SET-TEN-001
type: screen
product: platform-settings
module: tenants-and-environments
workspace: tenants
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-03
permissions:
  - perm.settings.tenant.read
  - perm.settings.tenant.manage
source-of-truth: screen
---
# Tenants & Environments

## 1. Objectif

Administrer tenants et environnements sans ambiguïté de scope.

## 2. Résultats utilisateur

L’utilisateur comprend la situation, prend la décision attendue et conserve le contexte du produit.

## 3. Points d’entrée

Navigation produit, lien profond, recherche globale, notification ou transition interproduits autorisée.

## 4. Points de sortie

Retour au contexte source, ouverture d’un inspecteur, navigation vers un écran lié ou transition interproduits explicitement confirmée.

## 5. Contexte

Tenant, environnement, période, objet actif, filtres sûrs et URL de retour sont visibles et préservés.

## 6. Structure de page

En-tête produit, navigation latérale, barre de contexte, zone principale et inspecteur canonique lorsque nécessaire.

## 7. Hiérarchie de l’information

La décision principale précède les détails; les informations secondaires sont révélées progressivement.

## 8. Actions principales

- Créer
- Suspendre
- Offboard
- Configurer

## 9. Actions secondaires

Copier un identifiant, ouvrir la source canonique, partager un lien profond et exporter uniquement avec permission.

## 10. Données et objets

- [tenant](../../../05-domain-model/objects/tenant.md)
- [environment](../../../05-domain-model/objects/environment.md)

## 11. Filtres et vues enregistrées

Les filtres sont URL-addressables. Les vues enregistrées utilisent la capacité canonique; la Work Queue suit `06-command/modules/incidents-and-work-queue/saved-views.md`.

## 12. Inspector

L’inspecteur suit exclusivement `03-design-system/components/inspector.md`; l’écran ne redéfinit ni sa structure ni ses états.

## 13. UX et interactions

Les panneaux n’effacent pas la position de la liste. Les actions à impact affichent cible, portée, effet, préconditions et retour arrière.

## 14. Clavier et accessibilité

Parcours clavier complet, focus visible, libellés textuels, alternatives aux graphes et respect de la réduction des animations.

## 15. Permissions

Permissions référencées: `perm.settings.tenant.read`, `perm.settings.tenant.manage`. La source unique est `14-security-permissions-and-trust/permission-model.md`.

## 16. Audit

Toute mutation enregistre acteur, tenant, objet, action, résultat, justification et identifiant de corrélation.

## 17. État Loading

Afficher le squelette de structure sans inventer de données; annoncer le chargement aux technologies d’assistance.

## 18. État Empty

Expliquer pourquoi aucune donnée n’est visible et proposer une action sûre ou un ajustement de filtre.

## 19. État Partial

Identifier les sources manquantes, la fraîcheur et les conséquences sur la décision.

## 20. État Error

Conserver les données valides, afficher l’erreur, l’identifiant de corrélation et une action de reprise sûre.

## 21. État Offline

Passer en lecture limitée lorsque possible, interdire les mutations non garanties et montrer la dernière synchronisation.

## 22. État Permission denied

Expliquer la capacité refusée sans révéler de données protégées et fournir le chemin de demande d’accès.

## 23. Comportement responsive

Préserver l’ordre de décision; les colonnes secondaires deviennent onglets ou panneaux sans masquer l’état actif.

## 24. Télémétrie produit

Mesurer ouverture, durée, erreurs, transitions et actions critiques sans enregistrer de secrets ni de contenu d’évidence.

## 25. Dépendances

Services d’objet, autorisation, audit, recherche, notifications et contrats d’implémentation du module.

## 26. Critères d’acceptation

Tous les six états obligatoires sont testés; les liens profonds survivent au rafraîchissement; les permissions sont vérifiées côté serveur; les transitions conservent le contexte.

## 27. Questions ouvertes

À compléter — contenu source non fourni dans le brief canonique.

## Transitions interproduits

- Contexte global vers tous les produits.
