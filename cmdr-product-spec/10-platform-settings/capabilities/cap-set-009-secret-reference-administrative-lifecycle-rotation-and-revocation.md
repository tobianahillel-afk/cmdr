---
id: CAP-SET-009
title: Secret Reference Administrative Lifecycle, Rotation and Revocation
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
delivery_status: defined
delivery_mode: planned
updated: 2026-08-13
requirement_ids: [REQ-PROD-006, REQ-PROD-009, REQ-PROD-010, REQ-PROD-011, REQ-PROD-012, REQ-AI-001, REQ-AI-003, REQ-AI-004, REQ-AI-009, REQ-AI-010]
open_decisions: [OPEN-013]
source-of-truth: canonical
---
# CAP-SET-009 — Secret Reference Administrative Lifecycle, Rotation and Revocation

## 1. Définition
Capability Platform Settings qui définit l’administration provider-neutral du `Secret Reference` canonique, de son lifecycle `pending/active/rotating/expired/revoked`, de la rotation/révocation de la **référence administrative** et de son audit, sans exposer ni administrer directement le secret sous-jacent.

## 2. Problème utilisateur
Un Platform Administrator doit pouvoir enregistrer et gérer une référence de secret tenant-scoped, suivre rotation, expiration et révocation, comprendre les refus Security et conserver une provenance exploitable sans lire la valeur du secret ni confondre Secret Reference avec Credential, API key, token, certificat ou opération provider externe.

## 3. Objectifs
Fournir un lifecycle administratif unique et auditable pour Secret Reference; garantir reference-only/no read-back; appliquer Tenant, permissions, SoD, step-up et privacy; distinguer rotation/révocation de la référence des opérations sur le secret externe; conserver un chemin manuel/déterministe complet et un handoff explicite lorsque l’opération sous-jacente appartient à un autre owner.

## 4. Non-objectifs
Ne pas créer de raw Secret, Credential, APIKey, Token, Certificate, External Account ou nouvelle relation Principal; ne pas générer/lire/écrire une valeur secrète; ne pas sélectionner Vault/KMS/HSM, algorithme crypto, provider/API/protocole, schéma physique, Screen ID ou Permission ID; ne pas revendiquer rotation/révocation externe du credential matériel.

## 5. Propriétaire
`Platform Settings Product Lead` est l’unique capability owner de l’administration `Secret Reference`. Security reste owner du Permission Model, protection des secrets, cryptographie/key-management, RBAC/ABAC, tenant isolation, SoD, step-up et autorisation. Les sources actuelles n’attribuent pas à Settings l’exécution de la rotation/révocation du secret matériel ou d’un credential provider.

## 6. Utilisateurs
Platform Administrator principal; Security Administrator pour contraintes et contrôles; consommateurs autorisés de références comme Studio, Endpoint ou Integration secondaires. Un consommateur d’une Secret Reference n’obtient ni sa valeur ni l’ownership administratif.

## 7. Conditions d’entrée
Tenant résolu avant la cible; acteur administratif identifiable; permission serveur actuelle; Secret Reference stable lorsqu’existante; règles Security disponibles; aucune opération ne dépend d’un secret brut exposé. Environment n’est utilisé que si une relation canonique explicite le fournit.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Secret Reference cible ou intention d’enregistrement | `SET-SEC-001` / administrateur | référence ou intention | oui | courante | aucune mutation |
| Tenant | contexte Settings | référence canonique | oui | courante | refus sans fallback |
| version/état de la référence | objet Secret Reference | état + version | oui pour mutation | courante | relecture requise |
| metadata de référence non secrète | objet/module | type de référence/locator opaque/policy reference selon source | selon opération | courante | validation refusée |
| permission/SoD/step-up context | Security | décision d’autorisation | oui pour mutation | courante | refus serveur |
| outcome externe sous-jacent | owner externe réellement compétent | succès/error + provenance, sans valeur | seulement si un handoff a été exécuté | observé | aucun succès déduit |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Secret Reference | Platform Settings | id, tenant-id, version, état, metadata non secrète, timestamps/provenance | read/manage selon action |
| Tenant | Platform Settings | id et scope | read |
| Environment | Platform Settings | référence uniquement si source canonique | read si applicable |
| Permission context | Security | RBAC/ABAC/SoD/step-up/tenant outcome | consume |
| Principal acteur | Platform Settings | référence d’acteur pour provenance uniquement | read selon autorisation |

