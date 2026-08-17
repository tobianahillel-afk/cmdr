---
id: CAP-INV-505
title: Source Reliability and Information Credibility Assessment
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
source-of-truth: canonical
---
# CAP-INV-505 — Source Reliability and Information Credibility Assessment

## 1. Définition
Évaluer séparément la fiabilité historique d’une source, la crédibilité d’une information, la corroboration, les contradictions et la confiance analytique, sans imposer d’échelle universelle ni transformer une claim en fait.

## 2. Problème utilisateur
Confondre source fiable, information crédible et vérité analytique masque les contradictions et crée une confiance injustifiée.

## 3. Objectifs
- créer des assessments sourcés et révisables.
- distinguer reliability, credibility, confidence, freshness et corroboration.
- exposer éléments favorables, défavorables, contradictions et limites.
- annoter, contester, revoir et superseder sans score opaque.

## 4. Non-objectifs
Aucune API, protocole, format d’échange, standard imposé, provider imposé, schéma physique, modèle de graphe, moteur de scoring, scraper, commande, code, collecte active, attribution automatique, Indicator déployé, watchlist active, règle Detection, blocage, réponse, partage externe, contenu 4B.3B.2, Cloud/Mobile Analysis ou réécriture détaillée d’écran.

## 5. Propriétaire
Investigate possède le contexte analytique et **Source Reliability and Information Credibility Assessment** comme concept fonctionnel. Shared conserve Entity, Graph, Timeline, Search, Object Linking, Versioning, Jobs, Notifications, Trace, Activity, Export, Reporting, Collaboration et Recovery. Command conserve Detection, Signal, Alert et Incident. Detection Engineering conserve Detection Content et son lifecycle. Settings conserve sources, providers, connectors, secrets, rétention, accès et health. Studio conserve Tool, Tool Call, Workflow, Automation Run et Human Gate. Govern conserve Decision, Approval, partage externe futur, Response Run et Result.

## 6. Utilisateurs
Principal : **Intelligence Reviewer**. Secondaires : Threat Intelligence Analyst, Intelligence Manager, Investigation Lead, Detection Engineer, SOC Analyst, Reviewer et Auditor autorisés selon le scope.

## 7. Conditions d’entrée
Tenant, environnement, période, source, versions, permissions, restrictions, handling markings, objectifs, owner et return origin sont explicites. Une absence produit un état incomplete, partial, blocked, restricted ou unknown ; elle n’est jamais remplacée par une donnée inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Source context | CAP-INV-504 / Settings | owner, access, history and limitations | oui | current source version | not-assessed |
| Information/material | CAP-INV-506 | claim, provenance and extraction context | oui | immutable/reference version | credibility unknown |
| Corroborating/contradicting sources | CAP-INV-503/514 | independent candidate evidence | non | declared periods | uncorroborated |
| Method and reviewer | human/process policy | functional method, author and scope | oui | assessment time | assessment invalid |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Intelligence Source projection | Settings/CAP-INV-504 | history and access context | lecture |
| Intelligence Material | Investigate/Artifact owner | claims and provenance | lecture |
| Sighting / Relationship / Contradiction | Investigate concepts | supporting and opposing context | lecture |
| Source Activity/Audit | Shared/Settings | quality history and changes | lecture |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Source Reliability Assessment | créer, revoir, contester, superseder | Investigate concept | reliability ≠ credibility |
| Information Credibility Assessment | créer, revoir, contester, superseder | Investigate concept | credibility ≠ correctness |
| Assessment rationale | annoter/versionner | Investigate | method and limits visible |

