---
id: CAP-CMD-204
title: Business Impact Context
product: command
module: risk-and-coverage
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids: [REQ-PROD-003, REQ-PROD-005, REQ-PROD-013, REQ-PROD-021]
open_decisions: [OPEN-013]
source-of-truth: canonical
---
# CAP-CMD-204 — Business Impact Context
## 1. Définition
Capture/projette l’impact métier d’un Incident avec type, portée, source, owner, certitude et date, puis le transmet à Govern.
## 2. Problème utilisateur
Un score/texte sans auteur influence priorité et autorité. Principal : Incident Commander ; secondaires : Business Owner, Service Owner, SOC L2.
## 3. Objectifs
Distinguer assumed/confirmed/disputed/unknown ; relier Services/owner/type ; conserver source/date/justification ; transmettre sans score opaque.
## 4. Non-objectifs
Calcul financier définitif, modification Catalog, remplacement Business Owner ou confirmation IA.
## 5. Propriétaire
Command possède l’impact Incident ; Catalog/Investigate/Govern gardent leurs objets.
## 6. Utilisateurs
Incident Commander principal ; Business/Service Owner et SOC L2 secondaires.
## 7. Conditions d’entrée
Incident accessible ; service confirmé ou absence documentée ; acteur/source.
## 8. Entrées fonctionnelles
Impact statement (human/source, type/extent/certainty), Service Context, technical context Incident/Case/Finding. Chaque source expose fraîcheur ; unknown crée follow-up.
## 9. Objets lus
Incident (Command), Service (Shared), Finding/Result (Investigate/Govern) en projection.
## 10. Objets créés ou modifiés
Incident : statement/certainty/source/owner (Command, C2). Action Request context : transmis à Govern sans Decision locale.
## 11. Fonctionnalités
Saisir impact ; confirmer/contester hypothèse ; lier Services/owner ; afficher source/fraîcheur ; préparer contexte Govern.
## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Saisir assumed | coordinateur | Incident | 2 | source/raison | assumed | OPEN-013 |
| Confirmer | Business Owner autorisé | Incident | 2 | données suffisantes | confirmed | OPEN-013 |
| Contester | rôle autorisé | Incident | 2 | motif | disputed | OPEN-013 |
| Transmettre | requester | Action Request context | 2 | action proposée | package | oui |
## 13. Automatisation et IA
IA peut proposer un résumé/impact assumed, jamais confirmed. Moteur peut agréger Services/indicateurs. Alternative : saisie et confirmation humaines.
## 14. États fonctionnels
`unknown`, `assumed`, `confirmed`, `disputed`, `outdated`, `superseded`.
## 15. États d’interface
Certitude/source toujours visibles ; stale/outdated distinct ; Permission denied masque seulement la projection protégée. Rendu DS.
## 16. Sorties
Business impact context vers Priority/Mission Control/Govern ; follow-up Task vers Business/Service Owner.
## 17. Transitions
Validation vers Business Owner/Task ; Action Request vers Govern avec Incident/Services/impact/urgence/source ; Decision reste Govern.
## 18. Dépendances
CAP-CMD-201/205/106, Catalog, Linking Service.
## 19. Source de vérité
Impact Incident : Command ; Service/Finding/Result : propriétaires sources.
## 20. Provenance et audit
Auteur/producteur, certitude, source, before/after, time et correlation ID.
## 21. Permissions fonctionnelles
Incident read/manage, future Business Owner confirmation, Govern Action Request create.
## 22. Limites et erreurs
Service inconnu, impact stale, owner absent, source conflictuelle ou refus ne doivent jamais produire confirmed.
## 23. Métriques
Incidents avec certitude ; délai assumed→confirmed/disputed ; Action Requests avec impact sourcé.
## 24. Classification de livraison
`defined` / `planned`, cible native ; preuve documentaire.
## 25. Critères d’acceptation
**Given** impact proposé par IA, **When** affiché, **Then** il reste assumed et source/producteur sont visibles.

**Given** Business Owner confirme, **When** action autorisée, **Then** before/after/auteur/date sont audités.

**Given** aucun modèle, **When** impact est saisi, **Then** voie humaine complète disponible.
## 26. Questions ouvertes
Permission de confirmation métier ? Types d’impact canoniques ? — requirements ci-dessus. `OPEN-013` reste ouverte.
## 27. Consommateurs documentaires
Incident Detail, Mission Control, Priority, Govern transition, parcours Phase 5, écrans Phase 6, objets Phase 7, permissions ultérieures.
