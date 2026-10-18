# E10-INV-690A — One full-stack owner-authorized DFIR integration candidate

## Actual GitHub source and scope

This is **integration of existing FUNCTIONAL CODE**, not a new SIEM engine, new persistence store or another duplicated test framework. Native GitHub compare, checked before branch creation, shows immutable PR #92 HEAD `fced2c5c89f0458247be9cc0644e516118af2ae6` → exact existing PR #109 HEAD `dbbeb85606d98057643daf306d962228db7aa32b` = **71 REAL pre-existing published commits ahead, 0 behind, 35 changed paths**. This integration branch was made directly at that already published SHA, preserving every reviewed parent SHA. This note is one additional **documentation-only** commit, with no runtime code modified.

This single integration PR reduces the required staged sequence after PR #92: full #93–#109 runtime code is reviewed/tested as ONE integrated diff, instead of requiring the separate ancestor #108 and #109 PRs to be merged in sequence. Original source PRs #93–#109 remain untouched, open and independently reviewable. Upstream PR #92, #81, foundational #36 and documentation #2 remain DRAFT/UNMERGED; this does **not** make the product delivered to `main`.

## Existing executable laboratory analyst journey

1. Operator selects local Wazuh 4.x JSONL and generic normalized JSONL feeds and imports them into existing synthetic, tenant/environment scoped Case/Timeline/typed-source owner graphs. Snapshot save is explicit, bounded, private, mode 0600 POSIX, atomic no-overwrite; the selected Case is minimized to readable entries and typed sources. No network ingestion daemon, IAM production or trusted signed chain-of-custody.
2. The analyst independently retains two SHA-256 digests of two saved Case snapshots. Before comparison, the first-party loader checks **the exact bytes later decoded**, independently validates principal, tenant, environment, Case, complete Timeline pages and typed sources for EACH owner.
3. Existing Event Search Jobs search authorized normalized feeds. Existing first-party Wazuh cohort collection derives 5-minute rule+agent windows and actual event pairs with same rule+agent and inclusive <=60 seconds UTC. These are retrospective **temporal co-occurrences**, not malware attribution or causation.
4. Existing private paired evidence CLI takes an explicitly selected pair ordinal (never an access grant) and a NEW destination, independently opens the selected sources from each Case with first-party authorization. Writes ONE private no-overwrite JSON dossier, with no raw preview, Case bearer or token on stdout.
5. Default **v2** includes Wazuh pair previews in private file and up to ten authorized generic context **metadata-only** events +/-60s from each side. Explicit opt-in **v3** includes independently reauthorized generic source previews in the same private file; no preview without explicit private report selection. Latest #109 reuses the SAME first-party generic owner inspection per selected context lead to avoid opening it twice in v3. The first and second owner/source boundaries remain independent.
6. Source revocation, wrong digest, denied principal/tenant/environment/Case, conflicting evidence or existing output destination MUST cause fail-closed behavior with no partial published dossier; tests must demonstrate this for BOTH owner sides.

## Integration acceptance and risk

- Run the EXACT integration PR HEAD `investigation-timeline-runtime` Linux Go+Node analyst browser, macOS and Windows Go jobs, plus mandatory `engineering-kernel` CHK-SECURITY-TESTS at its **UNCHANGED >=90% changed security-critical coverage** gate. Previous exact-source CI was runtime successful for PR #108 and #109, but **NOT equivalent** to this new consolidated full diff.
- Historical source CI: PR #108 exact HEAD `c5249657ac35ba2a718fc451ca85a1e9e4428c6c` runtime SUCCESS run 39163283597; security coverage FAIL 89.79% <90% run 39163283681. PR #109 exact HEAD `dbbeb85606d98057643daf306d962228db7aa32b` runtime SUCCESS run 39164604397; security coverage FAIL 87.76% <90% run 39164604386. Never transfer green from source PRs or lower/waive the gate. Upstream #92's own gates/reviews are still a separate prerequisite.
- Before any merge, check and close real R0/R1 source or tenant access defects, R2 functional data correctness, verify source and integrated L2 findings, mandatory CI and NATIVE branch protection/independent GitHub reviews. A ChatGPT scheduled L2 is **not** a native GitHub second reviewer.
- Do not add generic abstractions or nearly identical edge-case tests solely to manipulate a coverage percentage. Any genuine remaining sensitive path uncovered should be improved through a focused design correction or meaningful behavior test; after two unsuccessful attempts reconsider design.
- No live production ingestion, durable transactional Case store, real IAM, cryptographic custody, 24/7 SIEM, event-index service or GPU capability is claimed. Decision `OPEN-013` concerning mutation authorization remains open under canonical product specifications.
- No new paid APIs, dependencies, synthetic security pass claims or manual computer requirements.

## Next functional step once integration is healthy

Treat this full Case-local lab journey as a coherent **operator-ready testable milestone**, then decide the first genuine production boundary (canonical persistent Case/typed-source data store and tenant security adapter) from `cmdr-product-spec/**` before attempting any live ingestion or multitenant deployment. Do NOT turn this integration note into a new architecture or governance policy.
