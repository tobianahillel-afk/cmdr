---
id: CAP-INV-531
title: Indicator Operationalization and Detection or Command Handoff
product: investigate
module: threat-intelligence
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-06
requirement_ids: [REQ-PROD-014, REQ-PROD-019, REQ-PROD-020, REQ-PROD-055, REQ-INV-006, REQ-AI-002, REQ-SEC-001, REQ-SEC-002, REQ-UX-010]
open_decisions: [OPEN-013, OPEN-017, OPEN-018, OPEN-019]
source-of-truth: canonical
---
# CAP-INV-531 — Indicator Operationalization and Detection or Command Handoff

## 1. Définition
Préparer un Indicator Operationalization Package vers Detection Engineering, Command ou Settings avec sources, restrictions, confidence, contradictions, objectif, consommateur, expiration, retrait, monitoring et retour de statut, sans créer de règle, watchlist active, Signal, blocage ni modification runtime.

## 2. Problème utilisateur
Un package analytique peut être pris pour un Indicator déployé ou une Detection active.

## 3. Objectifs
Sélectionner Indicators autorisés ; définir objectif/destination/limites/expiry/retrait/monitoring ; préparer handoff ; suivre acceptation/refus/erreurs ; superseder ou retirer avant activation.

## 4. Non-objectifs
Aucune API, protocole, langage de règle, format, provider, runtime, activation, déploiement, Signal, Alert, blocage, réponse, partage externe, Cloud/Mobile ou écran détaillé.

## 5. Propriétaire
Investigate possède l’Operationalization Package. Detection Engineering possède Detection Content/Hypothesis/Coverage/Gap. Command possède Detection runtime, Signal, Alert et Incident. Settings possède targets/configurations. Govern possède l’autorité de production. Shared possède Linking/Trace/Versioning.

## 6. Utilisateurs
Principal : **Detection Engineer**. Secondaires : Threat Intelligence Analyst, SOC Analyst, Intelligence Manager, Incident Commander, Settings Admin, Reviewer et Auditor autorisés.

## 7. Conditions d’entrée
Indicators/version, sources, restrictions, confidence, contradictions, objectif, destination, expiry, owner, permissions, retrait et return origin sont explicites ; sinon package incomplete/blocked.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Authorized Indicators | CAP-INV-507 | selected versions and assessments | oui | current versions | no package |
| Watchlist/product/consumer context | CAP-INV-528/530 | objective, scope and restrictions | selon usage | current version | partial |
| Destination capabilities | Detection Engineering / Command / Settings | authorized consumer projection | oui | current projection | blocked |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Indicator / Sighting / Assessment | Investigate | sources, confidence, contradiction, expiry | lecture/lien |
| Detection Content/Hypothesis/Gap | Detection Engineering | consumer need projection | lecture/lien |
| Runtime/Signal/Incident/target status | Command / Settings | operational projection only | lecture limitée |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Indicator Operationalization Package | créer, revoir, versionner, retirer, superseder | Investigate concept | package ≠ deployed Indicator |
| Detection/Command/Settings handoff | préparer/lier | destination owner | aucune mutation silencieuse |
| Return disposition | enregistrer | Investigate / source owner | statut reçu, non inventé |

## 11. Fonctionnalités
Sélectionner Indicators ; définir objectif/destination/limites/expiry/retrait/monitoring ; préparer handoff ; recevoir statut, erreurs et refus ; retirer avant activation ; préserver retour/provenance.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| consulter et comparer | Detection Engineer | Indicators/destination projections | 0 | lecture autorisée | vue sourcée | non |
| produire comparaison/package preview | Detection Engineer | result/Tool Call | 1 | scope visible | résultat attribué | selon politique |
| créer/revoir/retirer package | Detection Engineer | Operationalization Package | 2 | owner/restrictions explicites | proposition réversible | OPEN-013/019 |
| déployer, activer ou modifier runtime | aucun rôle local | Detection/watchlist/runtime | 3 | owner/Govern requis | aucune exécution | obligatoire |
| supprimer historique | aucun rôle local | trace | 4 | interdit | refus audité | strict |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| préparer package | oui | formulaires/catalogues | oui | draft sourcé | checklist |
| comparer Indicators/consumers | oui | tables/diff | oui | suggestion incertaine | matrice |
| résumer limites/monitoring | oui | agrégations sourcées | oui | résumé attribué | timeline/revue |
| créer règle, Signal ou déployer | owner humain | contrôles | jamais autonome | jamais décisionnaire | Detection/Command/Govern |

