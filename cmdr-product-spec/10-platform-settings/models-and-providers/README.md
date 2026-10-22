---
id: settings-models-and-providers
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-13
source-of-truth: canonical
---
# Models & Providers

## Objectif
Administrer les Model Providers et la configuration de routage utilisables par Studio et les consommateurs autorisés, sans transférer à Platform Settings l'exécution d'inférence, de Tool, de Workflow, d'Automation Run ou d'un runtime provider.

## Périmètre
Platform Settings possède l'administration du canonical `Model Provider`, sa configuration tenant-scoped, son lifecycle, les validations locales déterministes, les projections de disponibilité/health sourcées et la configuration de routage. Le module ne crée ni canonical `Model`, ni `Routing Policy`, ni moteur d'inférence, ni provider SDK, ni gateway, ni failover automatique.

## Propriétaire fonctionnel
`Platform Settings Product Lead`.

## Objets concernés
- `Model Provider` — canonical, owner Platform Settings;
- `Tenant` — obligatoire;
- `Environment` — seulement lorsqu'une relation canonique le source;
- `Secret Reference` — dépendance de référence seulement lorsqu'elle est explicitement configurée/sourcée;
- `Integration` — objet distinct, aucune équivalence ni cardinalité obligatoire inventée.

Aucun canonical `Model`, `Model Catalog`, `Model Availability`, `Model Route`, `Routing Policy`, generic `Provider`, `Provider Instance`, `Provider Endpoint`, `Connection`, `Connector` ou `Credential` n'est introduit.

## Fonctionnalités
- Model Provider administrative lifecycle;
- configuration et validation locale déterministe;
- demande/handoff de test provider et projection d'un résultat réellement observé;
- projection de model availability et provider-local health lorsqu'une source autorisée les fournit;
- routing configuration, allowed-use eligibility, fallback constraints et source-backed cost/latency metadata;
- provenance des changements et de tout changement effectif de provider/model observé;
- aucun silent provider switch.

## Frontières d'exécution
`Tester` ne signifie pas que Settings possède le client/provider probe. `Model availability` et `Health` sont des projections tant qu'aucun exécuteur canonique n'est attribué. `Router` configure des contraintes et préférences; il ne crée pas un runtime routing engine. Une sélection effective reçue d'un runtime owner peut être projetée, mais Settings ne revendique pas l'exécution qui l'a produite.

Studio conserve Tool, Tool Call, Skill, Workflow, Human Gate, Automation Agent/Run et l'exécution agentique. Security conserve permissions, RBAC/ABAC, SoD, step-up et policy evaluation de sécurité. Govern conserve Policy/Decision/Approval/Response authority. `Policy` n'est pas redéfini ici.

## Secret Reference
Une relation vers `Secret Reference` reste reference-only et uniquement si elle est sourcée. Aucun secret brut, password, token, API key, private key ou credential n'est lu, stocké, journalisé ou injecté par cette spécification. Aucun Vault/KMS/HSM, algorithme cryptographique ou secret retrieval API n'est inventé.

## UX et interactions
Réutiliser `SET-MDL-001`; `SET-SEC-001`, `SET-HLT-001` et `SET-AUD-001` ne sont que des surfaces dépendantes/support lorsque nécessaire. Préserver Tenant/contexte, Inspector, deep links, états Loading/Empty/Partial/Error/Offline/Permission denied et provenance visible.

## Permissions
Réutiliser `perm.platform-settings.model-provider.read/manage`. Les aliases UI existants `perm.settings.model.read/manage` restent actifs là où déjà référencés. Aucun nouvel ID ni permission d'exécution n'est créé.

## États
Les états `Model Provider` restent exactement `configured`, `validating`, `active`, `degraded`, `disabled`. Les états d'interface n'ajoutent aucun état métier.

## Dépendances
- `05-domain-model/objects/model-provider.md`;
- `05-domain-model/objects/secret-reference.md` lorsque sourcé;
- Security permission/trust sources;
- Administrative Audit et Health comme consommateurs/projections;
- Studio comme consommateur de configuration, sans transfert de runtime ownership;
- `OPEN-008`, `OPEN-012`, `OPEN-013` restent ouverts.

## Critères d'acceptation
- Model Provider != Integration et aucun objet concurrent n'est créé;
- configuration provider != provider runtime;
- model metadata != canonical Model object;
- fallback configuration != automatic failover;
- aucune sélection/changement effectif de provider n'est présenté comme silencieux ou comme exécuté par Settings sans source;
- aucune disponibilité/health n'est fabriquée;
- aucune relation Secret Reference obligatoire n'est inventée;
- chemin manuel/déterministe disponible sans IA;
- documentary `defined/planned` ne signifie pas implémenté.

## Questions ouvertes
`OPEN-008`, `OPEN-012` et `OPEN-013` restent explicitement ouverts. Leur résolution n'est pas nécessaire pour une spécification provider-neutral, mais reste requise pour certaines décisions d'implémentation/support/autorité.
