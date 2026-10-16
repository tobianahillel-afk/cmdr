package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

func importedCaseDataset() timeline.MemoryTimelineData {
	data := demoData()
	data.Principals[0].ID = "investigator-42"
	data.Cases[0].ID = "case-imported"
	data.Cases[0].Readers = []string{"investigator-42"}
	for i := range data.Entries {
		if data.Entries[i].CaseID == "case-a" {
			data.Entries[i].CaseID = "case-imported"
			data.Entries[i].Readers = []string{"investigator-42"}
		}
	}
	for i := range data.Sources {
		if data.Sources[i].Identity.CaseID == "case-a" {
			data.Sources[i].Identity.CaseID = "case-imported"
			if len(data.Sources[i].Readers) != 0 {
				data.Sources[i].Readers = []string{"investigator-42"}
			}
		}
	}
	return data
}

func writeDemoDataset(t *testing.T, data timeline.MemoryTimelineData) string {
	t.Helper()
	payload, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "owners.json")
	if err := os.WriteFile(path, payload, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestImportedLocalCaseRunsFullAuthenticatedInvestigation(t *testing.T) {
	path := writeDemoDataset(t, importedCaseDataset())
	loaded, err := loadDemoDataset(path)
	if err != nil {
		t.Fatal(err)
	}
	handler, err := newDemoHandlerWithData(testBearer, "investigator-42", loaded)
	if err != nil {
		t.Fatal(err)
	}
	request := func(route, bearer string) (int, string) {
		t.Helper()
		req, err := http.NewRequest(http.MethodGet, "http://localhost"+route, nil)
		if err != nil {
			t.Fatal(err)
		}
		if bearer != "" {
			req.Header.Set("Authorization", "Bearer "+bearer)
		}
		rec := &bufferResponse{header: make(http.Header)}
		handler.ServeHTTP(rec, req)
		return rec.status, rec.body.String()
	}
	status, cases := request("/demo-cases", testBearer)
	if status != 200 || !strings.Contains(cases, "case-imported") ||
		strings.Contains(cases, "case-b") {
		t.Fatalf("imported Case discovery incorrect: %d %s", status, cases)
	}
	status, body := request("/case-timeline?case_id=case-imported&q="+url.QueryEscape("event-old")+"&limit=10", testBearer)
	if status != 200 || !strings.Contains(body, `"id":"event-old"`) ||
		strings.Contains(body, "demo-hidden-source") || strings.Contains(body, `"id":"event-new"`) {
		t.Fatalf("authorized imported search incorrect: %d %s", status, body)
	}
	for _, route := range []string{
		"/case-timeline?case_id=case-b",
		"/case-timeline?case_id=case-a",
		"/demo-cases?tenant=demo-tenant-b",
	} {
		status, body := request(route, testBearer)
		if status == 200 || strings.Contains(body, "event-b") || strings.Contains(body, "demo-incident-b") {
			t.Fatalf("unauthorized imported selection leaked: %d %s", status, body)
		}
	}
	status, body = request("/case-timeline?case_id=case-imported", "")
	if status != 401 || strings.Contains(body, "event-old") {
		t.Fatalf("anonymous imported Case was available: %d %s", status, body)
	}
	if _, err := newDemoHandlerWithData(testBearer, "analyst-a", loaded); err == nil {
		t.Fatal("principal missing from loaded owner snapshot accepted")
	}
}

type bufferResponse struct {
	header http.Header
	body bytes.Buffer
	status int
}

func (r *bufferResponse) Header() http.Header { return r.header }
func (r *bufferResponse) Write(data []byte) (int, error) {
	if r.status == 0 {
		r.status = 200
	}
	return r.body.Write(data)
}
func (r *bufferResponse) WriteHeader(status int) { r.status = status }

func TestDatasetRejectsMalformedAndOversizeInputs(t *testing.T) {
	base := writeDemoDataset(t, importedCaseDataset())
	raw, err := os.ReadFile(base)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name string
		input []byte
	}{
		{"trailing-object", append(append([]byte(nil), raw...), []byte("{}")...)},
		{"unknown-field", append([]byte(`{"unexpected":"value",`), raw[1:]...)},
		{"malformed", []byte(`{"Cases":[`)},
		{"null", []byte("null")},
		{"too-large", bytes.Repeat([]byte("x"), int(maxDemoDatasetBytes)+1)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := filepath.Join(t.TempDir(), "invalid.json")
			if err := os.WriteFile(p, tc.input, 0o600); err != nil {
				t.Fatal(err)
			}
			if _, err := loadDemoDataset(p); err == nil {
				t.Fatal("untrusted or malformed dataset accepted")
			}
		})
	}
	if _, err := loadDemoDataset(t.TempDir()); err == nil {
		t.Fatal("directory was accepted as dataset")
	}
}

func TestImportedDatasetCLIUsesRealLoopbackRoute(t *testing.T) {
	path := writeDemoDataset(t, importedCaseDataset())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log := &demoCapture{ready: make(chan struct{})}
	done := make(chan error, 1)
	go func() {
		done <- runDemoCLI(ctx, []string{"-principal=investigator-42", "-dataset="+path}, log)
	}()
	select {
	case <-log.ready:
	case err := <-done:
		t.Fatalf("imported dataset server exited: %v", err)
	case <-time.After(10 * time.Second):
		t.Fatal("imported dataset server did not start")
	}
	var endpoint, token string
	for _, line := range strings.Split(log.String(), "\n") {
		if strings.HasPrefix(line, "URL: ") {
			endpoint = strings.TrimPrefix(line, "URL: ")
		}
		if strings.HasPrefix(line, "Bearer token (this run only): ") {
			token = strings.TrimPrefix(line, "Bearer token (this run only): ")
		}
	}
	parsed, err := url.Parse(endpoint)
	if err != nil || parsed.Hostname() != "127.0.0.1" ||
		parsed.Query().Get("case_id") != "case-imported" || len(token) != 64 {
		t.Fatalf("incorrect private dataset endpoint: %q", endpoint)
	}
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{Timeout: 3 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	response.Body.Close()
	if err != nil || response.StatusCode != 200 || !bytes.Contains(body, []byte("event-old")) ||
		bytes.Contains(body, []byte("demo-hidden-source")) {
		t.Fatalf("imported Case HTTP flow failed: %v, %d, %s", err, response.StatusCode, body)
	}
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(6 * time.Second):
		t.Fatal("imported snapshot server did not shut down")
	}
}
