---
id: CAP-INV-397
title: Network Forensics Handoff to Evidence, Findings, Detection Engineering and Intelligence
product: investigate
module: analysis-workbench
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-05
requirement_ids:
  - REQ-INV-001
  - REQ-PROD-014
  - REQ-PROD-020
  - REQ-OBJ-004
  - REQ-AI-002
  - REQ-SEC-001
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
source-of-truth: canonical
---
# CAP-INV-397 — Network Forensics Handoff to Evidence, Findings, Detection Engineering and Intelligence

## 1. Définition
Sélectionner des packets, flows, sessions, conversations, transactions, observations DNS, métadonnées chiffrées, certificats, Entities, relations, anomalies, Derived Artifacts et contradictions afin de préparer des packages non qualifiés pour Evidence, Findings et futures phases Detection Engineering et Intelligence.

## 2. Problème utilisateur
Sans handoff explicite et sourcé, un paquet, domaine, certificat, anomalie ou objet extrait peut être qualifié prématurément comme Evidence, Finding, IOC, règle ou objet Intelligence.

## 3. Objectifs
- sélectionner les observations et objets pertinents avec leurs limites ;
- expliquer leur relation à la Hypothesis ;
- conserver les éléments favorables, contradictoires et la qualité de couverture ;
- préparer une Evidence candidate vers CAP-INV-107/108 ;
- préparer un Finding Draft vers CAP-INV-109 ;
- préparer un futur package Detection Engineering sans créer de règle ;
- préparer un futur package Intelligence sans créer d’objet canonique ;
- préserver provenance et retour au Workbench.

## 4. Non-objectifs
Ne pas qualifier automatiquement Evidence ou Finding, confirmer un IOC, créer/tester/déployer une règle Detection, créer un Indicator, Campaign ou Threat Actor, interagir avec une cible, rejouer du trafic, contourner un chiffrement, utiliser un secret, choisir un moteur, définir une API, une commande, un protocole interne ou un écran détaillé.

## 5. Propriétaire
Investigate / Analysis Workbench / Investigate Product Lead possède Evidence Candidate Package, Finding Draft et Future Handoff Package jusqu’à acceptation par la capability destination. CAP-INV-107/108/109 conserve la qualification Evidence/Finding. La future Phase 4B.3 possédera Detection Engineering et Intelligence. Studio, Settings, Shared, Command et Govern conservent leurs objets.

## 6. Utilisateurs
Principal : **Investigation Lead**. Secondaires : Network Forensics Analyst, Evidence Reviewer, Finding Owner, futurs Detection Engineer et Intelligence Analyst.

## 7. Conditions d’entrée
Case et Hypothesis accessibles ; observations sélectionnées ; capture, coverage, timebase, restrictions, provenance et reproductibilité visibles ; permission de préparer le handoff ; aucun état de destination créé silencieusement.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Observations et sélections Network | CAP-INV-380..396 | records, relations, anomalies et Artifacts | oui | versions de session | rester draft ou `partial` |
| Hypothesis et contexte Case | Investigate | question, scope et return origin | oui | état courant | bloquer le handoff |
| Coverage, timebase et contradictions | CAP-INV-383/393/394 | qualité et limites | oui | snapshot lié | `partial` et revue obligatoire |
| Provenance et reproductibilité | CAP-INV-396 | sources, Tools, paramètres et lacunes | oui | assessment courant | bloquer ou marquer non reproductible |
| Permissions et restrictions d’export | Security / Settings | minimisation et autorité | oui | réévaluées | masquer ou refuser |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case / Hypothesis | Investigate | objectif, reasoning, restrictions et return origin | consulter/lier |
| Packet, Flow, Session, Conversation et Transaction observations | Investigate concepts | sélection, source, qualité et limites | consulter/sélectionner |
| DNS, Certificate, Encrypted Traffic, Entity et Network Relationship observations | Investigate / Shared projections | valeurs, sources, confiance et contradictions | consulter/sélectionner |
| Network Anomaly / Timeline / Comparison / Derived Artifact | Investigate | disposition, time quality, différences et lineage | consulter/sélectionner |
| Reproducibility Assessment / Trace | Investigate / Shared | provenance et lacunes | consulter |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Evidence Candidate Package | créer, modifier, retirer ou soumettre | Investigate | ne devient Evidence qu’après CAP-INV-107/108 |
| Finding Draft | créer, modifier, retirer ou soumettre | Investigate | ne devient Finding qu’après CAP-INV-109 |
| Future Detection Engineering Package | préparer ou retirer | Investigate jusqu’au handoff | aucune règle créée, testée ou déployée |
| Future Intelligence Package | préparer ou retirer | Investigate jusqu’au handoff | aucun objet Intelligence canonique créé |
| Trace / Activity event | émettre | Shared | append-only avec destination et disposition |

