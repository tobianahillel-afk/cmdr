package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestUnexecutedRuntimeBlocksReportsOnlyRealChangedCoverageMisses(t *testing.T) {
	adapter := runtimeSecurityAdapter{
		BoundaryID:     "investigation-timeline-runtime",
		RuntimeRoot:    "product-runtime/investigation-timeline",
		ModuleIdentity: "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline",
	}
	const profile = "mode: atomic\n" +
		"github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline/cmd/case-timeline-demo/normalized_search.go:90.2,95.6 8 0\n" +
		"github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline/cmd/case-timeline-demo/normalized_search.go:12.1,18.3 2 1\n" +
		"github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline/cmd/case-timeline-demo/normalized_search.go:105.3,108.1 3 0\n" +
		"github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline/cmd/case-timeline-demo/other.go:30.1,31.3 25 0\n"
	path := "product-runtime/investigation-timeline/cmd/case-timeline-demo/normalized_search.go"
	got := unexecutedRuntimeBlocks([]byte(profile), []string{path}, adapter, 8)
	want := []string{path + ":90 (8 unexecuted statements)",
		path + ":105 (3 unexecuted statements)"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("changed security coverage diagnostic = %v, want %v", got, want)
	}
	// When this PR edits only the kernel, show module-level gaps rather
	// than claiming the existing parent PR's security debt has disappeared.
	got = unexecutedRuntimeBlocks([]byte(profile), []string{"tools/cmdr-dev/security_tests.go"}, adapter, 2)
	if len(got) != 2 || !strings.Contains(got[0], "other.go:30") ||
		!strings.Contains(got[1], "normalized_search.go:90") {
		t.Fatalf("unmodified-runtime fallback failed: %v", got)
	}
	if got = unexecutedRuntimeBlocks([]byte(profile), []string{path}, adapter, 0); len(got) != 0 {
		t.Fatalf("zero limit returned diagnostic blocks: %v", got)
	}
}

func TestUnexecutedRuntimeBlocksDoesNotLeakContentOrClaimCoverage(t *testing.T) {
	adapter := runtimeSecurityAdapter{
		BoundaryID:     "investigation-timeline-runtime",
		RuntimeRoot:    "product-runtime/investigation-timeline",
		ModuleIdentity: "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline",
	}
	const profile = "mode: atomic\n" +
		"github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline/cmd/owner.go:6.1,7.1 2 0\n" +
		"github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline/cmd/owner.go:8.1,10.1 3 1\n" +
		"bad record\n"
	before, err := parseRuntimeCoverageProfile([]byte(profile[:strings.Index(profile, "bad record")]), []string{
		"product-runtime/investigation-timeline/cmd/owner.go",
	}, adapter)
	if err != nil {
		t.Fatal(err)
	}
	_ = unexecutedRuntimeBlocks([]byte(profile), nil, adapter, 10)
	after, err := parseRuntimeCoverageProfile([]byte(profile[:strings.Index(profile, "bad record")]), []string{
		"product-runtime/investigation-timeline/cmd/owner.go",
	}, adapter)
	if err != nil || before != after {
		t.Fatalf("diagnostic modified scored coverage: before=%+v after=%+v err=%v", before, after, err)
	}
	for _, line := range unexecutedRuntimeBlocks([]byte(profile), nil, adapter, 10) {
		if strings.Contains(line, "Bearer") || strings.Contains(line, "preview") ||
			!strings.Contains(line, "cmd/owner.go:6") {
			t.Fatalf("diagnostic exposed unsafe or invalid information: %q", line)
		}
	}
}
