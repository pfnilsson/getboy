package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestHeadersInitialization tests that headers are initialized correctly
func TestHeadersInitialization(t *testing.T) {
	m := New().(model)

	if len(m.headers) != 1 {
		t.Errorf("expected 1 initial header row, got %d", len(m.headers))
	}

	if m.headers[0].key.Value() != "" {
		t.Error("initial header key should be empty")
	}

	if m.headers[0].value.Value() != "" {
		t.Error("initial header value should be empty")
	}
}

// TestHeadersNavigation tests j/k navigation between header rows
func TestHeadersNavigation(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.activeTab = tabHeaders
	m.editorPart = edHeaders
	m.headerIdx = 0

	// Add some header rows
	m.headers[0].key.SetValue("Content-Type")
	m.headers[0].value.SetValue("application/json")
	m.headers = append(m.headers, newHeaderRow())
	m.headers[1].key.SetValue("Authorization")
	m.headers[1].value.SetValue("Bearer token")
	m.headers = append(m.headers, newHeaderRow()) // Empty row at end

	// Navigate down with j
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(model)
	if m.headerIdx != 1 {
		t.Errorf("after 'j' headerIdx = %d, want 1", m.headerIdx)
	}

	// Navigate up with k
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = updated.(model)
	if m.headerIdx != 0 {
		t.Errorf("after 'k' headerIdx = %d, want 0", m.headerIdx)
	}
}

// TestHeadersFieldNavigation tests h/l navigation between key and value
func TestHeadersFieldNavigation(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.activeTab = tabHeaders
	m.editorPart = edHeaders
	m.headerField = headerKey

	// Navigate to value with l
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	m = updated.(model)
	if m.headerField != headerValue {
		t.Errorf("after 'l' headerField = %v, want headerValue", m.headerField)
	}

	// Navigate back to key with h
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	m = updated.(model)
	if m.headerField != headerKey {
		t.Errorf("after 'h' headerField = %v, want headerKey", m.headerField)
	}
}

// TestHeadersInsertMode tests typing in header fields
func TestHeadersInsertMode(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.activeTab = tabHeaders
	m.editorPart = edHeaders
	m.headerIdx = 0
	m.headerField = headerKey

	// Enter insert mode
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
	m = updated.(model)
	if !m.insertMode {
		t.Error("should be in insert mode after 'i'")
	}

	// Type a character in key field
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'X'}})
	m = updated.(model)
	if m.headers[0].key.Value() != "X" {
		t.Errorf("key value = %q, want 'X'", m.headers[0].key.Value())
	}
}

// TestHeadersTabBetweenFields tests tab navigation in insert mode
func TestHeadersTabBetweenFields(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.activeTab = tabHeaders
	m.editorPart = edHeaders
	m.headerIdx = 0
	m.headerField = headerKey
	m.insertMode = true
	m.applyFocus()

	// Add a second row so we can test tab navigation between rows
	m.addHeaderRow()
	m.headerIdx = 0
	m.headerField = headerKey

	// Tab should move from key to value
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(model)
	if m.headerField != headerValue {
		t.Errorf("after tab headerField = %v, want headerValue", m.headerField)
	}

	// Another tab should move to next row's key (row exists now)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(model)
	if m.headerIdx != 1 {
		t.Errorf("after second tab headerIdx = %d, want 1", m.headerIdx)
	}
	if m.headerField != headerKey {
		t.Errorf("after second tab headerField = %v, want headerKey", m.headerField)
	}
}

// TestHeadersRawModeToggle tests toggling raw mode
func TestHeadersRawModeToggle(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.activeTab = tabHeaders

	if m.headersRaw {
		t.Error("headersRaw should be false initially")
	}

	// Toggle raw mode with 'r'
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	m = updated.(model)
	if !m.headersRaw {
		t.Error("headersRaw should be true after 'r'")
	}

	// Toggle back
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	m = updated.(model)
	if m.headersRaw {
		t.Error("headersRaw should be false after second 'r'")
	}
}

// TestHeadersToRaw tests conversion to raw format
func TestHeadersToRaw(t *testing.T) {
	m := New().(model)
	m.headers[0].key.SetValue("Content-Type")
	m.headers[0].value.SetValue("application/json")
	m.headers = append(m.headers, newHeaderRow())
	m.headers[1].key.SetValue("Authorization")
	m.headers[1].value.SetValue("Bearer token")

	raw := m.headersToRaw()

	if !strings.Contains(raw, "Content-Type: application/json") {
		t.Error("raw should contain 'Content-Type: application/json'")
	}
	if !strings.Contains(raw, "Authorization: Bearer token") {
		t.Error("raw should contain 'Authorization: Bearer token'")
	}
}

// TestHeadersFromRaw tests parsing from raw format
func TestHeadersFromRaw(t *testing.T) {
	m := New().(model)

	raw := "Content-Type: application/json\nAuthorization: Bearer token"
	m.headersFromRaw(raw)

	if len(m.headers) != 2 {
		t.Errorf("expected 2 headers, got %d", len(m.headers))
	}

	if m.headers[0].key.Value() != "Content-Type" {
		t.Errorf("first header key = %q, want 'Content-Type'", m.headers[0].key.Value())
	}
	if m.headers[0].value.Value() != "application/json" {
		t.Errorf("first header value = %q, want 'application/json'", m.headers[0].value.Value())
	}

	if m.headers[1].key.Value() != "Authorization" {
		t.Errorf("second header key = %q, want 'Authorization'", m.headers[1].key.Value())
	}
	if m.headers[1].value.Value() != "Bearer token" {
		t.Errorf("second header value = %q, want 'Bearer token'", m.headers[1].value.Value())
	}
}

