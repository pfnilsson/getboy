package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestViewContainsPanes tests that the view contains all pane titles
func TestViewContainsPanes(t *testing.T) {
	m := New().(model)

	// Set size so layout works properly
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = updated.(model)

	view := m.View()

	// Check that all pane titles are present
	expectedTitles := []string{"Requests", "Request", "Response"}
	for _, title := range expectedTitles {
		if !strings.Contains(view, title) {
			t.Errorf("view does not contain pane title %q", title)
		}
	}

	// Check that status bar is present
	if !strings.Contains(view, "tab: panes") {
		t.Error("view does not contain status message")
	}
}

// TestViewSidebar tests the sidebar view rendering
func TestViewSidebar(t *testing.T) {
	m := New().(model)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = updated.(model)

	sidebarView := m.viewSidebar()

	// Should contain the pane title
	if !strings.Contains(sidebarView, "Requests") {
		t.Error("sidebar view does not contain 'Requests' title")
	}

	// Should contain the badge [1]
	if !strings.Contains(sidebarView, "[1]") {
		t.Error("sidebar view does not contain badge [1]")
	}
}

// TestViewEditor tests the editor view rendering
func TestViewEditor(t *testing.T) {
	m := New().(model)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = updated.(model)

	editorView := m.viewEditor()

	// Should contain the pane title
	if !strings.Contains(editorView, "Request") {
		t.Error("editor view does not contain 'Request' title")
	}

	// Should contain the badge [2]
	if !strings.Contains(editorView, "[2]") {
		t.Error("editor view does not contain badge [2]")
	}

	// Should contain tabs in title bar
	if !strings.Contains(editorView, "Overview") {
		t.Error("editor view does not contain 'Overview' tab")
	}
	if !strings.Contains(editorView, "Headers") {
		t.Error("editor view does not contain 'Headers' tab")
	}
	if !strings.Contains(editorView, "Body") {
		t.Error("editor view does not contain 'Body' tab")
	}

	// Should contain labels (when on overview tab)
	if m.activeTab == tabOverview {
		if !strings.Contains(editorView, "Method:") {
			t.Error("editor view does not contain 'Method:' label")
		}
		if !strings.Contains(editorView, "URL:") {
			t.Error("editor view does not contain 'URL:' label")
		}
	}
}

// TestViewResponse tests the response view rendering
func TestViewResponse(t *testing.T) {
	m := New().(model)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = updated.(model)

	responseView := m.viewResponse()

	// Should contain the pane title
	if !strings.Contains(responseView, "Response") {
		t.Error("response view does not contain 'Response' title")
	}

	// Should contain the badge [3]
	if !strings.Contains(responseView, "[3]") {
		t.Error("response view does not contain badge [3]")
	}
}

// TestViewLoading tests that loading state is shown
func TestViewLoading(t *testing.T) {
	m := New().(model)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = updated.(model)
	m.loading = true

	view := m.View()

	if !strings.Contains(view, "loading…") {
		t.Error("view does not show loading indicator when loading is true")
	}
}

// TestViewError tests that error messages are shown
func TestViewError(t *testing.T) {
	m := New().(model)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = updated.(model)
	m.err = &viewTestError{msg: "test error"}

	view := m.View()

	if !strings.Contains(view, "error: test error") {
		t.Error("view does not show error message when err is set")
	}
}

// viewTestError is a simple error implementation for testing
type viewTestError struct {
	msg string
}

func (e *viewTestError) Error() string {
	return e.msg
}

// TestViewInsertModeIndicator tests that insert mode shows in status bar
func TestViewInsertModeIndicator(t *testing.T) {
	m := New().(model)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = updated.(model)
	m.insertMode = true

	view := m.View()

	if !strings.Contains(view, "-- INSERT --") {
		t.Error("view does not show INSERT mode indicator when insertMode is true")
	}
}

// TestViewEditorSelectionIndicator tests that selection indicator shows in editor
func TestViewEditorSelectionIndicator(t *testing.T) {
	m := New().(model)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = updated.(model)
	m.pane = paneEditor
	m.insertMode = false

	// Test overview tab - method selection
	m.activeTab = tabOverview
	m.editorPart = edMethod
	editorView := m.viewOverviewTab()

	if !strings.Contains(editorView, "> Method:") {
		t.Error("editor view does not show selection indicator for method")
	}

	// Test overview tab - URL selection
	m.editorPart = edURL
	editorView = m.viewOverviewTab()

	if !strings.Contains(editorView, "> URL:") {
		t.Error("editor view does not show selection indicator for URL")
	}

	// Test body tab - body selection
	m.activeTab = tabBody
	m.editorPart = edBody
	bodyView := m.viewBodyTab()

	if !strings.Contains(bodyView, "> Body") {
		t.Error("body tab does not show selection indicator for body")
	}
}

// TestViewNoSelectionIndicatorInInsertMode tests that selection indicator is hidden in insert mode
func TestViewNoSelectionIndicatorInInsertMode(t *testing.T) {
	m := New().(model)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = updated.(model)
	m.pane = paneEditor
	m.activeTab = tabOverview
	m.insertMode = true
	m.editorPart = edMethod

	editorView := m.viewOverviewTab()

	// Should not contain selection indicator when in insert mode
	if strings.Contains(editorView, "> Method:") {
		t.Error("editor view should not show selection indicator in insert mode")
	}
}

// TestViewNoSelectionIndicatorInOtherPanes tests that selection indicator is hidden in other panes
func TestViewNoSelectionIndicatorInOtherPanes(t *testing.T) {
	m := New().(model)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = updated.(model)
	m.pane = paneSidebar
	m.activeTab = tabOverview
	m.insertMode = false
	m.editorPart = edMethod

	editorView := m.viewOverviewTab()

	// Should not contain selection indicator when not in editor pane
	if strings.Contains(editorView, "> Method:") {
		t.Error("editor view should not show selection indicator when not in editor pane")
	}
}
