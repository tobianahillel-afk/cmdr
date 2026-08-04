---
id: CAP-CMD-303
title: Improvement Actions
product: command
module: readiness-and-operations
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids: [REQ-PROD-009, REQ-PROD-013, REQ-PROD-021, REQ-OBJ-012]
open_decisions: [OPEN-013]
source-of-truth: canonical
---
# CAP-CMD-303 — Improvement Actions
## 1. Définition
Transforme Result, exercice, incident review ou gap de couverture en Task Command avec cause, owner, priorité, résultat attendu et vérification.
## 2. Problème utilisateur
Les leçons restent des recommandations sans owner ni échéance. Principal : Readiness Coordinator ; secondaires : Incident Commander, Task owner, Service Owner.
## 3. Objectifs
Réutiliser Task ; lier cause/expected result ; gérer owner/priority/blockers/due ; exiger vérification ou exception de clôture.
## 4. Non-objectifs
Ne pas créer backlog général, modifier source, évaluer un agent ou finaliser schéma Task.
## 5. Propriétaire
Command possède la Task ; Result/exercise/coverage restent leurs sources propriétaires.
## 6. Utilisateurs
Readiness Coordinator principal ; Incident Commander, Task owner, Service Owner secondaires.
## 7. Conditions d’entrée
Source ou justification, action attendue, owner/team ou raison unassigned.
## 8. Entrées fonctionnelles
Source gap (Result/exercise/coverage/review), expected outcome, owner/due/priority. Une action vague est refusée.
## 9. Objets lus
Result (Govern projection), Exercise/readiness/coverage (sources), Task (Command).
## 10. Objets créés ou modifiés
Task create/update/close/reopen (Command, C2) ; source relation via Linking Service, sans mutation source.
## 11. Fonctionnalités
Créer Task depuis gap ; assign/prioritize/track ; validation evidence ref ; clôture avec vérification/exception ; reopen si invalidation.
## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Créer | coordinateur | Task | 2 | source/expected outcome | open | OPEN-013 |
| Update/assign | owner/coordinateur | Task | 2 | version | updated | OPEN-013 |
| Soumettre vérification | owner | Task | 2 | outcome | verification | non |
| Clore/réouvrir | verifier/coordinateur | Task | 2 | preuve/raison | done/reopened | OPEN-013 |
## 13. Automatisation et IA
Règles détectent gaps ; moteur déduplique ; workflow notifie ; agent propose Task. Humain valide et ferme. Sans IA : création/validation manuelles.
## 14. États fonctionnels
`open`, `assigned`, `in-progress`, `blocked`, `verification`, `closed`, `reopened`, `cancelled`.
## 15. États d’interface
Partial montre evidence/source manquante ; conflit garde version ; Offline bloque clôture ; Permission denied masque source. Rendu DS.
## 16. Sorties
Improvement Task vers Work Queue ; verification outcome vers Readiness assessment, sourcé.
## 17. Transitions
Validation technique vers product owner/Studio/Settings ; failed verification rouvre Task.
## 18. Dépendances
CAP-CMD-107/301/302/203, Linking Service.
## 19. Source de vérité
Task et coordination : Command ; cause/evidence : source propriétaire.
## 20. Provenance et audit
Source, cause, expected result, owner, changes, evidence, verifier et correlation ID.
## 21. Permissions fonctionnelles
Task manage, source read, future verification permission. `OPEN-013` reportée.
## 22. Limites et erreurs
Source absente, expected outcome vague, evidence inaccessible, conflit ou permission refusée empêchent clôture confirmée.
## 23. Métriques
Gap→owned Task ; clôture avec vérification ; reopen après échec.
## 24. Classification de livraison
`defined` / `planned`, cible native ; preuve documentaire.
## 25. Critères d’acceptation
**Given** un gap de coverage, **When** Task créée, **Then** cause, owner, due et expected result sont liés.

**Given** vérification échoue, **When** résultat enregistré, **Then** Task est reopened avec raison/evidence.

**Given** aucun modèle, **When** action créée/fermée, **Then** voie manuelle complète disponible.
## 26. Questions ouvertes
Qui confirme vérification ? Quelles exceptions permettent clôture sans preuve ? — requirements ci-dessus. `OPEN-013` reste ouverte.
## 27. Consommateurs documentaires
Readiness, Work Queue, Result/exercise reviews, parcours Phase 5, écrans Phase 6, Task Phase 7, permissions ultérieures.
