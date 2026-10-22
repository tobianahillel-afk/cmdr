---
id: settings-sources-and-parsers
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-14
source-of-truth: canonical
---
# Sources & Parsers

## Objectif

Configurer Data Sources et Parsers.

## Périmètre

Module du produit 10-platform-settings. Les objets, permissions, composants et transitions partagés sont référencés et non redéfinis.

## Propriétaire fonctionnel

Platform Settings Product Lead.

## Objets concernés

- data-source
- parser
- integration

## Fonctionnalités

- Connection.
- Schema.
- Test events.
- Health.

## Frontières fonctionnelles source-auditées

Platform Settings possède l'administration de `Data Source` et `Parser`: configuration, cycle de vie, validation strictement locale, portée/fraîcheur, versionnement et projections sourcées. Cette administration ne transfère pas à Settings l'exécution d'acquisition, collection, ingestion, connecteur, probe de santé, moteur de parsing, normalisation, traitement d'événements, registre physique de schéma ou stockage.

`Data Source` et `Parser` restent deux objets canoniques distincts. `Integration` reste un objet distinct de `Data Source`. `Connection`, `Schema`, `Mapping` et `Collector` restent des termes fonctionnels ou de contrat lorsque sourcés; ce module ne les transforme pas en objets canoniques.

Les sources canoniques actuelles n'établissent aucune relation directe `Data Source` → `Parser`. Le module ne crée donc ni assignment, compatibilité canonique, route, sélection, fallback ni precedence entre ces objets.

## Tests et exécution

`Tester` ne signifie pas automatiquement validation locale sans effet. Lorsqu'un test nécessite un service externe, un connecteur, une collection, un moteur de parsing ou tout autre runtime, Settings ne possède que la demande administrative, les préconditions, le handoff et la projection du résultat observé lorsque ces éléments sont sourcés. L'exécuteur réel conserve son propriétaire canonique.

Une projection de santé est une vue Settings d'informations sourcées; elle n'est pas l'exécution d'un health probe.

## Frontières produit

- Investigate conserve Case, Evidence, Finding, analyse, Detection Engineering et Intelligence.
- Collection conserve acquisition/exécution/résultat et custody lorsque ces responsabilités s'appliquent.
- Endpoint reste le composant d'exécution endpoint distinct.
- CMDR Studio conserve Tool, Tool Call, Automation Run et le runtime agentique.
- Govern conserve Decision, Approval, Decision Authority, Response Run et Result.
- Security conserve le Permission Model, RBAC/ABAC, isolation tenant, SoD, step-up et enforcement.

## UX et interactions

- Conserver le contexte de liste, vue et objet.
- Utiliser l’Inspector canonique.
- Afficher les six états obligatoires.
- Préserver navigation clavier et liens profonds.

## Permissions

Voir `../../14-security-permissions-and-trust/permission-model.md` et le registre des permissions. Les familles existantes `perm.settings.source.*`, `perm.platform-settings.data-source.*` et `perm.platform-settings.parser.*` sont réutilisées sans création ni normalisation globale d'identifiants.

## États

Les états métier viennent des fichiers d’objets canoniques; la page ajoute uniquement Loading, Empty, Partial, Error, Offline et Permission denied.

## Dépendances

- 03-design-system/
- 04-experience-architecture/
- 05-domain-model/
- 17-implementation-contracts/

## Critères d’acceptation

- Aucune définition d’objet ou de permission locale.
- Tous les écrans du module ont un front matter et 27 sections.
- Les transitions sont auditées et idempotentes.
- Aucun runtime ou support de source/parser n'est présenté comme implémenté par la seule documentation de Settings.
- Aucune relation `Data Source` → `Parser` n'est inventée.

## Questions ouvertes

Les questions de format/version initiale de schéma, SLO/limites, support effectif des sources et autorité par défaut des mutations réversibles restent dans leurs sources/OPEN existantes; ce module ne les résout pas.
