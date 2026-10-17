package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

func TestOfflineCLIProducesAuthorizedCaseTriageWithoutListener(t *testing.T) {
	var out bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{"-principal=analyst-a", "-report-case=case-a"}, &out); err != nil {
		t.Fatal(err)
	}
	var report caseTriageReceipt
	if err := json.Unmarshal(out.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if report.Format != "cmdr-local-case-triage-v1" || !report.Synthetic ||
		report.CaseID != "case-a" || report.VisibleEntries != 2 ||
		report.Kinds["observed"] != 1 || report.Kinds["inferred"] != 1 ||
		report.UncertainEntries != 0 || report.DistinctSources != 2 ||
		report.FirstCertainUTC != "2026-10-01T10:00:00Z" ||
		report.LastCertainUTC != "2026-10-01T11:00:00Z" || len(report.TopSources) != 2 {
		t.Fatalf("incomplete authorized triage: %+v", report)
	}
	for _, private := range []string{
		testBearer, "Bearer", "demo-hidden-source", "demo-secret",
		"restricted", "event-new", "event-old", "demo-tenant-b",
		"Local analyst viewer:", "http://", "DO NOT DISPLAY",
	} {
		if strings.Contains(out.String(), private) {
			t.Fatalf("report contained credential/raw row/forbidden source %q: %s", private, out.String())
		}
	}
	if !strings.Contains(out.String(), "demo-evidence-a") {
		t.Fatalf("permitted sources absent from triage: %s", out.String())
	}
	var other bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{"-principal=analyst-b", "-report-case=case-b"}, &other); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(other.String(), "demo-evidence-a") || strings.Contains(other.String(), "case-a") {
		t.Fatalf("Case B report leaked Case A: %s", other.String())
	}
}

func TestOfflineCLIUsesActualOperatorDiskOwnerSnapshot(t *testing.T) {
	path := snapshotFixture(t)
	var out bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{
		"-principal=analyst-a", "-snapshot=" + path, "-report-case=forensic-case-on-disk",
	}, &out); err != nil {
		t.Fatal(err)
	}
	var got caseTriageReceipt
	if err := json.Unmarshal(out.Bytes(), &got); err != nil ||
		got.CaseID != "forensic-case-on-disk" || got.VisibleEntries != 2 {
		t.Fatalf("disk Case was not actually aggregated: %+v %v", got, err)
	}
	if strings.Contains(out.String(), "case-b") || strings.Contains(out.String(), "demo-hidden-source") {
		t.Fatalf("disk report leaked another tenant or denied source: %s", out.String())
	}
}

func TestOfflineCLIReadsAllAuthorizedPagesAndAggregatesKinds(t *testing.T) {
	data := demoData()
	data.Cases[0].EntryIDs = []string{"restricted"}
	for i := 0; i < 205; i++ {
		id := fmt.Sprintf("disk-entry-%03d", i)
		kind := "observed"
		if i%2 != 0 {
			kind = "inferred"
		}
		data.Cases[0].EntryIDs = append(data.Cases[0].EntryIDs, id)
		data.Entries = append(data.Entries, timeline.MemoryTimelineEntry{
			CaseID: "case-a",
			Value: timeline.Entry{
				ID: id, Tenant: "demo-tenant-a", Kind: kind,
				SourceTime: time.Date(2026, 10, 1, 10, 0, 0, 0, time.UTC).
					Add(time.Duration(i) * time.Second).Format(time.RFC3339),
				Timezone: "+00:00", Version: 1, SourceRef: "demo-evidence-a",
				Provenance: "operator-disk-telemetry",
			},
			Readers: []string{"analyst-a"},
		})
	}
	payload, err := json.Marshal(data)
	if err != nil { t.Fatal(err) }
	path := filepath.Join(t.TempDir(), "operator-owner.json")
	if err := os.WriteFile(path, payload, 0o600); err != nil { t.Fatal(err) }
	var output bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{
		"-snapshot=" + path, "-principal=analyst-a", "-report-case=case-a",
	}, &output); err != nil {
		t.Fatal(err)
	}
	var got caseTriageReceipt
	if err := json.Unmarshal(output.Bytes(), &got); err != nil { t.Fatal(err) }
	if got.VisibleEntries != 205 || got.Kinds["observed"] != 103 ||
		got.Kinds["inferred"] != 102 || got.DistinctSources != 1 ||
		got.TopSources[0].Reference != "demo-evidence-a" ||
		got.TopSources[0].Count != 205 ||
		got.FirstCertainUTC != "2026-10-01T10:00:00Z" ||
		got.LastCertainUTC != "2026-10-01T10:03:24Z" {
		t.Fatalf("incomplete multi-page operator Case report: %+v", got)
	}
	if strings.Contains(output.String(), "disk-entry-") ||
		strings.Contains(output.String(), "restricted") ||
		strings.Contains(output.String(), "demo-hidden-source") {
		t.Fatalf("receipt exposed raw Entry IDs or forbidden source: %s", output.String())
	}
}

