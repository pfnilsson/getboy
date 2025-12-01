package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestEnsureURL tests the URL normalization logic
func TestEnsureURL(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "already has https",
			input: "https://example.com",
			want:  "https://example.com",
		},
		{
			name:  "already has http",
			input: "http://example.com",
			want:  "http://example.com",
		},
		{
			name:  "no protocol",
			input: "example.com",
			want:  "https://example.com",
		},
		{
			name:  "with whitespace",
			input: "  example.com  ",
			want:  "https://example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := model{}
			got := m.ensureURL(tt.input)
			if got != tt.want {
				t.Errorf("ensureURL(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestPaneNavigation tests that pane focus cycles correctly
func TestPaneNavigation(t *testing.T) {
	m := New().(model)

	// Start at sidebar
	if m.pane != paneSidebar {
		t.Errorf("initial pane = %v, want %v", m.pane, paneSidebar)
	}

	// Next should go to editor
	m.nextPane()
	if m.pane != paneEditor {
		t.Errorf("after nextPane() = %v, want %v", m.pane, paneEditor)
	}

	// Next should go to response
	m.nextPane()
	if m.pane != paneResponse {
		t.Errorf("after nextPane() = %v, want %v", m.pane, paneResponse)
	}

	// Next should cycle back to sidebar
	m.nextPane()
	if m.pane != paneSidebar {
		t.Errorf("after nextPane() = %v, want %v", m.pane, paneSidebar)
	}

	// Test prevPane
	m.prevPane()
	if m.pane != paneResponse {
		t.Errorf("after prevPane() = %v, want %v", m.pane, paneResponse)
	}
}

// TestEditorNavigation tests that editor part focus cycles correctly within tabs
func TestEditorNavigation(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor

	// Overview tab: should cycle between method and URL
	m.activeTab = tabOverview
	m.editorPart = edMethod

	// Next should go to URL
	m.nextEditorPart()
	if m.editorPart != edURL {
		t.Errorf("after nextEditorPart() = %v, want %v", m.editorPart, edURL)
	}

	// Next should cycle back to method (only 2 fields in overview)
	m.nextEditorPart()
	if m.editorPart != edMethod {
		t.Errorf("after nextEditorPart() = %v, want %v", m.editorPart, edMethod)
	}

	// Body tab: should stay on body (only 1 field)
	m.activeTab = tabBody
	m.editorPart = edBody
	m.nextEditorPart()
	if m.editorPart != edBody {
		t.Errorf("in body tab, editorPart = %v, want %v", m.editorPart, edBody)
	}
}

// TestWindowSizeMsg tests that the model responds to window resize
func TestWindowSizeMsg(t *testing.T) {
	m := New().(model)

	msg := tea.WindowSizeMsg{Width: 120, Height: 40}
	updated, _ := m.Update(msg)
	m = updated.(model)

	if m.width != 120 {
		t.Errorf("width = %d, want 120", m.width)
	}
	if m.height != 40 {
		t.Errorf("height = %d, want 40", m.height)
	}
}

// TestLayoutCalculations tests the layout width calculations
func TestLayoutCalculations(t *testing.T) {
	m := New().(model)
	m.width = 100
	m.height = 40

	sidebarW := m.sidebarWidth()
	if sidebarW != 28 {
		t.Errorf("sidebarWidth() = %d, want 28", sidebarW)
	}

	rightW := m.rightPaneWidth()
	// width(100) - sidebar(28) - borders(4) = 68
	expectedRight := 68
	if rightW != expectedRight {
		t.Errorf("rightPaneWidth() = %d, want %d", rightW, expectedRight)
	}
}

// TestEditorJKNavigation tests that j/k keys navigate editor parts when not in insert mode
func TestEditorJKNavigation(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.activeTab = tabOverview
	m.insertMode = false

	// Start at method
	if m.editorPart != edMethod {
		t.Fatalf("initial editorPart = %v, want %v", m.editorPart, edMethod)
	}

	// Press 'j' should go to URL
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(model)
	if m.editorPart != edURL {
		t.Errorf("after 'j' editorPart = %v, want %v", m.editorPart, edURL)
	}

	// Press 'j' again should cycle back to method (only 2 fields in overview tab)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(model)
	if m.editorPart != edMethod {
		t.Errorf("after 'j' editorPart = %v, want %v", m.editorPart, edMethod)
	}

	// Press 'k' should go to URL
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = updated.(model)
	if m.editorPart != edURL {
		t.Errorf("after 'k' editorPart = %v, want %v", m.editorPart, edURL)
	}
}

// TestInsertMode tests entering and exiting insert mode
func TestInsertMode(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.insertMode = false

	// Press 'i' should enter insert mode
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
	m = updated.(model)
	if !m.insertMode {
		t.Error("after 'i' insertMode should be true")
	}

	// Press 'esc' should exit insert mode
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEscape})
	m = updated.(model)
	if m.insertMode {
		t.Error("after 'esc' insertMode should be false")
	}
}

// TestInsertModeBlocksNavigation tests that navigation keys are passed to input in insert mode
func TestInsertModeBlocksNavigation(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.editorPart = edURL
	m.insertMode = true

	// Press tab - should be passed to input, not change pane
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(model)

	// Should still be in insert mode and on editor pane
	if !m.insertMode {
		t.Error("insertMode should remain true - tab should go to input")
	}
	if m.pane != paneEditor {
		t.Errorf("pane = %v, want %v (should not change in insert mode)", m.pane, paneEditor)
	}
}

// TestNavigationExitsInsertMode tests that using nextPane/prevPane exits insert mode
func TestNavigationExitsInsertMode(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.insertMode = true

	// Use nextPane directly (this is what tab does when not in insert mode)
	m.nextPane()

	if m.insertMode {
		t.Error("insertMode should be false after nextPane()")
	}
	if m.pane != paneResponse {
		t.Errorf("pane = %v, want %v", m.pane, paneResponse)
	}
}

// TestNumberKeysInInsertMode tests that number keys are passed to input in insert mode
func TestNumberKeysInInsertMode(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.editorPart = edURL
	m.insertMode = true
	m.applyFocus() // Focus the URL input

	// Press '1' - should be passed to URL input, not change pane
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	m = updated.(model)

	// Should still be in insert mode and on editor pane
	if !m.insertMode {
		t.Error("insertMode should remain true - numbers should go to input")
	}
	if m.pane != paneEditor {
		t.Errorf("pane = %v, want %v (should not change in insert mode)", m.pane, paneEditor)
	}
	// The URL input should have received the '1'
	if m.url.Value() != "1" {
		t.Errorf("url.Value() = %q, want %q", m.url.Value(), "1")
	}
}

// TestInsertModeKeysPassToInput tests that keys are passed to inputs in insert mode
func TestInsertModeKeysPassToInput(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.editorPart = edURL
	m.insertMode = true
	m.applyFocus() // Focus the URL input

	// Type 'j' - should go to URL input, not navigate
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(model)

	// Should still be on URL (not moved to body)
	if m.editorPart != edURL {
		t.Errorf("editorPart = %v, want %v (should not navigate in insert mode)", m.editorPart, edURL)
	}

	// The URL input should have received the 'j'
	if m.url.Value() != "j" {
		t.Errorf("url.Value() = %q, want %q", m.url.Value(), "j")
	}
}

// TestEditorArrowNavigation tests that arrow keys navigate editor parts
func TestEditorArrowNavigation(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.insertMode = false

	// Start at method
	if m.editorPart != edMethod {
		t.Fatalf("initial editorPart = %v, want %v", m.editorPart, edMethod)
	}

	// Press 'down' should go to URL
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(model)
	if m.editorPart != edURL {
		t.Errorf("after 'down' editorPart = %v, want %v", m.editorPart, edURL)
	}

	// Press 'up' should go back to method
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(model)
	if m.editorPart != edMethod {
		t.Errorf("after 'up' editorPart = %v, want %v", m.editorPart, edMethod)
	}
}

// TestEditorTabNavigation tests that [ and ] navigate tabs
func TestEditorTabNavigation(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.insertMode = false

	// Start at overview tab
	if m.activeTab != tabOverview {
		t.Fatalf("initial activeTab = %v, want %v", m.activeTab, tabOverview)
	}

	// Press ']' should go to headers
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{']'}})
	m = updated.(model)
	if m.activeTab != tabHeaders {
		t.Errorf("after ']' activeTab = %v, want %v", m.activeTab, tabHeaders)
	}

	// Press '[' should go back to overview
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'['}})
	m = updated.(model)
	if m.activeTab != tabOverview {
		t.Errorf("after '[' activeTab = %v, want %v", m.activeTab, tabOverview)
	}
}

