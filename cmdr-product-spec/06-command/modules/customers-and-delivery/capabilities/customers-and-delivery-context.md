---
id: CAP-CMD-401
title: Customers and Delivery Context
product: command
module: customers-and-delivery
owner: Command Product Lead
status: draft
delivery_status: proposed
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids: [REQ-PROD-053, REQ-PROD-013, REQ-PROD-019, REQ-PROD-033]
open_decisions: [OPEN-006, OPEN-013]
source-of-truth: proposal
---

# CAP-CMD-401 — Customers and Delivery Context

## 1. Définition
Proposition deployment-dependent présentant contexte client, portefeuille, engagements, delivery et reporting depuis Command/Reporting Engine sans imposer un modèle MSP/MSSP.

## 2. Problème utilisateur
Certains déploiements coordonnent plusieurs clients et engagements ; d’autres sont internes. Sans distinction, les fonctions économiques sont présentées comme universelles.

## 3. Objectifs
Réutiliser les fonctions indépendantes du modèle économique, documenter internal/enterprise/MSP, consommer Reporting Engine/métriques et rendre les fonctions contractuelles explicitement activées.

## 4. Non-objectifs
N’impose pas multi-client, billing ou portail, ne duplique pas Reporting Engine et ne ferme pas `OPEN-006`.

## 5. Propriétaire
Command possède seulement le contexte local et les Tasks ; Reporting Engine et sources client/contrat gardent ownership.

## 6. Utilisateurs
Principal : Service Delivery Manager. Secondaires : Customer Success, Incident Commander, Business Owner et Internal Security Lead selon déploiement.

## 7. Conditions d’entrée
Deployment model déclaré, capability activée, audience/scope tenant-client et permissions ; `OPEN-006` reste ouverte.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Deployment model | architecture/deployment owner | internal, enterprise ou MSP/MSSP | oui | configuration courante | capability `disabled-for-deployment` ou proposed |
| Command projections | Incident, Task, Result et Readiness | contexte de delivery | non | fraîcheur visible par source | contexte `partial` |
| Engagement context | source customer/contract | SLA, obligations et audience | non | version effective | aucune obligation contractuelle inférée |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Incident / Task | Command | état, owner, SLA et outcome | consulter et relier |
| Result | Govern | outcome vérifié et residual risk | consulter en projection |
| Report conceptuel | Shared Reporting Engine | draft, snapshot, audience et publication | consulter et demander selon permission |
| Customer / engagement context | source deployment | identité, contrat et obligations | consulter en projection, tenant-scoped |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Report request/draft | demander ou contribuer via Reporting Engine | Shared Reporting Engine | Command ne possède pas Report ; aucun `report.md` créé |
| Delivery follow-up Task | créer ou modifier | Command | classe 2, engagement source lié |
| Customer / contract source | aucune mutation | source deployment | projection en lecture seule |

## 11. Fonctionnalités
Reporting interne sans Customer, portfolio/engagement optionnels, report request avec citations, Tasks delivery et séparation SLA opérationnel/contractuel.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter | rôle autorisé | projections | 0 | module activé | scope visible | non |
| Créer report draft | rôle report | Report | 2 | audience, scope et template | draft Shared | workflow de reporting |
| Créer follow-up | Delivery Manager | Task | 2 | obligation ou gap | Task | OPEN-013 |
| Publier ou exporter | rôle report | Report | 2 | review et permission | résultat Shared | pas d’autorité Command |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Agréger les citations | oui | oui | oui | non nécessaire | Reporting Engine et sélection humaine |
| Préparer un report | oui | templates/snapshots | oui | brouillon attribué | template, citations et édition humaine |
| Détecter une obligation à suivre | oui | règle contractuelle si configurée | oui | suggestion attribuée | revue de l’engagement et création manuelle |
| Publier ou exporter | oui | validation/redaction | workflow possible | jamais automatiquement | revue et action humaines |

## 14. États fonctionnels
`proposed`, `disabled-for-deployment`, `configured`, `partial`, `data-stale`, `audience-restricted`, `suspended`.

## 15. États d’interface
Disabled explique le deployment ; Partial/Stale montrent les sources ; Permission denied ne révèle pas client/contrat ; Offline bloque publication.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Delivery context | projection | rôles autorisés Command | tenant/client scope, source et fraîcheur visibles |
| Report request/draft | objet Shared | Reporting Engine et reviewer | citations, snapshot, audience et ownership Shared conservés |
| Delivery Task | Task Command | Work Queue | attribuée et liée à une source d’engagement |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Customers and Delivery | création de report | Reporting Engine | tenant/client scope, citations, audience et return origin | contexte delivery restauré |
| Customers and Delivery | obligation ou gap | Task Coordination | engagement ref, due, owner et résultat attendu | contexte delivery restauré |
| Customers and Delivery | inspection contractuelle | source customer/contract | référence, version et return origin | même scope restauré |

## 18. Dépendances
`OPEN-006`, Reporting Engine, Metrics Engine, Export Engine, Business Service Catalog, CAP-CMD-105 et CAP-CMD-303.

## 19. Source de vérité
Objets Command, Result Govern, Report Shared et sources customer/contract restent séparés.

## 20. Provenance et audit
Deployment model, audience, citations, snapshot, redaction, reviewer, publication/export et correlation ID.

## 21. Permissions fonctionnelles
Command read, Shared report create/review/publish/export, source customer/contract et tenant isolation.

## 22. Limites et erreurs
Module disabled, scope client ambigu, données contractuelles absentes/stale, audience refusée ou risque cross-tenant empêchent toute promesse implicite.

## 23. Métriques
Déploiements utilisant reporting générique versus client context, reports avec citations/snapshots et Tasks liées à un engagement ; aucune cible définitive.

## 24. Classification de livraison
`proposed` / `planned`, deployment-dependent. Aucune inclusion universelle sans décision `OPEN-006` et preuve.

## 25. Critères d’acceptation
**Given** un déploiement interne sans Customer, **When** un report est créé, **Then** Reporting Engine fonctionne et le module client reste disabled/proposed.

**Given** un MSP configuré, **When** le contexte s’ouvre, **Then** tenant, client, audience et source sont explicites et isolés.

**Given** aucun modèle, **When** un report est préparé, **Then** citations, templates et édition humaine suffisent.

## 26. Questions ouvertes
Quels scénarios activent le module, quels concepts sont génériques et comment séparer SLA contractuel ? — Requirement IDs ci-dessus ; `OPEN-006` et `OPEN-013` restent ouvertes.

## 27. Consommateurs documentaires
Customer Overview seulement si retenu, Reporting, Work Queue, parcours Phase 5, écrans Phase 6 conditionnels, phase Objets et permissions ultérieures.