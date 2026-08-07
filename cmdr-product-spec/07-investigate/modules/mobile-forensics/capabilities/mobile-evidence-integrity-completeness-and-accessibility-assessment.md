---
id: CAP-INV-705
title: Mobile Evidence Integrity, Completeness and Accessibility Assessment
product: investigate
module: mobile-forensics
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-07
requirement_ids: [REQ-INV-001, REQ-PROD-014, REQ-PROD-020, REQ-OBJ-004, REQ-SEC-001]
open_decisions: [OPEN-011, OPEN-013, OPEN-014]
source-of-truth: canonical
---
# CAP-INV-705 — Mobile Evidence Integrity, Completeness and Accessibility Assessment

## 1. Définition
Évaluer séparément l’intégrité, la complétude et l’accessibilité d’un Mobile Evidence Package ou d’une représentation d’extraction, sans confondre ces dimensions avec pertinence, exploitabilité ou qualification Evidence.

## 2. Problème utilisateur
Une source peut être intacte mais incomplète, complète mais inaccessible, chiffrée sans être corrompue, ou partiellement vérifiée. Un verdict unique masque ces limites et peut conduire à de fausses conclusions.

## 3. Objectifs
Distinguer `verified`, `partially-verified`, `unverified`, `incomplete`, `corrupted`, `encrypted`, `locked`, `restricted`, `unsupported`, `inaccessible`, `disputed`; examiner source/acquisition, checks disponibles, parties manquantes, fichiers invalides, transformations, contradictions et suitability for analysis.

## 4. Non-objectifs
Aucun algorithme/hash imposé, déchiffrement, unlock/bypass, qualification Evidence, interprétation de contenu, outil tiers, moteur, format ou machine d’état finale.

## 5. Propriétaire
Investigate possède les assessments analytiques. Collection/Trust conservent acquisition/custody sources; Evidence workflow qualifie Evidence; Security contrôle accès.

## 6. Utilisateurs
Principal : Evidence Reviewer. Secondaires : Mobile Forensics Analyst, DFIR Analyst, Investigation Lead et Sensitive Data Reviewer.

## 7. Conditions d’entrée
Package/représentation et Acquisition Context disponibles, source/custody refs accessibles selon permission, checks ou limitations déclarés, reviewer identifié.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Package/representation | source | material under review | oui | selected version | assessment impossible |
| Acquisition Context | CAP-INV-704 | method/custody/limits | oui | reviewed version | `unverified`/`incomplete` |
| Integrity checks/manifest | source/collector | verification inputs | non | acquisition version | `unverified` |
| Missing/invalid/encrypted areas | parser/source | coverage/access facts | non | analysis version | unknown, not complete |
| Permissions/restrictions | Security/source | accessibility boundary | oui | access time | `restricted`/`inaccessible` |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Mobile Evidence Package | source/Investigate concept | version/coverage | read |
| Acquisition Context | Investigate | method/custody/limits | read |
| Artifact/Evidence | Investigate | source/qualification context | read only; no auto qualification |
| Custody/verification events | Collection/Trust | checks, gaps, disputes | read |
| Permission/classification | Security | access boundary | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Integrity Assessment | create/review/dispute/supersede | Investigate concept | method/check/source explicit |
| Completeness Assessment | create/review/dispute/supersede | Investigate concept | missing areas explicit |
| Accessibility Assessment | create/review/dispute/supersede | Investigate concept | locked/encrypted/restricted/inaccessible distinct |
| Collection gap proposal | prepare | Collection | no acquisition execution |

