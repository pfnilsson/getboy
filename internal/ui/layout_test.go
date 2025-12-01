package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestLayoutDimensions tests that panes maintain correct dimensions
func TestLayoutDimensions(t *testing.T) {
	tests := []struct {
		name         string
		width        int
		height       int
		wantSidebar  int
		wantRightMin int  // Minimum right pane width
		allowOverflow bool // Allow total to exceed terminal width
	}{
		{
			name:         "standard terminal",
			width:        120,
			height:       40,
			wantSidebar:  28,
			wantRightMin: 30,
			allowOverflow: false,
		},
		{
			name:         "wide terminal",
			width:        200,
			height:       60,
			wantSidebar:  28,
			wantRightMin: 30,
			allowOverflow: false,
		},
		{
			name:         "narrow terminal",
			width:        80,
			height:       24,
			wantSidebar:  28,
			wantRightMin: 30,
			allowOverflow: false,
		},
		{
			name:         "minimum usable terminal",
			width:        66, // 28 + 30 + 4 + buffer
			height:       20,
			wantSidebar:  28,
			wantRightMin: 30,
			allowOverflow: false,
		},
		{
			name:         "below minimum terminal",
			width:        50,
			height:       20,
			wantSidebar:  28,
			wantRightMin: 30,
			allowOverflow: true, // App maintains minimum sizes even if terminal is too small
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New().(model)
			updated, _ := m.Update(tea.WindowSizeMsg{Width: tt.width, Height: tt.height})
			m = updated.(model)

			// Check sidebar width is constant
			sidebarW := m.sidebarWidth()
			if sidebarW != tt.wantSidebar {
				t.Errorf("sidebarWidth() = %d, want %d", sidebarW, tt.wantSidebar)
			}

			// Check right pane has minimum width
			rightW := m.rightPaneWidth()
			if rightW < tt.wantRightMin {
				t.Errorf("rightPaneWidth() = %d, want at least %d", rightW, tt.wantRightMin)
			}

			// Check that sidebar + right pane + borders = total width (approximately)
			// Account for 4 chars: 2 for sidebar borders + 2 for right pane borders
			totalUsed := sidebarW + rightW + 4
			expectedTotal := tt.width

			// Check for overflow unless explicitly allowed
			if !tt.allowOverflow {
				// Allow some tolerance for border rendering
				if totalUsed > expectedTotal+2 {
					t.Errorf("total width used (%d) exceeds terminal width (%d)", totalUsed, expectedTotal)
				}
			} else {
				// For below-minimum terminals, verify we maintain minimum sizes
				// even though it will overflow
				if rightW < tt.wantRightMin {
					t.Errorf("rightPaneWidth() = %d, should maintain minimum of %d even when overflowing",
						rightW, tt.wantRightMin)
				}
			}
		})
	}
}

// TestLayoutProportions tests that editor and response panes split height correctly
func TestLayoutProportions(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{"standard", 120, 40},
		{"tall", 100, 80},
		{"short", 120, 20},
		{"square", 60, 60},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New().(model)
			updated, _ := m.Update(tea.WindowSizeMsg{Width: tt.width, Height: tt.height})
			m = updated.(model)

			contentHeight := max(m.height-1, 6) // -1 for status bar

			// Check that body and viewport have reasonable heights
			bodyHeight := m.body.Height()
			viewportHeight := m.view.Height

			// Both should be positive
			if bodyHeight <= 0 {
				t.Errorf("body height = %d, want > 0", bodyHeight)
			}
			if viewportHeight <= 0 {
				t.Errorf("viewport height = %d, want > 0", viewportHeight)
			}

			// Together they should not exceed content height
			// Add some buffer for titles and borders
			maxAllowed := contentHeight
			if bodyHeight+viewportHeight > maxAllowed {
				t.Errorf("body (%d) + viewport (%d) = %d exceeds content height %d",
					bodyHeight, viewportHeight, bodyHeight+viewportHeight, maxAllowed)
			}
		})
	}
}

