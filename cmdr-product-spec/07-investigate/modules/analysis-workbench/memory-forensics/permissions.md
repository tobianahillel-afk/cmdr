---
id: investigate-memory-forensics-permissions
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
requirements:
  - REQ-SEC-001
  - REQ-SEC-002
  - REQ-PROD-060
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
---
# Functional permissions

| Besoin | Capability | Risque | Classe | Masquage / step-up | Séparation des tâches | Owner / phase |
|---|---|---|---:|---|---|---|
| Memory Image read/raw/restricted | 347/349 | contenu très sensible | 0 | classification et step-up possibles | acquisition/review distincts | Investigate/Security |
| Session create/update/close/reopen | 348 | contexte analytique | 2 | OPEN-013 | owner/contributor/reviewer | Investigate/Security |
| Platform/profile select | 350 | interprétation structurante | 2 | possible | analyst/reviewer | Investigate/Security |
| Process/thread/region/module/handle/network/kernel read | 351..359 | données système sensibles | 0 | least privilege | reviewer selon cas | Investigate/Security |
| Automated analysis | 351..361 | traitement borné | 1 | confirmation explicite | initiateur/reviewer | Studio/Investigate |
| Sensitive presence / masked read | 356 | donnée sensible | 0 | masquée par défaut | least privilege | Security |
| Sensitive reveal | 356 | exposition de secret | 0 | step-up probable | requester/reviewer distincts | Security |
| Sensitive copy | 356 | exfiltration | 1 | step-up fort | double contrôle | Security |
| Sensitive export | 356 | diffusion externe | 1/2 | policy stricte et step-up | approbateur distinct | Security/Shared |
| Derived Artifact extract/export | 360 | capture et diffusion | 1/2 | policy et step-up | extractor/reviewer/exporter | Investigate/Shared |
| Evidence/Finding/Detection handoff prepare | 362 | impact de conclusion | 2 | possible | author/reviewer | Investigate/future owner |
| Cross-tenant analysis | toutes | isolement tenant | 0/1/2 | contrôle obligatoire | séparation forte | Platform/Security |

La matrice atomique, les namespaces, RBAC/ABAC, règles de step-up et SoD finales sont reportés.
