# E10-INV-120A-SNAPSHOT — multi-feed Wazuh cohorts in analyst UI

This functional increment extends PR #49, **not** the oversized foundation PR #36. It turns the existing multi-file offline Wazuh import from a Case-only timeline into an integrated `/demo-wazuh-cohorts` browser investigation path for agent and rule counts.

## Reuse-first technical choice

Reuse the existing Wazuh 4.x parser, bounded 2–4 operator-local feed loader, in-memory Case/Timeline/typed-source owner, bearer-protected loopback route and single-feed rule/agent aggregate engine. The only necessary extension is to bind rule metadata to **feed-prefixed, normalized source IDs**, and provide that metadata map to the same existing authorization-first aggregator. No Elasticsearch/OpenSearch, new provider, index, custom auth model, queue, paid API or persistence is added.

## End-to-end behavior

Operator selects authorized snapshot and two distinct Wazuh files; Wazuh IDs may collide across files. After all files are validated and imported without mutation, the authenticated analyst viewer can request `/demo-wazuh-cohorts?case_id=<selected-case>`. The route checks the transient bearer and exact Case selector before the evaluator re-reads the authorized Case pages and typed source objects. The aggregate contains only allowed Wazuh agent IDs, rule IDs, numeric severities and time windows; never hostnames, private raw descriptions, or the bearer. The existing `/case-timeline` and `/demo-source-open` continue to work for selected entries.

The complete two-feed loopback test checks rule 5710 across both feeds, distinct agents, maximum numeric severity, denied bearer, foreign Case and duplicate selectors, source preview authorization and unchanged original files.

## Limits and security

**Offline analyst laboratory**, not production XDR/SIEM ingestion. Local operator-selected files, ephemeral loopback bearer, caller-selected memory provider, no persistent durable custody, no provider IAM integration, no real-time upstream Wazuh subscription. Metadata never grants Case, Timeline or source permission; only already-authorized entries contribute. `OPEN-013` remains unresolved.

Do not mark the lot VERIFIED or claim CI/L2 success until exact GitHub HEAD evidence confirms it. PR #49 had engineering-kernel FAILED at `ffa40a468e4eda5706963d0c134c46d2eb1cd82c` because changed security-critical coverage was 88.36% below the 90% gate; preserve that evidence, and reevaluate on the new SHA.
