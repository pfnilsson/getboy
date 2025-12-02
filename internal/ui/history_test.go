package ui

import (
	"os"
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestAddToHistory tests adding entries to history
func TestAddToHistory(t *testing.T) {
	t.Run("adds new entry", func(t *testing.T) {
		var entries []historyEntry
		entry := historyEntry{Method: "GET", URL: "https://example.com"}
		entries = addToHistory(entries, entry)

		if len(entries) != 1 {
			t.Fatalf("expected 1 entry, got %d", len(entries))
		}
		if entries[0].URL != "https://example.com" {
			t.Errorf("URL = %q, want %q", entries[0].URL, "https://example.com")
		}
	})

	t.Run("prepends new entries", func(t *testing.T) {
		entries := []historyEntry{{Method: "GET", URL: "https://first.com"}}
		entry := historyEntry{Method: "POST", URL: "https://second.com"}
		entries = addToHistory(entries, entry)

		if len(entries) != 2 {
			t.Fatalf("expected 2 entries, got %d", len(entries))
		}
		if entries[0].URL != "https://second.com" {
			t.Errorf("first URL = %q, want %q", entries[0].URL, "https://second.com")
		}
	})

	t.Run("removes duplicates and moves to front", func(t *testing.T) {
		entries := []historyEntry{
			{Method: "GET", URL: "https://first.com"},
			{Method: "GET", URL: "https://second.com"},
			{Method: "GET", URL: "https://third.com"},
		}
		// Re-add the second URL with same params
		entry := historyEntry{Method: "GET", URL: "https://second.com"}
		entries = addToHistory(entries, entry)

		if len(entries) != 3 {
			t.Fatalf("expected 3 entries, got %d", len(entries))
		}
		// https://second.com should now be first
		if entries[0].URL != "https://second.com" {
			t.Errorf("first URL = %q, want %q", entries[0].URL, "https://second.com")
		}
		if entries[1].URL != "https://first.com" {
			t.Errorf("second URL = %q, want %q", entries[1].URL, "https://first.com")
		}
	})

	t.Run("same URL with different body creates separate entries", func(t *testing.T) {
		entries := []historyEntry{
			{Method: "POST", URL: "https://api.com", Body: `{"v":1}`},
		}
		// Add same URL but different body
		entry := historyEntry{Method: "POST", URL: "https://api.com", Body: `{"v":2}`}
		entries = addToHistory(entries, entry)

		if len(entries) != 2 {
			t.Fatalf("expected 2 entries (different body), got %d", len(entries))
		}
		if entries[0].Body != `{"v":2}` {
			t.Errorf("first body = %q, want %q", entries[0].Body, `{"v":2}`)
		}
	})

	t.Run("same URL with different headers creates separate entries", func(t *testing.T) {
		entries := []historyEntry{
			{Method: "GET", URL: "https://api.com", Headers: map[string]string{"Auth": "token1"}},
		}
		// Add same URL but different headers
		entry := historyEntry{Method: "GET", URL: "https://api.com", Headers: map[string]string{"Auth": "token2"}}
		entries = addToHistory(entries, entry)

		if len(entries) != 2 {
			t.Fatalf("expected 2 entries (different headers), got %d", len(entries))
		}
		if entries[0].Headers["Auth"] != "token2" {
			t.Errorf("first auth = %q, want %q", entries[0].Headers["Auth"], "token2")
		}
	})

	t.Run("identical request is deduplicated", func(t *testing.T) {
		entries := []historyEntry{
			{Method: "POST", URL: "https://api.com", Body: `{"data":true}`, Headers: map[string]string{"Content-Type": "application/json"}},
		}
		// Add exact same request
		entry := historyEntry{Method: "POST", URL: "https://api.com", Body: `{"data":true}`, Headers: map[string]string{"Content-Type": "application/json"}}
		entries = addToHistory(entries, entry)

		if len(entries) != 1 {
			t.Fatalf("expected 1 entry (identical request), got %d", len(entries))
		}
	})

	t.Run("trims to max size", func(t *testing.T) {
		var entries []historyEntry
		// Add more than maxHistoryItems with unique URLs
		for i := range maxHistoryItems + 10 {
			entry := historyEntry{Method: "GET", URL: "https://example.com/page/" + string(rune(i))}
			entries = addToHistory(entries, entry)
		}

		if len(entries) != maxHistoryItems {
			t.Errorf("expected %d entries, got %d", maxHistoryItems, len(entries))
		}
	})
}

// TestHistoryToItems tests conversion of history entries to list items
func TestHistoryToItems(t *testing.T) {
	entries := []historyEntry{
		{Method: "GET", URL: "https://example.com/api", Body: "", Headers: nil},
		{Method: "POST", URL: "https://example.com/submit", Body: `{"key":"value"}`, Headers: map[string]string{"Content-Type": "application/json"}},
	}

	items := historyToItems(entries)

	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}

	// Check first item
	if items[0].method != "GET" {
		t.Errorf("first method = %q, want %q", items[0].method, "GET")
	}
	if items[0].url != "https://example.com/api" {
		t.Errorf("first url = %q, want %q", items[0].url, "https://example.com/api")
	}

	// Check second item
	if items[1].method != "POST" {
		t.Errorf("second method = %q, want %q", items[1].method, "POST")
	}
	if items[1].body != `{"key":"value"}` {
		t.Errorf("second body = %q, want %q", items[1].body, `{"key":"value"}`)
	}
	if items[1].headers["Content-Type"] != "application/json" {
		t.Errorf("second Content-Type = %q, want %q", items[1].headers["Content-Type"], "application/json")
	}
}