// TestEditorLeftRightNavigation tests that left/right arrow keys navigate tabs
func TestEditorLeftRightNavigation(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.insertMode = false

	// Start at overview tab
	if m.activeTab != tabOverview {
		t.Fatalf("initial activeTab = %v, want %v", m.activeTab, tabOverview)
	}

	// Press 'right' should go to headers
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRight})
	m = updated.(model)
	if m.activeTab != tabHeaders {
		t.Errorf("after 'right' activeTab = %v, want %v", m.activeTab, tabHeaders)
	}

	// Press 'left' should go back to overview
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyLeft})
	m = updated.(model)
	if m.activeTab != tabOverview {
		t.Errorf("after 'left' activeTab = %v, want %v", m.activeTab, tabOverview)
	}
}

// TestApplyFocusInInsertMode tests that applyFocus only focuses inputs in insert mode
func TestApplyFocusInInsertMode(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.editorPart = edURL

	// Not in insert mode - inputs should be blurred
	m.insertMode = false
	m.applyFocus()

	// Method is now a dropdown, not a text input - no focus check needed
	if m.url.Focused() {
		t.Error("url should not be focused when not in insert mode")
	}
	if m.body.Focused() {
		t.Error("body should not be focused when not in insert mode")
	}

	// In insert mode - selected input should be focused
	m.insertMode = true
	m.applyFocus()

	// Method is a dropdown, URL should be focused since editorPart is edURL
	if !m.url.Focused() {
		t.Error("url should be focused when editorPart is edURL and in insert mode")
	}
	if m.body.Focused() {
		t.Error("body should not be focused when editorPart is edURL")
	}
}

