---
id: CAP-INV-536
title: Intelligence Correction, Retraction and Product Supersession
product: investigate
module: threat-intelligence
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-06
requirement_ids: [REQ-PROD-014, REQ-PROD-019, REQ-PROD-020, REQ-PROD-055, REQ-INV-006, REQ-AI-002, REQ-SEC-001, REQ-SEC-002, REQ-UX-010]
open_decisions: [OPEN-013, OPEN-018, OPEN-019]
source-of-truth: canonical
---
# CAP-INV-536 — Intelligence Correction, Retraction and Product Supersession

## 1. Définition
Signaler et évaluer une erreur candidate, préparer correction, rétractation, nouvelle version, retrait d’usage actif, notifications et supersession sans supprimer l’historique ni exécuter un rappel externe.

## 2. Problème utilisateur
Une erreur peut être corrigée silencieusement, laissant consommateurs, Indicators ou conclusions sur une version obsolète.

## 3. Objectifs
- identifier version, conclusions, sources, relations et consommateurs affectés.
- préparer correction/rétractation, notifications et nouvelle version.
- conserver anciennes conclusions, justification, auteur et timestamps.

## 4. Non-objectifs
Aucune API, protocole, suppression d’historique, rappel externe exécuté, runtime mutation, règle, blocage, réponse, partage externe, Cloud/Mobile ou écran détaillé.

## 5. Propriétaire
Investigate possède Correction Record, Retraction Record et Product Supersession. Shared possède Versioning, Notifications, Trace et Recovery. Govern/destination owner possèdent rappel externe ou révocation gouvernée. Chaque consumer conserve ses objets.

## 6. Utilisateurs
Principal : **Intelligence Manager**. Secondaires : Analyst, Reviewer, Consumer Owner, Detection Engineer, SOC Analyst, Approver et Auditor autorisés.

## 7. Conditions d’entrée
Issue signalée, version affectée, sources, conclusions, consumers, restrictions, owner, reviewer, permissions et return origin sont explicites ; sinon assessing/partial/blocked.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Affected product/version | CAP-INV-526/529 | content and publication lineage | oui | immutable selected version | no correction |
| Error evidence or feedback | CAP-INV-532/534 / consumer | issue and supporting context | oui | current evidence | assessing |
| Consumers and downstream links | CAP-INV-529..531 / Shared Linking | affected access and handoffs | oui | current links | notification blocked |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Product/Publication/Assessment | Investigate | affected version and conclusions | lecture/lien |
| Indicator/Watchlist/Handoff projection | Investigate / destination owners | downstream impact | lecture limitée |
| Version/Trace/Notification | Shared | history and consumer lineage | lecture/lien |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Correction Record | créer, contester, versionner, superseder | Investigate concept | correction ≠ history deletion |
| Retraction Record | proposer, revoir, retirer, superseder | Investigate concept | retraction ≠ deletion |
| Product Supersession | créer/lier | Investigate using Shared Versioning | superseded ≠ false |

## 11. Fonctionnalités
Identifier erreur/version/consumers/conclusions/sources/relations ; préparer correction, rétractation, nouvelle version et notifications ; retirer de l’usage actif ; préserver historique/justification/auteur/timestamps ; préparer recall externe futur.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| consulter/comparer | Intelligence Manager | affected lineage | 0 | lecture autorisée | vue sourcée | non |
| analyser impact/diff | Intelligence Manager | result/Tool Call | 1 | version/consumers visibles | résultat attribué | selon politique |
| créer correction/retraction/supersession | Intelligence Manager | records locaux | 2 | owner/reviewer explicites | disposition réversible | OPEN-013/019 |
| exécuter rappel externe ou runtime revocation | aucun rôle local | external/runtime object | 3 | Govern/owner requis | aucune exécution | obligatoire |
| supprimer historique/provenance | aucun rôle local | trace | 4 | interdit | refus audité | strict |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| détecter versions/consumers affectés | oui | linking/diff | oui | suggestion sourcée | tables/timeline |
| préparer correction/retraction | oui | template versionné | oui | draft attribué | formulaire/checklist |
| préparer notifications | oui | règles déterministes | oui | résumé possible | notification template |
| supprimer ou rappeler extérieurement | owner humain | contrôles | jamais autonome | jamais décisionnaire | Govern/owner workflow |

