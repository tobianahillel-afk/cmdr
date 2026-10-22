---
id: CAP-INV-535
title: Intelligence Requirement Satisfaction and Collection Feedback
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
# CAP-INV-535 — Intelligence Requirement Satisfaction and Collection Feedback

## 1. Définition
Évaluer la satisfaction d’une Intelligence Requirement à partir des produits, réponses, gaps, sources, confidence et contradictions, puis préparer un Collection Feedback Package ou un besoin dérivé sans déclencher automatiquement une collecte.

## 2. Problème utilisateur
Une Requirement peut être clôturée parce qu’un produit existe, même si les questions restent partielles ou incertaines.

## 3. Objectifs
- relier produits et critères de satisfaction.
- qualifier réponses, questions non résolues, gaps et limitations.
- rouvrir, dériver, clôturer ou superseder et préparer les feedbacks de collecte.

## 4. Non-objectifs
Aucune API, protocole, provider, connecteur, collecte réelle, scraping, tâche automatique, règle, blocage, réponse, partage externe, Cloud/Mobile ou écran détaillé.

## 5. Propriétaire
Investigate possède Requirement Satisfaction Assessment et Collection Feedback Package. CAP-INV-502 conserve Intelligence Requirement. Settings possède les demandes de sources et configurations. Studio possède Workflow/Run. Les propriétaires de collecte conservent toute exécution.

## 6. Utilisateurs
Principal : **Intelligence Manager**. Secondaires : Threat Intelligence Analyst, Collection Manager, Settings Admin, Studio Builder, Consumer Reviewer et Auditor autorisés.

## 7. Conditions d’entrée
Requirement/version, critères, produits liés, questions, sources, confidence, contradictions, owner, permissions et return origin sont explicites ; sinon not-assessed/blocked/partial.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Intelligence Requirement and criteria | CAP-INV-502 | need and satisfaction criteria | oui | current version | no assessment |
| Product and analysis outputs | CAP-INV-519..534 | answers, gaps and limitations | oui | selected versions | unsatisfied/partial |
| Source/collection context | CAP-INV-504 / Settings / Studio | availability and possible next steps | non | current projection | feedback limited |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Intelligence Requirement | Investigate | question, criteria, priority, owner | lecture/lien |
| Product / Assessment / Feedback | Investigate | answer, confidence, gaps | lecture/lien |
| Source / Workflow projection | Settings / Studio | available capability only | lecture limitée |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Requirement Satisfaction Assessment | créer, contester, réviser, superseder | Investigate concept | satisfied ≠ absolute truth |
| Collection Feedback Package | créer, versionner, retirer, superseder | Investigate concept | package ≠ collection execution |
| Derived Requirement proposal | préparer/lier | CAP-INV-502 owner | proposal ≠ active priority |

## 11. Fonctionnalités
Relier Requirement/produits ; examiner critères, réponses, unresolved questions, gaps, sources, confidence et contradictions ; classer satisfaction ; rouvrir/dériver/clôturer ; préparer demande source, Settings ou Studio.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| consulter/comparer | Intelligence Manager | requirement/outputs | 0 | lecture autorisée | vue sourcée | non |
| calculer couverture bornée | Intelligence Manager | assessment result | 1 | critères explicites | résultat explicable | selon politique |
| créer assessment/package/proposal | Intelligence Manager | concepts locaux | 2 | owner/provenance explicites | version réversible | OPEN-013 |
| lancer collecte/configurer source | aucun rôle local | collection/settings | 3 | owner autorisé | aucune exécution | obligatoire |
| supprimer historique | aucun rôle local | trace | 4 | interdit | refus audité | strict |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| comparer critères/réponses | oui | matrice explicable | oui | suggestion sourcée | checklist |
| détecter gaps/questions | oui | diff/règles | oui | résumé incertain | tables/revue |
| préparer feedback/package | oui | template versionné | oui | draft attribué | formulaire/workflow |
| conclure vérité ou lancer collecte | owner humain | contrôles | jamais autonome | jamais décisionnaire | revue et handoff |

