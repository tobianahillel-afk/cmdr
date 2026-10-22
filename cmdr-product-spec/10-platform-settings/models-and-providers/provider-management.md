---
id: settings-models-and-providers-provider-management
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-13
source-of-truth: canonical
---
# Provider Management

## Objectif
Définir l'administration provider-neutral du canonical `Model Provider` sans attribuer à Platform Settings l'exécution technique du fournisseur.

## Périmètre
Le module administre la configuration et le lifecycle du Model Provider, réalise seulement des validations locales déterministes sans effet, prépare les handoffs nécessaires aux tests externes et projette uniquement des résultats, model availability et health réellement sourcés.

## Propriétaire fonctionnel
`Platform Settings Product Lead`.

## Objets concernés
- `Model Provider`;
- `Tenant` obligatoire;
- `Environment` uniquement lorsque sourcé;
- `Secret Reference` comme référence optionnelle/sourcée, jamais comme valeur sensible.

`Integration` reste un objet distinct. Aucun generic Provider, Model, Model Catalog, Credential, Connection ou Connector canonique n'est créé.

## Fonctionnalités
- lifecycle `configured / validating / active / degraded / disabled`;
- métadonnées/configuration administrative provider-neutral;
- validation locale déterministe;
- `Tester` = demande/préconditions/handoff + projection d'un outcome réellement observé, pas exécution technique externe;
- model availability = projection d'une observation sourcée, pas moteur de découverte;
- health = projection provider-local sourcée, pas moteur de sonde;
- activation/désactivation administrative avec provenance;
- Secret Reference uniquement lorsque la source établit cette dépendance.

## Data policy / Policy boundary
Les contraintes de données, residency, usage autorisé et sécurité sont consommées depuis leurs owners. Le canonical `Policy` reste Govern-owned; cette surface ne crée pas de Routing Policy ou Safety Policy canonique concurrente et n'évalue pas la Security policy à la place de Security.

## Runtime boundary
Provider configuration != provider runtime. Aucune passerelle d'inférence, librairie fournisseur, découverte de modèles, moteur de health, gestion technique de quotas, protocole externe ou moteur de sélection runtime n'est défini ici.

## UX et interactions
Réutiliser `SET-MDL-001`; afficher état, fraîcheur, provenance et cause des données partielles. `SET-HLT-001` peut consommer/projeter du health; `SET-AUD-001` reçoit la provenance; `SET-SEC-001` reste la surface des dépendances de connexion/référence.

## Permissions
`perm.platform-settings.model-provider.read/manage`; aliases UI existants `perm.settings.model.read/manage` lorsque déjà utilisés. Aucun nouvel ID, aucune permission execute.

## Dépendances
Model Provider object; Tenant; Security; Secret Reference seulement si sourcé; Administrative Audit; Health; Studio comme consommateur; `OPEN-008`, `OPEN-012`, `OPEN-013`.

## Critères d'acceptation
- aucun état autre que les cinq états canoniques;
- aucune disponibilité, health ou validation externe fabriquée;
- aucune donnée sensible dans audit ou contexte d'assistance;
- aucun transfert de runtime ownership;
- aucune cardinalité Model Provider→Secret Reference inventée;
- toute mutation administrative est tenant-scoped, versionnée et auditée;
- `defined/planned` ne prouve ni provider support ni implémentation.

## Questions ouvertes
La disponibilité/support réels et la portée provider/delivery restent sous `OPEN-008`/`OPEN-012`; l'autorité par défaut de certaines mutations Class 2 reste `OPEN-013`.
