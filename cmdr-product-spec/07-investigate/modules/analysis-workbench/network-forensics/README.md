---
id: investigate-network-forensics
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
requirements:
  - REQ-INV-001
  - REQ-PROD-014
  - REQ-PROD-020
  - REQ-AI-002
  - REQ-SEC-001
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
---
# Network Forensics

## Mission
Fournir dans Investigate un contexte Technical Workbench en lecture seule pour analyser des captures réseau, reconstruire des observations sourcées, corréler avec les autres analyses et préparer des handoffs sans acquisition active, interaction cible, qualification automatique ou dépendance à l’IA.

## Capability range
CAP-INV-380..397 — dix-huit capabilities `defined` et `planned`.

## Ownership
Investigate possède le contexte, la session fonctionnelle, les observations, interprétations, relations analytiques, Derived Artifacts et packages candidats. Collection/Endpoint Agent produit la capture et ses limitations. Settings administre capteurs, Fleet, Policies, providers, storage, retention, health, time synchronization et secrets. Studio possède Tool, Tool Call, Workflow et Automation Run. Command possède Detection, Signal, Alert et Incident. Govern possède l’autorité réelle. Shared possède Entity, Graph, Timeline et les mécanismes transversaux.

## Functional chain
Case / Collection Job / Artifact → Intake → Integrity/Scope → Coverage/Timebase → Packet/Frame → Flow/Session → Conversation/Protocol → DNS/Transactions/Encrypted Metadata/Transfers → Entities/Behavior/Timeline/Comparison → Extraction/Provenance → Evidence, Findings et futurs handoffs.

## Safety
Aucune génération ou injection de packet, interception active, scanning actif, replay, appel cible, déchiffrement non autorisé, extraction/utilisation de secret, exploit, règle Detection, objet Intelligence canonique, API, commande ou code.

## Workbench constraints
Un canvas principal, un Inspector droit, deux panneaux auxiliaires maximum, six onglets visibles maximum, console basse optionnelle, Artifact Explorer contextuel, Automation Tray fermé par défaut, payloads sensibles masqués, focus et return origin préservés.

## Delivery
Aucun moteur, produit tiers, protocole final, format, écran détaillé ou implémentation n’est revendiqué. La phase 4B.3 reste non commencée.
