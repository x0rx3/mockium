package comparer

import (
	"mockium/internal/defs"
	"mockium/pkg/core/model"
	"net/http"
	"reflect"
	"regexp"
)

func New() *comparer {
	return &comparer{}
}

type comparer struct{}

func (inst *comparer) Compare(expected, actual any) bool {
	switch exp := expected.(type) {
	case *regexp.Regexp:
		if str, ok := actual.(string); ok {
			return exp.MatchString(str)
		}
		return false
	case string:
		if exp == defs.AnyValuePlaceholder {
			return true
		}
		return exp == actual
	case []any:
		if aSlice, ok := actual.([]any); ok {
			return inst.compareSlices(exp, aSlice)
		}
		return false
	case map[string]any:
		if aMap, ok := actual.(map[string]any); ok {
			return inst.compareMaps(exp, aMap)
		}
		return false
	case int, int8, int16, int32, int64:
		expVal := reflect.ValueOf(exp).Int()
		switch actVal := actual.(type) {
		case int, int8, int16, int32, int64:
			return expVal == reflect.ValueOf(actVal).Int()
		case float32, float64:
			return float64(expVal) == reflect.ValueOf(actVal).Float()
		}
	case float32, float64:
		expVal := reflect.ValueOf(exp).Float()
		switch actVal := actual.(type) {
		case int, int8, int16, int32, int64:
			return expVal == float64(reflect.ValueOf(actVal).Int())
		case float32, float64:
			return expVal == reflect.ValueOf(actVal).Float()
		}
	case []*model.Cookie:
		if _, ok := actual.([]*http.Cookie); !ok {
			return false
		}
	}
	return expected == actual
}

// compareSlices performs a deep comparison of two slices using the Comparer's rules.
// The slices are considered equal if:
//   - They have the same length
//   - Each corresponding element matches according to Compare()
//
// Parameters:
//   - expected: The slice containing expected values/patterns
//   - actual: The slice to compare against
//
// Returns:
//   - true if all elements match, false otherwise
func (inst *comparer) compareSlices(expected, actual []any) bool {
	if len(expected) != len(actual) {
		return false
	}

	for i := range expected {
		if !inst.Compare(expected[i], actual[i]) {
			return false
		}
	}
	return true
}

// compareMaps performs a deep comparison of two maps using the Comparer's rules.
// The maps are considered equal if:
//   - All keys from expected exist in actual
//   - Each corresponding value matches according to Compare()
//
// Note: This is not a symmetric comparison - extra keys in actual are ignored.
//
// Parameters:
//   - expected: The map containing expected keys/patterns
//   - actual: The map to compare against
//
// Returns:
//   - true if all expected keys and values match, false otherwise
func (inst *comparer) compareMaps(expected, actual map[string]any) bool {
	for key, expVal := range expected {
		if actVal, exists := actual[key]; !exists || !inst.Compare(expVal, actVal) {
			return false
		}
	}
	return true
}
