package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"testing"
)

func TestCommandQueueProbeExecutesActualAuthorizedProjection(t *testing.T) {
	var output bytes.Buffer
	if err := run([]string{"-iterations", "10000", "-operation", "projection", "-mode", "latency"}, &output); err != nil {
		t.Fatalf("in-process Command projection failed: %v", err)
	}
	var receipt observation
	if err := json.Unmarshal(output.Bytes(), &receipt); err != nil {
		t.Fatalf("non-JSON performance receipt: %v", err)
	}
	if receipt.Operation != "projection" || receipt.Mode != "latency" || receipt.Operations != 10000 ||
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
		{"-mode", "resource", "-operation", "saved-view"},
		{"-mode", "resource", "-memory-limit-mib", "128"},
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
	if err := run([]string{"-iterations", "1"}, unavailableReceiptWriter{}); err == nil {
		t.Fatal("failed to report unusable benchmark receipt sink")
	}
}

func TestCommandQueueResourceProbeMeasuresRealProjectionHeap(t *testing.T) {
	var output bytes.Buffer
	if err:=run([]string{"-iterations","100000","-operation","projection","-mode","resource"},&output);err!=nil {
		t.Fatalf("Command resource probe failed: %v",err)
	}
	var sample observation
	if err:=json.Unmarshal(output.Bytes(),&sample);err!=nil {t.Fatal(err)}
	if sample.Mode!="resource" || sample.Operation!="projection" || sample.Operations!=100000 ||
		sample.PeakHeapBytes==0 || sample.TotalAllocationBytes==0 ||
		sample.PeakHeapBytes>128*1024*1024 || sample.ElapsedNS<=0 {
		t.Fatalf("Command probe did not actually record bounded runtime heap/allocations: %+v",sample)
	}
}

func TestCommandQueueResourceProbeDoesNotReturnPartialReceiptOnSinkFailure(t *testing.T) {
	err:=run([]string{"-mode","resource","-iterations","100000"},unavailableReceiptWriter{})
	if err==nil {t.Fatal("failed to reject unavailable sensitive benchmark result sink")}
}
