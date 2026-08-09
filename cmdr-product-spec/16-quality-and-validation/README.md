---
id: quality-readme
domain: 16-quality-and-validation
status: draft
owner: Quality Lead
updated: 2026-08-09
source-of-truth: canonical
---
# Quality and Validation

## Objectif

Définir les contrôles qui comparent l’arborescence réelle au manifeste attendu et valident liens, noms, front matter, sources, permissions, propriétaires, métriques et preuves de publication.

## Périmètre

Quality records evidence and verification stages. It does not own product behavior, objects, permissions, roadmap decisions, technology choices or implementation.

## Active evidence

- Validation status: `validation-status.md`.
- Structural baseline: `architecture-manifest.md` and `expected-path-manifest.md`.
- Command current revalidation: `reports/phase-4a-command-current-revalidation.md` — **PASS AFTER POST-PUBLICATION VERIFICATION, 60/60**.
- Phase 4B Investigate closure: `reports/phase-4b-investigate-capability-closure.md` — **PASS**.
- Cloud Analysis: `reports/phase-4b4a-cloud-analysis-capability-conformance.md` — **243/243 PASS**.
- Mobile Forensics: `reports/phase-4b4b-mobile-forensics-capability-conformance.md` — **250/250 PASS**.
- Govern GOV-1 dedicated verification: `reports/govern-gov1-post-publication-verification.md` — **180/180 PASS**.
- **Govern GOV-2:** `reports/govern-gov2-playbooks-response-runs-verification-rollback-capability-conformance.md` — pre-publication remote gates remain PENDING until the fifth functional commit is published and verified.

## GOV-2 quality contract

Parent roadmap: **Delivery Roadmap Phase 4 — Govern** (`roadmap-phase-4-govern`). Execution lot: **GOV-2 — Playbooks, Response Runs, Execution, Verification and Rollback**. GOV-2 is not a roadmap phase. `Phase 4C Govern` does not exist.

Required GOV-2 structural evidence:
- 17 unique `CAP-GOV-017..033` capability files;
- 459 numbered sections;
- 102 mandatory S8/S9/S10/S13/S16/S17 tables;
- zero empty/generic mandatory tables;
- zero duplicate/recycled IDs or owner conflicts;
- 190/190 gates only after remote verification;
- zero GOV-3 capability, detailed screen rewrite, provider/runtime selection, API/protocol, command or product implementation.

## Non-regression contract

GOV-2 closure must preserve:
- GOV-1 `CAP-GOV-001..016`, 16/432/96 and historical 180/180 PASS;
- Command 27 capabilities, 26 defined + 1 proposed, five Requirements ranges and `DEP-CMD-001..010`;
- Investigate 243 capabilities and Phase 4B PASS;
- Requirements 122 = 99 conform / 20 partial / 3 absent / 0 contradictory unless independently justified;
- 18 OPEN decisions unless an independently resolved source is found;
- root README and `main` unchanged;
- PR #2 open, Draft, unmerged.

## Validation rules

- planned capability != implementation evidence;
- technical executor output != canonical Result;
- post-publication PASS cannot be claimed before exact remote head, commit chain, PR, README, `main`, structural counts and non-regression are verified;
- historical evidence is preserved or extended additively rather than destructively condensed;
- a failed/pending gate remains visible; one FAIL makes GOV-2 PARTIAL;
- remote-dependent gates stay PENDING until the fifth functional publication and remote check;
- no-AI paths, source ownership, target/scope/expiry and sensitive-data boundaries are quality requirements, not optional prose.

## Counts prepared for final verification

- global capabilities: 303;
- Command / Investigate / Govern: 27 / 243 / 33;
- defined / proposed / planned: 301 / 2 / 303;
- GOV-2: 17 / 459 sections / 102 tables;
- Govern cumulative: 33 / 891 / 198;
- all products currently specified: 8181 sections / 1818 mandatory tables;
- GOV-3 capabilities: 0.

## Permissions / screens / implementation

GOV-2 validates functional permission needs only; final RBAC/ABAC remains future. Existing nine Govern screens remain; detailed rewrites/new Screen IDs = 0. No executable command, exploit/bypass, API/protocol, raw secret, provider/runtime or product code is produced.

## OPEN

OPEN-007, OPEN-008, OPEN-013, OPEN-015 and OPEN-019 remain explicitly open. GOV-2 creates/closes no decision.