## 11. Fonctionnalités
- sélectionner packets, flows, sessions, conversations et transactions ;
- sélectionner DNS, encrypted metadata, certificats, Entities et relations ;
- sélectionner anomalies, Derived Artifacts, contradictions et limites ;
- expliquer le lien à Hypothesis ;
- préparer et revoir les quatre types de package ;
- conserver source, coverage, timebase, Tool context et reproductibilité ;
- soumettre au propriétaire destination et revenir au Workbench ;
- fonctionner sans fournisseur de modèle.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Sélectionner et inspecter | Investigation Lead | observations et Artifacts | 0 | lecture autorisée | sélection sourcée | non |
| Produire un résumé ou package borné | Investigation Lead | candidate package | 1 | sources et paramètres explicites | brouillon attribué | non |
| Modifier, annoter ou retirer | Investigation Lead | candidate package | 2 | permission réversible | nouvelle version ou retrait tracé | OPEN-013 |
| Soumettre vers Evidence/Finding | Investigation Lead | package | 2 | critères et provenance présents | destination reçoit un candidat | owner destination |
| Préparer futur Detection/Intelligence handoff | Investigation Lead | future package | 2 | limites et phase future explicites | package uniquement | future owner |

Les classes 3 et 4 sont exclues. Aucune règle, action cible ou qualification canonique n’est exécutée.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Grouper les sélections | oui | filtres et règles | oui | suggestion attribuée | sélection manuelle et tables |
| Résumer les observations | oui | templates et agrégations | oui | résumé sourcé | rapport structuré déterministe |
| Vérifier la complétude du package | oui | checklist | oui | explication facultative | checklist et validation humaine |
| Proposer une relation à Hypothesis | oui | règles explicables | oui | suggestion avec incertitude | annotation manuelle |
| Qualifier Evidence, Finding, règle ou Intelligence | owner destination | contrôles seulement | workflow de revue | jamais autonome | capabilities propriétaires |

Toute automatisation expose initiateur, producteur/version, Automation Run, Tool Calls, sources, paramètres, statut, erreurs, incertitude et acceptation, modification ou rejet humain.

## 14. États fonctionnels
`draft`, `incomplete`, `ready-for-review`, `submitted`, `partial`, `blocked`, `accepted`, `rejected`, `disputed`, `superseded`, `withdrawn`. États fonctionnels, pas machine objet définitive.

## 15. États d’interface
Loading conserve sélection et Hypothesis ; Empty distingue aucune sélection et aucune donnée ; Partial nomme les limites ; Error conserve les éléments valides ; Offline permet la revue locale sans soumission ; Permission denied masque les valeurs sensibles ; Stale expose source ou assessment obsolète.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Evidence Candidate Package | candidat Investigate | CAP-INV-107/108 | aucune qualification automatique ; sources et limites visibles |
| Finding Draft | brouillon Investigate | CAP-INV-109 | éléments favorables et contradictoires conservés |
| Future Detection Engineering Package | package futur | Phase 4B.3 | aucune règle créée, testée ou déployée |
| Future Intelligence Package | package futur | Phase 4B.3 | aucun Indicator, Campaign ou Threat Actor créé |
| Handoff disposition event | Trace / Activity event | Case / Workbench | destination, auteur, statut et return origin visibles |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-380..396 | sélectionner un handoff | CAP-INV-397 | records, Hypothesis, contradictions, limits, provenance et permissions | Workbench avec sélection préservée |
| CAP-INV-397 | soumettre Evidence candidate | CAP-INV-107/108 | package, lineage, uncertainty, reviewer context et return origin | CAP-INV-397 / Workbench |
| CAP-INV-397 | soumettre Finding Draft | CAP-INV-109 | draft, supporting/contradictory context, sources et return origin | CAP-INV-397 / Workbench |
| CAP-INV-397 | préparer futur package | Phase 4B.3 | network knowledge, conditions, limitations, reproducibility et provenance | CAP-INV-397 / Workbench |

