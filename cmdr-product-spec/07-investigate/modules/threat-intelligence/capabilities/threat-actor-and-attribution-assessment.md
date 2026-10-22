---
id: CAP-INV-522
title: Threat Actor and Attribution Assessment
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
# CAP-INV-522 — Threat Actor and Attribution Assessment

## 1. Définition
Évaluer des Threat Entity Candidates et préparer une Attribution Assessment avec alternatives et incertitude, sans identité certaine ni attribution automatique.

## 2. Problème utilisateur
Des TTP, malware ou infrastructures partagés peuvent conduire à une attribution artificielle.

## 3. Objectifs
- examiner aliases, campaigns, infrastructure, malware, tools, TTP, targeting et temporalité.
- documenter alternatives, contradictions et niveaux d’incertitude.
- soumettre à revue, contester, retirer et superseder.

## 4. Non-objectifs
Aucune API, protocole, format d’échange, standard imposé, provider, schéma physique, moteur de scoring, algorithme d’attribution, scraper, commande, code produit, collecte active, action de réponse, partage externe réel, contenu Cloud Analysis, Mobile Forensics ou réécriture détaillée d’écran.

## 5. Propriétaire
Investigate possède le contexte analytique Threat Intelligence et les concepts fonctionnels locaux. Shared conserve Report, Entity, Graph, Timeline, Search, Linking, Versioning, Notifications, Collaboration, Trace, Activity, Export, Reporting et Recovery. Command conserve Detection runtime, Signal, Alert, Incident et dispositions opérationnelles. Detection Engineering conserve Detection Content, Hypothesis, Coverage, Gap et lifecycle. Settings conserve sources, providers, feeds, connecteurs, secrets, accès, destinations, stockage, rétention et health. Studio conserve Tool, Tool Call, Workflow, Automation Run, Automation Agent et Human Gate. Govern conserve Decision, Approval, Action Request, partage externe gouverné, Response Run et Result.

## 6. Utilisateurs
Principal : **Senior Intelligence Analyst**. Secondaires : Threat Intelligence Analyst, Senior Intelligence Analyst, Intelligence Manager, SOC Analyst, Detection Engineer, Incident Commander, Investigation Lead, Consumer Reviewer, Approver et Auditor autorisés selon le scope.

## 7. Conditions d’entrée
Tenant, environnement, période, owner, question, sources, versions, permissions, markings, restrictions, erreurs et return origin sont explicites. Toute absence produit un état incomplete, partial, blocked, restricted ou unknown ; aucune donnée n’est inventée.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Threat Entity Candidates | capability/source owner | versioned context | oui | selected version | état incomplete ou blocked |
| Fusion and competing assessments | linked owner | supporting context | selon scope | declared period | partialité visible |
| Restrictions et permissions | source owners / Security | markings, access and limits | oui | current decision | action refusée |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Threat Entity Candidate | source owner | authorized version and provenance | lecture/lien |
| Campaign/Infrastructure/Malware/TTP knowledge | source owner | context, alternatives and limits | lecture/comparaison |
| Source assessments | source owner | supporting projection | lecture limitée |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Threat Actor Assessment | créer, annoter, versionner, contester, superseder | Investigate concept | réversible, sourcé et non exécutoire |
| Attribution Assessment | créer ou mettre à jour | Investigate concept / destination owner | aucune mutation silencieuse |

## 11. Fonctionnalités
- examiner aliases, campaigns, infrastructure, malware, tools, TTP, targeting et temporalité.
- documenter alternatives, contradictions et niveaux d’incertitude.
- soumettre à revue, contester, retirer et superseder.
- conserver tenant, environnement, période, versions, sources, restrictions, erreurs, autorité et return origin.
- fonctionner sans IA ni fournisseur obligatoire.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| consulter, filtrer, rechercher, comparer | Senior Intelligence Analyst | sources et projections | 0 | lecture autorisée | vue sourcée et permission-aware | non |
| exécuter analyse, comparaison ou génération bornée | Senior Intelligence Analyst | résultat analytique / Tool Call | 1 | scope et restrictions visibles | résultat attribué, partialité et erreurs visibles | selon politique |
| créer, revoir, publier internement de façon réversible ou préparer un handoff | Senior Intelligence Analyst | concept Threat Actor and Attribution Assessment | 2 | owner, provenance et séparation explicites | version ou proposition non effective hors scope | OPEN-013/019 |
| partager extérieurement, activer, déployer ou modifier un runtime | aucun rôle local | objet externe/production | 3 | hors périmètre local | aucune exécution locale | obligatoire |
| détruire provenance ou historique | aucun rôle local | source/historique | 4 | interdit | refus audité | strict |

Investigate exécute uniquement les classes 0 à 2. Les classes 3 et 4 sont bloquées ou routées vers Govern et le propriétaire runtime.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| préparer Threat Actor and Attribution Assessment | oui | formulaires, matrices et règles explicables | oui | draft sourcé | formulaire structuré et checklist |
| comparer sources, versions ou alternatives | oui | diff, tables et agrégations | oui | assistance avec incertitude | comparateur et revue humaine |
| résumer gaps, contradictions et limites | oui | vues sourcées | oui | résumé attribué | timeline, matrice et Inspector |
| confirmer, approuver, partager extérieurement ou exécuter | humain autorisé / futur owner | contrôles seulement | jamais autonome | jamais décisionnaire | revue humaine et Govern |

Toute sortie automatisée expose initiateur, moteur ou agent et version, Tool Calls, Automation Run, sources, paramètres, timestamp, statut, erreurs, incertitude, owner humain et disposition humaine. Aucune fonction essentielle ne dépend d’un modèle et aucun chatbot n’est obligatoire.

