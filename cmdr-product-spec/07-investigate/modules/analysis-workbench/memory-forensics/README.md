---
id: investigate-memory-forensics-canonical
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
requirements:
  - REQ-INV-001
  - REQ-PROD-014
  - REQ-PROD-020
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
---
# Memory Forensics

## Mission
Fournir le contexte fonctionnel dans lequel un analyste autorisé ouvre une Memory Image, vérifie son acquisition et son intégrité, sélectionne un profil, reconstruit des observations mémoire, protège les données sensibles, extrait des Derived Artifacts et prépare des handoffs explicites.

## Périmètre de Phase 4B.2B.3A
CAP-INV-347..362 couvrent intake, Memory Forensics Session, integrity review, platform/profile, processus et threads, régions et mappings, modules et drivers, handles et IPC, network state présent en mémoire, sensitive material, anomalies, kernel state, timeline, extraction, provenance et handoff.

## Exclusions
Aucune acquisition, commande, méthode d’extraction de secret, moteur, plugin, offset, algorithme, API, protocole, format final, Disk/Filesystem/Network Forensics complète, Cloud, Mobile ou Detection Engineering.

## Ownership
Investigate possède Memory Image, contexte Forensics, sessions, observations, annotations, Derived Artifacts et candidates/drafts. Endpoint Agent contribue à l’acquisition. Settings possède Fleet, Policies, plateformes, stockage, rétention et santé. Studio possède Tools et Runs. Govern possède l’autorité sur les cibles réelles. Shared possède les mécanismes transversaux.

## Invariants
- Memory Image ≠ Memory Forensics Session ≠ Evidence.
- Memory Forensics ≠ Debugger Memory View.
- Une image partielle reste partielle.
- Une anomalie, connexion ou donnée sensible candidate ne devient pas automatiquement Finding, IOC ou credential valide.
- Les valeurs sensibles sont masquées par défaut.
- Toute fonction essentielle possède une voie déterministe et manuelle.
- Aucun écran détaillé n’est réécrit.

## Technical Workbench
Un canvas principal, un Inspector droit, deux panneaux auxiliaires maximum, six onglets techniques visibles maximum, console basse optionnelle, Artifact Explorer contextuel, Automation Tray fermé par défaut, historique/retour/focus conservés, resizers accessibles, graphes avec alternative structurée et valeurs sensibles masquées.

## Delivery
Seize capabilities `draft`, `defined` et `planned`; aucune implémentation ou intégration native n’est revendiquée.
