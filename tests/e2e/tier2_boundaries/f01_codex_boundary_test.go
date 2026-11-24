package tier2_boundaries

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

// T2.F1.1: Secret injection into DraftBinding fails validation without leakage.
func TestT2_F1_1_DraftBindingSecretInjectionValidation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Run secret disclosure adversarial suite on studio-catalog
	passed, output, err := harness.RunGoTest(
		ctx,
		"product-runtime/studio-catalog",
		"-run", "TestChallenge_DraftBindings_SecretNonDisclosure",
	)
	if !passed || err != nil {
		t.Fatalf("T2.F1.1 FAIL: DraftBinding secret non-disclosure verification failed: %v\nOutput: %s", err, output)
	}

	// Verify that literal injected secrets are never echoed in cleartext
	sensitivePatterns := []string{
		"synth-secret-bearer-token-live-abc12345",
		"super-secret-password-raw-987654321",
	}
	for _, pattern := range sensitivePatterns {
		if strings.Contains(output, pattern) {
			t.Fatalf("T2.F1.1 FAIL: cleartext sensitive token %q leaked in test execution output!", pattern)
		}
	}
	t.Log("T2.F1.1 PASS: secret injection into DraftBinding strictly fails validation without data leakage")
}

// T2.F1.2: Deep copy isolation (mutating source slice does not mutate target).
func TestT2_F1_2_DeepCopyIsolationSliceMutation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	passed, output, err := harness.RunGoTest(
		ctx,
		"product-runtime/studio-catalog",
		"-run", "TestChallenge_DeepCopy_NonMutation",
	)
	if !passed || err != nil {
		t.Fatalf("T2.F1.2 FAIL: Deep copy slice/map isolation regression test failed: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(output, "Modifying_Inspected_Contract_Cannot_Mutate_Stored_Contract") {
		t.Fatalf("T2.F1.2 FAIL: expected subtest Modifying_Inspected_Contract_Cannot_Mutate_Stored_Contract did not execute")
	}
	t.Log("T2.F1.2 PASS: deep copy isolation successfully prevents mutations from corrupting stored structures")
}

// T2.F1.3: Typed nil pointer handling does not panic or bypass nil checks.
func TestT2_F1_3_TypedNilPointerHandlingSafe(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	passed, output, err := harness.RunGoTest(
		ctx,
		"product-runtime/studio-catalog",
		"-run", "TestChallenge_DeepCopy_TypedNilPointers",
	)
	if !passed || err != nil {
		t.Fatalf("T2.F1.3 FAIL: Typed nil pointer challenge test failed: %v\nOutput: %s", err, output)
	}

	expectedSubtests := []string{
		"nil_int_ptr",
		"nil_string_ptr",
		"nil_bool_ptr",
		"nil_secretref_ptr",
		"nil_entityref_ptr",
		"nil_arbitrary_struct_ptr",
	}
	for _, sub := range expectedSubtests {
		if !strings.Contains(output, sub) {
			t.Fatalf("T2.F1.3 FAIL: expected typed nil subtest %q was not executed", sub)
		}
	}
	t.Log("T2.F1.3 PASS: typed nil pointer variants handled safely without panics or nil bypasses")
}

// T2.F1.4: Null map entries preservation in complex nested structures.
func TestT2_F1_4_NullMapEntriesPreservation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	passed, output, err := harness.RunGoTest(
		ctx,
		"product-runtime/studio-catalog",
		"-run", "TestChallenge_DeepCopy_MapWithNilValues",
	)
	if !passed || err != nil {
		t.Fatalf("T2.F1.4 FAIL: Map with nil values preservation test failed: %v\nOutput: %s", err, output)
	}

	if !strings.Contains(output, "CustomObject_with_nil_entries") ||
		!strings.Contains(output, "TypedPtrMap_with_typed_nil_and_non-nil_pointers") {
		t.Fatalf("T2.F1.4 FAIL: expected null map subtests did not all execute in output: %s", output)
	}
	t.Log("T2.F1.4 PASS: null/nil map entries are preserved across complex reflection structures")
}

// T2.F1.5: Optional default values fallback without crash when omitted.
func TestT2_F1_5_OptionalDefaultValuesFallback(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	passed, output, err := harness.RunGoTest(
		ctx,
		"product-runtime/studio-catalog",
		"-run", "TestChallenge_OptionalDefaults_StrictValidation",
	)
	if !passed || err != nil {
		t.Fatalf("T2.F1.5 FAIL: Optional default values validation failed: %v\nOutput: %s", err, output)
	}

	validationSubtests := []string{
		"Invalid_Integer_Default_Fails_Validation",
		"Disallowed_Enum_Member_Default_Fails_Validation",
		"Literal_Secret_in_Optional_Default_Fails_Validation",
		"Valid_Optional_Default_Passes_Validation",
	}
	for _, sub := range validationSubtests {
		if !strings.Contains(output, sub) {
			t.Fatalf("T2.F1.5 FAIL: expected optional default subtest %q not found in output", sub)
		}
	}
	t.Log("T2.F1.5 PASS: optional default values fallback adheres strictly to schema validation without crashes")
}
