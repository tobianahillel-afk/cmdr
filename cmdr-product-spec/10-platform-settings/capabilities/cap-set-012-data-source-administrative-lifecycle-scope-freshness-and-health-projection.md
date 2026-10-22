---
id: CAP-SET-012
type: capability
domain: 10-platform-settings
module: sources-and-parsers
status: draft
delivery_status: defined
delivery_mode: planned
owner: Platform Settings Product Lead
updated: 2026-08-14
source-of-truth: canonical
requirements:
  - REQ-PROD-003
  - REQ-PROD-004
  - REQ-PROD-005
  - REQ-PROD-006
  - REQ-PROD-008
  - REQ-PROD-009
  - REQ-PROD-010
  - REQ-PROD-011
  - REQ-PROD-012
  - REQ-AI-001
  - REQ-AI-002
  - REQ-AI-003
  - REQ-AI-004
  - REQ-AI-007
  - REQ-AI-008
  - REQ-AI-009
  - REQ-AI-010
  - REQ-AI-011
open_decisions:
  - OPEN-008
  - OPEN-013
---
# Data Source Administrative Lifecycle, Scope, Freshness and Health Projection

## 1. Objectif

Permettre à un administrateur Platform Settings de définir, inspecter et faire évoluer la configuration administrative d'un `Data Source` canonique, d'en comprendre la portée et la fraîcheur, de le désactiver de manière sûre et d'observer une santé sourcée sans attribuer à Settings l'exécution d'acquisition, de collection, d'ingestion ou de probe externe.

Cette capacité est documentaire et `planned`; elle ne constitue pas une preuve d'implémentation ou de disponibilité d'un adaptateur de source.

## 2. Résultats utilisateur

L'administrateur peut comprendre quel Data Source est configuré pour quel tenant, quel est son état canonique, quelles métadonnées de portée sont effectivement sourcées, quelle est la fraîcheur des observations disponibles, quelles validations locales bloquent une mutation et quelle provenance explique chaque changement.

Il peut également distinguer une configuration administrative valide d'un runtime réellement disponible: `configured` ou `active` ne garantit ni collecte effective ni complétude de données.

## 3. Périmètre

La capacité couvre exclusivement:
- création et modification de configuration administrative d'un Data Source;
- inspection et transition de son cycle de vie canonique;
- portée explicitement fournie par les sources canoniques;
- fraîcheur et état de santé comme projections sourcées;
- validation locale et déterministe de préconditions disponibles;
- désactivation administrative sûre et versionnée;
- audit, corrélation et provenance;
- handoff vers les propriétaires runtime lorsqu'une opération exige une exécution externe.

## 4. Hors périmètre

Sont explicitement exclus:
- acquisition et collection runtime;
- moteur d'ingestion;
- exécution de connecteur;
- health probe externe;
- event-processing engine;
- parser runtime ou sélection d'un Parser;
- normalisation, stockage ou schéma physique;
- support effectif d'un fournisseur, protocole, API ou type de source;
- création de `Connection`, `Collector`, `Schema`, `Mapping` ou autre nouvel objet canonique;
- création d'un objet ou d'une relation Data Source→Parser.

## 5. Propriété et frontières

Le propriétaire fonctionnel est **Platform Settings Product Lead**. L'objet canonique `OBJ-DATA_SOURCE` reste Platform Settings-owned. Security conserve Permission Model, RBAC/ABAC, isolation tenant, SoD, step-up et enforcement. Investigate/Collection, Endpoint, Studio, Govern et les autres propriétaires runtime conservent leurs responsabilités respectives.

`Integration` reste un objet distinct; une référence éventuelle n'en fait ni le Data Source lui-même ni un moteur de connexion. Une `Secret Reference` éventuelle reste reference-only.

## 6. Acteurs et rôles

Acteur principal: administrateur Platform Settings autorisé à lire ou gérer Data Source.

