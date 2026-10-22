---
id: CAP-INV-213
title: Collection Integrity and Custody
product: investigate
module: collection-and-live-response
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-014
  - REQ-PROD-020
  - REQ-OBJ-004
open_decisions:
  - OPEN-014
source-of-truth: canonical
---
# CAP-INV-213 — Collection Integrity and Custody

## 1. Définition
Décrire fonctionnellement origine, autorisation, transformations, transferts, accès, copies, vérification et custody sans choisir algorithme ou stockage.

## 2. Problème utilisateur
Sans Collection Integrity and Custody, l’utilisateur perd le lien entre le Case, l’Endpoint, l’autorité applicable, l’exécution locale et les résultats. Les états partiels ou offline peuvent alors être pris pour un succès et les objets peuvent être confondus.

## 3. Objectifs
- fournir demande, autorisation, acquisition, transformations, transferts, rétention;
- exposer cible, scope, fraîcheur, policy, permission et classe d’action;
- conserver erreurs, résultats partiels, provenance et retour au Case;
- produire custody/provenance events, verification status sans transférer l’ownership.

## 4. Non-objectifs
- ne pas administrer la Fleet ni les Endpoint Policies;
- ne pas définir protocole, API, commande, moteur, format, PKI, stockage ou plateforme supportée;
- ne pas créer automatiquement Evidence, Finding, Decision, Response Run ou Govern Result;
- ne pas commencer Analysis Workbench.

## 5. Propriétaire
Investigate / Collection and Live Response / Investigate Product Lead possède le contexte métier et les relations au Case. Platform Settings administre Fleet/Policies; Endpoint Agent exécute localement; Govern possède l’autorité risquée.

## 6. Utilisateurs
Principal : Evidence Reviewer / DFIR Analyst. Secondaires : Investigation Lead, Incident Commander, approbateur Govern ou Platform Administrator en consultation.

## 7. Conditions d’entrée
Tenant et environnement conservés, Case accessible, Endpoint résolu, demande et autorisation identifiées, acquisition ou résultat disponible, permissions de custody/review et politiques de rétention projetées.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Case et objectif | Investigate | contexte métier | oui | version courante | rester incomplete |
| Request/job/opération et autorisation | Investigate / Govern | origine et autorité | oui | versions visibles | custody non vérifiable |
| Endpoint/Agent et acquisition | Settings / Endpoint Agent | cible et méthode déclarée | oui | timestamps et statut visibles | partial/disputed |
| Artifact et transformations | Investigate | objet et dérivations | oui pour custody Artifact | versions courantes | lien manquant explicite |
| Accès, transfert, copie, rétention/export | owners respectifs | événements de custody | selon usage | ordonnés et sourcés | lacune visible |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case / Collection Request | Investigate | contexte, scope et autorisation | consulter |
| Collection Job / Endpoint Operation | concepts à formaliser | exécution et résultats | consulter |
| Endpoint / Endpoint Agent | partagé / Endpoint Agent | cible et méthode déclarée | consulter |
| Artifact / Evidence | Investigate | origine, versions, qualifications et relations | consulter/revoir |
| Action Request / Decision / Response Run / Result | Govern | autorité éventuelle | consulter |
| Export / retention / legal hold projections | Shared / Settings / Trust | restrictions | consulter |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Custody event/record conceptuel | créer, compléter, contester ou supersede | Investigate/Trust, modèle final futur | aucune suppression silencieuse ; auteur et timestamp obligatoires |
| Artifact provenance relation | créer/supersede | Investigate | acquisition, transformation et versions conservées |
| Verification/dispute status | enregistrer | Investigate | distinct de pertinence et Evidence qualification |
| Artifact/Evidence source | aucune réécriture | Investigate | références versionnées seulement |

## 11. Fonctionnalités
Revoir origine, cible, initiateur, autorisation, méthode déclarée, timestamps, transformations, transferts, accès, copies, dérivés, rétention, export et legal hold futur ; vérifier, contester et naviguer vers sources.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter la chaîne | reviewer | custody/provenance | 0 | read | chaîne sourcée | non |
| Ajouter une annotation/correction | reviewer autorisé | custody event | 2 | raison et version | nouvel événement, ancien conservé | OPEN-013 |
| Vérifier ou contester | reviewer | verification status | 2 | sources disponibles | statut attribué | non par défaut |
| Préparer export | utilisateur autorisé | export request | 1/2 | policy/rétention | handoff Export Engine | selon sensibilité |
| Relier à Evidence | Evidence Reviewer | relation | 2 | qualification explicite | CAP-INV-107/108 | non automatique |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Assembler les événements | oui | relations/IDs/timestamps | oui | non nécessaire | timeline déterministe |
| Signaler un maillon manquant | oui | règles | oui | explication | checklist de custody |
| Comparer copies/dérivés | oui | métadonnées déclarées | oui | résumé | inspection manuelle |
| Proposer contestation | oui | règles | oui | suggestion | revue humaine |
| Qualifier Evidence | humain explicite | jamais automatique | non | non | CAP-INV-107/108 |

