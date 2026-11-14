# Roadmap : Capacité Suivante CAP-STD-005 (Tool Output, Result and Error Contract)

**Capacité Canonique :** `CAP-STD-005`  
**Titre :** Tool Output, Result and Error Contract  
**Module Produit :** `cmdr-studio` / `studio-foundations`  
**Spécification Canonique :** `cmdr-product-spec/09-cmdr-studio/capabilities/cap-std-005-tool-output-result-and-error-contract.md`  
**Priorité d'Ordonnancement :** Immédiate (`reuse-first-v1` — suite logique de `CAP-STD-004`)  
**Frontière Architecturale :** `studio-catalog-runtime` (`product-runtime/studio-catalog/**`)

---

## 1. Contexte et Justification Produit

Avec la finalisation et l'intégration vérifiée de **`CAP-STD-004` (Tool Input, Parameter and Validation Contract)** sur `main` (PR #225, remédiation #228, ledger synchronisé dans #231), le système dispose d'un contrat rigoureux de validation d'entrées, d'isolation multi-tenant et de masquage de secrets à l'appel.

La capacité **`CAP-STD-005`** complète ce socle fondamental en régissant les **sorties techniques, projections structurées, avertissements et erreurs d'exécution**. Elle matérialise l'invariant central :
> *Un succès technique d'outil peut être partiel, obsolète ou limité et ne doit JAMAIS être automatiquement promu en outcome métier (Evidence, Finding, ou Result Govern).*

---

## 2. Exigences Produit et Invariants Canoniques

Conformément à `cap-std-005-tool-output-result-and-error-contract.md` (baseline `82e630fc`) :

1. **Attribution et Traçabilité :**
   - Chaque sortie technique doit être liée à une version d'outil exacte (`Tool version pinned`) et à une provenance d'appel traçable (`call provenance / correlation ID`).
   - Une absence de corrélation produit une erreur explicite d'attribution.
2. **États Fonctionnels Déterministes :**
   - 8 états d'exécution canoniques : `complete`, `partial`, `warning`, `failed`, `timeout`, `unsupported`, `permission-blocked`, `runtime-unavailable`.
   - En cas d'état `partial` ou `warning`, les sorties valides déjà acquises et la corrélation sont impérativement préservées.
3. **Séparation Définition / Outcome Métier :**
   - Studio possède le contrat de sortie (`Tool output contract`).
   - Le propriétaire du runtime produit la sortie technique brute/structurée.
   - Les modules consommateurs (Govern, Investigate, Command) qualifient indépendamment leurs propres objets sans requalification automatique par Studio.
4. **Non-Divulgation et Secrets :**
   - Aucune fuite de données protégées ou d'informations d'identification brutes dans les sorties d'erreurs ou les logs.
   - Masquage conforme aux règles établies sous `CAP-STD-004`.
5. **Déterminisme et Indépendance IA :**
   - La validation du contrat de sortie et la gestion des codes d'erreurs s'exécutent de façon 100% déterministe hors modèle de langage.

---

## 3. Plan d'Implémentation Technique

### Étape 1 : Modèles de Données et Types (`product-runtime/studio-catalog/tool_output.go`)
- Structs canoniques pour `ToolOutputContract`, `ToolExecutionOutcome`, `OutputFieldDefinition`, `TechnicalWarning`, `ExecutionError`.
- Énumération stricte des 8 états d'exécution.
- Structure de corrélation et métadonnées d'audit (sans dépendance externe).

### Étape 2 : Fonctions de Validation et de Projection
- `ValidateToolOutputContract(contract ToolOutputContract) error`
- `EvaluateTechnicalOutcome(rawOutput []byte, contract ToolOutputContract, correlationID string) (ToolExecutionOutcome, error)`
- Anti-oracle et contrôle d'isolation multi-tenant.

### Étape 3 : Tests Unitaires et Adversariaux (`product-runtime/studio-catalog/tool_output_test.go`)
- Tests de couverture nominale (tous types et états supportés).
- Tests de non-promotion en outcome métier.
- Tests adversariaux de troncature, timeout, caractères malveillants et absence de correlation ID.
- Tests de non-mutation et d'isolation des locataires.

### Étape 4 : Manifeste et Intégration
- Création du lot `E10-STD-005A-OUTPUT` sous `work/lots/`.
- Mise à jour du graphe de travail `work/graph.json`.
- Enregistrement dans le registre d'implémentation `engineering/implementation/ledger.json`.

---

## 4. Statut de Préparation

- **Spécification :** PRÊTE (conforme et figée dans `cmdr-product-spec/09-cmdr-studio/capabilities/cap-std-005-tool-output-result-and-error-contract.md`).
- **Dépendances amont :** SATISFAITES (`CAP-STD-004` intégré au commit `63598553`).
- **Outillage de gouvernance :** OPÉRATIONNEL (`MAINLINE_HEALTH_GATE` actif, 100% tests verts).
- **Prochaine action :** Créer la branche `feat/cap-std-005-tool-output-contract` depuis le HEAD validé de `main` dès lancement du prochain incrément.
