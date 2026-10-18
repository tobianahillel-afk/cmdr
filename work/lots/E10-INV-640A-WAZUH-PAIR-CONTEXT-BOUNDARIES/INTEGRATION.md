# E10-INV-640A — Correlated Wazuh-pair context, readable generic evidence, and exact boundaries

## Integrated analyst behavior

This lot finishes the *actual* two-Case investigation path introduced in PR #104: two independently SHA256-verified owner-scoped Cases containing Wazuh 4.x alerts plus normalized generic JSONL evidence, first-party Wazuh rule+agent correlation (same agent+rule within ±60s), an explicit one-file private dossier containing BOTH re-authorized Wazuh source previews plus up to ten nearest owner-authorized generic contextual event **references per Case**. Generic raw messages are never copied into that paired dossier. Entry references require their own fresh typed-source open to read the source; a matching timeline or numeric correlation is not an access grant.

**This change adds the missing end-to-end risk-boundary evidence, not an independent detector, new validation layer or invented production store.**

- Import a real feed of 16 generic events with before/after timestamps, two events outside ±60s, a symmetric tie at ±1s, another tie at ±2s. Assert deterministic top-ten *by absolute clock distance* (then earliest UTC), actual canonical imported IDs, *independent first-party typed-source open for each selected lead*, truncated=true when legitimate leads exceed ten.
- With a small independent feed, assert that exactly **−60s and +60s are INCLUDED** and −61s/+61s excluded. No off-by-one or Case catalog insertion-order shortcut.
- Save a Case with a Wazuh correlated pair and nearby generic source, then independently revoke the GENERIC typed-source reader in an owner-valid snapshot, retaining a correct new operator SHA-256. The paired Wazuh event sources remain authorized; the private correlated dossier is still produced but **MUST NEVER expose the denied generic preview, source-ref, user or contextual event ID**, and that Case context list is empty. This protects an important product distinction: a Wazuh pair authorization does not transfer read rights to adjacent generic evidence.
- Invalid Wazuh center clock and canceled request fail before producing contextual leads.
- All parent PR #104 behavioral tests (two sources, full private report, unrelated Cases, changed SHA, original generic feeds deleted, no half-dossier), Go+Node Linux and Go Windows/macOS still required.

No new IAM, ingestion daemon, ML engine, generic framework, remote API, database, custody signature, elevated owner grants or loosened CLI mode. The on-disk Case preview is sensitive local operator data, not production tenant IAM; the paired contextual event IDs are navigation hints and the near-time relation is **co-occurrence**, never causation.

## Consolidation and gating

Frozen parent PR #104 SHA `a47f48f9d943bc2d5d7285cfeed355ca956cdf1e`. Its actual Linux/macOS/Windows runtime CI succeeded but changed critical security coverage was 76.19% < mandatory 90.00%. This child deliberately targets a real previously untested authorization/regression scenario in the **same integrated product journey** and corrects NO security threshold. Rerun exact HEAD GitHub Actions; **do not merge if mandatory `engineering-kernel` remains red**, even when functional tests pass. The parent PR #103, #102 and foundational #36 remain draft. If the new delta is materially accepted, request ONE independent read-only L2, without duplicating already-completed PR #104 review or altering frozen parent SHA.