type rejectReportOutput struct{}
func (rejectReportOutput) Write([]byte) (int,error) { return 0, io.ErrClosedPipe }

func TestOfflineCLIRejectsUnauthorizedIncompleteAndFailedOutput(t *testing.T) {
	for _, args := range [][]string{
		{"-principal=analyst-a", "-report-case=case-b"},
		{"-principal=analyst-b", "-report-case=case-a"},
		{"-principal=analyst-a", "-report-case=missing"},
		{"-principal=invalid", "-report-case=case-a"},
		{"-principal=analyst-a", "-report-case=case-a", "-snapshot=relative.json"},
		{"-principal=analyst-a", "-report-case=" + strings.Repeat("x", 1100)},
	} {
		var out bytes.Buffer
		if err := runDemoCLI(context.Background(), args, &out); err == nil || out.Len() != 0 {
			t.Fatalf("unauthorized/invalid report yielded output: %v %s", err, out.String())
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var out bytes.Buffer
	if err := writeDemoCaseTriage(ctx, "analyst-a", "", "case-a", &out); !errors.Is(err, context.Canceled) || out.Len() != 0 {
		t.Fatalf("cancelled report leaked: %v %q", err, out.String())
	}
	if err := writeDemoCaseTriage(context.Background(), "analyst-a", "", "case-a", rejectReportOutput{}); !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("writer failure masked: %v", err)
	}
}

func TestOfflineCLIEmptyReadableCaseHasNoFalseEvidence(t *testing.T) {
	data := demoData()
	data.Cases[0].EntryIDs = []string{"restricted"}
	payload, err := json.Marshal(data)
	if err != nil { t.Fatal(err) }
	path := filepath.Join(t.TempDir(), "empty-owner.json")
	if err := os.WriteFile(path,payload,0o600); err != nil { t.Fatal(err) }
	// The Case can remain readable while every entry is denied.
	var out bytes.Buffer
	if err := writeDemoCaseTriage(context.Background(),"analyst-a",path,"case-a",&out);err != nil {
		t.Fatal(err)
	}
	var r caseTriageReceipt
	if err := json.Unmarshal(out.Bytes(),&r);err != nil {t.Fatal(err)}
	if r.VisibleEntries!=0 || len(r.TopSources)!=0 || r.FirstCertainUTC!="" ||
		r.UncertainEntries!=0 || r.DistinctSources!=0 {
		t.Fatalf("denied-only Case produced false evidence: %+v",r)
	}
}

func TestOfflineTriageFailsClosedOnOwnerAndPaginationResponses(t *testing.T) {
	entry := timeline.ProjectedEntry{
		ID: "visible", Kind: "observed", SourceRef: "authorized-source",
		Provenance: "owner-rechecked", Version: 1,
		SourceTime: "2026-10-01T10:00:00Z", Timezone: "+00:00",
	}
	base := caseTriagePage{Status: "allow", Entries: []timeline.ProjectedEntry{entry}}
	base.Page.Offset, base.Page.Limit, base.Page.TotalVisible = 0, 200, 1
	cases := []struct {
		name string
		mutate func(*caseTriagePage)
	}{
		{"total-too-high", func(p *caseTriagePage) {p.Page.TotalVisible=timeline.MaxEntries+1}},
		{"no-entries", func(p *caseTriagePage) {p.Entries=nil}},
		{"offset-mismatch", func(p *caseTriagePage) {p.Page.Offset=100}},
		{"limit-mismatch", func(p *caseTriagePage) {p.Page.Limit=201}},
		{"has-more-lie", func(p *caseTriagePage) {p.Page.HasMore=true}},
		{"duplicate-visible-IDs", func(p *caseTriagePage) {p.Entries=append(p.Entries,p.Entries[0]);p.Page.TotalVisible=2}},
		{"empty-visible-ID",func(p *caseTriagePage) {p.Entries[0].ID=""}},
		{"empty-source",func(p *caseTriagePage) {p.Entries[0].SourceRef=""}},
		{"empty-provenance",func(p *caseTriagePage) {p.Entries[0].Provenance=""}},
		{"invalid-version",func(p *caseTriagePage) {p.Entries[0].Version=0}},
		{"unsupported-kind",func(p *caseTriagePage) {p.Entries[0].Kind="forged"}},
		{"invalid-known-timestamp",func(p *caseTriagePage) {p.Entries[0].SourceTime="bad"}},
	}
	for _,tc:=range cases {
		t.Run(tc.name,func(t *testing.T){
			p:=base
			p.Entries=append([]timeline.ProjectedEntry(nil),base.Entries...)
			tc.mutate(&p)
			h:=http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
				w.Header().Set("Content-Type","application/json")
				if r.URL.Path=="/demo-cases" {_,_=w.Write([]byte(`{"cases":["case-a"]}`));return}
				_ = json.NewEncoder(w).Encode(p)
			})
			if _,err:=collectDemoCaseTriage(context.Background(),h,testBearer,"case-a");err==nil{
				t.Fatal("accepted malformed authorized page")
			}
		})
	}
}

