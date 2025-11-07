package ui

// viewSidebar renders the sidebar pane containing the list of saved requests.
func (m model) viewSidebar() string {
	sb := m.sidebar.View()
	sbBox := titledPane(
		sb,
		m.sidebarWidth(),
		m.pane == paneSidebar,
		paneBadge(1),
		"Requests",
	)
	return sbBox
}