// TestHeadersDeleteRow tests deleting a header row
func TestHeadersDeleteRow(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.activeTab = tabHeaders
	m.editorPart = edHeaders

	// Add some header rows
	m.headers[0].key.SetValue("Header1")
	m.headers[0].value.SetValue("Value1")
	m.headers = append(m.headers, newHeaderRow())
	m.headers[1].key.SetValue("Header2")
	m.headers[1].value.SetValue("Value2")
	m.headers = append(m.headers, newHeaderRow()) // Empty row

	m.headerIdx = 1 // Select second row

	// Delete with 'd'
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	m = updated.(model)

	if len(m.headers) != 2 {
		t.Errorf("expected 2 headers after delete, got %d", len(m.headers))
	}

	// First row should still be Header1
	if m.headers[0].key.Value() != "Header1" {
		t.Errorf("first header key = %q, want 'Header1'", m.headers[0].key.Value())
	}
}

// TestHeadersAddRow tests that 'a' adds a new header row
func TestHeadersAddRow(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.activeTab = tabHeaders
	m.editorPart = edHeaders

	// Initial state has one empty row
	if len(m.headers) != 1 {
		t.Errorf("expected 1 header, got %d", len(m.headers))
	}

	// Press 'a' to add a row
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	m = updated.(model)

	// Should now have 2 rows
	if len(m.headers) != 2 {
		t.Errorf("expected 2 headers after 'a', got %d", len(m.headers))
	}

	// Should be on the new row
	if m.headerIdx != 1 {
		t.Errorf("headerIdx = %d, want 1", m.headerIdx)
	}
}

// TestHeadersViewRendering tests that headers tab renders correctly
func TestHeadersViewRendering(t *testing.T) {
	m := New().(model)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = updated.(model)
	m.pane = paneEditor
	m.activeTab = tabHeaders

	m.headers[0].key.SetValue("Content-Type")
	m.headers[0].value.SetValue("application/json")

	view := m.viewHeadersTab()

	// Should contain the header
	if !strings.Contains(view, "Content-Type") {
		t.Error("view should contain 'Content-Type'")
	}
	if !strings.Contains(view, "application/json") {
		t.Error("view should contain 'application/json'")
	}

	// Should show mode toggle indicator
	if !strings.Contains(view, "Structured") || !strings.Contains(view, "Raw") {
		t.Error("view should contain mode toggle indicator [Structured / Raw]")
	}
}

// TestHeadersRawModeTextareaInput tests that keys go to textarea in raw mode
func TestHeadersRawModeTextareaInput(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.activeTab = tabHeaders
	m.editorPart = edHeaders
	m.headersRaw = true
	m.insertMode = true
	m.applyFocus()

	// Type a character - should go to raw textarea
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'X'}})
	m = updated.(model)

	if m.headersRawText.Value() != "X" {
		t.Errorf("headersRawText.Value() = %q, want 'X'", m.headersRawText.Value())
	}
}

// TestHeadersRawModeSync tests that content syncs when toggling modes
func TestHeadersRawModeSync(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.activeTab = tabHeaders

	// Set up structured headers
	m.headers[0].key.SetValue("Content-Type")
	m.headers[0].value.SetValue("application/json")
	m.headers = append(m.headers, newHeaderRow())
	m.headers[1].key.SetValue("Authorization")
	m.headers[1].value.SetValue("Bearer token")

	// Toggle to raw mode - should sync content
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	m = updated.(model)

	rawValue := m.headersRawText.Value()
	if !strings.Contains(rawValue, "Content-Type: application/json") {
		t.Error("raw textarea should contain 'Content-Type: application/json'")
	}
	if !strings.Contains(rawValue, "Authorization: Bearer token") {
		t.Error("raw textarea should contain 'Authorization: Bearer token'")
	}

	// Modify raw text
	m.headersRawText.SetValue("X-Custom: header\nX-Another: value")

	// Toggle back to structured mode - should parse raw text
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	m = updated.(model)

	if len(m.headers) != 2 {
		t.Errorf("expected 2 headers after parsing, got %d", len(m.headers))
	}
	if m.headers[0].key.Value() != "X-Custom" {
		t.Errorf("first header key = %q, want 'X-Custom'", m.headers[0].key.Value())
	}
	if m.headers[0].value.Value() != "header" {
		t.Errorf("first header value = %q, want 'header'", m.headers[0].value.Value())
	}
}

// TestHeadersRawModeFocus tests that applyFocus focuses raw textarea in raw mode
func TestHeadersRawModeFocus(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.activeTab = tabHeaders
	m.editorPart = edHeaders
	m.headersRaw = true
	m.insertMode = true

	m.applyFocus()

	if !m.headersRawText.Focused() {
		t.Error("headersRawText should be focused in raw mode with insert mode")
	}

	// Structured headers should not be focused
	for i, h := range m.headers {
		if h.key.Focused() {
			t.Errorf("header[%d].key should not be focused in raw mode", i)
		}
		if h.value.Focused() {
			t.Errorf("header[%d].value should not be focused in raw mode", i)
		}
	}
}

// TestHeadersRawModeNoTabNavigation tests that tab doesn't navigate fields in raw mode
func TestHeadersRawModeNoTabNavigation(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.activeTab = tabHeaders
	m.editorPart = edHeaders
	m.headersRaw = true
	m.insertMode = true
	m.headerField = headerKey
	m.applyFocus()

	// Press tab - should be passed to textarea, not change headerField
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(model)

	// Tab should insert into textarea (or be ignored by textarea)
	// but should NOT navigate between key/value fields
	// The headerField should remain unchanged (raw mode ignores structured navigation)
	if m.headerField != headerKey {
		t.Errorf("headerField changed in raw mode, got %v", m.headerField)
	}
}
