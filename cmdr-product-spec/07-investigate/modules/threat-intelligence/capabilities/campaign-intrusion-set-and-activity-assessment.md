---
id: CAP-INV-523
title: Campaign, Intrusion Set and Activity Assessment
product: investigate
module: threat-intelligence
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-06
requirement_ids:
  - REQ-PROD-014
  - REQ-PROD-019
  - REQ-PROD-020
  - REQ-PROD-055
  - REQ-INV-006
  - REQ-AI-002
  - REQ-SEC-001
  - REQ-SEC-002
  - REQ-UX-010
open_decisions:
  - OPEN-013
  - OPEN-018
  - OPEN-019
source-of-truth: canonical
---
# CAP-INV-523 — Campaign, Intrusion Set and Activity Assessment

## 1. Définition
Évaluer des Activity Cluster, Campaign et Intrusion Set Candidates avec périodes, Sightings, TTP, malware, infrastructure, targets autorisés, contradictions et regroupements alternatifs.

## 2. Problème utilisateur
Des co-occurrences peuvent être fusionnées en opération ou attribution certaine.

## 3. Objectifs
- produire un résultat versionné, sourcé, contestable et borné.
- préserver alternatives, contradictions, restrictions et return origin.
- permettre revue, retrait, correction et supersession sans effacement.

## 4. Non-objectifs
Aucune API, protocole, format d’échange, standard imposé, provider, schéma physique, moteur, algorithme opaque, scraper, commande, code produit, collecte active, déploiement, blocage, réponse, partage externe réel, contenu Cloud Analysis, Mobile Forensics ou réécriture détaillée d’écran.

## 5. Propriétaire
Investigate possède ce contexte analytique et ses concepts fonctionnels. Shared conserve Report, Reporting, Entity, Graph, Timeline, Search, Linking, Versioning, Notifications, Collaboration, Trace, Activity, Export et Recovery. Command conserve Detection runtime, Signal, Alert, Incident et dispositions. Detection Engineering conserve Detection Content, Hypothesis, Coverage, Gap et lifecycle. Settings conserve sources, providers, feeds, connecteurs, secrets, accès, destinations, stockage, rétention et health. Studio conserve Tool, Tool Call, Workflow, Automation Run, Automation Agent et Human Gate. Govern conserve Decision, Approval, Action Request, partage externe gouverné, Response Run et Result.

## 6. Utilisateurs
Principal : **Threat Intelligence Analyst**. Secondaires : Threat Intelligence Analyst, Senior Intelligence Analyst, Intelligence Manager, SOC Analyst, Detection Engineer, Incident Commander, Investigation Lead, Consumer Reviewer, Approver et Auditor autorisés.

## 7. Conditions d’entrée
Tenant, environnement, période, owner, question, sources, versions, permissions, markings, restrictions, erreurs et return origin sont explicites. Toute absence produit `incomplete`, `partial`, `blocked`, `restricted` ou `unknown`.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Activity/Campaign/Intrusion Set Candidates | CAP-INV-511 | candidate groupings | oui | selected/current version | partialité ou blocage visible |
| Sightings, TTP and linked knowledge | CAP-INV-508..514 | temporal and relationship context | selon scope | selected/current version | partialité ou blocage visible |
| Permissions, markings and restrictions | source owners / Security | access and handling | oui | current decision | action refusée |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Campaign candidates | Investigate | versions, periods and gaps | lecture/lien limité |
| Sightings/TTP/Malware/Infrastructure | Investigate | support and alternatives | lecture/lien limité |
| Trace and versions | Shared | lineage and comparison | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Campaign Assessment | créer, annoter, contester, versionner, retirer, superseder | Investigate concept | réversible et sourcé |
| Activity / Intrusion Set Assessment | créer ou mettre à jour | Investigate concept / destination owner | aucune mutation silencieuse |

## 11. Fonctionnalités
- Évaluer des Activity Cluster, Campaign et Intrusion Set Candidates avec périodes, Sightings, TTP, malware, infrastructure, targets autorisés, contradictions et regroupements alternatifs.
- exposer sources, hypothèses, limites, erreurs et autorité.
- conserver versions, consommateurs et contexte de retour.
- fonctionner intégralement sans IA.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| consulter, filtrer, rechercher, comparer | Threat Intelligence Analyst | sources et projections | 0 | lecture autorisée | vue sourcée | non |
| exécuter analyse/comparaison/génération bornée | Threat Intelligence Analyst | résultat / Tool Call | 1 | scope visible | résultat attribué avec erreurs | selon politique |
| créer, revoir, versionner ou préparer un handoff | Threat Intelligence Analyst | Campaign Assessment | 2 | owner et provenance explicites | version/proposition réversible | OPEN-013/019 |
| partager extérieurement, activer ou modifier un runtime | aucun rôle local | objet externe | 3 | hors périmètre | aucune exécution | obligatoire |
| détruire provenance ou historique | aucun rôle local | historique | 4 | interdit | refus audité | strict |

Investigate exécute seulement les classes 0 à 2 ; les classes 3 et 4 sont bloquées ou routées vers Govern et le propriétaire runtime.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| préparer Campaign, Intrusion Set and Activity Assessment | oui | formulaires, matrices, règles explicables | oui | draft sourcé | formulaire et checklist |
| comparer sources, versions et alternatives | oui | diff, tables, agrégations | oui | assistance incertaine | comparateur humain |
| résumer gaps, contradictions et limites | oui | vues sourcées | oui | résumé attribué | timeline, matrice, Inspector |
| confirmer, approuver, partager ou exécuter | humain autorisé | contrôles | jamais autonome | jamais décisionnaire | revue humaine et Govern |

