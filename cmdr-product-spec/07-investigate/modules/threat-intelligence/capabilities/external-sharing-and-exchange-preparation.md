---
id: CAP-INV-533
title: External Sharing and Exchange Preparation
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
# CAP-INV-533 — External Sharing and Exchange Preparation

## 1. Définition
Préparer un package de partage ou d’échange externe avec version, sources, classification, markings, releasability, minimisation, destinataire candidat, finalité, risques, approbateur requis et return origin, sans transmettre aucune information.

## 2. Problème utilisateur
Une préparation peut être prise pour une autorisation ou une transmission effectivement réalisée.

## 3. Objectifs
Identifier objet/version, sources, restrictions, données à masquer et destinataire ; documenter finalité, durée, risques, reviewer/approbateur ; préparer Action Request et futur retrait/rappel sans exécuter le partage.

## 4. Non-objectifs
Aucune API, protocole, format, STIX/TAXII imposé, provider, connecteur, commande, code, transmission, publication client/publique, collecte active, réponse, Cloud/Mobile ou écran détaillé.

## 5. Propriétaire
Investigate possède la préparation analytique. Govern possède Decision/Approval et autorité de partage. Settings possède destinations/intégrations/configurations. Shared possède Export/Reporting/Trace. L’owner de destination conserve l’exécution réelle.

## 6. Utilisateurs
Principal : **Intelligence Manager**. Secondaires : Analyst, Consumer Reviewer, Approver, Security/Privacy/Compliance Reviewer et Auditor autorisés.

## 7. Conditions d’entrée
Objet/version, tenant, sources, classifications, markings, restrictions, releasability candidate, destinataire, finalité, permissions, owner, reviewer, approbateur et return origin sont explicites ; toute lacune bloque ou rend partial.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Product/version or knowledge package | CAP-INV-526..529 | selected share candidate | oui | current version | no package |
| Handling and releasability context | CAP-INV-528 / source owners | markings, restrictions, redaction | oui | current assessment | blocked |
| Destination and authority context | Settings / Govern | candidate destination/approver | oui | current projection | incomplete |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Intelligence Product / Publication Record | Investigate | selected version/context | lecture/lien |
| Source restrictions/classifications | source owners / Security | handling projection | lecture limitée |
| Destination / Decision / Approval | Settings / Govern | candidate configuration/status | lecture/lien |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| External Sharing Preparation Package | créer, annoter, contester, versionner, retirer, superseder | Investigate concept | package ≠ transmission |
| Action Request | préparer/lier | Govern owner | request ≠ Decision/Approval |
| Future recall context | préparer | destination/Govern owner | aucun rappel exécuté |

## 11. Fonctionnalités
Sélectionner objet/version ; conserver sources/restrictions/markings ; définir redaction, destinataire, finalité, durée, risques et approbateur ; préparer Action Request ; retirer/superseder sans transmission.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| consulter/filtrer/comparer | Intelligence Manager | source/product/authority projections | 0 | lecture autorisée | vue sourcée | non |
| produire preview/redaction assessment | Intelligence Manager | result/Tool Call | 1 | restrictions visibles | résultat attribué | selon politique |
| créer package/Action Request | Intelligence Manager | preparation package | 2 | owner/reviewer/approver explicites | proposition réversible | OPEN-013/019 |
| transmettre/publier réellement | aucun rôle local | destination externe | 3 | Govern + owner requis | aucune exécution | obligatoire |
| partage irrévocable/destruction trace | aucun rôle local | données/provenance | 4 | interdit | refus audité | strict |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| préparer inventaire/destination | oui | formulaires/catalogues | oui | draft sourcé | checklist |
| identifier données à masquer | oui | règles/diff | oui | suggestion attribuée | matrice redaction |
| préparer risques/Action Request | oui | template | oui | résumé incertain | workflow humain |
| approuver/transmettre | Govern/owner humain | contrôles | jamais autonome | jamais décisionnaire | Decision/Approval |

Toute automatisation expose initiateur, moteur/version, Tool Calls, Run, sources, paramètres, erreurs, incertitude et disposition humaine.

## 14. États fonctionnels
`draft`, `incomplete`, `under-review`, `redaction-required`, `approval-required`, `ready-for-action-request`, `submitted`, `returned`, `withdrawn`, `superseded`, `expired`.

## 15. États d’interface
Loading conserve contexte ; Empty distingue absence/interdiction ; Partial nomme données/approbations manquantes ; Error conserve valide ; Offline stale/read-only ; Permission denied ne révèle rien ; Stale expose dates ; Conflict offre diff/recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| External Sharing Preparation Package | package versionné | Govern / Settings / future owner | aucune transmission, restrictions visibles |
| Action Request preparation | request context | Govern | request ≠ Approval |
| Recall/withdrawal context | future context | Govern / destination owner | aucun rappel exécuté |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-528/529 | external need identified | CAP-INV-533 | product version, markings, redaction, purpose | Plan/Publication |
| CAP-INV-533 | review complete | Govern Action Request | package, risks, approver, conditions, return origin | Preparation |
| Govern/Settings | request returned/decided | CAP-INV-533 | status, reasons, limits, authority | Action Request |

Chaque transition conserve owner, tenant, versions, permissions, restrictions, erreurs, autorité, provenance et return origin.

## 18. Dépendances
CAP-INV-526..529/536/537 ; Govern ; Settings ; Shared Export/Reporting/Trace ; OPEN-013/018/019.

## 19. Source de vérité
Investigate est source de la préparation uniquement. Govern est source de l’autorité et le propriétaire de destination de l’exécution. Approval ≠ transmission terminée.

## 20. Provenance et audit
Conserver objet/version, sources, restrictions, classifications, markings, redactions, destinataire, finalité, risques, reviewers, approbateurs, Action Request, décisions, erreurs, timestamps et return origin.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| package create/review | external disclosure | 2 | minimisation/redaction | OPEN-013/019 | author/reviewer | Investigate/source owners | Permissions |
| actual external transmission | regulated disclosure | 3/4 | destination policy | mandatory | requester/approver/executor | Govern/Settings/owner | Future |

## 22. Limites et erreurs
Preparation ≠ sharing ; Approval ≠ transmission ; releasable ≠ shared ; aucun protocole/standard/provider imposé ; stale authority, destination unavailable, tenant mismatch et rejected request restent visibles.

## 23. Métriques conceptuelles
Packages par état/destination ; packages complets ; transmissions exécutées par Investigate — cible zéro ; perte de provenance/permission implicite — cible zéro.

## 24. Classification de livraison
`defined` / `planned` ; aucun partage, intégration ou déploiement revendiqué.

## 25. Critères d’acceptation
### 1. Restriction
**Given** une source interdisant le partage  
**When** la préparation est ouverte  
**Then** elle reste blocked/redaction-required et aucune transmission n’est faite.

### 2. Action Request
**Given** un package complet exigeant Govern  
**When** la demande est préparée  
**Then** approbateur, risques et return origin sont conservés sans Approval/envoi implicite.

### 3. Sans IA
**Given** aucun modèle  
**When** le package est préparé  
**Then** formulaires, catalogues, matrices de redaction et revue humaine suffisent.

## 26. Questions ouvertes
OPEN-013/018/019 restent ouvertes ; aucune option ni standard n’est sélectionné.

## 27. Consommateurs documentaires
Threat Intelligence, Investigate, Govern, Settings, Shared, Security/Permissions, Quality, Technique et Roadmap. Aucun partage réel, Cloud ou Mobile n’est lancé.
