package ui

import (
	"os"
	"testing"
)

func TestExpandEnvVars(t *testing.T) {
	// Set up test environment variables
	os.Setenv("TEST_VAR", "test_value")
	os.Setenv("API_KEY", "secret123")
	os.Setenv("EMPTY_VAR", "")
	defer func() {
		os.Unsetenv("TEST_VAR")
		os.Unsetenv("API_KEY")
		os.Unsetenv("EMPTY_VAR")
	}()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no variables",
			input:    "plain text",
			expected: "plain text",
		},
		{
			name:     "single variable",
			input:    "Bearer ${API_KEY}",
			expected: "Bearer secret123",
		},
		{
			name:     "multiple variables",
			input:    "${TEST_VAR} and ${API_KEY}",
			expected: "test_value and secret123",
		},
		{
			name:     "variable at start",
			input:    "${TEST_VAR} suffix",
			expected: "test_value suffix",
		},
		{
			name:     "variable at end",
			input:    "prefix ${TEST_VAR}",
			expected: "prefix test_value",
		},
		{
			name:     "undefined variable becomes empty",
			input:    "prefix ${UNDEFINED_VAR} suffix",
			expected: "prefix  suffix",
		},
		{
			name:     "empty variable value",
			input:    "prefix ${EMPTY_VAR} suffix",
			expected: "prefix  suffix",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only variable",
			input:    "${API_KEY}",
			expected: "secret123",
		},
		{
			name:     "adjacent variables",
			input:    "${TEST_VAR}${API_KEY}",
			expected: "test_valuesecret123",
		},
		{
			name:     "variable with underscores",
			input:    "${TEST_VAR}",
			expected: "test_value",
		},
		{
			name:     "incomplete syntax not replaced",
			input:    "$API_KEY",
			expected: "$API_KEY",
		},
		{
			name:     "incomplete syntax with brace",
			input:    "${API_KEY",
			expected: "${API_KEY",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandEnvVars(tt.input)
			if result != tt.expected {
				t.Errorf("expandEnvVars(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestEnvVarPattern(t *testing.T) {
	// Test the regex pattern matches correctly
	tests := []struct {
		input   string
		matches []string
	}{
		{"${VAR}", []string{"${VAR}"}},
		{"${VAR_NAME}", []string{"${VAR_NAME}"}},
		{"${a}${b}", []string{"${a}", "${b}"}},
		{"no vars", nil},
		{"$VAR", nil},
		{"${}", nil}, // empty var name doesn't match (requires at least one char)
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			matches := envVarPattern.FindAllString(tt.input, -1)
			if len(matches) != len(tt.matches) {
				t.Errorf("FindAllString(%q) found %d matches, want %d", tt.input, len(matches), len(tt.matches))
				return
			}
			for i, m := range matches {
				if m != tt.matches[i] {
					t.Errorf("match[%d] = %q, want %q", i, m, tt.matches[i])
				}
			}
		})
	}
}
