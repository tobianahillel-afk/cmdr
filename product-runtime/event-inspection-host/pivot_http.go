package inspectionhost

import (
	"encoding/json"
	"io"
	"mime"
	"net/http"

	inspection "github.com/tobianahillel-afk/cmdr/product-runtime/event-inspection"
)

const pivotRoute = "/investigate/event-inspection/pivot"

// pivotHTTPRequest is an untrusted selection. Access decisions and the return
// destination are intentionally absent: both must be determined by the host.
type pivotHTTPRequest struct {
	Request
	Field        string `json:"field"`
	Value        string `json:"value"`
	ValuePresent bool   `json:"valuePresent"`
	TimeStart    string `json:"timeStart"`
	TimeEnd      string `json:"timeEnd"`
}

type browserPivotTime struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type browserPivotReturn struct {
	TenantRef     string `json:"tenantRef"`
	EnvironmentRef string `json:"environmentRef"`
	EventRef      string `json:"eventRef"`
	CorrelationID string `json:"correlationId"`
	Origin        string `json:"origin"`
}

type browserPivotDraft struct {
	Kind             string             `json:"kind"`
	TargetCapability string             `json:"targetCapability"`
	TenantRef        string             `json:"tenantRef"`
	EnvironmentRef   string             `json:"environmentRef"`
	EventRef         string             `json:"eventRef"`
	JobID            string             `json:"jobId"`
	RunID            string             `json:"runId"`
	CorrelationID    string             `json:"correlationId"`
	Field            string             `json:"field"`
	Value            string             `json:"value"`
	FieldOrigin      string             `json:"fieldOrigin"`
	SourceRef        string             `json:"sourceRef"`
	DataSourceRef    string             `json:"dataSourceRef"`
	TimeRange        browserPivotTime   `json:"timeRange"`
	ReturnContext    browserPivotReturn `json:"returnContext"`
	AuditRequired    bool               `json:"auditRequired"`
}

func (h *handler) servePivot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Referrer-Policy", "no-referrer")
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		unavailable(w, http.StatusMethodNotAllowed)
		return
	}
	if r.URL == nil || r.URL.RawQuery != "" || r.Body == nil {
		unavailable(w, http.StatusBadRequest)
		return
	}
	media, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || media != "application/json" {
		unavailable(w, http.StatusUnsupportedMediaType)
		return
	}
	var q pivotHTTPRequest
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxRequestBytes))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&q); err != nil {
		unavailable(w, http.StatusBadRequest)
		return
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF ||
		!validRef(q.TenantRef) || q.TenantRef == "*" ||
		!validRef(q.EnvironmentRef) || !validRef(q.EventRef) ||
		!validRef(q.JobID) || !validRef(q.RunID) ||
		!validRef(q.CorrelationID) || !validRef(q.SourceHint) ||
		!validField(q.Field) || !q.ValuePresent {
		unavailable(w, http.StatusBadRequest)
		return
	}
	if r.Context().Err() != nil {
		unavailable(w, http.StatusServiceUnavailable)
		return
	}
	ctx, projected, _, authorized := h.authorizedInspection(r, q.Request)
	if !authorized {
		unavailable(w, http.StatusNotFound)
		return
	}
	// The server has independently projected *only* the Telemetry owner's
	// field-authorized source/normalized/raw-derived/enriched maps. The core
	// additionally enforces exact map-key presence and selected value match.
	// No client "valueAuthorized" or arbitrary return-origin is accepted.
	result := inspection.PreparePivot(projected.Projection, inspection.PivotSelection{
		Field: q.Field, Value: q.Value, ValuePresent: q.ValuePresent,
		ValueAuthorized: true, TimeStart: q.TimeStart, TimeEnd: q.TimeEnd,
		ReturnOrigin: "event-search",
	})
	if !result.Allowed || result.Draft == nil || ctx.Err() != nil || r.Context().Err() != nil {
		unavailable(w, http.StatusNotFound)
		return
	}
	d := result.Draft
	response := browserPivotDraft{
		Kind: "pivot-draft", TargetCapability: d.TargetCapability,
		TenantRef: d.TenantRef, EnvironmentRef: d.EnvironmentRef,
		EventRef: d.EventRef, JobID: q.JobID, RunID: q.RunID,
		CorrelationID: d.CorrelationID, Field: d.Field, Value: d.Value,
		FieldOrigin: d.FieldOrigin, SourceRef: d.SourceRef,
		DataSourceRef: d.DataSourceRef,
		TimeRange: browserPivotTime{Start: d.TimeRange.Start, End: d.TimeRange.End},
		ReturnContext: browserPivotReturn{
			TenantRef: d.ReturnContext.TenantRef,
			EnvironmentRef: d.ReturnContext.EnvironmentRef,
			EventRef: d.ReturnContext.EventRef,
			CorrelationID: d.ReturnContext.CorrelationID,
			Origin: d.ReturnContext.Origin,
		},
		AuditRequired: result.AuditRequired,
	}
	data, err := json.Marshal(response)
	if err != nil || ctx.Err() != nil || r.Context().Err() != nil {
		unavailable(w, http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(append(data, '\n'))
}
