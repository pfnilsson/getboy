package ui

// viewResponse renders the response pane containing the HTTP response.
func (m model) viewResponse() string {
	content := m.view.View()

	respBox := titledPane(
		content,
		m.rightPaneWidth(),
		m.pane == paneResponse,
		paneBadge(3),
		"Response",
	)
	return respBox
}
