# Permanent Continuous Delivery Lifecycle Protocol

**Mandat :** `SELECT → BRANCH → IMPLEMENT → TEST → REVIEW → PR → MERGE → POST-MERGE VERIFY → CLEANUP → NEXT`  
**Autorité :** Architecture & Delivery Governance  
**Dépôt :** `tobianahillel-afk/cmdr`  
**Référence d'intégration unique :** `main`

---

## 1. Principes Directeurs et Invariants

1. **Branche maîtresse unique (`main`) :** `main` est l'unique source de vérité et la seule cible de développement et d'intégration. Aucune branche intermédiaire de longue durée (`develop`, `foundation`, etc.) n'est tolérée.
2. **Atomicité des incréments :** 1 branche courte = 1 incrément cohérent et délimité.
3. **Plafond d'en-cours :** Maximum 2 à 4 pull requests actives simultanément hors réconciliation historique.
4. **Protection stricte de `main` :**
   - Interdiction formelle des suppressions de branche et des force-pushes destructifs (`allow_force_pushes: false`, `allow_deletions: false`).
   - Statuts CI obligatoires (`full-kernel`, `kernel`).
   - Exécution systématique sans filtre de chemin bloquant (`paths` sans exclusion).
   - Résolution obligatoire de toutes les discussions de revue (`required_conversation_resolution: true`).
5. **Invariance du Product Spec :** `cmdr-product-spec/**` demeure strictement en lecture seule (0 octet de diff vs baseline canonique `82e630fce3a7f6236e03709116eca6dd0a1acd40`).
6. **Pureté de l'exécution :** Bibliothèque standard Go exclusivement dans `product-runtime/` sans dépendance tierce externe ni contournement d'architecture.

---

## 2. Déroulement Canonique des 10 Phases du Cycle

```
+----------+      +----------+      +-------------+      +----------+      +----------+
|  SELECT  | ---> |  BRANCH  | ---> |  IMPLEMENT  | ---> |   TEST   | ---> |  REVIEW  |
+----------+      +----------+      +-------------+      +----------+      +----------+
                                                                                 |
                                                                                 v
+----------+      +----------+      +-------------+      +----------+      +----------+
|   NEXT   | <--- | CLEANUP  | <--- | POST-MERGE  | <--- |  MERGE   | <--- |    PR    |
+----------+      +----------+      +-------------+      +----------+      +----------+
```

### Phase 1 — SELECT (Sélection et Cadrage)
- Consulter le graphe de travail (`work/graph.json`), l'état consolidé (`engineering/state/current.json`) et l'algorithme d'ordonnancement canonique `reuse-first-v1` (`tools/cmdr-dev implementation-wave`).
- Identifier la prochaine capacité admissible selon la priorité produit, la valeur métier, le déblocage des dépendances et l'absence de blocage décisionnel.
- Formuler les exigences canoniques (`REQ-*`), contrats applicables et critères d'acceptation.

### Phase 2 — BRANCH (Création de Branche Dédiée)
- Vérifier la propreté de l'arbre local (`git status`).
- Récupérer le dernier état validé de `origin/main` (`git pull origin main`).
- Créer une branche courte dédiée nommée selon la convention :
  - Fonctionnalité : `feat/<cap-id>-<slug>`
  - Correctif : `fix/<defect-id>-<slug>`
  - Tâche de gouvernance : `chore/<lot-id>-<slug>`

### Phase 3 — IMPLEMENT (Implémentation Minimale)
- Appliquer strictement le principe de changement minimal (aucun refactoring opportuniste non requis).
- Implémenter les types, contrats et fonctions de la capacité dans le respect des frontières architecturales (`boundaries.json`).
- Utiliser exclusivement la stdlib Go dans `product-runtime/`.

### Phase 4 — TEST (Validation Automatisée Exhaustive)
- Rédiger des tests automatisés déterministes couvrant le chemin nominal, les cas limites, les rejets stricts et l'isolation multi-tenant.
- Ajouter des suites de tests adversariaux synthétiques (sans secret réel).
- Exécuter la suite complète sans cache (`go test -count=1 ./...`, avec `-race` et `go vet`).
- Vérifier l'absence de régression sur les packages dépendants et l'outillage `cmdr-dev`.

### Phase 5 — REVIEW (Revue Adversariale Contradictoire)
- Soumettre le code à une revue technique rigoureuse (revue indépendante L2 par un sous-agent distinct ou revue contradictoire formelle).
- Vérifier les 4 invariants majeurs : non-divulgation de secrets (logs, DTOs, erreurs), immutabilité des contrats, validation défensive et isolation des contextes.
- Valider le respect du zéro diff sur `cmdr-product-spec/**`.

### Phase 6 — PR (Publication et Déclenchement CI)
- Pousser la branche sur `origin` sans force-push.
- Ouvrir une Pull Request ciblant `main` avec description structurée (contexte, modifications, contrats, preuves de validation).
- Laisser la CI GitHub Actions (`full-kernel`, `kernel`, matrices multi-OS) s'exécuter jusqu'au passage intégral au vert.

### Phase 7 — MERGE (Fusion Contrôlée sur Main)
- Vérifier que tous les statuts requis sont au vert (`pass`).
- Résoudre toute discussion ou remarque de revue active via l'API GitHub (`resolveReviewThread`).
- Exécuter la fusion squash sur `main` via `gh pr merge <number> --squash --delete-branch --admin` (ou équivalent GitHub API).

### Phase 8 — POST-MERGE VERIFY (Vérification Post-Fusion & Health Gate)
- Basculer sur `main` local et synchroniser : `git checkout main && git pull`.
- Vérifier la propreté de l'arbre de travail (`git status`).
- Exécuter le `MAINLINE_HEALTH_GATE` permanent via `cmdr-dev mainline-health-gate --root . --json`.
- Confirmer que le statut global retourné est strictement `MAIN_HEALTHY` (exit code 0). Si `MAIN_UNHEALTHY`, suspendre tout nouveau merge et corriger immédiatement.

### Phase 9 — CLEANUP (Assainissement Immédiat)
- Supprimer la branche locale (`git branch -d <branch>`).
- Élaguer les références distantes obsolètes (`git fetch origin --prune`).
- Vérifier que `origin/main` demeure l'unique branche active sur le dépôt distant.

### Phase 10 — NEXT (Traçabilité et Transition)
- Mettre à jour les registres d'ingénierie :
  - `work/reconciliation/prs.json` et `work/reconciliation/PR_CLOSURE_REGISTER.md`.
  - `engineering/implementation/ledger.json` (enregistrement de la preuve formelle).
  - `work/graph.json` (mise à jour du statut du nœud vers `VERIFIED` ou `ACCEPTED`).
- Sélectionner la capacité suivante et boucler le cycle.
