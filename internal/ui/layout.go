package ui

func (m *model) recomputeLayout() {
	contentHeight := max(m.height-2, 6)

	sidebarWidth := m.sidebarWidth()
	rightWidth := max(m.width-sidebarWidth-1, 30)

	editorHeight := contentHeight / 2
	respHeight := contentHeight - editorHeight
	if editorHeight < 5 {
		editorHeight = 5
	}
	if respHeight < 3 {
		respHeight = 3
	}

	m.sidebar.SetSize(sidebarWidth-2, contentHeight-2)
	m.method.Width = 8
	m.url.Width = rightWidth - 18
	m.body.SetWidth(rightWidth - 4)
	m.body.SetHeight(editorHeight - 4)
	m.view.Width = rightWidth - 4
	m.view.Height = respHeight - 3
}

func (m model) rightPaneWidth() int {
	sidebarW := m.sidebarWidth()
	width := max(m.width-sidebarW-1, 30)
	return width
}

func (m model) sidebarWidth() int {
	return 28
}
