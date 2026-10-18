package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Some output adapters, including imperfect integrations with external
// consumers, violate io.Writer's contract by reporting n < len(p) with no
// error. An authorized forensic receipt must never claim success then.
type silentlyShortReceiptWriter struct{}

func (silentlyShortReceiptWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	return len(p) / 2, nil
}

func TestForensicReceiptsDoNotReportSuccessOnSilentTruncation(t *testing.T) {
	const analyst = "analyst-a"
	mixedFile, mixedCase, mixedSHA := savedHeadlessMixedCase(t)
	for _, test := range []struct {
		name string
		run func(io.Writer) error
	}{
		{"headless authorized Event Search", func(out io.Writer) error {
			return runSavedCaseHeadlessHunt(context.Background(), analyst, mixedFile,
				mixedCase, mixedSHA, "source_ref", "all", "local-jsonl-source-",
				"2026-10-17T09:59:00Z", "2026-10-17T10:02:00Z", out)
		}},
		{"verified Case inventory", func(out io.Writer) error {
			return verifySavedNormalizedCase(context.Background(), analyst, mixedFile,
				mixedCase, mixedSHA, out)
		}},
		{"two independent Case comparison", func(out io.Writer) error {
			firstFile, firstCase, firstSHA, secondFile, secondCase, secondSHA :=
				buildTwoVerifiedLocalCaseArchives(t)
			return compareVerifiedSavedCases(context.Background(), analyst,
				firstFile, firstCase, firstSHA,
				secondFile, secondCase, secondSHA, out)
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			// Sanity check: this is a real authorized workflow, not an error
			// in the input Case that happens to hide a bad output writer.
			var complete bytes.Buffer
			if err := test.run(&complete); err != nil || !json.Valid(complete.Bytes()) {
				t.Fatalf("expected valid complete forensic receipt: %v %s", err, complete.String())
			}
			err := test.run(silentlyShortReceiptWriter{})
			if !errors.Is(err, io.ErrShortWrite) {
				t.Fatalf("partial forensic output incorrectly marked as complete: %v", err)
			}
		})
	}
}

func TestMixedAuthorizedEventSearchRejectsURLQuerySmugglingAndWrongMediaType(t *testing.T) {
	path, caseID, digest := savedHeadlessMixedCase(t)
	data, imported, err := loadAuthorizedImportedCase(context.Background(), "analyst-a", path, caseID, digest)
	if err != nil {
		t.Fatal(err)
	}
	token, err := randomBearer()
	if err != nil {
		t.Fatal(err)
	}
	owner, err := newDemoHandlerWithData(token, "analyst-a", data)
	if err != nil {
		t.Fatal(err)
	}
	handler := newDemoLocalJSONLSearchHandler(token, "analyst-a", caseID, data, imported, owner)
	selection := []byte(`{"case_id":"` + caseID + `","source_family":"all","field":"source_ref","value":"local-jsonl-source-","from":"2026-10-17T09:59:00Z","to":"2026-10-17T10:02:00Z"}`)
	execute := func(uri, mediaType string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, uri, bytes.NewReader(selection))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", mediaType)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, req)
		return response
	}
	// This search is authorized and completes with actual mixed-source
	// correlation when the sensitive selection lives ONLY in the POST body.
	allowed := execute("/demo-local-search", "application/json")
	if allowed.Code != http.StatusOK || !strings.Contains(allowed.Body.String(), "mixed_temporal_context") {
		t.Fatalf("authorized body-only mixed search unavailable: %d %s", allowed.Code, allowed.Body.String())
	}
	for _, tc := range []struct{uri, contentType string}{
		{"/demo-local-search?value=PRIVATE-GENERIC-NEAR", "application/json"},
		{"/demo-local-search?case_id="+caseID, "application/json"},
		{"/demo-local-search", "text/plain"},
		{"/demo-local-search", "application/x-www-form-urlencoded"},
	} {
		denied := execute(tc.uri, tc.contentType)
		if denied.Code == http.StatusOK || strings.Contains(denied.Body.String(), "local-jsonl-entry-") ||
			strings.Contains(denied.Body.String(), "PRIVATE-GENERIC-NEAR") ||
			strings.Contains(denied.Body.String(), "Bearer "+token) {
			t.Fatalf("sensitive URL/media-type request disclosed authorized search results: %s %d %q",
				tc.uri, denied.Code, denied.Body.String())
		}
	}
}
