---
id: ownership-register
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-OBJ-001
  - REQ-OBJ-002
  - REQ-OBJ-003
  - REQ-OBJ-004
  - REQ-OBJ-005
  - REQ-OBJ-006
  - REQ-OBJ-007
  - REQ-OBJ-008
  - REQ-OBJ-009
  - REQ-OBJ-010
  - REQ-OBJ-011
  - REQ-OBJ-012
---
# Registre de propriété canonique

## Règle

Le propriétaire définit la sémantique, les états et les changements d'un concept. Un consommateur peut afficher, filtrer, lier ou demander une transition autorisée, mais ne peut pas créer une seconde définition, machine d'état ou permission.

| Concept | Propriétaire | Raison | Fichier canonique | Consommateurs | Actions des consommateurs | Restrictions | Requirement IDs | Question |
|---|---|---|---|---|---|---|---|---|
| Incident | Command | situation et coordination opérationnelles | `../05-domain-model/objects/incident.md` | Investigate, Govern, Reporting | lire, lier Case/Decision/Result, projeter | pas de cycle de vie concurrent | REQ-OBJ-001 | — |
| Task | Command, sauf tâche spécialisée déclarée | distribution du travail opérationnel | `../05-domain-model/objects/task.md` | tous produits | créer selon workflow, affecter, lier | une tâche spécialisée indique son owner | REQ-PROD-009 | — |
| Case | Investigate | workspace d'investigation | `../05-domain-model/objects/case.md` | Command, Govern, Reporting | ouvrir, lire projection, lier | Command ne possède pas le Case | REQ-OBJ-002 | — |
| Hypothesis | Investigate | raisonnement analytique | `../05-domain-model/objects/hypothesis.md` | Studio comme support, Reporting | proposer, lire, relier | une proposition IA reste attribuée | REQ-OBJ-002 | — |
| Artifact | Investigate | élément technique acquis ou dérivé | `../05-domain-model/objects/artifact.md` | Command, Govern, Studio | prévisualiser, lier, traiter selon permission | ne devient Evidence qu'après workflow | REQ-PROD-061 | OPEN-014 |
| Evidence | Investigate | preuve avec provenance et intégrité | `../05-domain-model/objects/evidence.md` | Command, Govern, Reporting | lire projection, citer, vérifier autorisation | aucun consommateur ne modifie rétroactivement | REQ-OBJ-003 | — |
| Finding | Investigate | conclusion analytique soutenue | `../05-domain-model/objects/finding.md` | Command, Govern, Studio | lire, demander action, proposer enrichissement | ne vaut pas Decision ni Result | REQ-OBJ-004 | — |
| Action Request | Govern pour le lifecycle ; producteurs Command/Investigate | demande gouvernée précédant une décision | `../05-domain-model/objects/action-request.md` | Command, Investigate, Govern | créer, compléter, suivre | le producteur ne s'auto-approuve pas | REQ-PROD-004 | — |
| Decision | Govern | autorité et justification | `../05-domain-model/objects/decision.md` | Command, Investigate, Studio, Endpoint | fournir contexte, lire résultat | aucun agent ni produit consommateur ne décide implicitement | REQ-OBJ-005 | — |
| Approval | Govern | expression d'autorité participant à une Decision | `../05-domain-model/objects/approval.md` | Studio Human Gate, Audit | soumettre une validation, lire | Approval n'est pas synonyme de Decision | REQ-OBJ-006 | OPEN-007 |
| Response Run | Govern | exécution d'une réponse autorisée | `../05-domain-model/objects/response-run.md` | Command, Investigate, Endpoint Agent, Studio | observer, exécuter une étape autorisée, vérifier | ne pas confondre avec Automation Run | REQ-OBJ-007 | OPEN-015 |
| Result | Govern | résultat vérifié de réponse | `../05-domain-model/objects/result.md` | Command, Investigate, Reporting | consommer, relier, mesurer | ne remplace pas Finding | REQ-OBJ-007 | — |
| Endpoint | modèle partagé, administré par Platform Settings | représentation commune de la cible | fichier canonique planifié en Phase 7 | Command, Investigate, Govern, Endpoint Agent | afficher, cibler selon permissions | Settings administre, aucun produit ne redéfinit | REQ-PROD-017 | OPEN-008 |
| Endpoint Agent Fleet | Platform Settings | inventaire administratif, version, santé, policy | `../05-domain-model/objects/endpoint-agent-fleet.md` | Investigate, Govern, Command | lire projection, sélectionner endpoint | Investigate n'administre pas la flotte | REQ-OBJ-008 | — |
| Skill | CMDR Studio | connaissance ou méthode versionnée | `../05-domain-model/objects/skill.md` | Command, Investigate, Govern | invoquer selon permission | ne devient pas Tool | REQ-OBJ-009 | — |
| Tool | CMDR Studio | interface exécutable et gouvernée | fichier canonique planifié en Phase 7 | produits opérationnels | appeler par Tool Call autorisé | pas d'appel caché ou hors trace | REQ-OBJ-009 | — |
| Tool Call | CMDR Studio | tentative immuable d'exécuter un Tool | fichier canonique planifié en Phase 7 | Govern, Audit, Control Room | lire, interrompre selon contrat | absent du modèle détaillé actuel | REQ-OBJ-009 | — |
| Automation Agent | CMDR Studio | entité agentique configurée | `../05-domain-model/objects/automation-agent.md` | Command, Investigate, Govern | déclencher ou utiliser | ne possède ni objet métier ni décision | REQ-AI-002, REQ-OBJ-009 | — |
| Agent Team | CMDR Studio | composition coordonnée d'agents | `../05-domain-model/objects/agent-team.md` | produits opérationnels | déclencher selon workflow | aucune autorité implicite | REQ-OBJ-009 | — |
| Workflow | CMDR Studio | orchestration explicite | `../05-domain-model/objects/workflow.md` | tous produits | exécuter une version autorisée | ne remplace pas le workflow métier propriétaire | REQ-OBJ-009 | — |
| Human Gate | CMDR Studio | suspension d'une Automation Run | `../05-domain-model/objects/human-gate.md` | Govern, produits opérationnels | fournir validation humaine | relation à Decision encore ouverte | REQ-AI-004, REQ-OBJ-009 | OPEN-007 |
| Automation Run | CMDR Studio | exécution d'un Workflow ou Agent | fichier canonique planifié en Phase 7 | produits opérationnels, Audit | observer, interrompre, reprendre | ne devient Response Run que par transition gouvernée | REQ-PROD-062, REQ-OBJ-009 | OPEN-015 |
| Reporting Engine | Shared Capabilities | production et publication de rapports multi-produits | `../12-shared-capabilities/reporting-engine.md` | Command, Investigate, Govern | composer, revoir, publier selon permission | un rapport local ne recrée pas le moteur | REQ-OBJ-010 | — |
| Saved Views génériques | Shared Capabilities | persistence commune de vues | `../12-shared-capabilities/saved-view-engine.md` | tous produits | enregistrer et appliquer | ne contient pas les règles Work Queue | REQ-OBJ-011 | conflit à corriger Phase 3/6 |
| Work Queue Saved Views | Command | organisation du travail Command | `../06-command/modules/incidents-and-work-queue/saved-views.md` | Command | créer, partager, appliquer | non référencées comme source générique hors Command | REQ-OBJ-012, REQ-UX-008 | conflit à corriger Phase 6 |
| Permission Model | Security, Permissions and Trust | sémantique d'autorisation commune | `../14-security-permissions-and-trust/permission-model.md` | tous produits | référencer des permissions | aucune permission locale concurrente | REQ-PROD-006 | namespaces détaillés Phase 7 |
| Inspector | Design System | composant de contexte sélectionné | `../03-design-system/components/inspector.md` | tous produits | configurer sections et actions locales | pas de panneau droit concurrent | REQ-UX-002 | — |
| Context Bar component | Design System | anatomie et interaction visuelles | `../03-design-system/components/context-bar.md` | tous produits | afficher contexte | ne possède pas le modèle de propagation | REQ-UX-006 | — |
| Context propagation capability | Shared Capabilities / Experience Architecture | conservation et transport du contexte | `../04-experience-architecture/context-propagation.md` | tous produits | lire et transmettre contexte | ne redéfinit pas le composant visuel | REQ-UX-006 | — |

## Règles de changement

Un changement de propriétaire exige une ADR, la mise à jour de ce registre, des sources canoniques, des permissions, des écrans, parcours et contrats consommateurs.

## Critères d'acceptation

**Given** `Endpoint Agent Fleet`,  
**When** le registre est consulté,  
**Then** Platform Settings est le seul propriétaire, Investigate et Govern sont consommateurs, Endpoint Agent est le composant administré et aucun autre fichier actif ne revendique le lifecycle de la flotte.

**Given** une Evidence affichée dans Command,  
**When** l'utilisateur ouvre sa projection,  
**Then** la source Investigate reste identifiable et Command ne redéfinit ni schéma ni état.
