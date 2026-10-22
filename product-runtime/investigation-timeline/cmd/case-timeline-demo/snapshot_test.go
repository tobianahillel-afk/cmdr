package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

func snapshotFixture(t *testing.T) string {
	t.Helper()
	data := demoData()
	// The disk dataset is NOT the compiled Case fixture: the entire
	// canonical Case/Timeline/source relationship is changed together.
	data.Cases[0].ID = "forensic-case-on-disk"
	for i := range data.Entries {
		if data.Entries[i].CaseID == "case-a" {
			data.Entries[i].CaseID = "forensic-case-on-disk"
		}
	}
	for i := range data.Sources {
		if data.Sources[i].Identity.CaseID == "case-a" {
			data.Sources[i].Identity.CaseID = "forensic-case-on-disk"
		}
	}
	payload, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "owner-snapshot.json")
	if err := os.WriteFile(path, payload, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestSnapshotDiskToAuthenticatedCaseAndTimelineHTTP(t *testing.T) {
	path := snapshotFixture(t)
	data, err := loadDemoSnapshot(path)
	if err != nil {
		t.Fatal(err)
	}
	caseID, err := firstDemoAuthorizedCase(data, "analyst-a")
	if err != nil || caseID != "forensic-case-on-disk" {
		t.Fatalf("disk Case not selected: %q %v", caseID, err)
	}
	handler, err := newDemoHandlerWithData(testBearer, "analyst-a", data)
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(handler)
	defer server.Close()
	status, result := requestDemo(t, server, "/demo-cases", testBearer)
	if status != http.StatusOK || !strings.Contains(result, "forensic-case-on-disk") ||
		strings.Contains(result, "case-b") {
		t.Fatalf("Case discovery did not use disk ACL: HTTP %d %s", status, result)
	}
	status, result = requestDemo(t, server, "/case-timeline?case_id=forensic-case-on-disk&limit=20", testBearer)
	if status != http.StatusOK {
		t.Fatalf("disk Timeline HTTP %d %s", status, result)
	}
	var got struct {
		Entries []timeline.ProjectedEntry `json:"entries"`
	}
	if err := json.Unmarshal([]byte(result), &got); err != nil || len(got.Entries) != 2 ||
		got.Entries[0].ID != "event-old" || got.Entries[1].ID != "event-new" {
		t.Fatalf("disk entries not projected correctly: %+v %v", got, err)
	}
	for _, forbidden := range []string{"demo-secret", "demo-hidden-source", "demo-tenant-b", "restricted"} {
		if strings.Contains(result, forbidden) {
			t.Fatalf("disk HTTP leaked restricted source %s", forbidden)
		}
	}
	for _, req := range []struct {
		path, token string
		status int
	}{
		{"/case-timeline?case_id=case-a", testBearer, 404},
		{"/case-timeline?case_id=case-b", testBearer, 404},
		{"/case-timeline?case_id=forensic-case-on-disk", "", 401},
		{"/case-timeline?case_id=forensic-case-on-disk", "bad", 401},
	} {
		code, body := requestDemo(t, server, req.path, req.token)
		if code != req.status || strings.Contains(body, "event-new") ||
			strings.Contains(body, "demo-hidden-source") {
			t.Fatalf("denied disk request %q: %d %s", req.path, code, body)
		}
	}
}

func TestSnapshotDiskCannotReplaceOwnerAuthority(t *testing.T) {
	path := snapshotFixture(t)
	data, err := loadDemoSnapshot(path)
	if err != nil {
		t.Fatal(err)
	}
	// The disk cannot escalate access via a source permission flag.
	for i := range data.Entries {
		data.Entries[i].Value.SourcePermission = "allow"
	}
	h, err := newDemoHandlerWithData(testBearer, "analyst-a", data)
	if err != nil {
		t.Fatal(err)
	}
	s := httptest.NewServer(h)
	defer s.Close()
	_, body := requestDemo(t, s, "/case-timeline?case_id=forensic-case-on-disk", testBearer)
	if strings.Contains(body, "demo-hidden-source") || strings.Contains(body, "restricted") {
		t.Fatalf("disk authorization claim bypassed real owner ACL: %s", body)
	}
	// Contradictory owner membership is rejected before listener creation.
	data.Cases[0].EntryIDs = append(data.Cases[0].EntryIDs, "event-b")
	if _, err := newDemoHandlerWithData(testBearer, "analyst-a", data); err == nil {
		t.Fatal("cross-tenant disk catalog was accepted")
	}
}

func TestSnapshotDiskRejectsInvalidContentAndPaths(t *testing.T) {
	valid := snapshotFixture(t)
	data, err := os.ReadFile(valid)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name string
		contents []byte
	}{
		{"trailing-document", append(append([]byte(nil), data...), []byte("\n{}")...)},
		{"unknown-field", []byte(`{"Extra":"forged"}`)},
		{"malformed", []byte(`{"Cases":[`)},
		{"excess-size", bytes.Repeat([]byte("X"), maxDemoSnapshotBytes+1)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			file := filepath.Join(t.TempDir(), "bad.json")
			if err := os.WriteFile(file, tc.contents, 0o600); err != nil {
				t.Fatal(err)
			}
			if _, err := loadDemoSnapshot(file); err == nil {
				t.Fatal("unsafe snapshot accepted")
			}
		})
	}
	for _, path := range []string{"relative.json", t.TempDir(), filepath.Join(t.TempDir(), "missing")} {
		if _, err := loadDemoSnapshot(path); err == nil {
			t.Fatalf("unsafe or missing path accepted: %q", path)
		}
	}
	symlink := filepath.Join(t.TempDir(), "link.json")
	if err := os.Symlink(valid, symlink); err == nil {
		if _, err := loadDemoSnapshot(symlink); err == nil {
			t.Fatal("symlinked snapshot accepted")
		}
	}
}

func TestSnapshotCLISelectsDiskCaseAndRetainsDefaults(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var output bytes.Buffer
	if err := runDemoCLI(ctx, []string{"-principal=analyst-a", "-snapshot="+snapshotFixture(t)}, &output); !errors.Is(err, context.Canceled) || output.Len()!=0 {
		t.Fatalf("cancelled snapshot startup exposed listener/token: %v %s", err, output.String())
	}
	if got, err := loadDemoSnapshot(""); err != nil || len(got.Cases) != len(demoData().Cases) {
		t.Fatalf("default embedded fixture changed: %v", err)
	}
	data, err := loadDemoSnapshot(snapshotFixture(t))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := firstDemoAuthorizedCase(data, "analyst-b"); err != nil {
		t.Fatalf("second tenant's Case lost: %v", err)
	}
	if _, err := firstDemoAuthorizedCase(data, "unknown"); err == nil {
		t.Fatal("unknown snapshot principal selected a Case")
	}
}