// TestLayoutMinimumDimensions tests that components maintain minimum dimensions
func TestLayoutMinimumDimensions(t *testing.T) {
	// Test with very small window
	m := New().(model)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 40, Height: 10})
	m = updated.(model)

	// Sidebar should still have its width
	if m.sidebar.Width() < 0 {
		t.Errorf("sidebar width = %d, should not be negative", m.sidebar.Width())
	}

	// Right pane components should have minimum dimensions
	if m.url.Width < 0 {
		t.Errorf("url width = %d, should not be negative", m.url.Width)
	}

	if m.body.Width() < 0 {
		t.Errorf("body width = %d, should not be negative", m.body.Width())
	}

	if m.body.Height() < 0 {
		t.Errorf("body height = %d, should not be negative", m.body.Height())
	}

	if m.view.Height < 0 {
		t.Errorf("viewport height = %d, should not be negative", m.view.Height)
	}
}

// TestLayoutComponentWidths tests that all components fit within their panes
func TestLayoutComponentWidths(t *testing.T) {
	m := New().(model)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = updated.(model)

	rightWidth := m.rightPaneWidth()

	// URL should fit within right pane (with margin for label)
	// URL is now on its own line: "  URL:    " prefix (~14 chars) + url.Width
	urlTotal := 14 + m.url.Width

	if urlTotal > rightWidth+10 { // +10 buffer for padding
		t.Errorf("URL line width (%d) exceeds right pane width (%d)",
			urlTotal, rightWidth)
	}

	// Body should fit within right pane (accounting for padding)
	if m.body.Width() > rightWidth {
		t.Errorf("body width (%d) exceeds right pane width (%d)",
			m.body.Width(), rightWidth)
	}

	// Viewport should fit within right pane
	if m.view.Width > rightWidth {
		t.Errorf("viewport width (%d) exceeds right pane width (%d)",
			m.view.Width, rightWidth)
	}
}

// TestLayoutStability tests that layout remains consistent across multiple updates
func TestLayoutStability(t *testing.T) {
	m := New().(model)

	// Apply same window size twice
	msg := tea.WindowSizeMsg{Width: 120, Height: 40}

	updated1, _ := m.Update(msg)
	m1 := updated1.(model)

	updated2, _ := m1.Update(msg)
	m2 := updated2.(model)

	// Dimensions should be identical
	if m1.sidebar.Width() != m2.sidebar.Width() {
		t.Error("sidebar width changed on second update")
	}
	if m1.sidebar.Height() != m2.sidebar.Height() {
		t.Error("sidebar height changed on second update")
	}
	if m1.body.Width() != m2.body.Width() {
		t.Error("body width changed on second update")
	}
	if m1.body.Height() != m2.body.Height() {
		t.Error("body height changed on second update")
	}
	if m1.view.Width != m2.view.Width {
		t.Error("viewport width changed on second update")
	}
	if m1.view.Height != m2.view.Height {
		t.Error("viewport height changed on second update")
	}
}

// TestLayoutSidebarConstantWidth tests that sidebar width never changes
func TestLayoutSidebarConstantWidth(t *testing.T) {
	widths := []int{50, 80, 120, 150, 200}

	for _, width := range widths {
		m := New().(model)
		updated, _ := m.Update(tea.WindowSizeMsg{Width: width, Height: 40})
		m = updated.(model)

		sidebarW := m.sidebarWidth()
		if sidebarW != 28 {
			t.Errorf("at width %d: sidebarWidth() = %d, want 28", width, sidebarW)
		}
	}
}

