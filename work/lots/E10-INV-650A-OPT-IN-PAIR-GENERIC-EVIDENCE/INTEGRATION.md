# E10-INV-650A — One consent-gated private Wazuh pair + generic evidence dossier

## Real integrated analyst journey

The existing first-party cross-Case workflow independently verifies two saved Wazuh+normalized-JSONL Case owners, correlates Wazuh rule/agent events with an inclusive 60-second window, and opens the pair's **two sensitive Wazuh source previews** into a private atomic no-overwrite report. PR #105 added up to ten nearby authorized generic JSONL contextual **IDs and time offsets per Case**, but intentionally omitted their preview payloads. The analyst had to restart source inspection independently for each contextual lead.

This lot integrates those already existing components into ONE operator-consented evidence artifact; it does **not** implement a new detector, correlation algorithm, IAM, store or generic provider.

**Existing default unchanged:** CLI `-inspect-saved-pair 1 -private-report /absolute/new.json` continues to generate **v2**, containing both Wazuh previews and *metadata-only* generic context links. No new source disclosures.

**New explicit request:** Append `-include-pair-generic-evidence` to the SAME explicit `-inspect-saved-pair` and `-private-report /absolute/NEW.json` command. Only then write format **`cmdr-private-verified-wazuh-pair-v3`**, retaining v2 fields and including `first_generic_evidence` and `second_generic_evidence` as arrays of actual first-party authorized typed-source inspection receipts (maximum ten per selected Case). Their Entry IDs match the corresponding shortlisted `first_generic_context.entries` and `second_generic_context.entries` in order. No generic evidence beyond the existing inclusive ±60-second, observed certain-time / canonical imported / telemetry-event selection is included.

Every contextual ID is strictly a selector: re-open it through the **original Case's** `inspectDemoEntryWithData` first-party principal/Case/Timeline/typed-source owner. Source ACL denial or inconsistent return aborts before ANY private file is published. Each original Case snapshot's bytes and operator-held SHA-256, tenant/environment, principal and source permissions were independently verified by the existing loaders before this expansion. Wazuh pair correlation itself is unchanged. No multi-tenant owner merge or cross-Case permission reuse. Raw contextual text is sensitive and is written **ONLY** to the same operator-selected private, exclusive-create/no-overwrite, bounded, 0600-on-POSIX artifact; no stdout or browser path.

The full report can contain personal/sensitive event details; the operator must explicitly choose `-include-pair-generic-evidence`. This is a **local synthetic forensic convenience**, not authenticated source authenticity, a custody signature, a deployed SIEM, production IAM or proof of causal association. The source previews come from operator-controlled normalized local fixture owners, not continuous production ingestion. SHA-256 is an unsigned unexpected-change check only.

### Behavioral acceptance

- Construct independently saved Wazuh Case owners containing real generic JSONL evidence at +10 seconds in first, +2 seconds in second, and a +4-minute out-of-window decoy. Confirm SHA-256 and Case/typed-source authorization, then run ACTUAL CLI opt-in. v3 private output includes the two added sensitive nearby generic source previews, independently matching each Case and existing entry IDs, and excludes far preview, stdout bearer and loopback listener.
- Run same CLI WITHOUT opt-in: output remains v2, and no generic previews appear. This compatibility is deliberately retained rather than silently broadening disclosure.
- Using inclusion without an explicit valid 1..10 pair or without a new private destination MUST fail closed without stdout or file creation.
- Foreign principal, wrong Case and arbitrary Entry ID cannot reopen context through owner, even after earlier correlation. An altered second snapshot with an old SHA must fail before any first-side report appears.
- Maximum 10 per Case, existing source owner typed checks, no new dependencies or paid APIs. GitHub Actions need actual Go runtime Linux/macOS/Windows and browser tests plus mandatory security test coverage >=90% (NEVER lower/waive).
- Existing foundational PR #36 and upstream integration drafts are not production merged; canonical decisions remain open.

## Integration

Parent exact frozen PR #105 SHA `c41159123f4e1f16163b0f26c2e46ece867b7b92`, target new DRAFT child PR. One independent risk-based L2 only after stable published SHA/CI because *explicit preview expansion across separately owned Cases* is a substantive sensitive-data boundary. Do not edit the parent L2-frozen SHA or merge while mandatory checks are red.