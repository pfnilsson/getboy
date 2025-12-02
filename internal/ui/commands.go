package ui

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func doHTTP(method, url, body string, headers map[string]string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
		defer cancel()

		// Expand env vars in body
		expandedBody := expandEnvVars(body)

		var reader io.Reader
		if expandedBody != "" {
			reader = bytes.NewBufferString(expandedBody)
		}
		req, err := http.NewRequestWithContext(ctx, method, url, reader)
		if err != nil {
			return httpDoneMsg{Err: err}
		}

		// Set headers with env var expansion
		for k, v := range headers {
			req.Header.Set(k, expandEnvVars(v))
		}

		// Default Content-Type for body if not already set
		if expandedBody != "" && req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", "application/json")
		}

		client := &http.Client{Timeout: 12 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return httpDoneMsg{Err: err}
		}
		defer func() { _ = resp.Body.Close() }()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return httpDoneMsg{Err: err}
		}
		return httpDoneMsg{Status: resp.Status, Body: string(b)}
	}
}
