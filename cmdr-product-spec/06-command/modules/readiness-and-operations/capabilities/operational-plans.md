---
id: CAP-CMD-304
title: Operational Plans
product: command
module: readiness-and-operations
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids: [REQ-PROD-005, REQ-PROD-013, REQ-PROD-021, REQ-PROD-057]
open_decisions: [OPEN-010, OPEN-013]
source-of-truth: canonical
---
# CAP-CMD-304 — Operational Plans
## 1. Définition
Maintient des plans opérationnels versionnés décrivant scénarios, rôles, responsabilités, conditions d’activation, procédures liées, owner, statut et dernière revue, sans moteur d’exécution.
## 2. Problème utilisateur
Des procédures dispersées n’exposent ni owner, ni review date, ni activation. Principal : Readiness Coordinator ; secondaires : Incident Commander, Business Owner, Response Operator.
## 3. Objectifs
Rendre plans/owners/conditions visibles ; lier procédures/capabilities ; versionner revue/supersession ; permettre activation conceptuelle sans exécution.
## 4. Non-objectifs
Ne pas exécuter playbook/workflow, posséder les outils, remplacer Studio Workflow/Govern Playbook ou définir tous les champs.
## 5. Propriétaire
Command possède le plan opérationnel de coordination ; Studio/Govern et sources conservent leurs procédures exécutables.
## 6. Utilisateurs
Readiness Coordinator principal ; Incident Commander, Business Owner, Response Operator secondaires.
## 7. Conditions d’entrée
Scénario, owner, audience, scope, procédures de référence et permission.
## 8. Entrées fonctionnelles
Plan definition (owner/coordinator), procedure links (external/source catalog), capability requirements (Readiness inventory). Les références stale restent visibles.
## 9. Objets lus
Plan/Playbook refs (Command/Govern/Studio/source), Capability Readiness, Service/Incident projections.
## 10. Objets créés ou modifiés
Operational Plan record create/update/review/supersede (Command readiness, C2) ; Task review/gap (Command, C2).
## 11. Fonctionnalités
Créer/versionner ; définir activation/roles ; lier sans copier ; montrer review/gaps ; active/expired/superseded.
## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Créer/modifier | coordinateur | plan | 2 | owner/scenario | draft/version | OPEN-013 |
| Soumettre/revoir | owner/reviewer | plan | 2 | complet | review/active | OPEN-013 |
| Superséder | owner | plan | 2 | replacement lié | superseded | OPEN-013 |
| Créer review Task | coordinateur | Task | 2 | due/stale | Task | non |
## 13. Automatisation et IA
Règles vérifient review/completude ; moteur résout refs ; workflow notifie ; IA peut préparer un résumé, jamais activer. Sans IA : édition/revue manuelles.
## 14. États fonctionnels
`draft`, `review`, `active`, `expired`, `superseded`, `withdrawn`, `stale-references`.
## 15. États d’interface
Stale refs et owners manquants sont visibles ; Offline bloque activation/supersession ; Permission denied masque les procédures protégées. Rendu DS.
## 16. Sorties
Operational Plan vers IC/exercises ; review/gap Task vers owner avec raison/due.
## 17. Transitions
Exécution souhaitée vers Studio/Govern owner avec plan ref/Incident/conditions ; exercice vers CAP-CMD-302 avec scenario/roles/objectives.
## 18. Dépendances
CAP-CMD-301/302/305, Linking, Versioning, Notifications.
## 19. Source de vérité
Plan record : Command ; procédures/playbooks/workflows : propriétaires sources.
## 20. Provenance et audit
Version, author/reviewer, source refs, before/after, activation status, tenant et correlation ID.
## 21. Permissions fonctionnelles
Command read/coordinate + future review/activate permission. `OPEN-010/013` restent ouvertes.
## 22. Limites et erreurs
Procedure stale/inaccessible, capability gap, owner absent, conflit ou permission refusée empêchent active trompeur.
## 23. Métriques
Plans avec owner/review date ; expired plans liés à scénarios actifs ; gaps avant activation.
## 24. Classification de livraison
`defined` / `planned`, cible native ; aucun moteur d’exécution défini.
## 25. Critères d’acceptation
**Given** un plan complet, **When** reviewer autorisé l’active, **Then** version, owner, conditions et refs sont auditables.

**Given** une procédure stale, **When** le plan est consulté, **Then** état stale-references est visible et une Task peut être créée.

**Given** aucun modèle, **When** plan est créé/revu, **Then** édition et validation manuelles suffisent.
## 26. Questions ouvertes
Quel objet distingue Operational Plan/Govern Playbook ? Quelles transitions exigent reviewer distinct ? — requirements ci-dessus. `OPEN-010/013` ouvertes.
## 27. Consommateurs documentaires
Readiness, exercises, Incident coordination, parcours Phase 5, écrans Phase 6, objets Phase 7, permissions ultérieures.
