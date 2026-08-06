---
id: CAP-INV-530
title: Intelligence Watchlist Definition and Lifecycle
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
# CAP-INV-530 — Intelligence Watchlist Definition and Lifecycle

## 1. Définition
Définir fonctionnellement une Watchlist, son objectif, ses Indicators autorisés, sources, owner, consommateurs, scope, environnements, durée, conditions, restrictions, priorités, handoffs, versions, statut et révocation sans l’activer.

## 2. Problème utilisateur
Une définition analytique peut être prise pour une watchlist active, et un match pour un Signal ou une activité malveillante certaine.

## 3. Objectifs
Définir contenu/finalité/lifecycle ; préparer activation handoff ; recevoir une projection runtime ; expirer, révoquer, superseder et retirer sans mutation directe du runtime.

## 4. Non-objectifs
Aucune API, protocole, format, provider, activation, moteur de matching, Signal, Alert, règle, blocage, réponse, partage externe, Cloud/Mobile ou écran détaillé.

## 5. Propriétaire
Investigate possède Watchlist Definition et lifecycle assessment. Settings/runtime owner possède configuration/activation. Command possède Signal/Alert. Detection Engineering possède Detection Content. Govern possède l’autorité. Shared possède Trace/Versioning/Notifications.

## 6. Utilisateurs
Principal : **Threat Intelligence Analyst**. Secondaires : Intelligence Manager, Detection Engineer, SOC Analyst, Settings Admin, Reviewer et Auditor autorisés.

## 7. Conditions d’entrée
Objectif, Indicators, sources, confidence, restrictions, consommateurs, scope, environnements, durée, owner, permissions et return origin sont explicites ; sinon draft/incomplete/blocked.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Indicator Candidates or authorized Indicators | CAP-INV-507 | selected versioned items | oui | current versions | no definition |
| Audience, scope and handling | CAP-INV-528 / Settings | consumers and restrictions | oui | current policy | blocked |
| Runtime capability projection | Settings / runtime owner | possible targets and health | non | current projection | inactive definition |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Indicator / Sighting / lifecycle | Investigate | confidence, contradiction, recency | lecture/lien |
| Consumer/environment projection | Settings | authorized scope | lecture |
| Runtime status | runtime owner | observed activation/health only | lecture limitée |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Watchlist Definition | créer, revoir, versionner, expirer, révoquer, superseder, retirer | Investigate concept | definition ≠ active watchlist |
| Activation handoff | préparer | Settings/Govern/runtime owner | aucune activation locale |
| Lifecycle Assessment | créer, contester, réviser | Investigate concept | status observed, not invented |

## 11. Fonctionnalités
Définir objectif, contenu, portée, durée, conditions et priorités ; préparer handoff ; voir statut/erreurs ; retirer, expirer, révoquer, superseder ; préserver provenance.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| consulter et comparer | Analyst | definition/runtime projections | 0 | lecture autorisée | vue sourcée | non |
| produire comparaison/lifecycle assessment | Analyst | result/Tool Call | 1 | scope visible | résultat attribué | selon politique |
| créer/revoir definition et handoff | Analyst | Watchlist Definition | 2 | owner/restrictions explicites | version réversible | OPEN-013/019 |
| activer/désactiver runtime | aucun rôle local | active watchlist | 3 | owner/Govern requis | aucune exécution | obligatoire |
| supprimer historique | aucun rôle local | trace | 4 | interdit | refus audité | strict |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| préparer Watchlist Definition | oui | formulaires/règles | oui | draft sourcé | checklist |
| vérifier expiry et versions | oui | dates/diff | oui | résumé attribué | timeline |
| préparer handoff | oui | template | oui | suggestion incertaine | workflow humain |
| activer ou conclure malveillance | owner humain | contrôles | jamais autonome | jamais décisionnaire | Settings/Govern/Command |

Toute automatisation expose initiateur, moteur/version, sources, paramètres, erreurs et disposition humaine.

## 14. États fonctionnels
`draft`, `under-review`, `approved-for-handoff`, `inactive`, `activation-requested`, `active-observed`, `degraded`, `expired`, `revoked`, `superseded`, `retired`.

## 15. États d’interface
Loading conserve contexte ; Empty distingue aucune definition/interdiction ; Partial expose items manquants ; Error conserve le valide ; Offline stale/read-only ; Permission denied ne révèle rien ; Stale expose runtime date ; Conflict offre diff/recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Watchlist Definition | concept versionné | Settings/Govern/runtime owner | aucune activation implicite |
| Activation handoff | package | owner destination | sources, restrictions, expiry visibles |
| Lifecycle Assessment | assessment | CAP-INV-532/536/537 | statut reçu, non inventé |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-507 | watchlist need identified | CAP-INV-530 | Indicators, sources, confidence, expiry | Candidate |
| CAP-INV-530 | approved for handoff | Settings/Govern/runtime owner | definition, scope, restrictions, rollback | Definition |
| runtime owner | status observed | CAP-INV-530/532 | version, per-target status, errors, time | Runtime |

Chaque transition conserve owner, tenant, versions, permissions, restrictions, erreurs, autorité, provenance et return origin.

## 18. Dépendances
CAP-INV-507/517/528/531/532/536/537 ; Settings ; Govern ; Command ; Shared ; OPEN-013/018/019.

## 19. Source de vérité
Investigate est source de la Definition ; runtime owner est source de l’activation/santé ; Command reste source des Signals/Alerts.

## 20. Provenance et audit
Conserver definition, items, sources, scope, consumers, expiry, reviews, handoffs, runtime projections, erreurs, révocations, versions, auteurs et timestamps.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| definition create/review | unsupported monitoring | 2 | restricted Indicators masked | OPEN-013/019 | author/reviewer | Investigate | Permissions |
| activation/runtime change | production impact | 3 | target scope | mandatory | requester/approver/executor | Settings/Govern/runtime | Future |

## 22. Limites et erreurs
Definition ≠ active watchlist ; match ≠ Signal ni malveillance ; `active-observed` vient du runtime owner ; stale, degraded, rejected activation et version mismatch restent visibles.

## 23. Métriques conceptuelles
Definitions par état ; items expirés ; activations exécutées par Investigate — cible zéro ; matches convertis automatiquement en Signal — cible zéro.

## 24. Classification de livraison
`defined` / `planned` ; aucune watchlist active n’est revendiquée.

## 25. Critères d’acceptation
### 1. Activation demandée
**Given** une Definition approved-for-handoff  
**When** elle est transmise  
**Then** elle reste non active jusqu’au statut du runtime owner.

### 2. Match observé
**Given** un match runtime  
**When** il est consulté  
**Then** il n’est ni Signal ni preuve de malveillance et conserve contexte/source.

### 3. Sans IA
**Given** aucun modèle  
**When** la Definition est gérée  
**Then** formulaires, dates, diff, catalogues et revue humaine suffisent.

## 26. Questions ouvertes
OPEN-013/018/019 restent ouvertes ; activation, moteur et policy finale restent futurs.

## 27. Consommateurs documentaires
Threat Intelligence, Detection Engineering, Command, Settings, Govern, Shared, Permissions, Quality et Roadmap. Aucun Cloud/Mobile n’est lancé.
