---
id: CAP-INV-528
title: Handling, Releasability and Dissemination Planning
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
# CAP-INV-528 — Handling, Releasability and Dissemination Planning

## 1. Définition
Préparer les markings, restrictions, redactions, audiences internes, exclusions, règles de copie/export, expiration, révocation et le Dissemination Plan d’un produit Intelligence sans l’envoyer.

## 2. Problème utilisateur
Un marking ou une recommandation de release peut être interprété comme permission ou publication.

## 3. Objectifs
- préserver restrictions et définir données/annexes à masquer.
- définir audience, tenants, environnements, exclusions et releasability candidate.
- préparer le plan et une demande Govern lorsque nécessaire.

## 4. Non-objectifs
Aucune API, protocole, format, standard, provider, code, transmission externe, watchlist active, Indicator déployé, règle, blocage, réponse, Cloud/Mobile ou écran détaillé.

## 5. Propriétaire
Investigate possède Releasability Assessment et Dissemination Plan. Settings possède consumer groups, policies et destinations administrées. Govern possède Decision/Approval et partage gouverné. Shared possède Export/Reporting/Notifications/Trace. Les sources conservent leurs restrictions.

## 6. Utilisateurs
Principal : **Intelligence Manager**. Secondaires : Analyst, Consumer Reviewer, Approver, Security/Privacy Reviewer et Auditor autorisés.

## 7. Conditions d’entrée
Produit/version, sources, markings, restrictions, audience, tenants, environnements, owner, reviewers, permissions et return origin sont explicites ; sinon état partial, restricted ou blocked.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Product version and Release Recommendation | CAP-INV-526/527 | release candidate context | oui | selected version | blocked |
| Source markings and restrictions | source owners / CAP-INV-504 | handling constraints | oui | current versions | restricted |
| Consumer groups and policies | Platform Settings / Security | authorized audience projection | oui | current policy | no plan |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Intelligence Product Version | Investigate | reviewed content and lineage | lecture/lien |
| Source restrictions / classification | source owner / Security | handling constraints | lecture limitée |
| Consumer group / destination policy | Settings | authorized projection | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Releasability Assessment | créer, contester, versionner, superseder | Investigate concept | releasable ≠ shared |
| Dissemination Plan | créer, modifier, retirer, superseder | Investigate concept | plan ≠ publication |
| Govern request context | préparer | Govern owner | request ≠ Approval |

## 11. Fonctionnalités
Conserver markings/restrictions ; définir audience, exclusions, redaction, copy/export/share rules, expiry, revocation et reviewer ; préparer plan/Govern request ; préserver provenance et retour.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| consulter et comparer | Intelligence Manager | product/policies | 0 | lecture autorisée | vue sourcée | non |
| produire preview ou assessment borné | Intelligence Manager | redaction/releasability result | 1 | restrictions visibles | résultat attribué | selon politique |
| créer assessment/plan | Intelligence Manager | concepts locaux | 2 | owner/reviewer explicites | version réversible | OPEN-013/019 |
| partager extérieurement ou publier publiquement | aucun rôle local | destination externe | 3 | Govern requis | aucune exécution | obligatoire |
| détruire provenance | aucun rôle local | trace | 4 | interdit | refus audité | strict |

Investigate exécute uniquement les classes 0 à 2.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| préparer handling et redaction | oui | règles explicables | oui | suggestion sourcée | matrice/checklist |
| comparer audience et restrictions | oui | tables/diff | oui | résumé incertain | comparateur humain |
| préparer plan et demande | oui | template versionné | oui | draft attribué | formulaire/workflow |
| approuver ou partager | Govern humain | contrôles | jamais autonome | jamais décisionnaire | revue/Approval |

Toute sortie automatisée expose initiateur, moteur/version, Tool Calls, Run, sources, paramètres, erreurs, incertitude et disposition humaine.

## 14. États fonctionnels
`draft`, `eligibility-review`, `internally-releasable`, `restricted`, `approval-required`, `approved-for-internal-release`, `blocked`, `withdrawn`, `superseded`, `expired`.

## 15. États d’interface
Loading conserve contexte ; Empty distingue absence/interdiction ; Partial nomme les manques ; Error conserve le valide ; Offline stale/read-only ; Permission denied ne révèle rien ; Stale expose dates ; Conflict offre diff/recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Releasability Assessment | assessment versionné | CAP-INV-529/533 | restrictions et limites visibles |
| Dissemination Plan | plan versionné | CAP-INV-529/533/Govern | aucune transmission implicite |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-527 | release recommended | CAP-INV-528 | version, audience, markings, limits | Review |
| CAP-INV-528 | internal release eligible | CAP-INV-529 | audience, redaction, expiry, restrictions | Plan |
| CAP-INV-528 | external need | CAP-INV-533/Govern | candidate destination, risks, authority | Plan |

Chaque transition conserve owner, tenant, versions, permissions, restrictions, erreurs, autorité, provenance et return origin.

## 18. Dépendances
CAP-INV-504/525..527/529/533/536 ; Settings ; Govern ; Shared Export/Reporting ; OPEN-013/018/019.

## 19. Source de vérité
Investigate est source des assessments/plans ; les restrictions, policies, Decision et exécution restent chez leurs owners.

## 20. Provenance et audit
Conserver produit/version, sources, markings, restrictions, redactions, audiences, reviewers, décisions, erreurs, timestamps et retour. Correction par version/supersession, jamais effacement.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| releasability/plan create-review | disclosure | 2 | source data minimized | OPEN-013/019 | author/reviewer | Investigate/source owners | Permissions |
| external release | regulated disclosure | 3 | policy enforced | mandatory | requester/approver/executor | Govern/Settings | Future |

## 22. Limites et erreurs
Marking ≠ permission ; releasable ≠ partagé ; internal ≠ external ; product permission ≠ raw-source permission. Stale policy, tenant mismatch, rejected review et unavailable destination restent visibles.

## 23. Métriques conceptuelles
Plans par état ; plans avec restrictions/redactions complètes ; publication ou partage externe local — cible zéro ; permission implicite — cible zéro.

## 24. Classification de livraison
`defined` / `planned` ; preuve documentaire uniquement.

## 25. Critères d’acceptation
### 1. Restriction source
**Given** une source non releasable  
**When** le plan est préparé  
**Then** il reste restricted/blocked et aucune publication n’a lieu.

### 2. Audience interne
**Given** plusieurs groupes consommateurs  
**When** l’audience est définie  
**Then** tenants, exclusions et redactions restent explicites sans accès implicite aux sources.

### 3. Sans IA
**Given** aucun modèle  
**When** le plan est produit  
**Then** formulaires, matrices, règles et revue humaine suffisent.

## 26. Questions ouvertes
OPEN-013, OPEN-018 et OPEN-019 restent ouvertes ; aucune politique finale n’est sélectionnée.

## 27. Consommateurs documentaires
Threat Intelligence, Investigate, Settings, Govern, Shared, Permissions, Experience, Quality, Technique et Roadmap. Aucun Cloud/Mobile n’est lancé.
