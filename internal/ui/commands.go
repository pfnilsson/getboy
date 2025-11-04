package ui

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func doHTTP(method, url, body string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
		defer cancel()

		var reader io.Reader
		if body != "" {
			reader = bytes.NewBufferString(body)
		}
		req, err := http.NewRequestWithContext(ctx, method, url, reader)
		if err != nil {
			return httpDoneMsg{Err: err}
		}
		if body != "" {
			req.Header.Set("Content-Type", "application/json")
		}

		client := &http.Client{Timeout: 12 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return httpDoneMsg{Err: err}
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return httpDoneMsg{Err: err}
		}
		return httpDoneMsg{Status: resp.Status, Body: string(b)}
	}
}
