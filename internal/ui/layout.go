package ui

func (m *model) recomputeLayout() {
	contentHeight := m.contentHeight()

	sidebarWidth := m.sidebarWidth()
	// Account for 4 chars: 2 for sidebar borders + 2 for right pane borders
	rightWidth := max(m.width-sidebarWidth-4, 30)

	editorHeight := m.editorHeight()
	respHeight := m.responseHeight()

	m.sidebar.SetSize(sidebarWidth-2, contentHeight-2)
	m.url.Width = rightWidth - 14 // Account for "  URL:    " prefix
	m.body.SetWidth(rightWidth - 4)
	m.body.SetHeight(editorHeight - 4)
	m.view.Width = rightWidth - 4
	m.view.Height = respHeight - 3
}

// contentHeight returns the height available for panes (total minus status bar)
func (m model) contentHeight() int {
	return max(m.height-1, 6)
}

// editorHeight returns the height for the editor/request pane
func (m model) editorHeight() int {
	contentHeight := m.contentHeight()
	editorHeight := contentHeight / 2
	if editorHeight < 5 {
		editorHeight = 5
	}
	return editorHeight
}

// responseHeight returns the height for the response pane
func (m model) responseHeight() int {
	contentHeight := m.contentHeight()
	editorHeight := m.editorHeight()
	respHeight := contentHeight - editorHeight
	if respHeight < 3 {
		respHeight = 3
	}
	return respHeight
}

func (m model) rightPaneWidth() int {
	sidebarW := m.sidebarWidth()
	// Account for 4 chars: 2 for sidebar borders + 2 for right pane borders
	width := max(m.width-sidebarW-4, 30)
	return width
}

func (m model) sidebarWidth() int {
	return 28
}
