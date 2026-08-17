# Endpoint EPT-3 Changelog

## 2026-08-11 — build
- Started **EPT-3 — Local Detection and Endpoint Investigation** from exact EPT-2 verification baseline `5d7c037aff6004984416665e7e188a8700e62b2f`.
- Revalidated `CAP-EPT-031..046` as free/unreserved and retained exactly 16 independently justified capabilities after reading all 8 Detection and all 8 Investigation source documents.
- Added `CAP-EPT-031..046`, all `draft / defined / planned`: Detection Content consumption/eligibility, local evaluation/match/candidate/context/grouping/coverage, process/file/network/user/system contexts, local timeline/correlation, pivots, detection-to-investigation expansion, summary/handoff and provenance.
- Preserved EPT-1 190/190 and EPT-2 200/200; `CAP-EPT-001..030` are not modified.
- Preserved Investigate ownership of Detection Engineering/Case/Evidence/Finding, Command ownership of canonical Detection/Signal/Alert/Incident, Shared generic Search/Timeline/Linking/Correlation, Settings administration, Govern response authority and Studio Tool/Run semantics.
- `OPEN-008` and `OPEN-017` remain open; no platform support, detection runtime/language/model or engine is selected.
- Endpoint cumulative build target: **46 / 1242 / 276**. Global target: **431 capabilities / 429 defined / 2 proposed / 431 planned / 11637 sections / 2586 tables**.
- Build-time quality before publication: **204 PASS / 6 PENDING-BUILD-OR-REMOTE / 0 FAIL**; EPT-4..6 remain NOT STARTED.
- No acquisition, Live Response, containment, response execution, API/protocol, physical schema, final RBAC or product code is introduced.

## 2026-08-11 — post-publication verification
- Functional build SHA: `941bfb5da8a3598ca3dd79d135246b0a8865a31a`; baseline → build **5 ahead / 0 behind**, same merge base.
- Remote checks confirmed PR #2 open/Draft/unmerged, `main` unchanged, README branch/main unchanged and CI/status N/A.
- `CAP-EPT-031..046` remain **16 / 432 / 96 / >=48 GWT**; `CAP-EPT-001..030` remain intact and `CAP-EPT-047` is absent.
- OPEN remains 18 with OPEN-008 and OPEN-017 open; Endpoint Screen IDs remain 0; EPT-4..6 remain NOT STARTED.
- Final verdict becomes effective after publication/recheck of this documentation-only record: **EPT-3 PASS AFTER POST-PUBLICATION VERIFICATION — 210/210 PASS, 0 PENDING, 0 FAIL**.
- Exact final verification-record SHA is recorded in PR #2 after publication to avoid self-reference.