// TestTruncateURL tests URL truncation for display
func TestTruncateURL(t *testing.T) {
	tests := []struct {
		url    string
		maxLen int
		want   string
	}{
		{"short", 10, "short"},
		{"exactly10!", 10, "exactly10!"},
		{"this is a very long url", 15, "this is a ve..."},
		{"https://example.com/api/endpoint", 20, "https://example.c..."},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			got := truncateURL(tt.url, tt.maxLen)
			if got != tt.want {
				t.Errorf("truncateURL(%q, %d) = %q, want %q", tt.url, tt.maxLen, got, tt.want)
			}
		})
	}
}

// TestSidebarTabNavigation tests tab/shift-tab switches sidebar tabs
func TestSidebarTabNavigation(t *testing.T) {
	m := New().(model)
	m.pane = paneSidebar

	// Start on History tab
	if m.sidebarTab != sidebarHistory {
		t.Fatalf("initial sidebarTab = %v, want sidebarHistory", m.sidebarTab)
	}

	// Press tab - should go to Saved
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(model)
	if m.sidebarTab != sidebarSaved {
		t.Errorf("after tab sidebarTab = %v, want sidebarSaved", m.sidebarTab)
	}

	// Press tab again - should cycle back to History
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(model)
	if m.sidebarTab != sidebarHistory {
		t.Errorf("after second tab sidebarTab = %v, want sidebarHistory", m.sidebarTab)
	}

	// Press shift+tab - should go to Saved
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	m = updated.(model)
	if m.sidebarTab != sidebarSaved {
		t.Errorf("after shift+tab sidebarTab = %v, want sidebarSaved", m.sidebarTab)
	}
}

// TestSidebarRightLeftNavigation tests right/left arrows switch sidebar tabs
func TestSidebarRightLeftNavigation(t *testing.T) {
	m := New().(model)
	m.pane = paneSidebar

	// Press right - should go to Saved
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRight})
	m = updated.(model)
	if m.sidebarTab != sidebarSaved {
		t.Errorf("after right sidebarTab = %v, want sidebarSaved", m.sidebarTab)
	}

	// Press left - should go back to History
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyLeft})
	m = updated.(model)
	if m.sidebarTab != sidebarHistory {
		t.Errorf("after left sidebarTab = %v, want sidebarHistory", m.sidebarTab)
	}
}

// TestHistoryRecordedOnSend tests that history is recorded when sending a request
func TestHistoryRecordedOnSend(t *testing.T) {
	m := New().(model)
	// Clear any existing history from previous test runs
	m.history = nil
	m.pane = paneEditor
	m.url.SetValue("https://httpbin.org/get")
	m.setMethod("GET")

	// Press enter to send request (will also add to history)
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(model)

	// History should have exactly one entry
	if len(m.history) != 1 {
		t.Errorf("history length = %d, want 1", len(m.history))
	}

	// Check the entry was recorded correctly
	if m.history[0].Method != "GET" {
		t.Errorf("history method = %q, want %q", m.history[0].Method, "GET")
	}
	if m.history[0].URL != "https://httpbin.org/get" {
		t.Errorf("history URL = %q, want %q", m.history[0].URL, "https://httpbin.org/get")
	}
}

