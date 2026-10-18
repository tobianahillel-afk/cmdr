# E10-INV-710A — Single full-scope GitHub DFIR integration candidate

## Why this is a release-oriented integration rather than another micro-feature

The source HEAD `4d6c8057557a68657b9b64584c4f41bac17b3868` is the **already-published** functional PR #112 (exact Wazuh pair selection by two canonical Entry IDs). Native GitHub commit compare against the **older, immutable green foundation PR #81** `c173c34dbe224b6427cfb37ed880da9266504e1c` reports **113 existing real commits ahead / zero behind / 62 changed files**. This integration branch starts at #112's exact immutable HEAD rather than rewriting/cherry-picking previously audited commits. The ancestor L2-reviewed PR SHAs are not touched. This document adds only a delivery-oriented handoff decision, not another duplicate runtime algorithm.

## Real integrated local DFIR capabilities already in this GitHub diff

- Operator-selected OFFLINE, synthetic Wazuh 4.x and normalized generic JSONL feeds (1–4) feed a first-party Case/Timeline/source-owner graph, scoped by principal, tenant, environment, Case, readable Entry and typed Source. Real production multi-tenant IAM and streaming ingestion are NOT implemented.
- Explicit private Case-only 0600/no-overwrite snapshot write; exact-byte SHA-256 verification using a separately retained expected digest; authorized restart and offline JSON-only verification; no inference that SHA-256 is signed forensic custody.
- Authenticated Case Event Search Jobs with full authorized Timeline/source traversal and bounded, source-reopened generic text matching; source IDs retain owner provenance.
- Independent Case A / Case B authorization and chronology comparison. Existing actual Wazuh cross-Case rule/agent/time correlation, contextual generic events in bounded +/-60s windows, headless source-owner reauthorization, default private metadata-only dossier v2 and explicit private v3 generic previews; no public raw evidence / bearer in stdout.
- Most recent functional enhancement from PR #112: source-grounded analyst selection of a Wazuh pair by **both exact owner-authorized Entry IDs**, not an unstable ordinal, with re-computation of currently eligible pairs on re-open. This keeps the existing Case/source/SHA256 boundary and private no-replace writer. Source owner ACL must be checked AGAIN on every inspection.

## Full-scope verification and merge policy

- This PR's **base is #81**, not #110 and not main. Security checks must inspect this complete 113-commit functional change set, including older changes to normalized import, offline ownership, Event Search Job, Wazuh association, private dossiers, and cross-Case source authorization.
- The green small-diff CI of PR #110 and #112 **does not certify** the broad #81..#112 integration boundary. Prior PR #111 already attempted a similar broad scope and failed mandatory `CHK-SECURITY-TESTS` **89.27% < 90%** despite runtime Go+Node passing.
- Mandatory **90.00% changed security-critical code coverage**, first-party risk review, all native GitHub required checks/approvals, ancestor status, and explicit absence of R0/R1/R2 vulnerabilities are NON-WAIVABLE. A green security check here does NOT automatically clear any other unmerged ancestor checks or permit merging a chain of draft PRs.
- Runtime must be tested on Linux/Windows/macOS and browser JS as in the existing GitHub Actions. This branch should be treated as ONE auditable deliverable rather than another sequence of micro-PRs. No second-identity GitHub review can be simulated by a scheduled ChatGPT L2.
- Do not force merge this PR or claim readiness for production: local synthetic DFIR lacks durable transactional indexing, tenant IAM integration, live SIEM connectors/ingestion, trusted custody, and real SOC operational SLOs.

## Technical decision

Avoid a new detector/search framework, a second persistence engine, multiple new ACL adapters, or the creation of tests solely to influence coverage. If this full-scope CI fails, inspect its actual uncovered critical runtime statements and fix a *material* owner-access or evidence-integrity path rather than lowering gates. Previous source PRs, especially #111's earlier full-scope attempt, should not be blindly merged. Consolidating GitHub diffs is a review technique, NOT evidence that software is deployed.
