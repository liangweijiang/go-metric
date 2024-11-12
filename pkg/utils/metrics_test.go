package utils

import (
	"fmt"
	"testing"
)

func TestSanitizeMetricName(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"/cpu/classes/gc/mark/assist:cpu-seconds", ""},
		/*{"/start", "_start"},
		{"noChange", "noChange"},
		{"1start", "o_1start"},
		{"!special$", "o__special_"},
		{"with/slash_and_123", "with_slash_and_123"},*/
	}

	for _, tc := range testCases {
		result := SanitizeMetricName(tc.input)
		/*if result != tc.expected {
			t.Errorf("SanitizeMetricName(%q) = %q; want %q", tc.input, result, tc.expected)
		}*/
		fmt.Println(result)
	}
}
