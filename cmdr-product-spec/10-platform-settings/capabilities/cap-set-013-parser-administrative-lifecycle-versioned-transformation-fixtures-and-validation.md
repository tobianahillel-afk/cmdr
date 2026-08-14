---
id: CAP-SET-013
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
# Parser Administrative Lifecycle, Versioned Transformation, Fixtures and Validation

## 1. Objectif

Permettre à un administrateur Platform Settings de définir et gouverner la représentation administrative d'un `Parser` canonique: cycle de vie, versions, contrat de transformation, références de schéma, fixtures, erreurs/qualité, validation de configuration et rollback administratif source-backed, sans transformer Settings en moteur d'exécution du Parser.

Cette capability reste `planned`; sa conformité documentaire ne prouve ni parser runtime, ni plugin, ni normalisation, ni support d'un format ou d'une source.

## 2. Résultats utilisateur

L'administrateur peut comprendre quelle version de Parser est administrativement définie, quel input dialect et output-schema contract sont référencés, quelles fixtures sont associées à la définition, quels défauts de configuration sont détectables localement, et quelle provenance explique un changement de version ou d'état.

Il peut distinguer `testing` d'une exécution runtime réussie, et distinguer une validation locale d'un résultat produit par un moteur externe.

## 3. Périmètre

La capacité couvre:
- création et modification de la représentation administrative Parser;
- cycle de vie canonique;
- versionnement et compatibilité administrative limitée aux métadonnées explicitement sourcées;
- input dialect et output schema comme dimensions de contrat;
- fixtures comme données de test/référence administratives;
- erreurs et métadonnées de qualité;
- validation déterministe locale de configuration et de préconditions;
- rollback vers une version/configuration source-backed lorsque explicitement permis;
- audit, provenance et handoffs vers l'exécuteur runtime lorsqu'une exécution réelle est demandée.

## 4. Hors périmètre

Sont exclus:
- parser execution engine;
- sandbox ou plugin runtime;
- normalization engine;
- stream processor;
- schema-registry implementation;
- sélection automatique de Parser;
- relation ou assignment Data Source→Parser;
- precedence, routing ou fallback Parser;
- réécriture de données ingérées;
- sélection d'ECS, OCSF, CIM, OpenTelemetry ou d'un autre standard non sourcé;
- création d'objets `Schema`, `Mapping`, `ParserCompatibility`, `ParserAssignment` ou runtime.

## 5. Propriété et frontières

Le propriétaire fonctionnel est **Platform Settings Product Lead**. `OBJ-PARSER` reste Platform Settings-owned. Platform Architecture conserve le Parser Contract, Event Contract et leurs choix techniques futurs. Security conserve l'autorisation. Studio conserve Tool/Tool Call/Automation Run et le runtime agentique. Investigate/Collection, Endpoint, Govern et les autres domaines conservent leurs propres moteurs et autorités.

Le Parser administratif n'est ni un Tool Studio, ni un Detection rule, ni un runtime Collection, ni une preuve d'Evidence/Finding.

## 6. Acteurs et rôles

Acteur principal: administrateur Platform Settings autorisé à lire/gérer Parser.

Acteurs secondaires: administrateur/analyste autorisé inspectant une définition; Platform Architecture comme propriétaire des implementation contracts; service d'audit; runtime owner recevant un handoff; consommateurs autorisés de métadonnées de version/qualité.

Aucun acteur ne reçoit implicitement autorité de parser execution par la permission Settings `.manage`.

## 7. Préconditions

- tenant connu et appliqué avant résolution;
- permission réévaluée côté serveur;
- version attendue disponible pour toute mutation;
- état source compatible avec la transition demandée;
- toute référence d'input/output schema provient d'une source explicitement disponible;
- fixture et métadonnées de qualité ne contiennent pas de secret brut;
- une demande de test runtime est séparée de la validation locale;
- aucun lien Data Source→Parser n'est présumé.

## 8. Inputs

| Input | Source | Required | Validation / boundary |
|---|---|---|---|
| Parser identifier or creation intent | Settings administrative context | yes | immutable ID once established; tenant-scoped |
| Tenant context | canonical Tenant/session context | yes | tenant-first; no cross-tenant fallback |
| Expected Parser version | current Parser representation | for mutation | stale version rejected |
| Administrative version metadata | authorized operator/source | yes for versioned change | no runtime deployment implication |
| Input dialect reference/description | Parser Contract/source metadata | when defined | no unsourced standard selection |
| Output schema reference/contract | Parser Contract/source metadata | when defined | reference/contract only; schema implementation remains external |
| Fixture metadata/content reference | authorized fixture source | optional | test/reference material; no secret-value exposure |
| Error/quality metadata | sourced validation/runtime observation | optional | source and timestamp retained; not fabricated |
| Justification / correlation ID | operator/workflow context | for auditable mutation | preserved in provenance |