Toute automatisation expose initiateur, moteur/version, sources, paramètres, erreurs, incertitude et disposition humaine.

## 14. États fonctionnels
`not-assessed`, `unsatisfied`, `partially-satisfied`, `satisfied-with-limitations`, `satisfied`, `obsolete`, `reopened`, `blocked`, `disputed`, `superseded`.

## 15. États d’interface
Loading conserve contexte ; Empty distingue aucune réponse/interdiction ; Partial nomme gaps ; Error conserve valide ; Offline stale/read-only ; Permission denied ne révèle rien ; Stale expose versions ; Conflict offre diff/recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Requirement Satisfaction Assessment | assessment versionné | Requirement owner / CAP-INV-537 | réponses, gaps et limites visibles |
| Collection Feedback Package | package | CAP-INV-502/Settings/Studio/collection owner | aucune collecte automatique |
| Derived Requirement proposal | proposal | CAP-INV-502 | besoin non activé implicitement |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-502/526/534 | assessment requested | CAP-INV-535 | criteria, answers, feedback, gaps | source object |
| CAP-INV-535 | collection gap confirmed | CAP-INV-502/Settings/Studio | gap, source need, restrictions, return origin | Assessment |
| CAP-INV-535 | lifecycle complete | CAP-INV-537 | status, limitations, packages, decisions | Requirement |

Chaque transition conserve owner, tenant, versions, permissions, restrictions, erreurs, autorité, provenance et return origin.

## 18. Dépendances
CAP-INV-502/504/519..537 ; Settings ; Studio ; collection owners ; Shared Trace/Versioning ; OPEN-013/018/019.

## 19. Source de vérité
Investigate est source de l’assessment/package ; CAP-INV-502 reste source du Requirement ; Settings/Studio/collection owner restent sources de leur exécution.

## 20. Provenance et audit
Conserver Requirement/version, critères, produits, réponses, gaps, sources, confidence, contradictions, feedback, reviewers, dispositions, packages, erreurs, timestamps et return origin.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| satisfaction assess/reopen/close | premature closure or priority impact | 2 | restricted sources masked | OPEN-013 | assessor/reviewer | Investigate | Permissions |
| collection feedback handoff | unauthorized collection | 2/3 | scope minimized | possible | requester/executor | Investigate/destination owner | Permissions/Future |

## 22. Limites et erreurs
Satisfied ≠ vérité absolue ; produit existant ≠ Requirement satisfaite ; collection gap ≠ collection task ; stale sources, disputed criteria, partial answers et blocked destination restent visibles.

## 23. Métriques conceptuelles
Requirements par état ; gaps/packages ; clôtures sans critères — cible zéro ; collectes déclenchées automatiquement — cible zéro.

## 24. Classification de livraison
`defined` / `planned` ; aucune collecte ou configuration n’est revendiquée.

## 25. Critères d’acceptation
### 1. Satisfaction limitée
**Given** une réponse utile mais des gaps critiques  
**When** l’assessment est produit  
**Then** état satisfied-with-limitations/partial et gaps restent visibles.

### 2. Feedback de collecte
**Given** une question non résolue  
**When** le package est préparé  
**Then** il conserve scope/restrictions sans lancer de collecte.

### 3. Sans IA
**Given** aucun modèle  
**When** la satisfaction est évaluée  
**Then** matrices, critères, tables et revue humaine suffisent.

## 26. Questions ouvertes
OPEN-013/018/019 restent ouvertes ; aucun provider, workflow ou seuil universel n’est choisi.

## 27. Consommateurs documentaires
Threat Intelligence, Settings, Studio, Collection, Shared, Permissions, Quality et Roadmap. Aucun Cloud/Mobile n’est lancé.
