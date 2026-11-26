package harness

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

// StrictUnmarshal unmarshals JSON while disallowing any unknown fields.
func StrictUnmarshal(data []byte, v interface{}) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}

// GetJSONField traverses a dot-separated path (e.g. "dimensions.code.passed" or "items.0.name")
// in a JSON object or string and returns the target value.
func GetJSONField(jsonStr string, dotPath string) (interface{}, error) {
	var current interface{}
	decoder := json.NewDecoder(strings.NewReader(jsonStr))
	if err := decoder.Decode(&current); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	parts := strings.Split(dotPath, ".")
	for _, part := range parts {
		if current == nil {
			return nil, fmt.Errorf("nil value at path segment '%s'", part)
		}
		switch val := current.(type) {
		case map[string]interface{}:
			next, ok := val[part]
			if !ok {
				return nil, fmt.Errorf("field '%s' not found in JSON object", part)
			}
			current = next
		case []interface{}:
			idx, err := strconv.Atoi(part)
			if err != nil || idx < 0 || idx >= len(val) {
				return nil, fmt.Errorf("invalid array index '%s' in JSON array", part)
			}
			current = val[idx]
		default:
			return nil, fmt.Errorf("cannot traverse segment '%s' on non-container type %T", part, current)
		}
	}
	return current, nil
}

// AssertJSONFieldEqual verifies that the JSON field at dotPath equals expected.
func AssertJSONFieldEqual(t *testing.T, jsonStr, dotPath string, expected interface{}) {
	t.Helper()
	actual, err := GetJSONField(jsonStr, dotPath)
	if err != nil {
		t.Fatalf("path '%s' resolution failed: %v", dotPath, err)
	}

	// Normalize numbers for float64 vs int comparison
	if expNum, ok := toFloat64(expected); ok {
		if actNum, ok := toFloat64(actual); ok {
			if expNum != actNum {
				t.Fatalf("JSON path '%s' mismatch: expected %v (%T), got %v (%T)", dotPath, expected, expected, actual, actual)
			}
			return
		}
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("JSON path '%s' mismatch: expected %v, got %v", dotPath, expected, actual)
	}
}

// AssertJSONFieldExists verifies that a dotPath exists in the JSON document.
func AssertJSONFieldExists(t *testing.T, jsonStr, dotPath string) {
	t.Helper()
	_, err := GetJSONField(jsonStr, dotPath)
	if err != nil {
		t.Fatalf("expected JSON field '%s' to exist, but got error: %v", dotPath, err)
	}
}

func toFloat64(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case int:
		return float64(n), true
	case int8:
		return float64(n), true
	case int16:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case uint:
		return float64(n), true
	case uint8:
		return float64(n), true
	case uint16:
		return float64(n), true
	case uint32:
		return float64(n), true
	case uint64:
		return float64(n), true
	case float32:
		return float64(n), true
	case float64:
		return n, true
	default:
		return 0, false
	}
}
