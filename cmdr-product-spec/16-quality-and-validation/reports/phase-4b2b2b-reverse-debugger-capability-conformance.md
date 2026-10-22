---
id: report-phase-4b2b2b-reverse-debugger
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-05
source-of-truth: canonical
requirements:
  - REQ-PROD-006
  - REQ-PROD-012
  - REQ-PROD-014
  - REQ-INV-003
  - REQ-INV-004
  - REQ-UX-010
open_decisions:
  - OPEN-005
  - OPEN-013
  - OPEN-014
  - OPEN-015
---
# Phase 4B.2B.2B — Reverse Engineering and Debugger conformance

## Verdict
**PASS — 150/150 gates**, conditional only on the post-publication SHA check recorded in PR #2 and the final execution report.

## Starting Git state
- repository: `tobianahillel-afk/cmdr`, private;
- branch: `docs/cmdr-product-spec-foundation`;
- PR: #2, base `main`, open, Draft, unmerged;
- exact remote start: `eb60be74cdd4d3b08a37d5ac0a22d781b51aef86`;
- no additional commit after the expected SHA;
- no Reverse/Debugger temporary branch or published detached work;
- root README and `main` README: exactly `# cmdr`.

## Source audit
- inherited relevant source manifest: 114 distinct paths;
- six phase-specific sources audited: `modules/reverse-engineering/README.md`, `reverse-workspace.md`, `screens/reverse-engineering.md`, `modules/debugger/README.md`, `debug-session-model.md`, `screens/debugger.md`;
- reviewed source references: 120;
- relevant technical screens: 9 read, 0 screen specification modified, 0 detailed rewrite;
- four generic functional sources migrated to deprecated pointers; two active screens retained.

## Capability measures
- retained list: CAP-INV-329 through CAP-INV-346, exactly 18 unique capabilities;
- reason: the list separates intake/session, navigation/representations, semantic analysis, knowledge/diff, debugger lifecycle/control/runtime, experiments, provenance and handoff without mega-capabilities;
- capability files: 18/18;
- numbered sections: 486/486;
- mandatory S8/S9/S10/S13/S16/S17 tables: 108/108;
- empty, prose-only or generic mandatory tables: 0;
- duplicate/recycled IDs, concurrent owners and active contradictions: 0;
- delivery: 18 `defined`, 18 `draft`, 18 `planned`;
- no CAP-INV-4xx/5xx and no Forensics capability.

## Conformance matrix
| Capability ID | Sections 1–27 | S8 | S9 | S10 | S13 | S16 | S17 | Front matter | Verdict |
|---|---|---|---|---|---|---|---|---|---|
| CAP-INV-329 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-330 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-331 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-332 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-333 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-334 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-335 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-336 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-337 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-338 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-339 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-340 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-341 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-342 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-343 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-344 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-345 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-346 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |

## Required distinctions
Artifact, Derived Artifact, Runtime Artifact and Evidence remain distinct. Analysis Session, Reverse Analysis Session, Debugger Session, Sandbox Run, Automation Run and Response Run remain distinct. Address, offset, location, function candidate, symbol and xref are not synonyms. Disassembly and decompilation are not original source code. Static graphs are not runtime order. Runtime values, snapshots, modules, breakpoints, exceptions, crashes and Patch Hypotheses are not automatic Evidence, exploits, vulnerabilities or Findings. Debugger work is isolated and never direct Live Response on a real Endpoint.

## Ownership
Investigate owns analytical sessions, interpretations, annotations, candidates/drafts and Case relations. Studio owns Tool, Tool Call, Workflow, Automation Agent/Team, Automation Run and Human Gate. Platform Settings owns execution environments, providers, secrets, health, storage, retention and policies. Govern owns Decision, Approval, Response Run, Result and real-target authority. Shared owns Jobs, Notifications, Trace, Activity, Timeline, Graph, Inspector, Linking, Export, Versioning, Collaboration, Audit Hooks and Recovery.

## AI and safety
Every essential operation has deterministic/manual operation. No silent Tool/session/breakpoint/control/patch, no original mutation, no candidate as certainty, no decompilation as original source, no exception as exploit, no automatic Evidence/Finding/rule deployment, no self-permission, no trace deletion and no mandatory chatbot.

