# E4-SEC-002B-E2EREVAL — Four-module verified runtime toolchain transition

**Security owner:** E4-SEC-002B existing engineering sublot; bounded child E2EREVAL depends on E4-SEC-002B-SCAID IMPLEMENTED (PR #196). **Code owner:** existing `tools/cmdr-dev` security/attestation validation. **Git base:** PR #196 frozen HEAD `013a42c1804b67bca8c181fa0b22fed585769e30`, itself child of Command PR #195; Git parent is not product/business parent. **Product behavior:** CAP-INV-002, CAP-INV-004, CAP-INV-005 and CAP-INV-006 are unaffected; original old VERIFIED historical runtime/perf source SHA stay immutable.

## Root-cause and exact allowed identities
A new mandatory govulncheck Go check in PR #195 found one reachable stdlib vulnerability GO-2026-6809 in Go below 1.26.10, and its strict mandatory gate returned FAILURE. PR #196 changed nine Go module declarations to official patched Go 1.26.10 without changing product code, passed the Command and Timeline jobs, but its global CI correctly rejected historical E2E claims that the runtime source was byte-identical. The four legacy checkers are Event Inspection, Event Search, Hunt Management and Saved Query Assets. This child implements a **single common** pin-and-revalidate flow rather than four unsafe exceptions.

Every module must differ *solely* in go.mod compared with its own frozen evidence commit. Both old and new go.mod SHA-1 Git blobs are static and reviewed:
| Runtime | old go 1.26.8 blob | patched go 1.26.10 blob |
|---|---|---|
| Event Inspection | `a2b90dca906b971f167a9108b6436bf9580a82e7` | `9f45d10c3e854581bb835f326a8d3a5572dcea1f` |
| Event Search | `aa7adaad1e3d7b57fcc698c842fb965a23c30068` | `c954ca9b250acb749cda2f3bd6d98983777bd136` |
| Hunt Management | `e05aba665edb6051273695bae4f5464f40f93559` | `29d4f0aa7b40cd0c0e859f70e4d04a2d777a4e31` |
| Saved Query Assets | `6c1ebc18e40d413a998916d42194cc0230d33a1a` | `ac9ab9259a841ce23b1e9fa230341f6bd7c8f104` |

A source diff absent => historical runtime unchanged flag true. Exactly one verified go.mod upgrade => false with *fresh* `go test -count=1 ./...` and `go vet ./...` run with local patched compiler; existing independent contract/adversarial/performance and mandatory SCA gates still execute. Every other runtime content drift, changed go.mod, surprise provider or non-ancestor revision returns blocking ERROR. Security synthetic git tests check both the genuine byte transition and rejection of a future version or extra source.

## Remaining limits/rollout and review
CI MUST run on the exact final SHA. A green child branch fixes the inherited validation path, not necessarily the parent PR #196's separate red CI; rebase/consolidate safely before merging according to native approvals and mandatory branch protection. Do not forge parent check conclusions, old evidence, GitHub native APPROVED or production readiness. Runtime products retain synthetic/local owners and their own open decisions. No UI/DB/performance claims, no dependency suppression or paid API. The child changes only the engineering control-plane; any genuine new runtime code must obtain a separate material L2 and E2E proof.