Acteurs secondaires: opérateur ou analyste autorisé qui inspecte état/fraîcheur; systèmes source de health/freshness; service d'audit; propriétaire runtime recevant un handoff explicite.

Aucun rôle administratif Settings ne reçoit automatiquement autorité d'exécution, d'export ou de réponse Govern.

## 7. Préconditions

- le tenant est connu avant résolution du Data Source;
- l'acteur dispose de la permission requise au moment de l'opération;
- l'identifiant et la version de l'objet sont stables;
- toute portée ou référence externe affichée possède une source identifiable;
- une mutation vérifie les préconditions d'état et de version disponibles;
- un test impliquant un runtime externe ne peut pas être requalifié en validation locale.

## 8. Inputs

| Input | Source | Required | Validation / boundary |
|---|---|---|---|
| Data Source identifier or creation intent | Settings administrative context | yes | immutable ID once established; tenant-scoped |
| Tenant context | canonical Tenant/session context | yes | tenant-first resolution; no cross-tenant fallback |
| Expected object version | Data Source current representation | for mutation | reject stale version rather than overwrite silently |
| Administrative configuration metadata | authorized operator input | when changing configuration | validate only fields supported by current canonical/source contract |
| Scope metadata | sourced configuration/context | when available | source-dependent; must not invent Environment relation |
| Freshness observation | health/runtime source | optional | carries timestamp/source; absence is not currentness |
| Health observation | Health source/contract | optional | projection only; not evidence of a Settings probe |
| Integration or Secret Reference identifier | canonical reference source | optional and source-dependent | reference only; no raw secret read-back |
| Justification / correlation ID | operator/workflow context | for auditable mutation | retained in audit provenance |

## 9. Objects Read

| Object / concept | Read purpose | Authority boundary |
|---|---|---|
| `OBJ-DATA_SOURCE` | state, version, tenant, administrative metadata and provenance | primary canonical object |
| Tenant context | enforce tenant-first scope | reference only; no Tenant redefinition |
| `OBJ-INTEGRATION` | inspect a source-backed connection reference where present | Integration remains distinct from Data Source |
| `OBJ-SECRET_REFERENCE` | inspect reference metadata where present and authorized | never reveal underlying secret value |
| Health Contract projection | read status/freshness/details scope when sourced | Settings does not own the probe executor |
| Local audit/provenance record | reconstruct administrative change | audit reference does not create a new Evidence object |

## 10. Objects Created/Modified

| Object | Operation | Constraint |
|---|---|---|
| `OBJ-DATA_SOURCE` | create administrative representation | tenant mandatory; state begins only as allowed by canonical lifecycle |
| `OBJ-DATA_SOURCE` | update sourced configuration metadata | optimistic/version precondition; auditable |
| `OBJ-DATA_SOURCE` | transition lifecycle state | only canonical states/transitions; invalid transition rejected server-side |
| `OBJ-DATA_SOURCE` | disable/reactivate where source-backed | reversible/versioned administrative mutation; `OPEN-013` remains open |
| Audit/provenance record | append mutation outcome | actor, tenant, action, result, justification, version, correlation |
| None | Integration/Secret Reference/Parser mutation | this capability does not mutate those objects merely because it references them |

## 11. États et cycle de vie

Les états canoniques restent exactement:
- `configured`;
- `active`;
- `degraded`;
- `disabled`.

La capacité ne crée aucun cinquième état. Un état UI Loading/Empty/Partial/Error/Offline/Permission denied reste un état de présentation, pas un état Data Source. `degraded` ne prouve pas la cause; `active` ne prouve pas une ingestion complète.

## 12. Comportement fonctionnel

La création ou modification produit une nouvelle version administrative traçable. La validation locale peut refuser tenant incohérent, état invalide, référence introuvable selon les données disponibles ou version périmée. La désactivation administrative doit indiquer cible, effet administratif attendu et conséquences connues sans promettre l'arrêt d'un runtime non possédé.