Toute automatisation expose initiateur, moteur/version, sources, paramètres, erreurs, incertitude et disposition humaine.

## 14. États fonctionnels
`issue-reported`, `assessing`, `correction-required`, `correction-in-progress`, `corrected`, `retraction-proposed`, `retracted`, `superseded`, `disputed`, `withdrawn`.

## 15. États d’interface
Loading conserve contexte ; Empty distingue aucune issue/interdiction ; Partial expose consumers inconnus ; Error conserve valide ; Offline stale/read-only ; Permission denied ne révèle rien ; Stale expose versions ; Conflict offre diff/recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Correction Record | record versionné | CAP-INV-526/529/537 | ancienne version résoluble |
| Retraction Record | record versionné | consumers / Govern | retraction ≠ deletion |
| Product Supersession | version relation | Shared/consumers | historique et justification conservés |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-532/534 | error candidate identified | CAP-INV-536 | evidence, version, consumers, limits | source assessment |
| CAP-INV-536 | new version required | CAP-INV-526/527 | affected conclusions and correction rationale | Correction |
| CAP-INV-536 | withdrawal/retraction approved locally | CAP-INV-529 / Govern | affected publication, consumers, notifications, authority | Record |

Chaque transition conserve owner, tenant, versions, permissions, restrictions, erreurs, autorité, provenance et return origin.

## 18. Dépendances
CAP-INV-526..535/537 ; Shared Versioning/Notifications/Trace/Recovery ; Govern ; destination owners ; OPEN-013/018/019.

## 19. Source de vérité
Investigate est source des correction/retraction records ; Shared reste source du mécanisme de version ; chaque owner reste source de son retrait/rappel effectif.

## 20. Provenance et audit
Conserver issue, version, conclusions, sources, relations, Indicators/watchlists, consumers, analyses, décisions, notifications, auteurs, reviewers, timestamps, anciennes/nouvelles versions et return origin.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| correction/retraction create-review | reputational and consumer impact | 2 | restricted lineage masked | OPEN-013/019 | reporter/reviewer | Investigate | Permissions |
| external recall/runtime revocation | external/production impact | 3 | destination scope | mandatory | requester/approver/executor | Govern/owner | Future |

## 22. Limites et erreurs
Correction ≠ suppression ; retraction ≠ deletion ; superseded ≠ false ; notification ≠ external recall executed. Unknown consumers, stale links, rejected retraction et partial delivery restent visibles.

## 23. Métriques conceptuelles
Issues par état ; consumers notified ; silent corrections/history deletions — cible zéro ; external recalls executed locally — cible zéro.

## 24. Classification de livraison
`defined` / `planned` ; aucune révocation externe/runtime n’est revendiquée.

## 25. Critères d’acceptation
### 1. Correction
**Given** une conclusion erronée publiée  
**When** la correction est créée  
**Then** ancienne version, justification et consumers restent visibles.

### 2. Rétractation
**Given** une erreur majeure  
**When** la rétractation est proposée  
**Then** proposition, review et notification restent distinctes du rappel externe.

### 3. Sans IA
**Given** aucun modèle  
**When** l’impact est évalué  
**Then** diff, linking, tables, templates et revue humaine suffisent.

## 26. Questions ouvertes
OPEN-013/018/019 restent ouvertes ; aucune policy finale de recall ou revocation n’est imposée.

## 27. Consommateurs documentaires
Threat Intelligence, Investigate, Command, Detection Engineering, Settings, Govern, Shared, Permissions, Quality et Roadmap. Aucun Cloud/Mobile n’est lancé.
