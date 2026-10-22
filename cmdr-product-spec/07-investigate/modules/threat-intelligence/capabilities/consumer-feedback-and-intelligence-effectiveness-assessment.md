---
id: CAP-INV-534
title: Consumer Feedback and Intelligence Effectiveness Assessment
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
# CAP-INV-534 — Consumer Feedback and Intelligence Effectiveness Assessment

## 1. Définition
Recueillir et comparer le feedback des consommateurs, puis produire une Intelligence Effectiveness Assessment sourcée et bornée sur utilité, compréhension, pertinence, fraîcheur, actionnabilité, limites et usages déclarés.

## 2. Problème utilisateur
L’usage ou l’opinion d’un consommateur peut être traité comme valeur objective ou ground truth certaine.

## 3. Objectifs
- recueillir feedback, rôle, environnement et période.
- distinguer usage, utilité, limites et résultats déclarés.
- préparer correction, Requirement ou amélioration sans action automatique.

## 4. Non-objectifs
Aucune API, protocole, scoring opaque, enquête externe imposée, décision automatique, règle, blocage, réponse, partage externe, Cloud/Mobile ou écran détaillé.

## 5. Propriétaire
Investigate possède Feedback Record et Effectiveness Assessment. Command, Detection Engineering, Cases/Hunts et autres consommateurs conservent leurs objets et résultats. Shared possède metrics, notifications, collaboration, trace et versioning.

## 6. Utilisateurs
Principal : **Consumer Reviewer**. Secondaires : Intelligence Manager, Threat Intelligence Analyst, SOC Analyst, Detection Engineer, Investigation Lead, Incident Commander et Auditor autorisés.

## 7. Conditions d’entrée
Produit/publication/version, consommateur, période, rôle, contexte, permissions, restrictions et return origin sont explicites. Feedback anonyme, incomplet ou non vérifiable reste qualifié comme tel.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Product and publication version | CAP-INV-526/529 | evaluated content and delivery context | oui | selected version | no assessment |
| Consumer feedback | authorized consumer | perceived utility and limitations | oui | declared period | partial |
| Declared outcomes | Command / Cases / Hunts / Detection Engineering | usage and impact projections | non | source-owned period | value unknown |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Product / Publication / Access Record | Investigate | version, audience and restrictions | lecture/lien |
| Case/Hunt/Detection Gap/Decision projection | respective owner | declared supported outcome | lecture limitée |
| Metrics/Trace/Activity | Shared | usage and provenance context | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Consumer Feedback Record | créer, annoter, contester, versionner, superseder | Investigate concept | feedback ≠ ground truth |
| Intelligence Effectiveness Assessment | créer, réviser, retirer, superseder | Investigate concept | usage ≠ value |
| Improvement or correction proposal | préparer | destination owner | proposition ≠ action |

## 11. Fonctionnalités
Recueillir utilité, compréhension, pertinence, fraîcheur, actionnabilité, source quality, confidence, gaps, faux rapprochements, contradictions et usages déclarés ; comparer ; contester ; produire assessment ; préparer correction/Requirement.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| consulter et comparer | Consumer Reviewer | feedback/usage projections | 0 | lecture autorisée | vue sourcée | non |
| agréger ou analyser feedback | Consumer Reviewer | result/Tool Call | 1 | période et population visibles | résultat attribué | selon politique |
| créer feedback/assessment/proposal | Consumer Reviewer | concepts locaux | 2 | owner/provenance explicites | version réversible | OPEN-013/019 |
| modifier un objet consommateur ou agir | aucun rôle local | external object | 3 | owner requis | aucune exécution | obligatoire |
| supprimer feedback/historique | aucun rôle local | trace | 4 | interdit | refus audité | strict |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| collecter/classer feedback | oui | formulaires/catalogues | oui | suggestion sourcée | questionnaire structuré |
| comparer périodes/segments | oui | agrégations/diff | oui | résumé incertain | tables/metrics |
| préparer effectiveness summary | oui | template versionné | oui | draft attribué | matrice/revue |
| déclarer valeur ou vérité | humain autorisé | contrôles | jamais autonome | jamais décisionnaire | revue humaine |

