package ui

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/chroma/v2/quick"
)

func isLikelyJSON(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	if s[0] != '{' && s[0] != '[' {
		return false
	}
	var v any
	return json.Unmarshal([]byte(s), &v) == nil
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

func highlightJSON(s string) string {
	var buf bytes.Buffer
	if err := quick.Highlight(&buf, s, "json", bestFormatter(), chromaStyle()); err != nil {
		return s
	}
	return buf.String()
}

// renderResponse pretty-prints + highlights JSON when detected.
func renderResponse(body string) string {
	if isLikelyJSON(body) {
		return highlightJSON(prettyJSON(body))
	}
	return body
}
