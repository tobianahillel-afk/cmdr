---
id: CAP-CMD-305
title: Capability Readiness
product: command
module: readiness-and-operations
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids: [REQ-PROD-012, REQ-PROD-013, REQ-PROD-019, REQ-PROD-057]
open_decisions: [OPEN-010, OPEN-013]
source-of-truth: canonical
---
# CAP-CMD-305 — Capability Readiness
## 1. Définition
Évalue si une capability est available, partial, degraded, not-configured, not-tested, maintenance, planned ou unknown sans confondre readiness, document status et delivery mode.
## 2. Problème utilisateur
Une cible native peut être non configurée/non livrée dans un tenant. Principal : Readiness Coordinator ; secondaires : Incident Commander, Product Owner, Platform Administrator.
## 3. Objectifs
Séparer les statuts ; afficher scope/source/date ; lier blockers/evidence ; créer validation/remediation Task ou ouvrir owner.
## 4. Non-objectifs
Ne pas promouvoir delivery mode, remplacer Studio Assurance/Settings Health ou inventer une preuve.
## 5. Propriétaire
Command possède l’assessment tenant-scoped ; Capability Register et produits owners gardent delivery/config/health/evidence.
## 6. Utilisateurs
Readiness Coordinator principal ; Incident Commander, Product Owner, Platform Administrator secondaires.
## 7. Conditions d’entrée
Capability ID, tenant/env, sources delivery/config/health/test.
## 8. Entrées fonctionnelles
Capability Register (owner/status/target), delivery/config health (products/Settings), test/assurance result (exercise/Studio/validation). Absence produit planned/unknown/not-tested selon preuve.
## 9. Objets lus
Capability entry, health/assurance/exercise evidence, Task.
## 10. Objets créés ou modifiés
Readiness assessment status/rationale/evidence refs (Command, C2) ; remediation/validation Task (Command, C2).
## 11. Fonctionnalités
Résoudre ID/owner ; afficher delivery séparément ; assessment tenant ; lier evidence/blockers ; Task validation/remediation.
## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter | lecteur | Capability | 0 | registry entry | status/source | non |
| Update assessment | coordinateur | assessment | 2 | evidence/rationale | new assessment | OPEN-013 |
| Créer validation Task | coordinateur | Task | 2 | not-tested/unknown | Task | non |
| Ouvrir owner source | lecteur | product capability | 0 | permission | transition | non |
## 13. Automatisation et IA
Moteur agrège registry/health/evidence ; règles appliquent mapping ; IA résume seulement ; humain valide assessment. Sans IA : mêmes sources/mappings.
## 14. États fonctionnels
`available`, `partial`, `degraded`, `not-configured`, `not-tested`, `maintenance`, `planned`, `unknown`.
## 15. États d’interface
Chaque état expose scope/source/date ; Partial/Unknown jamais masqués ; Permission denied n’expose pas evidence protégée. Rendu DS.
## 16. Sorties
Capability readiness vers Readiness/Plans/Mission Control ; validation/remediation Task vers owner.
## 17. Transitions
Owner detail vers product/Settings/Studio avec Capability ID/evidence refs ; gap vers Improvement Actions.
## 18. Dépendances
Capability Register, CAP-CMD-301/303, Platform Health, Studio Assurance, Metrics Engine.
## 19. Source de vérité
Register possède delivery/owner ; sources possèdent health/evidence ; Command possède assessment tenant-scoped.
## 20. Provenance et audit
Capability ID, scope, evidence refs, assessor, mapping/version, rationale et date.
## 21. Permissions fonctionnelles
Command read/coordinate + source reads. `OPEN-010/013` reportées.
## 22. Limites et erreurs
Entry absente, evidence stale, tenant mismatch, source down ou permission refusée donnent unknown/partial, jamais available inventé.
## 23. Métriques
Capabilities avec assessment tenant ; durée not-tested ; assessments avec evidence.
## 24. Classification de livraison
`defined` / `planned`, cible native ; readiness ne modifie jamais delivery mode.
## 25. Critères d’acceptation
**Given** capability `delivery_mode: planned`, **When** assessment est consulté, **Then** delivery reste planned et readiness est séparée.

**Given** evidence de test expirée, **When** update est calculé, **Then** not-tested/partial est affiché avec source/date.

**Given** aucun modèle, **When** assessment est produit, **Then** registry, health et règles déterministes suffisent.
## 26. Questions ouvertes
Quelle source authoritative tenant availability ? Quels statuts normaliser en Phase Objets ? — requirements ci-dessus. `OPEN-010/013` restent ouvertes.
## 27. Consommateurs documentaires
Readiness, Mission Control, Plans, reports, parcours Phase 5, écrans Phase 6, objets Phase 7, permissions ultérieures.
