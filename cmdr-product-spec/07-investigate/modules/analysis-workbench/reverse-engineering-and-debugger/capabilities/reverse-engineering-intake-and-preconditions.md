---
id: CAP-INV-329
title: Reverse Engineering Intake and Preconditions
product: investigate
module: analysis-workbench
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-05
requirement_ids:
  - REQ-INV-003
  - REQ-PROD-014
  - REQ-PROD-020
  - REQ-PROD-052
  - REQ-SEC-001
  - REQ-UX-002
open_decisions:
  - OPEN-005
  - OPEN-013
  - OPEN-014
  - OPEN-015
source-of-truth: canonical
---

# CAP-INV-329 — Reverse Engineering Intake and Preconditions

## 1. Définition
Reverse Engineering Intake and Preconditions définit le comportement produit permettant à un analyste autorisé d’évaluer un Artifact, ses résultats statiques et dynamiques, ses restrictions et les Tools compatibles avant toute ouverture de Reverse Analysis Session.

## 2. Problème utilisateur
à un analyste autorisé d’évaluer un Artifact, ses résultats statiques et dynamiques, ses restrictions et les Tools compatibles avant toute ouverture de Reverse Analysis Session.

## 3. Objectifs
- ouvrir depuis Case, Static Analysis ou Dynamic Analysis sans exécution implicite
- rendre visibles type, architecture détectée ou déclarée, formats, restrictions et compatibilités
- définir objectif, scope, owner et return origin
- créer ou reprendre explicitement une Reverse Analysis Session

## 4. Non-objectifs
- lancer automatiquement un Tool ou une analyse
- choisir un moteur de désassemblage, de décompilation ou un debugger
- définir les formats internes d’adresse ou les algorithmes de détection

## 5. Propriétaire
Investigate possède le contexte analytique, les annotations et les relations au Case. Studio conserve les Tools, Tool Calls et Automation Runs; Platform Settings administre les environnements; Govern conserve l’autorité sur toute cible réelle; Shared conserve les mécanismes transversaux.

## 6. Utilisateurs
Malware Analyst principal; Reverse Engineer et Investigation Lead comme contributeurs ou reviewers.

## 7. Conditions d’entrée
- un Case ou un Artifact accessible
- provenance et restrictions disponibles ou explicitement manquantes
- permissions générales évaluées
- aucune cible réelle sélectionnée

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Artifact et lineage | CAP-INV-105/311/324 | source analytique | Oui | version active | bloquer l’intake si la source n’est pas accessible |
| Résultats statiques | CAP-INV-301..313 | observations et Derived Artifacts | Non | dernière version connue | continuer en signalant l’absence |
| Résultats dynamiques | CAP-INV-314..328 | Runtime Artifacts et Behavioral Observations | Non | session sélectionnée | continuer sans inventer de comportement |
| Compatibilité Tools/environnements | Studio / Platform Settings | capabilities, versions et health | Oui | courante | état tool-unavailable ou environment-unavailable |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | contexte, Hypothesis et return origin | lecture et liaison uniquement |
| Artifact / Derived Artifact / Runtime Artifact | Investigate | source, versions, restrictions et provenance | lecture; aucun original modifié |
| Tool / Tool Call / Automation Run | CMDR Studio | outil, version, exécution et attribution | lecture et sélection; lifecycle Studio |
| Analysis Session / Dynamic Analysis Session | Investigate | résultats statiques et dynamiques existants | lecture et liaison |
| Environment / Tool availability | Platform Settings / Studio | compatibilité, health et restrictions | projection en lecture seule |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Trace / Activity event | Shared | émission de l’action et de son résultat | append-only; aucune suppression locale |
| Reverse Intake | Investigate | create/update/resume | aucun Tool lancé à l’ouverture |
| Reverse Analysis Session link | Investigate | prepare or resume | création explicite uniquement |

## 11. Fonctionnalités
- afficher source, provenance, type, architecture, formats et restrictions
- présenter résultats statiques, dynamiques et Derived/Runtime Artifacts disponibles
- détecter format ou architecture non supporté sans simuler un résultat
- définir objectif, scope et owner
- créer ou reprendre explicitement la session et préserver le return origin

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspecter les préconditions | Malware Analyst | Reverse Intake | 0 | Artifact lisible | assessment visible | Non |
| Modifier objectif ou scope | Reverse Engineer | Reverse Intake | 2 | session non active | intake versionné | OPEN-013 |
| Créer ou reprendre la session | Malware Analyst | Reverse Analysis Session | 2 | intake ready | session liée, aucun Tool lancé | Non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Inspecter les préconditions Reverse | Oui | Oui | Oui | Oui | viewer déterministe et navigation manuelle de les préconditions Reverse |
| Proposer un Tool ou une voie compatible | Oui | Oui | Oui | Oui | règles, heuristiques visibles, comparaison et revue humaine |
| Préparer un handoff | Oui | Oui | Oui | Oui | sélection manuelle, checklist et validation humaine |

Toute sortie automatisée expose initiateur, producteur et version, Automation Run et Tool Calls lorsqu’ils existent, sources, paramètres, timestamp, statut, incertitude, owner humain et disposition acceptée, modifiée ou rejetée.

