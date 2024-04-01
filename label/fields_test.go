package label

import (
	"reflect"
	"testing"
)

func TestParseIndex(t *testing.T) {
	tests := []struct {
		name     string
		raw      string
		expected Index
	}{
		{
			name:     "Single value in raw string",
			raw:      "value",
			expected: Index{"value"},
		},
		{
			name:     "Single escape value in raw string",
			raw:      "\"value1.value2\"",
			expected: Index{"value1.value2"},
		},
		{
			name:     "Multiple values in raw string",
			raw:      "value1.value2.value3",
			expected: Index{"value1", "value2", "value3"},
		},
		{
			name:     "Multiple values with escape characters center",
			raw:      "value1.\"value2.value3\".value4",
			expected: Index{"value1", "value2.value3", "value4"},
		},
		{
			name:     "Multiple values with escape characters left",
			raw:      "\"value1.value2\".value3.value4",
			expected: Index{"value1.value2", "value3", "value4"},
		},
		{
			name:     "Multiple values with escape characters right",
			raw:      "value1.value2.\"value3.value4\"",
			expected: Index{"value1", "value2", "value3.value4"},
		},
		{
			name:     "Multiple values with muliple escape characters",
			raw:      "\"value1.value2\".value3.\"value4.value5\".value6.\"value7.value8\"",
			expected: Index{"value1.value2", "value3", "value4.value5", "value6", "value7.value8"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ParseIndex(test.raw)
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("Unexpected result. Got: %s, Want: %s", result, test.expected)
			}
		})
	}
}

func TestFindMap(t *testing.T) {
	tests := []struct {
		name     string
		index    Index
		input    map[string]interface{}
		expected interface{}
	}{
		{
			name:     "Empty index and nil input",
			index:    Index{},
			input:    nil,
			expected: nil,
		},
		{
			name:     "Empty index and non-nil input",
			index:    Index{},
			input:    map[string]interface{}{"key": "value"},
			expected: nil,
		},
		{
			name:     "Non-empty index and nil input",
			index:    Index{"key"},
			input:    nil,
			expected: nil,
		},
		{
			name:     "Non-empty index and non-nil input with missing key",
			index:    Index{"key"},
			input:    map[string]interface{}{"otherKey": "value"},
			expected: nil,
		},
		{
			name:     "Non-empty index and non-nil input with existing key",
			index:    Index{"key"},
			input:    map[string]interface{}{"key": "value"},
			expected: "value",
		},
		{
			name:     "Nested map with non-empty index",
			index:    Index{"key1", "key2"},
			input:    map[string]interface{}{"key1": map[string]interface{}{"key2": "value"}},
			expected: "value",
		},
		{
			name:     "Nested map with empty index",
			index:    Index{},
			input:    map[string]interface{}{"key1": map[string]interface{}{"key2": "value"}},
			expected: nil,
		},
		{
			name:     "Nested map with invalid index",
			index:    Index{"key1", "key2", "key3"},
			input:    map[string]interface{}{"key1": map[string]interface{}{"key2": "value"}},
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := FindMap(test.index, test.input)
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("Unexpected result. Got: %v, Want: %v", result, test.expected)
			}
		})
	}
}

func TestFieldExtractor_Extract(t *testing.T) {
	tests := []struct {
		name     string
		mappings map[string]string
		input    map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name:     "Empty mappings and nil input",
			mappings: map[string]string{},
			input:    nil,
			expected: map[string]interface{}{},
		},
		{
			name:     "Empty mappings and non-nil input",
			mappings: map[string]string{},
			input:    map[string]interface{}{"key": "value"},
			expected: map[string]interface{}{},
		},
		{
			name:     "Non-empty mappings and nil input",
			mappings: map[string]string{"key": "key1"},
			input:    nil,
			expected: map[string]interface{}{"key": nil},
		},
		{
			name:     "Non-empty mappings and non-nil input with missing key",
			mappings: map[string]string{"key": "key1"},
			input:    map[string]interface{}{"otherKey": "value"},
			expected: map[string]interface{}{"key": nil},
		},
		{
			name:     "Non-empty mappings and non-nil input with existing key",
			mappings: map[string]string{"key": "key1"},
			input:    map[string]interface{}{"key1": "value"},
			expected: map[string]interface{}{"key": "value"},
		},
		{
			name:     "Nested mappings with non-empty index",
			mappings: map[string]string{"key": "key1.key2"},
			input:    map[string]interface{}{"key1": map[string]interface{}{"key2": "value"}},
			expected: map[string]interface{}{"key": "value"},
		},
		{
			name:     "Nested mappings with empty index",
			mappings: map[string]string{},
			input:    map[string]interface{}{"key1": map[string]interface{}{"key2": "value"}},
			expected: map[string]interface{}{},
		},
		{
			name:     "Nested mappings with invalid index",
			mappings: map[string]string{"key": "key1.key2.key3"},
			input:    map[string]interface{}{"key1": map[string]interface{}{"key2": "value"}},
			expected: map[string]interface{}{"key": nil},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := NewFieldExtractor(test.mappings)
			result := f.Extract(test.input)
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("Unexpected result. Got: %v, Want: %v", result, test.expected)
			}
		})
	}
}