Toute automatisation expose initiateur, moteur/version, sources, paramètres, erreurs, incertitude et disposition humaine.

## 14. États fonctionnels
`draft`, `collecting`, `partial`, `ready-for-assessment`, `under-review`, `supported`, `mixed`, `inconclusive`, `disputed`, `superseded`, `withdrawn`.

## 15. États d’interface
Loading conserve contexte ; Empty distingue aucun feedback/interdiction ; Partial expose segments manquants ; Error conserve valide ; Offline stale/read-only ; Permission denied ne révèle rien ; Stale expose période ; Conflict offre diff/recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Consumer Feedback Record | record versionné | CAP-INV-535/536/537 | rôle, période et limites visibles |
| Effectiveness Assessment | assessment | Intelligence Manager / Requirement owner | usage et valeur distingués |
| Improvement proposal | handoff | CAP-INV-535/536 | aucune action automatique |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-529 | product consumed | CAP-INV-534 | version, consumer, access, period | Publication |
| CAP-INV-534 | requirement impact identified | CAP-INV-535 | feedback, outcomes, limits, confidence | Assessment |
| CAP-INV-534 | error or misleading conclusion identified | CAP-INV-536 | affected version, evidence, consumers | Feedback |

Chaque transition conserve owner, tenant, versions, permissions, restrictions, erreurs, autorité, provenance et return origin.

## 18. Dépendances
CAP-INV-529/532/535..537 ; Command ; Cases/Hunts ; Detection Engineering ; Shared Metrics/Trace ; OPEN-013/019.

## 19. Source de vérité
Investigate est source du feedback et de l’assessment ; chaque consommateur reste source de ses résultats. Une déclaration d’usage ne remplace jamais la preuve correspondante.

## 20. Provenance et audit
Conserver produit/version, consommateur/rôle, période, questions, réponses, outcomes déclarés, sources, méthodes, erreurs, contradictions, reviewers, timestamps et disposition.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| feedback create/read/compare | personnel or customer disclosure | 0/2 | identity/context scoped | possible | contributor/reviewer | Investigate/consumer owner | Permissions |
| outcome verification | cross-product access | 0/1 | source-owner policy | possible | assessor/source owner | respective owner | Permissions |

## 22. Limites et erreurs
Usage élevé ≠ valeur élevée ; feedback ≠ ground truth ; absence de feedback ≠ inutilité ; biais, faible échantillon, période courte, stale data et self-report restent visibles.

## 23. Métriques conceptuelles
Feedback par rôle/période ; assessments avec limites ; score opaque de valeur — cible zéro ; correction automatique — cible zéro.

## 24. Classification de livraison
`defined` / `planned` ; aucune mesure opérationnelle implémentée n’est revendiquée.

## 25. Critères d’acceptation
### 1. Usage élevé, faible utilité
**Given** beaucoup d’accès mais une utilité déclarée faible  
**When** l’assessment est créé  
**Then** usage et valeur restent distincts avec limites visibles.

### 2. Feedback contradictoire
**Given** deux groupes avec avis opposés  
**When** ils sont comparés  
**Then** les deux populations et explications restent visibles sans score unique.

### 3. Sans IA
**Given** aucun modèle  
**When** l’efficacité est évaluée  
**Then** formulaires, tables, metrics et revue humaine suffisent.

## 26. Questions ouvertes
OPEN-013/018/019 restent ouvertes ; aucune métrique universelle ni policy finale n’est imposée.

## 27. Consommateurs documentaires
Threat Intelligence, Investigate, Command, Cases/Hunts, Detection Engineering, Shared, Permissions, Quality et Roadmap. Aucun Cloud/Mobile n’est lancé.
