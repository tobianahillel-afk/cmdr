---
id: operating-model
domain: 01-product-vision
status: draft
owner: Head of Product
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-003
  - REQ-PROD-004
  - REQ-PROD-008
  - REQ-PROD-009
  - REQ-AI-007
  - REQ-SEC-001
  - REQ-SEC-002
---
# Modèle opérationnel

## Progression de travail

| Étape | Propriétaire principal | Objet / résultat | Handoff | Automatisation possible | Contrôle humain |
|---|---|---|---|---|---|
| Détection | moteur / Command | Detection, Signal, Alert | vers triage | règles déterministes, enrichissement | disposition et escalade |
| Triage | Command | Incident ou disposition | vers coordination / Case | résumé et priorisation proposés | owner et priorité |
| Coordination | Command | Incident, Tasks, situation | vers Investigate / Govern | handover préparé | Incident Commander |
| Investigation | Investigate | Case, Hypothesis | vers collecte / Finding | requêtes et pivots assistés | analyste responsable |
| Collecte | Investigate + Endpoint | Artifact / Evidence candidate | vers vérification | jobs gouvernés | scope et intégrité |
| Preuve | Investigate | Evidence | vers Finding | extraction déterministe/assistée | validation de provenance |
| Conclusion | Investigate | Finding | vers Action Request | proposition et rédaction | reviewer analytique |
| Demande | Command/Investigate → Govern | Action Request | vers review | complétude vérifiée | initiateur responsable |
| Décision | Govern | Decision, Approval | vers Run ou refus | policy evaluation | approbateur / autorité |
| Réponse | Govern + opérateur/Endpoint | Response Run | vers vérification | workflow et outils autorisés | interruption possible |
| Vérification | Govern / Investigate | Result | retour Command | collecte de preuve de résultat | acceptation du résultat |
| Amélioration | Command / Investigate / Studio | actions, règles, readiness | vers backlog | analyse de tendances | owner de capability |

## Ownership et handoff

Un handoff transmet : objets et identifiants, owner et prochaine action, contexte tenant/environnement, Evidence et incertitudes pertinentes, décisions en attente, contraintes de temps et lien de retour. Le destinataire n'acquiert pas la propriété des objets sources.

## Règles déterministes

Les règles, queries, policies et workflows déterministes restent visibles, versionnés et testables. Un agent peut les invoquer ou les proposer, pas les rendre inutiles ou invisibles.

## Automation Agents

Ils participent comme acteurs attribués, produisent Tool Calls et outputs traçables, ne s'auto-approuvent pas, peuvent être interrompus et repris, retournent le résultat au produit propriétaire et créent une Action Request lorsqu'une action exige Govern.

## Human Gates

Une Human Gate suspend une Automation Run. Elle ne devient pas automatiquement une Decision Govern. `OPEN-007` précisera les cas où la validation doit être matérialisée en Approval/Decision.

## Classification des actions

| Classe | Nature | Exemple | Gouvernance conceptuelle |
|---|---|---|---|
| 0 | observation | lire, rechercher, inspecter | permission et audit selon donnée |
| 1 | collecte | fichier, logs, mémoire | scope, intégrité et audit |
| 2 | modification réversible | configuration avec rollback | policy et autorité selon contexte |
| 3 | containment | isolation, quarantaine, kill, revoke | Govern requis par défaut |
| 4 | destructive / irréversible | suppression ou changement sans rollback fiable | Decision explicite et approbations |

`OPEN-013` doit préciser les valeurs par défaut de la classe 2.

## Escalade

Une escalade est déclenchée par risque, manque d'autorité, Evidence insuffisante, impact métier, erreur de scope, échec de rollback ou conflit de policy. Elle conserve le contexte et l'owner actuel jusqu'à acceptation du handoff.

## Audit

Chaque transition importante enregistre acteur, rôle, objet source, objet cible, action, résultat, justification, policy, automation/run et correlation ID lorsque ces éléments existent.

## Critère d'acceptation

**Given** une demande d'isolation issue d'un Finding,  
**When** elle traverse Govern,  
**Then** le scope, la classe 3, les Evidence, l'approbateur, le rollback, le Response Run, la vérification et le Result sont reliés ; l'analyste reste auteur du Finding et Govern propriétaire de la Decision.