Aucune valeur secrète, password, API key, token, private key ou certificat privé n’est une entrée lisible de cette capability.

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Secret Reference | enregistrement, metadata/lifecycle mutation, rotation/révocation de référence | Platform Settings | reference-only, états canoniques, Tenant obligatoire, version/provenance conservées |
| Administrative audit event | émission/référence | Platform Settings / Administrative Audit | metadata minimale, action/outcome/corrélation, jamais de secret |

Aucun objet raw Secret, Credential, APIKey, Token ou Certificate n’est créé. La table S10 ne sous-entend aucun vault write ou provider credential mutation.

## 11. Fonctionnalités
Enregistrer une Secret Reference, lire sa metadata autorisée, valider localement les préconditions, appliquer une transition de lifecycle de référence, représenter `rotating`, `expired` et `revoked`, préparer/handoff une opération externe lorsqu’une source en attribue l’exécution ailleurs, projeter un outcome réellement observé, appliquer access-policy reference seulement comme metadata sourcée et auditer toute opération.

## 12. Actions utilisateur
Class 0: inspecter metadata/état/provenance autorisés. Class 1: uniquement validation déterministe sans effet ni secret read-back. Class 2: enregistrement ou mutation réversible/versionnée de la Secret Reference lorsque les règles existantes le permettent, y compris transition administrative sourcée; `OPEN-013` reste ouvert pour la politique par défaut. Une rotation/révocation du secret externe n’est ni classée ni revendiquée ici faute d’owner/effect canonique attribué.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| expliquer état/metadata non secrète | oui | oui | oui | oui | état canonique + règles |
| identifier expiration/staleness | oui | oui | oui | oui | timestamps/état sourcés |
| préparer un plan de rotation | oui | oui | oui comme suggestion | oui, suggestion seulement | étapes/règles administratives |
| mutation de référence | oui | oui | non autonome par IA | suggestion seulement | action explicite autorisée |
| expliquer un outcome externe reçu | oui | oui pour faits bruts | oui | oui avec provenance/incertitude | outcome brut + provenance |

L’IA ne reçoit aucune valeur secrète, ne génère pas de production credential, ne tourne/révoque pas un secret externe de façon autonome et ne contourne pas RBAC/ABAC/SoD/step-up.

## 14. États fonctionnels
États `Secret Reference` exactement: `pending`, `active`, `rotating`, `expired`, `revoked`. `rotating` décrit le lifecycle de la référence administrative et ne prouve pas que Settings génère ou remplace le secret sous-jacent. `revoked` sur la Secret Reference ne prouve pas la révocation d’un token/API key/certificat externe.

## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied, Stale et Sensitive metadata unavailable exposent cause/fraîcheur sans révéler la valeur. Offline interdit toute mutation exigeant garantie serveur. Un outcome externe absent reste unknown/pending plutôt que supposé réussi.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Secret Reference state projection | référence canonique | Settings/consommateurs autorisés | état + Tenant + version, aucune valeur |
| rotation/revocation administrative outcome | résultat de mutation de référence | administrateur/audit | succès/refus + cause + corrélation |
| external-operation handoff | référence + contexte non secret | Security/provider/runtime owner lorsqu’il est canonique | aucune exécution matérielle revendiquée par Settings |
| observed external outcome projection | résultat/error sourcé | administrateur/audit | reflète uniquement un résultat reçu, sans secret |
| audit provenance | événement/référence | `SET-AUD-001` | metadata minimale et traçable |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| `SET-SEC-001` | register/inspect/rotate/revoke reference | Secret Reference owner surface | reference/tenant/state/version | projection mise à jour |
| Settings | mutation de référence | `SET-AUD-001` | acteur/cible/action/outcome/correlation id | audit immutable |
| Settings | opération sur secret externe réellement requise | owner externe non attribué par cette capability | Secret Reference ref + intention + contraintes, jamais valeur brute | outcome externe si produit |
| consommateur autorisé | resolve reference metadata | Settings/Security boundary | référence + tenant + permission context | metadata/deny, jamais secret read-back |

Aucun handoff ne transforme Secret Reference en credential matériel ni ne transfère son ownership administratif.

## 18. Dépendances
Objet `Secret Reference`, Tenant, Environment seulement si sourcé, `SET-SEC-001`, Secrets & Connections, Administrative Audit, Security Permission Model, ABAC, tenant isolation, SoD, step-up, secrets-and-key-management, audit/immutability, provenance/integrity, privacy/minimization et service identity lorsque réellement référencée. Studio/Endpoint/Integrations consomment des références sans ownership.

