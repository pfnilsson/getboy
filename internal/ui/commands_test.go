package ui

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// TestDoHTTP tests the HTTP command with a mock server
func TestDoHTTP(t *testing.T) {
	t.Run("successful GET request", func(t *testing.T) {
		// Create a test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				t.Errorf("expected GET request, got %s", r.Method)
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"message":"success"}`))
		}))
		defer server.Close()

		// Execute the command
		cmd := doHTTP("GET", server.URL, "", nil)
		msg := cmd()

		// Check the result
		result, ok := msg.(httpDoneMsg)
		if !ok {
			t.Fatalf("expected httpDoneMsg, got %T", msg)
		}

		if result.Err != nil {
			t.Errorf("unexpected error: %v", result.Err)
		}

		if result.Status != "200 OK" {
			t.Errorf("status = %q, want %q", result.Status, "200 OK")
		}

		if result.Body != `{"message":"success"}` {
			t.Errorf("body = %q, want %q", result.Body, `{"message":"success"}`)
		}
	})

	t.Run("successful POST request with body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Errorf("expected POST request, got %s", r.Method)
			}

			// Check Content-Type header
			ct := r.Header.Get("Content-Type")
			if ct != "application/json" {
				t.Errorf("Content-Type = %q, want %q", ct, "application/json")
			}

			// Echo back the request body
			var data map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
				t.Errorf("failed to decode body: %v", err)
			}

			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(data)
		}))
		defer server.Close()

		// Execute the command
		requestBody := `{"test":"data"}`
		cmd := doHTTP("POST", server.URL, requestBody, nil)
		msg := cmd()

		// Check the result
		result, ok := msg.(httpDoneMsg)
		if !ok {
			t.Fatalf("expected httpDoneMsg, got %T", msg)
		}

		if result.Err != nil {
			t.Errorf("unexpected error: %v", result.Err)
		}

		if result.Status != "200 OK" {
			t.Errorf("status = %q, want %q", result.Status, "200 OK")
		}
	})

	t.Run("handles invalid URL", func(t *testing.T) {
		cmd := doHTTP("GET", "://invalid-url", "", nil)
		msg := cmd()

		result, ok := msg.(httpDoneMsg)
		if !ok {
			t.Fatalf("expected httpDoneMsg, got %T", msg)
		}

		if result.Err == nil {
			t.Error("expected error for invalid URL, got nil")
		}
	})

	t.Run("handles 404 response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("not found"))
		}))
		defer server.Close()

		cmd := doHTTP("GET", server.URL, "", nil)
		msg := cmd()

		result, ok := msg.(httpDoneMsg)
		if !ok {
			t.Fatalf("expected httpDoneMsg, got %T", msg)
		}

		if result.Err != nil {
			t.Errorf("unexpected error: %v", result.Err)
		}

		if result.Status != "404 Not Found" {
			t.Errorf("status = %q, want %q", result.Status, "404 Not Found")
		}

		if result.Body != "not found" {
			t.Errorf("body = %q, want %q", result.Body, "not found")
		}
	})

	t.Run("sends custom headers", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth != "Bearer mytoken" {
				t.Errorf("Authorization = %q, want %q", auth, "Bearer mytoken")
			}
			custom := r.Header.Get("X-Custom-Header")
			if custom != "custom-value" {
				t.Errorf("X-Custom-Header = %q, want %q", custom, "custom-value")
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		headers := map[string]string{
			"Authorization":   "Bearer mytoken",
			"X-Custom-Header": "custom-value",
		}
		cmd := doHTTP("GET", server.URL, "", headers)
		msg := cmd()

		result, ok := msg.(httpDoneMsg)
		if !ok {
			t.Fatalf("expected httpDoneMsg, got %T", msg)
		}
		if result.Err != nil {
			t.Errorf("unexpected error: %v", result.Err)
		}
	})

	t.Run("expands env vars in headers", func(t *testing.T) {
		os.Setenv("TEST_TOKEN", "env_token_value")
		defer os.Unsetenv("TEST_TOKEN")

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth != "Bearer env_token_value" {
				t.Errorf("Authorization = %q, want %q", auth, "Bearer env_token_value")
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		headers := map[string]string{
			"Authorization": "Bearer ${TEST_TOKEN}",
		}
		cmd := doHTTP("GET", server.URL, "", headers)
		msg := cmd()

		result, ok := msg.(httpDoneMsg)
		if !ok {
			t.Fatalf("expected httpDoneMsg, got %T", msg)
		}
		if result.Err != nil {
			t.Errorf("unexpected error: %v", result.Err)
		}
	})

	t.Run("expands env vars in body", func(t *testing.T) {
		os.Setenv("TEST_VALUE", "expanded_value")
		defer os.Unsetenv("TEST_VALUE")

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var data map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
				t.Errorf("failed to decode body: %v", err)
			}
			if data["key"] != "expanded_value" {
				t.Errorf("body key = %q, want %q", data["key"], "expanded_value")
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		body := `{"key":"${TEST_VALUE}"}`
		cmd := doHTTP("POST", server.URL, body, nil)
		msg := cmd()

		result, ok := msg.(httpDoneMsg)
		if !ok {
			t.Fatalf("expected httpDoneMsg, got %T", msg)
		}
		if result.Err != nil {
			t.Errorf("unexpected error: %v", result.Err)
		}
	})

	t.Run("custom Content-Type header overrides default", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ct := r.Header.Get("Content-Type")
			if ct != "text/plain" {
				t.Errorf("Content-Type = %q, want %q", ct, "text/plain")
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		headers := map[string]string{
			"Content-Type": "text/plain",
		}
		cmd := doHTTP("POST", server.URL, "plain text body", headers)
		msg := cmd()

		result, ok := msg.(httpDoneMsg)
		if !ok {
			t.Fatalf("expected httpDoneMsg, got %T", msg)
		}
		if result.Err != nil {
			t.Errorf("unexpected error: %v", result.Err)
		}
	})
}
