---
id: CAP-INV-006
title: Saved Searches and Query Assets
product: investigate
module: signals-and-hunt
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-014
  - REQ-PROD-020
open_decisions: []
---
# CAP-INV-006 — Saved Searches and Query Assets

## 1. Définition
Versionner et partager l’intention de recherche Investigate sans la confondre avec une Saved View générique ou une Detection Rule.

## 2. Problème utilisateur
Les Queries utiles sont copiées sans auteur, prérequis, version ou validation. Saved Search, vue de résultats et règle de détection répondent à des objectifs distincts.

## 3. Objectifs
Sauvegarder, paramétrer, versionner, partager, dupliquer et déprécier des recherches ; afficher auteur, sources, champs, prérequis et dernière validation ; préserver les distinctions conceptuelles.

## 4. Non-objectifs
Ne pas posséder Saved Views ; ne pas créer ou déployer une Detection Rule ; ne pas masquer la Query ; ne pas garantir des résultats sur toutes les sources.

## 5. Propriétaire
Investigate / Signals and Hunt / Investigate Product Lead. Saved Views restent Shared ; Detection Rules restent Detection Engineering.

## 6. Utilisateurs
Principal : Threat Hunter. Secondaires : SOC Analyst et Investigation Lead reviewer. Detection Engineer futur est consommateur, non owner.

## 7. Conditions d’entrée
Query lisible et versionnée, sources/champs identifiés, auteur, scope de partage et état de validation explicites.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Query version | Shared Query | contenu exécutable | oui | version immuable | publication bloquée |
| Paramètres et variables | auteur | configuration | non | types/defaults visibles | paramètres marqués requis |
| Sources et champs | Data Sources | prérequis | oui | dernière validation | asset incompatible |
| Scope de partage | tenant policy | visibilité | oui | réévalué à l’accès | rester personnel |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Query | Shared Capabilities | version, paramètres, diagnostics | consulter, exécuter, référencer |
| Saved View | Shared Capabilities | présentation | consulter sans fusionner |
| Detection Rule | Investigate 4B.3 | référence éventuelle | consulter sans convertir |
| Case/Hunt | Investigate | contexte d’usage | relier l’asset |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Saved Search | créer, versionner, partager, déprécier | Investigate | conserve Query, paramètres, sources et validation |
| Query Asset | créer ou dupliquer | Investigate | distinct d’une Rule et d’une Saved View |
| Case/Hunt relation | lier | Investigate | référence, sans copie silencieuse |
| Saved View | aucune modification | Shared Capabilities | séparation obligatoire |

## 11. Fonctionnalités
Créer Saved Search personnelle/partagée, versionner, afficher prérequis, dupliquer avec lineage, déprécier sans supprimer, lier Case/Hunt et vérifier la compatibilité conceptuelle.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Sauvegarder | Threat Hunter | Saved Search | 2 | Query visible | draft versionné | OPEN-013 |
| Partager | Investigation Lead | Query Asset | 2 | scope et revue | asset partagé | OPEN-013 |
| Dupliquer | Analyst | Saved Search | 2 | lecture | nouveau lineage | OPEN-013 |
| Déprécier | Owner | Query Asset | 2 | raison/remplaçant | notice versionnée | OPEN-013 |
| Exécuter | Analyst | Search Job | 0 | sources/permissions | run Event Search | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Validation compatibilité | oui | oui | oui | explication | schéma et diagnostics |
| Suggestion paramètres | oui | oui | oui | oui | documentation et defaults |
| Recommandation dépréciation | oui | oui | oui | oui | règles champs/sources |
| Partage | oui | policy | workflow | non décisionnel | revue humaine |

## 14. États fonctionnels
`draft`, `validated`, `shared`, `incompatible`, `deprecated`, `superseded`, `archived`. États Draft, non machine finale.

## 15. États d’interface
Loading préserve le draft ; Empty distingue absence d’assets ; Partial nomme les prérequis manquants ; Offline reste en lecture ; Permission denied masque les assets non autorisés ; Stale expose la date de validation.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Saved Search | asset | Analysts/Hunts/Cases | versionné, paramétré, attribué |
| Validation record | événement | utilisateurs | source, date et résultat |
| Deprecation notice | migration record | consommateurs | raison et remplaçant |
| Execution handoff | Query/version | Event Search | asset original inchangé |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Query Authoring | sauvegarde | Saved Search | Query/version, paramètres, sources | éditeur |
| Saved Search | exécution | Event Search | version, valeurs, contexte | asset avec run lié |
| Saved Search | liaison | Hunt ou Case | asset, rôle, notes | asset |
| Query Asset | référence future | Detection Engineering | lineage et source | aucune Rule automatique |

## 18. Dépendances
Shared Query, Versioning, Saved Views, CAP-INV-002/003/005, Data Source schema, future 4B.3 et permissions create/share/deprecate.

## 19. Source de vérité
Query reste Shared ; Saved Search et Query Asset restent Investigate ; Saved View reste Shared ; Detection Rule reste 4B.3 ; la validation est datée.

## 20. Provenance et audit
Auteur, Query lineage, sources/champs, validation, partages/scopes, dépréciation/remplaçant et exécutions liées.

## 21. Permissions fonctionnelles
Create personal, share team/tenant, deprecate owner/reviewer, execute selon sources et cross-tenant copy interdit.

## 22. Limites et erreurs
Query inaccessible, source/champ incompatible, conflit de version, scope invalide, secret détecté, remplaçant inaccessible ou permission révoquée.

## 23. Métriques
Assets avec prérequis complets, réutilisation Hunt/Case, incompatibilités détectées avant run, délai de migration et réduction des copies ad hoc.

## 24. Classification de livraison
`defined` / `planned`, cible native. Promotion conditionnée par versioning, lineage, validation datée et scopes de partage.

## 25. Critères d’acceptation
**Given** une Query et une Saved View **When** la recherche est sauvegardée **Then** une Saved Search distincte est créée et aucune Rule ne naît.

**Given** un champ supprimé **When** la compatibilité est vérifiée **Then** l’asset devient incompatible avec raison et versions conservées.

**Given** un asset ancien **When** il est exécuté **Then** sources et permissions sont réévaluées et un nouveau Search Job est lié.

## 26. Questions ouvertes
Politique de revalidation, relation future avec Detection Rules et scopes d’approbation restent à préciser.

## 27. Consommateurs documentaires
Event Search, Hunt, Case Workspace, future Detection Engineering et future matrice de permissions.
