---
id: investigate-network-forensics-permissions
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
open_decisions:
  - OPEN-013
---
# Functional permissions

| Famille | Risque | Classes | Masquage / gate | Owner futur |
|---|---|---:|---|---|
| Capture/session read and manage | source et contexte sensible | 0,2 | tenant, scope, contributor rights | Investigate/Security |
| Packet metadata / raw packet | brut réseau | 0 | raw séparé de metadata | Security/Permissions |
| Payload presence / preview / read | données privées ou secrets | 0 | masqué par défaut, droits distincts | Security/Privacy |
| Authorized decrypted projection | contenu hautement sensible | 0/1 | autorisation, step-up et audit potentiels | Security/Govern |
| Payload copy/extract/export | exfiltration secondaire | 1/2 | permissions séparées, minimisation, SoD | Security/Export |
| Reconstruction/decode/comparison | volume et interprétation | 1 | Tool/version/scope explicites | Investigate/Studio |
| Entity relation/anomaly review | attribution erronée | 2 | review et contestation | Investigate/Shared |
| Evidence/Finding/future handoff | qualification | 2 | owner destination qualifie | Investigate |
| Cross-tenant analysis | isolation | 0/1/2 | interdit par défaut | Security |

Aucun namespace, RBAC/ABAC, permission atomique ou règle définitive de step-up n’est finalisé.