## 9. Objects Read

| Object / concept | Read purpose | Authority boundary |
|---|---|---|
| `OBJ-PARSER` | state, version, tenant, administrative transformation metadata | primary canonical object |
| Tenant context | tenant-first resolution | reference only |
| Parser Contract | input/output/version/fixtures/errors/quality semantics | contract owner remains Platform Architecture |
| Event Contract | understand schema-version/correlation context where relevant | no event-processing ownership transfer |
| Fixture reference/material metadata | inspect deterministic test definition | fixture is not a runtime result |
| Audit/provenance record | reconstruct administrative version/state change | not an Evidence object |

## 10. Objects Created/Modified

| Object | Operation | Constraint |
|---|---|---|
| `OBJ-PARSER` | create administrative Parser definition | tenant mandatory; canonical lifecycle only |
| `OBJ-PARSER` | update versioned transformation metadata | version precondition; auditable |
| `OBJ-PARSER` | transition state | only canonical states/transitions |
| `OBJ-PARSER` | select prior administrative version/configuration for rollback where source-backed | versioned/reversible administrative action; does not roll back processed data/runtime automatically |
| Audit/provenance record | append change/result | actor, tenant, object, action, result, justification, version, correlation |
| None | Data Source assignment/runtime/Schema object mutation | prohibited by this capability |

## 11. États et cycle de vie

Les états canoniques restent exactement:
- `draft`;
- `testing`;
- `active`;
- `degraded`;
- `retired`.

`testing` décrit un état canonique Parser, pas la preuve qu'un moteur a exécuté une fixture avec succès. `degraded` n'invente pas une cause. `retired` n'implique pas suppression physique ou retrait d'un artefact runtime sans preuve séparée.

## 12. Comportement fonctionnel

Chaque changement administratif produit une représentation versionnée et traçable. Les validations locales peuvent vérifier tenant, état, version attendue, présence/forme de références contractuelles disponibles et cohérence déterministe de fixture/configuration sans appeler de runtime.

Le rollback administratif peut désigner une version antérieure explicitement connue et autorisée; aucune bascule automatique de pipeline, retrait de déploiement ou retraitement de données n'est déduit.

## 13. Automation/AI

| Assistance | Allowed? | Boundary |
|---|---|---|
| Explain Parser state/version/contract metadata | yes | only sourced facts |
| Explain fixture or deterministic validation error | yes | no runtime execution |
| Summarize error/quality provenance | yes | uncertainty/source retained |
| Suggest non-authoritative mapping idea | yes | suggestion does not create schema fields or Mapping object |
| Invent schema fields/standard/compatibility | no | requires canonical source/decision |
| Execute parser/plugin/sandbox/normalizer | no | runtime outside Settings |
| Create Source→Parser assignment/selection/fallback | no | no canonical relation exists |
| Change active/retired version autonomously | no | AI is not administrative authority |

## 14. Sécurité et permissions

Lecture réutilise `perm.platform-settings.parser.read`; mutation réutilise `perm.platform-settings.parser.manage`. L'écran existant peut conserver `perm.settings.source.read/manage` comme aliases UI sans bulk normalization.

Tenant-first, RBAC/ABAC, server-side enforcement, SoD et step-up restent Security-owned. `.manage` n'implique pas `execute`, `export` ou `approve`, et n'autorise aucun moteur/parser runtime.

## 15. Audit et provenance

Toute création, mutation, transition et rollback administratif enregistre acteur, tenant, objet, action, résultat, justification, version et correlation ID. Un résultat de test externe projeté identifie son origine et ne devient pas un fait produit par Settings.

Les fixtures et métadonnées de qualité n'exposent pas de secret brut et conservent leur source/classification lorsqu'elle existe.

## 16. Outputs

| Output | Consumer | Semantics |
|---|---|---|
| Versioned Parser administrative definition | Settings / authorized consumers | configuration documentary/admin, not runtime artifact proof |
| Lifecycle state | Settings/authorized consumers | canonical Parser state only |
| Deterministic validation result | administrator | strictly local/no-effect result |
| Input/output contract projection | authorized consumers | source-backed references; no schema registry implied |
| Fixture definition/reference | administrator/runtime owner | test input definition, not successful execution proof |
| Error/quality projection | administrator/Health/support | source-attributed observation |
| Audit correlation | Administrative Audit | traceability of administrative change |

## 17. Transitions/Handoffs

