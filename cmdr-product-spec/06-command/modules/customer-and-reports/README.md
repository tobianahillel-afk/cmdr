---
id: command-legacy-customer-and-reports
domain: 06-command
status: deprecated
owner: Command Product Lead
updated: 2026-08-14
source-of-truth: migration
requirements:
  - REQ-PROD-053
open_decisions:
  - OPEN-013
  - OPEN-019
---
# Customer and Reports — deprecated module path

Remplacé par la source canonique `../customers-and-delivery/README.md`, `CAP-CMD-401` et l'architecture validée `../../00-governance/adr/ADR-0008-customers-mssp-delivery-deployment-and-cross-tenant-architecture.md`.

Reporting Engine reste Shared. `CMD-CRP-001` conserve son Screen ID pour stabilité de migration mais son module fonctionnel est désormais Customers & Delivery. Ce chemin reste historique/migration et ne redéfinit ni Customer, ni Tenant, ni portfolio, ni Report.
