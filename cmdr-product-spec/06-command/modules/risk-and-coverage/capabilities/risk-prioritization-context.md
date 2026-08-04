---
id: CAP-CMD-205
title: Risk Prioritization Context
product: command
module: risk-and-coverage
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids: [REQ-PROD-003, REQ-PROD-010, REQ-PROD-013, REQ-PROD-021]
open_decisions: [OPEN-013]
source-of-truth: canonical
---
# CAP-CMD-205 — Risk Prioritization Context
## 1. Définition
Combine explicitement impact, urgence, criticité Service, exposure, coverage, confiance, dépendances et SLA pour éclairer une priorité sans moteur définitif imposé.
## 2. Problème utilisateur
Un score unique masque facteurs, inconnues et responsabilités. Principal : Incident Commander ; secondaires : SOC L2, Business Owner, Risk stakeholder.
## 3. Objectifs
Afficher facteurs/sources ; indiquer missing/conflicts ; produire proposition distincte de priority ; permettre accept/reject humain audité.
## 4. Non-objectifs
Choisir algorithme final, score universel, modifier sources ou transformer proposal en Decision.
## 5. Propriétaire
Command possède le contexte/proposal ; sources gardent facteurs ; CAP-CMD-002 possède la mutation de priorité.
## 6. Utilisateurs
Incident Commander principal ; SOC L2, Business Owner, Risk stakeholder secondaires.
## 7. Conditions d’entrée
Work item, au moins un facteur, définitions/fraîcheur par facteur.
## 8. Entrées fonctionnelles
Impact/Urgency/SLA Command ; Service/Exposure/Coverage projections ; recommendation rule/engine/workflow/agent avec version/run. Une absence produit `partial`, jamais un facteur inventé.
## 9. Objets lus
Incident/Task (Command), Service/Exposure/Coverage (sources), Confidence/Severity (Signal/Alert/source).
## 10. Objets créés ou modifiés
Priority recommendation (proposal Command) ; Incident/Task seulement via acceptation dans CAP-CMD-002, classe 2.
## 11. Fonctionnalités
Afficher facteur/unité/source/freshness ; comparer action/inaction ; proposal déterministe ; proposal IA attribuée ; renvoi vers Priority Management.
## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspecter | lecteur | risk context | 0 | sources | facteurs | non |
| Recalcul déterministe | coordinateur | context | 0 | inputs/version | proposal | non |
| Accepter/rejeter | coordinateur | proposal | 2 | rationale visible | disposition | OPEN-013 |
| Ouvrir source | lecteur | source product | 0 | permission | transition | non |
## 13. Automatisation et IA
Humain décide ; règle/moteur calculent explicitement ; workflow peut préparer ; agent propose seulement ; Govern reçoit les high-effect actions. Sans IA : calcul déterministe et décision humaine.
## 14. États fonctionnels
`complete`, `partial`, `conflicting`, `proposal-available`, `proposal-rejected`, `unknown`.
## 15. États d’interface
Chaque facteur affiche source/fraîcheur ; Partial/Conflict ne sont jamais masqués ; Offline interdit acceptation stale. Rendu DS.
## 16. Sorties
Risk context vers Priority/Mission Control/Govern ; Priority proposal sourcée vers coordinateur.
## 17. Transitions
Accept proposal vers CAP-CMD-002 sans changer priority avant action ; high-effect vers Govern sans Decision locale.
## 18. Dépendances
CAP-CMD-002/104/105/201/202/203/204, Metrics Engine.
## 19. Source de vérité
Chaque facteur reste source-owned ; proposal est Command et ne devient pas effective.
## 20. Provenance et audit
Producer type, version/run, facteurs, missing data, rationale et disposition.
## 21. Permissions fonctionnelles
Command read ; coordinate pour disposition ; source-specific reads. `OPEN-013` reportée.
## 22. Limites et erreurs
Facteurs incompatibles/stale/unknown, source refusée ou moteur indisponible donnent Partial/Conflict, jamais score universel.
## 23. Métriques
Proposals avec facteurs/source/freshness ; accept/reject par producer ; contexts conflicting/unknown.
## 24. Classification de livraison
`defined` / `planned`, cible native ; moteur final non choisi.
## 25. Critères d’acceptation
**Given** facteurs partiels, **When** proposal calculée, **Then** missing data et facteurs sont visibles et priority reste inchangée.

**Given** proposal IA, **When** consultée, **Then** producer/run/rationale sont visibles et acceptation humaine requise.

**Given** aucun modèle, **When** recalcul demandé, **Then** moteur déterministe ou analyse manuelle fournit le contexte.
## 26. Questions ouvertes
Quel moteur ultérieur ? Quels facteurs comparables entre tenants ? — requirements ci-dessus. `OPEN-013` reste ouverte.
## 27. Consommateurs documentaires
Priority Management, Mission Control, Govern package, parcours Phase 5, écrans Phase 6, objets Phase 7, permissions ultérieures.