## 14. États fonctionnels
- draft
- incomplete
- ready
- unsupported
- restricted
- tool-unavailable
- environment-unavailable
- policy-blocked

Ces états sont fonctionnels et ne constituent pas une machine d’état objet définitive.

## 15. États d’interface
- **Loading** conserve le Workbench, l’Artifact, la sélection et le return origin.
- **Empty** explique l’absence de résultat sans simuler une analyse.
- **Partial** identifie les sources, vues ou événements manquants et leurs conséquences.
- **Error** conserve les résultats valides, l’erreur et une reprise sûre.
- **Offline** limite les mutations et affiche la dernière synchronisation.
- **Permission denied** ne révèle aucune donnée protégée.
- **Stale** distingue la dernière observation connue de l’état courant.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Assessment d’intake | Reverse Intake | CAP-INV-330 | préconditions, absences et restrictions explicites |
| Contexte de navigation | Context reference | Technical Workbench | Case, Artifact, objectif et return origin conservés |
| Blocage explicite | Policy/compatibility event | analyste et audit | aucune analyse fictive ni lancement silencieux |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Case | ouvrir Reverse Engineering | Reverse Intake | tenant, Case, Artifact, Hypothesis, objectif, return origin | retour au Case |
| Static Analysis | demande Reverse | Reverse Intake | Artifact, Analysis Session, résultats, Derived Artifacts, annotations | retour Static Analysis |
| Dynamic Analysis | demande Reverse | Reverse Intake | Artifact, Runtime Artifacts, Runs, observations, provenance | retour Dynamic Sandbox |
| Reverse Intake | validation humaine | Reverse Analysis Session | type, architecture, restrictions, objectif, Tools compatibles, owner | retour intake en erreur |

Chaque transition conserve tenant, environnement, Case, Artifact, sélection, permissions et return origin. Une erreur ne détruit pas la source ni les résultats déjà valides.

## 18. Dépendances
- CAP-INV-301,307,311,312,313
- CAP-INV-314,324,327,328
- CAP-INV-105 Artifact Management
- CMDR Studio Tool catalog
- Platform Settings environment and policy projections
- OPEN-005/013/014/015

## 19. Source de vérité
L’Artifact et sa provenance viennent d’Investigate; les résultats statiques/dynamiques restent dans leurs sessions propriétaires; Tool/version vient de Studio; environnement et policy viennent de Settings.

## 20. Provenance et audit
Conserver acteur, source d’entrée, version Artifact, compatibilités évaluées, restrictions, objectif, scope, décision create/resume et return origin.

## 21. Permissions fonctionnelles
| Besoin | Risque | Classe | Step-up éventuel | Séparation des tâches | Owner | Phase propriétaire |
|---|---|---:|---|---|---|---|
| Reverse Analysis read | exposition de contenu sensible | 0 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| Reverse Analysis Session create/resume | création de contexte analytique | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| sensitive Artifact read | accès à contenu restreint | 0 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |
| automated reverse analysis request | orchestration potentiellement coûteuse | 2 | selon sensibilité, environnement et policy | initiateur distinct du reviewer lorsque requis | Investigate / Security | Permissions phase |

La matrice atomique, les namespaces et le modèle RBAC/ABAC restent reportés à la phase Permissions.

## 22. Limites et erreurs
- format ou architecture inconnue visible
- Tool indisponible sans fallback fictif
- Artifact restreint ou cross-tenant bloqué
- aucune analyse automatique à l’ouverture
- résultats partiels conservés

## 23. Métriques
- taux d’intakes ready/blocked
- raisons unsupported/tool-unavailable
- temps entre intake et décision explicite
- retours au contexte source réussis

## 24. Classification de livraison
Delivery status `defined`; delivery mode `planned`. La promotion exige objets et permissions approuvés, Tools et environnements évalués, contrats techniques, tests de sécurité, écrans et release evidence. Aucun moteur, debugger, produit tiers ou implémentation n’est choisi.

## 25. Critères d’acceptation
### Scénario 1
**Given** un Artifact lié à un Case, architecture non supportée et aucun Tool compatible
**When** l’analyste ouvre l’intake
**Then** architecture et indisponibilité sont visibles; aucun Tool n’est lancé; résultats existants et return origin restent accessibles

### Scénario 2
**Given** un Artifact compatible sans fournisseur IA
**When** l’analyste prépare la session
**Then** les contrôles déterministes et manuels restent disponibles

### Scénario 3
**Given** un Artifact restreint sans permission
**When** l’intake évalue l’accès
**Then** aucun contenu protégé n’est révélé et la demande d’accès reste externe

## 26. Questions ouvertes
- OPEN-005 reste ouverte; aucune décision n’est fermée par cette spécification.
- OPEN-013 reste ouverte; aucune décision n’est fermée par cette spécification.
- OPEN-014 reste ouverte; aucune décision n’est fermée par cette spécification.
- OPEN-015 reste ouverte; aucune décision n’est fermée par cette spécification.

Les schémas, cardinalités, machines d’état, formats d’adresse et permissions atomiques sont reportés aux phases propriétaires.

## 27. Consommateurs documentaires
- Analysis Workbench module
- Reverse Engineering screen INV-REV-001
- Case Workspace
- Static Analysis and Dynamic Sandbox handoffs
- future object and permission specifications
