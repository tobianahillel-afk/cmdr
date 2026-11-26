package harness

import (
	"context"
	"errors"
	"testing"
)

func TestIsMainBranch_EnvCases(t *testing.T) {
	cases := []struct {
		name      string
		eventName string
		gitRef    string
		refName   string
		baseRef   string
		want      bool
	}{
		{
			name:      "CI_PushToMain_DetachedHEAD",
			eventName: "push",
			gitRef:    "refs/heads/main",
			refName:   "main",
			want:      true,
		},
		{
			name:      "CI_PushToFeatureBranch",
			eventName: "push",
			gitRef:    "refs/heads/feat/m3-mainline-health-gate-integration",
			refName:   "feat/m3-mainline-health-gate-integration",
			want:      false,
		},
		{
			name:      "CI_PushToTagNamedMain",
			eventName: "push",
			gitRef:    "refs/tags/main",
			refName:   "main",
			want:      false,
		},
		{
			name:      "CI_PullRequestMergeRef",
			eventName: "pull_request",
			gitRef:    "refs/pull/235/merge",
			refName:   "235/merge",
			baseRef:   "main",
			want:      false,
		},
		{
			name:      "CI_PullRequestSpoofedRef",
			eventName: "pull_request",
			gitRef:    "refs/heads/main",
			refName:   "main",
			baseRef:   "main",
			want:      false,
		},
		{
			name:      "CI_PullRequestTarget",
			eventName: "pull_request_target",
			gitRef:    "refs/heads/main",
			refName:   "main",
			want:      false,
		},
		{
			name:      "CI_WorkflowDispatchOnMain",
			eventName: "workflow_dispatch",
			gitRef:    "refs/heads/main",
			refName:   "main",
			want:      true,
		},
		{
			name:      "CI_WorkflowDispatchOnFeature",
			eventName: "workflow_dispatch",
			gitRef:    "refs/heads/feat/x",
			refName:   "feat/x",
			want:      false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("GITHUB_EVENT_NAME", tc.eventName)
			t.Setenv("GITHUB_REF", tc.gitRef)
			t.Setenv("GITHUB_REF_NAME", tc.refName)
			t.Setenv("GITHUB_BASE_REF", tc.baseRef)

			if got := IsMainBranch(context.Background()); got != tc.want {
				t.Fatalf("IsMainBranch() = %t, want %t (event=%q ref=%q refName=%q)",
					got, tc.want, tc.eventName, tc.gitRef, tc.refName)
			}
		})
	}
}

func TestIsMainBranchHelper_LocalBranchesAndErrors(t *testing.T) {
	cases := []struct {
		name          string
		eventName     string
		gitRef        string
		refName       string
		currentBranch string
		branchErr     error
		want          bool
	}{
		{"LocalMain", "", "", "", "main", nil, true},
		{"LocalFeature", "", "", "", "feat/m3", nil, false},
		{"LocalDetachedHEAD", "", "", "", "HEAD", nil, false},
		{"LocalEmptyBranch", "", "", "", "", nil, false},
		{"GitCurrentBranchError", "", "", "", "", errors.New("git failed"), false},
		{"PushMainDetachedHEAD", "push", "refs/heads/main", "main", "HEAD", nil, true},
		{"PushMainEmptyBranchWithErr", "push", "refs/heads/main", "main", "", errors.New("git failed"), true},
		{"RefNameNotMain", "", "", "feat/other", "main", nil, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isMainBranch(tc.eventName, tc.gitRef, tc.refName, tc.currentBranch, tc.branchErr); got != tc.want {
				t.Fatalf("isMainBranch(%q,%q,%q,%q,%v) = %t, want %t",
					tc.eventName, tc.gitRef, tc.refName, tc.currentBranch, tc.branchErr, got, tc.want)
			}
		})
	}
}