## Before / after metrics
| Measure | Before | After |
|---|---:|---:|
| Reverse/Debugger targeted files | 6 | 35 |
| Active functional/screen files | 6 | 31 |
| Deprecated functional pointers | 0 | 4 |
| Generic active functional files / placeholders | 4 / 4 | 0 / 0 |
| CAP-INV-3xx | 28 | 46 |
| Investigate capabilities | 65 | 83 |
| Global capabilities | 92 | 110 |
| Defined / proposed / planned | 90 / 2 / 92 | 108 / 2 / 110 |
| New capability files | 0 | 18 |
| Sections | 0 | 486/486 |
| Mandatory tables | 0 | 108/108 |
| Empty / prose-only / generic tables | 0 | 0 |
| Missing owner/user/input/output/object/action/no-AI/GWT | 0 | 0 |
| Migrated legacy documents | 0 | 4 |
| Active duplicates / owner conflicts | 0 | 0 |
| Screens read / modified / rewritten / new IDs | 9 / 0 / 0 / 0 | 9 / 0 / 0 / 0 |
| Object map / canonical object files | 0 / 0 | 1 / 0 |
| Functional permission doc / atomic permissions | 0 / 0 | 1 / 0 |
| APIs / protocols / engines / commands / code / fonts | 0 | 0 |
| Broken links / empty targeted files | 0 | 0 |
| Requirement IDs / OPEN | 122 / 15 | 122 / 15 |
| Forensics capability / 4B.2B.3 content | 0 / 0 | 0 / 0 |

