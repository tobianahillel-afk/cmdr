package tier1_features

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

// T1.F13.1: Verify go vet passes on all 12 Go runtime modules.
func TestT1_F13_1_AllRuntimeModulesGoVet(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	for _, mod := range harness.RuntimeGoModules {
		passed, out, err := harness.RunGoVet(ctx, mod)
		if !passed || err != nil {
			t.Fatalf("go vet failed on module %s: %v\nOutput: %s", mod, err, out)
		}
	}
	t.Logf("T1.F13.1 PASS: go vet passed cleanly on all %d runtime Go modules", len(harness.RuntimeGoModules))
}

// T1.F13.2: Verify go test -count=1 passes on all 12 Go runtime modules.
func TestT1_F13_2_AllRuntimeModulesGoTest(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	for _, mod := range harness.RuntimeGoModules {
		passed, out, err := harness.RunGoTest(ctx, mod)
		if !passed || err != nil {
			t.Fatalf("go test failed on module %s: %v\nOutput: %s", mod, err, out)
		}
	}
	t.Logf("T1.F13.2 PASS: go test passed on all %d runtime Go modules", len(harness.RuntimeGoModules))
}

// T1.F13.3: Verify tools/cmdr-dev unit tests pass.
func TestT1_F13_3_CmdrDevTestsPass(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	passed, out, err := harness.RunGoTest(ctx, "tools/cmdr-dev")
	if !passed || err != nil {
		t.Fatalf("cmdr-dev unit tests failed: %v\nOutput: %s", err, out)
	}
	t.Log("T1.F13.3 PASS: tools/cmdr-dev unit tests passed cleanly")
}

// T1.F13.4: Verify product-runtime/event-search-frontend tests pass.
func TestT1_F13_4_FrontendTestsPass(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	passed, out, err := harness.RunFrontendTest(ctx)
	if !passed || err != nil {
		t.Fatalf("event-search-frontend tests failed: %v\nOutput: %s", err, out)
	}
	t.Log("T1.F13.4 PASS: product-runtime/event-search-frontend tests passed cleanly")
}

// T1.F13.5: Verify zero race conditions during test campaign.
func TestT1_F13_5_ZeroRaceConditions(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// Run race test on core critical concurrent module studio-catalog
	res, outcome := harness.RunGoTestWithRace(
		ctx,
		"product-runtime/studio-catalog",
		"github.com/tobianahillel-afk/cmdr/product-runtime/studio-catalog",
		[]string{"TestAdversarial_ConcurrentStoreStress"},
	)

	switch outcome.Status {
	case harness.RaceCheckPass:
		t.Log("T1.F13.5 PASS: zero race conditions detected during verified concurrent execution")
	case harness.RaceCheckUnverified:
		t.Skipf("T1.F13.5 UNVERIFIED on %s: %s", runtime.GOOS, outcome.Reason)
	default:
		t.Fatalf("T1.F13.5 FAIL (%s):\n%s", outcome.Reason, res.Output)
	}
}
