package ui

import (
	"strings"
	"testing"
)

func TestDetectContentType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected contentType
	}{
		// JSON
		{"json object", `{"key": "value"}`, contentJSON},
		{"json array", `[1, 2, 3]`, contentJSON},
		{"json with whitespace", `  {"key": "value"}  `, contentJSON},
		{"invalid json object", `{key: value}`, contentUnknown},

		// HTML
		{"html doctype", `<!DOCTYPE html><html></html>`, contentHTML},
		{"html doctype lowercase", `<!doctype html><html></html>`, contentHTML},
		{"html tag", `<html><body></body></html>`, contentHTML},
		{"html with head", `<head><title>Test</title></head>`, contentHTML},
		{"html with body", `<body><p>Hello</p></body>`, contentHTML},
		{"html with div", `<div class="test">content</div>`, contentHTML},
		{"html with script", `<script>alert('hi')</script>`, contentHTML},

		// XML
		{"xml declaration", `<?xml version="1.0"?><root/>`, contentXML},
		{"xml with xmlns", `<root xmlns="http://example.com"><item/></root>`, contentXML},
		{"soap envelope", `<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"></soap:Envelope>`, contentXML},
		{"generic xml", `<root><child>value</child></root>`, contentXML},

		// Unknown
		{"empty", ``, contentUnknown},
		{"plain text", `Hello, World!`, contentUnknown},
		{"number", `12345`, contentUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detectContentType(tt.input)
			if result != tt.expected {
				t.Errorf("detectContentType(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRenderResponseJSON(t *testing.T) {
	input := `{"name":"test","value":123}`
	result := renderResponse(input)

	// Should be pretty-printed (contains newlines and indentation)
	if !strings.Contains(result, "\n") {
		t.Error("JSON should be pretty-printed with newlines")
	}
	if !strings.Contains(result, "  ") {
		t.Error("JSON should be pretty-printed with indentation")
	}
}

func TestRenderResponseHTML(t *testing.T) {
	input := `<!DOCTYPE html><html><body><h1>Hello</h1></body></html>`
	result := renderResponse(input)

	// Should contain the original content (highlighting adds ANSI codes)
	if !strings.Contains(result, "html") {
		t.Error("HTML response should contain original content")
	}
}

func TestRenderResponseXML(t *testing.T) {
	input := `<?xml version="1.0"?><root><item>test</item></root>`
	result := renderResponse(input)

	// Should contain the original content
	if !strings.Contains(result, "root") {
		t.Error("XML response should contain original content")
	}
}

func TestRenderResponsePlainText(t *testing.T) {
	input := `Just some plain text`
	result := renderResponse(input)

	// Should return unchanged
	if result != input {
		t.Errorf("Plain text should be unchanged, got %q", result)
	}
}

func TestPrettyJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:  "simple object",
			input: `{"a":"b"}`,
			expected: `{
  "a": "b"
}`,
		},
		{
			name:  "nested object",
			input: `{"outer":{"inner":"value"}}`,
			expected: `{
  "outer": {
    "inner": "value"
  }
}`,
		},
		{
			name:     "invalid json",
			input:    `not json`,
			expected: `not json`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := prettyJSON(tt.input)
			if result != tt.expected {
				t.Errorf("prettyJSON(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
