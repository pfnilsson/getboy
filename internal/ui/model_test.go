package ui

import (
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

// TestEditorNavigation tests that editor part focus cycles correctly
func TestEditorNavigation(t *testing.T) {
	m := New().(model)
	m.pane = paneEditor

	// Start at method
	if m.editorPart != edMethod {
		t.Errorf("initial editorPart = %v, want %v", m.editorPart, edMethod)
	}

	// Next should go to URL
	m.nextEditorPart()
	if m.editorPart != edURL {
		t.Errorf("after nextEditorPart() = %v, want %v", m.editorPart, edURL)
	}

	// Next should go to body
	m.nextEditorPart()
	if m.editorPart != edBody {
		t.Errorf("after nextEditorPart() = %v, want %v", m.editorPart, edBody)
	}

	// Next should cycle back to method
	m.nextEditorPart()
	if m.editorPart != edMethod {
		t.Errorf("after nextEditorPart() = %v, want %v", m.editorPart, edMethod)
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
