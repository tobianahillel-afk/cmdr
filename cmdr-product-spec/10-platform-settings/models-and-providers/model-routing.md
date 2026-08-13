---
id: settings-models-and-providers-model-routing
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-13
source-of-truth: canonical
---
# Model Routing

## Objectif
Définir la configuration administrative provider-neutral de routage de modèles sans créer ni moteur runtime, ni canonical Routing Policy, ni automatic failover.

## Périmètre
Platform Settings définit les contraintes et métadonnées de sélection autorisées pour les Model Providers: allowed-use eligibility, ordre/priorité ou conditions uniquement lorsqu'elles sont sourcées, contraintes de coût/latence seulement si les sources les fournissent, fallback configuration et provenance des changements. Une sélection effective peut être observée/projetée depuis son véritable runtime owner; Settings ne revendique pas son exécution.

## Propriétaire fonctionnel
`Platform Settings Product Lead`.

## Objets concernés
`Model Provider`, Tenant et références sourcées. Il n'existe pas de canonical `Model`, `Model Route`, `Routing Policy` ou `Provider Switch` créé par ce module. Le canonical `Policy` reste Govern-owned.

## Fonctionnalités
- allowed-use eligibility;
- configuration de préférences/priorités/conditions lorsque la source les définit;
- contraintes de coût/latence comme metadata d'entrée, jamais comme promesse de runtime;
- fallback configuration;
- validation locale déterministe de configuration;
- changements versionnés/audités;
- projection d'une effective provider/model selection réellement observée;
- no silent provider switch: tout changement effectif observé expose cause/source/provenance lorsqu'elles sont disponibles.

## Fallback boundary
Fallback configuration != automatic failover. Cette spécification ne crée ni failover engine, ni weighted/round-robin/load-balancing, ni optimiseur, ni quota/rate-limit runtime, ni inference gateway. Si le runtime sélectionne un autre provider/model, Settings peut projeter l'observation mais ne prétend pas avoir exécuté le switch.

## Runtime boundary
Studio conserve Tool/Tool Call/Skill/Workflow/Human Gate/Automation Run et l'exécution agentique. Le runtime provider/routing reste hors de cette capability tant qu'une source canonique ne l'attribue pas. Settings n'exécute pas d'appel fournisseur du seul fait qu'il configure le routage.

## UX et interactions
Réutiliser `SET-MDL-001`. Toute modification affiche Tenant, cible, contraintes, provenance et effet administratif attendu. Une effective selection inconnue reste unknown; aucune bascule silencieuse ou supposée n'est affichée comme fait.

## Permissions
Réutiliser `perm.platform-settings.model-provider.read/manage` et les aliases UI existants `perm.settings.model.read/manage`. Aucun nouvel ID ou droit execute.

## États
Le routage n'introduit aucun nouvel état Model Provider. Toute projection respecte les états canoniques du Model Provider et la fraîcheur de la source observée.

## Dépendances
Model Provider; Tenant; Security; data-policy/residency owners; Administrative Audit; Studio/runtime consumers; `OPEN-008`, `OPEN-012`, `OPEN-013`.

## Critères d'acceptation
- aucune Policy/Routing Policy canonique concurrente;
- aucune sélection effective inventée;
- aucun automatic failover implicite;
- tout changement administratif est versionné et audité;
- tout changement effectif observé est non-silencieux dans la projection;
- aucune disponibilité/support provider n'est déduite de la configuration;
- chemin manuel/déterministe sans IA;
- `defined/planned` ne signifie pas runtime implémenté.

## Questions ouvertes
La portée supportée, les choix provider/delivery et l'autorité par défaut des mutations Class 2 restent respectivement gouvernés par `OPEN-008`, `OPEN-012` et `OPEN-013`.
