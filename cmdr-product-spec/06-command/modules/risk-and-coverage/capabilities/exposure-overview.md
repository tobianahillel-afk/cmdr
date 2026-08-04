---
id: CAP-CMD-202
title: Exposure Overview
product: command
module: risk-and-coverage
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids: [REQ-PROD-006, REQ-PROD-013, REQ-PROD-032, REQ-PROD-037]
open_decisions: [OPEN-013]
source-of-truth: canonical
---
# CAP-CMD-202 — Exposure Overview
## 1. Définition
Présente les expositions utiles à la coordination avec source, état, fraîcheur et relations Service/Incident, sans scanner ni créer une source d’exposition.
## 2. Problème utilisateur
Les coordinateurs doivent comprendre le risque sans se substituer aux outils de vulnérabilité/posture. Principal : Incident Commander ; secondaires : SOC L2, Business Owner, owner externe. Sans capacité, exposition et contexte opérationnel restent séparés.
## 3. Objectifs
Agréger projections autorisées ; relier exposure/service/incident/coverage ; rendre source/fraîcheur visibles ; créer Task/lien, jamais modifier la source.
## 4. Non-objectifs
Scanner, définir un objet complet, confirmer vulnérabilité ou calculer un score universel.
## 5. Propriétaire
Command possède la projection opérationnelle ; la source externe/shared possède Exposure.
## 6. Utilisateurs
Incident Commander principal ; SOC L2, Business Owner, exposure owner externe secondaires.
## 7. Conditions d’entrée
Source configurée ou partial ; tenant/service ; permission source.
## 8. Entrées fonctionnelles
| Entrée | Source | Type | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Exposure projection | integration/source owner | summary | non | source timestamp | unavailable/unknown |
| Service links | Catalog | service/criticality | non | déclarée | exposure sans service confirmé |
| Incident links | Command | relations | non | version courante | aucune relation |
## 9. Objets lus
| Objet | Owner | Projection | Droit local |
|---|---|---|---|
| Exposure | external/shared | type/status/source/freshness | projection |
| Service | Shared | relation/criticité | projection |
| Incident | Command | relation/état | lecture/modification lien |
## 10. Objets créés ou modifiés
| Objet | Opération | Owner | Règle |
|---|---|---|---|
| Incident / Task | lien exposure/follow-up | Command | classe 2 |
| Exposure source | aucune | external | read-only |
## 11. Fonctionnalités
Filtrer par service/source/état/fraîcheur ; inspecter définition ; lier Incident ; créer Task ; comparer exposure/coverage sans fusion.
## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter | lecteur | projection | 0 | permission | détail | non |
| Lier Incident | coordinateur | relation | 2 | compatible | lien sourcé | OPEN-013 |
| Créer Task | coordinateur | Task | 2 | action attendue | Task | non |
| Ouvrir source | lecteur | source externe | 0 | deep link | transition/return | non |
## 13. Automatisation et IA
Matching peut être rule/deterministic ; agent suggère seulement ; humain confirme. Sans IA : filtres, sources et liens manuels.
## 14. États fonctionnels
`active`, `mitigated-source`, `accepted-source`, `unknown`, `stale`, `source-unavailable`, `unlinked`.
## 15. États d’interface
Partial/Stale/Unavailable sont explicites ; Offline bloque lien ; Permission denied masque la projection. Rendu DS.
## 16. Sorties
Exposure context sourcé vers Risk/Mission Control/Incident ; Task/lien vers Work Queue sans mutation source.
## 17. Transitions
Ouvrir source externe avec return origin ; besoin d’investigation vers Investigate avec refs, Case owner Investigate.
## 18. Dépendances
Catalog, Linking Service, Data Quality, CAP-CMD-201/203.
## 19. Source de vérité
Source externe possède Exposure ; Command possède seulement liens Incident/Task.
## 20. Provenance et audit
Source/version/time, lien, acteur, justification et correlation ID.
## 21. Permissions fonctionnelles
Lecture Command + source-specific ; manage Incident/Task pour lien. Atomisation reportée.
## 22. Limites et erreurs
Source down, projection stale, lien ambigu, service inconnu ou permission refusée ne deviennent jamais exposure confirmed.
## 23. Métriques
Projections avec source/fraîcheur ; liées à Service/Incident ; sources indisponibles. Pas de cible.
## 24. Classification de livraison
`defined` / `planned`, cible native ; aucun scanner prouvé.
## 25. Critères d’acceptation
**Given** une exposure externe autorisée, **When** elle est liée à un Incident, **Then** la relation est auditée et la source reste owner.

**Given** source indisponible, **When** la vue est ouverte, **Then** état/fraîcheur restent visibles et aucune donnée n’est inventée.

**Given** aucun modèle, **When** la capability est utilisée, **Then** filtres et actions manuelles suffisent.
## 26. Questions ouvertes
Quel objet/source canonique ? Quelles expositions sont pertinentes sans devenir une liste vulnérabilité ? — requirements ci-dessus. `OPEN-013` reste ouverte.
## 27. Consommateurs documentaires
Risk screens, Incident Detail, Investigate transition, parcours Phase 5, écrans Phase 6, objets Phase 7, permissions ultérieures.
