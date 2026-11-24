package tier1_features

import (
	"context"
	"testing"
	"time"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

// T1.F1.1: Verify PR #225 DraftBindings review thread is marked resolved in GitHub or local audit records.
func TestT1_F1_1_CodexPR225DraftBindingsThreadResolved(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := harness.NewGitHubClient()
	if err != nil {
		t.Fatalf("failed to initialize github client: %v", err)
	}

	resolved, isOffline, err := client.VerifyCodexPRResolution(ctx, 225)
	if err != nil {
		t.Fatalf("verification of PR #225 resolution failed: %v", err)
	}
	if !resolved {
		t.Fatalf("PR #225 DraftBindings thread is not marked resolved (offline=%v)", isOffline)
	}
	t.Logf("T1.F1.1 PASS: PR #225 review thread resolved (offline fallback: %v)", isOffline)
}

// T1.F1.2: Verify PR #228 secrets non-disclosure regression tests pass in studio-catalog.
func TestT1_F1_2_CodexPR228SecretsNonDisclosure(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	passed, output, err := harness.RunGoTest(ctx, "product-runtime/studio-catalog", "-run", "TestToolInput_AC2_SecretReferenceMasking_And_SettingsOwnership")
	if !passed || err != nil {
		t.Fatalf("PR #228 secrets non-disclosure regression test failed: %v\nOutput: %s", err, output)
	}
	t.Log("T1.F1.2 PASS: secrets non-disclosure regression suite passed with zero leaks")
}

// T1.F1.3: Verify optional default values validation regression test passes in studio-catalog.
func TestT1_F1_3_CodexOptionalDefaultValuesValidation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	passed, output, err := harness.RunGoTest(ctx, "product-runtime/studio-catalog", "-run", "TestChallenge_ProjectedDefaultBindings_ReferenceIDs")
	if !passed || err != nil {
		t.Fatalf("optional default values validation test failed: %v\nOutput: %s", err, output)
	}
	t.Log("T1.F1.3 PASS: optional default values regression test passed")
}

// T1.F1.4: Verify deep copying of mutable structures regression test passes in studio-catalog.
func TestT1_F1_4_CodexDeepCopyingMutableStructures(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	passed, output, err := harness.RunGoTest(ctx, "product-runtime/studio-catalog", "-run", "TestAdversarial_CloneDeepCopyImmunity")
	if !passed || err != nil {
		t.Fatalf("deep copy mutable structures regression test failed: %v\nOutput: %s", err, output)
	}
	t.Log("T1.F1.4 PASS: mutable structures deep copying regression test passed")
}

// T1.F1.5: Verify typed nil and null map entries preservation regression tests pass in studio-catalog.
func TestT1_F1_5_CodexTypedNilAndNullMapEntries(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	passed1, output1, err1 := harness.RunGoTest(ctx, "product-runtime/studio-catalog", "-run", "TestChallenge_DeepCopy_TypedNilPointers")
	if !passed1 || err1 != nil {
		t.Fatalf("typed nil pointers regression test failed: %v\nOutput: %s", err1, output1)
	}

	passed2, output2, err2 := harness.RunGoTest(ctx, "product-runtime/studio-catalog", "-run", "TestChallenge_DeepCopy_MapWithNilValues")
	if !passed2 || err2 != nil {
		t.Fatalf("map with nil values preservation test failed: %v\nOutput: %s", err2, output2)
	}

	t.Log("T1.F1.5 PASS: typed nil and null map entries preservation regression tests passed")
}
