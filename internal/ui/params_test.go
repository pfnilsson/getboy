package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestParamsInitialization tests that params are initialized correctly
func TestParamsInitialization(t *testing.T) {
	m := New().(model)

	if len(m.params) != 1 {
		t.Errorf("expected 1 initial param row, got %d", len(m.params))
	}

	if m.params[0].key.Value() != "" {
		t.Error("initial param key should be empty")
	}

	if m.params[0].value.Value() != "" {
		t.Error("initial param value should be empty")
	}
}

// TestParamsNavigation tests j/k navigation between param rows
func TestParamsNavigation(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.activeTab = tabParams
	m.editorPart = edParams
	m.paramIdx = 0

	// Add some param rows
	m.params[0].key.SetValue("page")
	m.params[0].value.SetValue("1")
	m.params = append(m.params, newParamRow())
	m.params[1].key.SetValue("limit")
	m.params[1].value.SetValue("10")
	m.params = append(m.params, newParamRow()) // Empty row at end

	// Navigate down with j
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(model)
	if m.paramIdx != 1 {
		t.Errorf("after 'j' paramIdx = %d, want 1", m.paramIdx)
	}

	// Navigate up with k
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = updated.(model)
	if m.paramIdx != 0 {
		t.Errorf("after 'k' paramIdx = %d, want 0", m.paramIdx)
	}
}

// TestParamsFieldNavigation tests l navigation to value field
func TestParamsFieldNavigation(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.activeTab = tabParams
	m.editorPart = edParams
	m.paramField = headerKey

	// Navigate to value with l
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	m = updated.(model)
	if m.paramField != headerValue {
		t.Errorf("after 'l' paramField = %v, want headerValue", m.paramField)
	}

	// 'h' now switches to Headers tab, use tab in insert mode to navigate fields
}

// TestParamsInsertMode tests typing in param fields
func TestParamsInsertMode(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.activeTab = tabParams
	m.editorPart = edParams
	m.paramIdx = 0
	m.paramField = headerKey

	// Enter insert mode
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
	m = updated.(model)
	if !m.insertMode {
		t.Error("should be in insert mode after 'i'")
	}

	// Type a character in key field
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	m = updated.(model)
	if m.params[0].key.Value() != "q" {
		t.Errorf("key value = %q, want 'q'", m.params[0].key.Value())
	}
}

// TestParamsTabBetweenFields tests tab navigation in insert mode
func TestParamsTabBetweenFields(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.activeTab = tabParams
	m.editorPart = edParams
	m.paramIdx = 0
	m.paramField = headerKey
	m.insertMode = true
	m.applyFocus()

	// Add a second row so we can test tab navigation between rows
	m.addParamRow()
	m.paramIdx = 0
	m.paramField = headerKey

	// Tab should move from key to value
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(model)
	if m.paramField != headerValue {
		t.Errorf("after tab paramField = %v, want headerValue", m.paramField)
	}

	// Another tab should move to next row's key (row exists now)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(model)
	if m.paramIdx != 1 {
		t.Errorf("after second tab paramIdx = %d, want 1", m.paramIdx)
	}
	if m.paramField != headerKey {
		t.Errorf("after second tab paramField = %v, want headerKey", m.paramField)
	}
}

// TestParamsDeleteRow tests deleting a param row
func TestParamsDeleteRow(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.activeTab = tabParams
	m.editorPart = edParams

	// Add some param rows
	m.params[0].key.SetValue("param1")
	m.params[0].value.SetValue("value1")
	m.params = append(m.params, newParamRow())
	m.params[1].key.SetValue("param2")
	m.params[1].value.SetValue("value2")
	m.params = append(m.params, newParamRow()) // Empty row

	m.paramIdx = 1 // Select second row

	// Delete with 'd'
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	m = updated.(model)

	if len(m.params) != 2 {
		t.Errorf("expected 2 params after delete, got %d", len(m.params))
	}

	// First row should still be param1
	if m.params[0].key.Value() != "param1" {
		t.Errorf("first param key = %q, want 'param1'", m.params[0].key.Value())
	}
}

// TestParamsAddRow tests that 'a' adds a new param row
func TestParamsAddRow(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.activeTab = tabParams
	m.editorPart = edParams

	// Initial state has one empty row
	if len(m.params) != 1 {
		t.Errorf("expected 1 param, got %d", len(m.params))
	}

	// Press 'a' to add a row
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	m = updated.(model)

	// Should now have 2 rows
	if len(m.params) != 2 {
		t.Errorf("expected 2 params after 'a', got %d", len(m.params))
	}

	// Should be on the new row
	if m.paramIdx != 1 {
		t.Errorf("paramIdx = %d, want 1", m.paramIdx)
	}
}

