package harness

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	StatusVerified   = "VERIFIED"
	StatusUnverified = "UNVERIFIED"
)

var ErrUnverified = errors.New("UNVERIFIED")

// OpenPRsResult captures the verified or unverified count of open GitHub PRs.
type OpenPRsResult struct {
	Count    int
	Verified bool
	Status   string
	Reason   string
}

// CodexResolutionResult captures the verified or unverified state of PR review-thread resolution.
type CodexResolutionResult struct {
	Resolved   bool
	Verified   bool
	Status     string
	MatchedSHA string
	Source     string
	Reason     string
}

// PRMetadata represents PR reconciliation records stored locally.
type PRMetadata struct {
	PRNumber             int    `json:"pr_number"`
	Title                string `json:"title,omitempty"`
	HeadBranch           string `json:"head_branch,omitempty"`
	BaseBranch           string `json:"base_branch,omitempty"`
	Status               string `json:"status,omitempty"`
	Disposition          string `json:"disposition,omitempty"`
	DispositionRationale string `json:"disposition_rationale,omitempty"`
	Commit               string `json:"commit,omitempty"`
	HeadSHA              string `json:"head_sha,omitempty"`
	MergeCommit          string `json:"merge_commit,omitempty"`
	SHA                  string `json:"sha,omitempty"`
}

// GitHubClient provides helpers to inspect PRs and review threads.
type GitHubClient struct {
	repoRoot string
}

// NewGitHubClient creates a GitHub client referencing the local repository root.
func NewGitHubClient() (*GitHubClient, error) {
	root, err := RepoRoot()
	if err != nil {
		return nil, err
	}
	return &GitHubClient{repoRoot: root}, nil
}

// ParseOpenPRsJSON parses the JSON output of `gh pr list --state open --json number`.
func ParseOpenPRsJSON(stdout string) (OpenPRsResult, error) {
	trimmed := strings.TrimSpace(stdout)
	if trimmed == "" {
		reason := "empty gh pr list response"
		return OpenPRsResult{Count: -1, Verified: false, Status: StatusUnverified, Reason: reason}, fmt.Errorf("%w: %s", ErrUnverified, reason)
	}
	var prs []map[string]interface{}
	if err := json.Unmarshal([]byte(trimmed), &prs); err != nil || prs == nil {
		reason := fmt.Sprintf("malformed gh pr list JSON: %v", err)
		if err == nil {
			reason = "malformed gh pr list JSON: expected array, got null"
		}
		return OpenPRsResult{Count: -1, Verified: false, Status: StatusUnverified, Reason: reason}, fmt.Errorf("%w: %s", ErrUnverified, reason)
	}
	for i, item := range prs {
		if item == nil {
			reason := fmt.Sprintf("malformed gh pr list JSON: item[%d] is null", i)
			return OpenPRsResult{Count: -1, Verified: false, Status: StatusUnverified, Reason: reason}, fmt.Errorf("%w: %s", ErrUnverified, reason)
		}
		num, ok := item["number"].(float64)
		if !ok || num <= 0 || num != float64(int(num)) {
			reason := fmt.Sprintf("malformed gh pr list JSON: item[%d] missing positive integer \"number\"", i)
			return OpenPRsResult{Count: -1, Verified: false, Status: StatusUnverified, Reason: reason}, fmt.Errorf("%w: %s", ErrUnverified, reason)
		}
	}
	return OpenPRsResult{
		Count:    len(prs),
		Verified: true,
		Status:   StatusVerified,
		Reason:   "live gh pr list verified",
	}, nil
}

func parseOpenPRsJSON(stdout string) (OpenPRsResult, error) {
	return ParseOpenPRsJSON(stdout)
}

// CountOpenPRs queries GitHub via `gh pr list --state open --json number`.
// Fails closed with ErrUnverified / StatusUnverified if gh CLI is unavailable or returns invalid JSON.
func (c *GitHubClient) CountOpenPRs(ctx context.Context) (OpenPRsResult, error) {
	res := RunCmd(ctx, c.repoRoot, nil, "gh", "pr", "list", "--state", "open", "--json", "number")
	if res.ExitCode == 0 && res.Err == nil {
		return ParseOpenPRsJSON(res.Stdout)
	}
	reason := strings.TrimSpace(res.Stderr)
	if reason == "" && res.Err != nil {
		reason = res.Err.Error()
	}
	if reason == "" {
		reason = fmt.Sprintf("gh pr list exited with code %d", res.ExitCode)
	}
	return OpenPRsResult{
		Count:    -1,
		Verified: false,
		Status:   StatusUnverified,
		Reason:   reason,
	}, fmt.Errorf("%w: %s", ErrUnverified, reason)
}