Les observations health/freshness sont affichées avec leur source et leur âge. Une absence d'observation produit une information manquante/partielle, jamais un `healthy` implicite.

## 13. Automation/AI

| Assistance | Allowed? | Boundary |
|---|---|---|
| Explain Data Source state/configuration | yes | must use sourced facts and expose uncertainty |
| Summarize freshness/health provenance | yes | may not fabricate observations |
| Explain deterministic validation failure | yes | no authority to override validation |
| Suggest a safe administrative next step | yes | suggestion only; permission and human/admin action still required |
| Invent scope, Environment relation or source support | no | unsourced data forbidden |
| Execute collection, ingestion, connector or probe | no | runtime belongs to separately sourced owner |
| Read/reveal Secret value | no | Secret Reference remains reference-only |
| Activate/disable autonomously | no | AI is not administrative authority |

## 14. Sécurité et permissions

Lecture réutilise `perm.platform-settings.data-source.read`; mutation réutilise `perm.platform-settings.data-source.manage`. `SET-SRC-001` peut continuer à référencer les aliases UI existants `perm.settings.source.read/manage` sans normalisation globale.

Tenant-first resolution, RBAC/ABAC, server-side enforcement, SoD et step-up s'appliquent selon Security. Lire n'implique ni exporter, ni exécuter, ni approuver. `.manage` ne confère pas un droit de probe, collection, ingestion ou response execution.

## 15. Audit et provenance

Toute mutation enregistre acteur, tenant, objet, action, résultat, justification, version et identifiant de corrélation. Tout refus d'accès ou cross-scope est audité selon Security.

Les health/freshness observations conservent source, timestamp/fraîcheur et détails de portée disponibles. La provenance distingue explicitement un résultat observé d'une validation locale Settings.

## 16. Outputs

| Output | Consumer | Semantics |
|---|---|---|
| Versioned Data Source administrative state | Settings / authorized consumers | canonical administrative representation, not runtime proof |
| Deterministic validation result | administrator | local/no-effect result only |
| Scope representation | authorized UI/consumer | only sourced dimensions; Environment remains source-dependent |
| Freshness projection | Settings/Health/consumers | age/source visible; stale/unknown preserved |
| Health projection | Settings/Health/Command or other authorized consumer | source-attributed status; no probe ownership transfer |
| Disable/change result | administrator/audit | administrative outcome only; runtime consequence not fabricated |
| Audit correlation | audit/supporting screens | immutable trace to administrative action |

## 17. Transitions/Handoffs

| From | To | Handoff | Ownership preserved |
|---|---|---|---|
| Sources & Parsers | Platform Health | sourced health/freshness projection | Health presentation does not become probe execution |
| Sources & Parsers | Integration/Secrets administration | reference inspection/navigation when sourced | Integration/Secret Reference remain separate objects |
| Sources & Parsers | runtime/collection owner | request/preconditions when external source test or execution is required | runtime executor remains external to this capability |
| Sources & Parsers | Administrative Audit | correlation/provenance link | audit remains supporting surface |
| Sources & Parsers | Shared Search / Command / Investigate | authorized source context/projection | consumer does not gain mutation permission or source access |
| Any consumer | Sources & Parsers | deep link/context return | tenant and source permission are re-evaluated |

## 18. Gestion des erreurs

Les erreurs de validation locale identifient le champ/précondition concerné sans inventer de correctif runtime. Une erreur externe projetée conserve son origine et son correlation ID lorsqu'ils existent. Une référence inaccessible est distinguée d'une référence inexistante si la sécurité empêche d'en révéler l'existence.

Aucun retry externe automatique n'est défini par cette capacité.

## 19. Partiel, stale et offline

Une health/freshness observation stale reste visible comme stale. Une observation manquante produit Partial/Unknown selon la source UI, jamais une valeur par défaut rassurante. En offline, les mutations non garanties sont interdites; la dernière représentation valide peut être lue avec son âge si la politique l'autorise.

