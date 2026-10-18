# E10-INV-450A — Private atomic forensic receipts from three existing analyst workflows

## Functional integration

Existing CMDR synthetic operator-local workflows can **verify** one authorized saved Case's typed source inventory, **hunt** a saved mixed Wazuh/generic JSONL Case using the real first-party Event Search Job with temporal leads, and **compare** two separately saved and independently SHA-256-verified Cases. Each previously printed an authorized JSON receipt to stdout. This is inconvenient for durable analyst handoffs and risks collecting sensitive source reference IDs in terminal transcripts/CI logs.

This delivery adds **one optional `-private-report /absolute/NEW/report.json` flag**, supported exclusively in those three already authorized modes. The existing owner/IAM simulation, source-open checks, Search Job execution, temporal correlation and compare calculation are unchanged. Every workflow runs COMPLETELY to an in-memory bounded JSON receipt **before** the existing private atomic publisher writes anything. Failure to authorize either Case, validate either independently retained digest, execute the full search, or serialize valid bounded JSON means **no output file, no receipt on stdout and no partial artifact**.

Reuse `writePrivateSnapshotAtomic` from the existing saved Case storage path: private 0600 temporary file (POSIX), sync and close, atomic hard link to a **new destination only**, unlink temp, refuse existing file/symlink/hardlink and unsupported filesystem publication rather than overwriting evidence. Windows ACL privacy depends on operator directory/OS. Destination is an explicit absolute path; no server endpoint can choose it.

### Operator CLI examples

Add the following flag to any of the existing commands:

```text
-private-report /absolute/new-only/private-hunt.json
```

- `-verify-import-snapshot-case CASE -snapshot CASE.json -expect-snapshot-sha256 64HEX -principal ANALYST`
- `-hunt-saved-case CASE -snapshot CASE.json -expect-snapshot-sha256 64HEX -hunt-field source_ref -hunt-value local-jsonl-source- -hunt-from RFC3339 -hunt-to RFC3339 -principal ANALYST`
- `-compare-saved-case CASE_A -snapshot A.json -expect-snapshot-sha256 A64HEX -compare-saved-with-case CASE_B -compare-saved-with-snapshot B.json -expect-other-snapshot-sha256 B64HEX -principal ANALYST`

When `-private-report` is omitted, existing normal stdout behavior remains unchanged. If supplied, stdout is deliberately empty even after success. The operator already knows the output path. The on-disk report may contain **sensitive authorized event/reference identifiers**; it is not a complete evidence archive and must not be posted to public GitHub issues or logs.

### Acceptance and security

- End-to-end real **private output** for all three previously implemented workflows, including Wazuh↔generic mixed-source Search Job and independent two-Case timeline comparison.
- Output must be complete and parse as one valid JSON document with the existing receipt format; source previews, hidden owner audit state, browser bearer and listener information are never included.
- Each output file is mode 0600 on POSIX, newly created and never overwritten. A second request with the same path fails and leaves original bytes unchanged. A failed principal, wrong hash, second unauthorized Case or incorrect mode produces no output and no private temporary artifact.
- Existing ACL, tenant, environment, source and SHA-256 checks are the ONLY sources of authority. This feature **does not** add a new IAM/authorization wrapper, a new search engine, a database, a paid API, a live collector or an authenticated forensic custody signature.
- Local synthetic operator-controlled Case owners ONLY. A SHA-256 verified against an independently retained operator value detects unexpected byte changes, not an adversary who can replace both file and expected value.

## GitHub disposition

Base is PR #85 frozen SHA `5246c03653c627f026d940378419a57a59dec1a5`, **NOT main**. Preserve prior L2-frozen SHA, and treat the 90% changed-security-critical coverage gate, independent review and parent approvals as mandatory. Never merge while the engine gate or any required parent approval is red. Parent PR #85 is known to have functional Actions SUCCESS with `CHK-SECURITY-TESTS` at **83.87% <90%**. This new integration must be assessed at its own exact published HEAD before assigning any green status. After functional and security CI stabilize, schedule **one** independent read-only L2 because the new private filesystem output is a material confidentiality boundary. Do not request duplicate L2 per test commit.
