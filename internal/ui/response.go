package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// viewResponse renders the response pane containing the HTTP response.
func (m model) viewResponse() string {
	resp := lipgloss.JoinVertical(lipgloss.Left, m.view.View())
	respBox := titledPane(
		resp,
		m.rightPaneWidth(),
		m.pane == paneResponse,
		paneBadge(3),
		"Response",
	)
	return respBox
}
