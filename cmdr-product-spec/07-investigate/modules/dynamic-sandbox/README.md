---
id: investigate-dynamic-sandbox
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-014
  - REQ-PROD-020
  - REQ-INV-005
open_decisions:
  - OPEN-005
  - OPEN-013
  - OPEN-014
  - OPEN-015
---
# Dynamic Sandbox and Behavioral Analysis

## Mission
Permettre une analyse dynamique explicite d’un Artifact dans un Sandbox Environment autorisé, isolé et administré par Platform Settings, puis organiser les observations, Runtime Artifacts, comparaisons, provenance et handoffs analytiques.

## Périmètre 4B.2B.2A
Intake, sélection d’environnement, Dynamic Analysis Session, Sandbox Run, timeline comportementale, processus, fichiers et système, réseau, persistance candidate, profils d’interaction, Runtime Artifacts, comparaison multi-Run, sécurité de l’analyse, provenance, reproductibilité et Evidence/Finding handoff.

## Exclusions
Reverse Engineering, désassemblage, décompilation, Debugger, breakpoints, patching, Memory/Disk/Network Forensics avancé, cloud/mobile analysis, moteur, hyperviseur, instrumentation, API, protocole, commande ou code.

## Ownership
Investigate possède le contexte analytique, les observations et relations au Case. Platform Settings administre les Sandbox Environments. Studio possède Tools et Automation Runs. Govern conserve l’autorité sur les cibles réelles. Shared possède les mécanismes transversaux.

## Invariants
Aucun Run silencieux ; aucun réseau réel silencieux ; aucun Runtime Artifact automatiquement Evidence ; aucune observation automatiquement Finding ; timeout ou absence d’observation ne vaut pas absence de menace.