## 11. Fonctionnalités
- créer des assessments sourcés et révisables.
- distinguer reliability, credibility, confidence, freshness et corroboration.
- exposer éléments favorables, défavorables, contradictions et limites.
- annoter, contester, revoir et superseder sans score opaque.
- conserver tenant, environnement, versions, sources, restrictions, erreurs, attribution et return origin.
- fonctionner sans fournisseur de modèle ni chatbot obligatoire.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| consulter, filtrer, rechercher, comparer | Intelligence Reviewer | sources et projections de Source Reliability and Information Credibility Assessment | 0 | lecture autorisée | vue sourcée et permission-aware | non |
| exécuter extraction, normalisation, assessment ou comparaison bornée | Intelligence Reviewer | résultat analytique / Tool Call | 1 | lancement explicite, scope et restrictions visibles | résultat attribué, partialité et erreurs visibles | selon politique |
| créer, annoter, contester, versionner, superseder ou préparer un handoff | Intelligence Reviewer | concept fonctionnel Source Reliability and Information Credibility Assessment | 2 | mutation réversible, owner et provenance explicites | nouvelle version ou proposition non effective | OPEN-013 |
| publier, partager, déployer, bloquer ou modifier une source administrative | aucun rôle local | objet externe ou production | 3 | hors périmètre ; future Decision/Approval | aucune exécution locale | obligatoire |
| supprimer irréversiblement ou détruire la provenance | aucun rôle local | connaissance/historique | 4 | interdit par défaut | refus audité | strict |

Investigate exécute uniquement les classes 0 à 2. Les classes 3 et 4 sont bloquées ou routées vers le futur owner/Govern ; aucune action réelle n’est réalisée dans cette phase.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| préparer ou compléter Source Reliability and Information Credibility Assessment | oui | formulaires, catalogues et règles explicables | oui | proposition sourcée | formulaire structuré et checklist |
| extraire, comparer ou détecter des lacunes | oui | parsers, diff et comparateurs déterministes | oui | assistance avec incertitude | tables, filtres, recherche et revue humaine |
| résumer sources, contradictions et limites | oui | agrégations sourcées | oui | résumé attribué | timeline, matrice et Inspector |
| confirmer, attribuer, fusionner ou publier | humain autorisé / future phase | contrôles seulement | non autonome | jamais décisionnaire | revue humaine et Govern lorsque requis |

Toute sortie automatisée expose initiateur, agent ou moteur et version, Automation Run, Tool Calls, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et acceptation, modification ou rejet. Aucun chatbot n’est obligatoire et aucune fonction essentielle ne dépend d’un modèle.

## 14. États fonctionnels
`not-assessed`, `assessing`, `reliable-candidate`, `mixed`, `unreliable-candidate`, `credibility-unknown`, `corroborated`, `contradicted`, `stale`, `disputed`, `superseded`. Ces états sont fonctionnels et versionnés ; ils ne constituent pas un schéma ou une machine d’état canonique finale.

## 15. États d’interface
Loading conserve le contexte et la source ; Empty distingue absence, interdiction et non-collecte ; Partial nomme les éléments manquants ; Error conserve les résultats valides ; Offline est stale/read-only ; Permission denied ne révèle aucune donnée protégée ; Stale conserve dates et consommateurs ; Conflict offre diff, versions et recovery. Le Design System possède le rendu.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Reliability assessment | assessment | CAP-INV-506..518 | source history and limits distinct |
| Credibility assessment | assessment | candidate capabilities | claim-specific support/contradiction visible |
| Dispute/review event | business event | reviewer/project | no automatic invalidation |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-504 | source selected | CAP-INV-505 | source history, owner, restrictions and freshness | Catalog |
| CAP-INV-506 | material parsed | CAP-INV-505 | claim, source, extraction and ambiguities | Material |
| CAP-INV-514/515 | corroboration or contradiction changes | CAP-INV-505 | relations, sources and assessment deltas | Assessment |
| CAP-INV-505 | assessment available | CAP-INV-507..518 | reliability, credibility, rationale and limits | Assessment |

Chaque transition conserve l’owner source et destination, tenant, environnement, versions, source, markings, permissions, restrictions, erreurs, autorité, provenance et return origin. Une transition n’étend jamais implicitement les droits.

## 18. Dépendances
CAP-INV-504/506/514/515; Shared Activity/Comparison; Security; OPEN-013/018. Les Shared Capabilities sont consommées sans redéfinition. `OPEN-018` couvre l’ontologie, la portabilité et l’interopérabilité futures sans sélectionner de standard ou protocole.

