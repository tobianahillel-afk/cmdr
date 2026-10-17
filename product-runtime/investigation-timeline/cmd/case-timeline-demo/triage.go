package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"time"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// caseTriageReceipt is an explicitly requested, unsigned LOCAL report.
// Only already-authorized entry aggregates are represented: no bearer,
// raw Timeline Entry, source preview, tenant or environment is exported.
type caseTriageReceipt struct {
	Format          string         `json:"format"`
	CaseID          string         `json:"case_id"`
	Synthetic       bool           `json:"synthetic"`
	VisibleEntries  int            `json:"visible_entries"`
	Kinds           map[string]int `json:"kinds"`
	UncertainEntries int           `json:"uncertain_entries"`
	FirstCertainUTC string         `json:"first_certain_utc,omitempty"`
	LastCertainUTC  string         `json:"last_certain_utc,omitempty"`
	DistinctSources int            `json:"distinct_source_references"`
	TopSources      []sourceCount  `json:"top_sources"`
	Visibility      string         `json:"visibility"`
}

type sourceCount struct {
	Reference string `json:"reference"`
	Count     int    `json:"count"`
}

type caseTriagePage struct {
	Status string                    `json:"status"`
	Entries []timeline.ProjectedEntry `json:"entries"`
	Page struct {
		Offset int `json:"offset"`
		Limit int `json:"limit"`
		TotalVisible int `json:"total_visible"`
		HasMore bool `json:"has_more"`
	} `json:"page"`
}

// queryDemoOwner performs a real auth/ACL handler dispatch in process, without
// starting a listener, returning credentials or trusting caller permission flags.
func queryDemoOwner(ctx context.Context, handler http.Handler, token, path string, out any) error {
	req := httptest.NewRequest(http.MethodGet, path, nil).WithContext(ctx)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if w.Code != http.StatusOK || w.Body.Len() > maxDemoSnapshotBytes+4096 {
		return errors.New("authorized Case request unavailable")
	}
	if err := json.Unmarshal(w.Body.Bytes(), out); err != nil {
		return errors.New("invalid authorized Case response")
	}
	return nil
}

// writeDemoCaseTriage reads the operator-selected snapshot and independently
// rechecks Case membership and all server-side Case/Entry/source permissions
// through the existing handlers. No data is written until all pages validate.
func writeDemoCaseTriage(ctx context.Context, principal, snapshot, caseID string, output io.Writer) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if caseID == "" || len(caseID) > 1024 || principal == "" {
		return errors.New("explicit Case and principal required for triage")
	}
	data, err := loadDemoSnapshot(snapshot)
	if err != nil {
		return err
	}
	token, err := randomBearer()
	if err != nil {
		return err
	}
	handler, err := newDemoHandlerWithData(token, principal, data)
	if err != nil {
		return err
	}
	result, err := collectDemoCaseTriage(ctx, handler, token, caseID)
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return json.NewEncoder(output).Encode(result)
}

// collectDemoCaseTriage is separated from the local file reader so its
// authorization and pagination failure behavior can be tested against
// realistic failing HTTP handlers. It never trusts the reported page counts.
func collectDemoCaseTriage(ctx context.Context, handler http.Handler, token, caseID string) (caseTriageReceipt, error) {
	var discovery struct { Cases []string `json:"cases"` }
	if err := queryDemoOwner(ctx, handler, token, "/demo-cases", &discovery); err != nil {
		return caseTriageReceipt{}, err
	}
	authorized := false
	for _, id := range discovery.Cases {
		if id == caseID { authorized = true; break }
	}
	if !authorized {
		return caseTriageReceipt{}, errors.New("Case not available to selected principal")
	}

	receipt := caseTriageReceipt{
		Format: "cmdr-local-case-triage-v1",
		CaseID: caseID, Synthetic: true,
		Kinds: make(map[string]int),
		TopSources: make([]sourceCount,0),
	}
	const pageSize = 200
	seen := make(map[string]bool)
	sources := make(map[string]int)
	total := -1
	var earliest, latest time.Time
	for offset := 0; ; offset += pageSize {
		if offset > timeline.MaxEntries {
			return caseTriageReceipt{}, errors.New("Case exceeds bounded triage capacity")
		}
		selection := url.Values{
			"case_id": {caseID}, "limit": {"200"},
			"offset": {fmt.Sprint(offset)},
		}
		var page caseTriagePage
		if err := queryDemoOwner(ctx, handler, token, "/case-timeline?"+selection.Encode(), &page); err != nil {
			return caseTriageReceipt{}, err
		}
		if total < 0 {
			total = page.Page.TotalVisible
			receipt.Visibility = page.Status
		}
		if page.Page.TotalVisible != total || total < 0 || total > timeline.MaxEntries ||
			page.Page.Offset != offset || page.Page.Limit != pageSize ||
			page.Status != receipt.Visibility ||
			len(page.Entries) > pageSize || offset+len(page.Entries) > total ||
			(len(page.Entries) == 0 && offset < total) ||
			page.Page.HasMore != (offset+len(page.Entries) < total) {
			return caseTriageReceipt{}, errors.New("inconsistent authorized Case pagination")
		}
		for _, entry := range page.Entries {
			if entry.ID == "" || seen[entry.ID] || entry.SourceRef == "" ||
				entry.Provenance == "" || entry.Version < 1 {
				return caseTriageReceipt{}, errors.New("invalid or duplicated authorized Timeline Entry")
			}
			switch entry.Kind {
			case "observed", "inferred", "user-action", "decision", "response":
			default:
				return caseTriageReceipt{}, errors.New("invalid authorized Timeline classification")
			}
			seen[entry.ID] = true
			receipt.Kinds[entry.Kind]++
			sources[entry.SourceRef]++
			if entry.TimeUncertain {
				receipt.UncertainEntries++
				continue
			}
			instant, err := time.Parse(time.RFC3339Nano, entry.SourceTime)
			if err != nil {
				return caseTriageReceipt{}, errors.New("invalid known Timeline timestamp")
			}
			instant = instant.UTC()
			if earliest.IsZero() || instant.Before(earliest) { earliest = instant }
			if latest.IsZero() || instant.After(latest) { latest = instant }
		}
		if !page.Page.HasMore {
			break
		}
	}
	if len(seen) != total || ctx.Err() != nil {
		return caseTriageReceipt{}, errors.New("incomplete or cancelled authorized Case triage")
	}
	receipt.VisibleEntries = total
	receipt.DistinctSources = len(sources)
	if !earliest.IsZero() {
		receipt.FirstCertainUTC = earliest.Format(time.RFC3339Nano)
		receipt.LastCertainUTC = latest.Format(time.RFC3339Nano)
	}
	for ref, count := range sources {
		receipt.TopSources = append(receipt.TopSources, sourceCount{Reference: ref, Count: count})
	}
	sort.Slice(receipt.TopSources,func(i,j int) bool {
		a,b:=receipt.TopSources[i],receipt.TopSources[j]
		if a.Count!=b.Count{return a.Count>b.Count}
		return a.Reference<b.Reference
	})
	if len(receipt.TopSources)>5 {receipt.TopSources=receipt.TopSources[:5]}
	return receipt, nil
}
