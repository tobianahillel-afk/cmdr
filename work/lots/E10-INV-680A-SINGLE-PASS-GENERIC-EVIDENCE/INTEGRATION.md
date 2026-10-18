# E10-INV-680A — One authorized source opening per generic DFIR lead

## Existing problem and focused integration

The integrated cross-Case Wazuh pair dossier in frozen PR #108 contains bounded generic +/-60 second contextual entries. The default private v2 dossier contains only metadata; the explicitly opted-in v3 dossier includes sensitive generic source previews. Previously, each selected contextual generic source was independently inspected by the first-party Case/Timeline/typed-source owner to authorize the context, then v3 **reopened the same source a second time** to include its preview. That doubled owner construction and source checks on the v3 hot path, and added a separate duplicate validation function.

This child consolidates both into the existing contextual selection operation: after independently verifying BOTH Case snapshots, tenant/environment/principal, SHA-256, Wazuh pair source owners and canonical imported IDs, select up to ten nearest certain-UTC generic events per Case as before. **Open every selected typed source once through `inspectDemoEntryWithData`**, check the same Case ID, selected Entry ID, synthetic provenance and `telemetry-event` type, then:
- v2 default: retain ONLY entry ID, UTC time, offset, truncated flag; discard the sensitive preview immediately;
- v3 ONLY after explicit existing `-include-pair-generic-evidence` plus private destination: retain exactly those already source-authorized inspection records in the same existing private no-overwrite, 0600 dossier, and no metadata/preview stdout.

Any failed owner inspection returns NO partial context/evidence. The existing atomic `writeAuthorizedReceiptPrivate` waits for independent inspection of both Cases before publishing a complete file. The set of candidate IDs, matching rules, +/-60s inclusivity, sorting, limit ten, default/v3 formats and ACLs are unchanged.

### Comparison of alternatives

- **Retain double re-opening:** simple but repeats first-party owner initialization on up to 20 generic events and maintains two authorization branches. Rejected as unnecessary runtime cost.
- **Skip owner reauthorization when returning previews:** unsafe; selector IDs and source metadata never grant access. Rejected.
- **One authorized inspection reused for v2/v3:** selected. No new handler, abstraction layer, dependency, detector, IAM or persistence engine.

### Integration acceptance

All preexisting Go end-to-end tests must remain green on exact HEAD: real Wazuh+normalized JSONL saved source files; independent SHA-256, tenant/environment and Case rights; v2 no generic previews; v3 explicit private previews only for up to ten source-authorized generic entries on each side; revoked first/second generic source never leaks, revoked Wazuh source blocks publication; 60-second bounds and ordering unchanged; no bearer or preview stdout, no overwrite or partial dossier. Linux Go+Node/browser, macOS and Windows via existing `investigation-timeline-runtime` Actions. **Unchanged** mandatory engineering-kernel `CHK-SECURITY-TESTS` >=90% changed security-critical coverage. Never treat green runtime as merge approval if engineering/ancestors fail.

### Integration status and limits

Child of immutable PR #108 SHA `c5249657ac35ba2a718fc451ca85a1e9e4428c6c`, based on consolidation PR #92 (also draft). No parent SHA is changed; #108 full 64-commit comparison failed at **89.79% <90%** even though runtime CI passed. This change is a hot-path simplification, NOT a release, a security-gate waiver, real tenant IAM, live provider, production storage, signed chain-of-custody or causal correlation. A single read-only L2 should target the narrowed functional/security delta after exact-head CI. The much larger PR #36 and product-spec PR #2 still need their independent integration approvals and product decisions.