Toute automatisation expose initiateur, moteur/version, Tool Calls, Run, sources, paramètres, erreurs, incertitude et disposition humaine.

## 14. États fonctionnels
`draft`, `review-ready`, `ready-for-handoff`, `handed-off`, `accepted`, `partially-accepted`, `rejected`, `withdrawn`, `superseded`, `expired`.

## 15. États d’interface
Loading conserve contexte ; Empty distingue aucune destination/interdiction ; Partial expose items refusés ; Error conserve succès valides ; Offline stale/read-only ; Permission denied ne révèle rien ; Stale expose versions ; Conflict offre diff/recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Indicator Operationalization Package | package versionné | Detection Engineering / Command / Settings | aucune activation implicite |
| Handoff disposition | event | CAP-INV-532/534/537 | statut, refus et erreurs visibles |
| Withdrawal/supersession context | lifecycle context | destination owner | retrait préparé, non exécuté silencieusement |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-507/530 | operational need identified | CAP-INV-531 | Indicators, sources, restrictions, objective | source context |
| CAP-INV-531 | handoff ready | Detection Engineering / Command / Settings | package, limits, expiry, monitoring, return origin | Package |
| destination owner | disposition available | CAP-INV-531/532 | accepted/refused items, errors, version, status | Handoff |

Chaque transition conserve owner, tenant, versions, permissions, restrictions, erreurs, autorité, provenance et return origin.

## 18. Dépendances
CAP-INV-507/513/517/528/530/532/534/537 ; CAP-INV-401..435 ; Command ; Settings ; Govern ; Shared ; OPEN-013/017/018/019.

## 19. Source de vérité
Investigate est source du Package ; les destinations sont sources de leur contenu/runtime/statut. Une acceptation ne signifie pas déploiement, match ou Signal.

## 20. Provenance et audit
Conserver Indicators, sources, assessments, objective, destination, expiry, monitoring, handoff, accepted/refused items, errors, withdrawals, versions, authors, timestamps et return origin.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| package create/review/handoff | unintended operational influence | 2 | restricted Indicators masked | OPEN-013/019 | author/reviewer | Investigate/destination | Permissions |
| deploy/activate/runtime modify | production impact | 3 | target scope | mandatory | requester/approver/executor | Govern/runtime owner | Future |

## 22. Limites et erreurs
Package ≠ règle, watchlist active, blocage ou deployed Indicator ; Indicator ≠ Detection Content ; aucun Signal créé ; rejected/partial/stale/version mismatch restent visibles.

## 23. Métriques conceptuelles
Packages par destination/disposition ; refused items ; déploiements ou Signals créés localement — cible zéro ; lineage perdu — cible zéro.

## 24. Classification de livraison
`defined` / `planned` ; aucune activation ou règle n’est revendiquée.

## 25. Critères d’acceptation
### 1. Handoff Detection
**Given** des Indicators avec contradictions  
**When** le package est transmis à Detection Engineering  
**Then** sources, contradictions et limits sont conservées et aucune règle n’est créée automatiquement.

### 2. Refus partiel
**Given** une destination refusant certains items  
**When** le retour est reçu  
**Then** chaque disposition et erreur reste traçable sans inventer un déploiement.

### 3. Sans IA
**Given** aucun modèle  
**When** le package est préparé  
**Then** formulaires, comparateurs, catalogues et revue humaine suffisent.

## 26. Questions ouvertes
OPEN-013/017/018/019 restent ouvertes ; runtime, langage et policy finale ne sont pas choisis.

## 27. Consommateurs documentaires
Threat Intelligence, Detection Engineering, Command, Settings, Govern, Shared, Permissions, Quality, Technique et Roadmap. Aucun Cloud/Mobile n’est lancé.