// TestEnterKeyExecutesRequest tests that enter key executes HTTP request in editor pane
func TestEnterKeyExecutesRequest(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.url.SetValue("https://example.com")

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(model)

	if !m.loading {
		t.Error("loading should be true after pressing enter")
	}
	if cmd == nil {
		t.Error("cmd should not be nil - should return HTTP command")
	}
}

// TestEnterKeyRequiresURL tests that enter key shows error when URL is empty
func TestEnterKeyRequiresURL(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.url.SetValue("")

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(model)

	if m.loading {
		t.Error("loading should be false when URL is empty")
	}
	if cmd != nil {
		t.Error("cmd should be nil when URL is empty")
	}
	if m.status != "Enter a URL first" {
		t.Errorf("status = %q, want %q", m.status, "Enter a URL first")
	}
}

// TestEnterKeyLoadsItemFromSidebar tests that enter key loads selected item in sidebar
func TestEnterKeyLoadsItemFromSidebar(t *testing.T) {
	m := New().(model)
	m.pane = paneSidebar

	// Press enter to load the first item
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(model)

	// Should have loaded the first example item
	if m.methodValue() != "GET" {
		t.Errorf("method = %q, want %q", m.methodValue(), "GET")
	}
	if m.url.Value() != "https://httpbin.org/get" {
		t.Errorf("url = %q, want %q", m.url.Value(), "https://httpbin.org/get")
	}
}

// TestHttpDoneMsg tests handling of HTTP response message
func TestHttpDoneMsg(t *testing.T) {
	t.Run("successful response", func(t *testing.T) {
		m := New().(model)
		m.loading = true

		updated, _ := m.Update(httpDoneMsg{
			Status: "200 OK",
			Body:   `{"result":"success"}`,
			Err:    nil,
		})
		m = updated.(model)

		if m.loading {
			t.Error("loading should be false after httpDoneMsg")
		}
		if m.status != "200 OK" {
			t.Errorf("status = %q, want %q", m.status, "200 OK")
		}
		if m.err != nil {
			t.Errorf("err should be nil, got %v", m.err)
		}
	})

	t.Run("error response", func(t *testing.T) {
		m := New().(model)
		m.loading = true

		testErr := &testError{msg: "connection failed"}
		updated, _ := m.Update(httpDoneMsg{
			Err: testErr,
		})
		m = updated.(model)

		if m.loading {
			t.Error("loading should be false after httpDoneMsg")
		}
		if m.status != "Request failed" {
			t.Errorf("status = %q, want %q", m.status, "Request failed")
		}
		if m.err != testErr {
			t.Errorf("err = %v, want %v", m.err, testErr)
		}
	})
}