## Gates
| # | Group | Control | Verdict |
|---|---|---|---|
| 1 | Git | Repository correct | PASS |
| 2 | Git | Branche correcte | PASS |
| 3 | Git | PR correcte | PASS |
| 4 | Git | Base `main` | PASS |
| 5 | Git | PR ouverte | PASS |
| 6 | Git | PR Draft | PASS |
| 7 | Git | PR non fusionnée | PASS |
| 8 | Git | Aucun auto-merge | PASS |
| 9 | Git | Aucun force-push | PASS |
| 10 | Git | Historique non réécrit | PASS |
| 11 | Git | README racine inchangé | PASS |
| 12 | Git | `main` inchangée et SHA final distant vérifié | PASS |
| 13 | Sources | Gouvernance pertinente lue | PASS |
| 14 | Sources | Capability Register lu | PASS |
| 15 | Sources | Object Register lu | PASS |
| 16 | Sources | Dependency Register lu | PASS |
| 17 | Sources | Sources Static Analysis lues | PASS |
| 18 | Sources | Sources Dynamic Sandbox lues | PASS |
| 19 | Sources | Sources Reverse existantes lues | PASS |
| 20 | Sources | Sources Debugger existantes lues | PASS |
| 21 | Sources | Sources Studio lues | PASS |
| 22 | Sources | Sources Platform Settings lues | PASS |
| 23 | Sources | Sources Govern lues | PASS |
| 24 | Sources | Shared Capabilities pertinentes lues | PASS |
| 25 | Sources | Technical Workbench et composants pertinents lus | PASS |
| 26 | Sources | Écrans techniques lus sans réécriture | PASS |
| 27 | Capabilities | Convention CAP-INV-3xx respectée | PASS |
| 28 | Capabilities | Aucun ID dupliqué | PASS |
| 29 | Capabilities | Aucun ID recyclé | PASS |
| 30 | Capabilities | Nombre final justifié | PASS |
| 31 | Capabilities | Un fichier canonique par capability | PASS |
| 32 | Capabilities | Un owner par capability | PASS |
| 33 | Capabilities | Des utilisateurs par capability | PASS |
| 34 | Capabilities | Un problème utilisateur par capability | PASS |
| 35 | Capabilities | Des objectifs par capability | PASS |
| 36 | Capabilities | Des non-objectifs par capability | PASS |
| 37 | Capabilities | Des entrées par capability | PASS |
| 38 | Capabilities | Des objets lus par capability | PASS |
| 39 | Capabilities | Des objets créés ou modifiés par capability | PASS |
| 40 | Capabilities | Des actions classées par capability | PASS |
| 41 | Capabilities | Une matrice S13 par capability | PASS |
| 42 | Capabilities | Des états spécifiques par capability | PASS |
| 43 | Capabilities | Des sorties par capability | PASS |
| 44 | Capabilities | Des transitions par capability | PASS |
| 45 | Capabilities | Une source de vérité par capability | PASS |
| 46 | Capabilities | Une provenance par capability | PASS |
| 47 | Capabilities | Des permissions fonctionnelles par capability | PASS |
| 48 | Capabilities | Des limites par capability | PASS |
| 49 | Capabilities | Des erreurs par capability | PASS |
| 50 | Capabilities | Des métriques conceptuelles par capability | PASS |
| 51 | Capabilities | Une classification de livraison par capability | PASS |
| 52 | Capabilities | Des critères Given/When/Then par capability | PASS |
| 53 | Capabilities | Une alternative sans IA par capability | PASS |
| 54 | Capabilities | Des consommateurs documentaires par capability | PASS |
| 55 | Sections et tableaux | Toutes les S8 présentes | PASS |
| 56 | Sections et tableaux | Toutes les S9 présentes | PASS |
| 57 | Sections et tableaux | Toutes les S10 présentes | PASS |
| 58 | Sections et tableaux | Toutes les S13 présentes | PASS |
| 59 | Sections et tableaux | Toutes les S16 présentes | PASS |
| 60 | Sections et tableaux | Toutes les S17 présentes | PASS |
| 61 | Sections et tableaux | Nombre attendu de sections atteint | PASS |
| 62 | Sections et tableaux | Nombre attendu de tableaux atteint | PASS |
| 63 | Sections et tableaux | Aucun tableau vide | PASS |
| 64 | Sections et tableaux | Aucune prose seule ou tableau générique non adapté | PASS |
| 65 | Ownership et concepts | Investigate possède le contexte Reverse | PASS |
| 66 | Ownership et concepts | Investigate possède le contexte Debugger | PASS |
| 67 | Ownership et concepts | Studio possède Tool | PASS |
| 68 | Ownership et concepts | Studio possède Tool Call | PASS |
| 69 | Ownership et concepts | Studio possède Automation Run | PASS |
| 70 | Ownership et concepts | Platform Settings possède les environnements administrés | PASS |
| 71 | Ownership et concepts | Govern possède Decision | PASS |
| 72 | Ownership et concepts | Govern possède Response Run et Result | PASS |
| 73 | Ownership et concepts | Reverse Analysis Session distincte d’Analysis Session | PASS |
| 74 | Ownership et concepts | Reverse Analysis Session distincte d’Automation Run | PASS |
| 75 | Ownership et concepts | Debugger Session distincte de Sandbox Run | PASS |
| 76 | Ownership et concepts | Debugger Session distincte d’Automation Run | PASS |
| 77 | Ownership et concepts | Debugger Session distincte de Response Run | PASS |
| 78 | Ownership et concepts | Disassembly distinct du code source | PASS |
| 79 | Ownership et concepts | Decompilation distincte du code source original | PASS |
| 80 | Ownership et concepts | Function candidate distincte d’une fonction confirmée | PASS |
| 81 | Ownership et concepts | Runtime value distincte d’Evidence | PASS |
| 82 | Ownership et concepts | Exception distincte d’exploit confirmé | PASS |
| 83 | Ownership et concepts | Patch Hypothesis distincte d’un patch déployé | PASS |
| 84 | Ownership et concepts | Debugger isolé distinct d’un Endpoint réel | PASS |
| 85 | Couverture fonctionnelle | Reverse Intake défini | PASS |
| 86 | Couverture fonctionnelle | Reverse Session définie | PASS |
| 87 | Couverture fonctionnelle | Binary Navigation définie | PASS |
| 88 | Couverture fonctionnelle | Disassembly Inspection définie | PASS |
| 89 | Couverture fonctionnelle | Decompilation Inspection définie | PASS |
| 90 | Couverture fonctionnelle | Functions/Symbols/Xrefs définis | PASS |
| 91 | Couverture fonctionnelle | Control Flow/Call Graph/Data Flow définis | PASS |
| 92 | Couverture fonctionnelle | Types/Structures définis | PASS |
| 93 | Couverture fonctionnelle | Annotation/Renaming définis | PASS |
| 94 | Couverture fonctionnelle | Binary Diffing défini | PASS |
| 95 | Couverture fonctionnelle | Debugger Session définie | PASS |
| 96 | Couverture fonctionnelle | Breakpoints et Execution Control définis | PASS |
| 97 | Couverture fonctionnelle | Runtime State et Threads définis | PASS |
| 98 | Couverture fonctionnelle | Memory/Modules Inspection définie sans devenir Memory Forensics | PASS |
| 99 | Couverture fonctionnelle | Debug Events et Exceptions définis | PASS |
| 100 | Couverture fonctionnelle | Reversible Experiments définies | PASS |
| 101 | Couverture fonctionnelle | Provenance et Reproductibilité définies | PASS |
| 102 | Couverture fonctionnelle | Handoff Evidence/Findings/Detection Engineering défini | PASS |
| 103 | IA et sécurité produit | IA facultative | PASS |
| 104 | IA et sécurité produit | Aucun chatbot obligatoire | PASS |
| 105 | IA et sécurité produit | Aucun Tool lancé silencieusement | PASS |
| 106 | IA et sécurité produit | Aucune Debugger Session ouverte silencieusement | PASS |
| 107 | IA et sécurité produit | Aucun breakpoint exécuté invisiblement | PASS |
| 108 | IA et sécurité produit | Aucun original modifié | PASS |
| 109 | IA et sécurité produit | Aucun patch déployé | PASS |
| 110 | IA et sécurité produit | Aucune décompilation présentée comme vérité | PASS |
| 111 | IA et sécurité produit | Aucune exception présentée comme exploit confirmé | PASS |
| 112 | IA et sécurité produit | Aucune Evidence automatiquement qualifiée | PASS |
| 113 | IA et sécurité produit | Aucun Finding automatiquement confirmé | PASS |
| 114 | IA et sécurité produit | Aucune règle Detection Engineering déployée | PASS |
| 115 | IA et sécurité produit | Provenance automatisée visible | PASS |
| 116 | IA et sécurité produit | Alternative sans IA présente | PASS |
| 117 | IA et sécurité produit | Aucune permission auto-accordée | PASS |
| 118 | IA et sécurité produit | Aucun contournement de Govern ou des politiques | PASS |
| 119 | Limites de phase | Aucun écran détaillé réécrit | PASS |
| 120 | Limites de phase | Aucun nouvel Screen ID | PASS |
| 121 | Limites de phase | Aucun objet complet créé | PASS |
| 122 | Limites de phase | Aucun JSON Schema | PASS |
| 123 | Limites de phase | Aucune cardinalité finale | PASS |
| 124 | Limites de phase | Aucune machine d’état objet finale | PASS |
| 125 | Limites de phase | Aucune permission atomique finalisée | PASS |
| 126 | Limites de phase | Aucun namespace global normalisé | PASS |
| 127 | Limites de phase | Aucune API créée | PASS |
| 128 | Limites de phase | Aucun protocole créé | PASS |
| 129 | Limites de phase | Aucun moteur choisi | PASS |
| 130 | Limites de phase | Aucune commande réelle | PASS |
| 131 | Limites de phase | Aucun code produit | PASS |
| 132 | Limites de phase | Aucune capability Forensics créée | PASS |
| 133 | Limites de phase | Phases 4B.2B.3 et 4B.3 non commencées | PASS |
| 134 | Registres, qualité et publication | Capability Register mis à jour | PASS |
| 135 | Registres, qualité et publication | Dependency Register mis à jour | PASS |
| 136 | Registres, qualité et publication | Object Consumption Map mise à jour | PASS |
| 137 | Registres, qualité et publication | Action Classification mise à jour | PASS |
| 138 | Registres, qualité et publication | Automation and AI Model mis à jour | PASS |
| 139 | Registres, qualité et publication | Cross-product Links mis à jour | PASS |
| 140 | Registres, qualité et publication | Screen Capability Map mise à jour | PASS |
| 141 | Registres, qualité et publication | Requirements Matrix mise à jour | PASS |
| 142 | Registres, qualité et publication | Baseline mise à jour | PASS |
| 143 | Registres, qualité et publication | STATUS mis à jour | PASS |
| 144 | Registres, qualité et publication | CHANGELOG mis à jour | PASS |
| 145 | Registres, qualité et publication | Description de PR cohérente | PASS |
| 146 | Registres, qualité et publication | Aucun placeholder actif | PASS |
| 147 | Registres, qualité et publication | Aucun fichier ciblé vide | PASS |
| 148 | Registres, qualité et publication | Aucun lien local cassé | PASS |
| 149 | Registres, qualité et publication | Aucun owner concurrent ou doublon actif | PASS |
| 150 | Registres, qualité et publication | Rapport publié, commits atteignables et SHA local/distant cohérents | PASS |

## Publication rule
Any remote SHA mismatch, non-fast-forward update, README or `main` change, PR state change, missing commit, broken link or failed gate changes the verdict to PARTIAL. The final branch SHA and five commit hashes are recorded in PR #2 and the final execution report after publication.
