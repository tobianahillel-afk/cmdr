---
id: CAP-INV-529
title: Internal Intelligence Publication and Consumer Access
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
# CAP-INV-529 — Internal Intelligence Publication and Consumer Access

## 1. Définition
Publier réversiblement une version Intelligence en interne vers des audiences autorisées et tracer accès, erreurs, suspension, retrait et supersession sans partage externe ni permission implicite sur les sources brutes.

## 2. Problème utilisateur
Une publication interne peut exposer trop de données ou être confondue avec une publication publique.

## 3. Objectifs
Vérifier quality/release/restrictions ; créer Internal Publication Record ; notifier les audiences autorisées ; suspendre, retirer et superseder avec historique.

## 4. Non-objectifs
Aucune API, protocole, format, provider, publication externe/publique/client, watchlist active, Indicator déployé, règle, blocage, réponse, Cloud/Mobile ou écran détaillé.

## 5. Propriétaire
Investigate possède Internal Publication Record et Consumer Access Record. Settings possède groupes/policies. Shared possède Notifications, Search, Reporting, Versioning, Trace et Recovery. Govern possède les releases gouvernées. Les sources conservent leurs permissions.

## 6. Utilisateurs
Principal : **Intelligence Manager**. Secondaires : Analyst, Consumer Reviewer, SOC Analyst, Detection Engineer, Incident Commander et Auditor autorisés.

## 7. Conditions d’entrée
Version quality-passed, Release Recommendation, Dissemination Plan, audience, restrictions, expiry, permissions, owner et return origin sont explicites ; sinon blocked/partial.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Quality-passed Product Version | CAP-INV-526/527 | reviewed release candidate | oui | selected version | blocked |
| Dissemination Plan | CAP-INV-528 | audience, redaction and expiry | oui | current plan | no publication |
| Consumer access policies | Settings / Security | authorized groups | oui | current policy | denied |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Product Version / Review | Investigate | content and disposition | lecture/lien |
| Consumer group / access policy | Settings | authorized projection | lecture |
| Notification / Trace / Version | Shared | delivery and history | lecture/lien |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Internal Publication Record | créer, suspendre, retirer, versionner, superseder | Investigate concept | internal ≠ external |
| Consumer Access Record | créer, révoquer, versionner | Investigate / Settings policy | accès produit ≠ accès source brute |

## 11. Fonctionnalités
Publier seulement vers audiences autorisées ; voir versions, consommateurs, accès et erreurs ; notifier ; suspendre, retirer, superseder ; préserver historique.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| consulter publications et accès | Intelligence Manager | records/projections | 0 | lecture autorisée | vue sourcée | non |
| préparer preview/notification | Intelligence Manager | publication result | 1 | plan valide | résultat attribué | selon politique |
| publier internement, suspendre, retirer | Intelligence Manager | publication/access records | 2 | policy/audience/expiry explicites | action réversible auditée | OPEN-013/019 |
| partager extérieurement/publiquement | aucun rôle local | externe | 3 | Govern requis | aucune exécution | obligatoire |
| supprimer historique | aucun rôle local | trace | 4 | interdit | refus audité | strict |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| vérifier éligibilité/audience | oui | policies/règles | oui | suggestion sourcée | checklist |
| générer preview/notification | oui | template versionné | oui | résumé attribué | modèle documentaire |
| tracer accès/erreurs | oui | événements déterministes | oui | synthèse possible | tables/timeline |
| décider accès externe | Govern humain | contrôles | jamais autonome | jamais décisionnaire | Approval humaine |

Toute automatisation expose initiateur, moteur/version, sources, paramètres, erreurs et disposition humaine.

## 14. États fonctionnels
`draft`, `eligibility-review`, `approved-for-internal-release`, `internally-published`, `partial`, `suspended`, `withdrawn`, `superseded`, `expired`, `failed`.

## 15. États d’interface
Loading conserve contexte ; Empty distingue aucune publication/interdiction ; Partial expose consommateurs non servis ; Error conserve succès valides ; Offline stale/read-only ; Permission denied ne révèle rien ; Stale expose version ; Conflict offre diff/recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Internal Publication Record | record versionné | authorized consumers / CAP-INV-534 | audience/restrictions visibles |
| Consumer Access Record | access event | Settings/Shared/Auditor | aucune permission source implicite |
| Notification event | notification | authorized consumers | notification ≠ Alert |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-528 | internal release eligible | CAP-INV-529 | version, audience, redaction, expiry | Plan |
| CAP-INV-529 | consumer uses product | CAP-INV-534 | publication, consumer, period, context | Publication |
| CAP-INV-536 | correction/retraction | CAP-INV-529 | affected version and notification | Correction |

Chaque transition conserve owner, tenant, versions, permissions, restrictions, erreurs, autorité, provenance et return origin.

## 18. Dépendances
CAP-INV-526..528/532/534/536/537 ; Settings ; Shared ; Govern ; OPEN-013/019.

## 19. Source de vérité
Investigate est source du record ; Settings reste source des policies/groupes ; les sources restent propriétaires de leurs données.

## 20. Provenance et audit
Conserver version, review, plan, audience, redactions, accès, notifications, erreurs, suspension, retrait, supersession, auteurs et timestamps. Aucune trace supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| internal publish/suspend/withdraw | internal disclosure | 2 | redacted product only | OPEN-013/019 | author/reviewer/publisher | Investigate/Settings | Permissions |
| raw-source access | sensitive disclosure | 0/3 | source policy | possible | consumer/source owner | source owner | Existing/Future |

## 22. Limites et erreurs
Internal ≠ external/public ; consumer access ≠ raw-source access ; recommendation ≠ publication ; partial delivery, denied consumer, stale policy et superseded version restent visibles.

## 23. Métriques conceptuelles
Publications par état/audience ; accès autorisés/refusés ; permissions implicites — cible zéro ; externe/public local — cible zéro.

## 24. Classification de livraison
`defined` / `planned` ; preuve documentaire uniquement.

## 25. Critères d’acceptation
### 1. Consommateur non autorisé
**Given** un groupe hors audience  
**When** la publication est exécutée  
**Then** l’accès est refusé sans révéler contenu ni sources.

### 2. Retrait
**Given** une version publiée corrigée  
**When** elle est retirée/supersedée  
**Then** historique, consommateurs et notifications restent traçables.

### 3. Sans IA
**Given** aucun modèle  
**When** la publication interne est préparée  
**Then** policies, checklists, templates et notifications déterministes suffisent.

## 26. Questions ouvertes
OPEN-013/018/019 restent ouvertes ; aucune policy finale ni partage externe n’est imposé.

## 27. Consommateurs documentaires
Threat Intelligence, Command, Detection Engineering, Settings, Shared, Govern, Permissions, Experience, Quality et Roadmap. Aucun Cloud/Mobile n’est lancé.