Chaque transition conserve tenant, Case, ownership, restrictions, erreurs, permissions et return origin.

## 18. Dépendances
CAP-INV-102/104/105/107/108/109/301..313/321/329..346/347..362/363..379/380..396 ; Shared Entity/Graph/Timeline/Trace/Activity/Export/Reporting ; Studio provenance ; future Phase 4B.3 ; OPEN-005/008/013/014/015.

## 19. Source de vérité
Investigate est source des packages candidats et de leur disposition avant acceptation. Evidence et Finding restent qualifiés par CAP-INV-107/108/109. Detection Engineering et Intelligence restent entièrement futurs. Les observations sources et objets Shared/Studio/Settings/Command/Govern restent chez leurs propriétaires.

## 20. Provenance et audit
Conserver Case, Hypothesis, capture/version, session, selected observations, coverage/timebase, contradictions, Derived Artifacts, Tools/Calls/Runs, parameters, access restrictions, reproducibility assessment, package versions, auteur, reviewer, destination, accept/reject/modify disposition et return origin.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Masquage | Step-up potentiel | Séparation | Owner | Phase |
|---|---|---:|---|---|---|---|---|
| Candidate package read | données réseau et reasoning | 0 | valeurs sensibles masquées | possible | viewer/reviewer | Investigate | Permissions |
| Candidate package create/update/withdraw | mutation analytique | 2 | minimisation | OPEN-013 | auteur/reviewer | Investigate | Permissions |
| Evidence/Finding submission | qualification potentielle | 2 | données classifiées protégées | possible | preparer/qualifier séparés | Investigate | Permissions |
| Future Detection/Intelligence package preparation | risque de surqualification | 2 | candidats et limites explicites | possible | Investigate/future owner | Investigate / future 4B.3 | Permissions |
| Package export | fuite | 1 | redaction et restrictions | possible | auteur/export reviewer | Shared Export / Security | Permissions |

La matrice atomique, namespaces, RBAC/ABAC et règles finales de step-up restent futurs.

## 22. Limites et erreurs
Sélection vide, provenance incomplète, coverage limitée, timebase incertaine, contradiction non traitée, payload restreint, permission révoquée, destination indisponible ou tenant mismatch. Packet ≠ Evidence ; flow ≠ Evidence ; domaine ≠ IOC confirmé ; certificat ≠ identité certaine ; anomalie ≠ Finding ; package Detection ≠ règle ; package Intelligence ≠ objet canonique.

## 23. Métriques
Packages par destination et état, packages rejetés ou modifiés, provenance complète, contradictions conservées, restrictions d’export, temps de revue et cas de qualification automatique — cible zéro. Aucune cible numérique définitive hors invariant de sécurité.

## 24. Classification de livraison
`defined` / `planned`. Preuve documentaire seulement ; aucune Phase 4B.3, règle Detection, objet Intelligence, moteur, API, protocole, commande, code ou écran détaillé.

## 25. Critères d’acceptation
### 1. Handoff Detection Engineering
**Given** un comportement réseau documenté avec sources et limites
**When** un futur package Detection est préparé
**Then** aucune règle n’est créée, testée ou déployée et Phase 4B.3 reste non commencée

### 2. Handoff Intelligence
**Given** un domaine, certificat ou endpoint candidat
**When** le handoff Intelligence est préparé
**Then** aucun Indicator, Campaign ou Threat Actor n’est créé et la provenance est conservée

### 3. Evidence candidate
**Given** un objet extrait partiel et une anomalie candidate
**When** un package Evidence est soumis
**Then** partialité, contradictions et source restent visibles et la qualification appartient à CAP-INV-107/108

### 4. Sans IA
**Given** aucun fournisseur de modèle
**When** un package est préparé
**Then** sélection manuelle, templates, checklists et revue humaine couvrent le workflow

## 26. Questions ouvertes
OPEN-005 moteurs ; OPEN-008 support ; OPEN-013 classe 2 ; OPEN-014 Artifact/Attachment ; OPEN-015 bridge des Runs. OPEN-011 et OPEN-012 restent ouvertes. Les objets, permissions atomiques, écrans et Phase 4B.3 restent futurs.

## 27. Consommateurs documentaires
Evidence Creation and Review, Finding Management, Case Workspace, Investigation Reporting, future Detection Engineering and Intelligence, phases Objets, Permissions, Journeys, Écrans, Technique et Validation.
