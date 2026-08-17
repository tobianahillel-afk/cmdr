---
id: investigate-network-forensics-sensitive-payload-handling
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
requirements:
  - REQ-SEC-001
  - REQ-SEC-002
---
# Sensitive payload handling

| Niveau | Information | Permission distincte | Classe | Contrôle |
|---:|---|---|---:|---|
| 1 | existence du payload | payload-presence read | 0 | aucune donnée révélée |
| 2 | metadata du payload | packet/transaction metadata read | 0 | minimisation |
| 3 | aperçu masqué | payload-preview | 0 | redaction et taille bornée |
| 4 | lecture | payload-read | 0 | permission et audit |
| 5 | projection déchiffrée autorisée | authorized-decrypted-read | 0/1 | autorisation, policy et step-up potentiel |
| 6 | copie | payload-copy | 2 | justification, SoD et audit |
| 7 | extraction | payload-extract | 1/2 | source bornée, Derived Artifact et lineage |
| 8 | export | payload-export | 1/2 | export policy, minimisation et destinataire |

Interdictions : révélation par défaut, copie libre, extraction/utilisation de secret, transmission à un modèle sans permission, contournement du chiffrement, replay ou appel cible. Un payload inaccessible reste distinct d’un payload absent.
