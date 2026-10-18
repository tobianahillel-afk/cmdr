# E10-INV-310A — Close legacy generic JSONL Case export cross-tenant leak

## Real R1 vulnerability and existing fix reused

The original `-snapshot OWNER -import-jsonl EVENTS -into-case CASE` CLI path, `writeDemoImportedJSONL`, still marshaled **the entire imported `MemoryTimelineData`** to stdout. Although `importLocalEvents` verifies that the selected principal can import into the named Case and HTTP navigation applies tenant/Case/typed-source ACLs, the original owner snapshot can contain **other tenants/Cases and same-Case source records explicitly denied to that principal**. Any terminal redirection or downstream consumer of stdout receives those bytes regardless of HTTP authorization. This is a substantive R1 confidentiality boundary, not a cosmetic/code-coverage issue.

The existing newer `scopedNormalizedCaseSnapshot` already solves this exact authorized-snapshot projection in the explicit private Case save and verified restore flows. Instead of adding another validator/provider, the legacy import path now invokes the SAME first-party owner projection **after complete validated import but before encoding any output**:
- one selected principal with tenant+environment checked against the selected Case;
- exactly that Case and its readable Timeline catalog entries;
- only typed-source records independently readable by the principal;
- no other principals, tenants, Cases, denied entry/source identifiers, source previews or other metadata;
- complete revalidation by the existing owner constructor before marshaling.

The format remains the existing reusable `MemoryTimelineData` JSON with the same imported stable Case/Entry/source IDs, timestamps and the principal's authorized source previews. The output is **not redacted from the selected operator's legitimate Case data**, but no longer constitutes a raw full-owner export. No new search engine, ACL model, database, dependency, API or material UI. Normalized JSONL import, online viewer, Wazuh Case imports and existing private no-overwrite saves remain unchanged.

## Verified behavior and limitations

The original disk-to-import-to-Timeline/source drilldown test is upgraded: output must contain precisely the selected principal and Case, four readable entries (two original readable Case records and two new normalized JSONL records), and **must not contain any known forbidden record of the fixture**, including another tenant's Case or denied same-Case source. The exported reduced owner snapshot still reloads and passes actual authenticated Timeline/source navigation; the original owner snapshot remains byte-for-byte unchanged. The existing CLI import test now verifies the restricted owner shape rather than assuming it copied every principal's records. Invalid file/Case/principal, cancellation, writer failure and identity-conflict tests continue.

This improves **one specific synthetic local operator export**; it does not turn the current in-memory demo principal into production IAM, nor does it remove older full-owner archives that users may already possess. CLI stdout still legitimately contains readable Case source previews, so operators should redirect it only to protected destinations. For a private no-overwrite saved artifact, use the already implemented `-save-import-snapshot`.

## Integration and merge rule

Draft child of frozen PR #69 SHA `cc07a9121bb653d36ad0024aafa03217b9e8c65c` (not main). Preserve every previously L2 frozen ancestor SHA. Check exact-head GitHub Actions (Go+Node Linux; Go Windows/macOS) and **unchanged mandatory security-critical coverage >=90%**. Existing PR #69 at 88.61% fails `CHK-SECURITY-TESTS`; this fix does not waive that requirement or claim that any ancestor can already merge. As a material data boundary fix, schedule only one targeted read-only L2 after stable publish, without repeating parent reviews.
