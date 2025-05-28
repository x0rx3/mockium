package comparer

import (
	"mockium/internal/service/constants"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComparer_compare(t *testing.T) {
	comparer := New()

	tests := []struct {
		name     string
		expected any
		actual   any
		want     bool
	}{
		{
			name:     "Simple string match",
			expected: "test",
			actual:   "test",
			want:     true,
		},
		{
			name:     "Simple string mismatch",
			expected: "test",
			actual:   "mismatch",
			want:     false,
		},
		{
			name:     "Regexp match",
			expected: regexp.MustCompile("^[a-z]+$"),
			actual:   "test",
			want:     true,
		},
		{
			name:     "Regexp mismatch",
			expected: regexp.MustCompile("^[0-9]+$"),
			actual:   "test",
			want:     false,
		},
		{
			name:     "Any value placeholder",
			expected: constants.AnyValuePlaceholder,
			actual:   "anything",
			want:     true,
		},
		{
			name:     "Nested map match",
			expected: map[string]any{"user": map[string]any{"name": "John"}},
			actual:   map[string]any{"user": map[string]any{"name": "John"}},
			want:     true,
		},
		{
			name:     "Nested map mismatch",
			expected: map[string]any{"user": map[string]any{"name": "John"}},
			actual:   map[string]any{"user": map[string]any{"name": "Alice"}},
			want:     false,
		},
		{
			name:     "Slice match",
			expected: []any{1, 2, 3},
			actual:   []any{1, 2, 3},
			want:     true,
		},
		{
			name:     "Slice length mismatch",
			expected: []any{1, 2},
			actual:   []any{1, 2, 3},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, comparer.Compare(tt.expected, tt.actual))
		})
	}
}
