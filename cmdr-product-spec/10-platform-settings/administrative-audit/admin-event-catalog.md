---
id: settings-administrative-audit-admin-event-catalog
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-12
source-of-truth: canonical
---
# Admin Event Catalog

## Objectif

Définir admin event catalog dans Platform Settings.

## Périmètre

Document canonique du domaine. Il définit uniquement son sujet et renvoie vers les autres sources de vérité pour les concepts partagés. Il reçoit les outcomes/provenance administratifs sans devenir owner des objets ou de l'autorisation.

## Propriétaire fonctionnel

Platform Settings Product Lead.

## Objets concernés

- Concepts du document
- Références canoniques liées

## Fonctionnalités

- Identity changes.
- Tenant changes.
- Secrets/connections.
- Agent policies.
- Retention.

Pour Identity Administration, les événements/provenance couvrent les mutations Principal/Role réellement sourcées, les validations/refus et les Access Review outcomes/handoffs. Une `revoke` disposition CAP-SET-007 n'est pas enregistrée comme suppression d'assignment tant qu'aucune mécanique canonique de retrait n'existe.

## UX et interactions

- Navigation par liens stables.
- Contenu lisible en thème clair et sombre.
- Aucune duplication des définitions externes.
- `SET-AUD-001` reste l'unique écran administratif d'audit existant; aucun nouvel Screen ID n'est créé.

## Permissions

Les modifications suivent le modèle défini dans `../../14-security-permissions-and-trust/permission-model.md` lorsque le document décrit une capacité exécutable. Aucun nouvel ID n'est introduit.

## États

Le statut documentaire suit `00-governance/document-status-model.md`; les états métier restent dans leurs sources canoniques.

## Dépendances

- `../../00-governance/source-of-truth-policy.md`
- `../../14-security-permissions-and-trust/audit-and-immutability.md`
- `../capabilities/cap-set-005-principal-administrative-lifecycle-and-identity-state.md`
- `../capabilities/cap-set-006-role-administrative-lifecycle-constraints-and-principal-relation-boundary.md`
- `../capabilities/cap-set-007-access-review-evidence-and-revocation-disposition.md`

## Critères d’acceptation

- Le document a un propriétaire unique.
- Les liens locaux sont valides.
- Les événements administratifs conservent acteur, Tenant, cible, action/outcome, justification et corrélation lorsque disponibles.
- Aucun événement ne transforme une disposition en mutation non sourcée.

## Questions ouvertes

- Les décisions globales existantes restent dans le registre OPEN; ce catalogue n'en résout aucune.
