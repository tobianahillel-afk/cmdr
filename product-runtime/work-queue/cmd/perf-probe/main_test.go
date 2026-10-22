package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestCommandQueueProbeExecutesActualAuthorizedProjection(t *testing.T) {
	var output bytes.Buffer
	if err := run([]string{"-iterations", "2000", "-operation", "projection", "-mode", "latency"}, &output); err != nil {
		t.Fatalf("in-process Command projection failed: %v", err)
	}
	var receipt observation
	if err := json.Unmarshal(output.Bytes(), &receipt); err != nil {
		t.Fatalf("non-JSON performance receipt: %v", err)
	}
	if receipt.Operation != "projection" || receipt.Mode != "latency" || receipt.Operations < 2000 || receipt.Operations > maxMeasuredOperations ||
		receipt.ElapsedNS <= 0 || receipt.NSPerOperation <= 0 {
		t.Fatalf("invalid measured Command runtime: %+v", receipt)
	}
	// Without any browser, database or provider network dependency, a valid
	// owner-authorized Incident must be projected every iteration.
}

func TestCommandQueueProbeRejectsUnsupportedModesBeforeReceipt(t *testing.T) {
	for _, args := range [][]string{
		{"-iterations", "0"},
		{"-iterations", "100001"},
		{"-iterations", "not-an-integer"},
		{"-mode", "resource"},
		{"-mode", "resource", "-memory-limit-mib", "63"},
		{"-mode", "latency", "-memory-limit-mib", "128"},
		{"-operation", "saved-view"},
		{"-operation", "projection", "unexpected-position"},
		{"-invalid-option"},
	} {
		var output bytes.Buffer
		if err := run(args, &output); err == nil || output.Len() != 0 {
			t.Fatalf("invalid projection request emitted a result: args=%v body=%q", args, output.String())
		}
	}
}

type unavailableReceiptWriter struct{}
func (unavailableReceiptWriter) Write([]byte) (int, error) {
	return 0, errors.New("operator evidence sink unavailable")
}
var _ io.Writer = unavailableReceiptWriter{}

func TestCommandQueueProbeSurfacesReceiptFailure(t *testing.T) {
	if err := run([]string{"-iterations", "100000"}, unavailableReceiptWriter{}); err == nil || !strings.Contains(err.Error(), "operator evidence sink unavailable") {
		t.Fatalf("failed to report unusable benchmark receipt sink: %v", err)
	}
}

func TestCommandQueueProbeResourceMeasuresGoHeap(t *testing.T) {
	var out bytes.Buffer
	if err:=run([]string{"-iterations","10000","-mode","resource","-memory-limit-mib","128"},&out);err!=nil{t.Fatal(err)}
	var v observation
	if err:=json.Unmarshal(out.Bytes(),&v);err!=nil{t.Fatal(err)}
	if v.Mode!="resource"||v.Operations<10000||v.Operations>maxMeasuredOperations||v.ElapsedNS<=0||v.PeakHeapBytes==0||v.TotalAllocationBytes==0 {t.Fatalf("invalid resource measurement: %+v",v)}
}