// TestParamsViewRendering tests that params tab renders correctly
func TestParamsViewRendering(t *testing.T) {
	m := New().(model)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = updated.(model)
	m.pane = paneEditor
	m.activeTab = tabParams

	m.params[0].key.SetValue("page")
	m.params[0].value.SetValue("1")

	view := m.viewParamsTab()

	// Should contain the param
	if !strings.Contains(view, "page") {
		t.Error("view should contain 'page'")
	}
	// The value should be visible
	if !strings.Contains(view, "1") {
		t.Error("view should contain '1'")
	}
}

// TestSyncParamsFromURL tests parsing URL query string into params
func TestSyncParamsFromURL(t *testing.T) {
	t.Run("single param", func(t *testing.T) {
		m := New().(model)
		m.url.SetValue("https://example.com?page=1")
		m.syncParamsFromURL()

		if len(m.params) != 1 {
			t.Fatalf("expected 1 param, got %d", len(m.params))
		}
		if m.params[0].key.Value() != "page" {
			t.Errorf("param key = %q, want 'page'", m.params[0].key.Value())
		}
		if m.params[0].value.Value() != "1" {
			t.Errorf("param value = %q, want '1'", m.params[0].value.Value())
		}
	})

	t.Run("multiple params", func(t *testing.T) {
		m := New().(model)
		m.url.SetValue("https://example.com?page=1&limit=10")
		m.syncParamsFromURL()

		if len(m.params) != 2 {
			t.Fatalf("expected 2 params, got %d", len(m.params))
		}

		// Check that both params exist (order may vary due to map iteration)
		foundPage, foundLimit := false, false
		for _, p := range m.params {
			if p.key.Value() == "page" && p.value.Value() == "1" {
				foundPage = true
			}
			if p.key.Value() == "limit" && p.value.Value() == "10" {
				foundLimit = true
			}
		}
		if !foundPage {
			t.Error("should have page=1 param")
		}
		if !foundLimit {
			t.Error("should have limit=10 param")
		}
	})

	t.Run("no query params", func(t *testing.T) {
		m := New().(model)
		m.url.SetValue("https://example.com")
		m.syncParamsFromURL()

		// Should reset to single empty row
		if len(m.params) != 1 {
			t.Fatalf("expected 1 param, got %d", len(m.params))
		}
		if m.params[0].key.Value() != "" {
			t.Errorf("param key = %q, want empty", m.params[0].key.Value())
		}
	})

	t.Run("empty URL", func(t *testing.T) {
		m := New().(model)
		m.params[0].key.SetValue("existing")
		m.url.SetValue("")
		m.syncParamsFromURL()

		// Should leave params unchanged
		if m.params[0].key.Value() != "existing" {
			t.Error("empty URL should not modify params")
		}
	})

	t.Run("URL encoded values", func(t *testing.T) {
		m := New().(model)
		m.url.SetValue("https://example.com?name=hello%20world")
		m.syncParamsFromURL()

		if len(m.params) != 1 {
			t.Fatalf("expected 1 param, got %d", len(m.params))
		}
		if m.params[0].value.Value() != "hello world" {
			t.Errorf("param value = %q, want 'hello world'", m.params[0].value.Value())
		}
	})
}

