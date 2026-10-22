---
id: CAP-EPT-096
title: Sensitive Material, Secret Reference and Credential Handling Boundary
product: endpoint-agent
module: security
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-12
requirement_ids: [REQ-PROD-005, REQ-PROD-006, REQ-PROD-012, REQ-PROD-017, REQ-PROD-018, REQ-PROD-019, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004, REQ-SEC-005]
open_decisions: [OPEN-008]
source-of-truth: canonical
---
# CAP-EPT-096 — Sensitive Material, Secret Reference and Credential Handling Boundary

## 1. Définition
Définir la consommation locale de Secret References et les faits de manipulation/redaction de matériel sensible sans transférer à Endpoint la gestion des credentials/secrets ni exposer de raw secret.

## 2. Problème utilisateur
Des opérations Endpoint peuvent nécessiter une référence sensible; copier la valeur, la journaliser ou confondre référence et secret casserait les frontières de sécurité.

## 3. Objectifs
Represent Secret Reference, credential/protected-config refs, local sensitive-material exposure boundary, masking, unavailable/expired/invalid projection if sourced, access restriction and provenance.

## 4. Non-objectifs
Aucun secret store, key format, raw secret, credential administration, KMS/HSM implementation, rotation engine, API or physical storage.

## 5. Propriétaire
Settings/Security own secret/credential administration and policy. Endpoint owns only local reference consumption and handling-state facts.

## 6. Utilisateurs
Endpoint Operator, Security Reviewer, Platform Administrator, Auditor, Govern/Studio consumers when referenced.

## 7. Conditions d’entrée
Authorized Secret Reference or protected-config ref, tenant/environment scope, operation context, access policy and masking requirements.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Secret Reference | Settings | reference only | conditionnel | current/versioned | operation blocked/restricted |
| access policy | Security/Settings | restriction | oui for sensitive use | current | deny |
| operation context | Endpoint capability | purpose/scope | oui | exact | no use |
| validity/expiry status | Settings/source | reference state | si sourcé | current | validity unknown |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Secret Reference | Settings | id/type/status only | restricted use/read |
| Endpoint Policy | Settings | allowed reference context | read projection |
| Permission Model | Security | access constraints | reference |
| Endpoint operation | Endpoint | purpose/correlation | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Sensitive-Material Handling State | derive/record | Endpoint | never raw value |
| Masking/Restriction State | derive | Endpoint | protected fields hidden |
| Secret/credential | aucune mutation | Settings/Security owner | reference only |

## 11. Fonctionnalités
Consume reference under scope; validate availability/status when sourced; enforce redaction/no-log semantics; expose unavailable/invalid/restricted; avoid raw material in outputs/audit; correlate use by reference.

## 12. Actions utilisateur
Inspect masked reference metadata Class 0; validate reference eligibility Class 1; use authorized reference as part of an independently authorized operation only. No secret reveal/edit/rotation action.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| validate reference state | oui | oui | oui | explain metadata only | deterministic status |
| apply masking | oui | oui | oui | non nécessaire | masking rules |
| reveal/invent/change secret | non | interdit | non | interdit | Settings/Security workflow |

## 14. États fonctionnels
`reference-available`, `reference-unavailable`, `reference-invalid`, `reference-expired`, `restricted`, `masked`, `use-eligible`, `use-denied`, `status-unknown`, `stale`.

## 15. États d’interface
No Screen ID. Raw values never appear; permission denial does not reveal protected metadata.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| handling state | Endpoint fact | operation/Quality | no raw secret |
| masked reference provenance | audit context | CAP-EPT-097 | ref id/status only |
| restriction/error | technical fact | operator/Settings | sensitive details minimized |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Settings Secret Reference | authorized operation | CAP-EPT-096 | reference/status/scope | Settings ownership retained |
| CAP-EPT-096 | eligible use | target Endpoint capability | opaque ref only | handling state retained |
| invalid/restricted | failure | Settings/Security consumer | ref/status/reason metadata | no secret content |

## 18. Dépendances
Endpoint secret-protection/secure-storage; Settings secret management; Security secrets/key/privacy; CAP-EPT-085/094/095/097/098; OPEN-008.

## 19. Source de vérité
Settings/Security source owns secrets and credentials; Endpoint is SOT only for local handling/use facts tied to opaque references.

## 20. Provenance et audit
Reference id/type/status, operation/Agent/tenant, access decision, masking state, timestamps and errors—never raw material.

## 21. Permissions fonctionnelles
Secret Reference use, masked metadata read, sensitive-state read, cross-tenant deny; no raw-secret read or final RBAC introduced.

## 22. Limites et erreurs
Secret Reference != raw secret; unavailable reference != authentication failure automatically; masked state != secure-storage proof; no raw credential logging.

## 23. Métriques
Unavailable/invalid/expired/restricted references, masking violations target zero, denied use and provenance completeness.

## 24. Classification de livraison
`draft / defined / planned`; no secret store/KMS/credential implementation.

## 25. Critères d’acceptation
**Given** an operation references a secret, **When** local audit is produced, **Then** only the opaque reference and handling status are retained.

**Given** a reference is unavailable, **When** operation eligibility is checked, **Then** the operation is blocked/restricted without inventing a credential failure.

**Given** AI is unavailable, **When** masking/status checks run, **Then** deterministic controls still work.

## 26. Questions ouvertes
OPEN-008 remains open; platform-specific secret storage/access mechanisms are not selected.

## 27. Consommateurs documentaires
EPT-6 Security/Update, Settings, Security Architecture, Studio/Govern references, Quality, registers and Roadmap.