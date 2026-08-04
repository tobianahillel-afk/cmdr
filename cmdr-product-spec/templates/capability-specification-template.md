---
id: template-capability-specification
domain: templates
status: draft
owner: Product Architecture
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-006
  - REQ-PROD-012
  - REQ-PROD-019
---

# Capability Specification Template

## Front matter obligatoire

```yaml
---
id: CAP-CMD-XXX
title:
product: command
module:
owner:
status: draft
delivery_status: defined|partial|planned|proposed|out-of-scope
delivery_mode: native|integrated|temporary-integration|planned|out-of-scope
last_updated: 2026-08-04
requirement_ids:
  - REQ-...
open_decisions:
  - OPEN-...
source-of-truth: canonical
---
```

`delivery_status` décrit la maturité fonctionnelle du document. `delivery_mode` décrit la preuve de livraison actuelle. Une capability peut donc être `defined` et `planned`.

## Convention d’identifiants

- `CAP-CMD-xxx` — Command ;
- `CAP-INV-xxx` — Investigate ;
- `CAP-GOV-xxx` — Govern ;
- `CAP-STD-xxx` — CMDR Studio ;
- `CAP-SET-xxx` — Platform Settings ;
- `CAP-EPT-xxx` — Endpoint Agent ;
- `CAP-SHR-xxx` — Shared Capabilities.

Un identifiant est unique, immuable, indépendant du chemin, inscrit dans le Capability Register et non recyclé après dépréciation. Les plages numériques peuvent regrouper les modules sans donner un sens fonctionnel au numéro.

## 1. Définition

Définition normative courte, suffisamment précise pour distinguer cette capability de ses voisines.

## 2. Problème utilisateur

Décrire la situation, les utilisateurs et la conséquence observable sans la capability.

## 3. Objectifs

Lister les résultats attendus, pas les composants d’interface.

## 4. Non-objectifs

Définir explicitement ce que la capability ne fait pas et les propriétaires concernés.

## 5. Propriétaire

Préciser produit, module, rôle responsable et justification de l’ownership.

## 6. Utilisateurs

Distinguer rôle principal, rôles secondaires et responsabilités.

## 7. Conditions d’entrée

Préciser contexte, objets, données, familles de permissions et états requis.

## 8. Entrées fonctionnelles

| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|

## 9. Objets lus

| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|

## 10. Objets créés ou modifiés

| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|

Une projection n’accorde jamais l’ownership ou le droit de modifier le cycle de vie de l’objet source.

## 11. Fonctionnalités

Chaque fonction possède un comportement observable et une frontière claire.

## 12. Actions utilisateur

| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|

Utiliser les classes 0 à 4. Une action de classe 3 ou 4 initiée depuis un produit non propriétaire devient une demande vers Govern, jamais une exécution locale.

## 13. Automatisation et IA

| Fonction | Humain | Règle | Moteur déterministe | Workflow | Agent | Govern | Alternative sans IA |
|---|---|---|---|---|---|---|---|

Toutes les fonctions essentielles possèdent une voie manuelle ou déterministe. Une suggestion d’agent reste attribuée et n’acquiert aucune autorité.

## 14. États fonctionnels

Définir les états du travail porté par la capability. Ne pas recopier automatiquement une machine d’état d’objet.

## 15. États d’interface

Décrire uniquement les implications fonctionnelles de Loading, Empty, Partial, Error, Offline, Permission denied et Stale lorsque pertinent. Le Design System reste propriétaire du rendu.

## 16. Sorties

| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|

## 17. Transitions

| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|

## 18. Dépendances

Lister capabilities, objets, Shared Capabilities, produits et décisions ouvertes.

## 19. Source de vérité

Distinguer données propriétaires, projections, sources externes et données dérivées.

## 20. Provenance et audit

Décrire acteur, moteur, source, changement, trace et conservation conceptuelle.

## 21. Permissions fonctionnelles

Référencer les familles existantes et signaler les besoins futurs. Ne pas créer la matrice atomique finale.

## 22. Limites et erreurs

Couvrir données absentes ou stale, objet inaccessible, dépendance indisponible, conflit, changement de tenant ou environnement et autorisation refusée.

## 23. Métriques

Définir des familles de métriques conceptuelles, sans cible numérique non approuvée.

## 24. Classification de livraison

Documenter delivery mode, preuve actuelle, cible, dépendances et conditions de promotion.

## 25. Critères d’acceptation

Au moins trois scénarios Given/When/Then spécifiques : voie nominale, erreur/permission et fonctionnement sans IA lorsque la capability peut utiliser une IA.

## 26. Questions ouvertes

Chaque question nomme owner, phase cible et Requirement IDs. Ne pas fermer une décision par simple rédaction.

## 27. Consommateurs documentaires

Lister modules, parcours futurs, écrans futurs, objets, permissions et contrats futurs.

## Contrôles

Une capability échoue à la revue si elle n’a pas d’owner, utilisateur, entrée, sortie, objet, action classée, alternative non-IA, critères spécifiques ou Requirement ID. Un texte identique répété entre plusieurs capabilities est un défaut de qualité.
