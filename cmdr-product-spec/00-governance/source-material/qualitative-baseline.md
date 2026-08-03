---
id: phase-0-qualitative-baseline
domain: 00-governance
status: draft
owner: QA and Traceability Lead
updated: 2026-08-03
source-of-truth: source-material
---
# Phase 0 Qualitative Baseline

## Independent baseline

| Metric | Baseline |
|---|---:|
| Markdown files before Phase 0 | 781 |
| Draft | 767 |
| In Review | 0 |
| Validated or approved | 0 |
| Implemented | 0 |
| Deprecated | 2 |
| Without status | 12 |
| Generic document skeleton | 559 files |
| Template-level screens | 61 |
| Insufficiently formalized objects | 57 |
| `À compléter` occurrences | 461 |
| Exact repeated placeholder | 289 |
| Files using one of three generic blocks | 677 / 781 |
| Estimated generic or repeated text | approximately 60% |

## Requirement baseline

| Metric | Count |
|---|---:|
| Requirement IDs | 122 |
| Decisions | 55 |
| Constraints | 50 |
| Preferences | 2 |
| Open questions | 15 |
| Conform | 0 |
| Partial | 68 |
| Absent | 47 |
| Contradictory | 7 |

## Principal contradictions

1. Current source-of-truth policy does not place sponsor source material above ADRs.
2. Studio permissions use both `perm.cmdr-studio.*` and `perm.studio.*`.
3. Settings permissions use both `perm.platform-settings.*` and `perm.settings.*`.
4. Work Queue Saved View language appears in unrelated product screens.
5. Tool Call and Automation Run are absent from the object model.
6. Endpoint Agent is described as complete while important capabilities remain planned or placeholder.
7. Screen files satisfy structural sections but not substantive screen behavior.

## Phase 0 quality gates

- source requirements have stable IDs;
- open decisions have owner, target phase and blocking state;
- no proposal is labeled decided;
- no product module, screen or detailed object is bulk-rewritten;
- repository root README remains exactly `# cmdr`;
- PR #2 remains open, draft and unmerged.

## Outcome

Phase 0 is **PASS for requirement capture and PARTIAL for current canonical coverage**. The repository remains **FAIL qualitative** until later product, UX and functional phases resolve the baseline anomalies.