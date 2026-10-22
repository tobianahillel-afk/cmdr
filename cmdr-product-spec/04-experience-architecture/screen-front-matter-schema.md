---
id: experience-screen-frontmatter
domain: 04-experience-architecture
status: draft
owner: UX Architecture Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-UX-010
---
# Front matter des écrans


Champs obligatoires : `id`, `type: screen`, `product`, `module`, `workspace`, `status`, `owner`, `updated`, `permissions`, `source-of-truth: screen`, `requirements`.

`id` est immuable ; status suit la gouvernance ; permissions sont des références, jamais des définitions. Une vue/mode/filtre ne reçoit pas un Screen ID. Les validateurs échouent sur champ absent, statut inconnu ou permission non cataloguée.
