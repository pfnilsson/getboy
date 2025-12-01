package ui

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/chroma/v2/quick"
)

// contentType represents detected content types for highlighting
type contentType int

const (
	contentUnknown contentType = iota
	contentJSON
	contentHTML
	contentXML
)

func detectContentType(s string) contentType {
	s = strings.TrimSpace(s)
	if s == "" {
		return contentUnknown
	}

	// Check for JSON (starts with { or [)
	if s[0] == '{' || s[0] == '[' {
		var v any
		if json.Unmarshal([]byte(s), &v) == nil {
			return contentJSON
		}
	}

	// Check for HTML
	lower := strings.ToLower(s)
	if strings.HasPrefix(lower, "<!doctype html") ||
		strings.HasPrefix(lower, "<html") ||
		strings.HasPrefix(s, "<!DOCTYPE html") {
		return contentHTML
	}

	// Check for XML (including XHTML)
	if strings.HasPrefix(s, "<?xml") ||
		strings.HasPrefix(lower, "<!doctype") ||
		(strings.HasPrefix(s, "<") && strings.Contains(s[:min(100, len(s))], "xmlns")) {
		return contentXML
	}

	// Simple heuristic: if it starts with < and contains common HTML tags
	if strings.HasPrefix(s, "<") {
		testStr := strings.ToLower(s[:min(500, len(s))])
		htmlTags := []string{"<head", "<body", "<div", "<span", "<p>", "<a ", "<script", "<style", "<meta", "<link"}
		for _, tag := range htmlTags {
			if strings.Contains(testStr, tag) {
				return contentHTML
			}
		}
		// If it looks like markup but not HTML, assume XML
		if strings.Contains(s, "</") {
			return contentXML
		}
	}

	return contentUnknown
}

func prettyJSON(s string) string {
	var v any
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return s
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	_ = enc.Encode(v)
	return strings.TrimRight(buf.String(), "\n")
}

func highlight(s string, lexer string) string {
	var buf bytes.Buffer
	if err := quick.Highlight(&buf, s, lexer, bestFormatter(), chromaStyle()); err != nil {
		return s
	}
	return buf.String()
}

// renderResponse detects content type and applies syntax highlighting.
func renderResponse(body string) string {
	switch detectContentType(body) {
	case contentJSON:
		return highlight(prettyJSON(body), "json")
	case contentHTML:
		return highlight(body, "html")
	case contentXML:
		return highlight(body, "xml")
	default:
		return body
	}
}