// TestMinimumTerminalSize documents the minimum recommended terminal size
func TestMinimumTerminalSize(t *testing.T) {
	// The app requires minimum terminal width to display properly
	// Sidebar (28) + Right pane minimum (30) + Borders (4) = 62 characters minimum
	const minWidth = 62
	const minHeight = 20

	m := New().(model)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: minWidth, Height: minHeight})
	m = updated.(model)

	// At minimum width, components should fit without overflow
	sidebarW := m.sidebarWidth()
	rightW := m.rightPaneWidth()
	totalUsed := sidebarW + rightW + 4

	if totalUsed > minWidth+2 { // +2 tolerance
		t.Errorf("at minimum width %d, total used %d exceeds available space", minWidth, totalUsed)
	}

	// Document this in the test output
	t.Logf("Minimum recommended terminal size: %dx%d", minWidth, minHeight)
	t.Logf("Layout: sidebar=%d, right=%d, borders=4, total=%d", sidebarW, rightW, totalUsed)
}

// TestEditorContentFitsInPane tests that editor content doesn't overflow
func TestEditorContentFitsInPane(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{"standard", 120, 40},
		{"tall", 100, 80},
		{"short", 120, 20},
		{"minimum", 66, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New().(model)
			updated, _ := m.Update(tea.WindowSizeMsg{Width: tt.width, Height: tt.height})
			m = updated.(model)

			contentHeight := max(m.height-1, 6)
			editorHeight := contentHeight / 2
			if editorHeight < 5 {
				editorHeight = 5
			}

			// Calculate the lines used in the editor view (tabs are now in title bar)
			methodURLLines := 1   // method + URL row
			bodyTitleLines := 1   // "Body" title
			bodyContentLines := m.body.Height()

			totalEditorLines := methodURLLines + bodyTitleLines + bodyContentLines

			// Total should not exceed allocated height (with some margin for borders/padding)
			if totalEditorLines > editorHeight {
				t.Errorf("editor content (%d lines) exceeds allocated height (%d lines)",
					totalEditorLines, editorHeight)
				t.Logf("  method/URL: %d, body title: %d, body content: %d",
					methodURLLines, bodyTitleLines, bodyContentLines)
			}

			// Body height should never be negative or zero
			if m.body.Height() < 1 {
				t.Errorf("body height = %d, should be at least 1", m.body.Height())
			}
		})
	}
}

// TestPaneHeightAlignment tests that sidebar height equals editor + response heights
func TestPaneHeightAlignment(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{"standard", 120, 40},
		{"tall", 100, 80},
		{"short", 120, 20},
		{"minimum", 66, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New().(model)
			updated, _ := m.Update(tea.WindowSizeMsg{Width: tt.width, Height: tt.height})
			m = updated.(model)

			sidebarHeight := m.contentHeight()
			editorHeight := m.editorHeight()
			responseHeight := m.responseHeight()

			// Sidebar should equal editor + response
			if editorHeight+responseHeight != sidebarHeight {
				t.Errorf("height mismatch: editor (%d) + response (%d) = %d, but sidebar = %d",
					editorHeight, responseHeight, editorHeight+responseHeight, sidebarHeight)
			}
		})
	}
}

// TestLayoutRightPaneScales tests that right pane scales with window width
func TestLayoutRightPaneScales(t *testing.T) {
	// Test that right pane gets wider as window gets wider
	m1 := New().(model)
	updated1, _ := m1.Update(tea.WindowSizeMsg{Width: 80, Height: 40})
	m1 = updated1.(model)

	m2 := New().(model)
	updated2, _ := m2.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m2 = updated2.(model)

	m3 := New().(model)
	updated3, _ := m3.Update(tea.WindowSizeMsg{Width: 160, Height: 40})
	m3 = updated3.(model)

	right1 := m1.rightPaneWidth()
	right2 := m2.rightPaneWidth()
	right3 := m3.rightPaneWidth()

	if right2 <= right1 {
		t.Errorf("right pane should grow: %d (120w) should be > %d (80w)", right2, right1)
	}

	if right3 <= right2 {
		t.Errorf("right pane should grow: %d (160w) should be > %d (120w)", right3, right2)
	}
}