// testError is a simple error implementation for testing
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

// TestQuitKeys tests that q and ctrl+c quit the application
func TestQuitKeys(t *testing.T) {
	tests := []struct {
		name string
		key  tea.KeyMsg
	}{
		{"q", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}},
		{"ctrl+c", tea.KeyMsg{Type: tea.KeyCtrlC}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New().(model)
			_, cmd := m.Update(tt.key)

			// Should return tea.Quit command
			if cmd == nil {
				t.Error("cmd should not be nil for quit key")
			}
		})
	}
}

// TestShiftTabNavigation tests that shift+tab navigates backwards through panes
func TestShiftTabNavigation(t *testing.T) {
	m := New().(model)
	// Start at sidebar (pane 0)
	if m.pane != paneSidebar {
		t.Fatalf("initial pane = %v, want %v", m.pane, paneSidebar)
	}

	// Shift+tab should go to response (pane 2, wrapping backwards)
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	m = updated.(model)
	if m.pane != paneResponse {
		t.Errorf("after shift+tab pane = %v, want %v", m.pane, paneResponse)
	}
}

// TestDefaultMethodIsGET tests that empty method defaults to GET
func TestDefaultMethodIsGET(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	// methodIdx defaults to 0 (GET)
	m.url.SetValue("https://example.com")

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(model)

	if !m.loading {
		t.Error("loading should be true")
	}
	if cmd == nil {
		t.Error("cmd should not be nil")
	}
	// Status should show GET (the default)
	if !strings.Contains(m.status, "GET") {
		t.Errorf("status = %q, should contain 'GET'", m.status)
	}
}

// TestMethodDropdownCycling tests that method dropdown cycles through options
func TestMethodDropdownCycling(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.setMethod("POST") // Start at POST
	m.url.SetValue("https://example.com")

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(model)

	// Status should show POST
	if !strings.Contains(m.status, "POST") {
		t.Errorf("status = %q, should contain 'POST'", m.status)
	}
}

// TestSetMethod tests setting method by name
func TestSetMethod(t *testing.T) {
	m := New().(model)

	// Set to POST
	m.setMethod("POST")
	if m.methodValue() != "POST" {
		t.Errorf("after setMethod(POST), methodValue() = %q, want POST", m.methodValue())
	}

	// Set to lowercase should work
	m.setMethod("delete")
	if m.methodValue() != "DELETE" {
		t.Errorf("after setMethod(delete), methodValue() = %q, want DELETE", m.methodValue())
	}

	// Set to unknown method should default to GET
	m.setMethod("UNKNOWN")
	if m.methodValue() != "GET" {
		t.Errorf("after setMethod(UNKNOWN), methodValue() = %q, want GET", m.methodValue())
	}
}

// TestMethodDropdownNavigation tests j/k navigation in method dropdown
func TestMethodDropdownNavigation(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.activeTab = tabOverview
	m.editorPart = edMethod
	m.insertMode = true

	// Start at GET (index 0)
	if m.methodValue() != "GET" {
		t.Fatalf("initial method = %q, want GET", m.methodValue())
	}

	// Press 'j' (down) should go to POST
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(model)
	if m.methodValue() != "POST" {
		t.Errorf("after 'j' method = %q, want POST", m.methodValue())
	}

	// Press 'k' (up) should go back to GET
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = updated.(model)
	if m.methodValue() != "GET" {
		t.Errorf("after 'k' method = %q, want GET", m.methodValue())
	}

	// Press 'k' again should wrap to OPTIONS (last item)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = updated.(model)
	if m.methodValue() != "OPTIONS" {
		t.Errorf("after 'k' (wrap) method = %q, want OPTIONS", m.methodValue())
	}
}

// TestInsertModeWithBodyInput tests insert mode with body input (textarea)
func TestInsertModeWithBodyInput(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.editorPart = edBody
	m.insertMode = true
	m.applyFocus()

	// Type '{' - should go to body input
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'{'}})
	m = updated.(model)

	if m.body.Value() != "{" {
		t.Errorf("body.Value() = %q, want %q", m.body.Value(), "{")
	}
}

