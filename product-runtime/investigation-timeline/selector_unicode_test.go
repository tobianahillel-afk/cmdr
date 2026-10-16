package investigationtimeline

import (
	"context"
	"testing"
)

func TestAuthorizationSelectorsRejectNonTextUnicode(t *testing.T) {
	selectors := []struct {
		name  string
		value string
	}{
		{"combining-grapheme-joiner", "x\u034fy"},
		{"variation-selector-16", "x\ufe0fy"},
		{"supplementary-variation-selector", "x\U000e0100y"},
		{"mongolian-variation-selector", "x\u180by"},
		{"invalid-utf8", "x" + string([]byte{0xff}) + "y"},
	}
	for _, tc := range selectors {
		t.Run(tc.name, func(t *testing.T) {
			if !invalidScope(tc.value) {
				t.Fatalf("invisible/invalid authority selector accepted: %q", tc.value)
			}
			for _, target := range []string{"tenant", "environment", "case", "timeline-id", "linked-source"} {
				t.Run(target, func(t *testing.T) {
					in := allowedFixture()
					caseID := "case-123"
					switch target {
					case "tenant":
						in.Scope.Tenant = tc.value
					case "environment":
						in.Scope.Environment = tc.value
					case "case":
						caseID = tc.value
					case "timeline-id":
						in.Entries[0].ID = tc.value
					case "linked-source":
						in.Entries[0].SourceRef = tc.value
					}
					checker := allowedAccessChecker()
					legacy := ProjectWithAccessChecker(context.Background(), in, caseID, checker)
					if target == "tenant" || target == "environment" || target == "case" {
						if legacy.Status != "deny" || len(checker.calls) != 0 {
							t.Fatalf("malformed scope reached authority: %+v calls=%+v", legacy, checker.calls)
						}
					} else if target == "timeline-id" {
						if len(legacy.Entries) != 0 || len(checker.calls) != 1 {
							t.Fatalf("malformed ID reached Timeline provider: %+v calls=%+v", legacy, checker.calls)
						}
					} else {
						if legacy.Status != "deny" || len(checker.calls) != 2 {
							t.Fatalf("malformed linked source reached owner provider: %+v calls=%+v", legacy, checker.calls)
						}
					}
					if target != "case" {
						projected := Project(in)
						if target == "tenant" || target == "environment" || target == "timeline-id" || target == "linked-source" {
							if projected.Status != "deny" || len(projected.Entries) != 0 {
								t.Fatalf("invalid selector reached projection: %+v", projected)
							}
						}
					}
			})
		}
	})
	}
}

func TestOwnerSelectorsFailBeforeUntrustedAuthority(t *testing.T) {
	bad := []string{"x\u034fy", "x\ufe0fy", "x\U000e0100y", "x\u180by", "x" + string([]byte{0xff})}
	for _, value := range bad {
		checker, owner := ownerFixture()
		got := ProjectFromOwnerSnapshots(context.Background(), allowedFixture().Scope, "case-123", []string{value}, checker, owner)
		if got.Status != "deny" || len(checker.calls) != 0 || len(owner.calls) != 0 {
			t.Fatalf("malformed selected ID crossed owner boundary: %+v calls=%+v", got, checker.calls)
		}
		checker, owner = ownerFixture()
		entry := owner.entries["x"]
		entry.SourceRef = value
		owner.entries["x"] = entry
		got = ProjectFromOwnerSnapshots(context.Background(), allowedFixture().Scope, "case-123", []string{"x"}, checker, owner)
		if got.Status != "deny" || len(checker.calls) != 2 || len(owner.calls) != 1 {
			t.Fatalf("malformed source selector reached linked-source authority: %+v calls=%+v", got, checker.calls)
		}
	}
}

func TestVisibleCombiningAccentRemainsValid(t *testing.T) {
	// Do not prohibit all combining marks: legitimate visible names/paths
	// may use Unicode combining accents rather than precomposed glyphs.
	if invalidScope("r\u0065\u0301fe\u0301rence") {
		t.Fatal("ordinary visible combining accents were overblocked")
	}
}
