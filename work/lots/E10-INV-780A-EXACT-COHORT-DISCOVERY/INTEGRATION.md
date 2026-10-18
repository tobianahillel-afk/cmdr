# E10-INV-780A — Find authorized Wazuh cohorts by exact rule or agent

**Product delta.** On two separately authorized SHA-256-verified persisted Cases, an analyst can restrict the existing metadata-only Wazuh cohort catalog with either `-catalog-rule-id RULE`, `-catalog-agent-id AGENT`, or both, together with `-list-saved-wazuh-cohorts-offset 0`. Filtering precedes sorting and pagination. A known rule or agent no longer requires paging through unrelated cohort labels; select a returned cohort and reuse the existing exact pair page, Entry-ID pivot and private dossier mechanisms with their separate authorization.

Existing output, no-filter behavior and source reauthorization are preserved. No new service, engine, dependency, privileged API or stdout source previews. The new selectors cannot accompany private dossier, pair pagination, exact-ID discovery or other CLI modes. The count represents time-correlated Wazuh observations, **not causality or evidence authorization**.

**Branch:** stacked on PR #124 / E10-INV-770A at exact parent `37c9a64714d5a0ae4a08691c69480126500b1bf8`; its and all earlier parent branches remain unchanged, DRAFT and unmerged into `main`. CI requires the actual child HEAD. A single context-separated READ-ONLY L2 should review the filter-to-authorization boundary, before/after page semantics and private output regressions once the feature HEAD is stable.
