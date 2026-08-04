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

`delivery_status` décrit la maturité fonctionnelle. `delivery_mode` décrit la preuve de livraison actuelle. Une capability peut être `defined` et `planned`.

## Convention d’identifiants

- `CAP-CMD-xxx` — Command ;
- `CAP-INV-xxx` — Investigate ;
- `CAP-GOV-xxx` — Govern ;
- `CAP-STD-xxx` — CMDR Studio ;
- `CAP-SET-xxx` — Platform Settings ;
- `CAP-EPT-xxx` — Endpoint Agent ;
- `CAP-SHR-xxx` — Shared Capabilities.

Un identifiant est unique, immuable, indépendant du chemin, inscrit dans le Capability Register et non recyclé.

## 1. Définition
Définition normative courte et distinctive.

## 2. Problème utilisateur
Situation, utilisateurs et conséquence observable sans la capability.

## 3. Objectifs
Résultats attendus, pas composants d’interface.

## 4. Non-objectifs
Frontières explicites et propriétaires concernés.

## 5. Propriétaire
Produit, module, rôle responsable et justification.

## 6. Utilisateurs
Rôle principal, rôles secondaires et responsabilités.

## 7. Conditions d’entrée
Contexte, objets, données, permissions générales et états requis.

## 8. Entrées fonctionnelles

| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|

Une ligne par entrée significative ; aucune cellule `N/A` sans justification.

## 9. Objets lus

| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|

Une projection n’accorde ni ownership ni mutation du cycle de vie source.

## 10. Objets créés ou modifiés

| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|

Même sans mutation, conserver une ligne explicite décrivant la projection en lecture seule ou l’événement produit.

## 11. Fonctionnalités
Chaque fonction possède un comportement observable et une frontière claire.

## 12. Actions utilisateur

| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|

Utiliser les classes 0 à 4. Une action de classe 3 ou 4 initiée hors Govern devient une demande, jamais une exécution locale.

## 13. Automatisation et IA

| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|

Ce tableau est obligatoire. Chaque fonction essentielle possède une alternative manuelle ou déterministe. Une matrice secondaire Humain/Règle/Moteur/Workflow/Agent/Govern peut compléter ce tableau mais ne le remplace pas.

## 14. États fonctionnels
États propres au travail porté par la capability, distincts des machines d’état d’objet.

## 15. États d’interface
Implications fonctionnelles de Loading, Empty, Partial, Error, Offline, Permission denied et Stale. Le Design System possède le rendu.

## 16. Sorties

| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|

## 17. Transitions

| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|

La transition conserve l’ownership du produit destination et le return origin.

## 18. Dépendances
Capabilities, objets, Shared Capabilities, produits et décisions ouvertes.

## 19. Source de vérité
Données propriétaires, projections, sources externes et données dérivées.

## 20. Provenance et audit
Acteur, moteur, source, changement, trace et conservation conceptuelle.

## 21. Permissions fonctionnelles
Familles existantes et lacunes futures, sans matrice atomique finale.

## 22. Limites et erreurs
Données absentes ou stale, objet inaccessible, dépendance indisponible, conflit, tenant/environnement et autorisation refusée.

## 23. Métriques
Métriques conceptuelles sans cible numérique non approuvée.

## 24. Classification de livraison
Delivery mode, preuve actuelle, cible, dépendances et conditions de promotion.

## 25. Critères d’acceptation
Au moins trois scénarios Given/When/Then spécifiques : nominal, erreur/permission et voie sans IA lorsque pertinente.

## 26. Questions ouvertes
Owner, phase cible et Requirement IDs ; aucune fermeture par simple rédaction.

## 27. Consommateurs documentaires
Modules, parcours, écrans, objets, permissions et contrats futurs.

## Contrôles

Une capability échoue si elle n’a pas 27 sections, front matter valide, owner, utilisateur, entrée, sortie, objet, action classée, alternative non-IA, critères spécifiques, Requirement ID ou l’un des six tableaux obligatoires. Les tableaux copiés sans adaptation sont un défaut.