## 19. Source de vérité
Investigate est source du contexte Threat Intelligence, des assessments et candidates locaux. Chaque objet canonique reste chez son owner. Une projection, extraction, relation, score, suggestion ou handoff ne remplace jamais sa source et ne transfère ni ownership ni permission.

## 20. Provenance et audit
Conserver Requirement, Knowledge Project, Case/Hunt/Incident/Detection/analysis origin, sources, access context, materials, Artifacts, extractions, Tools, Tool Calls, Automation Runs, candidates, relations, Sightings, assessments, contradictions, versions, supersessions, expirations, revocations, auteurs, reviewers, timestamps, paramètres, erreurs, restrictions, décisions humaines et return origin. Toute correction se fait par version ou supersession ; aucune trace n’est supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| Reliability assessment create/update | source reputation impact | 2 | source identity may be masked | OPEN-013 | assessor/reviewer | Investigate | Permissions |
| Credibility assessment review | knowledge classification impact | 0/2 | restricted material scoped | possible | analyst/reviewer | Investigate | Permissions |
| Method/rationale export | sensitive source exposure | 1 | redaction and markings | step-up possible | reviewer/auditor | Shared/Investigate | Permissions |

Les namespaces, permissions atomiques, RBAC/ABAC, step-up définitif et séparation finale des tâches restent reportés. La permission d’un Project ne remplace jamais celle de la source ou de la destination.

## 22. Limites et erreurs
- Source reliability ≠ information credibility ≠ correctness ≠ confidence.
- A reliable source may be wrong; an unreliable source may be correct.
- Number of sources does not prove independence.
- No universal scale or opaque score is imposed.
- Les états sont des projections fonctionnelles, pas une machine d’état objet définitive.
- Stale, partial, restricted, tenant mismatch, timeout, cancellation, source unavailable et version superseded restent visibles.
- Une sortie IA, un nombre de sources, un edge, un score ou une enrichment ne constitue jamais seul une vérité, une attribution, une Approval ou une action.

## 23. Métriques conceptuelles
- assessments with rationale/method/limits.
- disputed and superseded assessments.
- claims with corroborating and contradicting sources.
- opaque automatic truth scores — target zero.
- sorties automatisées avec initiateur, version, sources, paramètres, erreurs, incertitude et disposition humaine.
- permission auto-accordée, contradiction masquée et trace supprimée — cible zéro.

Aucun seuil universel, score opaque ou objectif quantitatif non approuvé n’est imposé.

## 24. Classification de livraison
`defined` / `planned` ; preuve documentaire uniquement. Aucun statut `validated`, `implemented`, `native`, `integrated`, `deployed`, `active` ou `operational` n’est revendiqué. Promotion conditionnée par les phases Objets, Permissions, Écrans, Technique et décisions ouvertes.

## 25. Critères d’acceptation
### 1. Source fiable, information contradictoire
**Given** une source historiquement fiable et plusieurs observations contraires  
**When** la credibility est évaluée  
**Then** reliability reste visible, credibility peut être contested et aucun fait n’est confirmé automatiquement.

### 2. Contexte absent
**Given** une information sans source résoluble  
**When** l’assessment est ouvert  
**Then** credibility reste unknown et la lacune est visible.

### 3. Sans IA
**Given** aucun modèle  
**When** l’assessment est réalisé  
**Then** matrices explicables, comparaisons et revue humaine suffisent.

## 26. Questions ouvertes
- OPEN-013 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-018 reste ouverte ; aucune option n’est sélectionnée.
- OPEN-009 reste la seule décision historiquement résolue.
- Ontologie finale, schémas, identifiants, cardinalités, modèles de graph, taxonomies, formats d’échange, permissions atomiques, contrats techniques, écrans détaillés et 4B.3B.2 restent futurs.

## 27. Consommateurs documentaires
Threat Intelligence Foundations, Investigate, Cases/Hunts/Evidence, Analysis Workbench, Detection Engineering, Command projections, Platform Settings, Studio, Govern, Shared, Objects, Permissions, Experience Architecture, Screens, Journeys, Quality, Technique et future 4B.3B.2. Le document ne lance ni Cloud Analysis ni Mobile Forensics.