## 11. Fonctionnalités
Review source, acquisition, available checks, invalid/missing parts, encryption/lock/restrictions, transformations and contradictions; compare versions; record suitability by analysis type without claiming content relevance; navigate to custody source.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspect checks/gaps | reviewer | source | 0 | read | evidence trust context | non |
| Verify/partially verify/dispute | reviewer | assessment | 2 | source+reason | versioned assessment | OPEN-013 |
| Compare versions | analyste | packages | 1 | read both | coverage diff | non |
| Mark accessibility limitation | reviewer | assessment | 2 | restriction source | explicit limitation | non |
| Prepare complementary collection | lead | gap package | 2 | gap | Collection handoff | future authority |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Validate manifest/checks | oui | oui | oui | non | deterministic validator |
| Compare coverage | oui | oui | oui | summary | table/diff |
| Identify missing/invalid items | oui | rules | oui | explanation | checklist |
| Summarize limitations | oui | structured | oui | yes | structured table |
| Qualify Evidence | humain owner workflow | no | no | interdit | Evidence Review |

## 14. États fonctionnels
Integrity/accessibility vocabulary: `unverified`, `verifying`, `verified`, `partially-verified`, `incomplete`, `corrupted`, `encrypted`, `locked`, `restricted`, `inaccessible`, `disputed`. Dimensions remain separate in storage/presentation.

## 15. États d’interface
Loading preserves source version; Empty explains absent checks; Partial lists verified/unverified regions; Error keeps successful checks; Offline allows review of existing records; Permission denied hides protected content; Stale marks superseded package/check.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Integrity Assessment | concept | CAP-INV-706..719 | verification scope/method/gaps visible |
| Completeness Assessment | concept | all Mobile analysis | missing regions explicit |
| Accessibility Assessment | concept | analysis/permissions | lock/encryption/restrictions distinct |
| Gap package | package | Collection | bounded missing source, no execution |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-INV-704 | trust review | CAP-INV-705 | package, acquisition, custody, restrictions | CAP-INV-704 |
| CAP-INV-705 | analysis allowed | CAP-INV-706..715 | assessments + limitations | Session |
| CAP-INV-705 | ambiguity relevant | CAP-INV-717 | contradictions/gaps | assessment |
| Missing area | collection need | CAP-INV-202 | gap, source, purpose, authority | assessment |
| Closure | provenance | CAP-INV-719 | checks, reviewers, versions | Session |

## 18. Dépendances
CAP-INV-213, CAP-INV-704/706..719, Evidence Trust, Security Permissions, Shared Versioning/Trace, OPEN-011/013/014.

## 19. Source de vérité
Source package and collector records remain owners of raw facts; Collection/Trust owns custody mechanisms; Investigate owns assessments; Evidence qualification remains Evidence workflow.

## 20. Provenance et audit
Record package/version, acquisition context, checks, algorithms only as declared source metadata (not selected standard), checked regions, gaps, invalid files, transformations, restrictions, permissions, reviewer, disputes, timestamps and supersession.

## 21. Permissions fonctionnelles
Package/manifest read, restricted metadata read, integrity verify/review/dispute, completeness/accessibility review, raw-region read separately, collection-gap prepare, provenance export.

## 22. Limites et erreurs
No manifest/check, corrupted subset, encrypted/locked/restricted area, contradictory source, unknown transformation, permission revocation or partial package preserves valid results and prevents complete/verified claims beyond evidence.

## 23. Métriques
Verified/partial/unverified packages, missing-area counts, restricted/inaccessible rates, disputes, re-verification after new version and analysis blocked by trust limitations.

## 24. Classification de livraison
`defined` / `planned`. No hash algorithm, format, acquisition engine, decryption tool or final Evidence model selected.

## 25. Critères d’acceptation
**Given** a package with valid checks but missing application areas **When** reviewed **Then** integrity may be verified for represented material while completeness remains incomplete.

**Given** encrypted content **When** inaccessible **Then** it is marked encrypted/inaccessible, not corrupted or absent, and no bypass is attempted.

**Given** no AI **When** review occurs **Then** deterministic checks, tables, diff and checklist provide the full assessment.

## 26. Questions ouvertes
OPEN-011 future acquisition/platform support; OPEN-013 reviewer mutations; OPEN-014 final package/Evidence/custody relations remain open.

## 27. Consommateurs documentaires
All Mobile analysis capabilities, Evidence Board, Collection, Trust, Security, Objects, Quality and Roadmap.