| From | To | Handoff | Ownership preserved |
|---|---|---|---|
| Sources & Parsers | Parser runtime owner | request/preconditions/fixture reference when real execution is required | engine remains external to Settings |
| Sources & Parsers | Platform Architecture contract | referenced input/output/version constraints | implementation contract owner unchanged |
| Sources & Parsers | Administrative Audit | version/state/change provenance | audit remains supporting surface |
| Sources & Parsers | authorized downstream consumer | Parser version/error/quality metadata | consumer gains no manage/execute right |
| Runtime source | Sources & Parsers | observed test/error/quality result | result remains source-attributed projection |
| Any consumer | Sources & Parsers | deep link/context return | tenant and permission re-evaluated |

## 18. Gestion des erreurs

Une invalidité locale produit un refus explicite, sans appeler le runtime pour « vérifier quand même ». Une erreur runtime reçue par projection garde son propriétaire/source. Un schéma ou une fixture inaccessible n'est pas traité comme inexistant lorsque la sécurité masque sa présence.

Aucun retry ou fallback Parser automatique n'est créé.

## 19. Partiel, stale et offline

Les erreurs/quality observations peuvent être partial/stale/unknown. Leur âge et source sont conservés. En offline, une définition mise en cache peut être lue si autorisée mais les mutations non garanties sont bloquées. Une fixture non disponible ne devient pas une fixture vide valide.

## 20. Tenant et environnement

`tenant-id` est obligatoire pour `OBJ-PARSER`. Toute lecture/mutation est tenant-scoped, sans fallback cross-tenant.

Environment peut être un contexte de navigation ou une métadonnée source-dependent, mais cette capability n'ajoute aucune relation Environment universelle au Parser.

## 21. Frontières Data Source et cross-product

Data Source et Parser restent distincts et sans relation canonique directe. Être dans le même module Settings n'autorise pas assignment, auto-selection, routing, fallback ou precedence.

Investigate/Detection Engineering conservent leurs règles et analyse; Collection conserve runtime d'acquisition; Endpoint conserve ses primitives; Studio conserve ses runs/outils; Govern conserve autorité; Shared conserve ses services génériques.

## 22. Frontière runtime

La capability ne définit ni parser engine, plugin runtime, sandbox, normalizer, stream processor, schema registry, deployment mechanism ni retraitement de données. `Test Parser` nécessitant une exécution réelle devient un handoff + résultat observé, jamais une exécution Settings.

`draft/defined/planned` ne signifie jamais implemented/native/integrated/deployed/available.

## 23. Contrats Schema, Parser et Event

Le Parser Contract conserve les dimensions Input dialect, Output schema, Version, Fixtures, Errors et Quality. Le format de schéma, sa version initiale et les SLO/limites restent ouverts.

L'Event Contract conserve Event envelope, Tenant, event/ingestion time, actor/source, schema version et correlation. Cette capability ne sélectionne ni ECS, OCSF, CIM, OpenTelemetry ni un autre standard, et n'implémente aucun event transformer.

## 24. Alternative déterministe / sans AI

Toutes les fonctions essentielles restent disponibles sans AI: inspection, version comparison, state/precondition validation, édition administrative autorisée, fixture inspection, rollback administratif source-backed, lecture d'erreurs/qualité sourcées et audit.

L'absence de Model Provider ne bloque aucune fonction essentielle.

## 25. Dépendances

Dépendances principales:
- `05-domain-model/objects/parser.md`;
- `10-platform-settings/sources-and-parsers/`;
- `10-platform-settings/sources-and-parsers/screens/sources-and-parsers.md`;
- `17-implementation-contracts/parser-contract.md`;
- `17-implementation-contracts/event-contract.md`;
- Permission Model;
- audit/provenance;
- runtime owner uniquement via handoff explicite.

## 26. Critères d'acceptation

**Given** un Parser tenant A en version N et un administrateur autorisé, **When** une modification de contrat/version valide est appliquée avec la version N, **Then** la nouvelle définition reste tenant-scoped, versionnée et auditée sans exécuter de parser engine.

**Given** une fixture et des références de contrat localement disponibles, **When** la validation déterministe est lancée, **Then** elle peut confirmer/refuser la cohérence de configuration mais ne prétend jamais avoir exécuté le Parser runtime.

**Given** qu'aucune relation Data Source→Parser n'est canonique, **When** une configuration tente de définir assignment, auto-selection, fallback ou precedence, **Then** la capacité refuse d'en faire une sémantique canonique et conserve ce besoin hors périmètre/différé.

**Given** un `Test Parser` qui nécessite plugin/sandbox/engine externe, **When** l'administrateur le demande, **Then** Settings transmet demande/préconditions/fixture reference et affiche seulement un résultat sourcé; aucune exécution runtime n'est attribuée à la capability.

## 27. Questions ouvertes et décisions différées

Le format de schéma, la version initiale, les SLO/limites, le runtime, la compatibilité/association Source↔Parser et le support effectif restent non résolus. `OPEN-008` et `OPEN-013` restent ouverts. Aucun standard, moteur, relation, fallback, precedence ou source supportée n'est choisi par cette définition.
