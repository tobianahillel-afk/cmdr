# E10-INV-200A — Authorized literal hunt over rendered Wazuh alert previews

## Actual functional increment

From the existing Wazuh 4.x local analyst viewer, a permitted operator can now select **Text in authorized preview**, enter a literal UTF-8 phrase of 3–128 bytes and precise `[from,to)` UTC search time, then **execute an actual CMDR first-party Event Search Job**. This is not merely another metadata/agent/rule cohort filter: it evaluates the owner-authorized source's **rendered alert summary text** (Wazuh rule description, or `full_log` fallback only when no rule description exists). The search ignores case and treats the term literally, not as a regex/DSL.

The HTTP path is the **existing** bearer-guarded `GET /demo-wazuh-search`, activated only in Wazuh import mode. The existing Event Search Orchestrator and its tenant/environment/source/time validation and audit remain in control. For each Wazuh Timeline entry within the time range, the backend first reauthorizes the **typed owner source-open**, and only then compares its bounded rendered Preview. Therefore:
- No raw JSONL source re-read, no raw log backend or direct full_log inspection from an untrusted request;
- No Wazuh metadata is trusted as an ACL, and cross-tenant, denied source, revoked bearer or partial owner traversal fails closed with no search output;
- Existing scope: **maximum 500 owner-authorized in-window candidate previews** and **maximum 200 SearchJob hits**, both strict complete-query limits. Overflow fails the entire job, never silently truncates it;
- Query results disclose only stable authorized source IDs, normalized rule/agent metadata, source UTC time and, if present, existing descriptive burst leads. They never copy rendered preview text or the bearer into the Job/metadata receipts; the existing separate source-navigation owner must authorize any explicit preview opening.
- QueryVersion distinctly records `wazuh-rendered-preview-literal-v1`, not the existing rule/agent pivot version. The literal input is not sent to raw source audit logs by the backend.
- The browser **does not enable JSON replay export for text hunts** to avoid storing potentially sensitive literal terms in local files. Existing rule/agent search replay remains supported.

## Tests and boundaries

Real disk-to-Case-to-authorized-source-to-Event Search Job scenario over two Wazuh feeds with original alert ID collisions, case-insensitive text match across both files, negative ordinary alert exclusion, empty completed nonmatch job, no source preview leakage, denial of foreign Case, missing/overlong/newline-containing literal, invalid scope/time, unsupported full_log selectors and unauthenticated bearer. Browser scenario verifies correct selected literal in request, search hits appear and no replay export is enabled. Existing Go/node tests Linux/macOS/Windows must remain green.

**Important exclusions:** laboratory-local **synthetic identity** and snapshot authorizations, not production IAM. Not a continuous Wazuh connector, full-text inverted index, raw log search, regex/DSL, durable Search Job, malware evidence/chain-of-custody storage, or signed forensic report. Rendering is truncated to the existing source preview; this capability cannot claim complete raw event content coverage. No paid API or new external dependencies.

## Technical choice and integration

Minimal amendment to existing first-party `wazuhPivotBackend` rather than a new provider, wrapper, validator pipeline, or search engine. Metadata filters still filter before typed-source open; the new preview mode opens the source **before** literal matching. One backend executes each Job; existing Case owner and source ACL remain authoritative. The parent PR #58 immutable L2 SHA `315f3019209f2cdae27ccbc013bddcfa466891e1` and all earlier reviewed SHAs stay unchanged. PR remains draft until exact-HEAD `investigation-timeline-runtime`, unchanged mandatory `engineering-kernel` coverage floor 90% and one independent read-only L2; never use this increment to waive prior coverage defects or force upstream merges.