func verifyGraphQLPRResolution(stdout string, expectedSHA string) (CodexResolutionResult, error) {
	exp := strings.TrimSpace(expectedSHA)
	if exp == "" {
		reason := "expectedSHA must be non-empty"
		return CodexResolutionResult{Verified: false, Status: StatusUnverified, Source: "github_graphql", Reason: reason}, fmt.Errorf("%w: %s", ErrUnverified, reason)
	}

	type gqlResp struct {
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
		Data struct {
			Repository struct {
				PullRequest *struct {
					State       string `json:"state"`
					MergeCommit *struct {
						OID string `json:"oid"`
					} `json:"mergeCommit"`
					HeadRefOID    string `json:"headRefOid"`
					ReviewThreads struct {
						PageInfo struct {
							HasNextPage bool `json:"hasNextPage"`
						} `json:"pageInfo"`
						Nodes []struct {
							IsResolved bool `json:"isResolved"`
						} `json:"nodes"`
					} `json:"reviewThreads"`
				} `json:"pullRequest"`
			} `json:"repository"`
		} `json:"data"`
	}

	var resp gqlResp
	if err := json.Unmarshal([]byte(strings.TrimSpace(stdout)), &resp); err != nil {
		reason := fmt.Sprintf("malformed GraphQL response: %v", err)
		return CodexResolutionResult{Verified: false, Status: StatusUnverified, Source: "github_graphql", Reason: reason}, fmt.Errorf("%w: %s", ErrUnverified, reason)
	}
	if len(resp.Errors) > 0 {
		reason := fmt.Sprintf("GraphQL errors: %s", resp.Errors[0].Message)
		return CodexResolutionResult{Verified: false, Status: StatusUnverified, Source: "github_graphql", Reason: reason}, fmt.Errorf("%w: %s", ErrUnverified, reason)
	}
	pr := resp.Data.Repository.PullRequest
	if pr == nil || strings.TrimSpace(pr.State) == "" {
		reason := "GraphQL response missing pullRequest state"
		return CodexResolutionResult{Verified: false, Status: StatusUnverified, Source: "github_graphql", Reason: reason}, fmt.Errorf("%w: %s", ErrUnverified, reason)
	}
	if !strings.EqualFold(strings.TrimSpace(pr.State), "MERGED") {
		reason := fmt.Sprintf("pullRequest is not MERGED (state=%q)", strings.TrimSpace(pr.State))
		return CodexResolutionResult{Verified: false, Status: StatusUnverified, Source: "github_graphql", Reason: reason}, fmt.Errorf("%w: %s", ErrUnverified, reason)
	}
	if pr.ReviewThreads.PageInfo.HasNextPage {
		reason := "reviewThreads pagination incomplete (hasNextPage=true)"
		return CodexResolutionResult{Verified: false, Status: StatusUnverified, Source: "github_graphql", Reason: reason}, fmt.Errorf("%w: %s", ErrUnverified, reason)
	}
	if len(pr.ReviewThreads.Nodes) == 0 {
		reason := "no reviewThreads returned in GraphQL response"
		return CodexResolutionResult{Verified: false, Status: StatusUnverified, Source: "github_graphql", Reason: reason}, fmt.Errorf("%w: %s", ErrUnverified, reason)
	}

	mergeOID := ""
	if pr.MergeCommit != nil {
		mergeOID = strings.TrimSpace(pr.MergeCommit.OID)
	}
	headOID := strings.TrimSpace(pr.HeadRefOID)

	matchedSHA := mergeOID
	if matchedSHA == "" {
		matchedSHA = headOID
	}

	switch {
	case strings.EqualFold(mergeOID, exp):
		matchedSHA = mergeOID
	case strings.EqualFold(headOID, exp):
		matchedSHA = headOID
	default:
		return CodexResolutionResult{
			Resolved:   false,
			Verified:   true,
			Status:     "SHA_MISMATCH",
			MatchedSHA: matchedSHA,
			Source:     "github_graphql",
			Reason:     fmt.Sprintf("expected SHA %s matched neither mergeCommit %q nor headRefOid %q", exp, mergeOID, headOID),
		}, nil
	}

	allResolved := true
	for _, node := range pr.ReviewThreads.Nodes {
		if !node.IsResolved {
			allResolved = false
			break
		}
	}

	reason := "all GraphQL review threads resolved"
	if !allResolved {
		reason = "one or more GraphQL review threads remain unresolved"
	}

	return CodexResolutionResult{
		Resolved:   allResolved,
		Verified:   true,
		Status:     StatusVerified,
		MatchedSHA: matchedSHA,
		Source:     "github_graphql",
		Reason:     reason,
	}, nil
}

