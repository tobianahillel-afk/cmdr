---
id: investigate-reverse-engineering-and-debugger
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
requirements:
  - REQ-INV-003
  - REQ-INV-004
  - REQ-PROD-014
  - REQ-PROD-020
open_decisions:
  - OPEN-005
  - OPEN-013
  - OPEN-014
  - OPEN-015
---
# Reverse Engineering and Debugger

## Mission
Permettre à un analyste autorisé d’interpréter un Artifact, de naviguer entre représentations statiques, de conduire une Debugger Session isolée et de préparer des handoffs qualifiables sans modifier l’original, viser un Endpoint réel ou choisir un moteur.

## Capabilities
CAP-INV-329 through CAP-INV-346 form one canonical functional module. Reverse Engineering owns analytical interpretation; Debugger owns isolated session context; Tool and Automation lifecycle remain Studio; environments remain Platform Settings; real-target authority remains Govern.

## Technical Workbench
One main canvas, one right Inspector, at most two auxiliary panels, at most six visible technical tabs, optional bottom Console, contextual Artifact Explorer, Automation Tray closed by default, keyboard focus restoration, structured alternatives to graphs and accessible resizers.

## Invariants
- disassembly and decompilation are representations, not original source code;
- candidate functions, symbols, types and names are not certainties;
- Debugger Session is not Sandbox Run, Automation Run or Response Run;
- runtime values, snapshots, exceptions and patch experiments are not Evidence or Findings automatically;
- original Artifacts are immutable;
- no direct debugging of a real Endpoint;
- all essential workflows operate without an AI provider.

## Delivery
All eighteen capabilities are `defined`, `draft` and `planned`. No engine, debugger, third-party product, API, protocol, command or product code is selected.
