# E10-INV-370A — One authorized Case with Wazuh and generic JSONL evidence

## Functional integration — no new engine

An analyst can now work with a mixed offline Case built **entirely using existing CLI import/snapshot capabilities**:

1. `-snapshot <owner> -into-case <CASE> -wazuh-alerts <wazuh.jsonl> -save-wazuh-case-snapshot <private-A.json>`
2. `-snapshot <private-A.json> -into-case <CASE> -import-jsonl <generic.jsonl> -save-import-snapshot <private-mixed.json>`
3. `-snapshot <private-mixed.json> -resume-import-search-case <CASE> [-expect-snapshot-sha256 <64-hex>]`

All files are operator-selected, local, bounded, and validated through the existing strict first-party JSONL/Wazuh importers and Case/Timeline/typed-source owner. Both imports persist in the same authorized private Case projection without raw input files being needed later. No network ingestion, index, service, new schema, external library, database, or API.

**Before this lot**, the existing first-party Event Search Job treated Wazuh metadata as an all-or-nothing Case property: if `len(wazuh) != len(imported)`, numeric rule/agent pivots failed completely and all Wazuh search chronology/bursts were omitted, even if the Search Job had already independently authorized real Wazuh sources. The Case owner model already supports optional Wazuh metadata *per typed source*, so this global assumption blocked a legitimate mixed-evidence Case.

**Now**, the same authorized bounded `POST /demo-local-search` can:
- Return generic JSONL and Wazuh alert source references together when their rendered previews match an analyst's query.
- Execute rule/agent numeric pivots against Wazuh sources that actually possess typed owner metadata, while excluding generic events from those specific predicates. A Case without any Wazuh metadata still refuses numeric Wazuh pivots.
- Derive Wazuh chronological hit metadata and 60-second burst leads exclusively from the Wazuh subset of an already COMPLETED and fully owner-authorized Search Job. Generic results are never misclassified as Wazuh evidence or burst members. If a search matches only generic events, it returns no Wazuh leads.
- Preserve metadata searches by canonical `entry_id`, `source_ref`, and normalized `provenance` from parent PR #76, and separately re-open each matched typed source through the Case owner. No new authorization or correlation algorithm.

Strict fail-closed boundaries are unchanged: Case/tenant/environment/principal owner authorization, exact source owner-open **before any field comparison**, complete paginated Timeline, source identity validation, 500 in-window candidates and 200 search refs, no partial output on failure, <=4KiB authenticated JSON POST, no raw preview or bearer in response. The Wazuh burst window is the single algorithm unified from PR #51/#58 in PR #75, not a new copy.

## Acceptance (real CI behavioral evidence)

The integration test creates a real normalized Wazuh feed with three same agent/rule alerts 40 seconds apart and a separate generic normalized JSONL evidence file 10 seconds after the first alert. It invokes BOTH existing save CLI modes, producing an actual mixed private Case archive, then creates the existing authenticated handler and executes REAL first-party Event Search Jobs:
- Text query matching all four returns four separately re-openable source refs, three source-backed Wazuh chronology entries and precisely one 3-alert Wazuh burst with no generic Entry ID.
- Wazuh `rule_id` and `agent_id` selectors return only the three true Wazuh alert refs in this mixed Case.
- A generic-only query returns exactly the generic source without invented Wazuh hits/bursts.
- Wrong bearer and Case mismatch refuse data; no raw previews or token appear in the Search Job/receipt; after deleting both original feed files the persisted Case still parses.

This is **offline synthetic investigation**, not live Wazuh collection, an IAM-production tenant authorization claim, causal incident verdict, signed forensic custody or durable indexed SearchJob database. Static Wazuh metadata are first-party owner-normalized and cannot replace per-source authorizations.

## Review and merge restrictions

Child of PR #76 frozen exact SHA `8caa21ba458e0c5fa5b3e817734beee2a9502b5e` (NOT main). Preserve parent #76 and all other L2-reviewed parent SHA values; PR #36 foundation and upstream stack remain DRAFT. Exact-head Linux Go+Node/browser, macOS Go, Windows Go and unchanged mandatory `engineering-kernel` changed-security-critical coverage **>=90.00%** are required. Parent #76 was 86.39% <90%, so no automatic merge. ONE independent READ-ONLY risk-based L2 on this material mixed-data source-owner boundary after verified delivery. No waiver of the threshold or new governance framework.