## 19. Source de vérité
`05-domain-model/objects/secret-reference.md` pour objet/états; `10-platform-settings/secrets-and-connections/secret-management.md` pour Reference only/Rotation/Access policy/Audit; `secrets-and-connections/README.md` et `SET-SEC-001` pour actions; Security `secrets-and-key-management.md`, Permission Model, SoD et step-up pour contraintes. `service-identity.md` mentionne credential rotation sous Security mais ne transfère pas ce runtime à Settings.

## 20. Provenance et audit
Toute création/enregistrement, modification, transition, rotation/révocation de référence, handoff externe, outcome reçu ou refus conserve acteur, Tenant, Secret Reference ref, état/version, justification pertinente, timestamps et correlation id selon les mécanismes canoniques. La valeur secrète et toute matière sensible non nécessaire sont exclues des logs, exemples, GWT, audit et contexte IA.

## 21. Permissions fonctionnelles
Consomme `perm.platform-settings.secret-reference.read` et `perm.platform-settings.secret-reference.manage` côté objet ainsi que `perm.settings.secret.read-metadata` / `perm.settings.secret.manage` côté `SET-SEC-001` lorsque applicables. Leur coexistence n’autorise aucun bulk rename. `.manage` n’est pas étendu à la lecture du secret ou à une opération provider externe non canonique.

## 22. Limites et erreurs
Secret Reference ≠ raw Secret ≠ Credential ≠ API Key ≠ Token ≠ Certificate ≠ Principal. Reference rotation ≠ underlying-secret rotation; reference revocation ≠ external credential revocation. Toute cible cross-tenant, version stale, transition invalide, SoD/step-up insatisfait ou permission absente est refusée sans fallback ni fuite. L’absence d’owner externe bloque l’exécution externe plutôt que de l’inventer.

## 23. Métriques
Références enregistrées, transitions administratives, rotations/révocations de référence acceptées/refusées, expirations, refus cross-tenant, échecs SoD/step-up, outcomes externes success/error/unknown seulement lorsqu’observés, erreurs de provenance. Aucune mesure ne compte ou expose des valeurs secrètes; aucune cible SLO/KPI ou efficacité crypto/provider n’est revendiquée.

## 24. Classification de livraison
`defined / planned`. Cette capability est documentaire et provider-neutral; elle ne prouve aucun Vault/KMS/HSM, provider API, rotation engine, credential lifecycle externe, release ou disponibilité runtime. Le PASS documentaire futur ne signifiera pas implémentation.

## 25. Critères d’acceptation
**Given** une Secret Reference tenant-scoped et une permission manage valide, **When** une mutation administrative sourcée est appliquée, **Then** seuls les états canoniques de la référence, sa version et sa provenance changent sans lire la valeur secrète.  
**Given** une Secret Reference `active`, **When** une rotation administrative de référence est engagée, **Then** l’état peut suivre le lifecycle sourcé sans prétendre que Settings a généré/écrit un secret externe.  
**Given** une Secret Reference `revoked`, **When** son état est affiché, **Then** l’interface ne prétend pas qu’un credential provider externe a été révoqué sans outcome canonique observé.  
**Given** une opération externe nécessaire mais aucun owner technique canonique attribué, **When** l’administrateur demande la rotation/révocation sous-jacente, **Then** Settings limite son contrat au handoff/résultat et n’invente ni vault call ni provider API.  
**Given** une cible d’un autre Tenant ou un contrôle SoD/step-up refusé, **When** une mutation est demandée, **Then** elle est refusée sans fallback ni fuite.  
**Given** l’IA indisponible, **When** la référence est administrée, **Then** lecture metadata, validation et mutation autorisée restent possibles manuellement/déterministiquement.

## 26. Questions ouvertes
`OPEN-013` reste ouvert pour la politique d’autorité par défaut des mutations Class 2. L’owner et les mécanismes de rotation/révocation du secret matériel ou du credential externe restent non attribués par cette capability; aucune nouvelle OPEN n’est créée automatiquement.

## 27. Consommateurs documentaires
Secrets & Connections, Integrations, Models & Providers, Sources & Parsers, Administrative Audit, Security, Studio, Endpoint, Command, Investigate, Govern, Shared, Quality et Roadmap. Les consommateurs reçoivent uniquement des références/projections autorisées et jamais la valeur secrète par ce contrat.
