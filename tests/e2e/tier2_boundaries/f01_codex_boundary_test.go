package tier2_boundaries

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

const studioCatalogPkg = "github.com/tobianahillel-afk/cmdr/product-runtime/studio-catalog"

// T2.F1.1: Secret injection into DraftBinding fails validation without leakage.
func TestT2_F1_1_DraftBindingSecretInjectionValidation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Run secret disclosure adversarial suite on studio-catalog
	const target = "TestChallenge_DraftBindings_SecretNonDisclosure"
	res, err := harness.RunGoTestVerified(
		ctx,
		"product-runtime/studio-catalog",
		studioCatalogPkg,
		[]string{target},
		"-run", "^"+target+"$",
	)
	if !res.Passed || err != nil {
		t.Fatalf("T2.F1.1 FAIL: DraftBinding secret non-disclosure verification failed: %v\nOutput: %s", err, res.Output)
	}

	// Verify that literal injected secrets are never echoed in cleartext
	sensitivePatterns := []string{
		"synth-secret-bearer-token-live-abc12345",
		"super-secret-password-raw-987654321",
	}
	for _, pattern := range sensitivePatterns {
		if strings.Contains(res.Output, pattern) {
			t.Fatalf("T2.F1.1 FAIL: cleartext sensitive token %q leaked in test execution output!", pattern)
		}
	}
	t.Log("T2.F1.1 PASS: secret injection into DraftBinding strictly fails validation without data leakage")
}

// T2.F1.2: Deep copy isolation (mutating source slice does not mutate target).
func TestT2_F1_2_DeepCopyIsolationSliceMutation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	const parent = "TestChallenge_DeepCopy_NonMutation"
	expectedTests := []string{
		parent,
		parent + "/Modifying_Inspected_Contract_Cannot_Mutate_Stored_Contract",
	}
	res, err := harness.RunGoTestVerified(
		ctx,
		"product-runtime/studio-catalog",
		studioCatalogPkg,
		expectedTests,
		"-run", "^"+parent+"$",
	)
	if !res.Passed || err != nil {
		t.Fatalf("T2.F1.2 FAIL: Deep copy slice/map isolation regression test failed: %v\nOutput: %s", err, res.Output)
	}
	t.Log("T2.F1.2 PASS: deep copy isolation successfully prevents mutations from corrupting stored structures")
}

// T2.F1.3: Typed nil pointer handling does not panic or bypass nil checks.
func TestT2_F1_3_TypedNilPointerHandlingSafe(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	const parent = "TestChallenge_DeepCopy_TypedNilPointers"
	expectedTests := []string{
		parent,
		parent + "/nil_int_ptr",
		parent + "/nil_string_ptr",
		parent + "/nil_bool_ptr",
		parent + "/nil_secretref_ptr",
		parent + "/nil_entityref_ptr",
		parent + "/nil_arbitrary_struct_ptr",
	}
	res, err := harness.RunGoTestVerified(
		ctx,
		"product-runtime/studio-catalog",
		studioCatalogPkg,
		expectedTests,
		"-run", "^"+parent+"$",
	)
	if !res.Passed || err != nil {
		t.Fatalf("T2.F1.3 FAIL: Typed nil pointer challenge test failed: %v\nOutput: %s", err, res.Output)
	}
	t.Log("T2.F1.3 PASS: typed nil pointer variants handled safely without panics or nil bypasses")
}

// T2.F1.4: Null map entries preservation in complex nested structures.
func TestT2_F1_4_NullMapEntriesPreservation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	const parent = "TestChallenge_DeepCopy_MapWithNilValues"
	expectedTests := []string{
		parent,
		parent + "/CustomObject_with_nil_entries",
		parent + "/TypedPtrMap_with_typed_nil_and_non-nil_pointers",
	}
	res, err := harness.RunGoTestVerified(
		ctx,
		"product-runtime/studio-catalog",
		studioCatalogPkg,
		expectedTests,
		"-run", "^"+parent+"$",
	)
	if !res.Passed || err != nil {
		t.Fatalf("T2.F1.4 FAIL: Map with nil values preservation test failed: %v\nOutput: %s", err, res.Output)
	}
	t.Log("T2.F1.4 PASS: null/nil map entries are preserved across complex reflection structures")
}

// T2.F1.5: Optional default values fallback without crash when omitted.
func TestT2_F1_5_OptionalDefaultValuesFallback(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	const parent = "TestChallenge_OptionalDefaults_StrictValidation"
	expectedTests := []string{
		parent,
		parent + "/Invalid_Integer_Default_Fails_Validation",
		parent + "/Disallowed_Enum_Member_Default_Fails_Validation",
		parent + "/Literal_Secret_in_Optional_Default_Fails_Validation",
		parent + "/Valid_Optional_Default_Passes_Validation",
	}
	res, err := harness.RunGoTestVerified(
		ctx,
		"product-runtime/studio-catalog",
		studioCatalogPkg,
		expectedTests,
		"-run", "^"+parent+"$",
	)
	if !res.Passed || err != nil {
		t.Fatalf("T2.F1.5 FAIL: Optional default values validation failed: %v\nOutput: %s", err, res.Output)
	}
	t.Log("T2.F1.5 PASS: optional default values fallback adheres strictly to schema validation without crashes")
}