## 20. Tenant et environnement

`tenant-id` est obligatoire pour `OBJ-DATA_SOURCE`. Toute résolution et mutation est tenant-scoped. Aucun fallback cross-tenant n'existe.

Environment peut être conservé comme contexte de navigation ou comme métadonnée lorsque la source canonique l'établit, mais cette capacité n'ajoute pas une relation Environment universelle au Data Source.

## 21. Frontières cross-product

Command peut consommer un contexte de source mais ne devient pas propriétaire du Data Source. Investigate peut consommer des événements/résultats et signaler un gap sans acquérir l'administration. Collection conserve l'acquisition/exécution/custody. Endpoint conserve son runtime. Studio conserve l'agentic runtime. Govern conserve l'autorité d'action. Shared conserve ses mécanismes transverses.

## 22. Frontière runtime

La capacité ne définit aucun collector, connector engine, ingestion engine, network client, health probe, parser runtime, stream processor, storage engine ou adapter fournisseur. Un `Test Source` externe est une demande/handoff avec résultat observé, pas une exécution Settings.

Le statut documentaire `defined/planned` n'est jamais présenté comme `implemented`, `native`, `integrated`, `deployed` ou `available`.

## 23. Contrats Health et Event

Le Health Contract est consommé pour `Component`, `Status`, `Capabilities affected`, `Freshness`, `SLO` et `Details scope` lorsqu'ils sont effectivement disponibles. Aucun SLO n'est inventé.

L'Event Contract peut fournir tenant, actor/source, timestamps, schema version et correlation à des consommateurs runtime; cette capacité ne construit pas l'envelope ni le moteur d'ingestion.

## 24. Alternative déterministe / sans AI

Toutes les fonctions essentielles sont réalisables sans AI: inspection, comparaison de version, validation de tenant/état, modification administrative autorisée, désactivation, lecture de fraîcheur/santé sourcée et consultation de provenance.

L'absence de Model Provider n'empêche aucune fonction essentielle de cette capacité.

## 25. Dépendances

Dépendances canoniques principales:
- `05-domain-model/objects/data-source.md`;
- `05-domain-model/objects/integration.md` lorsque référencé;
- `05-domain-model/objects/secret-reference.md` lorsque référencé;
- `10-platform-settings/sources-and-parsers/`;
- `10-platform-settings/sources-and-parsers/screens/sources-and-parsers.md`;
- `17-implementation-contracts/health-contract.md`;
- `17-implementation-contracts/event-contract.md`;
- Permission Model et audit/provenance.

## 26. Critères d'acceptation

**Given** un Data Source du tenant A en version N et un administrateur autorisé du tenant A, **When** une modification administrative valide est appliquée avec la version N, **Then** la nouvelle représentation reste tenant-scoped, versionnée et auditée sans créer d'objet Connection/Collector ni exécuter d'ingestion.

**Given** une tentative de mutation avec une version périmée ou un tenant différent, **When** la précondition est évaluée, **Then** la mutation est refusée côté serveur et le refus est audité sans fallback cross-tenant.

**Given** une observation de santé vieille de plusieurs intervalles ou sans timestamp fiable, **When** elle est affichée, **Then** la vue expose sa fraîcheur/incertitude et ne fabrique pas un état healthy actuel.

**Given** que `Test Source` nécessite un connecteur ou un probe externe, **When** l'administrateur le demande, **Then** Settings produit uniquement demande/préconditions/handoff et projection du résultat sourcé; aucune exécution externe n'est attribuée à cette capacité.

## 27. Questions ouvertes et décisions différées

`OPEN-008` reste ouvert pour disponibilité/support des plateformes et sources. `OPEN-013` reste ouvert pour l'autorité par défaut des mutations réversibles de classe 2. Les SLO/limites effectifs restent dans leurs contrats sources. Aucun standard, adaptateur, runtime, source supportée ou relation Source→Parser n'est décidé par cette capability definition.
