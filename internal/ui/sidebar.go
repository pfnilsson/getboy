package ui

// viewSidebar renders the sidebar pane containing the list of saved requests.
func (m model) viewSidebar() string {
	content := m.sidebar.View()

	sbBox := titledPane(
		content,
		m.sidebarWidth(),
		m.pane == paneSidebar,
		paneBadge(1),
		"Requests",
	)
	return sbBox
}
