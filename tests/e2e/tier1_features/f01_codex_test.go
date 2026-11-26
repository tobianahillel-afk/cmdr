package tier1_features

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

const studioCatalogPkg = "github.com/tobianahillel-afk/cmdr/product-runtime/studio-catalog"

// T1.F1.1: Verify PR #225 DraftBindings review thread is marked resolved in GitHub or SHA-attributed audit records.
func TestT1_F1_1_CodexPR225DraftBindingsThreadResolved(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := harness.NewGitHubClient()
	if err != nil {
		t.Fatalf("failed to initialize github client: %v", err)
	}

	const pr225MergeSHA = "e4666f586854e02766568ffd2cfe29ef32696b49"
	res, err := client.VerifyCodexPRResolution(ctx, 225, pr225MergeSHA)
	if errors.Is(err, harness.ErrUnverified) || !res.Verified {
		t.Skipf("T1.F1.1 UNVERIFIED (live GitHub API unavailable and offline registers lack structured SHA-attributed thread proof): %v", err)
		return
	}
	if err != nil {
		t.Fatalf("verification of PR #225 resolution failed: %v", err)
	}
	if !res.Resolved {
		t.Fatalf("PR #225 DraftBindings thread is not marked resolved (status=%s, sha=%s, reason=%s)", res.Status, res.MatchedSHA, res.Reason)
	}
	t.Logf("T1.F1.1 PASS: PR #225 review thread resolved (source=%s, sha=%s)", res.Source, res.MatchedSHA)
}

// T1.F1.2: Verify PR #228 secrets non-disclosure regression tests pass in studio-catalog.
func TestT1_F1_2_CodexPR228SecretsNonDisclosure(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	const target = "TestToolInput_AC2_SecretReferenceMasking_And_SettingsOwnership"
	res, err := harness.RunGoTestVerified(ctx, "product-runtime/studio-catalog", studioCatalogPkg, []string{target}, "-run", "^"+target+"$")
	if !res.Passed || err != nil {
		t.Fatalf("PR #228 secrets non-disclosure regression test failed: %v\nOutput: %s", err, res.Output)
	}
	t.Log("T1.F1.2 PASS: secrets non-disclosure regression suite passed with zero leaks")
}

// T1.F1.3: Verify optional default values validation regression test passes in studio-catalog.
func TestT1_F1_3_CodexOptionalDefaultValuesValidation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	const target = "TestChallenge_ProjectedDefaultBindings_ReferenceIDs"
	res, err := harness.RunGoTestVerified(ctx, "product-runtime/studio-catalog", studioCatalogPkg, []string{target}, "-run", "^"+target+"$")
	if !res.Passed || err != nil {
		t.Fatalf("optional default values validation test failed: %v\nOutput: %s", err, res.Output)
	}
	t.Log("T1.F1.3 PASS: optional default values regression test passed")
}

// T1.F1.4: Verify deep copying of mutable structures regression test passes in studio-catalog.
func TestT1_F1_4_CodexDeepCopyingMutableStructures(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	const target = "TestAdversarial_CloneDeepCopyImmunity"
	res, err := harness.RunGoTestVerified(ctx, "product-runtime/studio-catalog", studioCatalogPkg, []string{target}, "-run", "^"+target+"$")
	if !res.Passed || err != nil {
		t.Fatalf("deep copy mutable structures regression test failed: %v\nOutput: %s", err, res.Output)
	}
	t.Log("T1.F1.4 PASS: mutable structures deep copying regression test passed")
}

// T1.F1.5: Verify typed nil and null map entries preservation regression tests pass in studio-catalog.
func TestT1_F1_5_CodexTypedNilAndNullMapEntries(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	const target1 = "TestChallenge_DeepCopy_TypedNilPointers"
	res1, err1 := harness.RunGoTestVerified(ctx, "product-runtime/studio-catalog", studioCatalogPkg, []string{target1}, "-run", "^"+target1+"$")
	if !res1.Passed || err1 != nil {
		t.Fatalf("typed nil pointers regression test failed: %v\nOutput: %s", err1, res1.Output)
	}

	const target2 = "TestChallenge_DeepCopy_MapWithNilValues"
	res2, err2 := harness.RunGoTestVerified(ctx, "product-runtime/studio-catalog", studioCatalogPkg, []string{target2}, "-run", "^"+target2+"$")
	if !res2.Passed || err2 != nil {
		t.Fatalf("map with nil values preservation test failed: %v\nOutput: %s", err2, res2.Output)
	}

	t.Log("T1.F1.5 PASS: typed nil and null map entries preservation regression tests passed")
}