Toute sortie automatisée expose initiateur, moteur/agent/version, Tool Calls, Automation Run, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et acceptation/modification/rejet. Aucun chatbot ni modèle n’est obligatoire.

## 14. États fonctionnels
`draft`, `under-analysis`, `supported`, `weakly-supported`, `ambiguous`, `contradicted`, `inconclusive`, `disputed`, `withdrawn`, `superseded`. États fonctionnels versionnés, non machine d’état canonique finale.

## 15. États d’interface
Loading conserve le contexte ; Empty distingue absence/interdiction/non-collecte ; Partial nomme les manques ; Error conserve le valide ; Offline est stale/read-only ; Permission denied ne révèle rien ; Stale expose dates/consommateurs ; Conflict offre diff et recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Campaign Assessment | concept versionné | capability/owner suivant | sources, restrictions et incertitude visibles |
| Activity / Intrusion Set Assessment | événement ou package | reviewer/consumer | aucun effet externe implicite |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-511 | assessment requested | CAP-INV-523 | candidate groups, periods, sources and gaps | Candidate |
| CAP-INV-521 | fusion available | CAP-INV-523 | claims, dependencies and alternatives | Fusion |
| CAP-INV-523 | product planning requested | CAP-INV-525 | assessment, timeline, limits and implications | Assessment |

Chaque transition conserve ownership, tenant, environnement, versions, sources, markings, permissions, restrictions, erreurs, autorité, provenance et return origin ; aucun droit n’est étendu.

## 18. Dépendances
Threat Intelligence CAP-INV-501..537 selon les transitions ; owners Command, Detection Engineering, Settings, Studio, Govern et Shared ; OPEN-013/018/019 et, lorsqu’indiquées, OPEN-014/015. Aucun protocole, moteur, provider ou format externe n’est choisi.

## 19. Source de vérité
Investigate est source du contexte analytique local. Chaque objet canonique reste chez son owner. Projection, synthèse, recommandation, publication interne ou handoff ne remplace pas la source, ne transfère pas l’ownership et n’étend pas les permissions.

## 20. Provenance et audit
Conserver Intake, Requirement, Project, Handoff, Session, question, hypothèses, sources, materials, candidates, Sightings, relations, assessments, produits, versions, reviews, publications, consommateurs, Tools, Tool Calls, Automation Runs, décisions humaines, erreurs, restrictions, timestamps et return origin. Correction par version, retrait, rétractation ou supersession ; aucune trace supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| Campaign assessment create/update/review | false grouping and attribution | 2 | sources/identités/données restreintes masquées | OPEN-013 | auteur/reviewer | Investigate / source owner | Permissions |
| provenance read/export | cross-product disclosure | 0/1 | lineage restreinte | possible | viewer/auditor | Shared/source owners | Permissions |

Namespaces, permissions atomiques, RBAC/ABAC, step-up et séparation finale restent futurs. Permission produit ≠ permission source.

## 22. Limites et erreurs
- Activity Cluster ≠ Campaign ≠ Intrusion Set.
- Co-occurrence ≠ same operation.
- Campaign Assessment ≠ Incident or actor attribution.
- Stale, partial, restricted, tenant mismatch, timeout, cancellation, source/consumer unavailable et version superseded restent visibles.
- Sortie IA, score, nombre de sources, edge, usage ou feedback ne constitue jamais seul vérité, Approval ou action.

## 23. Métriques conceptuelles
- objets par état, owner, version et consumer.
- objets avec sources, alternatives, restrictions et disposition humaine.
- permissions silencieuses, historique perdu ou effets externes implicites — cible zéro.
Aucun seuil universel ni score opaque.

## 24. Classification de livraison
`defined` / `planned` ; preuve documentaire uniquement. Aucun statut `validated`, `implemented`, `native`, `integrated`, `deployed`, `active` ou `operational`.

## 25. Critères d’acceptation
### 1. Contexte incomplet
**Given** sources ou permissions manquantes  
**When** Threat Intelligence Analyst poursuit  
**Then** état incomplete/partial/blocked, lacunes visibles, aucune vérité/action inventée, return origin conservé.

### 2. Alternative ou restriction
**Given** éléments contradictoires, dépendants ou restreints  
**When** le concept est comparé ou revu  
**Then** alternatives, sources, restrictions et versions restent visibles, sans fusion ni permission silencieuse.

### 3. Sans IA
**Given** aucun modèle  
**When** le workflow est exécuté  
**Then** formulaires, matrices, tableaux, diff, timelines, recherche, notifications déterministes et revue humaine suffisent.

## 26. Questions ouvertes
OPEN-013 reste ouverte pour la classe 2 ; OPEN-018 pour ontologie/interoperability ; OPEN-019 pour dissemination, releasability, sharing et consumer access. Aucun choix imposé. Objets finaux, permissions atomiques, écrans, contrats techniques et implémentation restent futurs.

## 27. Consommateurs documentaires
Threat Intelligence Foundations et Analysis/Products, Investigate, Cases/Hunts/Evidence, Analysis Workbench, Detection Engineering, Command, Platform Settings, CMDR Studio, Govern, Shared, Objects, Permissions, Experience Architecture, Screens, Journeys, Quality, Technique et Roadmap. Aucun Cloud Analysis ni Mobile Forensics.
