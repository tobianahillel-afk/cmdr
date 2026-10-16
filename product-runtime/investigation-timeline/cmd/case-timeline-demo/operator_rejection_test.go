package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// A selected file is never an authentication oracle. All invalid operator
// principals and incomplete owner data must fail BEFORE a listener or bearer.
func TestOperatorSnapshotFailsClosedAtStartup(t *testing.T) {
	originalPath, err := filepath.Abs(filepath.Join("testdata", "example-owners.json"))
	if err != nil {
		t.Fatal(err)
	}
	original, err := loadDemoSnapshot(originalPath)
	if err != nil {
		t.Fatal(err)
	}
	write := func(t *testing.T, data timeline.MemoryTimelineData) string {
		t.Helper()
		encoded, err := json.Marshal(data)
		if err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(t.TempDir(), "operator.json")
		if err := os.WriteFile(path, encoded, 0o600); err != nil {
			t.Fatal(err)
		}
		return path
	}
	for _, tc := range []struct {
		name, principal string
		path string
		mutate func(*timeline.MemoryTimelineData)
	}{
		{name: "unknown-principal", principal: "not-a-member", path: originalPath},
		{name: "absent-file", principal: "investigator-42", path: filepath.Join(t.TempDir(), "not-found.json")},
		{name: "principal-without-case", principal: "investigator-42", mutate: func(d *timeline.MemoryTimelineData) { for i := range d.Cases { if d.Cases[i].Tenant == "local-tenant-a" { d.Cases[i].Readers = nil } } }},
		{name: "invalid-source-owner", principal: "investigator-42", mutate: func(d *timeline.MemoryTimelineData) { d.Sources[0].Readers = []string{"forged-principal"} }},
		{name: "cross-tenant-case-record", principal: "investigator-42", mutate: func(d *timeline.MemoryTimelineData) { d.Cases[0].EntryIDs = append(d.Cases[0].EntryIDs, "alert-99") }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := tc.path
			if tc.mutate != nil {
				copy := original
				copy.Cases = append([]timeline.MemoryTimelineCase(nil), original.Cases...)
				copy.Sources = append([]timeline.MemoryTimelineSource(nil), original.Sources...)
				tc.mutate(&copy)
				path = write(t, copy)
			}
			var visible bytes.Buffer
			if err := serveDemoWithSnapshot(context.Background(), tc.principal, path, &visible); err == nil || visible.Len() != 0 {
				t.Fatalf("invalid owner startup exposed a listener or bearer: %v %q", err, visible.String())
			}
		})
	}
	var output bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{"-principal=investigator-42"}, &output); err == nil || output.Len() != 0 {
		t.Fatalf("imported principal accepted without an explicit snapshot: %v %q", err, output.String())
	}
}

