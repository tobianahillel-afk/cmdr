package main

import (
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestBundledOfflineCaseExampleEndToEnd(t *testing.T) {
	data, err := loadDemoDataset(filepath.Join("testdata", "example-owners.json"))
	if err != nil {
		t.Fatal(err)
	}
	handler, err := newDemoHandlerWithData(testBearer, "investigator-42", data)
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(handler)
	defer server.Close()
	status, body := requestDemo(t, server, "/demo-cases", testBearer)
	if status != 200 || !strings.Contains(body, "case-incident-42") ||
		strings.Contains(body, "case-private-99") {
		t.Fatalf("imported Case discovery failed: %d %s", status, body)
	}
	status, body = requestDemo(t, server, "/case-timeline?case_id=case-incident-42&limit=20", testBearer)
	if status != 200 || !strings.Contains(body, `"id":"alert-1"`) ||
		!strings.Contains(body, `"id":"alert-2"`) ||
		strings.Contains(body, "alert-99") {
		t.Fatalf("imported read-only timeline failed: %d %s", status, body)
	}
	status, body = requestDemo(t, server, "/case-timeline?case_id=case-private-99", testBearer)
	if status != 404 || strings.Contains(body, "alert-99") {
		t.Fatalf("cross-tenant Case was available: %d %s", status, body)
	}
}