func TestOfflineTriageNoAuthorityAndMalformedHTTPNeverBecomesReport(t *testing.T) {
	for _,tc:=range []struct{name string; body string; code int}{
		{"unauthorized","{}",http.StatusUnauthorized},
		{"invalid-json","not-json",http.StatusOK},
		{"missing-case",`{"cases":["case-b"]}`,http.StatusOK},
		{"oversized-owner-response",strings.Repeat("X",maxDemoSnapshotBytes+5000),http.StatusOK},
	}{
		t.Run(tc.name,func(t *testing.T){
			h:=http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
				w.WriteHeader(tc.code)
				_,_=io.WriteString(w,tc.body)
			})
			if _,err:=collectDemoCaseTriage(context.Background(),h,testBearer,"case-a");err==nil{
				t.Fatal("unauthenticated or invalid owner response accepted")
			}
		})
	}
	ctx,cancel:=context.WithCancel(context.Background())
	cancel()
	h:=http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		_,_=w.Write([]byte(`{"cases":["case-a"]}`))
	})
	if _,err:=collectDemoCaseTriage(ctx,h,testBearer,"case-a");!errors.Is(err,context.Canceled){
		t.Fatalf("cancelled owner read was accepted: %v",err)
	}
}

func TestOfflineTriageUncertainAuthorizedEntryHasNoFabricatedTime(t *testing.T) {
	entry:=timeline.ProjectedEntry{ID:"uncertain",Kind:"inferred",
		SourceRef:"authorized-source",Provenance:"local-observation",Version:1,
		TimeUncertain:true, SourceTime:"not-an-instant"}
	p:=caseTriagePage{Status:"uncertain",Entries:[]timeline.ProjectedEntry{entry}}
	p.Page.Limit,p.Page.TotalVisible=200,1
	h:=http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		if r.URL.Path=="/demo-cases" {_,_=w.Write([]byte(`{"cases":["case-a"]}`));return}
		_ = json.NewEncoder(w).Encode(p)
	})
	r,err:=collectDemoCaseTriage(context.Background(),h,testBearer,"case-a")
	if err!=nil {t.Fatal(err)}
	if r.UncertainEntries!=1 || r.FirstCertainUTC!="" || r.LastCertainUTC!="" ||
		r.Kinds["inferred"]!=1 || r.Visibility!="uncertain" {
		t.Fatalf("uncertain data became asserted chronology: %+v",r)
	}
}

// A revoked or inconsistent *second page* must never produce a partial report.
func TestOfflineTriageRejectsBrokenSecondAuthorizedPage(t *testing.T) {
	var first caseTriagePage
	first.Status="allow"
	first.Page.Limit=200
	first.Page.TotalVisible=201
	first.Page.HasMore=true
	for i:=0;i<200;i++ {
		first.Entries=append(first.Entries,timeline.ProjectedEntry{
			ID:fmt.Sprintf("row-%03d",i), Kind:"observed",
			SourceRef:"authorized-source", Provenance:"owner",Version:1,
			Timezone:"+00:00",SourceTime:"2026-10-01T10:00:00Z",
		})
	}
	h:=http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		if r.URL.Path=="/demo-cases" {_,_=w.Write([]byte(`{"cases":["case-a"]}`));return}
		if r.URL.Query().Get("offset")=="0" {_=json.NewEncoder(w).Encode(first);return}
		w.WriteHeader(http.StatusNotFound)
		_,_=w.Write([]byte(`{"error":"request-unavailable"}`))
	})
	if _,err:=collectDemoCaseTriage(context.Background(),h,testBearer,"case-a");err==nil{
		t.Fatal("second-page ACL revocation led to partial successful report")
	}
	// A successful second page must also keep the same visible total.
	second:=caseTriagePage{Status:"allow",Entries:[]timeline.ProjectedEntry{
		{ID:"row-200",Kind:"observed",SourceRef:"authorized-source",Provenance:"owner",
		Version:1,Timezone:"+00:00",SourceTime:"2026-10-01T10:00:01Z"},
	}}
	second.Page.Offset,second.Page.Limit,second.Page.TotalVisible=200,200,202
	h=http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		if r.URL.Path=="/demo-cases" {_,_=w.Write([]byte(`{"cases":["case-a"]}`));return}
		if r.URL.Query().Get("offset")=="0" {_=json.NewEncoder(w).Encode(first);return}
		_ = json.NewEncoder(w).Encode(second)
	})
	if _,err:=collectDemoCaseTriage(context.Background(),h,testBearer,"case-a");err==nil{
		t.Fatal("changed total across Case pages accepted")
	}
}
