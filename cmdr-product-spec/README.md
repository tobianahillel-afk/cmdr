---
id: root-readme
domain: repository
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: navigation
---
# CMDR Product Specification

Ce répertoire est la source documentaire canonique de CMDR. Les produits canoniques sont **Command**, **Investigate**, **Govern**, **CMDR Studio**, **Platform Settings** et **Endpoint Agent**. Les intitulés marketing *Command Center*, *Investigation Lab* et *Response & Governance* restent autorisés dans l’interface, sans remplacer les noms de produits.

## Démarrage

- [Index complet](INDEX.md)
- [Statut](STATUS.md)
- [Historique](CHANGELOG.md)
- [Politique de source de vérité](00-governance/source-of-truth-policy.md)
- [Manifeste attendu](16-quality-and-validation/expected-path-manifest.md)
- [Registre de propriété](00-governance/ownership-register.md)

## Chaîne canonique

`Telemetry Event → Detection → Signal → Alert → Incident → Case → Evidence → Finding → Action Request → Decision → Response Run → Result`

Chaque objet possède une source unique dans `05-domain-model/objects/`. Les permissions vivent dans `14-security-permissions-and-trust/permission-model.md`, l’Inspector dans `03-design-system/components/inspector.md`, et les règles Page/Vue/Mode/Filtre dans `04-experience-architecture/page-view-mode-filter-rules.md`.