// TestSyncURLFromParams tests updating URL from params
func TestSyncURLFromParams(t *testing.T) {
	t.Run("single param", func(t *testing.T) {
		m := New().(model)
		m.url.SetValue("https://example.com")
		m.params[0].key.SetValue("page")
		m.params[0].value.SetValue("1")
		m.syncURLFromParams()

		url := m.url.Value()
		if !strings.Contains(url, "page=1") {
			t.Errorf("URL = %q, should contain 'page=1'", url)
		}
	})

	t.Run("multiple params", func(t *testing.T) {
		m := New().(model)
		m.url.SetValue("https://example.com")
		m.params[0].key.SetValue("page")
		m.params[0].value.SetValue("1")
		m.params = append(m.params, newParamRow())
		m.params[1].key.SetValue("limit")
		m.params[1].value.SetValue("10")
		m.syncURLFromParams()

		url := m.url.Value()
		if !strings.Contains(url, "page=1") {
			t.Errorf("URL = %q, should contain 'page=1'", url)
		}
		if !strings.Contains(url, "limit=10") {
			t.Errorf("URL = %q, should contain 'limit=10'", url)
		}
	})

	t.Run("empty key skipped", func(t *testing.T) {
		m := New().(model)
		m.url.SetValue("https://example.com")
		m.params[0].key.SetValue("")
		m.params[0].value.SetValue("orphan")
		m.params = append(m.params, newParamRow())
		m.params[1].key.SetValue("valid")
		m.params[1].value.SetValue("value")
		m.syncURLFromParams()

		url := m.url.Value()
		if strings.Contains(url, "orphan") {
			t.Errorf("URL = %q, should not contain empty key param", url)
		}
		if !strings.Contains(url, "valid=value") {
			t.Errorf("URL = %q, should contain 'valid=value'", url)
		}
	})

	t.Run("replaces existing query string", func(t *testing.T) {
		m := New().(model)
		m.url.SetValue("https://example.com?old=value")
		m.params[0].key.SetValue("new")
		m.params[0].value.SetValue("param")
		m.syncURLFromParams()

		url := m.url.Value()
		if strings.Contains(url, "old=value") {
			t.Errorf("URL = %q, should not contain old query string", url)
		}
		if !strings.Contains(url, "new=param") {
			t.Errorf("URL = %q, should contain 'new=param'", url)
		}
	})

	t.Run("URL encoding", func(t *testing.T) {
		m := New().(model)
		m.url.SetValue("https://example.com")
		m.params[0].key.SetValue("name")
		m.params[0].value.SetValue("hello world")
		m.syncURLFromParams()

		url := m.url.Value()
		// URL should encode the space
		if !strings.Contains(url, "name=hello") {
			t.Errorf("URL = %q, should contain encoded param", url)
		}
	})
}

// TestParamsTwoWaySync tests that entering params tab syncs from URL
func TestParamsTwoWaySync(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.activeTab = tabOverview

	// Set URL with query params
	m.url.SetValue("https://example.com?foo=bar")

	// Switch to params tab - should sync from URL
	m.nextTab() // params
	if m.activeTab != tabParams {
		t.Fatalf("activeTab = %v, want tabParams", m.activeTab)
	}

	// Params should be synced from URL
	if len(m.params) != 1 {
		t.Fatalf("expected 1 param, got %d", len(m.params))
	}
	if m.params[0].key.Value() != "foo" {
		t.Errorf("param key = %q, want 'foo'", m.params[0].key.Value())
	}
	if m.params[0].value.Value() != "bar" {
		t.Errorf("param value = %q, want 'bar'", m.params[0].value.Value())
	}
}

// TestEnterExitsInsertModeInParams tests that enter exits insert mode for params
func TestEnterExitsInsertModeInParams(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.activeTab = tabParams
	m.editorPart = edParams
	m.insertMode = true
	m.applyFocus()

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(model)

	if m.insertMode {
		t.Error("insertMode should be false after enter in params")
	}
}

// TestParamsUpDownInInsertMode tests up/down navigation in insert mode
func TestParamsUpDownInInsertMode(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.activeTab = tabParams
	m.editorPart = edParams
	m.insertMode = true

	// Add rows
	m.params[0].key.SetValue("first")
	m.params = append(m.params, newParamRow())
	m.params[1].key.SetValue("second")
	m.paramIdx = 0
	m.applyFocus()

	// Press down - should move to next row
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(model)
	if m.paramIdx != 1 {
		t.Errorf("after down paramIdx = %d, want 1", m.paramIdx)
	}

	// Press up - should move back
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(model)
	if m.paramIdx != 0 {
		t.Errorf("after up paramIdx = %d, want 0", m.paramIdx)
	}
}

// TestParamsShiftTabNavigation tests shift+tab backwards navigation in insert mode
func TestParamsShiftTabNavigation(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.activeTab = tabParams
	m.editorPart = edParams
	m.insertMode = true

	// Add rows
	m.params[0].key.SetValue("first")
	m.params = append(m.params, newParamRow())
	m.params[1].key.SetValue("second")
	m.paramIdx = 1
	m.paramField = headerKey
	m.applyFocus()

	// Press shift+tab - should move to previous row's value
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	m = updated.(model)
	if m.paramIdx != 0 {
		t.Errorf("after shift+tab paramIdx = %d, want 0", m.paramIdx)
	}
	if m.paramField != headerValue {
		t.Errorf("after shift+tab paramField = %v, want headerValue", m.paramField)
	}

	// Another shift+tab - should move to key
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	m = updated.(model)
	if m.paramField != headerKey {
		t.Errorf("after second shift+tab paramField = %v, want headerKey", m.paramField)
	}
}