// TestPrevEditorPart tests cycling backwards through editor parts within tabs
func TestPrevEditorPart(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.activeTab = tabOverview
	m.editorPart = edMethod

	// In overview tab, should wrap to URL (only 2 fields)
	m.prevEditorPart()
	if m.editorPart != edURL {
		t.Errorf("editorPart = %v, want %v", m.editorPart, edURL)
	}

	// Should go back to method
	m.prevEditorPart()
	if m.editorPart != edMethod {
		t.Errorf("editorPart = %v, want %v", m.editorPart, edMethod)
	}

	// In body tab, should stay on body (only 1 field)
	m.activeTab = tabBody
	m.editorPart = edBody
	m.prevEditorPart()
	if m.editorPart != edBody {
		t.Errorf("in body tab, editorPart = %v, want %v", m.editorPart, edBody)
	}
}

// TestTabCycling tests that tabs cycle correctly
func TestTabCycling(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor

	// Start at overview
	if m.activeTab != tabOverview {
		t.Fatalf("initial activeTab = %v, want %v", m.activeTab, tabOverview)
	}

	// Next tab
	m.nextTab()
	if m.activeTab != tabHeaders {
		t.Errorf("activeTab = %v, want %v", m.activeTab, tabHeaders)
	}

	m.nextTab()
	if m.activeTab != tabBody {
		t.Errorf("activeTab = %v, want %v", m.activeTab, tabBody)
	}

	// Should wrap to overview
	m.nextTab()
	if m.activeTab != tabOverview {
		t.Errorf("activeTab = %v, want %v", m.activeTab, tabOverview)
	}

	// Prev tab should wrap to body
	m.prevTab()
	if m.activeTab != tabBody {
		t.Errorf("activeTab = %v, want %v", m.activeTab, tabBody)
	}
}

// TestTabSwitchResetsEditorPart tests that switching tabs resets editorPart appropriately
func TestTabSwitchResetsEditorPart(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor
	m.activeTab = tabOverview
	m.editorPart = edURL // Start at URL

	// Switch to body tab - should set editorPart to edBody
	m.nextTab() // headers
	m.nextTab() // body
	if m.activeTab != tabBody {
		t.Fatalf("activeTab = %v, want %v", m.activeTab, tabBody)
	}
	if m.editorPart != edBody {
		t.Errorf("after switching to body tab, editorPart = %v, want %v", m.editorPart, edBody)
	}

	// Switch back to overview tab - should set editorPart to edMethod
	m.nextTab() // overview
	if m.activeTab != tabOverview {
		t.Fatalf("activeTab = %v, want %v", m.activeTab, tabOverview)
	}
	if m.editorPart != edMethod {
		t.Errorf("after switching to overview tab, editorPart = %v, want %v", m.editorPart, edMethod)
	}
}

// TestApplyFocusAllEditorParts tests applyFocus for all editor parts
func TestApplyFocusAllEditorParts(t *testing.T) {
	// Method is now a dropdown (not focusable), URL and Body are text inputs
	tests := []struct {
		part     editorFocus
		wantURL  bool
		wantBody bool
	}{
		{edMethod, false, false}, // Method is a dropdown, no focus
		{edURL, true, false},
		{edBody, false, true},
	}

	for _, tt := range tests {
		t.Run(string(rune('0'+tt.part)), func(t *testing.T) {
			m := New().(model)
			m.pane = paneEditor
			m.editorPart = tt.part
			m.insertMode = true
			m.applyFocus()

			if m.url.Focused() != tt.wantURL {
				t.Errorf("url.Focused() = %v, want %v", m.url.Focused(), tt.wantURL)
			}
			if m.body.Focused() != tt.wantBody {
				t.Errorf("body.Focused() = %v, want %v", m.body.Focused(), tt.wantBody)
			}
		})
	}
}

// TestResponsePaneKeyHandling tests that keys are passed to viewport in response pane
func TestResponsePaneKeyHandling(t *testing.T) {
	m := New().(model)
	m.pane = paneResponse

	// Keys should be passed to viewport (not cause errors)
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(model)

	// Should still be on response pane
	if m.pane != paneResponse {
		t.Errorf("pane = %v, want %v", m.pane, paneResponse)
	}
}
