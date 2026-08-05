---
id: investigate-network-forensics-scope
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
---
# Scope

## Inclus
Capture intake, integrity, scope, coverage, timebase, interfaces, packet/frame inspection, flow/session reconstruction, protocol/conversation analysis, DNS, Web/application transactions, encrypted metadata/certificates, transfer reconstruction, Entity relations, behavior candidates, Network Timeline, comparisons, extraction, provenance, reproducibility et handoffs.

## Exclus
Acquisition active, capteur, Fleet/Policy administration, live monitoring, Event Search concurrent, SIEM, packet generation/injection/crafting, active scanning/interception/MITM, replay, interaction cible, déchiffrement non autorisé, key extraction, credential use, exploit, Detection Engineering, Intelligence canonique, Cloud/Mobile, API, protocole interne, moteur, commande, code et écran détaillé.

## Invariants
- capture valide ≠ capture représentative ;
- packet loss ≠ absence de trafic ;
- observation ≠ conclusion ;
- payload inaccessible ≠ payload absent ;
- aucune fonction essentielle ne dépend d’un modèle.