Toute sortie automatisée expose initiateur, moteur/version, Automation Run/Tool Calls, sources, paramètres, timestamp, statut, incertitude, owner humain, disposition et trace.

## 14. États fonctionnels
`initiated`, `recording`, `partial`, `verified`, `disputed`, `superseded`, `retention-pending`, `export-restricted`. Ce ne sont pas des états finaux d’objet.

## 15. États d’interface
Loading conserve la chaîne ; Empty indique l’événement source manquant ; Partial marque les maillons ; Error garde les événements valides ; Offline autorise lecture du dernier snapshot ; Permission denied masque contenus sensibles ; Stale montre versions superseded.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Custody/provenance chain | événements/record conceptuel | Evidence Review et audit | ordre, acteurs, sources et lacunes visibles |
| Verification/dispute status | état fonctionnel | Case et Evidence | attribué, versionné et distinct de pertinence |
| Export/retention context | package | Export/Settings/Trust | restrictions conservées |
| Artifact source relations | relations | CAP-INV-105/107/112/114 | aucun changement silencieux de version |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Collection Job/Operation Result | sortie matérielle | CAP-INV-213 | request, target, acquisition, timestamps, transformations, Artifact | source |
| CAP-INV-213 | qualification requise | CAP-INV-107/108 | Artifact/source, custody, verification, contestations | custody view |
| CAP-INV-213 | préparer export | Export Engine | Artifact/version, classification, retention, restrictions | custody view |
| Evidence/Report | ouvrir source | CAP-INV-213 | relation, version, return origin | Evidence/Report |

## 18. Dépendances
CAP-INV-105/107/108/112/114/203/212/214, Endpoint Agent custody source, Shared Timeline/Linking/Export, retention/legal-hold owners et OPEN-014.

## 19. Source de vérité
Investigate conserve les événements métier et relations Artifact/Case ; Endpoint Agent reste source des événements locaux déclarés ; Shared/Trust possèdent les mécanismes d’audit/export ; Evidence qualification reste Investigate.

## 20. Provenance et audit
Request, autorisation, Case, Endpoint/Agent, méthode déclarée, timestamps, transformations, transferts, accès, copies, dérivés, reviewers, contestations, retention/export et correlation IDs.

## 21. Permissions fonctionnelles
Custody read/review/dispute, Artifact read/export, sensitive output, retention/legal hold projection, cross-tenant restrictions et separation of duties. Matrice atomique reportée.

## 22. Limites et erreurs
Source absente, timestamps incohérents, copie non liée, transformation inconnue, permission révoquée, Artifact superseded, export restreint ou legal hold inconnu rendent la chaîne partial/disputed sans effacer les données valides.

## 23. Métriques
Artifacts avec chaîne complète, maillons manquants, contestations, corrections, exports bloqués, relations cassées et navigations source réussies.

## 24. Classification de livraison
`defined` / `planned`. Aucun algorithme, manifeste, PKI, protocole ou stockage n’est choisi ; promotion dépend des phases Objets, Trust, Technique et OPEN-014.

## 25. Critères d’acceptation
**Given** un Artifact collecté **When** sa custody est ouverte **Then** request, autorisation, cible, acquisition, transferts, transformations et accès sont navigables ou explicitement manquants.

**Given** une correction **When** le reviewer la soumet **Then** l’entrée antérieure reste visible et la nouvelle version est attribuée.

**Given** aucun modèle IA **When** la chaîne est revue **Then** timeline, règles et checklist permettent la vérification complète.

## 26. Questions ouvertes
OPEN-014 reste ouverte. Le record de custody, intégrité technique, legal hold, rétention et permissions finales appartiennent aux phases Objets/Trust/Technique.

## 27. Consommateurs documentaires
Evidence Creation/Review, Artifact Management, Case Replay, Reporting, Export, Trust, Endpoint Agent custody source, phases Objets/Permissions/Technique.