func isValidHexSHA(sha string) bool {
	if len(sha) != 40 {
		return false
	}
	for i := 0; i < len(sha); i++ {
		c := sha[i]
		isHex := (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
		if !isHex {
			return false
		}
	}
	return true
}

func extractRecordSHA(item map[string]interface{}) string {
	for _, key := range []string{"merge_commit", "commit", "head_sha", "sha"} {
		if raw, ok := item[key].(string); ok {
			trimmed := strings.TrimSpace(raw)
			if isValidHexSHA(trimmed) {
				return trimmed
			}
		}
	}
	return ""
}

func verifyCodexPRRecord(item map[string]interface{}, prNumber int, expectedSHA string, commitExists func(string) bool) (CodexResolutionResult, error) {
	if item == nil {
		reason := fmt.Sprintf("PR #%d record missing from prs.json", prNumber)
		return CodexResolutionResult{Verified: false, Status: StatusUnverified, Source: "offline_prs_json", Reason: reason}, fmt.Errorf("%w: %s", ErrUnverified, reason)
	}

	exp := strings.TrimSpace(expectedSHA)
	if exp == "" {
		reason := fmt.Sprintf("PR #%d verification requires non-empty expectedSHA", prNumber)
		return CodexResolutionResult{Verified: false, Status: StatusUnverified, Source: "offline_prs_json", Reason: reason}, fmt.Errorf("%w: %s", ErrUnverified, reason)
	}

	statusRaw, _ := item["status"].(string)
	status := strings.TrimSpace(statusRaw)
	if !strings.EqualFold(status, "MERGED") {
		reason := fmt.Sprintf("PR #%d record in prs.json is not MERGED (status=%q)", prNumber, status)
		return CodexResolutionResult{Verified: false, Status: StatusUnverified, Source: "offline_prs_json", Reason: reason}, fmt.Errorf("%w: %s", ErrUnverified, reason)
	}

	dispRaw, _ := item["disposition"].(string)
	disp := strings.ToUpper(strings.TrimSpace(dispRaw))
	if disp == "" || disp == "REJECTED" || disp == "SUPERSEDED" || disp == "CLOSED" || strings.Contains(disp, "REJECT") || strings.Contains(disp, "SUPERSED") {
		reason := fmt.Sprintf("PR #%d record in prs.json has non-accepted disposition %q", prNumber, strings.TrimSpace(dispRaw))
		return CodexResolutionResult{Verified: false, Status: StatusUnverified, Source: "offline_prs_json", Reason: reason}, fmt.Errorf("%w: %s", ErrUnverified, reason)
	}

	revResolved, hasRevBool := item["review_threads_resolved"].(bool)
	codexResolved, hasCodexBool := item["codex_threads_resolved"].(bool)
	if !hasRevBool && !hasCodexBool {
		reason := fmt.Sprintf("PR #%d record in prs.json lacks explicit review_threads_resolved/codex_threads_resolved boolean", prNumber)
		return CodexResolutionResult{Verified: false, Status: StatusUnverified, Source: "offline_prs_json", Reason: reason}, fmt.Errorf("%w: %s", ErrUnverified, reason)
	}

	resolved := (hasRevBool && revResolved) || (hasCodexBool && codexResolved)
	if !resolved {
		reason := fmt.Sprintf("PR #%d record in prs.json explicitly marks review threads unresolved", prNumber)
		return CodexResolutionResult{Verified: false, Status: StatusUnverified, Source: "offline_prs_json", Reason: reason}, fmt.Errorf("%w: %s", ErrUnverified, reason)
	}

	recordSHA := extractRecordSHA(item)
	if recordSHA == "" {
		reason := fmt.Sprintf("PR #%d record in prs.json lacks structured 40-hex SHA attribution (merge_commit/commit/head_sha/sha)", prNumber)
		return CodexResolutionResult{Verified: false, Status: StatusUnverified, Source: "offline_prs_json", Reason: reason}, fmt.Errorf("%w: %s", ErrUnverified, reason)
	}

	if !strings.EqualFold(recordSHA, exp) {
		reason := fmt.Sprintf("PR #%d SHA attribution mismatch: expected %s, found %s", prNumber, exp, recordSHA)
		return CodexResolutionResult{Verified: false, Status: StatusUnverified, MatchedSHA: recordSHA, Source: "offline_prs_json", Reason: reason}, fmt.Errorf("%w: %s", ErrUnverified, reason)
	}

	if commitExists == nil || !commitExists(recordSHA) {
		reason := fmt.Sprintf("PR #%d attributed commit %s does not exist in local git history", prNumber, recordSHA)
		return CodexResolutionResult{Verified: false, Status: StatusUnverified, MatchedSHA: recordSHA, Source: "offline_prs_json", Reason: reason}, fmt.Errorf("%w: %s", ErrUnverified, reason)
	}

	return CodexResolutionResult{
		Resolved:   true,
		Verified:   true,
		Status:     StatusVerified,
		MatchedSHA: recordSHA,
		Source:     "offline_prs_json",
		Reason:     "explicit SHA-attributed thread resolution verified in prs.json",
	}, nil
}

// VerifyCodexPRResolution checks whether a PR's review threads are resolved and attributed to expectedSHA.
// Uses live GitHub GraphQL first; when live GraphQL fails to execute, offline prs.json requires
// MERGED status, accepted disposition, explicit boolean thread-resolution proof, structured SHA
// attribution matching expectedSHA, and local git commit existence.
func (c *GitHubClient) VerifyCodexPRResolution(ctx context.Context, prNumber int, expectedSHA string) (CodexResolutionResult, error) {
	exp := strings.TrimSpace(expectedSHA)
	if exp == "" {
		reason := fmt.Sprintf("PR #%d verification requires non-empty expectedSHA", prNumber)
		return CodexResolutionResult{Verified: false, Status: StatusUnverified, Source: "github_graphql", Reason: reason}, fmt.Errorf("%w: %s", ErrUnverified, reason)
	}

	query := fmt.Sprintf(`query { repository(owner:"tobianahillel-afk", name:"cmdr") { pullRequest(number:%d) { state mergeCommit { oid } headRefOid reviewThreads(first:100) { pageInfo { hasNextPage } nodes { isResolved } } } } }`, prNumber)
	res := RunCmd(ctx, c.repoRoot, nil, "gh", "api", "graphql", "-f", "query="+query)
	if res.ExitCode == 0 && res.Err == nil {
		return verifyGraphQLPRResolution(res.Stdout, exp)
	}

	prsPath := filepath.Join(c.repoRoot, "work", "reconciliation", "prs.json")
	data, err := os.ReadFile(prsPath)
	if err != nil {
		reason := fmt.Sprintf("failed to read prs.json: %v", err)
		return CodexResolutionResult{Verified: false, Status: StatusUnverified, Source: "offline_prs_json", Reason: reason}, fmt.Errorf("%w: %s", ErrUnverified, reason)
	}

	var prList []map[string]interface{}
	if err := json.Unmarshal(data, &prList); err != nil {
		reason := fmt.Sprintf("failed to parse prs.json as array: %v", err)
		return CodexResolutionResult{Verified: false, Status: StatusUnverified, Source: "offline_prs_json", Reason: reason}, fmt.Errorf("%w: %s", ErrUnverified, reason)
	}

	for _, item := range prList {
		if num, ok := item["pr_number"].(float64); ok && int(num) == prNumber {
			return verifyCodexPRRecord(item, prNumber, exp, func(sha string) bool {
				return GitCommitExistsInDir(ctx, c.repoRoot, sha)
			})
		}
	}

	reason := fmt.Sprintf("PR #%d not found in work/reconciliation/prs.json", prNumber)
	return CodexResolutionResult{
		Verified: false,
		Status:   StatusUnverified,
		Source:   "offline_prs_json",
		Reason:   reason,
	}, fmt.Errorf("%w: %s", ErrUnverified, reason)
}

// ClosureRegisterContainsPR checks whether closureText explicitly references PR prNumber
// using digit/alphanumeric boundaries so `#23` does not match `#230`.
func ClosureRegisterContainsPR(closureText string, prNumber int) bool {
	if prNumber <= 0 {
		return false
	}
	hashRe := regexp.MustCompile(fmt.Sprintf(`(?:^|[^0-9])#%d(?:[^0-9A-Za-z]|$)`, prNumber))
	if hashRe.MatchString(closureText) {
		return true
	}
	prRe := regexp.MustCompile(fmt.Sprintf(`(?i)(?:^|[^0-9A-Za-z])PR\s+%d(?:[^0-9A-Za-z]|$)`, prNumber))
	return prRe.MatchString(closureText)
}

// ValidatePRRegisterConsistency validates work/reconciliation/prs.json and cross-checks
// every record against work/reconciliation/PR_CLOSURE_REGISTER.md and local git history.
func ValidatePRRegisterConsistency(ctx context.Context, root string) error {
	prsPath := filepath.Join(root, "work", "reconciliation", "prs.json")
	data, err := os.ReadFile(prsPath)
	if err != nil {
		return fmt.Errorf("failed to read prs.json: %w", err)
	}

	var records []PRMetadata
	if err := json.Unmarshal(data, &records); err != nil {
		return fmt.Errorf("failed to parse prs.json: %w", err)
	}
	if len(records) == 0 {
		return fmt.Errorf("prs.json contains zero PR records")
	}

	closurePath := filepath.Join(root, "work", "reconciliation", "PR_CLOSURE_REGISTER.md")
	closureBytes, err := os.ReadFile(closurePath)
	if err != nil {
		return fmt.Errorf("failed to read PR_CLOSURE_REGISTER.md: %w", err)
	}
	closureText := string(closureBytes)

	seenPRs := make(map[int]bool, len(records))
	for _, item := range records {
		if item.PRNumber <= 0 {
			return fmt.Errorf("invalid PR number: %d", item.PRNumber)
		}
		if seenPRs[item.PRNumber] {
			return fmt.Errorf("duplicate PR #%d in prs.json", item.PRNumber)
		}
		seenPRs[item.PRNumber] = true

		if strings.TrimSpace(item.Title) == "" {
			return fmt.Errorf("PR #%d has empty title", item.PRNumber)
		}

		status := strings.TrimSpace(item.Status)
		if status != "MERGED" && status != "CLOSED" && status != "OPEN" {
			return fmt.Errorf("PR #%d has unknown status %q", item.PRNumber, item.Status)
		}

		disp := strings.TrimSpace(item.Disposition)
		rationale := strings.TrimSpace(item.DispositionRationale)

		if status == "CLOSED" || status == "MERGED" {
			if disp == "" {
				return fmt.Errorf("PR #%d missing disposition", item.PRNumber)
			}
			if rationale == "" {
				return fmt.Errorf("PR #%d missing disposition rationale", item.PRNumber)
			}
		}

		if status == "MERGED" && strings.EqualFold(disp, "REJECTED") {
			return fmt.Errorf("contradictory status: PR #%d marked MERGED but disposition is REJECTED", item.PRNumber)
		}
		if status == "OPEN" && disp != "" && !strings.EqualFold(disp, "IN_REVIEW") {
			return fmt.Errorf("contradictory status: open PR #%d cannot have terminal disposition %s", item.PRNumber, disp)
		}

		sha := strings.TrimSpace(item.MergeCommit)
		if sha == "" {
			sha = strings.TrimSpace(item.Commit)
		}
		if sha == "" {
			sha = strings.TrimSpace(item.HeadSHA)
		}
		if sha == "" {
			sha = strings.TrimSpace(item.SHA)
		}
		if sha != "" {
			if !isValidHexSHA(sha) {
				return fmt.Errorf("PR #%d has malformed commit SHA format %q", item.PRNumber, sha)
			}
			if !GitCommitExistsInDir(ctx, root, sha) {
				return fmt.Errorf("PR #%d references non-existent commit %s in git object storage", item.PRNumber, sha)
			}
		}

		if !ClosureRegisterContainsPR(closureText, item.PRNumber) {
			return fmt.Errorf("PR #%d present in prs.json is missing from PR_CLOSURE_REGISTER.md", item.PRNumber)
		}
	}

	return nil
}
