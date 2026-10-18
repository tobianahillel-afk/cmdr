# E10-INV-730A — Consolidate real source-authorized Wazuh pair discovery, paging and private drilldown

## Why this is an integration, not a new correlation engine

As of 18 October 2026, PR #114 had advanced from `9f43ceab9f96ac7648429f2361bf8b4324861786` to `80d3cf81e153935000f25c0c25ddad7853f83eca`: it added **source-authorized exact first-Entry pair discovery** beyond the original top-ten display. Both of its latest SHA workflows are green. Source PR #115 (`b261a400e2bcb53dba00922708dbd59c24893449`) developed **bounded page discovery of later paired Wazuh event leads** against the **older** #114 SHA. The GitHub PR #115 was `mergeable_state=dirty` relative to the updated #114 branch, despite its own exact-SHA workflows being green. Do not retarget the L2-frozen #115 head, overwrite either source PR or force a merge.

This new child branch starts at **exact latest green #114 HEAD** `80d3cf81e153935000f25c0c25ddad7853f83eca` and integrates the actual tested #115 algorithm/page API + tests + source note as a coherent product increment. This removes the stale-parent conflict without ever changing a reviewed ancestor SHA. Unlike a cherry-pick of the old `main.go`, its CLI integration retains **both** modern exact first-Entry discovery and paginated list mode.

## Operator journey, all first-party and offline

1. Import Wazuh 4.x events into **two independent operator-selected Case snapshots**, preserving existing original Entry IDs, canonical typed-source owners and operator-recorded file SHA-256 values. Both individual owners are revalidated against the caller principal, tenant, environment and source ACL before any metadata.
2. Run the existing `-compare-saved-case CASE_A ... -compare-saved-with-case CASE_B` with both snapshots and two expected SHA-256 values. Default metadata-only summary still presents **at most ten** earliest authorized 60-second rule+agent matched Wazuh leads.
3. Add `-list-saved-wazuh-pair-offset 0`, then `10`, `20`, and so on (up to bounded 5000 eligible offset) to retrieve the next ten or fewer source-linked pairs and a deterministic `has_more` flag, **without exposing raw preview, emitting any bearer or opening a listener**. An offset applies to authorized actually matched pairs only; it cannot be used to enumerate denied source entries.
4. Optionally resolve one Entry with `-find-saved-pair-first-entry FIRST_ENTRY`, including an Entry past the top ten, to recover its exact nearest compatible second Entry with the already present first-party 60-second matching rules.
5. For sensitive source inspection, perform a **SEPARATE** explicit request with `-inspect-saved-pair-first-entry FIRST -inspect-saved-pair-second-entry SECOND -private-report /absolute/new-file.json`. The original typed sources are independently reauthorized again and written only to the existing private, no-overwrite 0600 artifact. IDs and correlation timestamps are navigation hints, **not authorization** or causal attribution.

**Important integration guard:** page, exact-ID discovery and explicit private inspection are mutually exclusive. If a caller supplies a page flag together with either other mode, refuse the entire operation **before** any metadata or private output can be emitted. The old #115 snapshot did not know about #114's later discovery mode; resolving this conflict is the material R2 integration fix.

## Verification criteria / limits

- End-to-end with two independently saved synthetic Wazuh Cases (12 matching first Case alerts and one second Case alert), with independent matching SHA-256 values: page offset 0 yields 10 and `has_more`, offset 10 yields 2, offset 12 empty. Exact-ID discovery on an Entry from page 10 must return the SAME authorized second Entry, and a separate source-authorized private drilldown must reopen both raw Wazuh source descriptions without stdout preview or overwriting anything.
- Mixed page+find or page+private-inspection flags fail closed, including a **valid** later matched Entry ID; no accidental file, metadata, preview or bearer. Denied principals, missing Cases, wrong digests, tampered but parseable second snapshot, negative/oversized offsets and forbidden generic-preview flags also refuse without partial disclosure.
- The implementation reuses `sharedWazuhPairsPageFromObservations` from source #115 and `writeVerifiedSelectedWazuhPairLead` from #114; does not create a second matcher or parallel permission graph.
- This remains the synthetic operator-local first-party lab: no production IAM, streaming ingestion, signed chain of custody, mutable Case DB, cross-tenant access or paid API. A 60-second nearest-neighbor pair is temporal coincidence, not proven causality.

## Branches, checks, merging and L2

Base: frozen #114 `80d3cf81e153935000f25c0c25ddad7853f83eca`. Source #115 `b261a400e2bcb53dba00922708dbd59c24893449` remains untouched. The enormous #36 foundational draft and the entire earlier draft dependency chain are NOT promoted to main. Exact final HEAD Linux Go+Node browser, macOS and Windows plus mandatory `engineering-kernel` **changed-security-critical 90% floor** are prerequisites. No waiver, no merge while any required gate or parent L2 is unresolved. Request **ONE** independent read-only L2 on this coherent integration's exact SHA and diff only after material stable GitHub CI evidence. Do not re-review all parents or manufacture a second-identity GitHub native approval.