## 14. États fonctionnels
`draft`, `under-analysis`, `weakly-supported`, `supported`, `ambiguous`, `contradicted`, `inconclusive`, `disputed`, `withdrawn`, `superseded`. Ces états sont des projections fonctionnelles versionnées, pas une machine d’état canonique finale.

## 15. États d’interface
Loading conserve contexte et source ; Empty distingue absence, interdiction et non-collecte ; Partial nomme les éléments manquants ; Error conserve les résultats valides ; Offline est stale/read-only ; Permission denied ne révèle rien ; Stale expose dates et consommateurs ; Conflict offre diff, versions et recovery.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Attribution Assessment | versioned concept | downstream capability / owner | sources, versions, restrictions et incertitude visibles |
| Actor assessment review event | business event or package | reviewer / consumer | aucun effet externe implicite |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-508 | actor assessment requested | CAP-INV-522 | candidate, aliases, relations and restrictions | Candidate |
| CAP-INV-521 | fusion ready | CAP-INV-522 | corroboration, dependencies and contradictions | Fusion |
| CAP-INV-522 | product planning requested | CAP-INV-525 | assessment, alternatives, confidence and limits | Assessment |

Chaque transition conserve ownership, tenant, environnement, versions, sources, markings, permissions, restrictions, erreurs, autorité, provenance et return origin ; elle n’étend jamais les droits.

## 18. Dépendances
CAP-INV-508/511/512/514/515/520/521; Shared Entity/Graph; OPEN-013/018/019. Shared est consommé sans redéfinition. Aucun protocole, moteur, provider ou format externe n’est choisi.

## 19. Source de vérité
Investigate est source du contexte analytique et des concepts Threat Intelligence locaux. Chaque objet canonique reste chez son owner. Une projection, synthèse, score, recommandation, publication interne ou handoff ne remplace jamais sa source, ne transfère pas l’ownership et n’étend pas les permissions.

## 20. Provenance et audit
Conserver Intake, Requirement, Knowledge Project, Analysis Handoff Package, Session, question, hypothèses, sources, materials, candidates, Sightings, relations, assessments, produits, versions, reviews, publications, consommateurs, Tools, Tool Calls, Automation Runs, décisions humaines, erreurs, restrictions, timestamps et return origin. Toute correction se fait par nouvelle version, retrait, rétractation ou supersession ; aucune trace n’est supprimée.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| Attribution Assessment create/update/review | reputation and false attribution | 2 | identities and sources masked | OPEN-013 | author/reviewer | Investigate / source owner | Permissions |
| provenance read/export | cross-product disclosure | 0/1 | restricted lineage masked | possible | viewer/auditor | Shared/source owners | Permissions |

Les namespaces, permissions atomiques, RBAC/ABAC, step-up définitif et séparation finale restent reportés. La permission sur un produit ou Project ne remplace jamais celle de la source.

## 22. Limites et erreurs
- assessment ≠ identity certainty.
- shared infrastructure/TTP/malware ≠ shared actor.
- high confidence ≠ absolute fact.
- Stale, partial, restricted, tenant mismatch, timeout, cancellation, source unavailable, consumer unavailable et version superseded restent visibles.
- Une sortie IA, un score, un nombre de sources, un edge, un usage ou un feedback ne constitue jamais seul une vérité, une Approval ou une action.

## 23. Métriques conceptuelles
- actor assessments by disposition.
- alternatives per assessment.
- automatic attribution — target zero.
- sorties automatisées avec sources, paramètres, erreurs, incertitude et disposition humaine.
- permission auto-accordée, contradiction masquée ou trace supprimée — cible zéro.

Aucun seuil universel ni score opaque n’est imposé.

## 24. Classification de livraison
`defined` / `planned` ; preuve documentaire uniquement. Aucun statut `validated`, `implemented`, `native`, `integrated`, `deployed`, `active` ou `operational` n’est revendiqué.

## 25. Critères d’acceptation
### 1. Contexte incomplet
**Given** un contexte de Threat Actor and Attribution Assessment avec sources ou permissions manquantes  
**When** Senior Intelligence Analyst tente de poursuivre  
**Then** l’état reste incomplete/partial/blocked, les lacunes sont visibles, aucune vérité ou action n’est inventée et le return origin est conservé.

### 2. Alternative ou restriction
**Given** des éléments contradictoires, dépendants ou restreints  
**When** le concept est comparé ou revu  
**Then** alternatives, sources, restrictions et versions restent visibles, aucune fusion ou permission silencieuse n’a lieu.

### 3. Sans IA
**Given** aucun modèle disponible  
**When** le workflow est exécuté  
**Then** formulaires, matrices, tableaux, diff, timelines, recherche, notifications déterministes et revue humaine suffisent.

## 26. Questions ouvertes
OPEN-013 reste ouverte pour les actions de classe 2. OPEN-018 reste ouverte pour l’ontologie et l’interopérabilité. OPEN-019 reste ouverte pour dissemination, releasability, partage et accès consommateurs. Aucun choix n’est imposé. Les objets finaux, permissions atomiques, écrans, contrats techniques et implémentation restent futurs.

## 27. Consommateurs documentaires
Threat Intelligence Foundations et Analysis/Products, Investigate, Cases/Hunts/Evidence, Analysis Workbench, Detection Engineering, Command, Platform Settings, CMDR Studio, Govern, Shared, Objects, Permissions, Experience Architecture, Screens, Journeys, Quality, Technique et Roadmap. Le document ne lance ni Cloud Analysis ni Mobile Forensics.
