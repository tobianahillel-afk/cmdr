package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"context"

	queue "github.com/tobianahillel-afk/cmdr/product-runtime/work-queue"
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

 
// The benchmark receipt is security-relevant evidence: a missing permission,
// revoked Incident, or failed owner must NEVER be presented as a measured,
// authorized Command projection. The regular CLI cannot choose its owner;
// this test-only seam exercises actual queue.Project failures.
type failingProbeOwner struct { mode string }

func (o failingProbeOwner) List(ctx context.Context, tenant, env string) ([]queue.Item, error) {
	if o.mode == "owner-unavailable" { return nil, errors.New("Incident/Task source is offline") }
	return (syntheticOwner{}).List(ctx, tenant, env)
}

func (o failingProbeOwner) CanRead(ctx context.Context, principal, tenant, env string, kind queue.Kind, id string) (bool, error) {
	switch o.mode {
	case "permission-backend-failed":
		return false, errors.New("current principal authorization unavailable")
	case "permission-revoked":
		return false, nil
	default:
		return (syntheticOwner{}).CanRead(ctx, principal, tenant, env, kind, id)
	}
}

func TestMeasuredCommandProbeRejectsUnverifiedOwnerEvidence(t *testing.T) {
	for _, mode := range []string{"latency", "resource"} {
		for _, cause := range []string{
			"owner-unavailable", "permission-backend-failed", "permission-revoked",
		} {
			t.Run(mode+"/"+cause, func(t *testing.T) {
				var output bytes.Buffer
				err := runWithOwner([]string{"-mode", mode, "-iterations", "3"}, &output, failingProbeOwner{mode:cause})
				if err == nil || output.Len() != 0 {
					t.Fatalf("incomplete or unauthorized benchmark was reported as a success: %v %q", err, output.String())
				}
			})
		}
	}
	// A missing host-side authority is not interchangeable with an
	// unconditional synthetic Allow; the runtime itself must refuse it.
	var output bytes.Buffer
	if err := runWithOwner([]string{"-iterations", "3"}, &output, nil); err == nil || output.Len() != 0 {
		t.Fatalf("nil owner published a false performance measurement: %v %q", err, output.String())
	}
}


// Benchmark authorization must remain valid THROUGHOUT the measurement.
// Passing an initial warm-up projection cannot be used as authority after a
// second owner read revokes an Incident. No partial success receipt is valid.
type revokingProbeOwner struct { reads int }

func (o *revokingProbeOwner) List(ctx context.Context, tenant, env string) ([]queue.Item, error) {
	return (syntheticOwner{}).List(ctx, tenant, env)
}
func (o *revokingProbeOwner) CanRead(ctx context.Context, principal, tenant, env string, kind queue.Kind, id string) (bool, error) {
	o.reads++
	if o.reads > 1 {
		return false, nil
	}
	return (syntheticOwner{}).CanRead(ctx, principal, tenant, env, kind, id)
}
func TestCommandBenchmarkRevocationDuringSamplingFailsWithoutReceipt(t *testing.T) {
	for _, mode := range []string{"latency", "resource"} {
		t.Run(mode, func(t *testing.T) {
			owner := &revokingProbeOwner{}
			var output bytes.Buffer
			err := runWithOwner([]string{"-iterations","100","-mode",mode}, &output, owner)
			if err == nil || output.Len() != 0 || owner.reads < 2 {
				t.Fatalf("revocation after successful warm-up must abort benchmark: err=%v reads=%d output=%q",
					err, owner.reads, output.String())
			}
		})
	}
}