// TestHistoryWithHeaders tests that headers are recorded in history
func TestHistoryWithHeaders(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.url.SetValue("https://httpbin.org/post")
	m.setMethod("POST")
	m.headers[0].key.SetValue("Content-Type")
	m.headers[0].value.SetValue("application/json")
	m.body.SetValue(`{"test":true}`)

	// Send request
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(model)

	// Check headers were recorded
	if m.history[0].Headers["Content-Type"] != "application/json" {
		t.Errorf("history Content-Type = %q, want %q", m.history[0].Headers["Content-Type"], "application/json")
	}
	if m.history[0].Body != `{"test":true}` {
		t.Errorf("history body = %q, want %q", m.history[0].Body, `{"test":true}`)
	}
}

// TestLoadFromHistoryWithHeaders tests that headers are loaded from history
func TestLoadFromHistoryWithHeaders(t *testing.T) {
	m := New().(model)
	m.pane = paneSidebar

	// Add a history item with headers
	m.addToHistoryAndSave("POST", "https://api.example.com", `{"data":"test"}`, map[string]string{
		"Authorization": "Bearer token",
		"Content-Type":  "application/json",
	})

	// Press enter to load
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(model)

	// Check that headers were loaded
	foundAuth := false
	foundCT := false
	for _, h := range m.headers {
		if h.key.Value() == "Authorization" && h.value.Value() == "Bearer token" {
			foundAuth = true
		}
		if h.key.Value() == "Content-Type" && h.value.Value() == "application/json" {
			foundCT = true
		}
	}
	if !foundAuth {
		t.Error("Authorization header not loaded from history")
	}
	if !foundCT {
		t.Error("Content-Type header not loaded from history")
	}

	// Check body was loaded
	if m.body.Value() != `{"data":"test"}` {
		t.Errorf("body = %q, want %q", m.body.Value(), `{"data":"test"}`)
	}
}

// TestSaveAndLoadHistory tests that history can be saved and loaded from disk
func TestSaveAndLoadHistory(t *testing.T) {
	// Create a temp directory for testing
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	t.Setenv("HOME", tmpDir)
	defer func() {
		_ = os.Setenv("HOME", origHome)
	}()

	// Create test history
	entries := []historyEntry{
		{Method: "GET", URL: "https://test1.com"},
		{Method: "POST", URL: "https://test2.com", Body: `{"key":"value"}`},
	}

	// Save history
	err := saveHistory(entries)
	if err != nil {
		t.Fatalf("saveHistory failed: %v", err)
	}

	// Verify file exists
	historyPath := filepath.Join(tmpDir, ".getboy", "history.json")
	if _, err := os.Stat(historyPath); os.IsNotExist(err) {
		t.Fatal("history file was not created")
	}

	// Load history
	loaded, err := loadHistory()
	if err != nil {
		t.Fatalf("loadHistory failed: %v", err)
	}

	if len(loaded) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(loaded))
	}
	if loaded[0].URL != "https://test1.com" {
		t.Errorf("first URL = %q, want %q", loaded[0].URL, "https://test1.com")
	}
	if loaded[1].Body != `{"key":"value"}` {
		t.Errorf("second body = %q, want %q", loaded[1].Body, `{"key":"value"}`)
	}
}

// TestEmptySidebarShowsMessage tests that empty history shows a message
func TestEmptySidebarShowsMessage(t *testing.T) {
	m := New().(model)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = updated.(model)

	// With no history, sidebar should show empty message
	m.history = nil
	m.updateSidebarItems()

	view := m.viewSidebar()
	if !containsAny(view, "No history yet", "Send a request") {
		t.Error("empty sidebar should show helpful message")
	}
}

// TestSavedTabShowsComingSoon tests that Saved tab shows coming soon message
func TestSavedTabShowsComingSoon(t *testing.T) {
	m := New().(model)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = updated.(model)

	// Switch to Saved tab
	m.sidebarTab = sidebarSaved

	view := m.viewSidebar()
	if !containsAny(view, "coming soon") {
		t.Error("Saved tab should show 'coming soon' message")
	}
}

// containsAny returns true if s contains any of the substrings
func containsAny(s string, substrings ...string) bool {
	for _, sub := range substrings {
		if containsString(s, sub) {
			return true
		}
	}
	return false
}

// containsString is a simple substring check
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStringHelper(s, substr))
}

func containsStringHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
