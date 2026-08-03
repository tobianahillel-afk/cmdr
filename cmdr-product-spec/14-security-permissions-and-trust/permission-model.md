---
id: permission-model
domain: 14-security-permissions-and-trust
status: draft
owner: Security Architecture Lead
updated: 2026-08-03
source-of-truth: canonical
---
# Permission Model

## Modèle

CMDR combine:

- **RBAC** pour les capacités de base;
- **ABAC** pour tenant, environnement, ownership, classification, criticité, legal hold et contexte;
- **Decision Authority** pour l’approbation d’actions, indépendante du CRUD;
- **Separation of Duties** pour empêcher l’auto-approbation ou les combinaisons interdites;
- **Step-up Authentication** pour les actions à risque.

## Format

`perm.<product-or-domain>.<resource>.<action>`

## Règles

- Le tenant est appliqué avant toute résolution d’objet.
- Lire n’implique ni exporter, ni exécuter, ni approuver.
- L’administration de plateforme n’accorde pas automatiquement autorité de réponse.
- Les permissions UI ne remplacent jamais l’enforcement serveur.
- Les décisions d’accès, refus et tentatives cross-scope sont auditées.
- L’Action Request, la Decision et le Response Run réévaluent permissions et conditions à chaque étape.
- Les commandes Endpoint Agent exigent permission, politique, signature, cible et expiry.

## UX

Une action non pertinente peut être cachée. Une action pertinente mais interdite est désactivée avec raison, permission requise et chemin de demande d’accès sans révéler de donnée protégée.

## Sources liées

- Catalogue: `permission-catalog.md`
- Registre: `../00-governance/registers/permission-register.md`
- Autorité: `decision-authority.md`
- Séparation des tâches: `separation-of-duties.md`
